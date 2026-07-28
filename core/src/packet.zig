const std = @import("std");
const Io = std.Io;

const Book = @import("book.zig");
const Order = Book.Order;

const Packet = @This();

pub const Instruction = enum(u8) {
    ProcessOrder = 0,
    CancelOrder = 1,
    UpdateOrder = 2,
};

pub const Header = extern struct {
    version: u8,
    instruction: Instruction,
    length: u16,
    nonce: u32,
    timestamp: u64,
    source: u64,
    destination: u64,
};

pub fn recv(bytes: []const u8) struct { *const Header, []Order } {
    const size = @sizeOf(Header);
    const header: *const Header = @ptrCast(@alignCast(bytes[0..size]));
    const orders: []Order = @ptrCast(@alignCast(@constCast(bytes[size..])));
    return .{ header, orders };
}
