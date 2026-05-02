const std = @import("std");
const Io = std.Io;
const mem = std.mem;

const Map = @import("map.zig").Map;

const Book = @This();

allocator: mem.Allocator,
asks: Map(u64, *Level, 4),
bids: Map(u64, *Level, 4),

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

pub fn init(allocator: std.mem.Allocator) !*Book {
    const self = try allocator.create(Book);
    self.* = .{
        .asks = .empty,
        .bids = .empty,
        .allocator = allocator,
    };

    return self;
}

pub fn deinit(self: *Book) void {
    var bids_iter = self.bids.iterator();
    while (bids_iter.next()) |*entry| entry.value.deinit(self.allocator);

    var asks_iter = self.asks.iterator();
    while (asks_iter.next()) |*entry| entry.value.deinit(self.allocator);

    self.allocator.destroy(self);
}

pub fn newOrder(self: *Book, order: *const Order) !void {
    switch (order.mode) {
        .gtc => try gtc(self, order),
        .fok => try fok(self, order),
        .ioc => try ioc(self, order),
    }
}

pub fn updateOrder(_: *Book, _: *const Order) !void {}

pub fn cancelOrder(_: *Book, _: *const Order) !void {}

// Good Till Cancelled
fn gtc(self: *Book, order: *const Order) !void {
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
fn fok(_: *Book, order: *const Order) !void {
    switch (order.side) {
        .buy => {},
        .sell => {},
    }
}

// Immediate or Cancel
fn ioc(_: *Book, order: *const Order) !void {
    switch (order.side) {
        .buy => {},
        .sell => {},
    }
}
