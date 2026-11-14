use ratatui::{
    buffer::Buffer,
    layout::{Alignment, Constraint, Direction, Flex, Layout, Rect, Size},
    style::Stylize,
    text::Text,
    widgets::Widget,
};

use crate::configs::LOGO;

#[derive(Debug)]
pub struct Banner {
    logo: Text<'static>,
    tag: Text<'static>,
}

impl Banner {
    pub fn new() -> Self {
        let logo = Text::from(LOGO.trim_start_matches('\n'))
            .alignment(Alignment::Center)
            .blue();

        let tag = Text::from("TERMINAL-NATIVE TRADING TERMINAL")
            .alignment(Alignment::Center)
            .dark_gray();

        Self { logo, tag }
    }

    pub fn size(&self) -> Size {
        let width = (self.logo.width().max(self.tag.width())) as u16;
        let height = (self.logo.height() + 1 + self.tag.height()) as u16;
        return Size { width, height };
    }
}

impl Widget for &Banner {
    fn render(self, area: Rect, buffer: &mut Buffer) {
        let layout = Layout::default()
            .direction(Direction::Vertical)
            .constraints(vec![
                Constraint::Length(self.logo.height() as u16),
                Constraint::Length(self.tag.height() as u16),
            ])
            .flex(Flex::Center)
            // .spacing(1)
            .split(area);

        let _ = (&self.logo).render(layout[0], buffer);
        let _ = (&self.tag).render(layout[1], buffer);
    }
}
