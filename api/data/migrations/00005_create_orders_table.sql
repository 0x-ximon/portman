-- +goose Up
-- +goose StatementBegin
CREATE TYPE ORDER_SIDE AS ENUM ('BUY', 'SELL');
CREATE TYPE ORDER_MODE AS ENUM ('GTC', 'FOK', 'IOC');
CREATE TYPE ORDER_STATUS AS ENUM ('PENDING', 'REJECTED', 'PARTIAL', 'CANCELLED', 'FULFILLED');

CREATE TABLE orders (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(id) NOT NULL,
    ticker_id BIGINT REFERENCES tickers(id) NOT NULL,

    price NUMERIC NOT NULL,
    quantity NUMERIC NOT NULL,
    order_side ORDER_SIDE NOT NULL,
    order_mode ORDER_MODE NOT NULL,
    order_status ORDER_STATUS NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT validate_price CHECK (price > 0),
    CONSTRAINT validate_quantity CHECK (quantity > 0)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE orders;
DROP TYPE ORDER_STATUS;
DROP TYPE ORDER_MODE;
DROP TYPE ORDER_SIDE;
-- +goose StatementEnd
