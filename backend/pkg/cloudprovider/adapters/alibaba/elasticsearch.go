package alibaba

import (
	"context"

	"github.com/opencmp/opencmp/pkg/cloudprovider"
)

// CreateElasticsearchInstance 创建 Elasticsearch 实例
// Note: Requires elasticsearch SDK integration - currently returns stub
func (p *AlibabaProvider) CreateElasticsearchInstance(ctx context.Context, config cloudprovider.ElasticsearchConfig) (*cloudprovider.ElasticsearchInstance, error) {
	return nil, cloudprovider.NewCloudError(
		cloudprovider.ErrUnsupportedOperation,
		"Elasticsearch SDK not integrated - requires elasticsearch client",
		"",
	)
}

// DeleteElasticsearchInstance 删除 Elasticsearch 实例
func (p *AlibabaProvider) DeleteElasticsearchInstance(ctx context.Context, instanceID string) error {
	return cloudprovider.NewCloudError(
		cloudprovider.ErrUnsupportedOperation,
		"Elasticsearch SDK not integrated - requires elasticsearch client",
		instanceID,
	)
}

// StartElasticsearchInstance 启动 Elasticsearch 实例
func (p *AlibabaProvider) StartElasticsearchInstance(ctx context.Context, instanceID string) error {
	return cloudprovider.NewCloudError(
		cloudprovider.ErrUnsupportedOperation,
		"Elasticsearch SDK not integrated - requires elasticsearch client",
		instanceID,
	)
}

// StopElasticsearchInstance 停止 Elasticsearch 实例
func (p *AlibabaProvider) StopElasticsearchInstance(ctx context.Context, instanceID string) error {
	return cloudprovider.NewCloudError(
		cloudprovider.ErrUnsupportedOperation,
		"Elasticsearch SDK not integrated - requires elasticsearch client",
		instanceID,
	)
}

// ResizeElasticsearchInstance 调整 Elasticsearch 实例规格
func (p *AlibabaProvider) ResizeElasticsearchInstance(ctx context.Context, instanceID string, spec cloudprovider.ElasticsearchSpec) error {
	return cloudprovider.NewCloudError(
		cloudprovider.ErrUnsupportedOperation,
		"Elasticsearch SDK not integrated - requires elasticsearch client",
		instanceID,
	)
}

// ListElasticsearchInstances 列出 Elasticsearch 实例
func (p *AlibabaProvider) ListElasticsearchInstances(ctx context.Context, filter cloudprovider.ElasticsearchFilter) ([]*cloudprovider.ElasticsearchInstance, error) {
	// Return empty list with no error for listing operations
	// This allows frontend to display empty state rather than error
	return []*cloudprovider.ElasticsearchInstance{}, nil
}