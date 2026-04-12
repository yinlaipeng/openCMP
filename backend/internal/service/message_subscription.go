package service

import (
	"context"

	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
)

// MessageSubscriptionService 消息订阅服务
type MessageSubscriptionService struct {
	db *gorm.DB
}

// NewMessageSubscriptionService 创建消息订阅服务
func NewMessageSubscriptionService(db *gorm.DB) *MessageSubscriptionService {
	return &MessageSubscriptionService{db: db}
}

// CreateSubscription 创建订阅
func (s *MessageSubscriptionService) CreateSubscription(ctx context.Context, sub *model.MessageSubscription) error {
	return s.db.WithContext(ctx).Create(sub).Error
}

// GetSubscription 获取订阅
func (s *MessageSubscriptionService) GetSubscription(ctx context.Context, id uint) (*model.MessageSubscription, error) {
	var sub model.MessageSubscription
	err := s.db.WithContext(ctx).First(&sub, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &sub, nil
}

// GetUserSubscription 获取用户对某类型消息的订阅
func (s *MessageSubscriptionService) GetUserSubscription(ctx context.Context, userID, messageTypeID uint) (*model.MessageSubscription, error) {
	var sub model.MessageSubscription
	err := s.db.WithContext(ctx).
		Where("user_id = ? AND message_type_id = ?", userID, messageTypeID).
		First(&sub).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &sub, nil
}

// ListUserSubscriptions 列出用户的所有订阅
func (s *MessageSubscriptionService) ListUserSubscriptions(ctx context.Context, userID uint) ([]*model.MessageSubscription, error) {
	var subs []*model.MessageSubscription
	err := s.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&subs).Error
	return subs, err
}

// UpdateSubscription 更新订阅
func (s *MessageSubscriptionService) UpdateSubscription(ctx context.Context, sub *model.MessageSubscription) error {
	return s.db.WithContext(ctx).Save(sub).Error
}

// DeleteSubscription 删除订阅
func (s *MessageSubscriptionService) DeleteSubscription(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Delete(&model.MessageSubscription{}, id).Error
}

// SetSubscriptionChannels 设置订阅渠道
func (s *MessageSubscriptionService) SetSubscriptionChannels(ctx context.Context, userID, messageTypeID uint, channels map[string]bool) error {
	sub, err := s.GetUserSubscription(ctx, userID, messageTypeID)
	if err != nil {
		return err
	}

	if sub == nil {
		// 创建新订阅
		sub = &model.MessageSubscription{
			UserID:        userID,
			MessageTypeID: messageTypeID,
			Email:         channels["email"],
			SMS:           channels["sms"],
			Webhook:       channels["webhook"],
			Station:       channels["station"],
		}
		return s.CreateSubscription(ctx, sub)
	}

	// 更新现有订阅
	return s.db.WithContext(ctx).Model(sub).Updates(map[string]interface{}{
		"email":   channels["email"],
		"sms":     channels["sms"],
		"webhook": channels["webhook"],
		"station": channels["station"],
	}).Error
}

// ListMessageTypeSubscriptions 列出某消息类型的所有订阅
func (s *MessageSubscriptionService) ListMessageTypeSubscriptions(ctx context.Context, messageTypeID uint) ([]*model.MessageSubscription, error) {
	var subs []*model.MessageSubscription
	err := s.db.WithContext(ctx).
		Where("message_type_id = ?", messageTypeID).
		Find(&subs).Error
	return subs, err
}

// GetSubscribersByChannel 获取通过特定渠道订阅某消息类型的用户
func (s *MessageSubscriptionService) GetSubscribersByChannel(ctx context.Context, messageTypeID uint, channel string) ([]uint, error) {
	var userIDs []uint
	query := s.db.WithContext(ctx).
		Model(&model.MessageSubscription{}).
		Where("message_type_id = ?", messageTypeID)

	switch channel {
	case "email":
		query = query.Where("email = ?", true)
	case "sms":
		query = query.Where("sms = ?", true)
	case "webhook":
		query = query.Where("webhook = ?", true)
	case "station":
		query = query.Where("station = ?", true)
	}

	err := query.Pluck("user_id", &userIDs).Error
	return userIDs, err
}
