package cloudprovider

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterProvider(t *testing.T) {
	// 清理注册表
	providerRegistry = make(map[string]ProviderFactory)

	// 测试注册
	factory := func(config CloudAccountConfig) (ICloudProvider, error) {
		return nil, nil
	}

	RegisterProvider("test", factory)

	// 验证注册成功
	assert.Contains(t, providerRegistry, "test")
}

func TestGetProvider(t *testing.T) {
	// 清理注册表
	providerRegistry = make(map[string]ProviderFactory)

	// 注册测试提供者
	factory := func(config CloudAccountConfig) (ICloudProvider, error) {
		return &mockProvider{}, nil
	}
	RegisterProvider("test", factory)

	// 测试获取
	config := CloudAccountConfig{
		ID:           "1",
		Name:         "test",
		ProviderType: "test",
		Credentials:  map[string]string{},
		Region:       "",
	}

	provider, err := GetProvider("test", config)
	assert.NoError(t, err)
	assert.NotNil(t, provider)
}

func TestGetProvider_NotFound(t *testing.T) {
	// 清理注册表
	providerRegistry = make(map[string]ProviderFactory)

	config := CloudAccountConfig{
		ID:           "1",
		Name:         "test",
		ProviderType: "unknown",
		Credentials:  map[string]string{},
		Region:       "",
	}

	provider, err := GetProvider("unknown", config)
	assert.Error(t, err)
	assert.Nil(t, provider)
	assert.Contains(t, err.Error(), "not found")
}

func TestListProviders(t *testing.T) {
	// 清理注册表
	providerRegistry = make(map[string]ProviderFactory)

	// 注册测试提供者
	RegisterProvider("test1", func(config CloudAccountConfig) (ICloudProvider, error) { return nil, nil })
	RegisterProvider("test2", func(config CloudAccountConfig) (ICloudProvider, error) { return nil, nil })

	providers := ListProviders()
	assert.Len(t, providers, 2)
	assert.Contains(t, providers, "test1")
	assert.Contains(t, providers, "test2")
}

// mockProvider 用于测试
type mockProvider struct{}

func (m *mockProvider) CreateVM(ctx context.Context, config VMCreateConfig) (*VirtualMachine, error) {
	return nil, nil
}
func (m *mockProvider) DeleteVM(ctx context.Context, vmID string) error { return nil }
func (m *mockProvider) StartVM(ctx context.Context, vmID string) error  { return nil }
func (m *mockProvider) StopVM(ctx context.Context, vmID string) error   { return nil }
func (m *mockProvider) RebootVM(ctx context.Context, vmID string) error { return nil }
func (m *mockProvider) GetVMStatus(ctx context.Context, vmID string) (*VMStatus, error) {
	return nil, nil
}
func (m *mockProvider) ListVMs(ctx context.Context, filter VMListFilter) ([]*VirtualMachine, error) {
	return nil, nil
}
func (m *mockProvider) ListImages(ctx context.Context, filter ImageFilter) ([]*Image, error) {
	return nil, nil
}
func (m *mockProvider) GetImage(ctx context.Context, imageID string) (*Image, error) {
	return nil, nil
}
func (m *mockProvider) CreateKeypair(ctx context.Context, name, publicKey string) (*Keypair, error) {
	return nil, nil
}
func (m *mockProvider) DeleteKeypair(ctx context.Context, keypairID string) error { return nil }
func (m *mockProvider) ListKeypairs(ctx context.Context) ([]*Keypair, error)      { return nil, nil }
func (m *mockProvider) CreateVPC(ctx context.Context, config VPCConfig) (*VPC, error) {
	return nil, nil
}
func (m *mockProvider) DeleteVPC(ctx context.Context, vpcID string) error { return nil }
func (m *mockProvider) GetVPC(ctx context.Context, vpcID string) (*VPC, error) {
	return nil, nil
}
func (m *mockProvider) ListVPCs(ctx context.Context, filter VPCFilter) ([]*VPC, error) {
	return nil, nil
}
func (m *mockProvider) CreateSubnet(ctx context.Context, config SubnetConfig) (*Subnet, error) {
	return nil, nil
}
func (m *mockProvider) DeleteSubnet(ctx context.Context, subnetID string) error { return nil }
func (m *mockProvider) GetSubnet(ctx context.Context, subnetID string) (*Subnet, error) {
	return nil, nil
}
func (m *mockProvider) ListSubnets(ctx context.Context, filter SubnetFilter) ([]*Subnet, error) {
	return nil, nil
}
func (m *mockProvider) CreateSecurityGroup(ctx context.Context, config SGConfig) (*SecurityGroup, error) {
	return nil, nil
}
func (m *mockProvider) DeleteSecurityGroup(ctx context.Context, sgID string) error { return nil }
func (m *mockProvider) AuthorizeSecurityGroup(ctx context.Context, sgID string, rules []SGRule) error {
	return nil
}
func (m *mockProvider) RevokeSecurityGroup(ctx context.Context, sgID string, rules []SGRule) error {
	return nil
}
func (m *mockProvider) ListSecurityGroups(ctx context.Context, filter SGFilter) ([]*SecurityGroup, error) {
	return nil, nil
}
func (m *mockProvider) AllocateEIP(ctx context.Context, config EIPConfig) (*EIP, error) {
	return nil, nil
}
func (m *mockProvider) ReleaseEIP(ctx context.Context, eipID string) error { return nil }
func (m *mockProvider) AssociateEIP(ctx context.Context, eipID, resourceID string) error {
	return nil
}
func (m *mockProvider) DissociateEIP(ctx context.Context, eipID string) error { return nil }
func (m *mockProvider) ListEIPs(ctx context.Context, filter EIPFilter) ([]*EIP, error) {
	return nil, nil
}
func (m *mockProvider) CreateDisk(ctx context.Context, config DiskConfig) (*Disk, error) {
	return nil, nil
}
func (m *mockProvider) DeleteDisk(ctx context.Context, diskID string) error { return nil }
func (m *mockProvider) AttachDisk(ctx context.Context, diskID, vmID string) error {
	return nil
}
func (m *mockProvider) DetachDisk(ctx context.Context, diskID string) error { return nil }
func (m *mockProvider) ResizeDisk(ctx context.Context, diskID string, sizeGB int) error {
	return nil
}
func (m *mockProvider) CreateSnapshot(ctx context.Context, diskID string, name string) (*Snapshot, error) {
	return nil, nil
}
func (m *mockProvider) DeleteSnapshot(ctx context.Context, snapshotID string) error {
	return nil
}
func (m *mockProvider) ListDisks(ctx context.Context, filter DiskFilter) ([]*Disk, error) {
	return nil, nil
}
func (m *mockProvider) ListSnapshots(ctx context.Context, filter SnapshotFilter) ([]*Snapshot, error) {
	return nil, nil
}
func (m *mockProvider) GetCloudInfo() CloudInfo                    { return CloudInfo{} }
func (m *mockProvider) ListRegions() ([]*Region, error)            { return nil, nil }
func (m *mockProvider) ListZones(regionID string) ([]*Zone, error) { return nil, nil }
func (m *mockProvider) ListInstanceTypes(regionID string) ([]*InstanceType, error) {
	return nil, nil
}
