-- +goose Up
-- +goose StatementBegin
CREATE TABLE winners (
    id INT AUTO_INCREMENT PRIMARY KEY,
    employee_id VARCHAR(64) NOT NULL,
    name VARCHAR(150) NOT NULL,
    position VARCHAR(150) NOT NULL,
    branch VARCHAR(150) NOT NULL,
    employment_type VARCHAR(20) NOT NULL,
    prize_type ENUM('door', 'grand') NOT NULL,
    round_id VARCHAR(50) NOT NULL,
    round_label VARCHAR(100) NOT NULL,
    won_at DATETIME NOT NULL,
    UNIQUE KEY uniq_winner (employee_id, round_id, prize_type),
    KEY idx_winners_prize_type (prize_type),
    KEY idx_winners_won_at (won_at)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE winners;
-- +goose StatementEnd
