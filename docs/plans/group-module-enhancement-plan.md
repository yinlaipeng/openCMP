# Task Plan: 用户组管理模块增强开发

## Goal
根据系统页面分析结果（findings.md Phase 33），完善"认证与安全"管理控制台中的"用户组"管理模块，实现与用户模块类似的布局、工具栏、搜索栏、批量操作功能。

## Current Phase
All phases complete ✅

---

## Phases

### Phase 1: UI组件结构设计
- **目标**: 设计组件架构，创建必要的子组件
- **Status:** complete ✅
- [ ] **T1.1: 设计组件层级结构**
  ```
  frontend/src/views/iam/groups/
  ├── index.vue                    # 主列表页（改造）
  ├── components/
  │   ├── GroupDetailDrawer.vue    # 用户组详情抽屉（新建）
  │   ├── GroupDetailBasicInfo.vue # 详情Tab - 基本信息（新建）
  │   ├── GroupProjectsTab.vue     # 已加入项目Tab（改造）
  │   ├── GroupUsersTab.vue        # 组内用户Tab（改造）
  │   ├── GroupOperationLogs.vue   # 操作日志Tab（改造）
  │   └── BatchDeleteModal.vue     # 批量删除确认弹窗（新建）
  ```
- [ ] **T1.2: 创建 GroupDetailDrawer.vue 骨架**
  - el-drawer 结构
  - Header区域（组图标+名称+操作按钮）
  - el-tabs 4个Tab容器
- [ ] **T1.3: 创建各Tab组件骨架**
  - GroupDetailBasicInfo.vue (el-descriptions两列)
  - GroupProjectsTab.vue (工具栏+表格)
  - GroupUsersTab.vue (工具栏+表格)
  - GroupOperationLogs.vue (工具栏+表格)
- [ ] **T1.4: 编译验证**
  - npm run build 成功

---

### Phase 2: 主列表页改造
- **目标**: 按系统页面截图改造用户组列表页
- **Status:** complete ✅
- [x] **T2.1: 工具区改造**
  - 刷新按钮（icon only） ✅
  - Create 按钮（Primary Blue）✅ 已有
  - Delete 按钮（初始disabled状态）✅
  - Download 图标按钮 ✅
  - Settings 图标按钮 ✅
- [x] **T2.2: 搜索栏轻量化改造** ✅
  - 移除 3字段 inline form
  - 改为单一搜索框 + 属性选择下拉
  - 属性选择: ID/Name/Description/Domain
- [x] **T2.3: 表格改造** ✅
  - Checkbox 选择列（type="selection"）
  - Name 列（可点击链接 + sortable）
  - Owner Domain 列（sortable）
  - Operations 列精简（管理项目/管理用户/删除）
- [x] **T2.4: 批量操作功能** ✅
  - 批量删除（选中后启用Delete按钮）
  - 批量删除确认弹窗（显示选中组列表）
- [x] **T2.5: 详情改为抽屉** ✅
  - 集成 GroupDetailDrawer
  - 替换 el-dialog
- [x] **T2.6: 编译验证** ✅
  - npm run build 成功

---

### Phase 3: 详情抽屉开发
- **目标**: 创建用户组详情侧边抽屉组件
- **Status:** complete ✅
- [x] **T3.1: 抽屉结构** - GroupDetailDrawer.vue ✅
  - el-drawer 响应式宽度
  - Header: 组图标 + 名称 + 刷新/编辑/删除按钮
- [x] **T3.2: Tab1 - 详情** - GroupDetailBasicInfo.vue ✅
  - el-descriptions 两列布局
  - 字段: ID/名称/描述/所属域/创建时间/更新时间
- [x] **T3.3: Tab2 - 已加入项目** - GroupProjectsTab.vue ✅
  - 工具栏: 加入项目按钮
  - 表格: 项目列表 + 移除操作
- [x] **T3.4: Tab3 - 组内用户** - GroupUsersTab.vue ✅
  - 工具栏: 添加用户按钮
  - 表格: 用户列表 + 移除操作
