-- +migrate Up
CREATE TABLE IF NOT EXISTS `users`(
    id bigint AUTO_INCREMENT NOT NULL,
    auth0_id VARCHAR (255) NOT NULL,
    name VARCHAR (255) NOT NULL,
    email VARCHAR (255),
    zoom_token VARCHAR (1000),
    zoom_refresh_token VARCHAR (1000),
    -- デフォルトはCURRENT_TIMESTAMPが設定されている
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    -- 主キー（重複することはできない）
    PRIMARY KEY (id)
);

-- +migrate Down
DROP TABLE IF EXISTS `users`;