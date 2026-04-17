package service

import (
	"context"
	"regexp"
	"strconv"
	"strings"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
	"github.com/opencmp/opencmp/pkg/cloudprovider"
)

// ResourceMappingService 资源项目归属映射服务
type ResourceMappingService struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewResourceMappingService 创建资源映射服务
func NewResourceMappingService(db *gorm.DB, logger *zap.Logger) *ResourceMappingService {
	return &ResourceMappingService{
		db:     db,
		logger: logger,
	}
}

// ProjectAttributionResult 项目归属结果
type ProjectAttributionResult struct {
	ProjectID   uint   // 确定的项目ID
	ProjectName string // 项目名称
	MatchedRule string // 匹配的规则描述
	IsDefault   bool   // 是否为默认项目
}

// DetermineProjectAttribution 根据资源标签和同步策略确定项目归属
// 业务逻辑流程：
// 1. 获取云账号绑定的同步策略
// 2. 解析资源标签
// 3. 按优先级遍历规则的标签映射规则
// 4. 返回最高优先级匹配结果
func (s *ResourceMappingService) DetermineProjectAttribution(ctx context.Context, cloudAccountID uint, resourceTags map[string]string) (*ProjectAttributionResult, error) {
	// 1. 获取云账号信息
	var cloudAccount model.CloudAccount
	if err := s.db.WithContext(ctx).First(&cloudAccount, cloudAccountID).Error; err != nil {
		s.logger.Warn("获取云账号失败", zap.Uint("cloud_account_id", cloudAccountID), zap.Error(err))
		return s.getDefaultProject(ctx, cloudAccount.DomainID), nil
	}

	// 2. 检查云账号的资源分配方式
	if cloudAccount.ResourceAssignmentMethod == "manual_assignment" {
		// 手动分配模式，不使用标签映射
		return s.getDefaultProject(ctx, cloudAccount.DomainID), nil
	}

	// 3. 获取云账号绑定的同步策略（通过定时任务关联）
	var syncPolicy model.SyncPolicy
	err := s.db.WithContext(ctx).
		Joins("JOIN scheduled_tasks ON scheduled_tasks.sync_policy_id = sync_policies.id").
		Where("scheduled_tasks.cloud_account_id = ?", cloudAccountID).
		Where("scheduled_tasks.enabled = ?", true).
		Where("sync_policies.enabled = ?", true).
		Where("sync_policies.domain_id = ?", cloudAccount.DomainID).
		First(&syncPolicy).Error

	if err != nil {
		// 未找到同步策略，使用默认项目
		s.logger.Debug("未找到绑定同步策略", zap.Uint("cloud_account_id", cloudAccountID))
		return s.getDefaultProject(ctx, cloudAccount.DomainID), nil
	}

	// 4. 加载同步策略的规则（按优先级排序）
	var rules []model.Rule
	err = s.db.WithContext(ctx).
		Preload("Tags").
		Where("sync_policy_id = ?", syncPolicy.ID).
		Order("id ASC"). // 规则按创建顺序为优先级
		Find(&rules).Error

	if err != nil || len(rules) == 0 {
		return s.getDefaultProject(ctx, cloudAccount.DomainID), nil
	}

	// 5. 遍历规则进行标签匹配
	for _, rule := range rules {
		if s.matchRule(rule, resourceTags) {
			// 匹配成功，返回目标项目
			result := &ProjectAttributionResult{
				IsDefault:   false,
				MatchedRule: "Rule #" + strconv.FormatUint(uint64(rule.ID), 10),
			}

			if rule.ResourceMapping == "specify_project" && rule.TargetProjectID != nil {
				result.ProjectID = *rule.TargetProjectID
				result.ProjectName = rule.TargetProjectName
			} else if rule.ResourceMapping == "specify_name" {
				// 根据标签值查找项目名称
				projectName := s.extractProjectNameFromTags(resourceTags, rule.Tags)
				if projectName != "" {
					var project model.Project
					if err := s.db.WithContext(ctx).
						Where("name = ?", projectName).
						Where("domain_id = ?", cloudAccount.DomainID).
						First(&project).Error; err == nil {
						result.ProjectID = project.ID
						result.ProjectName = project.Name
					}
				}
			}

			if result.ProjectID > 0 {
				s.logger.Debug("标签匹配成功",
					zap.Uint("cloud_account_id", cloudAccountID),
					zap.Uint("rule_id", rule.ID),
					zap.Uint("project_id", result.ProjectID),
					zap.String("matched_rule", result.MatchedRule))
				return result, nil
			}
		}
	}

	// 6. 未匹配任何规则，使用默认项目
	return s.getDefaultProject(ctx, cloudAccount.DomainID), nil
}

// matchRule 检查资源标签是否匹配规则
func (s *ResourceMappingService) matchRule(rule model.Rule, resourceTags map[string]string) bool {
	if len(rule.Tags) == 0 {
		return false
	}

	switch rule.ConditionType {
	case "all_match":
		// 所有标签都必须匹配
		for _, ruleTag := range rule.Tags {
			if !s.matchTag(ruleTag, resourceTags) {
				return false
			}
		}
		return true

	case "any_match":
		// 任一标签匹配即可
		for _, ruleTag := range rule.Tags {
			if s.matchTag(ruleTag, resourceTags) {
				return true
			}
		}
		return false

	case "key_match":
		// 仅检查键是否存在
		for _, ruleTag := range rule.Tags {
			if _, exists := resourceTags[ruleTag.TagKey]; exists {
				return true
			}
		}
		return false

	default:
		return false
	}
}

