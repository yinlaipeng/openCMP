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
