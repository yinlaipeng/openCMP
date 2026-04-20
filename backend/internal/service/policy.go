package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
)

// PolicyService 策略服务
type PolicyService struct {
	db *gorm.DB
}

// NewPolicyService 创建策略服务
func NewPolicyService(db *gorm.DB) *PolicyService {
	return &PolicyService{db: db}
}

// ListPolicies 列出策略
func (s *PolicyService) ListPolicies(ctx context.Context, scope string, domainID *string, isSystem *bool, limit, offset int) ([]model.Policy, int64, error) {
	var policies []model.Policy
	var total int64

	query := s.db.WithContext(ctx).Model(&model.Policy{})

	// 按作用域筛选
	if scope != "" {
		query = query.Where("scope = ?", scope)
	}

	// 按域 ID 筛选
	if domainID != nil && *domainID != "" {
		query = query.Where("domain_id = ? OR scope = 'system'", *domainID)
	}

	// 按系统策略筛选
	if isSystem != nil {
		query = query.Where("is_system = ?", *isSystem)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&policies).Error; err != nil {
		return nil, 0, err
	}

	// 计算每个策略的是否可删除/更新
	for i := range policies {
		s.calculatePolicyPermissions(&policies[i])
	}

	return policies, total, nil
}

// GetPolicy 获取策略详情
func (s *PolicyService) GetPolicy(ctx context.Context, id string) (*model.Policy, error) {
	var policy model.Policy
	if err := s.db.WithContext(ctx).First(&policy, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	s.calculatePolicyPermissions(&policy)
	return &policy, nil
}

// CreatePolicy 创建策略
func (s *PolicyService) CreatePolicy(ctx context.Context, policy *model.Policy) error {
	// 生成策略 ID
	if policy.ID == "" {
		policy.ID = uuid.New().String()
	}

	// 设置默认值
	policy.Enabled = true
	policy.IsSystem = false
	policy.IsEmulated = false
	policy.IsPublic = true
	policy.PublicScope = "system"
	policy.UpdateVersion = 0

	return s.db.WithContext(ctx).Create(policy).Error
}

// UpdatePolicy 更新策略
func (s *PolicyService) UpdatePolicy(ctx context.Context, id string, updates map[string]interface{}) error {
	// 系统策略不可更新
	var policy model.Policy
	if err := s.db.WithContext(ctx).First(&policy, "id = ?", id).Error; err != nil {
		return err
	}

	if policy.IsSystem {
		return errors.New("不可更新系统策略")
	}

	updates["updated_at"] = time.Now()
	updates["update_version"] = policy.UpdateVersion + 1

	return s.db.WithContext(ctx).Model(&policy).Updates(updates).Error
}

// DeletePolicy 删除策略
func (s *PolicyService) DeletePolicy(ctx context.Context, id string) error {
	var policy model.Policy
	if err := s.db.WithContext(ctx).First(&policy, "id = ?", id).Error; err != nil {
		return err
	}

	// 系统策略不可删除
	if policy.IsSystem {
		return errors.New("不可删除系统策略")
	}

	// 检查是否有关联的角色
	var count int64
	if err := s.db.Model(&model.RolePolicy{}).Where("policy_id = ?", id).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return fmt.Errorf("策略已关联 %d 个角色，无法删除", count)
	}

	return s.db.WithContext(ctx).Delete(&policy).Error
}

// calculatePolicyPermissions 计算策略的权限（是否可删除/更新）
func (s *PolicyService) calculatePolicyPermissions(policy *model.Policy) {
	// 系统策略不可删除和更新
	if policy.IsSystem {
		policy.CanDelete = false
		policy.CanUpdate = false
		policy.DeleteFailReason = datatypes.JSON(`{"class":"ForbiddenError","code":403,"details":"不可删除系统策略定义"}`)
		return
	}

	// 检查是否有关联的角色
	var count int64
	s.db.Model(&model.RolePolicy{}).Where("policy_id = ?", policy.ID).Count(&count)

	if count > 0 {
		policy.CanDelete = false
		policy.DeleteFailReason = datatypes.JSON(fmt.Sprintf(`{"class":"NotEmptyError","code":406,"details":"policy is in associated with %d roles"}`, count))
	} else {
		policy.CanDelete = true
	}

	// 非系统策略可以更新
	policy.CanUpdate = true
}

