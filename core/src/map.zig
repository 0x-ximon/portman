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

    const Node = union((enum { internal, leaf })) {
        const Self = @This();

        internal: Internal,
        leaf: Leaf,

        const Internal = struct {
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
                .internal => |*node| {
                    var i: usize = 0;
                    while (i < node.count) : (i += 1) {
                        switch (compare(key, node.keys[i])) {
                            .gt => continue,
                            .lt => break,
                            .eq => break,
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
                .leaf => |*leaf| {
                    var i: usize = 0;
                    while (i < leaf.count) : (i += 1) {
                        switch (compare(key, leaf.entries[i].key)) {
                            .gt => continue,
                            .lt => break,
                            .eq => {
                                leaf.entries[i].value = value;
                                return false;
                            },
                        }
                    }

                    var j = leaf.count;
                    while (j > i) : (j -= 1) leaf.entries[j] = leaf.entries[j - 1];
                    leaf.entries[i] = .{ .key = key, .value = value };
                    leaf.count += 1;
                    return true;
                },

                .internal => |*internal| {
                    var i: usize = 0;
                    while (i < internal.count) : (i += 1) {
                        switch (compare(key, internal.keys[i])) {
                            .gt, .eq => continue,
                            .lt => break,
                        }
                    }

                    var child = internal.children[i] orelse unreachable;
                    if (child.full()) {
                        try self.split(allocator, child, i);
                        switch (compare(key, internal.keys[i])) {
                            // If the key is equal to or greater than current key, move to next child
                            .gt, .eq => {
                                i += 1;
                                child = internal.children[i] orelse unreachable;
                            },
                            // If the key is less than current key, stay on current child
                            .lt => {},
                        }
                    }

                    return try child.insert(allocator, key, value);
                },
            }

            return false;
        }

        fn delete(self: *Self, allocator: mem.Allocator, key: K) !bool {
            _ = self;
            _ = allocator;
            _ = key;

            return false;
        }

        fn split(self: *Self, allocator: mem.Allocator, child: *Self, index: usize) !void {
            const other = try allocator.create(Self);

            switch (child.*) {
                .leaf => |*leaf| {
                    // Child gets 4, Other gets 3, First of Other gets promoted to self
                    other.* = .{ .leaf = .{ .count = min_keys, .next = leaf.next } };
                    @memcpy(other.leaf.entries[0..min_keys], leaf.entries[M..max_keys]);

                    leaf.count = M;
                    leaf.next = other;

                    var i = self.internal.count;
                    while (i > index) : (i -= 1) {
                        self.internal.keys[i] = self.internal.keys[i - 1];
                        self.internal.children[i + 1] = self.internal.children[i];
                    }

                    self.internal.keys[i] = other.leaf.entries[0].key;
                    self.internal.children[i + 1] = other;
                    self.internal.count += 1;
                },

                .internal => {
                    // Child gets 3, Other gets 3, Middle gets promoted to self
                    other.* = .{ .internal = .{ .count = min_keys } };
                    @memcpy(other.internal.keys[0..min_keys], child.internal.keys[M..max_keys]);
                    @memcpy(other.internal.children[0..min_children], child.internal.children[min_children..max_children]);

                    child.internal.count = min_keys;

                    var i = self.internal.count;
                    while (i > index) : (i -= 1) {
                        self.internal.keys[i] = self.internal.keys[i - 1];
                        self.internal.children[i + 1] = self.internal.children[i];
                    }

                    // The key is still in the child, but it's no longer counted as child
                    self.internal.keys[i] = child.internal.keys[M - 1];
                    self.internal.children[i + 1] = other;
                    self.internal.count += 1;
                },
            }
        }

        fn merge(self: *Self, allocator: mem.Allocator, child: *Self, index: usize) !void {
            _ = self;
            _ = child;
            _ = index;
            _ = allocator;
        }

        fn traverse(self: *Self) *Self {
            var node = self;
            while (node.* == .internal) node = node.internal.children[0].?;
            return node;
        }

        fn full(self: *Self) bool {
            return switch (self.*) {
                .internal => |*node| node.count == max_keys,
                .leaf => |*node| node.count == max_keys,
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
                    .internal => unreachable,
                    .leaf => |leaf| {
                        if (self.index < leaf.count) {
                            const entry = leaf.entries[self.index];
                            self.index += 1;
                            return entry;
                        }

                        if (leaf.next) |next_node| {
                            self.node = next_node;
                            self.index = 0;
                            return next_node.leaf.entries[0];
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

        pub fn put(self: *Self, allocator: mem.Allocator, key: K, value: V) !void {
            if (self.root) |root| {
                if (root.full()) {
                    const node = try allocator.create(Node);
                    node.* = .{ .internal = .{ .count = 1 } };
                    node.internal.children[0] = root;

                    try node.split(allocator, root, 0);
                    self.root = node;
                }

                const inserted = try self.root.?.insert(allocator, key, value);
                if (inserted) self.size += 1;
            } else {
                const node = try allocator.create(Node);
                node.* = .{ .leaf = .{ .count = 1, .next = null } };
                node.leaf.entries[0] = .{ .key = key, .value = value };
                self.root = node;
                self.size = 1;
                return;
            }
        }

        pub fn iterator(self: *Self) Iterator {
            if (self.root) |root| {
                const first = root.traverse();
                return .{ .node = first, .index = 0 };
            } else return .{ .node = null, .index = 0 };
        }

        //     pub fn remove(self: *Self, allocator: mem.Allocator, key: K) !void {
        //         const root = self.root orelse return;

        //         const deleted = try root.delete(allocator, key);
        //         if (!deleted) return;
        //         self.size -= 1;

        //         if (root.count == 0) {
        //             switch (root.is_leaf) {
        //                 true => {
        //                     allocator.destroy(root);
        //                     self.root = null;
        //                     self.size = 0;
        //                 },
        //                 false => {
        //                     const node = root.children[0] orelse return;
        //                     allocator.destroy(root);
        //                     self.root = node;
        //                 },
        //             }
        //         }
        //     }

        //     const Node = struct {
        //         count: usize = 0,
        //         is_leaf: bool = false,
        //         entries: [max_entries]Entry = undefined,
        //         children: [max_children]?*Node = .{null} ** max_children,

        //         fn delete(self: *Node, allocator: mem.Allocator, key: K) !bool {
        //             // INVARIANT: All minimum nodes must have been merged before deletion
        //             switch (self.is_leaf) {
        //                 true => {
        //                     for (0..self.count) |i| {
        //                         switch (compare(key, self.entries[i].key)) {
        //                             .lt => break,
        //                             .gt => continue,
        //                             .eq => {
        //                                 for (i..self.count - 1) |j| self.entries[j] = self.entries[j + 1];
        //                                 self.count -= 1;
        //                                 return true;
        //                             },
        //                         }

        //                         return false;
        //                     }
        //                 },

        //                 false => {
        //                     var i: usize = 0;
        //                     while (i < self.count) : (i += 1) {
        //                         switch (compare(key, self.entries[i].key)) {
        //                             .lt => break,
        //                             .gt => continue,
        //                             .eq => {
        //                                 // TODO: Handle internal node deletion
        //                             },
        //                         }
        //                     }

        //                     var child = self.children[i] orelse unreachable;
        //                     if (child.count == min_entries) {
        //                         const left = if (i == 0) null else self.children[i - 1];
        //                         const right = if (i == self.count) null else self.children[i + 1];
        //                         try self.merge(allocator, child, left, right, i);
        //                     }

        //                     return try child.delete(allocator, key);
        //                 },
        //             }

        //             return false;
        //         }

        //         fn merge(self: *Node, allocator: mem.Allocator, child: *Node, left: ?*Node, right: ?*Node, index: usize) !void {
        //             _ = self;
        //             _ = child;
        //             _ = index;
        //             _ = allocator;

        //             // Borrow from left sibling if possible
        //             if (left != null and left.?.count > min_entries) {
        //                 const other = left.?;
        //                 _ = other;
        //                 return;
        //             }

        //             // Borrow from right sibling if possible
        //             if (right != null and right.?.count > min_entries) {
        //                 const other = right.?;
        //                 _ = other;
        //                 return;
        //             }

        //             // Merge with left sibling if possible
        //             // Merge with right sibling if possible
        //         }
        //     };
    };
}
