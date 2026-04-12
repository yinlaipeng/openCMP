package service

import (
	"context"
	"encoding/json"
	"strconv"

	"gorm.io/gorm"

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

// CreateVM 创建虚拟机
func (s *ComputeService) CreateVM(ctx context.Context, accountID uint, config cloudprovider.VMCreateConfig) (*cloudprovider.VirtualMachine, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return provider.CreateVM(ctx, config)
}

// ListVMs 列出虚拟机
func (s *ComputeService) ListVMs(ctx context.Context, accountID uint, filter cloudprovider.VMListFilter) ([]*cloudprovider.VirtualMachine, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return provider.ListVMs(ctx, filter)
}

// GetVM 获取虚拟机
func (s *ComputeService) GetVM(ctx context.Context, accountID uint, vmID string) (*cloudprovider.VirtualMachine, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	// 通过 ListVMs 过滤获取单个 VM
	vms, err := provider.ListVMs(ctx, cloudprovider.VMListFilter{})
	if err != nil {
		return nil, err
	}

	for _, vm := range vms {
		if vm.ID == vmID {
			return vm, nil
		}
	}

	return nil, cloudprovider.NewCloudError(cloudprovider.ErrResourceNotFound, "vm not found", vmID)
}

// DeleteVM 删除虚拟机
func (s *ComputeService) DeleteVM(ctx context.Context, accountID uint, vmID string) error {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return err
	}

	return provider.DeleteVM(ctx, vmID)
}

// StartVM 启动虚拟机
func (s *ComputeService) StartVM(ctx context.Context, accountID uint, vmID string) error {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return err
	}

	return provider.StartVM(ctx, vmID)
}

// StopVM 停止虚拟机
func (s *ComputeService) StopVM(ctx context.Context, accountID uint, vmID string) error {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return err
	}

	return provider.StopVM(ctx, vmID)
}

// RebootVM 重启虚拟机
func (s *ComputeService) RebootVM(ctx context.Context, accountID uint, vmID string) error {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return err
	}

	return provider.RebootVM(ctx, vmID)
}

// VMAction 虚拟机操作
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

// ListImages 列出镜像
func (s *ComputeService) ListImages(ctx context.Context, accountID uint, filter cloudprovider.ImageFilter) ([]*cloudprovider.Image, error) {
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

// GetVMDetails 获取虚拟机详细信息
func (s *ComputeService) GetVMDetails(ctx context.Context, accountID uint, vmID string) (*cloudprovider.VirtualMachine, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	// Check if provider supports GetVM directly
	if vmProvider, ok := provider.(interface {
		GetVM(ctx context.Context, vmID string) (*cloudprovider.VirtualMachine, error)
	}); ok {
		return vmProvider.GetVM(ctx, vmID)
	}

	// Fallback: Get all VMs and find the one we want
	vms, err := provider.ListVMs(ctx, cloudprovider.VMListFilter{})
	if err != nil {
		return nil, err
	}

	for _, vm := range vms {
		if vm.ID == vmID {
			return vm, nil
		}
	}

	return nil, cloudprovider.NewCloudError(cloudprovider.ErrResourceNotFound, "vm not found", vmID)
}

// GetVMSecurityGroups 获取虚拟机关联的安全组
func (s *ComputeService) GetVMSecurityGroups(ctx context.Context, accountID uint, vmID string) ([]*cloudprovider.SecurityGroup, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	vm, err := s.GetVMDetails(ctx, accountID, vmID)
	if err != nil {
		return nil, err
	}

	// Check if provider supports security groups
	securityGroupProvider, ok := provider.(cloudprovider.ISecurityGroup)
	if !ok {
		// Return empty list if not supported
		return []*cloudprovider.SecurityGroup{}, nil
	}

	// Get all security groups and filter by VM's security group IDs
	allSecurityGroups, err := securityGroupProvider.ListSecurityGroups(ctx, cloudprovider.SGFilter{}) // Use correct filter type
	if err != nil {
		// If listing fails, return empty list instead of error
		return []*cloudprovider.SecurityGroup{}, nil
	}

	var vmSecurityGroups []*cloudprovider.SecurityGroup
	for _, sg := range allSecurityGroups {
		for _, sgID := range vm.SecurityGroups {
			if sg.ID == sgID {
				vmSecurityGroups = append(vmSecurityGroups, sg)
			}
		}
	}

	return vmSecurityGroups, nil
}

// GetVMNetworkInfo 获取虚拟机网络信息
func (s *ComputeService) GetVMNetworkInfo(ctx context.Context, accountID uint, vmID string) (map[string]interface{}, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	vm, err := s.GetVMDetails(ctx, accountID, vmID)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"private_ip": vm.PrivateIP,
		"public_ip":  vm.PublicIP,
		"vm_id":      vmID,
	}

	// Try to get VPC and subnet details if provider supports it
	vpcProvider, okVPC := provider.(cloudprovider.IVPC)
	if okVPC {
		vpc, err := vpcProvider.GetVPC(ctx, vm.VPCID)
		if err == nil {
			result["vpc"] = vpc
		}
	}

	subnetProvider, okSubnet := provider.(cloudprovider.ISubnet)
	if okSubnet {
		subnet, err := subnetProvider.GetSubnet(ctx, vm.SubnetID)
		if err == nil {
			result["subnet"] = subnet
		}
	}

	return result, nil
}

