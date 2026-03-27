package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ============= 域管理 =============

// Domain 域（租户）
type Domain struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"uniqueIndex;not null;size:100" json:"name"`
	Description string         `gorm:"size:500" json:"description"`
	Enabled     bool           `gorm:"default:true" json:"enabled"`
	ParentID    *uint          `gorm:"index" json:"parent_id"` // 支持域层级
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Domain) TableName() string {
	return "domains"
}

// ============= 项目管理 =============

// Project 项目
type Project struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"uniqueIndex;not null;size:100" json:"name"`
	Description string         `gorm:"size:500" json:"description"`
	DomainID    uint           `gorm:"index;not null" json:"domain_id"`
	ParentID    *uint          `gorm:"index" json:"parent_id"` // 支持项目层级
	Enabled     bool           `gorm:"default:true" json:"enabled"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Project) TableName() string {
	return "projects"
}

// ============= 用户管理 =============

// User 用户
type User struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Name           string         `gorm:"uniqueIndex;not null;size:100" json:"name"`
	DisplayName    string         `gorm:"size:100" json:"display_name"`
	Email          string         `gorm:"size:255;index" json:"email"`
	Phone          string         `gorm:"size:20" json:"phone"`
	Password       string         `gorm:"size:255;not null" json:"-"` // 加密存储
	DomainID       uint           `gorm:"index;not null" json:"domain_id"`
	Enabled        bool           `gorm:"default:true" json:"enabled"`
	MFAEnabled     bool           `gorm:"default:false" json:"mfa_enabled"`
	MFASecret      string         `gorm:"size:255" json:"-"` // MFA 密钥
	LastLoginAt    *time.Time     `json:"last_login_at"`
	LastLoginIP    string         `gorm:"size:50" json:"last_login_ip"`
	PasswordExpire *time.Time     `json:"password_expire"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "users"
}

// ============= 用户组管理 =============

// Group 用户组
type Group struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"uniqueIndex;not null;size:100" json:"name"`
	Description string         `gorm:"size:500" json:"description"`
	DomainID    uint           `gorm:"index;not null" json:"domain_id"`
	ParentID    *uint          `gorm:"index" json:"parent_id"` // 支持组层级
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Group) TableName() string {
	return "groups"
}

// UserGroup 用户组关联
type UserGroup struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;uniqueIndex:user_group" json:"user_id"`
	GroupID   uint      `gorm:"index;uniqueIndex:user_group" json:"group_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (UserGroup) TableName() string {
	return "user_groups"
}

// ============= 角色管理 =============

// Role 角色
type Role struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"uniqueIndex;not null;size:100" json:"name"`
	DisplayName string         `gorm:"size:100" json:"display_name"`
	Description string         `gorm:"size:500" json:"description"`
	DomainID    uint           `gorm:"index" json:"domain_id"` // 为空表示系统角色
	Type        string         `gorm:"type:varchar(20);default:custom" json:"type"` // system/custom
	Enabled     bool           `gorm:"default:true" json:"enabled"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Role) TableName() string {
	return "roles"
}

// ============= 权限管理 =============

// Permission 权限
type Permission struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"uniqueIndex;not null;size:100" json:"name"`
	DisplayName string         `gorm:"size:100" json:"display_name"`
	Description string         `gorm:"size:500" json:"description"`
	Resource    string         `gorm:"type:varchar(50);index" json:"resource"` // 资源类型
	Action      string         `gorm:"type:varchar(50);index" json:"action"`   // 操作类型
	Type        string         `gorm:"type:varchar(20);default:custom" json:"type"` // system/custom
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Permission) TableName() string {
	return "permissions"
}

// RolePermission 角色权限关联
type RolePermission struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	RoleID       uint      `gorm:"index;uniqueIndex:role_permission" json:"role_id"`
	PermissionID uint      `gorm:"index;uniqueIndex:role_permission" json:"permission_id"`
	CreatedAt    time.Time `json:"created_at"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}

// ============= 角色分配 =============

// UserRole 用户角色关联（域级别）
type UserRole struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;uniqueIndex:user_role_domain" json:"user_id"`
	RoleID    uint      `gorm:"index;uniqueIndex:user_role_domain" json:"role_id"`
	DomainID  uint      `gorm:"index;uniqueIndex:user_role_domain" json:"domain_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (UserRole) TableName() string {
	return "user_roles"
}

// ProjectUserRole 项目用户角色关联
type ProjectUserRole struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;uniqueIndex:project_user_role" json:"user_id"`
	RoleID    uint      `gorm:"index;uniqueIndex:project_user_role" json:"role_id"`
	ProjectID uint      `gorm:"index;uniqueIndex:project_user_role" json:"project_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (ProjectUserRole) TableName() string {
	return "project_user_roles"
}

// GroupRole 组角色关联（域级别）
type GroupRole struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	GroupID  uint      `gorm:"index;uniqueIndex:group_role_domain" json:"group_id"`
	RoleID   uint      `gorm:"index;uniqueIndex:group_role_domain" json:"role_id"`
	DomainID uint      `gorm:"index;uniqueIndex:group_role_domain" json:"domain_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (GroupRole) TableName() string {
	return "group_roles"
}

