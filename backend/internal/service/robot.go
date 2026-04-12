package service

import (
	"context"
	"encoding/json"

	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
)

// RobotService 机器人服务
type RobotService struct {
	db *gorm.DB
}

// NewRobotService 创建机器人服务
func NewRobotService(db *gorm.DB) *RobotService {
	return &RobotService{db: db}
}

// CreateRobot 创建机器人
func (s *RobotService) CreateRobot(ctx context.Context, robot *model.Robot) error {
	return s.db.WithContext(ctx).Create(robot).Error
}

// GetRobot 获取机器人
func (s *RobotService) GetRobot(ctx context.Context, id uint) (*model.Robot, error) {
	var robot model.Robot
	err := s.db.WithContext(ctx).First(&robot, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &robot, nil
}

// GetRobotByName 根据名称获取机器人
func (s *RobotService) GetRobotByName(ctx context.Context, name string) (*model.Robot, error) {
	var robot model.Robot
	err := s.db.WithContext(ctx).Where("name = ?", name).First(&robot).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &robot, nil
}

// ListRobots 列出机器人
func (s *RobotService) ListRobots(ctx context.Context, robotType string, limit, offset int) ([]*model.Robot, int64, error) {
	var robots []*model.Robot
	var total int64

	query := s.db.Model(&model.Robot{})
	if robotType != "" {
		query = query.Where("type = ?", robotType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&robots).Error

	return robots, total, err
}

// UpdateRobot 更新机器人
func (s *RobotService) UpdateRobot(ctx context.Context, robot *model.Robot) error {
	return s.db.WithContext(ctx).Save(robot).Error
}

// DeleteRobot 删除机器人
func (s *RobotService) DeleteRobot(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Delete(&model.Robot{}, id).Error
}

// EnableRobot 启用机器人
func (s *RobotService) EnableRobot(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Model(&model.Robot{}).Where("id = ?", id).Update("enabled", true).Error
}

// DisableRobot 禁用机器人
func (s *RobotService) DisableRobot(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Model(&model.Robot{}).Where("id = ?", id).Update("enabled", false).Error
}

// GetMessageTypes 获取机器人订阅的消息类型
func (s *RobotService) GetMessageTypes(ctx context.Context, robotID uint) ([]string, error) {
	var robot model.Robot
	if err := s.db.WithContext(ctx).First(&robot, robotID).Error; err != nil {
		return nil, err
	}

	if robot.MessageTypes == nil {
		return []string{}, nil
	}

	var types []string
	if err := json.Unmarshal(robot.MessageTypes, &types); err != nil {
		return nil, err
	}

	return types, nil
}

// SetMessageTypes 设置机器人订阅的消息类型
func (s *RobotService) SetMessageTypes(ctx context.Context, robotID uint, types []string) error {
	typesJSON, err := json.Marshal(types)
	if err != nil {
		return err
	}

	return s.db.WithContext(ctx).Model(&model.Robot{}).Where("id = ?", robotID).Update("message_types", typesJSON).Error
}
