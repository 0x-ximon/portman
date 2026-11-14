use ratatui::{buffer::Buffer, layout::Rect, widgets::StatefulWidget};

use crate::{
    meta::Route,
    screens::{
        account_screen::AccountScreen, auth_screen::AuthScreen, config_screen::ConfigScreen,
        home_screen::HomeScreen, trades_screen::TradesScreen,
    },
    state::State,
};

#[derive(Debug)]
pub struct Body {
    auth_screen: AuthScreen,
    home_screen: HomeScreen,
    config_screen: ConfigScreen,
    trades_screen: TradesScreen,
    account_screen: AccountScreen,
}

impl Body {
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

impl StatefulWidget for &Body {
    type State = State;

    fn render(self, area: Rect, buffer: &mut Buffer, state: &mut Self::State) {
        match state.route {
            Route::Auth => self.auth_screen.render(area, buffer, state),
            Route::Home => self.home_screen.render(area, buffer, state),
            Route::Account => self.account_screen.render(area, buffer, state),
            Route::Trades => self.trades_screen.render(area, buffer, state),
            Route::Config => self.config_screen.render(area, buffer, state),
        }
    }
}
