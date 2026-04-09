package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
	"github.com/opencmp/opencmp/internal/service"
)

// RobotHandler 机器人 Handler
type RobotHandler struct {
	robotService *service.RobotService
	logger       *zap.Logger
}

// NewRobotHandler 创建机器人 Handler
func NewRobotHandler(db *gorm.DB, logger *zap.Logger) *RobotHandler {
	return &RobotHandler{
		robotService: service.NewRobotService(db),
		logger:       logger,
	}
}

// List 获取机器人列表
func (h *RobotHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	robotType := c.Query("type")

	robots, total, err := h.robotService.ListRobots(c.Request.Context(), robotType, pageSize, offset)
	if err != nil {
		h.logger.Error("failed to list robots", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": robots,
		"total": total,
	})
}

// Get 获取机器人详情
func (h *RobotHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	robot, err := h.robotService.GetRobot(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get robot", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if robot == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "robot not found"})
		return
	}

	c.JSON(http.StatusOK, robot)
}

// Create 创建机器人
func (h *RobotHandler) Create(c *gin.Context) {
	var robot model.Robot
	if err := c.ShouldBindJSON(&robot); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.robotService.CreateRobot(c.Request.Context(), &robot); err != nil {
		h.logger.Error("failed to create robot", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, robot)
}

// Update 更新机器人
func (h *RobotHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	robot, err := h.robotService.GetRobot(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get robot", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if robot == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "robot not found"})
		return
	}

	if err := c.ShouldBindJSON(robot); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	robot.ID = uint(id)

	if err := h.robotService.UpdateRobot(c.Request.Context(), robot); err != nil {
		h.logger.Error("failed to update robot", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, robot)
}

// Delete 删除机器人
func (h *RobotHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.robotService.DeleteRobot(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("failed to delete robot", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// Enable 启用机器人
func (h *RobotHandler) Enable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.robotService.EnableRobot(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("failed to enable robot", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "enabled"})
}

// Disable 禁用机器人
func (h *RobotHandler) Disable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.robotService.DisableRobot(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("failed to disable robot", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "disabled"})
}

// Test 测试机器人连通性
func (h *RobotHandler) Test(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	robot, err := h.robotService.GetRobot(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get robot", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if robot == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "robot not found"})
		return
	}

	// 简单的连通性测试：验证机器人存在且启用
	if !robot.Enabled {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "robot is disabled"})
		return
	}

	// 根据机器人类型进行具体的连通性测试
	testPassed := true
	testMessage := "robot test passed"

	switch robot.Type {
	case "dingtalk", "lark", "feishu", "wechat", "webhook":
		// 验证必要字段是否存在
		if robot.WebhookURL == "" {
			testPassed = false
			testMessage = "missing required webhook URL"
		} else {
			// 在实际应用中，这里应该尝试向Webhook发送测试消息
			testMessage = "webhook validation passed"
		}
	default:
		testMessage = "unsupported robot type: " + robot.Type
		testPassed = false
	}

	c.JSON(http.StatusOK, gin.H{"success": testPassed, "message": testMessage})
}