// ============= 认证源管理 =============

// AuthSource 认证源
type AuthSource struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"uniqueIndex;not null;size:100" json:"name"`
	Description  string         `gorm:"size:500" json:"description"`
	Type         string         `gorm:"type:varchar(20);not null" json:"type"` // ldap/oidc/saml/local
	Enabled      bool           `gorm:"default:true" json:"enabled"`
	Config       datatypes.JSON `gorm:"type:json" json:"config"` // 认证源配置（加密）
	AutoCreate   bool           `gorm:"default:false" json:"auto_create"` // 自动创建用户
	DefaultRole  *uint          `gorm:"index" json:"default_role_id"` // 默认角色
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (AuthSource) TableName() string {
	return "auth_sources"
}

// ============= 安全告警 =============

// SecurityAlert 安全告警
type SecurityAlert struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Type       string         `gorm:"type:varchar(50);index" json:"type"` // login_failed/password_expired/mfa_disabled等
	Level      string         `gorm:"type:varchar(20);default:medium" json:"level"` // low/medium/high/critical
	Title      string         `gorm:"size:200" json:"title"`
	Message    string         `gorm:"type:text" json:"message"`
	UserID     *uint          `gorm:"index" json:"user_id"`
	SourceIP   string         `gorm:"size:50" json:"source_ip"`
	Status     string         `gorm:"type:varchar(20);default:active" json:"status"` // active/resolved/ignored
	HandledAt  *time.Time     `json:"handled_at"`
	HandledBy  *uint          `json:"handled_by"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

func (SecurityAlert) TableName() string {
	return "security_alerts"
}

// ============= 消息中心 =============

// MessageType 消息类型
type MessageType struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"uniqueIndex;not null;size:100" json:"name"`
	DisplayName string         `gorm:"size:100" json:"display_name"`
	Description string         `gorm:"size:500" json:"description"`
	Enabled     bool           `gorm:"default:true" json:"enabled"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func (MessageType) TableName() string {
	return "message_types"
}

// Message 消息
type Message struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	TypeID     uint           `gorm:"index" json:"type_id"`
	Title      string         `gorm:"size:200" json:"title"`
	Content    string         `gorm:"type:text" json:"content"`
	Level      string         `gorm:"type:varchar(20);default:info" json:"level"` // info/warning/error
	SenderID   uint           `gorm:"index" json:"sender_id"`
	ReceiverID uint           `gorm:"index" json:"receiver_id"`
	Read       bool           `gorm:"default:false" json:"read"`
	ReadAt     *time.Time     `json:"read_at"`
	CreatedAt  time.Time      `json:"created_at"`
}

func (Message) TableName() string {
	return "messages"
}

// ============= 通知渠道 =============

// NotificationChannel 通知渠道
type NotificationChannel struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"uniqueIndex;not null;size:100" json:"name"`
	Type        string         `gorm:"type:varchar(20);not null" json:"type"` // email/sms/webhook/dingtalk/wechat
	Description string         `gorm:"size:500" json:"description"`
	Config      datatypes.JSON `gorm:"type:json" json:"config"` // 渠道配置（加密）
	Enabled     bool           `gorm:"default:true" json:"enabled"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (NotificationChannel) TableName() string {
	return "notification_channels"
}

// ============= 消息订阅 =============

// MessageSubscription 消息订阅
type MessageSubscription struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"index;uniqueIndex:user_message_type" json:"user_id"`
	MessageTypeID uint    `gorm:"index;uniqueIndex:user_message_type" json:"message_type_id"`
	Email       bool      `gorm:"default:true" json:"email"`
	SMS         bool      `gorm:"default:false" json:"sms"`
	Webhook     bool      `gorm:"default:false" json:"webhook"`
	Station     bool      `gorm:"default:true" json:"station"` // 站内信
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (MessageSubscription) TableName() string {
	return "message_subscriptions"
}

// ============= 机器人管理 =============

// Robot 机器人
type Robot struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"uniqueIndex;not null;size:100" json:"name"`
	Description string         `gorm:"size:500" json:"description"`
	Type        string         `gorm:"type:varchar(20);not null" json:"type"` // webhook/dingtalk/wechat/feishu
	WebhookURL  string         `gorm:"size:500" json:"webhook_url"`
	Secret      string         `gorm:"size:255" json:"-"` // 签名密钥
	Enabled     bool           `gorm:"default:true" json:"enabled"`
	MessageTypes datatypes.JSON `gorm:"type:json" json:"message_types"` // 订阅的消息类型
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Robot) TableName() string {
	return "robots"
}

// ============= 接收人管理 =============

// Receiver 接收人
type Receiver struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"uniqueIndex;not null;size:100" json:"name"`
	Email     string         `gorm:"size:255" json:"email"`
	Phone     string         `gorm:"size:20" json:"phone"`
	UserID    *uint          `gorm:"uniqueIndex" json:"user_id"` // 关联用户（可选）
	Enabled   bool           `gorm:"default:true" json:"enabled"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Receiver) TableName() string {
	return "receivers"
}
