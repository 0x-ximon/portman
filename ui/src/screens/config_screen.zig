const std = @import("std");

const vaxis = @import("vaxis");
const vxfw = vaxis.vxfw;

const ConfigScreen = @This();

content: vxfw.Text = .{ .text = "⚙️ Configuration Screen\n\n- Theme: Dark\n- Network: Mainnet" },

pub fn widget(self: *ConfigScreen) vxfw.Widget {
    return self.content.widget();
}
