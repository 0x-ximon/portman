use tonic::transport::Server;

use crate::server::{PormanOrdersServer, portman_server::orders_server::OrdersServer};

mod orders;
mod server;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let addr = "[::1]:50051".parse()?;
    let portman_order_server = PormanOrdersServer::default();

    Server::builder()
        .add_service(OrdersServer::new(portman_order_server))
        .serve(addr)
        .await?;

    Ok(())
}
