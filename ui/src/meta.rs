#[derive(Debug)]
pub struct User {}

#[derive(Debug, Clone, Copy, PartialEq)]
pub enum Action {
    Authenticate,
}

#[derive(Debug)]
pub enum Route {
    Auth,
    Home,
    Account,
    Trades,
    Config,
}
