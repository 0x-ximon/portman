#![allow(dead_code)]
#![allow(unused_variables)]

use std::sync::Arc;
use tonic::transport;
use tower_http::trace::TraceLayer;

use crate::{
    order_book::OrderBook, server::OrdersServer,
    server::proto::orders_service_server::OrdersServiceServer,
};

mod order_book;
mod server;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let addr = "[::1]:50051".parse()?;
    let book = Arc::new(OrderBook::default());
    let order_server = OrdersServer::new(book.clone());

    tracing_subscriber::fmt()
        .with_max_level(tracing::Level::INFO)
        .init();

    transport::Server::builder()
        .layer(TraceLayer::new_for_grpc())
        .add_service(OrdersServiceServer::new(order_server))
        .serve(addr)
        .await?;

    Ok(())
}
