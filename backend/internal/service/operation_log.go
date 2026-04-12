package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/opencmp/opencmp/internal/model"
	"github.com/opencmp/opencmp/pkg/utils/pagination"
	"gorm.io/gorm"
)

type OperationLogService struct {
	DB *gorm.DB
}

// NewOperationLogService creates a new instance of OperationLogService
func NewOperationLogService(db *gorm.DB) *OperationLogService {
	return &OperationLogService{
		DB: db,
	}
}

// CreateOperationLog creates a new operation log entry
func (s *OperationLogService) CreateOperationLog(log *model.OperationLog) error {
	if log == nil {
		return errors.New("operation log cannot be nil")
	}

	// Set default values if not provided
	if log.OperationTime.IsZero() {
		log.OperationTime = time.Now()
	}
	if log.RiskLevel == "" {
		log.RiskLevel = "medium"
	}
	if log.TimeType == "" {
		log.TimeType = "realtime"
	}
	if log.Result == "" {
		log.Result = "success"
	}

	return s.DB.Create(log).Error
}

// GetOperationLogs retrieves operation logs with pagination and filtering
func (s *OperationLogService) GetOperationLogs(filter map[string]interface{}, pg *pagination.Pagination) (*pagination.Pagination, error) {
	var logs []model.OperationLog

	query := s.DB.Model(&model.OperationLog{})

	// Apply filters
	for k, v := range filter {
		switch k {
		case "resource_name":
			query = query.Where("resource_name LIKE ?", fmt.Sprintf("%%%s%%", v))
		case "resource_type":
			query = query.Where("resource_type = ?", v)
		case "operation_type":
			query = query.Where("operation_type = ?", v)
		case "service_type":
			query = query.Where("service_type = ?", v)
		case "risk_level":
			query = query.Where("risk_level = ?", v)
		case "result":
			query = query.Where("result = ?", v)
		case "operator":
			query = query.Where("operator LIKE ?", fmt.Sprintf("%%%s%%", v))
		case "project_id":
			query = query.Where("project_id = ?", v)
		case "domain_id":
			query = query.Where("domain_id = ?", v)
		case "user_id":
			query = query.Where("user_id = ?", v)
		}
	}

	// Apply pagination
	total := int64(0)
	query.Count(&total)
	pg.Total = total

	result := query.Offset(int(pg.GetOffset())).Limit(int(pg.Limit)).Order("operation_time DESC").Find(&logs)

	if result.Error != nil {
		return nil, result.Error
	}

	pg.Items = logs
	return pg, nil
}

// GetOperationLogByID retrieves an operation log by its ID
func (s *OperationLogService) GetOperationLogByID(id uint) (*model.OperationLog, error) {
	var log model.OperationLog
	err := s.DB.First(&log, id).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

// GetResourceOperationLogs retrieves operation logs for a specific resource
func (s *OperationLogService) GetResourceOperationLogs(resourceType string, resourceID uint, pg *pagination.Pagination) (*pagination.Pagination, error) {
	var logs []model.OperationLog

	query := s.DB.Where("resource_type = ? AND resource_id = ?", resourceType, resourceID).Model(&model.OperationLog{})

	// Apply pagination
	total := int64(0)
	query.Count(&total)
	pg.Total = total

	result := query.Offset(int(pg.GetOffset())).Limit(int(pg.Limit)).Order("operation_time DESC").Find(&logs)

	if result.Error != nil {
		return nil, result.Error
	}

	pg.Items = logs
	return pg, nil
}
