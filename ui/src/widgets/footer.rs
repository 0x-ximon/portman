use ratatui::{
    buffer::Buffer,
    layout::Rect,
    style::Stylize,
    widgets::{Block, Borders, Paragraph, StatefulWidget, Widget},
};

use crate::state::State;

#[derive(Debug)]
pub struct Footer {}

impl Footer {
    pub fn new() -> Self {
        Self {}
    }
}

impl StatefulWidget for &Footer {
    type State = State;

    fn render(self, area: Rect, buffer: &mut Buffer, state: &mut Self::State) {
        let _ = area;
        let _ = buffer;
        let _ = state;

        let border = Block::new().borders(Borders::TOP).dark_gray();

        let content = Paragraph::new("[Esc] Quit").block(border);
        content.render(area, buffer);
    }
}
