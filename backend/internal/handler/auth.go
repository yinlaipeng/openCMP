package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
	"github.com/opencmp/opencmp/internal/service"
	"github.com/opencmp/opencmp/pkg/utils"
)

// AuthHandler 认证 Handler
type AuthHandler struct {
	userService       *service.UserService
	authSourceService *service.AuthSourceService
	roleService       *service.RoleService
	domainService     *service.DomainService
	projectService    *service.ProjectService
	logger            *zap.Logger
	jwtSecret         string
	jwtExpire         int
	db                *gorm.DB
}

// NewAuthHandler 创建认证 Handler
func NewAuthHandler(db *gorm.DB, logger *zap.Logger, jwtSecret string, jwtExpire int) *AuthHandler {
	return &AuthHandler{
		userService:       service.NewUserService(db),
		authSourceService: service.NewAuthSourceService(db),
		roleService:       service.NewRoleService(db),
		domainService:     service.NewDomainService(db),
		projectService:    service.NewProjectService(db),
		logger:            logger,
		jwtSecret:         jwtSecret,
		jwtExpire:         jwtExpire,
		db:                db,
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string      `json:"token"`
	User  *model.User `json:"user"`
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()

	// 统一认证入口：先本地，后 LDAP
	user, authSource, err := h.authSourceService.AuthenticateUser(ctx, req.Username, req.Password, nil)
	if err != nil {
		// 区分用户禁用和密码错误
		if err.Error() == "user is disabled" {
			c.JSON(http.StatusForbidden, gin.H{"error": "user is disabled"})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	// 获取用户的角色
	roleIDs, err := h.userService.GetUserRoleIDs(ctx, user.ID)
	if err != nil {
		h.logger.Error("failed to get user roles", zap.Error(err))
		// 如果无法获取角色，继续登录但不包含角色信息
		roleIDs = []uint{}
	}

	// 更新最后登录信息
	h.userService.UpdateLastLogin(ctx, user.ID, c.ClientIP())

	// 生成 JWT token
	// 如果是通过认证源（LDAP 等）登录，在 token 中携带 auth_source_id 和 domain_id
	var token string
	if authSource != nil {
		token, err = utils.GenerateTokenWithExtra(user.ID, user.Name, roleIDs, user.DomainID, authSource.ID, h.jwtSecret, h.jwtExpire)
	} else {
		token, err = utils.GenerateToken(user.ID, user.Name, roleIDs, h.jwtSecret, h.jwtExpire)
	}
	if err != nil {
		h.logger.Error("failed to generate token", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token: token,
		User:  user,
	})
}

// Logout 用户登出
func (h *AuthHandler) Logout(c *gin.Context) {
	// TODO: 实现 token 黑名单机制
	c.JSON(http.StatusOK, gin.H{"message": "logout success"})
}

// GetCurrentUser 获取当前用户信息
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	// 从 token 中获取用户 ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	user, err := h.userService.GetUser(c.Request.Context(), userID.(uint))
	if err != nil {
		h.logger.Error("failed to get user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user info"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// ChangePassword 修改密码
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取当前用户信息
	user, err := h.userService.GetUser(c.Request.Context(), userID.(uint))
	if err != nil {
		h.logger.Error("failed to get user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user info"})
		return
	}

	// 验证旧密码
	if !checkPassword(user.Password, req.OldPassword) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "old password is incorrect"})
		return
	}

	err = h.userService.UpdatePassword(c.Request.Context(), userID.(uint), req.NewPassword)
	if err != nil {
		h.logger.Error("failed to change password", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to change password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password changed"})
}

// checkPassword 验证密码
func checkPassword(stored, input string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(stored), []byte(input))
	return err == nil
}

// UpdateProfile 更新个人信息
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	var req struct {
		DisplayName  string `json:"display_name"`
		Email        string `json:"email"`
		Phone        string `json:"phone"`
		Remark       string `json:"remark"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取当前用户信息
	user, err := h.userService.GetUser(c.Request.Context(), userID.(uint))
	if err != nil {
		h.logger.Error("failed to get user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user info"})
		return
	}

	// 更新用户信息
	if req.DisplayName != "" {
		user.DisplayName = req.DisplayName
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Remark != "" {
		user.Remark = req.Remark
	}

	err = h.userService.UpdateUser(c.Request.Context(), user)
	if err != nil {
		h.logger.Error("failed to update profile", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "profile updated", "user": user})
}

// GetPermissions 获取用户权限列表
func (h *AuthHandler) GetPermissions(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	ctx := c.Request.Context()

	// 获取用户的角色
	roleIDs, err := h.userService.GetUserRoleIDs(ctx, userID.(uint))
	if err != nil {
		h.logger.Error("failed to get user roles", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user roles"})
		return
	}

	// 通过角色获取权限
	var permissions []string
	err = h.db.Table("role_permissions").
		Select("permission_name").
		Where("role_id IN ?", roleIDs).
		Pluck("permission_name", &permissions).Error

	if err != nil {
		h.logger.Error("failed to get permissions", zap.Error(err))
		permissions = []string{}
	}

	// 去重
	uniquePerms := make(map[string]bool)
	for _, p := range permissions {
		uniquePerms[p] = true
	}

	result := make([]string, 0, len(uniquePerms))
	for p := range uniquePerms {
		result = append(result, p)
	}

	c.JSON(http.StatusOK, gin.H{"permissions": result})
}

// GetRegions 获取可用区域列表
func (h *AuthHandler) GetRegions(c *gin.Context) {
	// 返回云账号关联的区域列表
	var regions []struct {
		ID        uint   `json:"id"`
		Name      string `json:"name"`
		Provider  string `json:"provider"`
		AccountID uint   `json:"account_id"`
	}

	err := h.db.Table("cloud_accounts").
		Select("id, name, provider").
		Where("status = ?", "enabled").
		Find(&regions).Error

	if err != nil {
		h.logger.Error("failed to get regions", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get regions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"regions": regions})
}

// GetStats 获取认证统计信息
func (h *AuthHandler) GetStats(c *gin.Context) {
	stats := struct {
		UserCount    int64 `json:"user_count"`
		DomainCount  int64 `json:"domain_count"`
		ProjectCount int64 `json:"project_count"`
		GroupCount   int64 `json:"group_count"`
		RoleCount    int64 `json:"role_count"`
	}{}

	// 统计用户数
	h.db.Model(&model.User{}).Count(&stats.UserCount)

	// 统计域数
	h.db.Model(&model.Domain{}).Count(&stats.DomainCount)

	// 统计项目数
	h.db.Model(&model.Project{}).Count(&stats.ProjectCount)

	// 统计组数
	h.db.Model(&model.Group{}).Count(&stats.GroupCount)

	// 统计角色数
	h.db.Model(&model.Role{}).Count(&stats.RoleCount)

	c.JSON(http.StatusOK, stats)
}

// GetScopedResources 获取 scoped 资源
func (h *AuthHandler) GetScopedResources(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	ctx := c.Request.Context()

	// 获取用户所属的项目
	user, err := h.userService.GetUser(ctx, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
		return
	}

	// 返回用户可访问的域和项目
	resources := struct {
		DomainID   uint   `json:"domain_id"`
		DomainName string `json:"domain_name"`
		Projects   []struct {
			ID   uint   `json:"id"`
			Name string `json:"name"`
		} `json:"projects"`
	}{
		DomainID:   user.DomainID,
		DomainName: "",
		Projects:   []struct {
			ID   uint   `json:"id"`
			Name string `json:"name"`
		}{},
	}

	// 获取域名称
	if user.DomainID > 0 {
		domain, err := h.domainService.GetDomain(ctx, user.DomainID)
		if err == nil {
			resources.DomainName = domain.Name
		}
	}

	// 获取用户关联的项目
	var projectIDs []uint
	h.db.Table("project_members").
		Where("user_id = ?", userID).
		Pluck("project_id", &projectIDs)

	for _, pid := range projectIDs {
		project, err := h.projectService.GetProject(ctx, pid)
		if err == nil {
			resources.Projects = append(resources.Projects, struct {
				ID   uint   `json:"id"`
				Name string `json:"name"`
			}{ID: project.ID, Name: project.Name})
		}
	}

	c.JSON(http.StatusOK, resources)
}

// GetScopedPolicyBindings 获取策略绑定
func (h *AuthHandler) GetScopedPolicyBindings(c *gin.Context) {
	category := c.Query("category")

	// 返回策略绑定配置
	bindings := []struct {
		Category string `json:"category"`
		Key      string `json:"key"`
		Value    string `json:"value"`
	}{}

	// 简化实现：返回默认配置
	defaultBindings := []string{
		"sub_hidden_menus",
		"server_hidden_columns",
		"disk_hidden_columns",
		"dashboard_hidden_actions",
		"navbar_hidden_items",
	}

	for _, cat := range defaultBindings {
		if category == "" || category == cat {
			bindings = append(bindings, struct {
				Category string `json:"category"`
				Key      string `json:"key"`
				Value    string `json:"value"`
			}{
				Category: cat,
				Key:      "default",
				Value:    "",
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{"bindings": bindings})
}

// GetCapabilities 获取系统能力配置
func (h *AuthHandler) GetCapabilities(c *gin.Context) {
	capabilities := struct {
		Features  []string `json:"features"`
		Version   string   `json:"version"`
		BuildTime string   `json:"build_time"`
	}{
		Features: []string{
			"multi_cloud",
			"iam",
			"monitoring",
			"finance",
			"message_center",
		},
		Version:   "1.0.0",
		BuildTime: "2026-04-22",
	}

	c.JSON(http.StatusOK, capabilities)
}
