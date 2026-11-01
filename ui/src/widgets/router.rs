use ratatui::{
    buffer::Buffer,
    layout::Rect,
    widgets::{StatefulWidget, Widget},
};

use crate::{
    screens::{
        account_screen::AccountScreen, auth_screen::AuthScreen, config_screen::ConfigScreen,
        home_screen::HomeScreen, trades_screen::TradesScreen,
    },
    state::{Routes, State},
};

#[derive(Debug)]
pub struct Router {
    auth_screen: AuthScreen,
    home_screen: HomeScreen,
    config_screen: ConfigScreen,
    trades_screen: TradesScreen,
    account_screen: AccountScreen,
}

impl Router {
    pub fn new() -> Self {
        Self {
            auth_screen: AuthScreen::new(),
            home_screen: HomeScreen::new(),
            config_screen: ConfigScreen::new(),
            trades_screen: TradesScreen::new(),
            account_screen: AccountScreen::new(),
        }
    }
}

impl StatefulWidget for &Router {
    type State = State;

    fn render(self, area: Rect, buf: &mut Buffer, state: &mut Self::State) {
        match state.active {
            Routes::Auth => self.auth_screen.render(area, buf),
            Routes::Home => self.home_screen.render(area, buf),
            Routes::Account => self.account_screen.render(area, buf),
            Routes::Trades => self.trades_screen.render(area, buf),
            Routes::Config => self.config_screen.render(area, buf),
        }
    }
}
