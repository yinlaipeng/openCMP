# openCMP IAM 模块前端开发指南

## 1. 项目结构

```
frontend/
├── src/
│   ├── api/                 # API 接口定义
│   │   └── iam.ts          # IAM 相关 API
│   ├── types/              # 类型定义
│   │   ├── index.ts        # 通用类型
│   │   └── iam.ts          # IAM 相关类型
│   ├── views/              # 页面组件
│   │   └── iam/            # IAM 模块页面
│   │       ├── domains/    # 域管理
│   │       ├── projects/   # 项目管理
│   │       ├── users/      # 用户管理
│   │       ├── groups/     # 用户组管理
│   │       ├── roles/      # 角色管理
│   │       ├── permissions/ # 权限管理
│   │       ├── policies/   # 策略管理
│   │       └── auth-sources/ # 认证源管理
│   └── utils/              # 工具函数
│       └── request.ts      # 请求封装
```

## 2. 类型定义

### 2.1 用户类型
```typescript
export interface User {
  id: number
  name: string
  display_name: string
  email: string
  phone: string
  enabled: boolean
  mfa_enabled: boolean
  domain_id: number
  created_at: string
  updated_at: string
}
```

### 2.2 角色类型
```typescript
export interface Role {
  id: number
  name: string
  display_name: string
  description: string
  type: 'system' | 'custom'
  enabled: boolean
  is_public: boolean
  domain_id: number
  created_at: string
  updated_at: string
}
```

### 2.3 权限类型
```typescript
export interface Permission {
  id: number
  name: string
  display_name: string
  resource: string
  action: string
  type: 'system' | 'custom'
  description: string
  domain_id: number
  created_at: string
  updated_at: string
}
```

### 2.4 策略类型
```typescript
export interface Policy {
  id: string
  name: string
  description: string
  scope: 'system' | 'domain' | 'project'
  type: 'system' | 'custom'
  enabled: boolean
  is_system: boolean
  is_public: boolean
  domain_id: string
  policy: Record<string, any>
  created_at: string
  updated_at: string
  can_update: boolean
  can_delete: boolean
  delete_fail_reason?: {
    details: string
  }
}
```

## 3. API 使用方法

### 3.1 导入 API 函数
```typescript
import {
  getUsers,
  createUser,
  updateUser,
  deleteUser,
  getRoles,
  createRole,
  updateRole,
  deleteRole,
  getPermissions,
  createPermission,
  updatePermission,
  deletePermission,
  getPolicies,
  createPolicy,
  updatePolicy,
  deletePolicy,
  getAuthSources,
  createAuthSource,
  updateAuthSource,
  deleteAuthSource,
  getDomains,
  createDomain,
  updateDomain,
  deleteDomain,
  getProjects,
  createProject,
  updateProject,
  deleteProject,
  getGroups,
  createGroup,
  updateGroup,
  deleteGroup
} from '@/api/iam'
```

### 3.2 获取用户列表
```typescript
const loadUsers = async () => {
  try {
    const res = await getUsers({
      limit: 20,
      offset: 0,
      domain_id: 1
    })
    console.log(res.items) // 用户列表
    console.log(res.total) // 总数
  } catch (e) {
    console.error(e)
  }
}
```

### 3.3 创建用户
```typescript
const createUserHandler = async () => {
  try {
    const newUser = await createUser({
      name: 'john_doe',
      display_name: 'John Doe',
      email: 'john@example.com',
      phone: '13800138000',
      password: 'password123',
      domain_id: 1
    })
    console.log(newUser)
  } catch (e) {
    console.error(e)
  }
}
```

## 4. 组件开发规范

