use ratatui::{
    buffer::Buffer,
    layout::{Constraint, Direction, Flex, Layout, Rect},
    style::Stylize,
    widgets::{Paragraph, StatefulWidget, Widget},
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
pub struct Banner {}

impl Banner {
    pub fn new() -> Self {
        Self {}
    }
}

impl Widget for &Banner {
    fn render(self, area: Rect, buffer: &mut Buffer) {
        let _ = area;
        let _ = buffer;
    }
}
