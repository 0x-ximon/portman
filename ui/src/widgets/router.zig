const std = @import("std");

const zz = @import("zigzag");
const lib = @import("lib");

pub const Router = @This();

content: []const u8,
allocator: std.mem.Allocator,

pub const Msg = union(enum) {
    click: zz.MouseEvent,
    key: zz.KeyEvent,
};

pub fn init(allocator: std.mem.Allocator) !*Router {
    const self = try allocator.create(Router);
    self.* = .{
        .content = "Router",
        .allocator = allocator,
    };
    return self;
}

pub fn deinit(self: *Router) void {
    self.allocator.destroy(self);
}

pub fn react(_: *Router, _: Msg) void {}

pub fn render(self: *Router, allocator: std.mem.Allocator, width: u16, height: u16) []const u8 {
    const style = (zz.Style{})
        .borderAll(zz.Border.rounded)
        .fg(zz.Color.blue())
        .width(width - 2)
        .height(height - 2);

    return style.render(allocator, self.content) catch "router render failed";
}

// const router_view = self.router.view(
//     ctx,
// );
// const router = router_style.render(allocator, router_view) catch "router style render failed";

// const hinter = hinter_style.render(allocator, self.hinter.view(ctx)) catch "hinter style render failed";

// active: Panels,
// navigator: zz.TabGroup,

// home_screen: ScreenA,
// info_screen: ScreenB,

// const Panels = union(enum) {
//     dashboard,
//     account,
//     monitor,
//     information,
// };

// .active = .dashboard,
// self.home_screen = .{};
// self.info_screen = .{};

// self.navigator = zz.TabGroup.init(ctx.persistent_allocator);

// _ = self.navigator.addTab(.{
//     .id = "home",
//     .title = "Home",
//     .route = .{
//         .ctx = &self.home_screen,
//         .key_fn = ScreenA.onKey,
//         .render_fn = ScreenA.render,
//         .on_enter_fn = ScreenA.onEnter,
//     },
// }) catch {};

// _ = self.navigator.addTab(.{
//     .id = "info",
//     .title = "Info",
//     .route = .{
//         .ctx = &self.info_screen,
//         .key_fn = ScreenB.onKey,
//         .render_fn = ScreenB.render,
//     },
// }) catch {};

// return .none;

// pub fn update(self: *Router, msg: Msg, _: *zz.Context) zz.Cmd(Msg) {
//     switch (msg) {
//         .key => |k| switch (k.key) {
//             .char => |c| switch (c) {
//                 '1' => self.active = .dashboard,
//                 '2' => self.active = .account,
//                 '3' => self.active = .monitor,
//                 '4' => self.active = .information,
//                 else => {},
//             },
//             else => {},
//         },

//         .click => {},
//     }

//     return .none;
// }

// pub fn view(_: *const Router, ctx: *const zz.Context, width: u16, height: u16) []const u8 {
//     const style = (zz.Style{})
//         .borderAll(zz.Border.rounded)
//         .fg(zz.Color.black())
//         .width(width)
//         .height(height);

//     const content = std.fmt.allocPrint(ctx.allocator, "Width: {d}\tHeight: {d}", .{
//         ctx.width, ctx.height,
//     }) catch "Error";

//     style.render(ctx.allocator, content) catch "router style render failed";

//     return content;

//     // const allocator = ctx.allocator;

//     // const body = self.navigator.viewWithContent(ctx.allocator, "No active route") catch "render error";

//     // const cols = zz.flex.layout(allocator, ctx.width, ctx.height, &.{
//     //     .{ .constraint = .fill },
//     //     .{ .constraint = .{ .percentage = 25 } },
//     // }, .{ .direction = .row, .gap = 1 }) catch return "layout error";
// }

// // const AccountScreen = @import("../screens/account_screen.zig");
// // const HomeScreen = @import("../screens/home_screen.zig");
// // const SettingsScreen = @import("../screens/settings_screen.zig");

