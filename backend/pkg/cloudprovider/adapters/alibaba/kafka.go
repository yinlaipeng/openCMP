package alibaba

import (
	"context"

	"github.com/opencmp/opencmp/pkg/cloudprovider"
)

// CreateKafkaInstance 创建 Kafka 实例
// Note: Requires alikafka SDK integration - currently returns stub
func (p *AlibabaProvider) CreateKafkaInstance(ctx context.Context, config cloudprovider.KafkaConfig) (*cloudprovider.KafkaInstance, error) {
	return nil, cloudprovider.NewCloudError(
		cloudprovider.ErrUnsupportedOperation,
		"Kafka SDK not integrated - requires alikafka client",
		"",
	)
}

// DeleteKafkaInstance 删除 Kafka 实例
func (p *AlibabaProvider) DeleteKafkaInstance(ctx context.Context, instanceID string) error {
	return cloudprovider.NewCloudError(
		cloudprovider.ErrUnsupportedOperation,
		"Kafka SDK not integrated - requires alikafka client",
		instanceID,
	)
}

// StartKafkaInstance 启动 Kafka 实例
func (p *AlibabaProvider) StartKafkaInstance(ctx context.Context, instanceID string) error {
	return cloudprovider.NewCloudError(
		cloudprovider.ErrUnsupportedOperation,
		"Kafka SDK not integrated - requires alikafka client",
		instanceID,
	)
}

// StopKafkaInstance 停止 Kafka 实例
func (p *AlibabaProvider) StopKafkaInstance(ctx context.Context, instanceID string) error {
	return cloudprovider.NewCloudError(
		cloudprovider.ErrUnsupportedOperation,
		"Kafka SDK not integrated - requires alikafka client",
		instanceID,
	)
}

// ResizeKafkaInstance 调整 Kafka 实例规格
func (p *AlibabaProvider) ResizeKafkaInstance(ctx context.Context, instanceID string, spec cloudprovider.KafkaSpec) error {
	return cloudprovider.NewCloudError(
		cloudprovider.ErrUnsupportedOperation,
		"Kafka SDK not integrated - requires alikafka client",
		instanceID,
	)
}

// ListKafkaInstances 列出 Kafka 实例
func (p *AlibabaProvider) ListKafkaInstances(ctx context.Context, filter cloudprovider.KafkaFilter) ([]*cloudprovider.KafkaInstance, error) {
	// Return empty list with no error for listing operations
	// This allows frontend to display empty state rather than error
	return []*cloudprovider.KafkaInstance{}, nil
}