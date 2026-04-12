package aws

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/opencmp/opencmp/pkg/cloudprovider"
)

// CreateVM 创建虚拟机
func (p *AWSProvider) CreateVM(ctx context.Context, config cloudprovider.VMCreateConfig) (*cloudprovider.VirtualMachine, error) {
	input := &ec2.RunInstancesInput{
		ImageId:      aws.String(config.ImageID),
		InstanceType: types.InstanceType(config.InstanceType),
		MinCount:     aws.Int32(1),
		MaxCount:     aws.Int32(1),
	}

	// 设置子网
	if config.SubnetID != "" {
		input.SubnetId = aws.String(config.SubnetID)
	}

	// 设置安全组
	if len(config.SecurityGroups) > 0 {
		var sgIds []string
		for _, sg := range config.SecurityGroups {
			sgIds = append(sgIds, sg)
		}
		input.SecurityGroupIds = sgIds
	}

	// 设置名称标签
	if config.Name != "" {
		input.TagSpecifications = []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeInstance,
				Tags: []types.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String(config.Name),
					},
				},
			},
		}
	}

	result, err := p.ec2Client.RunInstances(ctx, input)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to create instance",
			err.Error(),
		)
	}

	if len(result.Instances) == 0 {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"no instance created",
			"",
		)
	}

	instance := result.Instances[0]
	return &cloudprovider.VirtualMachine{
		ID:           *instance.InstanceId,
		Name:         config.Name,
		Status:       cloudprovider.VMStatusPending,
		InstanceType: string(instance.InstanceType),
		ImageID:      config.ImageID,
		VPCID:        config.VPCID,
		SubnetID:     config.SubnetID,
		CreatedAt:    time.Now(),
		RegionID:     p.regionID,
	}, nil
}

// DeleteVM 删除虚拟机
func (p *AWSProvider) DeleteVM(ctx context.Context, vmID string) error {
	input := &ec2.TerminateInstancesInput{
		InstanceIds: []string{vmID},
	}

	_, err := p.ec2Client.TerminateInstances(ctx, input)
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
func (p *AWSProvider) StartVM(ctx context.Context, vmID string) error {
	input := &ec2.StartInstancesInput{
		InstanceIds: []string{vmID},
	}

	_, err := p.ec2Client.StartInstances(ctx, input)
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
func (p *AWSProvider) StopVM(ctx context.Context, vmID string) error {
	input := &ec2.StopInstancesInput{
		InstanceIds: []string{vmID},
	}

	_, err := p.ec2Client.StopInstances(ctx, input)
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
func (p *AWSProvider) RebootVM(ctx context.Context, vmID string) error {
	input := &ec2.RebootInstancesInput{
		InstanceIds: []string{vmID},
	}

	_, err := p.ec2Client.RebootInstances(ctx, input)
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
func (p *AWSProvider) GetVMStatus(ctx context.Context, vmID string) (*cloudprovider.VMStatus, error) {
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []string{vmID},
	}

	result, err := p.ec2Client.DescribeInstances(ctx, input)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to describe instance status",
			err.Error(),
		)
	}

	if len(result.Reservations) == 0 || len(result.Reservations[0].Instances) == 0 {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrResourceNotFound,
			"instance not found",
			vmID,
		)
	}

	instance := result.Reservations[0].Instances[0]
	status := convertEC2Status(instance.State.Name)

	return &status, nil
}

// ListVMs 列出虚拟机
func (p *AWSProvider) ListVMs(ctx context.Context, filter cloudprovider.VMListFilter) ([]*cloudprovider.VirtualMachine, error) {
	input := &ec2.DescribeInstancesInput{}

	// 设置 VPC 过滤
	if filter.VPCID != "" {
		input.Filters = []types.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []string{filter.VPCID},
			},
		}
	}

	result, err := p.ec2Client.DescribeInstances(ctx, input)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to describe instances",
			err.Error(),
		)
	}

	var vms []*cloudprovider.VirtualMachine
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			vmStatus := convertEC2Status(instance.State.Name)

			var name string
			for _, tag := range instance.Tags {
				if *tag.Key == "Name" {
					name = *tag.Value
					break
				}
			}

			var privateIP, publicIP string
			if instance.PrivateIpAddress != nil {
				privateIP = *instance.PrivateIpAddress
			}
			if instance.PublicIpAddress != nil {
				publicIP = *instance.PublicIpAddress
			}

			var zoneID string
			if instance.Placement != nil && instance.Placement.AvailabilityZone != nil {
				zoneID = *instance.Placement.AvailabilityZone
			}

			var imageID string
			if instance.ImageId != nil {
				imageID = *instance.ImageId
			}

			vms = append(vms, &cloudprovider.VirtualMachine{
				ID:           *instance.InstanceId,
				Name:         name,
				Status:       vmStatus,
				InstanceType: string(instance.InstanceType),
				ImageID:      imageID,
				PrivateIP:    privateIP,
				PublicIP:     publicIP,
				RegionID:     p.regionID,
				ZoneID:       zoneID,
			})
		}
	}

	return vms, nil
}

