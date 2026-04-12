package tencent

import (
	"context"
	"time"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"github.com/opencmp/opencmp/pkg/cloudprovider"
)

// CreateVM 创建虚拟机
func (p *TencentProvider) CreateVM(ctx context.Context, config cloudprovider.VMCreateConfig) (*cloudprovider.VirtualMachine, error) {
	request := cvm.NewRunInstancesRequest()
	request.InstanceType = &config.InstanceType
	request.ImageId = &config.ImageID

	// 设置子网和安全组
	if config.SubnetID != "" {
		request.VirtualPrivateCloud = &cvm.VirtualPrivateCloud{
			VpcId:    &config.VPCID,
			SubnetId: &config.SubnetID,
		}
	}

	if len(config.SecurityGroups) > 0 {
		var sgIds []*string
		for _, sg := range config.SecurityGroups {
			sgIds = append(sgIds, &sg)
		}
		request.SecurityGroupIds = sgIds
	}

	request.InstanceName = &config.Name

	// 设置系统盘
	request.SystemDisk = &cvm.SystemDisk{
		DiskSize: common.Int64Ptr(int64(config.DiskSize)),
		DiskType: common.StringPtr("CLOUD_PREMIUM"),
	}

	response, err := p.cvmClient.RunInstances(request)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to create instance",
			err.Error(),
		)
	}

	if len(response.Response.InstanceIdSet) == 0 {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"no instance created",
			"",
		)
	}

	instanceID := response.Response.InstanceIdSet[0]

	return &cloudprovider.VirtualMachine{
		ID:           *instanceID,
		Name:         config.Name,
		Status:       cloudprovider.VMStatusPending,
		InstanceType: config.InstanceType,
		ImageID:      config.ImageID,
		VPCID:        config.VPCID,
		SubnetID:     config.SubnetID,
		CreatedAt:    time.Now(),
		RegionID:     p.regionID,
	}, nil
}

// DeleteVM 删除虚拟机
func (p *TencentProvider) DeleteVM(ctx context.Context, vmID string) error {
	request := cvm.NewTerminateInstancesRequest()
	request.InstanceIds = []*string{&vmID}

	_, err := p.cvmClient.TerminateInstances(request)
	if err != nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to terminate instance",
			err.Error(),
		)
	}

	return nil
}

// StartVM 启动虚拟机
func (p *TencentProvider) StartVM(ctx context.Context, vmID string) error {
	request := cvm.NewStartInstancesRequest()
	request.InstanceIds = []*string{&vmID}

	_, err := p.cvmClient.StartInstances(request)
	if err != nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to start instance",
			err.Error(),
		)
	}

	return nil
}

// StopVM 停止虚拟机
func (p *TencentProvider) StopVM(ctx context.Context, vmID string) error {
	request := cvm.NewStopInstancesRequest()
	request.InstanceIds = []*string{&vmID}
	request.ForceStop = common.BoolPtr(false)

	_, err := p.cvmClient.StopInstances(request)
	if err != nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to stop instance",
			err.Error(),
		)
	}

	return nil
}

// RebootVM 重启虚拟机
func (p *TencentProvider) RebootVM(ctx context.Context, vmID string) error {
	request := cvm.NewRebootInstancesRequest()
	request.InstanceIds = []*string{&vmID}

	_, err := p.cvmClient.RebootInstances(request)
	if err != nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to reboot instance",
			err.Error(),
		)
	}

	return nil
}

// GetVMStatus 获取虚拟机状态
func (p *TencentProvider) GetVMStatus(ctx context.Context, vmID string) (*cloudprovider.VMStatus, error) {
	request := cvm.NewDescribeInstancesStatusRequest()
	request.InstanceIds = []*string{&vmID}

	response, err := p.cvmClient.DescribeInstancesStatus(request)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to describe instance status",
			err.Error(),
		)
	}

	if len(response.Response.InstanceStatusSet) == 0 {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrResourceNotFound,
			"instance not found",
			vmID,
		)
	}

	status := response.Response.InstanceStatusSet[0].InstanceState
	vmStatus := convertCVMStatus(*status)

	return &vmStatus, nil
}

