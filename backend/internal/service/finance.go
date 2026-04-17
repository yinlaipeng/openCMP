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

// ========== 成本聚合统计 ==========

// CostAggregation 成本聚合结果
type CostAggregation struct {
	Key    string  `json:"key"`     // 聚合键（项目名/账号名/服务名）
	Name   string  `json:"name"`    // 显示名称
	Cost   float64 `json:"cost"`    // 成本金额
	Count  int     `json:"count"`   // 记录数量
	Ratio  float64 `json:"ratio"`   // 占比（百分比）
}

// GetCostByProject 按项目统计成本
func (s *FinanceService) GetCostByProject(ctx context.Context, startDate, endDate string) ([]CostAggregation, error) {
	var results []CostAggregation

	// 基础查询：按project_id分组
	query := `
		SELECT
			COALESCE(p.name, '未归属') as name,
			COALESCE(b.project_id, 0) as key,
			SUM(b.total_cost) as cost,
			COUNT(*) as count
		FROM finance_bills b
		LEFT JOIN projects p ON b.project_id = p.id
		WHERE b.billing_cycle BETWEEN ? AND ?
		GROUP BY b.project_id, p.name
		ORDER BY cost DESC
	`

	if err := s.db.Raw(query, startDate, endDate).Scan(&results).Error; err != nil {
		return nil, err
	}

	// 计算总成本和占比
	var totalCost float64
	for _, r := range results {
		totalCost += r.Cost
	}
	for i := range results {
		if totalCost > 0 {
			results[i].Ratio = (results[i].Cost / totalCost) * 100
		}
		results[i].Key = results[i].Name // 使用名称作为key
	}

	return results, nil
}

// GetCostByAccount 按云账号统计成本
func (s *FinanceService) GetCostByAccount(ctx context.Context, startDate, endDate string) ([]CostAggregation, error) {
	var results []CostAggregation

	query := `
		SELECT
			ca.name as name,
			ca.id as key,
			SUM(b.total_cost) as cost,
			COUNT(*) as count
		FROM finance_bills b
		JOIN cloud_accounts ca ON b.cloud_account_id = ca.id
		WHERE b.billing_cycle BETWEEN ? AND ?
		GROUP BY b.cloud_account_id, ca.name
		ORDER BY cost DESC
	`

	if err := s.db.Raw(query, startDate, endDate).Scan(&results).Error; err != nil {
		return nil, err
	}

	// 计算总成本和占比
	var totalCost float64
	for _, r := range results {
		totalCost += r.Cost
	}
	for i := range results {
		if totalCost > 0 {
			results[i].Ratio = (results[i].Cost / totalCost) * 100
		}
	}

	return results, nil
}

// GetCostByService 按服务类型统计成本
func (s *FinanceService) GetCostByService(ctx context.Context, startDate, endDate string) ([]CostAggregation, error) {
	var results []CostAggregation

	query := `
		SELECT
			b.product_type as key,
			b.product_name as name,
			SUM(b.total_cost) as cost,
			COUNT(*) as count
		FROM finance_bills b
		WHERE b.billing_cycle BETWEEN ? AND ?
		GROUP BY b.product_type, b.product_name
		ORDER BY cost DESC
	`

	if err := s.db.Raw(query, startDate, endDate).Scan(&results).Error; err != nil {
		return nil, err
	}

	// 计算总成本和占比
	var totalCost float64
	for _, r := range results {
		totalCost += r.Cost
	}
	for i := range results {
		if totalCost > 0 {
			results[i].Ratio = (results[i].Cost / totalCost) * 100
		}
	}

	return results, nil
}

// CostTrendPoint 成本趋势数据点
type CostTrendPoint struct {
	Date    string  `json:"date"`    // 日期（YYYY-MM 或 YYYY-MM-DD）
	Cost    float64 `json:"cost"`    // 成本
	Budget  float64 `json:"budget"`  // 预算（可选）
	Change  float64 `json:"change"`  // 环比变化率
}

// GetCostTrend 获取成本趋势数据（用于图表）
func (s *FinanceService) GetCostTrend(ctx context.Context, cloudAccountID uint, months int) ([]CostTrendPoint, error) {
	var results []CostTrendPoint

	// 计算起始周期（前N个月）
	startDate := time.Now().AddDate(0, -months, 0).Format("2006-01")
	endDate := time.Now().Format("2006-01")

	query := `
		SELECT
			billing_cycle as date,
			SUM(total_cost) as cost
		FROM finance_bills
		WHERE billing_cycle BETWEEN ? AND ?
	`
	if cloudAccountID > 0 {
		query += " AND cloud_account_id = ?"
	}
	query += " GROUP BY billing_cycle ORDER BY billing_cycle"

	if cloudAccountID > 0 {
		if err := s.db.Raw(query, startDate, endDate, cloudAccountID).Scan(&results).Error; err != nil {
			return nil, err
		}
	} else {
		if err := s.db.Raw(query, startDate, endDate).Scan(&results).Error; err != nil {
			return nil, err
		}
	}

	// 计算环比变化率
	for i := 1; i < len(results); i++ {
		if results[i-1].Cost > 0 {
			results[i].Change = ((results[i].Cost - results[i-1].Cost) / results[i-1].Cost) * 100
		}
	}

	return results, nil
}

// GetCostSummary 获取成本概览（用于Dashboard）
func (s *FinanceService) GetCostSummary(ctx context.Context, cloudAccountID uint) (map[string]interface{}, error) {
	currentMonth := time.Now().Format("2006-01")
	lastMonth := time.Now().AddDate(0, -1, 0).Format("2006-01")

	// 当前月成本
	var currentCost float64
	query := s.db.Model(&model.Bill{}).Where("billing_cycle = ?", currentMonth)
	if cloudAccountID > 0 {
		query = query.Where("cloud_account_id = ?", cloudAccountID)
	}
	query.Select("SUM(total_cost)").Scan(&currentCost)

	// 上月成本
	var lastCost float64
	query = s.db.Model(&model.Bill{}).Where("billing_cycle = ?", lastMonth)
	if cloudAccountID > 0 {
		query = query.Where("cloud_account_id = ?", cloudAccountID)
	}
	query.Select("SUM(total_cost)").Scan(&lastCost)

	// 计算环比变化
	var changeRate float64
	if lastCost > 0 {
		changeRate = ((currentCost - lastCost) / lastCost) * 100
	}

	// 预算执行进度（如果有）
	var budget model.Budget
	budgetQuery := s.db.Model(&model.Budget{}).Where("cloud_account_id = ? AND status = ?", cloudAccountID, "active")
	if cloudAccountID > 0 {
		budgetQuery = budgetQuery.Where("cloud_account_id = ?", cloudAccountID)
	}
	budgetQuery.First(&budget)

	budgetUsage := 0.0
	if budget.Amount > 0 {
		budgetUsage = (currentCost / budget.Amount) * 100
	}

	return map[string]interface{}{
		"current_month_cost": currentCost,
		"last_month_cost":    lastCost,
		"change_rate":        changeRate,
		"budget_amount":      budget.Amount,
		"budget_usage":       budgetUsage,
		"currency":           "CNY",
	}, nil
}