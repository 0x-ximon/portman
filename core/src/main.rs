use std::{env, sync::Arc};

use anyhow::Result;
use axum::{
    Router,
    routing::{get, post},
};
use sqlx::PgPool;
use tower_http::trace::TraceLayer;
use tracing_subscriber::{layer::SubscriberExt, util::SubscriberInitExt};

use crate::{handlers::TickerHandler, repositories::ticker_repo::StorageTickerRepository};

mod handlers;
mod models;
mod repositories;
mod services;

#[derive(Clone)]
struct AppState {
    ticker_repo: StorageTickerRepository,
}

#[tokio::main]
async fn main() -> Result<()> {
    dotenv::dotenv()?;

    let db = PgPool::connect(&env::var("DATABASE_URL")?).await?;
    let state = Arc::new(AppState {
        ticker_repo: StorageTickerRepository::new(db),
    });

    tracing_subscriber::registry()
        .with(
            tracing_subscriber::EnvFilter::try_from_default_env().unwrap_or_else(|_| {
                format!(
                    "{}=debug,tower_http=debug,axum::rejection=trace",
                    env!("CARGO_CRATE_NAME")
                )
                .into()
            }),
        )
        .with(tracing_subscriber::fmt::layer())
        .init();

    let app = Router::new()
        .route("/tickers/", post(TickerHandler::post_ticker))
        .route(
            "/tickers/{symbol}",
            get(TickerHandler::get_ticker).delete(TickerHandler::delete_ticker),
        )
        .layer(TraceLayer::new_for_http())
        .with_state(state);

    let listener = tokio::net::TcpListener::bind("0.0.0.0:3001").await.unwrap();
    axum::serve(listener, app).await.unwrap();

    Ok(())
}