// GetPolicyByScope 按作用域获取策略列表
func (s *PolicyService) GetPolicyByScope(ctx context.Context, scope string) ([]model.Policy, error) {
	var policies []model.Policy
	query := s.db.WithContext(ctx).Where("scope = ? AND enabled = ? AND deleted = ?", scope, true, false)

	if err := query.Order("name ASC").Find(&policies).Error; err != nil {
		return nil, err
	}

	return policies, nil
}

// BatchCreatePolicies 批量创建策略（用于初始化）
func (s *PolicyService) BatchCreatePolicies(ctx context.Context, policies []model.Policy) error {
	return s.db.WithContext(ctx).CreateInBatches(policies, 50).Error
}

// CheckPolicyExists 检查策略是否存在
func (s *PolicyService) CheckPolicyExists(ctx context.Context, name string) (bool, error) {
	var count int64
	if err := s.db.WithContext(ctx).Model(&model.Policy{}).Where("name = ?", name).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// AssignPolicyToRole 分配策略给角色
func (s *PolicyService) AssignPolicyToRole(ctx context.Context, roleID uint, policyID string) error {
	// 检查是否已关联
	var count int64
	if err := s.db.Model(&model.RolePolicy{}).Where("role_id = ? AND policy_id = ?", roleID, policyID).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return errors.New("策略已关联到该角色")
	}

	rolePolicy := &model.RolePolicy{
		RoleID:   roleID,
		PolicyID: policyID,
	}

	return s.db.Create(rolePolicy).Error
}

// GetRolePolicies 获取角色的策略列表
func (s *PolicyService) GetRolePolicies(ctx context.Context, roleID uint) ([]model.Policy, error) {
	var policies []model.Policy
	err := s.db.WithContext(ctx).
		Joins("JOIN role_policies ON policies.id = role_policies.policy_id").
		Where("role_policies.role_id = ?", roleID).
		Find(&policies).Error

	if err != nil {
		return nil, err
	}

	return policies, nil
}

// RevokePolicyFromRole 从角色撤销策略
func (s *PolicyService) RevokePolicyFromRole(ctx context.Context, roleID uint, policyID string) error {
	return s.db.WithContext(ctx).
		Where("role_id = ? AND policy_id = ?", roleID, policyID).
		Delete(&model.RolePolicy{}).Error
}

// CheckUserPermission 检查用户是否有指定权限（通过策略）
func (s *PolicyService) CheckUserPermission(ctx context.Context, userID uint, resource, action string) (bool, error) {
	// 获取用户的所有策略（通过角色）
	var policyIDs []string

	// 从用户直接拥有的角色获取策略
	err := s.db.WithContext(ctx).
		Model(&model.RolePolicy{}).
		Joins("JOIN user_roles ON role_policies.role_id = user_roles.role_id").
		Where("user_roles.user_id = ?", userID).
		Pluck("role_policies.policy_id", &policyIDs).Error

	if err != nil {
		return false, err
	}

	// 从用户在项目级别拥有的角色获取策略
	var projectPolicyIDs []string
	err = s.db.WithContext(ctx).
		Model(&model.RolePolicy{}).
		Joins("JOIN project_user_roles ON role_policies.role_id = project_user_roles.role_id").
		Where("project_user_roles.user_id = ?", userID).
		Pluck("role_policies.policy_id", &projectPolicyIDs).Error

	if err != nil {
		return false, err
	}

	policyIDs = append(policyIDs, projectPolicyIDs...)

	// 从用户所属组拥有的角色获取策略
	var groupPolicyIDs []string
	err = s.db.WithContext(ctx).
		Model(&model.RolePolicy{}).
		Joins("JOIN group_roles ON role_policies.role_id = group_roles.role_id").
		Joins("JOIN user_groups ON group_roles.group_id = user_groups.group_id").
		Where("user_groups.user_id = ?", userID).
		Pluck("role_policies.policy_id", &groupPolicyIDs).Error

	if err != nil {
		return false, err
	}

	policyIDs = append(policyIDs, groupPolicyIDs...)

	if len(policyIDs) == 0 {
		return false, nil
	}

	// 检查这些策略是否允许指定的资源和动作
	var policies []model.Policy
	err = s.db.WithContext(ctx).
		Where("id IN ?", policyIDs).
		Find(&policies).Error

	if err != nil {
		return false, err
	}

	// 检查每个策略是否允许请求的操作
	for _, policy := range policies {
		if s.evaluatePolicy(policy, resource, action) {
			return true, nil
		}
	}

	return false, nil
}

