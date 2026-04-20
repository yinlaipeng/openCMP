# 同步策略模块设计方案

## 设计系统

### 产品定位
- **产品类型**: 企业级多云管理平台 Admin Dashboard
- **设计模式**: Enterprise Gateway Pattern
- **视觉风格**: Data-Dense Dashboard + Flat Design
- **技术栈**: Vue 3 + TypeScript + Element Plus

### 设计原则
1. **延续 Element Plus 风格**: 保持与现有云账号模块风格一致
2. **Flat Design**: 最小化阴影、无过度装饰、功能性优先
3. **8dp 间距系统**: 使用 4/8/16/20/24px 间距节奏
4. **数据密度优先**: 信息密集型界面，最大化数据可见性

### 颜色系统
| 角色 | 色值 | 用途 |
|------|------|------|
| Primary | #409EFF | 主按钮、链接、激活状态 |
| Success | #22C55E | 启用状态、成功提示 |
| Warning | #E6A23C | 警告、中风险 |
| Danger | #EF4444 | 禁用状态、删除、高风险 |
| Info | #909399 | 辅助信息、禁用标签 |

---

## 页面结构设计

### 1. 同步策略列表页 (`sync-policies/index.vue`)

#### 页面结构
```
├── div.sync-policies-container { padding: 20px }
│   ├── div.page-header { flex, margin-bottom: 20px }
│   │   ├── h2 "同步策略"
│   │   └── div.toolbar { gap: 8px }
│   │       ├── el-button [刷新] (icon: Refresh)
│   │       ├── el-button [新建策略] (type: primary, icon: Plus)
│   │       ├── el-dropdown [批量操作] (icon: Operation)
│   │       │   ├── 批量启用
│   │       │   ├── 批量禁用
│   │       │   ├── 批量删除
│   │       ├── el-button [导出] (icon: Download)
│   │
│   ├── el-tabs [顶部分类Tab]
│   │   ├── el-tab-pane label="全部"
│   │   ├── el-tab-pane label="已启用"
│   │   ├── el-tab-pane label="已禁用"
│   │
│   ├── el-card.filter-card { margin-bottom: 20px }
│   │   └── el-form inline
│   │       ├── el-form-item [名称] (el-input + placeholder: "支持策略名称/ID搜索")
│   │       ├── el-form-item [规则类型] (el-select: 全部匹配/至少一个/Key匹配)
│   │       ├── el-form-item [状态] (el-select: 启用/禁用)
│   │       ├── el-button [查询]
│   │       ├── el-button [重置]
│   │
│   ├── el-table { width: 100%, row-key="id" }
│   │   ├── el-table-column [ID] width=80
│   │   ├── el-table-column [策略名称] min-width=150
│   │   │   └── 点击名称打开详情抽屉
│   │   ├── el-table-column [规则数量] width=100
│   │   ├── el-table-column [应用范围] width=120
│   │   ├── el-table-column [执行次数] width=100
│   │   ├── el-table-column [最后执行时间] width=160
│   │   ├── el-table-column [状态] width=100 (el-switch)
│   │   ├── el-table-column [创建时间] width=160
│   │   ├── el-table-column [操作] width=140 fixed="right"
│   │   │   ├── el-button [查看] (打开详情抽屉)
│   │   │   ├── el-dropdown [更多]
│   │   │   │   ├── 执行策略
│   │   │   │   ├── 编辑
│   │   │   │   ├── 复制
│   │   │   │   ├── ───────────
│   │   │   │   ├── 启用/禁用
│   │   │   │   ├── 删除
│   │
│   ├── el-pagination.pagination { margin-top: 20px, text-align: right }
```

#### 工具区设计
```vue
<div class="toolbar">
  <!-- 刷新按钮 -->
  <el-button size="default" @click="loadPolicies">
    <el-icon><Refresh /></el-icon>
    刷新
  </el-button>
  
  <!-- 新建策略 -->
  <el-button type="primary" @click="showCreateDialog">
    <el-icon><Plus /></el-icon>
    新建策略
  </el-button>
  
  <!-- 批量操作下拉 -->
  <el-dropdown trigger="click" @command="handleBatchCommand" :disabled="selectedIds.length === 0">
    <el-button>
      <el-icon><Operation /></el-icon>
      批量操作
      <el-icon class="el-icon--right"><ArrowDown /></el-icon>
    </el-button>
    <template #dropdown>
      <el-dropdown-menu>
        <el-dropdown-item command="batchEnable">批量启用</el-dropdown-item>
        <el-dropdown-item command="batchDisable">批量禁用</el-dropdown-item>
        <el-dropdown-item command="batchDelete" divided>批量删除</el-dropdown-item>
      </el-dropdown-menu>
    </template>
  </el-dropdown>
  
  <!-- 导出 -->
  <el-button @click="handleExport">
    <el-icon><Download /></el-icon>
    导出
  </el-button>
</div>
```

