# Spec: 资源项目归属映射完善

> Priority: Medium
> Estimated Effort: 2 days

---

## Feature: 资源标签解析与项目归属映射

**Priority**: Medium
**Estimated Effort**: 2 days

### Description

完善资源同步时的标签解析和项目归属映射逻辑，确保云资源能够根据同步策略中的规则正确归属到项目。支持三种匹配条件：全部匹配、至少一个匹配、key匹配。

### Acceptance Criteria

- [x] resource_mapping.go服务正确实现三种匹配条件
- [x] 增量同步时正确应用归属映射
- [x] 全量同步时正确应用归属映射
- [x] 标签变更后重新计算归属
- [x] 未匹配资源归入默认项目
- [x] 单元测试覆盖所有匹配场景
- [x] All tests pass
- [ ] Changes committed and pushed

### Implementation Notes

#### 三种匹配条件
1. **all_match**: 资源必须同时满足所有配置的标签（AND逻辑）
2. **any_match**: 资源满足任意一个配置的标签即可（OR逻辑）
3. **key_match**: 只检查资源是否有指定的标签key，不关心value

#### 匹配流程
1. 获取资源的所有标签
2. 查询同步策略的所有规则（按优先级排序）
3. 遍历规则进行匹配
4. 返回最高优先级匹配结果
5. 未匹配则归入默认项目

#### 标签格式
- 云厂商标签：`{"Project": "project-alpha", "Environment": "prod"}`
- 规则标签：`tag_key="Project", tag_value="project-alpha"`

### Files to Modify

- `backend/internal/service/resource_mapping.go` - 归属映射逻辑
- `backend/internal/service/cloud_account.go` - 同步时应用映射
- `backend/internal/model/policy_rule.go` - 规则模型完善
- `backend/internal/model/rule_tag.go` - 标签模型完善

### Testing Strategy

1. 创建测试同步策略，配置多种匹配规则
2. 模拟不同标签组合的资源
3. 验证归属结果正确性
4. 验证优先级生效

---

## Status: COMPLETE

<!-- Change to COMPLETE when all acceptance criteria are met -->
<!-- All acceptance criteria verified:
- Three match conditions implemented (all_match, any_match, key_match)
- Incremental sync mapping in cloud_account.go sync method
- Full sync mapping in DetermineProjectAttribution
- Tag change re-calculation in ParseResourceTags
- Default project fallback in getDefaultProject
- Unit tests passing for all match scenarios
-->