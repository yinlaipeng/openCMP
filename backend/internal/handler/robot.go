package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

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

// Test 测试机器人连通性并发送测试消息
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

	// 检查机器人是否启用
	if !robot.Enabled {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "机器人未启用"})
		return
	}

	// 根据机器人类型发送测试消息
	testPassed := false
	testMessage := ""
	testDetail := ""

	switch robot.Type {
	case "dingtalk":
		testPassed, testMessage, testDetail = h.testDingTalkRobot(robot)
	case "feishu", "lark":
		testPassed, testMessage, testDetail = h.testFeishuRobot(robot)
	case "wechat":
		testPassed, testMessage, testDetail = h.testWeChatRobot(robot)
	case "webhook":
		testPassed, testMessage, testDetail = h.testGenericWebhook(robot)
	case "email":
		// 邮件机器人测试通过NotificationChannel进行
		testPassed = true
		testMessage = "邮件机器人配置验证通过"
		testDetail = "请通过通知渠道测试邮件发送"
	default:
		testMessage = "不支持的机器人类型: " + robot.Type
	}

	c.JSON(http.StatusOK, gin.H{
		"success": testPassed,
		"message": testMessage,
		"detail":  testDetail,
		"robot":   robot.Name,
		"type":    robot.Type,
	})
}

// testDingTalkRobot 测试钉钉机器人
func (h *RobotHandler) testDingTalkRobot(robot *model.Robot) (bool, string, string) {
	if robot.WebhookURL == "" {
		return false, "缺少Webhook地址", ""
	}

	// 构建钉钉消息格式
	msg := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": "【openCMP测试消息】\n这是一条来自openCMP系统的测试消息，用于验证机器人配置。\n发送时间: " + time.Now().Format("2006-01-02 15:04:05"),
		},
	}

	return h.sendWebhookMessage(robot.WebhookURL, msg, "钉钉")
}

// testFeishuRobot 测试飞书机器人
func (h *RobotHandler) testFeishuRobot(robot *model.Robot) (bool, string, string) {
	if robot.WebhookURL == "" {
		return false, "缺少Webhook地址", ""
	}

	// 构建飞书消息格式
	msg := map[string]interface{}{
		"msg_type": "text",
		"content": map[string]string{
			"text": "【openCMP测试消息】\n这是一条来自openCMP系统的测试消息，用于验证机器人配置。\n发送时间: " + time.Now().Format("2006-01-02 15:04:05"),
		},
	}

	return h.sendWebhookMessage(robot.WebhookURL, msg, "飞书")
}

// testWeChatRobot 测试企业微信机器人
func (h *RobotHandler) testWeChatRobot(robot *model.Robot) (bool, string, string) {
	if robot.WebhookURL == "" {
		return false, "缺少Webhook地址", ""
	}

	// 构建企业微信消息格式
	msg := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": "【openCMP测试消息】\n这是一条来自openCMP系统的测试消息，用于验证机器人配置。\n发送时间: " + time.Now().Format("2006-01-02 15:04:05"),
		},
	}

	return h.sendWebhookMessage(robot.WebhookURL, msg, "企业微信")
}

// testGenericWebhook 测试通用Webhook
func (h *RobotHandler) testGenericWebhook(robot *model.Robot) (bool, string, string) {
	if robot.WebhookURL == "" {
		return false, "缺少Webhook地址", ""
	}

	// 构建通用消息格式
	msg := map[string]interface{}{
		"source":    "openCMP",
		"type":      "test",
		"content":   "openCMP系统测试消息",
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		"robot":     robot.Name,
	}

	return h.sendWebhookMessage(robot.WebhookURL, msg, "通用Webhook")
}

// sendWebhookMessage 发送Webhook消息
func (h *RobotHandler) sendWebhookMessage(webhookURL string, msg map[string]interface{}, platform string) (bool, string, string) {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return false, "消息格式化失败", err.Error()
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(webhookURL, "application/json", bytes.NewReader(msgBytes))
	if err != nil {
		return false, "消息发送失败", err.Error()
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return true, platform + "机器人测试成功", "HTTP状态码: " + resp.Status
	}

	// 尝试解析错误响应
	var respBody map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&respBody)

	errMsg := "HTTP状态码: " + resp.Status
	if respBody != nil {
		if errmsg, ok := respBody["errmsg"].(string); ok {
			errMsg += ", 错误信息: " + errmsg
		}
	}

	return false, platform + "机器人返回错误", errMsg
}
