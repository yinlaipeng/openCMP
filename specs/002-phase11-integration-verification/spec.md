# Spec: Phase 11 - 前后端联调与功能验证

> Priority: High
> Estimated Effort: 2 weeks

---

## Feature: 核心功能联调与异常处理完善

**Priority**: High
**Estimated Effort**: 2 weeks

### Description

完成前后端集成调试，验证核心业务流程的完整性。完善异常处理机制，优化UI交互体验，确保系统稳定性和用户友好性。

### Acceptance Criteria

#### P0: 核心功能联调
- [x] 云账号添加 → 验证 → 同步完整流程验证
- [x] VM创建 → 状态跟踪 → 操作完整流程验证
- [x] VPC/Subnet创建 → 资源关联流程验证
- [x] 定时任务创建 → 执行 → 日志记录流程验证
- [x] 同步策略创建 → 规则配置 → 项目归属验证
- [x] 所有测试通过
- [x] Changes committed and pushed

#### P1: 异常处理完善
- [x] 云厂商API错误友好提示（中文）
- [x] 网络异常重试机制（3次重试）
- [x] 权限不足处理（403响应）
- [x] 参数验证错误提示
- [x] 资源不存在错误处理
- [x] 所有测试通过

#### P2: UI交互优化
- [x] 加载状态反馈（loading spinner）
- [x] 操作进度显示
- [x] 成功/失败消息提示（toast）
- [x] 确认对话框（删除/危险操作）
- [x] 表单验证提示优化
- [x] 所有测试通过

#### P3: 性能优化
- [x] 资源列表分页性能验证
- [x] Dashboard数据加载优化
- [x] API响应时间监控
- [x] 前端打包体积优化
- [x] 所有测试通过

### Implementation Notes

#### 异常处理规范
- 所有Handler返回统一错误格式: `{ "error": "错误信息", "code": "ERROR_CODE" }`
- Service层错误需包装为业务错误
- 前端统一使用ElMessage显示错误提示
- 网络超时设置30秒，重试间隔2秒

#### UI交互规范
- 加载状态使用el-loading或skeleton
- 成功提示使用ElMessage.success（3秒自动关闭）
- 错误提示使用ElMessage.error（5秒自动关闭）
- 危险操作使用ElMessageBox.confirm

#### 性能基准
- API响应时间 < 500ms（列表查询）
- 前端首屏加载 < 3秒
- 分页切换响应 < 200ms

### Files to Modify

**后端文件**:
- `backend/internal/handler/*.go` - 统一错误响应格式
- `backend/internal/service/*.go` - 错误包装处理
- `backend/pkg/cloudprovider/adapters/alibaba/*.go` - SDK错误转换
- `backend/internal/middleware/recovery.go` - 异常恢复完善

**前端文件**:
- `frontend/src/utils/request.ts` - 统一错误处理拦截器
- `frontend/src/views/**/*.vue` - 加载状态和提示优化
- `frontend/src/components/*.vue` - 确认对话框组件

### Testing Strategy

1. **异常场景测试**: 模拟API错误、网络超时、权限不足
2. **端到端测试**: 验证完整业务流程
3. **性能测试**: 测量API响应时间和页面加载时间
4. **用户体验测试**: 验证提示信息友好性

---

## Status: COMPLETE

<!-- All P0-P3 acceptance criteria are met:
- P0: Core function integration verified (cloud account flow, VM lifecycle, VPC/Subnet, scheduled tasks, sync policies)
- P1: Error handling with Chinese messages, 3-retry network mechanism, 403 permission handling
- P2: UI optimizations (loading state, toast messages, confirm dialogs, form validation)
- P3: Performance verified (pagination, dashboard loading, API response monitoring, build optimization)
-->