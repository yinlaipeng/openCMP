package cloudprovider

import "context"

// IMonitor 监控接口
type IMonitor interface {
	// 获取实例监控指标
	GetInstanceMetrics(ctx context.Context, instanceID string, metrics []string, period int, startTime, endTime string) (*InstanceMetrics, error)
	// 获取实例列表的最新指标
	ListInstanceMetrics(ctx context.Context, instanceIDs []string, metrics []string) ([]*InstanceMetricData, error)
}

// InstanceMetrics 实例监控指标集合
type InstanceMetrics struct {
	InstanceID string                 `json:"instance_id"`
	Metrics    map[string]MetricData  `json:"metrics"`
}

// MetricData 单个指标数据
type MetricData struct {
	Name      string      `json:"name"`
	Value     float64     `json:"value"`
	Unit      string      `json:"unit"`
	Timestamp int64       `json:"timestamp"`
	Points    []MetricPoint `json:"points,omitempty"` // 历史数据点
}

// MetricPoint 指标数据点
type MetricPoint struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}

// InstanceMetricData 实例指标数据
type InstanceMetricData struct {
	InstanceID string            `json:"instance_id"`
	Metrics    map[string]float64 `json:"metrics"`
	Timestamp  int64             `json:"timestamp"`
}

// MetricFilter 指标查询过滤条件
type MetricFilter struct {
	InstanceID string   `json:"instance_id"`
	Metrics    []string `json:"metrics"`    // cpu_utilization, memory_utilization, disk_utilization, network_in_rate, network_out_rate
	Period     int      `json:"period"`     // 采集周期(秒): 60, 300, 900
	StartTime  string   `json:"start_time"` // 格式: 2023-01-01T00:00:00Z
	EndTime    string   `json:"end_time"`   // 格式: 2023-01-01T23:59:59Z
}