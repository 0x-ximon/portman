use anyhow::Result;
use app::App;

mod app;
mod contracts;
mod screens;
mod services;
mod widgets;

fn main() -> Result<()> {
    ratatui::run(|terminal| App::default().run(terminal))
}
