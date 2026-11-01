use ratatui::{
    buffer::Buffer,
    layout::Rect,
    widgets::{Block, BorderType, Borders, Paragraph, Widget},
};

// const LOGO: &str = r#"
// ██████╗  ██████╗ ██████╗ ████████╗███╗   ███╗ █████╗ ███╗   ██╗
// ██╔══██╗██╔═══██╗██╔══██╗╚══██╔══╝████╗ ████║██╔══██╗████╗  ██║
// ██████╔╝██║   ██║██████╔╝   ██║   ██╔████╔██║███████║██╔██╗ ██║
// ██╔═══╝ ██║   ██║██╔══██╗   ██║   ██║╚██╔╝██║██╔══██║██║╚██╗██║
// ██║     ╚██████╔╝██║  ██║   ██║   ██║ ╚═╝ ██║██║  ██║██║ ╚████║
// ╚═╝      ╚═════╝ ╚═╝  ╚═╝   ╚═╝   ╚═╝     ╚═╝╚═╝  ╚═╝╚═╝  ╚═══╝"#;
//
// let logo = Paragraph::new(LOGO.trim_start_matches('\n'))
//     .style(Style::new().fg(Color::Rgb(0xF2, 0xA9, 0x3B)))
//     .alignment(Alignment::Center);

#[derive(Debug)]
pub struct AuthScreen {}

impl AuthScreen {
    pub fn new() -> Self {
        Self {}
    }
}

impl Widget for &AuthScreen {
    fn render(self, area: Rect, buffer: &mut Buffer) {
        let border = Block::new()
            .borders(Borders::ALL)
            .border_type(BorderType::Rounded);

        let content = Paragraph::new("Auth").block(border);
        content.render(area, buffer);
    }
}
