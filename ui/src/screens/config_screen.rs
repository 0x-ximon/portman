use ratatui::{
    buffer::Buffer,
    layout::Rect,
    widgets::{Block, BorderType, Borders, Paragraph, Widget},
};

#[derive(Debug)]
pub struct ConfigScreen {}

impl ConfigScreen {
    pub fn new() -> Self {
        Self {}
    }
}

impl Widget for &ConfigScreen {
    fn render(self, area: Rect, buf: &mut Buffer) {
        let border = Block::new()
            .borders(Borders::ALL)
            .border_type(BorderType::Rounded);

        let content = Paragraph::new("Config").block(border);
        content.render(area, buf);
    }
}
