package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/handler"
	"github.com/opencmp/opencmp/internal/middleware"
	"github.com/opencmp/opencmp/internal/model"
	"github.com/opencmp/opencmp/internal/service"
	initutils "github.com/opencmp/opencmp/internal/utils"
	_ "github.com/opencmp/opencmp/pkg/cloudprovider/adapters/alibaba" // 引入阿里云适配器以注册
	pkgutils "github.com/opencmp/opencmp/pkg/utils"
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
	db, sqlDB, err := initDatabase(cfg.Database)
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
		&model.Permission{}, // 新增权限模型
		&model.RolePermission{},
		&model.UserRole{},
		&model.ProjectUserRole{},
		&model.GroupRole{},
		&model.GroupProject{},
		&model.AuthSource{},
		&model.SecurityAlert{},
		&model.MessageType{},
		&model.Message{},
		&model.NotificationChannel{},
		&model.MessageSubscription{},
		&model.Robot{},
		&model.Receiver{},
		&model.ReceiverChannel{},
		&model.SyncPolicy{},    // 添加同步策略模型
		&model.ScheduledTask{}, // 添加定时任务模型
		&model.OperationLog{},  // 添加操作日志模型
	); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	// 执行默认数据初始化SQL脚本
	if err := initutils.ExecuteInitDataScript(sqlDB); err != nil {
		logger.Warn("failed to execute initialization script", zap.Error(err))
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
			cloudAccountGroup.POST("/:id/test-connection", cloudAccountHandler.TestConnection)
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

			// 详细的 VM 信息和操作端点
			computeGroup.GET("/vms/:id/details", computeHandler.GetVMDetails)
			computeGroup.GET("/vms/:id/security-groups", computeHandler.GetVMSecurityGroups)
			computeGroup.GET("/vms/:id/networks", computeHandler.GetVMNetworkInfo)
			computeGroup.GET("/vms/:id/disks", computeHandler.GetVMDisks)
			computeGroup.GET("/vms/:id/snapshots", computeHandler.GetVMSnapshots)
			computeGroup.GET("/vms/:id/logs", computeHandler.GetVMOperationLogs)
			computeGroup.GET("/vms/:id/vnc", computeHandler.GetVNCInfo)

			computeGroup.GET("/images", computeHandler.ListImages)
		}

		// 主机模版路由
		hostTemplateService := service.NewHostTemplateService(db)
		hostTemplateHandler := handler.NewHostTemplateHandler(hostTemplateService)
		hostTemplateGroup := v1.Group("/host-templates")
		{
			hostTemplateGroup.POST("", hostTemplateHandler.CreateHostTemplate)
			hostTemplateGroup.GET("", hostTemplateHandler.ListHostTemplates)
			hostTemplateGroup.GET("/:id", hostTemplateHandler.GetHostTemplate)
			hostTemplateGroup.PUT("/:id", hostTemplateHandler.UpdateHostTemplate)
			hostTemplateGroup.DELETE("/:id", hostTemplateHandler.DeleteHostTemplate)
		}

			// 弹性伸缩组路由
			autoscalingGroupService := service.NewAutoscalingGroupService(db)
			autoscalingGroupHandler := handler.NewAutoscalingGroupHandler(autoscalingGroupService)
			autoscalingGroupGroup := v1.Group("/autoscaling-groups")
			{
				autoscalingGroupGroup.POST("", autoscalingGroupHandler.CreateAutoscalingGroup)
				autoscalingGroupGroup.GET("", autoscalingGroupHandler.ListAutoscalingGroups)
				autoscalingGroupGroup.GET("/:id", autoscalingGroupHandler.GetAutoscalingGroup)
				autoscalingGroupGroup.PUT("/:id", autoscalingGroupHandler.UpdateAutoscalingGroup)
				autoscalingGroupGroup.DELETE("/:id", autoscalingGroupHandler.DeleteAutoscalingGroup)
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

				// 地理资源路由
				networkGroup.GET("/regions", networkHandler.ListRegions)
				networkGroup.GET("/zones", networkHandler.ListZones)

				// 高级网络资源路由
				networkGroup.POST("/vpc-interconnects", networkHandler.CreateVPCInterconnect)
				networkGroup.GET("/vpc-interconnects", networkHandler.ListVPCInterconnects)
				networkGroup.DELETE("/vpc-interconnects/:id", networkHandler.DeleteVPCInterconnect)

				networkGroup.POST("/vpc-peerings", networkHandler.CreateVPCPeering)
				networkGroup.GET("/vpc-peerings", networkHandler.ListVPCPeerings)
				networkGroup.DELETE("/vpc-peerings/:id", networkHandler.DeleteVPCPeering)

				networkGroup.POST("/route-tables", networkHandler.CreateRouteTable)
				networkGroup.GET("/route-tables", networkHandler.ListRouteTables)
				networkGroup.DELETE("/route-tables/:id", networkHandler.DeleteRouteTable)

				networkGroup.POST("/l2-networks", networkHandler.CreateL2Network)
				networkGroup.GET("/l2-networks", networkHandler.ListL2Networks)
				networkGroup.DELETE("/l2-networks/:id", networkHandler.DeleteL2Network)
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
			authSourceGroup.POST("/:id/sync", authSourceHandler.Sync)
			authSourceGroup.POST("/test-ldap-users", authSourceHandler.TestLDAPUsers)
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
			domainGroup.GET("/:id/projects", domainHandler.GetProjects)
			domainGroup.GET("/:id/roles", domainHandler.GetRoles)
			domainGroup.GET("/:id/cloud-accounts", domainHandler.GetCloudAccounts)
			domainGroup.GET("/:id/operation-logs", domainHandler.GetOperationLogs)
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
			projectGroup.PUT("/:id/manager", projectHandler.SetManager)
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
			userGroup.POST("/:id/projects", userHandler.AssignUserToProject)
			userGroup.DELETE("/:id/projects", userHandler.RemoveUserFromProject)
			userGroup.GET("/:id/projects", userHandler.GetUserProjects)
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
			groupGroup.GET("/:id/projects", groupHandler.GetProjects)
			groupGroup.POST("/:id/projects", groupHandler.AddProject)
			groupGroup.DELETE("/:id/projects", groupHandler.RemoveProject)
		}

		// 角色路由
		roleHandler := handler.NewRoleHandler(db, logger)
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

		// 权限路由
		permissionHandler := handler.NewPermissionHandler(db, logger)
		permissionGroup := v1.Group("/permissions")
		{
			permissionGroup.GET("", permissionHandler.List)
			permissionGroup.GET("/:id", permissionHandler.Get)
			permissionGroup.POST("", permissionHandler.Create)
			permissionGroup.PUT("/:id", permissionHandler.Update)
			permissionGroup.DELETE("/:id", permissionHandler.Delete)
			permissionGroup.POST("/:id/enable", permissionHandler.Enable)
			permissionGroup.POST("/:id/disable", permissionHandler.Disable)
			permissionGroup.GET("/:id/roles", permissionHandler.GetRolePermissions)
			permissionGroup.POST("/role/:role_id/assign/:permission_id", permissionHandler.AssignToRole)
			permissionGroup.DELETE("/role/:role_id/remove/:permission_id", permissionHandler.RemoveFromRole)
			permissionGroup.GET("/resource-action", permissionHandler.GetPermissionsForResourceAction)
		}

		// IAM 综合管理路由
		iamHandler := handler.NewIAMHandler(db, logger)
		iamGroup := v1.Group("/iam")
		{
			iamGroup.GET("/users-with-roles", iamHandler.ListUsersWithRoles)
		}

		// 消息中心路由
		messageHandler := handler.NewMessageHandler(db, logger)
		messageGroup := v1.Group("/messages")
		{
			messageGroup.GET("", messageHandler.List)
			messageGroup.GET("/unread-count", messageHandler.GetUnreadCount)
			messageGroup.GET("/:id", messageHandler.Get)
			messageGroup.PUT("/:id/read", messageHandler.MarkRead)
			messageGroup.POST("/mark-all-read", messageHandler.MarkAllRead)
			messageGroup.DELETE("/:id", messageHandler.Delete)
		}

		// 通知渠道路由
		notifChannelHandler := handler.NewNotificationChannelHandler(db, logger)
		notifChannelGroup := v1.Group("/notification-channels")
		{
			notifChannelGroup.GET("", notifChannelHandler.List)
			notifChannelGroup.GET("/:id", notifChannelHandler.Get)
			notifChannelGroup.POST("", notifChannelHandler.Create)
			notifChannelGroup.PUT("/:id", notifChannelHandler.Update)
			notifChannelGroup.DELETE("/:id", notifChannelHandler.Delete)
			notifChannelGroup.POST("/:id/enable", notifChannelHandler.Enable)
			notifChannelGroup.POST("/:id/disable", notifChannelHandler.Disable)
			notifChannelGroup.POST("/:id/test", notifChannelHandler.Test)
		}

		// 机器人路由
		robotHandler := handler.NewRobotHandler(db, logger)
		robotGroup := v1.Group("/robots")
		{
			robotGroup.GET("", robotHandler.List)
			robotGroup.GET("/:id", robotHandler.Get)
			robotGroup.POST("", robotHandler.Create)
			robotGroup.PUT("/:id", robotHandler.Update)
			robotGroup.DELETE("/:id", robotHandler.Delete)
			robotGroup.POST("/:id/enable", robotHandler.Enable)
			robotGroup.POST("/:id/disable", robotHandler.Disable)
			robotGroup.POST("/:id/test", robotHandler.Test)
		}

		// 接收人路由
		receiverHandler := handler.NewReceiverHandler(db, logger)
		receiverGroup := v1.Group("/receivers")
		{
			receiverGroup.GET("", receiverHandler.List)
			receiverGroup.GET("/:id", receiverHandler.Get)
			receiverGroup.POST("", receiverHandler.Create)
			receiverGroup.PUT("/:id", receiverHandler.Update)
			receiverGroup.DELETE("/:id", receiverHandler.Delete)
			receiverGroup.POST("/:id/enable", receiverHandler.Enable)
			receiverGroup.POST("/:id/disable", receiverHandler.Disable)
			receiverGroup.GET("/:id/channels", receiverHandler.GetChannels)
			receiverGroup.POST("/:id/channels", receiverHandler.SetChannels)
			receiverGroup.GET("/:id/with-channels", receiverHandler.GetWithChannels) // Get receiver with notification channels
		}

		// 消息订阅路由
		subscriptionHandler := handler.NewMessageSubscriptionHandler(db, logger)
		subscriptionGroup := v1.Group("/subscriptions")
		{
			subscriptionGroup.GET("", subscriptionHandler.List)
			subscriptionGroup.GET("/:id", subscriptionHandler.Get)
			subscriptionGroup.POST("", subscriptionHandler.Create)
			subscriptionGroup.PUT("/:id", subscriptionHandler.Update)
			subscriptionGroup.DELETE("/:id", subscriptionHandler.Delete)
		}
		v1.GET("/message-types", subscriptionHandler.ListMessageTypes)

		// 同步策略路由
		syncPolicyHandler := handler.NewSyncPolicyHandler(db, logger)
		syncPolicyGroup := v1.Group("/sync-policies")
		{
			syncPolicyGroup.GET("", syncPolicyHandler.List)
			syncPolicyGroup.GET("/:id", syncPolicyHandler.Get)
			syncPolicyGroup.POST("", syncPolicyHandler.Create)
			syncPolicyGroup.PUT("/:id", syncPolicyHandler.Update)
			syncPolicyGroup.DELETE("/:id", syncPolicyHandler.Delete)
			syncPolicyGroup.POST("/:id/status", syncPolicyHandler.UpdateStatus)
		}

		// 定时任务路由
		scheduledTaskHandler := handler.NewScheduledTaskHandler(db, logger)
		scheduledTaskGroup := v1.Group("/scheduled-tasks")
		{
			scheduledTaskGroup.GET("", scheduledTaskHandler.List)
			scheduledTaskGroup.GET("/:id", scheduledTaskHandler.Get)
			scheduledTaskGroup.POST("", scheduledTaskHandler.Create)
			scheduledTaskGroup.PUT("/:id", scheduledTaskHandler.Update)
			scheduledTaskGroup.DELETE("/:id", scheduledTaskHandler.Delete)
			scheduledTaskGroup.POST("/:id/status", scheduledTaskHandler.UpdateStatus)
		}

		// 操作日志路由
		operationLogHandler := handler.NewOperationLogHandler(service.NewOperationLogService(db))
		operationLogGroup := v1.Group("/operation-logs")
		{
			operationLogGroup.GET("", operationLogHandler.GetOperationLogs)
			operationLogGroup.GET("/:resource_type/:resource_id", operationLogHandler.GetResourceOperationLogs)
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
			hashedPassword, err := pkgutils.HashPassword(adminUser.Password)
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
	} else {
		// 用户已存在：检查密码是否为有效 bcrypt 哈希，不是则修复（兼容旧数据）
		if len(existingUser.Password) < 4 || existingUser.Password[:4] != "$2a$" && existingUser.Password[:4] != "$2b$" {
			hashedPassword, err := pkgutils.HashPassword("admin123")
			if err != nil {
				return fmt.Errorf("failed to hash password: %v", err)
			}
			if err := db.Model(&existingUser).Update("password", hashedPassword).Error; err != nil {
				return fmt.Errorf("failed to fix admin password: %v", err)
			}
		}
	}

	return nil
}

func initDatabase(cfg DatabaseConfig) (*gorm.DB, *sql.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, err
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)

	return db, sqlDB, nil
}
