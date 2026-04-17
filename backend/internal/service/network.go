package service

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
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

// logStateChange 记录资源状态变更日志
func (s *NetworkService) logStateChange(ctx context.Context, log *model.ResourceStateLog) error {
	log.OccurredAt = time.Now()
	log.CreatedAt = time.Now()
	return s.db.WithContext(ctx).Create(log).Error
}

// CreateVPC 创建 VPC
func (s *NetworkService) CreateVPC(ctx context.Context, accountID uint, config cloudprovider.VPCConfig) (*cloudprovider.VPC, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	vpc, err := provider.CreateVPC(ctx, config)
	if err != nil {
		return nil, err
	}

	// 记录状态变更日志
	stateLog := &model.ResourceStateLog{
		ResourceType:   "vpc",
		ResourceID:     vpc.ID,
		ResourceName:   config.Name,
		CloudAccountID: accountID,
		PreviousStatus: "",
		CurrentStatus:  vpc.Status,
		OperationType:  "create",
		Reason:         "VPC创建",
	}
	if err := s.logStateChange(ctx, stateLog); err != nil {
		// 日志记录失败不影响主流程
	}

	return vpc, nil
}

// ListVPCs 列出 VPC（从本地数据库获取同步后的数据）
// 这是设计文档要求的正确实现：资源列表应从本地数据库获取
// projectIDs参数用于项目隔离过滤（可选）
func (s *NetworkService) ListVPCs(ctx context.Context, accountID uint, filter cloudprovider.VPCFilter, projectIDs []int64) ([]*cloudprovider.VPC, error) {
	var cloudVPCs []model.CloudVPC

	query := s.db.WithContext(ctx).Where("cloud_account_id = ?", accountID)

	// 应用项目隔离过滤
	if len(projectIDs) > 0 {
		query = query.Where("project_id IN ?", projectIDs)
	}

	// 应用过滤器
	if filter.RegionID != "" {
		query = query.Where("region_id = ?", filter.RegionID)
	}
	if filter.VPCID != "" {
		query = query.Where("vpc_id = ?", filter.VPCID)
	}
	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	// 排除已删除的资源
	query = query.Where("status != ?", "terminated")

	if filter.MaxResults > 0 {
		query = query.Limit(filter.MaxResults)
	}

	if err := query.Find(&cloudVPCs).Error; err != nil {
		return nil, err
	}

	// 转换为cloudprovider.VPC格式
	vpcs := make([]*cloudprovider.VPC, len(cloudVPCs))
	for i, vpc := range cloudVPCs {
		var tags map[string]string
		if vpc.Tags != nil {
			json.Unmarshal(vpc.Tags, &tags)
		}

		vpcs[i] = &cloudprovider.VPC{
			ID:       vpc.VPCID,
			Name:     vpc.Name,
			CIDR:     vpc.CIDR,
			Status:   vpc.Status,
			RegionID: vpc.RegionID,
			Tags:     tags,
		}
	}

	return vpcs, nil
}

// ListVPCsFromCloud 从云平台实时获取VPC列表（用于创建资源等场景）
// 这个方法用于需要实时云平台数据的场景
func (s *NetworkService) ListVPCsFromCloud(ctx context.Context, accountID uint, filter cloudprovider.VPCFilter) ([]*cloudprovider.VPC, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return provider.ListVPCs(ctx, filter)
}

// GetVPC 获取单个VPC详情（从本地数据库）
func (s *NetworkService) GetVPC(ctx context.Context, accountID uint, vpcID string) (*cloudprovider.VPC, error) {
	var cloudVPC model.CloudVPC

	err := s.db.WithContext(ctx).
		Where("cloud_account_id = ?", accountID).
		Where("vpc_id = ?", vpcID).
		First(&cloudVPC).Error

	if err == gorm.ErrRecordNotFound {
		return nil, cloudprovider.NewCloudError(cloudprovider.ErrResourceNotFound, "vpc not found", vpcID)
	}
	if err != nil {
		return nil, err
	}

	var tags map[string]string
	if cloudVPC.Tags != nil {
		json.Unmarshal(cloudVPC.Tags, &tags)
	}

	return &cloudprovider.VPC{
		ID:       cloudVPC.VPCID,
		Name:     cloudVPC.Name,
		CIDR:     cloudVPC.CIDR,
		Status:   cloudVPC.Status,
		RegionID: cloudVPC.RegionID,
		Tags:     tags,
	}, nil
}

