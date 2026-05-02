const std = @import("std");
const mem = std.mem;

const MapError = error{};

// An ordered map implementation using a B-tree => BTreeMap.
pub fn Map(comptime K: type, comptime V: type, comptime B: u8) type {
    return struct {
        const Self = @This();
        const BranchingFactor = B;

        pub fn init(allocator: mem.Allocator) !*Self {
            const self = try allocator.create(Self);
            return self;
        }

        pub fn deinit(self: *Self, allocator: mem.Allocator) void {
            allocator.destroy(self);
        }

        pub fn iterator(_: *Self) Iterator {
            return .{};
        }

        pub fn get(_: *Self, _: K) ?V {
            return null;
        }

        pub fn put(_: *Self, _: mem.Allocator, _: K, _: V) MapError!void {
            return;
        }

        pub const Iterator = struct {
            pub const Entry = struct { key: K, value: V };

            pub fn next(_: *Iterator) ?*Entry {
                return null;
            }
        };
    };
}
