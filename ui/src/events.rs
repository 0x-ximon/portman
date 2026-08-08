use std::{
    sync::mpsc,
    thread,
    time::{Duration, Instant},
};

use anyhow::Result;
use crossterm::event::{KeyEvent, MouseEvent};
use ratatui::crossterm::event::{self as terminal_event, Event as TerminalEvent};

#[derive(Debug)]
pub enum Event {
    Tick,
    Key(KeyEvent),
    Mouse(MouseEvent),
    Resize(u16, u16),
    Quit,
}

#[derive(Debug)]
pub struct Events {
    sender: mpsc::Sender<Event>,
    receiver: mpsc::Receiver<Event>,
    handler: thread::JoinHandle<()>,
}

impl Events {
    pub fn new(rate: u64) -> Self {
        let interval = Duration::from_millis(rate);
        let (sender, receiver) = mpsc::channel();
        let handler = {
            let sender = sender.clone();
            thread::spawn(move || {
                let mut tick = Instant::now();
                loop {
                    let timeout = interval.checked_sub(tick.elapsed()).unwrap_or(interval);
                    if terminal_event::poll(timeout).is_ok() {
                        if let Ok(event) = terminal_event::read() {
                            let _ = match event {
                                TerminalEvent::Key(k) => sender.send(Event::Key(k)),
                                TerminalEvent::Mouse(m) => sender.send(Event::Mouse(m)),
                                TerminalEvent::Resize(w, h) => sender.send(Event::Resize(w, h)),
                                _ => unimplemented!(),
                            };
                        }
                    }

                    if tick.elapsed() >= interval {
                        let _ = sender.send(Event::Tick);
                        tick = Instant::now();
                    }
                }
            })
        };

        Self {
            sender,
            receiver,
            handler,
        }
    }

    pub fn next(&mut self) -> Result<Event> {
        Ok(self.receiver.recv()?)
    }
}
