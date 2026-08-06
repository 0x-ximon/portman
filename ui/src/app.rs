use anyhow::Result;
use crossterm::event::{self, Event, KeyCode};
use ratatui::style::Style;
use ratatui::symbols;
use ratatui::widgets::{Block, Tabs};
use ratatui::DefaultTerminal;

#[derive(Debug, Default)]
pub struct App {}

impl App {
    pub fn run(&self, terminal: &mut DefaultTerminal) -> Result<()> {
        loop {
            terminal.draw(|frame| self.render(frame))?;
            match event::read()? {
                Event::Key(key) => {
                    if key.is_press() {
                        match key.code {
                            KeyCode::Char('q') => return Ok(()),
                            _ => {}
                        }
                    }
                }

                _ => {}
            }
        }
    }

    fn render(&self, frame: &mut ratatui::Frame) {
        let tabs = Tabs::new(vec!["Tab1", "Tab2", "Tab3", "Tab4"])
            .block(Block::bordered().title("Tabs"))
            .style(Style::default().white())
            .highlight_style(Style::default().yellow())
            .select(2)
            .divider(symbols::DOT)
            .padding("->", "<-");

        frame.render_widget(tabs, frame.area());
    }
}
