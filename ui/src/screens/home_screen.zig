const std = @import("std");

const vaxis = @import("vaxis");
const vxfw = vaxis.vxfw;

const TickerSection = struct {
    allocator: std.mem.Allocator,
    container: vxfw.Widget,

    content: [2]vxfw.FlexItem,
    view: vxfw.FlexRow,

    chart_content: vxfw.Text,
    chart_padding: vxfw.Padding,
    chart_border: vxfw.Border,
    chart_item: vxfw.FlexItem,

    book_content: vxfw.Text,
    book_padding: vxfw.Padding,
    book_border: vxfw.Border,
    book_item: vxfw.FlexItem,

    pub fn init(allocator: std.mem.Allocator) !*TickerSection {
        const self = try allocator.create(TickerSection);
        self.allocator = allocator;

        // Chart Initialization
        self.chart_content = .{ .text = 
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

        self.chart_padding = .{ .child = self.chart_content.widget(), .padding = .{
            .left = 1,
            .right = 1,
            .top = 1,
            .bottom = 1,
        } };

        self.chart_border = .{
            .child = self.chart_padding.widget(),
        };

        self.chart_item = .{
            .widget = self.chart_border.widget(),
            .flex = 2,
        };

        // Book Initialization
        self.book_content = .{ .text = 
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

        self.book_padding = .{ .child = self.book_content.widget(), .padding = .{
            .left = 1,
            .right = 1,
            .top = 1,
            .bottom = 1,
        } };

        self.book_border = .{
            .child = self.book_padding.widget(),
        };

        self.book_item = .{
            .widget = self.book_border.widget(),
        };

        self.content = .{ self.chart_item, self.book_item };
        self.view = .{ .children = &self.content };

        self.container = self.view.widget();
        return self;
    }

    pub fn deinit(self: *TickerSection) void {
        self.allocator.destroy(self);
    }

    pub fn widget(self: *TickerSection) vxfw.Widget {
        return self.container;
    }

    pub fn flex(self: *TickerSection) vxfw.FlexItem {
        return .{ .widget = self.container };
    }
};

// const MarketSection = struct {
//     allocator: std.mem.Allocator,
//     container: vxfw.Widget,

//     content: [2]vxfw.Widget,
//     padding: vxfw.Padding,
//     view: vxfw.ListView,
//     border: vxfw.Border,

//     pub fn init(allocator: std.mem.Allocator) !*MarketSection {
//         const self = try allocator.create(MarketSection);
//         self.allocator = allocator;

//         return self;
//     }

//     pub fn deinit(self: *MarketSection) void {
//         self.allocator.destroy(self);
//     }

//     pub fn widget(self: *MarketSection) vxfw.Widget {
//         return self.container;
//     }

//     pub fn flex(self: *MarketSection) vxfw.FlexItem {
//         return self.flex;
//     }
// };

const HomeScreen = @This();

allocator: std.mem.Allocator,
container: vxfw.Widget,

content: [2]vxfw.FlexItem,
view: vxfw.FlexColumn,

top_item: vxfw.FlexItem,
top_content: vxfw.Text,
bottom_item: vxfw.FlexItem,
bottom_content: vxfw.Text,

ticker_section: *TickerSection,
ticker_item: vxfw.FlexItem,
// market_section: *MarketSection,

pub fn init(allocator: std.mem.Allocator) !*HomeScreen {
    const self = try allocator.create(HomeScreen);
    self.allocator = allocator;

    self.ticker_section = try TickerSection.init(allocator);
    // self.market_section = try MarketSection.init(allocator);
    //

    self.ticker_item = self.ticker_section.flex();

    self.bottom_content = .{ .text = 
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
    self.bottom_item = .{ .widget = self.bottom_content.widget() };

    // BUG: Adding the ticker item to the struct causes it to crash
    self.content = .{ self.ticker_item, self.bottom_item };
    self.view = .{ .children = &self.content };

    self.container = self.view.widget();
    return self;
}

pub fn deinit(self: *HomeScreen) void {
    self.ticker_section.deinit();
    // self.market_section.deinit();

    self.allocator.destroy(self);
}

pub fn widget(self: *HomeScreen) vxfw.Widget {
    return self.container;
}

fn typeErasedDrawFn(ptr: *anyopaque, ctx: vxfw.DrawContext) std.mem.Allocator.Error!vxfw.Surface {
    const self: *HomeScreen = @ptrCast(@alignCast(ptr));
    return self.container.draw(ctx);
}

fn typeErasedEventHandler(ptr: *anyopaque, ctx: *vxfw.EventContext, event: vxfw.Event) anyerror!void {
    const self: *HomeScreen = @ptrCast(@alignCast(ptr));
    _ = ctx; // autofix
    _ = event; // autofix
    _ = self; // autofix

}
