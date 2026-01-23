-- +goose Up
-- +goose StatementBegin
CREATE TABLE guests (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    `table` VARCHAR(50) NULL,
    present_at DATETIME NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE guests;
-- +goose StatementEnd
