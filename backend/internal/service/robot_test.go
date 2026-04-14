package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"

	"github.com/opencmp/opencmp/internal/model"
	"github.com/opencmp/opencmp/testutils"
)

func TestRobotService(t *testing.T) {
	db := testutils.SetupTestDB()
	defer testutils.TeardownTestDB(db)

	service := NewRobotService(db)
	ctx := context.Background()

	// Create a test robot
	robot := &model.Robot{
		Name:         "Test Robot",
		Type:         "dingtalk",
		WebhookURL:   "https://oapi.dingtalk.com/robot/send?access_token=test",
		Description:  "Test dingtalk robot",
		Enabled:      true,
		MessageTypes: datatypes.JSON(`["alert", "notification"]`),
	}

	// Test CreateRobot
	err := service.CreateRobot(ctx, robot)
	assert.NoError(t, err)
	assert.NotZero(t, robot.ID)

	// Test GetRobot
	retrievedRobot, err := service.GetRobot(ctx, robot.ID)
	assert.NoError(t, err)
	assert.Equal(t, robot.Name, retrievedRobot.Name)
	assert.Equal(t, robot.Type, retrievedRobot.Type)

	// Test GetRobotByName
	robotByName, err := service.GetRobotByName(ctx, "Test Robot")
	assert.NoError(t, err)
	assert.NotNil(t, robotByName)
	assert.Equal(t, robot.ID, robotByName.ID)

	// Test ListRobots
	robots, total, err := service.ListRobots(ctx, "dingtalk", 10, 0)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, robots, 1)
	assert.Equal(t, robot.ID, robots[0].ID)

	// Test UpdateRobot
	robot.Description = "Updated description"
	err = service.UpdateRobot(ctx, robot)
	assert.NoError(t, err)

	updatedRobot, err := service.GetRobot(ctx, robot.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated description", updatedRobot.Description)

	// Test DisableRobot
	err = service.DisableRobot(ctx, robot.ID)
	assert.NoError(t, err)

	disabledRobot, err := service.GetRobot(ctx, robot.ID)
	assert.NoError(t, err)
	assert.False(t, disabledRobot.Enabled)

	// Test EnableRobot
	err = service.EnableRobot(ctx, robot.ID)
	assert.NoError(t, err)

	enabledRobot, err := service.GetRobot(ctx, robot.ID)
	assert.NoError(t, err)
	assert.True(t, enabledRobot.Enabled)

	// Test GetMessageTypes
	types, err := service.GetMessageTypes(ctx, robot.ID)
	assert.NoError(t, err)
	assert.Equal(t, []string{"alert", "notification"}, types)

	// Test SetMessageTypes
	newTypes := []string{"alert", "notification", "info"}
	err = service.SetMessageTypes(ctx, robot.ID, newTypes)
	assert.NoError(t, err)

	updatedTypes, err := service.GetMessageTypes(ctx, robot.ID)
	assert.NoError(t, err)
	assert.Equal(t, newTypes, updatedTypes)

	// Test DeleteRobot
	err = service.DeleteRobot(ctx, robot.ID)
	assert.NoError(t, err)

	// Verify deletion
	deletedRobot, err := service.GetRobot(ctx, robot.ID)
	assert.NoError(t, err)
	assert.Nil(t, deletedRobot)
}
