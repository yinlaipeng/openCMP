package tencent

import (
	"context"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/opencmp/opencmp/pkg/cloudprovider"
)

// CreateVPC 创建 VPC
func (p *TencentProvider) CreateVPC(ctx context.Context, config cloudprovider.VPCConfig) (*cloudprovider.VPC, error) {
	request := vpc.NewCreateVpcRequest()
	request.VpcName = &config.Name
	request.CidrBlock = &config.CIDR

	response, err := p.vpcClient.CreateVpc(request)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to create vpc",
			err.Error(),
		)
	}

	return &cloudprovider.VPC{
		ID:          *response.Response.Vpc.VpcId,
		Name:        config.Name,
		CIDR:        config.CIDR,
		Description: config.Description,
		Status:      "Pending",
		RegionID:    p.regionID,
	}, nil
}

// DeleteVPC 删除 VPC
func (p *TencentProvider) DeleteVPC(ctx context.Context, vpcID string) error {
	request := vpc.NewDeleteVpcRequest()
	request.VpcId = &vpcID

	_, err := p.vpcClient.DeleteVpc(request)
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
func (p *TencentProvider) GetVPC(ctx context.Context, vpcID string) (*cloudprovider.VPC, error) {
	request := vpc.NewDescribeVpcsRequest()
	request.VpcIds = []*string{&vpcID}

	response, err := p.vpcClient.DescribeVpcs(request)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to describe vpc",
			err.Error(),
		)
	}

	if len(response.Response.VpcSet) == 0 {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrResourceNotFound,
			"vpc not found",
			vpcID,
		)
	}

	vpcItem := response.Response.VpcSet[0]
	return &cloudprovider.VPC{
		ID:          *vpcItem.VpcId,
		Name:        *vpcItem.VpcName,
		CIDR:        *vpcItem.CidrBlock,
		Description: "",
		Status:      "Available",
		RegionID:    p.regionID,
	}, nil
}

// ListVPCs 列出 VPC
func (p *TencentProvider) ListVPCs(ctx context.Context, filter cloudprovider.VPCFilter) ([]*cloudprovider.VPC, error) {
	request := vpc.NewDescribeVpcsRequest()

	if filter.VPCID != "" {
		request.VpcIds = []*string{&filter.VPCID}
	}

	response, err := p.vpcClient.DescribeVpcs(request)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to describe vpcs",
			err.Error(),
		)
	}

	var vpcs []*cloudprovider.VPC
	for _, vpcItem := range response.Response.VpcSet {
		vpcs = append(vpcs, &cloudprovider.VPC{
			ID:       *vpcItem.VpcId,
			Name:     *vpcItem.VpcName,
			CIDR:     *vpcItem.CidrBlock,
			Status:   "Available",
			RegionID: p.regionID,
		})
	}

	return vpcs, nil
}

// CreateSubnet 创建子网
func (p *TencentProvider) CreateSubnet(ctx context.Context, config cloudprovider.SubnetConfig) (*cloudprovider.Subnet, error) {
	request := vpc.NewCreateSubnetRequest()
	request.VpcId = &config.VPCID
	request.SubnetName = &config.Name
	request.CidrBlock = &config.CIDR
	request.Zone = &config.ZoneID

	response, err := p.vpcClient.CreateSubnet(request)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to create subnet",
			err.Error(),
		)
	}

	return &cloudprovider.Subnet{
		ID:     *response.Response.Subnet.SubnetId,
		Name:   config.Name,
		VPCID:  config.VPCID,
		CIDR:   config.CIDR,
		ZoneID: config.ZoneID,
		Status: "Pending",
	}, nil
}

