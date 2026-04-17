package model

import (
	"time"
)

// ResourceStateLog 资源状态变更日志
type ResourceStateLog struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	ResourceType   string    `json:"resource_type" gorm:"column:resource_type;size:50;index"`      // vm, disk, vpc, etc.
	ResourceID     string    `json:"resource_id" gorm:"column:resource_id;size:255;index"`          // 云厂商资源ID
	ResourceName   string    `json:"resource_name" gorm:"column:resource_name;size:255"`            // 资源名称
	CloudAccountID uint      `json:"cloud_account_id" gorm:"column:cloud_account_id;index"`        // 云账号ID
	PreviousStatus string    `json:"previous_status" gorm:"column:previous_status;size:50"`         // 原状态
	CurrentStatus  string    `json:"current_status" gorm:"column:current_status;size:50;index"`     // 新状态
	OperationType  string    `json:"operation_type" gorm:"column:operation_type;size:50"`           // start, stop, reboot, delete, create
	Operator       string    `json:"operator" gorm:"column:operator;size:100"`                      // 操作人
	OperatorID     *uint     `json:"operator_id" gorm:"column:operator_id"`                         // 操作人ID
	ProjectID      *uint     `json:"project_id" gorm:"column:project_id;index"`                     // 项目ID
	Reason         string    `json:"reason" gorm:"column:reason;size:500"`                          // 变更原因
	Details        string    `json:"details" gorm:"column:details;type:text"`                       // 详细信息JSON
	OccurredAt     time.Time `json:"occurred_at" gorm:"column:occurred_at;index"`                   // 变更时间
	CreatedAt      time.Time `json:"created_at"`
}

func (ResourceStateLog) TableName() string {
	return "resource_state_logs"
}