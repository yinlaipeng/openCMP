package alibaba

import (
	"context"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/opencmp/opencmp/pkg/cloudprovider"
)

// ListRDSInstances 列出 RDS 实例
func (p *AlibabaProvider) ListRDSInstances(ctx context.Context, filter cloudprovider.RDSFilter) ([]*cloudprovider.RDSInstance, error) {
	if p.rdsClient == nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrUnsupportedOperation,
			"RDS client not initialized",
			"",
		)
	}

	request := rds.CreateDescribeDBInstancesRequest()
	request.Scheme = "https"
	request.RegionId = p.regionID

	if filter.InstanceID != "" {
		request.DBInstanceId = filter.InstanceID
	}
	if filter.Engine != "" {
		request.Engine = filter.Engine
	}
	if filter.Status != "" {
		request.DBInstanceStatus = filter.Status
	}

	response, err := p.rdsClient.DescribeDBInstances(request)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to describe rds instances",
			err.Error(),
		)
	}

	var instances []*cloudprovider.RDSInstance
	for _, db := range response.Items.DBInstance {
		// DBInstanceMemory is in MB, convert to GB for StorageSize
		storageSize := db.DBInstanceMemory / 1024
		port, _ := strconv.Atoi(db.DBInstanceNetType)

		instances = append(instances, &cloudprovider.RDSInstance{
			ID:             db.DBInstanceId,
			Name:           db.DBInstanceDescription,
			Engine:         db.Engine,
			EngineVersion:  db.EngineVersion,
			InstanceType:   db.DBInstanceClass,
			StorageSize:    storageSize,
			StorageType:    db.DBInstanceStorageType,
			Status:         db.DBInstanceStatus,
			VPCID:          db.VpcId,
			SubnetID:       db.VSwitchId,
			Endpoint:       db.ConnectionString,
			Port:           port,
			MasterUsername: "", // Not returned by DescribeDBInstances
			ZoneID:         db.ZoneId,
			CreatedAt:      time.Now(), // Use current time as approximation
		})
	}

	return instances, nil
}

// CreateRDSInstance 创建 RDS 实例
func (p *AlibabaProvider) CreateRDSInstance(ctx context.Context, config cloudprovider.RDSConfig) (*cloudprovider.RDSInstance, error) {
	if p.rdsClient == nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrUnsupportedOperation,
			"RDS client not initialized",
			"",
		)
	}

	request := rds.CreateCreateDBInstanceRequest()
	request.Scheme = "https"
	request.RegionId = p.regionID
	request.DBInstanceDescription = config.Name
	request.Engine = config.Engine
	request.EngineVersion = config.EngineVersion
	request.DBInstanceClass = config.InstanceType
	request.DBInstanceStorage = requests.NewInteger(config.StorageSize)
	request.DBInstanceStorageType = config.StorageType
	request.VPCId = config.VPCID
	request.VSwitchId = config.SubnetID
	request.ZoneId = config.ZoneID
	// Note: MasterUsername/MasterPassword needs to be set via separate API after creation
	// or using CreateAccount API

	response, err := p.rdsClient.CreateDBInstance(request)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to create rds instance",
			err.Error(),
		)
	}

	return &cloudprovider.RDSInstance{
		ID:            response.DBInstanceId,
		Name:          config.Name,
		Engine:        config.Engine,
		EngineVersion: config.EngineVersion,
		InstanceType:  config.InstanceType,
		StorageSize:   config.StorageSize,
		StorageType:   config.StorageType,
		Status:        "Creating",
		VPCID:         config.VPCID,
		SubnetID:      config.SubnetID,
		ZoneID:        config.ZoneID,
		CreatedAt:     time.Now(),
	}, nil
}

// DeleteRDSInstance 删除 RDS 实例
func (p *AlibabaProvider) DeleteRDSInstance(ctx context.Context, instanceID string) error {
	if p.rdsClient == nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrUnsupportedOperation,
			"RDS client not initialized",
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
			"failed to delete rds instance",
			err.Error(),
		)
	}

	return nil
}

