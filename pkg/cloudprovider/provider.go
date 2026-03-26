package cloudprovider

// ICloudProvider 云提供商总接口
type ICloudProvider interface {
	ICompute  // 计算资源
	INetwork  // 网络资源
	IStorage  // 存储资源
	IDatabase // 数据服务

	// 云厂商信息
	GetCloudInfo() CloudInfo
	// 区域和可用区列表
	ListRegions() ([]*Region, error)
	ListZones(regionID string) ([]*Zone, error)
	// 实例规格列表
	ListInstanceTypes(regionID string) ([]*InstanceType, error)
}
