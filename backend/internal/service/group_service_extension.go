package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/opencmp/opencmp/internal/model"
	"gorm.io/gorm"
)

// AssignRoleToGroup 为组分配角色
func (s *GroupService) AssignRoleToGroup(ctx context.Context, groupID, roleID uint) error {
	var group model.Group
	if err := s.db.WithContext(ctx).First(&group, groupID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("group with ID %d not found", groupID)
		}
		return err
	}

	var role model.Role
	if err := s.db.WithContext(ctx).First(&role, roleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("role with ID %d not found", roleID)
		}
		return err
	}

	// 检查是否已拥有该角色
	var existing model.GroupRole
	err := s.db.WithContext(ctx).Where("group_id = ? AND role_id = ?", groupID, roleID).First(&existing).Error
	if err == nil {
		return nil // already assigned
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	groupRole := &model.GroupRole{
		GroupID:  groupID,
		RoleID:   roleID,
		DomainID: group.DomainID,
	}
	return s.db.WithContext(ctx).Create(groupRole).Error
}

// RevokeRoleFromGroup 从组移除角色
func (s *GroupService) RevokeRoleFromGroup(ctx context.Context, groupID, roleID uint) error {
	var groupRole model.GroupRole
	if err := s.db.WithContext(ctx).Where("group_id = ? AND role_id = ?", groupID, roleID).First(&groupRole).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("group does not have the specified role")
		}
		return err
	}
	return s.db.WithContext(ctx).Delete(&groupRole).Error
}

// GetGroupRoles 获取组的所有角色
func (s *GroupService) GetGroupRoles(ctx context.Context, groupID uint) ([]*model.Role, error) {
	var roles []*model.Role
	err := s.db.WithContext(ctx).
		Joins("JOIN group_roles ON roles.id = group_roles.role_id").
		Where("group_roles.group_id = ?", groupID).
		Find(&roles).Error
	return roles, err
}
