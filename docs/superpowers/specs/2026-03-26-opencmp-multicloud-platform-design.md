# 多云管理平台设计文档

**日期**: 2026-03-26  
**状态**: 已批准  
**参考项目**: [Cloudpods](https://github.com/yunionio/cloudpods)

---

## 1. 概述

### 1.1 项目目标

构建一个基于 Go + Gin + Gorm 的单体架构多云管理平台，通过统一接口 + 适配器模式，实现对多个云厂商资源的统一管理和快速接入。

### 1.2 技术栈

- **语言**: Go
- **Web 框架**: Gin
- **ORM**: Gorm
- **架构**: 单体模块化（Monolithic Modular）
- **部署**: 支持 Docker/Kubernetes

### 1.3 首批支持的云厂商

- 阿里云 (Alibaba Cloud)
- 腾讯云 (Tencent Cloud)
- AWS (Amazon Web Services)
- Azure (Microsoft Azure)

---

## 2. 整体架构

```
┌──────────────────────────────────────────────────────────────────────┐
│                         API Layer (Gin)                               │
│                    RESTful API / HTTP Handlers                        │
├──────────────────────────────────────────────────────────────────────┤
│                        Service Layer                                  │
│           业务逻辑层 (资源管理/账户管理/任务编排/权限控制)              │
├──────────────────────────────────────────────────────────────────────┤
│                     Cloud Provider Layer                              │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐             │
│  │ Alibaba  │  │ Tencent  │  │   AWS    │  │  Azure   │             │
│  │ Adapter  │  │ Adapter  │  │ Adapter  │  │ Adapter  │             │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘             │
├──────────────────────────────────────────────────────────────────────┤
│                  Cloud Interface Layer (分层标准接口)                   │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐  ┌────────────┐     │
│  │   计算资源  │  │   网络资源  │  │   存储资源  │  │   数据服务  │     │
│  │ ICompute   │  │ INetwork   │  │ IStorage   │  │ IDatabase  │     │
│  └────────────┘  └────────────┘  └────────────┘  └────────────┘     │
├──────────────────────────────────────────────────────────────────────┤
│                      Data Layer (Gorm)                                │
│              MySQL/PostgreSQL / 云账户配置 / 资源元数据存储            │
└──────────────────────────────────────────────────────────────────────┘
```

---

## 3. 分层资源接口设计

### 3.1 计算资源接口 (`pkg/cloudprovider/interfaces_compute.go`)

```go
// 虚拟机管理
type IVirtualMachine interface {
    CreateVM(ctx context.Context, config VMCreateConfig) (*VirtualMachine, error)
    DeleteVM(ctx context.Context, vmID string) error
    StartVM(ctx context.Context, vmID string) error
    StopVM(ctx context.Context, vmID string) error
    RebootVM(ctx context.Context, vmID string) error
    GetVMStatus(ctx context.Context, vmID string) (*VMStatus, error)
    ResizeVM(ctx context.Context, vmID string, spec VMSpec) error
    ListVMs(ctx context.Context, filter VMListFilter) ([]*VirtualMachine, error)
}

// 镜像管理
type IImage interface {
    ListImages(ctx context.Context, filter ImageFilter) ([]*Image, error)
    GetImage(ctx context.Context, imageID string) (*Image, error)
    CreateImage(ctx context.Context, config ImageCreateConfig) (*Image, error)
    DeleteImage(ctx context.Context, imageID string) error
}

// 密钥对管理
type IKeypair interface {
    CreateKeypair(ctx context.Context, name, publicKey string) (*Keypair, error)
    DeleteKeypair(ctx context.Context, keypairID string) error
    ListKeypairs(ctx context.Context) ([]*Keypair, error)
}

// 计算资源总接口
type ICompute interface {
    IVirtualMachine
    IImage
    IKeypair
}
```

### 3.2 网络资源接口 (`pkg/cloudprovider/interfaces_network.go`)

```go
// VPC 管理
type IVPC interface {
    CreateVPC(ctx context.Context, config VPCConfig) (*VPC, error)
    DeleteVPC(ctx context.Context, vpcID string) error
    GetVPC(ctx context.Context, vpcID string) (*VPC, error)
    ListVPCs(ctx context.Context, filter VPCFilter) ([]*VPC, error)
}

// 子网管理
type ISubnet interface {
    CreateSubnet(ctx context.Context, config SubnetConfig) (*Subnet, error)
    DeleteSubnet(ctx context.Context, subnetID string) error
    GetSubnet(ctx context.Context, subnetID string) (*Subnet, error)
    ListSubnets(ctx context.Context, filter SubnetFilter) ([]*Subnet, error)
}

// 安全组管理
type ISecurityGroup interface {
    CreateSecurityGroup(ctx context.Context, config SGConfig) (*SecurityGroup, error)
    DeleteSecurityGroup(ctx context.Context, sgID string) error
    AuthorizeSecurityGroup(ctx context.Context, sgID string, rules []SGRule) error
    RevokeSecurityGroup(ctx context.Context, sgID string, rules []SGRule) error
    ListSecurityGroups(ctx context.Context, filter SGFilter) ([]*SecurityGroup, error)
}

// 弹性 IP 管理
type IEIP interface {
    AllocateEIP(ctx context.Context, config EIPConfig) (*EIP, error)
    ReleaseEIP(ctx context.Context, eipID string) error
    AssociateEIP(ctx context.Context, eipID, resourceID string) error
    DissociateEIP(ctx context.Context, eipID string) error
    ListEIPs(ctx context.Context, filter EIPFilter) ([]*EIP, error)
}

// 负载均衡管理
type ILoadBalancer interface {
    CreateLoadBalancer(ctx context.Context, config LBConfig) (*LoadBalancer, error)
    DeleteLoadBalancer(ctx context.Context, lbID string) error
    CreateListener(ctx context.Context, lbID string, config ListenerConfig) (*Listener, error)
    DeleteListener(ctx context.Context, listenerID string) error
    ListLoadBalancers(ctx context.Context, filter LBFilter) ([]*LoadBalancer, error)
}

// DNS 管理
type IDNS interface {
    CreateDNSZone(ctx context.Context, config DNSZoneConfig) (*DNSZone, error)
    DeleteDNSZone(ctx context.Context, zoneID string) error
    CreateDNSRecord(ctx context.Context, zoneID string, config DNSRecordConfig) (*DNSRecord, error)
    DeleteDNSRecord(ctx context.Context, recordID string) error
    ListDNSZones(ctx context.Context) ([]*DNSZone, error)
    ListDNSRecords(ctx context.Context, zoneID string) ([]*DNSRecord, error)
}

// 网络资源总接口
type INetwork interface {
    IVPC
    ISubnet
    ISecurityGroup
    IEIP
    ILoadBalancer
    IDNS
}
```

### 3.3 存储资源接口 (`pkg/cloudprovider/interfaces_storage.go`)

```go
// 云磁盘管理
type IDisk interface {
    CreateDisk(ctx context.Context, config DiskConfig) (*Disk, error)
    DeleteDisk(ctx context.Context, diskID string) error
    AttachDisk(ctx context.Context, diskID, vmID string) error
    DetachDisk(ctx context.Context, diskID string) error
    ResizeDisk(ctx context.Context, diskID string, sizeGB int) error
    CreateSnapshot(ctx context.Context, diskID string, name string) (*Snapshot, error)
    DeleteSnapshot(ctx context.Context, snapshotID string) error
    ListDisks(ctx context.Context, filter DiskFilter) ([]*Disk, error)
    ListSnapshots(ctx context.Context, filter SnapshotFilter) ([]*Snapshot, error)
}

// 对象存储管理
type IBucket interface {
    CreateBucket(ctx context.Context, name string, config BucketConfig) error
    DeleteBucket(ctx context.Context, name string) error
    ListBuckets(ctx context.Context) ([]*Bucket, error)
    PutObject(ctx context.Context, bucketName, objectKey string, data io.Reader) error
    GetObject(ctx context.Context, bucketName, objectKey string) (io.ReadCloser, error)
    DeleteObject(ctx context.Context, bucketName, objectKey string) error
    ListObjects(ctx context.Context, bucketName string, prefix string) ([]*Object, error)
}

// NAS 文件存储管理
type INAS interface {
    CreateFileSystem(ctx context.Context, config FSConfig) (*FileSystem, error)
    DeleteFileSystem(ctx context.Context, fsID string) error
    MountFileSystem(ctx context.Context, fsID string, config MountConfig) (*MountTarget, error)
    UnmountFileSystem(ctx context.Context, mountID string) error
    ListFileSystems(ctx context.Context) ([]*FileSystem, error)
}

// 存储资源总接口
type IStorage interface {
    IDisk
    IBucket
    INAS
}
```

### 3.4 数据服务接口 (`pkg/cloudprovider/interfaces_database.go`)

```go
// RDS 关系型数据库
type IRDS interface {
    CreateRDSInstance(ctx context.Context, config RDSConfig) (*RDSInstance, error)
    DeleteRDSInstance(ctx context.Context, instanceID string) error
    StartRDSInstance(ctx context.Context, instanceID string) error
    StopRDSInstance(ctx context.Context, instanceID string) error
    RebootRDSInstance(ctx context.Context, instanceID string) error
    ResizeRDSInstance(ctx context.Context, instanceID string, spec RDSpec) error
    CreateRDSBackup(ctx context.Context, instanceID string, name string) (*RDSBackup, error)
    RestoreRDSFromBackup(ctx context.Context, backupID string, config RDSConfig) (*RDSInstance, error)
    CreateRDSAccount(ctx context.Context, instanceID string, config RDSAccountConfig) error
    ListRDSInstances(ctx context.Context, filter RDSFilter) ([]*RDSInstance, error)
    ListRDSBackups(ctx context.Context, instanceID string) ([]*RDSBackup, error)
}

// 弹性缓存 (Redis/Memcached)
type IElasticCache interface {
    CreateCacheInstance(ctx context.Context, config CacheConfig) (*CacheInstance, error)
    DeleteCacheInstance(ctx context.Context, instanceID string) error
    RebootCacheInstance(ctx context.Context, instanceID string) error
    ResizeCacheInstance(ctx context.Context, instanceID string, spec CacheSpec) error
    CreateCacheBackup(ctx context.Context, instanceID string) (*CacheBackup, error)
    ListCacheInstances(ctx context.Context, filter CacheFilter) ([]*CacheInstance, error)
}

// 数据服务总接口
type IDatabase interface {
    IRDS
    IElasticCache
}
```

### 3.5 云提供商总接口

```go
// ICloudProvider 组合所有资源接口
type ICloudProvider interface {
    ICompute    // 计算资源
    INetwork    // 网络资源
    IStorage    // 存储资源
    IDatabase   // 数据服务
    
    // 云厂商信息
    GetCloudInfo() CloudInfo
    // 区域和可用区列表
    ListRegions() ([]*Region, error)
    ListZones(regionID string) ([]*Zone, error)
    // 实例规格列表
    ListInstanceTypes(regionID string) ([]*InstanceType, error)
}
```

---

## 4. 通用类型定义

### 4.1 云账户配置

```go
type CloudAccountConfig struct {
    ID           string            `json:"id"`
    Name         string            `json:"name"`
    ProviderType string            `json:"provider_type"` // alibaba/tencent/aws/azure
    Credentials  map[string]string `json:"credentials"`   // 加密存储
    Region       string            `json:"region"`        // 默认区域
}
```

### 4.2 虚拟机配置与结构

```go
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
```

### 4.3 网络资源类型

```go
type VPCConfig struct {
    Name        string            `json:"name"`
    CIDR        string            `json:"cidr"`
    Description string            `json:"description"`
    Tags        map[string]string `json:"tags"`
}

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

type SubnetConfig struct {
    Name        string            `json:"name"`
    VPCID       string            `json:"vpc_id"`
    CIDR        string            `json:"cidr"`
    ZoneID      string            `json:"zone_id"`
    Description string            `json:"description"`
    Tags        map[string]string `json:"tags"`
}

type SGRule struct {
    Direction   string `json:"direction"` // ingress/egress
    Protocol    string `json:"protocol"`  // tcp/udp/icmp/all
    PortRange   string `json:"port_range"`
    CIDR        string `json:"cidr"`
    Action      string `json:"action"`    // accept/drop
    Description string `json:"description"`
    Priority    int    `json:"priority"`
}
```

### 4.4 存储与数据库类型

```go
type DiskConfig struct {
    Name   string            `json:"name"`
    Size   int               `json:"size"` // GB
    Type   string            `json:"type"` // SSD/HDD/ESSD
    ZoneID string            `json:"zone_id"`
    Tags   map[string]string `json:"tags"`
}

type RDSConfig struct {
    Name           string            `json:"name"`
    Engine         string            `json:"engine"` // MySQL/PostgreSQL/SQLServer
    EngineVersion  string            `json:"engine_version"`
    InstanceType   string            `json:"instance_type"`
    StorageSize    int               `json:"storage_size"` // GB
    VPCID          string            `json:"vpc_id"`
    SubnetID       string            `json:"subnet_id"`
    MasterUsername string            `json:"master_username"`
    MasterPassword string            `json:"master_password"`
    Tags           map[string]string `json:"tags"`
}
```

### 4.5 区域与实例规格

```go
type Region struct {
    ID     string `json:"id"`
    Name   string `json:"name"`
    Status string `json:"status"`
}

type Zone struct {
    ID       string `json:"id"`
    Name     string `json:"name"`
    RegionID string `json:"region_id"`
    Status   string `json:"status"`
}

type InstanceType struct {
    Name             string `json:"name"`
    CPU              int    `json:"cpu"`    // 核心数
    Memory           int    `json:"memory"` // MB
    GPU              int    `json:"gpu"`
    Category         string `json:"category"` // 通用型/计算型/内存型
    SupportedZones   []string `json:"supported_zones"`
}

type CloudInfo struct {
    Provider string   `json:"provider"`
    Version  string   `json:"version"`
    Regions  int      `json:"regions"`
    Services []string `json:"services"`
}
```

---

## 5. 适配器注册机制

```go
package cloudprovider

import (
    "fmt"
    "sync"
)

type ProviderFactory func(config CloudAccountConfig) (ICloudProvider, error)

var (
    providerRegistry = make(map[string]ProviderFactory)
    registryMutex    sync.RWMutex
)

// 注册云厂商适配器
func RegisterProvider(providerType string, factory ProviderFactory) {
    registryMutex.Lock()
    defer registryMutex.Unlock()
    providerRegistry[providerType] = factory
}

// 获取云提供商实例
func GetProvider(providerType string, config CloudAccountConfig) (ICloudProvider, error) {
    registryMutex.RLock()
    defer registryMutex.RUnlock()
    
    factory, ok := providerRegistry[providerType]
    if !ok {
        return nil, fmt.Errorf("provider %s not found", providerType)
    }
    return factory(config)
}

// 列出所有已注册的云厂商
func ListProviders() []string {
    registryMutex.RLock()
    defer registryMutex.RUnlock()
    
    providers := make([]string, 0, len(providerRegistry))
    for p := range providerRegistry {
        providers = append(providers, p)
    }
    return providers
}
```

---

## 6. 项目目录结构

```
openCMP/
├── cmd/
│   └── server/
│       └── main.go              # 应用入口
├── internal/
│   ├── handler/                 # HTTP Handler (Gin)
│   │   ├── cloud_account.go     # 云账户管理 API
│   │   ├── compute.go           # 计算资源 API
│   │   ├── network.go           # 网络资源 API
│   │   ├── storage.go           # 存储资源 API
│   │   └── database.go          # 数据库服务 API
│   ├── service/                 # 业务逻辑层
│   │   ├── cloud_account.go
│   │   ├── compute.go
│   │   ├── network.go
│   │   ├── storage.go
│   │   └── database.go
│   ├── model/                   # 数据模型 (Gorm)
│   │   ├── cloud_account.go
│   │   ├── compute.go
│   │   ├── network.go
│   │   ├── storage.go
│   │   └── database.go
│   └── middleware/              # Gin 中间件
│       ├── auth.go
│       ├── logger.go
│       └── recovery.go
├── pkg/
│   └── cloudprovider/           # 云适配器层
│       ├── interfaces_compute.go    # 计算资源接口
│       ├── interfaces_network.go    # 网络资源接口
│       ├── interfaces_storage.go    # 存储资源接口
│       ├── interfaces_database.go   # 数据库服务接口
│       ├── registry.go              # 适配器注册
│       ├── types.go                 # 通用类型定义
│       └── adapters/
│           ├── alibaba/         # 阿里云适配器
│           │   ├── provider.go
│           │   ├── vm.go
│           │   ├── vpc.go
│           │   ├── disk.go
│           │   └── rds.go
│           ├── tencent/         # 腾讯云适配器
│           │   ├── provider.go
│           │   ├── vm.go
│           │   ├── vpc.go
│           │   ├── disk.go
│           │   └── rds.go
│           ├── aws/             # AWS 适配器
│           │   ├── provider.go
│           │   ├── ec2.go
│           │   ├── vpc.go
│           │   ├── ebs.go
│           │   └── rds.go
│           └── azure/           # Azure 适配器
│               ├── provider.go
│               ├── vm.go
│               ├── vnet.go
│               ├── disk.go
│               └── sql.go
├── configs/
│   └── config.yaml
├── scripts/
│   └── init.sql
├── docs/
│   └── superpowers/specs/
│       └── 2026-03-26-opencmp-multicloud-platform-design.md
├── go.mod
├── go.sum
└── README.md
```

---

## 7. 快速接入新云厂商流程

以接入华为云为例：

### 步骤 1: 创建适配器目录

```bash
mkdir pkg/cloudprovider/adapters/huawei
```

### 步骤 2: 实现适配器工厂函数

```go
// pkg/cloudprovider/adapters/huawei/provider.go
package huawei

import (
    "github.com/yourorg/openCMP/pkg/cloudprovider"
)

type HuaweiProvider struct {
    config cloudprovider.CloudAccountConfig
    client *huaweicloud.Client
}

func NewHuaweiProvider(config cloudprovider.CloudAccountConfig) (cloudprovider.ICloudProvider, error) {
    client, err := huaweicloud.NewClient(config.Credentials)
    if err != nil {
        return nil, err
    }
    return &HuaweiProvider{config: config, client: client}, nil
}

func init() {
    cloudprovider.RegisterProvider("huawei", NewHuaweiProvider)
}
```

### 步骤 3: 实现各资源接口

- `vm.go` - 实现 `IVirtualMachine`, `IImage`, `IKeypair`
- `vpc.go` - 实现 `IVPC`, `ISubnet`, `ISecurityGroup`, `IEIP`
- `disk.go` - 实现 `IDisk`
- `rds.go` - 实现 `IRDS`

### 步骤 4: 完成

无需修改其他代码，新云厂商即可通过 API 使用。

---

## 8. API 设计

### 8.1 云账户管理

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | `/api/v1/cloud-accounts` | 添加云账户 |
| GET | `/api/v1/cloud-accounts` | 列出云账户 |
| GET | `/api/v1/cloud-accounts/:id` | 获取云账户详情 |
| PUT | `/api/v1/cloud-accounts/:id` | 更新云账户 |
| DELETE | `/api/v1/cloud-accounts/:id` | 删除云账户 |
| POST | `/api/v1/cloud-accounts/:id/verify` | 验证云账户 |

### 8.2 计算资源

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | `/api/v1/vms` | 创建虚拟机 |
| GET | `/api/v1/vms` | 列出虚拟机 |
| GET | `/api/v1/vms/:id` | 获取虚拟机详情 |
| DELETE | `/api/v1/vms/:id` | 删除虚拟机 |
| POST | `/api/v1/vms/:id/action` | 虚拟机操作 (start/stop/reboot) |
| PUT | `/api/v1/vms/:id/resize` | 调整虚拟机规格 |
| GET | `/api/v1/images` | 列出镜像 |
| GET | `/api/v1/keypairs` | 列出密钥对 |

### 8.3 网络资源

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | `/api/v1/vpcs` | 创建 VPC |
| GET | `/api/v1/vpcs` | 列出 VPC |
| DELETE | `/api/v1/vpcs/:id` | 删除 VPC |
| POST | `/api/v1/subnets` | 创建子网 |
| GET | `/api/v1/subnets` | 列出子网 |
| DELETE | `/api/v1/subnets/:id` | 删除子网 |
| POST | `/api/v1/security-groups` | 创建安全组 |
| POST | `/api/v1/security-groups/:id/rules` | 添加安全组规则 |
| GET | `/api/v1/eips` | 列出弹性 IP |
| POST | `/api/v1/load-balancers` | 创建负载均衡 |
| GET | `/api/v1/dns-zones` | 列出 DNS 区域 |

### 8.4 存储资源

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | `/api/v1/disks` | 创建磁盘 |
| GET | `/api/v1/disks` | 列出磁盘 |
| POST | `/api/v1/disks/:id/attach` | 挂载磁盘 |
| POST | `/api/v1/disks/:id/detach` | 卸载磁盘 |
| DELETE | `/api/v1/disks/:id` | 删除磁盘 |
| POST | `/api/v1/snapshots` | 创建快照 |
| GET | `/api/v1/buckets` | 列出存储桶 |

### 8.5 数据库服务

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | `/api/v1/rds-instances` | 创建 RDS 实例 |
| GET | `/api/v1/rds-instances` | 列出 RDS 实例 |
| POST | `/api/v1/rds-instances/:id/action` | RDS 操作 (start/stop/reboot) |
| POST | `/api/v1/rds-instances/:id/backups` | 创建 RDS 备份 |
| GET | `/api/v1/cache-instances` | 列出弹性缓存 |

---

## 9. 数据模型设计

### 9.1 云账户表

```go
type CloudAccount struct {
    ID           uint           `gorm:"primaryKey" json:"id"`
    Name         string         `gorm:"uniqueIndex;not null" json:"name"`
    ProviderType string         `gorm:"type:varchar(20);not null" json:"provider_type"`
    Credentials  datatypes.JSON `gorm:"type:json" json:"credentials"` // 加密存储
    Status       string         `gorm:"type:varchar(20)" json:"status"` // active/inactive/error
    CreatedAt    time.Time      `json:"created_at"`
    UpdatedAt    time.Time      `json:"updated_at"`
}
```

### 9.2 资源统一元数据表

```go
type CloudResource struct {
    ID           uint      `gorm:"primaryKey" json:"id"`
    AccountID    uint      `gorm:"index" json:"account_id"` // 所属云账户
    ProviderType string    `gorm:"type:varchar(20)" json:"provider_type"`
    ResourceType string    `gorm:"type:varchar(50)" json:"resource_type"` // vm/vpc/subnet/disk/rds...
    CloudID      string    `gorm:"index" json:"cloud_id"` // 云厂商资源 ID
    Name         string    `json:"name"`
    Region       string    `json:"region"`
    Zone         string    `json:"zone"`
    Status       string    `json:"status"`
    Specs        datatypes.JSON `gorm:"type:json" json:"specs"` // 规格详情
    Tags         datatypes.JSON `gorm:"type:json" json:"tags"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}
