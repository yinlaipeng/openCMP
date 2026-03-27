package service

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
)

// UserService 用户服务
type UserService struct {
	db *gorm.DB
}

// NewUserService 创建用户服务
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// CreateUser 创建用户
func (s *UserService) CreateUser(ctx context.Context, user *model.User) error {
	return s.db.WithContext(ctx).Create(user).Error
}

// GetUser 获取用户
func (s *UserService) GetUser(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	err := s.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByName 根据名称获取用户
func (s *UserService) GetUserByName(ctx context.Context, name string) (*model.User, error) {
	var user model.User
	err := s.db.WithContext(ctx).Where("name = ?", name).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// ListUsers 列出用户
func (s *UserService) ListUsers(ctx context.Context, domainID *uint, limit, offset int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	query := s.db.Model(&model.User{})
	if domainID != nil {
		query = query.Where("domain_id = ?", *domainID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&users).Error

	return users, total, err
}

// UpdateUser 更新用户
func (s *UserService) UpdateUser(ctx context.Context, user *model.User) error {
	return s.db.WithContext(ctx).Save(user).Error
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Delete(&model.User{}, id).Error
}

// EnableUser 启用用户
func (s *UserService) EnableUser(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Update("enabled", true).Error
}

// DisableUser 禁用用户
func (s *UserService) DisableUser(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Update("enabled", false).Error
}

// UpdatePassword 更新密码
func (s *UserService) UpdatePassword(ctx context.Context, id uint, newPassword string) error {
	return s.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Update("password", newPassword).Error
}

// UpdateLastLogin 更新最后登录信息
func (s *UserService) UpdateLastLogin(ctx context.Context, id uint, ip string) error {
	now := time.Now()
	return s.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"last_login_at":  now,
		"last_login_ip":  ip,
	}).Error
}

// ResetPassword 重置密码
func (s *UserService) ResetPassword(ctx context.Context, id uint) (string, error) {
	// TODO: 生成随机密码并发送给用户
	return "", nil
}
