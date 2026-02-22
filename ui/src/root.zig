const std = @import("std");

const vaxis = @import("vaxis");
const vxfw = vaxis.vxfw;

// pub const SidePanel = struct {
//     items: [1]vxfw.Widget = undefined,
//     container: vxfw.ListView = undefined,
//     border: vxfw.Border = undefined,

//     title: vxfw.Text = undefined,

//     fn widget(self: *SidePanel) vxfw.Widget {
//         return self.border.widget();
//     }
// };

// pub const MainPanel = struct {
//     items: [2]vxfw.FlexItem = undefined,
//     container: vxfw.FlexColumn = undefined,
//     labels: [1]vxfw.Border.BorderLabel = undefined,

//     chart: Chart = undefined,
//     chart_item: vxfw.FlexItem = undefined,
//     chart_border: vxfw.Border = undefined,

//     content: vxfw.Text = undefined,
//     content_item: vxfw.FlexItem = undefined,
//     content_border: vxfw.Border = undefined,

//     fn widget(self: *MainPanel) vxfw.Widget {
//         return self.border.widget();
//     }
// };

// pub const Layout = struct {
//     last_width: ?u16 = 100,
//     split: vxfw.SplitView = undefined,
//     // side_panel: SidePanel = undefined,
//     // main_panel: MainPanel = undefined,

//     fn widget(self: *Layout) vxfw.Widget {
//         return self.split.widget();
//     }
// };

// pub const Model = struct {
//     layout: Layout = undefined,
//     router: Router = .Home,
//     allocator: std.mem.Allocator = undefined,

//     pub fn init(allocator: std.mem.Allocator) !*Model {
//         const self = try allocator.create(Model);
//         self.allocator = allocator;
//         return self;
//     }

//     pub fn deinit(self: *Model) void {
//         self.allocator.destroy(self);
//     }

//     pub fn widget(self: *Model) vxfw.Widget {
//         return .{
//             .userdata = self,
//             .drawFn = Model.typeErasedDrawFn,
//             // .eventHandler = Model.typeErasedEventHandler,
//         };
//     }

// fn typeErasedEventHandler(ptr: *anyopaque, ctx: *vxfw.EventContext, event: vxfw.Event) anyerror!void {
//     const self: *Model = @ptrCast(@alignCast(ptr));
//     switch (event) {
//         .init => {
//             // Side Panel
//             self.layout.side_panel.title = .{ .text = "Portman" };
//             self.layout.side_panel.items = .{
//                 self.layout.side_panel.title.widget(),
//             };

//             self.layout.side_panel.container = .{ .children = .{
//                 .slice = &self.layout.side_panel.items,
//             } };

//             self.layout.side_panel.border = .{
//                 .child = self.layout.side_panel.container.widget(),
//             };

//             // Main Panel
//             self.layout.main_panel.content = .{ .text =
//                 \\ Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididuntut labore et
//                 \\ dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip
//                 \\ ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore
//                 \\ eu fugiat nulla pariatur.
//                 \\
//                 \\ Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est
//                 \\ laborum. Curabitur pretium tincidunt lacus. Nulla gravida orci a odio. Nullam varius, turpis et commodo
//                 \\ pharetra, est eros bibendum elit, nec luctus magna felis sollicitudin mauris. Integer in mauris eu nibh
//                 \\ euismod gravida.
//             };
//             self.layout.main_panel.content_item = .{
//                 .{ .widget = self.layout.main_panel.content.widget() },
//             };

//             // TODO: Get Tick Data from Exchange API
//             self.layout.main_panel.chart = .{};
//             self.layout.main_panel.chart_item = .{
//                 .widget = self.layout.main_panel.chart.widget(),
//             };

//             self.layout.main_panel.container = .{
//                 .children = .{
//                     self.layout.main_panel.chart_item,
//                     self.layout.main_panel.content_item,
//                 },
//             };

//             self.layout.main_panel.labels = .{
//                 .{ .text = " Market View ", .alignment = .top_center },
//             };

//             self.layout.main_panel.border = .{
//                 .child = self.layout.main_panel.container.widget(),
//                 .labels = &self.layout.main_panel.labels,
//             };

//             // Layout
//             self.layout.split = .{
//                 .lhs = self.layout.side_panel.widget(),
//                 .rhs = self.layout.main_panel.widget(),
//                 .width = self.layout.last_width orelse 0,
//                 .style = .{ .invisible = true },
//             };
//         },

//         .key_press => |key| {
//             if (key.matches('c', .{ .ctrl = true })) {
//                 ctx.quit = true;
//                 return;
//             }
//         },

//         else => {},
//     }
// }

//     fn typeErasedDrawFn(ptr: *anyopaque, ctx: vxfw.DrawContext) std.mem.Allocator.Error!vxfw.Surface {
//         const self: *Model = @ptrCast(@alignCast(ptr));
//         const current_width = ctx.max.width orelse 100;

//         if (self.layout.last_width == null or self.layout.last_width.? != current_width) {
//             self.layout.split.width = @intCast((@as(u32, current_width) * 20) / 100);
//             self.layout.last_width = current_width;
//         }

//         return self.layout.widget().draw(ctx);
//     }
// };
