const std = @import("std");

const vaxis = @import("vaxis");
const vxfw = vaxis.vxfw;

pub const Candle = struct {
    high: f32,
    open: f32,
    close: f32,
    low: f32,
};

pub const Chart = @This();

grid: vxfw.SizedBox = undefined,
data: *[]Candle = undefined,

pub fn widget(self: *Chart) vxfw.Widget {
    return .{
        .userdata = self,
        .drawFn = typeErasedDrawFn,
    };
}

fn typeErasedDrawFn(ptr: *anyopaque, ctx: vxfw.DrawContext) std.mem.Allocator.Error!vxfw.Surface {
    const self: *Chart = @ptrCast(@alignCast(ptr));
    const size: vxfw.Size = .{ .height = 100, .width = 120 };
    return try vxfw.Surface.init(ctx.arena, self.grid.widget(), size);
}

pub fn generateCandles(allocator: std.mem.Allocator, amount: usize) !*[]Candle {
    var dummy_candles = try allocator.alloc(Candle, amount);
    var prng = std.Random.DefaultPrng.init(69);

    const random = prng.random();
    var last_close: f32 = 100.0;

    for (dummy_candles) |*c| {
        const change = (random.float(f32) - 0.45) * 10.0;
        c.open = last_close;
        c.close = last_close + change;

        c.high = @max(c.open, c.close) + random.float(f32) * 2.0;
        c.low = @min(c.open, c.close) - random.float(f32) * 2.0;
        last_close = c.close;
    }

    return &dummy_candles;
}
