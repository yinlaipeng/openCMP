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

// DatabaseService 数据库服务
type DatabaseService struct {
	db *gorm.DB
}

// NewDatabaseService 创建数据库服务
func NewDatabaseService(db *gorm.DB) *DatabaseService {
	return &DatabaseService{db: db}
}

// getProvider 获取云提供商
func (s *DatabaseService) getProvider(ctx context.Context, accountID uint) (cloudprovider.ICloudProvider, error) {
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

// logStateChange 记录资源状态变更日志
func (s *DatabaseService) logStateChange(ctx context.Context, log *model.ResourceStateLog) error {
	log.OccurredAt = time.Now()
	log.CreatedAt = time.Now()
	return s.db.WithContext(ctx).Create(log).Error
}

// CreateRDS 创建 RDS 实例
func (s *DatabaseService) CreateRDS(ctx context.Context, accountID uint, config cloudprovider.RDSConfig) (*cloudprovider.RDSInstance, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	dbProvider, ok := provider.(cloudprovider.IDatabase)
	if !ok {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrUnsupportedOperation,
			"provider does not support database operations",
			"",
		)
	}

	instance, err := dbProvider.CreateRDSInstance(ctx, config)
	if err != nil {
		return nil, err
	}

	// 记录状态变更日志
	stateLog := &model.ResourceStateLog{
		ResourceType:   "rds",
		ResourceID:     instance.ID,
		ResourceName:   config.Name,
		CloudAccountID: accountID,
		PreviousStatus: "",
		CurrentStatus:  instance.Status,
		OperationType:  "create",
		Reason:         "RDS实例创建，引擎:" + config.Engine,
	}
	if err := s.logStateChange(ctx, stateLog); err != nil {
		// 日志记录失败不影响主流程
	}

	return instance, nil
}

// ListRDS 列出 RDS 实例
func (s *DatabaseService) ListRDS(ctx context.Context, accountID uint, filter cloudprovider.RDSFilter) ([]*cloudprovider.RDSInstance, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	dbProvider, ok := provider.(cloudprovider.IDatabase)
	if !ok {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrUnsupportedOperation,
			"provider does not support database operations",
			"",
		)
	}

	return dbProvider.ListRDSInstances(ctx, filter)
}

// GetRDS 获取 RDS 实例详情
func (s *DatabaseService) GetRDS(ctx context.Context, accountID uint, instanceID string) (*cloudprovider.RDSInstance, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	dbProvider, ok := provider.(cloudprovider.IDatabase)
	if !ok {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrUnsupportedOperation,
			"provider does not support database operations",
			"",
		)
	}

	filter := cloudprovider.RDSFilter{
		InstanceID: instanceID,
	}

	instances, err := dbProvider.ListRDSInstances(ctx, filter)
	if err != nil {
		return nil, err
	}

	if len(instances) == 0 {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrResourceNotFound,
			"rds instance not found",
			instanceID,
		)
	}

	return instances[0], nil
}

// DeleteRDS 删除 RDS 实例
func (s *DatabaseService) DeleteRDS(ctx context.Context, accountID uint, instanceID string) error {
	// 先获取当前状态用于日志记录
	instance, err := s.GetRDS(ctx, accountID, instanceID)
	if err != nil {
		return err
	}

	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return err
	}

	dbProvider, ok := provider.(cloudprovider.IDatabase)
	if !ok {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrUnsupportedOperation,
			"provider does not support database operations",
			"",
		)
	}

	if err := dbProvider.DeleteRDSInstance(ctx, instanceID); err != nil {
		return err
	}

	// 记录状态变更日志
	stateLog := &model.ResourceStateLog{
		ResourceType:   "rds",
		ResourceID:     instanceID,
		ResourceName:   instance.Name,
		CloudAccountID: accountID,
		PreviousStatus: instance.Status,
		CurrentStatus:  "terminated",
		OperationType:  "delete",
		Reason:         "RDS实例删除",
	}
	if err := s.logStateChange(ctx, stateLog); err != nil {
		// 日志记录失败不影响主流程
	}

	return nil
}