```

---

## 10. 设计原则

### 10.1 接口抽象原则

- **通用性优先**: 接口定义覆盖主流云厂商的通用能力
- **渐进式扩展**: 特殊功能通过扩展接口实现，不破坏原有接口
- **错误统一**: 定义统一的错误类型，屏蔽云厂商错误差异

### 10.2 适配器实现原则

- **幂等性**: Create/Delete 操作支持幂等调用
- **状态同步**: List/Get 操作实时同步云上状态
- **超时处理**: 所有 API 调用设置合理超时时间
- **重试机制**: 对可重试错误自动重试

### 10.3 安全原则

- **凭证加密**: 云账户凭证加密存储
- **最小权限**: 云账户使用最小权限原则配置 RAM/ IAM 策略
- **审计日志**: 所有资源操作记录审计日志

---

## 11. 后续扩展方向

### 11.1 更多云厂商

- 华为云 (Huawei Cloud)
- 百度智能云 (Baidu Cloud)
- 京东云 (JD Cloud)
- 天翼云 (CTYun)
- 移动云 (ECloud)
- 联通云 (Unicom Cloud)
- Google Cloud Platform (GCP)
- Oracle Cloud

### 11.2 私有云/虚拟化

- VMware vSphere
- OpenStack
- ZStack
- KVM (轻量级私有云)
- 裸金属服务器 (IPMI/Redfish)

### 11.3 更多资源类型

- 容器服务 (K8s 集群)
- 消息队列 (MQ)
- 大数据服务 (EMR)
- AI/ML 服务
- IoT 平台

### 11.4 高级功能

- 成本分析与优化建议
- 资源编排 (Terraform/ROS 集成)
- 自动化运维 (定时任务/告警)
- 多云容灾与备份
- 合规检查与等保支持

---

## 12. 自审检查

- [x] 无 TBD/TODO 占位符
- [x] 接口定义内部一致
- [x] 范围聚焦于核心资源管理
- [x] 所有需求无歧义
