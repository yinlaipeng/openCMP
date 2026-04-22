package model

import (
	"time"

	"gorm.io/datatypes"
)

// CloudVM 云虚拟机（同步后的本地存储）
type CloudVM struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	CloudAccountID uint           `gorm:"index;not null" json:"cloud_account_id"`
	InstanceID     string         `gorm:"size:100;index" json:"instance_id"`       // 云平台实例ID
	Name           string         `gorm:"size:200" json:"name"`
	Status         string         `gorm:"size:50" json:"status"`                  // running/stopped/terminated
	InstanceType   string         `gorm:"size:100" json:"instance_type"`
	ImageID        string         `gorm:"size:100" json:"image_id"`
	OSName         string         `gorm:"size:100" json:"os_name"`
	VPCID          string         `gorm:"size:100" json:"vpc_id"`
	SubnetID       string         `gorm:"size:100" json:"subnet_id"`
	PrivateIP      string         `gorm:"size:50" json:"private_ip"`
	PublicIP       string         `gorm:"size:50" json:"public_ip"`
	RegionID       string         `gorm:"size:100" json:"region_id"`
	ZoneID         string         `gorm:"size:100" json:"zone_id"`
	ProjectID      uint           `gorm:"index" json:"project_id"`                // 项目归属
	Tags           datatypes.JSON `gorm:"type:json" json:"tags"`                  // 资源标签
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

func (CloudVM) TableName() string { return "sync_cloud_vms" }

// CloudVPC 云VPC（同步后的本地存储）
type CloudVPC struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	CloudAccountID uint           `gorm:"index;not null" json:"cloud_account_id"`
	VPCID          string         `gorm:"size:100;index" json:"vpc_id"`
	Name           string         `gorm:"size:200" json:"name"`
	CIDR           string         `gorm:"size:50" json:"cidr"`
	Status         string         `gorm:"size:50" json:"status"`
	RegionID       string         `gorm:"size:100" json:"region_id"`
	ProjectID      uint           `gorm:"index" json:"project_id"`
	Tags           datatypes.JSON `gorm:"type:json" json:"tags"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

func (CloudVPC) TableName() string { return "sync_cloud_vpcs" }

// CloudSubnet 云子网（同步后的本地存储）
type CloudSubnet struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	CloudAccountID uint           `gorm:"index;not null" json:"cloud_account_id"`
	SubnetID       string         `gorm:"size:100;index" json:"subnet_id"`
	Name           string         `gorm:"size:200" json:"name"`
	VPCID          string         `gorm:"size:100;index" json:"vpc_id"`
	VPC            string         `gorm:"size:100" json:"vpc"` // VPC名称
	CIDR           string         `gorm:"size:50" json:"cidr"`
	GuestIPStart   string         `gorm:"size:20" json:"guest_ip_start"` // IP起始地址
	GuestIPEnd     string         `gorm:"size:20" json:"guest_ip_end"`   // IP结束地址
	GuestIPMask    int            `json:"guest_ip_mask"`                 // 子网掩码
	GuestIP6Mask   string         `gorm:"size:50" json:"guest_ip6_mask"` // IPv6地址段
	GuestGateway   string         `gorm:"size:20" json:"guest_gateway"`  // 网关地址
	DNS            string         `gorm:"size:50" json:"dns"`            // DNS服务器
	IsAutoAlloc    bool           `gorm:"default:true" json:"is_auto_alloc"` // 自动分配IP
	IsClassic      bool           `gorm:"default:false" json:"is_classic"` // 是否经典网络类型
	SchedTag       string         `gorm:"size:100" json:"schedtag"` // 调度标签
	Ports          int            `json:"ports"`           // 总端口数
	PortsUsed      int            `json:"ports_used"`      // 已用端口数
	Ports6Used     int            `json:"ports6_used"`     // IPv6已用端口数
	Vnics          int            `json:"vnics"`           // 网卡数
	ZoneID         string         `gorm:"size:100" json:"zone_id"`
	Zone           string         `gorm:"size:100" json:"zone"` // 可用区名称
	RegionID       string         `gorm:"size:100;index" json:"region_id"`
	Region         string         `gorm:"size:100" json:"region"` // 区域名称
	Status         string         `gorm:"size:50" json:"status"`
	ProviderType   string         `gorm:"size:20" json:"provider_type"` // 云平台类型
	AccountName    string         `gorm:"size:100" json:"account_name"` // 云账号名称
	ProjectID      uint           `gorm:"index" json:"project_id"`
	ProjectName    string         `gorm:"size:100" json:"project_name"`
	Tags           datatypes.JSON `gorm:"type:json" json:"tags"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

