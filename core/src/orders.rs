use std::{
    collections::{BTreeMap, VecDeque},
    sync::{RwLock, RwLockWriteGuard, atomic::AtomicU64},
};

use rust_decimal::Decimal;

#[derive(Eq, PartialEq, Hash, Clone, Debug)]
pub struct Symbol(pub String);

pub type OrderPrice = Decimal;
pub type OrderQuantity = Decimal;

#[derive(Copy, Clone, Debug)]
pub enum OrderSide {
    Buy,
    Sell,
}

#[derive(PartialEq, Copy, Clone, Debug)]
pub enum OrderType {
    Market,
    Limit,
}

#[derive(PartialEq, Copy, Clone, Debug)]
pub enum OrderStatus {
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

pub type Trades = Vec<Order>;

#[derive(Debug)]
pub struct PriceLevel {
    pub price: OrderPrice,
    pub quantity: OrderQuantity,
    pub orders: VecDeque<Order>,
}

#[derive(Debug)]
pub struct OrderBook {
    pub bids: RwLock<BTreeMap<i64, PriceLevel>>,
    pub asks: RwLock<BTreeMap<i64, PriceLevel>>,
    pub liquidity: AtomicU64,
    pub precision: i32,
}

impl OrderBook {
    pub fn new(precision: i32) -> Self {
        OrderBook {
            bids: RwLock::new(BTreeMap::new()),
            asks: RwLock::new(BTreeMap::new()),
            liquidity: AtomicU64::new(0),
            precision,
        }
    }

    pub fn market_order(&self, order: Order) -> anyhow::Result<((), Decimal)> {
        match order.side {
            OrderSide::Buy => {
                let mut asks = self.asks.write().map_err(|e| {
                    anyhow::anyhow!("RwLock poisoned while acquiring asks write lock: {}", e)
                })?;
                let (trades, remaining_quantity) = self.process_market_order(order, &mut asks);
                Ok(((), remaining_quantity))
            }
            OrderSide::Sell => {
                let mut bids = self.bids.write().map_err(|e| {
                    anyhow::anyhow!("RwLock poisoned while acquiring bids write lock: {}", e)
                })?;
                let (trades, remaining_quantity) = self.process_market_order(order, &mut bids);
                Ok(((), remaining_quantity))
            }
        }
    }

    pub fn limit_order(&self, order: Order) -> anyhow::Result<()> {
        match order.side {
            OrderSide::Buy => {
                let mut bids = self.bids.write().map_err(|e| {
                    anyhow::anyhow!("RwLock poisoned while acquiring asks write lock: {}", e)
                })?;

                self.process_limit_order(order, &mut bids)
            }

            OrderSide::Sell => {
                let mut asks = self.asks.write().map_err(|e| {
                    anyhow::anyhow!("RwLock poisoned while acquiring bids write lock: {}", e)
                })?;

                self.process_limit_order(order, &mut asks)
            }
        }
    }

    fn process_limit_order(
        &self,
        order: Order,
        book: &mut RwLockWriteGuard<'_, BTreeMap<i64, PriceLevel>>,
    ) -> anyhow::Result<()> {
        let key = self.price_to_key(order.price);
        let level = book.entry(key).or_insert_with(|| PriceLevel {
            quantity: Decimal::ZERO,
            price: order.price,
            orders: VecDeque::new(),
        });

        level.quantity += order.quantity;
        level.orders.push_back(order);

        Ok(())
    }

    fn process_market_order(
        &self,
        order: Order,
        book: &mut RwLockWriteGuard<BTreeMap<i64, PriceLevel>>,
    ) -> (Vec<()>, Decimal) {
        let mut removable_keys: Vec<i64> = Vec::new();
        let mut remaining_quantity = order.quantity;
        let trades: Vec<()> = Vec::new();

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
                level.quantity -= fill_qty;

                if remaining_quantity == Decimal::ZERO {
                } else if opposite_order.quantity == Decimal::ZERO {
                    if let Some(completed_order) = level.orders.pop_front() {}
                }
            }

            if level.quantity == Decimal::ZERO {
                removable_keys.push(*key);
            }
        }

        for key in removable_keys {
            book.remove(&key);
        }

        (trades, remaining_quantity)
    }

    fn price_to_key(&self, price: OrderPrice) -> i64 {
        todo!()
    }
}
