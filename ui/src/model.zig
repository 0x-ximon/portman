const std = @import("std");

const lib = @import("lib");
const zz = @import("zigzag");

const Model = @This();

pub const Msg = union(enum) {
    click: zz.MouseEvent,
    key: zz.KeyEvent,
};

router: zz.SubProgram(lib.Router, Msg),
hinter: zz.SubProgram(lib.Hinter, Msg),

pub fn init(self: *Model, ctx: *zz.Context) zz.Cmd(Msg) {
    self.router = .{};
    _ = self.router.init(ctx);

    self.hinter = .{};
    _ = self.hinter.init(ctx);

    return .{ .set_title = lib.APP_NAME };
}

pub fn update(self: *Model, msg: Msg, ctx: *zz.Context) zz.Cmd(Msg) {
    switch (msg) {
        .key => |k| {
            switch (k.key) {
                .escape => return .quit,
                else => {},
            }

            self.router.update(.{ .key = k }, ctx);
            self.hinter.update(.{ .key = k }, ctx);
        },

        .click => {},
    }

    return .none;
}

pub fn view(self: *const Model, ctx: *const zz.Context) []const u8 {
    const allocator = ctx.allocator;

    const rows = zz.flex.layout(allocator, ctx.width, ctx.height, &.{
        .{ .constraint = .fill },
        .{ .constraint = .{ .percentage = 10 } },
    }, .{ .direction = .column }) catch return "layout error";

    var router_style = (zz.Style{})
        .borderAll(zz.Border.rounded)
        .background(zz.Colors.black)
        .width(rows[0].width)
        .height(rows[0].height);

    var hinter_style = (zz.Style{})
        .borderAll(zz.Border.rounded)
        .background(zz.Colors.black)
        .width(rows[1].width)
        .height(rows[1].height);

    const router = router_style.render(allocator, self.router.view(ctx)) catch "router style render failed";
    const hinter = hinter_style.render(allocator, self.hinter.view(ctx)) catch "hinter style render failed";

    return zz.join.horizontal(allocator, .top, &.{ router, hinter }) catch "join failed";
}
