package service

import (
	"context"
	"encoding/json"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
)

// SyncLogService 同步日志服务
type SyncLogService struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewSyncLogService 创建同步日志服务
func NewSyncLogService(db *gorm.DB, logger *zap.Logger) *SyncLogService {
	return &SyncLogService{
		db:     db,
		logger: logger,
	}
}

// StartSyncLog 创建并开始同步日志
func (s *SyncLogService) StartSyncLog(ctx context.Context, cloudAccountID uint, cloudAccountName string, syncMode model.SyncMode, resourceType model.SyncResourceType, triggeredBy model.SyncTriggerType, scheduledTaskID *uint) (*model.SyncLog, error) {
	// 获取云账户的域ID
	var cloudAccount model.CloudAccount
	if err := s.db.WithContext(ctx).Select("domain_id").First(&cloudAccount, cloudAccountID).Error; err != nil {
		return nil, err
	}

	syncLog := &model.SyncLog{
		CloudAccountID:    cloudAccountID,
		CloudAccountName:  cloudAccountName,
		SyncMode:          string(syncMode),
		ResourceType:      string(resourceType),
		SyncStartTime:     time.Now(),
		Status:            string(model.SyncLogStatusRunning),
		DomainID:          cloudAccount.DomainID,
		TriggeredBy:       string(triggeredBy),
		ScheduledTaskID:   scheduledTaskID,
	}

	if err := s.db.WithContext(ctx).Create(syncLog).Error; err != nil {
		return nil, err
	}

	s.logger.Info("开始同步日志记录",
		zap.Uint("sync_log_id", syncLog.ID),
		zap.Uint("cloud_account_id", cloudAccountID),
		zap.String("sync_mode", string(syncMode)),
		zap.String("resource_type", string(resourceType)))

	return syncLog, nil
}

// UpdateSyncLogProgress 更新同步进度
func (s *SyncLogService) UpdateSyncLogProgress(ctx context.Context, syncLogID uint, newCount, updatedCount, deletedCount, skippedCount, errorCount int) error {
	return s.db.WithContext(ctx).Model(&model.SyncLog{}).
		Where("id = ?", syncLogID).
		Updates(map[string]interface{}{
			"new_count":     newCount,
			"updated_count": updatedCount,
			"deleted_count": deletedCount,
			"skipped_count": skippedCount,
			"error_count":   errorCount,
			"total_count":   newCount + updatedCount + deletedCount + skippedCount,
		}).Error
}

// CompleteSyncLog 完成同步日志
func (s *SyncLogService) CompleteSyncLog(ctx context.Context, syncLogID uint, status model.SyncLogStatus, errorMessage string) error {
	now := time.Now()

	var syncLog model.SyncLog
	if err := s.db.WithContext(ctx).First(&syncLog, syncLogID).Error; err != nil {
		return err
	}

	duration := int(now.Sub(syncLog.SyncStartTime).Seconds())

	return s.db.WithContext(ctx).Model(&model.SyncLog{}).
		Where("id = ?", syncLogID).
		Updates(map[string]interface{}{
			"sync_end_time": now,
			"sync_duration": duration,
			"status":        string(status),
			"error_message": errorMessage,
		}).Error
}

// AddSyncLogDetail 添加同步日志详情
func (s *SyncLogService) AddSyncLogDetail(ctx context.Context, syncLogID uint, detail model.SyncLogDetail) error {
	// 获取现有详情
	var syncLog model.SyncLog
	if err := s.db.WithContext(ctx).Select("details").First(&syncLog, syncLogID).Error; err != nil {
		return err
	}

	// 解析现有详情
	var details []model.SyncLogDetail
	if syncLog.Details != "" {
		if err := json.Unmarshal([]byte(syncLog.Details), &details); err != nil {
			details = []model.SyncLogDetail{}
		}
	}

	// 添加新详情
	detail.Timestamp = time.Now().Format(time.RFC3339)
	details = append(details, detail)

	// 保存详情
	detailsJSON, err := json.Marshal(details)
	if err != nil {
		return err
	}

	return s.db.WithContext(ctx).Model(&model.SyncLog{}).
		Where("id = ?", syncLogID).
		Update("details", string(detailsJSON)).Error
}

// BatchAddSyncLogDetails 批量添加同步日志详情
func (s *SyncLogService) BatchAddSyncLogDetails(ctx context.Context, syncLogID uint, details []model.SyncLogDetail) error {
	if len(details) == 0 {
		return nil
	}

	// 设置时间戳
	now := time.Now().Format(time.RFC3339)
	for i := range details {
		details[i].Timestamp = now
	}

	// 获取现有详情
	var syncLog model.SyncLog
	if err := s.db.WithContext(ctx).Select("details").First(&syncLog, syncLogID).Error; err != nil {
		return err
	}

	// 解析现有详情
	var existingDetails []model.SyncLogDetail
	if syncLog.Details != "" {
		if err := json.Unmarshal([]byte(syncLog.Details), &existingDetails); err != nil {
			existingDetails = []model.SyncLogDetail{}
		}
	}

	// 合并详情
	existingDetails = append(existingDetails, details...)

	// 保存详情
	detailsJSON, err := json.Marshal(existingDetails)
	if err != nil {
		return err
	}

	return s.db.WithContext(ctx).Model(&model.SyncLog{}).
		Where("id = ?", syncLogID).
		Update("details", string(detailsJSON)).Error
}

