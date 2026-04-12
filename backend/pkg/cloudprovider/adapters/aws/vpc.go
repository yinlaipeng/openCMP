package aws

import (
	"context"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/opencmp/opencmp/pkg/cloudprovider"
)

// CreateVPC 创建 VPC
func (p *AWSProvider) CreateVPC(ctx context.Context, config cloudprovider.VPCConfig) (*cloudprovider.VPC, error) {
	input := &ec2.CreateVpcInput{
		CidrBlock: aws.String(config.CIDR),
	}

	// 设置名称标签
	if config.Name != "" {
		input.TagSpecifications = []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeVpc,
				Tags: []types.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String(config.Name),
					},
				},
			},
		}
	}

	result, err := p.ec2Client.CreateVpc(ctx, input)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to create vpc",
			err.Error(),
		)
	}

	return &cloudprovider.VPC{
		ID:          *result.Vpc.VpcId,
		Name:        config.Name,
		CIDR:        config.CIDR,
		Description: config.Description,
		Status:      "Pending",
		RegionID:    p.regionID,
	}, nil
}

// DeleteVPC 删除 VPC
func (p *AWSProvider) DeleteVPC(ctx context.Context, vpcID string) error {
	input := &ec2.DeleteVpcInput{
		VpcId: aws.String(vpcID),
	}

	_, err := p.ec2Client.DeleteVpc(ctx, input)
	if err != nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to delete vpc",
			err.Error(),
		)
	}

	return nil
}

// GetVPC 获取 VPC
func (p *AWSProvider) GetVPC(ctx context.Context, vpcID string) (*cloudprovider.VPC, error) {
	input := &ec2.DescribeVpcsInput{
		VpcIds: []string{vpcID},
	}

	result, err := p.ec2Client.DescribeVpcs(ctx, input)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to describe vpc",
			err.Error(),
		)
	}

	if len(result.Vpcs) == 0 {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrResourceNotFound,
			"vpc not found",
			vpcID,
		)
	}

	vpc := result.Vpcs[0]
	var name string
	for _, tag := range vpc.Tags {
		if *tag.Key == "Name" {
			name = *tag.Value
			break
		}
	}

	return &cloudprovider.VPC{
		ID:          *vpc.VpcId,
		Name:        name,
		CIDR:        *vpc.CidrBlock,
		Description: "",
		Status:      string(vpc.State),
		RegionID:    p.regionID,
	}, nil
}

// ListVPCs 列出 VPC
func (p *AWSProvider) ListVPCs(ctx context.Context, filter cloudprovider.VPCFilter) ([]*cloudprovider.VPC, error) {
	input := &ec2.DescribeVpcsInput{}

	if filter.VPCID != "" {
		input.VpcIds = []string{filter.VPCID}
	}

	result, err := p.ec2Client.DescribeVpcs(ctx, input)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to describe vpcs",
			err.Error(),
		)
	}

	var vpcs []*cloudprovider.VPC
	for _, vpc := range result.Vpcs {
		var name string
		for _, tag := range vpc.Tags {
			if *tag.Key == "Name" {
				name = *tag.Value
				break
			}
		}

		vpcs = append(vpcs, &cloudprovider.VPC{
			ID:       *vpc.VpcId,
			Name:     name,
			CIDR:     *vpc.CidrBlock,
			Status:   string(vpc.State),
			RegionID: p.regionID,
		})
	}

	return vpcs, nil
}

// CreateSubnet 创建子网
func (p *AWSProvider) CreateSubnet(ctx context.Context, config cloudprovider.SubnetConfig) (*cloudprovider.Subnet, error) {
	input := &ec2.CreateSubnetInput{
		VpcId:     aws.String(config.VPCID),
		CidrBlock: aws.String(config.CIDR),
	}

	if config.ZoneID != "" {
		input.AvailabilityZone = aws.String(config.ZoneID)
	}

	// 设置名称标签
	if config.Name != "" {
		input.TagSpecifications = []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeSubnet,
				Tags: []types.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String(config.Name),
					},
				},
			},
		}
	}

	result, err := p.ec2Client.CreateSubnet(ctx, input)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to create subnet",
			err.Error(),
		)
	}

	return &cloudprovider.Subnet{
		ID:     *result.Subnet.SubnetId,
		Name:   config.Name,
		VPCID:  config.VPCID,
		CIDR:   config.CIDR,
		ZoneID: config.ZoneID,
		Status: "Pending",
	}, nil
}

// DeleteSubnet 删除子网
func (p *AWSProvider) DeleteSubnet(ctx context.Context, subnetID string) error {
	input := &ec2.DeleteSubnetInput{
		SubnetId: aws.String(subnetID),
	}

	_, err := p.ec2Client.DeleteSubnet(ctx, input)
	if err != nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to delete subnet",
			err.Error(),
		)
	}

	return nil
}

