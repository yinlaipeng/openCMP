package alibaba

import (
	"encoding/json"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/opencmp/opencmp/pkg/cloudprovider"
)

// BillingClient 阿里云BSS客户端包装
type BillingClient struct {
	client *sdk.Client
}

// NewBillingClient 创建BSS客户端
func NewBillingClient(accessKeyID, accessKeySecret string) (*BillingClient, error) {
	cred := credentials.NewAccessKeyCredential(accessKeyID, accessKeySecret)
	client, err := sdk.NewClientWithOptions("cn-hangzhou", sdk.NewConfig(), cred)
	if err != nil {
		return nil, err
	}
	return &BillingClient{client: client}, nil
}

// ListBills 获取账单列表
func (p *AlibabaProvider) ListBills(billingCycle string) ([]*cloudprovider.BillItem, error) {
	req := requests.NewCommonRequest()
	req.Method = requests.POST
	req.Scheme = "https"
	req.Domain = "business.aliyuncs.com"
	req.Version = "2017-12-14"
	req.ApiName = "QueryBill"
	req.QueryParams["BillingCycle"] = billingCycle
	req.QueryParams["ProductCode"] = ""
	req.QueryParams["SubscriptionType"] = "All"
	req.QueryParams["PageNum"] = "1"
	req.QueryParams["PageSize"] = "100"

	resp, err := p.sendBillingRequest(req)
	if err != nil {
		return nil, err
	}

	return parseBillResponse(resp)
}

// ListOrders 获取订单列表
func (p *AlibabaProvider) ListOrders(status string) ([]*cloudprovider.OrderItem, error) {
	req := requests.NewCommonRequest()
	req.Method = requests.POST
	req.Scheme = "https"
	req.Domain = "business.aliyuncs.com"
	req.Version = "2017-12-14"
	req.ApiName = "QueryOrders"
	req.QueryParams["Status"] = status
	req.QueryParams["PageNum"] = "1"
	req.QueryParams["PageSize"] = "100"

	resp, err := p.sendBillingRequest(req)
	if err != nil {
		return nil, err
	}

	return parseOrderResponse(resp)
}

// ListRenewalResources 获取待续费资源列表
func (p *AlibabaProvider) ListRenewalResources(daysThreshold int) ([]*cloudprovider.RenewalResource, error) {
	req := requests.NewCommonRequest()
	req.Method = requests.POST
	req.Scheme = "https"
	req.Domain = "business.aliyuncs.com"
	req.Version = "2017-12-14"
	req.ApiName = "QueryRenewalResources"
	req.QueryParams["ProductType"] = ""
	req.QueryParams["PageNum"] = "1"
	req.QueryParams["PageSize"] = "100"

	resp, err := p.sendBillingRequest(req)
	if err != nil {
		return nil, err
	}

	return parseRenewalResponse(resp, daysThreshold)
}

// GetCostAnalysis 获取成本分析数据
func (p *AlibabaProvider) GetCostAnalysis(startDate, endDate string) ([]*cloudprovider.CostData, error) {
	req := requests.NewCommonRequest()
	req.Method = requests.POST
	req.Scheme = "https"
	req.Domain = "business.aliyuncs.com"
	req.Version = "2017-12-14"
	req.ApiName = "QueryInstanceBill"
	req.QueryParams["BillingCycle"] = startDate[:7]
	req.QueryParams["Granularity"] = "DAILY"
	req.QueryParams["PageNum"] = "1"
	req.QueryParams["PageSize"] = "100"

	resp, err := p.sendBillingRequest(req)
	if err != nil {
		return nil, err
	}

	return parseCostAnalysisResponse(resp)
}

// GetAccountBalance 获取账户余额
func (p *AlibabaProvider) GetAccountBalance() (*cloudprovider.AccountBalance, error) {
	req := requests.NewCommonRequest()
	req.Method = requests.POST
	req.Scheme = "https"
	req.Domain = "business.aliyuncs.com"
	req.Version = "2017-12-14"
	req.ApiName = "QueryAccountBalance"

	resp, err := p.sendBillingRequest(req)
	if err != nil {
		return nil, err
	}

	return parseBalanceResponse(resp)
}

