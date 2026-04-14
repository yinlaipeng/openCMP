package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/opencmp/opencmp/internal/model"
	"github.com/opencmp/opencmp/testutils"
)

func TestMessageSubscriptionService(t *testing.T) {
	db := testutils.SetupTestDB()
	defer testutils.TeardownTestDB(db)

	service := NewMessageSubscriptionService(db)
	ctx := context.Background()

	// First, create a message type
	messageType := &model.MessageType{
		Name:        "test_notification",
		DisplayName: "Test Notification",
		Description: "Test notification type",
		Enabled:     true,
	}
	db.Create(messageType)

	// Create a test subscription
	subscription := &model.MessageSubscription{
		UserID:        1,
		MessageTypeID: messageType.ID,
		Email:         true,
		SMS:           false,
		Webhook:       true,
		Station:       true,
	}

	// Test CreateSubscription
	err := service.CreateSubscription(ctx, subscription)
	assert.NoError(t, err)
	assert.NotZero(t, subscription.ID)

	// Test GetSubscription
	retrievedSub, err := service.GetSubscription(ctx, subscription.ID)
	assert.NoError(t, err)
	assert.Equal(t, subscription.UserID, retrievedSub.UserID)
	assert.Equal(t, subscription.Email, retrievedSub.Email)

	// Test GetUserSubscription
	userSub, err := service.GetUserSubscription(ctx, 1, messageType.ID)
	assert.NoError(t, err)
	assert.NotNil(t, userSub)
	assert.Equal(t, subscription.ID, userSub.ID)

	// Test ListUserSubscriptions
	userSubs, err := service.ListUserSubscriptions(ctx, 1)
	assert.NoError(t, err)
	assert.Len(t, userSubs, 1)
	assert.Equal(t, subscription.ID, userSubs[0].ID)

	// Test UpdateSubscription
	subscription.Email = false
	subscription.SMS = true
	err = service.UpdateSubscription(ctx, subscription)
	assert.NoError(t, err)

	updatedSub, err := service.GetSubscription(ctx, subscription.ID)
	assert.NoError(t, err)
	assert.False(t, updatedSub.Email)
	assert.True(t, updatedSub.SMS)

	// Test DeleteSubscription
	err = service.DeleteSubscription(ctx, subscription.ID)
	assert.NoError(t, err)

	// Verify deletion
	deletedSub, err := service.GetSubscription(ctx, subscription.ID)
	assert.NoError(t, err)
	assert.Nil(t, deletedSub)
}

func TestMessageSubscriptionService_SetSubscriptionChannels(t *testing.T) {
	db := testutils.SetupTestDB()
	defer testutils.TeardownTestDB(db)

	service := NewMessageSubscriptionService(db)
	ctx := context.Background()

	// Create a message type
	messageType := &model.MessageType{
		Name:        "alert_notification",
		DisplayName: "Alert Notification",
		Description: "Alert notification type",
		Enabled:     true,
	}
	db.Create(messageType)

	// Test SetSubscriptionChannels for a new subscription
	channels := map[string]bool{
		"email":   true,
		"sms":     false,
		"webhook": true,
		"station": false,
	}

	err := service.SetSubscriptionChannels(ctx, 2, messageType.ID, channels)
	assert.NoError(t, err)

	// Verify the subscription was created
	sub, err := service.GetUserSubscription(ctx, 2, messageType.ID)
	assert.NoError(t, err)
	assert.NotNil(t, sub)
	assert.True(t, sub.Email)
	assert.False(t, sub.SMS)
	assert.True(t, sub.Webhook)
	assert.False(t, sub.Station)

	// Update the subscription
	updatedChannels := map[string]bool{
		"email":   false,
		"sms":     true,
		"webhook": false,
		"station": true,
	}

	err = service.SetSubscriptionChannels(ctx, 2, messageType.ID, updatedChannels)
	assert.NoError(t, err)

	// Verify the subscription was updated
	updatedSub, err := service.GetUserSubscription(ctx, 2, messageType.ID)
	assert.NoError(t, err)
	assert.NotNil(t, updatedSub)
	assert.False(t, updatedSub.Email)
	assert.True(t, updatedSub.SMS)
	assert.False(t, updatedSub.Webhook)
	assert.True(t, updatedSub.Station)
}

