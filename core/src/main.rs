#![allow(dead_code)]
#![allow(unused_variables)]

use tonic::transport;
use tower_http::trace::TraceLayer;

use crate::server::{OrdersServer, proto::orders_service_server::OrdersServiceServer};

mod orders;
mod server;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let addr = "[::1]:50051".parse()?;

    tracing_subscriber::fmt()
        .with_max_level(tracing::Level::INFO)
        .init();

    transport::Server::builder()
        .layer(TraceLayer::new_for_grpc())
        .add_service(OrdersServiceServer::new(OrdersServer::new()))
        .serve(addr)
        .await?;

    Ok(())
}
