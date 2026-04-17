package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/middleware"
	"github.com/opencmp/opencmp/internal/service"
	"github.com/opencmp/opencmp/pkg/cloudprovider"
)

// ComputeHandler 计算资源 Handler
type ComputeHandler struct {
	service *service.ComputeService
	logger  *zap.Logger
}

// NewComputeHandler 创建计算资源 Handler
func NewComputeHandler(db *gorm.DB, logger *zap.Logger) *ComputeHandler {
	return &ComputeHandler{
		service: service.NewComputeService(db),
		logger:  logger,
	}
}

// CreateVMRequest 创建虚拟机请求
type CreateVMRequest struct {
	AccountID      uint              `json:"account_id" binding:"required"`
	Name           string            `json:"name" binding:"required"`
	InstanceType   string            `json:"instance_type" binding:"required"`
	ImageID        string            `json:"image_id" binding:"required"`
	VPCID          string            `json:"vpc_id" binding:"required"`
	SubnetID       string            `json:"subnet_id" binding:"required"`
	SecurityGroups []string          `json:"security_groups"`
	DiskSize       int               `json:"disk_size"`
	Keypair        string            `json:"keypair"`
	Tags           map[string]string `json:"tags"`
}

// VMActionRequest 虚拟机操作请求
type VMActionRequest struct {
	Action string `json:"action" binding:"required"` // start/stop/reboot
}

