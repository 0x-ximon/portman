use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug, Default, Clone)]
pub struct Ticker {
    symbol: TickerSymbol,
    name: TickerName,

    last: TickerValue,
    ask: TickerValue,
    bid: TickerValue,
}

#[derive(Serialize, Deserialize, Debug, Default, Hash, PartialEq, Eq, Clone)]
pub struct TickerSymbol(String);

impl TickerSymbol {
    pub fn new(quote: &str, base: &str) -> Self {
        Self {
            0: format!("{:?}/{:?}", quote.to_uppercase(), base.to_uppercase()),
        }
    }
}

#[derive(Serialize, Deserialize, Debug, Default, Clone)]
pub struct TickerName(String);

#[derive(Serialize, Deserialize, Debug, Default, Clone)]
pub struct TickerValue(f64);
