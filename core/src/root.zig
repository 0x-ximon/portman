const std = @import("std");
const testing = std.testing;

const Book = @import("book.zig");
const Packet = @import("packet.zig");

test "Integration" {
    const allocator = testing.allocator;
    const Context = struct {
        const Self = @This();
    };

    const ctx = Context{};
    _ = ctx;

    const book = try Book.init(allocator);
    defer book.deinit();
}
