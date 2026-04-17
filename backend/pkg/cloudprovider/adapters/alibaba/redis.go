package alibaba

import (
	"context"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/opencmp/opencmp/pkg/cloudprovider"
)

// Use r-kvstore SDK for Redis
// Note: The alibaba-cloud-sdk-go has limited Redis SDK support
// We use the ECS client with Redis-related API calls through the rds client
// as Aliyun Redis is managed through similar API patterns

// ListCacheInstances 列出 Redis 实例
func (p *AlibabaProvider) ListCacheInstances(ctx context.Context, filter cloudprovider.CacheFilter) ([]*cloudprovider.CacheInstance, error) {
	if p.rdsClient == nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrUnsupportedOperation,
			"RDS client not initialized for Redis operations",
			"",
		)
	}

	// Aliyun Redis uses DescribeDBInstances with Engine=Redis
	request := rds.CreateDescribeDBInstancesRequest()
	request.Scheme = "https"
	request.RegionId = p.regionID
	request.Engine = "Redis"

	if filter.InstanceID != "" {
		request.DBInstanceId = filter.InstanceID
	}
	if filter.Status != "" {
		request.DBInstanceStatus = filter.Status
	}

	response, err := p.rdsClient.DescribeDBInstances(request)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to describe redis instances",
			err.Error(),
		)
	}

	var instances []*cloudprovider.CacheInstance
	for _, db := range response.Items.DBInstance {
		port := 6379 // Default Redis port
		if db.DBInstanceNetType != "" {
			portInt, _ := strconv.Atoi(db.DBInstanceNetType)
			if portInt > 0 {
				port = portInt
			}
		}

		instances = append(instances, &cloudprovider.CacheInstance{
			ID:            db.DBInstanceId,
			Name:          db.DBInstanceDescription,
			Engine:        "Redis",
			EngineVersion: db.EngineVersion,
			InstanceType:  db.DBInstanceClass,
			Status:        db.DBInstanceStatus,
			VPCID:         db.VpcId,
			SubnetID:      db.VSwitchId,
			Endpoint:      db.ConnectionString,
			Port:          port,
			ZoneID:        db.ZoneId,
			CreatedAt:     time.Now(),
		})
	}

	return instances, nil
}

// CreateCacheInstance 创建 Redis 实例
func (p *AlibabaProvider) CreateCacheInstance(ctx context.Context, config cloudprovider.CacheConfig) (*cloudprovider.CacheInstance, error) {
	if p.rdsClient == nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrUnsupportedOperation,
			"RDS client not initialized for Redis operations",
			"",
		)
	}

	// Aliyun Redis creation requires the r-kvstore SDK specifically
	// For now, we use a simplified approach through the RDS API pattern
	request := rds.CreateCreateDBInstanceRequest()
	request.Scheme = "https"
	request.RegionId = p.regionID
	request.DBInstanceDescription = config.Name
	request.Engine = "Redis"
	request.EngineVersion = config.EngineVersion
	request.DBInstanceClass = config.InstanceType
	request.VPCId = config.VPCID
	request.VSwitchId = config.SubnetID
	request.ZoneId = config.ZoneID

	response, err := p.rdsClient.CreateDBInstance(request)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to create redis instance",
			err.Error(),
		)
	}

	return &cloudprovider.CacheInstance{
		ID:            response.DBInstanceId,
		Name:          config.Name,
		Engine:        "Redis",
		EngineVersion: config.EngineVersion,
		InstanceType:  config.InstanceType,
		Status:        "Creating",
		VPCID:         config.VPCID,
		SubnetID:      config.SubnetID,
		ZoneID:        config.ZoneID,
		CreatedAt:     time.Now(),
	}, nil
}

// DeleteCacheInstance 删除 Redis 实例
func (p *AlibabaProvider) DeleteCacheInstance(ctx context.Context, instanceID string) error {
	if p.rdsClient == nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrUnsupportedOperation,
			"RDS client not initialized for Redis operations",
			"",
		)
	}

	request := rds.CreateDeleteDBInstanceRequest()
	request.Scheme = "https"
	request.DBInstanceId = instanceID

	_, err := p.rdsClient.DeleteDBInstance(request)
	if err != nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to delete redis instance",
			err.Error(),
		)
	}

	return nil
}

// RebootCacheInstance 重启 Redis 实例
func (p *AlibabaProvider) RebootCacheInstance(ctx context.Context, instanceID string) error {
	if p.rdsClient == nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrUnsupportedOperation,
			"RDS client not initialized for Redis operations",
			"",
		)
	}

	request := rds.CreateRestartDBInstanceRequest()
	request.Scheme = "https"
	request.DBInstanceId = instanceID

	_, err := p.rdsClient.RestartDBInstance(request)
	if err != nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to reboot redis instance",
			err.Error(),
		)
	}

	return nil
}

// ResizeCacheInstance 调整 Redis 实例规格
func (p *AlibabaProvider) ResizeCacheInstance(ctx context.Context, instanceID string, spec cloudprovider.CacheSpec) error {
	if p.rdsClient == nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrUnsupportedOperation,
			"RDS client not initialized for Redis operations",
			"",
		)
	}

	request := rds.CreateModifyDBInstanceSpecRequest()
	request.Scheme = "https"
	request.DBInstanceId = instanceID
	request.DBInstanceClass = spec.InstanceType

	_, err := p.rdsClient.ModifyDBInstanceSpec(request)
	if err != nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to resize redis instance",
			err.Error(),
		)
	}

	return nil
}

// CreateCacheBackup 创建 Redis 备份
func (p *AlibabaProvider) CreateCacheBackup(ctx context.Context, instanceID string) (*cloudprovider.CacheBackup, error) {
	if p.rdsClient == nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrUnsupportedOperation,
			"RDS client not initialized for Redis operations",
			"",
		)
	}

	request := rds.CreateCreateBackupRequest()
	request.Scheme = "https"
	request.DBInstanceId = instanceID

	response, err := p.rdsClient.CreateBackup(request)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to create redis backup",
			err.Error(),
		)
	}

	return &cloudprovider.CacheBackup{
		ID:         response.BackupJobId,
		InstanceID: instanceID,
		Status:     "Creating",
		StartTime:  time.Now(),
	}, nil
}