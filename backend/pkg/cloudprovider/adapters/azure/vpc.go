package azure

import (
	"context"
	"fmt"
	"strings"
		"time"


	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v4"

	"github.com/opencmp/opencmp/pkg/cloudprovider"
)

// CreateVPC 创建虚拟网络 (Azure VNet)
func (p *AzureProvider) CreateVPC(ctx context.Context, config cloudprovider.VPCConfig) (*cloudprovider.VPC, error) {
	resourceGroup := "opencmp-resources"

	parameters := armnetwork.VirtualNetwork{
		Location: to.Ptr(p.location),
		Properties: &armnetwork.VirtualNetworkPropertiesFormat{
			AddressSpace: &armnetwork.AddressSpace{
				AddressPrefixes: []*string{
					to.Ptr(config.CIDR),
				},
			},
			Subnets: []*armnetwork.Subnet{
				{
					Name: to.Ptr("default-subnet"),
					Properties: &armnetwork.SubnetPropertiesFormat{
						AddressPrefixes: []*string{
							to.Ptr(config.CIDR),
						},
					},
				},
			},
		},
	}

	poller, err := p.vnetClient.BeginCreateOrUpdate(ctx, resourceGroup, fmt.Sprintf("eip-%d", time.Now().UnixNano()), parameters, nil)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to create VNet",
			err.Error(),
		)
	}

	resp, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to wait for VNet creation",
			err.Error(),
		)
	}

	return p.convertAzureVNetToCloudVPC(resp.VirtualNetwork), nil
}

// DeleteVPC 删除虚拟网络
func (p *AzureProvider) DeleteVPC(ctx context.Context, vpcID string) error {
	resourceGroup, vnetName := p.parseVNetID(vpcID)

	poller, err := p.vnetClient.BeginDelete(ctx, resourceGroup, vnetName, nil)
	if err != nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to delete VNet",
			err.Error(),
		)
	}

	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to wait for VNet deletion",
			err.Error(),
		)
	}

	return nil
}

// GetVPC 获取虚拟网络
func (p *AzureProvider) GetVPC(ctx context.Context, vpcID string) (*cloudprovider.VPC, error) {
	resourceGroup, vnetName := p.parseVNetID(vpcID)

	resp, err := p.vnetClient.Get(ctx, resourceGroup, vnetName, nil)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrResourceNotFound,
			"failed to get VNet",
			err.Error(),
		)
	}

	return p.convertAzureVNetToCloudVPC(resp.VirtualNetwork), nil
}

// ListVPCs 列出虚拟网络
func (p *AzureProvider) ListVPCs(ctx context.Context, filter cloudprovider.VPCFilter) ([]*cloudprovider.VPC, error) {
	resourceGroup := "opencmp-resources"

	pager := p.vnetClient.NewListPager(resourceGroup, nil)

	var vpcs []*cloudprovider.VPC
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, cloudprovider.NewCloudError(
				cloudprovider.ErrOperationFailed,
				"failed to list VNets",
				err.Error(),
			)
		}

		for _, vnet := range page.VirtualNetworkListResult.Value {
			vpcs = append(vpcs, p.convertAzureVNetToCloudVPC(*vnet))
		}
	}

	return vpcs, nil
}

// CreateSubnet 创建子网
func (p *AzureProvider) CreateSubnet(ctx context.Context, config cloudprovider.SubnetConfig) (*cloudprovider.Subnet, error) {
	resourceGroup, vnetName := p.parseVNetID(config.VPCID)

	parameters := armnetwork.Subnet{
		Properties: &armnetwork.SubnetPropertiesFormat{
			AddressPrefixes: []*string{
				to.Ptr(config.CIDR),
			},
		},
	}

	poller, err := p.subnetClient.BeginCreateOrUpdate(ctx, resourceGroup, vnetName, fmt.Sprintf("eip-%d", time.Now().UnixNano()), parameters, nil)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to create Subnet",
			err.Error(),
		)
	}

	resp, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to wait for Subnet creation",
			err.Error(),
		)
	}

	return p.convertAzureSubnetToCloudSubnet(resp.Subnet, config.VPCID), nil
}

