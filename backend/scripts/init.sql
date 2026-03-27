-- 创建数据库 (如果不存在)
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

-- 插入示例数据
INSERT INTO cloud_accounts (name, provider_type, status, description, credentials, created_at, updated_at) 
VALUES 
    ('阿里云测试账号', 'alibaba', 'active', '阿里云测试环境', '{"access_key_id":"test","access_key_secret":"test"}', NOW(), NOW()),
    ('腾讯云测试账号', 'tencent', 'active', '腾讯云测试环境', '{"secret_id":"test","secret_key":"test"}', NOW(), NOW())
ON DUPLICATE KEY UPDATE name=name;