// DeleteSubnet 删除子网
func (p *TencentProvider) DeleteSubnet(ctx context.Context, subnetID string) error {
	request := vpc.NewDeleteSubnetRequest()
	request.SubnetId = &subnetID

	_, err := p.vpcClient.DeleteSubnet(request)
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
func (p *TencentProvider) GetSubnet(ctx context.Context, subnetID string) (*cloudprovider.Subnet, error) {
	request := vpc.NewDescribeSubnetsRequest()
	request.SubnetIds = []*string{&subnetID}

	response, err := p.vpcClient.DescribeSubnets(request)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to describe subnet",
			err.Error(),
		)
	}

	if len(response.Response.SubnetSet) == 0 {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrResourceNotFound,
			"subnet not found",
			subnetID,
		)
	}

	subnet := response.Response.SubnetSet[0]
	return &cloudprovider.Subnet{
		ID:     *subnet.SubnetId,
		Name:   *subnet.SubnetName,
		VPCID:  *subnet.VpcId,
		CIDR:   *subnet.CidrBlock,
		ZoneID: *subnet.Zone,
		Status: "Available",
	}, nil
}

// ListSubnets 列出子网
func (p *TencentProvider) ListSubnets(ctx context.Context, filter cloudprovider.SubnetFilter) ([]*cloudprovider.Subnet, error) {
	request := vpc.NewDescribeSubnetsRequest()

	if filter.VPCID != "" {
		request.Filters = []*vpc.Filter{
			{
				Name:   common.StringPtr("vpc-id"),
				Values: []*string{&filter.VPCID},
			},
		}
	}

	response, err := p.vpcClient.DescribeSubnets(request)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to describe subnets",
			err.Error(),
		)
	}

	var subnets []*cloudprovider.Subnet
	for _, subnet := range response.Response.SubnetSet {
		subnets = append(subnets, &cloudprovider.Subnet{
			ID:     *subnet.SubnetId,
			Name:   *subnet.SubnetName,
			VPCID:  *subnet.VpcId,
			CIDR:   *subnet.CidrBlock,
			ZoneID: *subnet.Zone,
			Status: "Available",
		})
	}

	return subnets, nil
}

// CreateSecurityGroup 创建安全组
func (p *TencentProvider) CreateSecurityGroup(ctx context.Context, config cloudprovider.SGConfig) (*cloudprovider.SecurityGroup, error) {
	request := vpc.NewCreateSecurityGroupRequest()
	request.GroupName = &config.Name
	request.GroupDescription = &config.Description

	response, err := p.vpcClient.CreateSecurityGroup(request)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to create security group",
			err.Error(),
		)
	}

	return &cloudprovider.SecurityGroup{
		ID:          *response.Response.SecurityGroup.SecurityGroupId,
		Name:        config.Name,
		Description: config.Description,
		VPCID:       config.VPCID,
	}, nil
}

