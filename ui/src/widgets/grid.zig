const std = @import("std");

const vaxis = @import("vaxis");
const vxfw = vaxis.vxfw;

const Global = @import("../global.zig");

pub const Grid = @This();

allocator: std.mem.Allocator,

pub fn init(allocator: std.mem.Allocator) !*Grid {
    const self = try allocator.create(Grid);
    self.allocator = allocator;
    return self;
}

pub fn deinit(self: *Grid) void {
    self.allocator.destroy(self);
}

pub fn widget(self: *Grid) vxfw.Widget {
    return .{
        .userdata = self,
        .drawFn = Grid.typeErasedDrawFn,
    };
}

fn typeErasedDrawFn(ptr: *anyopaque, ctx: vxfw.DrawContext) std.mem.Allocator.Error!vxfw.Surface {
    const self: *Grid = @ptrCast(@alignCast(ptr));

    const width = ctx.max.width orelse ctx.min.width;
    const height = ctx.max.height orelse @max(ctx.min.height, 20);
    const actual_size = vxfw.Size{ .width = width, .height = height };
    const grid_color = vaxis.Color{ .index = @intFromEnum(Global.Color.bright_black) };

    const surface = try vxfw.Surface.init(ctx.arena, self.widget(), actual_size);

    var col: u16 = 0;
    while (col < width) : (col += 10) {
        var row: u16 = 0;
        while (row < height) : (row += 1) {
            surface.writeCell(col, row, .{
                .char = .{ .grapheme = "." },
                .style = .{ .fg = grid_color },
            });
        }
    }

    var row: u16 = 0;
    while (row < height) : (row += 5) {
        var c: u16 = 0;
        while (c < width) : (c += 1) {
            surface.writeCell(c, row, .{
                .char = .{ .grapheme = "." },
                .style = .{ .fg = grid_color },
            });
        }
    }

    return surface;
}

fn typeErasedEventHandler(ptr: *anyopaque, ctx: *vxfw.EventContext, event: vxfw.Event) anyerror!void {
    const self: *Grid = @ptrCast(@alignCast(ptr));
    _ = ctx; // autofix
    _ = event; // autofix
    _ = self; // autofix
}
