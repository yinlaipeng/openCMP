package tencent

import (
	"context"
	"io"

	"github.com/opencmp/opencmp/pkg/cloudprovider"
)

// TencentProvider 腾讯云提供商
type TencentProvider struct {
	config   cloudprovider.CloudAccountConfig
	regionID string
}

// NewTencentProvider 创建腾讯云提供商
func NewTencentProvider(config cloudprovider.CloudAccountConfig) (cloudprovider.ICloudProvider, error) {
	secretID := config.Credentials["secret_id"]
	secretKey := config.Credentials["secret_key"]
	regionID := config.Region

	// TODO: 初始化腾讯云 SDK 客户端
	_ = secretID
	_ = secretKey

	return &TencentProvider{
		config:   config,
		regionID: regionID,
	}, nil
}

// GetCloudInfo 获取云厂商信息
func (p *TencentProvider) GetCloudInfo() cloudprovider.CloudInfo {
	return cloudprovider.CloudInfo{
		Provider: "tencent",
		Version:  "1.0.0",
		Services: []string{"CVM", "VPC", "CDB", "COS"},
	}
}

func (p *TencentProvider) ListRegions() ([]*cloudprovider.Region, error) {
	return []*cloudprovider.Region{}, nil
}

func (p *TencentProvider) ListZones(regionID string) ([]*cloudprovider.Zone, error) {
	return []*cloudprovider.Zone{}, nil
}

func (p *TencentProvider) ListInstanceTypes(regionID string) ([]*cloudprovider.InstanceType, error) {
	return []*cloudprovider.InstanceType{}, nil
}

