package alibaba

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/opencmp/opencmp/pkg/cloudprovider"
)

// CreateVPC 创建 VPC
func (p *AlibabaProvider) CreateVPC(ctx context.Context, config cloudprovider.VPCConfig) (*cloudprovider.VPC, error) {
	request := vpc.CreateCreateVpcRequest()
	request.Scheme = "https"
	request.VpcName = config.Name
	request.CidrBlock = config.CIDR
	request.RegionId = p.regionID

	response, err := p.vpcClient.CreateVpc(request)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to create vpc",
			err.Error(),
		)
	}

	return &cloudprovider.VPC{
		ID:          response.VpcId,
		Name:        config.Name,
		CIDR:        config.CIDR,
		Description: config.Description,
		Status:      "Pending",
		RegionID:    p.regionID,
		CreatedAt:   time.Now(),
	}, nil
}

// DeleteVPC 删除 VPC
func (p *AlibabaProvider) DeleteVPC(ctx context.Context, vpcID string) error {
	request := vpc.CreateDeleteVpcRequest()
	request.Scheme = "https"
	request.VpcId = vpcID

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
func (p *AlibabaProvider) GetVPC(ctx context.Context, vpcID string) (*cloudprovider.VPC, error) {
	request := vpc.CreateDescribeVpcsRequest()
	request.Scheme = "https"
	request.VpcId = vpcID
	request.RegionId = p.regionID

	response, err := p.vpcClient.DescribeVpcs(request)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to describe vpc",
			err.Error(),
		)
	}

	if len(response.Vpcs.Vpc) == 0 {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrResourceNotFound,
			"vpc not found",
			vpcID,
		)
	}

	vpcItem := response.Vpcs.Vpc[0]
	return &cloudprovider.VPC{
		ID:          vpcItem.VpcId,
		Name:        vpcItem.VpcName,
		CIDR:        vpcItem.CidrBlock,
		Description: vpcItem.Description,
		Status:      vpcItem.Status,
		RegionID:    p.regionID,
	}, nil
}

// ListVPCs 列出 VPC
func (p *AlibabaProvider) ListVPCs(ctx context.Context, filter cloudprovider.VPCFilter) ([]*cloudprovider.VPC, error) {
	request := vpc.CreateDescribeVpcsRequest()
	request.Scheme = "https"
	request.RegionId = p.regionID

	if filter.VPCID != "" {
		request.VpcId = filter.VPCID
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
	for _, vpcItem := range response.Vpcs.Vpc {
		vpcs = append(vpcs, &cloudprovider.VPC{
			ID:       vpcItem.VpcId,
			Name:     vpcItem.VpcName,
			CIDR:     vpcItem.CidrBlock,
			Status:   vpcItem.Status,
			RegionID: p.regionID,
		})
	}

	return vpcs, nil
}

// CreateSubnet 创建子网
func (p *AlibabaProvider) CreateSubnet(ctx context.Context, config cloudprovider.SubnetConfig) (*cloudprovider.Subnet, error) {
	request := vpc.CreateCreateVSwitchRequest()
	request.Scheme = "https"
	request.VSwitchName = config.Name
	request.VpcId = config.VPCID
	request.CidrBlock = config.CIDR
	request.ZoneId = config.ZoneID

	response, err := p.vpcClient.CreateVSwitch(request)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to create vswitch",
			err.Error(),
		)
	}

	return &cloudprovider.Subnet{
		ID:     response.VSwitchId,
		Name:   config.Name,
		VPCID:  config.VPCID,
		CIDR:   config.CIDR,
		ZoneID: config.ZoneID,
		Status: "Pending",
	}, nil
}

// DeleteSubnet 删除子网
func (p *AlibabaProvider) DeleteSubnet(ctx context.Context, subnetID string) error {
	request := vpc.CreateDeleteVSwitchRequest()
	request.Scheme = "https"
	request.VSwitchId = subnetID

	_, err := p.vpcClient.DeleteVSwitch(request)
	if err != nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to delete vswitch",
			err.Error(),
		)
	}

	return nil
}

// GetSubnet 获取子网
func (p *AlibabaProvider) GetSubnet(ctx context.Context, subnetID string) (*cloudprovider.Subnet, error) {
	request := vpc.CreateDescribeVSwitchesRequest()
	request.Scheme = "https"
	request.VSwitchId = subnetID

	response, err := p.vpcClient.DescribeVSwitches(request)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to describe vswitch",
			err.Error(),
		)
	}

	if len(response.VSwitches.VSwitch) == 0 {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrResourceNotFound,
			"vswitch not found",
			subnetID,
		)
	}

	vswitch := response.VSwitches.VSwitch[0]
	return &cloudprovider.Subnet{
		ID:     vswitch.VSwitchId,
		Name:   vswitch.VSwitchName,
		VPCID:  vswitch.VpcId,
		CIDR:   vswitch.CidrBlock,
		ZoneID: vswitch.ZoneId,
		Status: vswitch.Status,
	}, nil
}

