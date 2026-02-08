-- +goose Up
-- +goose StatementBegin
INSERT INTO users (id, balance) VALUES (1, 1000.00)
ON CONFLICT (id) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM users WHERE id = 1;
-- +goose StatementEnd
