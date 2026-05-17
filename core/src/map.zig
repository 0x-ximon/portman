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
    const max_children: comptime_int = 2 * D;
    const max_entries: comptime_int = (2 * D) - 1;

    return struct {
        const Self = @This();

        root: ?*Node,
        size: u64,

        pub const empty: Self = .{ .root = null, .size = 0 };
        pub const Entry = struct { key: K, value: V };

        pub fn put(self: *Self, allocator: mem.Allocator, key: K, value: V) !void {
            if (self.root) |root| {
                switch (root.count == max_entries) {
                    true => {
                        const node = try allocator.create(Node);
                        node.* = .{ .count = 1, .is_leaf = false };
                        node.children[0] = root;
                        self.root = node;

                        try node.split(allocator, root, 0);
                        const inserted = try node.insert(allocator, key, value);
                        if (inserted) self.size += 1;
                    },
                    false => {
                        const inserted = try root.insert(allocator, key, value);
                        if (inserted) self.size += 1;
                    },
                }
            } else {
                const node = try allocator.create(Node);
                node.* = .{ .count = 1, .is_leaf = true };
                node.entries[0] = .{ .key = key, .value = value };
                self.root = node;
                self.size = 1;
            }
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
            entries: [max_entries]Entry = undefined,
            children: [max_children]?*Node = .{null} ** max_children,

            fn insert(self: *Node, allocator: mem.Allocator, key: K, value: V) !bool {
                switch (self.is_leaf) {
                    true => {
                        switch (self.count == max_entries) {
                            // Invariant: A full leaf node during insertion should not exist
                            true => unreachable,
                            false => {
                                // Ensure the key does not already exist in the leaf node
                                for (0..self.count) |i| {
                                    switch (compare(key, self.entries[i].key)) {
                                        .lt => break,
                                        .gt => continue,
                                        .eq => {
                                            self.entries[i].value = value;
                                            return false;
                                        },
                                    }
                                }

                                const entry: Entry = .{ .key = key, .value = value };
                                self.entries[self.count] = entry;
                                self.count += 1;

                                // Sort the entries using insertion sort
                                var i = self.count - 1;
                                while (i > 0) : (i -= 1)
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
                                    return false;
                                },
                            }
                        }

                        const child = self.children[i] orelse unreachable;
                        if (child.count == max_entries) try self.split(allocator, child, i);
                        return try self.insert(allocator, key, value);
                    },
                }

                return true;
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

            fn split(self: *Node, allocator: mem.Allocator, child: *Node, index: usize) !void {
                const other = try allocator.create(Node);
                for (0..D - 1) |i| other.entries[i] = child.entries[i + D];

                if (!child.is_leaf) {
                    for (0..D) |i| other.children[i] = child.children[i + D];
                }

                child.* = .{ .count = D - 1, .is_leaf = child.is_leaf };
                other.* = .{ .count = D - 1, .is_leaf = child.is_leaf };

                // Shift entries and children of parent
                var j = self.count - 1;
                while (j > index) : (j -= 1) self.entries[j] = self.entries[j - 1];

                var k = self.count;
                while (k > index) : (k -= 1) self.children[k] = self.children[k - 1];

                // Promote the middle entry and the new node to the parent
                self.entries[index] = child.entries[D - 1];
                self.children[index + 1] = other;
                self.count += 1;
            }

            fn traverse(_: *Node) void {}
        };
    };
}
