use ratatui::{
    buffer::Buffer,
    layout::{Alignment, Constraint, Direction, Flex, Layout, Rect, Size},
    style::{Style, Stylize},
    text::Text,
    widgets::{Block, Borders, Paragraph, StatefulWidget, Widget},
};

use crate::{meta::Action, state::State};

#[derive(Debug)]
pub struct Authenticator {
    heading: Text<'static>,
    subheading: Text<'static>,
}

use crate::configs::LOGO;

impl Authenticator {
    pub fn new() -> Self {
        let heading = Text::from("Portman signs you in with a passkey").centered();
        let subheading = Text::from("This will redirect you to your browser.")
            .centered()
            .dark_gray();

        Self {
            heading,
            subheading,
        }
    }

    pub fn size(&self) -> Size {
        let logo = Text::from(LOGO.trim_start_matches('\n'));

        let width = (logo.width()) as u16;
        let height = (self.heading.height() + self.subheading.height() + 8) as u16;
        return Size { width, height };
    }
}

impl StatefulWidget for &Authenticator {
    type State = State;

    fn render(self, area: Rect, buffer: &mut Buffer, state: &mut Self::State) {
        let block = Block::new()
            .title_alignment(Alignment::Left)
            .borders(Borders::NONE);

        let layout = Layout::default()
            .direction(Direction::Vertical)
            .constraints(vec![
                Constraint::Length(self.heading.height() as u16),
                Constraint::Length(self.subheading.height() as u16),
                Constraint::Length(1),
                Constraint::Length(3),
            ])
            .flex(Flex::Start)
            .split(block.inner(area));

        let hovered = state.hovered == Some(Action::Authenticate);
        let style = match hovered {
            true => Style::new().on_blue(),
            false => Style::new().on_light_blue(),
        };

        let button = Paragraph::new("\nSign in with Passkey")
            .style(style)
            .centered();

        (&self.heading).render(layout[0], buffer);
        (&self.subheading).render(layout[1], buffer);
        (button).render(layout[3], buffer);

        state.hitboxes.push((layout[3], Action::Authenticate));
        block.render(area, buffer);
    }
}