// sendBillingRequest 发送BSS请求
func (p *AlibabaProvider) sendBillingRequest(req *requests.CommonRequest) (*responses.CommonResponse, error) {
	accessKeyID := p.config.Credentials["access_key_id"]
	accessKeySecret := p.config.Credentials["access_key_secret"]

	cred := credentials.NewAccessKeyCredential(accessKeyID, accessKeySecret)
	client, err := sdk.NewClientWithOptions("cn-hangzhou", sdk.NewConfig(), cred)
	if err != nil {
		return nil, err
	}

	resp, err := client.ProcessCommonRequest(req)
	if err != nil {
		if se, ok := err.(*errors.ServerError); ok {
			return nil, cloudprovider.NewCloudError(
				cloudprovider.ErrProviderError,
				se.Error(),
				se.ErrorCode(),
			)
		}
		return nil, err
	}

	return resp, nil
}

// ========== 响应结构体定义 ==========

// BillListResponse 账单列表响应
type BillListResponse struct {
	Code      string `json:"Code"`
	Message   string `json:"Message"`
	Success   bool   `json:"Success"`
	Data      struct {
		PageNum    int `json:"PageNum"`
	 PageSize   int `json:"PageSize"`
	 TotalCount int `json:"TotalCount"`
	 Items      []BillItemRaw `json:"Items"`
	} `json:"Data"`
}

type BillItemRaw struct {
	BillingCycle     string  `json:"BillingCycle"`
	ProductCode      string  `json:"ProductCode"`
	ProductName      string  `json:"ProductName"`
	ProductType      string  `json:"ProductType"`
	InstanceID       string  `json:"InstanceId"`
	SubscriptionType string  `json:"SubscriptionType"`
	UsageAmount      float64 `json:"UsageAmount"`
	UnitPrice        float64 `json:"UnitPrice"`
	TotalCost        float64 `json:"TotalCost"`
	Currency         string  `json:"Currency"`
	PaymentTime      string  `json:"PaymentTime"`
	Status           string  `json:"Status"`
}

// OrderListResponse 订单列表响应
type OrderListResponse struct {
	Code      string `json:"Code"`
	Message   string `json:"Message"`
	Success   bool   `json:"Success"`
	Data      struct {
		PageNum    int `json:"PageNum"`
	 PageSize   int `json:"PageSize"`
	 TotalCount int `json:"TotalCount"`
	 Orders     []OrderItemRaw `json:"Orders"`
	} `json:"Data"`
}

type OrderItemRaw struct {
	OrderID         string  `json:"OrderId"`
	OrderType       string  `json:"OrderType"`
	ProductCode     string  `json:"ProductCode"`
	ProductName     string  `json:"ProductName"`
	ProductType     string  `json:"ProductType"`
	InstanceID      string  `json:"InstanceId"`
	SubscriptionType string `json:"SubscriptionType"`
	Amount          float64 `json:"Amount"`
	Currency        string  `json:"Currency"`
	Status          string  `json:"Status"`
	PaymentTime     string  `json:"PaymentTime"`
	EffectiveTime   string  `json:"EffectiveTime"`
	ExpireTime      string  `json:"ExpireTime"`
}

// RenewalListResponse 待续费资源响应
type RenewalListResponse struct {
	Code      string `json:"Code"`
	Message   string `json:"Message"`
	Success   bool   `json:"Success"`
	Data      struct {
		PageNum    int `json:"PageNum"`
 PageSize   int `json:"PageSize"`
 TotalCount int `json:"TotalCount"`
 Resources  []RenewalResourceRaw `json:"Resources"`
	} `json:"Data"`
}

