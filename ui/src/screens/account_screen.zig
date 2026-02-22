const std = @import("std");

const vaxis = @import("vaxis");
const vxfw = vaxis.vxfw;

const AccountScreen = @This();

content: vxfw.Text = .{ .text = "👤 Account Screen\n\n- Profit/Loss\n- Open Trades" },

pub fn widget(self: *AccountScreen) vxfw.Widget {
    return self.content.widget();
}
