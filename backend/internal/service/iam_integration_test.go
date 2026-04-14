package service

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
)

func setupIAMTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to test database")
	}

	err = db.AutoMigrate(
		&model.Domain{},
		&model.Project{},
		&model.User{},
		&model.Group{},
		&model.Role{},
		&model.Permission{},
		&model.UserRole{},
		&model.RolePermission{},
		&model.GroupRole{},
		&model.Policy{},
		&model.PolicyStatement{},
		&model.RolePolicy{},
	)
	if err != nil {
		panic("failed to migrate test database")
	}

	return db
}

func TestIAMModuleIntegration(t *testing.T) {
	db := setupIAMTestDB()

	// 初始化服务
	domainService := NewDomainService(db)
	userService := NewUserService(db)
	roleService := NewRoleService(db)
	permissionService := NewPermissionService(db)
	policyService := NewPolicyService(db)

	t.Run("Complete IAM flow test", func(t *testing.T) {
		ctx := context.Background()

		// 1. 创建域
		domain := &model.Domain{
			Name:        "test-domain",
			Description: "A test domain",
		}
		err := domainService.CreateDomain(ctx, domain)
		assert.NoError(t, err)
		assert.NotZero(t, domain.ID)

		// 2. 创建用户
		user := &model.User{
			Name:     "test-user",
			Email:    "test@example.com",
			DomainID: domain.ID,
			Password: "password123",
		}
		err = userService.CreateUser(ctx, user)
		assert.NoError(t, err)
		assert.NotZero(t, user.ID)

		// 3. 创建角色
		role := &model.Role{
			Name:        "test-role",
			Description: "A test role",
			DomainID:    domain.ID,
		}
		err = roleService.CreateRole(ctx, role)
		assert.NoError(t, err)
		assert.NotZero(t, role.ID)

		// 4. 创建权限
		permission := &model.Permission{
			Name:        "test-permission",
			Resource:    "vm",
			Action:      "create",
			Scope:       "domain",
			Description: "Permission to create VMs",
			DomainID:    &domain.ID,
		}
		err = permissionService.CreatePermission(ctx, permission)
		assert.NoError(t, err)
		assert.NotZero(t, permission.ID)

		// 5. 将权限分配给角色
		err = roleService.AssignPermissionToRole(ctx, role.ID, permission.ID)
		assert.NoError(t, err)

		// 6. 将角色分配给用户
		err = userService.AssignUserRole(ctx, user.ID, role.ID, domain.ID)
		assert.NoError(t, err)

		// 7. 验证用户权限
		userPermissions, err := permissionService.GetUserPermissions(ctx, user.ID)
		assert.NoError(t, err)
		assert.Len(t, userPermissions, 1)
		assert.Equal(t, permission.ID, userPermissions[0].ID)

		// 8. 创建策略
		policyData := map[string]interface{}{
			"Version": "2012-10-17",
			"Statement": []map[string]interface{}{
				{
					"Effect":   "Allow",
					"Resource": "vm:*",
					"Action":   []string{"*"},
				},
			},
		}
		policyJSON, _ := json.Marshal(policyData)
		policy := &model.Policy{
			Name:        "test-policy",
			Description: "A test policy",
			Scope:       "domain",
			Policy:      policyJSON,
		}
		err = policyService.CreatePolicy(ctx, policy)
		assert.NoError(t, err)
		assert.NotZero(t, policy.ID)

		// 9. 将策略分配给角色
		err = policyService.AssignPolicyToRole(ctx, role.ID, policy.ID)
		assert.NoError(t, err)

		// 10. 验证用户拥有角色
		userRoles, err := roleService.GetUserRoles(ctx, user.ID, domain.ID)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(userRoles), 1)
	})
}
