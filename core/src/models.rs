use serde::{Deserialize, Serialize};
use sqlx::prelude::FromRow;

#[derive(Serialize, Deserialize, Debug, Default, Clone, FromRow)]
pub struct Ticker {
    pub symbol: String,
    pub name: String,
    pub last: f32,
    pub ask: f32,
    pub bid: f32,
}
