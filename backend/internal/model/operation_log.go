package model

import (
	"time"
)

// OperationLog represents an operation log entry
type OperationLog struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	OperationTime time.Time `json:"operation_time" gorm:"column:operation_time;index"`
	ResourceName  string    `json:"resource_name" gorm:"column:resource_name;size:255;index"`
	ResourceType  string    `json:"resource_type" gorm:"column:resource_type;size:100;index"`
	ResourceID    uint      `json:"resource_id" gorm:"column:resource_id;index"`
	OperationType string    `json:"operation_type" gorm:"column:operation_type;size:100"`
	ServiceType   string    `json:"service_type" gorm:"column:service_type;size:100"`
	RiskLevel     string    `json:"risk_level" gorm:"column:risk_level;size:20;default:'medium'"`
	TimeType      string    `json:"time_type" gorm:"column:time_type;size:50;default:'realtime'"`
	Result        string    `json:"result" gorm:"column:result;size:50;default:'success'"`
	Operator      string    `json:"operator" gorm:"column:operator;size:100"`
	ProjectID     *uint     `json:"project_id" gorm:"column:project_id;index"`
	DomainID      *uint     `json:"domain_id,omitempty" gorm:"column:domain_id;index"`
	UserID        *uint     `json:"user_id,omitempty" gorm:"column:user_id;index"`
	CloudAccountID *uint    `json:"cloud_account_id,omitempty" gorm:"column:cloud_account_id;index"` // 新增：关联云账户
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (OperationLog) TableName() string {
	return "operation_logs"
}