// RDSAction 执行 RDS 操作
func (s *DatabaseService) RDSAction(ctx context.Context, accountID uint, instanceID string, action string) error {
	// 先获取当前状态用于日志记录
	instance, err := s.GetRDS(ctx, accountID, instanceID)
	if err != nil {
		return err
	}

	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return err
	}

	dbProvider, ok := provider.(cloudprovider.IDatabase)
	if !ok {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrUnsupportedOperation,
			"provider does not support database operations",
			"",
		)
	}

	prevStatus := instance.Status
	switch action {
	case "start":
		err = dbProvider.StartRDSInstance(ctx, instanceID)
	case "stop":
		err = dbProvider.StopRDSInstance(ctx, instanceID)
	case "reboot":
		err = dbProvider.RebootRDSInstance(ctx, instanceID)
	default:
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"invalid action: "+action,
			"",
		)
	}

	if err != nil {
		return err
	}

	// 记录状态变更日志
	newStatus := prevStatus
	switch action {
	case "start":
		newStatus = "Starting"
	case "stop":
		newStatus = "Stopping"
	case "reboot":
		newStatus = "Restarting"
	}

	stateLog := &model.ResourceStateLog{
		ResourceType:   "rds",
		ResourceID:     instanceID,
		ResourceName:   instance.Name,
		CloudAccountID: accountID,
		PreviousStatus: prevStatus,
		CurrentStatus:  newStatus,
		OperationType:  action,
		Reason:         "RDS实例" + action + "操作",
	}
	if err := s.logStateChange(ctx, stateLog); err != nil {
		// 日志记录失败不影响主流程
	}

	return nil
}

// ResizeRDS 调整 RDS 规格
func (s *DatabaseService) ResizeRDS(ctx context.Context, accountID uint, instanceID string, spec cloudprovider.RDSpec) error {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return err
	}

	dbProvider, ok := provider.(cloudprovider.IDatabase)
	if !ok {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrUnsupportedOperation,
			"provider does not support database operations",
			"",
		)
	}

	return dbProvider.ResizeRDSInstance(ctx, instanceID, spec)
}

// CreateRDSBackup 创建 RDS 备份
func (s *DatabaseService) CreateRDSBackup(ctx context.Context, accountID uint, instanceID string, name string) (*cloudprovider.RDSBackup, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	dbProvider, ok := provider.(cloudprovider.IDatabase)
	if !ok {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrUnsupportedOperation,
			"provider does not support database operations",
			"",
		)
	}

	return dbProvider.CreateRDSBackup(ctx, instanceID, name)
}

// ListRDSBackups 列出 RDS 备份
func (s *DatabaseService) ListRDSBackups(ctx context.Context, accountID uint, instanceID string) ([]*cloudprovider.RDSBackup, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	dbProvider, ok := provider.(cloudprovider.IDatabase)
	if !ok {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrUnsupportedOperation,
			"provider does not support database operations",
			"",
		)
	}

	return dbProvider.ListRDSBackups(ctx, instanceID)
}

// CreateCache 创建缓存实例
func (s *DatabaseService) CreateCache(ctx context.Context, accountID uint, config cloudprovider.CacheConfig) (*cloudprovider.CacheInstance, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	dbProvider, ok := provider.(cloudprovider.IDatabase)
	if !ok {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrUnsupportedOperation,
			"provider does not support database operations",
			"",
		)
	}

	instance, err := dbProvider.CreateCacheInstance(ctx, config)
	if err != nil {
		return nil, err
	}

	// 记录状态变更日志
	stateLog := &model.ResourceStateLog{
		ResourceType:   "redis",
		ResourceID:     instance.ID,
		ResourceName:   config.Name,
		CloudAccountID: accountID,
		PreviousStatus: "",
		CurrentStatus:  instance.Status,
		OperationType:  "create",
		Reason:         "Redis实例创建，版本:" + config.EngineVersion,
	}
	if err := s.logStateChange(ctx, stateLog); err != nil {
		// 日志记录失败不影响主流程
	}

	return instance, nil
}

// ListCache 列出缓存实例
func (s *DatabaseService) ListCache(ctx context.Context, accountID uint, filter cloudprovider.CacheFilter) ([]*cloudprovider.CacheInstance, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	dbProvider, ok := provider.(cloudprovider.IDatabase)
	if !ok {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrUnsupportedOperation,
			"provider does not support database operations",
			"",
		)
	}

	return dbProvider.ListCacheInstances(ctx, filter)
}

