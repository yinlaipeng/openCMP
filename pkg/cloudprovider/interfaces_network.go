package cloudprovider

import "context"

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

// INetwork 网络资源总接口
type INetwork interface {
	IVPC
	ISubnet
	ISecurityGroup
	IEIP
	ILoadBalancer
	IDNS
}
