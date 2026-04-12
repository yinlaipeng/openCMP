package model

import (
	"time"

	"gorm.io/gorm"
)

// AutoscalingGroup 弹性伸缩组
type AutoscalingGroup struct {
	ID              string         `gorm:"primaryKey;size:64" json:"id"`                      // 弹性伸缩组 ID
	Name            string         `gorm:"uniqueIndex;not null;size:100" json:"name"`       // 弹性伸缩组名称
	Description   string         `gorm:"size:500" json:"description"`                     // 描述
	Status          string         `gorm:"type:varchar(20);not null;default:'Inactive'" json:"status"` // 伸缩组状态
	HostTemplateID  string         `gorm:"size:64;index" json:"host_template_id"`          // 主机模版ID
	CurrentCapacity int            `gorm:"default:0" json:"current_capacity"`               // 当前实例数
	DesiredCapacity int            `gorm:"default:1" json:"desired_capacity"`               // 期望实例数
	MinSize         int            `gorm:"default:0" json:"min_size"`                       // 最小实例数
	MaxSize         int            `gorm:"default:10" json:"max_size"`                      // 最大实例数
	Platform        string         `gorm:"size:50" json:"platform"`                        // 平台
	ProjectID       string         `gorm:"size:64;index" json:"project_id"`                // 项目ID
	RegionID        string         `gorm:"size:64" json:"region_id"`                       // 区域ID
	ZoneID          string         `gorm:"size:64" json:"zone_id"`                         // 可用区ID
	Tags            string         `gorm:"type:text" json:"tags"`                          // 标签 (JSON格式)
	CreatedAt       time.Time      `gorm:"index" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"index" json:"updated_at"`
	UpdatedBy       string         `gorm:"size:64" json:"updated_by"`                       // 更新人ID
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (AutoscalingGroup) TableName() string {
	return "autoscaling_groups"
}