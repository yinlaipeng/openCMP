ALTER TABLE projects ADD COLUMN manager_id INT UNSIGNED DEFAULT NULL;
ALTER TABLE projects ADD INDEX idx_projects_manager_id (manager_id);