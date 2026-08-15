use ratatui::{
    buffer::Buffer,
    layout::Rect,
    widgets::{Block, BorderType, Borders, Paragraph, Widget},
};

#[derive(Debug)]
pub struct HomeScreen {}

impl HomeScreen {
    pub fn new() -> Self {
        Self {}
    }
}

impl Widget for &HomeScreen {
    fn render(self, area: Rect, buf: &mut Buffer) {
        let border = Block::new()
            .borders(Borders::ALL)
            .border_type(BorderType::Rounded);

        let content = Paragraph::new("Home").block(border);
        content.render(area, buf);
    }
}
