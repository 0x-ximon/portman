const std = @import("std");
const Io = std.Io;
const mem = std.mem;
const math = std.math;
const testing = std.testing;

const Self = @This();
const Map = @import("map.zig")
    .Map(u64, *Level, 32, compare);

fn compare(a: u64, b: u64) math.Order {
    return math.order(a, b);
}

allocator: mem.Allocator,
asks: Map,
bids: Map,

const Status = enum(u8) {
    pending = 0,
    filled = 1,
    partial = 2,
    cancelled = 3,
};

const Side = enum(u8) {
    buy = 0,
    sell = 1,
};

const Mode = enum(u8) {
    gtc = 0,
    fok = 1,
    ioc = 2,
};

const Flags = packed struct(u8) {
    post_only: bool = false,
    hidden: bool = false,
    _: u6 = 0,
};

pub const Order = extern struct {
    id: u64,
    price: u64,
    quantity: u64,
    asset: u16,
    price_precision: u8,
    qty_precision: u8,
    status: Status,
    side: Side,
    mode: Mode,
    flags: Flags,
};

const Errors = error{
    InsufficientLiquidity,
};

const Level = struct {
    liquidity: u64,
    orders: std.ArrayList(Order),

    pub fn init(allocator: mem.Allocator) !*Level {
        const self = try allocator.create(Level);
        self.* = .{
            .orders = .empty,
            .liquidity = 0,
        };

        return self;
    }

    pub fn deinit(self: *Level, allocator: mem.Allocator) void {
        self.orders.deinit(allocator);
        allocator.destroy(self);
    }
};

pub fn init(allocator: std.mem.Allocator) !*Self {
    const self = try allocator.create(Self);
    self.* = .{
        .asks = .init(),
        .bids = .init(),
        .allocator = allocator,
    };

    return self;
}

pub fn deinit(self: *Self) void {
    var asks_iter = self.asks.iter(.ascending);
    while (asks_iter.next()) |*entry| entry.value.deinit(self.allocator);

    var bids_iter = self.bids.iter(.descending);
    while (bids_iter.prev()) |*entry| entry.value.deinit(self.allocator);

    self.bids.deinit(self.allocator);
    self.asks.deinit(self.allocator);

    self.allocator.destroy(self);
}

pub fn processOrder(self: *Self, order: *Order) ![]Order {
    return switch (order.mode) {
        .gtc => try gtc(self, order),
        .fok => try fok(self, order),
        .ioc => try ioc(self, order),
    };
}

pub fn updateOrder(_: *Self, _: *Order) !void {}

pub fn cancelOrder(_: *Self, _: *Order) !void {}

// Good Till Cancelled
fn gtc(self: *Self, order: *Order) ![]Order {
    var settled = std.ArrayList(Order).empty;
    errdefer settled.deinit(self.allocator);

    const Entry = struct { key: u64, value: *Level };
    var cleared = std.ArrayList(Entry).empty;
    defer cleared.deinit(self.allocator);

    switch (order.side) {
        .buy => {
            var map = &self.asks;

            // Fill with whatever is available
            var iter = map.iter(.ascending);
            while (iter.next()) |*entry| {
                const price = entry.key;
                const level = entry.value;
                if (price > order.price) break;

                var indexes = std.ArrayList(usize).empty;
                defer indexes.deinit(self.allocator);

                try self.match(order, level, &indexes, &settled);
                level.orders.orderedRemoveMany(indexes.items);

                if (level.liquidity == 0) try cleared.append(self.allocator, .{ .key = price, .value = level });
                if (order.quantity == 0) break;
            }

            for (cleared.items) |entry| {
                const deleted = try map.rid(self.allocator, entry.key);
                if (!deleted) unreachable;
                entry.value.deinit(self.allocator);
            }

            // Order not fully filled, store in bids
            if (order.status != .filled) {
                map = &self.bids;

                var level = map.get(order.price) orelse blk: {
                    const l = try Level.init(self.allocator);
                    errdefer l.deinit(self.allocator);
                    try map.put(self.allocator, order.price, l);
                    break :blk l;
                };

                try level.orders.append(self.allocator, order.*);
                level.liquidity += order.quantity;
            }
        },
        .sell => {
            var map = &self.bids;

            var iter = map.iter(.descending);
            while (iter.prev()) |*entry| {
                const price = entry.key;
                const level = entry.value;
                if (price < order.price) break;

                var indexes = std.ArrayList(usize).empty;
                defer indexes.deinit(self.allocator);

                try self.match(order, level, &indexes, &settled);
                level.orders.orderedRemoveMany(indexes.items);

                if (level.liquidity == 0) try cleared.append(self.allocator, .{ .key = price, .value = level });
                if (order.quantity == 0) break;
            }

            for (cleared.items) |entry| {
                const deleted = try map.rid(self.allocator, entry.key);
                if (!deleted) unreachable;
                entry.value.deinit(self.allocator);
            }

            // Order not fully filled, store in asks
            if (order.status != .filled) {
                map = &self.asks;

                var level = map.get(order.price) orelse blk: {
                    const l = try Level.init(self.allocator);
                    errdefer l.deinit(self.allocator);
                    try map.put(self.allocator, order.price, l);
                    break :blk l;
                };

                try level.orders.append(self.allocator, order.*);
                level.liquidity += order.quantity;
            }
        },
    }

    return try settled.toOwnedSlice(self.allocator);
}