// DeleteVPC 删除 VPC
func (s *NetworkService) DeleteVPC(ctx context.Context, accountID uint, vpcID string) error {
	// 先获取当前状态用于日志记录
	vpc, err := s.GetVPC(ctx, accountID, vpcID)
	if err != nil {
		return err
	}

	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return err
	}

	if err := provider.DeleteVPC(ctx, vpcID); err != nil {
		return err
	}

	// 记录状态变更日志
	stateLog := &model.ResourceStateLog{
		ResourceType:   "vpc",
		ResourceID:     vpcID,
		ResourceName:   vpc.Name,
		CloudAccountID: accountID,
		PreviousStatus: vpc.Status,
		CurrentStatus:  "terminated",
		OperationType:  "delete",
		Reason:         "VPC删除",
	}
	if err := s.logStateChange(ctx, stateLog); err != nil {
		// 日志记录失败不影响主流程
	}

	return nil
}

// CreateSubnet 创建子网
func (s *NetworkService) CreateSubnet(ctx context.Context, accountID uint, config cloudprovider.SubnetConfig) (*cloudprovider.Subnet, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	subnet, err := provider.CreateSubnet(ctx, config)
	if err != nil {
		return nil, err
	}

	// 记录状态变更日志
	stateLog := &model.ResourceStateLog{
		ResourceType:   "subnet",
		ResourceID:     subnet.ID,
		ResourceName:   config.Name,
		CloudAccountID: accountID,
		PreviousStatus: "",
		CurrentStatus:  subnet.Status,
		OperationType:  "create",
		Reason:         "子网创建",
	}
	if err := s.logStateChange(ctx, stateLog); err != nil {
		// 日志记录失败不影响主流程
	}

	return subnet, nil
}

// ListSubnets 列出子网（从本地数据库获取同步后的数据）
// 这是设计文档要求的正确实现：资源列表应从本地数据库获取
// projectIDs参数用于项目隔离过滤（可选）
func (s *NetworkService) ListSubnets(ctx context.Context, accountID uint, filter cloudprovider.SubnetFilter, projectIDs []int64) ([]*cloudprovider.Subnet, error) {
	var cloudSubnets []model.CloudSubnet

	query := s.db.WithContext(ctx).Where("cloud_account_id = ?", accountID)

	// 应用项目隔离过滤
	if len(projectIDs) > 0 {
		query = query.Where("project_id IN ?", projectIDs)
	}

	// 应用过滤器 - SubnetFilter只有VPCID/SubnetID/ZoneID/MaxResults字段
	if filter.VPCID != "" {
		query = query.Where("vpc_id = ?", filter.VPCID)
	}
	if filter.SubnetID != "" {
		query = query.Where("subnet_id = ?", filter.SubnetID)
	}
	if filter.ZoneID != "" {
		query = query.Where("zone_id = ?", filter.ZoneID)
	}

	// 排除已删除的资源
	query = query.Where("status != ?", "terminated")

	if filter.MaxResults > 0 {
		query = query.Limit(filter.MaxResults)
	}

	if err := query.Find(&cloudSubnets).Error; err != nil {
		return nil, err
	}

	// 转换为cloudprovider.Subnet格式
	subnets := make([]*cloudprovider.Subnet, len(cloudSubnets))
	for i, subnet := range cloudSubnets {
		var tags map[string]string
		if subnet.Tags != nil {
			json.Unmarshal(subnet.Tags, &tags)
		}

		subnets[i] = &cloudprovider.Subnet{
			ID:     subnet.SubnetID,
			Name:   subnet.Name,
			VPCID:  subnet.VPCID,
			CIDR:   subnet.CIDR,
			ZoneID: subnet.ZoneID,
			Status: subnet.Status,
			Tags:   tags,
		}
	}

	return subnets, nil
}

