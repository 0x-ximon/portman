use ratatui::{
    buffer::Buffer,
    layout::Rect,
    widgets::{Block, BorderType, Borders, StatefulWidget, Widget},
};

use crate::state::State;

#[derive(Debug)]
pub struct Router {}

impl Router {
    pub fn new() -> Self {
        Self {}
    }
}

impl StatefulWidget for &Router {
    type State = State;

    fn render(self, area: Rect, buf: &mut Buffer, state: &mut Self::State) {
        let content = Block::new()
            .borders(Borders::ALL)
            .border_type(BorderType::Rounded);

        content.render(area, buf);
    }
}
