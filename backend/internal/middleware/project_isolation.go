package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
)

// ProjectIsolationMiddleware 项目隔离中间件
// 确保用户只能访问其所属项目的资源（除非是系统管理员）
func ProjectIsolationMiddleware(db *gorm.DB, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从context获取user_id（由AuthMiddleware注入）
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
			c.Abort()
			return
		}

		// 检查用户是否为系统管理员（跳过项目隔离）
		if isSystemAdmin(db, userID.(uint)) {
			// 系统管理员可以访问所有项目，设置标记
			c.Set("is_admin", true)
			c.Set("all_projects_visible", true)
			c.Next()
			return
		}

		// 获取用户所属的项目ID列表
		projectIDs := getUserProjectIDs(db, userID.(uint))

		// 检查是否为财务人员角色（可以查看全局账单数据）
		if isFinanceRole(db, userID.(uint)) {
			// 财务人员对账单相关API有全局访问权限
			if isFinanceAPI(c.Request.URL.Path) {
				c.Set("is_finance", true)
				c.Set("all_projects_visible", true)
				c.Next()
				return
			}
		}

		// 检查是否为运维工程师角色（只能查看所属项目）
		if len(projectIDs) == 0 {
			// 用户没有任何项目，只能访问公共资源
			logger.Warn("用户无项目归属",
				zap.Uint("user_id", userID.(uint)),
				zap.String("path", c.Request.URL.Path))

			// 对于资源类API，返回空数据而不是拒绝
			// 这可以通过service层检查project_ids是否为空来处理
		}

		// 注入项目ID列表到context
		c.Set("project_ids", projectIDs)
		c.Set("project_isolation", true)

		logger.Debug("项目隔离检查",
			zap.Uint("user_id", userID.(uint)),
			zap.Int64s("project_ids", projectIDs))

		c.Next()
	}
}

// getUserProjectIDs 获取用户所属的项目ID列表
func getUserProjectIDs(db *gorm.DB, userID uint) []int64 {
	var projectIDs []int64

	// 1. 通过 ProjectUserRole 直接获取用户的项目
	db.Table("project_user_roles").
		Where("user_id = ?", userID).
		Pluck("project_id", &projectIDs)

	// 2. 通过用户所属的组获取项目
	var groupIDs []uint
	db.Table("user_groups").
		Where("user_id = ?", userID).
		Pluck("group_id", &groupIDs)

	if len(groupIDs) > 0 {
		var groupProjectIDs []int64
		db.Table("group_projects").
			Where("group_id IN (?)", groupIDs).
			Pluck("project_id", &groupProjectIDs)

		// 合并项目ID列表（去重）
		projectIDs = mergeUnique(projectIDs, groupProjectIDs)
	}

	return projectIDs
}

// mergeUnique 合并两个int64列表并去重
func mergeUnique(a, b []int64) []int64 {
	seen := make(map[int64]bool)
	result := []int64{}

	for _, id := range a {
		if !seen[id] {
			seen[id] = true
			result = append(result, id)
		}
	}

	for _, id := range b {
		if !seen[id] {
			seen[id] = true
			result = append(result, id)
		}
	}

	return result
}

// isFinanceRole 检查用户是否为财务人员角色
func isFinanceRole(db *gorm.DB, userID uint) bool {
	var count int64
	err := db.Table("user_roles").
		Joins("JOIN roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ?", userID).
		Where("roles.name IN (?)", []string{"finance", "财务人员", "财务"}).
		Count(&count).Error

	if err != nil {
		return false
	}

	return count > 0
}

// isFinanceAPI 检查是否为财务相关API
func isFinanceAPI(path string) bool {
	financePaths := []string{
		"/api/v1/finance/bills",
		"/api/v1/finance/orders",
		"/api/v1/finance/renewals",
		"/api/v1/finance/account-balance",
		"/api/v1/finance/cost",
		"/api/v1/finance/anomalies",
		"/api/v1/finance/budgets",
	}

	for _, fp := range financePaths {
		if strings.HasPrefix(path, fp) {
			return true
		}
	}

	return false
}

