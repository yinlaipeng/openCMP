package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/service"
	"github.com/opencmp/opencmp/pkg/cloudprovider"
)

// NetworkHandler 网络资源 Handler
type NetworkHandler struct {
	service *service.NetworkService
	logger  *zap.Logger
}

// NewNetworkHandler 创建网络资源 Handler
func NewNetworkHandler(db *gorm.DB, logger *zap.Logger) *NetworkHandler {
	return &NetworkHandler{
		service: service.NewNetworkService(db),
		logger:  logger,
	}
}

// CreateVPCRequest 创建 VPC 请求
type CreateVPCRequest struct {
	AccountID   uint              `json:"account_id" binding:"required"`
	Name        string            `json:"name" binding:"required"`
	CIDR        string            `json:"cidr" binding:"required"`
	Description string            `json:"description"`
	Tags        map[string]string `json:"tags"`
}

// CreateVPC 创建 VPC
func (h *NetworkHandler) CreateVPC(c *gin.Context) {
	var req CreateVPCRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config := cloudprovider.VPCConfig{
		Name:        req.Name,
		CIDR:        req.CIDR,
		Description: req.Description,
		Tags:        req.Tags,
	}

	vpc, err := h.service.CreateVPC(c.Request.Context(), req.AccountID, config)
	if err != nil {
		h.logger.Error("failed to create vpc", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, vpc)
}

// ListVPCs 列出 VPC
func (h *NetworkHandler) ListVPCs(c *gin.Context) {
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

	filter := cloudprovider.VPCFilter{
		RegionID: c.Query("region_id"),
	}

	vpcs, err := h.service.ListVPCs(c.Request.Context(), uint(accountID), filter)
	if err != nil {
		h.logger.Error("failed to list vpcs", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": vpcs,
		"total": len(vpcs),
	})
}

// DeleteVPC 删除 VPC
func (h *NetworkHandler) DeleteVPC(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	vpcID := c.Param("id")

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	if err := h.service.DeleteVPC(c.Request.Context(), uint(accountID), vpcID); err != nil {
		h.logger.Error("failed to delete vpc", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "vpc deleted"})
}

// CreateSubnetRequest 创建子网请求
type CreateSubnetRequest struct {
	AccountID   uint              `json:"account_id" binding:"required"`
	Name        string            `json:"name" binding:"required"`
	VPCID       string            `json:"vpc_id" binding:"required"`
	CIDR        string            `json:"cidr" binding:"required"`
	ZoneID      string            `json:"zone_id" binding:"required"`
	Description string            `json:"description"`
	Tags        map[string]string `json:"tags"`
}

// CreateSubnet 创建子网
func (h *NetworkHandler) CreateSubnet(c *gin.Context) {
	var req CreateSubnetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config := cloudprovider.SubnetConfig{
		Name:   req.Name,
		VPCID:  req.VPCID,
		CIDR:   req.CIDR,
		ZoneID: req.ZoneID,
		Tags:   req.Tags,
	}

	subnet, err := h.service.CreateSubnet(c.Request.Context(), req.AccountID, config)
	if err != nil {
		h.logger.Error("failed to create subnet", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, subnet)
}

// ListSubnets 列出子网
func (h *NetworkHandler) ListSubnets(c *gin.Context) {
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

	filter := cloudprovider.SubnetFilter{
		VPCID:  c.Query("vpc_id"),
		ZoneID: c.Query("zone_id"),
	}

	subnets, err := h.service.ListSubnets(c.Request.Context(), uint(accountID), filter)
	if err != nil {
		h.logger.Error("failed to list subnets", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": subnets,
		"total": len(subnets),
	})
}

// CreateSecurityGroupRequest 创建安全组请求
type CreateSecurityGroupRequest struct {
	AccountID   uint              `json:"account_id" binding:"required"`
	Name        string            `json:"name" binding:"required"`
	Description string            `json:"description"`
	VPCID       string            `json:"vpc_id" binding:"required"`
	Tags        map[string]string `json:"tags"`
}

// CreateSecurityGroup 创建安全组
func (h *NetworkHandler) CreateSecurityGroup(c *gin.Context) {
	var req CreateSecurityGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config := cloudprovider.SGConfig{
		Name:        req.Name,
		Description: req.Description,
		VPCID:       req.VPCID,
		Tags:        req.Tags,
	}

	sg, err := h.service.CreateSecurityGroup(c.Request.Context(), req.AccountID, config)
	if err != nil {
		h.logger.Error("failed to create security group", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, sg)
}

// ListSecurityGroups 列出安全组
func (h *NetworkHandler) ListSecurityGroups(c *gin.Context) {
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

	filter := cloudprovider.SGFilter{
		VPCID: c.Query("vpc_id"),
	}

	sgs, err := h.service.ListSecurityGroups(c.Request.Context(), uint(accountID), filter)
	if err != nil {
		h.logger.Error("failed to list security groups", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": sgs,
		"total": len(sgs),
	})
}

// CreateEIPRequest 创建弹性 IP 请求
type CreateEIPRequest struct {
	AccountID uint              `json:"account_id" binding:"required"`
	Bandwidth int               `json:"bandwidth" binding:"required"`
	RegionID  string            `json:"region_id" binding:"required"`
	Tags      map[string]string `json:"tags"`
}

// CreateEIP 创建弹性 IP
func (h *NetworkHandler) CreateEIP(c *gin.Context) {
	var req CreateEIPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config := cloudprovider.EIPConfig{
		Bandwidth: req.Bandwidth,
		RegionID:  req.RegionID,
		Tags:      req.Tags,
	}

	eip, err := h.service.CreateEIP(c.Request.Context(), req.AccountID, config)
	if err != nil {
		h.logger.Error("failed to create eip", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, eip)
}

// ListEIPs 列出弹性 IP
func (h *NetworkHandler) ListEIPs(c *gin.Context) {
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

	filter := cloudprovider.EIPFilter{
		RegionID: c.Query("region_id"),
		Status:   c.Query("status"),
	}

	eips, err := h.service.ListEIPs(c.Request.Context(), uint(accountID), filter)
	if err != nil {
		h.logger.Error("failed to list eips", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": eips,
		"total": len(eips),
	})
}
