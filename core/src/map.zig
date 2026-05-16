const std = @import("std");
const mem = std.mem;
const math = std.math;

pub fn Compare(comptime T: type) type {
    return fn (a: T, b: T) math.Order;
}

/// An Ordered Map backed by a B-tree => BTreeMap.
/// K => Key type, V => Value type, D => Degree (min number of children per node)
/// compare => Comparison function for keys
pub fn Map(comptime K: type, comptime V: type, comptime D: usize, comptime compare: Compare(K)) type {
    return struct {
        const Self = @This();

        root: ?*Node,
        size: u64,

        pub const empty: Self = .{ .root = null, .size = 0 };
        pub const Entry = struct { key: K, value: V };

        pub fn put(self: *Self, allocator: mem.Allocator, key: K, value: V) !void {
            if (self.root) |root| {
                const node = try root.insert(allocator, key, value);
                if (node) |n| self.root = n;
            } else {
                const node = try allocator.create(Node);
                node.* = .{ .count = 1, .is_leaf = true };
                node.entries[0] = .{ .key = key, .value = value };
                self.root = node;
            }

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
            count: usize = 0,
            is_leaf: bool = false,
            entries: [2 * D]Entry = undefined,
            children: [(2 * D) + 1]?*Node = .{null} ** ((2 * D) + 1),

            fn insert(self: *Node, allocator: mem.Allocator, key: K, value: V) !?*Node {
                switch (self.is_leaf) {
                    true => {
                        switch (self.isFull()) {
                            // Invariant: A full leaf node during insertion should not exist
                            true => unreachable,
                            false => {
                                const entry: Entry = .{ .key = key, .value = value };
                                self.entries[self.count] = entry;
                                self.count += 1;

                                // Sort the entries using insertion sort
                                for (self.count - 1..1) |i|
                                    if (compare(self.entries[i].key, self.entries[i - 1].key) == .lt)
                                        std.mem.swap(Entry, &self.entries[i], &self.entries[i - 1]);
                            },
                        }
                    },

                    false => {
                        var i: usize = 0;
                        while (i < self.count) : (i += 1) {
                            switch (compare(key, self.entries[i].key)) {
                                .lt => break,
                                .gt => continue,
                                .eq => {
                                    self.entries[i].value = value;
                                    return null;
                                },
                            }
                        }

                        var child = self.children[i] orelse unreachable;
                        // TODO: Handle the splits
                        // if (child.isFull()) child = try self.split(allocator, child, i);
                        return try child.insert(allocator, key, value);
                    },
                }

                return null;
            }

            fn delete(_: *Node, _: mem.Allocator, _: K, _: ?*Node, _: ?usize) void {}

            fn search(self: *Node, key: K) ?Entry {
                var i: usize = 0;
                while (i < self.count) : (i += 1) {
                    switch (compare(key, self.entries[i].key)) {
                        .eq => return self.entries[i],
                        .gt => continue,
                        .lt => break,
                    }
                }

                if (self.is_leaf) return null;
                const child = self.children[i] orelse unreachable;
                return child.search(key);
            }

            fn split(_: *Node, allocator: mem.Allocator, _: *Node, _: usize) !void {
                _ = try allocator.create(Node);
            }

            fn traverse(_: *Node) void {}

            fn isFull(self: *Node) bool {
                return self.count == (2 * D);
            }
        };
    };
}
