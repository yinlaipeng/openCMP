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

// ProjectHandler 项目 Handler
type ProjectHandler struct {
	service *service.ProjectService
	db      *gorm.DB
	logger  *zap.Logger
}

// NewProjectHandler 创建项目 Handler
func NewProjectHandler(db *gorm.DB, logger *zap.Logger) *ProjectHandler {
	return &ProjectHandler{
		service: service.NewProjectService(db),
		db:      db,
		logger:  logger,
	}
}

// ProjectResponse 项目响应结构（带统计字段）
type ProjectResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	DomainID    uint   `json:"domain_id"`
	ManagerID   *uint  `json:"manager_id"`
	Enabled     bool   `json:"enabled"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	// 统计字段
	UserCount  int  `json:"user_count"`
	GroupCount int  `json:"group_count"`
	IsSystem   bool `json:"is_system"`
	CanDelete  bool `json:"can_delete"`
	CanUpdate  bool `json:"can_update"`
	Admin      string `json:"admin"`
	AdminID    string `json:"admin_id"`
}

// CreateProjectRequest 创建项目请求
type CreateProjectRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	DomainID    uint   `json:"domain_id" binding:"required"`
	ManagerID   *uint  `json:"manager_id"` // 项目管理员ID
	ParentID    *uint  `json:"parent_id"`
}

// JoinProjectRequest 加入项目请求
type JoinProjectRequest struct {
	Users []uint `json:"users"`
	Roles []uint `json:"roles"`
}

// List 列出项目
func (h *ProjectHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	details := c.Query("details") == "true"

	// 获取筛选参数
	var domainID *uint
	if d := c.Query("domain_id"); d != "" {
		id, _ := strconv.ParseUint(d, 10, 32)
		uid := uint(id)
		domainID = &uid
	}

	keyword := c.Query("keyword")
	name := c.Query("name")
	enabledStr := c.Query("enabled")

	var enabled *bool
	if enabledStr != "" {
		if enabledStr == "true" {
			enabled = func() *bool { b := true; return &b }()
		} else if enabledStr == "false" {
			enabled = func() *bool { b := false; return &b }()
		}
	}

	// 使用 name 或 keyword
	if name != "" && keyword == "" {
		keyword = name
	}

	projects, total, err := h.service.ListProjects(c.Request.Context(), domainID, keyword, enabled, limit, offset)
	if err != nil {
		h.logger.Error("failed to list projects", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 如果需要详细信息，添加统计字段
	if details {
		var responses []ProjectResponse
		for _, p := range projects {
			resp := h.buildProjectResponse(p)
			responses = append(responses, resp)
		}
		c.JSON(http.StatusOK, gin.H{
			"items": responses,
			"total": total,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": projects,
		"total": total,
	})
}

// buildProjectResponse 构建项目响应（带统计字段）
func (h *ProjectHandler) buildProjectResponse(p *model.Project) ProjectResponse {
	resp := ProjectResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		DomainID:    p.DomainID,
		ManagerID:   p.ManagerID,
		Enabled:     p.Enabled,
		CreatedAt:   p.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   p.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		IsSystem:    p.Name == "system",
		CanDelete:   p.Name != "system",
		CanUpdate:   true,
	}

	// 计算统计字段
	resp.UserCount = h.countProjectUsers(p.ID)
	resp.GroupCount = h.countProjectGroups(p.ID)

	// 获取管理员信息
	if p.ManagerID != nil {
		var user model.User
		if err := h.db.First(&user, *p.ManagerID).Error; err == nil {
			resp.Admin = user.Name
			resp.AdminID = strconv.FormatUint(uint64(user.ID), 10)
		}
	}

	return resp
}

// countProjectUsers 计算项目用户数量
func (h *ProjectHandler) countProjectUsers(projectID uint) int {
	var count int64
	// 通过 user_projects 关联表计算
	h.db.Table("user_projects").Where("project_id = ?", projectID).Count(&count)
	return int(count)
}

// countProjectGroups 计算项目组数量
func (h *ProjectHandler) countProjectGroups(projectID uint) int {
	var count int64
	// 通过 group_projects 关联表计算
	h.db.Table("group_projects").Where("project_id = ?", projectID).Count(&count)
	return int(count)
}

// Get 获取项目详情
func (h *ProjectHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	project, err := h.service.GetProject(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get project", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if project == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}

	c.JSON(http.StatusOK, project)
}

// Create 创建项目
func (h *ProjectHandler) Create(c *gin.Context) {
	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project := &model.Project{
		Name:        req.Name,
		Description: req.Description,
		DomainID:    req.DomainID,
		ManagerID:   req.ManagerID,
		ParentID:    req.ParentID,
		Enabled:     true,
	}

	if err := h.service.CreateProject(c.Request.Context(), project); err != nil {
		h.logger.Error("failed to create project", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, project)
}

// Update 更新项目
func (h *ProjectHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	project, err := h.service.GetProject(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get project", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if project == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}

	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project.Name = req.Name
	project.Description = req.Description
	project.DomainID = req.DomainID
	project.ManagerID = req.ManagerID
	project.ParentID = req.ParentID

	if err := h.service.UpdateProject(c.Request.Context(), project); err != nil {
		h.logger.Error("failed to update project", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, project)
}

// Delete 删除项目
func (h *ProjectHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.DeleteProject(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("failed to delete project", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// Enable 启用项目
func (h *ProjectHandler) Enable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.EnableProject(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("failed to enable project", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "enabled"})
}

// Disable 禁用项目
func (h *ProjectHandler) Disable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.DisableProject(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("failed to disable project", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "disabled"})
}

// Join 加入项目（分配用户和角色）
func (h *ProjectHandler) Join(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req JoinProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 为每个用户分配每个角色
	for _, userID := range req.Users {
		for _, roleID := range req.Roles {
			if err := h.service.AddUserToProject(c.Request.Context(), uint(id), userID, roleID); err != nil {
				h.logger.Error("failed to add user to project", zap.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "joined"})
}

// GetUsers 获取项目的用户列表
func (h *ProjectHandler) GetUsers(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	users, total, err := h.service.GetProjectUsers(c.Request.Context(), uint(id), limit, offset)
	if err != nil {
		h.logger.Error("failed to get project users", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": users,
		"total": total,
	})
}

// GetRoles 获取项目的角色列表
func (h *ProjectHandler) GetRoles(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	roles, total, err := h.service.GetProjectRoles(c.Request.Context(), uint(id), limit, offset)
	if err != nil {
		h.logger.Error("failed to get project roles", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": roles,
		"total": total,
	})
}

// RemoveUser 从项目移除用户
func (h *ProjectHandler) RemoveUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	userID := c.Query("user_id")
	roleID := c.Query("role_id")

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	uid64, _ := strconv.ParseUint(userID, 10, 32)
	uid := uint(uid64)
	var rid uint = 0
	if roleID != "" {
		rid64, _ := strconv.ParseUint(roleID, 10, 32)
		rid = uint(rid64)
	}

	if err := h.service.RemoveUserFromProject(c.Request.Context(), uint(id), uid, rid); err != nil {
		h.logger.Error("failed to remove user from project", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "removed"})
}

// SetManagerRequest 设置项目管理员请求
type SetManagerRequest struct {
	UserID uint `json:"user_id" binding:"required"`
}

// SetManager 设置项目管理员
func (h *ProjectHandler) SetManager(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req SetManagerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.SetProjectManager(c.Request.Context(), uint(id), req.UserID); err != nil {
		h.logger.Error("failed to set project manager", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "project manager set successfully"})
}
