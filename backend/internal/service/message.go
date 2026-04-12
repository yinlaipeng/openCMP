package service

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
)

// MessageService 消息服务
type MessageService struct {
	db *gorm.DB
}

// NewMessageService 创建消息服务
func NewMessageService(db *gorm.DB) *MessageService {
	return &MessageService{db: db}
}

// CreateMessage 创建消息
func (s *MessageService) CreateMessage(ctx context.Context, msg *model.Message) error {
	return s.db.WithContext(ctx).Create(msg).Error
}

// GetMessage 获取消息
func (s *MessageService) GetMessage(ctx context.Context, id uint) (*model.Message, error) {
	var msg model.Message
	err := s.db.WithContext(ctx).First(&msg, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &msg, nil
}

// ListMessages 列出消息
func (s *MessageService) ListMessages(ctx context.Context, receiverID uint, limit, offset int) ([]*model.Message, int64, error) {
	var messages []*model.Message
	var total int64

	query := s.db.Model(&model.Message{}).Where("receiver_id = ?", receiverID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&messages).Error

	return messages, total, err
}

// ListUnreadMessages 列出未读消息
func (s *MessageService) ListUnreadMessages(ctx context.Context, receiverID uint) ([]*model.Message, error) {
	var messages []*model.Message
	err := s.db.WithContext(ctx).
		Where("receiver_id = ? AND read = false", receiverID).
		Order("created_at DESC").
		Find(&messages).Error
	return messages, err
}

// MarkAsRead 标记消息为已读
func (s *MessageService) MarkAsRead(ctx context.Context, id uint) error {
	now := time.Now()
	return s.db.WithContext(ctx).Model(&model.Message{}).Where("id = ?", id).Updates(map[string]interface{}{
		"read":    true,
		"read_at": now,
	}).Error
}

// MarkAllAsRead 标记所有消息为已读
func (s *MessageService) MarkAllAsRead(ctx context.Context, receiverID uint) error {
	now := time.Now()
	return s.db.WithContext(ctx).Model(&model.Message{}).Where("receiver_id = ? AND read = false", receiverID).Updates(map[string]interface{}{
		"read":    true,
		"read_at": now,
	}).Error
}

// DeleteMessage 删除消息
func (s *MessageService) DeleteMessage(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Delete(&model.Message{}, id).Error
}

// GetUnreadCount 获取未读消息数量
func (s *MessageService) GetUnreadCount(ctx context.Context, receiverID uint) (int64, error) {
	var count int64
	err := s.db.WithContext(ctx).Model(&model.Message{}).Where("receiver_id = ? AND read = false", receiverID).Count(&count).Error
	return count, err
}

// SendMessage 发送消息
func (s *MessageService) SendMessage(ctx context.Context, title, content string, level string, receiverID uint, senderID uint) error {
	msg := &model.Message{
		Title:      title,
		Content:    content,
		Level:      level,
		ReceiverID: receiverID,
		SenderID:   senderID,
		Read:       false,
	}
	return s.CreateMessage(ctx, msg)
}

// BroadcastMessage 广播消息（发送给所有用户）
func (s *MessageService) BroadcastMessage(ctx context.Context, title, content string, level string, senderID uint) error {
	// 获取所有启用用户
	var users []model.User
	if err := s.db.WithContext(ctx).Where("enabled = ?", true).Find(&users).Error; err != nil {
		return err
	}

	for _, user := range users {
		msg := &model.Message{
			Title:      title,
			Content:    content,
			Level:      level,
			ReceiverID: user.ID,
			SenderID:   senderID,
			Read:       false,
		}
		if err := s.CreateMessage(ctx, msg); err != nil {
			continue // 忽略单个错误
		}
	}
	return nil
}