type RenewalResourceRaw struct {
	InstanceID    string  `json:"InstanceId"`
	InstanceName  string  `json:"InstanceName"`
	ProductCode   string  `json:"ProductCode"`
	ProductName   string  `json:"ProductName"`
	ProductType   string  `json:"ProductType"`
	RegionId      string  `json:"RegionId"`
	ExpireTime    string  `json:"ExpireTime"`
	RenewalPrice  float64 `json:"RenewalPrice"`
	Currency      string  `json:"Currency"`
	Status        string  `json:"Status"`
}

// BalanceResponse 账户余额响应
type BalanceResponse struct {
	Code      string `json:"Code"`
	Message   string `json:"Message"`
	Success   bool   `json:"Success"`
	Data      struct {
 AvailableAmount float64 `json:"AvailableAmount"`
 CashAmount      float64 `json:"CashAmount"`
 CreditAmount    float64 `json:"CreditAmount"`
 Currency        string  `json:"Currency"`
	} `json:"Data"`
}

// CostAnalysisResponse 成本分析响应
type CostAnalysisResponse struct {
	Code      string `json:"Code"`
	Message   string `json:"Message"`
	Success   bool   `json:"Success"`
	Data      struct {
		PageNum    int `json:"PageNum"`
 PageSize   int `json:"PageSize"`
 TotalCount int `json:"TotalCount"`
 Items      []CostItemRaw `json:"Items"`
	} `json:"Data"`
}

type CostItemRaw struct {
 BillingCycle string  `json:"BillingCycle"`
 BillingDate  string  `json:"BillingDate"`
 ProductCode  string  `json:"ProductCode"`
 ProductName  string  `json:"ProductName"`
 RegionId     string  `json:"RegionId"`
 InstanceId   string  `json:"InstanceId"`
 SubscriptionType string `json:"SubscriptionType"`
 Cost         float64 `json:"Cost"`
 Currency     string  `json:"Currency"`
}

// ========== 解析函数 ==========

func parseBillResponse(resp *responses.CommonResponse) ([]*cloudprovider.BillItem, error) {
	var billResp BillListResponse
	if err := json.Unmarshal(resp.GetHttpContentBytes(), &billResp); err != nil {
		// 解析失败，返回空列表
		return []*cloudprovider.BillItem{}, nil
	}

	if !billResp.Success {
		return []*cloudprovider.BillItem{}, nil
	}

	items := []*cloudprovider.BillItem{}
	for _, raw := range billResp.Data.Items {
		billingMethod := "按量付费"
		if raw.SubscriptionType == "Subscription" {
			billingMethod = "包年包月"
		}

		items = append(items, &cloudprovider.BillItem{
			BillingCycle:  raw.BillingCycle,
			ProductType:   raw.ProductType,
			ProductName:   raw.ProductName,
			InstanceID:    raw.InstanceID,
			UsageAmount:   raw.UsageAmount,
			UnitPrice:     raw.UnitPrice,
			TotalCost:     raw.TotalCost,
			Currency:      raw.Currency,
			BillingMethod: billingMethod,
			Status:        raw.Status,
		})
	}

	return items, nil
}

func parseOrderResponse(resp *responses.CommonResponse) ([]*cloudprovider.OrderItem, error) {
	var orderResp OrderListResponse
	if err := json.Unmarshal(resp.GetHttpContentBytes(), &orderResp); err != nil {
		return []*cloudprovider.OrderItem{}, nil
	}

	if !orderResp.Success {
		return []*cloudprovider.OrderItem{}, nil
	}

	items := []*cloudprovider.OrderItem{}
	for _, raw := range orderResp.Data.Orders {
		orderType := "购买"
		switch raw.OrderType {
		case "Renewal":
			orderType = "续费"
		case "Upgrade":
			orderType = "升级"
		case "Downgrade":
			orderType = "降级"
		case "Refund":
			orderType = "退款"
		}

		items = append(items, &cloudprovider.OrderItem{
			OrderID:       raw.OrderID,
			OrderType:     orderType,
			ProductType:   raw.ProductType,
			ProductName:   raw.ProductName,
			InstanceID:    raw.InstanceID,
			Amount:        raw.Amount,
			Currency:      raw.Currency,
			Status:        raw.Status,
			PaymentTime:   raw.PaymentTime,
			EffectiveTime: raw.EffectiveTime,
			ExpireTime:    raw.ExpireTime,
		})
	}

	return items, nil
}

