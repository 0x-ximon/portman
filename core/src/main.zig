const std = @import("std");
const net = std.Io.net;
const Io = std.Io;

const App = @import("app.zig");

pub fn main(init: std.process.Init) !void {
    const env = init.environ_map;
    const gpa = init.gpa;
    const io = init.io;

    const arena = init.arena;
    defer arena.deinit();

    const host = env.get("HOST") orelse "127.0.0.1";
    const port = try std.fmt.parseInt(u16, env.get("PORT") orelse "3002", 10);

    const app = try App.init(arena.allocator(), .{ .host = host, .port = port });
    defer app.deinit(io);

    std.debug.print("Portman Core listening on {s}:{d}\n", .{ host, port });
    try app.run(io, gpa, .{ .reuse_address = true });
}
