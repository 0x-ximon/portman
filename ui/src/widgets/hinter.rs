use ratatui::{
    buffer::Buffer,
    layout::Rect,
    widgets::{Block, BorderType, Borders, Paragraph, StatefulWidget, Widget},
};

use crate::state::State;

#[derive(Debug)]
pub struct Hinter {}

impl Hinter {
    pub fn new() -> Self {
        Self {}
    }
}

impl StatefulWidget for &Hinter {
    type State = State;

    fn render(self, area: Rect, buffer: &mut Buffer, state: &mut Self::State) {
        _ = state;

        let border = Block::new()
            .borders(Borders::ALL)
            .border_type(BorderType::Rounded);

        let content = Paragraph::new("Hinter").block(border);
        content.render(area, buffer);
    }
}