// // const Screen = union(enum) {
// //     home: *HomeScreen,
// //     account: *AccountScreen,
// //     settings: *SettingsScreen,
// // };

// // allocator: std.mem.Allocator,
// // active: Screen,

// // pub fn init(allocator: std.mem.Allocator) !*Router {
// //     const self = try allocator.create(Router);
// //     self.allocator = allocator;

// //     self.active = .{ .home = try HomeScreen.init(self.allocator) };
// //     return self;
// // }

// // pub fn deinit(self: *Router) void {
// //     switch (self.active) {
// //         .home => |screen| screen.deinit(),
// //         .account => |screen| screen.deinit(),
// //         .settings => |screen| screen.deinit(),
// //     }

// //     self.allocator.destroy(self);
// // }

// // pub fn widget(self: *Router) vxfw.Widget {
// //     return .{
// //         .userdata = self,
// //         .drawFn = typeErasedDrawFn,
// //         .eventHandler = typeErasedEventHandler,
// //     };
// // }

// // fn typeErasedDrawFn(ptr: *anyopaque, ctx: vxfw.DrawContext) std.mem.Allocator.Error!vxfw.Surface {
// //     const self: *Router = @ptrCast(@alignCast(ptr));
// //     return switch (self.active) {
// //         .home => |screen| screen.widget().draw(ctx),
// //         .account => |screen| screen.widget().draw(ctx),
// //         .settings => |screen| screen.widget().draw(ctx),
// //     };
// // }

// // fn typeErasedEventHandler(ptr: *anyopaque, ctx: *vxfw.EventContext, event: vxfw.Event) anyerror!void {
// //     const self: *Router = @ptrCast(@alignCast(ptr));
// //     // Pass events down to the active screen so it can handle its own inputs
// //     try switch (self.active) {
// //         .home => |screen| screen.widget().handleEvent(ctx, event),
// //         .account => |screen| screen.widget().handleEvent(ctx, event),
// //         .settings => |screen| screen.widget().handleEvent(ctx, event),
// //     };
// // }

// const ScreenA = struct {
//     visits: usize = 0,

//     fn onEnter(ctx: *anyopaque) void {
//         const self: *ScreenA = @ptrCast(@alignCast(ctx));
//         self.visits += 1;
//     }

//     fn onKey(_: *anyopaque, _: zz.KeyEvent) bool {
//         return false;
//     }

//     fn render(ctx: *anyopaque, allocator: std.mem.Allocator) ![]const u8 {
//         const self: *ScreenA = @ptrCast(@alignCast(ctx));
//         return std.fmt.allocPrint(
//             allocator,
//             "Home Screen\n\nVisits: {d}\n\nUse Left/Right or 1..9 to switch tabs.",
//             .{self.visits},
//         );
//     }
// };

// const ScreenB = struct {
//     count: i32 = 0,

//     fn onKey(ctx: *anyopaque, key: zz.KeyEvent) bool {
//         const self: *ScreenB = @ptrCast(@alignCast(ctx));
//         switch (key.key) {
//             .char => |c| switch (c) {
//                 '+', '=' => {
//                     self.count += 1;
//                     return true;
//                 },
//                 '-', '_' => {
//                     self.count -= 1;
//                     return true;
//                 },
//                 else => {},
//             },
//             .up => {
//                 self.count += 1;
//                 return true;
//             },
//             .down => {
//                 self.count -= 1;
//                 return true;
//             },
//             else => {},
//         }
//         return false;
//     }

//     fn render(ctx: *anyopaque, allocator: std.mem.Allocator) ![]const u8 {
//         const self: *ScreenB = @ptrCast(@alignCast(ctx));
//         return std.fmt.allocPrint(
//             allocator,
//             "Counter Screen\n\nCount: {d}\n\nPress + / - while this tab is active.",
//             .{self.count},
//         );
//     }
// };
