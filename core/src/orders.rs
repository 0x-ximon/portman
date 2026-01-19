use std::str::FromStr;
use std::{
    collections::{BTreeMap, VecDeque},
    sync::{RwLock, RwLockWriteGuard},
};

use anyhow::Ok;
use rust_decimal::Decimal;
use serde::{Deserialize, Serialize};

use crate::server::proto;

pub type Symbol = String;
pub type Precision = u32;
pub type OrderPrice = Decimal;
pub type OrderQuantity = Decimal;

#[derive(Copy, Clone, Debug, Deserialize, Serialize)]
pub enum OrderSide {
    Unknown,
    Buy,
    Sell,
}

#[derive(PartialEq, Copy, Clone, Debug, Deserialize, Serialize)]
pub enum OrderType {
    Unknown,
    Market,
    Limit,
}

#[derive(PartialEq, Copy, Clone, Debug, Deserialize, Serialize)]
pub enum OrderStatus {
    Unknown,
    Pending,
    Fulfilled,
    Cancelled,
    Rejected,
}

#[derive(Clone, Debug, Deserialize, Serialize)]
pub struct Order {
    pub id: i64,
    pub side: OrderSide,
    pub price: OrderPrice,
    pub r#type: OrderType,
    pub status: OrderStatus,
    pub filled: OrderQuantity,
    pub quantity: OrderQuantity,
}

impl TryFrom<proto::Order> for Order {
    type Error = anyhow::Error;

    fn try_from(recv_order: proto::Order) -> Result<Self, Self::Error> {
        Ok(Self {
            id: recv_order.id,

            filled: Decimal::ZERO,
            price: Decimal::from_str(&recv_order.price)?,
            quantity: Decimal::from_str(&recv_order.quantity)?,

            side: match recv_order.side() {
                proto::Side::Buy => OrderSide::Buy,
                proto::Side::Sell => OrderSide::Sell,
                proto::Side::Unspecified => OrderSide::Unknown,
            },

            r#type: match recv_order.r#type() {
                proto::Type::Market => OrderType::Market,
                proto::Type::Limit => OrderType::Limit,
                proto::Type::Unspecified => OrderType::Unknown,
            },

            status: match recv_order.status() {
                proto::Status::Pending => OrderStatus::Pending,
                proto::Status::Fulfilled => OrderStatus::Fulfilled,
                proto::Status::Cancelled => OrderStatus::Cancelled,
                proto::Status::Rejected => OrderStatus::Rejected,
                proto::Status::Unspecified => OrderStatus::Unknown,
            },
        })
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct OrderPayload {
    id: i64,
    status: String,
}

impl From<&Order> for OrderPayload {
    fn from(order: &Order) -> Self {
        OrderPayload {
            id: order.id,
            status: match order.status {
                OrderStatus::Pending => "PENDING".to_string(),
                OrderStatus::Fulfilled => "FULFILLED".to_string(),
                OrderStatus::Cancelled => "CANCELLED".to_string(),
                OrderStatus::Rejected => "REJECTED".to_string(),
                OrderStatus::Unknown => "UNKNOWN".to_string(),
            },
        }
    }
}

pub type Orders = Vec<Order>;

#[derive(Debug)]
pub struct Level {
    pub liquidity: Decimal,
    pub orders: VecDeque<Order>,
}

#[derive(Debug)]
pub struct OrderBook {
    pub bids: RwLock<BTreeMap<Decimal, Level>>,
    pub asks: RwLock<BTreeMap<Decimal, Level>>,
    pub quantity_precision: Precision,
    pub price_precision: Precision,
}

impl OrderBook {
    pub fn new(price_precision: Precision, quantity_precision: Precision) -> Self {
        OrderBook {
            bids: RwLock::new(BTreeMap::new()),
            asks: RwLock::new(BTreeMap::new()),
            quantity_precision,
            price_precision,
        }
    }

    pub fn market_order(&self, order: Order) -> anyhow::Result<Orders> {
        match order.side {
            OrderSide::Buy => {
                let mut asks = self.asks.write().map_err(|e| {
                    anyhow::anyhow!("RwLock poisoned while acquiring asks write lock: {}", e)
                })?;

                let orders = self.process_market_order(order, &mut asks)?;
                Ok(orders)
            }
            OrderSide::Sell => {
                let mut bids = self.bids.write().map_err(|e| {
                    anyhow::anyhow!("RwLock poisoned while acquiring bids write lock: {}", e)
                })?;

                let orders = self.process_market_order(order, &mut bids)?;
                Ok(orders)
            }
            OrderSide::Unknown => Err(anyhow::anyhow!("Unknown order side")),
        }
    }

    pub fn limit_order(&self, order: Order) -> anyhow::Result<()> {
        match order.side {
            OrderSide::Buy => {
                let mut bids = self.bids.write().map_err(|e| {
                    anyhow::anyhow!("RwLock poisoned while acquiring asks write lock: {}", e)
                })?;

                self.process_limit_order(order, &mut bids)?;
                Ok(())
            }

            OrderSide::Sell => {
                let mut asks = self.asks.write().map_err(|e| {
                    anyhow::anyhow!("RwLock poisoned while acquiring bids write lock: {}", e)
                })?;

                self.process_limit_order(order, &mut asks)?;
                Ok(())
            }

            OrderSide::Unknown => Err(anyhow::anyhow!("Unknown order side")),
        }
    }

    fn process_market_order(
        &self,
        mut order: Order,
        book: &mut RwLockWriteGuard<BTreeMap<Decimal, Level>>,
    ) -> anyhow::Result<Orders> {
        let mut removable_keys: Vec<Decimal> = Vec::new();
        let mut orders = Orders::new();
        println!("Before: {:?}", book);

        'outer: for (key, level) in book.iter_mut() {
            while let Some(opp_order) = level.orders.front_mut() {
                if order.filled == order.quantity {
                    break 'outer;
                }

                let remaining = order.quantity - order.filled;
                let opp_remaining = opp_order.quantity - opp_order.filled;
                let fill_qty = remaining.min(opp_remaining);

                order.filled += fill_qty;
                opp_order.filled += fill_qty;

                if opp_order.filled == opp_order.quantity {
                    if let Some(mut completed_order) = level.orders.pop_front() {
                        completed_order.status = OrderStatus::Fulfilled;
                        orders.push(completed_order);
                    }
                }

                level.liquidity -= fill_qty;
            }

            if level.liquidity == Decimal::ZERO {
                removable_keys.push(*key);
            }
        }

        if order.filled == order.quantity {
            order.status = OrderStatus::Fulfilled
        } else {
            // TODO: Handle partial fills
            order.status = OrderStatus::Rejected
        }
        orders.push(order);

        for key in removable_keys {
            book.remove(&key);
        }

        println!("After: {:?}", book);
        Ok(orders)
    }

    fn process_limit_order(
        &self,
        order: Order,
        book: &mut RwLockWriteGuard<'_, BTreeMap<Decimal, Level>>,
    ) -> anyhow::Result<()> {
        let level = book.entry(order.price).or_insert_with(|| Level {
            liquidity: Decimal::ZERO,
            orders: VecDeque::new(),
        });

        level.liquidity += order.quantity;
        level.orders.push_back(order);

        Ok(())
    }
}
