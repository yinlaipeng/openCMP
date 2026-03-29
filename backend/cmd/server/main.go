package main

import (
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

	// 初始化 Gin
	gin.SetMode(cfg.Server.Mode)
	r := gin.New()

	// 使用中间件
	r.Use(middleware.RecoveryMiddleware(logger))
	r.Use(middleware.LoggerMiddleware(logger))

	// API 路由组
	v1 := r.Group("/api/v1")
	v1.Use(middleware.AuthMiddleware())
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
	}

	// 健康检查（不需要认证）
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 认证路由（不需要认证）
	authHandler := handler.NewAuthHandler(db, logger)
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
