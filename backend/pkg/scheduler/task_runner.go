package scheduler

import (
	"context"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
	"github.com/opencmp/opencmp/internal/service"
)

// TaskRunner 任务执行器
type TaskRunner struct {
	db             *gorm.DB
	accountService *service.CloudAccountService
	taskService    *service.ScheduledTaskService
	logger         *zap.Logger
}

// NewTaskRunner 创建任务执行器
func NewTaskRunner(db *gorm.DB, logger *zap.Logger) *TaskRunner {
	return &TaskRunner{
		db:             db,
		accountService: service.NewCloudAccountService(db),
		taskService:    service.NewScheduledTaskService(db),
		logger:         logger,
	}
}

// Run 执行任务
func (r *TaskRunner) Run(taskID uint, taskType string, cloudAccountID *uint) {
	ctx := context.Background()
	startTime := time.Now()

	r.logger.Info("task execution started",
		zap.Uint("task_id", taskID),
		zap.String("task_type", taskType),
		zap.Uintp("cloud_account_id", cloudAccountID))

	// 更新任务执行状态（可选：添加last_run_time字段）
	// 这里我们先记录日志，不修改数据库状态

	var err error
	switch taskType {
	case "sync_cloud_account":
		err = r.runSyncCloudAccount(ctx, cloudAccountID)
	case "sync_billing":
		err = r.runSyncBilling(ctx, cloudAccountID)
	case "sync_renewals":
		err = r.runSyncRenewals(ctx, cloudAccountID)
	case "full_sync":
		err = r.runFullSync(ctx, cloudAccountID)
	default:
		r.logger.Warn("unknown task type", zap.String("task_type", taskType))
		return
	}

	duration := time.Since(startTime)

	if err != nil {
		r.logger.Error("task execution failed",
			zap.Uint("task_id", taskID),
			zap.String("task_type", taskType),
			zap.Duration("duration", duration),
			zap.Error(err))
	} else {
		r.logger.Info("task execution completed",
			zap.Uint("task_id", taskID),
			zap.String("task_type", taskType),
			zap.Duration("duration", duration))
	}

	// 更新任务的最后执行时间（可选）
	r.updateTaskLastRun(ctx, taskID, startTime, err)
}

// runSyncCloudAccount 执行云账号增量同步
func (r *TaskRunner) runSyncCloudAccount(ctx context.Context, cloudAccountID *uint) error {
	if cloudAccountID == nil {
		return r.syncAllActiveAccounts(ctx, "incremental")
	}

	account, err := r.accountService.GetCloudAccount(ctx, *cloudAccountID)
	if err != nil {
		return err
	}
	if account == nil {
		return gorm.ErrRecordNotFound
	}

	_, err = r.accountService.SyncResources(ctx, account)
	return err
}

// runFullSync 执行云账号全量同步
func (r *TaskRunner) runFullSync(ctx context.Context, cloudAccountID *uint) error {
	if cloudAccountID == nil {
		return r.syncAllActiveAccounts(ctx, "full")
	}

	account, err := r.accountService.GetCloudAccount(ctx, *cloudAccountID)
	if err != nil {
		return err
	}
	if account == nil {
		return gorm.ErrRecordNotFound
	}

	// 全量同步需要在SyncResources中实现差异逻辑
	// 这里暂时使用相同的SyncResources方法
	_, err = r.accountService.SyncResources(ctx, account)
	return err
}

// runSyncBilling 执行账单同步
func (r *TaskRunner) runSyncBilling(ctx context.Context, cloudAccountID *uint) error {
	// TODO: 实现账单同步逻辑
	// 需要调用FinanceHandler的SyncBills方法
	r.logger.Info("sync billing task executed", zap.Uintp("cloud_account_id", cloudAccountID))
	return nil
}

// runSyncRenewals 执行续费资源同步
func (r *TaskRunner) runSyncRenewals(ctx context.Context, cloudAccountID *uint) error {
	// TODO: 实现续费资源同步逻辑
	// 需要调用FinanceHandler的SyncRenewals方法
	r.logger.Info("sync renewals task executed", zap.Uintp("cloud_account_id", cloudAccountID))
	return nil
}

// syncAllActiveAccounts 同步所有活跃的云账号
func (r *TaskRunner) syncAllActiveAccounts(ctx context.Context, mode string) error {
	var accounts []*model.CloudAccount

	// 查询所有enabled状态的云账号
	if err := r.db.WithContext(ctx).
		Where("enabled = ?", true).
		Find(&accounts).Error; err != nil {
		return err
	}

	r.logger.Info("syncing all active accounts",
		zap.Int("count", len(accounts)),
		zap.String("mode", mode))

	for _, account := range accounts {
		_, err := r.accountService.SyncResources(ctx, account)
		if err != nil {
			r.logger.Error("failed to sync account",
				zap.Uint("account_id", account.ID),
				zap.String("account_name", account.Name),
				zap.Error(err))
			// 继续同步其他账号，不中断
		}
	}

	return nil
}

// updateTaskLastRun 更新任务的最后执行时间
func (r *TaskRunner) updateTaskLastRun(ctx context.Context, taskID uint, startTime time.Time, execErr error) {
	// 更新任务状态（可选字段）
	// 注意：ScheduledTask模型可能需要添加last_run_time字段

	// 这里只记录日志，实际更新需要扩展模型
	status := "success"
	if execErr != nil {
		status = "failed"
	}

	r.logger.Info("task run record",
		zap.Uint("task_id", taskID),
		zap.Time("start_time", startTime),
		zap.String("status", status))

	// 可选：写入操作日志
	// operationLog := &model.OperationLog{
	//     ResourceType: "scheduled_task",
	//     ResourceID:   taskID,
	//     Action:       "auto_execute",
	//     Status:       status,
	//     ...
	// }
	// r.db.Create(operationLog)
}