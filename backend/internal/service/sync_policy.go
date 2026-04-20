package service

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
)

// SyncPolicyService 同步策略服务
type SyncPolicyService struct {
	db *gorm.DB
}

// NewSyncPolicyService 创建同步策略服务
func NewSyncPolicyService(db *gorm.DB) *SyncPolicyService {
	return &SyncPolicyService{db: db}
}

// CreateSyncPolicy 创建同步策略
func (s *SyncPolicyService) CreateSyncPolicy(ctx context.Context, policy *model.SyncPolicy, rules []model.Rule, ruleTags []model.RuleTag) error {
	tx := s.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(policy).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 创建规则
	for i := range rules {
		rules[i].SyncPolicyID = policy.ID
		if err := tx.Create(&rules[i]).Error; err != nil {
			tx.Rollback()
			return err
		}

		// 为当前规则创建关联的标签
		for j := range ruleTags {
			if ruleTags[j].RuleID == 0 { // 如果还没有关联到规则，则根据顺序分配
				// 查找刚创建的规则来确认ID
				rule := model.Rule{}
				if err := tx.Where("sync_policy_id = ? AND condition_type = ? AND resource_mapping = ?",
					policy.ID, rules[i].ConditionType, rules[i].ResourceMapping).Last(&rule).Error; err != nil {
					// 如果无法通过其他字段定位，则按顺序关联
					rule = rules[i]
				}

				if rule.ID != 0 {
					ruleTags[j].RuleID = rule.ID
				}
			}

			if err := tx.Create(&ruleTags[j]).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit().Error
}

// GetSyncPolicy 获取同步策略
func (s *SyncPolicyService) GetSyncPolicy(ctx context.Context, id uint) (*model.SyncPolicy, error) {
	var policy model.SyncPolicy
	err := s.db.WithContext(ctx).Preload("Rules.Tags").First(&policy, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &policy, nil
}

// ListSyncPolicies 列出同步策略
func (s *SyncPolicyService) ListSyncPolicies(ctx context.Context, limit, offset int) ([]*model.SyncPolicy, int64, error) {
	var policies []*model.SyncPolicy
	var total int64

	if err := s.db.Model(&model.SyncPolicy{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := s.db.WithContext(ctx).
		Preload("Rules.Tags").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&policies).Error

	return policies, total, err
}

// UpdateSyncPolicy 更新同步策略
func (s *SyncPolicyService) UpdateSyncPolicy(ctx context.Context, policy *model.SyncPolicy, rules []model.Rule, ruleTags []model.RuleTag) error {
	tx := s.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Save(policy).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 删除旧的规则和标签
	if err := tx.Where("sync_policy_id = ?", policy.ID).Delete(&model.Rule{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("rule_id IN (SELECT id FROM sync_policy_rules WHERE sync_policy_id = ?)", policy.ID).Delete(&model.RuleTag{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 创建新的规则
	for i := range rules {
		rules[i].SyncPolicyID = policy.ID
		if err := tx.Create(&rules[i]).Error; err != nil {
			tx.Rollback()
			return err
		}

		// 为当前规则创建关联的标签
		for j := range ruleTags {
			if ruleTags[j].RuleID == 0 { // 如果还没有关联到规则，则根据顺序分配
				// 查找刚创建的规则来确认ID
				rule := model.Rule{}
				if err := tx.Where("sync_policy_id = ? AND condition_type = ? AND resource_mapping = ?",
					policy.ID, rules[i].ConditionType, rules[i].ResourceMapping).Last(&rule).Error; err != nil {
					// 如果无法通过其他字段定位，则按顺序关联
					rule = rules[i]
				}

				if rule.ID != 0 {
					ruleTags[j].RuleID = rule.ID
				}
			}

			if err := tx.Create(&ruleTags[j]).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit().Error
}

// DeleteSyncPolicy 删除同步策略
func (s *SyncPolicyService) DeleteSyncPolicy(ctx context.Context, id uint) error {
	tx := s.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除相关规则和标签
	if err := tx.Where("sync_policy_id = ?", id).Delete(&model.Rule{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("rule_id IN (SELECT id FROM sync_policy_rules WHERE sync_policy_id = ?)", id).Delete(&model.RuleTag{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(&model.SyncPolicy{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// ToggleSyncPolicyStatus 切换同步策略状态
func (s *SyncPolicyService) ToggleSyncPolicyStatus(ctx context.Context, id uint, enabled bool) error {
	policy, err := s.GetSyncPolicy(ctx, id)
	if err != nil {
		return err
	}
	if policy == nil {
		return gorm.ErrRecordNotFound
	}

	policy.Enabled = enabled
	return s.db.WithContext(ctx).Save(policy).Error
}

// ExecuteSyncPolicy 执行同步策略
func (s *SyncPolicyService) ExecuteSyncPolicy(ctx context.Context, policyID uint, cloudAccountID *uint, operator string) (*model.SyncPolicyExecutionLog, error) {
	policy, err := s.GetSyncPolicy(ctx, policyID)
	if err != nil {
		return nil, err
	}
	if policy == nil {
		return nil, gorm.ErrRecordNotFound
	}

	// 创建执行日志
	execLog := &model.SyncPolicyExecutionLog{
		SyncPolicyID:  policyID,
		ExecutionTime: time.Now(),
		TriggerType:   "manual",
		Operator:      operator,
		Result:        "success",
	}

	if cloudAccountID != nil {
		execLog.CloudAccountID = cloudAccountID
	}

	// 模拟执行结果（实际应调用资源匹配逻辑）
	execLog.ResourceCount = 0
	execLog.MatchedCount = 0
	execLog.MappedCount = 0
	execLog.Duration = 100 // 模拟耗时

	// 保存执行日志
	if err := s.db.WithContext(ctx).Create(execLog).Error; err != nil {
		return nil, err
	}

	return execLog, nil
}

// GetExecutionLogs 获取执行日志列表
func (s *SyncPolicyService) GetExecutionLogs(ctx context.Context, policyID uint, limit, offset int, resultFilter string) ([]*model.SyncPolicyExecutionLog, int64, error) {
	var logs []*model.SyncPolicyExecutionLog
	var total int64

	query := s.db.WithContext(ctx).Model(&model.SyncPolicyExecutionLog{}).Where("sync_policy_id = ?", policyID)
	if resultFilter != "" {
		query = query.Where("result = ?", resultFilter)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Limit(limit).Offset(offset).Order("execution_time DESC").Find(&logs).Error
	return logs, total, err
}

// GetMappingResults 获取映射结果列表
func (s *SyncPolicyService) GetMappingResults(ctx context.Context, policyID uint, limit, offset int, projectFilter string) ([]*model.SyncPolicyMappingResult, int64, error) {
	var results []*model.SyncPolicyMappingResult
	var total int64

	query := s.db.WithContext(ctx).Model(&model.SyncPolicyMappingResult{}).Where("sync_policy_id = ?", policyID)
	if projectFilter != "" {
		query = query.Where("target_project_id = ?", projectFilter)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Limit(limit).Offset(offset).Order("mapped_at DESC").Find(&results).Error
	return results, total, err
}

// BatchToggleStatus 批量切换状态
func (s *SyncPolicyService) BatchToggleStatus(ctx context.Context, ids []uint, enabled bool) (int, error) {
	result := s.db.WithContext(ctx).Model(&model.SyncPolicy{}).Where("id IN ?", ids).Update("enabled", enabled)
	return int(result.RowsAffected), result.Error
}

// BatchDelete 批量删除
func (s *SyncPolicyService) BatchDelete(ctx context.Context, ids []uint) (int, error) {
	tx := s.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除相关规则和标签
	for _, id := range ids {
		if err := tx.Where("sync_policy_id = ?", id).Delete(&model.Rule{}).Error; err != nil {
			tx.Rollback()
			return 0, err
		}
		if err := tx.Where("rule_id IN (SELECT id FROM sync_policy_rules WHERE sync_policy_id = ?)", id).Delete(&model.RuleTag{}).Error; err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// 删除策略
	result := tx.Where("id IN ?", ids).Delete(&model.SyncPolicy{})
	if result.Error != nil {
		tx.Rollback()
		return 0, result.Error
	}

	return int(result.RowsAffected), tx.Commit().Error
}

// ExportPolicies 导出策略列表
func (s *SyncPolicyService) ExportPolicies(ctx context.Context) ([]*model.SyncPolicy, error) {
	var policies []*model.SyncPolicy
	err := s.db.WithContext(ctx).Preload("Rules.Tags").Order("created_at DESC").Find(&policies).Error
	return policies, err
}
