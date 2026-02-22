const std = @import("std");
const Allocator = std.mem.Allocator;

const vaxis = @import("vaxis");
const vxfw = vaxis.vxfw;

pub const SidePanel = struct {
    border: vxfw.Border = undefined,
    labels: [1]vxfw.Border.BorderLabel = undefined,

    title: vxfw.Text = undefined,
    menu: vxfw.ListView = undefined,
    items: [1]vxfw.Widget = undefined,

    fn widget(self: *SidePanel) vxfw.Widget {
        return self.border.widget();
    }
};

pub const MainPanel = struct {
    border: vxfw.Border = undefined,
    labels: [1]vxfw.Border.BorderLabel = undefined,

    content: vxfw.Text = undefined,
    container: vxfw.SizedBox = undefined,

    fn widget(self: *MainPanel) vxfw.Widget {
        return self.border.widget();
    }
};

pub const Layout = struct {
    last_width: ?u16 = 100,
    split: vxfw.SplitView = undefined,
    side_panel: SidePanel = undefined,
    main_panel: MainPanel = undefined,

    fn widget(self: *Layout) vxfw.Widget {
        return self.split.widget();
    }
};

pub const Model = struct {
    layout: Layout = undefined,
    allocator: Allocator = undefined,

    pub fn init(allocator: std.mem.Allocator) !*Model {
        const self = try allocator.create(Model);
        self.allocator = allocator;
        return self;
    }

    pub fn deinit(self: *Model) void {
        self.allocator.destroy(self);
    }

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
                // Side Panel
                self.layout.side_panel.title = .{ .text = "Portman" };
                self.layout.side_panel.items = .{
                    self.layout.side_panel.title.widget(),
                };

                self.layout.side_panel.menu = .{ .children = .{
                    .slice = &self.layout.side_panel.items,
                } };

                self.layout.side_panel.labels = .{.{
                    .text = " Side Panel ",
                    .alignment = .top_center,
                }};

                self.layout.side_panel.border = .{
                    .child = self.layout.side_panel.menu.widget(),
                    .labels = &self.layout.side_panel.labels,
                };

                // Main Panel
                self.layout.main_panel.content = .{ .text = 
                    \\ Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididuntut labore et
                    \\ dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip
                    \\ ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore
                    \\ eu fugiat nulla pariatur.
                    \\
                    \\ Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est
                    \\ laborum. Curabitur pretium tincidunt lacus. Nulla gravida orci a odio. Nullam varius, turpis et commodo
                    \\ pharetra, est eros bibendum elit, nec luctus magna felis sollicitudin mauris. Integer in mauris eu nibh
                    \\ euismod gravida.
                };

                self.layout.main_panel.labels = .{
                    .{ .text = " Main Panel ", .alignment = .top_center },
                };

                self.layout.main_panel.container = .{
                    .child = self.layout.main_panel.content.widget(),
                    .size = .{ .width = 200, .height = 100 },
                };

                self.layout.main_panel.border = .{
                    .child = self.layout.main_panel.container.widget(),
                    .labels = &self.layout.main_panel.labels,
                };

                // Layout
                self.layout.split = .{
                    .lhs = self.layout.side_panel.widget(),
                    .rhs = self.layout.main_panel.widget(),
                    .width = self.layout.last_width orelse 0,
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
        const current_width = ctx.max.width orelse 100;

        if (self.layout.last_width == null or self.layout.last_width.? != current_width) {
            self.layout.split.width = @intCast((@as(u32, current_width) * 20) / 100);
            self.layout.last_width = current_width;
        }

        return self.layout.widget().draw(ctx);
    }
};
