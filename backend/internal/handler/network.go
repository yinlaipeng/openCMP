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

// ListRegions 获取区域列表
func (h *NetworkHandler) ListRegions(c *gin.Context) {
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

	regions, err := h.service.ListRegions(c.Request.Context(), uint(accountID))
	if err != nil {
		h.logger.Error("failed to list regions", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": regions,
		"total": len(regions),
	})
}

// ListZones 获取可用区列表
func (h *NetworkHandler) ListZones(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	regionID := c.Query("region_id")

	if accountIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account_id is required"})
		return
	}

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	zones, err := h.service.ListZones(c.Request.Context(), uint(accountID), regionID)
	if err != nil {
		h.logger.Error("failed to list zones", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": zones,
		"total": len(zones),
	})
}

// CreateVPCInterconnectRequest 创建 VPC 互联请求
type CreateVPCInterconnectRequest struct {
	AccountID uint              `json:"account_id" binding:"required"`
	Name      string            `json:"name" binding:"required"`
	Type      string            `json:"type" binding:"required"`
	Bandwidth int               `json:"bandwidth" binding:"required"`
	RegionID  string            `json:"region_id" binding:"required"`
	PeerRegion string           `json:"peer_region" binding:"required"`
	Tags      map[string]string `json:"tags"`
}

