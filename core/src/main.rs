#![allow(dead_code)]
#![allow(unused_variables)]

use std::sync::Arc;

use axum::{
    Router,
    routing::{get, post},
};
use sqlx::PgPool;

use crate::{handlers::TickerHandler, repositories::ticker_repo::StorageTickerRepository};

mod configs;
mod handlers;
mod models;
mod repositories;
mod services;

#[derive(Clone)]
struct AppState {
    ticker_repo: StorageTickerRepository,
}

#[tokio::main]
async fn main() {
    let db = match PgPool::connect_lazy("Postgres Connection Goes Here") {
        Err(e) => panic!("Connection Failed {:?}", e),
        Ok(v) => v,
    };

    let state = Arc::new(AppState {
        ticker_repo: StorageTickerRepository::new(db),
    });

    let app = Router::new()
        .route("/tickers", post(TickerHandler::post_ticker))
        .route(
            "/tickers/{symbol}",
            get(TickerHandler::get_ticker)
                .put(TickerHandler::put_ticker)
                .delete(TickerHandler::delete_ticker),
        )
        .with_state(state);

    let listener = tokio::net::TcpListener::bind("0.0.0.0:3001").await.unwrap();
    axum::serve(listener, app).await.unwrap();
}
