package service

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
)

// ProjectService 项目服务
type ProjectService struct {
	db *gorm.DB
}

// NewProjectService 创建项目服务
func NewProjectService(db *gorm.DB) *ProjectService {
	return &ProjectService{db: db}
}

// CreateProject 创建项目
func (s *ProjectService) CreateProject(ctx context.Context, project *model.Project) error {
	return s.db.WithContext(ctx).Create(project).Error
}

// GetProject 获取项目
func (s *ProjectService) GetProject(ctx context.Context, id uint) (*model.Project, error) {
	var project model.Project
	err := s.db.WithContext(ctx).First(&project, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &project, nil
}

// ListProjects 列出项目
func (s *ProjectService) ListProjects(ctx context.Context, domainID *uint, keyword string, enabled *bool, limit, offset int) ([]*model.Project, int64, error) {
	var projects []*model.Project
	var total int64

	query := s.db.Model(&model.Project{})

	// 按域筛选
	if domainID != nil {
		query = query.Where("domain_id = ?", *domainID)
	}

	// 关键词搜索（项目名称、描述）
	if keyword != "" {
		keyword = "%" + keyword + "%"
		query = query.Where("name LIKE ? OR description LIKE ?", keyword, keyword)
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
		Find(&projects).Error

	return projects, total, err
}

// UpdateProject 更新项目
func (s *ProjectService) UpdateProject(ctx context.Context, project *model.Project) error {
	return s.db.WithContext(ctx).Save(project).Error
}

// DeleteProject 删除项目
func (s *ProjectService) DeleteProject(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Delete(&model.Project{}, id).Error
}

// EnableProject 启用项目
func (s *ProjectService) EnableProject(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Model(&model.Project{}).Where("id = ?", id).Update("enabled", true).Error
}

// DisableProject 禁用项目
func (s *ProjectService) DisableProject(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Model(&model.Project{}).Where("id = ?", id).Update("enabled", false).Error
}

// GetProjectUsers 获取项目的用户列表
func (s *ProjectService) GetProjectUsers(ctx context.Context, projectID uint, limit, offset int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	if err := s.db.WithContext(ctx).
		Model(&model.User{}).
		Joins("JOIN project_user_roles ON project_user_roles.user_id = users.id").
		Where("project_user_roles.project_id = ?", projectID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := s.db.WithContext(ctx).
		Model(&model.User{}).
		Joins("JOIN project_user_roles ON project_user_roles.user_id = users.id").
		Where("project_user_roles.project_id = ?", projectID).
		Limit(limit).
		Offset(offset).
		Find(&users).Error

	return users, total, err
}

// GetProjectRoles 获取项目的角色列表
func (s *ProjectService) GetProjectRoles(ctx context.Context, projectID uint, limit, offset int) ([]*model.Role, int64, error) {
	var roles []*model.Role
	var total int64

	query := s.db.Model(&model.Role{}).
		Joins("JOIN project_user_roles ON project_user_roles.role_id = roles.id").
		Where("project_user_roles.project_id = ?", projectID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := s.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Find(&roles).Error

	return roles, total, err
}

// AddUserToProject 添加用户到项目
func (s *ProjectService) AddUserToProject(ctx context.Context, projectID, userID, roleID uint) error {
	// 检查是否已存在
	var count int64
	if err := s.db.Model(&model.ProjectUserRole{}).
		Where("project_id = ? AND user_id = ? AND role_id = ?", projectID, userID, roleID).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return gorm.ErrDuplicatedKey
	}

	pur := &model.ProjectUserRole{
		ProjectID: projectID,
		UserID:    userID,
		RoleID:    roleID,
	}
	return s.db.WithContext(ctx).Create(pur).Error
}

// SetProjectManager 设置项目管理员
func (s *ProjectService) SetProjectManager(ctx context.Context, projectID, userID uint) error {
	// 验证项目是否存在
	var project model.Project
	if err := s.db.WithContext(ctx).First(&project, projectID).Error; err != nil {
		return err
	}

	// 验证用户是否存在
	var user model.User
	if err := s.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		return err
	}

	// 检查用户是否属于项目所在的域
	if user.DomainID != project.DomainID {
		return fmt.Errorf("user does not belong to the project's domain")
	}

	// 更新项目管理员
	return s.db.WithContext(ctx).Model(&project).Update("manager_id", &userID).Error
}

// RemoveUserFromProject 从项目移除用户
func (s *ProjectService) RemoveUserFromProject(ctx context.Context, projectID, userID, roleID uint) error {
	query := s.db.WithContext(ctx).Where("project_id = ? AND user_id = ?", projectID, userID)
	if roleID > 0 {
		query = query.Where("role_id = ?", roleID)
	}
	return query.Delete(&model.ProjectUserRole{}).Error
}
