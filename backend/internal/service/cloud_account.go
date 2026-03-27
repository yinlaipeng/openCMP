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
