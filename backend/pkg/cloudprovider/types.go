package cloudprovider

import "time"

// CloudAccountConfig 云账户配置
type CloudAccountConfig struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	ProviderType string            `json:"provider_type"` // alibaba/tencent/aws/azure
	Credentials  map[string]string `json:"credentials"`   // 加密存储
	Region       string            `json:"region"`        // 默认区域
}

// VMStatus 虚拟机状态
type VMStatus string

const (
	VMStatusPending   VMStatus = "Pending"
	VMStatusRunning   VMStatus = "Running"
	VMStatusStopped   VMStatus = "Stopped"
	VMStatusStarting  VMStatus = "Starting"
	VMStatusStopping  VMStatus = "Stopping"
	VMStatusRebooting VMStatus = "Rebooting"
	VMStatusError     VMStatus = "Error"
	VMStatusDeleted   VMStatus = "Deleted"
)

// VirtualMachine 虚拟机
type VirtualMachine struct {
	ID                string            `json:"id"`
	Name              string            `json:"name"`
	Status            VMStatus          `json:"status"`
	InstanceType      string            `json:"instance_type"`
	ImageID           string            `json:"image_id"`
	OSName            string            `json:"os_name"`            // Operating system name
	BillingMethod     string            `json:"billing_method"`   // Billing method (pay-as-you-go, subscription, etc.)
	Platform          string            `json:"platform"`         // Cloud platform (alibaba, tencent, aws, azure)
	ProjectID         string            `json:"project_id"`       // Associated project ID
	VPCID             string            `json:"vpc_id"`
	SubnetID          string            `json:"subnet_id"`
	PrivateIP         string            `json:"private_ip"`
	PublicIP          string            `json:"public_ip"`
	DiskIDs           []string          `json:"disk_ids"`
	SecurityGroups    []string          `json:"security_groups"`
	SecurityGroupNames []string         `json:"security_group_names"` // Names of security groups
	Keypair           string            `json:"keypair"`
	Tags              map[string]string `json:"tags"`
	CreatedAt         time.Time         `json:"created_at"`
	RegionID          string            `json:"region_id"`
	ZoneID            string            `json:"zone_id"`
}

// VMCreateConfig 虚拟机创建配置
type VMCreateConfig struct {
	Name           string            `json:"name"`
	InstanceType   string            `json:"instance_type"`
	ImageID        string            `json:"image_id"`
	VPCID          string            `json:"vpc_id"`
	SubnetID       string            `json:"subnet_id"`
	SecurityGroups []string          `json:"security_groups"`
	DiskSize       int               `json:"disk_size"` // GB
	Keypair        string            `json:"keypair"`
	UserData       string            `json:"user_data"`
	Tags           map[string]string `json:"tags"`
}

// VPC 虚拟私有云
type VPC struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	CIDR        string            `json:"cidr"`
	Description string            `json:"description"`
	Status      string            `json:"status"`
	Tags        map[string]string `json:"tags"`
	CreatedAt   time.Time         `json:"created_at"`
	RegionID    string            `json:"region_id"`
}

// VPCConfig VPC 创建配置
type VPCConfig struct {
	Name        string            `json:"name"`
	CIDR        string            `json:"cidr"`
	Description string            `json:"description"`
	Tags        map[string]string `json:"tags"`
}

// Region 区域
type Region struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

// Zone 可用区
type Zone struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	RegionID string `json:"region_id"`
	Status   string `json:"status"`
}

// CloudInfo 云厂商信息
type CloudInfo struct {
	Provider string   `json:"provider"`
	Version  string   `json:"version"`
	Regions  int      `json:"regions"`
	Services []string `json:"services"`
}

// Subnet 子网
type Subnet struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	VPCID     string            `json:"vpc_id"`
	CIDR      string            `json:"cidr"`
	ZoneID    string            `json:"zone_id"`
	Status    string            `json:"status"`
	Tags      map[string]string `json:"tags"`
	CreatedAt time.Time         `json:"created_at"`
}

// SubnetConfig 子网创建配置
type SubnetConfig struct {
	Name   string            `json:"name"`
	VPCID  string            `json:"vpc_id"`
	CIDR   string            `json:"cidr"`
	ZoneID string            `json:"zone_id"`
	Tags   map[string]string `json:"tags"`
}

