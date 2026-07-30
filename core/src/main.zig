const std = @import("std");
const net = std.Io.net;
const Io = std.Io;

const App = @import("app.zig");

pub fn main(init: std.process.Init) !void {
    const env = init.environ_map;
    const gpa = init.gpa;
    const io = init.io;

    const queue_url = env.get("QUEUE_URL") orelse "nats://localhost:4222";
    const app = try App.init(gpa, io, queue_url);
    defer app.deinit();

    std.debug.print("Portman Core Starting\n", .{});
    try app.run();
}
