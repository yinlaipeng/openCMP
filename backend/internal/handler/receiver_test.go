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

func TestReceiverHandler(t *testing.T) {
	db := testutils.SetupTestDB()
	defer testutils.TeardownTestDB(db)

	logger, _ := zap.NewDevelopment()
	handler := NewReceiverHandler(db, logger)

	// Create a test receiver
	receiver := &model.Receiver{
		Name:      "Test Receiver",
		Email:     "test@example.com",
		Phone:     "1234567890",
		Enabled:   true,
	}
	db.Create(receiver)

	t.Run("Get", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: strconv.Itoa(int(receiver.ID))}}

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
		c.Request = httptest.NewRequest("GET", "/receivers", nil)

		handler.List(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Create", func(t *testing.T) {
		newReceiver := &model.Receiver{
			Name:      "New Receiver",
			Email:     "new@example.com",
			Phone:     "0987654321",
			Enabled:   true,
		}

		jsonBytes, _ := json.Marshal(newReceiver)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/receivers", bytes.NewReader(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Create(c)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Update", func(t *testing.T) {
		updateData := &model.Receiver{
			Name:      "Updated Receiver",
			Email:     "updated@example.com",
			Phone:     "1111111111",
			Enabled:   false,
		}

		jsonBytes, _ := json.Marshal(updateData)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: strconv.Itoa(int(receiver.ID))}}
		c.Request = httptest.NewRequest("PUT", "/receivers/"+strconv.Itoa(int(receiver.ID)), bytes.NewReader(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Update(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Delete", func(t *testing.T) {
		// Create another receiver to delete
		deleteReceiver := &model.Receiver{
			Name:      "Delete Receiver",
			Email:     "delete@example.com",
			Phone:     "2222222222",
			Enabled:   true,
		}
		db.Create(deleteReceiver)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: strconv.Itoa(int(deleteReceiver.ID))}}
		c.Request = httptest.NewRequest("DELETE", "/receivers/"+strconv.Itoa(int(deleteReceiver.ID)), nil)

		handler.Delete(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}