// SecurityGroup 安全组
type SecurityGroup struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	VPCID       string            `json:"vpc_id"`
	Rules       []SGRule          `json:"rules"`
	Tags        map[string]string `json:"tags"`
	CreatedAt   time.Time         `json:"created_at"`
}

// SGConfig 安全组创建配置
type SGConfig struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	VPCID       string            `json:"vpc_id"`
	Tags        map[string]string `json:"tags"`
}

// SGRule 安全组规则
type SGRule struct {
	Direction   string `json:"direction"` // ingress/egress
	Protocol    string `json:"protocol"`  // tcp/udp/icmp/all
	PortRange   string `json:"port_range"`
	CIDR        string `json:"cidr"`
	Action      string `json:"action"` // accept/drop
	Description string `json:"description"`
	Priority    int    `json:"priority"`
	RuleID      string `json:"rule_id,omitempty"`
}

// EIP 弹性公网 IP
type EIP struct {
	ID           string    `json:"id"`
	Address      string    `json:"address"`
	Bandwidth    int       `json:"bandwidth"` // Mbps
	Status       string    `json:"status"`    // available/associated
	ResourceID   string    `json:"resource_id"`
	ResourceType string    `json:"resource_type"`
	RegionID     string    `json:"region_id"`
	CreatedAt    time.Time `json:"created_at"`
}

// EIPConfig EIP 创建配置
type EIPConfig struct {
	Bandwidth int               `json:"bandwidth"`
	RegionID  string            `json:"region_id"`
	Tags      map[string]string `json:"tags"`
}

// Disk 云硬盘
type Disk struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Size      int               `json:"size"`   // GB
	Type      string            `json:"type"`   // SSD/HDD/ESSD
	Status    string            `json:"status"` // available/in-use
	VMID      string            `json:"vm_id"`
	ZoneID    string            `json:"zone_id"`
	Tags      map[string]string `json:"tags"`
	CreatedAt time.Time         `json:"created_at"`
}

// DiskConfig 云硬盘创建配置
type DiskConfig struct {
	Name   string            `json:"name"`
	Size   int               `json:"size"` // GB
	Type   string            `json:"type"`
	ZoneID string            `json:"zone_id"`
	Tags   map[string]string `json:"tags"`
}

// Snapshot 快照
type Snapshot struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	DiskID    string    `json:"disk_id"`
	Size      int       `json:"size"` // GB
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// SnapshotFilter 快照过滤条件
type SnapshotFilter struct {
	SnapshotID string
	DiskID     string
	MaxResults int
}

// InstanceType 实例规格
type InstanceType struct {
	Name           string   `json:"name"`
	CPU            int      `json:"cpu"`    // 核心数
	Memory         int      `json:"memory"` // MB
	GPU            int      `json:"gpu"`
	Category       string   `json:"category"` // 通用型/计算型/内存型
	SupportedZones []string `json:"supported_zones"`
}

// LoadBalancer 负载均衡
type LoadBalancer struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Type        string            `json:"type"`   // alb/nlb
	Scheme      string            `json:"scheme"` // internet-facing/internal
	Address     string            `json:"address"`
	VPCID       string            `json:"vpc_id"`
	SubnetID    string            `json:"subnet_id"`
	ListenerIDs []string          `json:"listener_ids"`
	Tags        map[string]string `json:"tags"`
	CreatedAt   time.Time         `json:"created_at"`
}

// LBConfig 负载均衡创建配置
type LBConfig struct {
	Name     string            `json:"name"`
	Type     string            `json:"type"`
	Scheme   string            `json:"scheme"`
	VPCID    string            `json:"vpc_id"`
	SubnetID string            `json:"subnet_id"`
	Tags     map[string]string `json:"tags"`
}

// Listener 监听器
type Listener struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	LBID           string `json:"lb_id"`
	Protocol       string `json:"protocol"` // HTTP/HTTPS/TCP/UDP
	Port           int    `json:"port"`
	BackendGroupID string `json:"backend_group_id"`
	Status         string `json:"status"`
}

