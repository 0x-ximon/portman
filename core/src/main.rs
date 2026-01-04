use std::sync::Arc;
use tonic::transport::Server;

use crate::{
    order_book::OrderBook, server::OrdersServer,
    server::proto::orders_service_server::OrdersServiceServer,
};

mod order_book;
mod server;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let book = Arc::new(OrderBook::default());
    let addr = "[::1]:50051".parse()?;

    let order_server = OrdersServer::new(book.clone());

    Server::builder()
        .add_service(OrdersServiceServer::new(order_server))
        .serve(addr)
        .await?;

    Ok(())
}
