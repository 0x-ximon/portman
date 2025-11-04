use std::{
    collections::{BTreeMap, VecDeque},
    sync::{RwLock, RwLockWriteGuard, atomic::AtomicU64},
};

pub type Price = f64;
pub type Quantity = f64;

#[derive(Copy, Clone)]
pub enum Side {
    Buy,
    Sell,
}

#[derive(PartialEq, Copy, Clone)]
pub enum OrderType {
    Market,
    Limit,
}

#[derive(Copy, Clone)]
pub struct Order {
    pub id: u64,
    pub side: Side,
    pub quantity: Quantity,
    pub price: Option<Price>,
    pub order_type: OrderType,
    pub user_id: u64,
}

pub struct Trade {
    pub order: Order,
    pub executed_price: Price,
}

struct PriceLevel {
    pub price: Price,
    pub quantity: Quantity,
    pub orders: VecDeque<Order>,
}

struct OrderBook {
    pub bids: RwLock<BTreeMap<u64, PriceLevel>>,
    pub asks: RwLock<BTreeMap<u64, PriceLevel>>,
    pub order_counter: AtomicU64,
    pub order_precision: u8,
}

impl OrderBook {
    pub fn new(precision: u8) -> Self {
        OrderBook {
            bids: RwLock::new(BTreeMap::new()),
            asks: RwLock::new(BTreeMap::new()),
            order_counter: AtomicU64::new(1),
            order_precision: precision,
        }
    }

    pub fn market_order(&self, order: &Order) -> anyhow::Result<(Vec<Trade>, f64)> {
        if order.order_type != OrderType::Market {
            return Err(anyhow::anyhow!(
                "Order type must be Market for processing market orders."
            ));
        }

        if order.quantity == 0.0 {
            return Err(anyhow::anyhow!("Order quantity cannot be zero."));
        }

        match order.side {
            Side::Buy => {
                let mut asks = self.asks.write().map_err(|e| {
                    anyhow::anyhow!("RwLock poisoned while acquiring asks write lock: {}", e)
                })?;
                let (trades, remaining_quantity) = self.process_market_order(&order, &mut asks);
                Ok((trades, remaining_quantity))
            }
            Side::Sell => {
                let mut bids = self.bids.write().map_err(|e| {
                    anyhow::anyhow!("RwLock poisoned while acquiring bids write lock: {}", e)
                })?;
                let (trades, remaining_quantity) = self.process_market_order(&order, &mut bids);
                Ok((trades, remaining_quantity))
            }
        }
    }

    fn process_market_order(
        &self,
        order: &Order,
        book: &mut RwLockWriteGuard<BTreeMap<u64, PriceLevel>>,
    ) -> (Vec<Trade>, f64) {
        let mut removable_keys: Vec<u64> = Vec::new();
        let mut remaining_quantity = order.quantity;
        let mut trades: Vec<Trade> = Vec::new();

        for (key, level) in book.iter_mut() {
            if remaining_quantity == 0.0 {
                break;
            }

            while let Some(opposite_order) = level.orders.front_mut() {
                if remaining_quantity == 0.0 {
                    break;
                }

                let fill_qty = opposite_order.quantity.min(remaining_quantity);

                remaining_quantity -= fill_qty;
                opposite_order.quantity -= fill_qty;
                level.quantity -= fill_qty;

                if remaining_quantity == 0.0 {
                    trades.push(Trade {
                        // TODO: Use VWAP for executed price
                        executed_price: level.price,
                        order: *order,
                    });
                } else if opposite_order.quantity == 0.0 {
                    if let Some(completed_order) = level.orders.pop_front() {
                        trades.push(Trade {
                            executed_price: level.price,
                            order: completed_order,
                        });
                    }
                }
            }

            if level.quantity == 0.0 {
                removable_keys.push(*key);
            }
        }

        for key in removable_keys {
            book.remove(&key);
        }

        (trades, remaining_quantity)
    }

    fn process_limit_order(&self, mut order: &Order) -> anyhow::Result<()> {
        match order.side {
            Side::Buy => {
                let mut asks = self.asks.write().unwrap();
                todo!()
            }
            Side::Sell => {
                let mut bids = self.bids.write().unwrap();
                todo!()
            }
        }

        //         match order.side {
        //             Side::Buy => {
        //                 let mut bids = self.bids.write().map_err(|e| format!("Lock error: {}", e))?;
        //                 let level = bids.entry(price_key).or_insert_with(|| PriceLevel {
        //                     price: order.price,
        //                     total_quantity: 0,
        //                 });
        //                 level.total_quantity += order.quantity;
        //             }
        //             Side::Sell => {
        //                 let mut asks = self.asks.write().map_err(|e| format!("Lock error: {}", e))?;
        //                 let level = asks.entry(price_key).or_insert_with(|| PriceLevel {
        //                     price: order.price,
        //                     total_quantity: 0,
        //                 });
        //                 level.total_quantity += order.quantity;
        //             }
        //         }
        //         println!("LIMIT ORDER ID {} added to the book at price {}.", order.id, order.price);
        //         Ok(())
    }

    fn price_to_key(self, price: Price) -> u64 {
        (price * 10f64.powi(self.order_precision as i32)).round() as u64
    }

    fn key_to_price(self, key: u64) -> Price {
        key as f64 / 10f64.powi(self.order_precision as i32)
    }

    fn get_next_id(&self) -> u64 {
        use std::sync::atomic::Ordering::Relaxed;
        self.order_counter.fetch_add(1, Relaxed)
    }
}
