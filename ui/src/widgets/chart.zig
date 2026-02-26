const std = @import("std");

const vaxis = @import("vaxis");
const vxfw = vaxis.vxfw;

const Color = enum(u24) {
    red = 0xFF0000,
    green = 0x00FF00,
    blue = 0x0000FF,
    white = 0xFFFFFF,
    black = 0x000000,
};

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

        var prng = std.Random.DefaultPrng.init(69);

        const random = prng.random();
        var last_close: f32 = 100.0;

        for (self.items) |*c| {
            const change = (random.float(f32) - 0.45) * 10.0;
            c.open = last_close;
            c.close = last_close + change;

            c.high = @max(c.open, c.close) + random.float(f32) * 2.0;
            c.low = @min(c.open, c.close) - random.float(f32) * 2.0;
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

chart_content: vxfw.Text,
chart_border: vxfw.Border,
chart_widget: vxfw.Widget,

pub fn init(allocator: std.mem.Allocator) !*Chart {
    const self = try allocator.create(Chart);
    self.allocator = allocator;

    self.candles = try Candles.init(allocator, 100);

    // Chart Initialization
    self.chart_content = .{ .text = 
        \\ Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et
        \\ dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut
        \\ aliquip ex ea commodo consequat.
    };

    self.chart_border = .{
        .labels = &.{.{ .text = "Ticker Chart", .alignment = .top_left }},
        .child = self.chart_content.widget(),
    };

    self.chart_widget = self.chart_border.widget();

    return self;
}

pub fn deinit(self: *Chart) void {
    self.candles.deinit(self.allocator);
    self.allocator.destroy(self);
}

pub fn widget(self: *Chart) vxfw.Widget {
    // return self.chart_widget;

    return .{
        .userdata = self,
        .drawFn = Chart.typeErasedDrawFn,
        .eventHandler = Chart.typeErasedEventHandler,
    };
}

fn typeErasedDrawFn(ptr: *anyopaque, ctx: vxfw.DrawContext) std.mem.Allocator.Error!vxfw.Surface {
    const self: *Chart = @ptrCast(@alignCast(ptr));
    const max_size = ctx.max.size();

    const surface = try vxfw.Surface.init(ctx.arena, self.widget(), max_size);
    // const subSurfaces = try ctx.arena.alloc(vxfw.SubSurface, 1);

    // Calculate Scaling
    var max_val: f32 = 0;
    var min_val: f32 = std.math.floatMax(f32);
    for (self.candles.items) |candle| {
        max_val = @max(max_val, candle.high);
        min_val = @min(min_val, candle.low);
    }
    const range = max_val - min_val;

    for (0..@min(self.candles.items.len, 20)) |x| {
        const candle = self.candles.items[x];
        const bullish = candle.close >= candle.open;

        // Map prices to Y coordinates (top is 0)
        const y_high = valToY(candle.high, min_val, range, max_size.height);
        const y_low = valToY(candle.low, min_val, range, max_size.height);
        const y_max_body = valToY(@max(candle.open, candle.close), min_val, range, max_size.height);
        const y_min_body = valToY(@min(candle.open, candle.close), min_val, range, max_size.height);

        const color = if (bullish)
            vaxis.Color.rgbFromUint(@intFromEnum(Chart.Color.green))
        else
            vaxis.Color.rgbFromUint(@intFromEnum(Chart.Color.red));

        // Draw Wick
        for (y_high..y_low + 1) |y| {
            if (y >= max_size.height) continue;
            surface.writeCell(@intCast(x), @intCast(y), .{
                .char = .{ .grapheme = "│" },
                .style = .{ .fg = color },
            });
        }

        // Draw Body
        for (y_max_body..y_min_body + 1) |y| {
            if (y >= max_size.height) continue;
            surface.writeCell(@intCast(x), @intCast(y), .{
                .char = .{ .grapheme = "█" },
                .style = .{ .fg = color },
            });
        }
    }

    return surface;
}

fn typeErasedEventHandler(ptr: *anyopaque, ctx: *vxfw.EventContext, event: vxfw.Event) anyerror!void {
    _ = ctx; // autofix
    const self: *Chart = @ptrCast(@alignCast(ptr));
    _ = self; // autofix
    switch (event) {
        else => {},
    }
}

fn valToY(val: f32, min: f32, range: f32, height: u16) u16 {
    if (range == 0) return height / 2;
    const pct = (val - min) / range;
    const y = @as(f32, @floatFromInt(height - 1)) * (1.0 - pct);
    return @intFromFloat(@round(y));
}