func TestMessageSubscriptionService_ListMessageTypeSubscriptions(t *testing.T) {
	db := testutils.SetupTestDB()
	defer testutils.TeardownTestDB(db)

	service := NewMessageSubscriptionService(db)
	ctx := context.Background()

	// Create a message type
	messageType := &model.MessageType{
		Name:        "news_notification",
		DisplayName: "News Notification",
		Description: "News notification type",
		Enabled:     true,
	}
	db.Create(messageType)

	// Create multiple subscriptions for the same message type
	sub1 := &model.MessageSubscription{
		UserID:        3,
		MessageTypeID: messageType.ID,
		Email:         true,
		Station:       true,
	}
	sub2 := &model.MessageSubscription{
		UserID:        4,
		MessageTypeID: messageType.ID,
		Email:         false,
		Station:       true,
	}

	service.CreateSubscription(ctx, sub1)
	service.CreateSubscription(ctx, sub2)

	// Test ListMessageTypeSubscriptions
	subs, err := service.ListMessageTypeSubscriptions(ctx, messageType.ID)
	assert.NoError(t, err)
	assert.Len(t, subs, 2)

	// Verify that both subscriptions belong to the correct message type
	userIDs := make([]uint, len(subs))
	for i, sub := range subs {
		userIDs[i] = sub.UserID
	}
	assert.Contains(t, userIDs, uint(3))
	assert.Contains(t, userIDs, uint(4))
}

func TestMessageSubscriptionService_GetSubscribersByChannel(t *testing.T) {
	db := testutils.SetupTestDB()
	defer testutils.TeardownTestDB(db)

	service := NewMessageSubscriptionService(db)
	ctx := context.Background()

	// Create a message type
	messageType := &model.MessageType{
		Name:        "important_notification",
		DisplayName: "Important Notification",
		Description: "Important notification type",
		Enabled:     true,
	}
	db.Create(messageType)

	// Create subscriptions with different channel settings
	emailSub := &model.MessageSubscription{
		UserID:        5,
		MessageTypeID: messageType.ID,
		Email:         true,
		SMS:           false,
		Station:       true,
	}
	webhookSub := &model.MessageSubscription{
		UserID:        6,
		MessageTypeID: messageType.ID,
		Email:         false,
		SMS:           false,
		Webhook:       true,
		Station:       false,
	}
	bothSub := &model.MessageSubscription{
		UserID:        7,
		MessageTypeID: messageType.ID,
		Email:         true,
		SMS:           false,
		Webhook:       true,
		Station:       true,
	}

	service.CreateSubscription(ctx, emailSub)
	service.CreateSubscription(ctx, webhookSub)
	service.CreateSubscription(ctx, bothSub)

	// Test GetSubscribersByChannel for email
	emailSubscribers, err := service.GetSubscribersByChannel(ctx, messageType.ID, "email")
	assert.NoError(t, err)
	assert.Len(t, emailSubscribers, 2) // users 5 and 7
	assert.Contains(t, emailSubscribers, uint(5))
	assert.Contains(t, emailSubscribers, uint(7))

	// Test GetSubscribersByChannel for webhook
	webhookSubscribers, err := service.GetSubscribersByChannel(ctx, messageType.ID, "webhook")
	assert.NoError(t, err)
	assert.Len(t, webhookSubscribers, 2) // users 6 and 7
	assert.Contains(t, webhookSubscribers, uint(6))
	assert.Contains(t, webhookSubscribers, uint(7))

	// Test GetSubscribersByChannel for station
	stationSubscribers, err := service.GetSubscribersByChannel(ctx, messageType.ID, "station")
	assert.NoError(t, err)
	assert.Len(t, stationSubscribers, 2) // users 5 and 7
	assert.Contains(t, stationSubscribers, uint(5))
	assert.Contains(t, stationSubscribers, uint(7))
}
