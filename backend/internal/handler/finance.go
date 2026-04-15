// backend/internal/handler/finance.go
package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
	"github.com/opencmp/opencmp/internal/service"
)

type FinanceHandler struct {
	service *service.FinanceService
	logger  *zap.Logger
}

func NewFinanceHandler(db *gorm.DB, logger *zap.Logger) *FinanceHandler {
	return &FinanceHandler{
		service: service.NewFinanceService(db),
		logger:  logger,
	}
}

// ========== 账单相关 ==========

// GetBills 获取账单列表
func (h *FinanceHandler) GetBills(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	cloudAccountID, _ := strconv.ParseUint(c.Query("cloud_account_id"), 10, 32)
	billingCycle := c.Query("billing_cycle")

	bills, total, err := h.service.GetBills(c.Request.Context(), uint(cloudAccountID), billingCycle, page, pageSize)
	if err != nil {
		h.logger.Error("failed to get bills", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     bills,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// SyncBills 同步账单数据
func (h *FinanceHandler) SyncBills(c *gin.Context) {
	var req struct {
		CloudAccountID uint `json:"cloud_account_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	count, err := h.service.SyncBills(c.Request.Context(), req.CloudAccountID)
	if err != nil {
		h.logger.Error("failed to sync bills", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "sync completed",
		"count":   count,
	})
}

// ExportBills 导出账单
func (h *FinanceHandler) ExportBills(c *gin.Context) {
	var req struct {
		CloudAccountID uint   `json:"cloud_account_id"`
		BillingCycle   string `json:"billing_cycle"`
		Format         string `json:"format"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	url, err := h.service.ExportBills(c.Request.Context(), req.CloudAccountID, req.BillingCycle, req.Format)
	if err != nil {
		h.logger.Error("failed to export bills", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"download_url": url})
}

// ========== 订单相关 ==========

// GetOrders 获取订单列表
func (h *FinanceHandler) GetOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	cloudAccountID, _ := strconv.ParseUint(c.Query("cloud_account_id"), 10, 32)
	status := c.Query("status")

	orders, total, err := h.service.GetOrders(c.Request.Context(), uint(cloudAccountID), status, page, pageSize)
	if err != nil {
		h.logger.Error("failed to get orders", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     orders,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// SyncOrders 同步订单数据
func (h *FinanceHandler) SyncOrders(c *gin.Context) {
	var req struct {
		CloudAccountID uint `json:"cloud_account_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	count, err := h.service.SyncOrders(c.Request.Context(), req.CloudAccountID)
	if err != nil {
		h.logger.Error("failed to sync orders", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "sync completed",
		"count":   count,
	})
}

// ========== 续费管理 ==========

// GetRenewals 获取待续费资源列表
func (h *FinanceHandler) GetRenewals(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	cloudAccountID, _ := strconv.ParseUint(c.Query("cloud_account_id"), 10, 32)
	daysThreshold, _ := strconv.Atoi(c.DefaultQuery("days_threshold", "30"))

	renewals, total, err := h.service.GetRenewals(c.Request.Context(), uint(cloudAccountID), daysThreshold, page, pageSize)
	if err != nil {
		h.logger.Error("failed to get renewals", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     renewals,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// ========== 成本分析 ==========

// GetCostAnalysis 获取成本分析数据
func (h *FinanceHandler) GetCostAnalysis(c *gin.Context) {
	cloudAccountID, _ := strconv.ParseUint(c.Query("cloud_account_id"), 10, 32)
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	data, err := h.service.GetCostAnalysis(c.Request.Context(), uint(cloudAccountID), startDate, endDate)
	if err != nil {
		h.logger.Error("failed to get cost analysis", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

// ========== 成本报告 ==========

// GetCostReports 获取成本报告列表
func (h *FinanceHandler) GetCostReports(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	reports, total, err := h.service.GetCostReports(c.Request.Context(), page, pageSize)
	if err != nil {
		h.logger.Error("failed to get cost reports", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     reports,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GenerateCostReport 生成成本报告
func (h *FinanceHandler) GenerateCostReport(c *gin.Context) {
	var req struct {
		CloudAccountID uint   `json:"cloud_account_id"`
		StartDate      string `json:"start_date"`
		EndDate        string `json:"end_date"`
		ReportType     string `json:"report_type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reportID, err := h.service.GenerateCostReport(c.Request.Context(), req.CloudAccountID, req.StartDate, req.EndDate, req.ReportType)
	if err != nil {
		h.logger.Error("failed to generate cost report", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"report_id": reportID,
		"message":   "report generated",
	})
}

// ========== 预算管理 ==========

// GetBudgets 获取预算列表
func (h *FinanceHandler) GetBudgets(c *gin.Context) {
	cloudAccountID, _ := strconv.ParseUint(c.Query("cloud_account_id"), 10, 32)

	budgets, err := h.service.GetBudgets(c.Request.Context(), uint(cloudAccountID))
	if err != nil {
		h.logger.Error("failed to get budgets", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, budgets)
}

// CreateBudget 创建预算
func (h *FinanceHandler) CreateBudget(c *gin.Context) {
	var budget model.Budget
	if err := c.ShouldBindJSON(&budget); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateBudget(c.Request.Context(), &budget); err != nil {
		h.logger.Error("failed to create budget", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, budget)
}

// UpdateBudget 更新预算
func (h *FinanceHandler) UpdateBudget(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var budget model.Budget
	if err := c.ShouldBindJSON(&budget); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	budget.ID = uint(id)

	if err := h.service.UpdateBudget(c.Request.Context(), &budget); err != nil {
		h.logger.Error("failed to update budget", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, budget)
}

// DeleteBudget 删除预算
func (h *FinanceHandler) DeleteBudget(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if err := h.service.DeleteBudget(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("failed to delete budget", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// ========== 异常监测 ==========

// GetAnomalies 获取异常列表
func (h *FinanceHandler) GetAnomalies(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	cloudAccountID, _ := strconv.ParseUint(c.Query("cloud_account_id"), 10, 32)
	status := c.Query("status")
	severity := c.Query("severity")

	anomalies, total, err := h.service.GetAnomalies(c.Request.Context(), uint(cloudAccountID), status, severity, page, pageSize)
	if err != nil {
		h.logger.Error("failed to get anomalies", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     anomalies,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// ResolveAnomaly 处理异常
func (h *FinanceHandler) ResolveAnomaly(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req struct {
		Resolution string `json:"resolution" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.ResolveAnomaly(c.Request.Context(), uint(id), req.Resolution); err != nil {
		h.logger.Error("failed to resolve anomaly", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "resolved"})
}

// ========== 续费同步 ==========

// SyncRenewals 同步待续费资源
func (h *FinanceHandler) SyncRenewals(c *gin.Context) {
	var req struct {
		CloudAccountID uint `json:"cloud_account_id" binding:"required"`
		DaysThreshold  int  `json:"days_threshold"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.DaysThreshold == 0 {
		req.DaysThreshold = 30
	}

	count, err := h.service.SyncRenewals(c.Request.Context(), req.CloudAccountID, req.DaysThreshold)
	if err != nil {
		h.logger.Error("failed to sync renewals", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "sync completed",
		"count":   count,
	})
}

// ========== 账户余额 ==========

// GetAccountBalance 获取账户余额
func (h *FinanceHandler) GetAccountBalance(c *gin.Context) {
	cloudAccountID, _ := strconv.ParseUint(c.Query("cloud_account_id"), 10, 32)

	balance, err := h.service.GetAccountBalance(c.Request.Context(), uint(cloudAccountID))
	if err != nil {
		h.logger.Error("failed to get account balance", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, balance)
}