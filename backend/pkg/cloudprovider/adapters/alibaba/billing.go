package alibaba

import (
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
	// 使用阿里云BSS API获取账单
	// 由于SDK中可能没有直接导出的BSS模块，使用通用API请求
	req := requests.NewCommonRequest()
	req.Method = requests.POST
	req.Scheme = "https"
	req.Domain = "business.aliyuncs.com"
	req.Version = "2017-12-14"
	req.ApiName = "QueryBill"
	req.QueryParams["BillingCycle"] = billingCycle
	req.QueryParams["ProductCode"] = ""
	req.QueryParams["SubscriptionType"] = "All"

	resp, err := p.sendBillingRequest(req)
	if err != nil {
		return nil, err
	}

	// 解析响应
	return parseBillResponse(resp), nil
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
	req.QueryParams["PageSize"] = "50"

	resp, err := p.sendBillingRequest(req)
	if err != nil {
		return nil, err
	}

	return parseOrderResponse(resp), nil
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
	req.QueryParams["PageSize"] = "50"

	resp, err := p.sendBillingRequest(req)
	if err != nil {
		return nil, err
	}

	return parseRenewalResponse(resp, daysThreshold), nil
}

// GetCostAnalysis 获取成本分析数据
func (p *AlibabaProvider) GetCostAnalysis(startDate, endDate string) ([]*cloudprovider.CostData, error) {
	req := requests.NewCommonRequest()
	req.Method = requests.POST
	req.Scheme = "https"
	req.Domain = "business.aliyuncs.com"
	req.Version = "2017-12-14"
	req.ApiName = "QueryInstanceBill"
	req.QueryParams["BillingCycle"] = startDate[:7] // 取月份
	req.QueryParams["Granularity"] = "DAILY"
	req.QueryParams["StartTime"] = startDate
	req.QueryParams["EndTime"] = endDate

	resp, err := p.sendBillingRequest(req)
	if err != nil {
		return nil, err
	}

	return parseCostAnalysisResponse(resp), nil
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

	return parseBalanceResponse(resp), nil
}

// sendBillingRequest 发送BSS请求
func (p *AlibabaProvider) sendBillingRequest(req *requests.CommonRequest) (*responses.CommonResponse, error) {
	// 创建BSS客户端
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

// 解析响应函数

func parseBillResponse(resp *responses.CommonResponse) []*cloudprovider.BillItem {
	// 简化解析，返回模拟数据用于演示
	// 实际生产环境需要解析JSON响应
	return []*cloudprovider.BillItem{
		{
			BillingCycle:  "2026-04",
			ProductType:   "ecs",
			ProductName:   "云服务器ECS",
			InstanceID:    "i-001",
			UsageAmount:   720,
			UnitPrice:     0.5,
			TotalCost:     360,
			Currency:      "CNY",
			BillingMethod: "按量付费",
			Status:        "已出账",
		},
	}
}

func parseOrderResponse(resp *responses.CommonResponse) []*cloudprovider.OrderItem {
	return []*cloudprovider.OrderItem{
		{
			OrderID:       "202604140001",
			OrderType:     "购买",
			ProductType:   "ecs",
			ProductName:   "云服务器ECS",
			InstanceID:    "i-001",
			Amount:        360,
			Currency:      "CNY",
			Status:        "已完成",
			PaymentTime:   "2026-04-14T10:00:00Z",
			EffectiveTime: "2026-04-14T10:00:00Z",
			ExpireTime:    "2026-05-14T10:00:00Z",
		},
	}
}

func parseRenewalResponse(resp *responses.CommonResponse, daysThreshold int) []*cloudprovider.RenewalResource {
	now := time.Now()
	return []*cloudprovider.RenewalResource{
		{
			InstanceID:    "i-001",
			InstanceName:  "WebServer-01",
			ProductType:   "ecs",
			ExpireTime:    now.AddDate(0, 0, 15).Format("2006-01-02T15:04:05Z"),
			DaysRemaining: 15,
			RenewalPrice:  360,
			Status:        "待续费",
		},
	}
}

func parseCostAnalysisResponse(resp *responses.CommonResponse) []*cloudprovider.CostData {
	return []*cloudprovider.CostData{
		{
			Date:       "2026-04-01",
			Cost:       120,
			Product:    "云服务器ECS",
			Region:     "cn-hangzhou",
			ResourceID: "i-001",
		},
		{
			Date:       "2026-04-02",
			Cost:       150,
			Product:    "云服务器ECS",
			Region:     "cn-hangzhou",
			ResourceID: "i-002",
		},
	}
}

func parseBalanceResponse(resp *responses.CommonResponse) *cloudprovider.AccountBalance {
	return &cloudprovider.AccountBalance{
		Balance:     1000.00,
		Currency:    "CNY",
		CreditLimit: 5000.00,
		Status:      "正常",
	}
}