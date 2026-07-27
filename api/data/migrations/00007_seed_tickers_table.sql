-- +goose Up
-- +goose StatementBegin
INSERT INTO tickers (symbol, base_asset, quote_asset, lot_size, tick_size, status)
VALUES
    ('BTCUSD',  (SELECT id FROM assets WHERE symbol = 'BTC'),  (SELECT id FROM assets WHERE symbol = 'USD'), to_precision(8),  to_precision(2), 'OPEN'),
    ('ETHUSD',  (SELECT id FROM assets WHERE symbol = 'ETH'),  (SELECT id FROM assets WHERE symbol = 'USD'), to_precision(18), to_precision(2), 'OPEN'),
    ('TSLAEUR', (SELECT id FROM assets WHERE symbol = 'TSLA'), (SELECT id FROM assets WHERE symbol = 'EUR'), to_precision(2),  to_precision(2), 'OPEN'),
    ('XAUEUR',  (SELECT id FROM assets WHERE symbol = 'XAU'),  (SELECT id FROM assets WHERE symbol = 'EUR'), to_precision(2),  to_precision(2), 'OPEN'),
    ('EURUSD',  (SELECT id FROM assets WHERE symbol = 'EUR'),  (SELECT id FROM assets WHERE symbol = 'USD'), to_precision(2),  to_precision(2), 'OPEN');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM tickers WHERE symbol IN ('BTCUSD', 'ETHUSD', 'TSLAEUR', 'XAUEUR', 'EURUSD');
-- +goose StatementEnd
