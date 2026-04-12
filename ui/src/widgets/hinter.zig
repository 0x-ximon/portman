const std = @import("std");

const zz = @import("zigzag");
const lib = @import("lib");

pub const Hinter = @This();

content: []const u8,
allocator: std.mem.Allocator,

pub const Msg = union(enum) {
    click: zz.MouseEvent,
    key: zz.KeyEvent,
};

pub fn init(allocator: std.mem.Allocator) !*Hinter {
    const self = try allocator.create(Hinter);
    self.* = .{
        .content = "Hinter",
        .allocator = allocator,
    };
    return self;
}

pub fn deinit(self: *Hinter) void {
    self.allocator.destroy(self);
}

pub fn react(_: *Hinter, _: Msg) void {}

pub fn render(self: *Hinter, allocator: std.mem.Allocator, width: u16, height: u16) []const u8 {
    const style = (zz.Style{})
        .borderAll(zz.Border.rounded)
        .fg(zz.Color.black())
        .width(width - 2)
        .height(height - 2);

    return style.render(allocator, self.content) catch "hinter render failed";
}

// pub fn update(_: *Hinter, msg: Msg, _: *zz.Context) zz.Cmd(Msg) {
//     switch (msg) {
//         .key => |k| switch (k.key) {
//             .char => |c| switch (c) {
//                 else => {},
//             },
//             else => {},
//         },

//         .click => {},
//     }

//     return .none;
// }

// pub fn view(_: *const Hinter, ctx: *const zz.Context) []const u8 {
//     const content = std.fmt.allocPrint(ctx.allocator, "Width: {d}\tHeight: {d}", .{
//         ctx.width, ctx.height,
//     }) catch "Error";

//     return content;
// }
