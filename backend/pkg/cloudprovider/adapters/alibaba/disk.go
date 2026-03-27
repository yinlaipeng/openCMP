package alibaba

import (
	"context"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/opencmp/opencmp/pkg/cloudprovider"
)

// CreateDisk 创建磁盘
func (p *AlibabaProvider) CreateDisk(ctx context.Context, config cloudprovider.DiskConfig) (*cloudprovider.Disk, error) {
	request := ecs.CreateCreateDiskRequest()
	request.Scheme = "https"
	request.RegionId = p.regionID
	request.Size = requests.NewInteger(config.Size)
	request.DiskName = config.Name
	request.ZoneId = config.ZoneID

	if config.Type != "" {
		request.DiskCategory = config.Type
	}

	response, err := p.ecsClient.CreateDisk(request)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to create disk",
			err.Error(),
		)
	}

	return &cloudprovider.Disk{
		ID:        response.DiskId,
		Name:      config.Name,
		Size:      config.Size,
		Type:      config.Type,
		Status:    "available",
		ZoneID:    config.ZoneID,
		CreatedAt: time.Now(),
	}, nil
}

// DeleteDisk 删除磁盘
func (p *AlibabaProvider) DeleteDisk(ctx context.Context, diskID string) error {
	request := ecs.CreateDeleteDiskRequest()
	request.Scheme = "https"
	request.DiskId = diskID

	_, err := p.ecsClient.DeleteDisk(request)
	if err != nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to delete disk",
			err.Error(),
		)
	}

	return nil
}

// AttachDisk 挂载磁盘
func (p *AlibabaProvider) AttachDisk(ctx context.Context, diskID, vmID string) error {
	request := ecs.CreateAttachDiskRequest()
	request.Scheme = "https"
	request.DiskId = diskID
	request.InstanceId = vmID

	_, err := p.ecsClient.AttachDisk(request)
	if err != nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to attach disk",
			err.Error(),
		)
	}

	return nil
}

// DetachDisk 卸载磁盘
func (p *AlibabaProvider) DetachDisk(ctx context.Context, diskID string) error {
	request := ecs.CreateDetachDiskRequest()
	request.Scheme = "https"
	request.DiskId = diskID

	_, err := p.ecsClient.DetachDisk(request)
	if err != nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to detach disk",
			err.Error(),
		)
	}

	return nil
}

// ResizeDisk 扩容磁盘
func (p *AlibabaProvider) ResizeDisk(ctx context.Context, diskID string, sizeGB int) error {
	request := ecs.CreateResizeDiskRequest()
	request.Scheme = "https"
	request.DiskId = diskID
	request.NewSize = requests.NewInteger(sizeGB)

	_, err := p.ecsClient.ResizeDisk(request)
	if err != nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to resize disk",
			err.Error(),
		)
	}

	return nil
}

// CreateSnapshot 创建快照
func (p *AlibabaProvider) CreateSnapshot(ctx context.Context, diskID string, name string) (*cloudprovider.Snapshot, error) {
	request := ecs.CreateCreateSnapshotRequest()
	request.Scheme = "https"
	request.DiskId = diskID
	request.SnapshotName = name

	response, err := p.ecsClient.CreateSnapshot(request)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to create snapshot",
			err.Error(),
		)
	}

	return &cloudprovider.Snapshot{
		ID:     response.SnapshotId,
		Name:   name,
		DiskID: diskID,
		Status: "creating",
	}, nil
}

// DeleteSnapshot 删除快照
func (p *AlibabaProvider) DeleteSnapshot(ctx context.Context, snapshotID string) error {
	request := ecs.CreateDeleteSnapshotRequest()
	request.Scheme = "https"
	request.SnapshotId = snapshotID

	_, err := p.ecsClient.DeleteSnapshot(request)
	if err != nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to delete snapshot",
			err.Error(),
		)
	}

	return nil
}

// ListDisks 列出磁盘
func (p *AlibabaProvider) ListDisks(ctx context.Context, filter cloudprovider.DiskFilter) ([]*cloudprovider.Disk, error) {
	request := ecs.CreateDescribeDisksRequest()
	request.Scheme = "https"
	request.RegionId = p.regionID

	if filter.DiskID != "" {
		request.DiskIds = filter.DiskID
	}
	if filter.VMID != "" {
		request.InstanceId = filter.VMID
	}

	response, err := p.ecsClient.DescribeDisks(request)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to describe disks",
			err.Error(),
		)
	}

	var disks []*cloudprovider.Disk
	for _, disk := range response.Disks.Disk {
		disks = append(disks, &cloudprovider.Disk{
			ID:     disk.DiskId,
			Name:   disk.DiskName,
			Size:   disk.Size,
			Type:   disk.Category,
			Status: disk.Status,
			VMID:   disk.InstanceId,
			ZoneID: disk.ZoneId,
		})
	}

	return disks, nil
}

// ListSnapshots 列出快照
func (p *AlibabaProvider) ListSnapshots(ctx context.Context, filter cloudprovider.SnapshotFilter) ([]*cloudprovider.Snapshot, error) {
	request := ecs.CreateDescribeSnapshotsRequest()
	request.Scheme = "https"

	if filter.SnapshotID != "" {
		request.SnapshotIds = filter.SnapshotID
	}
	if filter.DiskID != "" {
		request.DiskId = filter.DiskID
	}

	response, err := p.ecsClient.DescribeSnapshots(request)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to describe snapshots",
			err.Error(),
		)
	}

	var snapshots []*cloudprovider.Snapshot
	for _, snapshot := range response.Snapshots.Snapshot {
		size, _ := strconv.Atoi(snapshot.SourceDiskSize)
		snapshots = append(snapshots, &cloudprovider.Snapshot{
			ID:     snapshot.SnapshotId,
			Name:   snapshot.SnapshotName,
			DiskID: snapshot.SourceDiskId,
			Size:   size,
			Status: snapshot.Status,
		})
	}

	return snapshots, nil
}
