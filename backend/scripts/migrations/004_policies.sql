-- 策略表迁移脚本
-- 参考 OneCloud 策略管理设计

-- 创建 policies 表
CREATE TABLE IF NOT EXISTS `policies` (
  `id` VARCHAR(64) PRIMARY KEY COMMENT '策略 ID（UUID）',
  `name` VARCHAR(100) NOT NULL UNIQUE COMMENT '策略名称',
  `description` VARCHAR(500) COMMENT '策略描述',
  `scope` VARCHAR(20) NOT NULL COMMENT '作用域：system/domain/project',
  `domain_id` VARCHAR(64) COMMENT '域 ID',
  `project_id` VARCHAR(64) COMMENT '项目 ID',
  `policy` JSON NOT NULL COMMENT '策略内容（JSON 格式）',
  `is_system` BOOLEAN DEFAULT FALSE COMMENT '是否系统策略',
  `is_public` BOOLEAN DEFAULT FALSE COMMENT '是否公开',
  `is_emulated` BOOLEAN DEFAULT FALSE COMMENT '是否预置策略',
  `enabled` BOOLEAN DEFAULT TRUE COMMENT '是否启用',
  `pending_deleted` BOOLEAN DEFAULT FALSE COMMENT '是否待删除',
  `deleted` BOOLEAN DEFAULT FALSE COMMENT '是否已删除',
  `public_scope` VARCHAR(20) COMMENT '公开范围',
  `update_version` INT DEFAULT 0 COMMENT '更新版本',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` TIMESTAMP NULL,
  INDEX `idx_scope` (`scope`),
  INDEX `idx_domain_id` (`domain_id`),
  INDEX `idx_project_id` (`project_id`),
  INDEX `idx_is_system` (`is_system`),
  INDEX `idx_enabled` (`enabled`),
  INDEX `idx_created_at` (`created_at`),
  INDEX `idx_updated_at` (`updated_at`),
  INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='策略表';

-- 创建角色策略关联表
CREATE TABLE IF NOT EXISTS `role_policies` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `role_id` INT NOT NULL COMMENT '角色 ID',
  `policy_id` VARCHAR(64) NOT NULL COMMENT '策略 ID',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY `uk_role_policy` (`role_id`, `policy_id`),
  INDEX `idx_role_id` (`role_id`),
  INDEX `idx_policy_id` (`policy_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色策略关联表';