// GetVM 获取虚拟机详情
func (p *AWSProvider) GetVM(ctx context.Context, vmID string) (*cloudprovider.VirtualMachine, error) {
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []string{vmID},
	}

	result, err := p.ec2Client.DescribeInstances(ctx, input)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to describe instance",
			err.Error(),
		)
	}

	if len(result.Reservations) == 0 || len(result.Reservations[0].Instances) == 0 {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrResourceNotFound,
			"instance not found",
			vmID,
		)
	}

	instance := result.Reservations[0].Instances[0]
	vmStatus := convertEC2Status(instance.State.Name)

	var name string
	for _, tag := range instance.Tags {
		if *tag.Key == "Name" {
			name = *tag.Value
			break
		}
	}

	var privateIP, publicIP string
	if instance.PrivateIpAddress != nil {
		privateIP = *instance.PrivateIpAddress
	}
	if instance.PublicIpAddress != nil {
		publicIP = *instance.PublicIpAddress
	}

	var zoneID string
	if instance.Placement != nil && instance.Placement.AvailabilityZone != nil {
		zoneID = *instance.Placement.AvailabilityZone
	}

	var imageID string
	if instance.ImageId != nil {
		imageID = *instance.ImageId
	}

	var launchTime time.Time
	if instance.LaunchTime != nil {
		launchTime = *instance.LaunchTime
	}

	return &cloudprovider.VirtualMachine{
		ID:           *instance.InstanceId,
		Name:         name,
		Status:       vmStatus,
		InstanceType: string(instance.InstanceType),
		ImageID:      imageID,
		PrivateIP:    privateIP,
		PublicIP:     publicIP,
		RegionID:     p.regionID,
		ZoneID:       zoneID,
		CreatedAt:    launchTime,
	}, nil
}

// ListImages 列出镜像
func (p *AWSProvider) ListImages(ctx context.Context, filter cloudprovider.ImageFilter) ([]*cloudprovider.Image, error) {
	input := &ec2.DescribeImagesInput{
		Owners: []string{"amazon"},
	}

	result, err := p.ec2Client.DescribeImages(ctx, input)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to describe images",
			err.Error(),
		)
	}

	var images []*cloudprovider.Image
	for _, image := range result.Images {
		var name string
		if image.Name != nil {
			name = *image.Name
		}

		var osName string
		if image.PlatformDetails != nil {
			osName = *image.PlatformDetails
		}

		images = append(images, &cloudprovider.Image{
			ID:           *image.ImageId,
			Name:         name,
			Description:  *image.Description,
			OSName:       osName,
			Architecture: string(image.Architecture),
			Status:       string(image.State),
		})
	}

	return images, nil
}

// GetImage 获取镜像
func (p *AWSProvider) GetImage(ctx context.Context, imageID string) (*cloudprovider.Image, error) {
	input := &ec2.DescribeImagesInput{
		ImageIds: []string{imageID},
	}

	result, err := p.ec2Client.DescribeImages(ctx, input)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to describe image",
			err.Error(),
		)
	}

	if len(result.Images) == 0 {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrResourceNotFound,
			"image not found",
			imageID,
		)
	}

	image := result.Images[0]
	var name string
	if image.Name != nil {
		name = *image.Name
	}

	var osName string
	if image.PlatformDetails != nil {
		osName = *image.PlatformDetails
	}

	return &cloudprovider.Image{
		ID:           *image.ImageId,
		Name:         name,
		Description:  *image.Description,
		OSName:       osName,
		Architecture: string(image.Architecture),
		Status:       string(image.State),
	}, nil
}

// convertEC2Status 转换 AWS EC2 状态
func convertEC2Status(state types.InstanceStateName) cloudprovider.VMStatus {
	switch state {
	case types.InstanceStateNameRunning:
		return cloudprovider.VMStatusRunning
	case types.InstanceStateNameStopped:
		return cloudprovider.VMStatusStopped
	case types.InstanceStateNamePending:
		return cloudprovider.VMStatusPending
	case types.InstanceStateNameStopping:
		return cloudprovider.VMStatusStopping
	case types.InstanceStateNameShuttingDown:
		return cloudprovider.VMStatusStopping
	default:
		return cloudprovider.VMStatusPending
	}
}