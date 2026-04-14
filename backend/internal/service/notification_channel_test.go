package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"

	"github.com/opencmp/opencmp/internal/model"
	"github.com/opencmp/opencmp/testutils"
)

func TestNotificationChannelService(t *testing.T) {
	db := testutils.SetupTestDB()
	defer testutils.TeardownTestDB(db)

	service := NewNotificationChannelService(db)
	ctx := context.Background()

	// Create a test notification channel
	channel := &model.NotificationChannel{
		Name:        "Test Channel",
		Type:        "email",
		Description: "Test email channel",
		Config:      datatypes.JSON(`{"smtp_host":"smtp.example.com","smtp_port":587,"from_address":"test@example.com"}`),
		Enabled:     true,
	}

	// Test CreateChannel
	err := service.CreateChannel(ctx, channel)
	assert.NoError(t, err)
	assert.NotZero(t, channel.ID)

	// Test GetChannel
	retrievedChannel, err := service.GetChannel(ctx, channel.ID)
	assert.NoError(t, err)
	assert.Equal(t, channel.Name, retrievedChannel.Name)
	assert.Equal(t, channel.Type, retrievedChannel.Type)

	// Test GetChannelByName
	channelByName, err := service.GetChannelByName(ctx, "Test Channel")
	assert.NoError(t, err)
	assert.NotNil(t, channelByName)
	assert.Equal(t, channel.ID, channelByName.ID)

	// Test ListChannels
	channels, total, err := service.ListChannels(ctx, "email", 10, 0)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, channels, 1)
	assert.Equal(t, channel.ID, channels[0].ID)

	// Test UpdateChannel
	channel.Description = "Updated description"
	err = service.UpdateChannel(ctx, channel)
	assert.NoError(t, err)

	updatedChannel, err := service.GetChannel(ctx, channel.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated description", updatedChannel.Description)

	// Test DisableChannel
	err = service.DisableChannel(ctx, channel.ID)
	assert.NoError(t, err)

	disabledChannel, err := service.GetChannel(ctx, channel.ID)
	assert.NoError(t, err)
	assert.False(t, disabledChannel.Enabled)

	// Test EnableChannel
	err = service.EnableChannel(ctx, channel.ID)
	assert.NoError(t, err)

	enabledChannel, err := service.GetChannel(ctx, channel.ID)
	assert.NoError(t, err)
	assert.True(t, enabledChannel.Enabled)

	// Test DeleteChannel
	err = service.DeleteChannel(ctx, channel.ID)
	assert.NoError(t, err)

	// Verify deletion
	deletedChannel, err := service.GetChannel(ctx, channel.ID)
	assert.NoError(t, err)
	assert.Nil(t, deletedChannel)
}
