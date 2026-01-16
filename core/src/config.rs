use std::{env, net::SocketAddr};

use dotenv::dotenv;
use questdb::ingress::Sender;

pub struct Config {
    pub addr: SocketAddr,
    pub conn: Sender,
}

impl Config {
    pub fn new() -> Result<Self, Box<dyn std::error::Error>> {
        dotenv()?;

        let db_url = env::var("DB_URL")?;
        let conn = Sender::from_conf(db_url)?;

        let port = env::var("PORT").unwrap_or_else(|_| String::from("50051"));
        let host = env::var("HOST").unwrap_or_else(|_| String::from("[::1]"));

        let addr: SocketAddr = format!("{}:{}", host, port).parse()?;

        Ok(Config { addr, conn })
    }
}
