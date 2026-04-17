package model

import (
	"time"

	"gorm.io/datatypes"
)

// CloudDisk 云硬盘
type CloudDisk struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	CloudAccountID uint           `gorm:"index;not null" json:"cloud_account_id"`
	DiskID         string         `gorm:"size:100;uniqueIndex:disk_account_idx" json:"disk_id"`
	Name           string         `gorm:"size:200" json:"name"`
	Size           int            `json:"size"` // GB
	Type           string         `gorm:"size:50" json:"type"`
	Status         string         `gorm:"size:50" json:"status"`
	VMID           string         `gorm:"size:100" json:"vm_id"`
	ZoneID         string         `gorm:"size:100" json:"zone_id"`
	RegionID       string         `gorm:"size:100" json:"region_id"`
	ProviderType   string         `gorm:"size:20" json:"provider_type"`
	ProjectID      uint           `gorm:"index" json:"project_id"`
	Tags           datatypes.JSON `gorm:"type:json" json:"tags"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

// CloudSnapshot 云快照
type CloudSnapshot struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	CloudAccountID uint      `gorm:"index;not null" json:"cloud_account_id"`
	SnapshotID     string    `gorm:"size:100;uniqueIndex:snapshot_account_idx" json:"snapshot_id"`
	Name           string    `gorm:"size:200" json:"name"`
	DiskID         string    `gorm:"size:100;index" json:"disk_id"`
	Size           int       `json:"size"` // GB
	Status         string    `gorm:"size:50" json:"status"`
	RegionID       string    `gorm:"size:100" json:"region_id"`
	ProviderType   string    `gorm:"size:20" json:"provider_type"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// TableName methods
func (CloudDisk) TableName() string     { return "sync_cloud_disks" }
func (CloudSnapshot) TableName() string { return "sync_cloud_snapshots" }