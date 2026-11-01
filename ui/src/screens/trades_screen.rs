use ratatui::{
    buffer::Buffer,
    layout::Rect,
    widgets::{Block, BorderType, Borders, Paragraph, Widget},
};

#[derive(Debug)]
pub struct TradesScreen {}

impl TradesScreen {
    pub fn new() -> Self {
        Self {}
    }
}

impl Widget for &TradesScreen {
    fn render(self, area: Rect, buffer: &mut Buffer) {
        let border = Block::new()
            .borders(Borders::ALL)
            .border_type(BorderType::Rounded);

        let content = Paragraph::new("Trades").block(border);
        content.render(area, buffer);
    }
}
