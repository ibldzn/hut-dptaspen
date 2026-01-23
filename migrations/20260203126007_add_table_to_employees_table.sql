-- +goose Up
-- +goose StatementBegin
ALTER TABLE employees
ADD COLUMN `table` VARCHAR(50) NULL AFTER nip_alt;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE employees
DROP COLUMN `table`;
-- +goose StatementEnd
