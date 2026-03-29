package service

import (
	"context"

	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
)

// DomainService 域服务
type DomainService struct {
	db *gorm.DB
}

// NewDomainService 创建域服务
func NewDomainService(db *gorm.DB) *DomainService {
	return &DomainService{db: db}
}

// CreateDomain 创建域
func (s *DomainService) CreateDomain(ctx context.Context, domain *model.Domain) error {
	return s.db.WithContext(ctx).Create(domain).Error
}

// GetDomain 获取域
func (s *DomainService) GetDomain(ctx context.Context, id uint) (*model.Domain, error) {
	var domain model.Domain
	err := s.db.WithContext(ctx).First(&domain, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &domain, nil
}

// ListDomains 列出域
func (s *DomainService) ListDomains(ctx context.Context, keyword string, enabled *bool, limit, offset int) ([]*model.Domain, int64, error) {
	var domains []*model.Domain
	var total int64

	query := s.db.Model(&model.Domain{})

	// 关键词搜索（域名称、描述）
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
		Find(&domains).Error

	return domains, total, err
}

// UpdateDomain 更新域
func (s *DomainService) UpdateDomain(ctx context.Context, domain *model.Domain) error {
	return s.db.WithContext(ctx).Save(domain).Error
}

// DeleteDomain 删除域
func (s *DomainService) DeleteDomain(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Delete(&model.Domain{}, id).Error
}

// EnableDomain 启用域
func (s *DomainService) EnableDomain(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Model(&model.Domain{}).Where("id = ?", id).Update("enabled", true).Error
}

// DisableDomain 禁用域
func (s *DomainService) DisableDomain(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Model(&model.Domain{}).Where("id = ?", id).Update("enabled", false).Error
}

// GetDomainUsers 获取域的用户列表
func (s *DomainService) GetDomainUsers(ctx context.Context, domainID uint, limit, offset int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	if err := s.db.WithContext(ctx).
		Model(&model.User{}).
		Where("domain_id = ?", domainID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := s.db.WithContext(ctx).
		Model(&model.User{}).
		Where("domain_id = ?", domainID).
		Limit(limit).
		Offset(offset).
		Find(&users).Error

	return users, total, err
}

// GetDomainGroups 获取域的用户组列表
func (s *DomainService) GetDomainGroups(ctx context.Context, domainID uint, limit, offset int) ([]*model.Group, int64, error) {
	var groups []*model.Group
	var total int64

	if err := s.db.WithContext(ctx).
		Model(&model.Group{}).
		Where("domain_id = ?", domainID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := s.db.WithContext(ctx).
		Model(&model.Group{}).
		Where("domain_id = ?", domainID).
		Limit(limit).
		Offset(offset).
		Find(&groups).Error

	return groups, total, err
}

// GetDomainProjects 获取域的项目列表
func (s *DomainService) GetDomainProjects(ctx context.Context, domainID uint, limit, offset int) ([]*model.Project, int64, error) {
	var projects []*model.Project
	var total int64

	if err := s.db.WithContext(ctx).
		Model(&model.Project{}).
		Where("domain_id = ?", domainID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := s.db.WithContext(ctx).
		Model(&model.Project{}).
		Where("domain_id = ?", domainID).
		Limit(limit).
		Offset(offset).
		Find(&projects).Error

	return projects, total, err
}

// GetDomainRoles 获取域的角色列表
func (s *DomainService) GetDomainRoles(ctx context.Context, domainID uint, limit, offset int) ([]*model.Role, int64, error) {
	var roles []*model.Role
	var total int64

	query := s.db.Model(&model.Role{}).Where("domain_id = ? OR domain_id IS NULL", domainID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := s.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Find(&roles).Error

	return roles, total, err
}
