-- +goose Up
-- +goose StatementBegin
CREATE TABLE scan_events (
    id INT AUTO_INCREMENT PRIMARY KEY,
    scanner_id TINYINT NOT NULL,
    name VARCHAR(150) NOT NULL,
    scanned_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    KEY idx_scan_events_scanner (scanner_id, scanned_at)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE scan_events;
-- +goose StatementEnd