// DeleteSubnet 删除子网
func (p *AzureProvider) DeleteSubnet(ctx context.Context, subnetID string) error {
	resourceGroup, vnetName, subnetName := p.parseSubnetID(subnetID)

	poller, err := p.subnetClient.BeginDelete(ctx, resourceGroup, vnetName, subnetName, nil)
	if err != nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to delete Subnet",
			err.Error(),
		)
	}

	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to wait for Subnet deletion",
			err.Error(),
		)
	}

	return nil
}

// GetSubnet 获取子网
func (p *AzureProvider) GetSubnet(ctx context.Context, subnetID string) (*cloudprovider.Subnet, error) {
	resourceGroup, vnetName, subnetName := p.parseSubnetID(subnetID)

	resp, err := p.subnetClient.Get(ctx, resourceGroup, vnetName, subnetName, nil)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrResourceNotFound,
			"failed to get Subnet",
			err.Error(),
		)
	}

	vpcID := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualNetworks/%s",
		p.subscriptionID, resourceGroup, vnetName)

	return p.convertAzureSubnetToCloudSubnet(resp.Subnet, vpcID), nil
}

// ListSubnets 列出子网
func (p *AzureProvider) ListSubnets(ctx context.Context, filter cloudprovider.SubnetFilter) ([]*cloudprovider.Subnet, error) {
	var subnets []*cloudprovider.Subnet

	if filter.VPCID != "" {
		resourceGroup, vnetName := p.parseVNetID(filter.VPCID)

		pager := p.subnetClient.NewListPager(resourceGroup, vnetName, nil)
		for pager.More() {
			page, err := pager.NextPage(ctx)
			if err != nil {
				return nil, cloudprovider.NewCloudError(
					cloudprovider.ErrOperationFailed,
					"failed to list Subnets",
					err.Error(),
				)
			}

			for _, subnet := range page.SubnetListResult.Value {
				subnets = append(subnets, p.convertAzureSubnetToCloudSubnet(*subnet, filter.VPCID))
			}
		}
	}

	return subnets, nil
}

// CreateSecurityGroup 创建安全组 (Azure NSG)
func (p *AzureProvider) CreateSecurityGroup(ctx context.Context, config cloudprovider.SGConfig) (*cloudprovider.SecurityGroup, error) {
	resourceGroup := "opencmp-resources"

	parameters := armnetwork.SecurityGroup{
		Location: to.Ptr(p.location),
		Properties: &armnetwork.SecurityGroupPropertiesFormat{
			SecurityRules: []*armnetwork.SecurityRule{},
		},
	}

	poller, err := p.nsgClient.BeginCreateOrUpdate(ctx, resourceGroup, fmt.Sprintf("eip-%d", time.Now().UnixNano()), parameters, nil)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to create NSG",
			err.Error(),
		)
	}

	resp, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to wait for NSG creation",
			err.Error(),
		)
	}

	return p.convertAzureNSGToCloudSG(resp.SecurityGroup), nil
}

// DeleteSecurityGroup 删除安全组
func (p *AzureProvider) DeleteSecurityGroup(ctx context.Context, sgID string) error {
	resourceGroup, nsgName := p.parseNSGID(sgID)

	poller, err := p.nsgClient.BeginDelete(ctx, resourceGroup, nsgName, nil)
	if err != nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to delete NSG",
			err.Error(),
		)
	}

	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to wait for NSG deletion",
			err.Error(),
		)
	}

	return nil
}

