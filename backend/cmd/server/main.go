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
	); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
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
			authSourceGroup.POST("/:id/test", authSourceHandler.Test)
		}

		userHandler := handler.NewUserHandler(db, logger)
		userGroup := v1.Group("/users")
		{
			userGroup.GET("", userHandler.List)
			userGroup.GET("/:id", userHandler.Get)
			userGroup.POST("", userHandler.Create)
			userGroup.DELETE("/:id", userHandler.Delete)
			userGroup.POST("/:id/enable", userHandler.Enable)
			userGroup.POST("/:id/disable", userHandler.Disable)
		}

		roleHandler := handler.NewRoleHandler(db, logger)
		roleGroup := v1.Group("/roles")
		{
			roleGroup.GET("", roleHandler.List)
			roleGroup.GET("/:id", roleHandler.Get)
			roleGroup.POST("", roleHandler.Create)
			roleGroup.DELETE("/:id", roleHandler.Delete)
			roleGroup.GET("/permissions", roleHandler.ListPermissions)
			roleGroup.POST("/:id/permissions", roleHandler.AssignPermission)
			roleGroup.GET("/:id/permissions", roleHandler.GetPermissions)
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
