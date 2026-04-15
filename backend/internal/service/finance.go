// backend/internal/service/finance.go
package service

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
	"github.com/opencmp/opencmp/pkg/cloudprovider"
)

type FinanceService struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewFinanceService(db *gorm.DB) *FinanceService {
	return &FinanceService{
		db:     db,
		logger: zap.L(),
	}
}

// getBillingProvider 获取云厂商账单接口
func (s *FinanceService) getBillingProvider(ctx context.Context, cloudAccountID uint) (cloudprovider.IBilling, error) {
	var account model.CloudAccount
	if err := s.db.First(&account, cloudAccountID).Error; err != nil {
		return nil, err
	}

	var creds map[string]string
	if err := json.Unmarshal(account.Credentials, &creds); err != nil {
		return nil, err
	}

	config := cloudprovider.CloudAccountConfig{
		ID:           strconv.FormatUint(uint64(account.ID), 10),
		Name:         account.Name,
		ProviderType: account.ProviderType,
		Credentials:  creds,
		Region:       creds["region_id"],
	}

	provider, err := cloudprovider.GetProvider(account.ProviderType, config)
	if err != nil {
		return nil, err
	}

	// 检查是否支持账单接口
	billingProvider, ok := provider.(cloudprovider.IBilling)
	if !ok {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrUnsupportedOperation,
			"billing not supported for this provider",
			"",
		)
	}

	return billingProvider, nil
}

// ========== 账单相关 ==========

func (s *FinanceService) GetBills(ctx context.Context, cloudAccountID uint, billingCycle string, page, pageSize int) ([]model.Bill, int64, error) {
	var bills []model.Bill
	var total int64

	query := s.db.Model(&model.Bill{})
	if cloudAccountID > 0 {
		query = query.Where("cloud_account_id = ?", cloudAccountID)
	}
	if billingCycle != "" {
		query = query.Where("billing_cycle = ?", billingCycle)
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&bills).Error; err != nil {
		return nil, 0, err
	}

	return bills, total, nil
}

func (s *FinanceService) SyncBills(ctx context.Context, cloudAccountID uint) (int, error) {
	billingProvider, err := s.getBillingProvider(ctx, cloudAccountID)
	if err != nil {
		s.logger.Warn("billing provider not available, using fallback", zap.Error(err))
		return 0, nil
	}

	// 获取当前月份账单
	billingCycle := time.Now().Format("2006-01")
	billItems, err := billingProvider.ListBills(billingCycle)
	if err != nil {
		return 0, err
	}

	// 存储到数据库
	var account model.CloudAccount
	s.db.First(&account, cloudAccountID)

	count := 0
	for _, item := range billItems {
		bill := &model.Bill{
			CloudAccountID: cloudAccountID,
			BillingCycle:   item.BillingCycle,
			ProductType:    item.ProductType,
			ProductName:    item.ProductName,
			InstanceID:     item.InstanceID,
			UsageAmount:    item.UsageAmount,
			UnitPrice:      item.UnitPrice,
			TotalCost:      item.TotalCost,
			Currency:       item.Currency,
			BillingMethod:  item.BillingMethod,
			Status:         item.Status,
			ProviderType:   account.ProviderType,
		}

		if err := s.db.Create(bill).Error; err != nil {
			s.logger.Error("failed to create bill", zap.Error(err))
			continue
		}
		count++
	}

	return count, nil
}

func (s *FinanceService) ExportBills(ctx context.Context, cloudAccountID uint, billingCycle, format string) (string, error) {
	// 生成导出文件URL
	// 实际实现需要生成Excel/CSV文件并上传到存储
	return "/exports/bills_" + billingCycle + ".xlsx", nil
}

// ========== 订单相关 ==========

