package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
	"github.com/opencmp/opencmp/pkg/cloudprovider"
)

// StorageHandler 存储资源 Handler
type StorageHandler struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewStorageHandler 创建存储资源 Handler
func NewStorageHandler(db *gorm.DB, logger *zap.Logger) *StorageHandler {
	return &StorageHandler{db: db, logger: logger}
}

// ListCloudDisks 列出云硬盘
func (h *StorageHandler) ListCloudDisks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	cloudAccountID, _ := strconv.ParseUint(c.Query("cloud_account_id"), 10, 32)
	status := c.Query("status")

	var disks []model.CloudDisk
	var total int64

	query := h.db.Model(&model.CloudDisk{})
	if cloudAccountID > 0 {
		query = query.Where("cloud_account_id = ?", cloudAccountID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&disks).Error; err != nil {
		h.logger.Error("failed to list cloud disks", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     disks,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// CreateCloudDisk 创建云硬盘
func (h *StorageHandler) CreateCloudDisk(c *gin.Context) {
	cloudAccountID, _ := strconv.ParseUint(c.Query("cloud_account_id"), 10, 32)

	var req struct {
		Name   string `json:"name" binding:"required"`
		Size   int    `json:"size" binding:"required"`
		Type   string `json:"type"`
		ZoneID string `json:"zone_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取云账户
	var account model.CloudAccount
	if err := h.db.First(&account, cloudAccountID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cloud account not found"})
		return
	}

	// 获取云厂商适配器
	provider, err := h.getProvider(&account)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 调用云厂商创建磁盘
	config := cloudprovider.DiskConfig{
		Name:   req.Name,
		Size:   req.Size,
		Type:   req.Type,
		ZoneID: req.ZoneID,
	}

	disk, err := provider.CreateDisk(c.Request.Context(), config)
	if err != nil {
		h.logger.Error("failed to create cloud disk", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 存储到数据库
	cloudDisk := &model.CloudDisk{
		CloudAccountID: uint(cloudAccountID),
		DiskID:         disk.ID,
		Name:           disk.Name,
		Size:           disk.Size,
		Type:           disk.Type,
		Status:         disk.Status,
		ZoneID:         disk.ZoneID,
		ProviderType:   account.ProviderType,
	}

	if err := h.db.Create(cloudDisk).Error; err != nil {
		h.logger.Error("failed to save cloud disk", zap.Error(err))
	}

	c.JSON(http.StatusCreated, cloudDisk)
}

// DeleteCloudDisk 删除云硬盘
func (h *StorageHandler) DeleteCloudDisk(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var disk model.CloudDisk
	if err := h.db.First(&disk, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "disk not found"})
		return
	}

	// 获取云账户
	var account model.CloudAccount
	if err := h.db.First(&account, disk.CloudAccountID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cloud account not found"})
		return
	}

	// 获取云厂商适配器
	provider, err := h.getProvider(&account)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 调用云厂商删除磁盘
	err = provider.DeleteDisk(c.Request.Context(), disk.DiskID)
	if err != nil {
		h.logger.Error("failed to delete cloud disk", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 从数据库删除
	if err := h.db.Delete(&disk).Error; err != nil {
		h.logger.Error("failed to delete disk record", zap.Error(err))
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// AttachCloudDisk 挂载云硬盘
func (h *StorageHandler) AttachCloudDisk(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req struct {
		VMID string `json:"vm_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var disk model.CloudDisk
	if err := h.db.First(&disk, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "disk not found"})
		return
	}

	var account model.CloudAccount
	if err := h.db.First(&account, disk.CloudAccountID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cloud account not found"})
		return
	}

	provider, err := h.getProvider(&account)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = provider.AttachDisk(c.Request.Context(), disk.DiskID, req.VMID)
	if err != nil {
		h.logger.Error("failed to attach cloud disk", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 更新状态
	disk.Status = "in_use"
	disk.VMID = req.VMID
	h.db.Save(&disk)

	c.JSON(http.StatusOK, gin.H{"message": "attached"})
}

// DetachCloudDisk 卸载云硬盘
func (h *StorageHandler) DetachCloudDisk(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var disk model.CloudDisk
	if err := h.db.First(&disk, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "disk not found"})
		return
	}

	var account model.CloudAccount
	if err := h.db.First(&account, disk.CloudAccountID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cloud account not found"})
		return
	}

	provider, err := h.getProvider(&account)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = provider.DetachDisk(c.Request.Context(), disk.DiskID)
	if err != nil {
		h.logger.Error("failed to detach cloud disk", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	disk.Status = "available"
	disk.VMID = ""
	h.db.Save(&disk)

	c.JSON(http.StatusOK, gin.H{"message": "detached"})
}

// ResizeCloudDisk 扩容云硬盘
func (h *StorageHandler) ResizeCloudDisk(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req struct {
		NewSize int `json:"new_size" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var disk model.CloudDisk
	if err := h.db.First(&disk, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "disk not found"})
		return
	}

	var account model.CloudAccount
	if err := h.db.First(&account, disk.CloudAccountID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cloud account not found"})
		return
	}

	provider, err := h.getProvider(&account)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = provider.ResizeDisk(c.Request.Context(), disk.DiskID, req.NewSize)
	if err != nil {
		h.logger.Error("failed to resize cloud disk", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	disk.Size = req.NewSize
	h.db.Save(&disk)

	c.JSON(http.StatusOK, gin.H{"message": "resized", "new_size": req.NewSize})
}

// ListCloudSnapshots 列出云快照
func (h *StorageHandler) ListCloudSnapshots(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	cloudAccountID, _ := strconv.ParseUint(c.Query("cloud_account_id"), 10, 32)

	var snapshots []model.CloudSnapshot
	var total int64

	query := h.db.Model(&model.CloudSnapshot{})
	if cloudAccountID > 0 {
		query = query.Where("cloud_account_id = ?", cloudAccountID)
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&snapshots).Error; err != nil {
		h.logger.Error("failed to list cloud snapshots", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     snapshots,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// CreateCloudSnapshot 创建云快照
func (h *StorageHandler) CreateCloudSnapshot(c *gin.Context) {
	diskID, _ := strconv.ParseUint(c.Query("disk_id"), 10, 32)

	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var disk model.CloudDisk
	if err := h.db.First(&disk, diskID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "disk not found"})
		return
	}

	var account model.CloudAccount
	if err := h.db.First(&account, disk.CloudAccountID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cloud account not found"})
		return
	}

	provider, err := h.getProvider(&account)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	snapshot, err := provider.CreateSnapshot(c.Request.Context(), disk.DiskID, req.Name)
	if err != nil {
		h.logger.Error("failed to create cloud snapshot", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	cloudSnapshot := &model.CloudSnapshot{
		CloudAccountID: disk.CloudAccountID,
		SnapshotID:     snapshot.ID,
		Name:           snapshot.Name,
		DiskID:         disk.ID,
		Size:           snapshot.Size,
		Status:         snapshot.Status,
		ProviderType:   account.ProviderType,
	}

	if err := h.db.Create(cloudSnapshot).Error; err != nil {
		h.logger.Error("failed to save cloud snapshot", zap.Error(err))
	}

	c.JSON(http.StatusCreated, cloudSnapshot)
}

// DeleteCloudSnapshot 删除云快照
func (h *StorageHandler) DeleteCloudSnapshot(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var snapshot model.CloudSnapshot
	if err := h.db.First(&snapshot, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "snapshot not found"})
		return
	}

	var account model.CloudAccount
	if err := h.db.First(&account, snapshot.CloudAccountID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cloud account not found"})
		return
	}

	provider, err := h.getProvider(&account)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = provider.DeleteSnapshot(c.Request.Context(), snapshot.SnapshotID)
	if err != nil {
		h.logger.Error("failed to delete cloud snapshot", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.db.Delete(&snapshot)

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// SyncCloudDisks 同步云硬盘
func (h *StorageHandler) SyncCloudDisks(c *gin.Context) {
	cloudAccountID, _ := strconv.ParseUint(c.Query("cloud_account_id"), 10, 32)

	var account model.CloudAccount
	if err := h.db.First(&account, cloudAccountID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cloud account not found"})
		return
	}

	provider, err := h.getProvider(&account)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	disks, err := provider.ListDisks(c.Request.Context(), cloudprovider.DiskFilter{})
	if err != nil {
		h.logger.Error("failed to list disks from provider", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	count := 0
	for _, d := range disks {
		var existingDisk model.CloudDisk
		err := h.db.Where("disk_id = ? AND cloud_account_id = ?", d.ID, cloudAccountID).First(&existingDisk).Error

		if err == gorm.ErrRecordNotFound {
			// 创建新记录
			newDisk := &model.CloudDisk{
				CloudAccountID: uint(cloudAccountID),
				DiskID:         d.ID,
				Name:           d.Name,
				Size:           d.Size,
				Type:           d.Type,
				Status:         d.Status,
				VMID:           d.VMID,
				ZoneID:         d.ZoneID,
				ProviderType:   account.ProviderType,
			}
			if err := h.db.Create(newDisk).Error; err == nil {
				count++
			}
		} else {
			// 更新现有记录
			existingDisk.Name = d.Name
			existingDisk.Size = d.Size
			existingDisk.Status = d.Status
			existingDisk.VMID = d.VMID
			h.db.Save(&existingDisk)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "sync completed",
		"count":   count,
		"total":   len(disks),
	})
}

// getProvider 获取云厂商适配器
func (h *StorageHandler) getProvider(account *model.CloudAccount) (cloudprovider.ICloudProvider, error) {
	return getCloudProvider(account)
}

// getCloudProvider 辅助函数获取云厂商适配器
func getCloudProvider(account *model.CloudAccount) (cloudprovider.ICloudProvider, error) {
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

	return cloudprovider.GetProvider(account.ProviderType, config)
}