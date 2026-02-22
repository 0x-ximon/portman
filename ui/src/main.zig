const std = @import("std");

const vaxis = @import("vaxis");
const vxfw = vaxis.vxfw;

const Model = @import("model.zig");

pub fn main() !void {
    var debug_allocator = std.heap.DebugAllocator(.{}).init;
    defer {
        const check = debug_allocator.deinit();
        if (check != .ok) {
            std.debug.print("Memory leak detected: {}\n", .{check});
        }
    }

    const allocator = debug_allocator.allocator();
    var app = try vxfw.App.init(allocator);
    defer app.deinit();

    const model = try Model.init(allocator);
    defer model.deinit();

    try app.run(model.widget(), .{});
}