// StartRDSInstance 启动 RDS 实例
func (p *AlibabaProvider) StartRDSInstance(ctx context.Context, instanceID string) error {
	if p.rdsClient == nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrUnsupportedOperation,
			"RDS client not initialized",
			"",
		)
	}

	request := rds.CreateStartDBInstanceRequest()
	request.Scheme = "https"
	request.DBInstanceId = instanceID

	_, err := p.rdsClient.StartDBInstance(request)
	if err != nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to start rds instance",
			err.Error(),
		)
	}

	return nil
}

// StopRDSInstance 停止 RDS 实例
func (p *AlibabaProvider) StopRDSInstance(ctx context.Context, instanceID string) error {
	if p.rdsClient == nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrUnsupportedOperation,
			"RDS client not initialized",
			"",
		)
	}

	request := rds.CreateStopDBInstanceRequest()
	request.Scheme = "https"
	request.DBInstanceId = instanceID

	_, err := p.rdsClient.StopDBInstance(request)
	if err != nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to stop rds instance",
			err.Error(),
		)
	}

	return nil
}

// RebootRDSInstance 重启 RDS 实例
func (p *AlibabaProvider) RebootRDSInstance(ctx context.Context, instanceID string) error {
	if p.rdsClient == nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrUnsupportedOperation,
			"RDS client not initialized",
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
			"failed to reboot rds instance",
			err.Error(),
		)
	}

	return nil
}

// ResizeRDSInstance 调整 RDS 实例规格
func (p *AlibabaProvider) ResizeRDSInstance(ctx context.Context, instanceID string, spec cloudprovider.RDSpec) error {
	if p.rdsClient == nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrUnsupportedOperation,
			"RDS client not initialized",
			"",
		)
	}

	request := rds.CreateModifyDBInstanceSpecRequest()
	request.Scheme = "https"
	request.DBInstanceId = instanceID
	request.DBInstanceClass = spec.InstanceType
	if spec.StorageSize > 0 {
		request.DBInstanceStorage = requests.NewInteger(spec.StorageSize)
	}

	_, err := p.rdsClient.ModifyDBInstanceSpec(request)
	if err != nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to resize rds instance",
			err.Error(),
		)
	}

	return nil
}

// CreateRDSBackup 创建 RDS 备份
func (p *AlibabaProvider) CreateRDSBackup(ctx context.Context, instanceID string, name string) (*cloudprovider.RDSBackup, error) {
	if p.rdsClient == nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrUnsupportedOperation,
			"RDS client not initialized",
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
			"failed to create rds backup",
			err.Error(),
		)
	}

	return &cloudprovider.RDSBackup{
		ID:         response.BackupJobId,
		Name:       name,
		InstanceID: instanceID,
		Status:     "Creating",
		StartTime:  time.Now(),
	}, nil
}

// ListRDSBackups 列出 RDS 备份
func (p *AlibabaProvider) ListRDSBackups(ctx context.Context, instanceID string) ([]*cloudprovider.RDSBackup, error) {
	if p.rdsClient == nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrUnsupportedOperation,
			"RDS client not initialized",
			"",
		)
	}

	request := rds.CreateDescribeBackupsRequest()
	request.Scheme = "https"
	request.DBInstanceId = instanceID

	response, err := p.rdsClient.DescribeBackups(request)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to describe rds backups",
			err.Error(),
		)
	}

	var backups []*cloudprovider.RDSBackup
	for _, b := range response.Items.Backup {
		// BackupSize is already int64, no conversion needed
		startTime, _ := time.Parse(time.RFC3339, b.BackupStartTime)
		endTime, _ := time.Parse(time.RFC3339, b.BackupEndTime)

		backups = append(backups, &cloudprovider.RDSBackup{
			ID:         b.BackupId,
			Name:       b.BackupDBNames, // Use BackupDBNames as name
			InstanceID: instanceID,
			Status:     b.BackupStatus,
			Size:       b.BackupSize,
			StartTime:  startTime,
			EndTime:    endTime,
		})
	}

	return backups, nil
}

// RestoreRDSFromBackup 从备份恢复 RDS
func (p *AlibabaProvider) RestoreRDSFromBackup(ctx context.Context, backupID string, config cloudprovider.RDSConfig) (*cloudprovider.RDSInstance, error) {
	// 暂不支持，返回错误
	return nil, cloudprovider.NewCloudError(
		cloudprovider.ErrUnsupportedOperation,
		"RestoreRDSFromBackup not fully implemented",
		backupID,
	)
}