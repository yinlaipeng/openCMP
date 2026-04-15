package azure

import (
	"context"
	"io"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v4"

	"github.com/opencmp/opencmp/pkg/cloudprovider"
)

// AzureProvider Azure 提供商
type AzureProvider struct {
	config         cloudprovider.CloudAccountConfig
	subscriptionID string
	cred           *azidentity.ClientSecretCredential
	vmClient       *armcompute.VirtualMachinesClient
	vnetClient     *armnetwork.VirtualNetworksClient
	subnetClient   *armnetwork.SubnetsClient
	nsgClient      *armnetwork.SecurityGroupsClient
	ipClient       *armnetwork.PublicIPAddressesClient
	location       string
}

// NewAzureProvider 创建 Azure 提供商
func NewAzureProvider(config cloudprovider.CloudAccountConfig) (cloudprovider.ICloudProvider, error) {
	tenantID := config.Credentials["tenant_id"]
	clientID := config.Credentials["client_id"]
	clientSecret := config.Credentials["client_secret"]
	subscriptionID := config.Credentials["subscription_id"]
	location := config.Credentials["location"]
	if location == "" {
		location = "eastus" // 默认区域
	}

	// 创建 Azure 身份认证
	cred, err := azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, nil)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to create Azure credential",
			err.Error(),
		)
	}

	// 创建虚拟机客户端
	vmClient, err := armcompute.NewVirtualMachinesClient(subscriptionID, cred, nil)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to create VM client",
			err.Error(),
		)
	}

	// 创建虚拟网络客户端
	vnetClient, err := armnetwork.NewVirtualNetworksClient(subscriptionID, cred, nil)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to create VNet client",
			err.Error(),
		)
	}

	// 创建子网客户端
	subnetClient, err := armnetwork.NewSubnetsClient(subscriptionID, cred, nil)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to create Subnet client",
			err.Error(),
		)
	}

	// 创建安全组客户端
	nsgClient, err := armnetwork.NewSecurityGroupsClient(subscriptionID, cred, nil)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to create NSG client",
			err.Error(),
		)
	}

	// 创建公网IP客户端
	ipClient, err := armnetwork.NewPublicIPAddressesClient(subscriptionID, cred, nil)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to create PublicIP client",
			err.Error(),
		)
	}

	return &AzureProvider{
		config:         config,
		subscriptionID: subscriptionID,
		cred:           cred,
		vmClient:       vmClient,
		vnetClient:     vnetClient,
		subnetClient:   subnetClient,
		nsgClient:      nsgClient,
		ipClient:       ipClient,
		location:       location,
	}, nil
}

// GetCloudInfo 获取云厂商信息
func (p *AzureProvider) GetCloudInfo() cloudprovider.CloudInfo {
	return cloudprovider.CloudInfo{
		Provider: "azure",
		Version:  "1.0.0",
		Services: []string{"VM", "VNet", "Subnet", "NSG", "PublicIP"},
	}
}

// parseResourceID 从 Azure 资源 ID 解析资源组名称
// 格式: /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/virtualMachines/{vmName}
func parseResourceGroup(resourceID string) (string, error) {
	// 简化解析 - 从配置中获取或使用默认值
	return "opencmp-resources", nil // 默认资源组
}

// ListRegions 列出区域
func (p *AzureProvider) ListRegions() ([]*cloudprovider.Region, error) {
	// Azure 区域列表
	regions := []*cloudprovider.Region{
		{ID: "eastus", Name: "East US"},
		{ID: "eastus2", Name: "East US 2"},
		{ID: "westus", Name: "West US"},
		{ID: "westus2", Name: "West US 2"},
		{ID: "centralus", Name: "Central US"},
		{ID: "northcentralus", Name: "North Central US"},
		{ID: "southcentralus", Name: "South Central US"},
		{ID: "westcentralus", Name: "West Central US"},
		{ID: "canadacentral", Name: "Canada Central"},
		{ID: "canadaeast", Name: "Canada East"},
		{ID: "brazilsouth", Name: "Brazil South"},
		{ID: "northeurope", Name: "North Europe"},
		{ID: "westeurope", Name: "West Europe"},
		{ID: "uksouth", Name: "UK South"},
		{ID: "ukwest", Name: "UK West"},
		{ID: "francecentral", Name: "France Central"},
		{ID: "francesouth", Name: "France South"},
		{ID: "switzerlandnorth", Name: "Switzerland North"},
		{ID: "switzerlandwest", Name: "Switzerland West"},
		{ID: "germanynorth", Name: "Germany North"},
		{ID: "germanywestcentral", Name: "Germany West Central"},
		{ID: "norwayeast", Name: "Norway East"},
		{ID: "norwaywest", Name: "Norway West"},
		{ID: "swedencentral", Name: "Sweden Central"},
		{ID: "australiaeast", Name: "Australia East"},
		{ID: "australiacentral", Name: "Australia Central"},
		{ID: "australiacentral2", Name: "Australia Central 2"},
		{ID: "australiasoutheast", Name: "Australia Southeast"},
		{ID: "eastasia", Name: "East Asia"},
		{ID: "southeastasia", Name: "Southeast Asia"},
		{ID: "japaneast", Name: "Japan East"},
		{ID: "japanwest", Name: "Japan West"},
		{ID: "koreacentral", Name: "Korea Central"},
		{ID: "koreasouth", Name: "Korea South"},
		{ID: "indiacentral", Name: "India Central"},
		{ID: "indiawest", Name: "India West"},
		{ID: "southafricawest", Name: "South Africa West"},
		{ID: "southafricawest", Name: "South Africa West"},
		{ID: "uaenorth", Name: "UAE North"},
		{ID: "uaecentral", Name: "UAE Central"},
	}
	return regions, nil
}

