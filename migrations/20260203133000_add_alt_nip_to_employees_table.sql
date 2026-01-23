-- +goose Up
-- +goose StatementBegin
ALTER TABLE employees
ADD COLUMN nip_alt VARCHAR(20) UNIQUE NULL AFTER nip;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE employees
DROP COLUMN nip_alt;
-- +goose StatementEnd