// AuthorizeSecurityGroup 授权安全组规则
func (p *AzureProvider) AuthorizeSecurityGroup(ctx context.Context, sgID string, rules []cloudprovider.SGRule) error {
	resourceGroup, nsgName := p.parseNSGID(sgID)

	// 获取现有安全组
	nsg, err := p.nsgClient.Get(ctx, resourceGroup, nsgName, nil)
	if err != nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrResourceNotFound,
			"failed to get NSG",
			err.Error(),
		)
	}

	// 添加新规则
	for i, rule := range rules {
		azureRule := armnetwork.SecurityRule{
			Name: to.Ptr(fmt.Sprintf("rule-%d", i+1)),
			Properties: &armnetwork.SecurityRulePropertiesFormat{
				Direction:                  p.mapRuleDirection(rule.Direction),
				Protocol:                   p.mapRuleProtocol(rule.Protocol),
				SourceAddressPrefix:        to.Ptr(rule.CIDR),
				SourcePortRange:            to.Ptr(rule.PortRange),
				DestinationAddressPrefix:   to.Ptr(rule.CIDR),
				DestinationPortRange:       to.Ptr(rule.PortRange),
				Access:                     to.Ptr(armnetwork.SecurityRuleAccessAllow),
				Priority:                   to.Ptr(int32(100 + i)),
			},
		}
		nsg.SecurityGroup.Properties.SecurityRules = append(
			nsg.SecurityGroup.Properties.SecurityRules,
			&azureRule,
		)
	}

	poller, err := p.nsgClient.BeginCreateOrUpdate(ctx, resourceGroup, nsgName, nsg.SecurityGroup, nil)
	if err != nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to update NSG",
			err.Error(),
		)
	}

	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to wait for NSG update",
			err.Error(),
		)
	}

	return nil
}

// RevokeSecurityGroup 撤销安全组规则
func (p *AzureProvider) RevokeSecurityGroup(ctx context.Context, sgID string, rules []cloudprovider.SGRule) error {
	// Azure 需要删除特定规则，简化处理：重建安全组
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "revoke rules requires rule IDs", "")
}

// ListSecurityGroups 列出安全组
func (p *AzureProvider) ListSecurityGroups(ctx context.Context, filter cloudprovider.SGFilter) ([]*cloudprovider.SecurityGroup, error) {

	pager := p.nsgClient.NewListAllPager(nil)

	var sgs []*cloudprovider.SecurityGroup
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, cloudprovider.NewCloudError(
				cloudprovider.ErrOperationFailed,
				"failed to list NSGs",
				err.Error(),
			)
		}

		for _, nsg := range page.SecurityGroupListResult.Value {
			sgs = append(sgs, p.convertAzureNSGToCloudSG(*nsg))
		}
	}

	return sgs, nil
}

// AllocateEIP 分配公网 IP
func (p *AzureProvider) AllocateEIP(ctx context.Context, config cloudprovider.EIPConfig) (*cloudprovider.EIP, error) {
	resourceGroup := "opencmp-resources"

	parameters := armnetwork.PublicIPAddress{
		Location: to.Ptr(p.location),
		Properties: &armnetwork.PublicIPAddressPropertiesFormat{
			PublicIPAllocationMethod: to.Ptr(armnetwork.IPAllocationMethodStatic),
		},
	}

	poller, err := p.ipClient.BeginCreateOrUpdate(ctx, resourceGroup, fmt.Sprintf("eip-%d", time.Now().UnixNano()), parameters, nil)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to allocate PublicIP",
			err.Error(),
		)
	}

	resp, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to wait for PublicIP allocation",
			err.Error(),
		)
	}

	return p.convertAzurePIPToCloudEIP(resp.PublicIPAddress), nil
}

// ReleaseEIP 释放公网 IP
func (p *AzureProvider) ReleaseEIP(ctx context.Context, eipID string) error {
	resourceGroup, pipName := p.parsePIPID(eipID)

	poller, err := p.ipClient.BeginDelete(ctx, resourceGroup, pipName, nil)
	if err != nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to release PublicIP",
			err.Error(),
		)
	}

	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to wait for PublicIP release",
			err.Error(),
		)
	}

	return nil
}

// AssociateEIP 绑定公网 IP
func (p *AzureProvider) AssociateEIP(ctx context.Context, eipID, resourceID string) error {
	// Azure 需要更新网络接口的 PublicIPAddress 配置
	// 这需要额外的 NetworkInterface 客户端，简化处理
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "EIP association requires NIC update", "")
}

