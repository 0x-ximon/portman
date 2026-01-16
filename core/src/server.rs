use std::{
    collections::HashMap,
    str::FromStr,
    sync::{Arc, RwLock},
};

use proto::orders_service_server::OrdersService;
use questdb::ingress::Sender;
use rust_decimal::Decimal;
use tonic::{Request, Response, Status};

use crate::{
    orders,
    server::proto::{
        NewOrderBookRequest, NewOrderBookResponse, SubmitOrderRequest, SubmitOrderResponse,
    },
    store::Store,
};

pub mod proto {
    tonic::include_proto!("proto");
}

pub struct OrdersServer {
    pub connection: Option<RwLock<Sender>>,
    pub order_books: RwLock<HashMap<orders::Symbol, Arc<orders::OrderBook>>>,
}

impl OrdersServer {
    pub fn new(conn: Option<RwLock<Sender>>) -> Self {
        Self {
            connection: conn,
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
        let Some(ref conn) = self.connection else {
            return Err(Status::unavailable("Connection not available"));
        };

        let recv_order = request.into_inner();

        let order_books = self
            .order_books
            .read()
            .map_err(|e| Status::internal(format!("Failed to read order books: {}", e)))?;

        let symbol = recv_order.symbol.to_string();

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

                price: Decimal::from_str(&recv_order.price)
                    .map_err(|_| Status::invalid_argument("Price: precision loss or nan"))?
                    .round_dp(order_book.price_precision),

                quantity: Decimal::from_str(&recv_order.quantity)
                    .map_err(|_| Status::invalid_argument("Quantity: precision loss or nan"))?
                    .round_dp(order_book.quantity_precision),

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

            let store = Store::new(conn);

            match order.r#type {
                orders::OrderType::Market => {
                    let orders = order_book
                        .market_order(order)
                        .map_err(|e| Status::internal(format!("Order processing error: {}", e)))?;

                    store
                        .save_orders(&symbol, &orders)
                        .map_err(|e| Status::internal(format!("Failed to save orders: {}", e)))?;
                }

                orders::OrderType::Limit => {
                    order_book
                        .limit_order(order)
                        .map_err(|e| Status::internal(format!("Order processing error: {}", e)))?;
                }

                orders::OrderType::Unknown => {
                    return Err(Status::invalid_argument("Unknown order type"));
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

        let price_precision = recv_order_book.price_precision;
        let quantity_precision = recv_order_book.quantity_precision;
        let symbol = recv_order_book.symbol.to_owned();

        let order_book = Arc::new(orders::OrderBook::new(price_precision, quantity_precision));
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

#[cfg(test)]
mod server_tests {
    use super::*;

    #[test]
    fn test_submit_order() {
        tokio::runtime::Runtime::new().unwrap().block_on(async {
            let server = OrdersServer::new(None);
            let symbol = "BTC/USD".to_string();
            let precision = 2;

            let request = Request::new(NewOrderBookRequest {
                symbol: symbol.to_owned(),
                price_precision: precision,
                quantity_precision: precision,
            });
            let response = server.new_order_book(request).await.unwrap();

            assert_eq!(response.into_inner().result, proto::Result::Success as i32);

            let order_books = server.order_books.read().unwrap();
            assert!(order_books.contains_key(&symbol));

            let payload = orders::Order {
                id: 1,
                side: orders::OrderSide::Buy,
                price: Decimal::new(200504, 3),
                quantity: Decimal::new(1004, 3),
                r#type: orders::OrderType::Limit,
                status: orders::OrderStatus::Pending,
            };

            let request = Request::new(SubmitOrderRequest {
                id: payload.id,
                side: payload.side as i32,
                r#type: payload.r#type as i32,
                status: payload.status as i32,
                price: payload.price.to_string(),
                quantity: payload.quantity.to_string(),
                symbol: symbol.to_owned(),
            });

            let response = server.submit_order(request).await.unwrap();
            assert_eq!(response.into_inner().result, proto::Result::Success as i32);

            let order_book = order_books.get(&symbol).unwrap();
            let mut bids = order_book.bids.write().unwrap();
            let asks = order_book.asks.read().unwrap();

            assert_eq!(bids.len(), 1);
            assert_eq!(asks.len(), 0);

            let (_, mut level) = bids.pop_first().unwrap();
            let order = level.orders.pop_front().unwrap();

            assert_eq!(order.id, payload.id);
            assert_eq!(order.price, Decimal::new(20050, precision));
            assert_eq!(order.quantity, Decimal::new(100, precision));
        });
    }
}
