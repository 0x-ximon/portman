const std = @import("std");

const vaxis = @import("vaxis");
const vxfw = vaxis.vxfw;

const SettingsScreen = @This();

allocator: std.mem.Allocator,
content: vxfw.Text = .{ .text = "⚙️ Configuration Screen\n\n- Theme: Dark\n- Network: Mainnet" },

pub fn init(allocator: std.mem.Allocator) !*SettingsScreen {
    const self = try allocator.create(SettingsScreen);
    self.allocator = allocator;
    return self;
}

pub fn deinit(self: *SettingsScreen) void {
    self.allocator.destroy(self);
}

pub fn widget(self: *SettingsScreen) vxfw.Widget {
    return self.content.widget();
}
