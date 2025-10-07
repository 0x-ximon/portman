use std::error::Error;

use crate::models::{Ticker, TickerSymbol};

pub mod ticker_repo;

pub trait TickerRepository: Send + Sync + 'static {
    async fn find_ticker(&self, symbol: &TickerSymbol) -> Result<Ticker, Box<dyn Error>>;

    async fn create_ticker(
        &mut self,
        ticker: Ticker,
        symbol: TickerSymbol,
    ) -> Result<Ticker, Box<dyn Error>>;

    async fn delete_ticker(
        &mut self,
        ticker: Ticker,
        symbol: &TickerSymbol,
    ) -> Result<Ticker, Box<dyn Error>>;

    async fn update_ticker(
        &mut self,
        ticker: Ticker,
        symbol: &TickerSymbol,
    ) -> Result<Ticker, Box<dyn Error>>;
}
