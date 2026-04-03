package service

import (
	"context"
	"strings"

	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
)

// RoleService 角色服务
type RoleService struct {
	db *gorm.DB
}

// NewRoleService 创建角色服务
func NewRoleService(db *gorm.DB) *RoleService {
	return &RoleService{db: db}
}

// CreateRole 创建角色
func (s *RoleService) CreateRole(ctx context.Context, role *model.Role) error {
	return s.db.WithContext(ctx).Create(role).Error
}

// GetRole 获取角色
func (s *RoleService) GetRole(ctx context.Context, id uint) (*model.Role, error) {
	var role model.Role
	err := s.db.WithContext(ctx).First(&role, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

// ListRoles 列出角色
func (s *RoleService) ListRoles(ctx context.Context, domainID *uint, keyword, roleType string, enabled *bool, limit, offset int) ([]*model.Role, int64, error) {
	var roles []*model.Role
	var total int64

	query := s.db.Model(&model.Role{})
	if domainID != nil {
		query = query.Where("domain_id = ? OR domain_id IS NULL", *domainID)
	}

	// 关键词搜索（角色名、显示名、描述）
	if keyword != "" {
		keyword = "%" + keyword + "%"
		query = query.Where("name LIKE ? OR display_name LIKE ? OR description LIKE ?", keyword, keyword, keyword)
	}

	// 按类型筛选
	if roleType != "" {
		query = query.Where("type = ?", roleType)
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
		Find(&roles).Error

	return roles, total, err
}

// UpdateRole 更新角色
func (s *RoleService) UpdateRole(ctx context.Context, role *model.Role) error {
	return s.db.WithContext(ctx).Save(role).Error
}

// DeleteRole 删除角色
func (s *RoleService) DeleteRole(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Delete(&model.Role{}, id).Error
}

// AssignRoleToUser 分配角色给用户
func (s *RoleService) AssignRoleToUser(ctx context.Context, userID, roleID, domainID uint) error {
	ur := &model.UserRole{
		UserID:   userID,
		RoleID:   roleID,
		DomainID: domainID,
	}
	return s.db.WithContext(ctx).Create(ur).Error
}

// RevokeRoleFromUser 撤销用户角色
func (s *RoleService) RevokeRoleFromUser(ctx context.Context, userID, roleID, domainID uint) error {
	return s.db.WithContext(ctx).Where("user_id = ? AND role_id = ? AND domain_id = ?", userID, roleID, domainID).Delete(&model.UserRole{}).Error
}

// GetUserRoles 获取用户角色
func (s *RoleService) GetUserRoles(ctx context.Context, userID, domainID uint) ([]*model.Role, error) {
	var roles []*model.Role
	err := s.db.WithContext(ctx).
		Table("roles").
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ? AND user_roles.domain_id = ?", userID, domainID).
		Find(&roles).Error
	return roles, err
}

// AssignRoleToGroup 分配角色给用户组
func (s *RoleService) AssignRoleToGroup(ctx context.Context, groupID, roleID, domainID uint) error {
	gr := &model.GroupRole{
		GroupID:  groupID,
		RoleID:   roleID,
		DomainID: domainID,
	}
	return s.db.WithContext(ctx).Create(gr).Error
}

// GetGroupRoles 获取用户组角色
func (s *RoleService) GetGroupRoles(ctx context.Context, groupID, domainID uint) ([]*model.Role, error) {
	var roles []*model.Role
	err := s.db.WithContext(ctx).
		Table("roles").
		Joins("JOIN group_roles ON group_roles.role_id = roles.id").
		Where("group_roles.group_id = ? AND group_roles.domain_id = ?", groupID, domainID).
		Find(&roles).Error
	return roles, err
}

// ListPermissions 列出权限（支持筛选、搜索）
func (s *RoleService) ListPermissions(ctx context.Context, resource, action, permissionType, keyword string, limit, offset int) ([]*model.Permission, int64, error) {
	var permissions []*model.Permission
	var total int64

	query := s.db.WithContext(ctx).Model(&model.Permission{})

	// 按资源类型筛选
	if resource != "" {
		query = query.Where("resource = ?", resource)
	}

	// 按操作类型筛选
	if action != "" {
		query = query.Where("action = ?", action)
	}

	// 按权限类型筛选（system/custom）
	if permissionType != "" {
		query = query.Where("type = ?", permissionType)
	}

	// 关键词搜索（权限标识、显示名称、描述）
	if keyword != "" {
		keyword = "%" + strings.ToLower(keyword) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(display_name) LIKE ? OR LOWER(description) LIKE ?", keyword, keyword, keyword)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Limit(limit).
		Offset(offset).
		Order("resource, action").
		Find(&permissions).Error

	return permissions, total, err
}

// GetPermission 获取权限详情
func (s *RoleService) GetPermission(ctx context.Context, id uint) (*model.Permission, error) {
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

// CreatePermission 创建权限
func (s *RoleService) CreatePermission(ctx context.Context, permission *model.Permission) error {
	// 设置默认类型
	if permission.Type == "" {
		permission.Type = "custom"
	}
	return s.db.WithContext(ctx).Create(permission).Error
}

// UpdatePermission 更新权限
func (s *RoleService) UpdatePermission(ctx context.Context, permission *model.Permission) error {
	// 系统权限不可更新
	if permission.Type == "system" {
		return gorm.ErrInvalidData
	}
	return s.db.WithContext(ctx).Save(permission).Error
}

// DeletePermission 删除权限
func (s *RoleService) DeletePermission(ctx context.Context, id uint) error {
	// 检查是否是系统权限
	var permission model.Permission
	if err := s.db.WithContext(ctx).First(&permission, id).Error; err != nil {
		return err
	}

	if permission.Type == "system" {
		return gorm.ErrInvalidData
	}

	// 检查是否有关联的角色
	var count int64
	if err := s.db.Model(&model.RolePermission{}).Where("permission_id = ?", id).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return gorm.ErrForeignKeyViolated
	}

	return s.db.WithContext(ctx).Delete(&model.Permission{}, id).Error
}

// AssignPermissionToRole 分配权限给角色
func (s *RoleService) AssignPermissionToRole(ctx context.Context, roleID, permissionID uint) error {
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

// RevokePermissionFromRole 撤销角色权限
func (s *RoleService) RevokePermissionFromRole(ctx context.Context, roleID, permissionID uint) error {
	return s.db.WithContext(ctx).Where("role_id = ? AND permission_id = ?", roleID, permissionID).Delete(&model.RolePermission{}).Error
}

// GetRolePermissions 获取角色权限
func (s *RoleService) GetRolePermissions(ctx context.Context, roleID uint) ([]*model.Permission, error) {
	var permissions []*model.Permission
	err := s.db.WithContext(ctx).
		Table("permissions").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&permissions).Error
	return permissions, err
}

// GetPermissionResources 获取所有资源类型（用于筛选）
func (s *RoleService) GetPermissionResources(ctx context.Context) ([]string, error) {
	var resources []string
	err := s.db.WithContext(ctx).
		Model(&model.Permission{}).
		Distinct("resource").
		Order("resource").
		Pluck("resource", &resources).Error
	return resources, err
}

// GetPermissionActions 获取所有操作类型（用于筛选）
func (s *RoleService) GetPermissionActions(ctx context.Context) ([]string, error) {
	var actions []string
	err := s.db.WithContext(ctx).
		Model(&model.Permission{}).
		Distinct("action").
		Order("action").
		Pluck("action", &actions).Error
	return actions, err
}

// EnableRole 启用角色
func (s *RoleService) EnableRole(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Model(&model.Role{}).Where("id = ?", id).Update("enabled", true).Error
}

// DisableRole 禁用角色
func (s *RoleService) DisableRole(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Model(&model.Role{}).Where("id = ?", id).Update("enabled", false).Error
}

// MakeRolePublic 公开角色
func (s *RoleService) MakeRolePublic(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Model(&model.Role{}).Where("id = ?", id).Update("is_public", true).Error
}

// GetRoleUsers 获取角色的用户列表
func (s *RoleService) GetRoleUsers(ctx context.Context, roleID uint, limit, offset int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	if err := s.db.WithContext(ctx).
		Model(&model.User{}).
		Joins("JOIN user_roles ON user_roles.user_id = users.id").
		Where("user_roles.role_id = ?", roleID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := s.db.WithContext(ctx).
		Model(&model.User{}).
		Joins("JOIN user_roles ON user_roles.user_id = users.id").
		Where("user_roles.role_id = ?", roleID).
		Limit(limit).
		Offset(offset).
		Find(&users).Error

	return users, total, err
}

// GetRoleGroups 获取角色的用户组列表
func (s *RoleService) GetRoleGroups(ctx context.Context, roleID uint, limit, offset int) ([]*model.Group, int64, error) {
	var groups []*model.Group
	var total int64

	if err := s.db.WithContext(ctx).
		Model(&model.Group{}).
		Joins("JOIN group_roles ON group_roles.group_id = groups.id").
		Where("group_roles.role_id = ?", roleID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := s.db.WithContext(ctx).
		Model(&model.Group{}).
		Joins("JOIN group_roles ON group_roles.group_id = groups.id").
		Where("group_roles.role_id = ?", roleID).
		Limit(limit).
		Offset(offset).
		Find(&groups).Error

	return groups, total, err
}

// GetRolePolicies 获取角色的策略列表
func (s *RoleService) GetRolePolicies(ctx context.Context, roleID uint) ([]*model.Policy, error) {
	var policies []*model.Policy
	err := s.db.WithContext(ctx).
		Table("policies").
		Joins("JOIN role_policies ON role_policies.policy_id = policies.id").
		Where("role_policies.role_id = ?", roleID).
		Find(&policies).Error
	return policies, err
}

// AssignPolicyToRole 分配策略给角色
func (s *RoleService) AssignPolicyToRole(ctx context.Context, roleID uint, policyID string) error {
	// 检查是否已关联
	var count int64
	if err := s.db.Model(&model.RolePolicy{}).Where("role_id = ? AND policy_id = ?", roleID, policyID).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return gorm.ErrDuplicatedKey
	}

	rp := &model.RolePolicy{
		RoleID:   roleID,
		PolicyID: policyID,
	}
	return s.db.WithContext(ctx).Create(rp).Error
}

// RevokePolicyFromRole 从角色撤销策略
func (s *RoleService) RevokePolicyFromRole(ctx context.Context, roleID uint, policyID string) error {
	return s.db.WithContext(ctx).Where("role_id = ? AND policy_id = ?", roleID, policyID).Delete(&model.RolePolicy{}).Error
}

// CheckUserPermission 检查用户是否有指定权限
func (s *RoleService) CheckUserPermission(ctx context.Context, userID uint, resource, action string) (bool, error) {
	// 检查用户直接拥有的角色权限
	var hasPermission bool
	
	// 查询用户在域级别拥有的角色及其权限
	err := s.db.WithContext(ctx).
		Raw(`
			SELECT COUNT(*) > 0 
			FROM permissions p
			JOIN role_permissions rp ON p.id = rp.permission_id
			JOIN user_roles ur ON rp.role_id = ur.role_id
			WHERE ur.user_id = ? AND p.resource = ? AND p.action = ?
		`, userID, resource, action).
		Scan(&hasPermission).Error
	
	if err != nil {
		return false, err
	}
	
	if hasPermission {
		return true, nil
	}
	
	// 查询用户在项目级别拥有的角色及其权限
	err = s.db.WithContext(ctx).
		Raw(`
			SELECT COUNT(*) > 0 
			FROM permissions p
			JOIN role_permissions rp ON p.id = rp.permission_id
			JOIN project_user_roles pur ON rp.role_id = pur.role_id
			WHERE pur.user_id = ? AND p.resource = ? AND p.action = ?
		`, userID, resource, action).
		Scan(&hasPermission).Error
	
	if err != nil {
		return false, err
	}
	
	if hasPermission {
		return true, nil
	}
	
	// 检查用户所属组拥有的角色权限
	err = s.db.WithContext(ctx).
		Raw(`
			SELECT COUNT(*) > 0 
			FROM permissions p
			JOIN role_permissions rp ON p.id = rp.permission_id
			JOIN group_roles gr ON rp.role_id = gr.role_id
			JOIN user_groups ug ON gr.group_id = ug.group_id
			WHERE ug.user_id = ? AND p.resource = ? AND p.action = ?
		`, userID, resource, action).
		Scan(&hasPermission).Error
	
	if err != nil {
		return false, err
	}
	
	return hasPermission, nil
}
