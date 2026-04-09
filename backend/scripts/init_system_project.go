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

	// 初始化系统项目
	if err := initSystemProject(db); err != nil {
		log.Fatalf("failed to init system project: %v", err)
	}

	fmt.Println("System project initialized successfully!")
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

func initSystemProject(db *gorm.DB) error {
	ctx := context.Background()

	// 获取默认域（通常是ID为1的域）
	var defaultDomain model.Domain
	if err := db.Where("name = ?", "Default").First(&defaultDomain).Error; err != nil {
		return fmt.Errorf("failed to find default domain: %v", err)
	}

	// 获取管理员用户
	var adminUser model.User
	if err := db.Where("name = ?", "admin").First(&adminUser).Error; err != nil {
		return fmt.Errorf("failed to find admin user: %v", err)
	}

	// 获取管理员角色
	var adminRole model.Role
	if err := db.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		return fmt.Errorf("failed to find admin role: %v", err)
	}

	// 创建或更新系统项目
	projectService := service.NewProjectService(db)
	systemProject := &model.Project{
		Name:        "system",
		Description: "默认系统项目",
		DomainID:    defaultDomain.ID,
		Enabled:     true,
	}

	var existingProject model.Project
	if err := db.Where("name = ? AND domain_id = ?", "system", defaultDomain.ID).First(&existingProject).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			if err := projectService.CreateProject(ctx, systemProject); err != nil {
				return fmt.Errorf("failed to create system project: %v", err)
			}
			fmt.Println("Created system project")
		} else {
			return fmt.Errorf("failed to check system project: %v", err)
		}
	} else {
		fmt.Println("System project already exists")
		// Update with the existing project for further operations
		systemProject = &existingProject
	}

	// 尝试将管理员用户加入系统项目（如果还没加入的话）
	if err := projectService.AddUserToProject(ctx, systemProject.ID, adminUser.ID, adminRole.ID); err != nil {
		// 可能已经存在，输出信息但不视为错误
		fmt.Printf("Could not add admin user to system project (might already exist): %v\n", err)
	} else {
		fmt.Println("Added admin user to system project")
	}

	return nil
}