// ListZones 列出可用区
func (p *AzureProvider) ListZones(regionID string) ([]*cloudprovider.Zone, error) {
	// Azure 大多数区域没有显式的可用区概念
	return []*cloudprovider.Zone{
		{ID: "1", Name: "Zone 1", RegionID: regionID},
		{ID: "2", Name: "Zone 2", RegionID: regionID},
		{ID: "3", Name: "Zone 3", RegionID: regionID},
	}, nil
}

// ListInstanceTypes 列出实例类型
func (p *AzureProvider) ListInstanceTypes(regionID string) ([]*cloudprovider.InstanceType, error) {
	instanceTypes := []*cloudprovider.InstanceType{
		{Name: "Standard_B1s", CPU: 1, Memory: 1, Category: "Burstable"},
		{Name: "Standard_B1ms", CPU: 1, Memory: 2, Category: "Burstable"},
		{Name: "Standard_B2s", CPU: 2, Memory: 4, Category: "Burstable"},
		{Name: "Standard_D2s_v3", CPU: 2, Memory: 8, Category: "General Purpose"},
		{Name: "Standard_D4s_v3", CPU: 4, Memory: 16, Category: "General Purpose"},
		{Name: "Standard_D8s_v3", CPU: 8, Memory: 32, Category: "General Purpose"},
		{Name: "Standard_E2s_v3", CPU: 2, Memory: 16, Category: "Memory Optimized"},
		{Name: "Standard_E4s_v3", CPU: 4, Memory: 32, Category: "Memory Optimized"},
		{Name: "Standard_F2s_v2", CPU: 2, Memory: 4, Category: "Compute Optimized"},
		{Name: "Standard_F4s_v2", CPU: 4, Memory: 8, Category: "Compute Optimized"},
	}
	return instanceTypes, nil
}

func init() {
	cloudprovider.RegisterProvider("azure", NewAzureProvider)
}

// ============================================
// IStorage 接口 stub 实现 (待后续完善)
// ============================================

func (p *AzureProvider) CreateDisk(ctx context.Context, config cloudprovider.DiskConfig) (*cloudprovider.Disk, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "CreateDisk not implemented", "")
}

func (p *AzureProvider) DeleteDisk(ctx context.Context, diskID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "DeleteDisk not implemented", "")
}

func (p *AzureProvider) AttachDisk(ctx context.Context, diskID, vmID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "AttachDisk not implemented", "")
}

func (p *AzureProvider) DetachDisk(ctx context.Context, diskID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "DetachDisk not implemented", "")
}

func (p *AzureProvider) ResizeDisk(ctx context.Context, diskID string, sizeGB int) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "ResizeDisk not implemented", "")
}

func (p *AzureProvider) CreateSnapshot(ctx context.Context, diskID string, name string) (*cloudprovider.Snapshot, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "CreateSnapshot not implemented", "")
}

func (p *AzureProvider) DeleteSnapshot(ctx context.Context, snapshotID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "DeleteSnapshot not implemented", "")
}

func (p *AzureProvider) ListDisks(ctx context.Context, filter cloudprovider.DiskFilter) ([]*cloudprovider.Disk, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "ListDisks not implemented", "")
}

func (p *AzureProvider) ListSnapshots(ctx context.Context, filter cloudprovider.SnapshotFilter) ([]*cloudprovider.Snapshot, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "ListSnapshots not implemented", "")
}

func (p *AzureProvider) CreateBucket(ctx context.Context, name string, config cloudprovider.BucketConfig) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "CreateBucket not implemented", "")
}

func (p *AzureProvider) DeleteBucket(ctx context.Context, name string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "DeleteBucket not implemented", "")
}

func (p *AzureProvider) ListBuckets(ctx context.Context) ([]*cloudprovider.Bucket, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "ListBuckets not implemented", "")
}

