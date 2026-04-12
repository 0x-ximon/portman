const std = @import("std");

const zz = @import("zigzag");

pub const Hinter = @This();

content: []const u8,

pub const Msg = union(enum) {
    click: zz.MouseEvent,
    key: zz.KeyEvent,
};

pub fn init(self: *Hinter, _: *zz.Context) zz.Cmd(Msg) {
    self.content = "Hello Hinter";
    return .none;
}

pub fn deinit(_: *Hinter) void {}

pub fn update(_: *Hinter, msg: Msg, _: *zz.Context) zz.Cmd(Msg) {
    switch (msg) {
        .key => |k| switch (k.key) {
            .char => |c| switch (c) {
                else => {},
            },
            else => {},
        },

        .click => {},
    }

    return .none;
}

pub fn view(_: *const Hinter, _: *const zz.Context) []const u8 {
    return "hinter_view";
}
