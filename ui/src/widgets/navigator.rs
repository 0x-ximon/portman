use ratatui::{
    buffer::Buffer,
    layout::Rect,
    widgets::{Block, BorderType, Borders, Widget},
};

#[derive(Debug)]
pub struct Navigator {}

impl Navigator {
    pub fn new() -> Self {
        Self {}
    }
}

impl Widget for Navigator {
    fn render(self, area: Rect, buf: &mut Buffer) {
        let content = Block::new()
            .borders(Borders::NONE)
            .border_type(BorderType::Rounded);

        content.render(area, buf);
    }
}
