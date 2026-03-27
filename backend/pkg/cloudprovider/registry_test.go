package cloudprovider

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterProvider(t *testing.T) {
	// 清理注册表
	originalRegistry := providerRegistry
	providerRegistry = make(map[string]ProviderFactory)
	defer func() { providerRegistry = originalRegistry }()

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
	originalRegistry := providerRegistry
	providerRegistry = make(map[string]ProviderFactory)
	defer func() { providerRegistry = originalRegistry }()

	// 注册测试提供者
	expectedProvider := &testProvider{}
	factory := func(config CloudAccountConfig) (ICloudProvider, error) {
		return expectedProvider, nil
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
	assert.Equal(t, expectedProvider, provider)
}

func TestGetProvider_NotFound(t *testing.T) {
	// 清理注册表
	originalRegistry := providerRegistry
	providerRegistry = make(map[string]ProviderFactory)
	defer func() { providerRegistry = originalRegistry }()

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
	originalRegistry := providerRegistry
	providerRegistry = make(map[string]ProviderFactory)
	defer func() { providerRegistry = originalRegistry }()

	// 注册测试提供者
	RegisterProvider("test1", func(config CloudAccountConfig) (ICloudProvider, error) { return nil, nil })
	RegisterProvider("test2", func(config CloudAccountConfig) (ICloudProvider, error) { return nil, nil })

	providers := ListProviders()
	assert.Len(t, providers, 2)
	assert.Contains(t, providers, "test1")
	assert.Contains(t, providers, "test2")
}

// testProvider 用于测试的最小实现
type testProvider struct{}

func (p *testProvider) GetCloudInfo() CloudInfo                                    { return CloudInfo{} }
func (p *testProvider) ListRegions() ([]*Region, error)                            { return nil, nil }
func (p *testProvider) ListZones(regionID string) ([]*Zone, error)                 { return nil, nil }
func (p *testProvider) ListInstanceTypes(regionID string) ([]*InstanceType, error) { return nil, nil }

// ICompute
func (p *testProvider) CreateVM(ctx context.Context, config VMCreateConfig) (*VirtualMachine, error) {
	return nil, nil
}
func (p *testProvider) DeleteVM(ctx context.Context, vmID string) error { return nil }
func (p *testProvider) StartVM(ctx context.Context, vmID string) error  { return nil }
func (p *testProvider) StopVM(ctx context.Context, vmID string) error   { return nil }
func (p *testProvider) RebootVM(ctx context.Context, vmID string) error { return nil }
func (p *testProvider) GetVMStatus(ctx context.Context, vmID string) (*VMStatus, error) {
	return nil, nil
}
func (p *testProvider) ListVMs(ctx context.Context, filter VMListFilter) ([]*VirtualMachine, error) {
	return nil, nil
}
func (p *testProvider) ListImages(ctx context.Context, filter ImageFilter) ([]*Image, error) {
	return nil, nil
}
func (p *testProvider) GetImage(ctx context.Context, imageID string) (*Image, error) { return nil, nil }
func (p *testProvider) CreateKeypair(ctx context.Context, name, publicKey string) (*Keypair, error) {
	return nil, nil
}
func (p *testProvider) DeleteKeypair(ctx context.Context, keypairID string) error { return nil }
func (p *testProvider) ListKeypairs(ctx context.Context) ([]*Keypair, error)      { return nil, nil }

// INetwork
func (p *testProvider) CreateVPC(ctx context.Context, config VPCConfig) (*VPC, error) {
	return nil, nil
}
func (p *testProvider) DeleteVPC(ctx context.Context, vpcID string) error      { return nil }
func (p *testProvider) GetVPC(ctx context.Context, vpcID string) (*VPC, error) { return nil, nil }
func (p *testProvider) ListVPCs(ctx context.Context, filter VPCFilter) ([]*VPC, error) {
	return nil, nil
}
func (p *testProvider) CreateSubnet(ctx context.Context, config SubnetConfig) (*Subnet, error) {
	return nil, nil
}
func (p *testProvider) DeleteSubnet(ctx context.Context, subnetID string) error { return nil }
func (p *testProvider) GetSubnet(ctx context.Context, subnetID string) (*Subnet, error) {
	return nil, nil
}
func (p *testProvider) ListSubnets(ctx context.Context, filter SubnetFilter) ([]*Subnet, error) {
	return nil, nil
}
func (p *testProvider) CreateSecurityGroup(ctx context.Context, config SGConfig) (*SecurityGroup, error) {
	return nil, nil
}
func (p *testProvider) DeleteSecurityGroup(ctx context.Context, sgID string) error { return nil }
func (p *testProvider) AuthorizeSecurityGroup(ctx context.Context, sgID string, rules []SGRule) error {
	return nil
}
func (p *testProvider) RevokeSecurityGroup(ctx context.Context, sgID string, rules []SGRule) error {
	return nil
}
func (p *testProvider) ListSecurityGroups(ctx context.Context, filter SGFilter) ([]*SecurityGroup, error) {
	return nil, nil
}
func (p *testProvider) AllocateEIP(ctx context.Context, config EIPConfig) (*EIP, error) {
	return nil, nil
}
func (p *testProvider) ReleaseEIP(ctx context.Context, eipID string) error               { return nil }
func (p *testProvider) AssociateEIP(ctx context.Context, eipID, resourceID string) error { return nil }
func (p *testProvider) DissociateEIP(ctx context.Context, eipID string) error            { return nil }
func (p *testProvider) ListEIPs(ctx context.Context, filter EIPFilter) ([]*EIP, error) {
	return nil, nil
}
func (p *testProvider) CreateLoadBalancer(ctx context.Context, config LBConfig) (*LoadBalancer, error) {
	return nil, nil
}
func (p *testProvider) DeleteLoadBalancer(ctx context.Context, lbID string) error { return nil }
func (p *testProvider) CreateListener(ctx context.Context, lbID string, config ListenerConfig) (*Listener, error) {
	return nil, nil
}
func (p *testProvider) DeleteListener(ctx context.Context, listenerID string) error { return nil }
func (p *testProvider) ListLoadBalancers(ctx context.Context, filter LBFilter) ([]*LoadBalancer, error) {
	return nil, nil
}
func (p *testProvider) CreateDNSZone(ctx context.Context, config DNSZoneConfig) (*DNSZone, error) {
	return nil, nil
}
func (p *testProvider) DeleteDNSZone(ctx context.Context, zoneID string) error { return nil }
func (p *testProvider) CreateDNSRecord(ctx context.Context, zoneID string, config DNSRecordConfig) (*DNSRecord, error) {
	return nil, nil
}
func (p *testProvider) DeleteDNSRecord(ctx context.Context, recordID string) error { return nil }
func (p *testProvider) ListDNSZones(ctx context.Context) ([]*DNSZone, error)       { return nil, nil }
func (p *testProvider) ListDNSRecords(ctx context.Context, zoneID string) ([]*DNSRecord, error) {
	return nil, nil
}

// IStorage
func (p *testProvider) CreateDisk(ctx context.Context, config DiskConfig) (*Disk, error) {
	return nil, nil
}
func (p *testProvider) DeleteDisk(ctx context.Context, diskID string) error             { return nil }
func (p *testProvider) AttachDisk(ctx context.Context, diskID, vmID string) error       { return nil }
func (p *testProvider) DetachDisk(ctx context.Context, diskID string) error             { return nil }
func (p *testProvider) ResizeDisk(ctx context.Context, diskID string, sizeGB int) error { return nil }
func (p *testProvider) CreateSnapshot(ctx context.Context, diskID string, name string) (*Snapshot, error) {
	return nil, nil
}
func (p *testProvider) DeleteSnapshot(ctx context.Context, snapshotID string) error { return nil }
func (p *testProvider) ListDisks(ctx context.Context, filter DiskFilter) ([]*Disk, error) {
	return nil, nil
}
func (p *testProvider) ListSnapshots(ctx context.Context, filter SnapshotFilter) ([]*Snapshot, error) {
	return nil, nil
}
func (p *testProvider) CreateBucket(ctx context.Context, name string, config BucketConfig) error {
	return nil
}
func (p *testProvider) DeleteBucket(ctx context.Context, name string) error { return nil }
func (p *testProvider) ListBuckets(ctx context.Context) ([]*Bucket, error)  { return nil, nil }
func (p *testProvider) PutObject(ctx context.Context, bucketName, objectKey string, data io.Reader) error {
	return nil
}
func (p *testProvider) GetObject(ctx context.Context, bucketName, objectKey string) (io.ReadCloser, error) {
	return nil, nil
}
func (p *testProvider) DeleteObject(ctx context.Context, bucketName, objectKey string) error {
	return nil
}
func (p *testProvider) ListObjects(ctx context.Context, bucketName string, prefix string) ([]*Object, error) {
	return nil, nil
}
func (p *testProvider) CreateFileSystem(ctx context.Context, config FSConfig) (*FileSystem, error) {
	return nil, nil
}
func (p *testProvider) DeleteFileSystem(ctx context.Context, fsID string) error { return nil }
func (p *testProvider) MountFileSystem(ctx context.Context, fsID string, config MountConfig) (*MountTarget, error) {
	return nil, nil
}
func (p *testProvider) UnmountFileSystem(ctx context.Context, mountID string) error { return nil }
func (p *testProvider) ListFileSystems(ctx context.Context) ([]*FileSystem, error)  { return nil, nil }

// IDatabase
func (p *testProvider) CreateRDSInstance(ctx context.Context, config RDSConfig) (*RDSInstance, error) {
	return nil, nil
}
func (p *testProvider) DeleteRDSInstance(ctx context.Context, instanceID string) error { return nil }
func (p *testProvider) StartRDSInstance(ctx context.Context, instanceID string) error  { return nil }
func (p *testProvider) StopRDSInstance(ctx context.Context, instanceID string) error   { return nil }
func (p *testProvider) RebootRDSInstance(ctx context.Context, instanceID string) error { return nil }
func (p *testProvider) ResizeRDSInstance(ctx context.Context, instanceID string, spec RDSpec) error {
	return nil
}
func (p *testProvider) CreateRDSBackup(ctx context.Context, instanceID string, name string) (*RDSBackup, error) {
	return nil, nil
}
func (p *testProvider) RestoreRDSFromBackup(ctx context.Context, backupID string, config RDSConfig) (*RDSInstance, error) {
	return nil, nil
}
func (p *testProvider) ListRDSInstances(ctx context.Context, filter RDSFilter) ([]*RDSInstance, error) {
	return nil, nil
}
func (p *testProvider) ListRDSBackups(ctx context.Context, instanceID string) ([]*RDSBackup, error) {
	return nil, nil
}
func (p *testProvider) CreateCacheInstance(ctx context.Context, config CacheConfig) (*CacheInstance, error) {
	return nil, nil
}
func (p *testProvider) DeleteCacheInstance(ctx context.Context, instanceID string) error { return nil }
func (p *testProvider) RebootCacheInstance(ctx context.Context, instanceID string) error { return nil }
func (p *testProvider) ResizeCacheInstance(ctx context.Context, instanceID string, spec CacheSpec) error {
	return nil
}
func (p *testProvider) CreateCacheBackup(ctx context.Context, instanceID string) (*CacheBackup, error) {
	return nil, nil
}
func (p *testProvider) ListCacheInstances(ctx context.Context, filter CacheFilter) ([]*CacheInstance, error) {
	return nil, nil
}
