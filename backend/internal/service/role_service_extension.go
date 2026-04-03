package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/opencmp/opencmp/internal/model"
	"gorm.io/gorm"
)

// GetPermissionByName 根据名称获取权限
func (s *RoleService) GetPermissionByName(ctx context.Context, name string) (*model.Permission, error) {
	var permission model.Permission
	err := s.db.WithContext(ctx).Where("name = ?", name).First(&permission).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("permission with name '%s' not found", name)
		}
		return nil, err
	}
	return &permission, nil
}
