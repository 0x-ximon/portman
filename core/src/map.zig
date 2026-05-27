const std = @import("std");
const mem = std.mem;
const math = std.math;

pub fn Compare(comptime T: type) type {
    return fn (a: T, b: T) math.Order;
}

/// A Sorted Map backed by a B+ Tree.
///
/// `K` & `V` are the key and value types, `M` is the minimum M-way factor.
/// `compare` is the comparison function for keys.
///
/// This map does not handle memory cleanup for keys or values but only nodes.
pub fn Map(comptime K: type, comptime V: type, comptime M: usize, comptime compare: Compare(K)) type {
    // Assuming M = 4
    const min_keys = M - 1; // 3
    const max_keys = (2 * M) - 1; // 7
    const min_children = M; // 4
    const max_children = (2 * M); // 8

    const Entry = struct { key: K, value: V };

    const Node = union((enum { inner, leaf })) {
        const Self = @This();

        inner: Inner,
        leaf: Leaf,

        const Inner = struct {
            count: usize = 0, // current number of keys
            keys: [max_keys]K = undefined,
            children: [max_children]?*Self = .{null} ** max_children,
        };

        const Leaf = struct {
            count: usize = 0, // current number of entries
            next: ?*Self = null,
            entries: [max_keys]Entry = undefined,
        };

        fn search(self: *Self, key: K) ?Entry {
            switch (self.*) {
                .inner => |*node| {
                    var i: usize = 0;
                    while (i < node.count) : (i += 1) {
                        switch (compare(key, node.keys[i])) {
                            .gt, .eq => continue,
                            .lt => break,
                        }
                    }

                    const child = node.children[i] orelse unreachable;
                    return child.search(key);
                },

                .leaf => |*node| {
                    for (0..node.count) |i| {
                        switch (compare(key, node.entries[i].key)) {
                            .gt => continue,
                            .lt => break,
                            .eq => return node.entries[i],
                        }
                    }

                    return null;
                },
            }
        }

        fn insert(self: *Self, allocator: mem.Allocator, key: K, value: V) !bool {
            switch (self.*) {
                .leaf => |*node| {
                    var i: usize = 0;
                    while (i < node.count) : (i += 1) {
                        switch (compare(key, node.entries[i].key)) {
                            .gt => continue,
                            .lt => break,
                            .eq => {
                                node.entries[i].value = value;
                                return false;
                            },
                        }
                    }

                    var j = node.count;
                    while (j > i) : (j -= 1) node.entries[j] = node.entries[j - 1];
                    node.entries[i] = .{ .key = key, .value = value };
                    node.count += 1;
                    return true;
                },

                .inner => |*node| {
                    var i: usize = 0;
                    while (i < node.count) : (i += 1) {
                        switch (compare(key, node.keys[i])) {
                            .gt, .eq => continue,
                            .lt => break,
                        }
                    }

                    var child = node.children[i] orelse unreachable;
                    if (child.full()) child = try self.split(allocator, child, i, key);
                    return try child.insert(allocator, key, value);
                },
            }
        }

        fn delete(self: *Self, allocator: mem.Allocator, key: K) !bool {
            switch (self.*) {
                .leaf => |*node| {
                    for (0..node.count) |i| {
                        switch (compare(key, node.entries[i].key)) {
                            .gt => continue,
                            .lt => break,
                            .eq => {
                                for (i..node.count - 1) |j| node.entries[j] = node.entries[j + 1];
                                node.count -= 1;
                                return true;
                            },
                        }
                    }

                    return false;
                },
                .inner => |*node| {
                    var i: usize = 0;
                    while (i < node.count) : (i += 1) {
                        switch (compare(key, node.keys[i])) {
                            .gt, .eq => continue,
                            .lt => break,
                        }
                    }

                    var child = node.children[i] orelse unreachable;
                    if (child.underfull()) {
                        const left = if (i == 0) null else node.children[i - 1];
                        const right = if (i == node.count) null else node.children[i + 1];
                        child = try self.move(allocator, child, left, right, i);
                    }

                    return try child.delete(allocator, key);
                },
            }
        }

        fn split(self: *Self, allocator: mem.Allocator, child: *Self, index: usize, key: K) !*Self {
            const other = try allocator.create(Self);

            switch (child.*) {
                .leaf => |*node| {
                    // Child gets 4, Other gets 3, First of Other gets promoted to self
                    other.* = .{ .leaf = .{ .count = min_keys, .next = node.next } };
                    @memcpy(other.leaf.entries[0..min_keys], node.entries[M..max_keys]);

                    node.count = M;
                    node.next = other;

                    var i = self.inner.count;
                    while (i > index) : (i -= 1) {
                        self.inner.keys[i] = self.inner.keys[i - 1];
                        self.inner.children[i + 1] = self.inner.children[i];
                    }

                    self.inner.keys[i] = other.leaf.entries[0].key;
                    self.inner.children[i + 1] = other;
                    self.inner.count += 1;
                },

                .inner => {
                    // Child gets 3, Other gets 3, Middle gets promoted to self
                    other.* = .{ .inner = .{ .count = min_keys } };
                    @memcpy(other.inner.keys[0..min_keys], child.inner.keys[M..max_keys]);
                    @memcpy(other.inner.children[0..min_children], child.inner.children[min_children..max_children]);

                    child.inner.count = min_keys;

                    var i = self.inner.count;
                    while (i > index) : (i -= 1) {
                        self.inner.keys[i] = self.inner.keys[i - 1];
                        self.inner.children[i + 1] = self.inner.children[i];
                    }

                    // The key is still in the child, but it's no longer counted as child
                    self.inner.keys[i] = child.inner.keys[M - 1];
                    self.inner.children[i + 1] = other;
                    self.inner.count += 1;
                },
            }

            switch (compare(key, self.inner.keys[index])) {
                .gt, .eq => return self.inner.children[index + 1] orelse unreachable,
                .lt => return self.inner.children[index] orelse unreachable,
            }
        }

        fn move(self: *Self, allocator: mem.Allocator, child: *Self, left: ?*Self, right: ?*Self, index: usize) !*Self {
            _ = allocator;

            // Borrow from left if possible
            if (left) |other| {
                if (!other.underfull()) {
                    switch (child.*) {
                        .leaf => |*node| {
                            var i = node.count;
                            const last = other.leaf.entries[other.leaf.count - 1];
                            while (i > 0) : (i -= 1) node.entries[i] = node.entries[i - 1];

                            node.entries[0] = last;
                            node.count += 1;
                            other.leaf.count -= 1;

                            self.inner.keys[index - 1] = node.entries[0].key;
                        },

                        .inner => |*node| {
                            var i = node.count;
                            while (i > 0) : (i -= 1) node.keys[i] = node.keys[i - 1];
                            node.keys[0] = self.inner.keys[index - 1];

                            const last = other.inner.keys[other.inner.count - 1];
                            self.inner.keys[index - 1] = last;

                            var j = node.count + 1;
                            while (j > 0) : (j -= 1) node.children[j] = node.children[j - 1];
                            node.children[0] = other.inner.children[other.inner.count];
                            other.inner.children[other.inner.count] = null;

                            node.count += 1;
                            other.inner.count -= 1;
                        },
                    }

                    return child;
                }
            }

            // Borrow from right if possible
            if (right) |other| {
                if (!other.underfull()) {
                    switch (child.*) {
                        .inner => |*node| {
                            node.keys[node.count] = self.inner.keys[index];
                            node.children[node.count + 1] = other.inner.children[0];
                            self.inner.keys[index] = other.inner.keys[0];

                            for (1..other.inner.count) |i| other.inner.keys[i - 1] = other.inner.keys[i];
                            for (1..other.inner.count + 1) |i| other.inner.children[i - 1] = other.inner.children[i];
                            other.inner.children[other.inner.count] = null;

                            node.count += 1;
                            other.inner.count -= 1;
                        },
                        .leaf => |*node| {
                            const first = other.leaf.entries[0];
                            for (1..other.leaf.count) |i| other.leaf.entries[i - 1] = other.leaf.entries[i];
                            node.entries[node.count] = first;

                            node.count += 1;
                            other.leaf.count -= 1;

                            self.inner.keys[index] = other.leaf.entries[0].key;
                        },
                    }
                }

                return child;
            }

            // Merge with left if possible
            if (left) |other| {
                switch (child.*) {
                    .inner => |*node| {
                        _ = node;
                    },
                    .leaf => |*node| {
                        _ = node;
                    },
                }
                return other;
            }

            // Merge with right if possible
            if (right) |other| {
                _ = other;

                switch (child.*) {
                    .inner => |*node| {
                        _ = node;
                    },
                    .leaf => |*node| {
                        _ = node;
                    },
                }
                return child;
            }

            unreachable; // This should never happen
        }

        fn traverse(self: *Self) *Self {
            var node = self;
            while (node.* == .inner) node = node.inner.children[0] orelse unreachable;
            return node;
        }

        fn full(self: *Self) bool {
            return switch (self.*) {
                .inner => |*node| node.count == max_keys,
                .leaf => |*node| node.count == max_keys,
            };
        }

        fn underfull(self: *Self) bool {
            return switch (self.*) {
                .inner => |*node| node.count == min_keys,
                .leaf => |*node| node.count == min_keys,
            };
        }

        fn empty(self: *Self) bool {
            return switch (self.*) {
                .inner => |*node| node.count == 0,
                .leaf => |*node| node.count == 0,
            };
        }
    };

    const Iterator = struct {
        const Self = @This();

        index: usize,
        node: ?*Node,

        pub fn next(self: *Self) ?Entry {
            if (self.node) |node| {
                switch (node.*) {
                    .inner => unreachable,
                    .leaf => |leaf| {
                        if (self.index < leaf.count) {
                            const entry = leaf.entries[self.index];
                            self.index += 1;
                            return entry;
                        }

                        if (leaf.next) |next_node| {
                            const entry = next_node.leaf.entries[0];
                            self.node = next_node;
                            self.index = 1;
                            return entry;
                        }

                        self.node = null;
                        self.index = 0;
                        return null;
                    },
                }
            }

            return null;
        }
    };

    return struct {
        const Self = @This();

        root: ?*Node,
        size: u64,

        pub const empty: Self = .{ .root = null, .size = 0 };

        /// Returns the value associated with the given key, if it exists.
        pub fn get(self: *Self, key: K) ?V {
            if (self.root == null) return null;
            const entry = self.root.?.search(key);

            if (entry == null) return null;
            return entry.?.value;
        }

        /// Inserts or updates the value associated with the given key.
        pub fn put(self: *Self, allocator: mem.Allocator, key: K, value: V) !void {
            var root = self.root orelse {
                const node = try allocator.create(Node);
                node.* = .{ .leaf = .{ .count = 1, .next = null } };
                node.leaf.entries[0] = .{ .key = key, .value = value };
                self.root = node;
                self.size = 1;
                return;
            };

            if (root.full()) {
                const node = try allocator.create(Node);
                node.* = .{ .inner = .{ .count = 0 } };
                node.inner.children[0] = root;
                root = try node.split(allocator, root, 0, key);
                self.root = node;
            }

            const inserted = try root.insert(allocator, key, value);
            if (inserted) self.size += 1;
        }

        /// Removes the value associated with the given key, if it exists.
        pub fn rid(self: *Self, allocator: mem.Allocator, key: K) !void {
            const root = self.root orelse return;

            const deleted = try root.delete(allocator, key);
            if (!deleted) return;
            self.size -= 1;

            if (root.empty()) {
                switch (root.*) {
                    .leaf => {
                        allocator.destroy(root);
                        self.root = null;
                        self.size = 0;
                    },
                    .inner => |*inner| {
                        const node = inner.children[0] orelse return;
                        allocator.destroy(root);
                        self.root = node;
                    },
                }
            }
        }

        /// Returns whether the map contains the given key.
        pub fn has(self: Self, key: K) bool {
            if (self.root == null) return false;
            return self.root.?.search(key) != null;
        }

        /// Returns an iterator over the map's entries.
        pub fn iter(self: *Self) Iterator {
            if (self.root) |root| {
                const first = root.traverse();
                return .{ .node = first, .index = 0 };
            } else return .{ .node = null, .index = 0 };
        }
    };
}