// GetVMDisks 获取虚拟机关联的磁盘
func (s *ComputeService) GetVMDisks(ctx context.Context, accountID uint, vmID string) ([]*cloudprovider.Disk, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	// Check if provider supports disk operations
	diskProvider, ok := provider.(cloudprovider.IDisk)
	if !ok {
		// Return empty list if not supported
		return []*cloudprovider.Disk{}, nil
	}

	// List all disks and filter by VM association
	allDisks, err := diskProvider.ListDisks(ctx, cloudprovider.DiskFilter{VMID: vmID}) // Use correct filter type
	if err != nil {
		// If listing fails, return empty list instead of error
		return []*cloudprovider.Disk{}, nil
	}

	var vmDisks []*cloudprovider.Disk
	for _, disk := range allDisks {
		if disk.VMID == vmID {
			vmDisks = append(vmDisks, disk)
		}
	}

	return vmDisks, nil
}

// GetVMSnapshots 获取虚拟机相关的快照
func (s *ComputeService) GetVMSnapshots(ctx context.Context, accountID uint, vmID string) ([]*cloudprovider.Snapshot, error) {
	provider, err := s.getProvider(ctx, accountID)
	if err != nil {
		return nil, err
	}

	vm, err := s.GetVMDetails(ctx, accountID, vmID)
	if err != nil {
		return nil, err
	}

	// Check if provider supports storage operations (which includes snapshots)
	storageProvider, ok := provider.(cloudprovider.IStorage)
	if !ok {
		// If storage not supported, return empty list
		return []*cloudprovider.Snapshot{}, nil
	}

	// If we have storage provider, check if it also implements snapshot functionality
	var allSnapshots []*cloudprovider.Snapshot
	for _, diskID := range vm.DiskIDs {
		// Check if storage provider also has snapshot capability
		if diskSnapshotProvider, ok := storageProvider.(interface {
			ListDiskSnapshots(ctx context.Context, diskID string) ([]*cloudprovider.Snapshot, error)
		}); ok {
			diskSnapshots, err := diskSnapshotProvider.ListDiskSnapshots(ctx, diskID)
			if err == nil {
				allSnapshots = append(allSnapshots, diskSnapshots...)
			}
		}
	}

	return allSnapshots, nil
}

// GetVMOperationLogs 获取虚拟机操作日志
func (s *ComputeService) GetVMOperationLogs(ctx context.Context, accountID uint, vmID string) ([]map[string]interface{}, error) {
	// For now, return placeholder logs
	// In a real implementation, this would connect to a logging system
	// or retrieve historical operation records
	logs := []map[string]interface{}{
		{
			"id":         "log-001",
			"operation":  "created",
			"timestamp":  "2024-01-01T10:00:00Z",
			"status":     "success",
			"operator":   "system",
			"details":    "Virtual machine created successfully",
		},
		{
			"id":         "log-002",
			"operation":  "started",
			"timestamp":  "2024-01-01T10:05:00Z",
			"status":     "success",
			"operator":   "admin",
			"details":    "VM started by admin user",
		},
	}

	return logs, nil
}

// GetVNCInfo 获取虚拟机VNC连接信息
func (s *ComputeService) GetVNCInfo(ctx context.Context, accountID uint, vmID string) (map[string]interface{}, error) {
	// In a real implementation, this would return actual VNC connection details
	// based on the cloud provider's console access mechanism
	vncInfo := map[string]interface{}{
		"vm_id":        vmID,
		"console_url":  "", // This would be populated with actual console URL
		"connection_type": "web_console", // or vnc, spice, etc.
		"supports_copy_paste": true,
		"console_username": "", // Usually determined by cloud provider
		"access_token":   "", // Temporary access token if needed
		"expires_at":     "", // When the token expires
	}

	return vncInfo, nil
}

// ResetPassword 重置虚拟机密码
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

// UpdateVMConfig 更新虚拟机配置
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
