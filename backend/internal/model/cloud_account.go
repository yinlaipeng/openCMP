package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// CloudAccount 云账户配置
type CloudAccount struct {
	ID                       uint           `gorm:"primaryKey" json:"id"`
	Name                     string         `gorm:"uniqueIndex;not null;size:100" json:"name"`
	ProviderType             string         `gorm:"type:varchar(20);not null" json:"provider_type"` // alibaba/tencent/aws/azure
	Credentials              datatypes.JSON `gorm:"type:json" json:"credentials"`                   // 加密存储
	Status                   string         `gorm:"type:varchar(20);default:active" json:"status"`  // active/inactive/error
	Description              string         `gorm:"size:500" json:"description"`
	Enabled                  bool           `gorm:"default:true" json:"enabled"`
	HealthStatus             string         `gorm:"type:varchar(20);default:healthy" json:"health_status"` // healthy/unhealthy
	Balance                  float64        `gorm:"default:0.0" json:"balance"`
	AccountNumber            string         `gorm:"size:100" json:"account_number"`
	LastSync                 *time.Time     `json:"last_sync,omitempty"`
	SyncTime                 string         `gorm:"size:50" json:"sync_time"`
	DomainID                 uint           `gorm:"index;default:1" json:"domain_id"`                                // 默认分配到默认域
	ResourceAssignmentMethod string         `gorm:"size:50;default:'tag_mapping'" json:"resource_assignment_method"` // tag_mapping/project_mapping/manual_assignment
	CreatedAt                time.Time      `json:"created_at"`
	UpdatedAt                time.Time      `json:"updated_at"`
	DeletedAt                gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (CloudAccount) TableName() string {
	return "cloud_accounts"
}

// CloudAccountStatus 云账户状态
type CloudAccountStatus string

const (
	CloudAccountStatusActive   CloudAccountStatus = "active"
	CloudAccountStatusInactive CloudAccountStatus = "inactive"
	CloudAccountStatusError    CloudAccountStatus = "error"
)
