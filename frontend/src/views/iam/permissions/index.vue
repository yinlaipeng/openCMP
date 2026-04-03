<template>
  <div class="permissions-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <span class="title">权限管理</span>
            <el-tabs v-model="activeTab" class="permission-tabs" @tab-change="handleTabChange">
              <el-tab-pane label="全部" name="all">
                <span class="tab-count">({{ totalPermissions }})</span>
              </el-tab-pane>
              <el-tab-pane label="系统权限" name="system">
                <span class="tab-count">({{ systemCount }})</span>
              </el-tab-pane>
              <el-tab-pane label="自定义权限" name="custom">
                <span class="tab-count">({{ customCount }})</span>
              </el-tab-pane>
            </el-tabs>
          </div>
          <el-button type="primary" @click="handleCreate">
            <el-icon><Plus /></el-icon>
            添加权限
          </el-button>
        </div>
      </template>

      <el-form :inline="true" :model="queryForm" class="query-form">
        <el-form-item label="资源类型">
          <el-select
            v-model="queryForm.resource"
            placeholder="全部"
            clearable
            @change="loadPermissions"
          >
            <el-option label="全部" value="" />
            <el-option
              v-for="item in resources"
              :key="item"
              :label="getResourceLabel(item)"
              :value="item"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="操作类型">
          <el-select
            v-model="queryForm.action"
            placeholder="全部"
            clearable
            @change="loadPermissions"
          >
            <el-option label="全部" value="" />
            <el-option
              v-for="item in actions"
              :key="item"
              :label="getActionLabel(item)"
              :value="item"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="">
          <el-input
            v-model="queryForm.keyword"
            placeholder="搜索权限名称或描述"
            clearable
            @input="loadPermissions"
            style="width: 240px"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </el-form-item>
      </el-form>

      <el-table
        :data="displayedPermissions"
        v-loading="loading"
        row-key="id"
        :height="tableHeight"
      >
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="name" label="权限标识" min-width="200" show-overflow-tooltip />
        <el-table-column prop="display_name" label="显示名称" min-width="140" show-overflow-tooltip />
        <el-table-column prop="resource" label="资源类型" width="110">
          <template #default="{ row }">
            <el-tag :type="getResourceTagType(row.resource)" size="small">
              {{ getResourceLabel(row.resource) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="action" label="操作类型" width="110">
          <template #default="{ row }">
            <el-tag :type="getActionTagType(row.action)" size="small">
              {{ getActionLabel(row.action) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="type" label="类型" width="90">
          <template #default="{ row }">
            <el-tag :type="row.type === 'system' ? 'warning' : 'success'" size="small">
              {{ row.type === 'system' ? '系统' : '自定义' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
        <el-table-column label="操作" width="140" fixed="right">
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
          :page-sizes="[50, 100, 200]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="loadPermissions"
          @current-change="loadPermissions"
        />
      </div>
    </el-card>

    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑权限' : '添加权限'"
      width="600px"
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="120px">
        <el-form-item label="权限标识" prop="name">
          <el-input v-model="form.name" placeholder="例如：cloud_account:list" :disabled="isEdit" />
          <div class="form-tip">格式：资源类型：操作，如 cloud_account:list</div>
        </el-form-item>
        <el-form-item label="显示名称" prop="display_name">
          <el-input v-model="form.display_name" placeholder="例如：查看云账户" />
        </el-form-item>
        <el-form-item label="资源类型" prop="resource">
          <el-select
            v-model="form.resource"
            placeholder="请选择资源类型"
            style="width: 100%"
            @change="handleResourceChange"
          >
            <el-option
              v-for="item in resources"
              :key="item"
              :label="getResourceLabel(item)"
              :value="item"
            />
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
            <el-option label="启动" value="start" />
            <el-option label="停止" value="stop" />
            <el-option label="重启" value="restart" />
            <el-option label="授权" value="grant" />
            <el-option label="撤销" value="revoke" />
            <el-option label="同步" value="sync" />
            <el-option label="导入" value="import" />
            <el-option label="导出" value="export" />
          </el-select>
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <el-radio-group v-model="form.type">
            <el-radio label="system">系统</el-radio>
            <el-radio label="custom">自定义</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="3"
            placeholder="权限描述"
          />
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
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Permission } from '@/types/iam'
import { getPermissions, createPermission, updatePermission, deletePermission, getResources, getActions } from '@/api/iam'

const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref<FormInstance>()
const activeTab = ref('all')
const tableHeight = ref(500)

const queryForm = reactive({
  resource: '',
  action: '',
  keyword: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 50,
  total: 0
})

// 所有权限数据
const allPermissions = ref<Permission[]>([])

// 资源类型和操作类型列表
const resources = ref<string[]>([])
const actions = ref<string[]>([])

// 系统权限数量
const systemCount = computed(() => {
  return allPermissions.value.filter(p => p.type === 'system').length
})

// 自定义权限数量
const customCount = computed(() => {
  return allPermissions.value.filter(p => p.type === 'custom').length
})

// 总权限数
const totalPermissions = computed(() => {
  return allPermissions.value.length
})

// 根据当前选中的 tab 和筛选条件返回显示的权限
const displayedPermissions = computed(() => {
  let result = [...allPermissions.value]

  // 根据 tab 筛选
  if (activeTab.value === 'system') {
    result = result.filter(p => p.type === 'system')
  } else if (activeTab.value === 'custom') {
    result = result.filter(p => p.type === 'custom')
  }

  // 根据资源类型筛选
  if (queryForm.resource) {
    result = result.filter(p => p.resource === queryForm.resource)
  }

  // 根据操作类型筛选
  if (queryForm.action) {
    result = result.filter(p => p.action === queryForm.action)
  }

  // 根据关键词搜索
  if (queryForm.keyword) {
    const keyword = queryForm.keyword.toLowerCase()
    result = result.filter(p =>
      p.name.toLowerCase().includes(keyword) ||
      p.display_name.toLowerCase().includes(keyword) ||
      p.description.toLowerCase().includes(keyword)
    )
  }

  // 分页
  pagination.total = result.length
  const start = (pagination.page - 1) * pagination.pageSize
  const end = start + pagination.pageSize
  return result.slice(start, end)
})

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

// 资源类型映射
const resourceMap: Record<string, string> = {
  cloud_account: '云账户',
  vm: '虚拟机',
  image: '镜像',
  disk: '磁盘',
  vpc: 'VPC',
  subnet: '子网',
  security_group: '安全组',
  eip: '弹性 IP',
  loadbalancer: '负载均衡',
  database: '数据库',
  user: '用户',
  role: '角色',
  permission: '权限',
  policy: '策略',
  auth_source: '认证源',
  message: '消息',
  alert: '告警',
  system: '系统'
}

// 操作类型映射
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

const getResourceLabel = (resource: string) => {
  return resourceMap[resource] || resource
}

const getActionLabel = (action: string) => {
  return actionMap[action] || action
}

const getResourceTagType = (resource: string) => {
  const map: Record<string, any> = {
    cloud_account: 'primary',
    vm: 'success',
    image: 'info',
    disk: 'warning',
    vpc: 'primary',
    subnet: 'success',
    security_group: 'warning',
    eip: 'info',
    loadbalancer: 'primary',
    user: 'success',
    role: 'warning',
    permission: 'info',
    policy: 'primary',
    auth_source: 'success',
    message: 'info',
    alert: 'danger',
    system: 'warning'
  }
  return map[resource] || ''
}

const getActionTagType = (action: string) => {
  const map: Record<string, any> = {
    list: 'info',
    get: 'info',
    create: 'success',
    update: 'warning',
    delete: 'danger',
    action: 'warning',
    grant: 'success',
    revoke: 'danger',
    start: 'success',
    stop: 'warning',
    restart: 'warning',
    sync: 'primary',
    import: 'primary',
    export: 'primary',
    verify: 'success',
    share: 'success',
    bind: 'primary',
    unbind: 'warning',
    attach: 'success',
    detach: 'warning',
    snapshot: 'warning',
    backup: 'warning',
    restore: 'success',
    resize: 'warning',
    migrate: 'primary',
    rebuild: 'danger',
    reset: 'warning',
    enable: 'success',
    disable: 'danger',
    clone: 'primary',
    test: 'info',
    read: 'info',
    resolve: 'success',
    ignore: 'info'
  }
  return map[action] || ''
}

const loadPermissions = async () => {
  loading.value = true
  try {
    const res = await getPermissions({
      resource: queryForm.resource,
      action: queryForm.action,
      type: activeTab.value === 'all' ? '' : activeTab.value,
      keyword: queryForm.keyword,
      limit: pagination.pageSize,
      offset: (pagination.page - 1) * pagination.pageSize
    })

    if (res.items) {
      allPermissions.value = res.items
      pagination.total = res.total
    }
  } catch (e: any) {
    console.error(e)
    ElMessage.error('加载权限列表失败')
  } finally {
    loading.value = false
  }
}

const loadResourcesAndActions = async () => {
  try {
    const [resRes, actRes] = await Promise.all([getResources(), getActions()])
    resources.value = resRes.items || []
    actions.value = actRes.items || []
  } catch (e) {
    console.error(e)
  }
}

const handleTabChange = (tab: string) => {
  pagination.page = 1
  loadPermissions()
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

const handleEdit = (row: Permission) => {
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

const handleDelete = async (row: Permission) => {
  try {
    await ElMessageBox.confirm('确定要删除该权限吗？', '提示', { type: 'warning' })
    await deletePermission(row.id)
    ElMessage.success('删除成功')
    loadPermissions()
  } catch (e: any) {
    if (e !== 'cancel') {
      console.error(e)
    }
  }
}

const handleResourceChange = () => {
  // 资源类型改变时，可以重置操作类型
  form.action = ''
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitting.value = true
    try {
      if (isEdit.value) {
        await updatePermission(form.id, {
          display_name: form.display_name,
          description: form.description
        })
        ElMessage.success('更新成功')
      } else {
        await createPermission({
          name: form.name,
          display_name: form.display_name,
          resource: form.resource,
          action: form.action,
          type: form.type,
          description: form.description
        })
        ElMessage.success('创建成功')
      }
      dialogVisible.value = false
      loadPermissions()
    } catch (e: any) {
      console.error(e)
      ElMessage.error(isEdit.value ? '更新失败' : '创建失败')
    } finally {
      submitting.value = false
    }
  })
}

onMounted(() => {
  loadResourcesAndActions()
  loadPermissions()
})
</script>

<style scoped>
.permissions-page {
  height: 100%;
  padding: 20px;
}

.page-card {
  height: calc(100% - 40px);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left {
  display: flex;
  flex-direction: column;
  gap: 12px;
  flex: 1;
}

.title {
  font-size: 18px;
  font-weight: bold;
}

.permission-tabs {
  margin-top: 8px;
}

.tab-count {
  font-size: 12px;
  color: #999;
  margin-left: 4px;
}

:deep(.permission-tabs .el-tabs__item) {
  font-size: 14px;
}

.query-form {
  margin-bottom: 16px;
  margin-top: 16px;
}

.pagination {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}

.form-tip {
  font-size: 12px;
  color: #999;
  margin-top: 5px;
}
</style>
