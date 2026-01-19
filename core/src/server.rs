use std::{
    collections::HashMap,
    sync::{Arc, RwLock},
};

use async_nats::jetstream::Context;
use proto::orders_service_server::OrdersService;
use questdb::ingress::Sender;
use rust_decimal::Decimal;
use tonic::{Request, Response, Status};

use crate::{
    orders::{self},
    server::proto::{
        NewOrderBookRequest, NewOrderBookResponse, SubmitOrderRequest, SubmitOrderResponse,
    },
    store::Store,
};

pub mod proto {
    tonic::include_proto!("proto");
}

pub struct OrdersServer {
    pub publisher: Option<Context>,
    pub connection: Option<RwLock<Sender>>,
    pub order_books: RwLock<HashMap<orders::Symbol, Arc<orders::OrderBook>>>,
}

impl OrdersServer {
    pub fn new(conn: Option<RwLock<Sender>>, publisher: Option<Context>) -> Self {
        Self {
            publisher,
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

        let Some(ref publisher) = self.publisher else {
            return Err(Status::unavailable("Publisher not available"));
        };

        let inner = request.into_inner();
        let symbol = inner.symbol;

        let Some(recv_order) = inner.order else {
            return Err(Status::invalid_argument("order is missing"));
        };

        let order_book = {
            let order_books = self
                .order_books
                .read()
                .map_err(|e| Status::internal(format!("Failed to read order books: {}", e)))?;

            order_books.get(&symbol).cloned()
        };

        if let Some(order_book) = order_book {
            let mut order: orders::Order = recv_order
                .try_into()
                .map_err(|e| Status::invalid_argument(format!("Could not parse order: {}", e)))?;

            if order.quantity == Decimal::ZERO {
                return Err(Status::invalid_argument("Order quantity cannot be zero."));
            }

            order.price = order.price.round_dp(order_book.price_precision);
            order.quantity = order.quantity.round_dp(order_book.quantity_precision);

            let store = Store::new(conn);

            match order.r#type {
                orders::OrderType::Market => {
                    let orders = order_book
                        .market_order(order)
                        .map_err(|e| Status::internal(format!("Order processing error: {}", e)))?;

                    store
                        .save_orders(&symbol, &orders)
                        // TODO: Reinsert the orders into the order book if failed to save to timeseries DB
                        .map_err(|e| Status::internal(format!("Failed to save orders: {}", e)))?;

                    let orders_payload: Vec<orders::OrderPayload> =
                        orders.iter().map(|order| order.into()).collect();

                    // PERF: Switch from JSON to Protobuf
                    let payload = serde_json::to_vec(&orders_payload).map_err(|e| {
                        Status::internal(format!("Failed to serialize orders: {}", e))
                    })?;

                    let ack = publisher
                        .publish("orders.processed", payload.into())
                        .await
                        .map_err(|e| {
                            Status::internal(format!(
                                "Failed to publish order processed event: {}",
                                e
                            ))
                        })?;

                    ack.await.map_err(|e| {
                        Status::internal(format!("Acknowledgement not received: {}", e))
                    })?;
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

        if order_books.contains_key(&symbol) {
            return Err(Status::already_exists(format!(
                "{} order book already exists",
                symbol
            )));
        }

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
            let server = OrdersServer::new(None, None);
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

            // let payload = orders::Order {
            //     id: 1,
            //     side: orders::OrderSide::Buy,
            //     price: Decimal::new(200504, 3),
            //     quantity: Decimal::new(1004, 3),
            //     r#type: orders::OrderType::Limit,
            //     status: orders::OrderStatus::Pending,
            // };

            // let request = Request::new(SubmitOrderRequest {
            //     id: payload.id,
            //     side: payload.side as i32,
            //     r#type: payload.r#type as i32,
            //     status: payload.status as i32,
            //     price: payload.price.to_string(),
            //     quantity: payload.quantity.to_string(),
            //     symbol: symbol.to_owned(),
            // });

            // // TEST: Mock the Timeseries DB
            // let response = server.submit_order(request).await.unwrap();
            // assert_eq!(response.into_inner().result, proto::Result::Success as i32);

            // let order_book = order_books.get(&symbol).unwrap();
            // let mut bids = order_book.bids.write().unwrap();
            // let asks = order_book.asks.read().unwrap();

            // assert_eq!(bids.len(), 1);
            // assert_eq!(asks.len(), 0);

            // let (_, mut level) = bids.pop_first().unwrap();
            // let order = level.orders.pop_front().unwrap();

            // assert_eq!(order.id, payload.id);
            // assert_eq!(order.price, Decimal::new(20050, precision));
            // assert_eq!(order.quantity, Decimal::new(100, precision));
        });
    }
}
