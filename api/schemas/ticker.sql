CREATE TABLE 
  tickers (
    ID SERIAL PRIMARY KEY,
    symbol TEXT PRIMARY KEY NOT NULL,
    name TEXT NOT NULL,
    last REAL,
    ask REAL,
    bid REAL
  )