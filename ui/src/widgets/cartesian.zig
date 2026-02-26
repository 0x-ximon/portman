const std = @import("std");

const vaxis = @import("vaxis");
const vxfw = vaxis.vxfw;

pub const CartesianPlane = @This();

min_x: f32,
max_x: f32,
min_y: f32,
max_y: f32,
size: vxfw.Size,

/// Translates a data point (x, y) to a screen coordinate (col, row)
pub fn toScreen(self: CartesianPlane, x: f32, y: f32) struct { col: u16, row: u16 } {
    const range_x = self.max_x - self.min_x;
    const range_y = self.max_y - self.min_y;

    const pct_x = if (range_x == 0) 0.5 else (x - self.min_x) / range_x;
    const pct_y = if (range_y == 0) 0.5 else (y - self.min_y) / range_y;

    // X grows left to right
    const screen_x = @as(f32, @floatFromInt(self.size.width - 1)) * pct_x;
    // Y grows top to bottom in terminal, so we invert it (1.0 - pct_y)
    const screen_y = @as(f32, @floatFromInt(self.size.height - 1)) * (1.0 - pct_y);

    return .{
        .col = @intFromFloat(@round(screen_x)),
        .row = @intFromFloat(@round(screen_y)),
    };
}
