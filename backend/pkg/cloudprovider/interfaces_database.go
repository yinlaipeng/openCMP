package cloudprovider

import (
	"context"
	"time"
)

// IRDS RDS 关系型数据库管理接口
type IRDS interface {
	CreateRDSInstance(ctx context.Context, config RDSConfig) (*RDSInstance, error)
	DeleteRDSInstance(ctx context.Context, instanceID string) error
	StartRDSInstance(ctx context.Context, instanceID string) error
	StopRDSInstance(ctx context.Context, instanceID string) error
	RebootRDSInstance(ctx context.Context, instanceID string) error
	ResizeRDSInstance(ctx context.Context, instanceID string, spec RDSpec) error
	CreateRDSBackup(ctx context.Context, instanceID string, name string) (*RDSBackup, error)
	RestoreRDSFromBackup(ctx context.Context, backupID string, config RDSConfig) (*RDSInstance, error)
	ListRDSInstances(ctx context.Context, filter RDSFilter) ([]*RDSInstance, error)
	ListRDSBackups(ctx context.Context, instanceID string) ([]*RDSBackup, error)
}

// IElasticCache 弹性缓存管理接口
type IElasticCache interface {
	CreateCacheInstance(ctx context.Context, config CacheConfig) (*CacheInstance, error)
	DeleteCacheInstance(ctx context.Context, instanceID string) error
	RebootCacheInstance(ctx context.Context, instanceID string) error
	ResizeCacheInstance(ctx context.Context, instanceID string, spec CacheSpec) error
	CreateCacheBackup(ctx context.Context, instanceID string) (*CacheBackup, error)
	ListCacheInstances(ctx context.Context, filter CacheFilter) ([]*CacheInstance, error)
}

// IDatabase 数据服务总接口
type IDatabase interface {
	IRDS
	IElasticCache
}

// RDSInstance RDS 实例
type RDSInstance struct {
	ID             string            `json:"id"`
	Name           string            `json:"name"`
	Engine         string            `json:"engine"` // MySQL/PostgreSQL/SQLServer/Postgres
	EngineVersion  string            `json:"engine_version"`
	InstanceType   string            `json:"instance_type"`
	StorageSize    int               `json:"storage_size"` // GB
	StorageType    string            `json:"storage_type"` // SSD/ESSD
	Status         string            `json:"status"`
	VPCID          string            `json:"vpc_id"`
	SubnetID       string            `json:"subnet_id"`
	Endpoint       string            `json:"endpoint"`
	Port           int               `json:"port"`
	MasterUsername string            `json:"master_username"`
	Tags           map[string]string `json:"tags"`
	CreatedAt      time.Time         `json:"created_at"`
	ZoneID         string            `json:"zone_id"`
}

// RDSConfig RDS 创建配置
type RDSConfig struct {
	Name           string            `json:"name"`
	Engine         string            `json:"engine"`
	EngineVersion  string            `json:"engine_version"`
	InstanceType   string            `json:"instance_type"`
	StorageSize    int               `json:"storage_size"`
	StorageType    string            `json:"storage_type"`
	VPCID          string            `json:"vpc_id"`
	SubnetID       string            `json:"subnet_id"`
	ZoneID         string            `json:"zone_id"`
	MasterUsername string            `json:"master_username"`
	MasterPassword string            `json:"master_password"`
	Tags           map[string]string `json:"tags"`
}

// RDSpec RDS 规格
type RDSpec struct {
	InstanceType string `json:"instance_type"`
	StorageSize  int    `json:"storage_size"`
}

// RDSBackup RDS 备份
type RDSBackup struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	InstanceID string    `json:"instance_id"`
	Status     string    `json:"status"` // creating/available/failed
	Size       int64     `json:"size"`   // bytes
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
}

// RDSFilter RDS 列表过滤条件
type RDSFilter struct {
	InstanceID string
	Engine     string
	Status     string
	VPCID      string
	MaxResults int
}

// CacheInstance 弹性缓存实例
type CacheInstance struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	Engine        string            `json:"engine"` // Redis/Memcached
	EngineVersion string            `json:"engine_version"`
	InstanceType  string            `json:"instance_type"`
	Status        string            `json:"status"`
	VPCID         string            `json:"vpc_id"`
	SubnetID      string            `json:"subnet_id"`
	Endpoint      string            `json:"endpoint"`
	Port          int               `json:"port"`
	Tags          map[string]string `json:"tags"`
	CreatedAt     time.Time         `json:"created_at"`
	ZoneID        string            `json:"zone_id"`
}

// CacheConfig 缓存实例创建配置
type CacheConfig struct {
	Name          string            `json:"name"`
	Engine        string            `json:"engine"`
	EngineVersion string            `json:"engine_version"`
	InstanceType  string            `json:"instance_type"`
	VPCID         string            `json:"vpc_id"`
	SubnetID      string            `json:"subnet_id"`
	ZoneID        string            `json:"zone_id"`
	Tags          map[string]string `json:"tags"`
}

// CacheSpec 缓存规格
type CacheSpec struct {
	InstanceType string `json:"instance_type"`
}

// CacheBackup 缓存备份
type CacheBackup struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	InstanceID string    `json:"instance_id"`
	Status     string    `json:"status"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
}

// CacheFilter 缓存实例列表过滤条件
type CacheFilter struct {
	InstanceID string
	Engine     string
	Status     string
	MaxResults int
}
