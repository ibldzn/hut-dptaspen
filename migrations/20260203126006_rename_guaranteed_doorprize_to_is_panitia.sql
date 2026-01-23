-- +goose Up
-- +goose StatementBegin
ALTER TABLE employees RENAME COLUMN guaranteed_doorprize TO is_panitia;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE employees RENAME COLUMN is_panitia TO guaranteed_doorprize;
-- +goose StatementEnd
