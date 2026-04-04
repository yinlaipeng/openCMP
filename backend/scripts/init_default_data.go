//go:build ignore

package main

import (
	"context"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gopkg.in/yaml.v3"
	"os"

	"github.com/opencmp/opencmp/internal/model"
	"github.com/opencmp/opencmp/internal/service"
	"github.com/opencmp/opencmp/pkg/utils"
)

type Config struct {
	Database DatabaseConfig `yaml:"database"`
}

type DatabaseConfig struct {
	DSN          string `yaml:"dsn"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxOpenConns int    `yaml:"max_open_conns"`
}

func main() {
	configPath := "configs/config.yaml"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	// 加载配置
	cfg, err := loadConfig(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// 初始化数据库
	db, err := initDatabase(cfg.Database)
	if err != nil {
		log.Fatalf("failed to init database: %v", err)
	}

	// 初始化默认数据
	if err := initDefaultData(db); err != nil {
		log.Fatalf("failed to init default data: %v", err)
	}

	fmt.Println("Default data initialized successfully!")
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
			fmt.Println("Created default domain")
		} else {
			return fmt.Errorf("failed to check default domain: %v", err)
		}
	} else {
		defaultDomain = &existingDomain
		fmt.Println("Default domain already exists")
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
			fmt.Println("Created admin role")
		} else {
			return fmt.Errorf("failed to check admin role: %v", err)
		}
	} else {
		adminRole = &existingRole
		fmt.Println("Admin role already exists")
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
			fmt.Println("Created admin user")
			
			// 为管理员用户分配管理员角色
			userService := service.NewUserService(db)
			if err := userService.AssignUserRole(ctx, adminUser.ID, adminRole.ID, defaultDomain.ID); err != nil {
				return fmt.Errorf("failed to assign admin role to admin user: %v", err)
			}
			fmt.Println("Assigned admin role to admin user")
		} else {
			return fmt.Errorf("failed to check admin user: %v", err)
		}
	} else {
		fmt.Println("Admin user already exists")
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
					return fmt.Errorf("failed to create permission %s: %v", perm.Name, err)
				}
				fmt.Printf("Created permission: %s\n", perm.Name)
			} else {
				return fmt.Errorf("failed to check permission %s: %v", perm.Name, err)
			}
		} else {
			fmt.Printf("Permission %s already exists\n", perm.Name)
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
			fmt.Printf("Could not assign permission %s to admin role: %v\n", perm.Name, err)
		} else {
			fmt.Printf("Assigned permission %s to admin role\n", perm.Name)
		}
	}

	return nil
}