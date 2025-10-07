use std::sync::Arc;

use axum::{
    Json,
    extract::{Path, State},
};

use crate::{
    AppState,
    models::{Ticker, TickerSymbol},
    repositories::TickerRepository,
};

#[derive(Debug)]
pub struct TickerHandler {}

impl TickerHandler {
    pub async fn get_ticker(
        State(state): State<Arc<AppState>>,
        Path(symbol): Path<TickerSymbol>,
    ) -> Json<Ticker> {
        match state.ticker_repo.find_ticker(&symbol).await {
            Ok(ticker) => Json(ticker),
            Err(_) => Json(Ticker::default()),
        }
    }

    pub async fn post_ticker(
        State(state): State<Arc<AppState>>,
        Path(symbol): Path<TickerSymbol>,
        Json(payload): Json<Ticker>,
    ) -> Json<Ticker> {
        match state
            .ticker_repo
            .clone()
            .create_ticker(payload, symbol)
            .await
        {
            Ok(ticker) => Json(ticker),
            Err(_) => Json(Ticker::default()),
        }
    }

    pub async fn put_ticker(
        State(state): State<Arc<AppState>>,
        Path(symbol): Path<TickerSymbol>,
        Json(payload): Json<Ticker>,
    ) -> Json<Ticker> {
        match state
            .ticker_repo
            .clone()
            .update_ticker(payload, &symbol)
            .await
        {
            Ok(ticker) => Json(ticker),
            Err(_) => Json(Ticker::default()),
        }
    }

    pub async fn delete_ticker(
        State(state): State<Arc<AppState>>,
        Path(symbol): Path<TickerSymbol>,
    ) -> Json<Ticker> {
        Json(Ticker::default())
    }
}
