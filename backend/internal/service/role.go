package service

import (
	"context"

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
func (s *RoleService) ListRoles(ctx context.Context, domainID *uint, limit, offset int) ([]*model.Role, int64, error) {
	var roles []*model.Role
	var total int64

	query := s.db.Model(&model.Role{})
	if domainID != nil {
		query = query.Where("domain_id = ? OR domain_id IS NULL", *domainID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.WithContext(ctx).
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

// ListPermissions 列出权限
func (s *RoleService) ListPermissions(ctx context.Context, limit, offset int) ([]*model.Permission, int64, error) {
	var permissions []*model.Permission
	var total int64

	if err := s.db.Model(&model.Permission{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := s.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("resource, action").
		Find(&permissions).Error

	return permissions, total, err
}

// AssignPermissionToRole 分配权限给角色
func (s *RoleService) AssignPermissionToRole(ctx context.Context, roleID, permissionID uint) error {
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
