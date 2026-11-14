use ratatui::{
    buffer::Buffer,
    layout::{Constraint, Direction, Layout, Rect},
    widgets::StatefulWidget,
};

use crate::state::State;
use crate::widgets::{body::Body, footer::Footer, header::Header};

#[derive(Debug)]
pub struct Model {
    header: Header,
    body: Body,
    footer: Footer,
}

impl Model {
    pub fn new() -> Self {
        Self {
            header: Header::new(),
            body: Body::new(),
            footer: Footer::new(),
        }
    }
}

impl StatefulWidget for &Model {
    type State = State;

    fn render(self, area: Rect, buffer: &mut Buffer, state: &mut Self::State) {
        let layout = Layout::default()
            .direction(Direction::Vertical)
            .constraints(vec![
                Constraint::Percentage(6),
                Constraint::Percentage(88),
                Constraint::Percentage(6),
            ])
            .split(area);

        state.hitboxes.clear();
        let _ = &self.header.render(layout[0], buffer, state);
        let _ = &self.body.render(layout[1], buffer, state);
        let _ = &self.footer.render(layout[2], buffer, state);
    }
}