func parseRenewalResponse(resp *responses.CommonResponse, daysThreshold int) ([]*cloudprovider.RenewalResource, error) {
	var renewalResp RenewalListResponse
	if err := json.Unmarshal(resp.GetHttpContentBytes(), &renewalResp); err != nil {
		return []*cloudprovider.RenewalResource{}, nil
	}

	if !renewalResp.Success {
		return []*cloudprovider.RenewalResource{}, nil
	}

	items := []*cloudprovider.RenewalResource{}
	now := time.Now()

	for _, raw := range renewalResp.Data.Resources {
		// 解析到期时间
		var expireTime time.Time
		if t, err := time.Parse(time.RFC3339, raw.ExpireTime); err == nil {
			expireTime = t
		} else if t, err := time.Parse("2006-01-02T15:04:05Z", raw.ExpireTime); err == nil {
			expireTime = t
		} else {
			expireTime = now.AddDate(0, 0, 30) // 默认30天后
		}

		daysRemaining := int(expireTime.Sub(now).Hours() / 24)

		// 根据阈值筛选
		if daysThreshold > 0 && daysRemaining > daysThreshold {
			continue
		}

		items = append(items, &cloudprovider.RenewalResource{
			InstanceID:    raw.InstanceID,
			InstanceName:  raw.InstanceName,
			ProductType:   raw.ProductType,
			ExpireTime:    expireTime.Format("2006-01-02T15:04:05Z"),
			DaysRemaining: daysRemaining,
			RenewalPrice:  raw.RenewalPrice,
			Status:        raw.Status,
		})
	}

	return items, nil
}

func parseCostAnalysisResponse(resp *responses.CommonResponse) ([]*cloudprovider.CostData, error) {
	var costResp CostAnalysisResponse
	if err := json.Unmarshal(resp.GetHttpContentBytes(), &costResp); err != nil {
		return []*cloudprovider.CostData{}, nil
	}

	if !costResp.Success {
		return []*cloudprovider.CostData{}, nil
	}

	items := []*cloudprovider.CostData{}
	for _, raw := range costResp.Data.Items {
		items = append(items, &cloudprovider.CostData{
			Date:       raw.BillingDate,
			Cost:       raw.Cost,
			Product:    raw.ProductName,
			Region:     raw.RegionId,
			ResourceID: raw.InstanceId,
		})
	}

	return items, nil
}

func parseBalanceResponse(resp *responses.CommonResponse) (*cloudprovider.AccountBalance, error) {
	var balanceResp BalanceResponse
	if err := json.Unmarshal(resp.GetHttpContentBytes(), &balanceResp); err != nil {
		// 返回默认值
		return &cloudprovider.AccountBalance{
			Balance:  0,
			Currency: "CNY",
			Status:   "未知",
		}, nil
	}

	if !balanceResp.Success {
		return &cloudprovider.AccountBalance{
			Balance:  0,
			Currency: "CNY",
			Status:   "查询失败",
		}, nil
	}

	status := "正常"
	if balanceResp.Data.AvailableAmount < 0 {
		status = "欠费"
	} else if balanceResp.Data.AvailableAmount < 100 {
		status = "余额不足"
	}

	return &cloudprovider.AccountBalance{
		Balance:     balanceResp.Data.AvailableAmount,
		Currency:    balanceResp.Data.Currency,
		CreditLimit: balanceResp.Data.CreditAmount,
		Status:      status,
	}, nil
}