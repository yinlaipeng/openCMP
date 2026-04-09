package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
	"github.com/opencmp/opencmp/internal/service"
)

// IAMHandler IAM 模块 Handler
type IAMHandler struct {
	userService    *service.UserService
	roleService    *service.RoleService
	domainService  *service.DomainService
	projectService *service.ProjectService
	groupService   *service.GroupService
	policyService  *service.PolicyService
	logger         *zap.Logger
}

// NewIAMHandler 创建 IAM Handler
func NewIAMHandler(db *gorm.DB, logger *zap.Logger) *IAMHandler {
	return &IAMHandler{
		userService:    service.NewUserService(db),
		roleService:    service.NewRoleService(db),
		domainService:  service.NewDomainService(db),
		projectService: service.NewProjectService(db),
		groupService:   service.NewGroupService(db),
		logger:         logger,
	}
}

// ListUsersWithRoles 获取用户列表（带角色信息）
func (h *IAMHandler) ListUsersWithRoles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	users, total, err := h.userService.ListUsers(c.Request.Context(), nil, "", "", nil, pageSize, offset)
	if err != nil {
		h.logger.Error("failed to list users", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list users"})
		return
	}

	type UserWithRoles struct {
		model.User
		Roles []model.Role `json:"roles"`
	}

	var usersWithRoles []UserWithRoles
	for _, user := range users {
		roles, err := h.getEffectiveRolesForUser(c.Request.Context(), user.ID)
		if err != nil {
			h.logger.Error("failed to get user roles", zap.Uint("user_id", user.ID), zap.Error(err))
			continue
		}

		usersWithRoles = append(usersWithRoles, UserWithRoles{
			User:  *user,
			Roles: roles,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": usersWithRoles,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
			"pages":     (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// getEffectiveRolesForUser 获取用户的有效角色（包括通过组获得的角色）
func (h *IAMHandler) getEffectiveRolesForUser(ctx context.Context, userID uint) ([]model.Role, error) {
	var roles []model.Role

	// 获取用户直接拥有的角色
	userRoles, err := h.userService.GetUserRoles(ctx, userID, 0) // 0 表示所有域
	if err != nil {
		return nil, err
	}

	for _, role := range userRoles {
		roles = append(roles, *role)
	}

	// 获取用户通过组获得的角色
	userGroups, err := h.userService.GetUserGroups(ctx, userID)
	if err != nil {
		return nil, err
	}

	for _, group := range userGroups {
		// 获取组在各个域的角色
		groupRoles, err := h.getGroupRolesInDomains(ctx, group.ID)
		if err != nil {
			h.logger.Error("failed to get group roles", zap.Uint("group_id", group.ID), zap.Error(err))
			continue
		}
		roles = append(roles, groupRoles...)
	}

	return roles, nil
}

// getGroupRolesInDomains 获取组的所有角色
func (h *IAMHandler) getGroupRolesInDomains(ctx context.Context, groupID uint) ([]model.Role, error) {
	groupRoles, err := h.groupService.GetGroupRoles(ctx, groupID)
	if err != nil {
		return nil, err
	}
	roles := make([]model.Role, 0, len(groupRoles))
	for _, role := range groupRoles {
		roles = append(roles, *role)
	}
	return roles, nil
}
