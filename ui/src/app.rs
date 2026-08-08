use std::{io, panic};

use anyhow::Result;
use crossterm::event::KeyCode;
use ratatui::{
    backend::CrosstermBackend,
    crossterm::{
        event::{DisableMouseCapture, EnableMouseCapture},
        execute,
        terminal::{self, EnterAlternateScreen, LeaveAlternateScreen},
    },
};

use crate::{
    events::{Event, Events},
    model::Model,
};

pub type Terminal = ratatui::Terminal<CrosstermBackend<io::Stdout>>;

#[derive(Debug)]
pub struct AppState {}

#[derive(Debug)]
pub struct App {
    model: Model,
    events: Events,
    state: AppState,
    terminal: Terminal,
}

impl App {
    pub fn new(terminal: Terminal, events: Events) -> Self {
        let model = Model::new();
        let state = AppState {};

        Self {
            model,
            state,
            events,
            terminal,
        }
    }

    pub fn init(&mut self) -> Result<()> {
        terminal::enable_raw_mode()?;
        execute!(io::stdout(), EnterAlternateScreen, EnableMouseCapture)?;

        let panic_hook = panic::take_hook();
        panic::set_hook(Box::new(move |panic| {
            Self::reset().expect("failed to reset the terminal");
            panic_hook(panic);
        }));

        self.terminal.hide_cursor()?;
        self.terminal.clear()?;

        Ok(())
    }

    pub fn deinit(&mut self) -> Result<()> {
        Self::reset()?;
        self.terminal.show_cursor()?;
        Ok(())
    }

    pub fn run(&mut self) -> Result<()> {
        loop {
            self.terminal
                .draw(|frame| frame.render_widget(&self.model, frame.area()))?;

            match self.events.next()? {
                Event::Tick => {}
                Event::Key(k) => {
                    if k.code == KeyCode::Esc {
                        return Ok(());
                    }
                }
                Event::Mouse(m) => {}
                Event::Resize(w, h) => {}
                Event::Quit => return Ok(()),
            }
        }
    }

    fn reset() -> Result<()> {
        terminal::disable_raw_mode()?;
        execute!(io::stdout(), LeaveAlternateScreen, DisableMouseCapture)?;
        Ok(())
    }
}
