package service

import (
	"context"

	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
)

// PermissionService 权限服务
type PermissionService struct {
	db *gorm.DB
}

// NewPermissionService 创建权限服务
func NewPermissionService(db *gorm.DB) *PermissionService {
	return &PermissionService{db: db}
}

// CreatePermission 创建权限
func (s *PermissionService) CreatePermission(ctx context.Context, permission *model.Permission) error {
	return s.db.WithContext(ctx).Create(permission).Error
}

// GetPermission 获取权限
func (s *PermissionService) GetPermission(ctx context.Context, id uint) (*model.Permission, error) {
	var permission model.Permission
	err := s.db.WithContext(ctx).First(&permission, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &permission, nil
}

// GetPermissionByName 根据名称获取权限
func (s *PermissionService) GetPermissionByName(ctx context.Context, name string) (*model.Permission, error) {
	var permission model.Permission
	err := s.db.WithContext(ctx).Where("name = ?", name).First(&permission).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &permission, nil
}

// ListPermissions 列出权限
func (s *PermissionService) ListPermissions(ctx context.Context, domainID *uint, keyword, resource, action, scope string, enabled *bool, limit, offset int) ([]*model.Permission, int64, error) {
	var permissions []*model.Permission
	var total int64

	query := s.db.Model(&model.Permission{})

	// 如果指定了域，则只查询该域下的权限（或系统级别的权限）
	if domainID != nil {
		query = query.Where("domain_id = ? OR domain_id IS NULL", *domainID)
	}

	// 关键词搜索（权限名、描述）
	if keyword != "" {
		keyword = "%" + keyword + "%"
		query = query.Where("name LIKE ? OR description LIKE ?", keyword, keyword)
	}

	// 按资源筛选
	if resource != "" {
		query = query.Where("resource = ?", resource)
	}

	// 按操作筛选
	if action != "" {
		query = query.Where("action = ?", action)
	}

	// 按作用域筛选
	if scope != "" {
		query = query.Where("scope = ?", scope)
	}

	// 按状态筛选
	if enabled != nil {
		query = query.Where("enabled = ?", *enabled)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := s.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&permissions).Error

	return permissions, total, err
}

// UpdatePermission 更新权限
func (s *PermissionService) UpdatePermission(ctx context.Context, permission *model.Permission) error {
	return s.db.WithContext(ctx).Save(permission).Error
}

// DeletePermission 删除权限
func (s *PermissionService) DeletePermission(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Delete(&model.Permission{}, id).Error
}

// EnablePermission 启用权限
func (s *PermissionService) EnablePermission(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Model(&model.Permission{}).Where("id = ?", id).Update("enabled", true).Error
}

// DisablePermission 禁用权限
func (s *PermissionService) DisablePermission(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Model(&model.Permission{}).Where("id = ?", id).Update("enabled", false).Error
}

// GetPermissionsByRole 获取角色的权限列表
func (s *PermissionService) GetPermissionsByRole(ctx context.Context, roleID uint) ([]*model.Permission, error) {
	var permissions []*model.Permission
	err := s.db.WithContext(ctx).
		Table("permissions").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&permissions).Error
	return permissions, err
}

// AssignPermissionToRole 分配权限给角色
func (s *PermissionService) AssignPermissionToRole(ctx context.Context, roleID, permissionID uint) error {
	// 检查是否已关联
	var count int64
	if err := s.db.Model(&model.RolePermission{}).Where("role_id = ? AND permission_id = ?", roleID, permissionID).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return gorm.ErrDuplicatedKey
	}

	rp := &model.RolePermission{
		RoleID:       roleID,
		PermissionID: permissionID,
	}
	return s.db.WithContext(ctx).Create(rp).Error
}

// RevokePermissionFromRole 从角色撤销权限
func (s *PermissionService) RevokePermissionFromRole(ctx context.Context, roleID, permissionID uint) error {
	return s.db.WithContext(ctx).Where("role_id = ? AND permission_id = ?", roleID, permissionID).Delete(&model.RolePermission{}).Error
}

// GetPermissionsForResourceAction 获取指定资源和操作的权限列表
func (s *PermissionService) GetPermissionsForResourceAction(ctx context.Context, resource, action string) ([]*model.Permission, error) {
	var permissions []*model.Permission
	err := s.db.WithContext(ctx).
		Where("resource = ? AND action = ?", resource, action).
		Find(&permissions).Error
	return permissions, err
}

// GetRolePermissions 获取角色的所有权限
func (s *PermissionService) GetRolePermissions(ctx context.Context, roleID uint) ([]*model.Permission, error) {
	var permissions []*model.Permission
	err := s.db.WithContext(ctx).
		Table("permissions").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&permissions).Error
	return permissions, err
}

// GetUserPermissions 获取用户的直接权限（不包括继承的权限）
func (s *PermissionService) GetUserPermissions(ctx context.Context, userID uint) ([]*model.Permission, error) {
	var permissions []*model.Permission
	err := s.db.WithContext(ctx).
		Table("permissions").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ?", userID).
		Group("permissions.id").
		Find(&permissions).Error
	return permissions, err
}

// GetUserPermissionsInDomain 获取用户在特定域中的权限
func (s *PermissionService) GetUserPermissionsInDomain(ctx context.Context, userID, domainID uint) ([]*model.Permission, error) {
	var permissions []*model.Permission
	err := s.db.WithContext(ctx).
		Table("permissions").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ? AND user_roles.domain_id = ?", userID, domainID).
		Group("permissions.id").
		Find(&permissions).Error
	return permissions, err
}

// GetUserPermissionsInProject 获取用户在特定项目中的权限
func (s *PermissionService) GetUserPermissionsInProject(ctx context.Context, userID, projectID uint) ([]*model.Permission, error) {
	var permissions []*model.Permission
	err := s.db.WithContext(ctx).
		Table("permissions").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN project_user_roles ON project_user_roles.role_id = role_permissions.role_id").
		Where("project_user_roles.user_id = ? AND project_user_roles.project_id = ?", userID, projectID).
		Group("permissions.id").
		Find(&permissions).Error
	return permissions, err
}

// GetGroupPermissions 获取用户组的权限
func (s *PermissionService) GetGroupPermissions(ctx context.Context, groupID uint) ([]*model.Permission, error) {
	var permissions []*model.Permission
	err := s.db.WithContext(ctx).
		Table("permissions").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN group_roles ON group_roles.role_id = role_permissions.role_id").
		Where("group_roles.group_id = ?", groupID).
		Group("permissions.id").
		Find(&permissions).Error
	return permissions, err
}
