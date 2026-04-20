package model

import (
	"errors"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ValidProviderTypes 有效的云厂商类型
var ValidProviderTypes = []string{"alibaba", "tencent", "aws", "azure", "huawei", "google", "openstack", "vmware"}

// CloudAccount 云账户配置
type CloudAccount struct {
	ID                       uint           `gorm:"primaryKey" json:"id"`
	Name                     string         `gorm:"uniqueIndex;not null;size:100" json:"name"`
	ProviderType             string         `gorm:"type:varchar(20);not null" json:"provider_type"` // alibaba/tencent/aws/azure
	Credentials              datatypes.JSON `gorm:"type:json" json:"credentials"`                   // 加密存储
	Status                   string         `gorm:"type:varchar(20);default:active" json:"status"`  // active/inactive/error/connected/disconnected/checking
	Description              string         `gorm:"size:500" json:"description"`
	Enabled                  bool           `gorm:"default:true" json:"enabled"`
	HealthStatus             string         `gorm:"type:varchar(20);default:healthy" json:"health_status"` // healthy/unhealthy
	Balance                  float64        `gorm:"default:0.0" json:"balance"`
	AccountNumber            string         `gorm:"size:100" json:"account_number"`
	LastSync                 *time.Time     `json:"last_sync,omitempty"`
	SyncTime                 string         `gorm:"size:50" json:"sync_time"`
	DomainID                 uint           `gorm:"index;default:1" json:"domain_id"`                                // 默认分配到默认域
	SyncPolicyID             *uint          `gorm:"index" json:"sync_policy_id"`                                     // 绑定的同步策略ID
	ResourceAssignmentMethod string         `gorm:"size:50;default:'tag_mapping'" json:"resource_assignment_method"` // tag_mapping/project_mapping/manual_assignment
	LastConnectionCheckTime  *time.Time     `json:"last_connection_check_time,omitempty"`                           // 最近连接检测时间
	ConnectionCheckError     string         `gorm:"size:500" json:"connection_check_error"`                          // 连接检测失败原因
	CreatedAt                time.Time      `json:"created_at"`
	UpdatedAt                time.Time      `json:"updated_at"`
	DeletedAt                gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (CloudAccount) TableName() string {
	return "cloud_accounts"
}

// BeforeSave GORM钩子，保存前验证
func (a *CloudAccount) BeforeSave(tx *gorm.DB) error {
	// 验证ProviderType是否有效
	if !isValidProviderType(a.ProviderType) {
		return errors.New("invalid provider_type: must be one of alibaba, tencent, aws, azure, huawei, google, openstack, vmware")
	}
	return nil
}

// isValidProviderType 检查云厂商类型是否有效
func isValidProviderType(providerType string) bool {
	for _, valid := range ValidProviderTypes {
		if valid == providerType {
			return true
		}
	}
	return false
}

// CloudAccountStatus 云账户状态
type CloudAccountStatus string

const (
	CloudAccountStatusActive      CloudAccountStatus = "active"      // 业务激活（旧状态，兼容）
	CloudAccountStatusInactive    CloudAccountStatus = "inactive"    // 业务禁用（旧状态，兼容）
	CloudAccountStatusError       CloudAccountStatus = "error"       // 错误（旧状态，兼容）
	CloudAccountStatusConnected   CloudAccountStatus = "connected"   // 已连接（连接状态）
	CloudAccountStatusDisconnected CloudAccountStatus = "disconnected" // 连接断开
	CloudAccountStatusChecking    CloudAccountStatus = "checking"    // 检测中
)
