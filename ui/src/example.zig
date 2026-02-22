const std = @import("std");
const vaxis = @import("vaxis");
const vxfw = vaxis.vxfw;

pub const Candle = struct {
    open: f32,
    high: f32,
    low: f32,
    close: f32,
};

pub const Chart = struct {
    data: []const Candle,

    pub fn widget(self: Chart) vxfw.Widget {
        return .{
            .userdata = @constCast(&self),
            .drawFn = drawFn,
        };
    }

    fn drawFn(ptr: *anyopaque, ctx: vxfw.DrawContext) !vxfw.Surface {
        const self: *Chart = @ptrCast(@alignCast(ptr));
        const width = ctx.max.width orelse 20;
        const height = ctx.max.height orelse 10;

        var surface = try vxfw.Surface.init(ctx.allocator, width, height);

        if (self.data.len == 0) return surface;

        // Calculate scaling
        var max_val: f32 = 0;
        var min_val: f32 = 1e10;
        for (self.data) |c| {
            max_val = @max(max_val, c.high);
            min_val = @min(min_val, c.low);
        }
        const range = max_val - min_val;

        // Draw candles (one per column, up to width)
        for (0..@min(width, self.data.len)) |x| {
            const candle = self.data[x];
            const is_up = candle.close >= candle.open;

            // Map prices to Y coordinates (top is 0)
            const y_high = valToY(candle.high, min_val, range, height);
            const y_low = valToY(candle.low, min_val, range, height);
            const y_max_body = valToY(@max(candle.open, candle.close), min_val, range, height);
            const y_min_body = valToY(@min(candle.open, candle.close), min_val, range, height);

            const color: vaxis.Cell.Color = if (is_up) .green else .red;

            // Draw Wick
            for (y_high..y_low + 1) |y| {
                if (y >= height) continue;
                surface.writeCell(@intCast(x), @intCast(y), .{
                    .char = .{ .grapheme = "│" },
                    .style = .{ .fg = color },
                });
            }

            // Draw Body
            for (y_max_body..y_min_body + 1) |y| {
                if (y >= height) continue;
                surface.writeCell(@intCast(x), @intCast(y), .{
                    .char = .{ .grapheme = "█" },
                    .style = .{ .fg = color },
                });
            }
        }

        return surface;
    }

    fn valToY(val: f32, min: f32, range: f32, height: u16) u16 {
        if (range == 0) return height / 2;
        const pct = (val - min) / range;
        const y = @as(f32, @floatFromInt(height - 1)) * (1.0 - pct);
        return @intFromFloat(@round(y));
    }
};