func (p *AzureProvider) PutObject(ctx context.Context, bucketName, objectKey string, data io.Reader) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "PutObject not implemented", "")
}

func (p *AzureProvider) GetObject(ctx context.Context, bucketName, objectKey string) (io.ReadCloser, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "GetObject not implemented", "")
}

func (p *AzureProvider) DeleteObject(ctx context.Context, bucketName, objectKey string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "DeleteObject not implemented", "")
}

func (p *AzureProvider) ListObjects(ctx context.Context, bucketName string, prefix string) ([]*cloudprovider.Object, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "ListObjects not implemented", "")
}

func (p *AzureProvider) CreateFileSystem(ctx context.Context, config cloudprovider.FSConfig) (*cloudprovider.FileSystem, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "CreateFileSystem not implemented", "")
}

func (p *AzureProvider) DeleteFileSystem(ctx context.Context, fsID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "DeleteFileSystem not implemented", "")
}

func (p *AzureProvider) MountFileSystem(ctx context.Context, fsID string, config cloudprovider.MountConfig) (*cloudprovider.MountTarget, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "MountFileSystem not implemented", "")
}

func (p *AzureProvider) UnmountFileSystem(ctx context.Context, mountID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "UnmountFileSystem not implemented", "")
}

func (p *AzureProvider) ListFileSystems(ctx context.Context) ([]*cloudprovider.FileSystem, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "ListFileSystems not implemented", "")
}

// ============================================
// IDatabase 接口 stub 实现
// ============================================

func (p *AzureProvider) CreateRDSInstance(ctx context.Context, config cloudprovider.RDSConfig) (*cloudprovider.RDSInstance, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "CreateRDSInstance not implemented", "")
}

func (p *AzureProvider) DeleteRDSInstance(ctx context.Context, instanceID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "DeleteRDSInstance not implemented", "")
}

func (p *AzureProvider) StartRDSInstance(ctx context.Context, instanceID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "StartRDSInstance not implemented", "")
}

func (p *AzureProvider) StopRDSInstance(ctx context.Context, instanceID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "StopRDSInstance not implemented", "")
}

func (p *AzureProvider) RebootRDSInstance(ctx context.Context, instanceID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "RebootRDSInstance not implemented", "")
}

func (p *AzureProvider) ResizeRDSInstance(ctx context.Context, instanceID string, spec cloudprovider.RDSpec) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "ResizeRDSInstance not implemented", "")
}

func (p *AzureProvider) CreateRDSBackup(ctx context.Context, instanceID string, name string) (*cloudprovider.RDSBackup, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "CreateRDSBackup not implemented", "")
}

func (p *AzureProvider) RestoreRDSFromBackup(ctx context.Context, backupID string, config cloudprovider.RDSConfig) (*cloudprovider.RDSInstance, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "RestoreRDSFromBackup not implemented", "")
}

func (p *AzureProvider) ListRDSInstances(ctx context.Context, filter cloudprovider.RDSFilter) ([]*cloudprovider.RDSInstance, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "ListRDSInstances not implemented", "")
}

func (p *AzureProvider) ListRDSBackups(ctx context.Context, instanceID string) ([]*cloudprovider.RDSBackup, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "ListRDSBackups not implemented", "")
}

func (p *AzureProvider) CreateCacheInstance(ctx context.Context, config cloudprovider.CacheConfig) (*cloudprovider.CacheInstance, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "CreateCacheInstance not implemented", "")
}

func (p *AzureProvider) DeleteCacheInstance(ctx context.Context, instanceID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "DeleteCacheInstance not implemented", "")
}

func (p *AzureProvider) RebootCacheInstance(ctx context.Context, instanceID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "RebootCacheInstance not implemented", "")
}

func (p *AzureProvider) ResizeCacheInstance(ctx context.Context, instanceID string, spec cloudprovider.CacheSpec) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "ResizeCacheInstance not implemented", "")
}

func (p *AzureProvider) CreateCacheBackup(ctx context.Context, instanceID string) (*cloudprovider.CacheBackup, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "CreateCacheBackup not implemented", "")
}

func (p *AzureProvider) ListCacheInstances(ctx context.Context, filter cloudprovider.CacheFilter) ([]*cloudprovider.CacheInstance, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "ListCacheInstances not implemented", "")
}

// ============================================
// Advanced network methods stub 实现
// ============================================

func (p *AzureProvider) CreateLoadBalancer(ctx context.Context, config cloudprovider.LBConfig) (*cloudprovider.LoadBalancer, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "CreateLoadBalancer not implemented", "")
}

func (p *AzureProvider) DeleteLoadBalancer(ctx context.Context, lbID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "DeleteLoadBalancer not implemented", "")
}