// ProjectFilter 项目过滤辅助结构
type ProjectFilter struct {
	AllProjectsVisible bool
	ProjectIDs         []int64
}

// GetProjectFilter 从context获取项目过滤信息
func GetProjectFilter(c *gin.Context) ProjectFilter {
	filter := ProjectFilter{}

	if allVisible, exists := c.Get("all_projects_visible"); exists {
		filter.AllProjectsVisible = allVisible.(bool)
	}

	if projectIDs, exists := c.Get("project_ids"); exists {
		filter.ProjectIDs = projectIDs.([]int64)
	}

	return filter
}

// ApplyProjectFilter 应用项目过滤到查询
// 返回一个gorm ScopeFunc用于添加WHERE条件
func ApplyProjectFilter(c *gin.Context, db *gorm.DB, tableName string, projectField string) *gorm.DB {
	filter := GetProjectFilter(c)

	if filter.AllProjectsVisible {
		return db // 无需过滤，返回所有数据
	}

	if len(filter.ProjectIDs) == 0 {
		// 用户无项目，返回空结果
		return db.Where("1 = 0") // 强制返回空结果
	}

	// 应用项目过滤条件
	return db.Where(tableName+"."+projectField+" IN (?)", filter.ProjectIDs)
}

// ProjectIsolationScope 项目隔离Scope函数
// 用于在service层快速应用项目过滤
func ProjectIsolationScope(c *gin.Context, projectField string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		filter := GetProjectFilter(c)

		if filter.AllProjectsVisible {
			return db
		}

		if len(filter.ProjectIDs) == 0 {
			return db.Where("1 = 0")
		}

		return db.Where(projectField+" IN (?)", filter.ProjectIDs)
	}
}

// CheckProjectAccess 检查用户是否有访问特定项目的权限
func CheckProjectAccess(db *gorm.DB, userID uint, projectID uint) bool {
	// 系统管理员可以访问所有项目
	if isSystemAdmin(db, userID) {
		return true
	}

	// 检查用户是否属于该项目
	projectIDs := getUserProjectIDs(db, userID)
	for _, pid := range projectIDs {
		if pid == int64(projectID) {
			return true
		}
	}

	return false
}

// RequireProjectAccess 项目访问检查中间件
// 用于需要指定项目ID的API端点
func RequireProjectAccess(db *gorm.DB, logger *zap.Logger, projectIDParam string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
			c.Abort()
			return
		}

		// 系统管理员跳过检查
		if isSystemAdmin(db, userID.(uint)) {
			c.Next()
			return
		}

		// 获取请求中的项目ID
		projectIDStr := c.Param(projectIDParam)
		if projectIDStr == "" {
			projectIDStr = c.Query(projectIDParam)
		}
		if projectIDStr == "" {
			// 尝试从请求体获取（需要解析JSON）
			// 这里简化处理，如果路径和查询参数都没有，跳过检查
			c.Next()
			return
		}

		// 转换项目ID
		var projectID uint
		for _, c := range projectIDStr {
			if c >= '0' && c <= '9' {
				projectID = projectID*10 + uint(c-'0')
			} else {
				break
			}
		}

		// 检查访问权限
		if !CheckProjectAccess(db, userID.(uint), projectID) {
			logger.Warn("项目访问权限检查失败",
				zap.Uint("user_id", userID.(uint)),
				zap.Uint("project_id", projectID))

			c.JSON(http.StatusForbidden, gin.H{
				"error":     "无权访问该项目",
				"code":      "PROJECT_ACCESS_DENIED",
				"project_id": projectID,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetUserProjects 获取用户可访问的项目列表（完整信息）
func GetUserProjects(db *gorm.DB, userID uint) []model.Project {
	// 系统管理员可以看到所有项目
	if isSystemAdmin(db, userID) {
		var projects []model.Project
		db.Where("enabled = ?", true).Find(&projects)
		return projects
	}

	projectIDs := getUserProjectIDs(db, userID)
	if len(projectIDs) == 0 {
		return nil
	}

	var projects []model.Project
	db.Where("id IN (?)", projectIDs).
		Where("enabled = ?", true).
		Find(&projects)

	return projects
}