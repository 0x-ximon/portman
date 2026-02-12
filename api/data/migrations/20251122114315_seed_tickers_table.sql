-- +goose Up
-- +goose StatementBegin
INSERT INTO
    tickers (symbol, base, quote, status)
VALUES
    ('BTCUSDT', 'BTC', 'USDT', 'OPEN'),
    ('ETHUSDT', 'ETH', 'USDT', 'OPEN'),
    ('SOLUSDT', 'SOL', 'USDT', 'CLOSED');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM tickers WHERE symbol IN ('BTC', 'ETH', 'SOL');
-- +goose StatementEnd
