const std = @import("std");

const zz = @import("zigzag");

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
    var program = try zz.Program(Model).init(allocator);

    defer program.deinit();
    try program.run();
}
