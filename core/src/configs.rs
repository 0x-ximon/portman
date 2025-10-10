use std::env;

pub struct DbConfig {
    pub username: String,
    pub password: String,
    pub port: u16,
    pub host: String,
    pub db: String,
}

pub struct Config {
    pub db: DbConfig,
}

impl Config {
    pub fn new() -> Self {

        Self {
            db: DbConfig {
                username: ,
                password: todo!(),
                port: todo!(),
                host: todo!(),
                db: todo!(),
            },
        }
    }
}
