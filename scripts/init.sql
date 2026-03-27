-- 创建数据库
CREATE DATABASE IF NOT EXISTS opencmp DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE opencmp;

-- 云账户表
CREATE TABLE IF NOT EXISTS cloud_accounts (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    provider_type VARCHAR(20) NOT NULL,
    credentials JSON,
    status VARCHAR(20) DEFAULT 'active',
    description VARCHAR(500),
    created_at DATETIME(3),
    updated_at DATETIME(3),
    deleted_at DATETIME(3),
    INDEX idx_provider_type (provider_type),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
