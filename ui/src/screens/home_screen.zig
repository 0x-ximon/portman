const std = @import("std");

const vaxis = @import("vaxis");
const vxfw = vaxis.vxfw;

const MainPanel = struct {
    allocator: std.mem.Allocator,
    container: vxfw.Widget,

    content: [2]vxfw.FlexItem,
    view: vxfw.FlexColumn,

    chart_content: vxfw.Text,
    // chart_padding: vxfw.Padding,
    chart_border: vxfw.Border,
    chart_item: vxfw.FlexItem,

    indicators_content: vxfw.Text,
    // indicators_padding: vxfw.Padding,
    indicators_border: vxfw.Border,
    indicators_item: vxfw.FlexItem,

    pub fn init(allocator: std.mem.Allocator) !*MainPanel {
        const self = try allocator.create(MainPanel);
        self.allocator = allocator;

        // Chart Initialization
        self.chart_content = .{ .text = 
            \\ Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididuntut labore et
            \\ dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip
            \\ ex ea commodo consequat.
        };

        // self.chart_padding = .{
        //     .child = self.chart_content.widget(),
        //     .padding = .{ .left = 1, .right = 1, .top = 1, .bottom = 1 },
        // };

        self.chart_border = .{
            .child = self.chart_content.widget(),
        };

        self.chart_item = .{
            .widget = self.chart_border.widget(),
            .flex = 3,
        };

        // Indicators Initialization
        self.indicators_content = .{ .text = 
            \\ Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididuntut labore et
            \\ dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip
            \\ ex ea commodo consequat.
        };

        // self.indicators_padding = .{
        //     .child = self.indicators_content.widget(),
        //     .padding = .{ .left = 1, .right = 1, .top = 1, .bottom = 1 },
        // };

        self.indicators_border = .{
            .child = self.indicators_content.widget(),
        };

        self.indicators_item = .{
            .widget = self.indicators_border.widget(),
            .flex = 1,
        };

        self.content = .{ self.chart_item, self.indicators_item };
        self.view = .{ .children = &self.content };

        self.container = self.view.widget();
        return self;
    }

    pub fn deinit(self: *MainPanel) void {
        self.allocator.destroy(self);
    }

    pub fn widget(self: *MainPanel) vxfw.Widget {
        return self.container;
    }
};

const SidePanel = struct {
    allocator: std.mem.Allocator,
    container: vxfw.Widget,

    content: [2]vxfw.FlexItem,
    view: vxfw.FlexColumn,

    book_content: vxfw.Text,
    // book_padding: vxfw.Padding,
    book_border: vxfw.Border,
    book_item: vxfw.FlexItem,

    watchlist_content: vxfw.Text,
    // watchlist_padding: vxfw.Padding,
    watchlist_border: vxfw.Border,
    watchlist_item: vxfw.FlexItem,

    pub fn init(allocator: std.mem.Allocator) !*SidePanel {
        const self = try allocator.create(SidePanel);
        self.allocator = allocator;

        // Order Book Initialization
        self.book_content = .{ .text = 
            \\ Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididuntut labore et
            \\ dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip
            \\ ex ea commodo consequat.
        };

        // self.book_padding = .{
        //     .child = self.book_content.widget(),
        //     .padding = .{ .left = 1, .right = 1, .top = 1, .bottom = 1 },
        // };

        self.book_border = .{
            .child = self.book_content.widget(),
        };

        self.book_item = .{
            .widget = self.book_border.widget(),
            .flex = 3,
        };

        // Watchlist Initialization
        self.watchlist_content = .{ .text = 
            \\ Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididuntut labore et
            \\ dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip
            \\ ex ea commodo consequat.
        };

        // self.watchlist_padding = .{
        //     .child = self.watchlist_content.widget(),
        //     .padding = .{ .left = 1, .right = 1, .top = 1, .bottom = 1 },
        // };

        self.watchlist_border = .{
            .child = self.watchlist_content.widget(),
        };

        self.watchlist_item = .{
            .widget = self.watchlist_border.widget(),
            .flex = 1,
        };

        self.content = .{ self.book_item, self.watchlist_item };
        self.view = .{ .children = &self.content };

        self.container = self.view.widget();
        return self;
    }

    pub fn deinit(self: *SidePanel) void {
        self.allocator.destroy(self);
    }

    pub fn widget(self: *SidePanel) vxfw.Widget {
        return self.container;
    }
};

const HomeScreen = @This();

allocator: std.mem.Allocator,
container: vxfw.Widget,

// TODO: Implement Proper Responsive Layout with Flexible Widgets
// Currently Blocked by https://github.com/rockorager/libvaxis/issues/229

split: vxfw.SplitView,
width: ?u16 = null,

main_panel: *MainPanel,
side_panel: *SidePanel,

pub fn init(allocator: std.mem.Allocator) !*HomeScreen {
    const self = try allocator.create(HomeScreen);
    self.allocator = allocator;

    self.main_panel = try MainPanel.init(allocator);
    self.side_panel = try SidePanel.init(allocator);

    self.split = .{
        .lhs = self.main_panel.widget(),
        .rhs = self.side_panel.widget(),
        .style = .{ .invisible = true },
        .width = 100,
    };

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
        .drawFn = typeErasedDrawFn,
        .eventHandler = typeErasedEventHandler,
    };
}

fn typeErasedDrawFn(ptr: *anyopaque, ctx: vxfw.DrawContext) std.mem.Allocator.Error!vxfw.Surface {
    const self: *HomeScreen = @ptrCast(@alignCast(ptr));
    const current_width = ctx.max.width orelse 200;

    if (self.width == null or self.width.? != current_width) {
        self.split.width = @intCast((@as(u32, current_width) * 80) / 100);
        self.width = current_width;
    }

    return self.split.widget().draw(ctx);
}

fn typeErasedEventHandler(ptr: *anyopaque, ctx: *vxfw.EventContext, event: vxfw.Event) anyerror!void {
    const self: *HomeScreen = @ptrCast(@alignCast(ptr));
    switch (event) {
        else => {},
    }

    try self.split.widget().handleEvent(ctx, event);
}