// CreateVM 创建虚拟机
func (h *ComputeHandler) CreateVM(c *gin.Context) {
	var req CreateVMRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config := cloudprovider.VMCreateConfig{
		Name:           req.Name,
		InstanceType:   req.InstanceType,
		ImageID:        req.ImageID,
		VPCID:          req.VPCID,
		SubnetID:       req.SubnetID,
		SecurityGroups: req.SecurityGroups,
		DiskSize:       req.DiskSize,
		Keypair:        req.Keypair,
		Tags:           req.Tags,
	}

	vm, err := h.service.CreateVM(c.Request.Context(), req.AccountID, config)
	if err != nil {
		h.logger.Error("failed to create vm", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, vm)
}

// ListVMs 列出虚拟机
func (h *ComputeHandler) ListVMs(c *gin.Context) {
	// 获取项目过滤信息
	projectFilter := middleware.GetProjectFilter(c)
	var projectIDs []int64
	if !projectFilter.AllProjectsVisible {
		projectIDs = projectFilter.ProjectIDs
	}

	filter := cloudprovider.VMListFilter{
		VPCID:    c.Query("vpc_id"),
		SubnetID: c.Query("subnet_id"),
		RegionID: c.Query("region_id"),
	}

	// account_id 可选，如果不传则查询所有云账号的虚拟机
	accountIDStr := c.Query("account_id")
	if accountIDStr == "" {
		// 查询所有云账号的虚拟机
		vms, err := h.service.ListAllVMs(c.Request.Context(), filter, projectIDs)
		if err != nil {
			h.logger.Error("failed to list all vms", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"items": vms,
			"total": len(vms),
		})
		return
	}

	// 查询指定云账号的虚拟机
	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	vms, err := h.service.ListVMs(c.Request.Context(), uint(accountID), filter, projectIDs)
	if err != nil {
		h.logger.Error("failed to list vms", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": vms,
		"total": len(vms),
	})
}

// GetVM 获取虚拟机详情
func (h *ComputeHandler) GetVM(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	vmID := c.Param("id")

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	vm, err := h.service.GetVM(c.Request.Context(), uint(accountID), vmID)
	if err != nil {
		h.logger.Error("failed to get vm", zap.Error(err))
		if err.Error() == "vm not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "vm not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, vm)
}

// DeleteVM 删除虚拟机
func (h *ComputeHandler) DeleteVM(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	vmID := c.Param("id")

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	if err := h.service.DeleteVM(c.Request.Context(), uint(accountID), vmID); err != nil {
		h.logger.Error("failed to delete vm", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "vm deleted"})
}

// VMAction 虚拟机操作
func (h *ComputeHandler) VMAction(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	vmID := c.Param("id")

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	var req VMActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.VMAction(c.Request.Context(), uint(accountID), vmID, req.Action); err != nil {
		h.logger.Error("failed to perform vm action", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "action performed", "action": req.Action})
}

// GetVMDetails 获取虚拟机详细信息
func (h *ComputeHandler) GetVMDetails(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	vmID := c.Param("id")

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	vmDetails, err := h.service.GetVMDetails(c.Request.Context(), uint(accountID), vmID)
	if err != nil {
		h.logger.Error("failed to get vm details", zap.Error(err))
		if err.Error() == "vm not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "vm not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, vmDetails)
}

// GetVMSecurityGroups 获取虚拟机关联的安全组
func (h *ComputeHandler) GetVMSecurityGroups(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	vmID := c.Param("id")

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	securityGroups, err := h.service.GetVMSecurityGroups(c.Request.Context(), uint(accountID), vmID)
	if err != nil {
		h.logger.Error("failed to get vm security groups", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": securityGroups,
		"total": len(securityGroups),
	})
}

// GetVMNetworkInfo 获取虚拟机网络信息
func (h *ComputeHandler) GetVMNetworkInfo(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	vmID := c.Param("id")

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	networkInfo, err := h.service.GetVMNetworkInfo(c.Request.Context(), uint(accountID), vmID)
	if err != nil {
		h.logger.Error("failed to get vm network info", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, networkInfo)
}

// GetVMDisks 获取虚拟机关联的磁盘
func (h *ComputeHandler) GetVMDisks(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	vmID := c.Param("id")

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	disks, err := h.service.GetVMDisks(c.Request.Context(), uint(accountID), vmID)
	if err != nil {
		h.logger.Error("failed to get vm disks", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": disks,
		"total": len(disks),
	})
}

// GetVMSnapshots 获取虚拟机相关的快照
func (h *ComputeHandler) GetVMSnapshots(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	vmID := c.Param("id")

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	snapshots, err := h.service.GetVMSnapshots(c.Request.Context(), uint(accountID), vmID)
	if err != nil {
		h.logger.Error("failed to get vm snapshots", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": snapshots,
		"total": len(snapshots),
	})
}

// GetVMOperationLogs 获取虚拟机操作日志
func (h *ComputeHandler) GetVMOperationLogs(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	vmID := c.Param("id")

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	logs, err := h.service.GetVMOperationLogs(c.Request.Context(), uint(accountID), vmID)
	if err != nil {
		h.logger.Error("failed to get vm operation logs", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": logs,
		"total": len(logs),
	})
}

// GetVNCInfo 获取虚拟机VNC连接信息
func (h *ComputeHandler) GetVNCInfo(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	vmID := c.Param("id")

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	vncInfo, err := h.service.GetVNCInfo(c.Request.Context(), uint(accountID), vmID)
	if err != nil {
		h.logger.Error("failed to get vm vnc info", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, vncInfo)
}

// ResetPassword 重置虚拟机密码
func (h *ComputeHandler) ResetPassword(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	vmID := c.Param("id")

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	var req struct {
		Username    string `json:"username"`
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.ResetPassword(c.Request.Context(), uint(accountID), vmID, req.Username, req.NewPassword)
	if err != nil {
		h.logger.Error("failed to reset vm password", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password reset successfully"})
}

// UpdateVMConfig 更新虚拟机配置
func (h *ComputeHandler) UpdateVMConfig(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	vmID := c.Param("id")

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	var req struct {
		InstanceType string `json:"instance_type"`
		Name         string `json:"name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.UpdateVMConfig(c.Request.Context(), uint(accountID), vmID, req.InstanceType, req.Name)
	if err != nil {
		h.logger.Error("failed to update vm config", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "config updated successfully"})
}

// ListImages 列出镜像
func (h *ComputeHandler) ListImages(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	if accountIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account_id is required"})
		return
	}

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	filter := cloudprovider.ImageFilter{
		Platform:   c.Query("platform"),
		MaxResults: 100,
	}

	images, err := h.service.ListImages(c.Request.Context(), uint(accountID), filter)
	if err != nil {
		h.logger.Error("failed to list images", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": images,
		"total": len(images),
	})
}