// GetSubnet 获取子网
func (p *AWSProvider) GetSubnet(ctx context.Context, subnetID string) (*cloudprovider.Subnet, error) {
	input := &ec2.DescribeSubnetsInput{
		SubnetIds: []string{subnetID},
	}

	result, err := p.ec2Client.DescribeSubnets(ctx, input)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to describe subnet",
			err.Error(),
		)
	}

	if len(result.Subnets) == 0 {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrResourceNotFound,
			"subnet not found",
			subnetID,
		)
	}

	subnet := result.Subnets[0]
	var name string
	for _, tag := range subnet.Tags {
		if *tag.Key == "Name" {
			name = *tag.Value
			break
		}
	}

	var zoneID string
	if subnet.AvailabilityZone != nil {
		zoneID = *subnet.AvailabilityZone
	}

	return &cloudprovider.Subnet{
		ID:     *subnet.SubnetId,
		Name:   name,
		VPCID:  *subnet.VpcId,
		CIDR:   *subnet.CidrBlock,
		ZoneID: zoneID,
		Status: string(subnet.State),
	}, nil
}

// ListSubnets 列出子网
func (p *AWSProvider) ListSubnets(ctx context.Context, filter cloudprovider.SubnetFilter) ([]*cloudprovider.Subnet, error) {
	input := &ec2.DescribeSubnetsInput{}

	if filter.VPCID != "" {
		input.Filters = []types.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []string{filter.VPCID},
			},
		}
	}

	result, err := p.ec2Client.DescribeSubnets(ctx, input)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to describe subnets",
			err.Error(),
		)
	}

	var subnets []*cloudprovider.Subnet
	for _, subnet := range result.Subnets {
		var name string
		for _, tag := range subnet.Tags {
			if *tag.Key == "Name" {
				name = *tag.Value
				break
			}
		}

		var zoneID string
		if subnet.AvailabilityZone != nil {
			zoneID = *subnet.AvailabilityZone
		}

		subnets = append(subnets, &cloudprovider.Subnet{
			ID:     *subnet.SubnetId,
			Name:   name,
			VPCID:  *subnet.VpcId,
			CIDR:   *subnet.CidrBlock,
			ZoneID: zoneID,
			Status: string(subnet.State),
		})
	}

	return subnets, nil
}

// CreateSecurityGroup 创建安全组
func (p *AWSProvider) CreateSecurityGroup(ctx context.Context, config cloudprovider.SGConfig) (*cloudprovider.SecurityGroup, error) {
	input := &ec2.CreateSecurityGroupInput{
		GroupName:   aws.String(config.Name),
		Description: aws.String(config.Description),
	}

	if config.VPCID != "" {
		input.VpcId = aws.String(config.VPCID)
	}

	result, err := p.ec2Client.CreateSecurityGroup(ctx, input)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to create security group",
			err.Error(),
		)
	}

	return &cloudprovider.SecurityGroup{
		ID:          *result.GroupId,
		Name:        config.Name,
		Description: config.Description,
		VPCID:       config.VPCID,
	}, nil
}

// DeleteSecurityGroup 删除安全组
func (p *AWSProvider) DeleteSecurityGroup(ctx context.Context, sgID string) error {
	input := &ec2.DeleteSecurityGroupInput{
		GroupId: aws.String(sgID),
	}

	_, err := p.ec2Client.DeleteSecurityGroup(ctx, input)
	if err != nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to delete security group",
			err.Error(),
		)
	}

	return nil
}

// AuthorizeSecurityGroup 授权安全组规则
func (p *AWSProvider) AuthorizeSecurityGroup(ctx context.Context, sgID string, rules []cloudprovider.SGRule) error {
	for _, rule := range rules {
		input := &ec2.AuthorizeSecurityGroupIngressInput{
			GroupId:    aws.String(sgID),
			IpProtocol: aws.String(rule.Protocol),
			FromPort:   aws.Int32(int32(parsePortRange(rule.PortRange, true))),
			ToPort:     aws.Int32(int32(parsePortRange(rule.PortRange, false))),
		}

		if rule.CIDR != "" {
			input.CidrIp = aws.String(rule.CIDR)
		}

		_, err := p.ec2Client.AuthorizeSecurityGroupIngress(ctx, input)
		if err != nil {
			return cloudprovider.NewCloudError(
				cloudprovider.ErrOperationFailed,
				"failed to authorize security group rule",
				err.Error(),
			)
		}
	}

	return nil
}

// RevokeSecurityGroup 撤销安全组规则
func (p *AWSProvider) RevokeSecurityGroup(ctx context.Context, sgID string, rules []cloudprovider.SGRule) error {
	for _, rule := range rules {
		input := &ec2.RevokeSecurityGroupIngressInput{
			GroupId:    aws.String(sgID),
			IpProtocol: aws.String(rule.Protocol),
			FromPort:   aws.Int32(int32(parsePortRange(rule.PortRange, true))),
			ToPort:     aws.Int32(int32(parsePortRange(rule.PortRange, false))),
		}

		if rule.CIDR != "" {
			input.CidrIp = aws.String(rule.CIDR)
		}

		_, err := p.ec2Client.RevokeSecurityGroupIngress(ctx, input)
		if err != nil {
			return cloudprovider.NewCloudError(
				cloudprovider.ErrOperationFailed,
				"failed to revoke security group rule",
				err.Error(),
			)
		}
	}

	return nil
}

