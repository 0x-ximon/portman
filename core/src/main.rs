#![allow(dead_code)]
#![allow(unused_variables)]

use std::sync::{ RwLock};

use tonic::transport;
use tower_http::trace::TraceLayer;

use crate::server::{OrdersServer, proto::orders_service_server::OrdersServiceServer};

mod config;
mod orders;
mod server;
mod store;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let cfg = config::Config::new()?;
    let conn = RwLock::new(cfg.conn);
    let addr = cfg.addr;

    tracing_subscriber::fmt()
        .with_max_level(tracing::Level::INFO)
        .init();

    transport::Server::builder()
        .layer(TraceLayer::new_for_grpc())
        .add_service(OrdersServiceServer::new(OrdersServer::new(Some(conn))))
        .serve(addr)
        .await?;

    Ok(())
}
