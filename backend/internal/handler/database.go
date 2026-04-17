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

// DatabaseHandler 数据库资源 Handler
type DatabaseHandler struct {
	service *service.DatabaseService
	logger  *zap.Logger
}

// NewDatabaseHandler 创建数据库资源 Handler
func NewDatabaseHandler(db *gorm.DB, logger *zap.Logger) *DatabaseHandler {
	return &DatabaseHandler{
		service: service.NewDatabaseService(db),
		logger:  logger,
	}
}

// CreateRDSRequest 创建 RDS 请求
type CreateRDSRequest struct {
	AccountID      uint              `json:"account_id" binding:"required"`
	Name           string            `json:"name" binding:"required"`
	Engine         string            `json:"engine" binding:"required"`         // MySQL/PostgreSQL/SQLServer
	EngineVersion  string            `json:"engine_version" binding:"required"` // 5.7/8.0/12 etc.
	InstanceType   string            `json:"instance_type" binding:"required"`
	StorageSize    int               `json:"storage_size"` // GB
	StorageType    string            `json:"storage_type"` // SSD/ESSD
	VPCID          string            `json:"vpc_id" binding:"required"`
	SubnetID       string            `json:"subnet_id" binding:"required"`
	ZoneID         string            `json:"zone_id"`
	MasterUsername string            `json:"master_username"`
	MasterPassword string            `json:"master_password"`
	Tags           map[string]string `json:"tags"`
}

// CreateRDS 创建 RDS 实例
func (h *DatabaseHandler) CreateRDS(c *gin.Context) {
	var req CreateRDSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config := cloudprovider.RDSConfig{
		Name:           req.Name,
		Engine:         req.Engine,
		EngineVersion:  req.EngineVersion,
		InstanceType:   req.InstanceType,
		StorageSize:    req.StorageSize,
		StorageType:    req.StorageType,
		VPCID:          req.VPCID,
		SubnetID:       req.SubnetID,
		ZoneID:         req.ZoneID,
		MasterUsername: req.MasterUsername,
		MasterPassword: req.MasterPassword,
		Tags:           req.Tags,
	}

	instance, err := h.service.CreateRDS(c.Request.Context(), req.AccountID, config)
	if err != nil {
		h.logger.Error("failed to create rds instance", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, instance)
}

// ListRDS 列出 RDS 实例
func (h *DatabaseHandler) ListRDS(c *gin.Context) {
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

	filter := cloudprovider.RDSFilter{
		InstanceID: c.Query("instance_id"),
		Engine:     c.Query("engine"),
		Status:     c.Query("status"),
		VPCID:      c.Query("vpc_id"),
	}

	instances, err := h.service.ListRDS(c.Request.Context(), uint(accountID), filter)
	if err != nil {
		h.logger.Error("failed to list rds instances", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, instances)
}

// GetRDS 获取 RDS 实例详情
func (h *DatabaseHandler) GetRDS(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	instanceID := c.Param("id")

	if accountIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account_id is required"})
		return
	}

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	instance, err := h.service.GetRDS(c.Request.Context(), uint(accountID), instanceID)
	if err != nil {
		h.logger.Error("failed to get rds instance", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, instance)
}

// DeleteRDS 删除 RDS 实例
func (h *DatabaseHandler) DeleteRDS(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	instanceID := c.Param("id")

	if accountIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account_id is required"})
		return
	}

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	err = h.service.DeleteRDS(c.Request.Context(), uint(accountID), instanceID)
	if err != nil {
		h.logger.Error("failed to delete rds instance", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "rds instance deleted successfully"})
}

// RDSActionRequest RDS 操作请求
type RDSActionRequest struct {
	Action string `json:"action" binding:"required"` // start/stop/reboot
}

// RDSAction RDS 实例操作
func (h *DatabaseHandler) RDSAction(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	instanceID := c.Param("id")

	var req RDSActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if accountIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account_id is required"})
		return
	}

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	err = h.service.RDSAction(c.Request.Context(), uint(accountID), instanceID, req.Action)
	if err != nil {
		h.logger.Error("failed to execute rds action", zap.Error(err), zap.String("action", req.Action))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "rds action executed successfully"})
}

// ResizeRDSRequest RDS 调整规格请求
type ResizeRDSRequest struct {
	InstanceType string `json:"instance_type"`
	StorageSize  int    `json:"storage_size"`
}

// ResizeRDS 调整 RDS 规格
func (h *DatabaseHandler) ResizeRDS(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	instanceID := c.Param("id")

	var req ResizeRDSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if accountIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account_id is required"})
		return
	}

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	spec := cloudprovider.RDSpec{
		InstanceType: req.InstanceType,
		StorageSize:  req.StorageSize,
	}

	err = h.service.ResizeRDS(c.Request.Context(), uint(accountID), instanceID, spec)
	if err != nil {
		h.logger.Error("failed to resize rds instance", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "rds instance resized successfully"})
}

