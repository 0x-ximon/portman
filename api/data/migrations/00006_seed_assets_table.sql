-- +goose Up
-- +goose StatementBegin
INSERT INTO assets (name, symbol, decimals, kind)
VALUES
    ('Bitcoin', 'BTC', 8, 'CRYPTO'),
    ('Ethereum', 'ETH', 18, 'CRYPTO'),
    ('Euro', 'EUR', 2, 'FIAT'),
    ('Gold', 'XAU', 2, 'COMMODITY'),
    ('Tesla', 'TSLA', 2, 'STOCK'),
    ('Dollar', 'USD', 2, 'FIAT');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM assets WHERE symbol IN ('BTC', 'ETH', 'EUR', 'XAU', 'TSLA', 'USD');
-- +goose StatementEnd
