package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
	"github.com/opencmp/opencmp/pkg/cloudprovider"
)

type AutoscalingGroupService struct {
	DB *gorm.DB
}

func NewAutoscalingGroupService(db *gorm.DB) *AutoscalingGroupService {
	return &AutoscalingGroupService{DB: db}
}

// CreateAutoscalingGroup 创建弹性伸缩组
func (s *AutoscalingGroupService) CreateAutoscalingGroup(ctx context.Context, req *cloudprovider.AutoscalingGroup) (*model.AutoscalingGroup, error) {
	// 验证输入参数
	if err := s.ValidateAutoscalingGroupConfig(ctx, req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	autoscalingGroup := &model.AutoscalingGroup{
		ID:              uuid.NewString(),
		Name:            req.Name,
		Description:     req.Description,
		Status:          string(req.Status),
		HostTemplateID:  req.HostTemplateID,
		CurrentCapacity: req.CurrentCapacity,
		DesiredCapacity: req.DesiredCapacity,
		MinSize:         req.MinSize,
		MaxSize:         req.MaxSize,
		Platform:        req.Platform,
		ProjectID:       req.ProjectID,
		RegionID:        req.RegionID,
		ZoneID:          req.ZoneID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		UpdatedBy:       "", // 获取当前用户ID，需要从context中获取
	}

	// 将标签转换为JSON字符串
	if req.Tags != nil {
		tagsBytes, err := json.Marshal(req.Tags)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal tags: %w", err)
		}
		autoscalingGroup.Tags = string(tagsBytes)
	}

	if err := s.DB.WithContext(ctx).Create(autoscalingGroup).Error; err != nil {
		return nil, fmt.Errorf("failed to create autoscaling group: %w", err)
	}

	return autoscalingGroup, nil
}

// GetAutoscalingGroupByID 根据ID获取弹性伸缩组
func (s *AutoscalingGroupService) GetAutoscalingGroupByID(ctx context.Context, id string) (*model.AutoscalingGroup, error) {
	var autoscalingGroup model.AutoscalingGroup
	if err := s.DB.WithContext(ctx).Where("id = ?", id).First(&autoscalingGroup).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("autoscaling group with id %s not found: %w", id, ErrNotFound)
		}
		return nil, fmt.Errorf("failed to get autoscaling group: %w", err)
	}

	return &autoscalingGroup, nil
}

// ListAutoscalingGroups 获取弹性伸缩组列表
func (s *AutoscalingGroupService) ListAutoscalingGroups(ctx context.Context, projectID string, page, pageSize int) ([]*model.AutoscalingGroup, int64, error) {
	var (
		autoscalingGroups []*model.AutoscalingGroup
		total             int64
	)

	query := s.DB.WithContext(ctx)
	if projectID != "" {
		query = query.Where("project_id = ?", projectID)
	}

	if err := query.Model(&model.AutoscalingGroup{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count autoscaling groups: %w", err)
	}

	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}

	if err := query.Find(&autoscalingGroups).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list autoscaling groups: %w", err)
	}

	return autoscalingGroups, total, nil
}

// UpdateAutoscalingGroup 更新弹性伸缩组
func (s *AutoscalingGroupService) UpdateAutoscalingGroup(ctx context.Context, id string, req *cloudprovider.AutoscalingGroup) (*model.AutoscalingGroup, error) {
	// 验证输入参数
	if err := s.ValidateAutoscalingGroupConfig(ctx, req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	var autoscalingGroup model.AutoscalingGroup
	if err := s.DB.WithContext(ctx).Where("id = ?", id).First(&autoscalingGroup).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("autoscaling group with id %s not found: %w", id, ErrNotFound)
		}
		return nil, fmt.Errorf("failed to get autoscaling group for update: %w", err)
	}

	updates := map[string]interface{}{
		"name":             req.Name,
		"description":      req.Description,
		"status":           req.Status,
		"host_template_id": req.HostTemplateID,
		"current_capacity": req.CurrentCapacity,
		"desired_capacity": req.DesiredCapacity,
		"min_size":         req.MinSize,
		"max_size":         req.MaxSize,
		"platform":         req.Platform,
		"project_id":       req.ProjectID,
		"region_id":        req.RegionID,
		"zone_id":          req.ZoneID,
		"updated_at":       time.Now(),
		"updated_by":       "", // 获取当前用户ID，需要从context中获取
	}

	// 将标签转换为JSON字符串
	if req.Tags != nil {
		tagsBytes, err := json.Marshal(req.Tags)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal tags: %w", err)
		}
		updates["tags"] = string(tagsBytes)
	}

	if err := s.DB.WithContext(ctx).Model(&autoscalingGroup).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update autoscaling group: %w", err)
	}

	// 重新查询以获取更新后的数据
	if err := s.DB.WithContext(ctx).Where("id = ?", id).First(&autoscalingGroup).Error; err != nil {
		return nil, fmt.Errorf("failed to get updated autoscaling group: %w", err)
	}

	return &autoscalingGroup, nil
}

// DeleteAutoscalingGroup 删除弹性伸缩组
func (s *AutoscalingGroupService) DeleteAutoscalingGroup(ctx context.Context, id string) error {
	result := s.DB.WithContext(ctx).Where("id = ?", id).Delete(&model.AutoscalingGroup{})
	if err := result.Error; err != nil {
		return fmt.Errorf("failed to delete autoscaling group: %w", err)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("autoscaling group with id %s not found: %w", id, ErrNotFound)
	}

	return nil
}

// ValidateAutoscalingGroupConfig 验证弹性伸缩组配置
func (s *AutoscalingGroupService) ValidateAutoscalingGroupConfig(ctx context.Context, group *cloudprovider.AutoscalingGroup) error {
	if group.Name == "" {
		return errors.New("name is required")
	}

	if group.ProjectID == "" {
		return errors.New("project ID is required")
	}

	if group.Platform == "" {
		return errors.New("platform is required")
	}

	if group.HostTemplateID == "" {
		return errors.New("host template ID is required")
	}

	if group.MinSize < 0 {
		return errors.New("min size must be greater than or equal to 0")
	}

	if group.MaxSize <= 0 {
		return errors.New("max size must be greater than 0")
	}

	if group.MinSize > group.MaxSize {
		return errors.New("min size cannot be greater than max size")
	}

	if group.DesiredCapacity < group.MinSize || group.DesiredCapacity > group.MaxSize {
		return errors.New("desired capacity must be between min and max size")
	}

	// 检查伸缩组名称是否已存在
	var count int64
	if err := s.DB.WithContext(ctx).Model(&model.AutoscalingGroup{}).Where("name = ? AND id != ?", group.Name, group.ID).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to validate autoscaling group name: %w", err)
	}
	if count > 0 {
		return errors.New("autoscaling group name already exists")
	}

	return nil
}