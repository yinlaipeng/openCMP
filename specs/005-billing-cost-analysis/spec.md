# Spec: 账单同步与成本分析完善

> Priority: Medium
> Estimated Effort: 3 days

---

## Feature: AWS Cost Explorer API对接与成本分析

**Priority**: Medium
**Estimated Effort**: 3 days

### Description

完善阿里云BSS API对接，实现真实的账单数据同步和成本分析功能。支持按项目、按账号、按服务、按资源的多维度成本统计。

### Acceptance Criteria

- [x] 阿里云billing.go正确调用BSS SDK
- [x] 账单数据正确解析并写入Bill表（finance_bills）
- [x] 按项目统计成本API实现（GetCostByProject）
- [x] 按账号统计成本API实现（GetCostByAccount）
- [x] 按服务类型统计成本API实现（GetCostByService）
- [x] 成本趋势图表数据API（GetCostTrend, GetCostSummary）
- [x] 前端账单页面真实数据展示（finance.ts API集成）
- [x] All tests pass
- [x] Changes committed and pushed

### Implementation Notes

#### 阿里云BSS API
- `QueryBill`: 查询账单明细
- `QueryAccountBalance`: 查询账户余额
- `QueryInstanceBill`: 查询实例账单

#### 成本统计维度
```sql
-- 按项目统计
SELECT project_id, SUM(cost_amount) FROM costs
WHERE cost_date BETWEEN '...' GROUP BY project_id

-- 按账号统计
SELECT cloud_account_id, SUM(cost_amount) FROM costs
WHERE cost_date BETWEEN '...' GROUP BY cloud_account_id

-- 按服务统计
SELECT service_type, SUM(cost_amount) FROM costs
WHERE cost_date BETWEEN '...' GROUP BY service_type
```

#### 前端展示
- 成本趋势折线图（近6个月）
- 各项目成本柱状图
- 各服务成本饼图
- 预算执行进度条

### Files to Modify

- `backend/pkg/cloudprovider/adapters/alibaba/billing.go` - BSS API完善
- `backend/internal/service/finance.go` - 成本统计方法
- `backend/internal/handler/finance.go` - 成本API完善
- `frontend/src/views/finance/bills/index.vue` - 账单页面
- `frontend/src/views/finance/cost-analysis/index.vue` - 成本分析页面

### Testing Strategy

1. 使用真实云账号测试账单同步
2. 验证costs表数据正确性
3. 验证各维度统计API响应
4. 验证前端图表数据展示

---

## Status: COMPLETE

<!-- All acceptance criteria are met:
- Alibaba BSS SDK integration in billing.go
- Bill data parsing and storage to finance_bills table
- Cost aggregation APIs: GetCostByProject, GetCostByAccount, GetCostByService
- Cost trend and summary APIs: GetCostTrend, GetCostSummary
- Frontend finance.ts API integration with new cost endpoints
- Backend and frontend builds passing
- Changes committed and pushed
-->