// DissociateEIP 解绑公网 IP
func (p *AzureProvider) DissociateEIP(ctx context.Context, eipID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "EIP dissociation requires NIC update", "")
}

// ListEIPs 列出公网 IP
func (p *AzureProvider) ListEIPs(ctx context.Context, filter cloudprovider.EIPFilter) ([]*cloudprovider.EIP, error) {
	resourceGroup := "opencmp-resources"

	pager := p.ipClient.NewListPager(resourceGroup, nil)

	var eips []*cloudprovider.EIP
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, cloudprovider.NewCloudError(
				cloudprovider.ErrOperationFailed,
				"failed to list PublicIPs",
				err.Error(),
			)
		}

		for _, pip := range page.PublicIPAddressListResult.Value {
			eips = append(eips, p.convertAzurePIPToCloudEIP(*pip))
		}
	}

	return eips, nil
}

// 解析函数
func (p *AzureProvider) parseVNetID(vnetID string) (resourceGroup, vnetName string) {
	if strings.Contains(vnetID, "/virtualNetworks/") {
		parts := strings.Split(vnetID, "/")
		for i, part := range parts {
			if part == "resourceGroups" && i+1 < len(parts) {
				resourceGroup = parts[i+1]
			}
			if part == "virtualNetworks" && i+1 < len(parts) {
				vnetName = parts[i+1]
			}
		}
	} else if strings.Contains(vnetID, ":") {
		parts := strings.Split(vnetID, ":")
		resourceGroup = parts[0]
		vnetName = parts[1]
	} else {
		resourceGroup = "opencmp-resources"
		vnetName = vnetID
	}
	return resourceGroup, vnetName
}

func (p *AzureProvider) parseSubnetID(subnetID string) (resourceGroup, vnetName, subnetName string) {
	if strings.Contains(subnetID, "/subnets/") {
		parts := strings.Split(subnetID, "/")
		for i, part := range parts {
			if part == "resourceGroups" && i+1 < len(parts) {
				resourceGroup = parts[i+1]
			}
			if part == "virtualNetworks" && i+1 < len(parts) {
				vnetName = parts[i+1]
			}
			if part == "subnets" && i+1 < len(parts) {
				subnetName = parts[i+1]
			}
		}
	} else if strings.Contains(subnetID, ":") {
		parts := strings.Split(subnetID, ":")
		if len(parts) >= 3 {
			resourceGroup = parts[0]
			vnetName = parts[1]
			subnetName = parts[2]
		}
	} else {
		resourceGroup = "opencmp-resources"
		vnetName = "default-vnet"
		subnetName = subnetID
	}
	return resourceGroup, vnetName, subnetName
}

func (p *AzureProvider) parseNSGID(nsgID string) (resourceGroup, nsgName string) {
	if strings.Contains(nsgID, "/networkSecurityGroups/") {
		parts := strings.Split(nsgID, "/")
		for i, part := range parts {
			if part == "resourceGroups" && i+1 < len(parts) {
				resourceGroup = parts[i+1]
			}
			if part == "networkSecurityGroups" && i+1 < len(parts) {
				nsgName = parts[i+1]
			}
		}
	} else {
		resourceGroup = "opencmp-resources"
		nsgName = nsgID
	}
	return resourceGroup, nsgName
}

func (p *AzureProvider) parsePIPID(pipID string) (resourceGroup, pipName string) {
	if strings.Contains(pipID, "/publicIPAddresses/") {
		parts := strings.Split(pipID, "/")
		for i, part := range parts {
			if part == "resourceGroups" && i+1 < len(parts) {
				resourceGroup = parts[i+1]
			}
			if part == "publicIPAddresses" && i+1 < len(parts) {
				pipName = parts[i+1]
			}
		}
	} else {
		resourceGroup = "opencmp-resources"
		pipName = pipID
	}
	return resourceGroup, pipName
}

