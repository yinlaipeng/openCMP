package service

import (
	"context"

	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
)

// GroupService 用户组服务
type GroupService struct {
	db *gorm.DB
}

// NewGroupService 创建用户组服务
func NewGroupService(db *gorm.DB) *GroupService {
	return &GroupService{db: db}
}

// CreateGroup 创建用户组
func (s *GroupService) CreateGroup(ctx context.Context, group *model.Group) error {
	return s.db.WithContext(ctx).Create(group).Error
}

// GetGroup 获取用户组
func (s *GroupService) GetGroup(ctx context.Context, id uint) (*model.Group, error) {
	var group model.Group
	err := s.db.WithContext(ctx).First(&group, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &group, nil
}

// ListGroups 列出用户组
func (s *GroupService) ListGroups(ctx context.Context, domainID *uint, limit, offset int) ([]*model.Group, int64, error) {
	var groups []*model.Group
	var total int64

	query := s.db.Model(&model.Group{})
	if domainID != nil {
		query = query.Where("domain_id = ?", *domainID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := s.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&groups).Error

	return groups, total, err
}

// DeleteGroup 删除用户组
func (s *GroupService) DeleteGroup(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Delete(&model.Group{}, id).Error
}

// UpdateGroup 更新用户组
func (s *GroupService) UpdateGroup(ctx context.Context, group *model.Group) error {
	return s.db.WithContext(ctx).Save(group).Error
}

// GetGroupUsers 获取用户组的用户列表
func (s *GroupService) GetGroupUsers(ctx context.Context, groupID uint, limit, offset int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	if err := s.db.WithContext(ctx).
		Model(&model.User{}).
		Joins("JOIN user_groups ON user_groups.user_id = users.id").
		Where("user_groups.group_id = ?", groupID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := s.db.WithContext(ctx).
		Model(&model.User{}).
		Joins("JOIN user_groups ON user_groups.user_id = users.id").
		Where("user_groups.group_id = ?", groupID).
		Limit(limit).
		Offset(offset).
		Find(&users).Error

	return users, total, err
}

// AddUserToGroup 添加用户到用户组
func (s *GroupService) AddUserToGroup(ctx context.Context, userID, groupID uint) error {
	// 检查是否已存在
	var count int64
	if err := s.db.Model(&model.UserGroup{}).Where("user_id = ? AND group_id = ?", userID, groupID).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return gorm.ErrDuplicatedKey
	}

	ug := &model.UserGroup{
		UserID:  userID,
		GroupID: groupID,
	}
	return s.db.WithContext(ctx).Create(ug).Error
}

// RemoveUserFromGroup 从用户组移除用户
func (s *GroupService) RemoveUserFromGroup(ctx context.Context, userID, groupID uint) error {
	return s.db.WithContext(ctx).Where("user_id = ? AND group_id = ?", userID, groupID).Delete(&model.UserGroup{}).Error
}
