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
	ID             string            `json:"id"`
	Name           string            `json:"name"`
	Status         VMStatus          `json:"status"`
	InstanceType   string            `json:"instance_type"`
	ImageID        string            `json:"image_id"`
	VPCID          string            `json:"vpc_id"`
	SubnetID       string            `json:"subnet_id"`
	PrivateIP      string            `json:"private_ip"`
	PublicIP       string            `json:"public_ip"`
	DiskIDs        []string          `json:"disk_ids"`
	SecurityGroups []string          `json:"security_groups"`
	Keypair        string            `json:"keypair"`
	Tags           map[string]string `json:"tags"`
	CreatedAt      time.Time         `json:"created_at"`
	RegionID       string            `json:"region_id"`
	ZoneID         string            `json:"zone_id"`
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

// SubnetFilter 子网过滤条件
type SubnetFilter struct {
	VPCID      string
	ZoneID     string
	Tags       map[string]string
	MaxResults int
}

// VPCFilter VPC 过滤条件
type VPCFilter struct {
	Tags       map[string]string
	RegionID   string
	MaxResults int
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
	Direction       string `json:"direction"`  // ingress/egress
	Protocol        string `json:"protocol"`   // tcp/udp/icmp
	PortRange       string `json:"port_range"` // 80/80, 1-65535
	CIDRBlock       string `json:"cidr_block"` // 0.0.0.0/0
	Policy          string `json:"policy"`     // accept/drop
	Priority        int    `json:"priority"`   // 1-100
	Description     string `json:"description"`
	SecurityGroupID string `json:"security_group_id"`
}

// SGFilter 安全组过滤条件
type SGFilter struct {
	VPCID      string
	Tags       map[string]string
	MaxResults int
}

// EIP 弹性公网 IP
type EIP struct {
	ID             string            `json:"id"`
	Name           string            `json:"name"`
	IPAddress      string            `json:"ip_address"`
	Status         string            `json:"status"`
	AssociatedWith string            `json:"associated_with"` // 关联的资源 ID
	Tags           map[string]string `json:"tags"`
	CreatedAt      time.Time         `json:"created_at"`
}

// EIPConfig EIP 创建配置
type EIPConfig struct {
	Name      string            `json:"name"`
	Bandwidth int               `json:"bandwidth"` // Mbps
	Tags      map[string]string `json:"tags"`
}

// EIPFilter EIP 过滤条件
type EIPFilter struct {
	Status     string
	Tags       map[string]string
	MaxResults int
}

// Disk 云硬盘
type Disk struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Size         int               `json:"size"` // GB
	Type         string            `json:"type"`
	Status       string            `json:"status"`
	AttachedVMID string            `json:"attached_vm_id"`
	Tags         map[string]string `json:"tags"`
	CreatedAt    time.Time         `json:"created_at"`
	ZoneID       string            `json:"zone_id"`
}

// DiskConfig 云硬盘创建配置
type DiskConfig struct {
	Name   string            `json:"name"`
	Size   int               `json:"size"` // GB
	Type   string            `json:"type"`
	ZoneID string            `json:"zone_id"`
	Tags   map[string]string `json:"tags"`
}

// DiskFilter 云硬盘过滤条件
type DiskFilter struct {
	ZoneID     string
	Status     string
	Tags       map[string]string
	MaxResults int
}

// Snapshot 快照
type Snapshot struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	DiskID    string            `json:"disk_id"`
	Size      int               `json:"size"` // GB
	Status    string            `json:"status"`
	Tags      map[string]string `json:"tags"`
	CreatedAt time.Time         `json:"created_at"`
}

// SnapshotFilter 快照过滤条件
type SnapshotFilter struct {
	DiskID     string
	Tags       map[string]string
	MaxResults int
}

// InstanceType 实例类型
type InstanceType struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	CPU            int      `json:"cpu"`    // 核心数
	Memory         int      `json:"memory"` // MB
	Generation     string   `json:"generation"`
	Family         string   `json:"family"`
	SupportedZones []string `json:"supported_zones"`
}