// 转换函数
func (p *AzureProvider) convertAzureVNetToCloudVPC(vnet armnetwork.VirtualNetwork) *cloudprovider.VPC {
	result := &cloudprovider.VPC{
		ID:       *vnet.ID,
		Name:     *vnet.Name,
		RegionID: *vnet.Location,
	}

	if vnet.Properties != nil && vnet.Properties.AddressSpace != nil {
		for _, prefix := range vnet.Properties.AddressSpace.AddressPrefixes {
			if prefix != nil {
				result.CIDR = *prefix
				break
			}
		}
	}

	if vnet.Properties != nil && vnet.Properties.ProvisioningState != nil {
		result.Status = string(*vnet.Properties.ProvisioningState)
	}

	return result
}

func (p *AzureProvider) convertAzureSubnetToCloudSubnet(subnet armnetwork.Subnet, vpcID string) *cloudprovider.Subnet {
	result := &cloudprovider.Subnet{
		ID:     *subnet.ID,
		Name:   *subnet.Name,
		VPCID:  vpcID,
		Status: "Available",
	}

	if subnet.Properties != nil {
		for _, prefix := range subnet.Properties.AddressPrefixes {
			if prefix != nil {
				result.CIDR = *prefix
				break
			}
		}
	}

	return result
}

func (p *AzureProvider) convertAzureNSGToCloudSG(nsg armnetwork.SecurityGroup) *cloudprovider.SecurityGroup {
	result := &cloudprovider.SecurityGroup{
		ID:       *nsg.ID,
		Name:     *nsg.Name,
	}


		if nsg.Properties != nil {
			for _, rule := range nsg.Properties.SecurityRules {
				if rule.Properties != nil {
					sgRule := cloudprovider.SGRule{
						Direction:  string(*rule.Properties.Direction),
						Protocol:   string(*rule.Properties.Protocol),
					}
					if rule.Properties.SourcePortRange != nil {
						sgRule.PortRange = *rule.Properties.SourcePortRange
					}
					if rule.Properties.DestinationPortRange != nil {
						sgRule.PortRange = *rule.Properties.DestinationPortRange
					}
					if rule.Properties.SourceAddressPrefix != nil {
						sgRule.CIDR = *rule.Properties.SourceAddressPrefix
					}
					result.Rules = append(result.Rules, sgRule)
				}
			}
		}


	return result
}

func (p *AzureProvider) convertAzurePIPToCloudEIP(pip armnetwork.PublicIPAddress) *cloudprovider.EIP {
	result := &cloudprovider.EIP{
		ID:       *pip.ID,
		RegionID: *pip.Location,
	}

	if pip.Properties != nil {
		if pip.Properties.IPAddress != nil {
			result.Address = *pip.Properties.IPAddress
		}
		if pip.Properties.ProvisioningState != nil {
			result.Status = string(*pip.Properties.ProvisioningState)
		}
	}

	return result
}

// 映射函数
func (p *AzureProvider) mapRuleDirection(direction string) *armnetwork.SecurityRuleDirection {
	if direction == "inbound" {
		return to.Ptr(armnetwork.SecurityRuleDirectionInbound)
	}
	return to.Ptr(armnetwork.SecurityRuleDirectionOutbound)
}

func (p *AzureProvider) mapRuleProtocol(protocol string) *armnetwork.SecurityRuleProtocol {
	switch protocol {
	case "tcp":
		return to.Ptr(armnetwork.SecurityRuleProtocolTCP)
	case "udp":
		return to.Ptr(armnetwork.SecurityRuleProtocolUDP)
	case "icmp":
		return to.Ptr(armnetwork.SecurityRuleProtocolIcmp)
	default:
		return to.Ptr(armnetwork.SecurityRuleProtocolAsterisk)
	}
}

func (p *AzureProvider) mapCIDRs(cidrs []string) []*string {
	result := make([]*string, len(cidrs))
	for i, cidr := range cidrs {
		result[i] = to.Ptr(cidr)
	}
	return result
}