// DeleteCache 删除缓存实例
func (s *DatabaseService) DeleteCache(ctx context.Context, accountID uint, instanceID string) error {
	// 先获取当前状态用于日志记录
	filter := cloudprovider.CacheFilter{InstanceID: instanceID}
	instances, err := s.ListCache(ctx, accountID, filter)
	if err != nil {
		return err
	}
	if len(instances) == 0 {
		return cloudprovider.NewCloudError(cloudprovider.ErrResourceNotFound, "cache instance not found", instanceID)
	}
	instance := instances[0]

	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return err
	}

	dbProvider, ok := provider.(cloudprovider.IDatabase)
	if !ok {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrUnsupportedOperation,
			"provider does not support database operations",
			"",
		)
	}

	if err := dbProvider.DeleteCacheInstance(ctx, instanceID); err != nil {
		return err
	}

	// 记录状态变更日志
	stateLog := &model.ResourceStateLog{
		ResourceType:   "redis",
		ResourceID:     instanceID,
		ResourceName:   instance.Name,
		CloudAccountID: accountID,
		PreviousStatus: instance.Status,
		CurrentStatus:  "terminated",
		OperationType:  "delete",
		Reason:         "Redis实例删除",
	}
	if err := s.logStateChange(ctx, stateLog); err != nil {
		// 日志记录失败不影响主流程
	}

	return nil
}

// CacheAction 执行缓存实例操作
func (s *DatabaseService) CacheAction(ctx context.Context, accountID uint, instanceID string, action string) error {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return err
	}

	dbProvider, ok := provider.(cloudprovider.IDatabase)
	if !ok {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrUnsupportedOperation,
			"provider does not support database operations",
			"",
		)
	}

	switch action {
	case "reboot":
		return dbProvider.RebootCacheInstance(ctx, instanceID)
	default:
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"invalid action: "+action,
			"",
		)
	}
}

// ResizeCache 调整缓存实例规格
func (s *DatabaseService) ResizeCache(ctx context.Context, accountID uint, instanceID string, spec cloudprovider.CacheSpec) error {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return err
	}

	dbProvider, ok := provider.(cloudprovider.IDatabase)
	if !ok {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrUnsupportedOperation,
			"provider does not support database operations",
			"",
		)
	}

	return dbProvider.ResizeCacheInstance(ctx, instanceID, spec)
}

// CreateCacheBackup 创建缓存备份
func (s *DatabaseService) CreateCacheBackup(ctx context.Context, accountID uint, instanceID string) (*cloudprovider.CacheBackup, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	dbProvider, ok := provider.(cloudprovider.IDatabase)
	if !ok {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrUnsupportedOperation,
			"provider does not support database operations",
			"",
		)
	}

	return dbProvider.CreateCacheBackup(ctx, instanceID)
}

// ListRDSSKUs 列出 RDS SKU 规格
// 参考 Cloudpods API: GET /api/v2/dbinstance_skus
func (s *DatabaseService) ListRDSSKUs(ctx context.Context, accountID uint, filter cloudprovider.RDSKUFilter) ([]*cloudprovider.RDSInstanceSKU, error) {
	// 获取云账号信息以确定 provider
	accountService := NewCloudAccountService(s.db)
	account, err := accountService.GetCloudAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}

	// 基于云账号 provider 返回 SKU 数据
	// 实际应调用云厂商 API，这里返回基础 SKU 作为参考
	skus := getMockRDSSKUs(account.ProviderType, filter)
	return skus, nil
}

// ListCacheSKUs 列出 Redis 缓存 SKU 规格
// 参考 Cloudpods API: GET /api/v2/elasticcacheskus
func (s *DatabaseService) ListCacheSKUs(ctx context.Context, accountID uint, filter cloudprovider.CacheSKUFilter) ([]*cloudprovider.CacheInstanceSKU, error) {
	// 获取云账号信息以确定 provider
	accountService := NewCloudAccountService(s.db)
	account, err := accountService.GetCloudAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}

	// 基于云账号 provider 返回 SKU 数据
	skus := getMockCacheSKUs(account.ProviderType, filter)
	return skus, nil
}

