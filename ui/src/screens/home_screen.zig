const std = @import("std");

const vaxis = @import("vaxis");
const vxfw = vaxis.vxfw;

const HomeScreen = @This();

content: vxfw.Text = .{ .text = "🏠 Home Screen\n\n- Candlestick Chart\n- Order Book" },

pub fn widget(self: *HomeScreen) vxfw.Widget {
    return self.content.widget();
}