func (CloudSubnet) TableName() string { return "sync_cloud_subnets" }

// CloudSecurityGroup 云安全组（同步后的本地存储）
type CloudSecurityGroup struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	CloudAccountID  uint           `gorm:"index;not null" json:"cloud_account_id"`
	SecurityGroupID string         `gorm:"size:100;uniqueIndex:sg_account_idx" json:"security_group_id"`
	Name            string         `gorm:"size:200" json:"name"`
	Description     string         `gorm:"size:500" json:"description"`
	VPCID           string         `gorm:"size:100;index" json:"vpc_id"`
	VPC             string         `gorm:"size:100" json:"vpc"` // VPC名称
	GuestCnt        int            `gorm:"default:0" json:"guest_cnt"` // 关联实例数
	PublicScope     string         `gorm:"size:20;default:'system'" json:"public_scope"` // 共享范围
	IsPublic        bool           `gorm:"default:false" json:"is_public"`
	IsSystem        bool           `gorm:"default:false" json:"is_system"`
	IsDefaultVPC    bool           `gorm:"default:false" json:"is_default_vpc"`
	RegionID        string         `gorm:"size:100;index" json:"region_id"`
	Region          string         `gorm:"size:100" json:"region"` // 区域名称
	Status          string         `gorm:"size:50;default:'ready'" json:"status"`
	ProviderType    string         `gorm:"size:20" json:"provider_type"` // 云平台类型
	AccountName     string         `gorm:"size:100" json:"account_name"` // 云账号名称
	ProjectID       uint           `gorm:"index" json:"project_id"`
	ProjectName     string         `gorm:"size:100" json:"project_name"`
	Tags            datatypes.JSON `gorm:"type:json" json:"tags"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

func (CloudSecurityGroup) TableName() string { return "sync_cloud_security_groups" }

// CloudEIP 云弹性公网IP（同步后的本地存储）
type CloudEIP struct {
	ID                  uint           `gorm:"primaryKey" json:"id"`
	CloudAccountID      uint           `gorm:"index;not null" json:"cloud_account_id"`
	EIPID               string         `gorm:"size:100;uniqueIndex:eip_account_idx" json:"eip_id"`
	Name                string         `gorm:"size:200" json:"name"`
	Address             string         `gorm:"size:50" json:"address"` // IP地址
	Bandwidth           int            `json:"bandwidth"` // Mbps
	BillingMethod       string         `gorm:"size:20" json:"billing_method"` // 计费方式
	ResourceID          string         `gorm:"size:100" json:"resource_id"`     // 绑定的资源ID
	ResourceType        string         `gorm:"size:50" json:"resource_type"`    // 绑定的资源类型
	ResourceName        string         `gorm:"size:100" json:"resource_name"`   // 绑定的资源名称
	RegionID            string         `gorm:"size:100;index" json:"region_id"`
	Region              string         `gorm:"size:100" json:"region"` // 区域名称
	Status              string         `gorm:"size:50" json:"status"` // available/in-use
	ProviderType        string         `gorm:"size:20" json:"provider_type"` // 云平台类型
	AccountName         string         `gorm:"size:100" json:"account_name"` // 云账号名称
	ProjectID           uint           `gorm:"index" json:"project_id"`
	ProjectName         string         `gorm:"size:100" json:"project_name"`
	Tags                datatypes.JSON `gorm:"type:json" json:"tags"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
}

