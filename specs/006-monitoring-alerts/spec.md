# Spec: 监控告警模块完善

> Priority: Low
> Estimated Effort: 2 days

---

## Feature: 监控数据采集与告警通知

**Priority**: Low
**Estimated Effort**: 2 days

### Description

完善监控模块的实时数据采集和告警通知功能，支持VM监控指标查询、告警策略配置、多渠道通知发送。

### Acceptance Criteria

- [ ] VM监控指标数据采集API实现（当前返回模拟数据，需集成阿里云云监控SDK）
- [x] 告警策略CRUD API完整实现
- [x] 告警历史记录查询API实现
- [x] 通知机器人测试发送功能（钉钉/飞书/企微/Webhook真实发送测试消息）
- [x] 消息订阅规则配置实现
- [ ] 前端监控页面真实数据展示
- [x] All tests pass
- [ ] Changes committed and pushed

### Implementation Notes

#### 监控指标
- CPU使用率
- 内存使用率
- 网络流量
- 磁盘IO
- 进程状态

#### 告警渠道
- 钉钉机器人（Webhook + 签名）
- 企业微信机器人
- Slack Webhook
- 邮件通知

#### 告警规则
```json
{
  "metric": "cpu_usage",
  "threshold": 80,
  "operator": ">",
  "duration": 5,  // 持续5分钟
  "notify_channels": ["dingtalk", "email"]
}
```

### Files to Modify

- `backend/internal/service/monitor.go` - 监控服务完善
- `backend/internal/handler/monitor.go` - 监控API完善
- `backend/internal/service/notification.go` - 通知服务
- `frontend/src/views/monitoring/resources/vms/index.vue` - VM监控页面
- `frontend/src/views/monitoring/alerts/policies/index.vue` - 告警策略页面

### Testing Strategy

1. 验证监控指标采集
2. 创建告警策略测试触发
3. 测试通知发送（钉钉/企微）
4. 验证告警历史记录

---

## Status: PENDING

<!-- Change to COMPLETE when all acceptance criteria are met -->