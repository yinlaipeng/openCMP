package alibaba

import (
	"context"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/opencmp/opencmp/pkg/cloudprovider"
)

// CreateVM 创建虚拟机
func (p *AlibabaProvider) CreateVM(ctx context.Context, config cloudprovider.VMCreateConfig) (*cloudprovider.VirtualMachine, error) {
	request := ecs.CreateCreateInstanceRequest()
	request.Scheme = "https"
	request.InstanceType = config.InstanceType
	request.ImageId = config.ImageID
	request.VSwitchId = config.SubnetID
	if len(config.SecurityGroups) > 0 {
		request.SecurityGroupId = config.SecurityGroups[0]
	}
	request.InstanceName = config.Name

	// 创建系统盘
	request.SystemDiskSize = requests.NewInteger(config.DiskSize)

	response, err := p.ecsClient.CreateInstance(request)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to create instance",
			err.Error(),
		)
	}

	instanceID := response.InstanceId

	// 启动实例
	startRequest := ecs.CreateStartInstanceRequest()
	startRequest.Scheme = "https"
	startRequest.InstanceId = instanceID
	_, err = p.ecsClient.StartInstance(startRequest)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to start instance",
			err.Error(),
		)
	}

	return &cloudprovider.VirtualMachine{
		ID:           instanceID,
		Name:         config.Name,
		Status:       cloudprovider.VMStatusStarting,
		InstanceType: config.InstanceType,
		ImageID:      config.ImageID,
		VPCID:        config.VPCID,
		SubnetID:     config.SubnetID,
		CreatedAt:    time.Now(),
		RegionID:     p.regionID,
	}, nil
}

// DeleteVM 删除虚拟机
func (p *AlibabaProvider) DeleteVM(ctx context.Context, vmID string) error {
	request := ecs.CreateDeleteInstanceRequest()
	request.Scheme = "https"
	request.InstanceId = vmID

	_, err := p.ecsClient.DeleteInstance(request)
	if err != nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to delete instance",
			err.Error(),
		)
	}

	return nil
}

// StartVM 启动虚拟机
func (p *AlibabaProvider) StartVM(ctx context.Context, vmID string) error {
	request := ecs.CreateStartInstanceRequest()
	request.Scheme = "https"
	request.InstanceId = vmID

	_, err := p.ecsClient.StartInstance(request)
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
func (p *AlibabaProvider) StopVM(ctx context.Context, vmID string) error {
	request := ecs.CreateStopInstanceRequest()
	request.Scheme = "https"
	request.InstanceId = vmID

	_, err := p.ecsClient.StopInstance(request)
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
func (p *AlibabaProvider) RebootVM(ctx context.Context, vmID string) error {
	request := ecs.CreateRebootInstanceRequest()
	request.Scheme = "https"
	request.InstanceId = vmID

	_, err := p.ecsClient.RebootInstance(request)
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
func (p *AlibabaProvider) GetVMStatus(ctx context.Context, vmID string) (*cloudprovider.VMStatus, error) {
	request := ecs.CreateDescribeInstanceStatusRequest()
	request.Scheme = "https"
	request.InstanceId = &[]string{vmID}

	response, err := p.ecsClient.DescribeInstanceStatus(request)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to get instance status",
			err.Error(),
		)
	}

	if len(response.InstanceStatuses.InstanceStatus) == 0 {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrResourceNotFound,
			"instance not found",
			vmID,
		)
	}

	status := response.InstanceStatuses.InstanceStatus[0].Status
	vmStatus := convertECSStatus(status)

	return &vmStatus, nil
}

// ListVMs 列出虚拟机
func (p *AlibabaProvider) ListVMs(ctx context.Context, filter cloudprovider.VMListFilter) ([]*cloudprovider.VirtualMachine, error) {
	request := ecs.CreateDescribeInstancesRequest()
	request.Scheme = "https"
	request.RegionId = p.regionID

	if filter.VPCID != "" {
		request.VpcId = filter.VPCID
	}

	response, err := p.ecsClient.DescribeInstances(request)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to describe instances",
			err.Error(),
		)
	}

	var vms []*cloudprovider.VirtualMachine
	for _, instance := range response.Instances.Instance {
		vmStatus := convertECSStatus(instance.Status)
		var privateIP, publicIP string
		if len(instance.InnerIpAddress.IpAddress) > 0 {
			privateIP = instance.InnerIpAddress.IpAddress[0]
		}
		if len(instance.PublicIpAddress.IpAddress) > 0 {
			publicIP = instance.PublicIpAddress.IpAddress[0]
		}
		vms = append(vms, &cloudprovider.VirtualMachine{
			ID:           instance.InstanceId,
			Name:         instance.InstanceName,
			Status:       vmStatus,
			InstanceType: instance.InstanceType,
			ImageID:      instance.ImageId,
			PrivateIP:    privateIP,
			PublicIP:     publicIP,
			RegionID:     p.regionID,
			ZoneID:       instance.ZoneId,
		})
	}

	return vms, nil
}

