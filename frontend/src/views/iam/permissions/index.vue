<template>
  <div class="permissions-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">权限管理</span>
          <el-button type="primary" @click="handleCreate">
            <el-icon><Plus /></el-icon>
            添加权限
          </el-button>
        </div>
      </template>
      
      <el-form :inline="true" :model="queryForm" class="query-form">
        <el-form-item label="资源类型">
          <el-select v-model="queryForm.resource" placeholder="全部" clearable>
            <el-option label="云账户" value="cloud_account" />
            <el-option label="虚拟机" value="vm" />
            <el-option label="网络" value="network" />
            <el-option label="存储" value="storage" />
            <el-option label="用户" value="user" />
            <el-option label="角色" value="role" />
          </el-select>
        </el-form-item>
        <el-form-item label="权限类型">
          <el-select v-model="queryForm.type" placeholder="全部" clearable>
            <el-option label="系统" value="system" />
            <el-option label="自定义" value="custom" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadPermissions">查询</el-button>
        </el-form-item>
      </el-form>
      
      <el-table :data="permissions" v-loading="loading" row-key="id">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="权限标识" width="200" />
        <el-table-column prop="display_name" label="显示名称" width="150" />
        <el-table-column prop="resource" label="资源类型" width="120">
          <template #default="{ row }">
            <el-tag>{{ getResourceName(row.resource) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="action" label="操作类型" width="120">
          <template #default="{ row }">
            <el-tag :type="getActionType(row.action)">{{ getActionName(row.action) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="row.type === 'system' ? 'warning' : ''">
              {{ row.type === 'system' ? '系统' : '自定义' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="描述" />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button 
              v-if="row.type === 'custom'" 
              size="small" 
              type="danger" 
              @click="handleDelete(row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <div class="pagination">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :total="pagination.total"
          :page-sizes="[20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="loadPermissions"
          @current-change="loadPermissions"
        />
      </div>
    </el-card>
    
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑权限' : '添加权限'" width="600px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="120px">
        <el-form-item label="权限标识" prop="name">
          <el-input v-model="form.name" placeholder="例如：cloud_account:list" />
          <div class="form-tip">格式：资源类型：操作，如 cloud_account:list</div>
        </el-form-item>
        <el-form-item label="显示名称" prop="display_name">
          <el-input v-model="form.display_name" placeholder="例如：查看云账户" />
        </el-form-item>
        <el-form-item label="资源类型" prop="resource">
          <el-select v-model="form.resource" placeholder="请选择资源类型" style="width: 100%">
            <el-option label="云账户" value="cloud_account" />
            <el-option label="虚拟机" value="vm" />
            <el-option label="网络" value="network" />
            <el-option label="存储" value="storage" />
            <el-option label="数据库" value="database" />
            <el-option label="用户" value="user" />
            <el-option label="角色" value="role" />
            <el-option label="认证源" value="auth_source" />
            <el-option label="消息" value="message" />
          </el-select>
        </el-form-item>
        <el-form-item label="操作类型" prop="action">
          <el-select v-model="form.action" placeholder="请选择操作类型" style="width: 100%">
            <el-option label="查看列表" value="list" />
            <el-option label="查看详情" value="get" />
            <el-option label="创建" value="create" />
            <el-option label="更新" value="update" />
            <el-option label="删除" value="delete" />
            <el-option label="操作" value="action" />
            <el-option label="授权" value="grant" />
            <el-option label="撤销" value="revoke" />
          </el-select>
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="form.description" type="textarea" :rows="3" placeholder="权限描述" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'

const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref<FormInstance>()

const queryForm = reactive({
  resource: '',
  type: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

const permissions = ref<any[]>([])

const form = reactive({
  id: 0,
  name: '',
  display_name: '',
  resource: '',
  action: '',
  type: 'custom',
  description: ''
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入权限标识', trigger: 'blur' }],
  display_name: [{ required: true, message: '请输入显示名称', trigger: 'blur' }],
  resource: [{ required: true, message: '请选择资源类型', trigger: 'change' }],
  action: [{ required: true, message: '请选择操作类型', trigger: 'change' }]
}

const getResourceName = (resource: string) => {
  const map: Record<string, string> = {
    cloud_account: '云账户',
    vm: '虚拟机',
    network: '网络',
    storage: '存储',
    database: '数据库',
    user: '用户',
    role: '角色',
    auth_source: '认证源',
    message: '消息'
  }
  return map[resource] || resource
}

const getActionName = (action: string) => {
  const map: Record<string, string> = {
    list: '查看列表',
    get: '查看详情',
    create: '创建',
    update: '更新',
    delete: '删除',
    action: '操作',
    grant: '授权',
    revoke: '撤销'
  }
  return map[action] || action
}

const getActionType = (action: string) => {
  const map: Record<string, any> = {
    list: 'info',
    get: 'info',
    create: 'success',
    update: 'warning',
    delete: 'danger',
    action: 'warning',
    grant: 'success',
    revoke: 'danger'
  }
  return map[action] || ''
}

const loadPermissions = async () => {
  loading.value = true
  try {
    // 模拟数据
    permissions.value = [
      { id: 1, name: 'cloud_account:list', display_name: '查看云账户', resource: 'cloud_account', action: 'list', type: 'system', description: '查看云账户列表' },
      { id: 2, name: 'cloud_account:create', display_name: '创建云账户', resource: 'cloud_account', action: 'create', type: 'system', description: '创建新的云账户' },
      { id: 3, name: 'cloud_account:update', display_name: '更新云账户', resource: 'cloud_account', action: 'update', type: 'system', description: '更新云账户信息' },
      { id: 4, name: 'cloud_account:delete', display_name: '删除云账户', resource: 'cloud_account', action: 'delete', type: 'system', description: '删除云账户' },
      { id: 5, name: 'vm:list', display_name: '查看虚拟机', resource: 'vm', action: 'list', type: 'system', description: '查看虚拟机列表' },
      { id: 6, name: 'vm:create', display_name: '创建虚拟机', resource: 'vm', action: 'create', type: 'system', description: '创建虚拟机' },
      { id: 7, name: 'vm:delete', display_name: '删除虚拟机', resource: 'vm', action: 'delete', type: 'system', description: '删除虚拟机' },
      { id: 8, name: 'vm:action', display_name: '操作虚拟机', resource: 'vm', action: 'action', type: 'system', description: '启动/停止/重启虚拟机' },
      { id: 9, name: 'user:list', display_name: '查看用户', resource: 'user', action: 'list', type: 'system', description: '查看用户列表' },
      { id: 10, name: 'user:create', display_name: '创建用户', resource: 'user', action: 'create', type: 'system', description: '创建新用户' },
      { id: 11, name: 'user:update', display_name: '更新用户', resource: 'user', action: 'update', type: 'system', description: '更新用户信息' },
      { id: 12, name: 'user:delete', display_name: '删除用户', resource: 'user', action: 'delete', type: 'system', description: '删除用户' },
      { id: 13, name: 'role:list', display_name: '查看角色', resource: 'role', action: 'list', type: 'system', description: '查看角色列表' },
      { id: 14, name: 'role:grant', display_name: '角色授权', resource: 'role', action: 'grant', type: 'system', description: '为角色分配权限' },
      { id: 15, name: 'auth_source:list', display_name: '查看认证源', resource: 'auth_source', action: 'list', type: 'system', description: '查看认证源列表' },
      { id: 16, name: 'auth_source:create', display_name: '创建认证源', resource: 'auth_source', action: 'create', type: 'system', description: '创建新的认证源' },
      { id: 17, name: 'message:list', display_name: '查看消息', resource: 'message', action: 'list', type: 'system', description: '查看消息列表' },
      { id: 18, name: 'alert:resolve', display_name: '处理告警', resource: 'alert', action: 'resolve', type: 'system', description: '处理安全告警' }
    ]
    
    // 过滤
    if (queryForm.resource) {
      permissions.value = permissions.value.filter(p => p.resource === queryForm.resource)
    }
    if (queryForm.type) {
      permissions.value = permissions.value.filter(p => p.type === queryForm.type)
    }
    
    pagination.total = permissions.value.length
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const handleCreate = () => {
  isEdit.value = false
  form.id = 0
  form.name = ''
  form.display_name = ''
  form.resource = ''
  form.action = ''
  form.type = 'custom'
  form.description = ''
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  isEdit.value = true
  form.id = row.id
  form.name = row.name
  form.display_name = row.display_name
  form.resource = row.resource
  form.action = row.action
  form.type = row.type
  form.description = row.description
  dialogVisible.value = true
}

const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm('确定要删除该权限吗？', '提示', { type: 'warning' })
    ElMessage.success('删除成功')
    loadPermissions()
  } catch (e: any) {
    if (e !== 'cancel') {
      console.error(e)
    }
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitting.value = true
    try {
      ElMessage.success(isEdit.value ? '更新成功' : '创建成功')
      dialogVisible.value = false
      loadPermissions()
    } catch (e) {
      console.error(e)
    } finally {
      submitting.value = false
    }
  })
}

onMounted(() => {
  loadPermissions()
})
</script>

<style scoped>
.permissions-page {
  height: 100%;
}

.page-card {
  height: 100%;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title {
  font-size: 18px;
  font-weight: bold;
}

.query-form {
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.form-tip {
  font-size: 12px;
  color: #999;
  margin-top: 5px;
}
</style>
