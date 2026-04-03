package service

import (
	"context"

	"github.com/opencmp/opencmp/internal/model"
)

// GetUserPermissions 获取用户通过角色获得的所有权限
func (s *UserService) GetUserPermissions(ctx context.Context, userID uint) ([]*model.Permission, error) {
	var permissions []*model.Permission
	err := s.db.WithContext(ctx).
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Joins("JOIN user_roles ON role_permissions.role_id = user_roles.role_id").
		Where("user_roles.user_id = ?", userID).
		Distinct("permissions.*").
		Find(&permissions).Error
	return permissions, err
}
