const std = @import("std");

const vaxis = @import("vaxis");
const vxfw = vaxis.vxfw;

const CartesianPlane = @import("../cartesian.zig");
const Global = @import("../global.zig");

pub const Candle = struct {
    high: f32,
    open: f32,
    close: f32,
    low: f32,
};

// TODO: Delete this struct as Candles should be streamed from a source
pub const Candles = struct {
    items: []Candle,

    fn init(allocator: std.mem.Allocator, amount: usize) !*Candles {
        const self = try allocator.create(Candles);
        self.items = try allocator.alloc(Candle, amount);

        var seed: u64 = undefined;
        try std.posix.getrandom(std.mem.asBytes(&seed));
        var prng = std.Random.DefaultPrng.init(seed);

        const random = prng.random();
        var last_close: f32 = 100.0;

        for (self.items) |*c| {
            const percent = (random.float(f32) - 0.5) * 10.0;
            const change = percent * last_close / 100.0;

            c.open = last_close;
            c.close = c.open + change;

            c.high = @max(c.open, c.close) +
                if (random.boolean()) (random.float(f32) * last_close * 0.05) else 0;

            c.low = @min(c.open, c.close) -
                if (random.boolean()) (random.float(f32) * last_close * 0.05) else 0;

            last_close = c.close;
        }

        return self;
    }

    pub fn deinit(self: *Candles, allocator: std.mem.Allocator) void {
        allocator.free(self.items);
        allocator.destroy(self);
    }
};

pub const Chart = @This();

allocator: std.mem.Allocator,
candles: *Candles,

pub fn init(allocator: std.mem.Allocator) !*Chart {
    const self = try allocator.create(Chart);
    self.allocator = allocator;

    self.candles = try Candles.init(allocator, 200);
    return self;
}

pub fn deinit(self: *Chart) void {
    self.candles.deinit(self.allocator);
    self.allocator.destroy(self);
}

pub fn widget(self: *Chart) vxfw.Widget {
    return .{
        .userdata = self,
        .drawFn = Chart.typeErasedDrawFn,
        .eventHandler = Chart.typeErasedEventHandler,
    };
}

fn typeErasedDrawFn(ptr: *anyopaque, ctx: vxfw.DrawContext) std.mem.Allocator.Error!vxfw.Surface {
    const self: *Chart = @ptrCast(@alignCast(ptr));

    const width = ctx.max.width orelse ctx.min.width;
    const height = ctx.max.height orelse @max(ctx.min.height, 10);
    const actual_size = vxfw.Size{ .width = width, .height = height };

    const surface = try vxfw.Surface.init(ctx.arena, self.widget(), actual_size);

    // Calculate Bounds
    var max_val: f32 = 0;
    var min_val: f32 = std.math.floatMax(f32);
    for (self.candles.items) |candle| {
        max_val = @max(max_val, candle.high);
        min_val = @min(min_val, candle.low);
    }

    // Define Margins and Spacing
    const candle_spacing = 1;
    const candle_width = 1;

    const stride = candle_width + candle_spacing;
    const max_visible_candles = width / stride;

    // Initialize our new helper
    const plane = CartesianPlane{
        .min_x = 0,
        .max_x = @floatFromInt(max_visible_candles),
        .min_y = min_val,
        .max_y = max_val,
        .size = actual_size,
    };

    // 3. Render Candles with gaps
    for (0..@min(self.candles.items.len, max_visible_candles)) |i| {
        const candle = self.candles.items[i];
        const bullish = candle.close >= candle.open;

        // Calculate X based on stride
        const screen_x = @as(u16, @intCast(i)) * stride;

        // Map Y coordinates using the Plane helper
        const y_high = plane.toScreen(0, candle.high).row;
        const y_low = plane.toScreen(0, candle.low).row;
        const y_open = plane.toScreen(0, candle.open).row;
        const y_close = plane.toScreen(0, candle.close).row;

        const y_max_body = @min(y_open, y_close);
        const y_min_body = @max(y_open, y_close);

        const color_idx = if (bullish) Global.Color.blue else Global.Color.bright_red;
        const cell_color = vaxis.Color{ .index = @intFromEnum(color_idx) };

        // Draw Wick (ensure y_high to y_low is drawn)
        var y: u16 = y_high;
        while (y <= y_low) : (y += 1) {
            surface.writeCell(screen_x, y, .{
                .char = .{ .grapheme = "│" },
                .style = .{ .fg = cell_color },
            });
        }

        // Draw Body (█)
        y = y_max_body;
        while (y <= y_min_body) : (y += 1) {
            surface.writeCell(screen_x, y, .{
                .char = .{ .grapheme = "█" },
                .style = .{ .fg = cell_color },
            });
        }
    }

    return surface;
}

fn typeErasedEventHandler(ptr: *anyopaque, ctx: *vxfw.EventContext, event: vxfw.Event) anyerror!void {
    const self: *Chart = @ptrCast(@alignCast(ptr));
    _ = event; // autofix
    _ = self; // autofix
    _ = ctx; // autofix
}

fn valToY(val: f32, min: f32, range: f32, height: u16) u16 {
    if (range == 0) return height / 2;
    const pct = (val - min) / range;
    const y = @as(f32, @floatFromInt(height - 1)) * (1.0 - pct);
    return @intFromFloat(@round(y));
}
