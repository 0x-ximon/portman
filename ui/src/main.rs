use crate::events::Events;
use anyhow::Result;
use app::App;
use ratatui::{backend::CrosstermBackend, Terminal};
use std::io;

mod app;
mod events;
// mod contracts;
// mod screens;
// mod services;
// mod widgets;

fn main() -> Result<()> {
    let stdout = io::stdout();
    let backend = CrosstermBackend::new(stdout);
    let terminal = Terminal::new(backend)?;
    let events = Events::new(250);
    let mut app = App::new(terminal, events);

    app.init()?;
    app.run()?;
    app.deinit()?;

    Ok(())
}
