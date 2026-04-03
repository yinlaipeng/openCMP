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
