const std = @import("std");

const vaxis = @import("vaxis");
const vxfw = vaxis.vxfw;

pub const Model = struct {
    title: vxfw.Text,

    // BUG: Hovering over the split view crashes the program
    layout: vxfw.SplitView,
    children: [1]vxfw.SubSurface = undefined,

    allocator: std.mem.Allocator,

    pub fn init(allocator: std.mem.Allocator, title: []const u8) !*Model {
        // TODO: Refactor the widget creation into separate methods
        const left_content = vxfw.Text{ .text = "Left Side" };
        const right_content = vxfw.Text{ .text = "Right Side" };

        const left_widget = vxfw.SizedBox{
            .size = .{ .width = 100, .height = 200 },
            .child = left_content.widget(),
        };

        const right_widget = vxfw.SizedBox{
            .size = .{ .width = 80, .height = 200 },
            .child = right_content.widget(),
        };

        const lhs = vxfw.Border{
            .labels = &.{.{ .text = "Ticker", .alignment = .top_center }},
            .child = left_widget.widget(),
        };

        const rhs = vxfw.Border{
            .labels = &.{.{ .text = "Order Book", .alignment = .top_center }},
            .child = right_widget.widget(),
        };

        const self = try allocator.create(Model);
        self.* = .{
            .allocator = allocator,
            .title = .{ .text = title },
            .layout = .{
                .lhs = lhs.widget(),
                .rhs = rhs.widget(),
                .width = 100,
            },
        };

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
        _ = self;

        switch (event) {
            .init => {
                // TODO: Implement ticker initialization by reading last saved state in cache
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

// pub const App = struct {
//     // 1. Store your widgets in the App struct so their memory addresses remain stable
//     term_view: vxfw.Text, // Using Text as a placeholder for your Terminal widget
//     rhs_view: vxfw.Text,
//     border: Border,
//     split_view: SplitView,

//     pub fn init() App {
//         var app: App = undefined;

//         // Initialize your child widgets [cite: 77, 78]
//         app.term_view = .{ .text = "Terminal running here..." };
//         app.rhs_view = .{ .text = "Side panel (30%)" };

//         // 2. Wrap the terminal view in a border [cite: 30]
//         app.border = .{
//             .child = app.term_view.widget(), [cite: 3, 30]
//             .style = .{}, // Optional: Add custom styling here [cite: 3]
//         };

//         // 3. Initialize the SplitView
//         // We set an initial dummy width; it will be overwritten on the first render/resize
//         app.split_view = .{
//             .lhs = app.border.widget(), [cite: 37, 79]
//             .rhs = app.rhs_view.widget(), [cite: 37, 79]
//             .width = 70, // Target width to draw at [cite: 37, 79]
//             .constrain = .lhs, // We constrain the left side so we control its exact width
//         };

//         return app;
//     }

//     /// Call this inside your main vaxis event loop whenever you receive a `.winsize` event
//     pub fn handleResize(self: *App, window_cols: u16) void {
//         // Calculate 70% of the total available columns
//         const lhs_width: u16 = @intCast((@as(u32, window_cols) * 70) / 100);

//         // Update the SplitView's absolute width dynamically [cite: 37, 57]
//         self.split_view.width = lhs_width;
//     }

//     /// Returns the root widget to be drawn by vxfw
//     pub fn widget(self: *App) vxfw.Widget {
//         return self.split_view.widget(); [cite: 37, 80]
//     }
// };

// // pub const App = struct {
// //     term_text: vxfw.Text,
// //     rhs_text: vxfw.Text,
// //     border: Border,
// //     split_view: SplitView,

// //     // We use this to track window resizes
// //     last_screen_width: ?u16 = null,

// //     pub fn init() App {
// //         var app: App = .{
// //             .term_text = .{ .text = "Terminal View (LHS)" },
// //             .rhs_text = .{ .text = "Side Panel (RHS)" },
// //             .border = undefined,
// //             .split_view = undefined,
// //         };

// //         // Wrap the terminal text in the Border widget
// //         app.border = .{
// //             .child = app.term_text.widget(),
// //             // Optional: specify styling here
// //             .style = .{},
// //         };

// //         // Initialize the SplitView
// //         app.split_view = .{
// //             .lhs = app.border.widget(),
// //             .rhs = app.rhs_text.widget(),
// //             .width = 0, // This gets calculated dynamically on the first draw
// //             .constrain = .lhs,
// //         };

// //         return app;
// //     }

// //     /// Expose the App as a vxfw.Widget
// //     pub fn widget(self: *App) vxfw.Widget {
// //         return .{
// //             .userdata = self,
// //             .eventHandler = typeErasedEventHandler,
// //             .drawFn = typeErasedDrawFn,
// //         };
// //     }

// //     fn typeErasedEventHandler(ptr: *anyopaque, ctx: *vxfw.EventContext, event: vxfw.Event) anyerror!void {
// //         const self: *App = @ptrCast(@alignCast(ptr));
// //         // Route events down to the split_view so dragging the separator works
// //         try self.split_view.widget().handleEvent(ctx, event);
// //     }

// //     fn typeErasedDrawFn(ptr: *anyopaque, ctx: vxfw.DrawContext) std.mem.Allocator.Error!vxfw.Surface {
// //         const self: *App = @ptrCast(@alignCast(ptr));
// //         const current_width = ctx.max.width orelse 100;

// //         // Recalculate the 70% split ONLY if the window has resized.
// //         // This ensures that if the user clicks and drags the separator,
// //         // we don't instantly overwrite their manual adjustments on the next frame.
// //         if (self.last_screen_width == null or self.last_screen_width.? != current_width) {
// //             self.split_view.width = @intCast((@as(u32, current_width) * 70) / 100);
// //             self.last_screen_width = current_width;
// //         }

// //         // Draw the split view with the updated context
// //         return self.split_view.widget().draw(ctx);
// //     }
// // };
