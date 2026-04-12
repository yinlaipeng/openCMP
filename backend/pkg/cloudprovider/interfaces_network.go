package cloudprovider

import (
	"context"
	"time"
)

// IVPC VPC 管理接口
type IVPC interface {
	CreateVPC(ctx context.Context, config VPCConfig) (*VPC, error)
	DeleteVPC(ctx context.Context, vpcID string) error
	GetVPC(ctx context.Context, vpcID string) (*VPC, error)
	ListVPCs(ctx context.Context, filter VPCFilter) ([]*VPC, error)
}

// VPCFilter VPC 列表过滤条件
type VPCFilter struct {
	VPCID      string
	Name       string
	Status     string
	Tags       map[string]string
	RegionID   string
	MaxResults int
	NextToken  string
}

// ISubnet 子网管理接口
type ISubnet interface {
	CreateSubnet(ctx context.Context, config SubnetConfig) (*Subnet, error)
	DeleteSubnet(ctx context.Context, subnetID string) error
	GetSubnet(ctx context.Context, subnetID string) (*Subnet, error)
	ListSubnets(ctx context.Context, filter SubnetFilter) ([]*Subnet, error)
}

// SubnetFilter 子网列表过滤条件
type SubnetFilter struct {
	VPCID      string
	SubnetID   string
	ZoneID     string
	MaxResults int
}

// ISecurityGroup 安全组管理接口
type ISecurityGroup interface {
	CreateSecurityGroup(ctx context.Context, config SGConfig) (*SecurityGroup, error)
	DeleteSecurityGroup(ctx context.Context, sgID string) error
	AuthorizeSecurityGroup(ctx context.Context, sgID string, rules []SGRule) error
	RevokeSecurityGroup(ctx context.Context, sgID string, rules []SGRule) error
	ListSecurityGroups(ctx context.Context, filter SGFilter) ([]*SecurityGroup, error)
}

// SGFilter 安全组列表过滤条件
type SGFilter struct {
	VPCID      string
	SGID       string
	Name       string
	MaxResults int
}

// IEIP 弹性 IP 管理接口
type IEIP interface {
	AllocateEIP(ctx context.Context, config EIPConfig) (*EIP, error)
	ReleaseEIP(ctx context.Context, eipID string) error
	AssociateEIP(ctx context.Context, eipID, resourceID string) error
	DissociateEIP(ctx context.Context, eipID string) error
	ListEIPs(ctx context.Context, filter EIPFilter) ([]*EIP, error)
}

// EIPFilter 弹性 IP 列表过滤条件
type EIPFilter struct {
	EIPID      string
	Address    string
	Status     string
	RegionID   string
	MaxResults int
}

// ILoadBalancer 负载均衡管理接口
type ILoadBalancer interface {
	CreateLoadBalancer(ctx context.Context, config LBConfig) (*LoadBalancer, error)
	DeleteLoadBalancer(ctx context.Context, lbID string) error
	CreateListener(ctx context.Context, lbID string, config ListenerConfig) (*Listener, error)
	DeleteListener(ctx context.Context, listenerID string) error
	ListLoadBalancers(ctx context.Context, filter LBFilter) ([]*LoadBalancer, error)
}

// LBFilter 负载均衡列表过滤条件
type LBFilter struct {
	LBID       string
	Name       string
	VPCID      string
	MaxResults int
}

// IDNS DNS 管理接口
type IDNS interface {
	CreateDNSZone(ctx context.Context, config DNSZoneConfig) (*DNSZone, error)
	DeleteDNSZone(ctx context.Context, zoneID string) error
	CreateDNSRecord(ctx context.Context, zoneID string, config DNSRecordConfig) (*DNSRecord, error)
	DeleteDNSRecord(ctx context.Context, recordID string) error
	ListDNSZones(ctx context.Context) ([]*DNSZone, error)
	ListDNSRecords(ctx context.Context, zoneID string) ([]*DNSRecord, error)
}

