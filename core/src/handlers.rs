use std::sync::Arc;

use axum::{
    Json,
    extract::{Path, State},
};

use crate::{AppState, models::Ticker, repositories::TickerRepository};

#[derive(Debug)]
pub struct TickerHandler {}

impl TickerHandler {
    pub async fn get_ticker(
        State(state): State<Arc<AppState>>,
        Path(symbol): Path<String>,
    ) -> Json<Ticker> {
        match state.ticker_repo.find_ticker(&symbol).await {
            Ok(ticker) => Json(ticker),
            Err(_) => Json(Ticker::default()),
        }
    }

    pub async fn post_ticker(
        State(state): State<Arc<AppState>>,
        Json(payload): Json<Ticker>,
    ) -> Json<String> {
        match state.ticker_repo.clone().create_ticker(payload).await {
            Ok(symbol) => Json(symbol),
            Err(_) => Json(String::new()),
        }
    }

    pub async fn delete_ticker(
        State(state): State<Arc<AppState>>,
        Path(symbol): Path<String>,
    ) -> Json<Ticker> {
        match state.ticker_repo.clone().delete_ticker(&symbol).await {
            Ok(symbol) => Json(symbol),
            Err(_) => Json(Ticker::default()),
        }
    }
}
