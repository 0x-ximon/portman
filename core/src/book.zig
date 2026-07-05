const std = @import("std");
const Io = std.Io;
const mem = std.mem;
const math = std.math;

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
    open = 0,
    filled = 1,
    cancelled = 2,
};

const Side = enum(u8) {
    buy = 0,
    sell = 1,
};

const Mode = enum(u8) {
    gtc = 0,
    fok = 1,
    ioc = 2,
    aon = 3,
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
    var bids_iter = self.bids.iter(.ascending);
    while (bids_iter.next()) |*entry| entry.value.deinit(self.allocator);

    var asks_iter = self.asks.iter(.descending);
    while (asks_iter.next()) |*entry| entry.value.deinit(self.allocator);

    self.bids.deinit(self.allocator);
    self.asks.deinit(self.allocator);

    self.allocator.destroy(self);
}

pub fn newOrder(self: *Self, order: *Order) !void {
    switch (order.mode) {
        .gtc => try gtc(self, order),
        .fok => {
            const orders = try fok(self, order);
            _ = orders;
        },
        .ioc => try ioc(self, order),
        .aon => try aon(self, order),
    }
}

pub fn updateOrder(_: *Self, _: *Order) !void {}

pub fn cancelOrder(_: *Self, _: *Order) !void {}

// Good Till Cancelled
fn gtc(self: *Self, order: *Order) !void {
    const map = switch (order.side) {
        .buy => &self.bids,
        .sell => &self.asks,
    };

    var level = map.get(order.price) orelse blk: {
        const l = try Level.init(self.allocator);
        try map.put(self.allocator, order.price, l);
        break :blk l;
    };

    try level.orders.append(self.allocator, order.*);
    level.liquidity += order.quantity;
}

// Fill or Kill
fn fok(self: *Self, order: *Order) ![]Order {
    switch (order.side) {
        .buy => {
            const map = &self.asks;
            var available: u64 = 0;

            // Pass 1: Verify liquidity
            var iter = map.iter(.ascending);
            while (iter.next()) |*entry| {
                if (available >= order.quantity) break;
                if (entry.key > order.price) {
                    order.status = .cancelled;
                    return Errors.InsufficientLiquidity;
                }

                available += entry.value.liquidity;
            }

            if (available < order.quantity) {
                order.status = .cancelled;
                return Errors.InsufficientLiquidity;
            }

            // Pass 2: Fill the order
            iter = map.iter(.ascending);
            const settled = std.ArrayList(Order).empty;

            while (iter.next()) |*entry| {
                var indexes = std.ArrayList(usize).empty;
                for (0.., entry.value.orders.items) |index, *item| {
                    if (math.order(item.quantity, order.quantity) == .lt) {
                        order.quantity -= item.quantity;
                        item.quantity = 0;
                        indexes.append(self.allocator, index) catch {};
                    } else {
                        item.quantity -= order.quantity;
                        order.status = .filled;
                        order.quantity = 0;
                    }
                }

                entry.value.orders.orderedRemoveMany(indexes.items);
                if (order.status == .filled) break;
            }

            return settled.items;
        },

        .sell => {
            const map = &self.bids;
            var available: u64 = 0;

            // Pass 1: Verify liquidity
            var iter = map.iter(.descending);
            while (iter.prev()) |*entry| {
                if (available >= order.quantity) break;
                if (entry.key < order.price) {
                    order.status = .cancelled;
                    return Errors.InsufficientLiquidity;
                }

                available += entry.value.liquidity;
            }

            if (available < order.quantity) {
                order.status = .cancelled;
                return Errors.InsufficientLiquidity;
            }

            // Pass 2: Fill the order
            iter = map.iter(.descending);
            const settled = std.ArrayList(Order).empty;

            while (iter.next()) |*entry| {
                var indexes = std.ArrayList(usize).empty;
                for (0.., entry.value.orders.items) |index, *item| {
                    if (math.order(item.quantity, order.quantity) == .lt) {
                        order.quantity -= item.quantity;
                        item.quantity = 0;
                        indexes.append(self.allocator, index) catch {};
                    } else {
                        item.quantity -= order.quantity;
                        order.status = .filled;
                        order.quantity = 0;
                    }
                }

                entry.value.orders.orderedRemoveMany(indexes.items);
                if (order.status == .filled) break;
            }

            return settled.items;
        },
    }
}

// Immediate or Cancel
fn ioc(_: *Self, order: *const Order) !void {
    switch (order.side) {
        .buy => {},
        .sell => {},
    }
}

// All or None
fn aon(_: *Self, order: *const Order) !void {
    switch (order.side) {
        .buy => {},
        .sell => {},
    }
}
