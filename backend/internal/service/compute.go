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

// ComputeService 计算资源服务
type ComputeService struct {
	db *gorm.DB
}

// NewComputeService 创建计算资源服务
func NewComputeService(db *gorm.DB) *ComputeService {
	return &ComputeService{db: db}
}

// getProvider 获取云提供商
func (s *ComputeService) getProvider(ctx context.Context, accountID uint) (cloudprovider.ICloudProvider, error) {
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
func (s *ComputeService) logStateChange(ctx context.Context, log *model.ResourceStateLog) error {
	log.OccurredAt = time.Now()
	log.CreatedAt = time.Now()
	return s.db.WithContext(ctx).Create(log).Error
}

// CreateVM 创建虚拟机（调用云平台API）
func (s *ComputeService) CreateVM(ctx context.Context, accountID uint, config cloudprovider.VMCreateConfig) (*cloudprovider.VirtualMachine, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	vm, err := provider.CreateVM(ctx, config)
	if err != nil {
		return nil, err
	}

	// 记录状态变更日志
	stateLog := &model.ResourceStateLog{
		ResourceType:   "vm",
		ResourceID:     vm.ID,
		ResourceName:   config.Name,
		CloudAccountID: accountID,
		PreviousStatus: "",
		CurrentStatus:  string(vm.Status),
		OperationType:  "create",
		Reason:         "虚拟机创建",
	}
	if err := s.logStateChange(ctx, stateLog); err != nil {
		// 日志记录失败不影响主流程
	}

	return vm, nil
}

// ListVMs 列出虚拟机（从本地数据库获取同步后的数据）
// 这是设计文档要求的正确实现：资源列表应从本地数据库获取
// projectIDs参数用于项目隔离过滤（可选）
func (s *ComputeService) ListVMs(ctx context.Context, accountID uint, filter cloudprovider.VMListFilter, projectIDs []int64) ([]*cloudprovider.VirtualMachine, error) {
	var cloudVMs []model.CloudVM

	query := s.db.WithContext(ctx).Where("cloud_account_id = ?", accountID)

	// 应用项目隔离过滤
	if len(projectIDs) > 0 {
		query = query.Where("project_id IN ?", projectIDs)
	}

	// 应用过滤器
	if filter.VPCID != "" {
		query = query.Where("vpc_id = ?", filter.VPCID)
	}
	if filter.SubnetID != "" {
		query = query.Where("subnet_id = ?", filter.SubnetID)
	}
	if filter.RegionID != "" {
		query = query.Where("region_id = ?", filter.RegionID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	// 排除已终止的资源
	query = query.Where("status != ?", "terminated")

	if filter.MaxResults > 0 {
		query = query.Limit(filter.MaxResults)
	}

	if err := query.Find(&cloudVMs).Error; err != nil {
		return nil, err
	}

	// 转换为cloudprovider.VirtualMachine格式
	vms := make([]*cloudprovider.VirtualMachine, len(cloudVMs))
	for i, vm := range cloudVMs {
		var tags map[string]string
		if vm.Tags != nil {
			json.Unmarshal(vm.Tags, &tags)
		}

		vms[i] = &cloudprovider.VirtualMachine{
			ID:           vm.InstanceID,
			Name:         vm.Name,
			Status:       cloudprovider.VMStatus(vm.Status),
			InstanceType: vm.InstanceType,
			ImageID:      vm.ImageID,
			OSName:       vm.OSName,
			VPCID:        vm.VPCID,
			SubnetID:     vm.SubnetID,
			PrivateIP:    vm.PrivateIP,
			PublicIP:     vm.PublicIP,
			RegionID:     vm.RegionID,
			ZoneID:       vm.ZoneID,
			Tags:         tags,
		}
	}

	return vms, nil
}

// ListAllVMs 列出所有云账号的虚拟机（从本地数据库获取同步后的数据）
// 用于默认展示所有虚拟机的场景
func (s *ComputeService) ListAllVMs(ctx context.Context, filter cloudprovider.VMListFilter, projectIDs []int64) ([]*cloudprovider.VirtualMachine, error) {
	var cloudVMs []model.CloudVM

	query := s.db.WithContext(ctx)

	// 应用项目隔离过滤
	if len(projectIDs) > 0 {
		query = query.Where("project_id IN ?", projectIDs)
	}

	// 应用过滤器
	if filter.VPCID != "" {
		query = query.Where("vpc_id = ?", filter.VPCID)
	}
	if filter.SubnetID != "" {
		query = query.Where("subnet_id = ?", filter.SubnetID)
	}
	if filter.RegionID != "" {
		query = query.Where("region_id = ?", filter.RegionID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	// 排除已终止的资源
	query = query.Where("status != ?", "terminated")

	if filter.MaxResults > 0 {
		query = query.Limit(filter.MaxResults)
	}

	if err := query.Find(&cloudVMs).Error; err != nil {
		return nil, err
	}

	// 获取所有涉及的云账号ID
	accountIDs := make([]uint, 0)
	for _, vm := range cloudVMs {
		accountIDs = append(accountIDs, vm.CloudAccountID)
	}

	// 查询云账号信息
	accountMap := make(map[uint]*model.CloudAccount)
	if len(accountIDs) > 0 {
		var accounts []model.CloudAccount
		s.db.WithContext(ctx).Where("id IN ?", accountIDs).Find(&accounts)
		for i := range accounts {
			accountMap[accounts[i].ID] = &accounts[i]
		}
	}

	// 转换为cloudprovider.VirtualMachine格式
	vms := make([]*cloudprovider.VirtualMachine, len(cloudVMs))
	for i, vm := range cloudVMs {
		var tags map[string]string
		if vm.Tags != nil {
			json.Unmarshal(vm.Tags, &tags)
		}

		vms[i] = &cloudprovider.VirtualMachine{
			ID:             vm.InstanceID,
			Name:           vm.Name,
			Status:         cloudprovider.VMStatus(vm.Status),
			InstanceType:   vm.InstanceType,
			ImageID:        vm.ImageID,
			OSName:         vm.OSName,
			VPCID:          vm.VPCID,
			SubnetID:       vm.SubnetID,
			PrivateIP:      vm.PrivateIP,
			PublicIP:       vm.PublicIP,
			RegionID:       vm.RegionID,
			ZoneID:         vm.ZoneID,
			Tags:           tags,
			CloudAccountID: vm.CloudAccountID,
		}

		// 附加云账号信息
		if account, ok := accountMap[vm.CloudAccountID]; ok {
			vms[i].Platform = account.ProviderType
			vms[i].AccountName = account.Name
		}
	}

	return vms, nil
}

// ListVMsFromCloud 从云平台实时获取虚拟机列表（用于创建资源等场景）
// 这个方法用于需要实时云平台数据的场景，如资源创建前的查询
func (s *ComputeService) ListVMsFromCloud(ctx context.Context, accountID uint, filter cloudprovider.VMListFilter) ([]*cloudprovider.VirtualMachine, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return provider.ListVMs(ctx, filter)
}

// GetVM 获取虚拟机（从本地数据库）
func (s *ComputeService) GetVM(ctx context.Context, accountID uint, vmID string) (*cloudprovider.VirtualMachine, error) {
	var cloudVM model.CloudVM

	err := s.db.WithContext(ctx).
		Where("cloud_account_id = ?", accountID).
		Where("instance_id = ?", vmID).
		First(&cloudVM).Error

	if err == gorm.ErrRecordNotFound {
		return nil, cloudprovider.NewCloudError(cloudprovider.ErrResourceNotFound, "vm not found", vmID)
	}
	if err != nil {
		return nil, err
	}

	var tags map[string]string
	if cloudVM.Tags != nil {
		json.Unmarshal(cloudVM.Tags, &tags)
	}

	return &cloudprovider.VirtualMachine{
		ID:           cloudVM.InstanceID,
		Name:         cloudVM.Name,
		Status:       cloudprovider.VMStatus(cloudVM.Status),
		InstanceType: cloudVM.InstanceType,
		ImageID:      cloudVM.ImageID,
		OSName:       cloudVM.OSName,
		VPCID:        cloudVM.VPCID,
		SubnetID:     cloudVM.SubnetID,
		PrivateIP:    cloudVM.PrivateIP,
		PublicIP:     cloudVM.PublicIP,
		RegionID:     cloudVM.RegionID,
		ZoneID:       cloudVM.ZoneID,
		Tags:         tags,
	}, nil
}

// DeleteVM 删除虚拟机（调用云平台API）
func (s *ComputeService) DeleteVM(ctx context.Context, accountID uint, vmID string) error {
	// 先获取当前状态用于日志记录
	vm, err := s.GetVM(ctx, accountID, vmID)
	if err != nil {
		return err
	}

	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return err
	}

	if err := provider.DeleteVM(ctx, vmID); err != nil {
		return err
	}

	// 记录状态变更日志
	stateLog := &model.ResourceStateLog{
		ResourceType:   "vm",
		ResourceID:     vmID,
		ResourceName:   vm.Name,
		CloudAccountID: accountID,
		PreviousStatus: string(vm.Status),
		CurrentStatus:  "terminated",
		OperationType:  "delete",
		Reason:         "虚拟机删除",
	}
	if err := s.logStateChange(ctx, stateLog); err != nil {
		// 日志记录失败不影响主流程
	}

	return nil
}

// StartVM 启动虚拟机（调用云平台API）
func (s *ComputeService) StartVM(ctx context.Context, accountID uint, vmID string) error {
	// 先获取当前状态用于日志记录
	vm, err := s.GetVM(ctx, accountID, vmID)
	if err != nil {
		return err
	}

	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return err
	}

	if err := provider.StartVM(ctx, vmID); err != nil {
		return err
	}

	// 记录状态变更日志
	stateLog := &model.ResourceStateLog{
		ResourceType:   "vm",
		ResourceID:     vmID,
		ResourceName:   vm.Name,
		CloudAccountID: accountID,
		PreviousStatus: string(vm.Status),
		CurrentStatus:  "Starting",
		OperationType:  "start",
		Reason:         "虚拟机启动",
	}
	if err := s.logStateChange(ctx, stateLog); err != nil {
		// 日志记录失败不影响主流程
	}

	return nil
}

// StopVM 停止虚拟机（调用云平台API）
func (s *ComputeService) StopVM(ctx context.Context, accountID uint, vmID string) error {
	// 先获取当前状态用于日志记录
	vm, err := s.GetVM(ctx, accountID, vmID)
	if err != nil {
		return err
	}

	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return err
	}

	if err := provider.StopVM(ctx, vmID); err != nil {
		return err
	}

	// 记录状态变更日志
	stateLog := &model.ResourceStateLog{
		ResourceType:   "vm",
		ResourceID:     vmID,
		ResourceName:   vm.Name,
		CloudAccountID: accountID,
		PreviousStatus: string(vm.Status),
		CurrentStatus:  "Stopping",
		OperationType:  "stop",
		Reason:         "虚拟机停止",
	}
	if err := s.logStateChange(ctx, stateLog); err != nil {
		// 日志记录失败不影响主流程
	}

	return nil
}

// RebootVM 重启虚拟机（调用云平台API）
func (s *ComputeService) RebootVM(ctx context.Context, accountID uint, vmID string) error {
	// 先获取当前状态用于日志记录
	vm, err := s.GetVM(ctx, accountID, vmID)
	if err != nil {
		return err
	}

	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return err
	}

	if err := provider.RebootVM(ctx, vmID); err != nil {
		return err
	}

	// 记录状态变更日志
	stateLog := &model.ResourceStateLog{
		ResourceType:   "vm",
		ResourceID:     vmID,
		ResourceName:   vm.Name,
		CloudAccountID: accountID,
		PreviousStatus: string(vm.Status),
		CurrentStatus:  "Starting",
		OperationType:  "reboot",
		Reason:         "虚拟机重启",
	}
	if err := s.logStateChange(ctx, stateLog); err != nil {
		// 日志记录失败不影响主流程
	}

	return nil
}

// VMAction 虚拟机操作（调用云平台API）
func (s *ComputeService) VMAction(ctx context.Context, accountID uint, vmID string, action string) error {
	switch action {
	case "start":
		return s.StartVM(ctx, accountID, vmID)
	case "stop":
		return s.StopVM(ctx, accountID, vmID)
	case "reboot":
		return s.RebootVM(ctx, accountID, vmID)
	default:
		return cloudprovider.NewCloudError(
			cloudprovider.ErrUnsupportedOperation,
			"unsupported action",
			action,
		)
	}
}

// ListImages 列出镜像（从本地数据库获取同步后的数据）
func (s *ComputeService) ListImages(ctx context.Context, accountID uint, filter cloudprovider.ImageFilter) ([]*cloudprovider.Image, error) {
	var cloudImages []model.CloudImage

	query := s.db.WithContext(ctx).Where("cloud_account_id = ?", accountID)

	if filter.Platform != "" {
		query = query.Where("os_name LIKE ?", "%"+filter.Platform+"%")
	}

	query = query.Where("status != ?", "terminated")

	if err := query.Limit(filter.MaxResults).Find(&cloudImages).Error; err != nil {
		return nil, err
	}

	images := make([]*cloudprovider.Image, len(cloudImages))
	for i, img := range cloudImages {
		var tags map[string]string
		if img.Tags != nil {
			json.Unmarshal(img.Tags, &tags)
		}

		images[i] = &cloudprovider.Image{
			ID:           img.ImageID,
			Name:         img.Name,
			OSName:       img.OSName,
			OSVersion:    img.OSVersion,
			Architecture: img.Architecture,
			Status:       img.Status,
			Tags:         tags,
		}
	}

	return images, nil
}

// ListImagesFromCloud 从云平台实时获取镜像列表（用于创建资源前的查询）
func (s *ComputeService) ListImagesFromCloud(ctx context.Context, accountID uint, filter cloudprovider.ImageFilter) ([]*cloudprovider.Image, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return provider.ListImages(ctx, filter)
}

// GetImage 获取镜像
func (s *ComputeService) GetImage(ctx context.Context, accountID uint, imageID string) (*cloudprovider.Image, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return provider.GetImage(ctx, imageID)
}

// GetVMDetails 获取虚拟机详细信息（从本地数据库）
func (s *ComputeService) GetVMDetails(ctx context.Context, accountID uint, vmID string) (*cloudprovider.VirtualMachine, error) {
	return s.GetVM(ctx, accountID, vmID)
}

// GetVMSecurityGroups 获取虚拟机关联的安全组（从本地数据库）
func (s *ComputeService) GetVMSecurityGroups(ctx context.Context, accountID uint, vmID string) ([]*cloudprovider.SecurityGroup, error) {
	// 先获取VM以获取其VPCID用于过滤安全组
	vm, err := s.GetVM(ctx, accountID, vmID)
	if err != nil {
		return nil, err
	}

	var cloudSGs []model.CloudSecurityGroup
	query := s.db.WithContext(ctx).
		Where("cloud_account_id = ?", accountID)

	// 根据VM的VPCID过滤安全组
	if vm.VPCID != "" {
		query = query.Where("vpc_id = ?", vm.VPCID)
	}

	query.Find(&cloudSGs)

	sgs := make([]*cloudprovider.SecurityGroup, 0)
	for _, sg := range cloudSGs {
		sgs = append(sgs, &cloudprovider.SecurityGroup{
			ID:          sg.SecurityGroupID,
			Name:        sg.Name,
			Description: sg.Description,
			VPCID:       sg.VPCID,
		})
	}

	return sgs, nil
}

// GetVMNetworkInfo 获取虚拟机网络信息（从本地数据库）
func (s *ComputeService) GetVMNetworkInfo(ctx context.Context, accountID uint, vmID string) (map[string]interface{}, error) {
	vm, err := s.GetVM(ctx, accountID, vmID)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"private_ip": vm.PrivateIP,
		"public_ip":  vm.PublicIP,
		"vm_id":      vmID,
		"vpc_id":     vm.VPCID,
		"subnet_id":  vm.SubnetID,
	}

	// 获取VPC和子网详情
	var vpc model.CloudVPC
	if err := s.db.WithContext(ctx).
		Where("cloud_account_id = ?", accountID).
		Where("vpc_id = ?", vm.VPCID).
		First(&vpc).Error; err == nil {
		result["vpc"] = map[string]interface{}{
			"id":     vpc.VPCID,
			"name":   vpc.Name,
			"cidr":   vpc.CIDR,
			"status": vpc.Status,
		}
	}

	var subnet model.CloudSubnet
	if err := s.db.WithContext(ctx).
		Where("cloud_account_id = ?", accountID).
		Where("subnet_id = ?", vm.SubnetID).
		First(&subnet).Error; err == nil {
		result["subnet"] = map[string]interface{}{
			"id":     subnet.SubnetID,
			"name":   subnet.Name,
			"cidr":   subnet.CIDR,
			"zone":   subnet.ZoneID,
			"status": subnet.Status,
		}
	}

	return result, nil
}

