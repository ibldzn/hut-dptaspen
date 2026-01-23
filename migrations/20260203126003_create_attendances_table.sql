-- +goose Up
-- +goose StatementBegin
CREATE TABLE attendances (
    id INT AUTO_INCREMENT PRIMARY KEY,
    person_type ENUM('employee', 'guest') NOT NULL,
    person_id INT NULL,
    guest_name VARCHAR(150) NULL,
    present_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uniq_attendance_employee (person_type, person_id),
    UNIQUE KEY uniq_attendance_guest (person_type, guest_name),
    KEY idx_attendances_present_at (present_at)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE attendances;
-- +goose StatementEnd
