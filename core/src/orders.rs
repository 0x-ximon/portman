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
}

#[derive(Clone, Debug, Deserialize, Serialize)]
pub struct Order {
    pub id: i64,
    pub side: OrderSide,
    pub price: OrderPrice,
    pub r#type: OrderType,
    pub status: OrderStatus,
    pub quantity: OrderQuantity,
}

impl TryFrom<proto::Order> for Order {
    type Error = anyhow::Error;

    fn try_from(recv_order: proto::Order) -> Result<Self, Self::Error> {
        Ok(Self {
            id: recv_order.id,

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
                proto::Status::Unspecified => OrderStatus::Unknown,
            },
        })
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
        let mut remaining_quantity = order.quantity;
        let mut orders = Orders::new();

        for (key, level) in book.iter_mut() {
            if remaining_quantity == Decimal::ZERO {
                break;
            }

            while let Some(opposite_order) = level.orders.front_mut() {
                if remaining_quantity == Decimal::ZERO {
                    break;
                }

                let fill_qty = opposite_order.quantity.min(remaining_quantity);

                remaining_quantity -= fill_qty;
                opposite_order.quantity -= fill_qty;
                level.liquidity -= fill_qty;

                if remaining_quantity == Decimal::ZERO {
                    // BUG: Include opposite orders when done
                } else if opposite_order.quantity == Decimal::ZERO {
                    if let Some(completed_order) = level.orders.pop_front() {
                        orders.push(completed_order);
                    }
                }
            }

            if level.liquidity == Decimal::ZERO {
                removable_keys.push(*key);
            }
        }

        order.status = OrderStatus::Fulfilled;
        orders.push(order);

        for key in removable_keys {
            book.remove(&key);
        }

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
