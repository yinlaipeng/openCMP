package service

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/opencmp/opencmp/internal/model"
	"github.com/opencmp/opencmp/pkg/cloudprovider"
)

// GetResourceStatsFromCloud 从云厂商获取资源统计
func GetResourceStatsFromCloud(ctx context.Context, account *model.CloudAccount) (map[string]interface{}, error) {
	var creds map[string]string
	if err := json.Unmarshal(account.Credentials, &creds); err != nil {
		return nil, err
	}

	config := cloudprovider.CloudAccountConfig{
		ID:           strconv.FormatUint(uint64(account.ID), 10),
		Name:         account.Name,
		ProviderType: account.ProviderType,
		Credentials:  creds,
		Region:       creds["region_id"],
	}

	provider, err := cloudprovider.GetProvider(account.ProviderType, config)
	if err != nil {
		return nil, err
	}

	stats := map[string]int{
		"vms": 0, "rds": 0, "redis": 0, "buckets": 0,
		"eips": 0, "public_ips": 0, "snapshots": 0,
		"vpcs": 0, "subnets": 0, "total_ips": 0,
		"vms_running": 0, "eips_bound": 0, "ips_used": 0,
		"disks": 0, "disks_mounted": 0,
	}

	// 获取虚拟机
	vms, err := provider.ListVMs(ctx, cloudprovider.VMListFilter{})
	if err == nil {
		stats["vms"] = len(vms)
		for _, vm := range vms {
			if vm.Status == "running" {
				stats["vms_running"]++
			}
		}
	}

	// 获取 VPC
	vpcs, err := provider.ListVPCs(ctx, cloudprovider.VPCFilter{})
	if err == nil {
		stats["vpcs"] = len(vpcs)
	}

	// 获取子网
	subnets, err := provider.ListSubnets(ctx, cloudprovider.SubnetFilter{})
	if err == nil {
		stats["subnets"] = len(subnets)
		// TODO: 计算总IP数
	}

	// 获取 EIP
	eips, err := provider.ListEIPs(ctx, cloudprovider.EIPFilter{})
	if err == nil {
		stats["eips"] = len(eips)
		for _, eip := range eips {
			if eip.Status == "bound" || eip.Status == "attached" {
				stats["eips_bound"]++
			}
		}
	}

	// 获取安全组
	_, err = provider.ListSecurityGroups(ctx, cloudprovider.SGFilter{})
	if err == nil {
		// 不计入统计
	}

	// 获取镜像
	_, err = provider.ListImages(ctx, cloudprovider.ImageFilter{})
	if err == nil {
		// 不计入统计
	}

	// 计算使用率
	usageRates := map[string]int{
		"vm_running_rate":  0,
		"disk_mounted_rate": 0,
		"eip_bound_rate":   0,
		"ip_used_rate":     0,
	}

	if stats["vms"] > 0 {
		usageRates["vm_running_rate"] = int(float64(stats["vms_running"]) / float64(stats["vms"]) * 100)
	}
	if stats["eips"] > 0 {
		usageRates["eip_bound_rate"] = int(float64(stats["eips_bound"]) / float64(stats["eips"]) * 100)
	}

	return map[string]interface{}{
		"resources":   stats,
		"usage_rates": usageRates,
	}, nil
}