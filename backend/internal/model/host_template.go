package model

import (
	"encoding/json"
	"time"

	"github.com/opencmp/opencmp/pkg/cloudprovider"
	"gorm.io/gorm"
)

// HostTemplate 主机模版
type HostTemplate struct {
	ID             string                         `gorm:"primaryKey;size:64" json:"id"`                    // 主机模版 ID
	Name           string                         `gorm:"uniqueIndex;not null;size:100" json:"name"`      // 模版名称
	Description  string                         `gorm:"size:500" json:"description"`                   // 模版描述
	Status         string                         `gorm:"type:varchar(20);not null;default:'Draft'" json:"status"` // 模版状态
	InstanceType   string                         `gorm:"size:100" json:"instance_type"`                // 实例规格
	CPUArch        string                         `gorm:"size:20" json:"cpu_arch"`                       // CPU架构
	MemorySize     int                            `gorm:"default:0" json:"memory_size"`                  // 内存大小(MB)
	CPUCount       int                            `gorm:"default:0" json:"cpu_count"`                      // CPU核心数
	DiskSize       int                            `gorm:"default:0" json:"disk_size"`                      // 磁盘大小(GB)
	ImageID        string                         `gorm:"size:64" json:"image_id"`                         // 镜像ID
	OSName         string                         `gorm:"size:100" json:"os_name"`                         // 操作系统名称
	OSVersion      string                         `gorm:"size:50" json:"os_version"`                       // 操作系统版本
	VPCID          string                         `gorm:"size:64" json:"vpc_id"`                           // VPC ID
	SubnetID       string                         `gorm:"size:64" json:"subnet_id"`                        // 子网ID
	BillingMethod  string                         `gorm:"size:50" json:"billing_method"`                 // 计费方式
	Platform       string                         `gorm:"size:50" json:"platform"`                         // 平台
	ProjectID      string                         `gorm:"size:64;index" json:"project_id"`               // 项目ID
	RegionID       string                         `gorm:"size:64" json:"region_id"`                        // 区域ID
	ZoneID         string                         `gorm:"size:64" json:"zone_id"`                          // 可用区ID
	Tags           string                         `gorm:"type:text" json:"tags"`                           // 标签 (stored as JSON string in db)
	CreatedAt      time.Time                      `gorm:"index" json:"created_at"`
	UpdatedAt      time.Time                      `gorm:"index" json:"updated_at"`
	UpdatedBy      string                         `gorm:"size:64" json:"updated_by"`                       // 更新人ID
	DeletedAt      gorm.DeletedAt                 `gorm:"index" json:"-"`
}

// ConvertToCloudProvider converts the model to the cloudprovider type
func (ht *HostTemplate) ConvertToCloudProvider() *cloudprovider.HostTemplate {
	var tags map[string]string
	if ht.Tags != "" {
		json.Unmarshal([]byte(ht.Tags), &tags)
	}

	status := cloudprovider.HostTemplateStatus(ht.Status)
	if status == "" {
		status = "Draft"
	}

	return &cloudprovider.HostTemplate{
		ID:           ht.ID,
		Name:         ht.Name,
		Description:  ht.Description,
		Status:       status,
		InstanceType: ht.InstanceType,
		CPUArch:      ht.CPUArch,
		MemorySize:   ht.MemorySize,
		CPUCount:     ht.CPUCount,
		DiskSize:     ht.DiskSize,
		ImageID:      ht.ImageID,
		OSName:       ht.OSName,
		OSVersion:    ht.OSVersion,
		VPCID:        ht.VPCID,
		SubnetID:     ht.SubnetID,
		BillingMethod: ht.BillingMethod,
		Platform:     ht.Platform,
		ProjectID:    ht.ProjectID,
		RegionID:     ht.RegionID,
		ZoneID:       ht.ZoneID,
		Tags:         tags,
		CreatedAt:    ht.CreatedAt,
		UpdatedAt:    ht.UpdatedAt,
	}
}

// ConvertFromCloudProvider converts from the cloudprovider type to the model
func (ht *HostTemplate) ConvertFromCloudProvider(cpHostTemplate *cloudprovider.HostTemplate) {
	if cpHostTemplate == nil {
		return
	}

	ht.ID = cpHostTemplate.ID
	ht.Name = cpHostTemplate.Name
	ht.Description = cpHostTemplate.Description
	ht.Status = string(cpHostTemplate.Status)
	ht.InstanceType = cpHostTemplate.InstanceType
	ht.CPUArch = cpHostTemplate.CPUArch
	ht.MemorySize = cpHostTemplate.MemorySize
	ht.CPUCount = cpHostTemplate.CPUCount
	ht.DiskSize = cpHostTemplate.DiskSize
	ht.ImageID = cpHostTemplate.ImageID
	ht.OSName = cpHostTemplate.OSName
	ht.OSVersion = cpHostTemplate.OSVersion
	ht.VPCID = cpHostTemplate.VPCID
	ht.SubnetID = cpHostTemplate.SubnetID
	ht.BillingMethod = cpHostTemplate.BillingMethod
	ht.Platform = cpHostTemplate.Platform
	ht.ProjectID = cpHostTemplate.ProjectID
	ht.RegionID = cpHostTemplate.RegionID
	ht.ZoneID = cpHostTemplate.ZoneID

	if cpHostTemplate.Tags != nil {
		tagsBytes, _ := json.Marshal(cpHostTemplate.Tags)
		ht.Tags = string(tagsBytes)
	} else {
		ht.Tags = ""
	}
}

func (HostTemplate) TableName() string {
	return "host_templates"
}