// GetVM 获取虚拟机
func (p *AlibabaProvider) GetVM(ctx context.Context, vmID string) (*cloudprovider.VirtualMachine, error) {
	request := ecs.CreateDescribeInstancesRequest()
	request.Scheme = "https"
	request.RegionId = p.regionID
	request.InstanceIds = "[\"" + vmID + "\"]"

	response, err := p.ecsClient.DescribeInstances(request)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to describe instance",
			err.Error(),
		)
	}

	if len(response.Instances.Instance) == 0 {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrResourceNotFound,
			"instance not found",
			vmID,
		)
	}

	instance := response.Instances.Instance[0]
	vmStatus := convertECSStatus(instance.Status)
	var privateIP, publicIP string
	if len(instance.InnerIpAddress.IpAddress) > 0 {
		privateIP = instance.InnerIpAddress.IpAddress[0]
	}
	if len(instance.PublicIpAddress.IpAddress) > 0 {
		publicIP = instance.PublicIpAddress.IpAddress[0]
	}

	return &cloudprovider.VirtualMachine{
		ID:           instance.InstanceId,
		Name:         instance.InstanceName,
		Status:       vmStatus,
		InstanceType: instance.InstanceType,
		ImageID:      instance.ImageId,
		PrivateIP:    privateIP,
		PublicIP:     publicIP,
		RegionID:     p.regionID,
		ZoneID:       instance.ZoneId,
		CreatedAt:    time.Now(), // Use current time as approximation
	}, nil
}

// ResetVMPassword 重置虚拟机密码
func (p *AlibabaProvider) ResetVMPassword(ctx context.Context, vmID, username, newPassword string) error {
	// Aliyun ECS ResetPassword requires the instance to be stopped
	// For now, we'll return unsupported operation since it requires changing instance state
	return cloudprovider.NewCloudError(
		cloudprovider.ErrUnsupportedOperation,
		"ResetVMPassword not fully implemented for Alibaba Cloud (requires instance to be stopped)",
		vmID,
	)
}

// UpdateVMConfig 更新虚拟机配置
func (p *AlibabaProvider) UpdateVMConfig(ctx context.Context, vmID, instanceType, name string) error {
	// This is a simplified implementation - in practice, changing instance type requires stopping the instance
	if instanceType != "" {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrUnsupportedOperation,
			"UpdateVMConfig (instance type change) not fully implemented for Alibaba Cloud (requires instance to be stopped)",
			vmID,
		)
	}

	// Only allow name changes
	if name != "" {
		request := ecs.CreateModifyInstanceAttributeRequest()
		request.Scheme = "https"
		request.InstanceId = vmID
		request.InstanceName = name

		_, err := p.ecsClient.ModifyInstanceAttribute(request)
		if err != nil {
			return cloudprovider.NewCloudError(
				cloudprovider.ErrOperationFailed,
				"failed to update instance name",
				err.Error(),
			)
		}
	}

	return nil
}

// convertECSStatus 转换 ECS 状态
func convertECSStatus(status string) cloudprovider.VMStatus {
	switch status {
	case "Running":
		return cloudprovider.VMStatusRunning
	case "Stopped":
		return cloudprovider.VMStatusStopped
	case "Starting":
		return cloudprovider.VMStatusStarting
	case "Stopping":
		return cloudprovider.VMStatusStopping
	default:
		return cloudprovider.VMStatusPending
	}
}

// ListImages 列出镜像
func (p *AlibabaProvider) ListImages(ctx context.Context, filter cloudprovider.ImageFilter) ([]*cloudprovider.Image, error) {
	request := ecs.CreateDescribeImagesRequest()
	request.Scheme = "https"
	request.RegionId = p.regionID

	if filter.Platform != "" {
		request.OSType = filter.Platform
	}

	response, err := p.ecsClient.DescribeImages(request)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to describe images",
			err.Error(),
		)
	}

	var images []*cloudprovider.Image
	for _, image := range response.Images.Image {
		images = append(images, &cloudprovider.Image{
			ID:          image.ImageId,
			Name:        image.ImageName,
			Description: image.Description,
			OSName:      image.OSName,
			Status:      image.Status,
			Size:        int64(image.Size),
		})
	}

	return images, nil
}

// GetImage 获取镜像
func (p *AlibabaProvider) GetImage(ctx context.Context, imageID string) (*cloudprovider.Image, error) {
	request := ecs.CreateDescribeImagesRequest()
	request.Scheme = "https"
	request.ImageId = imageID

	response, err := p.ecsClient.DescribeImages(request)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to describe image",
			err.Error(),
		)
	}

	if len(response.Images.Image) == 0 {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrResourceNotFound,
			"image not found",
			imageID,
		)
	}

	image := response.Images.Image[0]
	return &cloudprovider.Image{
		ID:          image.ImageId,
		Name:        image.ImageName,
		Description: image.Description,
		OSName:      image.OSName,
		Status:      image.Status,
		Size:        int64(image.Size),
	}, nil
}

// CreateKeypair 创建密钥对
func (p *AlibabaProvider) CreateKeypair(ctx context.Context, name, publicKey string) (*cloudprovider.Keypair, error) {
	return nil, cloudprovider.NewCloudError(
		cloudprovider.ErrUnsupportedOperation,
		"CreateKeypair not implemented",
		"",
	)
}

// DeleteKeypair 删除密钥对
func (p *AlibabaProvider) DeleteKeypair(ctx context.Context, keypairID string) error {
	return cloudprovider.NewCloudError(
		cloudprovider.ErrUnsupportedOperation,
		"DeleteKeypair not implemented",
		"",
	)
}

// ListKeypairs 列出密钥对
func (p *AlibabaProvider) ListKeypairs(ctx context.Context) ([]*cloudprovider.Keypair, error) {
	return nil, cloudprovider.NewCloudError(
		cloudprovider.ErrUnsupportedOperation,
		"ListKeypairs not implemented",
		"",
	)
}
