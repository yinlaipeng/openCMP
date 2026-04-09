package service

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
	"github.com/opencmp/opencmp/pkg/utils"
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
	// 对密码进行哈希处理
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	
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
func (s *UserService) ListUsers(ctx context.Context, domainID *uint, keyword, email string, enabled *bool, limit, offset int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	query := s.db.Model(&model.User{})

	// 按域筛选
	if domainID != nil {
		query = query.Where("domain_id = ?", *domainID)
	}

	// 按用户名关键词筛选
	if keyword != "" {
		keyword = "%" + keyword + "%"
		query = query.Where("name LIKE ? OR display_name LIKE ?", keyword, keyword)
	}

	// 按邮箱筛选
	if email != "" {
		query = query.Where("email = ?", email)
	}

	// 按状态筛选
	if enabled != nil {
		query = query.Where("enabled = ?", *enabled)
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
	// 对新密码进行哈希处理
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}
	
	return s.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Update("password", hashedPassword).Error
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

// ResetUserPassword 重置用户密码
func (s *UserService) ResetUserPassword(ctx context.Context, id uint, newPassword string) error {
	// 对新密码进行哈希处理
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}
	
	return s.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Update("password", hashedPassword).Error
}

// GetUserRoles 获取用户角色
func (s *UserService) GetUserRoles(ctx context.Context, userID, domainID uint) ([]*model.Role, error) {
	var roles []*model.Role
	err := s.db.WithContext(ctx).
		Table("roles").
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ? AND user_roles.domain_id = ?", userID, domainID).
		Find(&roles).Error
	return roles, err
}

// GetUserGroups 获取用户组
func (s *UserService) GetUserGroups(ctx context.Context, userID uint) ([]*model.Group, error) {
	var groups []*model.Group
	err := s.db.WithContext(ctx).
		Table("groups").
		Joins("JOIN user_groups ON user_groups.group_id = groups.id").
		Where("user_groups.user_id = ?", userID).
		Find(&groups).Error
	return groups, err
}

// AddUserToGroup 添加用户到用户组
func (s *UserService) AddUserToGroup(ctx context.Context, userID, groupID uint) error {
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
func (s *UserService) RemoveUserFromGroup(ctx context.Context, userID, groupID uint) error {
	return s.db.WithContext(ctx).Where("user_id = ? AND group_id = ?", userID, groupID).Delete(&model.UserGroup{}).Error
}

// AssignUserRole 分配角色给用户
func (s *UserService) AssignUserRole(ctx context.Context, userID, roleID, domainID uint) error {
	// 检查是否已存在
	var count int64
	if err := s.db.Model(&model.UserRole{}).Where("user_id = ? AND role_id = ? AND domain_id = ?", userID, roleID, domainID).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return gorm.ErrDuplicatedKey
	}

	ur := &model.UserRole{
		UserID:   userID,
		RoleID:   roleID,
		DomainID: domainID,
	}
	return s.db.WithContext(ctx).Create(ur).Error
}

// RevokeUserRole 撤销用户角色
func (s *UserService) RevokeUserRole(ctx context.Context, userID, roleID, domainID uint) error {
	return s.db.WithContext(ctx).Where("user_id = ? AND role_id = ? AND domain_id = ?", userID, roleID, domainID).Delete(&model.UserRole{}).Error
}

// GetUserRoleIDs 获取用户的角色ID列表
func (s *UserService) GetUserRoleIDs(ctx context.Context, userID uint) ([]uint, error) {
	var roleIDs []uint
	err := s.db.WithContext(ctx).
		Model(&model.UserRole{}).
		Where("user_id = ?", userID).
		Pluck("role_id", &roleIDs).Error
	if err != nil {
		return nil, err
	}

	// 也获取项目级别的角色
	var projectRoleIDs []uint
	err = s.db.WithContext(ctx).
		Model(&model.ProjectUserRole{}).
		Where("user_id = ?", userID).
		Pluck("role_id", &projectRoleIDs).Error
	if err != nil {
		return nil, err
	}

	// 合并两个切片
	roleIDs = append(roleIDs, projectRoleIDs...)

	return roleIDs, nil
}

// AssignUserToProject 将用户分配到项目
func (s *UserService) AssignUserToProject(ctx context.Context, userID, projectID, roleID uint) error {
	// 检查是否已存在
	var count int64
	if err := s.db.Model(&model.ProjectUserRole{}).Where("user_id = ? AND project_id = ? AND role_id = ?", userID, projectID, roleID).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return gorm.ErrDuplicatedKey
	}

	pur := &model.ProjectUserRole{
		UserID:    userID,
		ProjectID: projectID,
		RoleID:    roleID,
	}
	return s.db.WithContext(ctx).Create(pur).Error
}

// RemoveUserFromProject 将用户从项目中移除
func (s *UserService) RemoveUserFromProject(ctx context.Context, userID, projectID uint) error {
	return s.db.WithContext(ctx).Where("user_id = ? AND project_id = ?", userID, projectID).Delete(&model.ProjectUserRole{}).Error
}

// GetUserProjects 获取用户所属的项目列表
func (s *UserService) GetUserProjects(ctx context.Context, userID uint) ([]*model.Project, error) {
	var projects []*model.Project
	err := s.db.WithContext(ctx).
		Table("projects").
		Joins("JOIN project_user_roles ON project_user_roles.project_id = projects.id").
		Where("project_user_roles.user_id = ?", userID).
		Find(&projects).Error
	return projects, err
}