// CreateRDSBackup 创建 RDS 备份
func (h *DatabaseHandler) CreateRDSBackup(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	instanceID := c.Param("id")
	name := c.Query("name")

	if accountIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account_id is required"})
		return
	}

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	backup, err := h.service.CreateRDSBackup(c.Request.Context(), uint(accountID), instanceID, name)
	if err != nil {
		h.logger.Error("failed to create rds backup", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, backup)
}

// ListRDSBackups 列出 RDS 备份
func (h *DatabaseHandler) ListRDSBackups(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	instanceID := c.Param("id")

	if accountIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account_id is required"})
		return
	}

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	backups, err := h.service.ListRDSBackups(c.Request.Context(), uint(accountID), instanceID)
	if err != nil {
		h.logger.Error("failed to list rds backups", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, backups)
}

// CreateCacheRequest 创建 Redis 请求
type CreateCacheRequest struct {
	AccountID      uint              `json:"account_id" binding:"required"`
	Name           string            `json:"name" binding:"required"`
	Engine         string            `json:"engine" binding:"required"`         // Redis/Memcached
	EngineVersion  string            `json:"engine_version" binding:"required"` // 5.0/6.0/7.0
	InstanceType   string            `json:"instance_type" binding:"required"`
	VPCID          string            `json:"vpc_id" binding:"required"`
	SubnetID       string            `json:"subnet_id" binding:"required"`
	ZoneID         string            `json:"zone_id"`
	Tags           map[string]string `json:"tags"`
}

// CreateCache 创建缓存实例
func (h *DatabaseHandler) CreateCache(c *gin.Context) {
	var req CreateCacheRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config := cloudprovider.CacheConfig{
		Name:          req.Name,
		Engine:        req.Engine,
		EngineVersion: req.EngineVersion,
		InstanceType:  req.InstanceType,
		VPCID:         req.VPCID,
		SubnetID:      req.SubnetID,
		ZoneID:        req.ZoneID,
		Tags:          req.Tags,
	}

	instance, err := h.service.CreateCache(c.Request.Context(), req.AccountID, config)
	if err != nil {
		h.logger.Error("failed to create cache instance", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, instance)
}

// ListCache 列出缓存实例
func (h *DatabaseHandler) ListCache(c *gin.Context) {
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

	filter := cloudprovider.CacheFilter{
		InstanceID: c.Query("instance_id"),
		Engine:     c.Query("engine"),
		Status:     c.Query("status"),
	}

	instances, err := h.service.ListCache(c.Request.Context(), uint(accountID), filter)
	if err != nil {
		h.logger.Error("failed to list cache instances", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, instances)
}

// DeleteCache 删除缓存实例
func (h *DatabaseHandler) DeleteCache(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	instanceID := c.Param("id")

	if accountIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account_id is required"})
		return
	}

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	err = h.service.DeleteCache(c.Request.Context(), uint(accountID), instanceID)
	if err != nil {
		h.logger.Error("failed to delete cache instance", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "cache instance deleted successfully"})
}

// CacheActionRequest 缓存实例操作请求
type CacheActionRequest struct {
	Action string `json:"action" binding:"required"` // reboot
}

// CacheAction 缓存实例操作
func (h *DatabaseHandler) CacheAction(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	instanceID := c.Param("id")

	var req CacheActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if accountIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account_id is required"})
		return
	}

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	err = h.service.CacheAction(c.Request.Context(), uint(accountID), instanceID, req.Action)
	if err != nil {
		h.logger.Error("failed to execute cache action", zap.Error(err), zap.String("action", req.Action))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "cache action executed successfully"})
}

// ResizeCacheRequest 缓存实例调整规格请求
type ResizeCacheRequest struct {
	InstanceType string `json:"instance_type"`
}

// ResizeCache 调整缓存实例规格
func (h *DatabaseHandler) ResizeCache(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	instanceID := c.Param("id")

	var req ResizeCacheRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if accountIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account_id is required"})
		return
	}

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	spec := cloudprovider.CacheSpec{
		InstanceType: req.InstanceType,
	}

	err = h.service.ResizeCache(c.Request.Context(), uint(accountID), instanceID, spec)
	if err != nil {
		h.logger.Error("failed to resize cache instance", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "cache instance resized successfully"})
}

// CreateCacheBackup 创建缓存备份
func (h *DatabaseHandler) CreateCacheBackup(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	instanceID := c.Param("id")

	if accountIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account_id is required"})
		return
	}

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	backup, err := h.service.CreateCacheBackup(c.Request.Context(), uint(accountID), instanceID)
	if err != nil {
		h.logger.Error("failed to create cache backup", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, backup)
}