package azure

import (
	"context"
	"io"

	"github.com/opencmp/opencmp/pkg/cloudprovider"
)

// AzureProvider Azure 提供商
type AzureProvider struct {
	config         cloudprovider.CloudAccountConfig
	subscriptionID string
}

// NewAzureProvider 创建 Azure 提供商
func NewAzureProvider(config cloudprovider.CloudAccountConfig) (cloudprovider.ICloudProvider, error) {
	tenantID := config.Credentials["tenant_id"]
	clientID := config.Credentials["client_id"]
	clientSecret := config.Credentials["client_secret"]
	subscriptionID := config.Credentials["subscription_id"]

	// TODO: 初始化 Azure SDK 客户端
	_ = tenantID
	_ = clientID
	_ = clientSecret

	return &AzureProvider{
		config:         config,
		subscriptionID: subscriptionID,
	}, nil
}

// GetCloudInfo 获取云厂商信息
func (p *AzureProvider) GetCloudInfo() cloudprovider.CloudInfo {
	return cloudprovider.CloudInfo{
		Provider: "azure",
		Version:  "1.0.0",
		Services: []string{"VM", "VNet", "SQL", "Blob"},
	}
}

func (p *AzureProvider) ListRegions() ([]*cloudprovider.Region, error) {
	return []*cloudprovider.Region{}, nil
}

func (p *AzureProvider) ListZones(regionID string) ([]*cloudprovider.Zone, error) {
	return []*cloudprovider.Zone{}, nil
}

func (p *AzureProvider) ListInstanceTypes(regionID string) ([]*cloudprovider.InstanceType, error) {
	return []*cloudprovider.InstanceType{}, nil
}

