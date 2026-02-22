const std = @import("std");

const vaxis = @import("vaxis");
const vxfw = vaxis.vxfw;

pub const Chart = @import("components/chart.zig");
pub const Router = @import("router.zig");

const Model = @This();

allocator: std.mem.Allocator,
side_panel_text: vxfw.Text,
router: Router,
split: vxfw.SplitView,

pub fn init(allocator: std.mem.Allocator) !*Model {
    const self = try allocator.create(Model);

    // Start the app on the Home screen
    self.router = .{ .active = .{ .home = .{} } };

    // Placeholder for the side panel navigation
    self.side_panel_text = .{ .text = "Navigation:\n[1] Home\n[2] Account\n[3] Config" };

    self.split = .{
        .lhs = self.side_panel_text.widget(),
        .rhs = self.router.widget(), // The Router goes here!
        .width = 20,
        .constrain = .lhs,
    };

    self.allocator = allocator;
    return self;
}

pub fn deinit(self: *Model) void {
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
    return self.split.widget().draw(ctx);
}

fn typeErasedEventHandler(ptr: *anyopaque, ctx: *vxfw.EventContext, event: vxfw.Event) anyerror!void {
    const self: *Model = @ptrCast(@alignCast(ptr));

    // Handle global navigation hotkeys
    if (event == .key_press) {
        const key = event.key_press;
        if (key.matches('c', .{ .ctrl = true })) {
            ctx.quit = true;
            return;
        } else if (key.matches('1', .{})) {
            self.router.active = .{ .home = .{} };
            ctx.redraw = true;
        } else if (key.matches('2', .{})) {
            self.router.active = .{ .account = .{} };
            ctx.redraw = true;
        } else if (key.matches('3', .{})) {
            self.router.active = .{ .config = .{} };
            ctx.redraw = true;
        }
    }

    // Ensure the split view and active screen get mouse/key events
    try self.split.widget().handleEvent(ctx, event);
}
