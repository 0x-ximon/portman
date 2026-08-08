use std::{io, panic};

use crate::events::{Event, Events};
use anyhow::Result;
use crossterm::event::KeyCode;
use ratatui::{
    backend::CrosstermBackend,
    crossterm::{
        event::{DisableMouseCapture, EnableMouseCapture},
        execute,
        terminal::{self, EnterAlternateScreen, LeaveAlternateScreen},
    },
    layout::Alignment,
    widgets::{Block, BorderType, Borders},
    Frame,
};

pub type Terminal = ratatui::Terminal<CrosstermBackend<io::Stdout>>;

#[derive(Debug)]
pub struct App {
    terminal: Terminal,
    events: Events,
}

impl App {
    pub fn new(terminal: Terminal, events: Events) -> Self {
        Self { terminal, events }
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
            self.terminal.draw(|frame| Self::draw(frame))?;
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

    fn draw(frame: &mut Frame) {
        let widget = Block::new()
            .borders(Borders::ALL)
            .border_type(BorderType::Rounded)
            .title("Portman")
            .title_alignment(Alignment::Left);

        frame.render_widget(widget, frame.area());
    }

    fn reset() -> Result<()> {
        terminal::disable_raw_mode()?;
        execute!(io::stdout(), LeaveAlternateScreen, DisableMouseCapture)?;
        Ok(())
    }
}
