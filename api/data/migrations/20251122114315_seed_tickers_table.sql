-- +goose Up
-- +goose StatementBegin
INSERT INTO
    tickers (symbol, base, quote, status)
VALUES
    ('BTC/USDT', 'BTC', 'USDT', 'OPEN'),
    ('ETH/USDT', 'ETH', 'USDT', 'OPEN'),
    ('SOL/USDT', 'SOL', 'USDT', 'CLOSED');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM tickers WHERE symbol IN ('BTC', 'ETH', 'SOL');
-- +goose StatementEnd
