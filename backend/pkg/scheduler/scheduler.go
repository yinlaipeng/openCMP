package scheduler

import (
	"context"
	"sync"
	"time"

	"github.com/opencmp/opencmp/internal/model"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Scheduler 定时任务调度器
type Scheduler struct {
	cron    *cron.Cron
	db      *gorm.DB
	runner  *TaskRunner
	logger  *zap.Logger
	mu      sync.RWMutex
	entries map[uint]cron.EntryID // task_id -> entry_id
}

// NewScheduler 创建调度器
func NewScheduler(db *gorm.DB, logger *zap.Logger) *Scheduler {
	// 使用标准cron解析器，支持秒级精度可选
	c := cron.New(cron.WithSeconds(), cron.WithLocation(time.Local))

	return &Scheduler{
		cron:    c,
		db:      db,
		runner:  NewTaskRunner(db, logger),
		logger:  logger,
		entries: make(map[uint]cron.EntryID),
	}
}

// Start 启动调度器
func (s *Scheduler) Start() {
	s.cron.Start()
	s.logger.Info("scheduler started")

	// 加载所有active状态的定时任务
	s.LoadActiveTasks()
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	s.cron.Stop()
	s.logger.Info("scheduler stopped")
}

// LoadActiveTasks 加载所有活跃的定时任务
func (s *Scheduler) LoadActiveTasks() {
	ctx := context.Background()
	var tasks []*model.ScheduledTask

	// 查询所有active状态的任务
	if err := s.db.WithContext(ctx).
		Where("status = ?", model.ScheduledTaskStatusActive).
		Find(&tasks).Error; err != nil {
		s.logger.Error("failed to load active tasks", zap.Error(err))
		return
	}

	for _, task := range tasks {
		if err := s.AddTask(task); err != nil {
			s.logger.Error("failed to add task to scheduler",
				zap.Uint("task_id", task.ID),
				zap.String("task_name", task.Name),
				zap.Error(err))
		} else {
			s.logger.Info("task added to scheduler",
				zap.Uint("task_id", task.ID),
				zap.String("task_name", task.Name),
				zap.String("frequency", task.Frequency),
				zap.String("trigger_time", task.TriggerTime))
		}
	}

	s.logger.Info("active tasks loaded", zap.Int("count", len(tasks)))
}

// AddTask 添加任务到调度器
func (s *Scheduler) AddTask(task *model.ScheduledTask) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查任务是否已在调度器中
	if _, exists := s.entries[task.ID]; exists {
		return nil // 已存在，跳过
	}

	// 解析cron表达式
	cronExpr := s.parseFrequencyToCron(task.Frequency, task.TriggerTime)

	// 添加任务
	entryID, err := s.cron.AddFunc(cronExpr, func() {
		s.runner.Run(task.ID, task.Type, task.CloudAccountID)
	})
	if err != nil {
		return err
	}

	s.entries[task.ID] = entryID
	return nil
}

// RemoveTask 从调度器移除任务
func (s *Scheduler) RemoveTask(taskID uint) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if entryID, exists := s.entries[taskID]; exists {
		s.cron.Remove(entryID)
		delete(s.entries, taskID)
		s.logger.Info("task removed from scheduler", zap.Uint("task_id", taskID))
	}
}

// UpdateTask 更新任务（先移除再添加）
func (s *Scheduler) UpdateTask(task *model.ScheduledTask) error {
	s.RemoveTask(task.ID)

	if task.Status == string(model.ScheduledTaskStatusActive) {
		return s.AddTask(task)
	}
	return nil
}

// parseFrequencyToCron 将频率转换为cron表达式
// 支持的频率格式：
// - once: 单次执行（需要TriggerTime为具体时间）
// - daily: 每天执行
// - weekly: 每周执行
// - monthly: 每月执行
// - custom: 自定义cron表达式
// - hourly: 每小时执行
func (s *Scheduler) parseFrequencyToCron(frequency, triggerTime string) string {
	// triggerTime 格式可能是 "HH:MM" 或完整cron表达式

	switch frequency {
	case "once":
		// 单次执行：解析triggerTime作为具体执行时间
		// triggerTime 格式: "2026-04-15 10:00" 或 "HH:MM"
		// 对于单次任务，我们使用 @at 语法（cron库不支持，改为延迟执行）
		// 这里简化处理，使用触发时间的时分作为每天执行的时间点
		// 实际单次任务应该在执行后自动禁用
		hour, min := s.parseTime(triggerTime)
		return s.formatCron(hour, min, "*", "*", "*")

	case "hourly":
		// 每小时执行：0 * * * *
		return "0 * * * *"

	case "daily":
		// 每天执行：triggerTime格式为 "HH:MM"
		hour, min := s.parseTime(triggerTime)
		return s.formatCron(hour, min, "*", "*", "*")

	case "weekly":
		// 每周执行：triggerTime格式为 "HH:MM day" (day为周一=1到周日=7)
		hour, min, day := s.parseWeeklyTime(triggerTime)
		return s.formatCron(hour, min, "*", "*", day)

	case "monthly":
		// 每月执行：triggerTime格式为 "HH:MM day" (day为每月几号)
		hour, min, day := s.parseMonthlyTime(triggerTime)
		return s.formatCron(hour, min, day, "*", "*")

	case "custom":
		// 自定义cron表达式，直接使用triggerTime
		return triggerTime

	default:
		// 默认每小时执行
		return "0 * * * *"
	}
}