#### 搜索提示设计
```vue
<el-form-item label="名称">
  <el-input 
    v-model="queryForm.name" 
    placeholder="支持策略名称、ID搜索"
    clearable
  >
    <template #prefix>
      <el-icon><Search /></el-icon>
    </template>
  </el-input>
</el-form-item>
```

---

### 2. 同步策略详情抽屉 (`SyncPolicyDetailDrawer.vue`)

#### 抽屉结构
```
├── el-drawer
│   ├── size="60%"
│   ├── direction="rtl"
│   ├── withHeader=false (自定义顶部)
│   │
│   ├── div.drawer-header { flex, padding: 16px 20px }
│   │   ├── div.policy-icon
│   │   │   └── el-avatar :size=48
│   │   │       └── el-icon :size=32 <Setting />
│   │   ├── div.policy-info
│   │   │   ├── h3 {策略名称}
│   │   │   └── div.policy-tags { gap: 8px }
│   │   │       ├── el-tag [应用范围]
│   │   │       ├── el-tag [规则数量]
│   │   │       ├── el-switch [启用状态] (inline)
│   │   ├── div.quick-actions { gap: 8px }
│   │   │   ├── el-button [执行] (type: primary, icon: PlayOne)
│   │   │   ├── el-button [编辑] (icon: Edit)
│   │   │   ├── el-dropdown [更多]
│   │   │   │   ├── 复制策略
│   │   │   │   ├── 查看日志
│   │   │   │   ├── ───────────
│   │   │   │   ├── 删除
│   │
│   ├── el-tabs type="border-card" v-model="activeTab"
│   │   ├── el-tab-pane label="规则概览" name="rulesOverview"
│   │   │   └── RulesOverviewTab
│   │   ├── el-tab-pane label="执行日志" name="executionLogs"
│   │   │   └── ExecutionLogsTab
│   │   ├── el-tab-pane label="映射结果" name="mappingResults"
│   │   │   └── MappingResultsTab
```

#### 顶部区域设计
```vue
<div class="drawer-header">
  <!-- 策略图标 -->
  <div class="policy-icon">
    <el-avatar :size="48" :style="{ backgroundColor: '#409EFF' }">
      <el-icon :size="32"><Setting /></el-icon>
    </el-avatar>
  </div>
  
  <!-- 策略信息 -->
  <div class="policy-info">
    <h3>{{ policy?.name }}</h3>
    <div class="policy-tags">
      <el-tag size="small">{{ getScopeText(policy?.scope) }}</el-tag>
      <el-tag size="small" type="info">{{ policy?.rules?.length || 0 }} 条规则</el-tag>
      <!-- 启用状态开关 -->
      <el-switch 
        v-model="policy.enabled"
        size="small"
        @change="handleToggleStatus"
      />
    </div>
  </div>
  
  <!-- 快捷操作 -->
  <div class="quick-actions">
    <el-button size="small" type="primary" @click="handleExecute">
      <el-icon><PlayOne /></el-icon>
      执行
    </el-button>
    <el-button size="small" @click="handleEdit">
      <el-icon><Edit /></el-icon>
      编辑
    </el-button>
    <el-dropdown trigger="click" @command="handleQuickCommand">
      <el-button size="small">
        更多 <el-icon class="el-icon--right"><ArrowDown /></el-icon>
      </el-button>
      <template #dropdown>
        <el-dropdown-menu>
          <el-dropdown-item command="copy">复制策略</el-dropdown-item>
          <el-dropdown-item command="logs">查看日志</el-dropdown-item>
          <el-dropdown-item command="delete" divided>
            <span style="color: var(--el-color-danger)">删除</span>
          </el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>
  </div>
</div>
```

---

### 3. Tab 内容设计

