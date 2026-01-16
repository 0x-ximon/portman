use std::sync::RwLock;

use questdb::ingress::{Sender, TimestampNanos};

use crate::orders::Orders;

pub struct Store<'a> {
    timeseries_conn: &'a RwLock<Sender>,
}

impl<'a> Store<'a> {
    pub fn new(timeseries_conn: &'a RwLock<Sender>) -> Self {
        Self { timeseries_conn }
    }

    pub fn save_orders(&self, symbol: &str, orders: &Orders) -> anyhow::Result<()> {
        let mut conn = self.timeseries_conn.write().map_err(|e| {
            anyhow::anyhow!("RwLock poisoned while acquiring conn write lock: {}", e)
        })?;

        let mut buffer = conn.new_buffer();

        for order in orders {
            let side = match order.side {
                crate::orders::OrderSide::Buy => "buy",
                crate::orders::OrderSide::Sell => "sell",
                _ => "unknown",
            };

            buffer
                .table("trades")?
                .symbol("symbol", symbol)?
                .symbol("side", side)?
                .column_f64("price", order.price.try_into()?)?
                .column_f64("quantity", order.quantity.try_into()?)?
                .at(TimestampNanos::now())?;
        }

        conn.flush(&mut buffer)
            .map_err(|e| anyhow::anyhow!("Failed to flush buffer: {}", e))?;

        Ok(())
    }
}
