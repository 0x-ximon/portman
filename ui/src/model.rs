use ratatui::{
    buffer::Buffer,
    layout::{Constraint, Direction, Layout, Rect},
    widgets::Widget,
};

use crate::widgets::{navigator::Navigator, router::Router};

#[derive(Debug)]
pub struct Model {
    navigator: Navigator,
    router: Router,
}

impl Model {
    pub fn new() -> Self {
        let navigator = Navigator::new();
        let router = Router::new();
        Self { navigator, router }
    }
}

impl Widget for &Model {
    fn render(self, area: Rect, buf: &mut Buffer) {
        let layout = Layout::default()
            .direction(Direction::Vertical)
            .constraints(vec![Constraint::Percentage(10), Constraint::Percentage(90)])
            .split(area);

        let _ = &self.navigator.render(layout[0], buf);
        let _ = &self.router.render(layout[1], buf);
    }
}
