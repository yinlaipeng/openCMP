package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
)

// PermissionMiddleware 通用权限检查中间件
// 基于RBAC模型，检查用户是否有访问特定资源的权限
func PermissionMiddleware(db *gorm.DB, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从context获取user_id（由AuthMiddleware注入）
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
			c.Abort()
			return
		}

		// 安全类型转换
		uid, ok := userID.(uint)
		if !ok {
			logger.Error("invalid user_id type in context",
				zap.Any("user_id", userID),
				zap.String("type", "expected uint"))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "内部错误：无效的用户ID类型"})
			c.Abort()
			return
		}

		// 从context获取domain_id（可选，由AuthMiddleware注入）
		domainID, _ := c.Get("domain_id")

		// 解析请求路径，提取资源和操作
		resource, action := parseRequestPath(c.Request.URL.Path, c.Request.Method)

		// 如果无法解析资源，跳过权限检查（可能是公共API）
		if resource == "" {
			c.Next()
			return
		}

		// 检查用户是否为系统管理员（管理员拥有所有权限）
		if isSystemAdmin(db, uid) {
			c.Next()
			return
		}

		// 查询用户权限
		hasPermission := checkUserPermission(db, uid, domainID, resource, action)

		if !hasPermission {
			logger.Warn("权限检查失败",
				zap.Uint("user_id", uid),
				zap.String("resource", resource),
				zap.String("action", action),
				zap.String("path", c.Request.URL.Path))

			c.JSON(http.StatusForbidden, gin.H{
				"error":    "权限不足",
				"code":     "PERMISSION_DENIED",
				"resource": resource,
				"action":   action,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// parseRequestPath 解析请求路径，提取资源和操作
// 例如: /api/v1/users -> resource="user", GET -> action="list"
// 例如: /api/v1/users/123 -> resource="user", GET -> action="get"
func parseRequestPath(path, method string) (string, string) {
	// 移除前缀 /api/v1
	path = strings.TrimPrefix(path, "/api/v1")
	path = strings.TrimPrefix(path, "/")

	// 分割路径
	parts := strings.Split(path, "/")
	if len(parts) == 0 {
		return "", ""
	}

	// 第一个部分通常是资源类型
	resource := parts[0]

	// 转换为单数形式（简化处理）
	resource = toSingular(resource)

	// 根据HTTP方法和路径结构确定操作
	action := parseAction(method, parts)

	return resource, action
}

// parseAction 根据HTTP方法和路径结构确定操作
func parseAction(method string, parts []string) string {
	switch method {
	case "GET":
		// 如果路径有ID参数，是获取单个资源
		if len(parts) > 1 && parts[1] != "" {
			return "get"
		}
		return "list"
	case "POST":
		// 如果是执行类操作（如/verify, /sync, /execute）
		if len(parts) > 1 {
			suffix := parts[len(parts)-1]
			switch suffix {
			case "verify", "test", "sync", "execute", "enable", "disable":
				return suffix
			}
		}
		return "create"
	case "PUT", "PATCH":
		return "update"
	case "DELETE":
		return "delete"
	default:
		return method
	}
}

// toSingular 转换为单数形式（简化处理）
func toSingular(plural string) string {
	// 常见的复数形式转换
	singulars := map[string]string{
		"users":           "user",
		"groups":          "group",
		"roles":           "role",
		"permissions":     "permission",
		"domains":         "domain",
		"projects":        "project",
		"vms":             "vm",
		"vpcs":            "vpc",
		"subnets":         "subnet",
		"eips":            "eip",
		"images":          "image",
		"messages":        "message",
		"robots":          "robot",
		"receivers":       "receiver",
		"subscriptions":   "subscription",
		"channels":        "channel",
		"bills":           "bill",
		"orders":          "order",
		"budgets":         "budget",
		"anomalies":       "anomaly",
		"logs":            "log",
		"accounts":        "account",
		"policies":        "policy",
		"tasks":           "task",
		"alerts":          "alert",
		"disks":           "disk",
		"snapshots":       "snapshot",
		"auth-sources":    "auth_source",
		"cloud-accounts":  "cloud_account",
		"sync-policies":   "sync_policy",
		"scheduled-tasks": "scheduled_task",
		"operation-logs":  "operation_log",
		"notification-channels": "notification_channel",
		"security-groups":      "security_group",
		"cloud-users":          "cloud_user",
		"cloud-user-groups":    "cloud_user_group",
		"cloud-projects":       "cloud_project",
		"cloud-disks":          "cloud_disk",
		"cloud-snapshots":      "cloud_snapshot",
		"autoscaling-groups":   "autoscaling_group",
		"host-templates":       "host_template",
		"vpc-interconnects":    "vpc_interconnect",
		"vpc-peerings":         "vpc_peering",
		"route-tables":         "route_table",
		"l2-networks":          "l2_network",
		"renewals":             "renewal",
	}

	if singular, ok := singulars[plural]; ok {
		return singular
	}

	// 默认处理：去掉末尾的s
	if strings.HasSuffix(plural, "s") && len(plural) > 1 {
		return plural[:len(plural)-1]
	}

	return plural
}

// isSystemAdmin 检查用户是否为系统管理员
func isSystemAdmin(db *gorm.DB, userID uint) bool {
	// 查询用户是否拥有admin角色
	var count int64
	err := db.Table("user_roles").
		Joins("JOIN roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ?", userID).
		Where("roles.name = ?", "admin").
		Where("roles.type = ?", "system").
		Count(&count).Error

	if err != nil {
		return false
	}

	return count > 0
}

// checkUserPermission 检查用户是否有特定资源和操作的权限
func checkUserPermission(db *gorm.DB, userID uint, domainID interface{}, resource, action string) bool {
	// 查询用户的角色ID列表
	var roleIDs []uint
	err := db.Table("user_roles").
		Where("user_id = ?", userID).
		Pluck("role_id", &roleIDs).Error

	if err != nil || len(roleIDs) == 0 {
		return false
	}

	// 查询这些角色拥有的权限
	var count int64
	err = db.Table("role_permissions").
		Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id IN (?)", roleIDs).
		Where("permissions.resource = ?", resource).
		Where("permissions.action = ?", action).
		Where("permissions.enabled = ?", true).
		Count(&count).Error

	if err != nil {
		return false
	}

	return count > 0
}

// PermissionCheckFunc 创建自定义权限检查中间件
// 用于特定API端点的权限检查
func PermissionCheckFunc(db *gorm.DB, logger *zap.Logger, resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
			c.Abort()
			return
		}

		// 系统管理员拥有所有权限
		if isSystemAdmin(db, userID.(uint)) {
			c.Next()
			return
		}

		domainID, _ := c.Get("domain_id")
		if !checkUserPermission(db, userID.(uint), domainID, resource, action) {
			c.JSON(http.StatusForbidden, gin.H{
				"error":    "权限不足",
				"resource": resource,
				"action":   action,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetUserRoleIDs 获取用户的角色ID列表（辅助函数）
func GetUserRoleIDs(db *gorm.DB, userID uint) []uint {
	var roleIDs []uint
	db.Table("user_roles").
		Where("user_id = ?", userID).
		Pluck("role_id", &roleIDs)
	return roleIDs
}

// GetUserPermissions 获取用户的权限列表（辅助函数）
func GetUserPermissions(db *gorm.DB, userID uint) []model.Permission {
	roleIDs := GetUserRoleIDs(db, userID)
	if len(roleIDs) == 0 {
		return nil
	}

	var permissions []model.Permission
	db.Table("permissions").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id IN (?)", roleIDs).
		Where("permissions.enabled = ?", true).
		Find(&permissions)

	return permissions
}