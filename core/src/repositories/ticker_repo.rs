use sqlx::PgPool;
use std::error::Error;

use crate::{
    models::{Ticker, TickerSymbol},
    repositories::TickerRepository,
};

#[derive(Debug, Clone)]
pub struct StorageTickerRepository {
    db: PgPool,
}

impl StorageTickerRepository {
    pub fn new(db: PgPool) -> Self {
        Self { db }
    }
}

impl TickerRepository for StorageTickerRepository {
    async fn find_ticker(&self, symbol: &TickerSymbol) -> Result<Ticker, Box<dyn Error>> {
        todo!()
    }

    async fn create_ticker(
        &mut self,
        ticker: Ticker,
        symbol: TickerSymbol,
    ) -> Result<Ticker, Box<dyn Error>> {
        todo!()
    }

    async fn delete_ticker(
        &mut self,
        ticker: Ticker,
        symbol: &TickerSymbol,
    ) -> Result<Ticker, Box<dyn Error>> {
        todo!()
    }

    async fn update_ticker(
        &mut self,
        ticker: Ticker,
        symbol: &TickerSymbol,
    ) -> Result<Ticker, Box<dyn Error>> {
        todo!()
    }
}