// parseTime 解析时间 "HH:MM"
func (s *Scheduler) parseTime(timeStr string) (int, int) {
	// 简化解析，默认值
	if timeStr == "" {
		return 0, 0
	}

	// 支持多种格式
	// "10:00" -> hour=10, min=0
	// "10:30" -> hour=10, min=30

	parts := splitTime(timeStr)
	if len(parts) >= 2 {
		hour := parseInt(parts[0], 0)
		min := parseInt(parts[1], 0)
		return hour, min
	}

	return 0, 0
}

// parseWeeklyTime 解析每周时间 "HH:MM day"
func (s *Scheduler) parseWeeklyTime(timeStr string) (int, int, string) {
	parts := splitTime(timeStr)
	if len(parts) >= 3 {
		hour := parseInt(parts[0], 0)
		min := parseInt(parts[1], 0)
		day := parts[2] // day可能是数字或星期名称
		return hour, min, convertWeekDay(day)
	}
	return 0, 0, "*"
}

// parseMonthlyTime 解析每月时间 "HH:MM day"
func (s *Scheduler) parseMonthlyTime(timeStr string) (int, int, string) {
	parts := splitTime(timeStr)
	if len(parts) >= 3 {
		hour := parseInt(parts[0], 0)
		min := parseInt(parts[1], 0)
		day := parts[2] // day是每月几号
		return hour, min, day
	}
	return 0, 0, "1"
}

// formatCron 格式化cron表达式
func (s *Scheduler) formatCron(hour, min int, dom, month, dow string) string {
	// cron格式: 秒 分 时 日 月 周
	// 使用 WithSeconds 的cron需要6个字段
	return formatInt(min) + " " + formatInt(hour) + " " + dom + " " + month + " " + dow + " *"
}

// 辅助函数
func splitTime(timeStr string) []string {
	result := []string{}
	for _, sep := range []string{" ", ":"} {
		if parts := splitBy(timeStr, sep); len(parts) > 1 {
			result = parts
			break
		}
	}
	if len(result) == 0 && timeStr != "" {
		result = []string{timeStr}
	}
	return result
}

func splitBy(s, sep string) []string {
	if s == "" {
		return nil
	}
	for i := 0; i < len(s); i++ {
		if s[i:i+1] == sep {
			return []string{s[:i], s[i+1:]}
		}
	}
	return nil
}

func parseInt(s string, defaultVal int) int {
	val := defaultVal
	for _, c := range s {
		if c >= '0' && c <= '9' {
			val = val*10 + int(c-'0')
		} else {
			break
		}
	}
	return val
}

func formatInt(i int) string {
	if i < 10 {
		return "0" + intToStr(i)
	}
	return intToStr(i)
}

func intToStr(i int) string {
	if i == 0 {
		return "0"
	}
	result := ""
	for i > 0 {
		result = string('0'+i%10) + result
		i /= 10
	}
	return result
}

func convertWeekDay(day string) string {
	// 转换星期格式
	// 输入可能是: 1-7 (周一到周日), mon-sun, monday-sunday
	dayMap := map[string]string{
		"1":       "MON",
		"2":       "TUE",
		"3":       "WED",
		"4":       "THU",
		"5":       "FRI",
		"6":       "SAT",
		"7":       "SUN",
		"0":       "SUN",
		"mon":     "MON",
		"tue":     "TUE",
		"wed":     "WED",
		"thu":     "THU",
		"fri":     "FRI",
		"sat":     "SAT",
		"sun":     "SUN",
		"monday":  "MON",
		"tuesday": "TUE",
		"wednesday": "WED",
		"thursday": "THU",
		"friday":  "FRI",
		"saturday": "SAT",
		"sunday":  "SUN",
	}

	if mapped, ok := dayMap[day]; ok {
		return mapped
	}
	return "*"
}

// GetScheduledTaskIDs 获取已调度的任务ID列表
func (s *Scheduler) GetScheduledTaskIDs() []uint {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ids := make([]uint, 0, len(s.entries))
	for id := range s.entries {
		ids = append(ids, id)
	}
	return ids
}

// GetNextRunTime 获取任务的下次执行时间
func (s *Scheduler) GetNextRunTime(taskID uint) time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if entryID, exists := s.entries[taskID]; exists {
		entry := s.cron.Entry(entryID)
		return entry.Next
	}
	return time.Time{}
}