// ListenerConfig 监听器创建配置
type ListenerConfig struct {
	Name           string `json:"name"`
	Protocol       string `json:"protocol"`
	Port           int    `json:"port"`
	BackendGroupID string `json:"backend_group_id"`
}

// DNSZone DNS 区域
type DNSZone struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// DNSZoneConfig DNS 区域创建配置
type DNSZoneConfig struct {
	Name string `json:"name"`
}

// DNSRecord DNS 记录
type DNSRecord struct {
	ID     string `json:"id"`
	ZoneID string `json:"zone_id"`
	Name   string `json:"name"`
	Type   string `json:"type"` // A/AAAA/CNAME/MX/TXT/SRV
	Value  string `json:"value"`
	TTL    int    `json:"ttl"`
	Status string `json:"status"`
}

// DNSRecordConfig DNS 记录创建配置
type DNSRecordConfig struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
	TTL   int    `json:"ttl"`
}

// HostTemplateStatus 主机模版状态
type HostTemplateStatus string

const (
	HostTemplateStatusActive   HostTemplateStatus = "Active"
	HostTemplateStatusInactive HostTemplateStatus = "Inactive"
	HostTemplateStatusDraft    HostTemplateStatus = "Draft"
)

// HostTemplate 主机模版
type HostTemplate struct {
	ID             string            `json:"id"`
	Name           string            `json:"name"`
	Description    string            `json:"description"`
	Status         HostTemplateStatus `json:"status"`            // 模版状态
	InstanceType   string            `json:"instance_type"`      // 实例规格
	CPUArch        string            `json:"cpu_arch"`          // CPU架构
	MemorySize     int               `json:"memory_size"`       // 内存大小(MB)
	CPUCount       int               `json:"cpu_count"`         // CPU核心数
	DiskSize       int               `json:"disk_size"`         // 磁盘大小(GB)
	ImageID        string            `json:"image_id"`          // 镜像ID
	OSName         string            `json:"os_name"`           // 操作系统名称
	OSVersion      string            `json:"os_version"`        // 操作系统版本
	VPCID          string            `json:"vpc_id"`            // VPC ID
	SubnetID       string            `json:"subnet_id"`         // 子网ID
	BillingMethod  string            `json:"billing_method"`    // 计费方式
	Platform       string            `json:"platform"`          // 平台
	ProjectID      string            `json:"project_id"`        // 项目ID
	RegionID       string            `json:"region_id"`         // 区域ID
	ZoneID         string            `json:"zone_id"`           // 可用区ID
	Tags           map[string]string `json:"tags"`              // 标签
	CreatedAt      time.Time         `json:"created_at"`        // 创建时间
	UpdatedAt      time.Time         `json:"updated_at"`        // 更新时间
}

// ASGStatus 弹性伸缩组状态
type ASGStatus string

const (
	ASGStatusActive    ASGStatus = "Active"
	ASGStatusInactive  ASGStatus = "Inactive"
	ASGStatusDeleting  ASGStatus = "Deleting"
	ASGStatusError     ASGStatus = "Error"
)

// AutoscalingGroup 弹性伸缩组
type AutoscalingGroup struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`                      // 弹性伸缩组名称
	Description     string            `json:"description"`               // 描述
	Status          ASGStatus         `json:"status"`                    // 伸缩组状态
	HostTemplateID  string            `json:"host_template_id"`          // 主机模版ID
	CurrentCapacity int               `json:"current_capacity"`          // 当前实例数
	DesiredCapacity int               `json:"desired_capacity"`          // 期望实例数
	MinSize         int               `json:"min_size"`                  // 最小实例数
	MaxSize         int               `json:"max_size"`                  // 最大实例数
	Platform        string            `json:"platform"`                  // 平台
	ProjectID       string            `json:"project_id"`                // 项目ID
	RegionID        string            `json:"region_id"`                 // 区域ID
	ZoneID          string            `json:"zone_id"`                   // 可用区ID
	Tags            map[string]string `json:"tags"`                      // 标签
	CreatedAt       time.Time         `json:"created_at"`                // 创建时间
	UpdatedAt       time.Time         `json:"updated_at"`                // 更新时间
}
