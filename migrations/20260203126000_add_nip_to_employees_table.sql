-- +goose Up
-- +goose StatementBegin
-- unique nip
ALTER TABLE employees
ADD COLUMN nip VARCHAR(20) UNIQUE NOT NULL AFTER id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE employees
DROP COLUMN nip;
-- +goose StatementEnd
