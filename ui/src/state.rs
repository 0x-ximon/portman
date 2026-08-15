#[derive(Debug)]
pub enum Routes {
    Home,
    Account,
    Config,
}

#[derive(Debug)]
pub struct State {
    pub active: Routes,
}

impl State {
    pub fn new() -> Self {
        Self {
            active: Routes::Home,
        }
    }

    pub fn set_active(&mut self, route: Routes) {
        self.active = route;
    }
}