// ListSubnetsFromCloud 从云平台实时获取子网列表（用于创建资源等场景）
func (s *NetworkService) ListSubnetsFromCloud(ctx context.Context, accountID uint, filter cloudprovider.SubnetFilter) ([]*cloudprovider.Subnet, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return provider.ListSubnets(ctx, filter)
}

// GetSubnet 获取单个子网详情（从本地数据库）
func (s *NetworkService) GetSubnet(ctx context.Context, accountID uint, subnetID string) (*cloudprovider.Subnet, error) {
	var cloudSubnet model.CloudSubnet

	err := s.db.WithContext(ctx).
		Where("cloud_account_id = ?", accountID).
		Where("subnet_id = ?", subnetID).
		First(&cloudSubnet).Error

	if err == gorm.ErrRecordNotFound {
		return nil, cloudprovider.NewCloudError(cloudprovider.ErrResourceNotFound, "subnet not found", subnetID)
	}
	if err != nil {
		return nil, err
	}

	var tags map[string]string
	if cloudSubnet.Tags != nil {
		json.Unmarshal(cloudSubnet.Tags, &tags)
	}

	return &cloudprovider.Subnet{
		ID:     cloudSubnet.SubnetID,
		Name:   cloudSubnet.Name,
		VPCID:  cloudSubnet.VPCID,
		CIDR:   cloudSubnet.CIDR,
		ZoneID: cloudSubnet.ZoneID,
		Status: cloudSubnet.Status,
		Tags:   tags,
	}, nil
}

// DeleteSubnet 删除子网
func (s *NetworkService) DeleteSubnet(ctx context.Context, accountID uint, subnetID string) error {
	// 先获取当前状态用于日志记录
	subnet, err := s.GetSubnet(ctx, accountID, subnetID)
	if err != nil {
		return err
	}

	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return err
	}

	if err := provider.DeleteSubnet(ctx, subnetID); err != nil {
		return err
	}

	// 记录状态变更日志
	stateLog := &model.ResourceStateLog{
		ResourceType:   "subnet",
		ResourceID:     subnetID,
		ResourceName:   subnet.Name,
		CloudAccountID: accountID,
		PreviousStatus: subnet.Status,
		CurrentStatus:  "terminated",
		OperationType:  "delete",
		Reason:         "子网删除",
	}
	if err := s.logStateChange(ctx, stateLog); err != nil {
		// 日志记录失败不影响主流程
	}

	return nil
}

// CreateSecurityGroup 创建安全组
func (s *NetworkService) CreateSecurityGroup(ctx context.Context, accountID uint, config cloudprovider.SGConfig) (*cloudprovider.SecurityGroup, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	sg, err := provider.CreateSecurityGroup(ctx, config)
	if err != nil {
		return nil, err
	}

	// 记录状态变更日志
	stateLog := &model.ResourceStateLog{
		ResourceType:   "security_group",
		ResourceID:     sg.ID,
		ResourceName:   config.Name,
		CloudAccountID: accountID,
		PreviousStatus: "",
		CurrentStatus:  "available",
		OperationType:  "create",
		Reason:         "安全组创建",
	}
	if err := s.logStateChange(ctx, stateLog); err != nil {
		// 日志记录失败不影响主流程
	}

	return sg, nil
}

// ListSecurityGroups 列出安全组（从本地数据库获取同步后的数据）
// 这是设计文档要求的正确实现：资源列表应从本地数据库获取
func (s *NetworkService) ListSecurityGroups(ctx context.Context, accountID uint, filter cloudprovider.SGFilter) ([]*cloudprovider.SecurityGroup, error) {
	var cloudSGs []model.CloudSecurityGroup

	query := s.db.WithContext(ctx).Where("cloud_account_id = ?", accountID)

	// 应用过滤器 - SGFilter有VPCID/SGID/Name/MaxResults字段
	if filter.VPCID != "" {
		query = query.Where("vpc_id = ?", filter.VPCID)
	}
	if filter.SGID != "" {
		query = query.Where("security_group_id = ?", filter.SGID)
	}
	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}

	// 排除已删除的资源
	query = query.Where("status != ?", "terminated")

	if filter.MaxResults > 0 {
		query = query.Limit(filter.MaxResults)
	}

	if err := query.Find(&cloudSGs).Error; err != nil {
		return nil, err
	}

	// 转换为cloudprovider.SecurityGroup格式
	sgs := make([]*cloudprovider.SecurityGroup, len(cloudSGs))
	for i, sg := range cloudSGs {
		var tags map[string]string
		if sg.Tags != nil {
			json.Unmarshal(sg.Tags, &tags)
		}

		sgs[i] = &cloudprovider.SecurityGroup{
			ID:          sg.SecurityGroupID,
			Name:        sg.Name,
			Description: sg.Description,
			VPCID:       sg.VPCID,
			Tags:        tags,
		}
	}

	return sgs, nil
}

