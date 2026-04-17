package service

import (
	"context"
	"encoding/json"
	"strconv"

	"gorm.io/gorm"

	"github.com/opencmp/opencmp/pkg/cloudprovider"
)

// MiddlewareService 中间件服务
type MiddlewareService struct {
	db *gorm.DB
}

// NewMiddlewareService 创建中间件服务
func NewMiddlewareService(db *gorm.DB) *MiddlewareService {
	return &MiddlewareService{db: db}
}

// getProvider 获取云提供商
func (s *MiddlewareService) getProvider(ctx context.Context, accountID uint) (cloudprovider.ICloudProvider, error) {
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

// CreateKafka 创建 Kafka 实例
func (s *MiddlewareService) CreateKafka(ctx context.Context, accountID uint, config cloudprovider.KafkaConfig) (*cloudprovider.KafkaInstance, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	mwProvider, ok := provider.(cloudprovider.IMiddleware)
	if !ok {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrUnsupportedOperation,
			"provider does not support middleware operations",
			"",
		)
	}

	return mwProvider.CreateKafkaInstance(ctx, config)
}

// ListKafka 列出 Kafka 实例
func (s *MiddlewareService) ListKafka(ctx context.Context, accountID uint, filter cloudprovider.KafkaFilter) ([]*cloudprovider.KafkaInstance, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	mwProvider, ok := provider.(cloudprovider.IMiddleware)
	if !ok {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrUnsupportedOperation,
			"provider does not support middleware operations",
			"",
		)
	}

	return mwProvider.ListKafkaInstances(ctx, filter)
}

// GetKafka 获取 Kafka 实例详情
func (s *MiddlewareService) GetKafka(ctx context.Context, accountID uint, instanceID string) (*cloudprovider.KafkaInstance, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	mwProvider, ok := provider.(cloudprovider.IMiddleware)
	if !ok {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrUnsupportedOperation,
			"provider does not support middleware operations",
			"",
		)
	}

	filter := cloudprovider.KafkaFilter{
		InstanceID: instanceID,
	}

	instances, err := mwProvider.ListKafkaInstances(ctx, filter)
	if err != nil {
		return nil, err
	}

	if len(instances) == 0 {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrResourceNotFound,
			"kafka instance not found",
			instanceID,
		)
	}

	return instances[0], nil
}

// DeleteKafka 删除 Kafka 实例
func (s *MiddlewareService) DeleteKafka(ctx context.Context, accountID uint, instanceID string) error {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return err
	}

	mwProvider, ok := provider.(cloudprovider.IMiddleware)
	if !ok {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrUnsupportedOperation,
			"provider does not support middleware operations",
			"",
		)
	}

	return mwProvider.DeleteKafkaInstance(ctx, instanceID)
}

// KafkaAction 执行 Kafka 操作
func (s *MiddlewareService) KafkaAction(ctx context.Context, accountID uint, instanceID string, action string) error {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return err
	}

	mwProvider, ok := provider.(cloudprovider.IMiddleware)
	if !ok {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrUnsupportedOperation,
			"provider does not support middleware operations",
			"",
		)
	}

	switch action {
	case "start":
		return mwProvider.StartKafkaInstance(ctx, instanceID)
	case "stop":
		return mwProvider.StopKafkaInstance(ctx, instanceID)
	default:
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"invalid action: "+action,
			"",
		)
	}
}

// ResizeKafka 调整 Kafka 规格
func (s *MiddlewareService) ResizeKafka(ctx context.Context, accountID uint, instanceID string, spec cloudprovider.KafkaSpec) error {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return err
	}

	mwProvider, ok := provider.(cloudprovider.IMiddleware)
	if !ok {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrUnsupportedOperation,
			"provider does not support middleware operations",
			"",
		)
	}

	return mwProvider.ResizeKafkaInstance(ctx, instanceID, spec)
}

// CreateElasticsearch 创建 Elasticsearch 实例
func (s *MiddlewareService) CreateElasticsearch(ctx context.Context, accountID uint, config cloudprovider.ElasticsearchConfig) (*cloudprovider.ElasticsearchInstance, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	mwProvider, ok := provider.(cloudprovider.IMiddleware)
	if !ok {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrUnsupportedOperation,
			"provider does not support middleware operations",
			"",
		)
	}

	return mwProvider.CreateElasticsearchInstance(ctx, config)
}

// ListElasticsearch 列出 Elasticsearch 实例
func (s *MiddlewareService) ListElasticsearch(ctx context.Context, accountID uint, filter cloudprovider.ElasticsearchFilter) ([]*cloudprovider.ElasticsearchInstance, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	mwProvider, ok := provider.(cloudprovider.IMiddleware)
	if !ok {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrUnsupportedOperation,
			"provider does not support middleware operations",
			"",
		)
	}

	return mwProvider.ListElasticsearchInstances(ctx, filter)
}

// GetElasticsearch 获取 Elasticsearch 实例详情
func (s *MiddlewareService) GetElasticsearch(ctx context.Context, accountID uint, instanceID string) (*cloudprovider.ElasticsearchInstance, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	mwProvider, ok := provider.(cloudprovider.IMiddleware)
	if !ok {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrUnsupportedOperation,
			"provider does not support middleware operations",
			"",
		)
	}

	filter := cloudprovider.ElasticsearchFilter{
		InstanceID: instanceID,
	}

	instances, err := mwProvider.ListElasticsearchInstances(ctx, filter)
	if err != nil {
		return nil, err
	}

	if len(instances) == 0 {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrResourceNotFound,
			"elasticsearch instance not found",
			instanceID,
		)
	}

	return instances[0], nil
}

// DeleteElasticsearch 删除 Elasticsearch 实例
func (s *MiddlewareService) DeleteElasticsearch(ctx context.Context, accountID uint, instanceID string) error {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return err
	}

	mwProvider, ok := provider.(cloudprovider.IMiddleware)
	if !ok {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrUnsupportedOperation,
			"provider does not support middleware operations",
			"",
		)
	}

	return mwProvider.DeleteElasticsearchInstance(ctx, instanceID)
}

// ElasticsearchAction 执行 Elasticsearch 操作
func (s *MiddlewareService) ElasticsearchAction(ctx context.Context, accountID uint, instanceID string, action string) error {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return err
	}

	mwProvider, ok := provider.(cloudprovider.IMiddleware)
	if !ok {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrUnsupportedOperation,
			"provider does not support middleware operations",
			"",
		)
	}

	switch action {
	case "start":
		return mwProvider.StartElasticsearchInstance(ctx, instanceID)
	case "stop":
		return mwProvider.StopElasticsearchInstance(ctx, instanceID)
	default:
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"invalid action: "+action,
			"",
		)
	}
}

// ResizeElasticsearch 调整 Elasticsearch 规格
func (s *MiddlewareService) ResizeElasticsearch(ctx context.Context, accountID uint, instanceID string, spec cloudprovider.ElasticsearchSpec) error {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return err
	}

	mwProvider, ok := provider.(cloudprovider.IMiddleware)
	if !ok {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrUnsupportedOperation,
			"provider does not support middleware operations",
			"",
		)
	}

	return mwProvider.ResizeElasticsearchInstance(ctx, instanceID, spec)
}