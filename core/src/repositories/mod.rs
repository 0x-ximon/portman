use anyhow::Result;

use crate::models::Ticker;

pub mod ticker_repo;

pub trait TickerRepository: Send + Sync + 'static {
    async fn find_ticker(&self, symbol: &str) -> Result<Ticker>;
    async fn create_ticker(&mut self, ticker: Ticker) -> Result<String>;
    async fn delete_ticker(&mut self, symbol: &str) -> Result<Ticker>;
}