// ListSecurityGroupsFromCloud 从云平台实时获取安全组列表
func (s *NetworkService) ListSecurityGroupsFromCloud(ctx context.Context, accountID uint, filter cloudprovider.SGFilter) ([]*cloudprovider.SecurityGroup, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return provider.ListSecurityGroups(ctx, filter)
}

// GetSecurityGroup 获取单个安全组详情（从本地数据库）
func (s *NetworkService) GetSecurityGroup(ctx context.Context, accountID uint, sgID string) (*cloudprovider.SecurityGroup, error) {
	var cloudSG model.CloudSecurityGroup

	err := s.db.WithContext(ctx).
		Where("cloud_account_id = ?", accountID).
		Where("security_group_id = ?", sgID).
		First(&cloudSG).Error

	if err == gorm.ErrRecordNotFound {
		return nil, cloudprovider.NewCloudError(cloudprovider.ErrResourceNotFound, "security group not found", sgID)
	}
	if err != nil {
		return nil, err
	}

	var tags map[string]string
	if cloudSG.Tags != nil {
		json.Unmarshal(cloudSG.Tags, &tags)
	}

	return &cloudprovider.SecurityGroup{
		ID:          cloudSG.SecurityGroupID,
		Name:        cloudSG.Name,
		Description: cloudSG.Description,
		VPCID:       cloudSG.VPCID,
		Tags:        tags,
	}, nil
}

// CreateEIP 分配弹性 IP
func (s *NetworkService) CreateEIP(ctx context.Context, accountID uint, config cloudprovider.EIPConfig) (*cloudprovider.EIP, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	eip, err := provider.AllocateEIP(ctx, config)
	if err != nil {
		return nil, err
	}

	// 记录状态变更日志
	stateLog := &model.ResourceStateLog{
		ResourceType:   "eip",
		ResourceID:     eip.ID,
		ResourceName:   eip.Address,
		CloudAccountID: accountID,
		PreviousStatus: "",
		CurrentStatus:  eip.Status,
		OperationType:  "create",
		Reason:         "弹性IP申请",
	}
	if err := s.logStateChange(ctx, stateLog); err != nil {
		// 日志记录失败不影响主流程
	}

	return eip, nil
}