// GetVMDisks 获取虚拟机关联的磁盘（从本地数据库）
func (s *ComputeService) GetVMDisks(ctx context.Context, accountID uint, vmID string) ([]*cloudprovider.Disk, error) {
	var cloudDisks []model.CloudDisk

	s.db.WithContext(ctx).
		Where("cloud_account_id = ?", accountID).
		Where("vm_id = ?", vmID).
		Where("status != ?", "terminated").
		Find(&cloudDisks)

	disks := make([]*cloudprovider.Disk, len(cloudDisks))
	for i, disk := range cloudDisks {
		var tags map[string]string
		if disk.Tags != nil {
			json.Unmarshal(disk.Tags, &tags)
		}

		disks[i] = &cloudprovider.Disk{
			ID:     disk.DiskID,
			Name:   disk.Name,
			Size:   disk.Size,
			Type:   disk.Type,
			Status: disk.Status,
			VMID:   disk.VMID,
			ZoneID: disk.ZoneID,
			Tags:   tags,
		}
	}

	return disks, nil
}

// GetVMSnapshots 获取虚拟机相关的快照（从本地数据库）
func (s *ComputeService) GetVMSnapshots(ctx context.Context, accountID uint, vmID string) ([]*cloudprovider.Snapshot, error) {
	// 先获取VM关联的磁盘
	disks, err := s.GetVMDisks(ctx, accountID, vmID)
	if err != nil {
		return nil, err
	}

	diskIDs := make([]string, len(disks))
	for i, d := range disks {
		diskIDs[i] = d.ID
	}

	if len(diskIDs) == 0 {
		return []*cloudprovider.Snapshot{}, nil
	}

	var cloudSnapshots []model.CloudSnapshot
	s.db.WithContext(ctx).
		Where("cloud_account_id = ?", accountID).
		Where("disk_id IN ?", diskIDs).
		Where("status != ?", "terminated").
		Find(&cloudSnapshots)

	snapshots := make([]*cloudprovider.Snapshot, len(cloudSnapshots))
	for i, snap := range cloudSnapshots {
		snapshots[i] = &cloudprovider.Snapshot{
			ID:     snap.SnapshotID,
			Name:   snap.Name,
			DiskID: snap.DiskID,
			Size:   snap.Size,
			Status: snap.Status,
		}
	}

	return snapshots, nil
}

