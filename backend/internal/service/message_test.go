package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
	"github.com/opencmp/opencmp/testutils"
)

func TestMessageService(t *testing.T) {
	db := testutils.SetupTestDB()
	defer testutils.TeardownTestDB(db)

	service := NewMessageService(db)

	ctx := context.Background()

	// Create a test message
	msg := &model.Message{
		Title:      "Test Message",
		Content:    "This is a test message",
		Level:      "info",
		SenderID:   1,
		ReceiverID: 1,
		Read:       false,
	}

	// Test CreateMessage
	err := service.CreateMessage(ctx, msg)
	assert.NoError(t, err)
	assert.NotZero(t, msg.ID)

	// Test GetMessage
	retrievedMsg, err := service.GetMessage(ctx, msg.ID)
	assert.NoError(t, err)
	assert.Equal(t, msg.Title, retrievedMsg.Title)
	assert.Equal(t, msg.Content, retrievedMsg.Content)

	// Test ListMessages
	messages, total, err := service.ListMessages(ctx, 1, 10, 0)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, messages, 1)
	assert.Equal(t, msg.ID, messages[0].ID)

	// Test MarkAsRead
	err = service.MarkAsRead(ctx, msg.ID)
	assert.NoError(t, err)

	// Verify the message is marked as read
	updatedMsg, err := service.GetMessage(ctx, msg.ID)
	assert.NoError(t, err)
	assert.True(t, updatedMsg.Read)

	// Test MarkAllAsRead
	err = service.MarkAllAsRead(ctx, 1)
	assert.NoError(t, err)

	// Test GetUnreadCount
	count, err := service.GetUnreadCount(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)

	// Test DeleteMessage
	err = service.DeleteMessage(ctx, msg.ID)
	assert.NoError(t, err)

	// Verify deletion
	deletedMsg, err := service.GetMessage(ctx, msg.ID)
	assert.NoError(t, err)
	assert.Nil(t, deletedMsg)
}

func TestMessageService_ListUnreadMessages(t *testing.T) {
	db := testutils.SetupTestDB()
	defer testutils.TeardownTestDB(db)

	service := NewMessageService(db)
	ctx := context.Background()

	// Create two messages for user 1, one read and one unread
	msg1 := &model.Message{
		Title:      "Unread Message",
		Content:    "This is unread",
		Level:      "info",
		SenderID:   1,
		ReceiverID: 1,
		Read:       false,
	}
	msg2 := &model.Message{
		Title:      "Read Message",
		Content:    "This is read",
		Level:      "info",
		SenderID:   1,
		ReceiverID: 1,
		Read:       true,
	}

	service.CreateMessage(ctx, msg1)
	service.CreateMessage(ctx, msg2)

	// Test ListUnreadMessages
	unreadMessages, err := service.ListUnreadMessages(ctx, 1)
	assert.NoError(t, err)
	assert.Len(t, unreadMessages, 1)
	assert.Equal(t, msg1.Title, unreadMessages[0].Title)
}

func TestMessageService_SendMessage(t *testing.T) {
	db := testutils.SetupTestDB()
	defer testutils.TeardownTestDB(db)

	service := NewMessageService(db)
	ctx := context.Background()

	// Test SendMessage
	err := service.SendMessage(ctx, "Send Test", "Test content", "info", 1, 1)
	assert.NoError(t, err)

	// Verify message was created
	messages, total, err := service.ListMessages(ctx, 1, 10, 0)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, "Send Test", messages[0].Title)
}