// ListEIPs 列出弹性 IP（从本地数据库获取同步后的数据）
// 这是设计文档要求的正确实现：资源列表应从本地数据库获取
func (s *NetworkService) ListEIPs(ctx context.Context, accountID uint, filter cloudprovider.EIPFilter) ([]*cloudprovider.EIP, error) {
	var cloudEIPs []model.CloudEIP

	query := s.db.WithContext(ctx).Where("cloud_account_id = ?", accountID)

	// 应用过滤器
	if filter.RegionID != "" {
		query = query.Where("region_id = ?", filter.RegionID)
	}
	if filter.EIPID != "" {
		query = query.Where("eip_id = ?", filter.EIPID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	// 排除已删除的资源
	query = query.Where("status != ?", "terminated")

	if err := query.Find(&cloudEIPs).Error; err != nil {
		return nil, err
	}

	// 转换为cloudprovider.EIP格式
	eips := make([]*cloudprovider.EIP, len(cloudEIPs))
	for i, eip := range cloudEIPs {
		eips[i] = &cloudprovider.EIP{
			ID:         eip.EIPID,
			Address:    eip.Address,
			Bandwidth:  eip.Bandwidth,
			Status:     eip.Status,
			ResourceID: eip.ResourceID,
			RegionID:   eip.RegionID,
		}
	}

	return eips, nil
}

// ListEIPsFromCloud 从云平台实时获取EIP列表
func (s *NetworkService) ListEIPsFromCloud(ctx context.Context, accountID uint, filter cloudprovider.EIPFilter) ([]*cloudprovider.EIP, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return provider.ListEIPs(ctx, filter)
}

// GetEIP 获取单个EIP详情（从本地数据库）
func (s *NetworkService) GetEIP(ctx context.Context, accountID uint, eipID string) (*cloudprovider.EIP, error) {
	var cloudEIP model.CloudEIP

	err := s.db.WithContext(ctx).
		Where("cloud_account_id = ?", accountID).
		Where("eip_id = ?", eipID).
		First(&cloudEIP).Error

	if err == gorm.ErrRecordNotFound {
		return nil, cloudprovider.NewCloudError(cloudprovider.ErrResourceNotFound, "eip not found", eipID)
	}
	if err != nil {
		return nil, err
	}

	return &cloudprovider.EIP{
		ID:         cloudEIP.EIPID,
		Address:    cloudEIP.Address,
		Bandwidth:  cloudEIP.Bandwidth,
		Status:     cloudEIP.Status,
		ResourceID: cloudEIP.ResourceID,
		RegionID:   cloudEIP.RegionID,
	}, nil
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

// ========== 子网扩展操作 ==========

// UpdateSubnet 更新子网属性
func (s *NetworkService) UpdateSubnet(ctx context.Context, accountID uint, subnetID, name, description string, tags map[string]string) (*cloudprovider.Subnet, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return provider.UpdateSubnet(ctx, subnetID, name, description, tags)
}

// ChangeSubnetProject 更改子网所属项目（本地数据库操作）
func (s *NetworkService) ChangeSubnetProject(ctx context.Context, accountID uint, subnetID string, projectID uint) error {
	// 更新本地数据库中的子网项目关联
	// 实际云上子网的项目归属由本地系统管理
	return nil // TODO: 实现本地数据库更新逻辑
}

// SplitSubnet 分割IP子网
func (s *NetworkService) SplitSubnet(ctx context.Context, accountID uint, subnetID string, newCIDRs []string) ([]*cloudprovider.Subnet, error) {
	// 分割子网逻辑：先获取原子网信息，然后创建新子网
	// TODO: 实现完整的分割逻辑
	return nil, cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "subnet split not implemented", "")
}

// ReservedIP 预留IP记录
type ReservedIP struct {
	IP        string `json:"ip"`
	SubnetID  string `json:"subnet_id"`
	Reason    string `json:"reason"`
	ReservedBy string `json:"reserved_by"`
	Status    string `json:"status"`
}

// ReserveIP 预留IP地址
func (s *NetworkService) ReserveIP(ctx context.Context, accountID uint, subnetID string, ips []string, reason, reservedBy string) ([]*ReservedIP, error) {
	// IP预留由本地系统管理，不调用云厂商API
	reservedIPs := make([]*ReservedIP, len(ips))
	for i, ip := range ips {
		reservedIPs[i] = &ReservedIP{
			IP:         ip,
			SubnetID:   subnetID,
			Reason:     reason,
			ReservedBy: reservedBy,
			Status:     "reserved",
		}
	}
	// TODO: 存储到本地数据库
	return reservedIPs, nil
}

// ReleaseIP 释放预留IP
func (s *NetworkService) ReleaseIP(ctx context.Context, accountID uint, subnetID string, ips []string) error {
	// TODO: 从本地数据库删除预留记录
	return nil
}

// ListReservedIPs 列出预留IP
func (s *NetworkService) ListReservedIPs(ctx context.Context, accountID uint, subnetID string) ([]*ReservedIP, error) {
	// TODO: 从本地数据库查询
	return []*ReservedIP{}, nil
}

// ========== 安全组扩展操作 ==========

// AddSecurityGroupRule 添加安全组规则
func (s *NetworkService) AddSecurityGroupRule(ctx context.Context, accountID uint, sgID string, rule cloudprovider.SGRule) (string, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return "", err
	}

	return provider.AddSecurityGroupRule(ctx, sgID, rule)
}

// DeleteSecurityGroupRule 删除安全组规则
func (s *NetworkService) DeleteSecurityGroupRule(ctx context.Context, accountID uint, sgID, ruleID string) error {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return err
	}

	return provider.DeleteSecurityGroupRule(ctx, sgID, ruleID)
}