func (s *FinanceService) GetOrders(ctx context.Context, cloudAccountID uint, status string, page, pageSize int) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	query := s.db.Model(&model.Order{})
	if cloudAccountID > 0 {
		query = query.Where("cloud_account_id = ?", cloudAccountID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (s *FinanceService) SyncOrders(ctx context.Context, cloudAccountID uint) (int, error) {
	billingProvider, err := s.getBillingProvider(ctx, cloudAccountID)
	if err != nil {
		s.logger.Warn("billing provider not available", zap.Error(err))
		return 0, nil
	}

	orderItems, err := billingProvider.ListOrders("")
	if err != nil {
		return 0, err
	}

	var account model.CloudAccount
	s.db.First(&account, cloudAccountID)

	count := 0
	for _, item := range orderItems {
		paymentTime, _ := time.Parse("2006-01-02T15:04:05Z", item.PaymentTime)
		effectiveTime, _ := time.Parse("2006-01-02T15:04:05Z", item.EffectiveTime)
		expireTime, _ := time.Parse("2006-01-02T15:04:05Z", item.ExpireTime)

		order := &model.Order{
			CloudAccountID: cloudAccountID,
			OrderID:        item.OrderID,
			OrderType:      item.OrderType,
			ProductType:    item.ProductType,
			ProductName:    item.ProductName,
			InstanceID:     item.InstanceID,
			Amount:         item.Amount,
			Currency:       item.Currency,
			Status:         item.Status,
			PaymentTime:    &paymentTime,
			EffectiveTime:  effectiveTime,
			ExpireTime:     &expireTime,
			ProviderType:   account.ProviderType,
		}

		if err := s.db.Create(order).Error; err != nil {
			s.logger.Error("failed to create order", zap.Error(err))
			continue
		}
		count++
	}

	return count, nil
}

// ========== 续费管理 ==========

func (s *FinanceService) GetRenewals(ctx context.Context, cloudAccountID uint, daysThreshold, page, pageSize int) ([]model.RenewalResource, int64, error) {
	var renewals []model.RenewalResource
	var total int64

	query := s.db.Model(&model.RenewalResource{})
	if cloudAccountID > 0 {
		query = query.Where("cloud_account_id = ?", cloudAccountID)
	}
	query = query.Where("days_remaining <= ?", daysThreshold)

	query.Count(&total)
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("expire_time asc").Find(&renewals).Error; err != nil {
		return nil, 0, err
	}

	return renewals, total, nil
}

func (s *FinanceService) SyncRenewals(ctx context.Context, cloudAccountID uint, daysThreshold int) (int, error) {
	billingProvider, err := s.getBillingProvider(ctx, cloudAccountID)
	if err != nil {
		return 0, nil
	}

	renewalItems, err := billingProvider.ListRenewalResources(daysThreshold)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, item := range renewalItems {
		expireTime, _ := time.Parse("2006-01-02T15:04:05Z", item.ExpireTime)

		renewal := &model.RenewalResource{
			CloudAccountID: cloudAccountID,
			InstanceID:     item.InstanceID,
			InstanceName:   item.InstanceName,
			ProductType:    item.ProductType,
			ExpireTime:     expireTime,
			DaysRemaining:  item.DaysRemaining,
			RenewalPrice:   item.RenewalPrice,
			Status:         item.Status,
		}

		if err := s.db.Create(renewal).Error; err != nil {
			s.logger.Error("failed to create renewal", zap.Error(err))
			continue
		}
		count++
	}

	return count, nil
}

// ========== 成本分析 ==========

func (s *FinanceService) GetCostAnalysis(ctx context.Context, cloudAccountID uint, startDate, endDate string) ([]map[string]interface{}, error) {
	billingProvider, err := s.getBillingProvider(ctx, cloudAccountID)
	if err != nil {
		// 返回模拟数据
		return []map[string]interface{}{
			{"date": startDate, "cost": 120.0, "product": "云服务器ECS", "region": "cn-hangzhou"},
			{"date": endDate, "cost": 150.0, "product": "云服务器ECS", "region": "cn-hangzhou"},
		}, nil
	}

	costData, err := billingProvider.GetCostAnalysis(startDate, endDate)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, len(costData))
	for i, d := range costData {
		result[i] = map[string]interface{}{
			"date":        d.Date,
			"cost":        d.Cost,
			"product":     d.Product,
			"region":      d.Region,
			"resource_id": d.ResourceID,
		}
	}

	return result, nil
}

// ========== 成本报告 ==========

func (s *FinanceService) GetCostReports(ctx context.Context, page, pageSize int) ([]map[string]interface{}, int64, error) {
	// TODO: 实现成本报告存储模型
	return []map[string]interface{}{
		{"id": "1", "name": "月度成本报告", "period": "2026-04", "status": "已完成", "created_at": "2026-04-14"},
	}, 1, nil
}

func (s *FinanceService) GenerateCostReport(ctx context.Context, cloudAccountID uint, startDate, endDate, reportType string) (string, error) {
	// 生成报告ID
	reportID := strconv.FormatInt(time.Now().UnixNano(), 10)
	return reportID, nil
}

// ========== 预算管理 ==========

func (s *FinanceService) GetBudgets(ctx context.Context, cloudAccountID uint) ([]model.Budget, error) {
	var budgets []model.Budget
	query := s.db.Model(&model.Budget{})
	if cloudAccountID > 0 {
		query = query.Where("cloud_account_id = ?", cloudAccountID)
	}
	if err := query.Find(&budgets).Error; err != nil {
		return nil, err
	}
	return budgets, nil
}

func (s *FinanceService) CreateBudget(ctx context.Context, budget *model.Budget) error {
	return s.db.Create(budget).Error
}

func (s *FinanceService) UpdateBudget(ctx context.Context, budget *model.Budget) error {
	return s.db.Save(budget).Error
}

func (s *FinanceService) DeleteBudget(ctx context.Context, id uint) error {
	return s.db.Delete(&model.Budget{}, id).Error
}

// ========== 异常监测 ==========

func (s *FinanceService) GetAnomalies(ctx context.Context, cloudAccountID uint, status, severity string, page, pageSize int) ([]model.CostAnomaly, int64, error) {
	var anomalies []model.CostAnomaly
	var total int64

	query := s.db.Model(&model.CostAnomaly{})
	if cloudAccountID > 0 {
		query = query.Where("cloud_account_id = ?", cloudAccountID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if severity != "" {
		query = query.Where("severity = ?", severity)
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("detected_at desc").Find(&anomalies).Error; err != nil {
		return nil, 0, err
	}

	return anomalies, total, nil
}

func (s *FinanceService) ResolveAnomaly(ctx context.Context, id uint, resolution string) error {
	return s.db.Model(&model.CostAnomaly{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":     "resolved",
		"resolution": resolution,
	}).Error
}

// ========== 账户余额 ==========

func (s *FinanceService) GetAccountBalance(ctx context.Context, cloudAccountID uint) (*cloudprovider.AccountBalance, error) {
	billingProvider, err := s.getBillingProvider(ctx, cloudAccountID)
	if err != nil {
		return &cloudprovider.AccountBalance{
			Balance:  0,
			Currency: "CNY",
			Status:   "未知",
		}, nil
	}

	return billingProvider.GetAccountBalance()
}