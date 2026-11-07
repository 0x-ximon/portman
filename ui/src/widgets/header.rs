use ratatui::{
    buffer::Buffer,
    layout::{Constraint, Direction, Flex, Layout, Rect},
    style::Stylize,
    widgets::{Block, Borders, Padding, Paragraph, StatefulWidget, Widget},
};

use crate::{
    configs::APP_VERSION,
    state::{Routes, State},
    widgets::router::Router,
};

#[derive(Debug)]
pub struct Header {
    router: Router,
}

impl Header {
    pub fn new() -> Self {
        Self {
            router: Router::new(),
        }
    }
}

impl StatefulWidget for &Header {
    type State = State;

    fn render(self, area: Rect, buffer: &mut Buffer, state: &mut Self::State) {
        let block = Block::new()
            .padding(Padding::symmetric(2, 0))
            .borders(Borders::BOTTOM)
            .dark_gray();

        match state.active {
            Routes::Auth => {
                let layout = Layout::default()
                    .direction(Direction::Horizontal)
                    .constraints(vec![Constraint::Percentage(20), Constraint::Percentage(20)])
                    .flex(Flex::SpaceBetween)
                    .split(block.inner(area));

                let text = format!("portman v{}", APP_VERSION);
                Paragraph::new(text)
                    .dark_gray()
                    .left_aligned()
                    .render(layout[0], buffer);

                let text = "Signed Out";
                Paragraph::new(text)
                    .dark_gray()
                    .right_aligned()
                    .render(layout[1], buffer);
            }

            Routes::Home => todo!(),
            Routes::Account => todo!(),
            Routes::Trades => todo!(),
            Routes::Config => todo!(),
        };

        block.render(area, buffer);
    }
}
