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

func TestNotificationChannelHandler(t *testing.T) {
	db := testutils.SetupTestDB()
	defer testutils.TeardownTestDB(db)

	logger, _ := zap.NewDevelopment()
	handler := NewNotificationChannelHandler(db, logger)

	// Create a test notification channel
	channel := &model.NotificationChannel{
		Name:        "Test Channel",
		Type:        "email",
		Description: "Test description",
		Config:      []byte(`{"smtp_host":"smtp.test.com"}`),
		Enabled:     true,
	}
	db.Create(channel)

	t.Run("Get", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: strconv.Itoa(int(channel.ID))}}

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
		c.Request = httptest.NewRequest("GET", "/notification-channels", nil)

		handler.List(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Create", func(t *testing.T) {
		newChannel := &model.NotificationChannel{
			Name:        "New Channel",
			Type:        "webhook",
			Description: "New test channel",
			Config:      []byte(`{"url":"http://webhook.test"}`),
			Enabled:     true,
		}

		jsonBytes, _ := json.Marshal(newChannel)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/notification-channels", bytes.NewReader(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Create(c)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Update", func(t *testing.T) {
		updateData := &model.NotificationChannel{
			Name:        "Updated Channel",
			Type:        "email",
			Description: "Updated description",
			Config:      []byte(`{"smtp_host":"updated.smtp.com"}`),
			Enabled:     false,
		}

		jsonBytes, _ := json.Marshal(updateData)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: strconv.Itoa(int(channel.ID))}}
		c.Request = httptest.NewRequest("PUT", "/notification-channels/"+strconv.Itoa(int(channel.ID)), bytes.NewReader(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Update(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Enable", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: strconv.Itoa(int(channel.ID))}}
		c.Request = httptest.NewRequest("POST", "/notification-channels/"+strconv.Itoa(int(channel.ID))+"/enable", nil)

		handler.Enable(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Disable", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: strconv.Itoa(int(channel.ID))}}
		c.Request = httptest.NewRequest("POST", "/notification-channels/"+strconv.Itoa(int(channel.ID))+"/disable", nil)

		handler.Disable(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Delete", func(t *testing.T) {
		// Create another channel to delete
		deleteChannel := &model.NotificationChannel{
			Name:        "Delete Channel",
			Type:        "sms",
			Description: "To be deleted",
			Config:      []byte(`{"api_key":"test"}`),
			Enabled:     true,
		}
		db.Create(deleteChannel)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: strconv.Itoa(int(deleteChannel.ID))}}
		c.Request = httptest.NewRequest("DELETE", "/notification-channels/"+strconv.Itoa(int(deleteChannel.ID)), nil)

		handler.Delete(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}