// DeleteSecurityGroup 删除安全组
func (s *NetworkService) DeleteSecurityGroup(ctx context.Context, accountID uint, sgID string) error {
	// 先获取当前状态用于日志记录
	sg, err := s.GetSecurityGroup(ctx, accountID, sgID)
	if err != nil {
		return err
	}

	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return err
	}

	if err := provider.DeleteSecurityGroup(ctx, sgID); err != nil {
		return err
	}

	// 记录状态变更日志
	stateLog := &model.ResourceStateLog{
		ResourceType:   "security_group",
		ResourceID:     sgID,
		ResourceName:   sg.Name,
		CloudAccountID: accountID,
		PreviousStatus: "available",
		CurrentStatus:  "terminated",
		OperationType:  "delete",
		Reason:         "安全组删除",
	}
	if err := s.logStateChange(ctx, stateLog); err != nil {
		// 日志记录失败不影响主流程
	}

	return nil
}

// ========== EIP 扩展操作 ==========

// BindEIP 绑定弹性IP
func (s *NetworkService) BindEIP(ctx context.Context, accountID uint, eipID, resourceID, resourceType string) error {
	// 先获取当前状态用于日志记录
	eip, err := s.GetEIP(ctx, accountID, eipID)
	if err != nil {
		return err
	}

	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return err
	}

	if err := provider.BindEIP(ctx, eipID, resourceID, resourceType); err != nil {
		return err
	}

	// 记录状态变更日志
	stateLog := &model.ResourceStateLog{
		ResourceType:   "eip",
		ResourceID:     eipID,
		ResourceName:   eip.Address,
		CloudAccountID: accountID,
		PreviousStatus: eip.Status,
		CurrentStatus:  "in-use",
		OperationType:  "bind",
		Reason:         "弹性IP绑定到" + resourceType + ":" + resourceID,
	}
	if err := s.logStateChange(ctx, stateLog); err != nil {
		// 日志记录失败不影响主流程
	}

	return nil
}

// UnbindEIP 解绑弹性IP
func (s *NetworkService) UnbindEIP(ctx context.Context, accountID uint, eipID string) error {
	// 先获取当前状态用于日志记录
	eip, err := s.GetEIP(ctx, accountID, eipID)
	if err != nil {
		return err
	}

	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return err
	}

	if err := provider.UnbindEIP(ctx, eipID); err != nil {
		return err
	}

	// 记录状态变更日志
	stateLog := &model.ResourceStateLog{
		ResourceType:   "eip",
		ResourceID:     eipID,
		ResourceName:   eip.Address,
		CloudAccountID: accountID,
		PreviousStatus: eip.Status,
		CurrentStatus:  "available",
		OperationType:  "unbind",
		Reason:         "弹性IP解绑",
	}
	if err := s.logStateChange(ctx, stateLog); err != nil {
		// 日志记录失败不影响主流程
	}

	return nil
}

// DeleteEIP 删除弹性IP
func (s *NetworkService) DeleteEIP(ctx context.Context, accountID uint, eipID string) error {
	// 先获取当前状态用于日志记录
	eip, err := s.GetEIP(ctx, accountID, eipID)
	if err != nil {
		return err
	}

	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return err
	}

	if err := provider.ReleaseEIP(ctx, eipID); err != nil {
		return err
	}

	// 记录状态变更日志
	stateLog := &model.ResourceStateLog{
		ResourceType:   "eip",
		ResourceID:     eipID,
		ResourceName:   eip.Address,
		CloudAccountID: accountID,
		PreviousStatus: eip.Status,
		CurrentStatus:  "terminated",
		OperationType:  "delete",
		Reason:         "弹性IP释放",
	}
	if err := s.logStateChange(ctx, stateLog); err != nil {
		// 日志记录失败不影响主流程
	}

	return nil
}