### 4.1 页面组件结构
```vue
<template>
  <!-- 页面头部 -->
  <el-card>
    <template #header>
      <div class="card-header">
        <span class="title">用户管理</span>
        <el-button type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>
          新建用户
        </el-button>
      </div>
    </template>

    <!-- 筛选表单 -->
    <el-form :inline="true" :model="filterForm">
      <el-form-item label="用户名">
        <el-input v-model="filterForm.name" placeholder="请输入用户名" />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="loadUsers">查询</el-button>
      </el-form-item>
    </el-form>

    <!-- 数据表格 -->
    <el-table :data="users" v-loading="loading">
      <el-table-column prop="name" label="用户名" />
      <el-table-column prop="email" label="邮箱" />
      <el-table-column label="操作">
        <template #default="{ row }">
          <el-button size="small" @click="handleEdit(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <el-pagination
      v-model:current-page="pagination.page"
      v-model:page-size="pagination.limit"
      :total="pagination.total"
    />
  </el-card>

  <!-- 编辑对话框 -->
  <el-dialog v-model="dialogVisible" title="编辑用户">
    <el-form :model="form" :rules="rules" ref="formRef">
      <el-form-item label="用户名" prop="name">
        <el-input v-model="form.name" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="dialogVisible = false">取消</el-button>
      <el-button type="primary" @click="handleSubmit">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { User } from '@/types/iam'
import { getUsers, createUser, updateUser, deleteUser } from '@/api/iam'

const users = ref<User[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)
const currentUserId = ref(0)
const formRef = ref<FormInstance>()

const filterForm = reactive({
  name: ''
})

const pagination = reactive({
  page: 1,
  limit: 20,
  total: 0
})

const form = reactive({
  name: '',
  display_name: '',
  email: '',
  phone: '',
  password: '',
  domain_id: 1
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入用户名', trigger: 'blur' }]
}

const loadUsers = async () => {
  loading.value = true
  try {
    const res = await getUsers({
      limit: pagination.limit,
      offset: (pagination.page - 1) * pagination.limit,
      keyword: filterForm.name
    })
    users.value = res.items || []
    pagination.total = res.total || 0
  } catch (e: any) {
    ElMessage.error(e.message || '加载用户列表失败')
  } finally {
    loading.value = false
  }
}

const handleCreate = () => {
  isEdit.value = false
  form.name = ''
  form.display_name = ''
  form.email = ''
  form.phone = ''
  form.password = ''
  form.domain_id = 1
  dialogVisible.value = true
}

const handleEdit = (row: User) => {
  isEdit.value = true
  currentUserId.value = row.id
  form.name = row.name
  form.display_name = row.display_name
  form.email = row.email
  form.phone = row.phone || ''
  form.password = ''
  form.domain_id = row.domain_id
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    try {
      if (isEdit.value) {
        await updateUser(currentUserId.value, form)
        ElMessage.success('更新成功')
      } else {
        await createUser(form)
        ElMessage.success('创建成功')
      }
      dialogVisible.value = false
      loadUsers()
    } catch (e: any) {
      ElMessage.error(e.message || (isEdit.value ? '更新失败' : '创建失败'))
    }
  })
}

onMounted(() => {
  loadUsers()
})
</script>
```

## 5. 权限控制

### 5.1 路由权限控制
在 `router.ts` 中添加权限检查：

```typescript
// 路由守卫中添加权限检查
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  
  if (to.path === '/login') {
    next()
  } else {
    if (token) {
      // 检查用户是否有访问该路由的权限
      // 可以通过后端API检查或本地存储的权限信息
      next()
    } else {
      next('/login')
    }
  }
})
```

### 5.2 按钮权限控制
在组件中根据权限显示按钮：

```vue
<template>
  <el-table :data="data">
    <el-table-column label="操作">
      <template #default="{ row }">
        <el-button 
          v-if="hasPermission('user:update')" 
          size="small" 
          @click="handleEdit(row)">
          编辑
        </el-button>
        <el-button 
          v-if="hasPermission('user:delete')" 
          size="small" 
          type="danger" 
          @click="handleDelete(row)">
          删除
        </el-button>
      </template>
    </el-table-column>
  </el-table>
</template>

<script setup lang="ts">
// 检查用户是否有指定权限
const hasPermission = (permission: string): boolean => {
  // 从本地存储或全局状态获取用户权限列表
  const userPermissions = JSON.parse(localStorage.getItem('userPermissions') || '[]')
  return userPermissions.includes(permission)
}
</script>
```

