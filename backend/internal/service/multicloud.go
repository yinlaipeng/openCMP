package service

import (
	"context"

	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
)

// CloudAccessGroupService 云访问组服务
type CloudAccessGroupService struct {
	db *gorm.DB
}

// NewCloudAccessGroupService 创建云访问组服务
func NewCloudAccessGroupService(db *gorm.DB) *CloudAccessGroupService {
	return &CloudAccessGroupService{db: db}
}

// CreateCloudAccessGroup 创建云访问组
func (s *CloudAccessGroupService) CreateCloudAccessGroup(ctx context.Context, group *model.CloudAccessGroup) error {
	return s.db.WithContext(ctx).Create(group).Error
}

// GetCloudAccessGroup 获取云访问组
func (s *CloudAccessGroupService) GetCloudAccessGroup(ctx context.Context, id uint) (*model.CloudAccessGroup, error) {
	var group model.CloudAccessGroup
	err := s.db.WithContext(ctx).First(&group, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &group, nil
}

// ListCloudAccessGroups 列出云访问组
func (s *CloudAccessGroupService) ListCloudAccessGroups(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]*model.CloudAccessGroup, int64, error) {
	var groups []*model.CloudAccessGroup
	var total int64

	query := s.db.Model(&model.CloudAccessGroup{})

	if name, ok := filters["name"].(string); ok && name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if platform, ok := filters["platform"].(string); ok && platform != "" {
		query = query.Where("platform = ?", platform)
	}
	if domainID, ok := filters["domain_id"].(uint); ok && domainID > 0 {
		query = query.Where("domain_id = ?", domainID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&groups).Error

	return groups, total, err
}

// UpdateCloudAccessGroup 更新云访问组
func (s *CloudAccessGroupService) UpdateCloudAccessGroup(ctx context.Context, id uint, updates map[string]interface{}) error {
	return s.db.WithContext(ctx).Model(&model.CloudAccessGroup{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteCloudAccessGroup 删除云访问组
func (s *CloudAccessGroupService) DeleteCloudAccessGroup(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Delete(&model.CloudAccessGroup{}, id).Error
}

// BatchDeleteCloudAccessGroups 批量删除云访问组
func (s *CloudAccessGroupService) BatchDeleteCloudAccessGroups(ctx context.Context, ids []uint) error {
	return s.db.WithContext(ctx).Delete(&model.CloudAccessGroup{}, ids).Error
}

// ProxySettingService 代理设置服务
type ProxySettingService struct {
	db *gorm.DB
}

// NewProxySettingService 创建代理设置服务
func NewProxySettingService(db *gorm.DB) *ProxySettingService {
	return &ProxySettingService{db: db}
}

// CreateProxySetting 创建代理设置
func (s *ProxySettingService) CreateProxySetting(ctx context.Context, proxy *model.ProxySetting) error {
	return s.db.WithContext(ctx).Create(proxy).Error
}

// GetProxySetting 获取代理设置
func (s *ProxySettingService) GetProxySetting(ctx context.Context, id uint) (*model.ProxySetting, error) {
	var proxy model.ProxySetting
	err := s.db.WithContext(ctx).First(&proxy, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &proxy, nil
}

// ListProxySettings 列出代理设置
func (s *ProxySettingService) ListProxySettings(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]*model.ProxySetting, int64, error) {
	var proxies []*model.ProxySetting
	var total int64

	query := s.db.Model(&model.ProxySetting{})

	if name, ok := filters["name"].(string); ok && name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if domainID, ok := filters["domain_id"].(uint); ok && domainID > 0 {
		query = query.Where("domain_id = ?", domainID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&proxies).Error

	return proxies, total, err
}

// UpdateProxySetting 更新代理设置
func (s *ProxySettingService) UpdateProxySetting(ctx context.Context, id uint, updates map[string]interface{}) error {
	return s.db.WithContext(ctx).Model(&model.ProxySetting{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteProxySetting 删除代理设置
func (s *ProxySettingService) DeleteProxySetting(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Delete(&model.ProxySetting{}, id).Error
}

// BatchDeleteProxySettings 批量删除代理设置
func (s *ProxySettingService) BatchDeleteProxySettings(ctx context.Context, ids []uint) error {
	return s.db.WithContext(ctx).Delete(&model.ProxySetting{}, ids).Error
}

// SetProxySettingSharing 设置代理共享
func (s *ProxySettingService) SetProxySettingSharing(ctx context.Context, id uint, sharedScope string) error {
	return s.db.WithContext(ctx).Model(&model.ProxySetting{}).Where("id = ?", id).Update("shared_scope", sharedScope).Error
}