func (p *AzureProvider) CreateListener(ctx context.Context, lbID string, config cloudprovider.ListenerConfig) (*cloudprovider.Listener, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "CreateListener not implemented", "")
}

func (p *AzureProvider) DeleteListener(ctx context.Context, listenerID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "DeleteListener not implemented", "")
}

func (p *AzureProvider) ListLoadBalancers(ctx context.Context, filter cloudprovider.LBFilter) ([]*cloudprovider.LoadBalancer, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "ListLoadBalancers not implemented", "")
}

func (p *AzureProvider) CreateDNSZone(ctx context.Context, config cloudprovider.DNSZoneConfig) (*cloudprovider.DNSZone, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "CreateDNSZone not implemented", "")
}

func (p *AzureProvider) DeleteDNSZone(ctx context.Context, zoneID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "DeleteDNSZone not implemented", "")
}

func (p *AzureProvider) CreateDNSRecord(ctx context.Context, zoneID string, config cloudprovider.DNSRecordConfig) (*cloudprovider.DNSRecord, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "CreateDNSRecord not implemented", "")
}

func (p *AzureProvider) DeleteDNSRecord(ctx context.Context, recordID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "DeleteDNSRecord not implemented", "")
}

func (p *AzureProvider) ListDNSZones(ctx context.Context) ([]*cloudprovider.DNSZone, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "ListDNSZones not implemented", "")
}

func (p *AzureProvider) ListDNSRecords(ctx context.Context, zoneID string) ([]*cloudprovider.DNSRecord, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "ListDNSRecords not implemented", "")
}

func (p *AzureProvider) CreateVPCInterconnect(ctx context.Context, config cloudprovider.VPCInterconnectConfig) (*cloudprovider.VPCInterconnect, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "CreateVPCInterconnect not implemented", "")
}

func (p *AzureProvider) DeleteVPCInterconnect(ctx context.Context, interconnectID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "DeleteVPCInterconnect not implemented", "")
}

func (p *AzureProvider) ListVPCInterconnects(ctx context.Context, filter cloudprovider.VPCInterconnectFilter) ([]*cloudprovider.VPCInterconnect, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "ListVPCInterconnects not implemented", "")
}

func (p *AzureProvider) CreateVPCPeering(ctx context.Context, config cloudprovider.VPCPeeringConfig) (*cloudprovider.VPCPeering, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "CreateVPCPeering not implemented", "")
}

func (p *AzureProvider) DeleteVPCPeering(ctx context.Context, peeringID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "DeleteVPCPeering not implemented", "")
}

func (p *AzureProvider) ListVPCPeerings(ctx context.Context, filter cloudprovider.VPCPeeringFilter) ([]*cloudprovider.VPCPeering, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "ListVPCPeerings not implemented", "")
}

func (p *AzureProvider) CreateRouteTable(ctx context.Context, config cloudprovider.RouteTableConfig) (*cloudprovider.RouteTable, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "CreateRouteTable not implemented", "")
}

func (p *AzureProvider) DeleteRouteTable(ctx context.Context, routeTableID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "DeleteRouteTable not implemented", "")
}

func (p *AzureProvider) ListRouteTables(ctx context.Context, filter cloudprovider.RouteTableFilter) ([]*cloudprovider.RouteTable, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "ListRouteTables not implemented", "")
}

func (p *AzureProvider) CreateL2Network(ctx context.Context, config cloudprovider.L2NetworkConfig) (*cloudprovider.L2Network, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "CreateL2Network not implemented", "")
}

func (p *AzureProvider) DeleteL2Network(ctx context.Context, l2NetworkID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "DeleteL2Network not implemented", "")
}

func (p *AzureProvider) ListL2Networks(ctx context.Context, filter cloudprovider.L2NetworkFilter) ([]*cloudprovider.L2Network, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "ListL2Networks not implemented", "")
}

// ============================================
// Extended network methods (未实现)
// ============================================

func (p *AzureProvider) UpdateSubnet(ctx context.Context, subnetID, name, description string, tags map[string]string) (*cloudprovider.Subnet, error) {
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "UpdateSubnet not implemented", "")
}

func (p *AzureProvider) AddSecurityGroupRule(ctx context.Context, sgID string, rule cloudprovider.SGRule) (string, error) {
	return "", cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "AddSecurityGroupRule not implemented", "")
}

func (p *AzureProvider) DeleteSecurityGroupRule(ctx context.Context, sgID, ruleID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "DeleteSecurityGroupRule not implemented", "")
}

func (p *AzureProvider) BindEIP(ctx context.Context, eipID, resourceID, resourceType string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "BindEIP not implemented", "")
}

func (p *AzureProvider) UnbindEIP(ctx context.Context, eipID string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "UnbindEIP not implemented", "")
}