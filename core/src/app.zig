const std = @import("std");
const mem = std.mem;
const Io = std.Io;
const net = Io.net;

const Book = @import("book.zig");
const Order = Book.Order;
const Ticker = Book.Ticker;
const Packet = @import("packet.zig");
const Header = Packet.Header;

const App = @This();

address: net.IpAddress,
server: ?net.Server,
market: Market,

const Market = std.AutoHashMap(Ticker, *Book);

const AppError = error{
    ServerNotInitialized,
};

pub const AppConfig = struct {
    host: []const u8,
    port: u16,
};

pub fn init(allocator: mem.Allocator, config: AppConfig) !*App {
    const self = try allocator.create(App);
    self.* = .{
        .address = try .parse(config.host, config.port),
        .market = .init(allocator),
        .server = null,
    };

    return self;
}

pub fn deinit(self: *App, io: Io) void {
    var iter = self.market.valueIterator();
    while (iter.next()) |book| book.*.deinit();

    if (self.server) |*s| s.deinit(io);
    self.market.deinit();
}

pub fn run(self: *App, io: Io, allocator: mem.Allocator, options: net.IpAddress.ListenOptions) !void {
    self.server = try self.address.listen(io, options);
    while (true) {
        self.handle(io, allocator) catch |err| {
            std.log.err("something went wrong: {any}", .{err});
        };
    }
}

fn handle(self: *App, io: Io, allocator: mem.Allocator) !void {
    if (self.server == null) {
        return AppError.ServerNotInitialized;
    }

    const stream = try self.server.?.accept(io);
    defer stream.socket.close(io);

    var buffer: [1024 * 64]u8 align(8) = undefined;
    var buf_reader = stream.reader(io, &buffer);
    const reader = &buf_reader.interface;

    const header = try Packet.info(reader);
    const orders = try Packet.recv(reader, header.length);

    for (orders) |*order| {
        const book = self.market.get(order.ticker) orelse blk: {
            const b = try Book.init(allocator);
            try self.market.put(order.ticker, b);
            break :blk b;
        };

        switch (header.instruction) {
            .CancelOrder => try book.cancelOrder(order),
            .UpdateOrder => try book.updateOrder(order),
            .ProcessOrder => _ = try book.processOrder(order),
        }

        std.log.info("order processed: {any}", .{order});
    }
}
