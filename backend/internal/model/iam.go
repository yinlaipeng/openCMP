package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ============= 策略管理 =============

// Policy 策略（参考 OneCloud 设计）
type Policy struct {
	ID               string         `gorm:"primaryKey;size:64" json:"id"`                 // 策略 ID（UUID）
	Name             string         `gorm:"uniqueIndex;not null;size:100" json:"name"`    // 策略名称
	Description      string         `gorm:"size:500" json:"description"`                  // 策略描述
	Scope            string         `gorm:"type:varchar(20);not null;index" json:"scope"` // system/domain/project
	DomainID         string         `gorm:"size:64;index" json:"domain_id"`               // 域 ID
	ProjectID        string         `gorm:"size:64;index" json:"project_id"`              // 项目 ID
	Policy           datatypes.JSON `gorm:"type:json;not null" json:"policy"`             // 策略内容（JSON 格式）
	IsSystem         bool           `gorm:"default:false;index" json:"is_system"`         // 是否系统策略
	IsPublic         bool           `gorm:"default:false" json:"is_public"`               // 是否公开
	IsEmulated       bool           `gorm:"default:false" json:"is_emulated"`             // 是否预置策略
	Enabled          bool           `gorm:"default:true;index" json:"enabled"`            // 是否启用
	CanDelete        bool           `gorm:"-" json:"can_delete"`                          // 是否可删除（计算字段）
	CanUpdate        bool           `gorm:"-" json:"can_update"`                          // 是否可更新（计算字段）
	DeleteFailReason datatypes.JSON `gorm:"type:json" json:"delete_fail_reason"`          // 删除失败原因
	PendingDeleted   bool           `gorm:"default:false" json:"pending_deleted"`         // 是否待删除
	Deleted          bool           `gorm:"default:false" json:"deleted"`                 // 是否已删除
	PublicScope      string         `gorm:"type:varchar(20)" json:"public_scope"`         // 公开范围
	UpdateVersion    int            `gorm:"default:0" json:"update_version"`              // 更新版本
	CreatedAt        time.Time      `gorm:"index" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"index" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Policy) TableName() string {
	return "policies"
}

// RolePolicy 角色策略关联
type RolePolicy struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	RoleID    uint      `gorm:"index;uniqueIndex:role_policy" json:"role_id"`
	PolicyID  string    `gorm:"index;uniqueIndex:role_policy;size:64" json:"policy_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (RolePolicy) TableName() string {
	return "role_policies"
}

// PolicyStatement 策略语句
type PolicyStatement struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	PolicyID   string         `gorm:"index;size:64" json:"policy_id"`
	Effect     string         `gorm:"type:varchar(20);not null" json:"effect"` // Allow/Deny
	Resource   string         `gorm:"size:255" json:"resource"`                // 资源标识，支持通配符
	Actions    datatypes.JSON `gorm:"type:json" json:"actions"`                // 操作列表，JSON数组
	Conditions datatypes.JSON `gorm:"type:json" json:"conditions"`             // 条件，JSON对象
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (PolicyStatement) TableName() string {
	return "policy_statements"
}

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
	ParentID    *uint          `gorm:"index" json:"parent_id"`  // 支持项目层级
	ManagerID   *uint          `gorm:"index" json:"manager_id"` // 项目管理员ID
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
	Remark         string         `gorm:"size:255" json:"remark"` // 用户备注
	Email          string         `gorm:"size:255;index" json:"email"`
	Phone          string         `gorm:"size:20" json:"phone"`
	Password       string         `gorm:"size:255;not null" json:"-"` // 加密存储
	DomainID       uint           `gorm:"index;not null" json:"domain_id"`
	Enabled        bool           `gorm:"default:true" json:"enabled"`
	ConsoleLogin   bool           `gorm:"default:true" json:"console_login"` // 是否允许控制台登录
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
	DomainID    uint           `gorm:"index" json:"domain_id"`                      // 为空表示系统角色
	Type        string         `gorm:"type:varchar(20);default:custom" json:"type"` // system/custom
	Enabled     bool           `gorm:"default:true" json:"enabled"`
	IsPublic    bool           `gorm:"default:false" json:"is_public"` // 是否公开
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Role) TableName() string {
	return "roles"
}

// Permission 权限
type Permission struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"uniqueIndex;not null;size:100" json:"name"`    // 权限名称，如 "user.create"
	DisplayName string         `gorm:"size:100" json:"display_name"`                 // 显示名称
	Description string         `gorm:"size:500" json:"description"`                  // 权限描述
	Type        string         `gorm:"type:varchar(20);default:custom" json:"type"`  // 类型：system/custom
	Resource    string         `gorm:"size:100;not null;index" json:"resource"`      // 资源类型，如 "user", "vm", "project"
	Action      string         `gorm:"size:100;not null;index" json:"action"`        // 操作类型，如 "list", "create", "update", "delete"
	Scope       string         `gorm:"type:varchar(20);not null;index" json:"scope"` // system/domain/project
	DomainID    *uint          `gorm:"index" json:"domain_id"`                       // 域 ID（仅当 scope=domain 时）
	ProjectID   *uint          `gorm:"index" json:"project_id"`                      // 项目 ID（仅当 scope=project 时）
	Enabled     bool           `gorm:"default:true" json:"enabled"`                  // 是否启用
	IsPublic    bool           `gorm:"default:false" json:"is_public"`               // 是否公开
	Conditions  datatypes.JSON `gorm:"type:json" json:"conditions"`                  // 条件限制
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
	ID        uint      `gorm:"primaryKey" json:"id"`
	GroupID   uint      `gorm:"index;uniqueIndex:group_role_domain" json:"group_id"`
	RoleID    uint      `gorm:"index;uniqueIndex:group_role_domain" json:"role_id"`
	DomainID  uint      `gorm:"index;uniqueIndex:group_role_domain" json:"domain_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (GroupRole) TableName() string {
	return "group_roles"
}

