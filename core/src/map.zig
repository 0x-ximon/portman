const std = @import("std");
const mem = std.mem;

// An ordered map implementation using a B-tree => BTreeMap.
pub fn Map(comptime K: type, comptime V: type, comptime _: u8) type {
    return struct {
        const Self = @This();

        root: ?*Node = null,

        pub const empty: Self = .{ .root = null };

        pub fn get(_: *Self, _: K) ?V {
            return null;
        }

        pub fn put(_: *Self, _: mem.Allocator, _: K, _: V) !void {
            return;
        }

        pub fn iterator(_: *Self) Iterator {
            return .{};
        }

        pub const Iterator = struct {
            pub const Entry = struct { key: K, value: V };

            pub fn next(_: *Iterator) ?Entry {
                return null;
            }
        };

        pub const Node = struct {};
    };
}
