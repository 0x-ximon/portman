use std::{env, net::SocketAddr};

use async_nats::jetstream::Context;
use dotenv::dotenv;
use questdb::ingress::Sender;

pub struct Config {
    pub addr: SocketAddr,
    pub db_conn: Sender,
    pub jetstream: Context,
}

impl Config {
    pub async fn new() -> Result<Self, Box<dyn std::error::Error>> {
        dotenv()?;

        let db_url = env::var("DB_URL")?;
        let db_conn = Sender::from_conf(db_url)?;

        let port = env::var("PORT").unwrap_or_else(|_| String::from("50051"));
        let host = env::var("HOST").unwrap_or_else(|_| String::from("[::1]"));

        let addr: SocketAddr = format!("{}:{}", host, port).parse()?;

        let nats_url =
            env::var("NATS_URL").unwrap_or_else(|_| String::from("nats://localhost:4222"));
        let nats_client = async_nats::connect(nats_url).await?;
        let jetstream = async_nats::jetstream::new(nats_client);

        Ok(Config {
            addr,
            db_conn,
            jetstream,
        })
    }
}
