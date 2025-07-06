-- +goose Up
-- +goose StatementBegin
CREATE TABLE segment_rules (
    id CHAR(36) NOT NULL PRIMARY KEY,
    seg_id CHAR(36) NOT NULL,
    attribute VARCHAR(255) NOT NULL,
    operator VARCHAR(255) NOT NULL,
    value VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    INDEX idx_deleted_at (deleted_at),
    CONSTRAINT fk_segment_rules_segment
        FOREIGN KEY (seg_id)
        REFERENCES segments(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE segment_rules;
-- +goose StatementEnd
