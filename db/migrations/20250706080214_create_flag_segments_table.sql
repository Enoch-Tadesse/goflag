-- +goose Up
-- +goose StatementBegin
CREATE TABLE flag_segments (
    id CHAR(36) NOT NULL PRIMARY KEY,
    flag_id CHAR(36) NOT NULL,
    seg_id CHAR(36) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    INDEX idx_flag_id (flag_id),
    INDEX idx_seg_id (seg_id),
    INDEX idx_deleted_at (deleted_at),
    CONSTRAINT fk_flag_segments_flag FOREIGN KEY (flag_id) REFERENCES flags(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_flag_segments_segment FOREIGN KEY (seg_id) REFERENCES segments(id) ON DELETE CASCADE ON UPDATE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE flag_segments;
-- +goose StatementEnd
