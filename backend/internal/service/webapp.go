package service

import (
	"context"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
)

// WebappService 应用程序服务 Service
type WebappService struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewWebappService 创建 Webapp Service
func NewWebappService(db *gorm.DB) *WebappService {
	return &WebappService{
		db:     db,
		logger: zap.L(),
	}
}

// WebappFilter 筛选条件
type WebappFilter struct {
	Name           string
	Status         string
	Platform       string
	CloudAccountID string
	ProjectID      string
	Stack          string
}

// List 获取应用程序服务列表
func (s *WebappService) List(ctx context.Context, filter WebappFilter, page, pageSize int) ([]model.WebappInstance, int64, error) {
	var instances []model.WebappInstance
	var total int64

	query := s.db.Model(&model.WebappInstance{})

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
	if filter.ProjectID != "" {
		query = query.Where("project_id = ?", filter.ProjectID)
	}
	if filter.Stack != "" {
		query = query.Where("stack = ?", filter.Stack)
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

// GetByID 根据ID获取应用程序服务
func (s *WebappService) GetByID(ctx context.Context, id uint) (*model.WebappInstance, error) {
	var instance model.WebappInstance
	if err := s.db.First(&instance, id).Error; err != nil {
		return nil, err
	}
	return &instance, nil
}

// Create 创建应用程序服务
func (s *WebappService) Create(ctx context.Context, instance *model.WebappInstance) error {
	return s.db.Create(instance).Error
}

// Update 更新应用程序服务
func (s *WebappService) Update(ctx context.Context, instance *model.WebappInstance) error {
	return s.db.Save(instance).Error
}

// Delete 删除应用程序服务
func (s *WebappService) Delete(ctx context.Context, id uint) error {
	return s.db.Delete(&model.WebappInstance{}, id).Error
}

// BatchDelete 批量删除
func (s *WebappService) BatchDelete(ctx context.Context, ids []uint) error {
	return s.db.Delete(&model.WebappInstance{}, ids).Error
}

// SyncStatus 同步状态
func (s *WebappService) SyncStatus(ctx context.Context, id uint) error {
	var instance model.WebappInstance
	if err := s.db.First(&instance, id).Error; err != nil {
		return err
	}

	now := time.Now()
	instance.SyncTime = &now
	instance.Status = "normal"

	return s.db.Save(&instance).Error
}