// ListSecurityGroups 列出安全组
func (p *AWSProvider) ListSecurityGroups(ctx context.Context, filter cloudprovider.SGFilter) ([]*cloudprovider.SecurityGroup, error) {
	input := &ec2.DescribeSecurityGroupsInput{}

	if filter.VPCID != "" {
		input.Filters = []types.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []string{filter.VPCID},
			},
		}
	}

	result, err := p.ec2Client.DescribeSecurityGroups(ctx, input)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to describe security groups",
			err.Error(),
		)
	}

	var sgs []*cloudprovider.SecurityGroup
	for _, sg := range result.SecurityGroups {
		sgs = append(sgs, &cloudprovider.SecurityGroup{
			ID:          *sg.GroupId,
			Name:        *sg.GroupName,
			Description: *sg.Description,
			VPCID:       *sg.VpcId,
		})
	}

	return sgs, nil
}

// AllocateEIP 分配弹性 IP
func (p *AWSProvider) AllocateEIP(ctx context.Context, config cloudprovider.EIPConfig) (*cloudprovider.EIP, error) {
	input := &ec2.AllocateAddressInput{
		Domain: types.DomainTypeVpc,
	}

	result, err := p.ec2Client.AllocateAddress(ctx, input)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to allocate eip",
			err.Error(),
		)
	}

	return &cloudprovider.EIP{
		ID:        *result.AllocationId,
		Address:   *result.PublicIp,
		Bandwidth: config.Bandwidth,
		Status:    "available",
		RegionID:  p.regionID,
	}, nil
}

// ReleaseEIP 释放弹性 IP
func (p *AWSProvider) ReleaseEIP(ctx context.Context, eipID string) error {
	input := &ec2.ReleaseAddressInput{
		AllocationId: aws.String(eipID),
	}

	_, err := p.ec2Client.ReleaseAddress(ctx, input)
	if err != nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to release eip",
			err.Error(),
		)
	}

	return nil
}

// AssociateEIP 关联弹性 IP
func (p *AWSProvider) AssociateEIP(ctx context.Context, eipID, resourceID string) error {
	input := &ec2.AssociateAddressInput{
		AllocationId: aws.String(eipID),
		InstanceId:   aws.String(resourceID),
	}

	_, err := p.ec2Client.AssociateAddress(ctx, input)
	if err != nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to associate eip",
			err.Error(),
		)
	}

	return nil
}

// DissociateEIP 解绑弹性 IP
func (p *AWSProvider) DissociateEIP(ctx context.Context, eipID string) error {
	// 首先需要找到关联 ID
	describeInput := &ec2.DescribeAddressesInput{
		AllocationIds: []string{eipID},
	}

	result, err := p.ec2Client.DescribeAddresses(ctx, describeInput)
	if err != nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to describe address",
			err.Error(),
		)
	}

	if len(result.Addresses) == 0 || result.Addresses[0].AssociationId == nil {
		return nil // 已经解绑
	}

	input := &ec2.DisassociateAddressInput{
		AssociationId: result.Addresses[0].AssociationId,
	}

	_, err = p.ec2Client.DisassociateAddress(ctx, input)
	if err != nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to disassociate eip",
			err.Error(),
		)
	}

	return nil
}

// ListEIPs 列出弹性 IP
func (p *AWSProvider) ListEIPs(ctx context.Context, filter cloudprovider.EIPFilter) ([]*cloudprovider.EIP, error) {
	input := &ec2.DescribeAddressesInput{}

	if filter.EIPID != "" {
		input.AllocationIds = []string{filter.EIPID}
	}

	result, err := p.ec2Client.DescribeAddresses(ctx, input)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to describe eip addresses",
			err.Error(),
		)
	}

	var eips []*cloudprovider.EIP
	for _, addr := range result.Addresses {
		var resourceID string
		if addr.InstanceId != nil {
			resourceID = *addr.InstanceId
		}

		eips = append(eips, &cloudprovider.EIP{
			ID:         *addr.AllocationId,
			Address:    *addr.PublicIp,
			Bandwidth:  0,
			Status:     "available",
			ResourceID: resourceID,
			RegionID:   p.regionID,
		})
	}

	return eips, nil
}

// parsePortRange 解析端口范围 (格式: "80-80" 或 "80")
func parsePortRange(portRange string, isFrom bool) int {
	// 简单实现，假设格式为 "from-to" 或 单个端口
	if portRange == "" {
		return 0
	}

	// 尝试分割
	parts := []string{portRange}
	if strings.Contains(portRange, "-") {
		parts = strings.Split(portRange, "-")
	}

	if len(parts) == 1 {
		port, _ := strconv.Atoi(parts[0])
		return port
	}

	if isFrom {
		port, _ := strconv.Atoi(parts[0])
		return port
	}
	port, _ := strconv.Atoi(parts[1])
	return port
}