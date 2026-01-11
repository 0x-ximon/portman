use std::{
    collections::{BTreeMap, VecDeque},
    sync::{RwLock, RwLockWriteGuard},
};

use anyhow::Ok;
use rust_decimal::Decimal;

pub type Symbol = String;
pub type Precision = u32;
pub type OrderPrice = Decimal;
pub type OrderQuantity = Decimal;

#[derive(Copy, Clone, Debug)]
pub enum OrderSide {
    Unknown,
    Buy,
    Sell,
}

#[derive(PartialEq, Copy, Clone, Debug)]
pub enum OrderType {
    Unknown,
    Market,
    Limit,
}

#[derive(PartialEq, Copy, Clone, Debug)]
pub enum OrderStatus {
    Unknown,
    Pending,
    Fulfilled,
    Cancelled,
}

#[derive(Clone, Debug)]
pub struct Order {
    pub id: i64,
    pub side: OrderSide,
    pub price: OrderPrice,
    pub r#type: OrderType,
    pub status: OrderStatus,
    pub quantity: OrderQuantity,
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
