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

-- 域表
CREATE TABLE IF NOT EXISTS domains (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description VARCHAR(500),
    enabled BOOLEAN DEFAULT TRUE,
    parent_id BIGINT UNSIGNED,
    created_at DATETIME(3),
    updated_at DATETIME(3),
    deleted_at DATETIME(3)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    display_name VARCHAR(100),
    email VARCHAR(255),
    phone VARCHAR(20),
    password VARCHAR(255) NOT NULL,
    domain_id BIGINT UNSIGNED NOT NULL,
    enabled BOOLEAN DEFAULT TRUE,
    mfa_enabled BOOLEAN DEFAULT FALSE,
    mfa_secret VARCHAR(255),
    last_login_at DATETIME(3),
    last_login_ip VARCHAR(50),
    password_expire DATETIME(3),
    created_at DATETIME(3),
    updated_at DATETIME(3),
    deleted_at DATETIME(3),
    INDEX idx_domain_id (domain_id),
    INDEX idx_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 角色表
CREATE TABLE IF NOT EXISTS roles (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    display_name VARCHAR(100),
    description VARCHAR(500),
    domain_id BIGINT UNSIGNED,
    type VARCHAR(20) DEFAULT 'custom',
    enabled BOOLEAN DEFAULT TRUE,
    created_at DATETIME(3),
    updated_at DATETIME(3),
    deleted_at DATETIME(3)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 权限表
CREATE TABLE IF NOT EXISTS permissions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    display_name VARCHAR(100),
    description VARCHAR(500),
    resource VARCHAR(50),
    action VARCHAR(50),
    type VARCHAR(20) DEFAULT 'custom',
    created_at DATETIME(3),
    updated_at DATETIME(3),
    deleted_at DATETIME(3)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 插入默认数据

-- 默认域
INSERT INTO domains (id, name, description, enabled, created_at, updated_at) 
VALUES (1, 'Default', '默认域', TRUE, NOW(), NOW())
ON DUPLICATE KEY UPDATE name=name;

-- 超级管理员用户 (密码：admin123)
INSERT INTO users (name, display_name, email, password, domain_id, enabled, created_at, updated_at) 
VALUES ('admin', '超级管理员', 'admin@example.com', 'admin123', 1, TRUE, NOW(), NOW())
ON DUPLICATE KEY UPDATE name=name;

-- 系统角色
INSERT INTO roles (name, display_name, description, domain_id, type, enabled, created_at, updated_at) 
VALUES 
    ('admin', '系统管理员', '系统管理员角色', 1, 'system', TRUE, NOW(), NOW()),
    ('user', '普通用户', '普通用户角色', 1, 'system', TRUE, NOW(), NOW())
ON DUPLICATE KEY UPDATE name=name;

-- 默认权限
INSERT INTO permissions (name, display_name, resource, action, type, created_at, updated_at) 
VALUES 
    ('cloud_account:list', '查看云账户', 'cloud_account', 'list', 'system', NOW(), NOW()),
    ('cloud_account:create', '创建云账户', 'cloud_account', 'create', 'system', NOW(), NOW()),
    ('cloud_account:update', '更新云账户', 'cloud_account', 'update', 'system', NOW(), NOW()),
    ('cloud_account:delete', '删除云账户', 'cloud_account', 'delete', 'system', NOW(), NOW()),
    ('vm:list', '查看虚拟机', 'vm', 'list', 'system', NOW(), NOW()),
    ('vm:create', '创建虚拟机', 'vm', 'create', 'system', NOW(), NOW()),
    ('vm:delete', '删除虚拟机', 'vm', 'delete', 'system', NOW(), NOW()),
    ('vm:action', '操作虚拟机', 'vm', 'action', 'system', NOW(), NOW()),
    ('user:list', '查看用户', 'user', 'list', 'system', NOW(), NOW()),
    ('user:create', '创建用户', 'user', 'create', 'system', NOW(), NOW()),
    ('user:update', '更新用户', 'user', 'update', 'system', NOW(), NOW()),
    ('user:delete', '删除用户', 'user', 'delete', 'system', NOW(), NOW())
ON DUPLICATE KEY UPDATE name=name;
