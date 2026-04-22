package model

import (
	"encoding/json"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// WAFInstance WAF策略实例
type WAFInstance struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Name           string         `gorm:"size:200;not null" json:"name"`
	Status         string         `gorm:"size:20;default:'normal'" json:"status"` // normal, creating, deleting, error
	Type           string         `gorm:"size:50" json:"type"`                    // WAF类型
	Platform       string         `gorm:"size:20" json:"platform"`                // alibaba, tencent, aws, azure
	CloudAccountID uint           `gorm:"index" json:"cloud_account_id"`
	CloudAccount   *CloudAccount  `gorm:"foreignKey:CloudAccountID" json:"cloud_account,omitempty"`
	DomainID       uint           `gorm:"index;default:1" json:"domain_id"`
	RegionID       string         `gorm:"size:50" json:"region_id"`
	ExternalID     string         `gorm:"size:100" json:"external_id"`            // 云厂商侧ID
	Tags           datatypes.JSON `gorm:"type:json" json:"tags"`
	Description    string         `gorm:"size:500" json:"description"`
	Enabled        bool           `gorm:"default:true" json:"enabled"`
	SyncTime       *time.Time     `json:"sync_time"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (WAFInstance) TableName() string {
	return "waf_instances"
}

// WebappInstance 应用程序服务实例
type WebappInstance struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Name           string         `gorm:"size:200;not null" json:"name"`
	Status         string         `gorm:"size:20;default:'normal'" json:"status"` // normal, running, stopped, error
	Stack          string         `gorm:"size:50" json:"stack"`                   // 技术栈: tomcat, nginx, nodejs
	OsType         string         `gorm:"size:50" json:"os_type"`                 // 操作系统类型
	IpAddr         string         `gorm:"size:50" json:"ip_addr"`                 // IP地址
	Domain         string         `gorm:"size:100" json:"domain"`                 // 域名
	ServerFarm     string         `gorm:"size:100" json:"server_farm"`            // 服务器组
	Platform       string         `gorm:"size:20" json:"platform"`                // alibaba, tencent, aws, azure
	CloudAccountID uint           `gorm:"index" json:"cloud_account_id"`
	CloudAccount   *CloudAccount  `gorm:"foreignKey:CloudAccountID" json:"cloud_account,omitempty"`
	RegionID       string         `gorm:"size:50" json:"region_id"`
	ProjectID      uint           `gorm:"index" json:"project_id"`
	ExternalID     string         `gorm:"size:100" json:"external_id"`            // 云厂商侧ID
	Tags           datatypes.JSON `gorm:"type:json" json:"tags"`
	Description    string         `gorm:"size:500" json:"description"`
	Enabled        bool           `gorm:"default:true" json:"enabled"`
	SyncTime       *time.Time     `json:"sync_time"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (WebappInstance) TableName() string {
	return "webapp_instances"
}

// MapToJSON 将 map 转换为 JSON
func MapToJSON(m map[string]string) datatypes.JSON {
	if m == nil {
		return nil
	}
	data, _ := json.Marshal(m)
	return data
}