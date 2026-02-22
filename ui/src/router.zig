const std = @import("std");

const vaxis = @import("vaxis");
const vxfw = vaxis.vxfw;

const AccountScreen = @import("screens/account_screen.zig");
const HomeScreen = @import("screens/home_screen.zig");
const ConfigScreen = @import("screens/config_screen.zig");

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
    config: ConfigScreen,
};

pub const Router = @This();

active: Screen,

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
        .config => |*screen| screen.widget().draw(ctx),
    };
}

fn typeErasedEventHandler(ptr: *anyopaque, ctx: *vxfw.EventContext, event: vxfw.Event) anyerror!void {
    const self: *Router = @ptrCast(@alignCast(ptr));
    // Pass events down to the active screen so it can handle its own inputs
    try switch (self.active) {
        .home => |*screen| screen.widget().handleEvent(ctx, event),
        .account => |*screen| screen.widget().handleEvent(ctx, event),
        .config => |*screen| screen.widget().handleEvent(ctx, event),
    };
}
