package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
)

// ImageHandler 镜像 Handler
type ImageHandler struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewImageHandler 创建镜像 Handler
func NewImageHandler(db *gorm.DB, logger *zap.Logger) *ImageHandler {
	return &ImageHandler{
		db:     db,
		logger: logger,
	}
}

// ImageResponse 镜像响应结构
type ImageResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Format      string `json:"format"`
	OsName      string `json:"os_name"`
	OsVersion   string `json:"os_version"`
	Size        int64  `json:"size"`
	CpuArch     string `json:"cpu_arch"`
	ImageType   string `json:"image_type"`
	ShareScope  string `json:"share_scope"`
	Platform    string `json:"platform"`
	AccountName string `json:"account_name"`
	RegionID    string `json:"region_id"`
	ProjectID   string `json:"project_id"`
	ProjectName string `json:"project_name"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// ImageListParams 镜像列表查询参数
type ImageListParams struct {
	Page         int    `form:"page"`
	PageSize     int    `form:"page_size"`
	Details      bool   `form:"details"`
	IsGuestImage bool   `form:"is_guest_image"`
	ImageType    string `form:"image_type"`
	Name         string `form:"name"`
	OsName       string `form:"os_name"`
	Platform     string `form:"platform"`
	Format       string `form:"format"`
	Status       string `form:"status"`
	CpuArch      string `form:"cpu_arch"`
	AccountID    uint   `form:"account_id"`
	RegionID     string `form:"region_id"`
	ProjectID    string `form:"project_id"`
}

// ListImages 获取镜像列表
func (h *ImageHandler) ListImages(c *gin.Context) {
	var params ImageListParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 默认分页参数
	if params.Page == 0 {
		params.Page = 1
	}
	if params.PageSize == 0 {
		params.PageSize = 20
	}

	// 查询数据库
	var images []model.Image
	var total int64

	query := h.db.Model(&model.Image{})

	// 应用过滤条件
	if params.Name != "" {
		query = query.Where("name LIKE ?", "%"+params.Name+"%")
	}
	if params.OsName != "" {
		query = query.Where("os_name = ?", params.OsName)
	}
	if params.Platform != "" {
		query = query.Where("platform = ?", params.Platform)
	}
	if params.Format != "" {
		query = query.Where("format = ?", params.Format)
	}
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.CpuArch != "" {
		query = query.Where("cpu_arch = ?", params.CpuArch)
	}
	if params.ImageType != "" {
		query = query.Where("image_type = ?", params.ImageType)
	}
	if params.AccountID > 0 {
		query = query.Where("cloud_account_id = ?", params.AccountID)
	}
	if params.RegionID != "" {
		query = query.Where("region_id = ?", params.RegionID)
	}
	if params.ProjectID != "" {
		query = query.Where("project_id = ?", params.ProjectID)
	}

	// 计算总数
	query.Count(&total)

	// 分页查询
	offset := (params.Page - 1) * params.PageSize
	if err := query.Offset(offset).Limit(params.PageSize).Find(&images).Error; err != nil {
		h.logger.Error("failed to list images", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取镜像列表失败"})
		return
	}

	// 转换为响应格式
	var responses []ImageResponse
	for _, img := range images {
		// 获取云账号名称
		var accountName string
		if img.CloudAccountID > 0 {
			var account model.CloudAccount
			if err := h.db.Select("name").First(&account, img.CloudAccountID).Error; err == nil {
				accountName = account.Name
			}
		}

		responses = append(responses, ImageResponse{
			ID:          img.ID,
			Name:        img.Name,
			Description: img.Description,
			Status:      img.Status,
			Format:      img.Format,
			OsName:      img.OsName,
			OsVersion:   img.OsVersion,
			Size:        img.Size,
			CpuArch:     img.CpuArch,
			ImageType:   img.ImageType,
			ShareScope:  img.ShareScope,
			Platform:    img.Platform,
			AccountName: accountName,
			RegionID:    img.RegionID,
			ProjectID:   img.ProjectID,
			CreatedAt:   img.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   img.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"items": responses,
		"total": total,
		"pagination": gin.H{
			"page":        params.Page,
			"page_size":   params.PageSize,
			"total":       total,
			"total_pages": (total + int64(params.PageSize) - 1) / int64(params.PageSize),
		},
	})
}

// GetImage 获取单个镜像详情
func (h *ImageHandler) GetImage(c *gin.Context) {
	imageID := c.Param("id")

	var image model.Image
	if err := h.db.First(&image, "id = ?", imageID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "镜像不存在"})
		} else {
			h.logger.Error("failed to get image", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取镜像详情失败"})
		}
		return
	}

	// 获取云账号名称
	var accountName string
	if image.CloudAccountID > 0 {
		var account model.CloudAccount
		if err := h.db.Select("name").First(&account, image.CloudAccountID).Error; err == nil {
			accountName = account.Name
		}
	}

	c.JSON(http.StatusOK, ImageResponse{
		ID:          image.ID,
		Name:        image.Name,
		Description: image.Description,
		Status:      image.Status,
		Format:      image.Format,
		OsName:      image.OsName,
		OsVersion:   image.OsVersion,
		Size:        image.Size,
		CpuArch:     image.CpuArch,
		ImageType:   image.ImageType,
		ShareScope:  image.ShareScope,
		Platform:    image.Platform,
		AccountName: accountName,
		RegionID:    image.RegionID,
		ProjectID:   image.ProjectID,
		CreatedAt:   image.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   image.UpdatedAt.Format("2006-01-02 15:04:05"),
	})
}

// CreateImageRequest 创建镜像请求
type CreateImageRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	OsName      string   `json:"os_name" binding:"required"`
	OsVersion   string   `json:"os_version"`
	Format      string   `json:"format" binding:"required"`
	CpuArch     string   `json:"cpu_arch"`
	ProjectID   string   `json:"project_id"`
	Tags        []string `json:"tags"`
}

// CreateImage 创建镜像
func (h *ImageHandler) CreateImage(c *gin.Context) {
	var req CreateImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 创建镜像记录
	image := model.Image{
		Name:        req.Name,
		Description: req.Description,
		Status:      "Creating",
		Format:      req.Format,
		OsName:      req.OsName,
		OsVersion:   req.OsVersion,
		CpuArch:     req.CpuArch,
		ImageType:   "system",
		ShareScope:  "private",
		ProjectID:   req.ProjectID,
	}

	if err := h.db.Create(&image).Error; err != nil {
		h.logger.Error("failed to create image", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建镜像失败"})
		return
	}

	c.JSON(http.StatusCreated, ImageResponse{
		ID:          image.ID,
		Name:        image.Name,
		Description: image.Description,
		Status:      image.Status,
		Format:      image.Format,
		OsName:      image.OsName,
		OsVersion:   image.OsVersion,
		Size:        image.Size,
		CpuArch:     image.CpuArch,
		ImageType:   image.ImageType,
		ShareScope:  image.ShareScope,
		Platform:    image.Platform,
		RegionID:    image.RegionID,
		ProjectID:   image.ProjectID,
		CreatedAt:   image.CreatedAt.Format("2006-01-02 15:04:05"),
	})
}

// UpdateImageRequest 更新镜像请求
type UpdateImageRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	OsName      string `json:"os_name"`
	ShareScope  string `json:"share_scope"`
}

// UpdateImage 更新镜像信息
func (h *ImageHandler) UpdateImage(c *gin.Context) {
	imageID := c.Param("id")

	var req UpdateImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var image model.Image
	if err := h.db.First(&image, "id = ?", imageID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "镜像不存在"})
		} else {
			h.logger.Error("failed to find image", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查找镜像失败"})
		}
		return
	}

	// 更新字段
	if req.Name != "" {
		image.Name = req.Name
	}
	if req.Description != "" {
		image.Description = req.Description
	}
	if req.OsName != "" {
		image.OsName = req.OsName
	}
	if req.ShareScope != "" {
		image.ShareScope = req.ShareScope
	}

	if err := h.db.Save(&image).Error; err != nil {
		h.logger.Error("failed to update image", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新镜像失败"})
		return
	}

	// 获取云账号名称
	var accountName string
	if image.CloudAccountID > 0 {
		var account model.CloudAccount
		if err := h.db.Select("name").First(&account, image.CloudAccountID).Error; err == nil {
			accountName = account.Name
		}
	}

	c.JSON(http.StatusOK, ImageResponse{
		ID:          image.ID,
		Name:        image.Name,
		Description: image.Description,
		Status:      image.Status,
		Format:      image.Format,
		OsName:      image.OsName,
		OsVersion:   image.OsVersion,
		Size:        image.Size,
		CpuArch:     image.CpuArch,
		ImageType:   image.ImageType,
		ShareScope:  image.ShareScope,
		Platform:    image.Platform,
		AccountName: accountName,
		RegionID:    image.RegionID,
		ProjectID:   image.ProjectID,
		UpdatedAt:   image.UpdatedAt.Format("2006-01-02 15:04:05"),
	})
}

// DeleteImage 删除镜像
func (h *ImageHandler) DeleteImage(c *gin.Context) {
	imageID := c.Param("id")

	var image model.Image
	if err := h.db.First(&image, "id = ?", imageID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "镜像不存在"})
		} else {
			h.logger.Error("failed to find image", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查找镜像失败"})
		}
		return
	}

	if err := h.db.Delete(&image).Error; err != nil {
		h.logger.Error("failed to delete image", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除镜像失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// ImageBatchDeleteRequest 批量删除镜像请求
type ImageBatchDeleteRequest struct {
	IDs []string `json:"ids" binding:"required"`
}

// BatchDeleteImages 批量删除镜像
func (h *ImageHandler) BatchDeleteImages(c *gin.Context) {
	var req ImageBatchDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	success := 0
	failed := 0

	for _, id := range req.IDs {
		if err := h.db.Delete(&model.Image{}, "id = ?", id).Error; err != nil {
			failed++
		} else {
			success++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"total":   len(req.IDs),
		"success": success,
		"failed":  failed,
		"message": "批量删除完成",
	})
}

// SyncImages 同步镜像
func (h *ImageHandler) SyncImages(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	if accountIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少 account_id 参数"})
		return
	}

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 account_id"})
		return
	}

	// TODO: 实现从云账号同步镜像的逻辑
	// 这里应该调用 cloudprovider 的接口获取镜像列表并更新数据库

	c.JSON(http.StatusOK, gin.H{
		"message":    "镜像同步任务已启动",
		"account_id": accountID,
	})
}