// VPCInterconnect represents a connection between VPCs or between a VPC and an on-premises network
type VPCInterconnect struct {
	ID                string            `json:"id"`
	Name              string            `json:"name"`
	Type              string            `json:"type"` //专线, VPN, Direct Connect, etc.
	Bandwidth         int               `json:"bandwidth"` // Mbps
	Status            string            `json:"status"`
	RegionID          string            `json:"region_id"`
	PeerRegion        string            `json:"peer_region"`
	VPCID             string            `json:"vpc_id"`
	Domain            string            `json:"domain,omitempty"`
	Platform          string            `json:"platform,omitempty"`
	Account           string            `json:"account,omitempty"`
	Description       string            `json:"description,omitempty"`
	Tags              map[string]string `json:"tags"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
}

// VPCPeering represents a connection between two VPCs
type VPCPeering struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Status      string            `json:"status"`
	LocalVPCID  string            `json:"local_vpc_id"`
	PeerVPCID   string            `json:"peer_vpc_id"`
	PeerAccount string            `json:"peer_account"`
	Domain      string            `json:"domain,omitempty"`
	Platform    string            `json:"platform,omitempty"`
	Account     string            `json:"account,omitempty"`
	Description string            `json:"description,omitempty"`
	Tags        map[string]string `json:"tags"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// Route represents a route in a route table
type Route struct {
	DestinationCIDR string `json:"destination_cidr"`
	Target          string `json:"target"` // VPC Gateway, Internet Gateway, VPC Peering, etc.
	Status          string `json:"status"`
}

// RouteTable represents a route table
type RouteTable struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	VPCID       string            `json:"vpc_id"`
	Status      string            `json:"status"`
	Routes      []Route           `json:"routes"`
	Domain      string            `json:"domain,omitempty"`
	Platform    string            `json:"platform,omitempty"`
	Account     string            `json:"account,omitempty"`
	RegionID    string            `json:"region_id,omitempty"`
	Description string            `json:"description,omitempty"`
	Tags        map[string]string `json:"tags"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// L2Network represents a Layer 2 network
type L2Network struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	VLANID      int               `json:"vlan_id"`
	Status      string            `json:"status"`
	VPCID       string            `json:"vpc_id"`
	Bandwidth   int               `json:"bandwidth"` // Mbps
	NetworkCount int              `json:"network_count,omitempty"`
	Platform    string            `json:"platform,omitempty"`
	Domain      string            `json:"domain,omitempty"`
	RegionID    string            `json:"region_id,omitempty"`
	Description string            `json:"description,omitempty"`
	Tags        map[string]string `json:"tags"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// VPCInterconnectConfig represents the configuration for a VPC interconnect
type VPCInterconnectConfig struct {
	Name       string            `json:"name"`
	Type       string            `json:"type"`       //专线, VPN, Direct Connect, etc.
	Bandwidth  int               `json:"bandwidth"`  // Mbps
	RegionID   string            `json:"region_id"`
	PeerRegion string            `json:"peer_region"`
	Tags       map[string]string `json:"tags"`
}

// VPCInterconnectFilter represents filters for listing VPC interconnects
type VPCInterconnectFilter struct {
	Type       string `json:"type,omitempty"`
	Status     string `json:"status,omitempty"`
	RegionID   string `json:"region_id,omitempty"`
	PeerRegion string `json:"peer_region,omitempty"`
}

// VPCPeeringConfig represents the configuration for a VPC peering connection
type VPCPeeringConfig struct {
	Name         string            `json:"name"`
	LocalVPCID   string            `json:"local_vpc_id"`
	PeerVPCID    string            `json:"peer_vpc_id"`
	PeerAccount  string            `json:"peer_account"`
	Tags         map[string]string `json:"tags"`
}

// VPCPeeringFilter represents filters for listing VPC peerings
type VPCPeeringFilter struct {
	Status      string `json:"status,omitempty"`
	LocalVPCID  string `json:"local_vpc_id,omitempty"`
	PeerVPCID   string `json:"peer_vpc_id,omitempty"`
	PeerAccount string `json:"peer_account,omitempty"`
}

// RouteTableConfig represents the configuration for a route table
type RouteTableConfig struct {
	Name   string            `json:"name"`
	VPCID  string            `json:"vpc_id"`
	Routes []Route           `json:"routes"`
	Tags   map[string]string `json:"tags"`
}

// RouteTableFilter represents filters for listing route tables
type RouteTableFilter struct {
	Status string `json:"status,omitempty"`
	VPCID  string `json:"vpc_id,omitempty"`
}

// L2NetworkConfig represents the configuration for an L2 network
type L2NetworkConfig struct {
	Name    string            `json:"name"`
	VLANID  int               `json:"vlan_id"`
	VPCID   string            `json:"vpc_id"`
	Subnets []string          `json:"subnets"`
	Tags    map[string]string `json:"tags"`
}

// L2NetworkFilter represents filters for listing L2 networks
type L2NetworkFilter struct {
	Status   string `json:"status,omitempty"`
	VPCID    string `json:"vpc_id,omitempty"`
	VLANID   int    `json:"vlan_id,omitempty"`
}

// INetwork 网络资源总接口
type INetwork interface {
	IVPC
	ISubnet
	ISecurityGroup
	IEIP
	ILoadBalancer
	IDNS
	// Geography methods
	ListRegions() ([]*Region, error)
	ListZones(regionID string) ([]*Zone, error)
	// Advanced network methods
	CreateVPCInterconnect(ctx context.Context, config VPCInterconnectConfig) (*VPCInterconnect, error)
	DeleteVPCInterconnect(ctx context.Context, interconnectID string) error
	ListVPCInterconnects(ctx context.Context, filter VPCInterconnectFilter) ([]*VPCInterconnect, error)
	CreateVPCPeering(ctx context.Context, config VPCPeeringConfig) (*VPCPeering, error)
	DeleteVPCPeering(ctx context.Context, peeringID string) error
	ListVPCPeerings(ctx context.Context, filter VPCPeeringFilter) ([]*VPCPeering, error)
	CreateRouteTable(ctx context.Context, config RouteTableConfig) (*RouteTable, error)
	DeleteRouteTable(ctx context.Context, routeTableID string) error
	ListRouteTables(ctx context.Context, filter RouteTableFilter) ([]*RouteTable, error)
	CreateL2Network(ctx context.Context, config L2NetworkConfig) (*L2Network, error)
	DeleteL2Network(ctx context.Context, l2NetworkID string) error
	ListL2Networks(ctx context.Context, filter L2NetworkFilter) ([]*L2Network, error)
}