#### Tab 1: 规则概览 (`RulesOverviewTab.vue`)
```
├── div.rules-overview-tab { padding: 16px }
│   ├── div.toolbar { margin-bottom: 16px }
│   │   ├── el-button [添加规则] (icon: Plus)
│   │   ├── el-button [刷新] (icon: Refresh)
│   │
│   ├── el-collapse v-model="expandedRules"
│   │   └── el-collapse-item v-for="rule in rules"
│   │       ├── template #title
│   │       │   ├── div.rule-title { flex }
│   │       │   │   ├── span "规则 {index + 1}"
│   │       │   │   ├── el-tag [条件类型]
│   │       │   │   ├── el-tag [映射方式]
│   │       │   │   ├── el-button [删除] (icon: Delete, link)
│   │       │
│   │       ├── div.rule-content { padding: 16px }
│   │       │   ├── el-form
│   │       │   │   ├── el-form-item [条件类型]
│   │       │   │   ├── el-form-item [资源映射]
│   │       │   │   ├── el-form-item [目标项目] (if specify_project)
│   │       │   │   ├── el-form-item [匹配标签]
│   │       │   │   │   ├── div.tags-list
│   │       │   │   │   │   ├── el-tag v-for="tag" closable
│   │       │   │   │   └── el-button [添加标签]
```

#### Tab 2: 执行日志 (`ExecutionLogsTab.vue`)
```
├── div.execution-logs-tab { padding: 16px }
│   ├── div.toolbar { margin-bottom: 16px }
│   │   ├── el-select [结果筛选] (成功/失败)
│   │   ├── el-date-picker [时间范围]
│   │   ├── el-button [刷新]
│   │
│   ├── el-table
│   │   ├── el-table-column [执行时间] width=160
│   │   ├── el-table-column [触发方式] width=100 (手动/自动)
│   │   ├── el-table-column [资源数量] width=100
│   │   ├── el-table-column [匹配数量] width=100
│   │   ├── el-table-column [执行结果] width=80 (el-tag)
│   │   ├── el-table-column [耗时] width=80
│   │   ├── el-table-column [操作] width=80
│   │   │   └── el-button [查看详情]
│   │
│   ├── el-pagination
```

#### Tab 3: 映射结果 (`MappingResultsTab.vue`)
```
├── div.mapping-results-tab { padding: 16px }
│   ├── div.toolbar
│   │   ├── el-input [搜索资源]
│   │   ├── el-select [项目筛选]
│   │   ├── el-button [刷新]
│   │
│   ├── el-table
│   │   ├── el-table-column [资源名称] min-width=150
│   │   ├── el-table-column [资源类型] width=100
│   │   ├── el-table-column [云账号] width=120
│   │   ├── el-table-column [匹配规则] width=100
│   │   ├── el-table-column [映射项目] width=120
│   │   ├── el-table-column [匹配标签] min-width=200
│   │   │   └── el-tag v-for="tag"
│   │   ├── el-table-column [映射时间] width=160
│   │
│   ├── el-pagination
```

---

### 4. 规则编辑器优化