// ListVMs 列出虚拟机
func (p *TencentProvider) ListVMs(ctx context.Context, filter cloudprovider.VMListFilter) ([]*cloudprovider.VirtualMachine, error) {
	request := cvm.NewDescribeInstancesRequest()

	if filter.VPCID != "" {
		request.Filters = []*cvm.Filter{
			{
				Name:   common.StringPtr("vpc-id"),
				Values: []*string{&filter.VPCID},
			},
		}
	}

	response, err := p.cvmClient.DescribeInstances(request)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to describe instances",
			err.Error(),
		)
	}

	var vms []*cloudprovider.VirtualMachine
	for _, instance := range response.Response.InstanceSet {
		vmStatus := convertCVMStatus(*instance.InstanceState)

		var privateIP, publicIP string
		if len(instance.PrivateIpAddresses) > 0 {
			privateIP = *instance.PrivateIpAddresses[0]
		}
		if len(instance.PublicIpAddresses) > 0 {
			publicIP = *instance.PublicIpAddresses[0]
		}

		var zoneID string
		if instance.Placement != nil && instance.Placement.Zone != nil {
			zoneID = *instance.Placement.Zone
		}

		vms = append(vms, &cloudprovider.VirtualMachine{
			ID:           *instance.InstanceId,
			Name:         *instance.InstanceName,
			Status:       vmStatus,
			InstanceType: *instance.InstanceType,
			ImageID:      *instance.ImageId,
			PrivateIP:    privateIP,
			PublicIP:     publicIP,
			RegionID:     p.regionID,
			ZoneID:       zoneID,
		})
	}

	return vms, nil
}

// GetVM 获取虚拟机详情
func (p *TencentProvider) GetVM(ctx context.Context, vmID string) (*cloudprovider.VirtualMachine, error) {
	request := cvm.NewDescribeInstancesRequest()
	request.InstanceIds = []*string{&vmID}

	response, err := p.cvmClient.DescribeInstances(request)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to describe instance",
			err.Error(),
		)
	}

	if len(response.Response.InstanceSet) == 0 {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrResourceNotFound,
			"instance not found",
			vmID,
		)
	}

	instance := response.Response.InstanceSet[0]
	vmStatus := convertCVMStatus(*instance.InstanceState)

	var privateIP, publicIP string
	if len(instance.PrivateIpAddresses) > 0 {
		privateIP = *instance.PrivateIpAddresses[0]
	}
	if len(instance.PublicIpAddresses) > 0 {
		publicIP = *instance.PublicIpAddresses[0]
	}

	var zoneID string
	if instance.Placement != nil && instance.Placement.Zone != nil {
		zoneID = *instance.Placement.Zone
	}

	return &cloudprovider.VirtualMachine{
		ID:           *instance.InstanceId,
		Name:         *instance.InstanceName,
		Status:       vmStatus,
		InstanceType: *instance.InstanceType,
		ImageID:      *instance.ImageId,
		PrivateIP:    privateIP,
		PublicIP:     publicIP,
		RegionID:     p.regionID,
		ZoneID:       zoneID,
		CreatedAt:    time.Now(),
	}, nil
}

// ListImages 列出镜像
func (p *TencentProvider) ListImages(ctx context.Context, filter cloudprovider.ImageFilter) ([]*cloudprovider.Image, error) {
	request := cvm.NewDescribeImagesRequest()

	response, err := p.cvmClient.DescribeImages(request)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to describe images",
			err.Error(),
		)
	}

	var images []*cloudprovider.Image
	for _, image := range response.Response.ImageSet {
		images = append(images, &cloudprovider.Image{
			ID:           *image.ImageId,
			Name:         *image.ImageName,
			Description:  *image.ImageDescription,
			OSName:       *image.OsName,
			Architecture: *image.Architecture,
			Status:       *image.ImageState,
		})
	}

	return images, nil
}

// GetImage 获取镜像
func (p *TencentProvider) GetImage(ctx context.Context, imageID string) (*cloudprovider.Image, error) {
	request := cvm.NewDescribeImagesRequest()
	request.ImageIds = []*string{&imageID}

	response, err := p.cvmClient.DescribeImages(request)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to describe image",
			err.Error(),
		)
	}

	if len(response.Response.ImageSet) == 0 {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrResourceNotFound,
			"image not found",
			imageID,
		)
	}

	image := response.Response.ImageSet[0]
	return &cloudprovider.Image{
		ID:           *image.ImageId,
		Name:         *image.ImageName,
		Description:  *image.ImageDescription,
		OSName:       *image.OsName,
		Architecture: *image.Architecture,
		Status:       *image.ImageState,
	}, nil
}

// convertCVMStatus 换腾讯云状态
func convertCVMStatus(status string) cloudprovider.VMStatus {
	switch status {
	case "RUNNING":
		return cloudprovider.VMStatusRunning
	case "STOPPED":
		return cloudprovider.VMStatusStopped
	case "STARTING":
		return cloudprovider.VMStatusStarting
	case "STOPPING":
		return cloudprovider.VMStatusStopping
	default:
		return cloudprovider.VMStatusPending
	}
}