## 6. 国际化支持

### 6.1 资源类型映射
```typescript
const resourceMap: Record<string, string> = {
  cloud_account: '云账户',
  vm: '虚拟机',
  image: '镜像',
  disk: '磁盘',
  vpc: 'VPC',
  subnet: '子网',
  security_group: '安全组',
  eip: '弹性 IP',
  user: '用户',
  role: '角色',
  permission: '权限',
  policy: '策略',
  auth_source: '认证源',
  message: '消息',
  alert: '告警',
  system: '系统'
}

const getResourceLabel = (resource: string) => {
  return resourceMap[resource] || resource
}
```

### 6.2 操作类型映射
```typescript
const actionMap: Record<string, string> = {
  list: '查看列表',
  get: '查看详情',
  create: '创建',
  update: '更新',
  delete: '删除',
  action: '操作',
  start: '启动',
  stop: '停止',
  restart: '重启',
  grant: '授权',
  revoke: '撤销',
  sync: '同步',
  import: '导入',
  export: '导出',
  verify: '验证',
  share: '共享',
  bind: '绑定',
  unbind: '解绑',
  attach: '挂载',
  detach: '卸载',
  snapshot: '快照',
  backup: '备份',
  restore: '恢复',
  resize: '扩容',
  migrate: '迁移',
  rebuild: '重建',
  reset: '重置',
  enable: '启用',
  disable: '禁用',
  clone: '克隆',
  test: '测试',
  read: '已读',
  resolve: '处理',
  ignore: '忽略'
}

const getActionLabel = (action: string) => {
  return actionMap[action] || action
}
```

## 7. 最佳实践

### 7.1 错误处理
```typescript
const loadUsers = async () => {
  loading.value = true
  try {
    const res = await getUsers(params)
    users.value = res.items || []
    pagination.total = res.total || 0
  } catch (e: any) {
    // 统一错误处理
    if (e.response?.status === 401) {
      // 未授权，跳转到登录页
      localStorage.removeItem('token')
      router.push('/login')
    } else if (e.response?.status === 403) {
      // 权限不足
      ElMessage.error('您没有权限执行此操作')
    } else {
      // 其他错误
      ElMessage.error(e.message || '加载失败')
    }
  } finally {
    loading.value = false
  }
}
```

### 7.2 表单验证
```typescript
const rules: FormRules = {
  name: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度应在3-20之间', trigger: 'blur' },
    { pattern: /^[a-zA-Z0-9_-]+$/, message: '用户名只能包含字母、数字、下划线和横线', trigger: 'blur' }
  ],
  email: [
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  password: [
    { min: 8, message: '密码长度至少8位', trigger: 'blur' },
    { pattern: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)/, message: '密码必须包含大小写字母和数字', trigger: 'blur' }
  ]
}
```

### 7.3 分页处理
```typescript
const loadUsers = async () => {
  loading.value = true
  try {
    const params = {
      limit: pagination.limit,
      offset: (pagination.page - 1) * pagination.limit,
      ...filterForm
    }
    
    const res = await getUsers(params)
    users.value = res.items || []
    pagination.total = res.total || 0
  } catch (e: any) {
    ElMessage.error(e.message || '加载失败')
  } finally {
    loading.value = false
  }
}

// 分页变化时重新加载数据
const handlePageChange = () => {
  loadUsers()
}
```

## 8. 常见问题

### 8.1 如何添加新的IAM实体？
1. 在 `types/iam.ts` 中定义类型
2. 在 `api/iam.ts` 中添加API函数
3. 在 `views/iam/` 中创建页面组件
4. 在 `router.ts` 中添加路由

### 8.2 如何处理权限检查？
1. 在登录时获取用户权限列表
2. 在路由守卫中检查权限
3. 在组件中根据权限显示/隐藏元素
4. 在API调用时处理权限不足错误

### 8.3 如何处理大数据量表格？
1. 使用分页加载数据
2. 实现虚拟滚动（Element Plus Table 支持）
3. 添加搜索和过滤功能
4. 优化API查询性能