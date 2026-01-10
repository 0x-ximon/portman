use std::{
    collections::HashMap,
    str::FromStr,
    sync::{Arc, RwLock},
};

use proto::orders_service_server::OrdersService;
use rust_decimal::Decimal;
use tonic::{Request, Response, Status};

use crate::{
    orders::{self},
    server::proto::{
        NewOrderBookRequest, NewOrderBookResponse, SubmitOrderRequest, SubmitOrderResponse,
    },
};

pub mod proto {
    tonic::include_proto!("proto");
}

pub struct OrdersServer {
    pub order_books: RwLock<HashMap<orders::Symbol, Arc<orders::OrderBook>>>,
}

impl OrdersServer {
    pub fn new() -> Self {
        Self {
            order_books: RwLock::new(HashMap::new()),
        }
    }
}

#[tonic::async_trait]
impl OrdersService for OrdersServer {
    async fn submit_order(
        &self,
        request: Request<SubmitOrderRequest>,
    ) -> Result<Response<SubmitOrderResponse>, Status> {
        let recv_order = request.into_inner();

        let order_books = self
            .order_books
            .read()
            .map_err(|e| Status::internal(format!("Failed to read order books: {}", e)))?;

        let symbol = orders::Symbol(recv_order.symbol.to_string());

        if let Some(order_book) = order_books.get(&symbol) {
            let order = orders::Order {
                id: recv_order.id,

                side: match recv_order.side() {
                    proto::Side::Buy => orders::OrderSide::Buy,
                    proto::Side::Sell => orders::OrderSide::Sell,
                    proto::Side::Unspecified => {
                        return Err(Status::invalid_argument("Side must be specified"));
                    }
                },

                r#type: match recv_order.r#type() {
                    proto::Type::Market => orders::OrderType::Market,
                    proto::Type::Limit => orders::OrderType::Limit,
                    proto::Type::Unspecified => {
                        return Err(Status::invalid_argument("Type must be specified"));
                    }
                },

                price: Decimal::from_str(&recv_order.price).map_err(|_| {
                    Status::invalid_argument("Invalid price: precision loss or nan")
                })?,

                quantity: Decimal::from_str(&recv_order.quantity).map_err(|_| {
                    Status::invalid_argument("Invalid quantity: precision loss or nan")
                })?,

                status: match recv_order.status() {
                    proto::Status::Pending => orders::OrderStatus::Pending,
                    proto::Status::Fulfilled => orders::OrderStatus::Fulfilled,
                    proto::Status::Cancelled => orders::OrderStatus::Cancelled,
                    proto::Status::Unspecified => {
                        return Err(Status::invalid_argument("Status must be specified"));
                    }
                },
            };

            if order.quantity == Decimal::ZERO {
                return Err(Status::invalid_argument("Order quantity cannot be zero."));
            }

            match order.r#type {
                orders::OrderType::Market => {
                    order_book
                        .market_order(order)
                        .map_err(|e| Status::internal(format!("Order processing error: {}", e)))?;
                }

                orders::OrderType::Limit => {
                    order_book
                        .limit_order(order)
                        .map_err(|e| Status::internal(format!("Order processing error: {}", e)))?;
                }
            }
        } else {
            return Err(Status::internal(format!(
                "Order book not found for symbol {:?}",
                symbol
            )));
        }

        Ok(Response::new(SubmitOrderResponse {
            result: proto::Result::Success as i32,
        }))
    }

    async fn new_order_book(
        &self,
        request: Request<NewOrderBookRequest>,
    ) -> Result<Response<NewOrderBookResponse>, Status> {
        let recv_order_book = request.into_inner();

        let symbol = orders::Symbol(recv_order_book.symbol.to_owned());
        let precision = recv_order_book.precision;

        let order_book = Arc::new(orders::OrderBook::new(precision));
        let mut order_books = self
            .order_books
            .write()
            .map_err(|e| Status::internal(format!("Failed to write order books: {}", e)))?;

        order_books.insert(symbol, order_book);

        Ok(Response::new(NewOrderBookResponse {
            result: proto::Result::Success as i32,
        }))
    }
}
