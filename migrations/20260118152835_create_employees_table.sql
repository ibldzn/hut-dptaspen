-- +goose Up
-- +goose StatementBegin
CREATE TABLE employees (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    position VARCHAR(100) NOT NULL,
    branch_office VARCHAR(100) NOT NULL,
    employment_type ENUM('Organik', 'TAD') NOT NULL,
    is_excluded BOOLEAN DEFAULT FALSE,
    guaranteed_doorprize BOOLEAN DEFAULT FALSE,
    present_at DATETIME NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE employees;
-- +goose StatementEnd