// GroupProject 组项目关联
type GroupProject struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	GroupID   uint      `gorm:"index;uniqueIndex:group_project" json:"group_id"`
	ProjectID uint      `gorm:"index;uniqueIndex:group_project" json:"project_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (GroupProject) TableName() string {
	return "group_projects"
}

// ============= 认证源管理 =============

// AuthSource 认证源
type AuthSource struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"uniqueIndex;not null;size:100" json:"name"`
	Description string         `gorm:"size:500" json:"description"`
	Type        string         `gorm:"type:varchar(20);not null" json:"type"`        // ldap/oidc/saml/local
	Scope       string         `gorm:"type:varchar(20);default:system" json:"scope"` // system/domain
	DomainID    *uint          `gorm:"index" json:"domain_id"`                       // 域 ID（仅当 scope=domain 时）
	Enabled     bool           `gorm:"default:true" json:"enabled"`
	Config      datatypes.JSON `gorm:"type:json" json:"config"`          // 认证源配置（加密）
	AutoCreate  bool           `gorm:"default:false" json:"auto_create"` // 自动创建用户
	DefaultRole *uint          `gorm:"index" json:"default_role_id"`     // 默认角色
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (AuthSource) TableName() string {
	return "auth_sources"
}

// ============= 安全告警 =============

