package model

import (
	"time"

	"gorm.io/gorm"
)

// CloudAccessGroup 云访问组模型 (用于多云管理页面的云用户组)
type CloudAccessGroup struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"size:100;not null" json:"name"`
	Status       string         `gorm:"size:20;default:'active'" json:"status"`
	Permissions  string         `gorm:"size:500" json:"permissions"`
	Platform     string         `gorm:"size:50" json:"platform"`
	CloudAccounts string        `gorm:"size:500" json:"cloud_accounts"`
	SharedScope  string         `gorm:"size:20;default:'private'" json:"shared_scope"`
	DomainID     uint           `json:"domain_id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (CloudAccessGroup) TableName() string {
	return "cloud_access_groups"
}

// ProxySetting 代理设置模型
type ProxySetting struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"size:100;not null" json:"name"`
	HttpsProxy   string         `gorm:"size:500" json:"https_proxy"`
	HttpProxy    string         `gorm:"size:500" json:"http_proxy"`
	NoProxy      string         `gorm:"size:500" json:"no_proxy"`
	SharedScope  string         `gorm:"size:20;default:'private'" json:"shared_scope"`
	DomainID     uint           `json:"domain_id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (ProxySetting) TableName() string {
	return "proxy_settings"
}