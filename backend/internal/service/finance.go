// backend/internal/service/finance.go
package service

import (
	"context"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
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
	// TODO: 实现云厂商 API 调用同步账单
	// 当前返回模拟数据
	return 0, nil
}

func (s *FinanceService) ExportBills(ctx context.Context, cloudAccountID uint, billingCycle, format string) (string, error) {
	// TODO: 实现账单导出功能
	return "", nil
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
	// TODO: 实现云厂商 API 调用同步订单
	return 0, nil
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

// ========== 成本分析 ==========

func (s *FinanceService) GetCostAnalysis(ctx context.Context, cloudAccountID uint, startDate, endDate string) ([]map[string]interface{}, error) {
	// TODO: 实现成本分析聚合查询
	return []map[string]interface{}{}, nil
}

// ========== 成本报告 ==========

func (s *FinanceService) GetCostReports(ctx context.Context, page, pageSize int) ([]map[string]interface{}, int64, error) {
	// TODO: 实现成本报告列表
	return []map[string]interface{}{}, 0, nil
}

func (s *FinanceService) GenerateCostReport(ctx context.Context, cloudAccountID uint, startDate, endDate, reportType string) (string, error) {
	// TODO: 实现成本报告生成
	return "", nil
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