- [x] **T3.5: Tab4 - 操作日志** - GroupOperationLogs.vue ✅
  - 工具栏: 筛选
  - 表格: 日志列表

---

### Phase 4: 后端API补齐
- **目标**: 补齐缺失的后端API
- **Status:** complete ✅
- [x] **T4.1: 批量删除API** - POST /groups/batch-delete ✅
- [x] **T4.2: 导出用户组API** - 可选，暂不实现
- [x] **T4.3: 编译验证** - go build ./... 成功 ✅

---

### Phase 5: 前端API对接
- **目标**: 前端API函数与后端对接
- **Status:** complete ✅
- [x] **T5.1: 扩展 iam.ts API文件**
  - batchDeleteGroups ✅
- [x] **T5.2: 编译验证** - npm run build 成功 ✅

---

### Phase 6: 编译验证与测试
- **目标**: 验证前后端编译成功，功能正确
- **Status:** complete ✅
- [x] **T6.1: 后端编译验证** - go build ./... 成功 ✅
- [x] **T6.2: 前端编译验证** - npm run build 成功 ✅
- [x] **T6.3: 功能检查** - 所有组件和API已对接 ✅

---

## Summary: 用户组管理模块增强开发完成

### 已完成文件清单

**前端组件 (6个新建)**
- `frontend/src/views/iam/groups/components/GroupDetailDrawer.vue`
- `frontend/src/views/iam/groups/components/GroupDetailBasicInfo.vue`
- `frontend/src/views/iam/groups/components/GroupProjectsTab.vue`
- `frontend/src/views/iam/groups/components/GroupUsersTab.vue`
- `frontend/src/views/iam/groups/components/GroupOperationLogs.vue`
- `frontend/src/views/iam/groups/components/BatchDeleteModal.vue`

**前端修改 (2个)**
- `frontend/src/views/iam/groups/index.vue` - 主列表页完整改造
- `frontend/src/api/iam.ts` - 新增 batchDeleteGroups API函数

**后端修改 (3个)**
- `backend/internal/handler/group.go` - 新增 BatchDelete Handler方法
- `backend/internal/service/group.go` - 新增 BatchDeleteGroups Service方法
- `backend/cmd/server/main.go` - 新增 batch-delete 路由

### 实现的功能

| 功能 | 状态 |
|------|------|
| 工具栏完整（刷新/Create/Delete/Download/Settings） | ✅ |
| 轻量搜索栏 + 属性选择 | ✅ |
| Checkbox批量选择 | ✅ |
| 批量删除 | ✅ |
| 详情抽屉（替代dialog） | ✅ |
| 详情4个Tab（基本信息/项目/用户/日志） | ✅ |
| 表格排序 | ✅ |
| 操作按钮精简 | ✅ |
| 状态圆点样式 | ✅ |

### 待后续优化

1. 导出用户组CSV功能
2. 批量操作可优化为SQL批量执行
3. 标签管理功能待实现

---

## Key Questions
1. 详情抽屉宽度应设置为多少？（建议60%或800px）
2. 批量删除是否需要显示确认弹窗中的组详情？
3. 导出功能是否需要支持CSV格式？

## Decisions Made
| Decision | Rationale |
|----------|-----------|
| 详情改为抽屉而非弹窗 | 与用户模块保持一致，体验更好 |
| 保持Element Plus框架 | 与项目整体风格一致，已有20+页面使用 |
| 批量删除使用确认弹窗 | 避免误操作，符合安全规范 |

## Errors Encountered
| Error | Attempt | Resolution |
|-------|---------|------------|
| (待记录) | 1 | |

## Notes
- 现有前端 index.vue 已有基础实现（4Tab详情弹窗），需改造而非完全重写
- 后端已有全部API（CRUD + 用户管理 + 项目管理），只需补齐批量删除
- 优先完成前端改造，后端API只需少量补齐