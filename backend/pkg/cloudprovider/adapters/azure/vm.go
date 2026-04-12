package azure

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"

	"github.com/opencmp/opencmp/pkg/cloudprovider"
)

// CreateVM 创建虚拟机
func (p *AzureProvider) CreateVM(ctx context.Context, config cloudprovider.VMCreateConfig) (*cloudprovider.VirtualMachine, error) {
	resourceGroup := "opencmp-resources" // 默认资源组

	// 构建 VM 创建参数
	parameters := armcompute.VirtualMachine{
		Location: to.Ptr(p.location),
		Properties: &armcompute.VirtualMachineProperties{
			HardwareProfile: &armcompute.HardwareProfile{
				VMSize: to.Ptr(armcompute.VirtualMachineSizeTypes(config.InstanceType)),
			},
			StorageProfile: &armcompute.StorageProfile{
				ImageReference: &armcompute.ImageReference{
					Publisher: to.Ptr("Canonical"),
					Offer:     to.Ptr("UbuntuServer"),
					SKU:       to.Ptr("18.04-LTS"),
					Version:   to.Ptr("latest"),
				},
				OSDisk: &armcompute.OSDisk{
					Name:         to.Ptr(fmt.Sprintf("%s-osdisk", config.Name)),
					CreateOption: to.Ptr(armcompute.DiskCreateOptionTypesFromImage),
					Caching:      to.Ptr(armcompute.CachingTypesReadWrite),
					ManagedDisk: &armcompute.ManagedDiskParameters{
						StorageAccountType: to.Ptr(armcompute.StorageAccountTypesStandardLRS),
					},
					DiskSizeGB: to.Ptr(int32(config.DiskSize)),
				},
			},
			OSProfile: &armcompute.OSProfile{
				ComputerName:  to.Ptr(config.Name),
				AdminUsername: to.Ptr("azureuser"),
				LinuxConfiguration: &armcompute.LinuxConfiguration{
					DisablePasswordAuthentication: to.Ptr(true),
					SSH: &armcompute.SSHConfiguration{
						PublicKeys: []*armcompute.SSHPublicKey{
							{
								Path:    to.Ptr("/home/azureuser/.ssh/authorized_keys"),
								KeyData: to.Ptr(config.Keypair),
							},
						},
					},
				},
			},
			NetworkProfile: &armcompute.NetworkProfile{
				NetworkInterfaces: []*armcompute.NetworkInterfaceReference{
					{
						ID: to.Ptr(config.SubnetID),
						Properties: &armcompute.NetworkInterfaceReferenceProperties{
							Primary: to.Ptr(true),
						},
					},
				},
			},
		},
		Tags: map[string]*string{
			"Name": to.Ptr(config.Name),
		},
	}

	// 添加用户自定义标签
	if config.Tags != nil {
		for k, v := range config.Tags {
			parameters.Tags[k] = to.Ptr(v)
		}
	}

	poller, err := p.vmClient.BeginCreateOrUpdate(ctx, resourceGroup, config.Name, parameters, nil)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to create VM",
			err.Error(),
		)
	}

	resp, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to wait for VM creation",
			err.Error(),
		)
	}

	return p.convertAzureVMToCloudVM(resp.VirtualMachine), nil
}

// DeleteVM 删除虚拟机
func (p *AzureProvider) DeleteVM(ctx context.Context, vmID string) error {
	resourceGroup, vmName := p.parseVMID(vmID)

	poller, err := p.vmClient.BeginDelete(ctx, resourceGroup, vmName, nil)
	if err != nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to delete VM",
			err.Error(),
		)
	}

	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to wait for VM deletion",
			err.Error(),
		)
	}

	return nil
}

// StartVM 启动虚拟机
func (p *AzureProvider) StartVM(ctx context.Context, vmID string) error {
	resourceGroup, vmName := p.parseVMID(vmID)

	poller, err := p.vmClient.BeginStart(ctx, resourceGroup, vmName, nil)
	if err != nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to start VM",
			err.Error(),
		)
	}

	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to wait for VM start",
			err.Error(),
		)
	}

	return nil
}

// StopVM 停止虚拟机
func (p *AzureProvider) StopVM(ctx context.Context, vmID string) error {
	resourceGroup, vmName := p.parseVMID(vmID)

	poller, err := p.vmClient.BeginPowerOff(ctx, resourceGroup, vmName, nil)
	if err != nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to stop VM",
			err.Error(),
		)
	}

	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to wait for VM stop",
			err.Error(),
		)
	}

	return nil
}

// RebootVM 重启虚拟机
func (p *AzureProvider) RebootVM(ctx context.Context, vmID string) error {
	resourceGroup, vmName := p.parseVMID(vmID)

	poller, err := p.vmClient.BeginRestart(ctx, resourceGroup, vmName, nil)
	if err != nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to reboot VM",
			err.Error(),
		)
	}

	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		return cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to wait for VM reboot",
			err.Error(),
		)
	}

	return nil
}

