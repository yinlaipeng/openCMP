package cloudprovider

// IBilling 账单和成本管理接口
type IBilling interface {
	// 获取账单列表
	ListBills(billingCycle string) ([]*BillItem, error)
	// 获取订单列表
	ListOrders(status string) ([]*OrderItem, error)
	// 获取待续费资源列表
	ListRenewalResources(daysThreshold int) ([]*RenewalResource, error)
	// 获取成本分析数据
	GetCostAnalysis(startDate, endDate string) ([]*CostData, error)
	// 获取账户余额
	GetAccountBalance() (*AccountBalance, error)
}

// BillItem 账单项
type BillItem struct {
	BillingCycle   string  `json:"billing_cycle"`
	ProductType    string  `json:"product_type"`
	ProductName    string  `json:"product_name"`
	InstanceID     string  `json:"instance_id"`
	UsageAmount    float64 `json:"usage_amount"`
	UnitPrice      float64 `json:"unit_price"`
	TotalCost      float64 `json:"total_cost"`
	Currency       string  `json:"currency"`
	BillingMethod  string  `json:"billing_method"`
	Status         string  `json:"status"`
}

// OrderItem 订单项
type OrderItem struct {
	OrderID       string  `json:"order_id"`
	OrderType     string  `json:"order_type"`
	ProductType   string  `json:"product_type"`
	ProductName   string  `json:"product_name"`
	InstanceID    string  `json:"instance_id"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	Status        string  `json:"status"`
	PaymentTime   string  `json:"payment_time"`
	EffectiveTime string  `json:"effective_time"`
	ExpireTime    string  `json:"expire_time"`
}

// RenewalResource 待续费资源
type RenewalResource struct {
	InstanceID    string  `json:"instance_id"`
	InstanceName  string  `json:"instance_name"`
	ProductType   string  `json:"product_type"`
	ExpireTime    string  `json:"expire_time"`
	DaysRemaining int     `json:"days_remaining"`
	RenewalPrice  float64 `json:"renewal_price"`
	Status        string  `json:"status"`
}

// CostData 成本分析数据
type CostData struct {
	Date       string  `json:"date"`
	Cost       float64 `json:"cost"`
	Product    string  `json:"product"`
	Region     string  `json:"region"`
	ResourceID string  `json:"resource_id"`
}

// AccountBalance 账户余额
type AccountBalance struct {
	Balance     float64 `json:"balance"`
	Currency    string  `json:"currency"`
	CreditLimit float64 `json:"credit_limit"`
	Status      string  `json:"status"`
}