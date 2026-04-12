package service

import (
	"context"
	"fmt"

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

// isDefaultDomain 检查是否为默认域
func (s *DomainService) isDefaultDomain(domain *model.Domain) bool {
	return domain.Name == "Default" || domain.Name == "System" || domain.Name == "default" || domain.Name == "system"
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

// DomainWithAuthSourceCount 带认证源数量的域信息
type DomainWithAuthSourceCount struct {
	model.Domain
	AuthSourceCount int `json:"auth_source_count"`
}

// ListDomains 列出域
func (s *DomainService) ListDomains(ctx context.Context, keyword string, enabled *bool, limit, offset int) ([]*DomainWithAuthSourceCount, int64, error) {
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

	if err != nil {
		return nil, 0, err
	}

	// Convert to enhanced format with auth source count
	var enhancedDomains []*DomainWithAuthSourceCount
	for _, domain := range domains {
		count, _ := s.getAuthSourceCountForDomain(ctx, domain.ID)
		enhancedDomain := &DomainWithAuthSourceCount{
			Domain:          *domain,
			AuthSourceCount: count,
		}
		enhancedDomains = append(enhancedDomains, enhancedDomain)
	}

	return enhancedDomains, total, nil
}

// getAuthSourceCountForDomain 获取指定域的认证源数量
func (s *DomainService) getAuthSourceCountForDomain(ctx context.Context, domainID uint) (int, error) {
	var count int64
	err := s.db.Model(&model.AuthSource{}).
		Where("domain_id = ?", domainID).
		Count(&count).Error

	return int(count), err
}

// UpdateDomain 更新域
func (s *DomainService) UpdateDomain(ctx context.Context, domain *model.Domain) error {
	if s.isDefaultDomain(domain) {
		return fmt.Errorf("default domain cannot be updated")
	}
	return s.db.WithContext(ctx).Save(domain).Error
}

// DeleteDomain 删除域
func (s *DomainService) DeleteDomain(ctx context.Context, id uint) error {
	domain, err := s.GetDomain(ctx, id)
	if err != nil {
		return err
	}
	if domain == nil {
		return gorm.ErrRecordNotFound
	}
	if s.isDefaultDomain(domain) {
		return fmt.Errorf("default domain cannot be deleted")
	}
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

// GetDomainCloudAccounts 获取域的云账号列表
func (s *DomainService) GetDomainCloudAccounts(ctx context.Context, domainID uint, limit, offset int) ([]*model.CloudAccount, int64, error) {
	// 当前云账号表没有直接关联到域的字段
	// 可能通过项目间接关联，或需要扩展云账号模型增加域字段
	// 暂时返回空列表，可根据实际需求调整
	var cloudAccounts []*model.CloudAccount
	var total int64
	return cloudAccounts, total, nil
}

// GetDomainOperationLogs 获取域的操作日志列表
func (s *DomainService) GetDomainOperationLogs(ctx context.Context, domainID uint, limit, offset int) ([]*model.OperationLog, int64, error) {
	var logs []*model.OperationLog
	var total int64

	query := s.db.Model(&model.OperationLog{}).Where("domain_id = ?", domainID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// ListDomainsWithFilters 列出域（支持筛选）
func (s *DomainService) ListDomainsWithFilters(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*DomainWithAuthSourceCount, int64, error) {
	var domains []*model.Domain
	var total int64

	query := s.db.Model(&model.Domain{})

	// 关键词搜索（域名称、描述）
	if keyword, ok := filters["keyword"].(string); ok && keyword != "" {
		keyword = "%" + keyword + "%"
		query = query.Where("name LIKE ? OR description LIKE ?", keyword, keyword)
	}

	// 按状态筛选
	if enabled, ok := filters["enabled"].(bool); ok {
		query = query.Where("enabled = ?", enabled)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := s.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&domains).Error

	if err != nil {
		return nil, 0, err
	}

	// Convert to enhanced format with auth source count
	var enhancedDomains []*DomainWithAuthSourceCount
	for _, domain := range domains {
		count, _ := s.getAuthSourceCountForDomain(ctx, domain.ID)
		enhancedDomain := &DomainWithAuthSourceCount{
			Domain:          *domain,
			AuthSourceCount: count,
		}
		enhancedDomains = append(enhancedDomains, enhancedDomain)
	}

	return enhancedDomains, total, nil
}
