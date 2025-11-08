use std::sync::Arc;

use tonic::{Request, Response, Status};

use portman_server::orders_server::Orders;
use portman_server::{Order, Trade};

use crate::orders::{OrderBook, OrderType, Side};

pub mod portman_server {
    tonic::include_proto!("orders_package");
}

#[derive(Debug, Clone)]
pub struct PortmanOrdersServer {
    pub order_book: Arc<OrderBook>,
}

impl PortmanOrdersServer {
    pub fn new(order_book: Arc<OrderBook>) -> Self {
        Self { order_book }
    }
}

#[tonic::async_trait]
impl Orders for PortmanOrdersServer {
    async fn submit_order(&self, request: Request<Order>) -> Result<Response<Trade>, Status> {
        let recv_order = request.into_inner();

        let order = crate::orders::Order {
            id: 0,
            side: match recv_order.side {
                0 => Side::Buy,
                1 => Side::Sell,
                _ => return Err(Status::invalid_argument("Invalid side value")),
            },

            price: recv_order.price,
            quantity: recv_order.quantity,

            order_type: match recv_order.order_type {
                0 => OrderType::Market,
                1 => OrderType::Limit,
                _ => return Err(Status::invalid_argument("Invalid order type")),
            },

            user_id: recv_order.user_id,
        };

        match order.order_type {
            OrderType::Market => {
                self.order_book
                    .market_order(order)
                    .map_err(|e| Status::internal(format!("Order processing error: {}", e)))?;

                todo!()
            }

            OrderType::Limit => {
                self.order_book
                    .market_order(order)
                    .map_err(|e| Status::internal(format!("Order processing error: {}", e)))?;

                todo!()
            }
        };
    }
}
