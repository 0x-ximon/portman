const std = @import("std");

const lib = @import("lib");
const zz = @import("zigzag");

const panels = union(enum) {
    dashboard,
    account,
    settings,
    scripts,
    help,
};

const Model = @This();

count: i32,
active: panels,

pub const Msg = union(enum) {
    key: zz.KeyEvent,
    click: zz.MouseEvent,
};

pub fn init(self: *Model, _: *zz.Context) zz.Cmd(Msg) {
    self.* = .{ .count = 0, .active = .dashboard };
    return .{ .set_title = lib.APP_NAME };
}

pub fn view(_: *const Model, ctx: *const zz.Context) []const u8 {
    const allocator = ctx.allocator;

    const rows = zz.flex.layout(allocator, ctx.width, ctx.height, &.{
        .{ .constraint = .{ .percentage = 10 } },
        .{ .constraint = .fill },
        .{ .constraint = .{ .percentage = 10 } },
    }, .{ .direction = .column }) catch return "layout error";

    const cols = zz.flex.layout(allocator, rows[1].width, rows[1].height, &.{
        .{ .constraint = .fill },
        .{ .constraint = .{ .percentage = 25 } },
    }, .{ .direction = .row, .gap = 1 }) catch return "layout error";

    const header = render(allocator, "Header", rows[0].width, rows[0].height);
    const footer = render(allocator, "Footer", rows[2].width, rows[2].height);

    const left = render(allocator, "Left", cols[0].width, cols[0].height);
    const right = render(allocator, "Right", cols[1].width, cols[1].height);
    const main = zz.join.horizontal(allocator, .middle, &.{ left, right }) catch "main render error";

    const layout = zz.join.vertical(allocator, .center, &.{ header, main, footer }) catch "render error";
    return layout;
}

pub fn render(allocator: std.mem.Allocator, content: []const u8, w: u16, h: u16) []const u8 {
    var s = zz.Style{};
    s = s.borderAll(zz.Border.rounded);

    // Account for border (2 cells each side)
    const inner_w: u16 = if (w > 4) w - 4 else 1;
    const inner_h: u16 = if (h > 2) h - 2 else 1;
    s = s.width(inner_w);
    s = s.height(inner_h);

    return s.render(allocator, content) catch content;
}

pub fn update(self: *Model, msg: Msg, _: *zz.Context) zz.Cmd(Msg) {
    switch (msg) {
        .key => |k| switch (k.key) {
            .char => |c| switch (c) {
                '1' => self.active = .dashboard,
                '2' => self.active = .account,
                '3' => self.active = .settings,
                '4' => self.active = .scripts,
                '5' => self.active = .help,
                else => {},
            },
            // .tab => {
            //     self.active = (self.active + 1) % 5;
            // },
            .escape => return .quit,
            else => {},
        },

        .click => {},
    }

    return .none;
}

// allocator: std.mem.Allocator,
// split: vxfw.SplitView,
// width: ?u16 = null,

// navigator: *lib.Navigator,
// router: *lib.Router,

// pub fn init(allocator: std.mem.Allocator) !*Model {
//     const self = try allocator.create(Model);
//     self.allocator = allocator;

//     self.navigator = try lib.Navigator.init(self.allocator);
//     self.router = try lib.Router.init(self.allocator);

//     self.split = .{
//         .lhs = self.navigator.widget(),
//         .rhs = self.router.widget(),
//         .style = .{ .invisible = true },
//         .width = 100,
//     };

//     return self;
// }

// pub fn deinit(self: *Model) void {
//     self.navigator.deinit();
//     self.router.deinit();

//     self.allocator.destroy(self);
// }

// pub fn widget(self: *Model) vxfw.Widget {
//     return .{
//         .userdata = self,
//         .drawFn = typeErasedDrawFn,
//         .eventHandler = typeErasedEventHandler,
//     };
// }

// fn typeErasedDrawFn(ptr: *anyopaque, ctx: vxfw.DrawContext) std.mem.Allocator.Error!vxfw.Surface {
//     const self: *Model = @ptrCast(@alignCast(ptr));
//     const current_width = ctx.max.width orelse 100;

//     if (self.width == null or self.width.? != current_width) {
//         self.split.width = @intCast((@as(u32, current_width) * 15) / 100);
//         self.width = current_width;
//     }

//     return self.split.widget().draw(ctx);
// }

// fn typeErasedEventHandler(ptr: *anyopaque, ctx: *vxfw.EventContext, event: vxfw.Event) anyerror!void {
//     const self: *Model = @ptrCast(@alignCast(ptr));
//     switch (event) {
//         .init => {},
//         .key_press => |key| {
//             if (key.matches('c', .{ .ctrl = true })) {
//                 ctx.quit = true;
//                 return;
//             }

//             if (key.matches('1', .{})) {
//                 self.router.active = .{ .home = try .init(self.allocator) };
//                 ctx.redraw = true;
//             } else if (key.matches('2', .{})) {
//                 self.router.active = .{ .account = try .init(self.allocator) };
//                 ctx.redraw = true;
//             } else if (key.matches('3', .{})) {
//                 self.router.active = .{ .settings = try .init(self.allocator) };
//                 ctx.redraw = true;
//             }
//         },

//         else => {},
//     }

//     // Ensure the split view and active screen get mouse/key events
//     try self.split.widget().handleEvent(ctx, event);
// }
