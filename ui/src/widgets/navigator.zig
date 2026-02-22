const std = @import("std");

const vaxis = @import("vaxis");
const vxfw = vaxis.vxfw;

pub const Navigator = @This();

allocator: std.mem.Allocator,
container: vxfw.Widget,

content: [4]vxfw.Widget,
padding: vxfw.Padding,
view: vxfw.ListView,
border: vxfw.Border,

// TODO: Replace with ghost buttons
title: vxfw.Text,
home: vxfw.Text,
account: vxfw.Text,
config: vxfw.Text,

pub fn init(allocator: std.mem.Allocator) !*Navigator {
    const self = try allocator.create(Navigator);
    self.allocator = allocator;

    self.title = .{ .text = "Portman" };
    self.home = .{ .text = "Home" };
    self.account = .{ .text = "Account" };
    self.config = .{ .text = "Config" };

    self.content = .{
        self.title.widget(),
        self.home.widget(),
        self.account.widget(),
        self.config.widget(),
    };

    self.view = .{
        .children = .{ .slice = &self.content },
    };

    self.padding = .{ .child = self.view.widget(), .padding = .{
        .left = 1,
        .right = 1,
        .top = 1,
        .bottom = 1,
    } };

    self.border = .{
        .child = self.padding.widget(),
    };

    self.container = self.border.widget();
    return self;
}

pub fn deinit(self: *Navigator) void {
    self.allocator.destroy(self);
}

pub fn widget(self: *Navigator) vxfw.Widget {
    return .{
        .userdata = self,
        .drawFn = typeErasedDrawFn,
        .eventHandler = typeErasedEventHandler,
    };
}

fn typeErasedDrawFn(ptr: *anyopaque, ctx: vxfw.DrawContext) std.mem.Allocator.Error!vxfw.Surface {
    const self: *Navigator = @ptrCast(@alignCast(ptr));
    return self.container.draw(ctx);
}

fn typeErasedEventHandler(ptr: *anyopaque, ctx: *vxfw.EventContext, event: vxfw.Event) anyerror!void {
    const self: *Navigator = @ptrCast(@alignCast(ptr));
    _ = self; // autofix
    _ = ctx; // autofix

    switch (event) {
        else => {},
    }
}
