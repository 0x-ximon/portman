-- +goose Up
-- +goose StatementBegin
INSERT INTO 
    users (first_name, last_name, email_address, wallet_address, password, role)
VALUES
    ('John', 'Carmack', 'john.carmack@email.com', '0x00', '123456##', 'ADMIN'),
    ('Jane', 'Doe', 'jane.doe@email.com', '0x01', '123456##', 'BOT'),
    ('Bob', 'Smith', 'bob.smith@email.com', '0x02', '123456##', 'USER'),
    ('Alice', 'Johnson', 'alice.johnson@email.com', '0x03', '123456##', 'USER');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM users WHERE email_address IN ('john.carmack@email.com', 'jane.doe@email.com', 'bob.smith@email.com');
-- +goose StatementEnd
