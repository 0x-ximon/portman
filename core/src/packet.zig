const std = @import("std");
const Io = std.Io;

const Book = @import("book.zig");
const Order = Book.Order;

const Packet = @This();

pub const Instruction = enum(u8) {
    NewOrder = 0,
    CancelOrder = 1,
    UpdateOrder = 2,
};

pub const Header = extern struct {
    version: u8,
    instruction: Instruction,
    length: u16,
    nonce: u32,
    timestamp: u64,
    src_id: u64,
    dest_id: u64,
};

pub fn info(reader: *Io.Reader) !*const Header {
    const bytes = try reader.take(@sizeOf(Header));
    const header: *const Header = @ptrCast(@alignCast(bytes.ptr));
    return header;
}

pub fn recv(reader: *Io.Reader, length: usize) ![]Order {
    const bytes = try reader.take(length);
    const count = length / @sizeOf(Order);
    const orders = @as([*]Order, @ptrCast(@alignCast(bytes.ptr)))[0..count];
    return orders;
}

pub fn send(writer: Io.Writer) void {
    _ = writer;
}