// GetVMStatus 获取虚拟机状态
func (p *AzureProvider) GetVMStatus(ctx context.Context, vmID string) (*cloudprovider.VMStatus, error) {
	resourceGroup, vmName := p.parseVMID(vmID)

	resp, err := p.vmClient.InstanceView(ctx, resourceGroup, vmName, nil)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrOperationFailed,
			"failed to get VM status",
			err.Error(),
		)
	}

	status := cloudprovider.VMStatus("Unknown")

	if len(resp.Statuses) > 0 {
		for _, s := range resp.Statuses {
			if s.Code != nil {
				code := *s.Code
				if strings.HasPrefix(code, "PowerState/") {
					powerState := strings.TrimPrefix(code, "PowerState/")
					status = p.mapPowerState(powerState)
					break
				}
			}
		}
	}

	return &status, nil
}

// GetVM 获取单个虚拟机
func (p *AzureProvider) GetVM(ctx context.Context, vmID string) (*cloudprovider.VirtualMachine, error) {
	resourceGroup, vmName := p.parseVMID(vmID)

	resp, err := p.vmClient.Get(ctx, resourceGroup, vmName, nil)
	if err != nil {
		return nil, cloudprovider.NewCloudError(
			cloudprovider.ErrResourceNotFound,
			"failed to get VM",
			err.Error(),
		)
	}

	return p.convertAzureVMToCloudVM(resp.VirtualMachine), nil
}

// ResetVMPassword 重置 VM 密码
func (p *AzureProvider) ResetVMPassword(ctx context.Context, vmID, username, newPassword string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "ResetVMPassword not implemented", "")
}

// UpdateVMConfig 更新 VM 配置
func (p *AzureProvider) UpdateVMConfig(ctx context.Context, vmID, instanceType, name string) error {
	return cloudprovider.NewCloudError(cloudprovider.ErrUnsupportedOperation, "UpdateVMConfig not implemented", "")
}

// ListVMs 列出虚拟机
func (p *AzureProvider) ListVMs(ctx context.Context, filter cloudprovider.VMListFilter) ([]*cloudprovider.VirtualMachine, error) {
	resourceGroup := "opencmp-resources"

 pager := p.vmClient.NewListPager(resourceGroup, nil)

	var vms []*cloudprovider.VirtualMachine
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, cloudprovider.NewCloudError(
				cloudprovider.ErrOperationFailed,
				"failed to list VMs",
				err.Error(),
			)
		}

		for _, vm := range page.VirtualMachineListResult.Value {
			vms = append(vms, p.convertAzureVMToCloudVM(*vm))
		}
	}

	return vms, nil
}

// ListImages 列出镜像
func (p *AzureProvider) ListImages(ctx context.Context, filter cloudprovider.ImageFilter) ([]*cloudprovider.Image, error) {
	// Azure 常用镜像列表
	images := []*cloudprovider.Image{
		{ID: "Canonical:UbuntuServer:18.04-LTS:latest", Name: "Ubuntu 18.04 LTS", OSName: "Ubuntu", OSVersion: "18.04"},
		{ID: "Canonical:UbuntuServer:20.04-LTS:latest", Name: "Ubuntu 20.04 LTS", OSName: "Ubuntu", OSVersion: "20.04"},
		{ID: "Canonical:UbuntuServer:22.04-LTS:latest", Name: "Ubuntu 22.04 LTS", OSName: "Ubuntu", OSVersion: "22.04"},
		{ID: "RedHat:RHEL:7-LVM:latest", Name: "Red Hat Enterprise Linux 7", OSName: "RHEL", OSVersion: "7"},
		{ID: "RedHat:RHEL:8-LVM:latest", Name: "Red Hat Enterprise Linux 8", OSName: "RHEL", OSVersion: "8"},
		{ID: "SUSE:SLES:12-SP5:latest", Name: "SUSE Linux Enterprise Server 12 SP5", OSName: "SLES", OSVersion: "12"},
		{ID: "MicrosoftWindowsServer:WindowsServer:2019-Datacenter:latest", Name: "Windows Server 2019 Datacenter", OSName: "Windows", OSVersion: "2019"},
		{ID: "MicrosoftWindowsServer:WindowsServer:2022-Datacenter:latest", Name: "Windows Server 2022 Datacenter", OSName: "Windows", OSVersion: "2022"},
		{ID: "Debian:debian-10:10:latest", Name: "Debian 10", OSName: "Debian", OSVersion: "10"},
		{ID: "Debian:debian-11:11:latest", Name: "Debian 11", OSName: "Debian", OSVersion: "11"},
		{ID: "CentOS:CentOS:7_9:latest", Name: "CentOS 7.9", OSName: "CentOS", OSVersion: "7.9"},
	}
	return images, nil
}

