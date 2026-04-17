package cloudprovider

import (
	"context"
)

// IKafka Kafka消息队列管理接口
type IKafka interface {
	CreateKafkaInstance(ctx context.Context, config KafkaConfig) (*KafkaInstance, error)
	DeleteKafkaInstance(ctx context.Context, instanceID string) error
	StartKafkaInstance(ctx context.Context, instanceID string) error
	StopKafkaInstance(ctx context.Context, instanceID string) error
	ResizeKafkaInstance(ctx context.Context, instanceID string, spec KafkaSpec) error
	ListKafkaInstances(ctx context.Context, filter KafkaFilter) ([]*KafkaInstance, error)
}

// IElasticsearch Elasticsearch管理接口
type IElasticsearch interface {
	CreateElasticsearchInstance(ctx context.Context, config ElasticsearchConfig) (*ElasticsearchInstance, error)
	DeleteElasticsearchInstance(ctx context.Context, instanceID string) error
	StartElasticsearchInstance(ctx context.Context, instanceID string) error
	StopElasticsearchInstance(ctx context.Context, instanceID string) error
	ResizeElasticsearchInstance(ctx context.Context, instanceID string, spec ElasticsearchSpec) error
	ListElasticsearchInstances(ctx context.Context, filter ElasticsearchFilter) ([]*ElasticsearchInstance, error)
}

// IMiddleware 中间件服务总接口
type IMiddleware interface {
	IKafka
	IElasticsearch
}

// KafkaSpec Kafka规格
type KafkaSpec struct {
	Storage   int `json:"storage"`    // GB
	Bandwidth int `json:"bandwidth"`  // MB/s
}

// ElasticsearchSpec Elasticsearch规格
type ElasticsearchSpec struct {
	NodeCount int    `json:"node_count"`
	Storage   int    `json:"storage"` // GB
}