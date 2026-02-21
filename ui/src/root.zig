const std = @import("std");
const vaxis = @import("vaxis");
const vxfw = vaxis.vxfw;

const Text = vxfw.Text;
// const ListView = vxfw.ListView;
// const Widget = vxfw.Widget;

// const Cell = vaxis.Cell;
// const TextInput = vaxis.widgets.TextInput;
// const border = vaxis.widgets.border;

pub const Model = struct {
    layout: vxfw.SplitView = undefined,
    children: [1]vxfw.SubSurface = undefined,

    pub fn widget(self: *Model) vxfw.Widget {
        return .{
            .userdata = self,
            .drawFn = Model.typeErasedDrawFn,
            .eventHandler = Model.typeErasedEventHandler,
        };
    }

    fn typeErasedEventHandler(ptr: *anyopaque, ctx: *vxfw.EventContext, event: vxfw.Event) anyerror!void {
        const self: *Model = @ptrCast(@alignCast(ptr));
        switch (event) {
            .init => {
                // TODO: Implement proper initialization by reading last saved state in cache
                const lhs: Text = .{ .text = "Ticker" };
                const rhs: Text = .{ .text = "Order Book" };

                self.layout.lhs = lhs.widget();
                self.layout = .{
                    .lhs = lhs.widget(),
                    .rhs = rhs.widget(),
                    .width = 100,
                };
            },

            .key_press => |key| {
                if (key.matches('c', .{ .ctrl = true })) {
                    ctx.quit = true;
                    return;
                }
            },

            else => {},
        }
    }

    fn typeErasedDrawFn(ptr: *anyopaque, ctx: vxfw.DrawContext) std.mem.Allocator.Error!vxfw.Surface {
        const self: *Model = @ptrCast(@alignCast(ptr));
        const surf = try self.layout.widget().draw(ctx);

        self.children[0] = .{
            .surface = surf,
            .origin = .{ .row = 0, .col = 0 },
        };

        return .{
            .size = ctx.max.size(),
            .widget = self.widget(),
            .buffer = &.{},
            .children = &self.children,
        };
    }
};
