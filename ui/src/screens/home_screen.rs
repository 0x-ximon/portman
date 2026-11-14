use ratatui::{buffer::Buffer, layout::Rect, widgets::StatefulWidget};

use crate::state::State;

#[derive(Debug)]
pub struct HomeScreen {}

impl HomeScreen {
    pub fn new() -> Self {
        Self {}
    }
}

impl StatefulWidget for &HomeScreen {
    type State = State;

    fn render(self, area: Rect, buffer: &mut Buffer, state: &mut Self::State) {
        _ = area;
        _ = buffer;
        _ = state;
    }
}
