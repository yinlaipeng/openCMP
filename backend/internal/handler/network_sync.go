package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
)

// NetworkSyncHandler 网络同步资源 Handler
type NetworkSyncHandler struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewNetworkSyncHandler 创建网络同步资源 Handler
func NewNetworkSyncHandler(db *gorm.DB, logger *zap.Logger) *NetworkSyncHandler {
	return &NetworkSyncHandler{db: db, logger: logger}
}

// ========== SecurityGroup Handlers ==========

// ListSecurityGroups 列出安全组
func (h *NetworkSyncHandler) ListSecurityGroups(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	cloudAccountID, _ := strconv.ParseUint(c.Query("cloud_account_id"), 10, 32)
	status := c.Query("status")
	name := c.Query("name")
	region := c.Query("region")

	var securityGroups []model.CloudSecurityGroup
	var total int64

	query := h.db.Model(&model.CloudSecurityGroup{})
	if cloudAccountID > 0 {
		query = query.Where("cloud_account_id = ?", cloudAccountID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if region != "" {
		query = query.Where("region LIKE ?", "%"+region+"%")
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&securityGroups).Error; err != nil {
		h.logger.Error("failed to list security groups", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     securityGroups,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetSecurityGroup 获取安全组详情
func (h *NetworkSyncHandler) GetSecurityGroup(c *gin.Context) {
	id := c.Param("id")

	var securityGroup model.CloudSecurityGroup
	if err := h.db.First(&securityGroup, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "security group not found"})
		return
	}

	c.JSON(http.StatusOK, securityGroup)
}

// CreateSecurityGroup 创建安全组
func (h *NetworkSyncHandler) CreateSecurityGroup(c *gin.Context) {
	var req struct {
		CloudAccountID uint   `json:"cloud_account_id" binding:"required"`
		Name           string `json:"name" binding:"required"`
		Description    string `json:"description"`
		RegionID       string `json:"region_id" binding:"required"`
		VPCID          string `json:"vpc_id"`
		ProjectID      uint   `json:"project_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	securityGroup := &model.CloudSecurityGroup{
		CloudAccountID:  req.CloudAccountID,
		SecurityGroupID: "sg-" + strconv.FormatInt(time.Now().UnixNano(), 10),
		Name:            req.Name,
		Description:     req.Description,
		Status:          "creating",
		RegionID:        req.RegionID,
		VPCID:           req.VPCID,
		ProjectID:       req.ProjectID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := h.db.Create(securityGroup).Error; err != nil {
		h.logger.Error("failed to create security group", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, securityGroup)
}

// DeleteSecurityGroup 删除安全组
func (h *NetworkSyncHandler) DeleteSecurityGroup(c *gin.Context) {
	id := c.Param("id")

	var securityGroup model.CloudSecurityGroup
	if err := h.db.First(&securityGroup, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "security group not found"})
		return
	}

	if err := h.db.Delete(&securityGroup).Error; err != nil {
		h.logger.Error("failed to delete security group", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// BatchDeleteSecurityGroups 批量删除安全组
func (h *NetworkSyncHandler) BatchDeleteSecurityGroups(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := h.db.Where("id IN ?", req.IDs).Delete(&model.CloudSecurityGroup{})
	c.JSON(http.StatusOK, gin.H{
		"total":   len(req.IDs),
		"success": result.RowsAffected,
		"failed":  len(req.IDs) - int(result.RowsAffected),
		"message": "batch delete completed",
	})
}

// ========== Network (IP子网) Handlers ==========

// ListNetworks 列出IP子网
func (h *NetworkSyncHandler) ListNetworks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	cloudAccountID, _ := strconv.ParseUint(c.Query("cloud_account_id"), 10, 32)
	status := c.Query("status")
	name := c.Query("name")
	region := c.Query("region")
	vpcId := c.Query("vpc_id")

	var networks []model.CloudSubnet
	var total int64

	query := h.db.Model(&model.CloudSubnet{})
	if cloudAccountID > 0 {
		query = query.Where("cloud_account_id = ?", cloudAccountID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if region != "" {
		query = query.Where("region LIKE ?", "%"+region+"%")
	}
	if vpcId != "" {
		query = query.Where("vpc_id = ?", vpcId)
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&networks).Error; err != nil {
		h.logger.Error("failed to list networks", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     networks,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetNetwork 获取IP子网详情
func (h *NetworkSyncHandler) GetNetwork(c *gin.Context) {
	id := c.Param("id")

	var network model.CloudSubnet
	if err := h.db.First(&network, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "network not found"})
		return
	}

	c.JSON(http.StatusOK, network)
}

// CreateNetwork 创建IP子网
func (h *NetworkSyncHandler) CreateNetwork(c *gin.Context) {
	var req struct {
		CloudAccountID uint   `json:"cloud_account_id" binding:"required"`
		Name           string `json:"name" binding:"required"`
		VPCID          string `json:"vpc_id" binding:"required"`
		GuestIPStart   string `json:"guest_ip_start" binding:"required"`
		GuestIPEnd     string `json:"guest_ip_end" binding:"required"`
		GuestIPMask    int    `json:"guest_ip_mask" binding:"required"`
		GuestGateway   string `json:"guest_gateway"`
		ZoneID         string `json:"zone_id" binding:"required"`
		RegionID       string `json:"region_id" binding:"required"`
		ProjectID      uint   `json:"project_id"`
		IsAutoAlloc    bool   `json:"is_auto_alloc"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	network := &model.CloudSubnet{
		CloudAccountID: req.CloudAccountID,
		SubnetID:       "net-" + strconv.FormatInt(time.Now().UnixNano(), 10),
		Name:           req.Name,
		VPCID:          req.VPCID,
		GuestIPStart:   req.GuestIPStart,
		GuestIPEnd:     req.GuestIPEnd,
		GuestIPMask:    req.GuestIPMask,
		GuestGateway:   req.GuestGateway,
		ZoneID:         req.ZoneID,
		RegionID:       req.RegionID,
		ProjectID:      req.ProjectID,
		IsAutoAlloc:    req.IsAutoAlloc,
		Status:         "creating",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := h.db.Create(network).Error; err != nil {
		h.logger.Error("failed to create network", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, network)
}

// DeleteNetwork 删除IP子网
func (h *NetworkSyncHandler) DeleteNetwork(c *gin.Context) {
	id := c.Param("id")

	var network model.CloudSubnet
	if err := h.db.First(&network, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "network not found"})
		return
	}

	if err := h.db.Delete(&network).Error; err != nil {
		h.logger.Error("failed to delete network", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// BatchDeleteNetworks 批量删除IP子网
func (h *NetworkSyncHandler) BatchDeleteNetworks(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := h.db.Where("id IN ?", req.IDs).Delete(&model.CloudSubnet{})
	c.JSON(http.StatusOK, gin.H{
		"total":   len(req.IDs),
		"success": result.RowsAffected,
		"failed":  len(req.IDs) - int(result.RowsAffected),
		"message": "batch delete completed",
	})
}

// ========== EIP Handlers ==========

// ListEIPs 列出弹性公网IP
func (h *NetworkSyncHandler) ListEIPs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	cloudAccountID, _ := strconv.ParseUint(c.Query("cloud_account_id"), 10, 32)
	status := c.Query("status")
	name := c.Query("name")
	region := c.Query("region")

	var eips []model.CloudEIP
	var total int64

	query := h.db.Model(&model.CloudEIP{})
	if cloudAccountID > 0 {
		query = query.Where("cloud_account_id = ?", cloudAccountID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if region != "" {
		query = query.Where("region LIKE ?", "%"+region+"%")
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&eips).Error; err != nil {
		h.logger.Error("failed to list eips", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     eips,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetEIP 获取弹性公网IP详情
func (h *NetworkSyncHandler) GetEIP(c *gin.Context) {
	id := c.Param("id")

	var eip model.CloudEIP
	if err := h.db.First(&eip, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "eip not found"})
		return
	}

	c.JSON(http.StatusOK, eip)
}

// CreateEIP 创建弹性公网IP
func (h *NetworkSyncHandler) CreateEIP(c *gin.Context) {
	var req struct {
		CloudAccountID uint   `json:"cloud_account_id" binding:"required"`
		Name           string `json:"name" binding:"required"`
		Bandwidth      int    `json:"bandwidth" binding:"required"`
		BillingMethod  string `json:"billing_method"`
		RegionID       string `json:"region_id" binding:"required"`
		ProjectID      uint   `json:"project_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	eip := &model.CloudEIP{
		CloudAccountID: req.CloudAccountID,
		EIPID:          "eip-" + strconv.FormatInt(time.Now().UnixNano(), 10),
		Name:           req.Name,
		Bandwidth:      req.Bandwidth,
		BillingMethod:  req.BillingMethod,
		RegionID:       req.RegionID,
		ProjectID:      req.ProjectID,
		Status:         "creating",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := h.db.Create(eip).Error; err != nil {
		h.logger.Error("failed to create eip", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, eip)
}

// BindEIP 绑定弹性公网IP
func (h *NetworkSyncHandler) BindEIP(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		ResourceID   string `json:"resource_id" binding:"required"`
		ResourceType string `json:"resource_type"`
		ResourceName string `json:"resource_name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var eip model.CloudEIP
	if err := h.db.First(&eip, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "eip not found"})
		return
	}

	eip.ResourceID = req.ResourceID
	eip.ResourceType = req.ResourceType
	eip.ResourceName = req.ResourceName
	eip.Status = "in-use"
	eip.UpdatedAt = time.Now()
	h.db.Save(&eip)

	c.JSON(http.StatusOK, gin.H{"message": "eip bound"})
}

// UnbindEIP 解绑弹性公网IP
func (h *NetworkSyncHandler) UnbindEIP(c *gin.Context) {
	id := c.Param("id")

	var eip model.CloudEIP
	if err := h.db.First(&eip, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "eip not found"})
		return
	}

	eip.ResourceID = ""
	eip.ResourceType = ""
	eip.ResourceName = ""
	eip.Status = "available"
	eip.UpdatedAt = time.Now()
	h.db.Save(&eip)

	c.JSON(http.StatusOK, gin.H{"message": "eip unbound"})
}

// DeleteEIP 删除弹性公网IP
func (h *NetworkSyncHandler) DeleteEIP(c *gin.Context) {
	id := c.Param("id")

	var eip model.CloudEIP
	if err := h.db.First(&eip, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "eip not found"})
		return
	}

	if err := h.db.Delete(&eip).Error; err != nil {
		h.logger.Error("failed to delete eip", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// BatchDeleteEIPs 批量删除弹性公网IP
func (h *NetworkSyncHandler) BatchDeleteEIPs(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := h.db.Where("id IN ?", req.IDs).Delete(&model.CloudEIP{})
	c.JSON(http.StatusOK, gin.H{
		"total":   len(req.IDs),
		"success": result.RowsAffected,
		"failed":  len(req.IDs) - int(result.RowsAffected),
		"message": "batch delete completed",
	})
}

// ========== KeyPair Handlers ==========

// ListKeyPairs 列出密钥
func (h *NetworkSyncHandler) ListKeyPairs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	cloudAccountID, _ := strconv.ParseUint(c.Query("cloud_account_id"), 10, 32)
	status := c.Query("status")
	name := c.Query("name")
	region := c.Query("region")

	var keypairs []model.KeyPair
	var total int64

	query := h.db.Model(&model.KeyPair{})
	if cloudAccountID > 0 {
		query = query.Where("cloud_account_id = ?", cloudAccountID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if region != "" {
		query = query.Where("region LIKE ?", "%"+region+"%")
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&keypairs).Error; err != nil {
		h.logger.Error("failed to list keypairs", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     keypairs,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetKeyPair 获取密钥详情
func (h *NetworkSyncHandler) GetKeyPair(c *gin.Context) {
	id := c.Param("id")

	var keypair model.KeyPair
	if err := h.db.First(&keypair, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "keypair not found"})
		return
	}

	c.JSON(http.StatusOK, keypair)
}

// CreateKeyPair 创建密钥
func (h *NetworkSyncHandler) CreateKeyPair(c *gin.Context) {
	var req struct {
		CloudAccountID uint   `json:"cloud_account_id" binding:"required"`
		Name           string `json:"name" binding:"required"`
		PublicKey      string `json:"public_key"`
		RegionID       string `json:"region_id" binding:"required"`
		ProjectID      uint   `json:"project_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	keypair := &model.KeyPair{
		CloudAccountID: req.CloudAccountID,
		KeyPairID:      "kp-" + strconv.FormatInt(time.Now().UnixNano(), 10),
		Name:           req.Name,
		PublicKey:      req.PublicKey,
		RegionID:       req.RegionID,
		ProjectID:      req.ProjectID,
		Status:         "available",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := h.db.Create(keypair).Error; err != nil {
		h.logger.Error("failed to create keypair", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, keypair)
}

// DeleteKeyPair 删除密钥
func (h *NetworkSyncHandler) DeleteKeyPair(c *gin.Context) {
	id := c.Param("id")

	var keypair model.KeyPair
	if err := h.db.First(&keypair, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "keypair not found"})
		return
	}

	if err := h.db.Delete(&keypair).Error; err != nil {
		h.logger.Error("failed to delete keypair", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// BatchDeleteKeyPairs 批量删除密钥
func (h *NetworkSyncHandler) BatchDeleteKeyPairs(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := h.db.Where("id IN ?", req.IDs).Delete(&model.KeyPair{})
	c.JSON(http.StatusOK, gin.H{
		"total":   len(req.IDs),
		"success": result.RowsAffected,
		"failed":  len(req.IDs) - int(result.RowsAffected),
		"message": "batch delete completed",
	})
}

// ========== NAT Gateway Handlers ==========

// ListNATGateways 列出NAT网关
func (h *NetworkSyncHandler) ListNATGateways(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	cloudAccountID, _ := strconv.ParseUint(c.Query("cloud_account_id"), 10, 32)
	status := c.Query("status")
	name := c.Query("name")
	region := c.Query("region")
	natType := c.Query("nat_type")

	var natGateways []model.CloudNATGateway
	var total int64

	query := h.db.Model(&model.CloudNATGateway{})
	if cloudAccountID > 0 {
		query = query.Where("cloud_account_id = ?", cloudAccountID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if region != "" {
		query = query.Where("region LIKE ?", "%"+region+"%")
	}
	if natType != "" {
		query = query.Where("nat_type = ?", natType)
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&natGateways).Error; err != nil {
		h.logger.Error("failed to list nat gateways", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     natGateways,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetNATGateway 获取NAT网关详情
func (h *NetworkSyncHandler) GetNATGateway(c *gin.Context) {
	id := c.Param("id")

	var natGateway model.CloudNATGateway
	if err := h.db.First(&natGateway, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "nat gateway not found"})
		return
	}

	// 获取关联的规则
	var rules []model.CloudNATRule
	h.db.Where("nat_gateway_id = ?", natGateway.ID).Find(&rules)

	c.JSON(http.StatusOK, gin.H{
		"nat_gateway": natGateway,
		"rules":       rules,
	})
}

// CreateNATGateway 创建NAT网关
func (h *NetworkSyncHandler) CreateNATGateway(c *gin.Context) {
	var req struct {
		CloudAccountID uint   `json:"cloud_account_id" binding:"required"`
		Name           string `json:"name" binding:"required"`
		Description    string `json:"description"`
		NatType        string `json:"nat_type"` // 公网NAT/私网NAT
		Specification  string `json:"specification"`
		BillingMethod  string `json:"billing_method"` // Postpaid/Prepaid
		VpcID          string `json:"vpc_id" binding:"required"`
		SubnetID       string `json:"subnet_id"`
		EipID          string `json:"eip_id"` // 绑定的EIP ID
		RegionID       string `json:"region_id" binding:"required"`
		ProjectID      uint   `json:"project_id"`
		Tags           []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"tags"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查询 EIP 地址
	var eipAddress string
	if req.EipID != "" {
		var eip model.CloudEIP
		if err := h.db.Where("eip_id = ?", req.EipID).First(&eip).Error; err == nil {
			eipAddress = eip.Address
		}
	}

	natGateway := &model.CloudNATGateway{
		CloudAccountID: req.CloudAccountID,
		NatGatewayID:   "nat-" + strconv.FormatInt(time.Now().UnixNano(), 10),
		Name:           req.Name,
		Description:    req.Description,
		NatType:        req.NatType,
		Specification:  req.Specification,
		BillingMethod:  req.BillingMethod,
		VpcID:          req.VpcID,
		SubnetID:       req.SubnetID,
		EipID:          req.EipID,
		EipAddress:     eipAddress,
		RegionID:       req.RegionID,
		ProjectID:      req.ProjectID,
		Status:         "creating",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := h.db.Create(natGateway).Error; err != nil {
		h.logger.Error("failed to create nat gateway", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, natGateway)
}

// UpdateNATGateway 更新NAT网关
func (h *NetworkSyncHandler) UpdateNATGateway(c *gin.Context) {
	id := c.Param("id")

	var natGateway model.CloudNATGateway
	if err := h.db.First(&natGateway, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "nat gateway not found"})
		return
	}

	var req struct {
		Name          string `json:"name"`
		Description   string `json:"description"`
		Specification string `json:"specification"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != "" {
		natGateway.Name = req.Name
	}
	if req.Description != "" {
		natGateway.Description = req.Description
	}
	if req.Specification != "" {
		natGateway.Specification = req.Specification
	}
	natGateway.UpdatedAt = time.Now()

	h.db.Save(&natGateway)
	c.JSON(http.StatusOK, natGateway)
}

// DeleteNATGateway 删除NAT网关
func (h *NetworkSyncHandler) DeleteNATGateway(c *gin.Context) {
	id := c.Param("id")

	var natGateway model.CloudNATGateway
	if err := h.db.First(&natGateway, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "nat gateway not found"})
		return
	}

	// 删除关联的规则
	h.db.Where("nat_gateway_id = ?", natGateway.ID).Delete(&model.CloudNATRule{})

	if err := h.db.Delete(&natGateway).Error; err != nil {
		h.logger.Error("failed to delete nat gateway", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// BatchDeleteNATGateways 批量删除NAT网关
func (h *NetworkSyncHandler) BatchDeleteNATGateways(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 删除关联的规则
	h.db.Where("nat_gateway_id IN ?", req.IDs).Delete(&model.CloudNATRule{})

	result := h.db.Where("id IN ?", req.IDs).Delete(&model.CloudNATGateway{})
	c.JSON(http.StatusOK, gin.H{
		"total":   len(req.IDs),
		"success": result.RowsAffected,
		"failed":  len(req.IDs) - int(result.RowsAffected),
		"message": "batch delete completed",
	})
}

// ========== NAT Rule Handlers ==========

// ListNATRules 列出NAT规则
func (h *NetworkSyncHandler) ListNATRules(c *gin.Context) {
	natGatewayID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var rules []model.CloudNATRule
	var total int64

	query := h.db.Model(&model.CloudNATRule{})
	if natGatewayID > 0 {
		query = query.Where("nat_gateway_id = ?", natGatewayID)
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&rules).Error; err != nil {
		h.logger.Error("failed to list nat rules", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     rules,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// CreateNATRule 创建NAT规则
func (h *NetworkSyncHandler) CreateNATRule(c *gin.Context) {
	natGatewayID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req struct {
		RuleType     string `json:"rule_type" binding:"required"` // SNAT/DNAT
		Name         string `json:"name"`
		ExternalIP   string `json:"external_ip" binding:"required"`
		ExternalPort string `json:"external_port"`
		InternalIP   string `json:"internal_ip" binding:"required"`
		InternalPort string `json:"internal_port"`
		Protocol     string `json:"protocol"` // TCP/UDP/ALL
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rule := &model.CloudNATRule{
		NatGatewayID: uint(natGatewayID),
		RuleID:       "rule-" + strconv.FormatInt(time.Now().UnixNano(), 10),
		RuleType:     req.RuleType,
		Name:         req.Name,
		ExternalIP:   req.ExternalIP,
		ExternalPort: req.ExternalPort,
		InternalIP:   req.InternalIP,
		InternalPort: req.InternalPort,
		Protocol:     req.Protocol,
		Status:       "active",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := h.db.Create(rule).Error; err != nil {
		h.logger.Error("failed to create nat rule", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 更新 NAT Gateway 的规则计数
	var ruleCount int64
	h.db.Model(&model.CloudNATRule{}).Where("nat_gateway_id = ? AND rule_type = ?", natGatewayID, req.RuleType).Count(&ruleCount)
	fieldName := req.RuleType + "_table_entries"
		h.db.Model(&model.CloudNATGateway{}).Where("id = ?", natGatewayID).UpdateColumn(fieldName, ruleCount)

		c.JSON(http.StatusCreated, rule)
}

// DeleteNATRule 删除NAT规则
func (h *NetworkSyncHandler) DeleteNATRule(c *gin.Context) {
	id := c.Param("rule_id")

	var rule model.CloudNATRule
	if err := h.db.First(&rule, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "rule not found"})
		return
	}

	if err := h.db.Delete(&rule).Error; err != nil {
		h.logger.Error("failed to delete nat rule", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 更新 NAT Gateway 的规则计数
	if rule.RuleType == "SNAT" {
		h.db.Model(&model.CloudNATGateway{}).Where("id = ?", rule.NatGatewayID).
			UpdateColumn("snat_table_entries", gorm.Expr("snat_table_entries - 1"))
	} else {
		h.db.Model(&model.CloudNATGateway{}).Where("id = ?", rule.NatGatewayID).
			UpdateColumn("dnat_table_entries", gorm.Expr("dnat_table_entries - 1"))
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
// ========== IPv6 Gateway Handlers ==========

// ListIPv6Gateways 列出IPv6网关
func (h *NetworkSyncHandler) ListIPv6Gateways(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	cloudAccountID, _ := strconv.ParseUint(c.Query("cloud_account_id"), 10, 32)
	status := c.Query("status")
	name := c.Query("name")
	region := c.Query("region")

	var ipv6Gateways []model.CloudIPv6Gateway
	var total int64

	query := h.db.Model(&model.CloudIPv6Gateway{})
	if cloudAccountID > 0 {
		query = query.Where("cloud_account_id = ?", cloudAccountID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if region != "" {
		query = query.Where("region LIKE ?", "%"+region+"%")
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&ipv6Gateways).Error; err != nil {
		h.logger.Error("failed to list ipv6 gateways", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     ipv6Gateways,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetIPv6Gateway 获取IPv6网关详情
func (h *NetworkSyncHandler) GetIPv6Gateway(c *gin.Context) {
	id := c.Param("id")

	var ipv6Gateway model.CloudIPv6Gateway
	if err := h.db.First(&ipv6Gateway, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ipv6 gateway not found"})
		return
	}

	c.JSON(http.StatusOK, ipv6Gateway)
}

// CreateIPv6Gateway 创建IPv6网关
func (h *NetworkSyncHandler) CreateIPv6Gateway(c *gin.Context) {
	var req struct {
		CloudAccountID uint   `json:"cloud_account_id" binding:"required"`
		Name           string `json:"name" binding:"required"`
		VpcID          string `json:"vpc_id" binding:"required"`
		Specification  string `json:"specification"`
		Ipv6Cidr       string `json:"ipv6_cidr"`
		RegionID       string `json:"region_id" binding:"required"`
		ProjectID      uint   `json:"project_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ipv6Gateway := &model.CloudIPv6Gateway{
		CloudAccountID: req.CloudAccountID,
		Ipv6GatewayID:  "ipv6gw-" + strconv.FormatInt(time.Now().UnixNano(), 10),
		Name:           req.Name,
		VpcID:          req.VpcID,
		Specification:  req.Specification,
		Ipv6Cidr:       req.Ipv6Cidr,
		RegionID:       req.RegionID,
		ProjectID:      req.ProjectID,
		Status:         "creating",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := h.db.Create(ipv6Gateway).Error; err != nil {
		h.logger.Error("failed to create ipv6 gateway", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ipv6Gateway)
}

// DeleteIPv6Gateway 删除IPv6网关
func (h *NetworkSyncHandler) DeleteIPv6Gateway(c *gin.Context) {
	id := c.Param("id")

	var ipv6Gateway model.CloudIPv6Gateway
	if err := h.db.First(&ipv6Gateway, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ipv6 gateway not found"})
		return
	}

	if err := h.db.Delete(&ipv6Gateway).Error; err != nil {
		h.logger.Error("failed to delete ipv6 gateway", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// BatchDeleteIPv6Gateways 批量删除IPv6网关
func (h *NetworkSyncHandler) BatchDeleteIPv6Gateways(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := h.db.Where("id IN ?", req.IDs).Delete(&model.CloudIPv6Gateway{})
	c.JSON(http.StatusOK, gin.H{
		"total":   len(req.IDs),
		"success": result.RowsAffected,
		"failed":  len(req.IDs) - int(result.RowsAffected),
		"message": "batch delete completed",
	})
}

// ========== DNS Zone Handlers ==========

// ListDNSZones 列出DNS解析Zone
func (h *NetworkSyncHandler) ListDNSZones(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	cloudAccountID, _ := strconv.ParseUint(c.Query("cloud_account_id"), 10, 32)
	status := c.Query("status")
	name := c.Query("name")
	region := c.Query("region")

	var dnsZones []model.CloudDNSZone
	var total int64

	query := h.db.Model(&model.CloudDNSZone{})
	if cloudAccountID > 0 {
		query = query.Where("cloud_account_id = ?", cloudAccountID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if region != "" {
		query = query.Where("region LIKE ?", "%"+region+"%")
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&dnsZones).Error; err != nil {
		h.logger.Error("failed to list dns zones", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     dnsZones,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetDNSZone 获取DNS Zone详情
func (h *NetworkSyncHandler) GetDNSZone(c *gin.Context) {
	id := c.Param("id")

	var dnsZone model.CloudDNSZone
	if err := h.db.First(&dnsZone, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "dns zone not found"})
		return
	}

	// 获取关联的DNS记录
	var records []model.CloudDNSRecord
	h.db.Where("dns_zone_id = ?", dnsZone.ID).Find(&records)

	c.JSON(http.StatusOK, gin.H{
		"dns_zone": dnsZone,
		"records":  records,
	})
}

// CreateDNSZone 创建DNS Zone
func (h *NetworkSyncHandler) CreateDNSZone(c *gin.Context) {
	var req struct {
		CloudAccountID uint   `json:"cloud_account_id" binding:"required"`
		Name           string `json:"name" binding:"required"`
		AttributionScope string `json:"attribution_scope"`
		RegionID       string `json:"region_id"`
		ProjectID      uint   `json:"project_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dnsZone := &model.CloudDNSZone{
		CloudAccountID:   req.CloudAccountID,
		DnsZoneID:        "dns-" + strconv.FormatInt(time.Now().UnixNano(), 10),
		Name:             req.Name,
		AttributionScope: req.AttributionScope,
		RegionID:         req.RegionID,
		ProjectID:        req.ProjectID,
		Status:           "creating",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := h.db.Create(dnsZone).Error; err != nil {
		h.logger.Error("failed to create dns zone", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dnsZone)
}

// DeleteDNSZone 删除DNS Zone
func (h *NetworkSyncHandler) DeleteDNSZone(c *gin.Context) {
	id := c.Param("id")

	var dnsZone model.CloudDNSZone
	if err := h.db.First(&dnsZone, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "dns zone not found"})
		return
	}

	// 删除关联的DNS记录
	h.db.Where("dns_zone_id = ?", dnsZone.ID).Delete(&model.CloudDNSRecord{})

	if err := h.db.Delete(&dnsZone).Error; err != nil {
		h.logger.Error("failed to delete dns zone", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// BatchDeleteDNSZones 批量删除DNS Zone
func (h *NetworkSyncHandler) BatchDeleteDNSZones(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 删除关联的DNS记录
	h.db.Where("dns_zone_id IN ?", req.IDs).Delete(&model.CloudDNSRecord{})

	result := h.db.Where("id IN ?", req.IDs).Delete(&model.CloudDNSZone{})
	c.JSON(http.StatusOK, gin.H{
		"total":   len(req.IDs),
		"success": result.RowsAffected,
		"failed":  len(req.IDs) - int(result.RowsAffected),
		"message": "batch delete completed",
	})
}

// ========== DNS Record Handlers ==========

// ListDNSRecords 列出DNS记录
func (h *NetworkSyncHandler) ListDNSRecords(c *gin.Context) {
	dnsZoneID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var records []model.CloudDNSRecord
	var total int64

	query := h.db.Model(&model.CloudDNSRecord{})
	if dnsZoneID > 0 {
		query = query.Where("dns_zone_id = ?", dnsZoneID)
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&records).Error; err != nil {
		h.logger.Error("failed to list dns records", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     records,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// CreateDNSRecord 创建DNS记录
func (h *NetworkSyncHandler) CreateDNSRecord(c *gin.Context) {
	dnsZoneID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req struct {
		Name     string `json:"name" binding:"required"`
		Type     string `json:"type" binding:"required"` // A/CNAME/MX/TXT等
		Value    string `json:"value" binding:"required"`
		TTL      int    `json:"ttl"`
		Priority int    `json:"priority"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	record := &model.CloudDNSRecord{
		DnsZoneID: uint(dnsZoneID),
		RecordID:  "rec-" + strconv.FormatInt(time.Now().UnixNano(), 10),
		Name:      req.Name,
		Type:      req.Type,
		Value:     req.Value,
		TTL:       req.TTL,
		Priority:  req.Priority,
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.db.Create(record).Error; err != nil {
		h.logger.Error("failed to create dns record", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, record)
}

// DeleteDNSRecord 删除DNS记录
func (h *NetworkSyncHandler) DeleteDNSRecord(c *gin.Context) {
	id := c.Param("record_id")

	var record model.CloudDNSRecord
	if err := h.db.First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
		return
	}

	if err := h.db.Delete(&record).Error; err != nil {
		h.logger.Error("failed to delete dns record", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// ========== Region Handlers ==========

// ListRegions 列出云区域
func (h *NetworkSyncHandler) ListRegions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	cloudAccountID, _ := strconv.ParseUint(c.Query("cloud_account_id"), 10, 32)
	name := c.Query("name")

	var regions []model.CloudRegion
	var total int64

	query := h.db.Model(&model.CloudRegion{})
	if cloudAccountID > 0 {
		query = query.Where("cloud_account_id = ?", cloudAccountID)
	}
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&regions).Error; err != nil {
		h.logger.Error("failed to list regions", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     regions,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// ========== Zone Handlers ==========

// ListZones 列出可用区
func (h *NetworkSyncHandler) ListZones(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	cloudAccountID, _ := strconv.ParseUint(c.Query("cloud_account_id"), 10, 32)
	regionID := c.Query("region_id")
	name := c.Query("name")

	var zones []model.CloudZone
	var total int64

	query := h.db.Model(&model.CloudZone{})
	if cloudAccountID > 0 {
		query = query.Where("cloud_account_id = ?", cloudAccountID)
	}
	if regionID != "" {
		query = query.Where("region_id = ?", regionID)
	}
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&zones).Error; err != nil {
		h.logger.Error("failed to list zones", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     zones,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// ========== VPC Handlers ==========

// ListVPCsSync 列出VPC（从本地数据库）
func (h *NetworkSyncHandler) ListVPCsSync(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	cloudAccountID, _ := strconv.ParseUint(c.Query("cloud_account_id"), 10, 32)
	status := c.Query("status")
	name := c.Query("name")
	region := c.Query("region")

	var vpcs []model.CloudVPC
	var total int64

	query := h.db.Model(&model.CloudVPC{})
	if cloudAccountID > 0 {
		query = query.Where("cloud_account_id = ?", cloudAccountID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if region != "" {
		query = query.Where("region_id LIKE ?", "%"+region+"%")
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&vpcs).Error; err != nil {
		h.logger.Error("failed to list vpcs", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     vpcs,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// ========== Global VPC Handlers ==========

// ListGlobalVPCs 列出全局VPC
func (h *NetworkSyncHandler) ListGlobalVPCs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	cloudAccountID, _ := strconv.ParseUint(c.Query("cloud_account_id"), 10, 32)
	status := c.Query("status")
	name := c.Query("name")

	var globalVPCs []model.CloudGlobalVPC
	var total int64

	query := h.db.Model(&model.CloudGlobalVPC{})
	if cloudAccountID > 0 {
		query = query.Where("cloud_account_id = ?", cloudAccountID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&globalVPCs).Error; err != nil {
		h.logger.Error("failed to list global vpcs", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     globalVPCs,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// ========== VPC Interconnect Handlers ==========

// ListVPCInterconnects 列出VPC互通
func (h *NetworkSyncHandler) ListVPCInterconnects(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	cloudAccountID, _ := strconv.ParseUint(c.Query("cloud_account_id"), 10, 32)
	status := c.Query("status")
	name := c.Query("name")

	var interconnects []model.CloudVPCInterconnect
	var total int64

	query := h.db.Model(&model.CloudVPCInterconnect{})
	if cloudAccountID > 0 {
		query = query.Where("cloud_account_id = ?", cloudAccountID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&interconnects).Error; err != nil {
		h.logger.Error("failed to list vpc interconnects", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     interconnects,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// ========== L2 Network Handlers ==========

// ListL2Networks 列出二层网络
func (h *NetworkSyncHandler) ListL2Networks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	cloudAccountID, _ := strconv.ParseUint(c.Query("cloud_account_id"), 10, 32)
	status := c.Query("status")
	name := c.Query("name")
	region := c.Query("region")

	var l2Networks []model.CloudL2Network
	var total int64

	query := h.db.Model(&model.CloudL2Network{})
	if cloudAccountID > 0 {
		query = query.Where("cloud_account_id = ?", cloudAccountID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if region != "" {
		query = query.Where("region_id LIKE ?", "%"+region+"%")
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&l2Networks).Error; err != nil {
		h.logger.Error("failed to list l2 networks", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     l2Networks,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// ========== Route Table Handlers ==========

// ListRouteTables 列出路由表
func (h *NetworkSyncHandler) ListRouteTables(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	cloudAccountID, _ := strconv.ParseUint(c.Query("cloud_account_id"), 10, 32)
	status := c.Query("status")
	name := c.Query("name")
	vpcID := c.Query("vpc_id")

	var routeTables []model.CloudRouteTable
	var total int64

	query := h.db.Model(&model.CloudRouteTable{})
	if cloudAccountID > 0 {
		query = query.Where("cloud_account_id = ?", cloudAccountID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if vpcID != "" {
		query = query.Where("vpc_id = ?", vpcID)
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&routeTables).Error; err != nil {
		h.logger.Error("failed to list route tables", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     routeTables,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// ========== Load Balancer Instance Handlers ==========

// ListLBInstances 列出负载均衡实例
func (h *NetworkSyncHandler) ListLBInstances(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	cloudAccountID, _ := strconv.ParseUint(c.Query("cloud_account_id"), 10, 32)
	status := c.Query("status")
	name := c.Query("name")
	region := c.Query("region")

	var lbInstances []model.CloudLBInstance
	var total int64

	query := h.db.Model(&model.CloudLBInstance{})
	if cloudAccountID > 0 {
		query = query.Where("cloud_account_id = ?", cloudAccountID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if region != "" {
		query = query.Where("region_id LIKE ?", "%"+region+"%")
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&lbInstances).Error; err != nil {
		h.logger.Error("failed to list lb instances", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     lbInstances,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// ========== Load Balancer ACL Handlers ==========

// ListLBACLs 列出负载均衡ACL
func (h *NetworkSyncHandler) ListLBACLs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	cloudAccountID, _ := strconv.ParseUint(c.Query("cloud_account_id"), 10, 32)
	name := c.Query("name")

	var lbACLs []model.CloudLBACL
	var total int64

	query := h.db.Model(&model.CloudLBACL{})
	if cloudAccountID > 0 {
		query = query.Where("cloud_account_id = ?", cloudAccountID)
	}
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&lbACLs).Error; err != nil {
		h.logger.Error("failed to list lb acls", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     lbACLs,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// ========== Load Balancer Certificate Handlers ==========

// ListLBCertificates 列出负载均衡证书
func (h *NetworkSyncHandler) ListLBCertificates(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	cloudAccountID, _ := strconv.ParseUint(c.Query("cloud_account_id"), 10, 32)
	name := c.Query("name")

	var lbCertificates []model.CloudLBCertificate
	var total int64

	query := h.db.Model(&model.CloudLBCertificate{})
	if cloudAccountID > 0 {
		query = query.Where("cloud_account_id = ?", cloudAccountID)
	}
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&lbCertificates).Error; err != nil {
		h.logger.Error("failed to list lb certificates", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     lbCertificates,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// ========== CDN Domain Handlers ==========

// ListCDNDomains 列出CDN域名
func (h *NetworkSyncHandler) ListCDNDomains(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	cloudAccountID, _ := strconv.ParseUint(c.Query("cloud_account_id"), 10, 32)
	status := c.Query("status")
	domain := c.Query("domain")

	var cdnDomains []model.CloudCDNDomain
	var total int64

	query := h.db.Model(&model.CloudCDNDomain{})
	if cloudAccountID > 0 {
		query = query.Where("cloud_account_id = ?", cloudAccountID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if domain != "" {
		query = query.Where("domain LIKE ?", "%"+domain+"%")
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&cdnDomains).Error; err != nil {
		h.logger.Error("failed to list cdn domains", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     cdnDomains,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
