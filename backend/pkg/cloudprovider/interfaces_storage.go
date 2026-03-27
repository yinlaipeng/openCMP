package cloudprovider

import (
	"context"
	"io"
	"time"
)

// IDisk 云磁盘管理接口
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

// DiskFilter 磁盘列表过滤条件
type DiskFilter struct {
	DiskID     string
	VMID       string
	ZoneID     string
	Status     string
	MaxResults int
}

// IBucket 对象存储管理接口
type IBucket interface {
	CreateBucket(ctx context.Context, name string, config BucketConfig) error
	DeleteBucket(ctx context.Context, name string) error
	ListBuckets(ctx context.Context) ([]*Bucket, error)
	PutObject(ctx context.Context, bucketName, objectKey string, data io.Reader) error
	GetObject(ctx context.Context, bucketName, objectKey string) (io.ReadCloser, error)
	DeleteObject(ctx context.Context, bucketName, objectKey string) error
	ListObjects(ctx context.Context, bucketName string, prefix string) ([]*Object, error)
}

// INAS NAS 文件存储管理接口
type INAS interface {
	CreateFileSystem(ctx context.Context, config FSConfig) (*FileSystem, error)
	DeleteFileSystem(ctx context.Context, fsID string) error
	MountFileSystem(ctx context.Context, fsID string, config MountConfig) (*MountTarget, error)
	UnmountFileSystem(ctx context.Context, mountID string) error
	ListFileSystems(ctx context.Context) ([]*FileSystem, error)
}

// IStorage 存储资源总接口
type IStorage interface {
	IDisk
	IBucket
	INAS
}

// Bucket 存储桶
type Bucket struct {
	Name      string            `json:"name"`
	Region    string            `json:"region"`
	ACL       string            `json:"acl"` // private/public-read/public-read-write
	Tags      map[string]string `json:"tags"`
	CreatedAt time.Time         `json:"created_at"`
}

// BucketConfig 存储桶创建配置
type BucketConfig struct {
	Region string            `json:"region"`
	ACL    string            `json:"acl"`
	Tags   map[string]string `json:"tags"`
}

// Object 对象存储中的对象
type Object struct {
	Key          string    `json:"key"`
	Size         int64     `json:"size"`
	ETag         string    `json:"etag"`
	LastModified time.Time `json:"last_modified"`
	StorageClass string    `json:"storage_class"`
}

// FileSystem NAS 文件系统
type FileSystem struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Type        string            `json:"type"` // NFS/SMB
	Size        int               `json:"size"` // GB
	UsedSize    int               `json:"used_size"`
	Status      string            `json:"status"`
	VPCID       string            `json:"vpc_id"`
	MountTarget string            `json:"mount_target"`
	Tags        map[string]string `json:"tags"`
	CreatedAt   time.Time         `json:"created_at"`
}

// FSConfig 文件系统创建配置
type FSConfig struct {
	Name  string            `json:"name"`
	Type  string            `json:"type"`
	Size  int               `json:"size"`
	VPCID string            `json:"vpc_id"`
	Tags  map[string]string `json:"tags"`
}

// MountTarget 挂载点
type MountTarget struct {
	ID       string `json:"id"`
	FSID     string `json:"fs_id"`
	VPCID    string `json:"vpc_id"`
	SubnetID string `json:"subnet_id"`
	Address  string `json:"address"`
	Status   string `json:"status"`
}

// MountConfig 挂载点创建配置
type MountConfig struct {
	VPCID    string `json:"vpc_id"`
	SubnetID string `json:"subnet_id"`
}
