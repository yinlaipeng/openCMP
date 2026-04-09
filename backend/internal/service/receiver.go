package service

import (
	"context"

	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
)

// ReceiverService 接收人服务
type ReceiverService struct {
	db *gorm.DB
}

// NewReceiverService 创建接收人服务
func NewReceiverService(db *gorm.DB) *ReceiverService {
	return &ReceiverService{db: db}
}

// CreateReceiver 创建接收人
func (s *ReceiverService) CreateReceiver(ctx context.Context, receiver *model.Receiver) error {
	return s.db.WithContext(ctx).Create(receiver).Error
}

// GetReceiver 获取接收人
func (s *ReceiverService) GetReceiver(ctx context.Context, id uint) (*model.Receiver, error) {
	var receiver model.Receiver
	err := s.db.WithContext(ctx).First(&receiver, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &receiver, nil
}

// GetReceiverByName 根据名称获取接收人
func (s *ReceiverService) GetReceiverByName(ctx context.Context, name string) (*model.Receiver, error) {
	var receiver model.Receiver
	err := s.db.WithContext(ctx).Where("name = ?", name).First(&receiver).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &receiver, nil
}

// ListReceivers 列出接收人
func (s *ReceiverService) ListReceivers(ctx context.Context, limit, offset int) ([]*model.Receiver, int64, error) {
	var receivers []*model.Receiver
	var total int64

	query := s.db.Model(&model.Receiver{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Preload("Domain"). // Preload domain info
		Find(&receivers).Error

	return receivers, total, err
}

// UpdateReceiver 更新接收人
func (s *ReceiverService) UpdateReceiver(ctx context.Context, receiver *model.Receiver) error {
	return s.db.WithContext(ctx).Save(receiver).Error
}

// DeleteReceiver 删除接收人
func (s *ReceiverService) DeleteReceiver(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Delete(&model.Receiver{}, id).Error
}

// EnableReceiver 启用接收人
func (s *ReceiverService) EnableReceiver(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Model(&model.Receiver{}).Where("id = ?", id).Update("enabled", true).Error
}

// DisableReceiver 禁用接收人
func (s *ReceiverService) DisableReceiver(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Model(&model.Receiver{}).Where("id = ?", id).Update("enabled", false).Error
}

// GetReceiverByUserID 根据用户ID获取接收人
func (s *ReceiverService) GetReceiverByUserID(ctx context.Context, userID uint) (*model.Receiver, error) {
	var receiver model.Receiver
	err := s.db.WithContext(ctx).Where("user_id = ?", userID).First(&receiver).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &receiver, nil
}

// CreateReceiverFromUser 从用户创建接收人
func (s *ReceiverService) CreateReceiverFromUser(ctx context.Context, user *model.User) (*model.Receiver, error) {
	// 检查是否已存在
	existing, err := s.GetReceiverByUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}

	// 创建新接收人
	receiver := &model.Receiver{
		Name:     user.Name, // Using Name instead of DisplayName
		Email:    user.Email,
		Phone:    user.Phone,
		UserID:   &user.ID,
		DomainID: user.DomainID, // Set the domain from the user
		Enabled:  true,
	}

	if err := s.CreateReceiver(ctx, receiver); err != nil {
		return nil, err
	}

	return receiver, nil
}

// GetNotificationChannelsByReceiver 获取接收人的通知渠道
func (s *ReceiverService) GetNotificationChannelsByReceiver(ctx context.Context, receiverID uint) ([]*model.NotificationChannel, error) {
	var channels []*model.NotificationChannel

	err := s.db.WithContext(ctx).
		Model(&model.NotificationChannel{}).
		Joins("JOIN receiver_channels ON notification_channels.id = receiver_channels.notification_channel_id").
		Where("receiver_channels.receiver_id = ? AND receiver_channels.enabled = ?", receiverID, true).
		Find(&channels).Error

	return channels, err
}

// SetNotificationChannelsForReceiver 设置接收人的通知渠道
func (s *ReceiverService) SetNotificationChannelsForReceiver(ctx context.Context, receiverID uint, channelIDs []uint) error {
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 先删除现有的关联
	if err := tx.Where("receiver_id = ?", receiverID).Delete(&model.ReceiverChannel{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 添加新的关联
	for _, channelID := range channelIDs {
		rc := &model.ReceiverChannel{
			ReceiverID:          receiverID,
			NotificationChannelID: channelID,
			Enabled:             true,
		}
		if err := tx.Create(rc).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// GetReceiverWithChannels 获取接收人及其通知渠道
func (s *ReceiverService) GetReceiverWithChannels(ctx context.Context, id uint) (*model.Receiver, error) {
	var receiver model.Receiver
	err := s.db.WithContext(ctx).
		Preload("NotificationChannels").
		First(&receiver, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &receiver, nil
}