// GetVMOperationLogs 获取虚拟机操作日志（从本地数据库）
func (s *ComputeService) GetVMOperationLogs(ctx context.Context, accountID uint, vmID string) ([]map[string]interface{}, error) {
	var logs []model.OperationLog

	s.db.WithContext(ctx).
		Where("resource_type = ?", "vm").
		Where("resource_name = ?", vmID).
		Order("operation_time DESC").
		Limit(20).
		Find(&logs)

	result := make([]map[string]interface{}, len(logs))
	for i, log := range logs {
		result[i] = map[string]interface{}{
			"id":         log.ID,
			"operation":  log.OperationType,
			"timestamp":  log.OperationTime,
			"status":     log.Result,
			"operator":   log.Operator,
			"details":    log.ResourceName,
		}
	}

	return result, nil
}

// GetVNCInfo 获取虚拟机VNC连接信息（调用云平台API）
func (s *ComputeService) GetVNCInfo(ctx context.Context, accountID uint, vmID string) (map[string]interface{}, error) {
	_, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	// VNC需要实时云平台数据 - TODO: 实现真实的VNC连接获取逻辑
	return map[string]interface{}{
		"vm_id":               vmID,
		"console_url":         "",
		"connection_type":     "web_console",
		"supports_copy_paste": true,
	}, nil
}

