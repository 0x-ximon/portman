const std = @import("std");

const vaxis = @import("vaxis");
const vxfw = vaxis.vxfw;

const Chart = @import("../widgets/chart.zig");

const MainPanel = struct {
    allocator: std.mem.Allocator,
    container: vxfw.Widget,

    chart: *Chart,

    indicators_content: vxfw.Text,
    indicators_border: vxfw.Border,
    indicators_widget: vxfw.Widget,

    pub fn init(allocator: std.mem.Allocator) !*MainPanel {
        const self = try allocator.create(MainPanel);
        self.allocator = allocator;

        self.chart = try Chart.init(allocator);

        // Indicators Initialization
        self.indicators_content = .{ .text = 
            \\ Lorem ipsum dolor sit amet, consectetur adipiscing elit.
        };

        self.indicators_border = .{
            .labels = &.{.{ .text = "Technical Indicators", .alignment = .top_left }},
            .child = self.indicators_content.widget(),
        };

        self.indicators_widget = self.indicators_border.widget();

        return self;
    }

    pub fn deinit(self: *MainPanel) void {
        self.chart.deinit();

        self.allocator.destroy(self);
    }

    pub fn widget(self: *MainPanel) vxfw.Widget {
        return .{
            .userdata = self,
            .drawFn = MainPanel.typeErasedDrawFn,
            .eventHandler = MainPanel.typeErasedEventHandler,
        };
    }

    fn typeErasedDrawFn(ptr: *anyopaque, ctx: vxfw.DrawContext) std.mem.Allocator.Error!vxfw.Surface {
        const self: *MainPanel = @ptrCast(@alignCast(ptr));
        const max_size = ctx.max.size();

        const column: vxfw.FlexColumn = .{
            .children = &.{
                .init(self.chart.widget(), 2),
                .init(self.indicators_widget, 1),
            },
        };

        const children = try ctx.arena.alloc(vxfw.SubSurface, 1);
        children[0] = .{
            .origin = .{ .row = 0, .col = 0 },
            .surface = try column.draw(ctx.withConstraints(
                .{ .width = max_size.width, .height = 0 },
                ctx.max,
            )),
        };

        return .{
            .buffer = &.{},
            .size = max_size,
            .children = children,
            .widget = self.widget(),
        };
    }

    fn typeErasedEventHandler(ptr: *anyopaque, ctx: *vxfw.EventContext, event: vxfw.Event) anyerror!void {
        const self: *MainPanel = @ptrCast(@alignCast(ptr));
        try self.chart.widget().handleEvent(ctx, event);
    }
};

const SidePanel = struct {
    allocator: std.mem.Allocator,
    container: vxfw.Widget,

    book_content: vxfw.Text,
    book_border: vxfw.Border,
    book_widget: vxfw.Widget,

    watchlist_content: vxfw.Text,
    watchlist_border: vxfw.Border,
    watchlist_widget: vxfw.Widget,

    pub fn init(allocator: std.mem.Allocator) !*SidePanel {
        const self = try allocator.create(SidePanel);
        self.allocator = allocator;

        // Order Book Initialization
        self.book_content = .{ .text = 
            \\ Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et
            \\ dolore magna aliqua.
        };

        self.book_border = .{
            .labels = &.{.{ .text = "Order Book", .alignment = .top_left }},
            .child = self.book_content.widget(),
        };

        self.book_widget = self.book_border.widget();

        // Watchlist Initialization
        self.watchlist_content = .{ .text = 
            \\ Lorem ipsum dolor sit amet, consectetur adipiscing elit.
        };

        self.watchlist_border = .{
            .labels = &.{.{ .text = "Watchlist", .alignment = .top_left }},
            .child = self.watchlist_content.widget(),
        };

        self.watchlist_widget = self.watchlist_border.widget();

        return self;
    }

    pub fn deinit(self: *SidePanel) void {
        self.allocator.destroy(self);
    }

    pub fn widget(self: *SidePanel) vxfw.Widget {
        return .{
            .userdata = self,
            .drawFn = SidePanel.typeErasedDrawFn,
            .eventHandler = SidePanel.typeErasedEventHandler,
        };
    }

    fn typeErasedDrawFn(ptr: *anyopaque, ctx: vxfw.DrawContext) std.mem.Allocator.Error!vxfw.Surface {
        const self: *SidePanel = @ptrCast(@alignCast(ptr));
        const max_size = ctx.max.size();

        const column: vxfw.FlexColumn = .{
            .children = &.{
                .init(self.book_widget, 2),
                .init(self.watchlist_widget, 1),
            },
        };

        const children = try ctx.arena.alloc(vxfw.SubSurface, 1);
        children[0] = .{
            .origin = .{ .row = 0, .col = 0 },
            .surface = try column.draw(ctx.withConstraints(
                .{ .width = max_size.width, .height = 0 },
                ctx.max,
            )),
        };

        return .{
            .buffer = &.{},
            .size = max_size,
            .children = children,
            .widget = self.widget(),
        };
    }

    fn typeErasedEventHandler(ptr: *anyopaque, ctx: *vxfw.EventContext, event: vxfw.Event) anyerror!void {
        _ = ctx; // autofix
        const self: *SidePanel = @ptrCast(@alignCast(ptr));
        _ = self; // autofix
        switch (event) {
            else => {},
        }
    }
};

const HomeScreen = @This();

allocator: std.mem.Allocator,
main_panel: *MainPanel,
side_panel: *SidePanel,

pub fn init(allocator: std.mem.Allocator) !*HomeScreen {
    const self = try allocator.create(HomeScreen);
    self.allocator = allocator;

    self.main_panel = try MainPanel.init(allocator);
    self.side_panel = try SidePanel.init(allocator);

    return self;
}

pub fn deinit(self: *HomeScreen) void {
    self.main_panel.deinit();
    self.side_panel.deinit();

    self.allocator.destroy(self);
}

pub fn widget(self: *HomeScreen) vxfw.Widget {
    return .{
        .userdata = self,
        .drawFn = HomeScreen.typeErasedDrawFn,
        .eventHandler = HomeScreen.typeErasedEventHandler,
    };
}

fn typeErasedDrawFn(ptr: *anyopaque, ctx: vxfw.DrawContext) std.mem.Allocator.Error!vxfw.Surface {
    const self: *HomeScreen = @ptrCast(@alignCast(ptr));
    const max_size = ctx.max.size();

    const row: vxfw.FlexRow = .{
        .children = &.{
            .init(self.main_panel.widget(), 3),
            .init(self.side_panel.widget(), 1),
        },
    };

    const children = try ctx.arena.alloc(vxfw.SubSurface, 1);
    children[0] = .{
        .origin = .{ .row = 0, .col = 0 },
        .surface = try row.draw(ctx.withConstraints(
            .{ .width = 0, .height = max_size.height },
            ctx.max,
        )),
    };

    return .{
        .buffer = &.{},
        .size = max_size,
        .children = children,
        .widget = self.widget(),
    };
}

fn typeErasedEventHandler(ptr: *anyopaque, ctx: *vxfw.EventContext, event: vxfw.Event) anyerror!void {
    const self: *HomeScreen = @ptrCast(@alignCast(ptr));
    try self.main_panel.widget().handleEvent(ctx, event);
    try self.side_panel.widget().handleEvent(ctx, event);
}
