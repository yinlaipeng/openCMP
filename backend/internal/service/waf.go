package service

import (
	"context"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
)

// WAFService WAF策略 Service
type WAFService struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewWAFService 创建 WAF Service
func NewWAFService(db *gorm.DB) *WAFService {
	return &WAFService{
		db:     db,
		logger: zap.L(),
	}
}

// WAFFilter WAF筛选条件
type WAFFilter struct {
	Name           string
	Status         string
	Platform       string
	CloudAccountID string
	DomainID       string
}

// List 获取 WAF 实例列表
func (s *WAFService) List(ctx context.Context, filter WAFFilter, page, pageSize int) ([]model.WAFInstance, int64, error) {
	var instances []model.WAFInstance
	var total int64

	query := s.db.Model(&model.WAFInstance{})

	// 应用筛选条件
	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Platform != "" {
		query = query.Where("platform = ?", filter.Platform)
	}
	if filter.CloudAccountID != "" {
		query = query.Where("cloud_account_id = ?", filter.CloudAccountID)
	}
	if filter.DomainID != "" {
		query = query.Where("domain_id = ?", filter.DomainID)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&instances).Error; err != nil {
		return nil, 0, err
	}

	return instances, total, nil
}

// GetByID 根据ID获取 WAF 实例
func (s *WAFService) GetByID(ctx context.Context, id uint) (*model.WAFInstance, error) {
	var instance model.WAFInstance
	if err := s.db.First(&instance, id).Error; err != nil {
		return nil, err
	}
	return &instance, nil
}

// Create 创建 WAF 实例
func (s *WAFService) Create(ctx context.Context, instance *model.WAFInstance) error {
	return s.db.Create(instance).Error
}

// Update 更新 WAF 实例
func (s *WAFService) Update(ctx context.Context, instance *model.WAFInstance) error {
	return s.db.Save(instance).Error
}

// Delete 删除 WAF 实例
func (s *WAFService) Delete(ctx context.Context, id uint) error {
	return s.db.Delete(&model.WAFInstance{}, id).Error
}

// BatchDelete 批量删除 WAF 实例
func (s *WAFService) BatchDelete(ctx context.Context, ids []uint) error {
	return s.db.Delete(&model.WAFInstance{}, ids).Error
}

// SyncStatus 同步状态
func (s *WAFService) SyncStatus(ctx context.Context, id uint) error {
	var instance model.WAFInstance
	if err := s.db.First(&instance, id).Error; err != nil {
		return err
	}

	now := time.Now()
	instance.SyncTime = &now
	instance.Status = "normal"

	return s.db.Save(&instance).Error
}