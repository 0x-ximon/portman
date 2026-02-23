const std = @import("std");

const vaxis = @import("vaxis");
const vxfw = vaxis.vxfw;

const AccountScreen = @import("../screens/account_screen.zig");
const HomeScreen = @import("../screens/home_screen.zig");
const SettingsScreen = @import("../screens/settings_screen.zig");

// Tentative UI Screens
//  - Home
//      - Candlestick Chart
//      - Order Book / Market Depth
//      - Technical Indicator / Trading Volume
//      - Market Changes / News Feed
//
//  - Account
//      - Profit/Loss
//      - Open Trades
//      - Trades History
//      - Balances
//
//  - Config
//      - Theme
//      - Network

const Screen = union(enum) {
    home: HomeScreen,
    account: AccountScreen,
    settings: SettingsScreen,
};

pub const Router = @This();

allocator: std.mem.Allocator,
active: Screen,

pub fn init(allocator: std.mem.Allocator) !*Router {
    const self = try allocator.create(Router);
    self.allocator = allocator;

    self.active = .{ .home = .{} };
    return self;
}

pub fn deinit(self: *Router) void {
    self.allocator.destroy(self);
}

pub fn widget(self: *Router) vxfw.Widget {
    return .{
        .userdata = self,
        .drawFn = typeErasedDrawFn,
        .eventHandler = typeErasedEventHandler,
    };
}

fn typeErasedDrawFn(ptr: *anyopaque, ctx: vxfw.DrawContext) std.mem.Allocator.Error!vxfw.Surface {
    const self: *Router = @ptrCast(@alignCast(ptr));
    return switch (self.active) {
        .home => |*screen| screen.widget().draw(ctx),
        .account => |*screen| screen.widget().draw(ctx),
        .settings => |*screen| screen.widget().draw(ctx),
    };
}

fn typeErasedEventHandler(ptr: *anyopaque, ctx: *vxfw.EventContext, event: vxfw.Event) anyerror!void {
    const self: *Router = @ptrCast(@alignCast(ptr));
    // Pass events down to the active screen so it can handle its own inputs
    try switch (self.active) {
        .home => |*screen| screen.widget().handleEvent(ctx, event),
        .account => |*screen| screen.widget().handleEvent(ctx, event),
        .settings => |*screen| screen.widget().handleEvent(ctx, event),
    };
}
