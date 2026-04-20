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

// GroupHandler 用户组 Handler
type GroupHandler struct {
	service *service.GroupService
	db      *gorm.DB
	logger  *zap.Logger
}

// NewGroupHandler 创建用户组 Handler
func NewGroupHandler(db *gorm.DB, logger *zap.Logger) *GroupHandler {
	return &GroupHandler{
		service: service.NewGroupService(db),
		db:      db,
		logger:  logger,
	}
}

// GroupResponse 用户组响应结构（带统计字段）
type GroupResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	DomainID    uint   `json:"domain_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	// 统计字段
	UserCount    int  `json:"user_count"`
	ProjectCount int  `json:"project_count"`
	IsSSO        bool `json:"is_sso"`
	CanDelete    bool `json:"can_delete"`
	CanUpdate    bool `json:"can_update"`
}

// CreateGroupRequest 创建用户组请求
type CreateGroupRequest struct {
	Name        string `json:"name" binding:"required"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	DomainID    uint   `json:"domain_id" binding:"required"`
}

// List 列出用户组
func (h *GroupHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	details := c.Query("details") == "true"

	var domainID *uint
	if d := c.Query("domain_id"); d != "" {
		id, _ := strconv.ParseUint(d, 10, 32)
		uid := uint(id)
		domainID = &uid
	}

	keyword := c.Query("keyword")

	groups, total, err := h.service.ListGroups(c.Request.Context(), domainID, keyword, limit, offset)
	if err != nil {
		h.logger.Error("failed to list groups", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 如果需要详细信息，添加统计字段
	if details {
		var responses []GroupResponse
		for _, g := range groups {
			resp := h.buildGroupResponse(g)
			responses = append(responses, resp)
		}
		c.JSON(http.StatusOK, gin.H{
			"items": responses,
			"total": total,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": groups,
		"total": total,
	})
}

// buildGroupResponse 构建用户组响应（带统计字段）
func (h *GroupHandler) buildGroupResponse(g *model.Group) GroupResponse {
	resp := GroupResponse{
		ID:          g.ID,
		Name:        g.Name,
		Description: g.Description,
		DomainID:    g.DomainID,
		CreatedAt:   g.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   g.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		IsSSO:       false,
		CanDelete:   true,
		CanUpdate:   true,
	}

	// 计算统计字段
	resp.UserCount = h.countGroupUsers(g.ID)
	resp.ProjectCount = h.countGroupProjects(g.ID)

	return resp
}

// countGroupUsers 计算用户组用户数量
func (h *GroupHandler) countGroupUsers(groupID uint) int {
	var count int64
	h.db.Table("user_groups").Where("group_id = ?", groupID).Count(&count)
	return int(count)
}

// countGroupProjects 计算用户组项目数量
func (h *GroupHandler) countGroupProjects(groupID uint) int {
	var count int64
	h.db.Table("group_projects").Where("group_id = ?", groupID).Count(&count)
	return int(count)
}

// Get 获取用户组详情
func (h *GroupHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	group, err := h.service.GetGroup(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get group", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if group == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
		return
	}

	c.JSON(http.StatusOK, group)
}

// Create 创建用户组
func (h *GroupHandler) Create(c *gin.Context) {
	var req CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	group := &model.Group{
		Name:        req.Name,
		Description: req.Description,
		DomainID:    req.DomainID,
	}

	if err := h.service.CreateGroup(c.Request.Context(), group); err != nil {
		h.logger.Error("failed to create group", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, group)
}

// Update 更新用户组
func (h *GroupHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	group, err := h.service.GetGroup(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get group", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if group == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
		return
	}

	var req CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	group.Description = req.Description

	if err := h.service.UpdateGroup(c.Request.Context(), group); err != nil {
		h.logger.Error("failed to update group", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, group)
}

// Delete 删除用户组
func (h *GroupHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.DeleteGroup(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("failed to delete group", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// GetUsers 获取用户组的用户列表
func (h *GroupHandler) GetUsers(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	users, total, err := h.service.GetGroupUsers(c.Request.Context(), uint(id), limit, offset)
	if err != nil {
		h.logger.Error("failed to get group users", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": users,
		"total": total,
	})
}

// AddUser 添加用户到用户组
func (h *GroupHandler) AddUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		UserID uint `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AddUserToGroup(c.Request.Context(), req.UserID, uint(id)); err != nil {
		h.logger.Error("failed to add user to group", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "added"})
}

// RemoveUser 从用户组移除用户
func (h *GroupHandler) RemoveUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	uid, _ := strconv.ParseUint(userID, 10, 32)

	if err := h.service.RemoveUserFromGroup(c.Request.Context(), uint(uid), uint(id)); err != nil {
		h.logger.Error("failed to remove user from group", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "removed"})
}

// AddProject 添加项目到用户组
func (h *GroupHandler) AddProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		ProjectID uint `json:"project_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AddGroupToProject(c.Request.Context(), uint(id), req.ProjectID); err != nil {
		h.logger.Error("failed to add group to project", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "added"})
}

// RemoveProject 从项目移除用户组
func (h *GroupHandler) RemoveProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	projectID := c.Query("project_id")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project_id is required"})
		return
	}

	pid, _ := strconv.ParseUint(projectID, 10, 32)

	if err := h.service.RemoveGroupFromProject(c.Request.Context(), uint(id), uint(pid)); err != nil {
		h.logger.Error("failed to remove group from project", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "removed"})
}

// GetProjects 获取用户组的项目列表
func (h *GroupHandler) GetProjects(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	projects, total, err := h.service.GetGroupProjects(c.Request.Context(), uint(id), limit, offset)
	if err != nil {
		h.logger.Error("failed to get group projects", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": projects,
		"total": total,
	})
}

// BatchDeleteRequest 批量删除请求
type BatchDeleteRequest struct {
	GroupIDs []uint `json:"group_ids" binding:"required"`
}

// BatchDelete 批量删除用户组
func (h *GroupHandler) BatchDelete(c *gin.Context) {
	var req BatchDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.GroupIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "group_ids is required"})
		return
	}

	if err := h.service.BatchDeleteGroups(c.Request.Context(), req.GroupIDs); err != nil {
		h.logger.Error("failed to batch delete groups", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "deleted",
		"count":   len(req.GroupIDs),
	})
}
