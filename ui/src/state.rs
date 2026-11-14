use ratatui::layout::{Position, Rect};

use crate::meta::{Action, Route, User};

#[derive(Debug)]
pub struct State {
    pub route: Route,
    pub user: Option<User>,
    pub hovered: Option<Action>,
    pub hitboxes: Vec<(Rect, Action)>,
}

impl State {
    pub fn new() -> Self {
        Self {
            user: None,
            hovered: None,
            hitboxes: vec![],
            route: Route::Auth,
        }
    }

    pub fn hit(&self, position: Position) -> Option<Action> {
        self.hitboxes
            .iter()
            .rev()
            .find(|(r, _)| r.contains(position))
            .map(|(_, a)| *a)
    }

    pub fn navigate(&mut self, route: Route) {
        self.route = route;
    }
}