// matchTag 检查单个标签是否匹配
// 支持精确匹配和正则表达式匹配
func (s *ResourceMappingService) matchTag(ruleTag model.RuleTag, resourceTags map[string]string) bool {
	value, exists := resourceTags[ruleTag.TagKey]
	if !exists {
		return false
	}

	// 检查是否为正则表达式（包含特殊字符）
	if strings.ContainsAny(ruleTag.TagValue, ".*+?[]()|^$") {
		// 正则匹配
		matched, err := regexp.MatchString(ruleTag.TagValue, value)
		if err != nil {
			s.logger.Warn("正则匹配失败", zap.String("pattern", ruleTag.TagValue), zap.Error(err))
			return false
		}
		return matched
	}

	// 精确匹配
	return value == ruleTag.TagValue
}

// extractProjectNameFromTags 从标签中提取项目名称
func (s *ResourceMappingService) extractProjectNameFromTags(resourceTags map[string]string, ruleTags []model.RuleTag) string {
	for _, rt := range ruleTags {
		if value, exists := resourceTags[rt.TagKey]; exists {
			return value
		}
	}
	return ""
}

// getDefaultProject 获取域的默认项目
func (s *ResourceMappingService) getDefaultProject(ctx context.Context, domainID uint) *ProjectAttributionResult {
	// 查找域的默认项目（名称包含"默认"或"Default"的项目）
	var project model.Project
	err := s.db.WithContext(ctx).
		Where("domain_id = ?", domainID).
		Where("enabled = ?", true).
		Where("name LIKE ?", "%默认%").
		First(&project).Error

	if err != nil {
		// 没找到，查找名称包含Default的项目
		err = s.db.WithContext(ctx).
			Where("domain_id = ?", domainID).
			Where("enabled = ?", true).
			Where("name LIKE ?", "%Default%").
			First(&project).Error
	}

	if err != nil {
		// 没有默认项目，查找域的第一个可用项目
		err = s.db.WithContext(ctx).
			Where("domain_id = ?", domainID).
			Where("enabled = ?", true).
			First(&project).Error
	}

	if err != nil {
		s.logger.Warn("未找到域的可用项目", zap.Uint("domain_id", domainID))
		return &ProjectAttributionResult{
			ProjectID:   0,
			ProjectName: "",
			IsDefault:   true,
			MatchedRule: "无项目",
		}
	}

	return &ProjectAttributionResult{
		ProjectID:   project.ID,
		ProjectName: project.Name,
		IsDefault:   true,
		MatchedRule: "默认项目",
	}
}

// ParseResourceTags 从云资源中提取标签（统一接口）
func (s *ResourceMappingService) ParseResourceTags(resource interface{}) map[string]string {
	switch r := resource.(type) {
	case *cloudprovider.VirtualMachine:
		return r.Tags
	case *cloudprovider.VPC:
		return r.Tags
	case *cloudprovider.Subnet:
		return r.Tags
	case *cloudprovider.SecurityGroup:
		return r.Tags
	case *cloudprovider.Disk:
		return r.Tags
	case *cloudprovider.LoadBalancer:
		return r.Tags
	case *cloudprovider.HostTemplate:
		return r.Tags
	case *cloudprovider.AutoscalingGroup:
		return r.Tags
	default:
		return map[string]string{}
	}
}

// BatchDetermineProjectAttribution 批量确定资源项目归属
// 用于同步时批量处理资源，减少数据库查询次数
func (s *ResourceMappingService) BatchDetermineProjectAttribution(ctx context.Context, cloudAccountID uint, resources []interface{}) map[string]*ProjectAttributionResult {
	results := make(map[string]*ProjectAttributionResult)

	for _, resource := range resources {
		resourceID := s.getResourceID(resource)
		if resourceID == "" {
			continue
		}

		tags := s.ParseResourceTags(resource)
		attribution, err := s.DetermineProjectAttribution(ctx, cloudAccountID, tags)
		if err != nil {
			s.logger.Warn("确定项目归属失败", zap.String("resource_id", resourceID), zap.Error(err))
			continue
		}

		results[resourceID] = attribution
	}

	return results
}

// getResourceID 从资源中提取ID
func (s *ResourceMappingService) getResourceID(resource interface{}) string {
	switch r := resource.(type) {
	case *cloudprovider.VirtualMachine:
		return r.ID
	case *cloudprovider.VPC:
		return r.ID
	case *cloudprovider.Subnet:
		return r.ID
	case *cloudprovider.SecurityGroup:
		return r.ID
	case *cloudprovider.Disk:
		return r.ID
	case *cloudprovider.LoadBalancer:
		return r.ID
	case *cloudprovider.HostTemplate:
		return r.ID
	case *cloudprovider.AutoscalingGroup:
		return r.ID
	default:
		return ""
	}
}