use ratatui::{
    buffer::Buffer,
    layout::Rect,
    widgets::{Block, BorderType, Borders, Paragraph, StatefulWidget, Widget},
};

use crate::{meta::Route, state::State};

#[derive(Debug)]
pub struct Router {}

impl Router {
    pub fn new() -> Self {
        Self {}
    }
}

impl StatefulWidget for &Router {
    type State = State;

    fn render(self, area: Rect, buffer: &mut Buffer, state: &mut Self::State) {
        let text = match state.route {
            Route::Home => "Home",
            Route::Config => "Config",
            Route::Account => "Account",
            _ => "Unknown",
        };

        let border = Block::new()
            .borders(Borders::ALL)
            .border_type(BorderType::Rounded);

        let content = Paragraph::new(text).block(border);
        content.render(area, buffer);
    }
}