// ResetPassword 重置虚拟机密码（调用云平台API）
func (s *ComputeService) ResetPassword(ctx context.Context, accountID uint, vmID, username, newPassword string) error {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return err
	}

	// Check if the provider supports password reset
	if vmProvider, ok := provider.(interface {
		ResetVMPassword(ctx context.Context, vmID, username, newPassword string) error
	}); ok {
		return vmProvider.ResetVMPassword(ctx, vmID, username, newPassword)
	}

	// If not supported, return appropriate error
	return cloudprovider.NewCloudError(
		cloudprovider.ErrUnsupportedOperation,
		"password reset not supported by this provider",
		vmID,
	)
}

// UpdateVMConfig 更新虚拟机配置（调用云平台API）
func (s *ComputeService) UpdateVMConfig(ctx context.Context, accountID uint, vmID, instanceType, name string) error {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return err
	}

	// Check if the provider supports VM config update
	if vmProvider, ok := provider.(interface {
		UpdateVMConfig(ctx context.Context, vmID, instanceType, name string) error
	}); ok {
		return vmProvider.UpdateVMConfig(ctx, vmID, instanceType, name)
	}

	// If not supported, return appropriate error
	return cloudprovider.NewCloudError(
		cloudprovider.ErrUnsupportedOperation,
		"config update not supported by this provider",
		vmID,
	)
}
