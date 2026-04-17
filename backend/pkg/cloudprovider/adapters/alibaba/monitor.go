package alibaba

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/opencmp/opencmp/pkg/cloudprovider"
)

// 监控指标名称映射
var metricNameMapping = map[string]string{
	"cpu_usage":          "CPUUtilization",
	"cpu_utilization":    "CPUUtilization",
	"memory_usage":       "memory_usedutilization",
	"memory_utilization": "memory_usedutilization",
	"disk_usage":         "diskusage_utilization",
	"disk_utilization":   "diskusage_utilization",
	"network_in":         "VPC_PublicIP_InputRate",
	"network_in_rate":    "VPC_PublicIP_InputRate",
	"network_out":        "VPC_PublicIP_OutputRate",
	"network_out_rate":   "VPC_PublicIP_OutputRate",
}

// GetInstanceMetrics 获取实例监控指标
func (p *AlibabaProvider) GetInstanceMetrics(ctx context.Context, instanceID string, metrics []string, period int, startTime, endTime string) (*cloudprovider.InstanceMetrics, error) {
	// 创建云监控客户端
	cmsClient, err := cms.NewClientWithAccessKey(p.regionID, p.config.Credentials["access_key_id"], p.config.Credentials["access_key_secret"])
	if err != nil {
		return nil, cloudprovider.NewCloudError(cloudprovider.ErrProviderError, "failed to create cms client", err.Error())
	}

	result := &cloudprovider.InstanceMetrics{
		InstanceID: instanceID,
		Metrics:    make(map[string]cloudprovider.MetricData),
	}

	// 查询每个指标
	for _, metricName := range metrics {
		aliyunMetricName := metricNameMapping[metricName]
		if aliyunMetricName == "" {
			aliyunMetricName = metricName // 使用原始名称
		}

		request := cms.CreateDescribeMetricListRequest()
		request.Scheme = "https"
		request.Namespace = "acs_ecs_dashboard"
		request.MetricName = aliyunMetricName
		request.Period = strconv.Itoa(period)
		request.StartTime = startTime
		request.EndTime = endTime
		request.Length = "100"
		request.Dimensions = fmt.Sprintf(`[{"instanceId":"%s"}]`, instanceID)

		response, err := cmsClient.DescribeMetricList(request)
		if err != nil {
			// 权限不足或指标不存在，返回空值
			result.Metrics[metricName] = cloudprovider.MetricData{
				Name:      metricName,
				Value:     0,
				Unit:      getMetricUnit(aliyunMetricName),
				Timestamp: time.Now().Unix(),
			}
			continue
		}

		// 解析响应数据
		datapoints := parseDatapoints(response.Datapoints)
		if len(datapoints) > 0 {
			// 获取最新的数据点
			latestPoint := datapoints[len(datapoints)-1]
			value := 0.0
			if v, ok := latestPoint["Value"]; ok {
				switch vv := v.(type) {
				case float64:
					value = vv
				case int:
					value = float64(vv)
				case string:
					value, _ = strconv.ParseFloat(vv, 64)
				}
			}

			// 构建历史数据点
			points := make([]cloudprovider.MetricPoint, 0, len(datapoints))
			for _, dp := range datapoints {
				ptValue := 0.0
				if v, ok := dp["Value"]; ok {
					switch vv := v.(type) {
					case float64:
						ptValue = vv
					case int:
						ptValue = float64(vv)
					case string:
						ptValue, _ = strconv.ParseFloat(vv, 64)
					}
				}
				ts := int64(0)
				if t, ok := dp["timestamp"].(float64); ok {
					ts = int64(t)
				}
				points = append(points, cloudprovider.MetricPoint{
					Timestamp: ts,
					Value:     ptValue,
				})
			}

			// 获取timestamp
			ts := time.Now().Unix()
			if t, ok := latestPoint["timestamp"].(float64); ok {
				ts = int64(t)
			}

			result.Metrics[metricName] = cloudprovider.MetricData{
				Name:      metricName,
				Value:     value,
				Unit:      getMetricUnit(aliyunMetricName),
				Timestamp: ts,
				Points:    points,
			}
		} else {
			// 无数据，返回默认值
			result.Metrics[metricName] = cloudprovider.MetricData{
				Name:      metricName,
				Value:     0,
				Unit:      getMetricUnit(aliyunMetricName),
				Timestamp: time.Now().Unix(),
			}
		}
	}

	return result, nil
}

// ListInstanceMetrics 批量获取实例指标
func (p *AlibabaProvider) ListInstanceMetrics(ctx context.Context, instanceIDs []string, metrics []string) ([]*cloudprovider.InstanceMetricData, error) {
	// 创建云监控客户端
	cmsClient, err := cms.NewClientWithAccessKey(p.regionID, p.config.Credentials["access_key_id"], p.config.Credentials["access_key_secret"])
	if err != nil {
		return nil, cloudprovider.NewCloudError(cloudprovider.ErrProviderError, "failed to create cms client", err.Error())
	}

	results := make([]*cloudprovider.InstanceMetricData, 0, len(instanceIDs))

	// 为每个实例查询指标
	for _, instanceID := range instanceIDs {
		data := &cloudprovider.InstanceMetricData{
			InstanceID: instanceID,
			Metrics:    make(map[string]float64),
			Timestamp:  time.Now().Unix(),
		}

		for _, metricName := range metrics {
			aliyunMetricName := metricNameMapping[metricName]
			if aliyunMetricName == "" {
				aliyunMetricName = metricName
			}

			request := cms.CreateDescribeMetricLastRequest()
			request.Scheme = "https"
			request.Namespace = "acs_ecs_dashboard"
			request.MetricName = aliyunMetricName
			request.Dimensions = fmt.Sprintf(`[{"instanceId":"%s"}]`, instanceID)

			response, err := cmsClient.DescribeMetricLast(request)
			if err != nil {
				data.Metrics[metricName] = 0
				continue
			}

			datapoints := parseDatapoints(response.Datapoints)
			if len(datapoints) > 0 {
				latestPoint := datapoints[len(datapoints)-1]
				value := 0.0
				if v, ok := latestPoint["Value"]; ok {
					switch vv := v.(type) {
					case float64:
						value = vv
					case int:
						value = float64(vv)
					case string:
						value, _ = strconv.ParseFloat(vv, 64)
					}
				}
				data.Metrics[metricName] = value
			} else {
				data.Metrics[metricName] = 0
			}
		}

		results = append(results, data)
	}

	return results, nil
}

// parseDatapoints 解析阿里云返回的Datapoints JSON字符串
func parseDatapoints(datapointsStr string) []map[string]interface{} {
	if datapointsStr == "" {
		return nil
	}

	var datapoints []map[string]interface{}
	if err := json.Unmarshal([]byte(datapointsStr), &datapoints); err != nil {
		return nil
	}

	return datapoints
}

// getMetricUnit 获取指标单位
func getMetricUnit(metricName string) string {
	switch metricName {
	case "VPC_PublicIP_InputRate", "VPC_PublicIP_OutputRate":
		return "KB/s"
	default:
		return "percent"
	}
}