// ICompute
func (p *TencentProvider) CreateVM(ctx context.Context, config cloudprovider.VMCreateConfig) (*cloudprovider.VirtualMachine, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) DeleteVM(ctx context.Context, vmID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) StartVM(ctx context.Context, vmID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) StopVM(ctx context.Context, vmID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) RebootVM(ctx context.Context, vmID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) GetVMStatus(ctx context.Context, vmID string) (*cloudprovider.VMStatus, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) ListVMs(ctx context.Context, filter cloudprovider.VMListFilter) ([]*cloudprovider.VirtualMachine, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) ListImages(ctx context.Context, filter cloudprovider.ImageFilter) ([]*cloudprovider.Image, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) GetImage(ctx context.Context, imageID string) (*cloudprovider.Image, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) CreateKeypair(ctx context.Context, name, publicKey string) (*cloudprovider.Keypair, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) DeleteKeypair(ctx context.Context, keypairID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) ListKeypairs(ctx context.Context) ([]*cloudprovider.Keypair, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

// INetwork
func (p *TencentProvider) CreateVPC(ctx context.Context, config cloudprovider.VPCConfig) (*cloudprovider.VPC, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) DeleteVPC(ctx context.Context, vpcID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) GetVPC(ctx context.Context, vpcID string) (*cloudprovider.VPC, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) ListVPCs(ctx context.Context, filter cloudprovider.VPCFilter) ([]*cloudprovider.VPC, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) CreateSubnet(ctx context.Context, config cloudprovider.SubnetConfig) (*cloudprovider.Subnet, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) DeleteSubnet(ctx context.Context, subnetID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) GetSubnet(ctx context.Context, subnetID string) (*cloudprovider.Subnet, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) ListSubnets(ctx context.Context, filter cloudprovider.SubnetFilter) ([]*cloudprovider.Subnet, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) CreateSecurityGroup(ctx context.Context, config cloudprovider.SGConfig) (*cloudprovider.SecurityGroup, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) DeleteSecurityGroup(ctx context.Context, sgID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) AuthorizeSecurityGroup(ctx context.Context, sgID string, rules []cloudprovider.SGRule) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) RevokeSecurityGroup(ctx context.Context, sgID string, rules []cloudprovider.SGRule) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) ListSecurityGroups(ctx context.Context, filter cloudprovider.SGFilter) ([]*cloudprovider.SecurityGroup, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) AllocateEIP(ctx context.Context, config cloudprovider.EIPConfig) (*cloudprovider.EIP, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) ReleaseEIP(ctx context.Context, eipID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) AssociateEIP(ctx context.Context, eipID, resourceID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) DissociateEIP(ctx context.Context, eipID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) ListEIPs(ctx context.Context, filter cloudprovider.EIPFilter) ([]*cloudprovider.EIP, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) CreateLoadBalancer(ctx context.Context, config cloudprovider.LBConfig) (*cloudprovider.LoadBalancer, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) DeleteLoadBalancer(ctx context.Context, lbID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) CreateListener(ctx context.Context, lbID string, config cloudprovider.ListenerConfig) (*cloudprovider.Listener, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) DeleteListener(ctx context.Context, listenerID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) ListLoadBalancers(ctx context.Context, filter cloudprovider.LBFilter) ([]*cloudprovider.LoadBalancer, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) CreateDNSZone(ctx context.Context, config cloudprovider.DNSZoneConfig) (*cloudprovider.DNSZone, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) DeleteDNSZone(ctx context.Context, zoneID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) CreateDNSRecord(ctx context.Context, zoneID string, config cloudprovider.DNSRecordConfig) (*cloudprovider.DNSRecord, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) DeleteDNSRecord(ctx context.Context, recordID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) ListDNSZones(ctx context.Context) ([]*cloudprovider.DNSZone, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) ListDNSRecords(ctx context.Context, zoneID string) ([]*cloudprovider.DNSRecord, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

// IStorage
func (p *TencentProvider) CreateDisk(ctx context.Context, config cloudprovider.DiskConfig) (*cloudprovider.Disk, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) DeleteDisk(ctx context.Context, diskID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) AttachDisk(ctx context.Context, diskID, vmID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) DetachDisk(ctx context.Context, diskID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) ResizeDisk(ctx context.Context, diskID string, sizeGB int) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) CreateSnapshot(ctx context.Context, diskID string, name string) (*cloudprovider.Snapshot, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) DeleteSnapshot(ctx context.Context, snapshotID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) ListDisks(ctx context.Context, filter cloudprovider.DiskFilter) ([]*cloudprovider.Disk, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) ListSnapshots(ctx context.Context, filter cloudprovider.SnapshotFilter) ([]*cloudprovider.Snapshot, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) CreateBucket(ctx context.Context, name string, config cloudprovider.BucketConfig) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) DeleteBucket(ctx context.Context, name string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) ListBuckets(ctx context.Context) ([]*cloudprovider.Bucket, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) PutObject(ctx context.Context, bucketName, objectKey string, data io.Reader) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) GetObject(ctx context.Context, bucketName, objectKey string) (io.ReadCloser, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) DeleteObject(ctx context.Context, bucketName, objectKey string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) ListObjects(ctx context.Context, bucketName string, prefix string) ([]*cloudprovider.Object, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) CreateFileSystem(ctx context.Context, config cloudprovider.FSConfig) (*cloudprovider.FileSystem, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) DeleteFileSystem(ctx context.Context, fsID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) MountFileSystem(ctx context.Context, fsID string, config cloudprovider.MountConfig) (*cloudprovider.MountTarget, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) UnmountFileSystem(ctx context.Context, mountID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) ListFileSystems(ctx context.Context) ([]*cloudprovider.FileSystem, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

// IDatabase
func (p *TencentProvider) CreateRDSInstance(ctx context.Context, config cloudprovider.RDSConfig) (*cloudprovider.RDSInstance, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) DeleteRDSInstance(ctx context.Context, instanceID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) StartRDSInstance(ctx context.Context, instanceID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) StopRDSInstance(ctx context.Context, instanceID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) RebootRDSInstance(ctx context.Context, instanceID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) ResizeRDSInstance(ctx context.Context, instanceID string, spec cloudprovider.RDSpec) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) CreateRDSBackup(ctx context.Context, instanceID string, name string) (*cloudprovider.RDSBackup, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) RestoreRDSFromBackup(ctx context.Context, backupID string, config cloudprovider.RDSConfig) (*cloudprovider.RDSInstance, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) ListRDSInstances(ctx context.Context, filter cloudprovider.RDSFilter) ([]*cloudprovider.RDSInstance, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) ListRDSBackups(ctx context.Context, instanceID string) ([]*cloudprovider.RDSBackup, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) CreateCacheInstance(ctx context.Context, config cloudprovider.CacheConfig) (*cloudprovider.CacheInstance, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) DeleteCacheInstance(ctx context.Context, instanceID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) RebootCacheInstance(ctx context.Context, instanceID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) ResizeCacheInstance(ctx context.Context, instanceID string, spec cloudprovider.CacheSpec) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) CreateCacheBackup(ctx context.Context, instanceID string) (*cloudprovider.CacheBackup, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func (p *TencentProvider) ListCacheInstances(ctx context.Context, filter cloudprovider.CacheFilter) ([]*cloudprovider.CacheInstance, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "not implemented", "")
}

func init() {
	cloudprovider.RegisterProvider("tencent", NewTencentProvider)
}