func (CloudEIP) TableName() string { return "sync_cloud_eips" }

// CloudImage 云镜像（同步后的本地存储）
type CloudImage struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	CloudAccountID uint           `gorm:"index;not null" json:"cloud_account_id"`
	ImageID        string         `gorm:"size:100;index" json:"image_id"`
	Name           string         `gorm:"size:200" json:"name"`
	OSName         string         `gorm:"size:100" json:"os_name"`
	OSVersion      string         `gorm:"size:50" json:"os_version"`
	Architecture   string         `gorm:"size:50" json:"architecture"` // CPU架构
	Status         string         `gorm:"size:50" json:"status"`
	Tags           datatypes.JSON `gorm:"type:json" json:"tags"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

func (CloudImage) TableName() string { return "sync_cloud_images" }

// CloudRDS 云RDS数据库（同步后的本地存储）
type CloudRDS struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	CloudAccountID uint           `gorm:"index;not null" json:"cloud_account_id"`
	RDSID          string         `gorm:"size:100;index" json:"rds_id"`
	Name           string         `gorm:"size:200" json:"name"`
	Engine         string         `gorm:"size:50" json:"engine"`       // MySQL/PostgreSQL
	EngineVersion  string         `gorm:"size:50" json:"engine_version"`
	InstanceType   string         `gorm:"size:100" json:"instance_type"`
	Status         string         `gorm:"size:50" json:"status"`
	RegionID       string         `gorm:"size:100" json:"region_id"`
	ZoneID         string         `gorm:"size:100" json:"zone_id"`
	ProjectID      uint           `gorm:"index" json:"project_id"`
	Tags           datatypes.JSON `gorm:"type:json" json:"tags"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

func (CloudRDS) TableName() string { return "sync_cloud_rds" }

// CloudRedis 云Redis（同步后的本地存储）
type CloudRedis struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	CloudAccountID uint           `gorm:"index;not null" json:"cloud_account_id"`
	RedisID        string         `gorm:"size:100;index" json:"redis_id"`
	Name           string         `gorm:"size:200" json:"name"`
	InstanceType   string         `gorm:"size:100" json:"instance_type"`
	Status         string         `gorm:"size:50" json:"status"`
	RegionID       string         `gorm:"size:100" json:"region_id"`
	ZoneID         string         `gorm:"size:100" json:"zone_id"`
	ProjectID      uint           `gorm:"index" json:"project_id"`
	Tags           datatypes.JSON `gorm:"type:json" json:"tags"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

func (CloudRedis) TableName() string { return "sync_cloud_redis" }

// KeyPair SSH密钥（同步后的本地存储）
type KeyPair struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	CloudAccountID uint           `gorm:"index;not null" json:"cloud_account_id"`
	KeyPairID      string         `gorm:"size:100;uniqueIndex:keypair_account_idx" json:"keypair_id"`
	Name           string         `gorm:"size:200" json:"name"`
	Description    string         `gorm:"size:500" json:"description"`
	PublicKey      string         `gorm:"size:2000" json:"public_key"`
	Fingerprint    string         `gorm:"size:100" json:"fingerprint"`
	KeyPairType    string         `gorm:"size:20;default:'ssh'" json:"keypair_type"` // ssh/x509
	PublicScope    string         `gorm:"size:20;default:'system'" json:"public_scope"`
	GuestCnt       int            `gorm:"default:0" json:"guest_cnt"` // 关联实例数
	RegionID       string         `gorm:"size:100;index" json:"region_id"`
	Region         string         `gorm:"size:100" json:"region"`
	Status         string         `gorm:"size:50;default:'ready'" json:"status"`
	ProviderType   string         `gorm:"size:20" json:"provider_type"`
	AccountName    string         `gorm:"size:100" json:"account_name"`
	ProjectID      uint           `gorm:"index" json:"project_id"`
	ProjectName    string         `gorm:"size:100" json:"project_name"`
	Tags           datatypes.JSON `gorm:"type:json" json:"tags"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

func (KeyPair) TableName() string { return "sync_cloud_keypairs" }

