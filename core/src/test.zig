const std = @import("std");
const testing = std.testing;

pub const Book = @import("book.zig");
pub const Map = @import("map.zig");
pub const Packet = @import("packet.zig");

test "Semantic Analysis Discovery" {
    testing.refAllDecls(Book);
    testing.refAllDecls(Map);
    testing.refAllDecls(Packet);
}
