package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
)

// PermissionMiddleware 权限检查中间件
func PermissionMiddleware(logger *zap.Logger, db *gorm.DB, requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		// 检查用户权限
		hasPermission, err := checkUserPermission(db, userID.(uint), requiredPermission)
		if err != nil {
			logger.Error("Error checking user permission", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			c.Abort()
			return
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("Permission '%s' required", requiredPermission)})
			c.Abort()
			return
		}

		c.Next()
	}
}

// checkUserPermission 检查用户是否具有指定权限
func checkUserPermission(db *gorm.DB, userID uint, permissionName string) (bool, error) {
	// 检查用户是否为超级管理员（拥有所有权限）
	var adminRole model.Role
	if err := db.Where("name = ? AND type = ?", "admin", "system").First(&adminRole).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果没有找到管理员角色，则继续检查普通权限
		} else {
			return false, err
		}
	} else {
		// 检查用户是否拥有管理员角色
		var userRole model.UserRole
		if err := db.Where("user_id = ? AND role_id = ?", userID, adminRole.ID).First(&userRole).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return false, err
			}
		} else {
			// 用户是管理员，拥有所有权限
			return true, nil
		}
	}

	// 检查用户直接拥有的权限
	if hasDirectPermission(db, userID, permissionName) {
		return true, nil
	}

	// 检查用户所属组的权限
	if hasGroupPermission(db, userID, permissionName) {
		return true, nil
	}

	// 检查用户角色的权限
	if hasRolePermission(db, userID, permissionName) {
		return true, nil
	}

	return false, nil
}

// hasDirectPermission 检查用户是否直接拥有权限（未来可能支持用户直接权限）
func hasDirectPermission(db *gorm.DB, userID uint, permissionName string) bool {
	// 当前版本不支持用户直接权限，仅通过角色分配权限
	// 如果需要支持用户直接权限，可以在此处实现
	return false
}

// hasGroupPermission 检查用户所属组是否拥有权限
func hasGroupPermission(db *gorm.DB, userID uint, permissionName string) bool {
	// 获取用户信息以确定其组
	var user model.User
	if err := db.Select("id, group_id").First(&user, userID).Error; err != nil {
		return false
	}

	if user.GroupID == 0 {
		// 用户不属于任何组
		return false
	}

	// 获取组的角色
	var groupRoles []model.GroupRole
	if err := db.Where("group_id = ?", user.GroupID).Find(&groupRoles).Error; err != nil {
		return false
	}

	// 检查组角色是否拥有权限
	for _, groupRole := range groupRoles {
		if hasRolePermissionByID(db, groupRole.RoleID, permissionName) {
			return true
		}
	}

	return false
}

// hasRolePermission 检查用户角色是否拥有权限
func hasRolePermission(db *gorm.DB, userID uint, permissionName string) bool {
	// 获取用户的角色
	var userRoles []model.UserRole
	if err := db.Where("user_id = ?", userID).Find(&userRoles).Error; err != nil {
		return false
	}

	// 检查用户角色是否拥有权限
	for _, userRole := range userRoles {
		if hasRolePermissionByID(db, userRole.RoleID, permissionName) {
			return true
		}
	}

	return false
}

// hasRolePermissionByID 检查指定角色是否拥有权限
func hasRolePermissionByID(db *gorm.DB, roleID uint, permissionName string) bool {
	// 检查角色-权限关联
	var rolePermission model.RolePermission
	if err := db.Joins("JOIN permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ? AND permissions.name = ?", roleID, permissionName).
		First(&rolePermission).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 检查角色是否通过策略拥有权限
			return hasRolePermissionViaPolicy(db, roleID, permissionName)
		}
		return false
	}

	return true
}

// hasRolePermissionViaPolicy 检查角色是否通过策略拥有权限
func hasRolePermissionViaPolicy(db *gorm.DB, roleID uint, permissionName string) bool {
	// 解析权限名称 (格式: resource:action)
	parts := strings.Split(permissionName, ":")
	if len(parts) < 2 {
		return false
	}
	
	resource := parts[0]
	action := parts[1]

	// 获取角色关联的策略
	var rolePolicies []model.RolePolicy
	if err := db.Where("role_id = ?", roleID).Find(&rolePolicies).Error; err != nil {
		return false
	}

	// 检查策略是否允许该权限
	for _, rolePolicy := range rolePolicies {
		if hasPermissionViaPolicyID(db, rolePolicy.PolicyID, resource, action) {
			return true
		}
	}

	return false
}

// hasPermissionViaPolicyID 检查通过指定策略ID是否拥有权限
func hasPermissionViaPolicyID(db *gorm.DB, policyID uint, resource, action string) bool {
	var policy model.Policy
	if err := db.First(&policy, policyID).Error; err != nil {
		return false
	}

	// 解析策略语句
	var statements []model.PolicyStatement
	if err := db.Where("policy_id = ?", policyID).Find(&statements).Error; err != nil {
		return false
	}

	for _, stmt := range statements {
		// 检查资源匹配
		resourceMatch := stmt.Resource == "*" || stmt.Resource == resource || 
			(strings.HasSuffix(stmt.Resource, ":*") && strings.HasPrefix(resource, strings.TrimSuffix(stmt.Resource, ":*")))
		
		// 检查动作匹配
		actionMatch := containsString(stmt.Actions, action) || containsString(stmt.Actions, "*")

		if resourceMatch && actionMatch && stmt.Effect == "Allow" {
			return true
		}
	}

	return false
}

// containsString 检查字符串切片是否包含指定字符串
func containsString(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// RequirePermissions 需要任意一个权限（OR关系）
func RequireAnyPermission(logger *zap.Logger, db *gorm.DB, permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		// 检查用户是否拥有任一所需权限
		for _, permission := range permissions {
			hasPermission, err := checkUserPermission(db, userID.(uint), permission)
			if err != nil {
				logger.Error("Error checking user permission", zap.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				c.Abort()
				return
			}

			if hasPermission {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Required permission not granted"})
		c.Abort()
	}
}

// RequireAllPermissions 需要所有权限（AND关系）
func RequireAllPermissions(logger *zap.Logger, db *gorm.DB, permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		// 检查用户是否拥有所有所需权限
		for _, permission := range permissions {
			hasPermission, err := checkUserPermission(db, userID.(uint), permission)
			if err != nil {
				logger.Error("Error checking user permission", zap.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				c.Abort()
				return
			}

			if !hasPermission {
				c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("Permission '%s' required", permission)})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}