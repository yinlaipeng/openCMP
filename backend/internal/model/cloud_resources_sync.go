package model

import (
	"time"

	"gorm.io/datatypes"
)

// CloudVM 云虚拟机（同步后的本地存储）
type CloudVM struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	CloudAccountID uint           `gorm:"index;not null" json:"cloud_account_id"`
	InstanceID     string         `gorm:"size:100;index" json:"instance_id"`       // 云平台实例ID
	Name           string         `gorm:"size:200" json:"name"`
	Status         string         `gorm:"size:50" json:"status"`                  // running/stopped/terminated
	InstanceType   string         `gorm:"size:100" json:"instance_type"`
	ImageID        string         `gorm:"size:100" json:"image_id"`
	OSName         string         `gorm:"size:100" json:"os_name"`
	VPCID          string         `gorm:"size:100" json:"vpc_id"`
	SubnetID       string         `gorm:"size:100" json:"subnet_id"`
	PrivateIP      string         `gorm:"size:50" json:"private_ip"`
	PublicIP       string         `gorm:"size:50" json:"public_ip"`
	RegionID       string         `gorm:"size:100" json:"region_id"`
	ZoneID         string         `gorm:"size:100" json:"zone_id"`
	ProjectID      uint           `gorm:"index" json:"project_id"`                // 项目归属
	Tags           datatypes.JSON `gorm:"type:json" json:"tags"`                  // 资源标签
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

func (CloudVM) TableName() string { return "sync_cloud_vms" }

// CloudVPC 云VPC（同步后的本地存储）
type CloudVPC struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	CloudAccountID uint           `gorm:"index;not null" json:"cloud_account_id"`
	VPCID          string         `gorm:"size:100;index" json:"vpc_id"`
	Name           string         `gorm:"size:200" json:"name"`
	CIDR           string         `gorm:"size:50" json:"cidr"`
	Status         string         `gorm:"size:50" json:"status"`
	RegionID       string         `gorm:"size:100" json:"region_id"`
	ProjectID      uint           `gorm:"index" json:"project_id"`
	Tags           datatypes.JSON `gorm:"type:json" json:"tags"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

func (CloudVPC) TableName() string { return "sync_cloud_vpcs" }

// CloudSubnet 云子网（同步后的本地存储）
type CloudSubnet struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	CloudAccountID uint           `gorm:"index;not null" json:"cloud_account_id"`
	SubnetID       string         `gorm:"size:100;index" json:"subnet_id"`
	Name           string         `gorm:"size:200" json:"name"`
	VPCID          string         `gorm:"size:100;index" json:"vpc_id"`
	CIDR           string         `gorm:"size:50" json:"cidr"`
	ZoneID         string         `gorm:"size:100" json:"zone_id"`
	Status         string         `gorm:"size:50" json:"status"`
	ProjectID      uint           `gorm:"index" json:"project_id"`
	Tags           datatypes.JSON `gorm:"type:json" json:"tags"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

func (CloudSubnet) TableName() string { return "sync_cloud_subnets" }

// CloudSecurityGroup 云安全组（同步后的本地存储）
type CloudSecurityGroup struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	CloudAccountID uint           `gorm:"index;not null" json:"cloud_account_id"`
	SecurityGroupID string        `gorm:"size:100;index" json:"security_group_id"`
	Name           string         `gorm:"size:200" json:"name"`
	Description    string         `gorm:"size:500" json:"description"`
	VPCID          string         `gorm:"size:100;index" json:"vpc_id"`
	Status         string         `gorm:"size:50;default:'available'" json:"status"` // 默认可用状态
	ProjectID      uint           `gorm:"index" json:"project_id"`
	Tags           datatypes.JSON `gorm:"type:json" json:"tags"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

func (CloudSecurityGroup) TableName() string { return "sync_cloud_security_groups" }

// CloudEIP 云弹性公网IP（同步后的本地存储）
type CloudEIP struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	CloudAccountID uint      `gorm:"index;not null" json:"cloud_account_id"`
	EIPID          string    `gorm:"size:100;index" json:"eip_id"`
	Address        string    `gorm:"size:50" json:"address"`
	Bandwidth      int       `json:"bandwidth"` // Mbps
	Status         string    `gorm:"size:50" json:"status"`
	ResourceID     string    `gorm:"size:100" json:"resource_id"`     // 绑定的资源ID
	ResourceType   string    `gorm:"size:50" json:"resource_type"`    // 绑定的资源类型
	RegionID       string    `gorm:"size:100" json:"region_id"`
	ProjectID      uint      `gorm:"index" json:"project_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (CloudEIP) TableName() string { return "sync_cloud_eips" }

// CloudImage 云镜像（同步后的本地存储）
type CloudImage struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	CloudAccountID uint           `gorm:"index;not null" json:"cloud_account_id"`
	ImageID        string         `gorm:"size:100;index" json:"image_id"`
	Name           string         `gorm:"size:200" json:"name"`
	OSName         string         `gorm:"size:100" json:"os_name"`
	OSVersion      string         `gorm:"size:50" json:"os_version"`
	Architecture   string         `gorm:"size:50" json:"architecture"` // CPU架构
	Status         string         `gorm:"size:50" json:"status"`
	Tags           datatypes.JSON `gorm:"type:json" json:"tags"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

func (CloudImage) TableName() string { return "sync_cloud_images" }

// CloudRDS 云RDS数据库（同步后的本地存储）
type CloudRDS struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	CloudAccountID uint           `gorm:"index;not null" json:"cloud_account_id"`
	RDSID          string         `gorm:"size:100;index" json:"rds_id"`
	Name           string         `gorm:"size:200" json:"name"`
	Engine         string         `gorm:"size:50" json:"engine"`       // MySQL/PostgreSQL
	EngineVersion  string         `gorm:"size:50" json:"engine_version"`
	InstanceType   string         `gorm:"size:100" json:"instance_type"`
	Status         string         `gorm:"size:50" json:"status"`
	RegionID       string         `gorm:"size:100" json:"region_id"`
	ZoneID         string         `gorm:"size:100" json:"zone_id"`
	ProjectID      uint           `gorm:"index" json:"project_id"`
	Tags           datatypes.JSON `gorm:"type:json" json:"tags"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

func (CloudRDS) TableName() string { return "sync_cloud_rds" }

// CloudRedis 云Redis（同步后的本地存储）
type CloudRedis struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	CloudAccountID uint           `gorm:"index;not null" json:"cloud_account_id"`
	RedisID        string         `gorm:"size:100;index" json:"redis_id"`
	Name           string         `gorm:"size:200" json:"name"`
	InstanceType   string         `gorm:"size:100" json:"instance_type"`
	Status         string         `gorm:"size:50" json:"status"`
	RegionID       string         `gorm:"size:100" json:"region_id"`
	ZoneID         string         `gorm:"size:100" json:"zone_id"`
	ProjectID      uint           `gorm:"index" json:"project_id"`
	Tags           datatypes.JSON `gorm:"type:json" json:"tags"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

func (CloudRedis) TableName() string { return "sync_cloud_redis" }