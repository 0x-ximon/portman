const std = @import("std");
const mem = std.mem;
const Io = std.Io;

const nats = @import("nats");

const Book = @import("book.zig");
const Order = Book.Order;
const Ticker = Book.Ticker;
const Packet = @import("packet.zig");
const Header = Packet.Header;

const App = @This();

allocator: mem.Allocator,
queue: *nats.Client,
market: Market,
nonce: u32,

const Market = std.AutoHashMap(Ticker, *Book);

const AppError = error{
    ServerNotInitialized,
};

pub fn init(allocator: mem.Allocator, io: Io, queue_url: []const u8) !*App {
    const self = try allocator.create(App);
    const queue = try nats.Client.connect(allocator, io, queue_url, .{});

    self.* = .{
        .market = .init(allocator),
        .allocator = allocator,
        .queue = queue,
        .nonce = 0,
    };

    return self;
}

pub fn deinit(self: *App) void {
    var iter = self.market.valueIterator();
    while (iter.next()) |book| book.*.deinit();

    self.market.deinit();
    self.queue.deinit();
}

pub fn run(self: *App) !void {
    const subscription = try self.queue.subscribe("orders.match", nats.MsgHandler.init(App, self));
    defer subscription.deinit();
    while (true) {}
}

pub fn onMessage(self: *App, msg: *const nats.Message) void {
    self.nonce += 1;
    const header, const orders = Packet.recv(msg.data);
    const payload = self.handle(header, orders) catch |err| {
        std.log.err("Error: {s}", .{@errorName(err)});
        return;
    };

    if (payload.len > 0) {
        self.queue.publish("orders.processed", std.mem.sliceAsBytes(payload)) catch |err| {
            std.log.err("Error: {s}", .{@errorName(err)});
        };
    }
}

pub fn handle(self: *App, header: *const Header, orders: []Order) ![]Order {
    // PERF: Use a more performant data structure
    var payload = std.ArrayList(Order).empty;
    errdefer payload.deinit(self.allocator);

    for (orders) |*order| {
        const book = self.market.get(order.ticker) orelse blk: {
            const b = try Book.init(self.allocator);
            try self.market.put(order.ticker, b);
            break :blk b;
        };

        switch (header.instruction) {
            .CancelOrder => try book.cancelOrder(order),
            .UpdateOrder => try book.updateOrder(order),
            .ProcessOrder => {
                const processed = try book.processOrder(order);
                for (processed) |p| try payload.append(self.allocator, p);
            },
        }
    }

    return payload.toOwnedSlice(self.allocator);
}