// Fill or Kill
fn fok(self: *Self, order: *Order) ![]Order {
    var settled = std.ArrayList(Order).empty;
    errdefer settled.deinit(self.allocator);

    const Entry = struct { key: u64, value: *Level };
    var cleared = std.ArrayList(Entry).empty;
    defer cleared.deinit(self.allocator);

    switch (order.side) {
        .buy => {
            const map = &self.asks;

            // Pass 1: Verify liquidity
            var available: u64 = 0;
            var iter = map.iter(.ascending);
            while (iter.next()) |*entry| {
                const price = entry.key;
                const level = entry.value;

                if (price > order.price) {
                    order.status = .cancelled;
                    return Errors.InsufficientLiquidity;
                }

                available += level.liquidity;
                if (available >= order.quantity) break;
            }

            if (available < order.quantity) {
                order.status = .cancelled;
                return Errors.InsufficientLiquidity;
            }

            // Pass 2: Fill the order
            iter = map.iter(.ascending);
            while (iter.next()) |*entry| {
                const price = entry.key;
                const level = entry.value;

                var indexes = std.ArrayList(usize).empty;
                defer indexes.deinit(self.allocator);

                try self.match(order, level, &indexes, &settled);
                level.orders.orderedRemoveMany(indexes.items);

                if (level.liquidity == 0) try cleared.append(self.allocator, .{ .key = price, .value = level });
                if (order.quantity == 0) break;
            }

            for (cleared.items) |entry| {
                const deleted = try map.rid(self.allocator, entry.key);
                if (!deleted) unreachable;
                entry.value.deinit(self.allocator);
            }
        },

        .sell => {
            const map = &self.bids;

            // Pass 1: Verify liquidity
            var available: u64 = 0;
            var iter = map.iter(.descending);
            while (iter.prev()) |*entry| {
                const price = entry.key;
                const level = entry.value;

                if (price < order.price) {
                    order.status = .cancelled;
                    return Errors.InsufficientLiquidity;
                }

                available += level.liquidity;
                if (available >= order.quantity) break;
            }

            if (available < order.quantity) {
                order.status = .cancelled;
                return Errors.InsufficientLiquidity;
            }

            // Pass 2: Fill the order
            iter = map.iter(.descending);
            while (iter.prev()) |*entry| {
                const price = entry.key;
                const level = entry.value;

                var indexes = std.ArrayList(usize).empty;
                defer indexes.deinit(self.allocator);

                try self.match(order, level, &indexes, &settled);
                level.orders.orderedRemoveMany(indexes.items);

                if (level.liquidity == 0) try cleared.append(self.allocator, .{ .key = price, .value = level });
                if (order.quantity == 0) break;
            }

            for (cleared.items) |entry| {
                const deleted = try map.rid(self.allocator, entry.key);
                if (!deleted) unreachable;
                entry.value.deinit(self.allocator);
            }
        },
    }

    return try settled.toOwnedSlice(self.allocator);
}

// Immediate or Cancel
fn ioc(self: *Self, order: *Order) ![]Order {
    var settled = std.ArrayList(Order).empty;
    errdefer settled.deinit(self.allocator);

    const Entry = struct { key: u64, value: *Level };
    var cleared = std.ArrayList(Entry).empty;
    defer cleared.deinit(self.allocator);

    switch (order.side) {
        .buy => {
            const map = &self.asks;

            var iter = map.iter(.ascending);
            while (iter.next()) |*entry| {
                const price = entry.key;
                const level = entry.value;

                if (price > order.price) break;

                var indexes = std.ArrayList(usize).empty;
                defer indexes.deinit(self.allocator);

                try self.match(order, level, &indexes, &settled);
                level.orders.orderedRemoveMany(indexes.items);

                if (level.liquidity == 0) try cleared.append(self.allocator, .{ .key = price, .value = level });
                if (order.quantity == 0) break;
            }

            for (cleared.items) |entry| {
                const deleted = try map.rid(self.allocator, entry.key);
                if (!deleted) unreachable;
                entry.value.deinit(self.allocator);
            }
        },

        .sell => {
            const map = &self.bids;
            var iter = map.iter(.descending);

            while (iter.prev()) |*entry| {
                const price = entry.key;
                const level = entry.value;

                if (price < order.price) break;

                var indexes = std.ArrayList(usize).empty;
                defer indexes.deinit(self.allocator);

                try self.match(order, level, &indexes, &settled);
                level.orders.orderedRemoveMany(indexes.items);

                if (level.liquidity == 0) try cleared.append(self.allocator, .{ .key = price, .value = level });
                if (order.quantity == 0) break;
            }

            for (cleared.items) |entry| {
                const deleted = try map.rid(self.allocator, entry.key);
                if (!deleted) unreachable;
                entry.value.deinit(self.allocator);
            }
        },
    }

    return settled.toOwnedSlice(self.allocator);
}

