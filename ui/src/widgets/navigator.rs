use ratatui::{
    buffer::Buffer,
    layout::Rect,
    widgets::{Block, BorderType, Borders, Paragraph, StatefulWidget, Widget},
};

use crate::state::{Routes, State};

#[derive(Debug)]
pub struct Navigator {}

impl Navigator {
    pub fn new() -> Self {
        Self {}
    }
}

impl StatefulWidget for &Navigator {
    type State = State;

    fn render(self, area: Rect, buf: &mut Buffer, state: &mut Self::State) {
        let text = match state.active {
            Routes::Home => "Home",
            Routes::Config => "Config",
            Routes::Account => "Account",
        };

        let border = Block::new()
            .borders(Borders::ALL)
            .border_type(BorderType::Rounded);

        let content = Paragraph::new(text).block(border);
        content.render(area, buf);
    }
}
