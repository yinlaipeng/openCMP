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

// CloudAccountResourcesHandler 云账户资源 Handler
type CloudAccountResourcesHandler struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewCloudAccountResourcesHandler 创建云账户资源 Handler
func NewCloudAccountResourcesHandler(db *gorm.DB, logger *zap.Logger) *CloudAccountResourcesHandler {
	return &CloudAccountResourcesHandler{db: db, logger: logger}
}

// GetResourceStats 获取资源统计
func (h *CloudAccountResourcesHandler) GetResourceStats(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	// 获取云账户
	var account model.CloudAccount
	if err := h.db.First(&account, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
		return
	}

	// 从同步后的数据库获取资源统计
	stats := h.getSyncedResourceStats(uint(id))

	c.JSON(http.StatusOK, stats)
}

// getSyncedResourceStats 从数据库获取同步后的资源统计
func (h *CloudAccountResourcesHandler) getSyncedResourceStats(cloudAccountID uint) map[string]interface{} {
	// 统计各类资源数量
	var vmCount, vpcCount, subnetCount, sgCount, eipCount, imageCount, diskCount, snapshotCount, rdsCount, redisCount int64

	h.db.Model(&model.CloudVM{}).Where("cloud_account_id = ? AND status != ?", cloudAccountID, "terminated").Count(&vmCount)
	h.db.Model(&model.CloudVPC{}).Where("cloud_account_id = ? AND status != ?", cloudAccountID, "terminated").Count(&vpcCount)
	h.db.Model(&model.CloudSubnet{}).Where("cloud_account_id = ? AND status != ?", cloudAccountID, "terminated").Count(&subnetCount)
	h.db.Model(&model.CloudSecurityGroup{}).Where("cloud_account_id = ? AND status != ?", cloudAccountID, "terminated").Count(&sgCount)
	h.db.Model(&model.CloudEIP{}).Where("cloud_account_id = ? AND status != ?", cloudAccountID, "terminated").Count(&eipCount)
	h.db.Model(&model.CloudImage{}).Where("cloud_account_id = ?", cloudAccountID).Count(&imageCount)
	h.db.Model(&model.CloudDisk{}).Where("cloud_account_id = ? AND status != ?", cloudAccountID, "terminated").Count(&diskCount)
	h.db.Model(&model.CloudSnapshot{}).Where("cloud_account_id = ?", cloudAccountID).Count(&snapshotCount)
	h.db.Model(&model.CloudRDS{}).Where("cloud_account_id = ? AND status != ?", cloudAccountID, "terminated").Count(&rdsCount)
	h.db.Model(&model.CloudRedis{}).Where("cloud_account_id = ? AND status != ?", cloudAccountID, "terminated").Count(&redisCount)

	// 统计运行中的虚拟机
	var vmRunningCount int64
	h.db.Model(&model.CloudVM{}).Where("cloud_account_id = ? AND status = ?", cloudAccountID, "running").Count(&vmRunningCount)

	// 统计已挂载的磁盘
	var diskMountedCount int64
	h.db.Model(&model.CloudDisk{}).Where("cloud_account_id = ? AND vm_id IS NOT NULL AND vm_id != ''", cloudAccountID).Count(&diskMountedCount)

	// 统计已绑定的EIP
	var eipBoundCount int64
	h.db.Model(&model.CloudEIP{}).Where("cloud_account_id = ? AND resource_id IS NOT NULL AND resource_id != ''", cloudAccountID).Count(&eipBoundCount)

	return map[string]interface{}{
		"resources": map[string]int64{
			"vms":             vmCount,
			"vms_running":     vmRunningCount,
			"vpcs":            vpcCount,
			"subnets":         subnetCount,
			"security_groups": sgCount,
			"eips":            eipCount,
			"eips_bound":      eipBoundCount,
			"images":          imageCount,
			"disks":           diskCount,
			"disks_mounted":   diskMountedCount,
			"snapshots":       snapshotCount,
			"rds":             rdsCount,
			"redis":           redisCount,
		},
		"usage_rates": map[string]float64{
			"vm_running_rate":   float64(vmRunningCount) / float64(max(vmCount, 1)) * 100,
			"disk_mounted_rate": float64(diskMountedCount) / float64(max(diskCount, 1)) * 100,
			"eip_bound_rate":    float64(eipBoundCount) / float64(max(eipCount, 1)) * 100,
		},
	}
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

// GetPermissions 获取权限列表
func (h *CloudAccountResourcesHandler) GetPermissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	// 获取云账户
	var account model.CloudAccount
	if err := h.db.First(&account, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
		return
	}

	// 获取权限（模拟数据，后续可调用云厂商API）
	permissions := []map[string]interface{}{
		{"name": "AliyunECSFullAccess", "description": "云服务器ECS完全管理权限", "granted": true},
		{"name": "AliyunVPCFullAccess", "description": "专有网络VPC完全管理权限", "granted": true},
		{"name": "AliyunSLBFullAccess", "description": "负载均衡SLB完全管理权限", "granted": true},
		{"name": "AliyunRDSFullAccess", "description": "云数据库RDS完全管理权限", "granted": false},
		{"name": "AliyunOSSFullAccess", "description": "对象存储OSS完全管理权限", "granted": true},
	}

	c.JSON(http.StatusOK, gin.H{"permissions": permissions})
}

