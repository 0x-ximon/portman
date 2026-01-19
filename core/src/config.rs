use std::{env, net};

use async_nats::jetstream;
use dotenv::dotenv;
use questdb::ingress;

pub struct Config {
    pub addr: net::SocketAddr,
    pub db_conn: ingress::Sender,
    pub publisher: jetstream::Context,
}

impl Config {
    pub async fn new() -> Result<Self, Box<dyn std::error::Error>> {
        dotenv()?;

        let db_url = env::var("DB_URL")?;
        let db_conn = ingress::Sender::from_conf(db_url)?;

        let port = env::var("PORT").unwrap_or_else(|_| String::from("50051"));
        let host = env::var("HOST").unwrap_or_else(|_| String::from("[::1]"));

        let addr: net::SocketAddr = format!("{}:{}", host, port).parse()?;

        let nats_url =
            env::var("NATS_URL").unwrap_or_else(|_| String::from("nats://localhost:4222"));
        let nats_client = async_nats::connect(nats_url).await?;
        let publisher = jetstream::new(nats_client);

        Ok(Config {
            addr,
            db_conn,
            publisher,
        })
    }
}
