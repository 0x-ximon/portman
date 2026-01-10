use std::{str::FromStr, sync::Arc};

use proto::orders_service_server::OrdersService;
use rust_decimal::Decimal;
use tonic::{Request, Response, Status};

use crate::{
    order_book,
    server::proto::{SubmitOrderRequest, SubmitOrderResponse},
};

pub mod proto {
    tonic::include_proto!("proto");
}

pub struct OrdersServer {
    pub order_book: Arc<order_book::OrderBook>,
}

impl OrdersServer {
    pub fn new(order_book: Arc<order_book::OrderBook>) -> Self {
        Self { order_book }
    }
}

#[tonic::async_trait]
impl OrdersService for OrdersServer {
    async fn submit_order(
        &self,
        request: Request<SubmitOrderRequest>,
    ) -> Result<Response<SubmitOrderResponse>, Status> {
        let recv_order = request.into_inner();

        // 3. Construct Internal Order Struct
        let order = crate::order_book::Order {
            id: recv_order.id,

            side: match recv_order.side() {
                proto::Side::Buy => order_book::OrderSide::Buy,
                proto::Side::Sell => order_book::OrderSide::Sell,
                proto::Side::Unspecified => {
                    return Err(Status::invalid_argument("Side must be specified"));
                }
            },

            r#type: match recv_order.r#type() {
                proto::Type::Market => order_book::OrderType::Market,
                proto::Type::Limit => order_book::OrderType::Limit,
                proto::Type::Unspecified => {
                    return Err(Status::invalid_argument("Type must be specified"));
                }
            },

            price: Decimal::from_str(&recv_order.price)
                .map_err(|_| Status::invalid_argument("Invalid price: precision loss or nan"))?,

            quantity: Decimal::from_str(&recv_order.quantity)
                .map_err(|_| Status::invalid_argument("Invalid quantity: precision loss or nan"))?,

            status: match recv_order.status() {
                proto::Status::Pending => order_book::OrderStatus::Pending,
                proto::Status::Fulfilled => order_book::OrderStatus::Fulfilled,
                proto::Status::Cancelled => order_book::OrderStatus::Cancelled,
                proto::Status::Unspecified => {
                    return Err(Status::invalid_argument("Status must be specified"));
                }
            },
        };

        if order.quantity == Decimal::ZERO {
            return Err(Status::invalid_argument("Order quantity cannot be zero."));
        }

        match order.r#type {
            order_book::OrderType::Market => {
                // TODO: Handle case of remaining quantity
                // let (_, _) = self
                //     .order_book
                //     .market_order(order)
                //     .map_err(|e| Status::internal(format!("Order processing error: {}", e)))?;

                // Ok(Response::new(Trades {
                //     trades: trades
                //         .iter()
                //         .map(|t| Trade {
                //             id: t.id,
                //             executed_price: t.executed_price,
                //             order: Some(Order {
                //                 id: t.order.id,
                //                 side: match t.order.side {
                //                     Side::Buy => 0,
                //                     Side::Sell => 1,
                //                 },
                //                 price: t.order.price,
                //                 quantity: t.order.quantity,
                //                 order_type: match t.order.order_type {
                //                     OrderType::Market => 0,
                //                     OrderType::Limit => 1,
                //                 },
                //                 user_id: t.order.user_id,
                //             }),
                //         })
                //         .collect(),
                // }))
            }

            order_book::OrderType::Limit => {
                // self.order_book
                //     .limit_order(order)
                //     .map_err(|e| Status::internal(format!("Order processing error: {}", e)))?;

                // Ok(Response::new(Trades { trades: vec![] }))
            }
        }

        Ok(Response::new(SubmitOrderResponse {
            result: proto::Result::Success as i32,
        }))
    }
}
