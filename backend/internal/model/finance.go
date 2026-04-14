// backend/internal/model/finance.go
package model

import (
	"time"
)

// Bill 账单记录
type Bill struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	CloudAccountID uint      `gorm:"index;not null" json:"cloud_account_id"`
	BillingCycle   string    `gorm:"type:varchar(20);not null" json:"billing_cycle"` // 2026-04
	ProductType    string    `gorm:"type:varchar(50)" json:"product_type"`
	ProductName    string    `gorm:"size:200" json:"product_name"`
	InstanceID     string    `gorm:"size:100" json:"instance_id"`
	UsageAmount    float64   `json:"usage_amount"`
	UnitPrice      float64   `json:"unit_price"`
	TotalCost      float64   `json:"total_cost"`
	Currency       string    `gorm:"type:varchar(10);default:'CNY'" json:"currency"`
	BillingMethod  string    `gorm:"type:varchar(20)" json:"billing_method"`
	Status         string    `gorm:"type:varchar(20)" json:"status"`
	ProviderType   string    `gorm:"type:varchar(20)" json:"provider_type"`
	CreatedAt      time.Time `json:"created_at"`
}

// Order 订单记录
type Order struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	CloudAccountID uint       `gorm:"index;not null" json:"cloud_account_id"`
	OrderID        string     `gorm:"uniqueIndex;size:100" json:"order_id"`
	OrderType      string     `gorm:"type:varchar(20)" json:"order_type"`
	ProductType    string     `gorm:"type:varchar(50)" json:"product_type"`
	ProductName    string     `gorm:"size:200" json:"product_name"`
	InstanceID     string     `gorm:"size:100" json:"instance_id"`
	Amount         float64    `json:"amount"`
	Currency       string     `gorm:"type:varchar(10);default:'CNY'" json:"currency"`
	Status         string     `gorm:"type:varchar(20)" json:"status"`
	PaymentTime    *time.Time `json:"payment_time"`
	EffectiveTime  time.Time  `json:"effective_time"`
	ExpireTime     *time.Time `json:"expire_time"`
	ProviderType   string     `gorm:"type:varchar(20)" json:"provider_type"`
	CreatedAt      time.Time  `json:"created_at"`
}

// Budget 预算配置
type Budget struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	CloudAccountID uint      `gorm:"index" json:"cloud_account_id"`
	Name           string    `gorm:"size:100;not null" json:"name"`
	Type           string    `gorm:"type:varchar(20)" json:"type"`
	Amount         float64   `json:"amount"`
	AlertThreshold float64   `json:"alert_threshold"`
	CurrentUsage   float64   `json:"current_usage"`
	Status         string    `gorm:"type:varchar(20);default:'active'" json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// CostAnomaly 成本异常记录
type CostAnomaly struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	CloudAccountID uint      `gorm:"index;not null" json:"cloud_account_id"`
	AnomalyType    string    `gorm:"type:varchar(50)" json:"anomaly_type"`
	DetectedAt     time.Time `json:"detected_at"`
	Period         string    `gorm:"size:20" json:"period"`
	ExpectedCost   float64   `json:"expected_cost"`
	ActualCost     float64   `json:"actual_cost"`
	DeviationRate  float64   `json:"deviation_rate"`
	Severity       string    `gorm:"type:varchar(20)" json:"severity"`
	Status         string    `gorm:"type:varchar(20)" json:"status"`
	Resolution     string    `gorm:"size:500" json:"resolution"`
	CreatedAt      time.Time `json:"created_at"`
}

// RenewalResource 待续费资源
type RenewalResource struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	CloudAccountID uint      `gorm:"index;not null" json:"cloud_account_id"`
	InstanceID     string    `gorm:"size:100" json:"instance_id"`
	InstanceName   string    `gorm:"size:200" json:"instance_name"`
	ProductType    string    `gorm:"type:varchar(50)" json:"product_type"`
	ExpireTime     time.Time `json:"expire_time"`
	DaysRemaining  int       `json:"days_remaining"`
	RenewalPrice   float64   `json:"renewal_price"`
	Status         string    `gorm:"type:varchar(20)" json:"status"`
}

// TableName 方法
func (Bill) TableName() string         { return "finance_bills" }
func (Order) TableName() string        { return "finance_orders" }
func (Budget) TableName() string       { return "finance_budgets" }
func (CostAnomaly) TableName() string  { return "finance_cost_anomalies" }
func (RenewalResource) TableName() string { return "finance_renewal_resources" }