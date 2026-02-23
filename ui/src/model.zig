const std = @import("std");
const lib = @import("lib");

const vaxis = @import("vaxis");
const vxfw = vaxis.vxfw;

const Model = @This();

allocator: std.mem.Allocator,
split: vxfw.SplitView,
width: ?u16 = null,

navigator: *lib.Navigator,
router: *lib.Router,

pub fn init(allocator: std.mem.Allocator) !*Model {
    const self = try allocator.create(Model);
    self.allocator = allocator;

    self.navigator = try lib.Navigator.init(self.allocator);
    self.router = try lib.Router.init(self.allocator);

    return self;
}

pub fn deinit(self: *Model) void {
    self.navigator.deinit();
    self.router.deinit();

    self.allocator.destroy(self);
}

pub fn widget(self: *Model) vxfw.Widget {
    return .{
        .userdata = self,
        .drawFn = typeErasedDrawFn,
        .eventHandler = typeErasedEventHandler,
    };
}

fn typeErasedDrawFn(ptr: *anyopaque, ctx: vxfw.DrawContext) std.mem.Allocator.Error!vxfw.Surface {
    const self: *Model = @ptrCast(@alignCast(ptr));
    const current_width = ctx.max.width orelse 100;

    if (self.width == null or self.width.? != current_width) {
        self.split.width = @intCast((@as(u32, current_width) * 15) / 100);
        self.width = current_width;
    }

    return self.split.widget().draw(ctx);
}

fn typeErasedEventHandler(ptr: *anyopaque, ctx: *vxfw.EventContext, event: vxfw.Event) anyerror!void {
    const self: *Model = @ptrCast(@alignCast(ptr));
    switch (event) {
        .init => {
            self.split = .{
                .lhs = self.navigator.widget(),
                .rhs = self.router.widget(),
                .style = .{ .invisible = true },
                .width = 100,
            };
        },

        .key_press => |key| {
            if (key.matches('c', .{ .ctrl = true })) {
                ctx.quit = true;
                return;
            }

            if (key.matches('1', .{})) {
                self.router.active = .{ .home = try .init(self.allocator) };
                ctx.redraw = true;
            } else if (key.matches('2', .{})) {
                self.router.active = .{ .account = try .init(self.allocator) };
                ctx.redraw = true;
            } else if (key.matches('3', .{})) {
                self.router.active = .{ .settings = try .init(self.allocator) };
                ctx.redraw = true;
            }
        },

        else => {},
    }

    // Ensure the split view and active screen get mouse/key events
    try self.split.widget().handleEvent(ctx, event);
}