// ListSubscriptions 列出订阅
func (h *CloudAccountResourcesHandler) ListSubscriptions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var subscriptions []model.CloudSubscription
	if err := h.db.Where("cloud_account_id = ?", id).Find(&subscriptions).Error; err != nil {
		h.logger.Error("failed to list subscriptions", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": subscriptions, "total": len(subscriptions)})
}

// CreateSubscription 创建订阅
func (h *CloudAccountResourcesHandler) CreateSubscription(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		Name             string `json:"name" binding:"required"`
		SubscriptionID   string `json:"subscription_id" binding:"required"`
		DomainID         uint   `json:"domain_id"`
		DefaultProjectID *uint  `json:"default_project_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subscription := &model.CloudSubscription{
		CloudAccountID:   uint(id),
		Name:             req.Name,
		SubscriptionID:   req.SubscriptionID,
		DomainID:         req.DomainID,
		DefaultProjectID: req.DefaultProjectID,
		Enabled:          true,
		Status:           "normal",
		SyncStatus:       "completed",
	}

	if err := h.db.Create(subscription).Error; err != nil {
		h.logger.Error("failed to create subscription", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, subscription)
}

// ListCloudUsers 列出云用户
func (h *CloudAccountResourcesHandler) ListCloudUsers(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var users []model.CloudUser
	if err := h.db.Where("cloud_account_id = ?", id).Find(&users).Error; err != nil {
		h.logger.Error("failed to list cloud users", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": users, "total": len(users)})
}

// CreateCloudUser 创建云用户
func (h *CloudAccountResourcesHandler) CreateCloudUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		Username     string `json:"username" binding:"required"`
		ConsoleLogin bool   `json:"console_login"`
		Status       string `json:"status"`
		Password     string `json:"password"`
		LoginURL     string `json:"login_url"`
		LocalUserID  *uint  `json:"local_user_id"`
		Platform     string `json:"platform"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &model.CloudUser{
		CloudAccountID: uint(id),
		Username:       req.Username,
		ConsoleLogin:   req.ConsoleLogin,
		Status:         req.Status,
		Password:       req.Password,
		LoginURL:       req.LoginURL,
		LocalUserID:    req.LocalUserID,
		Platform:       req.Platform,
	}

	if user.Status == "" {
		user.Status = "normal"
	}

	if err := h.db.Create(user).Error; err != nil {
		h.logger.Error("failed to create cloud user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// UpdateCloudUser 更新云用户
func (h *CloudAccountResourcesHandler) UpdateCloudUser(c *gin.Context) {
	accountId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account id"})
		return
	}

	userId, err := strconv.ParseUint(c.Param("uid"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var user model.CloudUser
	if err := h.db.Where("id = ? AND cloud_account_id = ?", userId, accountId).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	var req struct {
		Username     string `json:"username"`
		ConsoleLogin *bool  `json:"console_login"`
		Status       string `json:"status"`
		Password     string `json:"password"`
		LoginURL     string `json:"login_url"`
		LocalUserID  *uint  `json:"local_user_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Username != "" {
		user.Username = req.Username
	}
	if req.ConsoleLogin != nil {
		user.ConsoleLogin = *req.ConsoleLogin
	}
	if req.Status != "" {
		user.Status = req.Status
	}
	if req.Password != "" {
		user.Password = req.Password
	}
	if req.LoginURL != "" {
		user.LoginURL = req.LoginURL
	}
	if req.LocalUserID != nil {
		user.LocalUserID = req.LocalUserID
	}

	if err := h.db.Save(&user).Error; err != nil {
		h.logger.Error("failed to update cloud user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteCloudUser 删除云用户
func (h *CloudAccountResourcesHandler) DeleteCloudUser(c *gin.Context) {
	accountId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account id"})
		return
	}

	userId, err := strconv.ParseUint(c.Param("uid"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	if err := h.db.Where("id = ? AND cloud_account_id = ?", userId, accountId).Delete(&model.CloudUser{}).Error; err != nil {
		h.logger.Error("failed to delete cloud user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// ListCloudUserGroups 列出云用户组
func (h *CloudAccountResourcesHandler) ListCloudUserGroups(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var groups []model.CloudUserGroup
	if err := h.db.Where("cloud_account_id = ?", id).Find(&groups).Error; err != nil {
		h.logger.Error("failed to list cloud user groups", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": groups, "total": len(groups)})
}

// CreateCloudUserGroup 创建云用户组
func (h *CloudAccountResourcesHandler) CreateCloudUserGroup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		Name       string `json:"name" binding:"required"`
		Status     string `json:"status"`
		Permissions string `json:"permissions"`
		Platform   string `json:"platform"`
		DomainID   uint   `json:"domain_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	group := &model.CloudUserGroup{
		CloudAccountID: uint(id),
		Name:           req.Name,
		Status:         req.Status,
		Permissions:    req.Permissions,
		Platform:       req.Platform,
		DomainID:       req.DomainID,
	}

	if group.Status == "" {
		group.Status = "normal"
	}

	if err := h.db.Create(group).Error; err != nil {
		h.logger.Error("failed to create cloud user group", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, group)
}

// UpdateCloudUserGroup 更新云用户组
func (h *CloudAccountResourcesHandler) UpdateCloudUserGroup(c *gin.Context) {
	accountId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account id"})
		return
	}

	groupId, err := strconv.ParseUint(c.Param("gid"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}

	var group model.CloudUserGroup
	if err := h.db.Where("id = ? AND cloud_account_id = ?", groupId, accountId).First(&group).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
		return
	}

	var req struct {
		Name        string `json:"name"`
		Status      string `json:"status"`
		Permissions string `json:"permissions"`
		DomainID    *uint  `json:"domain_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != "" {
		group.Name = req.Name
	}
	if req.Status != "" {
		group.Status = req.Status
	}
	if req.Permissions != "" {
		group.Permissions = req.Permissions
	}
	if req.DomainID != nil {
		group.DomainID = *req.DomainID
	}

	if err := h.db.Save(&group).Error; err != nil {
		h.logger.Error("failed to update cloud user group", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, group)
}

// DeleteCloudUserGroup 删除云用户组
func (h *CloudAccountResourcesHandler) DeleteCloudUserGroup(c *gin.Context) {
	accountId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account id"})
		return
	}

	groupId, err := strconv.ParseUint(c.Param("gid"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}

	if err := h.db.Where("id = ? AND cloud_account_id = ?", groupId, accountId).Delete(&model.CloudUserGroup{}).Error; err != nil {
		h.logger.Error("failed to delete cloud user group", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// ListCloudProjects 列出云上项目
func (h *CloudAccountResourcesHandler) ListCloudProjects(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var projects []model.CloudProject
	if err := h.db.Where("cloud_account_id = ?", id).Find(&projects).Error; err != nil {
		h.logger.Error("failed to list cloud projects", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": projects, "total": len(projects)})
}

// CreateCloudProject 创建云上项目
func (h *CloudAccountResourcesHandler) CreateCloudProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		Name           string `json:"name" binding:"required"`
		SubscriptionID *uint  `json:"subscription_id"`
		Status         string `json:"status"`
		Tags           string `json:"tags"`
		DomainID       uint   `json:"domain_id"`
		LocalProjectID *uint  `json:"local_project_id"`
		Priority       int    `json:"priority"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project := &model.CloudProject{
		CloudAccountID: uint(id),
		Name:           req.Name,
		SubscriptionID: req.SubscriptionID,
		Status:         req.Status,
		Tags:           req.Tags,
		DomainID:       req.DomainID,
		LocalProjectID: req.LocalProjectID,
		Priority:       req.Priority,
	}

	if project.Status == "" {
		project.Status = "normal"
	}

	if err := h.db.Create(project).Error; err != nil {
		h.logger.Error("failed to create cloud project", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, project)
}

// UpdateCloudProject 更新云上项目
func (h *CloudAccountResourcesHandler) UpdateCloudProject(c *gin.Context) {
	accountId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account id"})
		return
	}

	projectId, err := strconv.ParseUint(c.Param("pid"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	var project model.CloudProject
	if err := h.db.Where("id = ? AND cloud_account_id = ?", projectId, accountId).First(&project).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}

	var req struct {
		Name           string `json:"name"`
		Status         string `json:"status"`
		Tags           string `json:"tags"`
		LocalProjectID *uint  `json:"local_project_id"`
		Priority       *int   `json:"priority"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != "" {
		project.Name = req.Name
	}
	if req.Status != "" {
		project.Status = req.Status
	}
	if req.Tags != "" {
		project.Tags = req.Tags
	}
	if req.LocalProjectID != nil {
		project.LocalProjectID = req.LocalProjectID
	}
	if req.Priority != nil {
		project.Priority = *req.Priority
	}

	if err := h.db.Save(&project).Error; err != nil {
		h.logger.Error("failed to update cloud project", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, project)
}

// DeleteCloudProject 删除云上项目
func (h *CloudAccountResourcesHandler) DeleteCloudProject(c *gin.Context) {
	accountId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account id"})
		return
	}

	projectId, err := strconv.ParseUint(c.Param("pid"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	if err := h.db.Where("id = ? AND cloud_account_id = ?", projectId, accountId).Delete(&model.CloudProject{}).Error; err != nil {
		h.logger.Error("failed to delete cloud project", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// MapCloudProjectToLocal 映射云上项目到本地项目
func (h *CloudAccountResourcesHandler) MapCloudProjectToLocal(c *gin.Context) {
	accountId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account id"})
		return
	}

	projectId, err := strconv.ParseUint(c.Param("pid"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	var project model.CloudProject
	if err := h.db.Where("id = ? AND cloud_account_id = ?", projectId, accountId).First(&project).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}

	var req struct {
		LocalProjectID uint `json:"local_project_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project.LocalProjectID = &req.LocalProjectID
	if err := h.db.Save(&project).Error; err != nil {
		h.logger.Error("failed to map cloud project to local", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "mapped successfully", "local_project_id": req.LocalProjectID})
}

// UpdateSubscription 更新订阅
func (h *CloudAccountResourcesHandler) UpdateSubscription(c *gin.Context) {
	accountId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account id"})
		return
	}

	subId, err := strconv.ParseUint(c.Param("sid"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription id"})
		return
	}

	var subscription model.CloudSubscription
	if err := h.db.Where("id = ? AND cloud_account_id = ?", subId, accountId).First(&subscription).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
		return
	}

	var req struct {
		Name             string `json:"name"`
		SubscriptionID   string `json:"subscription_id"`
		Enabled          *bool  `json:"enabled"`
		DefaultProjectID *uint  `json:"default_project_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != "" {
		subscription.Name = req.Name
	}
	if req.SubscriptionID != "" {
		subscription.SubscriptionID = req.SubscriptionID
	}
	if req.Enabled != nil {
		subscription.Enabled = *req.Enabled
	}
	if req.DefaultProjectID != nil {
		subscription.DefaultProjectID = req.DefaultProjectID
	}

	if err := h.db.Save(&subscription).Error; err != nil {
		h.logger.Error("failed to update subscription", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subscription)
}

// DeleteSubscription 删除订阅
func (h *CloudAccountResourcesHandler) DeleteSubscription(c *gin.Context) {
	accountId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account id"})
		return
	}

	subId, err := strconv.ParseUint(c.Param("sid"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription id"})
		return
	}

	if err := h.db.Where("id = ? AND cloud_account_id = ?", subId, accountId).Delete(&model.CloudSubscription{}).Error; err != nil {
		h.logger.Error("failed to delete subscription", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// ToggleSubscription 切换订阅启用状态
func (h *CloudAccountResourcesHandler) ToggleSubscription(c *gin.Context) {
	accountId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account id"})
		return
	}

	subId, err := strconv.ParseUint(c.Param("sid"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription id"})
		return
	}

	var subscription model.CloudSubscription
	if err := h.db.Where("id = ? AND cloud_account_id = ?", subId, accountId).First(&subscription).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
		return
	}

	var req struct {
		Enabled bool `json:"enabled" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subscription.Enabled = req.Enabled
	if err := h.db.Save(&subscription).Error; err != nil {
		h.logger.Error("failed to toggle subscription", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "status updated", "enabled": req.Enabled})
}

// SyncSubscription 同步订阅
func (h *CloudAccountResourcesHandler) SyncSubscription(c *gin.Context) {
	accountId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account id"})
		return
	}

	subId, err := strconv.ParseUint(c.Param("sid"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription id"})
		return
	}

	var subscription model.CloudSubscription
	if err := h.db.Where("id = ? AND cloud_account_id = ?", subId, accountId).First(&subscription).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
		return
	}

	// TODO: 实现真实的同步逻辑
	now := time.Now()
	subscription.SyncTime = &now
	subscription.SyncStatus = "completed"
	subscription.SyncDuration = 5

	if err := h.db.Save(&subscription).Error; err != nil {
		h.logger.Error("failed to sync subscription", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "sync completed", "subscription": subscription})
}

// UpdateSubscriptionProject 更新订阅的默认项目
func (h *CloudAccountResourcesHandler) UpdateSubscriptionProject(c *gin.Context) {
	accountId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account id"})
		return
	}

	subId, err := strconv.ParseUint(c.Param("sid"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription id"})
		return
	}

	var subscription model.CloudSubscription
	if err := h.db.Where("id = ? AND cloud_account_id = ?", subId, accountId).First(&subscription).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
		return
	}

	var req struct {
		ProjectID uint `json:"project_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subscription.DefaultProjectID = &req.ProjectID
	if err := h.db.Save(&subscription).Error; err != nil {
		h.logger.Error("failed to update subscription project", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "project updated", "project_id": req.ProjectID})
}

// ListOperationLogs 列出操作日志
func (h *CloudAccountResourcesHandler) ListOperationLogs(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize

	var logs []model.OperationLog
	var total int64

	h.db.Model(&model.OperationLog{}).Where("cloud_account_id = ?", id).Count(&total)
	if err := h.db.Where("cloud_account_id = ?", id).Order("operation_time DESC").Limit(pageSize).Offset(offset).Find(&logs).Error; err != nil {
		h.logger.Error("failed to list operation logs", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": logs, "total": total, "page": page, "page_size": pageSize})
}

// 辅助函数
func getDefaultResources() map[string]int {
	return map[string]int{
		"vms": 0, "rds": 0, "redis": 0, "buckets": 0,
		"eips": 0, "public_ips": 0, "snapshots": 0,
		"vpcs": 0, "subnets": 0, "total_ips": 0,
		"vms_running": 0, "eips_bound": 0, "ips_used": 0,
		"disks": 0, "disks_mounted": 0,
	}
}

func getDefaultUsageRates() map[string]int {
	return map[string]int{
		"vm_running_rate":  0,
		"disk_mounted_rate": 0,
		"eip_bound_rate":   0,
		"ip_used_rate":     0,
	}
}