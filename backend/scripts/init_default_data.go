//go:build ignore

package main

import (
	"context"
	"fmt"
	"log"

	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
		Description: "系统管理员角色",
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

	roleService := service.NewRoleService(db)
	return nil
}
