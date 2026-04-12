package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
	"github.com/opencmp/opencmp/testutils"
)

func TestMessageSubscriptionHandler(t *testing.T) {
	db := testutils.SetupTestDB()
	defer testutils.TeardownTestDB(db)

	logger, _ := zap.NewDevelopment()
	handler := NewMessageSubscriptionHandler(db, logger)

	// Create a test message type
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
	db.Create(subscription)

	t.Run("Get", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: strconv.Itoa(int(subscription.ID))}}

		handler.Get(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Get_NotFound", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "999"}}

		handler.Get(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("List", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/subscriptions?user_id=1", nil)

		handler.List(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Create", func(t *testing.T) {
		newSub := &model.MessageSubscription{
			UserID:        2,
			MessageTypeID: messageType.ID,
			Email:         false,
			SMS:           true,
			Webhook:       false,
			Station:       true,
		}

		jsonBytes, _ := json.Marshal(newSub)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/subscriptions", bytes.NewReader(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Create(c)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Update", func(t *testing.T) {
		updateData := &model.MessageSubscription{
			UserID:        1,
			MessageTypeID: messageType.ID,
			Email:         false,
			SMS:           true,
			Webhook:       false,
			Station:       false,
		}

		jsonBytes, _ := json.Marshal(updateData)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: strconv.Itoa(int(subscription.ID))}}
		c.Request = httptest.NewRequest("PUT", "/subscriptions/"+strconv.Itoa(int(subscription.ID)), bytes.NewReader(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Update(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Delete", func(t *testing.T) {
		// Create another subscription to delete
		deleteSub := &model.MessageSubscription{
			UserID:        3,
			MessageTypeID: messageType.ID,
			Email:         true,
			SMS:           false,
			Webhook:       false,
			Station:       true,
		}
		db.Create(deleteSub)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: strconv.Itoa(int(deleteSub.ID))}}
		c.Request = httptest.NewRequest("DELETE", "/subscriptions/"+strconv.Itoa(int(deleteSub.ID)), nil)

		handler.Delete(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("ListMessageTypes", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/message-types", nil)

		handler.ListMessageTypes(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
