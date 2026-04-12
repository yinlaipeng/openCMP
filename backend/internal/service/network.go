package service

import (
	"context"
	"encoding/json"
	"strconv"

	"gorm.io/gorm"

	"github.com/opencmp/opencmp/pkg/cloudprovider"
)

// NetworkService 网络资源服务
type NetworkService struct {
	db *gorm.DB
}

// NewNetworkService 创建网络资源服务
func NewNetworkService(db *gorm.DB) *NetworkService {
	return &NetworkService{db: db}
}

// getProvider 获取云提供商
func (s *NetworkService) getProvider(ctx context.Context, accountID uint) (cloudprovider.ICloudProvider, error) {
	accountService := NewCloudAccountService(s.db)
	account, err := accountService.GetCloudAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, cloudprovider.NewCloudError(cloudprovider.ErrResourceNotFound, "account not found", "")
	}

	var creds map[string]string
	if err := json.Unmarshal(account.Credentials, &creds); err != nil {
		return nil, err
	}

	providerConfig := cloudprovider.CloudAccountConfig{
		ID:           strconv.Itoa(int(account.ID)),
		Name:         account.Name,
		ProviderType: account.ProviderType,
		Credentials:  creds,
		Region:       "",
	}

	return cloudprovider.GetProvider(account.ProviderType, providerConfig)
}

// CreateVPC 创建 VPC
func (s *NetworkService) CreateVPC(ctx context.Context, accountID uint, config cloudprovider.VPCConfig) (*cloudprovider.VPC, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return provider.CreateVPC(ctx, config)
}

// ListVPCs 列出 VPC
func (s *NetworkService) ListVPCs(ctx context.Context, accountID uint, filter cloudprovider.VPCFilter) ([]*cloudprovider.VPC, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return provider.ListVPCs(ctx, filter)
}

// DeleteVPC 删除 VPC
func (s *NetworkService) DeleteVPC(ctx context.Context, accountID uint, vpcID string) error {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return err
	}

	return provider.DeleteVPC(ctx, vpcID)
}

// CreateSubnet 创建子网
func (s *NetworkService) CreateSubnet(ctx context.Context, accountID uint, config cloudprovider.SubnetConfig) (*cloudprovider.Subnet, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return provider.CreateSubnet(ctx, config)
}

// ListSubnets 列出子网
func (s *NetworkService) ListSubnets(ctx context.Context, accountID uint, filter cloudprovider.SubnetFilter) ([]*cloudprovider.Subnet, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return provider.ListSubnets(ctx, filter)
}

// DeleteSubnet 删除子网
func (s *NetworkService) DeleteSubnet(ctx context.Context, accountID uint, subnetID string) error {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return err
	}

	return provider.DeleteSubnet(ctx, subnetID)
}

// CreateSecurityGroup 创建安全组
func (s *NetworkService) CreateSecurityGroup(ctx context.Context, accountID uint, config cloudprovider.SGConfig) (*cloudprovider.SecurityGroup, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return provider.CreateSecurityGroup(ctx, config)
}

// ListSecurityGroups 列出安全组
func (s *NetworkService) ListSecurityGroups(ctx context.Context, accountID uint, filter cloudprovider.SGFilter) ([]*cloudprovider.SecurityGroup, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return provider.ListSecurityGroups(ctx, filter)
}

// CreateEIP 分配弹性 IP
func (s *NetworkService) CreateEIP(ctx context.Context, accountID uint, config cloudprovider.EIPConfig) (*cloudprovider.EIP, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return provider.AllocateEIP(ctx, config)
}

// ListEIPs 列出弹性 IP
func (s *NetworkService) ListEIPs(ctx context.Context, accountID uint, filter cloudprovider.EIPFilter) ([]*cloudprovider.EIP, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return provider.ListEIPs(ctx, filter)
}

// ListRegions 列出区域
func (s *NetworkService) ListRegions(ctx context.Context, accountID uint) ([]*cloudprovider.Region, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return provider.ListRegions()
}

// ListZones 列出可用区
func (s *NetworkService) ListZones(ctx context.Context, accountID uint, regionID string) ([]*cloudprovider.Zone, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return provider.ListZones(regionID)
}

// CreateVPCInterconnect 创建 VPC 互联
func (s *NetworkService) CreateVPCInterconnect(ctx context.Context, accountID uint, config cloudprovider.VPCInterconnectConfig) (*cloudprovider.VPCInterconnect, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return provider.CreateVPCInterconnect(ctx, config)
}

// ListVPCInterconnects 列出 VPC 互联
func (s *NetworkService) ListVPCInterconnects(ctx context.Context, accountID uint, filter cloudprovider.VPCInterconnectFilter) ([]*cloudprovider.VPCInterconnect, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return provider.ListVPCInterconnects(ctx, filter)
}

// DeleteVPCInterconnect 删除 VPC 互联
func (s *NetworkService) DeleteVPCInterconnect(ctx context.Context, accountID uint, interconnectID string) error {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return err
	}

	return provider.DeleteVPCInterconnect(ctx, interconnectID)
}

// CreateVPCPeering 创建 VPC 对等连接
func (s *NetworkService) CreateVPCPeering(ctx context.Context, accountID uint, config cloudprovider.VPCPeeringConfig) (*cloudprovider.VPCPeering, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return provider.CreateVPCPeering(ctx, config)
}

// ListVPCPeerings 列出 VPC 对等连接
func (s *NetworkService) ListVPCPeerings(ctx context.Context, accountID uint, filter cloudprovider.VPCPeeringFilter) ([]*cloudprovider.VPCPeering, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return provider.ListVPCPeerings(ctx, filter)
}

// DeleteVPCPeering 删除 VPC 对等连接
func (s *NetworkService) DeleteVPCPeering(ctx context.Context, accountID uint, peeringID string) error {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return err
	}

	return provider.DeleteVPCPeering(ctx, peeringID)
}

// CreateRouteTable 创建路由表
func (s *NetworkService) CreateRouteTable(ctx context.Context, accountID uint, config cloudprovider.RouteTableConfig) (*cloudprovider.RouteTable, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return provider.CreateRouteTable(ctx, config)
}

// ListRouteTables 列出路由表
func (s *NetworkService) ListRouteTables(ctx context.Context, accountID uint, filter cloudprovider.RouteTableFilter) ([]*cloudprovider.RouteTable, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return provider.ListRouteTables(ctx, filter)
}

// DeleteRouteTable 删除路由表
func (s *NetworkService) DeleteRouteTable(ctx context.Context, accountID uint, routeTableID string) error {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return err
	}

	return provider.DeleteRouteTable(ctx, routeTableID)
}

// CreateL2Network 创建二层网络
func (s *NetworkService) CreateL2Network(ctx context.Context, accountID uint, config cloudprovider.L2NetworkConfig) (*cloudprovider.L2Network, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return provider.CreateL2Network(ctx, config)
}

// ListL2Networks 列出二层网络
func (s *NetworkService) ListL2Networks(ctx context.Context, accountID uint, filter cloudprovider.L2NetworkFilter) ([]*cloudprovider.L2Network, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return provider.ListL2Networks(ctx, filter)
}

// DeleteL2Network 删除二层网络
func (s *NetworkService) DeleteL2Network(ctx context.Context, accountID uint, l2NetworkID string) error {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return err
	}

	return provider.DeleteL2Network(ctx, l2NetworkID)
}
