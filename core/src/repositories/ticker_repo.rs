use anyhow::Result;
use sqlx::PgPool;

use crate::{models::Ticker, repositories::TickerRepository};

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
    async fn find_ticker(&self, symbol: &str) -> Result<Ticker> {
        let r = sqlx::query!(r#"SELECT * FROM Tickers WHERE symbol = $1;"#, symbol)
            .fetch_one(&self.db)
            .await?;

        Ok(Ticker {
            symbol: r.symbol,
            name: r.name,
            bid: r.bid,
            last: r.last,
            ask: r.ask,
        })
    }

    async fn create_ticker(&mut self, ticker: Ticker) -> Result<String> {
        let r = sqlx::query!(
            r#"INSERT INTO Tickers (symbol, name, last, ask, bid) VALUES ($1, $2, $3, $4, $5) RETURNING symbol"#,
            &ticker.symbol,
            &ticker.name,
            ticker.last,
            ticker.ask,
            ticker.bid
        )
        .fetch_one(&self.db)
        .await?;

        Ok(r.symbol)
    }

    async fn delete_ticker(&mut self, symbol: &str) -> Result<Ticker> {
        let r = sqlx::query!(
            r#"DELETE FROM Tickers where symbol = $1 RETURNING symbol, name, last, ask, bid"#,
            &symbol
        )
        .fetch_one(&self.db)
        .await?;

        Ok(Ticker {
            symbol: r.symbol,
            name: r.name,
            bid: r.bid,
            last: r.last,
            ask: r.ask,
        })
    }
}
