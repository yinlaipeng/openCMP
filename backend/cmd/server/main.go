package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/handler"
	"github.com/opencmp/opencmp/internal/middleware"
	"github.com/opencmp/opencmp/internal/model"
	"github.com/opencmp/opencmp/internal/service"
	"github.com/opencmp/opencmp/pkg/utils"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Auth     AuthConfig     `yaml:"auth"`
	Log      LogConfig      `yaml:"log"`
}

type ServerConfig struct {
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

type DatabaseConfig struct {
	Driver       string `yaml:"driver"`
	DSN          string `yaml:"dsn"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxOpenConns int    `yaml:"max_open_conns"`
}

type AuthConfig struct {
	JWTSecret        string `yaml:"jwt_secret"`
	TokenExpireHours int    `yaml:"token_expire_hours"`
}

type LogConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

func main() {
	configPath := flag.String("config", "configs/config.yaml", "config file path")
	flag.Parse()

	// 加载配置
	cfg, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// 初始化日志
	logger, err := initLogger(cfg.Log)
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	defer logger.Sync()

	// 初始化数据库
	db, err := initDatabase(cfg.Database)
	if err != nil {
		log.Fatalf("failed to init database: %v", err)
	}

	// 自动迁移表结构
	if err := db.AutoMigrate(
		&model.CloudAccount{},
		&model.Domain{},
		&model.Project{},
		&model.User{},
		&model.Group{},
		&model.UserGroup{},
		&model.Role{},
		&model.Permission{},
		&model.RolePermission{},
		&model.UserRole{},
		&model.ProjectUserRole{},
		&model.GroupRole{},
		&model.AuthSource{},
		&model.SecurityAlert{},
		&model.MessageType{},
		&model.Message{},
		&model.NotificationChannel{},
		&model.MessageSubscription{},
		&model.Robot{},
		&model.Receiver{},
		&model.Policy{},
		&model.RolePolicy{},
	); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	// 初始化内置认证源
	if err := initBuiltInAuthSource(db); err != nil {
		logger.Warn("failed to init built-in auth source", zap.Error(err))
	}

	// 初始化默认数据
	if err := initDefaultData(db); err != nil {
		logger.Warn("failed to init default data", zap.Error(err))
	}

	// 初始化 Gin
	gin.SetMode(cfg.Server.Mode)
	r := gin.New()

	// 使用中间件
	r.Use(middleware.RecoveryMiddleware(logger))
	r.Use(middleware.LoggerMiddleware(logger))

	// API 路由组
	v1 := r.Group("/api/v1")
	v1.Use(func(c *gin.Context) {
		c.Set("jwt_secret", cfg.Auth.JWTSecret)
		c.Next()
	})
	v1.Use(middleware.AuthMiddleware(logger))
	{
		// 云账户路由
		cloudAccountHandler := handler.NewCloudAccountHandler(db, logger)
		cloudAccountGroup := v1.Group("/cloud-accounts")
		{
			cloudAccountGroup.POST("", cloudAccountHandler.Create)
			cloudAccountGroup.GET("", cloudAccountHandler.List)
			cloudAccountGroup.GET("/:id", cloudAccountHandler.Get)
			cloudAccountGroup.PUT("/:id", cloudAccountHandler.Update)
			cloudAccountGroup.DELETE("/:id", cloudAccountHandler.Delete)
			cloudAccountGroup.POST("/:id/verify", cloudAccountHandler.Verify)
		}

		// 计算资源路由
		computeHandler := handler.NewComputeHandler(db, logger)
		computeGroup := v1.Group("/compute")
		{
			computeGroup.POST("/vms", computeHandler.CreateVM)
			computeGroup.GET("/vms", computeHandler.ListVMs)
			computeGroup.GET("/vms/:id", computeHandler.GetVM)
			computeGroup.DELETE("/vms/:id", computeHandler.DeleteVM)
			computeGroup.POST("/vms/:id/action", computeHandler.VMAction)
			computeGroup.GET("/images", computeHandler.ListImages)
		}

		// 网络资源路由
		networkHandler := handler.NewNetworkHandler(db, logger)
		networkGroup := v1.Group("/network")
		{
			networkGroup.POST("/vpcs", networkHandler.CreateVPC)
			networkGroup.GET("/vpcs", networkHandler.ListVPCs)
			networkGroup.DELETE("/vpcs/:id", networkHandler.DeleteVPC)

			networkGroup.POST("/subnets", networkHandler.CreateSubnet)
			networkGroup.GET("/subnets", networkHandler.ListSubnets)

			networkGroup.POST("/security-groups", networkHandler.CreateSecurityGroup)
			networkGroup.GET("/security-groups", networkHandler.ListSecurityGroups)

			networkGroup.POST("/eips", networkHandler.CreateEIP)
			networkGroup.GET("/eips", networkHandler.ListEIPs)
		}

		// IAM 认证与安全路由
		authSourceHandler := handler.NewAuthSourceHandler(db, logger)
		authSourceGroup := v1.Group("/auth-sources")
		{
			authSourceGroup.GET("", authSourceHandler.List)
			authSourceGroup.GET("/:id", authSourceHandler.Get)
			authSourceGroup.POST("", authSourceHandler.Create)
			authSourceGroup.PUT("/:id", authSourceHandler.Update)
			authSourceGroup.DELETE("/:id", authSourceHandler.Delete)
			authSourceGroup.POST("/:id/test", authSourceHandler.Test)
			authSourceGroup.POST("/:id/enable", authSourceHandler.Enable)
			authSourceGroup.POST("/:id/disable", authSourceHandler.Disable)
		}

		// 域管理路由
		domainHandler := handler.NewDomainHandler(db, logger)
		domainGroup := v1.Group("/domains")
		{
			domainGroup.GET("", domainHandler.List)
			domainGroup.GET("/:id", domainHandler.Get)
			domainGroup.POST("", domainHandler.Create)
			domainGroup.PUT("/:id", domainHandler.Update)
			domainGroup.DELETE("/:id", domainHandler.Delete)
			domainGroup.POST("/:id/enable", domainHandler.Enable)
			domainGroup.POST("/:id/disable", domainHandler.Disable)
			domainGroup.GET("/:id/users", domainHandler.GetUsers)
			domainGroup.GET("/:id/groups", domainHandler.GetGroups)
			domainGroup.GET("/:id/projects", domainHandler.GetProjects)
			domainGroup.GET("/:id/roles", domainHandler.GetRoles)
		}

		// 项目管理路由
		projectHandler := handler.NewProjectHandler(db, logger)
		projectGroup := v1.Group("/projects")
		{
			projectGroup.GET("", projectHandler.List)
			projectGroup.GET("/:id", projectHandler.Get)
			projectGroup.POST("", projectHandler.Create)
			projectGroup.PUT("/:id", projectHandler.Update)
			projectGroup.DELETE("/:id", projectHandler.Delete)
			projectGroup.POST("/:id/enable", projectHandler.Enable)
			projectGroup.POST("/:id/disable", projectHandler.Disable)
			projectGroup.POST("/:id/join", projectHandler.Join)
			projectGroup.GET("/:id/users", projectHandler.GetUsers)
			projectGroup.GET("/:id/roles", projectHandler.GetRoles)
			projectGroup.DELETE("/:id/users", projectHandler.RemoveUser)
		}

		userHandler := handler.NewUserHandler(db, logger)
		userGroup := v1.Group("/users")
		{
			userGroup.GET("", userHandler.List)
			userGroup.GET("/:id", userHandler.Get)
			userGroup.POST("", userHandler.Create)
			userGroup.PUT("/:id", userHandler.Update)
			userGroup.DELETE("/:id", userHandler.Delete)
			userGroup.POST("/:id/enable", userHandler.Enable)
			userGroup.POST("/:id/disable", userHandler.Disable)
			userGroup.POST("/:id/reset-password", userHandler.ResetPassword)
			userGroup.GET("/:id/roles", userHandler.GetRoles)
			userGroup.POST("/:id/roles", userHandler.AssignRole)
			userGroup.DELETE("/:id/roles", userHandler.RevokeRole)
			userGroup.GET("/:id/groups", userHandler.GetGroups)
			userGroup.POST("/:id/groups", userHandler.JoinGroup)
			userGroup.DELETE("/:id/groups", userHandler.LeaveGroup)
		}

		// 用户组路由
		groupHandler := handler.NewGroupHandler(db, logger)
		groupGroup := v1.Group("/groups")
		{
			groupGroup.GET("", groupHandler.List)
			groupGroup.GET("/:id", groupHandler.Get)
			groupGroup.POST("", groupHandler.Create)
			groupGroup.PUT("/:id", groupHandler.Update)
			groupGroup.DELETE("/:id", groupHandler.Delete)
			groupGroup.GET("/:id/users", groupHandler.GetUsers)
			groupGroup.POST("/:id/users", groupHandler.AddUser)
			groupGroup.DELETE("/:id/users", groupHandler.RemoveUser)
		}

		roleHandler := handler.NewRoleHandler(db, logger)
		policyHandler := handler.NewPolicyHandler(db, logger)
		roleGroup := v1.Group("/roles")
		{
			roleGroup.GET("", roleHandler.List)
			roleGroup.GET("/:id", roleHandler.Get)
			roleGroup.POST("", roleHandler.Create)
			roleGroup.PUT("/:id", roleHandler.Update)
			roleGroup.DELETE("/:id", roleHandler.Delete)
			roleGroup.POST("/:id/enable", roleHandler.Enable)
			roleGroup.POST("/:id/disable", roleHandler.Disable)
			roleGroup.POST("/:id/public", roleHandler.MakePublic)
			roleGroup.GET("/:id/users", roleHandler.GetRoleUsers)
			roleGroup.GET("/:id/groups", roleHandler.GetRoleGroups)
		}

		// 权限管理路由
		permissionGroup := v1.Group("/permissions")
		{
			permissionGroup.GET("", roleHandler.ListPermissions)
			permissionGroup.GET("/:id", roleHandler.GetPermission)
			permissionGroup.POST("", roleHandler.CreatePermission)
			permissionGroup.PUT("/:id", roleHandler.UpdatePermission)
			permissionGroup.DELETE("/:id", roleHandler.DeletePermission)
			permissionGroup.GET("/resources", roleHandler.ListResources)
			permissionGroup.GET("/actions", roleHandler.ListActions)
		}

		// 角色权限关联路由
		rolePermissionGroup := v1.Group("/roles")
		{
			rolePermissionGroup.GET("/:id/permissions", roleHandler.GetPermissions)
			rolePermissionGroup.POST("/:id/permissions", roleHandler.AssignPermission)
			rolePermissionGroup.DELETE("/:id/permissions", roleHandler.RevokePermission)
			rolePermissionGroup.GET("/:id/policies", roleHandler.GetRolePolicies)
			rolePermissionGroup.POST("/:id/policies", roleHandler.AssignPolicyToRole)
			rolePermissionGroup.DELETE("/:id/policies", roleHandler.RevokePolicyFromRole)
		}

		// 策略管理路由
		policyGroup := v1.Group("/policies")
		{
			policyGroup.GET("", policyHandler.List)
			policyGroup.GET("/:id", policyHandler.GetPolicy)
			policyGroup.POST("", policyHandler.CreatePolicy)
			policyGroup.PUT("/:id", policyHandler.UpdatePolicy)
			policyGroup.DELETE("/:id", policyHandler.DeletePolicy)
		}

		// IAM 综合管理路由
		iamHandler := handler.NewIAMHandler(db, logger)
		iamGroup := v1.Group("/iam")
		{
			iamGroup.GET("/users-with-roles", iamHandler.ListUsersWithRoles)
			iamGroup.GET("/users/:id/permissions", iamHandler.GetUserPermissions)
			iamGroup.POST("/check-permission", iamHandler.CheckPermission)
		}
	}

	// 健康检查（不需要认证）
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 认证路由（不需要认证）
	authHandler := handler.NewAuthHandler(db, logger, cfg.Auth.JWTSecret, cfg.Auth.TokenExpireHours)
	r.POST("/api/v1/auth/login", authHandler.Login)

	// 启动服务
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logger.Info("starting server", zap.String("addr", addr))
	if err := r.Run(addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func initLogger(cfg LogConfig) (*zap.Logger, error) {
	var config zap.Config
	if cfg.Format == "json" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	switch cfg.Level {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	}

	return config.Build()
}

// initBuiltInAuthSource 初始化内置认证源
func initBuiltInAuthSource(db *gorm.DB) error {
	// 检查是否已存在系统认证源
	var count int64
	if err := db.Model(&model.AuthSource{}).Where("type = ?", "local").Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil // 已存在，不需要创建
	}

	// 创建系统认证源（本地认证/SQL 认证）
	authSource := &model.AuthSource{
		Name:        "系统认证",
		Description: "系统内置的本地用户认证（SQL 认证）",
		Type:        "local",
		Scope:       "system",
		DomainID:    nil,
		Enabled:     true,
		AutoCreate:  false,
		Config:      nil, // 本地认证不需要额外配置
	}

	return db.Create(authSource).Error
}

// initDefaultData 初始化默认数据
func initDefaultData(db *gorm.DB) error {
	ctx := context.Background()

	// 创建默认域
	defaultDomain := &model.Domain{
		Name:        "Default",
		Description: "默认域",
		Enabled:     true,
	}
	
	var existingDomain model.Domain
	if err := db.Where("name = ?", "Default").First(&existingDomain).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			if err := db.Create(defaultDomain).Error; err != nil {
				return fmt.Errorf("failed to create default domain: %v", err)
			}
		} else {
			return fmt.Errorf("failed to check default domain: %v", err)
		}
	} else {
		defaultDomain = &existingDomain
	}

	// 创建管理员角色
	adminRole := &model.Role{
		Name:        "admin",
		DisplayName: "系统管理员",
		Description: "系统管理员角色，拥有所有权限",
		Type:        "system",
		Enabled:     true,
		IsPublic:    true,
	}
	
	var existingRole model.Role
	if err := db.Where("name = ?", "admin").First(&existingRole).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			if err := db.Create(adminRole).Error; err != nil {
				return fmt.Errorf("failed to create admin role: %v", err)
			}
		} else {
			return fmt.Errorf("failed to check admin role: %v", err)
		}
	} else {
		adminRole = &existingRole
	}

	// 创建默认管理员用户
	adminUser := &model.User{
		Name:        "admin",
		DisplayName: "管理员",
		Email:       "admin@example.com",
		Password:    "admin123", // 密码将在创建时被哈希
		DomainID:    defaultDomain.ID,
		Enabled:     true,
	}
	
	var existingUser model.User
	if err := db.Where("name = ?", "admin").First(&existingUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 对密码进行哈希处理
			hashedPassword, err := utils.HashPassword(adminUser.Password)
			if err != nil {
				return fmt.Errorf("failed to hash password: %v", err)
			}
			adminUser.Password = hashedPassword
			
			if err := db.Create(adminUser).Error; err != nil {
				return fmt.Errorf("failed to create admin user: %v", err)
			}
			
			// 为管理员用户分配管理员角色
			userService := service.NewUserService(db)
			if err := userService.AssignUserRole(ctx, adminUser.ID, adminRole.ID, defaultDomain.ID); err != nil {
				return fmt.Errorf("failed to assign admin role to admin user: %v", err)
			}
		} else {
			return fmt.Errorf("failed to check admin user: %v", err)
		}
	}

	// 创建一些基本权限
	permissions := []model.Permission{
		{Name: "user:list", DisplayName: "用户列表", Description: "查看用户列表", Resource: "user", Action: "list", Type: "system"},
		{Name: "user:create", DisplayName: "创建用户", Description: "创建新用户", Resource: "user", Action: "create", Type: "system"},
		{Name: "user:update", DisplayName: "更新用户", Description: "更新用户信息", Resource: "user", Action: "update", Type: "system"},
		{Name: "user:delete", DisplayName: "删除用户", Description: "删除用户", Resource: "user", Action: "delete", Type: "system"},
		{Name: "role:list", DisplayName: "角色列表", Description: "查看角色列表", Resource: "role", Action: "list", Type: "system"},
		{Name: "role:create", DisplayName: "创建角色", Description: "创建新角色", Resource: "role", Action: "create", Type: "system"},
		{Name: "role:update", DisplayName: "更新角色", Description: "更新角色信息", Resource: "role", Action: "update", Type: "system"},
		{Name: "role:delete", DisplayName: "删除角色", Description: "删除角色", Resource: "role", Action: "delete", Type: "system"},
		{Name: "domain:list", DisplayName: "域列表", Description: "查看域列表", Resource: "domain", Action: "list", Type: "system"},
		{Name: "domain:create", DisplayName: "创建域", Description: "创建新域", Resource: "domain", Action: "create", Type: "system"},
		{Name: "domain:update", DisplayName: "更新域", Description: "更新域信息", Resource: "domain", Action: "update", Type: "system"},
		{Name: "domain:delete", DisplayName: "删除域", Description: "删除域", Resource: "domain", Action: "delete", Type: "system"},
		{Name: "project:list", DisplayName: "项目列表", Description: "查看项目列表", Resource: "project", Action: "list", Type: "system"},
		{Name: "project:create", DisplayName: "创建项目", Description: "创建新项目", Resource: "project", Action: "create", Type: "system"},
		{Name: "project:update", DisplayName: "更新项目", Description: "更新项目信息", Resource: "project", Action: "update", Type: "system"},
		{Name: "project:delete", DisplayName: "删除项目", Description: "删除项目", Resource: "project", Action: "delete", Type: "system"},
		{Name: "cloud-account:list", DisplayName: "云账户列表", Description: "查看云账户列表", Resource: "cloud-account", Action: "list", Type: "system"},
		{Name: "cloud-account:create", DisplayName: "创建云账户", Description: "创建新云账户", Resource: "cloud-account", Action: "create", Type: "system"},
		{Name: "cloud-account:update", DisplayName: "更新云账户", Description: "更新云账户信息", Resource: "cloud-account", Action: "update", Type: "system"},
		{Name: "cloud-account:delete", DisplayName: "删除云账户", Description: "删除云账户", Resource: "cloud-account", Action: "delete", Type: "system"},
	}

	roleService := service.NewRoleService(db)
	for _, perm := range permissions {
		var existingPerm model.Permission
		if err := db.Where("name = ?", perm.Name).First(&existingPerm).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := roleService.CreatePermission(ctx, &perm); err != nil {
					// 可能已经存在，忽略错误
					continue
				}
			} else {
				return fmt.Errorf("failed to check permission %s: %v", perm.Name, err)
			}
		}
	}

	// 为管理员角色分配所有权限
	for _, perm := range permissions {
		var existingPerm model.Permission
		if err := db.Where("name = ?", perm.Name).First(&existingPerm).Error; err != nil {
			return fmt.Errorf("failed to get permission %s: %v", perm.Name, err)
		}
		
		if err := roleService.AssignPermissionToRole(ctx, adminRole.ID, existingPerm.ID); err != nil {
			// 可能已经分配了，忽略错误
			continue
		}
	}

	return nil
}

func initDatabase(cfg DatabaseConfig) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)

	return db, nil
}
