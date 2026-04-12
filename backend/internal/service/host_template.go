package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
	"github.com/opencmp/opencmp/pkg/cloudprovider"
)

type HostTemplateService struct {
	DB *gorm.DB
}

func NewHostTemplateService(db *gorm.DB) *HostTemplateService {
	return &HostTemplateService{DB: db}
}

// CreateHostTemplate 创建主机模版
func (s *HostTemplateService) CreateHostTemplate(ctx context.Context, req *cloudprovider.HostTemplate) (*model.HostTemplate, error) {
	hostTemplate := &model.HostTemplate{
		ID: uuid.NewString(),
	}
	hostTemplate.ConvertFromCloudProvider(req)
	hostTemplate.CreatedAt = time.Now()
	hostTemplate.UpdatedAt = time.Now()
	hostTemplate.UpdatedBy = "" // 获取当前用户ID，需要从context中获取

	if err := s.DB.WithContext(ctx).Create(hostTemplate).Error; err != nil {
		return nil, fmt.Errorf("failed to create host template: %w", err)
	}

	return hostTemplate, nil
}

// GetHostTemplateByID 根据ID获取主机模版
func (s *HostTemplateService) GetHostTemplateByID(ctx context.Context, id string) (*model.HostTemplate, error) {
	var hostTemplate model.HostTemplate
	if err := s.DB.WithContext(ctx).Where("id = ?", id).First(&hostTemplate).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("host template with id %s not found: %w", id, ErrNotFound)
		}
		return nil, fmt.Errorf("failed to get host template: %w", err)
	}

	return &hostTemplate, nil
}

// ListHostTemplates 获取主机模版列表
func (s *HostTemplateService) ListHostTemplates(ctx context.Context, projectID string, page, pageSize int) ([]*model.HostTemplate, int64, error) {
	var (
		hostTemplates []*model.HostTemplate
		total         int64
	)

	query := s.DB.WithContext(ctx)
	if projectID != "" {
		query = query.Where("project_id = ?", projectID)
	}

	if err := query.Model(&model.HostTemplate{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count host templates: %w", err)
	}

	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}

	if err := query.Find(&hostTemplates).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list host templates: %w", err)
	}

	return hostTemplates, total, nil
}

// UpdateHostTemplate 更新主机模版
func (s *HostTemplateService) UpdateHostTemplate(ctx context.Context, id string, req *cloudprovider.HostTemplate) (*model.HostTemplate, error) {
	var hostTemplate model.HostTemplate
	if err := s.DB.WithContext(ctx).Where("id = ?", id).First(&hostTemplate).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("host template with id %s not found: %w", id, ErrNotFound)
		}
		return nil, fmt.Errorf("failed to get host template for update: %w", err)
	}

	// Update the fields from the request
	hostTemplate.ConvertFromCloudProvider(req)
	hostTemplate.UpdatedAt = time.Now()
	hostTemplate.UpdatedBy = "" // 获取当前用户ID，需要从context中获取

	if err := s.DB.WithContext(ctx).Model(&hostTemplate).Updates(map[string]interface{}{
		"name":           hostTemplate.Name,
		"description":    hostTemplate.Description,
		"status":         hostTemplate.Status,
		"instance_type":  hostTemplate.InstanceType,
		"cpu_arch":       hostTemplate.CPUArch,
		"memory_size":    hostTemplate.MemorySize,
		"cpu_count":      hostTemplate.CPUCount,
		"disk_size":      hostTemplate.DiskSize,
		"image_id":       hostTemplate.ImageID,
		"os_name":        hostTemplate.OSName,
		"os_version":     hostTemplate.OSVersion,
		"vpc_id":         hostTemplate.VPCID,
		"subnet_id":      hostTemplate.SubnetID,
		"billing_method": hostTemplate.BillingMethod,
		"platform":       hostTemplate.Platform,
		"project_id":     hostTemplate.ProjectID,
		"region_id":      hostTemplate.RegionID,
		"zone_id":        hostTemplate.ZoneID,
		"tags":           hostTemplate.Tags,
		"updated_at":     hostTemplate.UpdatedAt,
		"updated_by":     hostTemplate.UpdatedBy,
	}).Error; err != nil {
		return nil, fmt.Errorf("failed to update host template: %w", err)
	}

	// 重新查询以获取更新后的数据
	if err := s.DB.WithContext(ctx).Where("id = ?", id).First(&hostTemplate).Error; err != nil {
		return nil, fmt.Errorf("failed to get updated host template: %w", err)
	}

	return &hostTemplate, nil
}

// DeleteHostTemplate 删除主机模版
func (s *HostTemplateService) DeleteHostTemplate(ctx context.Context, id string) error {
	result := s.DB.WithContext(ctx).Where("id = ?", id).Delete(&model.HostTemplate{})
	if err := result.Error; err != nil {
		return fmt.Errorf("failed to delete host template: %w", err)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("host template with id %s not found: %w", id, ErrNotFound)
	}

	return nil
}

// ValidateHostTemplateConfig 验证主机模版配置
func (s *HostTemplateService) ValidateHostTemplateConfig(ctx context.Context, template *cloudprovider.HostTemplate) error {
	if template.Name == "" {
		return errors.New("template name is required")
	}

	if template.InstanceType == "" {
		return errors.New("instance type is required")
	}

	if template.ImageID == "" {
		return errors.New("image ID is required")
	}

	if template.ProjectID == "" {
		return errors.New("project ID is required")
	}

	if template.Platform == "" {
		return errors.New("platform is required")
	}

	// 检查模板名称是否已存在
	var count int64
	if err := s.DB.WithContext(ctx).Model(&model.HostTemplate{}).Where("name = ? AND id != ?", template.Name, template.ID).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to validate template name: %w", err)
	}
	if count > 0 {
		return errors.New("template name already exists")
	}

	return nil
}

// Global variables for errors
var (
	ErrNotFound = errors.New("record not found")
)