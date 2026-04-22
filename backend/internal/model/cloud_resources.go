package model

import (
	"time"

	"gorm.io/gorm"
)

// CloudSubscription 云订阅
type CloudSubscription struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	CloudAccountID   uint           `gorm:"index;not null" json:"cloud_account_id"`
	Name             string         `gorm:"size:200;not null" json:"name"`
	SubscriptionID   string         `gorm:"size:100;not null" json:"subscription_id"`
	Enabled          bool           `gorm:"default:true" json:"enabled"`
	Status           string         `gorm:"size:20;default:'normal'" json:"status"`
	SyncTime         *time.Time     `json:"sync_time"`
	SyncDuration     int            `gorm:"default:0" json:"sync_duration"` // 秒
	SyncStatus       string         `gorm:"size:20;default:'completed'" json:"sync_status"`
	DomainID         uint           `gorm:"index;default:1" json:"domain_id"`
	DefaultProjectID *uint          `json:"default_project_id"`
	SyncPolicyID     *uint          `json:"sync_policy_id"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

func (CloudSubscription) TableName() string {
	return "cloud_subscriptions"
}

// CloudUser 云用户
type CloudUser struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	CloudAccountID uint           `gorm:"index;not null" json:"cloud_account_id"`
	Username       string         `gorm:"size:100;not null" json:"username"`
	ConsoleLogin   bool           `gorm:"default:false" json:"console_login"`
	Status         string         `gorm:"size:20;default:'normal'" json:"status"`
	Password       string         `gorm:"size:255" json:"-"` // 加密存储，不返回前端
	LoginURL       string         `gorm:"size:255" json:"login_url"`
	LocalUserID    *uint          `json:"local_user_id"`
	Platform       string         `gorm:"size:20" json:"platform"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (CloudUser) TableName() string {
	return "cloud_users"
}

// CloudUserGroup 云用户组
type CloudUserGroup struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	CloudAccountID uint           `gorm:"index;not null" json:"cloud_account_id"`
	Name           string         `gorm:"size:200;not null" json:"name"`
	Status         string         `gorm:"size:20;default:'normal'" json:"status"`
	Permissions    string         `gorm:"size:500" json:"permissions"`
	Platform       string         `gorm:"size:20" json:"platform"`
	DomainID       uint           `gorm:"index;default:1" json:"domain_id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (CloudUserGroup) TableName() string {
	return "cloud_user_groups"
}

// CloudProject 云上项目
type CloudProject struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	CloudAccountID uint           `gorm:"index;not null" json:"cloud_account_id"`
	Name           string         `gorm:"size:200;not null" json:"name"`
	SubscriptionID *uint          `json:"subscription_id"`
	Status         string         `gorm:"size:20;default:'normal'" json:"status"`
	Tags           string         `gorm:"size:500" json:"tags"` // JSON格式存储标签
	DomainID       uint           `gorm:"index;default:1" json:"domain_id"`
	LocalProjectID *uint          `json:"local_project_id"`
	Priority       int            `gorm:"default:0" json:"priority"`
	SyncTime       *time.Time     `json:"sync_time"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (CloudProject) TableName() string {
	return "cloud_projects"
}

// Image 镜像
type Image struct {
	ID              string         `gorm:"primaryKey;size:64" json:"id"`
	Name            string         `gorm:"size:200;not null" json:"name"`
	Description     string         `gorm:"size:500" json:"description"`
	Status          string         `gorm:"size:20;default:'Creating'" json:"status"`
	Format          string         `gorm:"size:20" json:"format"`       // qcow2, raw, vhd, iso
	OsName          string         `gorm:"size:50" json:"os_name"`      // CentOS, Ubuntu, Windows
	OsVersion       string         `gorm:"size:20" json:"os_version"`   // 20.04, 7.9
	Size            int64          `json:"size"`                        // bytes
	CpuArch         string         `gorm:"size:20;default:'x86_64'" json:"cpu_arch"` // x86_64, arm64
	ImageType       string         `gorm:"size:20;default:'system'" json:"image_type"` // system, custom
	ShareScope      string         `gorm:"size:20;default:'private'" json:"share_scope"` // private, public
	CloudAccountID  uint           `gorm:"index" json:"cloud_account_id"`
	ProjectID       string         `gorm:"size:64" json:"project_id"`
	ExternalID      string         `gorm:"size:100" json:"external_id"` // 云厂商镜像ID
	Platform        string         `gorm:"size:20" json:"platform"`     // alibaba, tencent, aws, azure
	RegionID        string         `gorm:"size:64" json:"region_id"`
	Checksum        string         `gorm:"size:128" json:"checksum"`    // MD5/SHA256
	MinDiskSize     int            `json:"min_disk_size"`               // 最小磁盘大小 GB
	MinMemorySize   int            `json:"min_memory_size"`             // 最小内存大小 MB
	IsPublic        bool           `gorm:"default:false" json:"is_public"`
	Tags            string         `gorm:"size:500" json:"tags"`        // JSON 格式
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Image) TableName() string {
	return "images"
}