package service

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
	"github.com/opencmp/opencmp/pkg/cloudprovider"
)

// MonitorService 监控服务
type MonitorService struct {
	db *gorm.DB
}

// NewMonitorService 创建监控服务
func NewMonitorService(db *gorm.DB) *MonitorService {
	return &MonitorService{db: db}
}

// getProvider 获取云提供商
func (s *MonitorService) getProvider(ctx context.Context, accountID uint) (cloudprovider.ICloudProvider, error) {
	accountService := NewCloudAccountService(s.db)
	account, err := accountService.GetCloudAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, cloudprovider.NewCloudError(cloudprovider.ErrResourceNotFound, "account not found", "")
	}

	var creds map[string]string
	if err := json.Unmarshal(account.Credentials, &creds); err != nil {
		return nil, err
	}

	providerConfig := cloudprovider.CloudAccountConfig{
		ID:           strconv.Itoa(int(account.ID)),
		Name:         account.Name,
		ProviderType: account.ProviderType,
		Credentials:  creds,
		Region:       "",
	}

	return cloudprovider.GetProvider(account.ProviderType, providerConfig)
}

// SyncMonitorResources 同步监控资源
func (s *MonitorService) SyncMonitorResources(ctx context.Context, accountID uint) ([]model.MonitorResource, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	// 获取VM列表作为监控资源
	vms, err := provider.ListVMs(ctx, cloudprovider.VMListFilter{})
	if err != nil {
		return nil, err
	}

	var resources []model.MonitorResource
	for _, vm := range vms {
		resource := model.MonitorResource{
			ResourceID:    vm.ID,
			ResourceName:  vm.Name,
			ResourceType:  "vm",
			MonitorStatus: "正常",
			AccountID:     accountID,
			Platform:      provider.GetCloudInfo().Provider,
			Region:        vm.RegionID,
			LastSyncAt:    time.Now(),
		}

		// 检查是否已存在
		var existing model.MonitorResource
		result := s.db.Where("resource_id = ? AND account_id = ?", vm.ID, accountID).First(&existing)
		if result.Error == nil {
			// 更现存在记录
			existing.ResourceName = vm.Name
			existing.MonitorStatus = "正常"
			existing.LastSyncAt = time.Now()
			s.db.Save(&existing)
			resources = append(resources, existing)
		} else if result.Error == gorm.ErrRecordNotFound {
			// 创建新记录
			s.db.Create(&resource)
			resources = append(resources, resource)
		}
	}

	return resources, nil
}

// GetResourceMetrics 获取资源监控指标
func (s *MonitorService) GetResourceMetrics(ctx context.Context, accountID uint, resourceID string) (map[string]interface{}, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	// 定义要查询的指标
	metrics := []string{"cpu_usage", "memory_usage", "disk_usage", "network_in", "network_out"}

	// 设置时间范围 (最近10分钟)
	endTime := time.Now().Format("2006-01-02T15:04:05Z")
	startTime := time.Now().Add(-10 * time.Minute).Format("2006-01-02T15:04:05Z")

	// 调用云厂商监控API获取真实数据
	instanceMetrics, err := provider.GetInstanceMetrics(ctx, resourceID, metrics, 60, startTime, endTime)
	if err != nil {
		// 如果获取失败，返回空数据而不是错误，让前端能够正常显示
		return s.getDefaultMetrics(resourceID), nil
	}

	// 转换为前端期望的格式
	result := map[string]interface{}{
		"resource_id": resourceID,
		"timestamp":   time.Now().Unix(),
	}

	for metricName, metricData := range instanceMetrics.Metrics {
		result[metricName] = map[string]interface{}{
			"value":     metricData.Value,
			"unit":      metricData.Unit,
			"timestamp": metricData.Timestamp,
		}
	}

	return result, nil
}

// getDefaultMetrics 获取默认指标数据（当云厂商API不可用时）
func (s *MonitorService) getDefaultMetrics(resourceID string) map[string]interface{} {
	return map[string]interface{}{
		"resource_id": resourceID,
		"timestamp":   time.Now().Unix(),
		"cpu_usage": map[string]interface{}{
			"value":     0,
			"unit":      "percent",
			"timestamp": time.Now().Unix(),
		},
		"memory_usage": map[string]interface{}{
			"value":     0,
			"unit":      "percent",
			"timestamp": time.Now().Unix(),
		},
		"disk_usage": map[string]interface{}{
			"value":     0,
			"unit":      "percent",
			"timestamp": time.Now().Unix(),
		},
		"network_in": map[string]interface{}{
			"value":     0,
			"unit":      "KB/s",
			"timestamp": time.Now().Unix(),
		},
		"network_out": map[string]interface{}{
			"value":     0,
			"unit":      "KB/s",
			"timestamp": time.Now().Unix(),
		},
	}
}

// CreateAlertPolicy 创建告警策略
func (s *MonitorService) CreateAlertPolicy(policy model.AlertPolicy) (*model.AlertPolicy, error) {
	policy.CreatedAt = time.Now()
	policy.UpdatedAt = time.Now()

	if err := s.db.Create(&policy).Error; err != nil {
		return nil, err
	}

	return &policy, nil
}

// ListAlertPolicies 列出告警策略
func (s *MonitorService) ListAlertPolicies(filter map[string]interface{}) ([]model.AlertPolicy, error) {
	var policies []model.AlertPolicy

	query := s.db.Model(&model.AlertPolicy{})

	for key, value := range filter {
		if value != "" && value != nil {
			query = query.Where(key+" = ?", value)
		}
	}

	if err := query.Find(&policies).Error; err != nil {
		return nil, err
	}

	return policies, nil
}

// UpdateAlertPolicy 更新告警策略
func (s *MonitorService) UpdateAlertPolicy(id uint, updates map[string]interface{}) (*model.AlertPolicy, error) {
	var policy model.AlertPolicy
	if err := s.db.First(&policy, id).Error; err != nil {
		return nil, err
	}

	updates["updated_at"] = time.Now()
	if err := s.db.Model(&policy).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &policy, nil
}

// DeleteAlertPolicy 删除告警策略
func (s *MonitorService) DeleteAlertPolicy(id uint) error {
	return s.db.Delete(&model.AlertPolicy{}, id).Error
}

// ListAlertHistory 列出告警历史
func (s *MonitorService) ListAlertHistory(filter map[string]interface{}, page, pageSize int) ([]model.AlertHistory, int64, error) {
	var history []model.AlertHistory
	var total int64

	query := s.db.Model(&model.AlertHistory{})

	for key, value := range filter {
		if value != "" && value != nil {
			query = query.Where(key+" = ?", value)
		}
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	query = query.Order("triggered_at desc").Offset(offset).Limit(pageSize)

	if err := query.Find(&history).Error; err != nil {
		return nil, 0, err
	}

	return history, total, nil
}