// DeleteSecurityGroup 删除安全组
func (p *TencentProvider) DeleteSecurityGroup(ctx context.Context, sgID string) error {
	request := vpc.NewDeleteSecurityGroupRequest()
	request.SecurityGroupId = &sgID

	_, err := p.vpcClient.DeleteSecurityGroup(request)
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
func (p *TencentProvider) AuthorizeSecurityGroup(ctx context.Context, sgID string, rules []cloudprovider.SGRule) error {
	for _, rule := range rules {
		request := vpc.NewCreateSecurityGroupPoliciesRequest()
		request.SecurityGroupId = &sgID

		policy := &vpc.SecurityGroupPolicy{
			Protocol:          &rule.Protocol,
			Port:              &rule.PortRange,
			CidrBlock:         &rule.CIDR,
			Action:            &rule.Action,
			PolicyDescription: common.StringPtr(""),
		}

		request.SecurityGroupPolicySet = &vpc.SecurityGroupPolicySet{
			Ingress: []*vpc.SecurityGroupPolicy{policy},
		}

		_, err := p.vpcClient.CreateSecurityGroupPolicies(request)
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
func (p *TencentProvider) RevokeSecurityGroup(ctx context.Context, sgID string, rules []cloudprovider.SGRule) error {
	for _, rule := range rules {
		request := vpc.NewDeleteSecurityGroupPoliciesRequest()
		request.SecurityGroupId = &sgID

		policy := &vpc.SecurityGroupPolicy{
			Protocol:  &rule.Protocol,
			Port:      &rule.PortRange,
			CidrBlock: &rule.CIDR,
			Action:    &rule.Action,
		}

		request.SecurityGroupPolicySet = &vpc.SecurityGroupPolicySet{
			Ingress: []*vpc.SecurityGroupPolicy{policy},
		}

		_, err := p.vpcClient.DeleteSecurityGroupPolicies(request)
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
func (p *TencentProvider) ListSecurityGroups(ctx context.Context, filter cloudprovider.SGFilter) ([]*cloudprovider.SecurityGroup, error) {
	request := vpc.NewDescribeSecurityGroupsRequest()

	if filter.VPCID != "" {
		request.Filters = []*vpc.Filter{
			{
				Name:   common.StringPtr("vpc-id"),
				Values: []*string{&filter.VPCID},
			},
		}
	}

	response, err := p.vpcClient.DescribeSecurityGroups(request)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to describe security groups",
			err.Error(),
		)
	}

	var sgs []*cloudprovider.SecurityGroup
	for _, sg := range response.Response.SecurityGroupSet {
		var desc string
		if sg.SecurityGroupDesc != nil {
			desc = *sg.SecurityGroupDesc
		}
		sgs = append(sgs, &cloudprovider.SecurityGroup{
			ID:          *sg.SecurityGroupId,
			Name:        *sg.SecurityGroupName,
			Description: desc,
		})
	}

	return sgs, nil
}

// AllocateEIP 分配弹性 IP
func (p *TencentProvider) AllocateEIP(ctx context.Context, config cloudprovider.EIPConfig) (*cloudprovider.EIP, error) {
	request := vpc.NewAllocateAddressesRequest()

	request.AddressCount = common.Int64Ptr(1)
	request.InternetMaxBandwidthOut = common.Int64Ptr(int64(config.Bandwidth))
	request.AddressType = common.StringPtr("WANIP")

	response, err := p.vpcClient.AllocateAddresses(request)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to allocate eip",
			err.Error(),
		)
	}

	if len(response.Response.AddressSet) == 0 {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"no eip allocated",
			"",
		)
	}

	eipAddr := response.Response.AddressSet[0]
	return &cloudprovider.EIP{
		ID:        *eipAddr,
		Bandwidth: config.Bandwidth,
		Status:    "available",
		RegionID:  p.regionID,
	}, nil
}

// ReleaseEIP 释放弹性 IP
func (p *TencentProvider) ReleaseEIP(ctx context.Context, eipID string) error {
	request := vpc.NewReleaseAddressesRequest()
	request.AddressIds = []*string{&eipID}

	_, err := p.vpcClient.ReleaseAddresses(request)
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
func (p *TencentProvider) AssociateEIP(ctx context.Context, eipID, resourceID string) error {
	request := vpc.NewAssociateAddressRequest()
	request.AddressId = &eipID
	request.InstanceId = &resourceID

	_, err := p.vpcClient.AssociateAddress(request)
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
func (p *TencentProvider) DissociateEIP(ctx context.Context, eipID string) error {
	request := vpc.NewDisassociateAddressRequest()
	request.AddressId = &eipID

	_, err := p.vpcClient.DisassociateAddress(request)
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
func (p *TencentProvider) ListEIPs(ctx context.Context, filter cloudprovider.EIPFilter) ([]*cloudprovider.EIP, error) {
	request := vpc.NewDescribeAddressesRequest()

	if filter.EIPID != "" {
		request.AddressIds = []*string{&filter.EIPID}
	}

	response, err := p.vpcClient.DescribeAddresses(request)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to describe eip addresses",
			err.Error(),
		)
	}

	var eips []*cloudprovider.EIP
	for _, addr := range response.Response.AddressSet {
		var publicIP string
		if addr.AddressIp != nil {
			publicIP = *addr.AddressIp
		}
		var resourceID string
		if addr.InstanceId != nil {
			resourceID = *addr.InstanceId
		}
		eips = append(eips, &cloudprovider.EIP{
			ID:         *addr.AddressId,
			Address:    publicIP,
			Bandwidth:  0, // 腾讯云 SDK 不直接返回带宽
			Status:     *addr.AddressStatus,
			ResourceID: resourceID,
			RegionID:   p.regionID,
		})
	}

	return eips, nil
}