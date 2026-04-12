package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
	"github.com/opencmp/opencmp/testutils"
)

func TestReceiverService(t *testing.T) {
	db := testutils.SetupTestDB()
	defer testutils.TeardownTestDB(db)

	service := NewReceiverService(db)
	ctx := context.Background()

	// Create a test receiver
	receiver := &model.Receiver{
		Name:    "Test Receiver",
		Email:   "test@example.com",
		Phone:   "1234567890",
		Enabled: true,
	}

	// Test CreateReceiver
	err := service.CreateReceiver(ctx, receiver)
	assert.NoError(t, err)
	assert.NotZero(t, receiver.ID)

	// Test GetReceiver
	retrievedReceiver, err := service.GetReceiver(ctx, receiver.ID)
	assert.NoError(t, err)
	assert.Equal(t, receiver.Name, retrievedReceiver.Name)
	assert.Equal(t, receiver.Email, retrievedReceiver.Email)

	// Test GetReceiverByName
	receiverByName, err := service.GetReceiverByName(ctx, "Test Receiver")
	assert.NoError(t, err)
	assert.NotNil(t, receiverByName)
	assert.Equal(t, receiver.ID, receiverByName.ID)

	// Test ListReceivers
	receivers, total, err := service.ListReceivers(ctx, 10, 0)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, receivers, 1)
	assert.Equal(t, receiver.ID, receivers[0].ID)

	// Test UpdateReceiver
	receiver.Email = "updated@example.com"
	err = service.UpdateReceiver(ctx, receiver)
	assert.NoError(t, err)

	updatedReceiver, err := service.GetReceiver(ctx, receiver.ID)
	assert.NoError(t, err)
	assert.Equal(t, "updated@example.com", updatedReceiver.Email)

	// Test DisableReceiver
	err = service.DisableReceiver(ctx, receiver.ID)
	assert.NoError(t, err)

	disabledReceiver, err := service.GetReceiver(ctx, receiver.ID)
	assert.NoError(t, err)
	assert.False(t, disabledReceiver.Enabled)

	// Test EnableReceiver
	err = service.EnableReceiver(ctx, receiver.ID)
	assert.NoError(t, err)

	enabledReceiver, err := service.GetReceiver(ctx, receiver.ID)
	assert.NoError(t, err)
	assert.True(t, enabledReceiver.Enabled)

	// Test DeleteReceiver
	err = service.DeleteReceiver(ctx, receiver.ID)
	assert.NoError(t, err)

	// Verify deletion
	deletedReceiver, err := service.GetReceiver(ctx, receiver.ID)
	assert.NoError(t, err)
	assert.Nil(t, deletedReceiver)
}

func TestReceiverService_CreateFromUser(t *testing.T) {
	db := testutils.SetupTestDB()
	defer testutils.TeardownTestDB(db)

	service := NewReceiverService(db)
	ctx := context.Background()

	// First, create a user to associate with the receiver
	user := &model.User{
		Name:        "testuser",
		DisplayName: "Test User",
		Email:       "testuser@example.com",
		Phone:       "9876543210",
		DomainID:    1,
		Enabled:     true,
		Password:    "hashed_password",
	}
	db.Create(user)

	// Test CreateReceiverFromUser
	receiver, err := service.CreateReceiverFromUser(ctx, user)
	assert.NoError(t, err)
	assert.NotNil(t, receiver)
	assert.Equal(t, user.DisplayName, receiver.Name)
	assert.Equal(t, user.Email, receiver.Email)
	assert.Equal(t, user.Phone, receiver.Phone)
	assert.Equal(t, user.ID, *receiver.UserID)

	// Test that creating from the same user returns the existing receiver
	existingReceiver, err := service.CreateReceiverFromUser(ctx, user)
	assert.NoError(t, err)
	assert.Equal(t, receiver.ID, existingReceiver.ID)
}

func TestReceiverService_GetByUserID(t *testing.T) {
	db := testutils.SetupTestDB()
	defer testutils.TeardownTestDB(db)

	service := NewReceiverService(db)
	ctx := context.Background()

	// Create a user first
	user := &model.User{
		Name:        "testuser2",
		DisplayName: "Test User 2",
		Email:       "testuser2@example.com",
		DomainID:    1,
		Enabled:     true,
		Password:    "hashed_password",
	}
	db.Create(user)

	// Create a receiver linked to the user
	receiver := &model.Receiver{
		Name:    "Linked Receiver",
		Email:   "linked@example.com",
		UserID:  &user.ID,
		Enabled: true,
	}
	err := service.CreateReceiver(ctx, receiver)
	assert.NoError(t, err)

	// Test GetReceiverByUserID
	foundReceiver, err := service.GetReceiverByUserID(ctx, user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, foundReceiver)
	assert.Equal(t, receiver.ID, foundReceiver.ID)
}
