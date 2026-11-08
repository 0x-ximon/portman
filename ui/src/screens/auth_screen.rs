use ratatui::{
    buffer::Buffer,
    layout::{Constraint, Direction, Layout, Rect, Size},
    widgets::{StatefulWidget, Widget},
};

use crate::{
    state::State,
    widgets::{authenticator::Authenticator, banner::Banner},
};

#[derive(Debug)]
pub struct AuthScreen {
    banner: Banner,
    authenticator: Authenticator,
}

impl AuthScreen {
    pub fn new() -> Self {
        Self {
            banner: Banner::new(),
            authenticator: Authenticator::new(),
        }
    }
}

impl StatefulWidget for &AuthScreen {
    type State = State;

    fn render(self, area: Rect, buffer: &mut Buffer, state: &mut Self::State) {
        let banner_size = self.banner.size();
        let signin_size = self.authenticator.size();
        let size: Size = Size {
            width: banner_size.width.max(signin_size.width),
            height: banner_size.height + signin_size.height + 2,
        };

        let area = area.centered(
            Constraint::Length(size.width),
            Constraint::Length(size.height),
        );

        let layout = Layout::default()
            .direction(Direction::Vertical)
            .constraints(vec![
                Constraint::Length(banner_size.height),
                Constraint::Length(signin_size.height),
            ])
            .spacing(4)
            .split(area);

        let _ = &self.banner.render(layout[0], buffer);
        let _ = &self.authenticator.render(layout[1], buffer, state);
    }
}