// GetSyncLogs 获取同步日志列表
func (s *SyncLogService) GetSyncLogs(ctx context.Context, cloudAccountID uint, limit int) ([]model.SyncLog, error) {
	var logs []model.SyncLog
	query := s.db.WithContext(ctx).
		Where("cloud_account_id = ?", cloudAccountID).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&logs).Error; err != nil {
		return nil, err
	}

	return logs, nil
}

// GetLatestSyncLog 获取最新同步日志
func (s *SyncLogService) GetLatestSyncLog(ctx context.Context, cloudAccountID uint) (*model.SyncLog, error) {
	var log model.SyncLog
	err := s.db.WithContext(ctx).
		Where("cloud_account_id = ?", cloudAccountID).
		Order("created_at DESC").
		First(&log).Error

	if err != nil {
		return nil, err
	}

	return &log, nil
}

// GetSyncLogByID 获取同步日志详情
func (s *SyncLogService) GetSyncLogByID(ctx context.Context, syncLogID uint) (*model.SyncLog, error) {
	var log model.SyncLog
	err := s.db.WithContext(ctx).First(&log, syncLogID).Error
	if err != nil {
		return nil, err
	}

	return &log, nil
}

// GetSyncLogsByDateRange 获取指定时间范围的同步日志
func (s *SyncLogService) GetSyncLogsByDateRange(ctx context.Context, domainID uint, startTime, endTime time.Time) ([]model.SyncLog, error) {
	var logs []model.SyncLog
	err := s.db.WithContext(ctx).
		Where("domain_id = ?", domainID).
		Where("sync_start_time >= ?", startTime).
		Where("sync_start_time <= ?", endTime).
		Order("sync_start_time DESC").
		Find(&logs).Error

	if err != nil {
		return nil, err
	}

	return logs, nil
}

// GetSyncStatistics 获取同步统计信息
func (s *SyncLogService) GetSyncStatistics(ctx context.Context, cloudAccountID uint, days int) (*SyncStatistics, error) {
	startTime := time.Now().AddDate(0, 0, -days)

	var stats SyncStatistics
	stats.CloudAccountID = cloudAccountID
	stats.PeriodDays = days

	// 总同步次数
	err := s.db.WithContext(ctx).Model(&model.SyncLog{}).
		Where("cloud_account_id = ?", cloudAccountID).
		Where("sync_start_time >= ?", startTime).
		Count(&stats.TotalSyncCount).Error
	if err != nil {
		return nil, err
	}

	// 成功次数
	err = s.db.WithContext(ctx).Model(&model.SyncLog{}).
		Where("cloud_account_id = ?", cloudAccountID).
		Where("sync_start_time >= ?", startTime).
		Where("status = ?", string(model.SyncLogStatusSuccess)).
		Count(&stats.SuccessCount).Error
	if err != nil {
		return nil, err
	}

	// 失败次数
	err = s.db.WithContext(ctx).Model(&model.SyncLog{}).
		Where("cloud_account_id = ?", cloudAccountID).
		Where("sync_start_time >= ?", startTime).
		Where("status IN ?", []string{string(model.SyncLogStatusFailed), string(model.SyncLogStatusPartialFail)}).
		Count(&stats.FailureCount).Error
	if err != nil {
		return nil, err
	}

	// 平均同步耗时
	var avgDurationFloat float64
	err = s.db.WithContext(ctx).Model(&model.SyncLog{}).
		Where("cloud_account_id = ?", cloudAccountID).
		Where("sync_start_time >= ?", startTime).
		Where("status = ?", string(model.SyncLogStatusSuccess)).
		Select("COALESCE(AVG(sync_duration), 0)").
		Scan(&avgDurationFloat).Error
	if err != nil {
		return nil, err
	}
	stats.AvgDuration = int(avgDurationFloat)

	// 资源统计
	var resourceStats struct {
		TotalNew     int
		TotalUpdated int
		TotalDeleted int
	}
	err = s.db.WithContext(ctx).Model(&model.SyncLog{}).
		Where("cloud_account_id = ?", cloudAccountID).
		Where("sync_start_time >= ?", startTime).
		Select("COALESCE(SUM(new_count), 0) as total_new, COALESCE(SUM(updated_count), 0) as total_updated, COALESCE(SUM(deleted_count), 0) as total_deleted").
		Scan(&resourceStats).Error
	if err != nil {
		return nil, err
	}
	stats.TotalNew = resourceStats.TotalNew
	stats.TotalUpdated = resourceStats.TotalUpdated
	stats.TotalDeleted = resourceStats.TotalDeleted

	return &stats, nil
}

// SyncStatistics 同步统计信息
type SyncStatistics struct {
	CloudAccountID uint   `json:"cloud_account_id"`
	PeriodDays      int    `json:"period_days"`
	TotalSyncCount  int64  `json:"total_sync_count"`
	SuccessCount    int64  `json:"success_count"`
	FailureCount    int64  `json:"failure_count"`
	AvgDuration     int    `json:"avg_duration"` // 平均耗时（秒）
	TotalNew        int    `json:"total_new"`
	TotalUpdated    int    `json:"total_updated"`
	TotalDeleted    int    `json:"total_deleted"`
}

// DeleteOldSyncLogs 删除旧的同步日志（数据清理）
func (s *SyncLogService) DeleteOldSyncLogs(ctx context.Context, retentionDays int) error {
	cutoffTime := time.Now().AddDate(0, 0, -retentionDays)

	result := s.db.WithContext(ctx).
		Where("sync_start_time < ?", cutoffTime).
		Delete(&model.SyncLog{})

	if result.Error != nil {
		return result.Error
	}

	s.logger.Info("清理旧同步日志",
		zap.Int64("deleted_count", result.RowsAffected),
		zap.Int("retention_days", retentionDays))

	return nil
}