// getMockRDSSKUs 获取模拟 RDS SKU 数据（实际应调用云厂商 API）
func getMockRDSSKUs(provider string, filter cloudprovider.RDSKUFilter) []*cloudprovider.RDSInstanceSKU {
	// 基于过滤条件返回 SKU
	skus := []*cloudprovider.RDSInstanceSKU{
		{
			ID:           "rds-sku-1",
			Provider:     provider,
			Engine:       "MySQL",
			EngineVersion: "5.7",
			Category:     "ha",
			StorageType:  "cloud_essd",
			CPU:          2,
			MemoryMB:     4096,
			InstanceType: "mysql.n2.medium.1c",
			Price:        0.5,
			RegionID:     filter.RegionID,
		},
		{
			ID:           "rds-sku-2",
			Provider:     provider,
			Engine:       "MySQL",
			EngineVersion: "8.0",
			Category:     "ha",
			StorageType:  "cloud_essd",
			CPU:          4,
			MemoryMB:     8192,
			InstanceType: "mysql.n4.large.1c",
			Price:        1.0,
			RegionID:     filter.RegionID,
		},
		{
			ID:           "rds-sku-3",
			Provider:     provider,
			Engine:       "PostgreSQL",
			EngineVersion: "14",
			Category:     "ha",
			StorageType:  "cloud_essd",
			CPU:          2,
			MemoryMB:     4096,
			InstanceType: "pg.n2.medium.1c",
			Price:        0.6,
			RegionID:     filter.RegionID,
		},
	}

	// 应用过滤条件
	result := []*cloudprovider.RDSInstanceSKU{}
	for _, sku := range skus {
		if filter.Engine != "" && sku.Engine != filter.Engine {
			continue
		}
		if filter.EngineVersion != "" && sku.EngineVersion != filter.EngineVersion {
			continue
		}
		if filter.Category != "" && sku.Category != filter.Category {
			continue
		}
		if filter.StorageType != "" && sku.StorageType != filter.StorageType {
			continue
		}
		result = append(result, sku)
	}

	return result
}

// getMockCacheSKUs 获取模拟 Redis SKU 数据（实际应调用云厂商 API）
func getMockCacheSKUs(provider string, filter cloudprovider.CacheSKUFilter) []*cloudprovider.CacheInstanceSKU {
	skus := []*cloudprovider.CacheInstanceSKU{
		{
			ID:             "cache-sku-1",
			Provider:       provider,
			Engine:         "Redis",
			EngineVersion:  "5.0",
			NodeType:       "standalone",
			PerformanceType: "standard",
			MemoryMB:       1024,
			InstanceType:   "redis.standard.small.default",
			Price:          0.1,
			RegionID:       filter.RegionID,
		},
		{
			ID:             "cache-sku-2",
			Provider:       provider,
			Engine:         "Redis",
			EngineVersion:  "6.0",
			NodeType:       "ha",
			PerformanceType: "standard",
			MemoryMB:       2048,
			InstanceType:   "redis.standard.medium.default",
			Price:          0.2,
			RegionID:       filter.RegionID,
		},
		{
			ID:             "cache-sku-3",
			Provider:       provider,
			Engine:         "Redis",
			EngineVersion:  "7.0",
			NodeType:       "cluster",
			PerformanceType: "performance",
			MemoryMB:       4096,
			InstanceType:   "redis.performance.large.default",
			Price:          0.5,
			RegionID:       filter.RegionID,
		},
	}

	// 应用过滤条件
	result := []*cloudprovider.CacheInstanceSKU{}
	for _, sku := range skus {
		if filter.Engine != "" && sku.Engine != filter.Engine {
			continue
		}
		if filter.EngineVersion != "" && sku.EngineVersion != filter.EngineVersion {
			continue
		}
		if filter.NodeType != "" && sku.NodeType != filter.NodeType {
			continue
		}
		if filter.PerformanceType != "" && sku.PerformanceType != filter.PerformanceType {
			continue
		}
		result = append(result, sku)
	}

	return result
}