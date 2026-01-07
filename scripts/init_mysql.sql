CREATE DATABASE IF NOT EXISTS notification CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE notification;

CREATE TABLE IF NOT EXISTS notification_tasks (
    id VARCHAR(36) PRIMARY KEY,
    target_address VARCHAR(255) NOT NULL,
    content TEXT,
    semantic VARCHAR(50) NOT NULL,
    platform VARCHAR(100) NOT NULL,
    max_retries INT NOT NULL DEFAULT 3,
    deadline DATETIME,
    callback_url VARCHAR(500),
    callback_mq VARCHAR(255),
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    retry_count INT NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
