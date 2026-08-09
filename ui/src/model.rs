use ratatui::{
    buffer::Buffer,
    layout::{Constraint, Direction, Layout, Rect},
    widgets::Widget,
};

use crate::widgets::{navigator::Navigator, router::Router};

#[derive(Debug)]
pub struct Model {}

impl Model {
    pub fn new() -> Self {
        Self {}
    }
}

impl Widget for &Model {
    fn render(self, area: Rect, buf: &mut Buffer) {
        let layout = Layout::default()
            .direction(Direction::Vertical)
            .constraints(vec![Constraint::Percentage(10), Constraint::Percentage(90)])
            .split(area);

        let navigator = Navigator::new();
        let router = Router::new();

        navigator.render(layout[0], buf);
        router.render(layout[1], buf);
    }
}