// ListSubnets 列出子网
func (p *AlibabaProvider) ListSubnets(ctx context.Context, filter cloudprovider.SubnetFilter) ([]*cloudprovider.Subnet, error) {
	request := vpc.CreateDescribeVSwitchesRequest()
	request.Scheme = "https"

	if filter.VPCID != "" {
		request.VpcId = filter.VPCID
	}

	response, err := p.vpcClient.DescribeVSwitches(request)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to describe vswitches",
			err.Error(),
		)
	}

	var subnets []*cloudprovider.Subnet
	for _, vswitch := range response.VSwitches.VSwitch {
		subnets = append(subnets, &cloudprovider.Subnet{
			ID:     vswitch.VSwitchId,
			Name:   vswitch.VSwitchName,
			VPCID:  vswitch.VpcId,
			CIDR:   vswitch.CidrBlock,
			ZoneID: vswitch.ZoneId,
			Status: vswitch.Status,
		})
	}

	return subnets, nil
}

// CreateSecurityGroup 创建安全组
func (p *AlibabaProvider) CreateSecurityGroup(ctx context.Context, config cloudprovider.SGConfig) (*cloudprovider.SecurityGroup, error) {
	request := ecs.CreateCreateSecurityGroupRequest()
	request.Scheme = "https"
	request.SecurityGroupName = config.Name
	request.Description = config.Description
	request.VpcId = config.VPCID

	response, err := p.ecsClient.CreateSecurityGroup(request)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to create security group",
			err.Error(),
		)
	}

	return &cloudprovider.SecurityGroup{
		ID:          response.SecurityGroupId,
		Name:        config.Name,
		Description: config.Description,
		VPCID:       config.VPCID,
	}, nil
}

// DeleteSecurityGroup 删除安全组
func (p *AlibabaProvider) DeleteSecurityGroup(ctx context.Context, sgID string) error {
	request := ecs.CreateDeleteSecurityGroupRequest()
	request.Scheme = "https"
	request.SecurityGroupId = sgID

	_, err := p.ecsClient.DeleteSecurityGroup(request)
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
func (p *AlibabaProvider) AuthorizeSecurityGroup(ctx context.Context, sgID string, rules []cloudprovider.SGRule) error {
	for _, rule := range rules {
		request := ecs.CreateAuthorizeSecurityGroupRequest()
		request.Scheme = "https"
		request.SecurityGroupId = sgID
		request.IpProtocol = rule.Protocol
		request.PortRange = rule.PortRange
		request.SourceCidrIp = rule.CIDR
		request.Policy = rule.Action

		_, err := p.ecsClient.AuthorizeSecurityGroup(request)
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
func (p *AlibabaProvider) RevokeSecurityGroup(ctx context.Context, sgID string, rules []cloudprovider.SGRule) error {
	for _, rule := range rules {
		request := ecs.CreateRevokeSecurityGroupRequest()
		request.Scheme = "https"
		request.SecurityGroupId = sgID
		request.IpProtocol = rule.Protocol
		request.PortRange = rule.PortRange
		request.SourceCidrIp = rule.CIDR
		request.Policy = rule.Action

		_, err := p.ecsClient.RevokeSecurityGroup(request)
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
func (p *AlibabaProvider) ListSecurityGroups(ctx context.Context, filter cloudprovider.SGFilter) ([]*cloudprovider.SecurityGroup, error) {
	request := ecs.CreateDescribeSecurityGroupsRequest()
	request.Scheme = "https"
	request.RegionId = p.regionID

	if filter.VPCID != "" {
		request.VpcId = filter.VPCID
	}

	response, err := p.ecsClient.DescribeSecurityGroups(request)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to describe security groups",
			err.Error(),
		)
	}

	var sgs []*cloudprovider.SecurityGroup
	for _, sg := range response.SecurityGroups.SecurityGroup {
		sgs = append(sgs, &cloudprovider.SecurityGroup{
			ID:          sg.SecurityGroupId,
			Name:        sg.SecurityGroupName,
			Description: sg.Description,
			VPCID:       sg.VpcId,
		})
	}

	return sgs, nil
}

// AllocateEIP 分配弹性 IP
func (p *AlibabaProvider) AllocateEIP(ctx context.Context, config cloudprovider.EIPConfig) (*cloudprovider.EIP, error) {
	request := vpc.CreateAllocateEipAddressRequest()
	request.Scheme = "https"
	request.RegionId = p.regionID
	request.Bandwidth = strconv.Itoa(config.Bandwidth)

	response, err := p.vpcClient.AllocateEipAddress(request)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to allocate eip",
			err.Error(),
		)
	}

	return &cloudprovider.EIP{
		ID:        response.AllocationId,
		Address:   response.EipAddress,
		Bandwidth: config.Bandwidth,
		Status:    "available",
		RegionID:  p.regionID,
	}, nil
}

// ReleaseEIP 释放弹性 IP
func (p *AlibabaProvider) ReleaseEIP(ctx context.Context, eipID string) error {
	request := vpc.CreateReleaseEipAddressRequest()
	request.Scheme = "https"
	request.AllocationId = eipID

	_, err := p.vpcClient.ReleaseEipAddress(request)
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
func (p *AlibabaProvider) AssociateEIP(ctx context.Context, eipID, resourceID string) error {
	request := vpc.CreateAssociateEipAddressRequest()
	request.Scheme = "https"
	request.AllocationId = eipID
	request.InstanceId = resourceID

	_, err := p.vpcClient.AssociateEipAddress(request)
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
func (p *AlibabaProvider) DissociateEIP(ctx context.Context, eipID string) error {
	request := vpc.CreateUnassociateEipAddressRequest()
	request.Scheme = "https"
	request.AllocationId = eipID

	_, err := p.vpcClient.UnassociateEipAddress(request)
	if err != nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to unassociate eip",
			err.Error(),
		)
	}

	return nil
}

// ListEIPs 列出弹性 IP
func (p *AlibabaProvider) ListEIPs(ctx context.Context, filter cloudprovider.EIPFilter) ([]*cloudprovider.EIP, error) {
	request := vpc.CreateDescribeEipAddressesRequest()
	request.Scheme = "https"
	request.RegionId = p.regionID

	if filter.EIPID != "" {
		request.AllocationId = filter.EIPID
	}

	response, err := p.vpcClient.DescribeEipAddresses(request)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to describe eip addresses",
			err.Error(),
		)
	}

	var eips []*cloudprovider.EIP
	for _, eip := range response.EipAddresses.EipAddress {
		bandwidth, _ := strconv.Atoi(eip.Bandwidth)
		eips = append(eips, &cloudprovider.EIP{
			ID:         eip.AllocationId,
			Address:    eip.IpAddress,
			Bandwidth:  bandwidth,
			Status:     eip.Status,
			ResourceID: eip.InstanceId,
			RegionID:   p.regionID,
		})
	}

	return eips, nil
}

