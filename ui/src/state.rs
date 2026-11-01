#[derive(Debug)]
pub struct User {}

#[derive(Debug)]
pub enum Routes {
    Auth,
    Home,
    Account,
    Trades,
    Config,
}

#[derive(Debug)]
pub struct State {
    pub active: Routes,
    pub user: Option<User>,
}

impl State {
    pub fn new() -> Self {
        Self {
            active: Routes::Auth,
            user: None,
        }
    }

    pub fn set_active(&mut self, route: Routes) {
        self.active = route;
    }
}
