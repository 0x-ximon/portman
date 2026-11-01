use ratatui::{buffer::Buffer, layout::Rect, widgets::StatefulWidget};

use crate::state::State;

#[derive(Debug)]
pub struct Body {}

impl Body {
    pub fn new() -> Self {
        Self {}
    }
}

impl StatefulWidget for &Body {
    type State = State;

    fn render(self, area: Rect, buffer: &mut Buffer, state: &mut Self::State) {
        let _ = area;
        let _ = buffer;
        let _ = state;

        //         let border = Block::new()
        //             .borders(Borders::ALL)
        //             .border_type(BorderType::Rounded);
        //
        //         let content = Paragraph::new(text).block(border);
        //         content.render(area, buf);
    }
}
