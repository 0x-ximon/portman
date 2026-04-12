const std = @import("std");

const lib = @import("lib");
const zz = @import("zigzag");

const Model = @This();

router: *lib.Router,
hinter: *lib.Hinter,

pub const Msg = union(enum) {
    click: zz.MouseEvent,
    key: zz.KeyEvent,
};

pub fn init(self: *Model, ctx: *zz.Context) zz.Cmd(Msg) {
    self.router = lib.Router.init(ctx.persistent_allocator) catch |err| {
        std.log.err("router init failed: {}", .{err});
        return .quit;
    };

    self.hinter = lib.Hinter.init(ctx.persistent_allocator) catch |err| {
        std.log.err("hinter init failed: {}", .{err});
        return .quit;
    };

    return .{ .set_title = "Portman" };
}

pub fn deinit(self: *Model) void {
    self.router.deinit();
    self.hinter.deinit();
}

pub fn update(self: *Model, msg: Msg, _: *zz.Context) zz.Cmd(Msg) {
    switch (msg) {
        .key => |k| {
            switch (k.key) {
                .escape => return .quit,
                else => {},
            }

            self.router.react(.{ .key = k });
            self.hinter.react(.{ .key = k });
        },

        .click => |c| {
            self.router.react(.{ .click = c });
            self.hinter.react(.{ .click = c });
        },
    }

    return .none;
}

pub fn view(self: *const Model, ctx: *const zz.Context) []const u8 {
    const allocator = ctx.allocator;

    const rows = zz.flex.layout(allocator, ctx.width, ctx.height, &.{
        .{ .constraint = .fill },
        .{ .constraint = .{ .percentage = 10 } },
    }, .{ .direction = .column }) catch return "layout error";

    const router = self.router.render(allocator, rows[0].width, rows[0].height);
    const hinter = self.hinter.render(allocator, rows[1].width, rows[1].height);

    return zz.join.horizontal(allocator, .top, &.{ router, hinter }) catch "join failed";
}
