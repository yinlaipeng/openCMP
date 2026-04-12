package service

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
)

type PermissionServiceTestSuite struct {
	suite.Suite
	service *PermissionService
	db      *gorm.DB
	mock    sqlmock.Sqlmock
}

func (suite *PermissionServiceTestSuite) SetupTest() {
	db, mock, err := NewMockDB()
	assert.NoError(suite.T(), err)
	suite.db = db
	suite.mock = mock
	suite.service = NewPermissionService(db)
}

func (suite *PermissionServiceTestSuite) TestCreatePermission() {
	permission := &model.Permission{
		Name:        "user.list",
		Description: "List users permission",
		Resource:    "user",
		Action:      "list",
		Scope:       "system",
		Enabled:     true,
	}

	suite.mock.ExpectBegin()
	suite.mock.ExpectExec("INSERT INTO `permissions`").
		WithArgs(permission.Name, permission.Description, permission.Resource, permission.Action, permission.Scope, permission.Enabled, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectCommit()

	err := suite.service.CreatePermission(context.Background(), permission)

	assert.NoError(suite.T(), err)
	assert.NotZero(suite.T(), permission.ID)
}

func (suite *PermissionServiceTestSuite) TestGetPermission() {
	expectedPermission := &model.Permission{
		ID:          1,
		Name:        "user.list",
		Description: "List users permission",
		Resource:    "user",
		Action:      "list",
		Scope:       "system",
		Enabled:     true,
	}

	rows := sqlmock.NewRows([]string{"id", "name", "description", "resource", "action", "scope", "domain_id", "enabled", "created_at", "updated_at"}).
		AddRow(expectedPermission.ID, expectedPermission.Name, expectedPermission.Description, expectedPermission.Resource, expectedPermission.Action, expectedPermission.Scope, expectedPermission.DomainID, expectedPermission.Enabled, expectedPermission.CreatedAt, expectedPermission.UpdatedAt)

	suite.mock.ExpectQuery("^SELECT (.+) FROM `permissions`").
		WithArgs(1).
		WillReturnRows(rows)

	permission, err := suite.service.GetPermission(context.Background(), 1)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedPermission.ID, permission.ID)
	assert.Equal(suite.T(), expectedPermission.Name, permission.Name)
}

func (suite *PermissionServiceTestSuite) TestGetPermission_NotFound() {
	suite.mock.ExpectQuery("^SELECT (.+) FROM `permissions`").
		WithArgs(999).
		WillReturnError(gorm.ErrRecordNotFound)

	permission, err := suite.service.GetPermission(context.Background(), 999)

	assert.NoError(suite.T(), err)
	assert.Nil(suite.T(), permission)
}

func (suite *PermissionServiceTestSuite) TestListPermissions() {
	expectedPermissions := []*model.Permission{
		{
			ID:          1,
			Name:        "user.list",
			Description: "List users permission",
			Resource:    "user",
			Action:      "list",
			Scope:       "system",
			Enabled:     true,
		},
		{
			ID:          2,
			Name:        "user.create",
			Description: "Create users permission",
			Resource:    "user",
			Action:      "create",
			Scope:       "system",
			Enabled:     true,
		},
	}

	countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
	suite.mock.ExpectQuery("^SELECT count(.+) FROM `permissions`").
		WillReturnRows(countRows)

	rows := sqlmock.NewRows([]string{"id", "name", "description", "resource", "action", "scope", "domain_id", "enabled", "created_at", "updated_at"}).
		AddRow(1, "user.list", "List users permission", "user", "list", "system", nil, true, "2023-01-01 00:00:00", "2023-01-01 00:00:00").
		AddRow(2, "user.create", "Create users permission", "user", "create", "system", nil, true, "2023-01-01 00:00:00", "2023-01-01 00:00:00")

	suite.mock.ExpectQuery("^SELECT (.+) FROM `permissions`").
		WillReturnRows(rows)

	permissions, total, err := suite.service.ListPermissions(context.Background(), nil, "", "", "", nil, 10, 0)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(2), total)
	assert.Len(suite.T(), permissions, 2)
}

func (suite *PermissionServiceTestSuite) TestUpdatePermission() {
	permission := &model.Permission{
		ID:          1,
		Name:        "user.list.updated",
		Description: "Updated description",
		Resource:    "user",
		Action:      "list",
		Scope:       "system",
		Enabled:     false,
	}

	suite.mock.ExpectBegin()
	suite.mock.ExpectExec("UPDATE `permissions` SET").
		WithArgs(permission.Name, permission.Description, permission.Resource, permission.Action, permission.Scope, permission.DomainID, permission.Enabled, sqlmock.AnyArg(), permission.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectCommit()

	err := suite.service.UpdatePermission(context.Background(), permission)

	assert.NoError(suite.T(), err)
}

func (suite *PermissionServiceTestSuite) TestDeletePermission() {
	suite.mock.ExpectBegin()
	suite.mock.ExpectExec("DELETE FROM `permissions`").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectCommit()

	err := suite.service.DeletePermission(context.Background(), 1)

	assert.NoError(suite.T(), err)
}

func (suite *PermissionServiceTestSuite) TestAssignPermissionToRole() {
	suite.mock.ExpectQuery("^SELECT count(.+) FROM `role_permissions`").
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	suite.mock.ExpectBegin()
	suite.mock.ExpectExec("INSERT INTO `role_permissions`").
		WithArgs(1, 1, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectCommit()

	err := suite.service.AssignPermissionToRole(context.Background(), 1, 1)

	assert.NoError(suite.T(), err)
}

func (suite *PermissionServiceTestSuite) TestAssignPermissionToRole_Duplicate() {
	suite.mock.ExpectQuery("^SELECT count(.+) FROM `role_permissions`").
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	err := suite.service.AssignPermissionToRole(context.Background(), 1, 1)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), gorm.ErrDuplicatedKey, err)
}

func (suite *PermissionServiceTestSuite) TestRevokePermissionFromRole() {
	suite.mock.ExpectBegin()
	suite.mock.ExpectExec("DELETE FROM `role_permissions`").
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectCommit()

	err := suite.service.RevokePermissionFromRole(context.Background(), 1, 1)

	assert.NoError(suite.T(), err)
}

func (suite *PermissionServiceTestSuite) TestGetPermissionsByRole() {
	rows := sqlmock.NewRows([]string{"id", "name", "description", "resource", "action", "scope", "domain_id", "enabled", "created_at", "updated_at"}).
		AddRow(1, "user.list", "List users permission", "user", "list", "system", nil, true, "2023-01-01 00:00:00", "2023-01-01 00:00:00").
		AddRow(2, "user.create", "Create users permission", "user", "create", "system", nil, true, "2023-01-01 00:00:00", "2023-01-01 00:00:00")

	suite.mock.ExpectQuery("^SELECT (.+) FROM `permissions`").
		WithArgs(1).
		WillReturnRows(rows)

	permissions, err := suite.service.GetPermissionsByRole(context.Background(), 1)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), permissions, 2)
}

func TestPermissionServiceTestSuite(t *testing.T) {
	suite.Run(t, new(PermissionServiceTestSuite))
}