// SecurityAlert 安全告警
type SecurityAlert struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Type      string     `gorm:"type:varchar(50);index" json:"type"`           // login_failed/password_expired/mfa_disabled等
	Level     string     `gorm:"type:varchar(20);default:medium" json:"level"` // low/medium/high/critical
	Title     string     `gorm:"size:200" json:"title"`
	Message   string     `gorm:"type:text" json:"message"`
	UserID    *uint      `gorm:"index" json:"user_id"`
	SourceIP  string     `gorm:"size:50" json:"source_ip"`
	Status    string     `gorm:"type:varchar(20);default:active" json:"status"` // active/resolved/ignored
	HandledAt *time.Time `json:"handled_at"`
	HandledBy *uint      `json:"handled_by"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (SecurityAlert) TableName() string {
	return "security_alerts"
}

// ============= 消息中心 =============

// MessageType 消息类型（参考 Cloudpods Topic 设计）
type MessageType struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Name          string    `gorm:"uniqueIndex;not null;size:100" json:"name"`           // 消息类型名称（如 user_lock）
	DisplayName   string    `gorm:"size:100" json:"display_name"`                         // 显示名称
	Description   string    `gorm:"size:500" json:"description"`                          // 描述
	Type          string    `gorm:"type:varchar(30);default:'resource'" json:"type"`     // 类型：security/resource/automated_process
	Enabled       bool      `gorm:"default:true" json:"enabled"`                          // 是否启用
	TitleCN       string    `gorm:"type:text" json:"title_cn"`                            // 中文标题模板
	TitleEN       string    `gorm:"type:text" json:"title_en"`                            // 英文标题模板
	ContentCN     string    `gorm:"type:text" json:"content_cn"`                          // 中文内容模板
	ContentEN     string    `gorm:"type:text" json:"content_en"`                          // 英文内容模板
	ResourceTypes string    `gorm:"type:text" json:"resource_types"`                      // 关联资源类型（JSON数组）
	GroupKeys     string    `gorm:"type:text" json:"group_keys"`                          // 分组键（JSON数组）
	AdvanceDays   string    `gorm:"type:text" json:"advance_days"`                        // 提前天数（JSON数组，资源到期提醒）
	IsSystem      bool      `gorm:"default:false" json:"is_system"`                       // 是否系统内置
	CanDelete     bool      `gorm:"default:true" json:"can_delete"`                       // 是否可删除
	CanUpdate     bool      `gorm:"default:true" json:"can_update"`                       // 是否可更新
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (MessageType) TableName() string {
	return "message_types"
}

// TopicReceiver 消息类型接收人（参考 Cloudpods）
type TopicReceiver struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	TopicID       uint      `gorm:"index;not null" json:"topic_id"`                     // 消息类型ID
	ReceiverType  string    `gorm:"type:varchar(20);not null" json:"receiver_type"`     // 接收人类型：user/group/role
	ReceiverID    uint      `gorm:"not null" json:"receiver_id"`                        // 接收人ID
	ReceiverName  string    `gorm:"size:100" json:"receiver_name"`                      // 接收人名称（缓存字段）
	Inbox         bool      `gorm:"default:true" json:"inbox"`                          // 站内信渠道
	Email         bool      `gorm:"default:false" json:"email"`                         // 邮件渠道
	Wechat        bool      `gorm:"default:false" json:"wechat"`                        // 企业微信渠道
	Dingtalk      bool      `gorm:"default:false" json:"dingtalk"`                      // 钉钉渠道
	Webhook       bool      `gorm:"default:false" json:"webhook"`                       // Webhook渠道
	Enabled       bool      `gorm:"default:true" json:"enabled"`                        // 是否启用
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (TopicReceiver) TableName() string {
	return "topic_receivers"
}

// Message 消息
type Message struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	TypeID     uint       `gorm:"index" json:"type_id"`
	Title      string     `gorm:"size:200" json:"title"`
	Content    string     `gorm:"type:text" json:"content"`
	Level      string     `gorm:"type:varchar(20);default:info" json:"level"` // info/warning/error
	SenderID   uint       `gorm:"index" json:"sender_id"`
	ReceiverID uint       `gorm:"index" json:"receiver_id"`
	Read       bool       `gorm:"default:false" json:"read"`
	ReadAt     *time.Time `json:"read_at"`
	CreatedAt  time.Time  `json:"created_at"`
}

func (Message) TableName() string {
	return "messages"
}

// ============= 通知渠道 =============

// NotificationChannel 通知渠道
type NotificationChannel struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"uniqueIndex;not null;size:100" json:"name"`
	Type        string         `gorm:"type:varchar(20);not null" json:"type"` // email/sms/webhook/dingtalk/wechat/feishu/lark
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
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"index;uniqueIndex:user_message_type" json:"user_id"`
	MessageTypeID uint      `gorm:"index;uniqueIndex:user_message_type" json:"message_type_id"`
	Email         bool      `gorm:"default:true" json:"email"`
	SMS           bool      `gorm:"default:false" json:"sms"`
	Webhook       bool      `gorm:"default:false" json:"webhook"`
	Station       bool      `gorm:"default:true" json:"station"` // 站内信
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (MessageSubscription) TableName() string {
	return "message_subscriptions"
}

// ============= 机器人管理 =============

// Robot 机器人
type Robot struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"uniqueIndex;not null;size:100" json:"name"`
	Description  string         `gorm:"size:500" json:"description"`
	Type         string         `gorm:"type:varchar(20);not null;index" json:"type"` // webhook/dingtalk/feishu/workwx
	WebhookURL   string         `gorm:"size:500" json:"webhook_url"`                 // 机器人Webhook地址
	Secret       string         `gorm:"size:255" json:"secret,omitempty"`            // 签名密钥（可选）
	Header       string         `gorm:"size:500" json:"header,omitempty"`            // Webhook类型：请求头（可选）
	Body         string         `gorm:"type:text" json:"body,omitempty"`             // Webhook类型：请求体模板（可选）
	MsgKey       string         `gorm:"size:100" json:"msg_key,omitempty"`           // Webhook类型：消息键（可选）
	SecretKey    string         `gorm:"size:255" json:"secret_key,omitempty"`        // Webhook类型：密钥（可选）
	Enabled      bool           `gorm:"default:true;index" json:"enabled"`
	Status       string         `gorm:"type:varchar(20);default:'ready'" json:"status"` // ready/creating/updating/error
	ProjectID    *uint          `gorm:"index" json:"project_id"`                      // 所属项目（可选）
	DomainID     *uint          `gorm:"index" json:"domain_id"`                       // 所属域（可选）
	SharedScope  string         `gorm:"type:varchar(20)" json:"shared_scope"`         // 共享范围
	MessageTypes datatypes.JSON `gorm:"type:json" json:"message_types"`               // 订阅的消息类型
	CreatedAt    time.Time      `gorm:"index" json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// Associations
	Domain  *Domain  `gorm:"foreignKey:DomainID" json:"domain,omitempty"`
	Project *Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
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
	DomainID  uint           `gorm:"index" json:"domain_id"`     // 所属域
	Enabled   bool           `gorm:"default:true" json:"enabled"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Associations
	Domain               Domain                 `gorm:"foreignKey:DomainID" json:"domain"`
	NotificationChannels []*NotificationChannel `gorm:"many2many:receiver_channels;" json:"notification_channels"`
}

// ReceiverChannel 接收人通知渠道关联
type ReceiverChannel struct {
	ID                    uint      `gorm:"primaryKey" json:"id"`
	ReceiverID            uint      `gorm:"index;uniqueIndex:idx_receiver_channel" json:"receiver_id"`
	NotificationChannelID uint      `gorm:"index;uniqueIndex:idx_receiver_channel" json:"notification_channel_id"`
	Enabled               bool      `gorm:"default:true" json:"enabled"`
	CreatedAt             time.Time `json:"created_at"`
}

func (ReceiverChannel) TableName() string {
	return "receiver_channels"
}
