use ratatui::{
    buffer::Buffer,
    layout::Rect,
    widgets::{StatefulWidget, Widget},
};

use crate::{
    screens::{
        account_screen::AccountScreen, config_screen::ConfigScreen, home_screen::HomeScreen,
    },
    state::{Routes, State},
};

#[derive(Debug)]
pub struct Router {
    home_screen: HomeScreen,
    config_screen: ConfigScreen,
    account_screen: AccountScreen,
}

impl Router {
    pub fn new() -> Self {
        Self {
            home_screen: HomeScreen::new(),
            config_screen: ConfigScreen::new(),
            account_screen: AccountScreen::new(),
        }
    }
}

impl StatefulWidget for &Router {
    type State = State;

    fn render(self, area: Rect, buf: &mut Buffer, state: &mut Self::State) {
        match state.active {
            Routes::Home => self.home_screen.render(area, buf),
            Routes::Config => self.config_screen.render(area, buf),
            Routes::Account => self.account_screen.render(area, buf),
        }
    }
}
