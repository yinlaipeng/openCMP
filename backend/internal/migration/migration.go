package migration

import (
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
)

// Migrate 执行数据库迁移
func Migrate(db *gorm.DB) error {
	// 自动迁移表结构
	return db.AutoMigrate(
		// 基础模型
		&model.CloudAccount{},
		&model.Domain{},
		&model.Project{},

		// IAM 模型
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

		// 安全和通知模型
		&model.SecurityAlert{},
		&model.MessageType{},
		&model.Message{},
		&model.NotificationChannel{},
		&model.MessageSubscription{},
		&model.Robot{},
		&model.Receiver{},
		&model.ReceiverChannel{},

		// 多云管理模型
		&model.SyncPolicy{},    // 添加同步策略模型
		&model.ScheduledTask{}, // 添加定时任务模型

		// 操作日志模型
		&model.OperationLog{},

		// 主机模版模型
		&model.HostTemplate{},

		// 弹性伸缩组模型
		&model.AutoscalingGroup{},
	)
}
