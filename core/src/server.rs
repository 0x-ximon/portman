use tonic::{Request, Response, Status};

use portman_server::orders_server::Orders;
use portman_server::{Order, Trade};

pub mod portman_server {
    tonic::include_proto!("orders_package");
}

#[derive(Debug, Default)]
pub struct PormanOrdersServer {}

#[tonic::async_trait]
impl Orders for PormanOrdersServer {
    async fn submit_order(&self, request: Request<Order>) -> Result<Response<Trade>, Status> {
        todo!()
    }
}
