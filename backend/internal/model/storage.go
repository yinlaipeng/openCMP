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
	MaxIOPS        int            `json:"max_iops"`
	DiskFormat     string         `gorm:"size:50" json:"disk_format"`
	Type           string         `gorm:"size:50" json:"type"`
	StorageType    string         `gorm:"size:50" json:"storage_type"`
	MediumType     string         `gorm:"size:50" json:"medium_type"`
	Status         string         `gorm:"size:50" json:"status"`
	VMID           string         `gorm:"size:100" json:"vm_id"`
	VMName         string         `gorm:"size:200" json:"vm_name"`
	DeviceName     string         `gorm:"size:100" json:"device_name"`
	PrimaryStorage string         `gorm:"size:100" json:"primary_storage"`
	ZoneID         string         `gorm:"size:100" json:"zone_id"`
	RegionID       string         `gorm:"size:100" json:"region_id"`
	BillingType    string         `gorm:"size:50" json:"billing_type"`
	ShutdownReset  bool           `json:"shutdown_reset"`
	ProviderType   string         `gorm:"size:20" json:"provider_type"`
	AccountName    string         `gorm:"size:200" json:"account_name"`
	ProjectID      uint           `gorm:"index" json:"project_id"`
	ProjectName    string         `gorm:"size:200" json:"project_name"`
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
	DiskName       string    `gorm:"size:200" json:"disk_name"`
	DiskType       string    `gorm:"size:50" json:"disk_type"`
	Size           int       `json:"size"` // GB
	Status         string    `gorm:"size:50" json:"status"`
	Progress       int       `json:"progress"`
	VMID           string    `gorm:"size:100" json:"vm_id"`
	VMName         string    `gorm:"size:200" json:"vm_name"`
	ProviderType   string    `gorm:"size:20" json:"provider_type"`
	AccountName    string    `gorm:"size:200" json:"account_name"`
	RegionID       string    `gorm:"size:100" json:"region_id"`
	Description    string    `gorm:"size:500" json:"description"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// InstanceSnapshot 主机快照
type InstanceSnapshot struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	CloudAccountID  uint           `gorm:"index;not null" json:"cloud_account_id"`
	SnapshotID      string         `gorm:"size:100;uniqueIndex:instance_snapshot_idx" json:"snapshot_id"`
	Name            string         `gorm:"size:200" json:"name"`
	InstanceID      string         `gorm:"size:100;index" json:"instance_id"`
	InstanceName    string         `gorm:"size:200" json:"instance_name"`
	DiskSnapshots   int            `json:"disk_snapshots"`
	MemorySnapshot  bool           `json:"memory_snapshot"`
	CPUArch         string         `gorm:"size:20" json:"cpu_arch"`
	Size            int            `json:"size"` // GB
	Status          string         `gorm:"size:50" json:"status"`
	ProviderType    string         `gorm:"size:20" json:"provider_type"`
	AccountName     string         `gorm:"size:200" json:"account_name"`
	RegionID        string         `gorm:"size:100" json:"region_id"`
	Description     string         `gorm:"size:500" json:"description"`
	Tags            datatypes.JSON `gorm:"type:json" json:"tags"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

// SnapshotPolicy 自动快照策略
type SnapshotPolicy struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	CloudAccountID  uint      `gorm:"index;not null" json:"cloud_account_id"`
	PolicyID        string    `gorm:"size:100;uniqueIndex:policy_account_idx" json:"policy_id"`
	Name            string    `gorm:"size:200" json:"name"`
	Status          string    `gorm:"size:50" json:"status"` // active/inactive
	ResourceType    string    `gorm:"size:50" json:"resource_type"` // disk/instance
	ScheduleType    string    `gorm:"size:50" json:"schedule_type"` // daily/weekly/monthly
	ExecuteTime     string    `gorm:"size:20" json:"execute_time"` // HH:mm
	WeekDay         string    `gorm:"size:10" json:"week_day"` // 1-7
	MonthDay        string    `gorm:"size:10" json:"month_day"` // 1-28
	RetentionDays   int       `json:"retention_days"`
	AssociatedCount int       `json:"associated_count"`
	ProviderType    string    `gorm:"size:20" json:"provider_type"`
	AccountName     string    `gorm:"size:200" json:"account_name"`
	RegionID        string    `gorm:"size:100" json:"region_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// TableName methods
func (CloudDisk) TableName() string         { return "sync_cloud_disks" }
func (CloudSnapshot) TableName() string     { return "sync_cloud_snapshots" }
func (InstanceSnapshot) TableName() string  { return "sync_instance_snapshots" }
func (SnapshotPolicy) TableName() string    { return "sync_snapshot_policies" }