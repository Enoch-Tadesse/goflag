-- +goose Up
-- +goose StatementBegin
CREATE TABLE flag_users (
    id CHAR(36) NOT NULL PRIMARY KEY,
    flag_id CHAR(36) NOT NULL,
    user_id CHAR(36) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    INDEX idx_flag_id (flag_id),
    INDEX idx_user_id (user_id),
    INDEX idx_deleted_at (deleted_at),
    CONSTRAINT fk_flag_users_flag FOREIGN KEY (flag_id) REFERENCES flags(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_flag_usrers_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE flag_users;
-- +goose StatementEnd
