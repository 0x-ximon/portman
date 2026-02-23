const std = @import("std");

const vaxis = @import("vaxis");
const vxfw = vaxis.vxfw;

const AccountScreen = @This();

allocator: std.mem.Allocator,
content: vxfw.Text = .{ .text = "👤 Account Screen\n\n- Profit/Loss\n- Open Trades" },

pub fn init(allocator: std.mem.Allocator) !*AccountScreen {
    const self = try allocator.create(AccountScreen);
    self.allocator = allocator;
    return self;
}

pub fn deinit(self: *AccountScreen) void {
    self.allocator.destroy(self);
}

pub fn widget(self: *AccountScreen) vxfw.Widget {
    return self.content.widget();
}
