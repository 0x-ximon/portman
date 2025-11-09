use std::sync::Arc;
use tonic::transport::Server;

use crate::{
    orders::OrderBook,
    server::{PortmanOrdersServer, portman_server::orders_server::OrdersServer},
};

mod orders;
mod server;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let order_book = Arc::new(OrderBook::default());

    let addr = "[::1]:50051".parse()?;
    let portman_order_server = PortmanOrdersServer::new(order_book.clone());

    Server::builder()
        .add_service(OrdersServer::new(portman_order_server))
        .serve(addr)
        .await?;

    Ok(())
}