fn match(
    self: *Self,
    order: *Order,
    level: *Level,
    indexes: *std.ArrayList(usize),
    settled: *std.ArrayList(Order),
) !void {
    for (0.., level.orders.items) |index, *item| {
        switch (math.order(item.quantity, order.quantity)) {
            .lt => {
                level.liquidity -= item.quantity;
                order.quantity -= item.quantity;

                item.quantity = 0;
                item.status = .filled;
                order.status = .partial;

                try indexes.append(self.allocator, index);
                try settled.append(self.allocator, item.*);
            },
            .eq => {
                level.liquidity -= order.quantity;
                order.status = .filled;
                item.status = .filled;
                order.quantity = 0;
                item.quantity = 0;

                try indexes.append(self.allocator, index);
                try settled.append(self.allocator, item.*);
                try settled.append(self.allocator, order.*);

                break;
            },

            .gt => {
                level.liquidity -= order.quantity;
                item.quantity -= order.quantity;

                order.quantity = 0;
                item.status = .partial;
                order.status = .filled;

                try settled.append(self.allocator, item.*);
                try settled.append(self.allocator, order.*);
                break;
            },
        }
    }
}

const Context = struct {
    const Ctx = @This();

    const E = error{
        CrossedBook,
        LiquidityMismatch,
        DeadOrderResting,
        InvalidSettledState,
    };

    fn bookInvariantsHold(book: *Self) !void {
        try verifyLiquidity(&book.bids);
        try verifyLiquidity(&book.asks);

        var bid_iter = book.bids.iter(.descending);
        var ask_iter = book.asks.iter(.ascending);

        var best_bid = bid_iter.prev();
        var best_ask = ask_iter.next();

        while (best_bid != null and best_ask != null) {
            const bid = best_bid orelse unreachable;
            const ask = best_ask orelse unreachable;

            // Bids should never be greater than Ask
            if (bid.key >= ask.key) return E.CrossedBook;

            best_bid = bid_iter.prev();
            best_ask = ask_iter.next();
        }
    }

    fn verifyLiquidity(map: *Map) !void {
        var iter = map.iter(.ascending);
        while (iter.next()) |entry| {
            const level = entry.value;
            var expected_liquidity: u64 = 0;

            for (level.orders.items) |o| {
                // A resting order should never be fully filled or have 0 quantity
                if (o.quantity == 0 or o.status == .filled) return E.DeadOrderResting;
                expected_liquidity += o.quantity;
            }

            // The cached level liquidity must exactly equal the sum of resting quantities
            if (expected_liquidity != level.liquidity) return E.LiquidityMismatch;
        }
    }

    fn totalLiquidity(map: *Map) u64 {
        var sum: u64 = 0;
        var iter = map.iter(.ascending);
        while (iter.next()) |entry| sum += entry.value.liquidity;
        return sum;
    }
};

fn testBook(_: Context, smith: *std.testing.Smith) !void {
    const allocator = testing.allocator;
    const book = try init(allocator);
    defer book.deinit();

    while (!smith.eos()) {
        var order: Order = .{
            .status = .pending,
            .qty_precision = 3,
            .price_precision = 3,
            .id = smith.value(u64),
            .side = smith.value(Side),
            .mode = smith.value(Mode),
            .flags = smith.value(Flags),
            .asset = smith.valueRangeAtMost(u16, 1000, 1010),
            .price = smith.valueRangeAtMost(u64, 100_000, 150_000),
            .quantity = smith.valueRangeAtMost(u64, 1000, 1_000_000_000),
        };

        const bids_before = Context.totalLiquidity(&book.bids);
        const asks_before = Context.totalLiquidity(&book.asks);
        const q_start = order.quantity;

        if (book.processOrder(&order)) |settled| {
            defer allocator.free(settled);

            const bids_after = Context.totalLiquidity(&book.bids);
            const asks_after = Context.totalLiquidity(&book.asks);

            const q_rem = order.quantity;
            const q_traded = q_start - q_rem;

            switch (order.side) {
                .buy => {
                    try testing.expectEqual(asks_before - asks_after, q_traded);

                    if (order.mode == .gtc and order.status != .filled)
                        try testing.expectEqual(bids_after - bids_before, q_rem)
                    else
                        try testing.expectEqual(bids_before, bids_after);
                },
                .sell => {
                    try testing.expectEqual(bids_before - bids_after, q_traded);

                    if (order.mode == .gtc and order.status != .filled)
                        try testing.expectEqual(asks_after - asks_before, q_rem)
                    else
                        try testing.expectEqual(asks_before, asks_after);
                },
            }
        } else |err| {
            try testing.expectEqual(order.status, .cancelled);
            try testing.expect(err == Errors.InsufficientLiquidity);

            try testing.expectEqual(bids_before, Context.totalLiquidity(&book.bids));
            try testing.expectEqual(asks_before, Context.totalLiquidity(&book.asks));
        }

        try Context.bookInvariantsHold(book);
    }
}

test "Book: Fuzz" {
    try testing.fuzz(Context{}, testBook, .{});
}
