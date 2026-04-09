-- SQL script to initialize built-in default data during application startup
-- This script creates essential system entities like default domain, admin user,
-- system project, and built-in authentication source

-- Insert default domain with ID 1
INSERT IGNORE INTO domains (id, name, description, enabled, created_at, updated_at)
VALUES (1, 'Default', 'Default system domain', true, NOW(), NOW());

-- Insert admin role with ID 1
INSERT IGNORE INTO roles (id, name, description, domain_id, created_at, updated_at)
VALUES (1, 'admin', 'Administrator role with full permissions', 1, NOW(), NOW());

-- Insert system role with ID 2
INSERT IGNORE INTO roles (id, name, description, domain_id, created_at, updated_at)
VALUES (2, 'system', 'System role for automated processes', 1, NOW(), NOW());

-- Insert admin user with ID 1
INSERT IGNORE INTO users (id, name, display_name, password, email, phone, domain_id, enabled, created_at, updated_at)
VALUES (1, 'admin', 'Administrator', '$2a$10$8K1TKxmOH6Q/QOA.fZt2yelEShHJXQCBzLr38lczfTnpkDqJF9vOu', 'admin@example.com', '', 1, true, NOW(), NOW()); -- password: admin123

-- Insert built-in authentication source
INSERT IGNORE INTO auth_sources (id, name, type, config, domain_id, enabled, created_at, updated_at)
VALUES (1, '系统认证', 'local', '{}', 1, true, NOW(), NOW());

-- Insert system project
INSERT IGNORE INTO projects (id, name, description, domain_id, enabled, created_at, updated_at)
VALUES (1, 'system', 'Default system project', 1, true, NOW(), NOW());

-- Insert system project
INSERT IGNORE INTO projects (id, name, description, domain_id, enabled, created_at, updated_at)
VALUES (1, 'system', 'Default system project', 1, true, NOW(), NOW());

-- Associate admin user with admin role in default domain
INSERT IGNORE INTO user_roles (user_id, role_id, domain_id, created_at)
VALUES (1, 1, 1, NOW());

-- Associate admin user with system project
INSERT IGNORE INTO project_user_roles (project_id, user_id, role_id, created_at)
VALUES (1, 1, 1, NOW());

-- Insert default message types
INSERT IGNORE INTO message_types (id, name, code, description, enabled, created_at, updated_at)
VALUES
(1, '系统消息', 'system', '系统自动生成的消息', true, NOW(), NOW()),
(2, '安全告警', 'alert', '安全相关的告警消息', true, NOW(), NOW()),
(3, '操作通知', 'operation', '用户操作结果的通知', true, NOW(), NOW()),
(4, '审批通知', 'approval', '工作流审批相关的通知', true, NOW(), NOW());

-- Insert default notification channels
INSERT IGNORE INTO notification_channels (id, name, type, config, enabled, created_at, updated_at)
VALUES
(1, '系统内置邮件', 'email', '{"smtp_server":"localhost","port":587,"username":"","password":""}', true, NOW(), NOW()),
(2, '系统内置企业微信', 'wechat_work', '{"corp_id":"","secret":"","agent_id":""}', true, NOW(), NOW()),
(3, '系统内置钉钉', 'dingtalk', '{"webhook_url":"","secret":""}', true, NOW(), NOW());