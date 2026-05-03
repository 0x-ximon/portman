const std = @import("std");
const mem = std.mem;
const math = std.math;

pub fn Compare(comptime T: type) type {
    return fn (a: T, b: T) math.Order;
}

/// An ordered map implementation using a B-tree => BTreeMap.
/// Based on Chapter 18 of The Introduction to Algorithms by CLRS.
/// K => Key type, V => Value type, D => Degree (min number of children per node)
/// compare => Comparison function for keys
pub fn Map(comptime K: type, comptime V: type, comptime D: u8, comptime compare: Compare(K)) type {
    return struct {
        const Self = @This();

        root: ?*Node,
        size: u64,

        pub const empty: Self = .{ .root = null, .size = 0 };
        pub const Entry = struct { key: K, value: V };

        pub fn put(self: *Self, allocator: mem.Allocator, key: K, value: V) !void {
            if (self.root) |root| {
                root.insert(allocator, key, value);
                return;
            }

            const node = try allocator.create(Node);
            node.* = .{ .count = 1, .is_leaf = true };
            node.entries[0] = .{ .key = key, .value = value };

            self.root = node;
            self.size = 1;
        }

        pub fn get(self: *Self, key: K) ?V {
            if (self.root == null) return null;
            const entry = self.root.?.search(key);

            if (entry == null) return null;
            return entry.?.value;
        }

        pub fn contains(self: Self, key: K) bool {
            if (self.root == null) return false;
            return self.root.?.search(key) != null;
        }

        pub fn iterator(_: *Self) Iterator {
            return .{};
        }

        const Iterator = struct {
            pub fn next(_: *Iterator) ?Entry {
                return null;
            }
        };

        const Node = struct {
            count: u8 = 0,
            is_leaf: bool = false,
            entries: [(2 * D) - 1]Entry = undefined,
            children: [2 * D]?*Node = .{null} ** (2 * D),

            pub fn insert(self: *Node, allocator: mem.Allocator, key: K, value: V) void {
                _ = self;
                _ = allocator;
                _ = key;
                _ = value;
            }

            pub fn search(self: *Node, key: K) ?Entry {
                var i: usize = 0;
                while (i < self.count) : (i += 1) {
                    switch (compare(key, self.entries[i].key)) {
                        .eq => return self.entries[i],
                        .gt => continue,
                        .lt => break,
                    }
                }

                if (self.is_leaf or self.children[i] == null) return null;
                return self.children[i].?.search(key);
            }

            pub fn traverse(_: *Node) void {}
        };
    };
}
