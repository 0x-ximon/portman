const std = @import("std");

const vaxis = @import("vaxis");
const vxfw = vaxis.vxfw;

const Layered = @This();

below: vxfw.Widget,
above: vxfw.Widget,

pub fn widget(self: *Layered) vxfw.Widget {
    return .{
        .userdata = self,
        .drawFn = Layered.typeErasedDrawFn,
        .eventHandler = Layered.typeErasedEventHandler,
    };
}

fn typeErasedDrawFn(ptr: *anyopaque, ctx: vxfw.DrawContext) std.mem.Allocator.Error!vxfw.Surface {
    const self: *Layered = @ptrCast(@alignCast(ptr));

    const width = ctx.max.width orelse ctx.min.width;
    const height = ctx.max.height orelse @max(ctx.min.height, 10);
    const actual_size = vxfw.Size{ .width = width, .height = height };

    const children = try ctx.arena.alloc(vxfw.SubSurface, 2);

    children[0] = .{
        .z_index = 0,
        .origin = .{ .row = 0, .col = 0 },
        .surface = try self.below.draw(ctx),
    };

    children[1] = .{
        .z_index = 1,
        .origin = .{ .row = 0, .col = 0 },
        .surface = try self.above.draw(ctx),
    };

    return .{
        .buffer = &.{},
        .size = actual_size,
        .children = children,
        .widget = self.widget(),
    };
}

fn typeErasedEventHandler(ptr: *anyopaque, ctx: *vxfw.EventContext, event: vxfw.Event) anyerror!void {
    const self: *Layered = @ptrCast(@alignCast(ptr));
    try self.above.handleEvent(ctx, event);
    try self.below.handleEvent(ctx, event);
}
