-- SQL script to delete all built-in system permissions
-- These are typically permissions with type='system' which are the built-in permissions

-- First, let's see what permissions we have
SELECT id, name, display_name, resource, action, type
FROM permissions
WHERE type = 'system'
ORDER BY id;

-- Delete all system permissions (built-in permissions)
DELETE FROM role_permissions WHERE permission_id IN (
    SELECT id FROM permissions WHERE type = 'system'
);

DELETE FROM permissions WHERE type = 'system';

-- Verify deletion
SELECT COUNT(*) as remaining_system_permissions
FROM permissions
WHERE type = 'system';