// evaluatePolicy 评估策略是否允许特定资源和动作
func (s *PolicyService) evaluatePolicy(policy model.Policy, resource, action string) bool {
	var policyMap map[string]interface{}
	if err := json.Unmarshal(policy.Policy, &policyMap); err != nil {
		// 如果无法解析策略，则视为不允许
		return false
	}

	// 获取策略语句
	statements, ok := policyMap["statement"].([]interface{})
	if !ok {
		// 尝试其他可能的键名
		if stmt, exists := policyMap["statements"]; exists {
			statements, ok = stmt.([]interface{})
		}
		if !ok {
			return false
		}
	}

	// 遍历每个策略语句
	for _, stmt := range statements {
		stmtMap, ok := stmt.(map[string]interface{})
		if !ok {
			continue
		}

		// 检查效果（允许还是拒绝）
		effect, ok := stmtMap["effect"].(string)
		if !ok || effect != "allow" {
			continue
		}

		// 检查资源
		resources, hasResources := stmtMap["resource"]
		if hasResources {
			resourceList, ok := resources.([]interface{})
			if !ok {
				// 如果不是数组，尝试字符串
				resourceStr, ok := resources.(string)
				if !ok || !matchResource(resourceStr, resource) {
					continue
				}
			} else {
				matched := false
				for _, res := range resourceList {
					resStr, ok := res.(string)
					if ok && matchResource(resStr, resource) {
						matched = true
						break
					}
				}
				if !matched {
					continue
				}
			}
		}

		// 检查动作
		actions, hasActions := stmtMap["action"]
		if hasActions {
			actionList, ok := actions.([]interface{})
			if !ok {
				// 如果不是数组，尝试字符串
				actionStr, ok := actions.(string)
				if !ok || !matchAction(actionStr, action) {
					continue
				}
			} else {
				matched := false
				for _, act := range actionList {
					actStr, ok := act.(string)
					if ok && matchAction(actStr, action) {
						matched = true
						break
					}
				}
				if !matched {
					continue
				}
			}
		}

		// 如果资源和动作都匹配，则允许
		return true
	}

	return false
}

// matchResource 检查资源是否匹配
func matchResource(policyResource, requestedResource string) bool {
	// 支持通配符匹配
	if policyResource == "*" || policyResource == requestedResource {
		return true
	}

	// 检查通配符模式，如 "user:*"
	if len(policyResource) > 2 && policyResource[len(policyResource)-2:] == ":*" {
		prefix := policyResource[:len(policyResource)-2]
		if requestedResource == prefix || len(requestedResource) > len(prefix) &&
			requestedResource[:len(prefix)] == prefix && requestedResource[len(prefix)] == ':' {
			return true
		}
	}

	return false
}

// matchAction 检查动作是否匹配
func matchAction(policyAction, requestedAction string) bool {
	// 支持通配符匹配
	if policyAction == "*" || policyAction == requestedAction {
		return true
	}

	// 检查通配符模式，如 "user:*"
	if len(policyAction) > 2 && policyAction[len(policyAction)-2:] == ":*" {
		prefix := policyAction[:len(policyAction)-2]
		if requestedAction == prefix || len(requestedAction) > len(prefix) &&
			requestedAction[:len(prefix)] == prefix && requestedAction[len(prefix)] == ':' {
			return true
		}
	}

	return false
}

// SetPolicyEnabled 设置策略启用状态
func (s *PolicyService) SetPolicyEnabled(ctx context.Context, id string, enabled bool) error {
	return s.db.WithContext(ctx).Model(&model.Policy{}).Where("id = ?", id).Update("enabled", enabled).Error
}

// GetPolicyRoles 获取策略关联的角色列表
func (s *PolicyService) GetPolicyRoles(ctx context.Context, policyID string) ([]model.Role, error) {
	var roles []model.Role
	err := s.db.WithContext(ctx).
		Joins("JOIN role_policies ON roles.id = role_policies.role_id").
		Where("role_policies.policy_id = ?", policyID).
		Find(&roles).Error

	if err != nil {
		return nil, err
	}

	return roles, nil
}