// UpdateSubnet 更新子网属性
func (p *AlibabaProvider) UpdateSubnet(ctx context.Context, subnetID, name, description string, tags map[string]string) (*cloudprovider.Subnet, error) {
	// 修改 VSwitch 名称
	if name != "" {
		request := vpc.CreateModifyVSwitchAttributeRequest()
		request.Scheme = "https"
		request.VSwitchId = subnetID
		request.VSwitchName = name
		request.Description = description

		_, err := p.vpcClient.ModifyVSwitchAttribute(request)
		if err != nil {
			return nil, cloudprovider.NewCloudError(
				cloudprovider.ErrOperationFailed,
				"failed to modify vswitch attribute",
				err.Error(),
			)
		}
	}

	// 获取更新后的子网信息
	return p.GetSubnet(ctx, subnetID)
}

// AddSecurityGroupRule 添加安全组规则
func (p *AlibabaProvider) AddSecurityGroupRule(ctx context.Context, sgID string, rule cloudprovider.SGRule) (string, error) {
	request := ecs.CreateAuthorizeSecurityGroupRequest()
	request.Scheme = "https"
	request.SecurityGroupId = sgID
	request.IpProtocol = rule.Protocol
	request.PortRange = rule.PortRange
	request.SourceCidrIp = rule.CIDR
	request.Policy = rule.Action
	request.Description = rule.Description

	_, err := p.ecsClient.AuthorizeSecurityGroup(request)
	if err != nil {
		return "", cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to authorize security group rule",
			err.Error(),
		)
	}

	// 阿里云没有返回规则ID，返回一个基于规则的标识
	return sgID + "-" + rule.Protocol + "-" + rule.PortRange, nil
}

// DeleteSecurityGroupRule 删除安全组规则
func (p *AlibabaProvider) DeleteSecurityGroupRule(ctx context.Context, sgID, ruleID string) error {
	// 需要先获取规则详情才能删除
	request := ecs.CreateRevokeSecurityGroupRequest()
	request.Scheme = "https"
	request.SecurityGroupId = sgID

	// ruleID 格式为 "sgId-protocol-portRange"，解析出协议和端口范围
	parts := strings.Split(ruleID, "-")
	if len(parts) >= 3 {
		request.IpProtocol = parts[1]
		request.PortRange = parts[2]
	}

	_, err := p.ecsClient.RevokeSecurityGroup(request)
	if err != nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to revoke security group rule",
			err.Error(),
		)
	}

	return nil
}

// BindEIP 绑定弹性IP
func (p *AlibabaProvider) BindEIP(ctx context.Context, eipID, resourceID, resourceType string) error {
	request := vpc.CreateAssociateEipAddressRequest()
	request.Scheme = "https"
	request.AllocationId = eipID
	request.InstanceId = resourceID
	// 阿里云支持: EcsInstance, SlbInstance, NatGateway, HaVip
	// resourceType 会被自动处理

	_, err := p.vpcClient.AssociateEipAddress(request)
	if err != nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to bind eip",
			err.Error(),
		)
	}

	return nil
}

// UnbindEIP 解绑弹性IP
func (p *AlibabaProvider) UnbindEIP(ctx context.Context, eipID string) error {
	request := vpc.CreateUnassociateEipAddressRequest()
	request.Scheme = "https"
	request.AllocationId = eipID

	_, err := p.vpcClient.UnassociateEipAddress(request)
	if err != nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to unbind eip",
			err.Error(),
		)
	}

	return nil
}
