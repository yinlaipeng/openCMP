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
	"github.com/opencmp/opencmp/internal/service"
	"github.com/opencmp/opencmp/testutils"
)

func TestMessageHandler(t *testing.T) {
	db := testutils.SetupTestDB()
	defer testutils.TeardownTestDB(db)

	logger, _ := zap.NewDevelopment()
	handler := NewMessageHandler(db, logger)

	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

	// Create a test message
	msg := &model.Message{
		Title:      "Test Message",
		Content:    "Test content",
		Level:      "info",
		SenderID:   1,
		ReceiverID: 1,
		Read:       false,
	}
	db.Create(msg)

	t.Run("Get", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: strconv.Itoa(int(msg.ID))}}

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
		c.Request = httptest.NewRequest("GET", "/messages?user_id=1", nil)

		handler.List(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("MarkRead", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: strconv.Itoa(int(msg.ID))}}

		handler.MarkRead(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Delete", func(t *testing.T) {
		// Create another message to delete
		deleteMsg := &model.Message{
			Title:      "Delete Test",
			Content:    "To be deleted",
			Level:      "info",
			SenderID:   1,
			ReceiverID: 1,
			Read:       false,
		}
		db.Create(deleteMsg)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: strconv.Itoa(int(deleteMsg.ID))}}

		handler.Delete(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GetUnreadCount", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/messages/unread-count?user_id=1", nil)

		handler.GetUnreadCount(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}