// ICompute
func (p *AzureProvider) CreateVM(ctx context.Context, config cloudprovider.VMCreateConfig) (*cloudprovider.VirtualMachine, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) DeleteVM(ctx context.Context, vmID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) StartVM(ctx context.Context, vmID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) StopVM(ctx context.Context, vmID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) RebootVM(ctx context.Context, vmID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) GetVMStatus(ctx context.Context, vmID string) (*cloudprovider.VMStatus, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) ListVMs(ctx context.Context, filter cloudprovider.VMListFilter) ([]*cloudprovider.VirtualMachine, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) ListImages(ctx context.Context, filter cloudprovider.ImageFilter) ([]*cloudprovider.Image, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) GetImage(ctx context.Context, imageID string) (*cloudprovider.Image, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) CreateKeypair(ctx context.Context, name, publicKey string) (*cloudprovider.Keypair, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) DeleteKeypair(ctx context.Context, keypairID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) ListKeypairs(ctx context.Context) ([]*cloudprovider.Keypair, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

// INetwork
func (p *AzureProvider) CreateVPC(ctx context.Context, config cloudprovider.VPCConfig) (*cloudprovider.VPC, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) DeleteVPC(ctx context.Context, vpcID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) GetVPC(ctx context.Context, vpcID string) (*cloudprovider.VPC, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) ListVPCs(ctx context.Context, filter cloudprovider.VPCFilter) ([]*cloudprovider.VPC, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) CreateSubnet(ctx context.Context, config cloudprovider.SubnetConfig) (*cloudprovider.Subnet, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) DeleteSubnet(ctx context.Context, subnetID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) GetSubnet(ctx context.Context, subnetID string) (*cloudprovider.Subnet, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) ListSubnets(ctx context.Context, filter cloudprovider.SubnetFilter) ([]*cloudprovider.Subnet, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) CreateSecurityGroup(ctx context.Context, config cloudprovider.SGConfig) (*cloudprovider.SecurityGroup, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) DeleteSecurityGroup(ctx context.Context, sgID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) AuthorizeSecurityGroup(ctx context.Context, sgID string, rules []cloudprovider.SGRule) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) RevokeSecurityGroup(ctx context.Context, sgID string, rules []cloudprovider.SGRule) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) ListSecurityGroups(ctx context.Context, filter cloudprovider.SGFilter) ([]*cloudprovider.SecurityGroup, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) AllocateEIP(ctx context.Context, config cloudprovider.EIPConfig) (*cloudprovider.EIP, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) ReleaseEIP(ctx context.Context, eipID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) AssociateEIP(ctx context.Context, eipID, resourceID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) DissociateEIP(ctx context.Context, eipID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) ListEIPs(ctx context.Context, filter cloudprovider.EIPFilter) ([]*cloudprovider.EIP, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) CreateLoadBalancer(ctx context.Context, config cloudprovider.LBConfig) (*cloudprovider.LoadBalancer, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) DeleteLoadBalancer(ctx context.Context, lbID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) CreateListener(ctx context.Context, lbID string, config cloudprovider.ListenerConfig) (*cloudprovider.Listener, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) DeleteListener(ctx context.Context, listenerID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) ListLoadBalancers(ctx context.Context, filter cloudprovider.LBFilter) ([]*cloudprovider.LoadBalancer, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) CreateDNSZone(ctx context.Context, config cloudprovider.DNSZoneConfig) (*cloudprovider.DNSZone, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) DeleteDNSZone(ctx context.Context, zoneID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) CreateDNSRecord(ctx context.Context, zoneID string, config cloudprovider.DNSRecordConfig) (*cloudprovider.DNSRecord, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) DeleteDNSRecord(ctx context.Context, recordID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) ListDNSZones(ctx context.Context) ([]*cloudprovider.DNSZone, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) ListDNSRecords(ctx context.Context, zoneID string) ([]*cloudprovider.DNSRecord, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

// IStorage
func (p *AzureProvider) CreateDisk(ctx context.Context, config cloudprovider.DiskConfig) (*cloudprovider.Disk, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) DeleteDisk(ctx context.Context, diskID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) AttachDisk(ctx context.Context, diskID, vmID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) DetachDisk(ctx context.Context, diskID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) ResizeDisk(ctx context.Context, diskID string, sizeGB int) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) CreateSnapshot(ctx context.Context, diskID string, name string) (*cloudprovider.Snapshot, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) DeleteSnapshot(ctx context.Context, snapshotID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) ListDisks(ctx context.Context, filter cloudprovider.DiskFilter) ([]*cloudprovider.Disk, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) ListSnapshots(ctx context.Context, filter cloudprovider.SnapshotFilter) ([]*cloudprovider.Snapshot, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) CreateBucket(ctx context.Context, name string, config cloudprovider.BucketConfig) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) DeleteBucket(ctx context.Context, name string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) ListBuckets(ctx context.Context) ([]*cloudprovider.Bucket, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) PutObject(ctx context.Context, bucketName, objectKey string, data io.Reader) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) GetObject(ctx context.Context, bucketName, objectKey string) (io.ReadCloser, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) DeleteObject(ctx context.Context, bucketName, objectKey string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) ListObjects(ctx context.Context, bucketName string, prefix string) ([]*cloudprovider.Object, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) CreateFileSystem(ctx context.Context, config cloudprovider.FSConfig) (*cloudprovider.FileSystem, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) DeleteFileSystem(ctx context.Context, fsID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) MountFileSystem(ctx context.Context, fsID string, config cloudprovider.MountConfig) (*cloudprovider.MountTarget, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) UnmountFileSystem(ctx context.Context, mountID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) ListFileSystems(ctx context.Context) ([]*cloudprovider.FileSystem, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

// IDatabase
func (p *AzureProvider) CreateRDSInstance(ctx context.Context, config cloudprovider.RDSConfig) (*cloudprovider.RDSInstance, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) DeleteRDSInstance(ctx context.Context, instanceID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) StartRDSInstance(ctx context.Context, instanceID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) StopRDSInstance(ctx context.Context, instanceID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) RebootRDSInstance(ctx context.Context, instanceID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) ResizeRDSInstance(ctx context.Context, instanceID string, spec cloudprovider.RDSpec) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) CreateRDSBackup(ctx context.Context, instanceID string, name string) (*cloudprovider.RDSBackup, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) RestoreRDSFromBackup(ctx context.Context, backupID string, config cloudprovider.RDSConfig) (*cloudprovider.RDSInstance, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) ListRDSInstances(ctx context.Context, filter cloudprovider.RDSFilter) ([]*cloudprovider.RDSInstance, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) ListRDSBackups(ctx context.Context, instanceID string) ([]*cloudprovider.RDSBackup, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) CreateCacheInstance(ctx context.Context, config cloudprovider.CacheConfig) (*cloudprovider.CacheInstance, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) DeleteCacheInstance(ctx context.Context, instanceID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) RebootCacheInstance(ctx context.Context, instanceID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) ResizeCacheInstance(ctx context.Context, instanceID string, spec cloudprovider.CacheSpec) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) CreateCacheBackup(ctx context.Context, instanceID string) (*cloudprovider.CacheBackup, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *AzureProvider) ListCacheInstances(ctx context.Context, filter cloudprovider.CacheFilter) ([]*cloudprovider.CacheInstance, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func init() {
	cloudprovider.RegisterProvider("azure", NewAzureProvider)
}
