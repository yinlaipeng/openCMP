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