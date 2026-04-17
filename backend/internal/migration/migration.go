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
		&model.Policy{},
		&model.PolicyStatement{},
		&model.RolePolicy{},

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
		&model.SyncPolicy{},    // 同步策略模型
		&model.Rule{},          // 同步规则模型
		&model.RuleTag{},       // 规则标签模型
		&model.ScheduledTask{}, // 定时任务模型
		&model.SyncLog{},       // 同步日志模型

		// 云账户详情子页面模型
		&model.CloudSubscription{},
		&model.CloudUser{},
		&model.CloudUserGroup{},
		&model.CloudProject{},

		// 同步云资源模型
		&model.CloudVM{},
		&model.CloudVPC{},
		&model.CloudSubnet{},
		&model.CloudSecurityGroup{},
		&model.CloudEIP{},
		&model.CloudImage{},
		&model.CloudDisk{},
		&model.CloudSnapshot{},
		&model.CloudRDS{},
		&model.CloudRedis{},

		// 操作日志模型
		&model.OperationLog{},

		// 主机模版模型
		&model.HostTemplate{},

		// 弹性伸缩组模型
		&model.AutoscalingGroup{},

		// 财务模型
		&model.Budget{},
		&model.Bill{},
		&model.Order{},
		&model.RenewalResource{},
		&model.CostAnomaly{},

		// 监控模型
		&model.AlertPolicy{},
		&model.AlertHistory{},
		&model.MonitorResource{},

		// 资源状态变更日志
		&model.ResourceStateLog{},
	)
}
