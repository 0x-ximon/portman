use ratatui::{
    buffer::Buffer,
    layout::Rect,
    widgets::{Block, BorderType, Borders, Widget},
};

#[derive(Debug)]
pub struct Router {}

impl Router {
    pub fn new() -> Self {
        Self {}
    }
}

impl Widget for Router {
    fn render(self, area: Rect, buf: &mut Buffer) {
        let content = Block::new()
            .borders(Borders::NONE)
            .border_type(BorderType::Rounded);

        content.render(area, buf);
    }
}
