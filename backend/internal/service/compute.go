package service

import (
	"context"
	"encoding/json"

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
		ID:           string(account.ID),
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
