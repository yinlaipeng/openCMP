package service

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
	"github.com/opencmp/opencmp/pkg/cloudprovider"
)

// CloudAccountService 云账户服务
type CloudAccountService struct {
	db *gorm.DB
}

// NewCloudAccountService 创建云账户服务
func NewCloudAccountService(db *gorm.DB) *CloudAccountService {
	return &CloudAccountService{db: db}
}

// CreateCloudAccount 创建云账户
func (s *CloudAccountService) CreateCloudAccount(ctx context.Context, account *model.CloudAccount) error {
	return s.db.WithContext(ctx).Create(account).Error
}

// GetCloudAccount 获取云账户
func (s *CloudAccountService) GetCloudAccount(ctx context.Context, id uint) (*model.CloudAccount, error) {
	var account model.CloudAccount
	err := s.db.WithContext(ctx).First(&account, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}

// ListCloudAccounts 列出云账户
func (s *CloudAccountService) ListCloudAccounts(ctx context.Context, limit, offset int) ([]*model.CloudAccount, int64, error) {
	var accounts []*model.CloudAccount
	var total int64

	if err := s.db.Model(&model.CloudAccount{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := s.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&accounts).Error

	return accounts, total, err
}

// UpdateCloudAccount 更新云账户
func (s *CloudAccountService) UpdateCloudAccount(ctx context.Context, account *model.CloudAccount) error {
	return s.db.WithContext(ctx).Save(account).Error
}

// DeleteCloudAccount 删除云账户
func (s *CloudAccountService) DeleteCloudAccount(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Delete(&model.CloudAccount{}, id).Error
}

// VerifyCloudAccount 验证云账户
func (s *CloudAccountService) VerifyCloudAccount(ctx context.Context, account *model.CloudAccount) (bool, error) {
	var creds map[string]string
	if err := json.Unmarshal(account.Credentials, &creds); err != nil {
		return false, err
	}

	config := cloudprovider.CloudAccountConfig{
		ID:           strconv.Itoa(int(account.ID)),
		Name:         account.Name,
		ProviderType: account.ProviderType,
		Credentials:  creds,
		Region:       creds["region_id"], // 使用凭证中的region_id，如果没有则为空，适配器会使用默认值
	}

	provider, err := cloudprovider.GetProvider(account.ProviderType, config)
	if err != nil {
		return false, err
	}

	// 尝试获取云厂商信息来验证连接
	cloudInfo := provider.GetCloudInfo()
	if cloudInfo.Provider == "" {
		return false, nil
	}

	return true, nil
}

// TestConnection 测试云账户连接
func (s *CloudAccountService) TestConnection(ctx context.Context, account *model.CloudAccount) (bool, error) {
	var creds map[string]string
	if err := json.Unmarshal(account.Credentials, &creds); err != nil {
		return false, err
	}

	config := cloudprovider.CloudAccountConfig{
		ID:           strconv.Itoa(int(account.ID)),
		Name:         account.Name,
		ProviderType: account.ProviderType,
		Credentials:  creds,
		Region:       creds["region_id"], // 使用凭证中的region_id，如果没有则为空，适配器会使用默认值
	}

	provider, err := cloudprovider.GetProvider(account.ProviderType, config)
	if err != nil {
		return false, err
	}

	// 尝试获取云厂商信息来测试连接
	cloudInfo := provider.GetCloudInfo()
	if cloudInfo.Provider == "" {
		return false, nil
	}

	return true, nil
}

// SyncResources 同步云账户资源
// 返回同步的资源统计（VM数量、VPC数量等）
func (s *CloudAccountService) SyncResources(ctx context.Context, account *model.CloudAccount) (map[string]int, error) {
	var creds map[string]string
	if err := json.Unmarshal(account.Credentials, &creds); err != nil {
		return nil, err
	}

	config := cloudprovider.CloudAccountConfig{
		ID:           strconv.Itoa(int(account.ID)),
		Name:         account.Name,
		ProviderType: account.ProviderType,
		Credentials:  creds,
		Region:       creds["region_id"],
	}

	provider, err := cloudprovider.GetProvider(account.ProviderType, config)
	if err != nil {
		return nil, err
	}

	stats := map[string]int{}

	// 同步虚拟机
	vms, err := provider.ListVMs(ctx, cloudprovider.VMListFilter{})
	if err == nil {
		stats["vms"] = len(vms)
	}

	// 同步 VPC
	vpcs, err := provider.ListVPCs(ctx, cloudprovider.VPCFilter{})
	if err == nil {
		stats["vpcs"] = len(vpcs)
	}

	// 同步子网
	subnets, err := provider.ListSubnets(ctx, cloudprovider.SubnetFilter{})
	if err == nil {
		stats["subnets"] = len(subnets)
	}

	// 同步安全组
	sgs, err := provider.ListSecurityGroups(ctx, cloudprovider.SGFilter{})
	if err == nil {
		stats["security_groups"] = len(sgs)
	}

	// 同步 EIP
	eips, err := provider.ListEIPs(ctx, cloudprovider.EIPFilter{})
	if err == nil {
		stats["eips"] = len(eips)
	}

	// 同步镜像
	images, err := provider.ListImages(ctx, cloudprovider.ImageFilter{})
	if err == nil {
		stats["images"] = len(images)
	}

	// 更新同步状态
	account.Status = string(model.CloudAccountStatusActive)
	s.db.WithContext(ctx).Save(account)

	return stats, nil
}

// VerifyCredentials 实际验证云账户凭证（通过调用真实 API）
func (s *CloudAccountService) VerifyCredentials(ctx context.Context, account *model.CloudAccount) (bool, string, error) {
	var creds map[string]string
	if err := json.Unmarshal(account.Credentials, &creds); err != nil {
		return false, "invalid credentials format", err
	}

	config := cloudprovider.CloudAccountConfig{
		ID:           strconv.Itoa(int(account.ID)),
		Name:         account.Name,
		ProviderType: account.ProviderType,
		Credentials:  creds,
		Region:       creds["region_id"],
	}

	provider, err := cloudprovider.GetProvider(account.ProviderType, config)
	if err != nil {
		return false, "failed to initialize provider: " + err.Error(), err
	}

	// 尝试实际调用云厂商 API 来验证凭证有效性
	// 不同云厂商使用不同的验证方式
	switch account.ProviderType {
	case "alibaba":
		// 尝试列出区域来验证
		regions, err := provider.ListRegions()
		if err != nil {
			return false, "failed to list regions: " + err.Error(), err
		}
		if len(regions) == 0 {
			return false, "no regions returned", nil
		}
		return true, "credentials verified, " + strconv.Itoa(len(regions)) + " regions available", nil

	case "tencent", "aws", "azure":
		// 尝试列出镜像来验证
		images, err := provider.ListImages(ctx, cloudprovider.ImageFilter{})
		if err != nil {
			return false, "failed to list images: " + err.Error(), err
		}
		return true, "credentials verified, " + strconv.Itoa(len(images)) + " images available", nil

	default:
		// 对于其他云厂商，使用基本验证
		cloudInfo := provider.GetCloudInfo()
		if cloudInfo.Provider == "" {
			return false, "provider info not available", nil
		}
		return true, "provider initialized", nil
	}
}

// TestConnectionWithCredentials 使用新凭证测试连接
func (s *CloudAccountService) TestConnectionWithCredentials(ctx context.Context, account *model.CloudAccount, accessKeyId, accessKeySecret string) (bool, string, []string, error) {
	// 构建新的凭证配置
	creds := map[string]string{
		"access_key_id":     accessKeyId,
		"access_key_secret": accessKeySecret,
		"region_id":         "cn-hangzhou", // 使用默认区域测试
	}

	config := cloudprovider.CloudAccountConfig{
		ID:           strconv.Itoa(int(account.ID)),
		Name:         account.Name,
		ProviderType: account.ProviderType,
		Credentials:  creds,
		Region:       "cn-hangzhou",
	}

	provider, err := cloudprovider.GetProvider(account.ProviderType, config)
	if err != nil {
		return false, "failed to initialize provider: " + err.Error(), []string{}, err
	}

	// 尝试列出区域来验证凭证有效性
	regions, err := provider.ListRegions()
	if err != nil {
		return false, "failed to list regions: " + err.Error(), []string{}, err
	}

	// 将区域转换为字符串数组
	regionNames := []string{}
	for _, region := range regions {
		regionNames = append(regionNames, region.Name)
	}

	if len(regions) == 0 {
		return false, "no regions returned", regionNames, nil
	}

	return true, "连接成功，" + strconv.Itoa(len(regions)) + " 个区域可用", regionNames, nil
}
