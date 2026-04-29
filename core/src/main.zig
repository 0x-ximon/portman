const std = @import("std");
const net = std.Io.net;
const Io = std.Io;

const Order = extern struct {
    order_id: u64,
    price: i64,
    quantity: u32,
    side: u8,
    _padding: [3]u8, // Matches the Go padding
};

pub fn main(init: std.process.Init) !void {
    const env = init.environ_map;
    const io = init.io;

    const host = env.get("HOST") orelse "127.0.0.1";
    const port = try std.fmt.parseInt(u16, env.get("PORT") orelse "3002", 10);

    const address = try net.IpAddress.parse(host, port);
    var server = try address.listen(io, .{});
    defer server.deinit(io);

    std.debug.print("Portman Core listening on {f}\n", .{server.socket.address});

    while (true) {
        const stream = try server.accept(io);
        defer stream.socket.close(io);

        var buffer: [1024 * 64]u8 = undefined;
        var buf_reader = stream.reader(io, &buffer);
        var reader = &buf_reader.interface;

        const msg = reader.takeDelimiterExclusive('\n') catch |err| blk: {
            if (err == error.EndOfStream) {
                break :blk "Connection Closed";
            } else return err;
        };
        std.log.info("{s}", .{msg});
    }
}
