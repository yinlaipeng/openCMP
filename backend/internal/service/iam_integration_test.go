package service

import (
	"context"
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
		// 1. 创建域
		domain := &model.Domain{
			Name:        "test-domain",
			Description: "A test domain",
		}
		err := domainService.CreateDomain(domain)
		assert.NoError(t, err)
		assert.NotZero(t, domain.ID)

		// 2. 创建用户
		user := &model.User{
			Name:     "test-user",
			Email:    "test@example.com",
			Username: "testuser",
			DomainID: domain.ID,
		}
		err = userService.CreateUser(user, "password123")
		assert.NoError(t, err)
		assert.NotZero(t, user.ID)

		// 3. 创建角色
		role := &model.Role{
			Name:        "test-role",
			Description: "A test role",
			DomainID:    domain.ID,
		}
		err = roleService.CreateRole(role)
		assert.NoError(t, err)
		assert.NotZero(t, role.ID)

		// 4. 创建权限
		permission := &model.Permission{
			Name:        "test-permission",
			Resource:    "vm",
			Action:      "create",
			Description: "Permission to create VMs",
			DomainID:    domain.ID,
		}
		err = permissionService.CreatePermission(permission)
		assert.NoError(t, err)
		assert.NotZero(t, permission.ID)

		// 5. 将权限分配给角色
		err = roleService.AssignPermissionToRole(nil, role.ID, permission.ID)
		assert.NoError(t, err)

		// 6. 将角色分配给用户
		ctx := context.Background()
		err = userService.AssignUserRole(ctx, user.ID, role.ID, domain.ID)
		assert.NoError(t, err)

		// 7. 验证用户权限
		userPermissions, err := userService.GetUserPermissions(ctx, user.ID)
		assert.NoError(t, err)
		assert.Len(t, userPermissions, 1)
		assert.Equal(t, permission.ID, userPermissions[0].ID)

		// 8. 创建策略
		policy := &model.Policy{
			Name:        "test-policy",
			Description: "A test policy",
			DomainID:    domain.ID,
		}
		statements := []model.PolicyStatement{
			{
				Effect:   "Allow",
				Resource: "vm:*",
				Actions:  []string{"*"},
			},
		}
		err = policyService.CreatePolicy(ctx, policy, statements)
		assert.NoError(t, err)
		assert.NotZero(t, policy.ID)

		// 9. 将策略分配给角色
		err = policyService.AssignPolicyToRole(ctx, role.ID, policy.ID)
		assert.NoError(t, err)

		// 10. 验证策略生效
		hasPermission, err := policyService.EvaluatePolicy(ctx, user.ID, "user", "vm", "delete")
		assert.NoError(t, err)
		assert.True(t, hasPermission)
	})
}
