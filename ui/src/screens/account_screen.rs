use ratatui::{buffer::Buffer, layout::Rect, widgets::StatefulWidget};

use crate::state::State;

#[derive(Debug)]
pub struct AccountScreen {}

impl AccountScreen {
    pub fn new() -> Self {
        Self {}
    }
}

impl StatefulWidget for &AccountScreen {
    type State = State;

    fn render(self, area: Rect, buffer: &mut Buffer, state: &mut Self::State) {
        _ = area;
        _ = buffer;
        _ = state;
    }
}
