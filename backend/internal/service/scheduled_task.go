package service

import (
	"context"

	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
)

// ScheduledTaskService 定时同步任务服务
type ScheduledTaskService struct {
	db *gorm.DB
}

// NewScheduledTaskService 创建定时同步任务服务
func NewScheduledTaskService(db *gorm.DB) *ScheduledTaskService {
	return &ScheduledTaskService{db: db}
}

// CreateScheduledTask 创建定时同步任务
func (s *ScheduledTaskService) CreateScheduledTask(ctx context.Context, task *model.ScheduledTask) error {
	return s.db.WithContext(ctx).Create(task).Error
}

// GetScheduledTask 获取定时同步任务
func (s *ScheduledTaskService) GetScheduledTask(ctx context.Context, id uint) (*model.ScheduledTask, error) {
	var task model.ScheduledTask
	err := s.db.WithContext(ctx).First(&task, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &task, nil
}

// ListScheduledTasks 列出定时同步任务
func (s *ScheduledTaskService) ListScheduledTasks(ctx context.Context, limit, offset int) ([]*model.ScheduledTask, int64, error) {
	var tasks []*model.ScheduledTask
	var total int64

	if err := s.db.Model(&model.ScheduledTask{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := s.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&tasks).Error

	return tasks, total, err
}

// ListScheduledTasksByAccount 列出指定云账号的定时同步任务
func (s *ScheduledTaskService) ListScheduledTasksByAccount(ctx context.Context, cloudAccountID uint, limit, offset int) ([]*model.ScheduledTask, int64, error) {
	var tasks []*model.ScheduledTask
	var total int64

	query := s.db.Model(&model.ScheduledTask{}).Where("cloud_account_id = ? OR cloud_account_id IS NULL", cloudAccountID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := s.db.WithContext(ctx).
		Where("cloud_account_id = ? OR cloud_account_id IS NULL", cloudAccountID).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&tasks).Error

	return tasks, total, err
}

// UpdateScheduledTask 更新定时同步任务
func (s *ScheduledTaskService) UpdateScheduledTask(ctx context.Context, task *model.ScheduledTask) error {
	return s.db.WithContext(ctx).Save(task).Error
}

// DeleteScheduledTask 删除定时同步任务
func (s *ScheduledTaskService) DeleteScheduledTask(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Delete(&model.ScheduledTask{}, id).Error
}

// UpdateScheduledTaskStatus 更新定时同步任务状态
func (s *ScheduledTaskService) UpdateScheduledTaskStatus(ctx context.Context, id uint, status string) error {
	return s.db.WithContext(ctx).Model(&model.ScheduledTask{}).Where("id = ?", id).Update("status", status).Error
}