#### 可视化规则展示
```vue
<!-- 规则卡片设计 -->
<el-card class="rule-card" shadow="hover">
  <template #header>
    <div class="rule-header">
      <span class="rule-index">规则 {{ index + 1 }}</span>
      <div class="rule-tags">
        <el-tag size="small" type="primary">
          {{ getConditionText(rule.condition_type) }}
        </el-tag>
        <el-tag size="small" type="success">
          {{ getMappingText(rule.resource_mapping) }}
        </el-tag>
      </div>
      <el-button 
        type="danger" 
        size="small" 
        link
        @click="removeRule(index)"
      >
        <el-icon><Delete /></el-icon>
      </el-button>
    </div>
  </template>
  
  <!-- 规则内容 -->
  <div class="rule-content">
    <!-- 条件类型选择 -->
    <el-form-item label="条件类型">
      <el-select v-model="rule.condition_type" style="width: 100%">
        <el-option 
          label="全部匹配标签（AND）" 
          value="all_match"
        >
          <span>全部匹配标签</span>
          <span style="color: #909399; margin-left: 8px">所有标签必须匹配</span>
        </el-option>
        <el-option 
          label="至少一个标签（OR）" 
          value="any_match"
        >
          <span>至少一个标签</span>
          <span style="color: #909399; margin-left: 8px">任一标签匹配即可</span>
        </el-option>
        <el-option 
          label="根据标签Key匹配" 
          value="key_match"
        >
          <span>根据标签Key匹配</span>
          <span style="color: #909399; margin-left: 8px">只匹配Key，忽略Value</span>
        </el-option>
      </el-select>
    </el-form-item>
    
    <!-- 资源映射选择 -->
    <el-form-item label="资源映射">
      <el-select v-model="rule.resource_mapping" style="width: 100%">
        <el-option label="指定项目" value="specify_project" />
        <el-option label="根据名称映射" value="specify_name" />
      </el-select>
    </el-form-item>
    
    <!-- 目标项目选择器 -->
    <el-form-item 
      label="目标项目" 
      v-if="rule.resource_mapping === 'specify_project'"
    >
      <el-select 
        v-model="rule.target_project_id"
        placeholder="请选择目标项目"
        filterable
        style="width: 100%"
      >
        <el-option 
          v-for="project in projects"
          :key="project.id"
          :label="project.name"
          :value="project.id"
        />
      </el-select>
    </el-form-item>
    
    <!-- 标签编辑区 -->
    <el-form-item label="匹配标签">
      <div class="tags-editor">
        <!-- 标签列表 -->
        <div class="tags-list">
          <div 
            v-for="(tag, tagIndex) in rule.tags"
            :key="tagIndex"
            class="tag-row"
          >
            <el-input 
              v-model="tag.tag_key"
              placeholder="标签Key"
              style="width: 140px"
            />
            <span style="margin: 0 8px">=</span>
            <el-input 
              v-model="tag.tag_value"
              placeholder="标签Value"
              style="width: 140px"
            />
            <el-button 
              type="danger"
              size="small"
              link
              @click="removeTag(ruleIndex, tagIndex)"
            >
              <el-icon><Close /></el-icon>
            </el-button>
          </div>
        </div>
        
        <!-- 添加标签按钮 -->
        <el-button 
          size="small"
          @click="addTag(ruleIndex)"
          style="margin-top: 8px"
        >
          <el-icon><Plus /></el-icon>
          添加标签
        </el-button>
      </div>
    </el-form-item>
  </div>
</el-card>
```

---

## 新增后端API设计

### 1. 执行策略API
```go
// POST /sync-policies/:id/execute
func (h *SyncPolicyHandler) Execute(c *gin.Context) {
    // 执行策略，返回执行ID
}
```

### 2. 执行日志API
```go
// GET /sync-policies/:id/execution-logs
func (h *SyncPolicyHandler) GetExecutionLogs(c *gin.Context) {
    // 返回策略执行日志列表
}
```

### 3. 映射结果API
```go
// GET /sync-policies/:id/mapping-results
func (h *SyncPolicyHandler) GetMappingResults(c *gin.Context) {
    // 返回策略映射结果列表
}
```

### 4. 批量操作API
```go
// POST /sync-policies/batch-enable
func (h *SyncPolicyHandler) BatchEnable(c *gin.Context) {
    // 批量启用策略
}

// POST /sync-policies/batch-disable
func (h *SyncPolicyHandler) BatchDisable(c *gin.Context) {
    // 批量禁用策略
}

// POST /sync-policies/batch-delete
func (h *SyncPolicyHandler) BatchDelete(c *gin.Context) {
    // 批量删除策略
}
```

---

## 实施顺序

### P0: 列表页基础功能完善
1. 添加工具区完整按钮
2. 添加顶部Tab分类
3. 添加搜索提示文案
4. 添加表格选择列和批量操作
5. 点击名称打开详情抽屉
6. 更多菜单分组优化

### P1: 详情抽屉改造
1. 创建 SyncPolicyDetailDrawer.vue
2. 实现顶部区域（图标/名称/启停开关/快捷操作）
3. 创建 RulesOverviewTab.vue
4. 创建 ExecutionLogsTab.vue
5. 创建 MappingResultsTab.vue

### P2: 规则编辑器优化
1. 规则卡片可视化设计
2. 条件类型选择器带说明
3. 标签编辑器行式布局

### P3: 后端API补齐
1. 执行策略API
2. 执行日志API和模型
3. 映射结果API
4. 批量操作API
5. 导出API

---

## 风险与注意事项

1. **详情弹窗改为抽屉**: 需重构现有 showDetailDialog 相关逻辑
2. **执行日志表**: 需新建 SyncPolicyExecutionLog 模型
3. **映射结果**: 需关联 CloudProject.LocalProjectID，可能需新建表
4. **批量操作**: 需后端支持事务和并发控制
5. **前端API文件**: 需在 sync-policy.ts 中新增 API 函数