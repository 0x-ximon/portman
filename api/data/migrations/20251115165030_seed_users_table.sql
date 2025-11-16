-- +goose Up
-- +goose StatementBegin
INSERT INTO 
    users (first_name, last_name, phone_number, email_address, wallet_address, password, role)
VALUES
    ('John', 'Carmack', '1234567890', 'john.carmack@email.com', '0x00', '123456##', 'ADMIN'),
    ('Alice', 'Johnson', '7777777777', 'alice.johnson@email.com', '0x03', '123456##', 'USER'),
    ('Bob', 'Smith', '5555555555', 'bob.smith@email.com', '0x02', '123456##', 'USER'),
    ('Jane', 'Doe', '9876543210', 'jane.doe@email.com', '0x01', '123456##', 'BOT');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM users WHERE email_address IN ('john.carmack@email.com', 'jane.doe@email.com', 'bob.smith@email.com');
-- +goose StatementEnd