// GetImage 获取镜像
func (p *AzureProvider) GetImage(ctx context.Context, imageID string) (*cloudprovider.Image, error) {
	images, err := p.ListImages(ctx, cloudprovider.ImageFilter{})
	if err != nil {
		return nil, err
	}

	for _, img := range images {
		if img.ID == imageID {
			return img, nil
		}
	}

	return nil, cloudprovider.NewCloudError(
		cloudprovider.ErrResourceNotFound,
		"image not found",
		imageID,
	)
}

// CreateKeypair 创建密钥对 (Azure 不支持，返回空实现)
func (p *AzureProvider) CreateKeypair(ctx context.Context, name, publicKey string) (*cloudprovider.Keypair, error) {
	return &cloudprovider.Keypair{
		ID:        name,
		Name:      name,
		PublicKey: publicKey,
	}, nil
}

// DeleteKeypair 删除密钥对
func (p *AzureProvider) DeleteKeypair(ctx context.Context, keypairID string) error {
	// Azure SSH 密钥在 VM 中管理，不需要单独删除
	return nil
}

// ListKeypairs 列出密钥对
func (p *AzureProvider) ListKeypairs(ctx context.Context) ([]*cloudprovider.Keypair, error) {
	return []*cloudprovider.Keypair{}, nil
}

// parseVMID 解析 VM ID
// 格式可以是: "/subscriptions/xxx/resourceGroups/yyy/providers/Microsoft.Compute/virtualMachines/vmName"
// 或者简化的: "vmName" 或 "resourceGroup:vmName"
func (p *AzureProvider) parseVMID(vmID string) (resourceGroup, vmName string) {
	if strings.Contains(vmID, "/resourceGroups/") {
		// 完整资源 ID
		parts := strings.Split(vmID, "/")
		for i, part := range parts {
			if part == "resourceGroups" && i+1 < len(parts) {
				resourceGroup = parts[i+1]
			}
			if part == "virtualMachines" && i+1 < len(parts) {
				vmName = parts[i+1]
			}
		}
	} else if strings.Contains(vmID, ":") {
		// 简化格式: "resourceGroup:vmName"
		parts := strings.Split(vmID, ":")
		resourceGroup = parts[0]
		vmName = parts[1]
	} else {
		// 仅 VM 名称
		resourceGroup = "opencmp-resources"
		vmName = vmID
	}
	return resourceGroup, vmName
}

// convertAzureVMToCloudVM 转换 Azure VM 到通用 VM 类型
func (p *AzureProvider) convertAzureVMToCloudVM(vm armcompute.VirtualMachine) *cloudprovider.VirtualMachine {
	result := &cloudprovider.VirtualMachine{
		ID:           *vm.ID,
		Name:         *vm.Name,
		InstanceType: string(*vm.Properties.HardwareProfile.VMSize),
		RegionID:     *vm.Location,
	}

	// 获取状态
	if vm.Properties != nil && vm.Properties.ProvisioningState != nil {
		result.Status = cloudprovider.VMStatus(p.mapProvisioningState(*vm.Properties.ProvisioningState))
	}

	// 获取 IP 地址 (需要从网络接口获取，这里简化处理)
	if vm.Properties != nil && vm.Properties.NetworkProfile != nil {
		for _, nic := range vm.Properties.NetworkProfile.NetworkInterfaces {
			if nic.ID != nil {
				result.SubnetID = *nic.ID
			}
		}
	}

	// 获取操作系统类型
	if vm.Properties != nil && vm.Properties.StorageProfile != nil && vm.Properties.StorageProfile.OSDisk != nil {
		if vm.Properties.StorageProfile.OSDisk.OSType != nil {
			result.OSName = string(*vm.Properties.StorageProfile.OSDisk.OSType)
		}
	}

	// 获取标签
	if vm.Tags != nil {
		result.Tags = make(map[string]string)
		for k, v := range vm.Tags {
			if v != nil {
				result.Tags[k] = *v
			}
		}
	}

	return result
}

// mapProvisioningState 映射 Provisioning 状态
func (p *AzureProvider) mapProvisioningState(state string) string {
	switch state {
	case "Succeeded":
		return "Running"
	case "Creating":
		return "Pending"
	case "Updating":
		return "Updating"
	case "Deleting":
		return "Terminating"
	case "Failed":
		return "Error"
	default:
		return "Unknown"
	}
}

// mapPowerState 映射 Power 状态
func (p *AzureProvider) mapPowerState(state string) cloudprovider.VMStatus {
	switch state {
	case "running":
		return cloudprovider.VMStatusRunning
	case "stopped":
		return cloudprovider.VMStatusStopped
	case "deallocated":
		return cloudprovider.VMStatusStopped
	case "starting":
		return cloudprovider.VMStatusStarting
	case "stopping":
		return cloudprovider.VMStatusStopping
	default:
		return cloudprovider.VMStatus("Unknown")
	}
}