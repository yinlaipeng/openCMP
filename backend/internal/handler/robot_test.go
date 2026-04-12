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

func TestRobotHandler(t *testing.T) {
	db := testutils.SetupTestDB()
	defer testutils.TeardownTestDB(db)

	logger, _ := zap.NewDevelopment()
	handler := NewRobotHandler(db, logger)

	// Create a test robot
	robot := &model.Robot{
		Name:        "Test Robot",
		Type:        "dingtalk",
		WebhookURL:  "https://oapi.dingtalk.com/robot/send?access_token=test",
		Description: "Test dingtalk robot",
		Enabled:     true,
	}
	db.Create(robot)

	t.Run("Get", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: strconv.Itoa(int(robot.ID))}}

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
		c.Request = httptest.NewRequest("GET", "/robots", nil)

		handler.List(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Create", func(t *testing.T) {
		newRobot := &model.Robot{
			Name:        "New Robot",
			Type:        "wechat",
			WebhookURL:  "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=test",
			Description: "New test robot",
			Enabled:     true,
		}

		jsonBytes, _ := json.Marshal(newRobot)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/robots", bytes.NewReader(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Create(c)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Update", func(t *testing.T) {
		updateData := &model.Robot{
			Name:        "Updated Robot",
			Type:        "webhook",
			WebhookURL:  "https://webhook.test/updated",
			Description: "Updated description",
			Enabled:     false,
		}

		jsonBytes, _ := json.Marshal(updateData)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: strconv.Itoa(int(robot.ID))}}
		c.Request = httptest.NewRequest("PUT", "/robots/"+strconv.Itoa(int(robot.ID)), bytes.NewReader(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Update(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Enable", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: strconv.Itoa(int(robot.ID))}}
		c.Request = httptest.NewRequest("POST", "/robots/"+strconv.Itoa(int(robot.ID))+"/enable", nil)

		handler.Enable(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Disable", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: strconv.Itoa(int(robot.ID))}}
		c.Request = httptest.NewRequest("POST", "/robots/"+strconv.Itoa(int(robot.ID))+"/disable", nil)

		handler.Disable(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Delete", func(t *testing.T) {
		// Create another robot to delete
		deleteRobot := &model.Robot{
			Name:        "Delete Robot",
			Type:        "feishu",
			WebhookURL:  "https://open.feishu.cn/open-apis/bot/v2/hook/test",
			Description: "To be deleted",
			Enabled:     true,
		}
		db.Create(deleteRobot)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: strconv.Itoa(int(deleteRobot.ID))}}
		c.Request = httptest.NewRequest("DELETE", "/robots/"+strconv.Itoa(int(deleteRobot.ID)), nil)

		handler.Delete(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
