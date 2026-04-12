package cloudprovider

import (
	"context"
	"time"
)

// IVirtualMachine 虚拟机管理接口
type IVirtualMachine interface {
	CreateVM(ctx context.Context, config VMCreateConfig) (*VirtualMachine, error)
	DeleteVM(ctx context.Context, vmID string) error
	StartVM(ctx context.Context, vmID string) error
	StopVM(ctx context.Context, vmID string) error
	RebootVM(ctx context.Context, vmID string) error
	GetVMStatus(ctx context.Context, vmID string) (*VMStatus, error)
	ListVMs(ctx context.Context, filter VMListFilter) ([]*VirtualMachine, error)
	GetVM(ctx context.Context, vmID string) (*VirtualMachine, error) // Added this method
	ResetVMPassword(ctx context.Context, vmID, username, newPassword string) error // Added for password reset
	UpdateVMConfig(ctx context.Context, vmID, instanceType, name string) error // Added for config update
}

// VMListFilter 虚拟机列表过滤条件
type VMListFilter struct {
	VPCID      string
	SubnetID   string
	Status     VMStatus
	Tags       map[string]string
	RegionID   string
	ZoneID     string
	MaxResults int
	NextToken  string
}

// IImage 镜像管理接口
type IImage interface {
	ListImages(ctx context.Context, filter ImageFilter) ([]*Image, error)
	GetImage(ctx context.Context, imageID string) (*Image, error)
}

// Image 镜像
type Image struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	OSName       string            `json:"os_name"`
	OSVersion    string            `json:"os_version"`
	Architecture string            `json:"architecture"` // x86_64, arm64
	Status       string            `json:"status"`
	Size         int64             `json:"size"` // bytes
	Tags         map[string]string `json:"tags"`
}

// ImageFilter 镜像过滤条件
type ImageFilter struct {
	Platform     string // Windows, Linux, Ubuntu, CentOS
	Architecture string
	Status       string
	MaxResults   int
}

// IKeypair 密钥对管理接口
type IKeypair interface {
	CreateKeypair(ctx context.Context, name, publicKey string) (*Keypair, error)
	DeleteKeypair(ctx context.Context, keypairID string) error
	ListKeypairs(ctx context.Context) ([]*Keypair, error)
}

// Keypair 密钥对
type Keypair struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	PublicKey   string    `json:"public_key"`
	Fingerprint string    `json:"fingerprint"`
	CreatedAt   time.Time `json:"created_at"`
}

// ICompute 计算资源总接口
type ICompute interface {
	IVirtualMachine
	IImage
	IKeypair
}