// CreateVPCInterconnect 创建 VPC 互联
func (h *NetworkHandler) CreateVPCInterconnect(c *gin.Context) {
	var req CreateVPCInterconnectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config := cloudprovider.VPCInterconnectConfig{
		Name:       req.Name,
		Type:       req.Type,
		Bandwidth:  req.Bandwidth,
		RegionID:   req.RegionID,
		PeerRegion: req.PeerRegion,
		Tags:       req.Tags,
	}

	interconnect, err := h.service.CreateVPCInterconnect(c.Request.Context(), req.AccountID, config)
	if err != nil {
		h.logger.Error("failed to create vpc interconnect", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, interconnect)
}

// ListVPCInterconnects 列出 VPC 互联
func (h *NetworkHandler) ListVPCInterconnects(c *gin.Context) {
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

	filter := cloudprovider.VPCInterconnectFilter{
		RegionID: c.Query("region_id"),
	}

	interconnects, err := h.service.ListVPCInterconnects(c.Request.Context(), uint(accountID), filter)
	if err != nil {
		h.logger.Error("failed to list vpc interconnects", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": interconnects,
		"total": len(interconnects),
	})
}

// CreateVPCPeeringRequest 创建 VPC 对等连接请求
type CreateVPCPeeringRequest struct {
	AccountID   uint              `json:"account_id" binding:"required"`
	Name        string            `json:"name" binding:"required"`
	LocalVPCID  string            `json:"local_vpc_id" binding:"required"`
	PeerVPCID   string            `json:"peer_vpc_id" binding:"required"`
	PeerAccount string            `json:"peer_account" binding:"required"`
	Tags        map[string]string `json:"tags"`
}

// CreateVPCPeering 创建 VPC 对等连接
func (h *NetworkHandler) CreateVPCPeering(c *gin.Context) {
	var req CreateVPCPeeringRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config := cloudprovider.VPCPeeringConfig{
		Name:         req.Name,
		LocalVPCID:   req.LocalVPCID,
		PeerVPCID:    req.PeerVPCID,
		PeerAccount:  req.PeerAccount,
		Tags:         req.Tags,
	}

	peering, err := h.service.CreateVPCPeering(c.Request.Context(), req.AccountID, config)
	if err != nil {
		h.logger.Error("failed to create vpc peering", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, peering)
}

// ListVPCPeerings 列出 VPC 对等连接
func (h *NetworkHandler) ListVPCPeerings(c *gin.Context) {
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

	filter := cloudprovider.VPCPeeringFilter{
		LocalVPCID: c.Query("local_vpc_id"),
	}

	peerings, err := h.service.ListVPCPeerings(c.Request.Context(), uint(accountID), filter)
	if err != nil {
		h.logger.Error("failed to list vpc peerings", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": peerings,
		"total": len(peerings),
	})
}

// CreateRouteTableRequest 创建路由表请求
type CreateRouteTableRequest struct {
	AccountID uint              `json:"account_id" binding:"required"`
	Name      string            `json:"name" binding:"required"`
	VPCID     string            `json:"vpc_id" binding:"required"`
	Routes    []cloudprovider.Route `json:"routes"`
	Tags      map[string]string `json:"tags"`
}

// CreateRouteTable 创建路由表
func (h *NetworkHandler) CreateRouteTable(c *gin.Context) {
	var req CreateRouteTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config := cloudprovider.RouteTableConfig{
		Name:   req.Name,
		VPCID:  req.VPCID,
		Routes: req.Routes,
		Tags:   req.Tags,
	}

	routeTable, err := h.service.CreateRouteTable(c.Request.Context(), req.AccountID, config)
	if err != nil {
		h.logger.Error("failed to create route table", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, routeTable)
}

// ListRouteTables 列出路由表
func (h *NetworkHandler) ListRouteTables(c *gin.Context) {
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

	filter := cloudprovider.RouteTableFilter{
		VPCID: c.Query("vpc_id"),
	}

	routeTables, err := h.service.ListRouteTables(c.Request.Context(), uint(accountID), filter)
	if err != nil {
		h.logger.Error("failed to list route tables", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": routeTables,
		"total": len(routeTables),
	})
}

// CreateL2NetworkRequest 创建二层网络请求
type CreateL2NetworkRequest struct {
	AccountID uint              `json:"account_id" binding:"required"`
	Name      string            `json:"name" binding:"required"`
	VLANID    int               `json:"vlan_id" binding:"required"`
	VPCID     string            `json:"vpc_id" binding:"required"`
	Subnets   []string          `json:"subnets"`
	Tags      map[string]string `json:"tags"`
}

// CreateL2Network 创建二层网络
func (h *NetworkHandler) CreateL2Network(c *gin.Context) {
	var req CreateL2NetworkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config := cloudprovider.L2NetworkConfig{
		Name:    req.Name,
		VLANID:  req.VLANID,
		VPCID:   req.VPCID,
		Subnets: req.Subnets,
		Tags:    req.Tags,
	}

	l2Network, err := h.service.CreateL2Network(c.Request.Context(), req.AccountID, config)
	if err != nil {
		h.logger.Error("failed to create l2 network", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, l2Network)
}

// ListL2Networks 列出二层网络
func (h *NetworkHandler) ListL2Networks(c *gin.Context) {
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

	filter := cloudprovider.L2NetworkFilter{
		VPCID: c.Query("vpc_id"),
	}

	l2Networks, err := h.service.ListL2Networks(c.Request.Context(), uint(accountID), filter)
	if err != nil {
		h.logger.Error("failed to list l2 networks", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": l2Networks,
		"total": len(l2Networks),
	})
}

// DeleteVPCInterconnect 删除 VPC 互联
func (h *NetworkHandler) DeleteVPCInterconnect(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	interconnectID := c.Param("id")

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	if err := h.service.DeleteVPCInterconnect(c.Request.Context(), uint(accountID), interconnectID); err != nil {
		h.logger.Error("failed to delete vpc interconnect", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "vpc interconnect deleted"})
}

// DeleteVPCPeering 删除 VPC 对等连接
func (h *NetworkHandler) DeleteVPCPeering(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	peeringID := c.Param("id")

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	if err := h.service.DeleteVPCPeering(c.Request.Context(), uint(accountID), peeringID); err != nil {
		h.logger.Error("failed to delete vpc peering", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "vpc peering deleted"})
}

// DeleteRouteTable 删除路由表
func (h *NetworkHandler) DeleteRouteTable(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	routeTableID := c.Param("id")

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	if err := h.service.DeleteRouteTable(c.Request.Context(), uint(accountID), routeTableID); err != nil {
		h.logger.Error("failed to delete route table", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "route table deleted"})
}

// DeleteL2Network 删除二层网络
func (h *NetworkHandler) DeleteL2Network(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	l2NetworkID := c.Param("id")

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	if err := h.service.DeleteL2Network(c.Request.Context(), uint(accountID), l2NetworkID); err != nil {
		h.logger.Error("failed to delete l2 network", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "l2 network deleted"})
}

// ========== 子网扩展操作 ==========

// UpdateSubnetRequest 更新子网请求
type UpdateSubnetRequest struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Tags        map[string]string `json:"tags"`
}

// UpdateSubnet 更新子网属性
func (h *NetworkHandler) UpdateSubnet(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	subnetID := c.Param("id")

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	var req UpdateSubnetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subnet, err := h.service.UpdateSubnet(c.Request.Context(), uint(accountID), subnetID, req.Name, req.Description, req.Tags)
	if err != nil {
		h.logger.Error("failed to update subnet", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subnet)
}

// DeleteSubnet 删除子网
func (h *NetworkHandler) DeleteSubnet(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	subnetID := c.Param("id")

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	if err := h.service.DeleteSubnet(c.Request.Context(), uint(accountID), subnetID); err != nil {
		h.logger.Error("failed to delete subnet", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "subnet deleted"})
}

// ChangeSubnetProjectRequest 更改子网项目请求
type ChangeSubnetProjectRequest struct {
	ProjectID uint `json:"project_id" binding:"required"`
}

// ChangeSubnetProject 更改子网所属项目
func (h *NetworkHandler) ChangeSubnetProject(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	subnetID := c.Param("id")

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	var req ChangeSubnetProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.ChangeSubnetProject(c.Request.Context(), uint(accountID), subnetID, req.ProjectID); err != nil {
		h.logger.Error("failed to change subnet project", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "project changed", "project_id": req.ProjectID})
}

// SplitSubnetRequest 分割子网请求
type SplitSubnetRequest struct {
	NewCIDRs []string `json:"new_cidrs" binding:"required"`
}

// SplitSubnet 分割IP子网
func (h *NetworkHandler) SplitSubnet(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	subnetID := c.Param("id")

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	var req SplitSubnetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newSubnets, err := h.service.SplitSubnet(c.Request.Context(), uint(accountID), subnetID, req.NewCIDRs)
	if err != nil {
		h.logger.Error("failed to split subnet", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "subnet split", "new_subnets": newSubnets})
}

// ReserveIPRequest 预留IP请求
type ReserveIPRequest struct {
	IPs       []string `json:"ips" binding:"required"`
	Reason    string   `json:"reason"`
	ReservedBy string  `json:"reserved_by"`
}

// ReserveIP 预留IP地址
func (h *NetworkHandler) ReserveIP(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	subnetID := c.Param("id")

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	var req ReserveIPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reservedIPs, err := h.service.ReserveIP(c.Request.Context(), uint(accountID), subnetID, req.IPs, req.Reason, req.ReservedBy)
	if err != nil {
		h.logger.Error("failed to reserve ip", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ip reserved", "reserved_ips": reservedIPs})
}

// ReleaseIPRequest 释放预留IP请求
type ReleaseIPRequest struct {
	IPs []string `json:"ips" binding:"required"`
}

// ReleaseIP 释放预留IP地址
func (h *NetworkHandler) ReleaseIP(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	subnetID := c.Param("id")

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	var req ReleaseIPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.ReleaseIP(c.Request.Context(), uint(accountID), subnetID, req.IPs); err != nil {
		h.logger.Error("failed to release ip", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ip released"})
}

// ListReservedIPs 列出预留IP
func (h *NetworkHandler) ListReservedIPs(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	subnetID := c.Query("subnet_id")

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	reservedIPs, err := h.service.ListReservedIPs(c.Request.Context(), uint(accountID), subnetID)
	if err != nil {
		h.logger.Error("failed to list reserved ips", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": reservedIPs, "total": len(reservedIPs)})
}

// ========== 安全组扩展操作 ==========

// AddSecurityGroupRuleRequest 添加安全组规则请求
type AddSecurityGroupRuleRequest struct {
	Direction   string `json:"direction" binding:"required"` // ingress/egress
	Protocol    string `json:"protocol" binding:"required"`
	PortRange   string `json:"port_range"`
	CIDR        string `json:"cidr" binding:"required"`
	Action      string `json:"action" binding:"required"`
	Description string `json:"description"`
	Priority    int    `json:"priority"`
}

// AddSecurityGroupRule 添加安全组规则
func (h *NetworkHandler) AddSecurityGroupRule(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	sgID := c.Param("id")

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	var req AddSecurityGroupRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rule := cloudprovider.SGRule{
		Direction:   req.Direction,
		Protocol:    req.Protocol,
		PortRange:   req.PortRange,
		CIDR:        req.CIDR,
		Action:      req.Action,
		Description: req.Description,
		Priority:    req.Priority,
	}

	ruleID, err := h.service.AddSecurityGroupRule(c.Request.Context(), uint(accountID), sgID, rule)
	if err != nil {
		h.logger.Error("failed to add security group rule", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "rule added", "rule_id": ruleID})
}

// DeleteSecurityGroupRule 删除安全组规则
func (h *NetworkHandler) DeleteSecurityGroupRule(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	sgID := c.Param("id")
	ruleID := c.Param("rule_id")

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	if err := h.service.DeleteSecurityGroupRule(c.Request.Context(), uint(accountID), sgID, ruleID); err != nil {
		h.logger.Error("failed to delete security group rule", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "rule deleted"})
}

// DeleteSecurityGroup 删除安全组
func (h *NetworkHandler) DeleteSecurityGroup(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	sgID := c.Param("id")

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	if err := h.service.DeleteSecurityGroup(c.Request.Context(), uint(accountID), sgID); err != nil {
		h.logger.Error("failed to delete security group", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "security group deleted"})
}

// ========== EIP 扩展操作 ==========

// BindEIPRequest 绑定EIP请求
type BindEIPRequest struct {
	ResourceID   string `json:"resource_id" binding:"required"`
	ResourceType string `json:"resource_type" binding:"required"` // vm/nat_gateway/slb
}

// BindEIP 绑定弹性IP
func (h *NetworkHandler) BindEIP(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	eipID := c.Param("id")

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	var req BindEIPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.BindEIP(c.Request.Context(), uint(accountID), eipID, req.ResourceID, req.ResourceType); err != nil {
		h.logger.Error("failed to bind eip", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "eip bound", "resource_id": req.ResourceID})
}

// UnbindEIP 解绑弹性IP
func (h *NetworkHandler) UnbindEIP(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	eipID := c.Param("id")

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	if err := h.service.UnbindEIP(c.Request.Context(), uint(accountID), eipID); err != nil {
		h.logger.Error("failed to unbind eip", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "eip unbound"})
}

// DeleteEIP 删除弹性IP
func (h *NetworkHandler) DeleteEIP(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	eipID := c.Param("id")

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	if err := h.service.DeleteEIP(c.Request.Context(), uint(accountID), eipID); err != nil {
		h.logger.Error("failed to delete eip", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "eip deleted"})
}
