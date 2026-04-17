package service

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"go.uber.org/zap"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
	"github.com/opencmp/opencmp/pkg/cloudprovider"
)

// CloudAccountService 云账户服务
type CloudAccountService struct {
	db *gorm.DB
}

// NewCloudAccountService 创建云账户服务
func NewCloudAccountService(db *gorm.DB) *CloudAccountService {
	return &CloudAccountService{db: db}
}

// convertTagsToJSON 将map[string]string转换为datatypes.JSON
func convertTagsToJSON(tags map[string]string) datatypes.JSON {
	if tags == nil {
		return datatypes.JSON{}
	}
	tagsJSON, _ := json.Marshal(tags)
	return datatypes.JSON(tagsJSON)
}

// CreateCloudAccount 创建云账户
func (s *CloudAccountService) CreateCloudAccount(ctx context.Context, account *model.CloudAccount) error {
	return s.db.WithContext(ctx).Create(account).Error
}

// GetCloudAccount 获取云账户
func (s *CloudAccountService) GetCloudAccount(ctx context.Context, id uint) (*model.CloudAccount, error) {
	var account model.CloudAccount
	err := s.db.WithContext(ctx).First(&account, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}

// ListCloudAccounts 列出云账户
func (s *CloudAccountService) ListCloudAccounts(ctx context.Context, limit, offset int) ([]*model.CloudAccount, int64, error) {
	var accounts []*model.CloudAccount
	var total int64

	if err := s.db.Model(&model.CloudAccount{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := s.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&accounts).Error

	return accounts, total, err
}

// UpdateCloudAccount 更新云账户
func (s *CloudAccountService) UpdateCloudAccount(ctx context.Context, account *model.CloudAccount) error {
	return s.db.WithContext(ctx).Save(account).Error
}

// DeleteCloudAccount 删除云账户
func (s *CloudAccountService) DeleteCloudAccount(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Delete(&model.CloudAccount{}, id).Error
}

// VerifyCloudAccount 验证云账户
func (s *CloudAccountService) VerifyCloudAccount(ctx context.Context, account *model.CloudAccount) (bool, error) {
	var creds map[string]string
	if err := json.Unmarshal(account.Credentials, &creds); err != nil {
		return false, err
	}

	config := cloudprovider.CloudAccountConfig{
		ID:           strconv.Itoa(int(account.ID)),
		Name:         account.Name,
		ProviderType: account.ProviderType,
		Credentials:  creds,
		Region:       creds["region_id"], // 使用凭证中的region_id，如果没有则为空，适配器会使用默认值
	}

	provider, err := cloudprovider.GetProvider(account.ProviderType, config)
	if err != nil {
		return false, err
	}

	// 尝试获取云厂商信息来验证连接
	cloudInfo := provider.GetCloudInfo()
	if cloudInfo.Provider == "" {
		return false, nil
	}

	return true, nil
}

// TestConnection 测试云账户连接
func (s *CloudAccountService) TestConnection(ctx context.Context, account *model.CloudAccount) (bool, error) {
	var creds map[string]string
	if err := json.Unmarshal(account.Credentials, &creds); err != nil {
		return false, err
	}

	config := cloudprovider.CloudAccountConfig{
		ID:           strconv.Itoa(int(account.ID)),
		Name:         account.Name,
		ProviderType: account.ProviderType,
		Credentials:  creds,
		Region:       creds["region_id"], // 使用凭证中的region_id，如果没有则为空，适配器会使用默认值
	}

	provider, err := cloudprovider.GetProvider(account.ProviderType, config)
	if err != nil {
		return false, err
	}

	// 尝试获取云厂商信息来测试连接
	cloudInfo := provider.GetCloudInfo()
	if cloudInfo.Provider == "" {
		return false, nil
	}

	return true, nil
}

// SyncResources 同步云账户资源
// 返回同步的资源统计（VM数量、VPC数量等）
func (s *CloudAccountService) SyncResources(ctx context.Context, account *model.CloudAccount) (map[string]int, error) {
	return s.SyncResourcesWithMode(ctx, account, model.SyncModeIncremental, model.SyncTriggerManual, nil, nil)
}

// SyncResourcesWithMode 同步云账户资源（支持指定同步模式和资源类型）
// syncMode: incremental（增量）或 full（全量）
// triggeredBy: manual（手动）或 scheduled（定时任务）
// scheduledTaskID: 定时任务ID（如果通过定时任务触发）
// resourceTypes: 要同步的资源类型列表，空表示同步所有类型
func (s *CloudAccountService) SyncResourcesWithMode(ctx context.Context, account *model.CloudAccount, syncMode model.SyncMode, triggeredBy model.SyncTriggerType, scheduledTaskID *uint, resourceTypes []string) (map[string]int, error) {
	var creds map[string]string
	if err := json.Unmarshal(account.Credentials, &creds); err != nil {
		return nil, err
	}

	config := cloudprovider.CloudAccountConfig{
		ID:           strconv.Itoa(int(account.ID)),
		Name:         account.Name,
		ProviderType: account.ProviderType,
		Credentials:  creds,
		Region:       creds["region_id"],
	}

	provider, err := cloudprovider.GetProvider(account.ProviderType, config)
	if err != nil {
		return nil, err
	}

	// 初始化服务
	logger := zap.NewNop() // 使用默认logger
	mappingService := NewResourceMappingService(s.db, logger)
	syncLogService := NewSyncLogService(s.db, logger)

	stats := map[string]int{}
	totalNew := 0
	totalUpdated := 0
	totalDeleted := 0
	totalSkipped := 0
	totalErrors := 0

	// Helper function to check if a resource type should be synced
	shouldSync := func(resourceType string) bool {
		if len(resourceTypes) == 0 {
			return true // No filter means sync all
		}
		for _, rt := range resourceTypes {
			if rt == resourceType || rt == "all" {
				return true
			}
		}
		return false
	}

	// 同步虚拟机
	if shouldSync("vm") {
		vms, err := provider.ListVMs(ctx, cloudprovider.VMListFilter{})
		if err == nil {
			newCount, updatedCount, deletedCount, skippedCount, errorCount := s.syncVMs(ctx, account, vms, syncMode, mappingService)
			stats["vms"] = len(vms)
			stats["vms_new"] = newCount
			stats["vms_updated"] = updatedCount
			stats["vms_deleted"] = deletedCount
			stats["vms_skipped"] = skippedCount
			stats["vms_errors"] = errorCount
			totalNew += newCount
			totalUpdated += updatedCount
			totalDeleted += deletedCount
			totalSkipped += skippedCount
			totalErrors += errorCount
		}
	}

	// 同步 VPC
	if shouldSync("vpc") {
		vpcs, err := provider.ListVPCs(ctx, cloudprovider.VPCFilter{})
		if err == nil {
			newCount, updatedCount, deletedCount, skippedCount, errorCount := s.syncVPCs(ctx, account, vpcs, syncMode, mappingService)
			stats["vpcs"] = len(vpcs)
			stats["vpcs_new"] = newCount
			stats["vpcs_updated"] = updatedCount
			stats["vpcs_deleted"] = deletedCount
			stats["vpcs_skipped"] = skippedCount
			stats["vpcs_errors"] = errorCount
			totalNew += newCount
			totalUpdated += updatedCount
			totalDeleted += deletedCount
			totalSkipped += skippedCount
			totalErrors += errorCount
		}
	}

	// 同步子网
	if shouldSync("subnet") {
		subnets, err := provider.ListSubnets(ctx, cloudprovider.SubnetFilter{})
		if err == nil {
			newCount, updatedCount, deletedCount, skippedCount, errorCount := s.syncSubnets(ctx, account, subnets, syncMode, mappingService)
			stats["subnets"] = len(subnets)
			stats["subnets_new"] = newCount
			stats["subnets_updated"] = updatedCount
			stats["subnets_deleted"] = deletedCount
			stats["subnets_skipped"] = skippedCount
			stats["subnets_errors"] = errorCount
			totalNew += newCount
			totalUpdated += updatedCount
			totalDeleted += deletedCount
			totalSkipped += skippedCount
			totalErrors += errorCount
		}
	}

	// 同步安全组
	if shouldSync("security_group") {
		sgs, err := provider.ListSecurityGroups(ctx, cloudprovider.SGFilter{})
		if err == nil {
			newCount, updatedCount, deletedCount, skippedCount, errorCount := s.syncSecurityGroups(ctx, account, sgs, syncMode, mappingService)
			stats["security_groups"] = len(sgs)
			stats["security_groups_new"] = newCount
			stats["security_groups_updated"] = updatedCount
			stats["security_groups_deleted"] = deletedCount
			stats["security_groups_skipped"] = skippedCount
			stats["security_groups_errors"] = errorCount
			totalNew += newCount
			totalUpdated += updatedCount
			totalDeleted += deletedCount
			totalSkipped += skippedCount
			totalErrors += errorCount
		}
	}

	// 同步 EIP
	if shouldSync("eip") {
		eips, err := provider.ListEIPs(ctx, cloudprovider.EIPFilter{})
		if err == nil {
			newCount, updatedCount, deletedCount, skippedCount, errorCount := s.syncEIPs(ctx, account, eips, syncMode, mappingService)
			stats["eips"] = len(eips)
			stats["eips_new"] = newCount
			stats["eips_updated"] = updatedCount
			stats["eips_deleted"] = deletedCount
			stats["eips_skipped"] = skippedCount
			stats["eips_errors"] = errorCount
			totalNew += newCount
			totalUpdated += updatedCount
			totalDeleted += deletedCount
			totalSkipped += skippedCount
			totalErrors += errorCount
		}
	}

	// 同步镜像
	if shouldSync("image") {
		images, err := provider.ListImages(ctx, cloudprovider.ImageFilter{})
		if err == nil {
			newCount, updatedCount, deletedCount, skippedCount, errorCount := s.syncImages(ctx, account, images, syncMode)
			stats["images"] = len(images)
			stats["images_new"] = newCount
			stats["images_updated"] = updatedCount
			stats["images_deleted"] = deletedCount
			stats["images_skipped"] = skippedCount
			stats["images_errors"] = errorCount
			totalNew += newCount
			totalUpdated += updatedCount
			totalDeleted += deletedCount
			totalSkipped += skippedCount
			totalErrors += errorCount
		}
	}

	// 同步云硬盘（如果provider支持IStorage接口）
	if shouldSync("disk") {
		disks, err := provider.ListDisks(ctx, cloudprovider.DiskFilter{})
		if err == nil {
			newCount, updatedCount, deletedCount, skippedCount, errorCount := s.syncDisks(ctx, account, disks, syncMode, mappingService)
			stats["disks"] = len(disks)
			stats["disks_new"] = newCount
			stats["disks_updated"] = updatedCount
			stats["disks_deleted"] = deletedCount
			stats["disks_skipped"] = skippedCount
			stats["disks_errors"] = errorCount
			totalNew += newCount
			totalUpdated += updatedCount
			totalDeleted += deletedCount
			totalSkipped += skippedCount
			totalErrors += errorCount
		}
	}

	// 同步云快照（如果provider支持IStorage接口）
	if shouldSync("snapshot") {
		snapshots, err := provider.ListSnapshots(ctx, cloudprovider.SnapshotFilter{})
		if err == nil {
			newCount, updatedCount, deletedCount, skippedCount, errorCount := s.syncSnapshots(ctx, account, snapshots, syncMode)
			stats["snapshots"] = len(snapshots)
			stats["snapshots_new"] = newCount
			stats["snapshots_updated"] = updatedCount
			stats["snapshots_deleted"] = deletedCount
			stats["snapshots_skipped"] = skippedCount
			stats["snapshots_errors"] = errorCount
			totalNew += newCount
			totalUpdated += updatedCount
			totalDeleted += deletedCount
			totalSkipped += skippedCount
			totalErrors += errorCount
		}
	}

	// 同步RDS数据库（如果provider支持IDatabase接口）
	if shouldSync("rds") {
		rdsInstances, err := provider.ListRDSInstances(ctx, cloudprovider.RDSFilter{})
		if err == nil {
			newCount, updatedCount, deletedCount, skippedCount, errorCount := s.syncRDS(ctx, account, rdsInstances, syncMode, mappingService)
			stats["rds"] = len(rdsInstances)
			stats["rds_new"] = newCount
			stats["rds_updated"] = updatedCount
			stats["rds_deleted"] = deletedCount
			stats["rds_skipped"] = skippedCount
			stats["rds_errors"] = errorCount
			totalNew += newCount
			totalUpdated += updatedCount
			totalDeleted += deletedCount
			totalSkipped += skippedCount
			totalErrors += errorCount
		}
	}

	// 同步Redis缓存（如果provider支持IDatabase接口）
	if shouldSync("redis") {
		redisInstances, err := provider.ListCacheInstances(ctx, cloudprovider.CacheFilter{})
		if err == nil {
			newCount, updatedCount, deletedCount, skippedCount, errorCount := s.syncRedis(ctx, account, redisInstances, syncMode, mappingService)
			stats["redis"] = len(redisInstances)
			stats["redis_new"] = newCount
			stats["redis_updated"] = updatedCount
			stats["redis_deleted"] = deletedCount
			stats["redis_skipped"] = skippedCount
			stats["redis_errors"] = errorCount
			totalNew += newCount
			totalUpdated += updatedCount
			totalDeleted += deletedCount
			totalSkipped += skippedCount
			totalErrors += errorCount
		}
	}

	// 创建同步日志（汇总）
	syncLog, _ := syncLogService.StartSyncLog(ctx, account.ID, account.Name, syncMode, "all", triggeredBy, scheduledTaskID)
	if syncLog != nil {
		syncLogService.UpdateSyncLogProgress(ctx, syncLog.ID, totalNew, totalUpdated, totalDeleted, totalSkipped, totalErrors)
		status := model.SyncLogStatusSuccess
		if totalErrors > 0 {
			status = model.SyncLogStatusPartialFail
		}
		syncLogService.CompleteSyncLog(ctx, syncLog.ID, status, "")
	}

	// 更新同步状态
	account.Status = string(model.CloudAccountStatusActive)
	now := time.Now()
	account.LastSync = &now
	s.db.WithContext(ctx).Save(account)

	return stats, nil
}

// syncVMs 同步虚拟机资源
func (s *CloudAccountService) syncVMs(ctx context.Context, account *model.CloudAccount, vms []*cloudprovider.VirtualMachine, syncMode model.SyncMode, mappingService *ResourceMappingService) (int, int, int, int, int) {
	newCount := 0
	updatedCount := 0
	deletedCount := 0
	skippedCount := 0
	errorCount := 0

	for _, vm := range vms {
		// 确定项目归属
		attribution, err := mappingService.DetermineProjectAttribution(ctx, account.ID, vm.Tags)
		if err != nil {
			errorCount++
			continue
		}

		// 查询本地是否存在该资源
		var existingVM model.CloudVM
		err = s.db.WithContext(ctx).
			Where("cloud_account_id = ?", account.ID).
			Where("instance_id = ?", vm.ID).
			First(&existingVM).Error

		if err == gorm.ErrRecordNotFound {
			// 新资源，创建
			tagsJSON, _ := json.Marshal(vm.Tags)
			newVM := model.CloudVM{
				CloudAccountID: account.ID,
				InstanceID:     vm.ID,
				Name:           vm.Name,
				Status:         string(vm.Status),
				InstanceType:   vm.InstanceType,
				ImageID:        vm.ImageID,
				OSName:         vm.OSName,
				VPCID:          vm.VPCID,
				SubnetID:       vm.SubnetID,
				PrivateIP:      vm.PrivateIP,
				PublicIP:       vm.PublicIP,
				RegionID:       vm.RegionID,
				ZoneID:         vm.ZoneID,
				ProjectID:      attribution.ProjectID,
				Tags:           datatypes.JSON(tagsJSON),
			}
			if s.db.WithContext(ctx).Create(&newVM).Error == nil {
				newCount++
			} else {
				errorCount++
			}
		} else if err == nil {
			// 已存在，根据同步模式处理
			if syncMode == model.SyncModeFull {
				// 全量同步：更新状态
				tagsJSON, _ := json.Marshal(vm.Tags)
				changes := map[string]interface{}{
					"status":        string(vm.Status),
					"instance_type": vm.InstanceType,
					"private_ip":    vm.PrivateIP,
					"public_ip":     vm.PublicIP,
					"project_id":    attribution.ProjectID,
					"tags":          datatypes.JSON(tagsJSON),
				}
				if s.db.WithContext(ctx).Model(&existingVM).Updates(changes).Error == nil {
					updatedCount++
				} else {
					errorCount++
				}
			} else {
				// 增量同步：跳过已存在的资源
				skippedCount++
			}
		} else {
			errorCount++
		}
	}

	// 全量同步：标记云平台已删除的资源
	if syncMode == model.SyncModeFull {
		var localVMs []model.CloudVM
		s.db.WithContext(ctx).
			Where("cloud_account_id = ?", account.ID).
			Where("status != ?", "terminated").
			Find(&localVMs)

		cloudVMIDs := make(map[string]bool)
		for _, vm := range vms {
			cloudVMIDs[vm.ID] = true
		}

		for _, localVM := range localVMs {
			if !cloudVMIDs[localVM.InstanceID] {
				// 云平台已删除，标记为terminated
				s.db.WithContext(ctx).Model(&localVM).Update("status", "terminated")
				deletedCount++
			}
		}
	}

	return newCount, updatedCount, deletedCount, skippedCount, errorCount
}

// syncVPCs 同步VPC资源
func (s *CloudAccountService) syncVPCs(ctx context.Context, account *model.CloudAccount, vpcs []*cloudprovider.VPC, syncMode model.SyncMode, mappingService *ResourceMappingService) (int, int, int, int, int) {
	newCount := 0
	updatedCount := 0
	deletedCount := 0
	skippedCount := 0
	errorCount := 0

	for _, vpc := range vpcs {
		attribution, err := mappingService.DetermineProjectAttribution(ctx, account.ID, vpc.Tags)
		if err != nil {
			errorCount++
			continue
		}

		var existingVPC model.CloudVPC
		err = s.db.WithContext(ctx).
			Where("cloud_account_id = ?", account.ID).
			Where("vpc_id = ?", vpc.ID).
			First(&existingVPC).Error

		if err == gorm.ErrRecordNotFound {
			newVPC := model.CloudVPC{
				CloudAccountID: account.ID,
				VPCID:          vpc.ID,
				Name:           vpc.Name,
				CIDR:           vpc.CIDR,
				Status:         vpc.Status,
				RegionID:       vpc.RegionID,
				ProjectID:      attribution.ProjectID,
				Tags:           convertTagsToJSON(vpc.Tags),
			}
			if s.db.WithContext(ctx).Create(&newVPC).Error == nil {
				newCount++
			} else {
				errorCount++
			}
		} else if err == nil {
			if syncMode == model.SyncModeFull {
				changes := map[string]interface{}{
					"status":     vpc.Status,
					"cidr":       vpc.CIDR,
					"project_id": attribution.ProjectID,
					"tags":       convertTagsToJSON(vpc.Tags),
				}
				if s.db.WithContext(ctx).Model(&existingVPC).Updates(changes).Error == nil {
					updatedCount++
				} else {
					errorCount++
				}
			} else {
				skippedCount++
			}
		} else {
			errorCount++
		}
	}

	if syncMode == model.SyncModeFull {
		var localVPCs []model.CloudVPC
		s.db.WithContext(ctx).
			Where("cloud_account_id = ?", account.ID).
			Where("status != ?", "terminated").
			Find(&localVPCs)

		cloudVPCIDs := make(map[string]bool)
		for _, vpc := range vpcs {
			cloudVPCIDs[vpc.ID] = true
		}

		for _, localVPC := range localVPCs {
			if !cloudVPCIDs[localVPC.VPCID] {
				s.db.WithContext(ctx).Model(&localVPC).Update("status", "terminated")
				deletedCount++
			}
		}
	}

	return newCount, updatedCount, deletedCount, skippedCount, errorCount
}

// syncSubnets 同步子网资源
func (s *CloudAccountService) syncSubnets(ctx context.Context, account *model.CloudAccount, subnets []*cloudprovider.Subnet, syncMode model.SyncMode, mappingService *ResourceMappingService) (int, int, int, int, int) {
	newCount := 0
	updatedCount := 0
	deletedCount := 0
	skippedCount := 0
	errorCount := 0

	for _, subnet := range subnets {
		attribution, err := mappingService.DetermineProjectAttribution(ctx, account.ID, subnet.Tags)
		if err != nil {
			errorCount++
			continue
		}

		var existingSubnet model.CloudSubnet
		err = s.db.WithContext(ctx).
			Where("cloud_account_id = ?", account.ID).
			Where("subnet_id = ?", subnet.ID).
			First(&existingSubnet).Error

		if err == gorm.ErrRecordNotFound {
			newSubnet := model.CloudSubnet{
				CloudAccountID: account.ID,
				SubnetID:       subnet.ID,
				Name:           subnet.Name,
				VPCID:          subnet.VPCID,
				CIDR:           subnet.CIDR,
				ZoneID:         subnet.ZoneID,
				Status:         subnet.Status,
				ProjectID:      attribution.ProjectID,
				Tags:           convertTagsToJSON(subnet.Tags),
			}
			if s.db.WithContext(ctx).Create(&newSubnet).Error == nil {
				newCount++
			} else {
				errorCount++
			}
		} else if err == nil {
			if syncMode == model.SyncModeFull {
				changes := map[string]interface{}{
					"status":     subnet.Status,
					"cidr":       subnet.CIDR,
					"project_id": attribution.ProjectID,
					"tags":       convertTagsToJSON(subnet.Tags),
				}
				if s.db.WithContext(ctx).Model(&existingSubnet).Updates(changes).Error == nil {
					updatedCount++
				} else {
					errorCount++
				}
			} else {
				skippedCount++
			}
		} else {
			errorCount++
		}
	}

	if syncMode == model.SyncModeFull {
		var localSubnets []model.CloudSubnet
		s.db.WithContext(ctx).
			Where("cloud_account_id = ?", account.ID).
			Where("status != ?", "terminated").
			Find(&localSubnets)

		cloudSubnetIDs := make(map[string]bool)
		for _, subnet := range subnets {
			cloudSubnetIDs[subnet.ID] = true
		}

		for _, localSubnet := range localSubnets {
			if !cloudSubnetIDs[localSubnet.SubnetID] {
				s.db.WithContext(ctx).Model(&localSubnet).Update("status", "terminated")
				deletedCount++
			}
		}
	}

	return newCount, updatedCount, deletedCount, skippedCount, errorCount
}

// syncSecurityGroups 同步安全组资源
func (s *CloudAccountService) syncSecurityGroups(ctx context.Context, account *model.CloudAccount, sgs []*cloudprovider.SecurityGroup, syncMode model.SyncMode, mappingService *ResourceMappingService) (int, int, int, int, int) {
	newCount := 0
	updatedCount := 0
	deletedCount := 0
	skippedCount := 0
	errorCount := 0

	for _, sg := range sgs {
		attribution, err := mappingService.DetermineProjectAttribution(ctx, account.ID, sg.Tags)
		if err != nil {
			errorCount++
			continue
		}

		var existingSG model.CloudSecurityGroup
		err = s.db.WithContext(ctx).
			Where("cloud_account_id = ?", account.ID).
			Where("security_group_id = ?", sg.ID).
			First(&existingSG).Error

		if err == gorm.ErrRecordNotFound {
			newSG := model.CloudSecurityGroup{
				CloudAccountID:     account.ID,
				SecurityGroupID:    sg.ID,
				Name:               sg.Name,
				Description:        sg.Description,
				VPCID:              sg.VPCID,
				Status:             "available", // SecurityGroup没有Status字段，使用默认值
				ProjectID:          attribution.ProjectID,
				Tags:               convertTagsToJSON(sg.Tags),
			}
			if s.db.WithContext(ctx).Create(&newSG).Error == nil {
				newCount++
			} else {
				errorCount++
			}
		} else if err == nil {
			if syncMode == model.SyncModeFull {
				changes := map[string]interface{}{
					"project_id": attribution.ProjectID,
					"tags":       convertTagsToJSON(sg.Tags),
				}
				if s.db.WithContext(ctx).Model(&existingSG).Updates(changes).Error == nil {
					updatedCount++
				} else {
					errorCount++
				}
			} else {
				skippedCount++
			}
		} else {
			errorCount++
		}
	}

	if syncMode == model.SyncModeFull {
		var localSGs []model.CloudSecurityGroup
		s.db.WithContext(ctx).
			Where("cloud_account_id = ?", account.ID).
			Where("status != ?", "terminated").
			Find(&localSGs)

		cloudSGIDs := make(map[string]bool)
		for _, sg := range sgs {
			cloudSGIDs[sg.ID] = true
		}

		for _, localSG := range localSGs {
			if !cloudSGIDs[localSG.SecurityGroupID] {
				s.db.WithContext(ctx).Model(&localSG).Update("status", "terminated")
				deletedCount++
			}
		}
	}

	return newCount, updatedCount, deletedCount, skippedCount, errorCount
}

// syncEIPs 同步EIP资源
func (s *CloudAccountService) syncEIPs(ctx context.Context, account *model.CloudAccount, eips []*cloudprovider.EIP, syncMode model.SyncMode, mappingService *ResourceMappingService) (int, int, int, int, int) {
	newCount := 0
	updatedCount := 0
	deletedCount := 0
	skippedCount := 0
	errorCount := 0

	for _, eip := range eips {
		// EIP没有Tags字段，使用默认项目归属
		attribution, err := mappingService.DetermineProjectAttribution(ctx, account.ID, map[string]string{})
		if err != nil {
			errorCount++
			continue
		}

		var existingEIP model.CloudEIP
		err = s.db.WithContext(ctx).
			Where("cloud_account_id = ?", account.ID).
			Where("eip_id = ?", eip.ID).
			First(&existingEIP).Error

		if err == gorm.ErrRecordNotFound {
			newEIP := model.CloudEIP{
				CloudAccountID: account.ID,
				EIPID:          eip.ID,
				Address:        eip.Address,
				Bandwidth:      eip.Bandwidth,
				Status:         eip.Status,
				ResourceID:     eip.ResourceID,
				ResourceType:   eip.ResourceType,
				RegionID:       eip.RegionID,
				ProjectID:      attribution.ProjectID,
			}
			if s.db.WithContext(ctx).Create(&newEIP).Error == nil {
				newCount++
			} else {
				errorCount++
			}
		} else if err == nil {
			if syncMode == model.SyncModeFull {
				changes := map[string]interface{}{
					"status":      eip.Status,
					"bandwidth":   eip.Bandwidth,
					"resource_id": eip.ResourceID,
					"project_id":  attribution.ProjectID,
				}
				if s.db.WithContext(ctx).Model(&existingEIP).Updates(changes).Error == nil {
					updatedCount++
				} else {
					errorCount++
				}
			} else {
				skippedCount++
			}
		} else {
			errorCount++
		}
	}

	if syncMode == model.SyncModeFull {
		var localEIPs []model.CloudEIP
		s.db.WithContext(ctx).
			Where("cloud_account_id = ?", account.ID).
			Where("status != ?", "terminated").
			Find(&localEIPs)

		cloudEIPIDs := make(map[string]bool)
		for _, eip := range eips {
			cloudEIPIDs[eip.ID] = true
		}

		for _, localEIP := range localEIPs {
			if !cloudEIPIDs[localEIP.EIPID] {
				s.db.WithContext(ctx).Model(&localEIP).Update("status", "terminated")
				deletedCount++
			}
		}
	}

	return newCount, updatedCount, deletedCount, skippedCount, errorCount
}

// syncImages 同步镜像资源（镜像不需要项目归属映射）
func (s *CloudAccountService) syncImages(ctx context.Context, account *model.CloudAccount, images []*cloudprovider.Image, syncMode model.SyncMode) (int, int, int, int, int) {
	newCount := 0
	updatedCount := 0
	deletedCount := 0
	skippedCount := 0
	errorCount := 0

	for _, img := range images {
		var existingImg model.CloudImage
		err := s.db.WithContext(ctx).
			Where("cloud_account_id = ?", account.ID).
			Where("image_id = ?", img.ID).
			First(&existingImg).Error

		if err == gorm.ErrRecordNotFound {
			newImg := model.CloudImage{
				CloudAccountID: account.ID,
				ImageID:        img.ID,
				Name:           img.Name,
				OSName:         img.OSName,
				OSVersion:      img.OSVersion,
				Architecture:   img.Architecture,
				Status:         img.Status,
				Tags:           convertTagsToJSON(img.Tags),
			}
			if s.db.WithContext(ctx).Create(&newImg).Error == nil {
				newCount++
			} else {
				errorCount++
			}
		} else if err == nil {
			if syncMode == model.SyncModeFull {
				changes := map[string]interface{}{
					"status":      img.Status,
					"architecture": img.Architecture,
					"tags":        convertTagsToJSON(img.Tags),
				}
				if s.db.WithContext(ctx).Model(&existingImg).Updates(changes).Error == nil {
					updatedCount++
				} else {
					errorCount++
				}
			} else {
				skippedCount++
			}
		} else {
			errorCount++
		}
	}

	if syncMode == model.SyncModeFull {
		var localImages []model.CloudImage
		s.db.WithContext(ctx).
			Where("cloud_account_id = ?", account.ID).
			Where("status != ?", "terminated").
			Find(&localImages)

		cloudImageIDs := make(map[string]bool)
		for _, img := range images {
			cloudImageIDs[img.ID] = true
		}

		for _, localImg := range localImages {
			if !cloudImageIDs[localImg.ImageID] {
				s.db.WithContext(ctx).Model(&localImg).Update("status", "terminated")
				deletedCount++
			}
		}
	}

	return newCount, updatedCount, deletedCount, skippedCount, errorCount
}

// syncDisks 同步云硬盘资源
func (s *CloudAccountService) syncDisks(ctx context.Context, account *model.CloudAccount, disks []*cloudprovider.Disk, syncMode model.SyncMode, mappingService *ResourceMappingService) (int, int, int, int, int) {
	newCount := 0
	updatedCount := 0
	deletedCount := 0
	skippedCount := 0
	errorCount := 0

	for _, disk := range disks {
		attribution, err := mappingService.DetermineProjectAttribution(ctx, account.ID, disk.Tags)
		if err != nil {
			errorCount++
			continue
		}

		var existingDisk model.CloudDisk
		err = s.db.WithContext(ctx).
			Where("cloud_account_id = ?", account.ID).
			Where("disk_id = ?", disk.ID).
			First(&existingDisk).Error

		if err == gorm.ErrRecordNotFound {
			newDisk := model.CloudDisk{
				CloudAccountID: account.ID,
				DiskID:         disk.ID,
				Name:           disk.Name,
				Size:           disk.Size,
				Type:           disk.Type,
				Status:         disk.Status,
				VMID:           disk.VMID,
				ZoneID:         disk.ZoneID,
				RegionID:       disk.ZoneID[:len(disk.ZoneID)-1], // 从ZoneID提取RegionID
				ProviderType:   account.ProviderType,
				ProjectID:      attribution.ProjectID,
				Tags:           convertTagsToJSON(disk.Tags),
			}
			if s.db.WithContext(ctx).Create(&newDisk).Error == nil {
				newCount++
			} else {
				errorCount++
			}
		} else if err == nil {
			if syncMode == model.SyncModeFull {
				changes := map[string]interface{}{
					"status":     disk.Status,
					"vm_id":      disk.VMID,
					"project_id": attribution.ProjectID,
					"tags":       convertTagsToJSON(disk.Tags),
				}
				if s.db.WithContext(ctx).Model(&existingDisk).Updates(changes).Error == nil {
					updatedCount++
				} else {
					errorCount++
				}
			} else {
				skippedCount++
			}
		} else {
			errorCount++
		}
	}

	if syncMode == model.SyncModeFull {
		var localDisks []model.CloudDisk
		s.db.WithContext(ctx).
			Where("cloud_account_id = ?", account.ID).
			Where("status != ?", "terminated").
			Find(&localDisks)

		cloudDiskIDs := make(map[string]bool)
		for _, disk := range disks {
			cloudDiskIDs[disk.ID] = true
		}

		for _, localDisk := range localDisks {
			if !cloudDiskIDs[localDisk.DiskID] {
				s.db.WithContext(ctx).Model(&localDisk).Update("status", "terminated")
				deletedCount++
			}
		}
	}

	return newCount, updatedCount, deletedCount, skippedCount, errorCount
}

// syncSnapshots 同步云快照资源
func (s *CloudAccountService) syncSnapshots(ctx context.Context, account *model.CloudAccount, snapshots []*cloudprovider.Snapshot, syncMode model.SyncMode) (int, int, int, int, int) {
	newCount := 0
	updatedCount := 0
	deletedCount := 0
	skippedCount := 0
	errorCount := 0

	for _, snapshot := range snapshots {
		var existingSnapshot model.CloudSnapshot
		err := s.db.WithContext(ctx).
			Where("cloud_account_id = ?", account.ID).
			Where("snapshot_id = ?", snapshot.ID).
			First(&existingSnapshot).Error

		if err == gorm.ErrRecordNotFound {
			newSnapshot := model.CloudSnapshot{
				CloudAccountID: account.ID,
				SnapshotID:     snapshot.ID,
				Name:           snapshot.Name,
				DiskID:         snapshot.DiskID,
				Size:           snapshot.Size,
				Status:         snapshot.Status,
				ProviderType:   account.ProviderType,
			}
			if s.db.WithContext(ctx).Create(&newSnapshot).Error == nil {
				newCount++
			} else {
				errorCount++
			}
		} else if err == nil {
			if syncMode == model.SyncModeFull {
				changes := map[string]interface{}{
					"status": snapshot.Status,
				}
				if s.db.WithContext(ctx).Model(&existingSnapshot).Updates(changes).Error == nil {
					updatedCount++
				} else {
					errorCount++
				}
			} else {
				skippedCount++
			}
		} else {
			errorCount++
		}
	}

	if syncMode == model.SyncModeFull {
		var localSnapshots []model.CloudSnapshot
		s.db.WithContext(ctx).
			Where("cloud_account_id = ?", account.ID).
			Where("status != ?", "terminated").
			Find(&localSnapshots)

		cloudSnapshotIDs := make(map[string]bool)
		for _, snapshot := range snapshots {
			cloudSnapshotIDs[snapshot.ID] = true
		}

		for _, localSnapshot := range localSnapshots {
			if !cloudSnapshotIDs[localSnapshot.SnapshotID] {
				s.db.WithContext(ctx).Model(&localSnapshot).Update("status", "terminated")
				deletedCount++
			}
		}
	}

	return newCount, updatedCount, deletedCount, skippedCount, errorCount
}

// syncRDS 同步RDS数据库资源
func (s *CloudAccountService) syncRDS(ctx context.Context, account *model.CloudAccount, rdsInstances []*cloudprovider.RDSInstance, syncMode model.SyncMode, mappingService *ResourceMappingService) (int, int, int, int, int) {
	newCount := 0
	updatedCount := 0
	deletedCount := 0
	skippedCount := 0
	errorCount := 0

	for _, rds := range rdsInstances {
		attribution, err := mappingService.DetermineProjectAttribution(ctx, account.ID, rds.Tags)
		if err != nil {
			errorCount++
			continue
		}

		var existingRDS model.CloudRDS
		err = s.db.WithContext(ctx).
			Where("cloud_account_id = ?", account.ID).
			Where("rds_id = ?", rds.ID).
			First(&existingRDS).Error

		if err == gorm.ErrRecordNotFound {
			newRDS := model.CloudRDS{
				CloudAccountID: account.ID,
				RDSID:          rds.ID,
				Name:           rds.Name,
				Engine:         rds.Engine,
				EngineVersion:  rds.EngineVersion,
				InstanceType:   rds.InstanceType,
				Status:         rds.Status,
				ZoneID:         rds.ZoneID,
				ProjectID:      attribution.ProjectID,
				Tags:           convertTagsToJSON(rds.Tags),
			}
			if s.db.WithContext(ctx).Create(&newRDS).Error == nil {
				newCount++
			} else {
				errorCount++
			}
		} else if err == nil {
			if syncMode == model.SyncModeFull {
				changes := map[string]interface{}{
					"status":      rds.Status,
					"project_id":  attribution.ProjectID,
					"tags":        convertTagsToJSON(rds.Tags),
				}
				if s.db.WithContext(ctx).Model(&existingRDS).Updates(changes).Error == nil {
					updatedCount++
				} else {
					errorCount++
				}
			} else {
				skippedCount++
			}
		} else {
			errorCount++
		}
	}

	if syncMode == model.SyncModeFull {
		var localRDS []model.CloudRDS
		s.db.WithContext(ctx).
			Where("cloud_account_id = ?", account.ID).
			Where("status != ?", "terminated").
			Find(&localRDS)

		cloudRDSIDs := make(map[string]bool)
		for _, rds := range rdsInstances {
			cloudRDSIDs[rds.ID] = true
		}

		for _, localRDS := range localRDS {
			if !cloudRDSIDs[localRDS.RDSID] {
				s.db.WithContext(ctx).Model(&localRDS).Update("status", "terminated")
				deletedCount++
			}
		}
	}

	return newCount, updatedCount, deletedCount, skippedCount, errorCount
}

// syncRedis 同步Redis缓存资源
func (s *CloudAccountService) syncRedis(ctx context.Context, account *model.CloudAccount, redisInstances []*cloudprovider.CacheInstance, syncMode model.SyncMode, mappingService *ResourceMappingService) (int, int, int, int, int) {
	newCount := 0
	updatedCount := 0
	deletedCount := 0
	skippedCount := 0
	errorCount := 0

	for _, redis := range redisInstances {
		attribution, err := mappingService.DetermineProjectAttribution(ctx, account.ID, redis.Tags)
		if err != nil {
			errorCount++
			continue
		}

		var existingRedis model.CloudRedis
		err = s.db.WithContext(ctx).
			Where("cloud_account_id = ?", account.ID).
			Where("redis_id = ?", redis.ID).
			First(&existingRedis).Error

		if err == gorm.ErrRecordNotFound {
			newRedis := model.CloudRedis{
				CloudAccountID: account.ID,
				RedisID:        redis.ID,
				Name:           redis.Name,
				InstanceType:   redis.InstanceType,
				Status:         redis.Status,
				ZoneID:         redis.ZoneID,
				ProjectID:      attribution.ProjectID,
				Tags:           convertTagsToJSON(redis.Tags),
			}
			if s.db.WithContext(ctx).Create(&newRedis).Error == nil {
				newCount++
			} else {
				errorCount++
			}
		} else if err == nil {
			if syncMode == model.SyncModeFull {
				changes := map[string]interface{}{
					"status":      redis.Status,
					"project_id":  attribution.ProjectID,
					"tags":        convertTagsToJSON(redis.Tags),
				}
				if s.db.WithContext(ctx).Model(&existingRedis).Updates(changes).Error == nil {
					updatedCount++
				} else {
					errorCount++
				}
			} else {
				skippedCount++
			}
		} else {
			errorCount++
		}
	}

	if syncMode == model.SyncModeFull {
		var localRedis []model.CloudRedis
		s.db.WithContext(ctx).
			Where("cloud_account_id = ?", account.ID).
			Where("status != ?", "terminated").
			Find(&localRedis)

		cloudRedisIDs := make(map[string]bool)
		for _, redis := range redisInstances {
			cloudRedisIDs[redis.ID] = true
		}

		for _, localRedis := range localRedis {
			if !cloudRedisIDs[localRedis.RedisID] {
				s.db.WithContext(ctx).Model(&localRedis).Update("status", "terminated")
				deletedCount++
			}
		}
	}

	return newCount, updatedCount, deletedCount, skippedCount, errorCount
}

// VerifyCredentials 实际验证云账户凭证（通过调用真实 API）
func (s *CloudAccountService) VerifyCredentials(ctx context.Context, account *model.CloudAccount) (bool, string, error) {
	var creds map[string]string
	if err := json.Unmarshal(account.Credentials, &creds); err != nil {
		return false, "invalid credentials format", err
	}

	config := cloudprovider.CloudAccountConfig{
		ID:           strconv.Itoa(int(account.ID)),
		Name:         account.Name,
		ProviderType: account.ProviderType,
		Credentials:  creds,
		Region:       creds["region_id"],
	}

	provider, err := cloudprovider.GetProvider(account.ProviderType, config)
	if err != nil {
		return false, "failed to initialize provider: " + err.Error(), err
	}

	// 尝试实际调用云厂商 API 来验证凭证有效性
	// 不同云厂商使用不同的验证方式
	switch account.ProviderType {
	case "alibaba":
		// 尝试列出区域来验证
		regions, err := provider.ListRegions()
		if err != nil {
			return false, "failed to list regions: " + err.Error(), err
		}
		if len(regions) == 0 {
			return false, "no regions returned", nil
		}
		return true, "credentials verified, " + strconv.Itoa(len(regions)) + " regions available", nil

	case "tencent", "aws", "azure":
		// 尝试列出镜像来验证
		images, err := provider.ListImages(ctx, cloudprovider.ImageFilter{})
		if err != nil {
			return false, "failed to list images: " + err.Error(), err
		}
		return true, "credentials verified, " + strconv.Itoa(len(images)) + " images available", nil

	default:
		// 对于其他云厂商，使用基本验证
		cloudInfo := provider.GetCloudInfo()
		if cloudInfo.Provider == "" {
			return false, "provider info not available", nil
		}
		return true, "provider initialized", nil
	}
}

// TestConnectionWithCredentials 使用新凭证测试连接
func (s *CloudAccountService) TestConnectionWithCredentials(ctx context.Context, account *model.CloudAccount, accessKeyId, accessKeySecret string) (bool, string, []string, error) {
	// 构建新的凭证配置
	creds := map[string]string{
		"access_key_id":     accessKeyId,
		"access_key_secret": accessKeySecret,
		"region_id":         "cn-hangzhou", // 使用默认区域测试
	}

	config := cloudprovider.CloudAccountConfig{
		ID:           strconv.Itoa(int(account.ID)),
		Name:         account.Name,
		ProviderType: account.ProviderType,
		Credentials:  creds,
		Region:       "cn-hangzhou",
	}

	provider, err := cloudprovider.GetProvider(account.ProviderType, config)
	if err != nil {
		return false, "failed to initialize provider: " + err.Error(), []string{}, err
	}

	// 尝试列出区域来验证凭证有效性
	regions, err := provider.ListRegions()
	if err != nil {
		return false, "failed to list regions: " + err.Error(), []string{}, err
	}

	// 将区域转换为字符串数组
	regionNames := []string{}
	for _, region := range regions {
		regionNames = append(regionNames, region.Name)
	}

	if len(regions) == 0 {
		return false, "no regions returned", regionNames, nil
	}

	return true, "连接成功，" + strconv.Itoa(len(regions)) + " 个区域可用", regionNames, nil
}