// CloudNATGateway 云NAT网关（同步后的本地存储）
type CloudNATGateway struct {
	ID                  uint           `gorm:"primaryKey" json:"id"`
	CloudAccountID      uint           `gorm:"index;not null" json:"cloud_account_id"`
	NatGatewayID        string         `gorm:"size:100;uniqueIndex:nat_account_idx" json:"nat_gateway_id"`
	Name                string         `gorm:"size:200" json:"name"`
	Description         string         `gorm:"size:500" json:"description"`
	NatType             string         `gorm:"size:20" json:"nat_type"` // 公网NAT/私网NAT
	Specification       string         `gorm:"size:100" json:"specification"` // 规格
	BillingMethod       string         `gorm:"size:20" json:"billing_method"` // Postpaid/Prepaid
	VpcID               string         `gorm:"size:100" json:"vpc_id"`
	VpcName             string         `gorm:"size:100" json:"vpc_name"`
	SubnetID            string         `gorm:"size:100" json:"subnet_id"`
	SubnetName          string         `gorm:"size:100" json:"subnet_name"`
	EipID               string         `gorm:"size:100" json:"eip_id"`
	EipAddress          string         `gorm:"size:50" json:"eip_address"`
	SnatTableEntries    int            `gorm:"default:0" json:"snat_table_entries"` // SNAT规则数
	DnatTableEntries    int            `gorm:"default:0" json:"dnat_table_entries"` // DNAT规则数
	OwnerDomain         string         `gorm:"size:100" json:"owner_domain"` // 所属域
	RegionID            string         `gorm:"size:100;index" json:"region_id"`
	Region              string         `gorm:"size:100" json:"region"`
	Status              string         `gorm:"size:50" json:"status"` // Available/Creating/Deleting/Error
	ProviderType        string         `gorm:"size:20" json:"provider_type"`
	AccountName         string         `gorm:"size:100" json:"account_name"`
	ProjectID           uint           `gorm:"index" json:"project_id"`
	ProjectName         string         `gorm:"size:100" json:"project_name"`
	Tags                datatypes.JSON `gorm:"type:json" json:"tags"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
}

func (CloudNATGateway) TableName() string { return "sync_cloud_nat_gateways" }

// CloudNATRule NAT规则（同步后的本地存储）
type CloudNATRule struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	NatGatewayID   uint           `gorm:"index;not null" json:"nat_gateway_id"`
	RuleID         string         `gorm:"size:100;index" json:"rule_id"`
	RuleType       string         `gorm:"size:20" json:"rule_type"` // SNAT/DNAT
	Name           string         `gorm:"size:200" json:"name"`
	ExternalIP     string         `gorm:"size:50" json:"external_ip"` // 外部IP
	ExternalPort   string         `gorm:"size:20" json:"external_port"` // 外部端口
	InternalIP     string         `gorm:"size:50" json:"internal_ip"` // 内部IP
	InternalPort   string         `gorm:"size:20" json:"internal_port"` // 内部端口
	Protocol       string         `gorm:"size:20" json:"protocol"` // TCP/UDP/ALL
	Status         string         `gorm:"size:50" json:"status"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

func (CloudNATRule) TableName() string { return "sync_cloud_nat_rules" }

// CloudIPv6Gateway IPv6网关（同步后的本地存储）
type CloudIPv6Gateway struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	CloudAccountID uint           `gorm:"index;not null" json:"cloud_account_id"`
	Ipv6GatewayID  string         `gorm:"size:100;uniqueIndex:ipv6gw_account_idx" json:"ipv6_gateway_id"`
	Name           string         `gorm:"size:200" json:"name"`
	VpcID          string         `gorm:"size:100" json:"vpc_id"`
	VpcName        string         `gorm:"size:100" json:"vpc_name"`
	Specification  string         `gorm:"size:100" json:"specification"` // 规格
	Ipv6Cidr       string         `gorm:"size:50" json:"ipv6_cidr"` // IPv6地址段
	RegionID       string         `gorm:"size:100;index" json:"region_id"`
	Region         string         `gorm:"size:100" json:"region"`
	Status         string         `gorm:"size:50" json:"status"` // Available/Creating/Deleting/Error
	ProviderType   string         `gorm:"size:20" json:"provider_type"`
	AccountName    string         `gorm:"size:100" json:"account_name"`
	ProjectID      uint           `gorm:"index" json:"project_id"`
	ProjectName    string         `gorm:"size:100" json:"project_name"`
	Tags           datatypes.JSON `gorm:"type:json" json:"tags"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

func (CloudIPv6Gateway) TableName() string { return "sync_cloud_ipv6_gateways" }

// CloudDNSZone DNS解析Zone（同步后的本地存储）
type CloudDNSZone struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	CloudAccountID uint           `gorm:"index;not null" json:"cloud_account_id"`
	DnsZoneID      string         `gorm:"size:100;uniqueIndex:dns_account_idx" json:"dns_zone_id"`
	Name           string         `gorm:"size:200" json:"name"`
	VpcCount       int            `gorm:"default:0" json:"vpc_count"` // 关联VPC数
	AttributionScope string        `gorm:"size:20" json:"attribution_scope"` // 归属范围
	RegionID       string         `gorm:"size:100;index" json:"region_id"`
	Region         string         `gorm:"size:100" json:"region"`
	Status         string         `gorm:"size:50" json:"status"` // Available/Creating/Deleting/Error
	ProviderType   string         `gorm:"size:20" json:"provider_type"`
	AccountName    string         `gorm:"size:100" json:"account_name"`
	ProjectID      uint           `gorm:"index" json:"project_id"`
	ProjectName    string         `gorm:"size:100" json:"project_name"`
	Tags           datatypes.JSON `gorm:"type:json" json:"tags"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

func (CloudDNSZone) TableName() string { return "sync_cloud_dns_zones" }

// CloudDNSRecord DNS解析记录
type CloudDNSRecord struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	DnsZoneID  uint      `gorm:"index;not null" json:"dns_zone_id"`
	RecordID   string    `gorm:"size:100;index" json:"record_id"`
	Name       string    `gorm:"size:200" json:"name"` // 记录名
	Type       string    `gorm:"size:20" json:"type"` // A/CNAME/MX/TXT/NS等
	Value      string    `gorm:"size:500" json:"value"` // 记录值
	TTL        int       `json:"ttl"` // TTL时间
	Priority   int       `json:"priority"` // 优先级（MX记录）
	Status     string    `gorm:"size:50" json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (CloudDNSRecord) TableName() string { return "sync_cloud_dns_records" }

// ========== Region/Zone Models ==========

// CloudRegion 云区域（同步后的本地存储）
type CloudRegion struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	CloudAccountID uint           `gorm:"index;not null" json:"cloud_account_id"`
	RegionID       string         `gorm:"size:100;uniqueIndex:region_account_idx" json:"region_id"`
	Name           string         `gorm:"size:200" json:"name"`
	Enabled        bool           `gorm:"default:true" json:"enabled"`
	ZoneCount      int            `gorm:"default:0" json:"zone_count"`
	VPCCount       int            `gorm:"default:0" json:"vpc_count"`
	ServerCount    int            `gorm:"default:0" json:"server_count"`
	ProviderType   string         `gorm:"size:20" json:"provider_type"`
	AccountName    string         `gorm:"size:100" json:"account_name"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

func (CloudRegion) TableName() string { return "sync_cloud_regions" }

// CloudZone 云可用区（同步后的本地存储）
type CloudZone struct {
	ID                   uint           `gorm:"primaryKey" json:"id"`
	CloudAccountID       uint           `gorm:"index;not null" json:"cloud_account_id"`
	ZoneID               string         `gorm:"size:100;uniqueIndex:zone_account_idx" json:"zone_id"`
	Name                 string         `gorm:"size:200" json:"name"`
	RegionID             string         `gorm:"size:100;index" json:"region_id"`
	Region               string         `gorm:"size:100" json:"region"`
	L2NetworkCount       int            `gorm:"default:0" json:"l2_network_count"`
	HostCount            int            `gorm:"default:0" json:"host_count"`
	EnabledHostCount     int            `gorm:"default:0" json:"enabled_host_count"`
	PhysicalHostCount    int            `gorm:"default:0" json:"physical_host_count"`
	EnabledPhysicalHostCount int        `gorm:"default:0" json:"enabled_physical_host_count"`
	Status               string         `gorm:"size:50;default:'ready'" json:"status"`
	ProviderType         string         `gorm:"size:20" json:"provider_type"`
	AccountName          string         `gorm:"size:100" json:"account_name"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
}

func (CloudZone) TableName() string { return "sync_cloud_zones" }

// ========== Global VPC / VPC Interconnect / L2 Network / Route Table Models ==========

// CloudGlobalVPC 全局VPC（同步后的本地存储）
type CloudGlobalVPC struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	CloudAccountID uint           `gorm:"index;not null" json:"cloud_account_id"`
	GlobalVpcID    string         `gorm:"size:100;uniqueIndex:gvpc_account_idx" json:"global_vpc_id"`
	Name           string         `gorm:"size:200" json:"name"`
	Description    string         `gorm:"size:500" json:"description"`
	VPCCount       int            `gorm:"default:0" json:"vpc_count"`
	Status         string         `gorm:"size:50;default:'ready'" json:"status"`
	ProviderType   string         `gorm:"size:20" json:"provider_type"`
	AccountName    string         `gorm:"size:100" json:"account_name"`
	OwnerDomain    string         `gorm:"size:100" json:"owner_domain"`
	ProjectID      uint           `gorm:"index" json:"project_id"`
	ProjectName    string         `gorm:"size:100" json:"project_name"`
	Tags           datatypes.JSON `gorm:"type:json" json:"tags"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

func (CloudGlobalVPC) TableName() string { return "sync_cloud_global_vpcs" }

// CloudVPCInterconnect VPC互联（同步后的本地存储）
type CloudVPCInterconnect struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	CloudAccountID   uint           `gorm:"index;not null" json:"cloud_account_id"`
	InterconnectID   string         `gorm:"size:100;uniqueIndex:interconn_account_idx" json:"interconnect_id"`
	Name             string         `gorm:"size:200" json:"name"`
	VPCCount         int            `gorm:"default:0" json:"vpc_count"`
	AttributionScope string         `gorm:"size:50" json:"attribution_scope"`
	InterconnectType string         `gorm:"size:50" json:"interconnect_type"` // CCN/DCN等
	Bandwidth        int            `gorm:"default:0" json:"bandwidth"` // Mbps
	Status           string         `gorm:"size:50;default:'ready'" json:"status"`
	ProviderType     string         `gorm:"size:20" json:"provider_type"`
	AccountName      string         `gorm:"size:100" json:"account_name"`
	Description      string         `gorm:"size:500" json:"description"`
	ProjectID        uint           `gorm:"index" json:"project_id"`
	ProjectName      string         `gorm:"size:100" json:"project_name"`
	Tags             datatypes.JSON `gorm:"type:json" json:"tags"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

func (CloudVPCInterconnect) TableName() string { return "sync_cloud_vpc_interconnects" }

// CloudL2Network 二层网络（同步后的本地存储）
type CloudL2Network struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	CloudAccountID uint           `gorm:"index;not null" json:"cloud_account_id"`
	L2NetworkID    string         `gorm:"size:100;uniqueIndex:l2net_account_idx" json:"l2_network_id"`
	Name           string         `gorm:"size:200" json:"name"`
	VPCID          string         `gorm:"size:100" json:"vpc_id"`
	VPC            string         `gorm:"size:100" json:"vpc"` // VPC名称
	Bandwidth      int            `gorm:"default:0" json:"bandwidth"` // Mbps
	VlanID         string         `gorm:"size:50" json:"vlan_id"`
	MTU            int            `gorm:"default:1500" json:"mtu"`
	NetworkCount   int            `gorm:"default:0" json:"network_count"`
	Status         string         `gorm:"size:50;default:'ready'" json:"status"`
	ProviderType   string         `gorm:"size:20" json:"provider_type"`
	AccountName    string         `gorm:"size:100" json:"account_name"`
	OwnerDomain    string         `gorm:"size:100" json:"owner_domain"`
	RegionID       string         `gorm:"size:100;index" json:"region_id"`
	Region         string         `gorm:"size:100" json:"region"`
	ProjectID      uint           `gorm:"index" json:"project_id"`
	ProjectName    string         `gorm:"size:100" json:"project_name"`
	Tags           datatypes.JSON `gorm:"type:json" json:"tags"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

func (CloudL2Network) TableName() string { return "sync_cloud_l2_networks" }

// CloudRouteTable 路由表（同步后的本地存储）
type CloudRouteTable struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	CloudAccountID uint           `gorm:"index;not null" json:"cloud_account_id"`
	RouteTableID   string         `gorm:"size:100;uniqueIndex:rt_account_idx" json:"route_table_id"`
	Name           string         `gorm:"size:200" json:"name"`
	VPCID          string         `gorm:"size:100" json:"vpc_id"`
	VPC            string         `gorm:"size:100" json:"vpc"` // VPC名称
	RouteCount     int            `gorm:"default:0" json:"route_count"`
	Status         string         `gorm:"size:50;default:'ready'" json:"status"`
	ProviderType   string         `gorm:"size:20" json:"provider_type"`
	AccountName    string         `gorm:"size:100" json:"account_name"`
	OwnerDomain    string         `gorm:"size:100" json:"owner_domain"`
	RegionID       string         `gorm:"size:100;index" json:"region_id"`
	Region         string         `gorm:"size:100" json:"region"`
	ProjectID      uint           `gorm:"index" json:"project_id"`
	ProjectName    string         `gorm:"size:100" json:"project_name"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

func (CloudRouteTable) TableName() string { return "sync_cloud_route_tables" }

// ========== Load Balancer Models ==========

// CloudLBInstance 负载均衡实例（同步后的本地存储）
type CloudLBInstance struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	CloudAccountID  uint           `gorm:"index;not null" json:"cloud_account_id"`
	LbID            string         `gorm:"size:100;uniqueIndex:lb_account_idx" json:"lb_id"`
	Name            string         `gorm:"size:200" json:"name"`
	Address         string         `gorm:"size:50" json:"address"` // LB地址
	AddressType     string         `gorm:"size:50" json:"address_type"` // intranet/internet
	Specification   string         `gorm:"size:100" json:"specification"`
	Bandwidth       int            `gorm:"default:0" json:"bandwidth"` // Mbps
	VPCID           string         `gorm:"size:100" json:"vpc_id"`
	VPC             string         `gorm:"size:100" json:"vpc"`
	SecurityGroupID string         `gorm:"size:100" json:"security_group_id"`
	SecurityGroup   string         `gorm:"size:100" json:"security_group"`
	ChargingMethod  string         `gorm:"size:50" json:"charging_method"` // Prepaid/Postpaid
	ListenerCount   int            `gorm:"default:0" json:"listener_count"`
	BackendCount    int            `gorm:"default:0" json:"backend_count"`
	ProviderType    string         `gorm:"size:20" json:"provider_type"`
	AccountName     string         `gorm:"size:100" json:"account_name"`
	ProjectID       uint           `gorm:"index" json:"project_id"`
	ProjectName     string         `gorm:"size:100" json:"project_name"`
	RegionID        string         `gorm:"size:100;index" json:"region_id"`
	Region          string         `gorm:"size:100" json:"region"`
	Status          string         `gorm:"size:50" json:"status"`
	Tags            datatypes.JSON `gorm:"type:json" json:"tags"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

func (CloudLBInstance) TableName() string { return "sync_cloud_lb_instances" }

// CloudLBACL 负载均衡访问控制（同步后的本地存储）
type CloudLBACL struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	CloudAccountID uint           `gorm:"index;not null" json:"cloud_account_id"`
	AclID          string         `gorm:"size:100;uniqueIndex:acl_account_idx" json:"acl_id"`
	Name           string         `gorm:"size:200" json:"name"`
	AddressSource  string         `gorm:"size:500" json:"address_source"` // IP地址列表
	Remarks        string         `gorm:"size:500" json:"remarks"`
	ListenerCount  int            `gorm:"default:0" json:"listener_count"`
	Status         string         `gorm:"size:50;default:'ready'" json:"status"`
	ProviderType   string         `gorm:"size:20" json:"provider_type"`
	AccountName    string         `gorm:"size:100" json:"account_name"`
	SharedScope    string         `gorm:"size:50" json:"shared_scope"`
	RegionID       string         `gorm:"size:100;index" json:"region_id"`
	Region         string         `gorm:"size:100" json:"region"`
	ProjectID      uint           `gorm:"index" json:"project_id"`
	ProjectName    string         `gorm:"size:100" json:"project_name"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

func (CloudLBACL) TableName() string { return "sync_cloud_lb_acls" }

// CloudLBCertificate 负载均衡证书（同步后的本地存储）
type CloudLBCertificate struct {
	ID                    uint           `gorm:"primaryKey" json:"id"`
	CloudAccountID        uint           `gorm:"index;not null" json:"cloud_account_id"`
	CertID                string         `gorm:"size:100;uniqueIndex:cert_account_idx" json:"cert_id"`
	Name                  string         `gorm:"size:200" json:"name"`
	DomainName            string         `gorm:"size:200" json:"domain_name"`
	Expiration            string         `gorm:"size:50" json:"expiration"` // 过期时间
	SubjectAlternativeNames string       `gorm:"size:500" json:"subject_alternative_names"` // SAN域名列表
	ListenerCount         int            `gorm:"default:0" json:"listener_count"`
	Status                string         `gorm:"size:50;default:'ready'" json:"status"`
	SharedScope           string         `gorm:"size:50" json:"shared_scope"`
	ProjectID             uint           `gorm:"index" json:"project_id"`
	ProjectName           string         `gorm:"size:100" json:"project_name"`
	ProviderType          string         `gorm:"size:20" json:"provider_type"`
	AccountName           string         `gorm:"size:100" json:"account_name"`
	RegionID              string         `gorm:"size:100;index" json:"region_id"`
	Region                string         `gorm:"size:100" json:"region"`
	Tags                  datatypes.JSON `gorm:"type:json" json:"tags"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
}

func (CloudLBCertificate) TableName() string { return "sync_cloud_lb_certificates" }

// ========== CDN Model ==========

// CloudCDNDomain CDN域名（同步后的本地存储）
type CloudCDNDomain struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	CloudAccountID uint           `gorm:"index;not null" json:"cloud_account_id"`
	CdnID          string         `gorm:"size:100;uniqueIndex:cdn_account_idx" json:"cdn_id"`
	Name           string         `gorm:"size:200" json:"name"` // 域名
	CNAME          string         `gorm:"size:200" json:"cname"` // CNAME地址
	Area           string         `gorm:"size:50" json:"area"` // mainland/overseas/global
	ServiceType    string         `gorm:"size:50" json:"service_type"` // web/download/video
	OriginType     string         `gorm:"size:50" json:"origin_type"` // ip/domain/oss
	OriginAddress  string         `gorm:"size:500" json:"origin_address"` // 源站地址
	Status         string         `gorm:"size:50" json:"status"` // online/offline/configuring
	ProviderType   string         `gorm:"size:20" json:"provider_type"`
	AccountName    string         `gorm:"size:100" json:"account_name"`
	ProjectID      uint           `gorm:"index" json:"project_id"`
	ProjectName    string         `gorm:"size:100" json:"project_name"`
	Tags           datatypes.JSON `gorm:"type:json" json:"tags"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

func (CloudCDNDomain) TableName() string { return "sync_cloud_cdn_domains" }