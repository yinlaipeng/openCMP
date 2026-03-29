<template>
  <div class="policy-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <span class="title">策略管理</span>
            <el-tabs v-model="activeScope" class="policy-tabs" @tab-change="handleScopeChange">
              <el-tab-pane label="全部" name="all">
                <span class="tab-count">({{ totalPolicies }})</span>
              </el-tab-pane>
              <el-tab-pane label="系统策略" name="system">
                <span class="tab-count">({{ systemPolicies.length }})</span>
              </el-tab-pane>
              <el-tab-pane label="自定义策略" name="custom">
                <span class="tab-count">({{ customPolicies.length }})</span>
              </el-tab-pane>
            </el-tabs>
          </div>
          <el-button type="primary" @click="handleCreate">
            <el-icon><Plus /></el-icon>
            创建策略
          </el-button>
        </div>
      </template>

      <el-form :inline="true" :model="queryForm" class="query-form">
        <el-form-item label="作用域">
          <el-select v-model="queryForm.scope" placeholder="全部" clearable @change="loadPolicies">
            <el-option label="全部" value="" />
            <el-option label="系统级 (system)" value="system" />
            <el-option label="域级 (domain)" value="domain" />
            <el-option label="项目级 (project)" value="project" />
          </el-select>
        </el-form-item>
        <el-form-item label="">
          <el-input
            v-model="queryForm.keyword"
            placeholder="搜索策略名称或描述"
            clearable
            @input="loadPolicies"
            style="width: 240px"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </el-form-item>
      </el-form>

      <el-table :data="displayedPolicies" v-loading="loading" row-key="id" :height="tableHeight">
        <el-table-column prop="id" label="策略 ID" width="180" />
        <el-table-column prop="name" label="策略名称" min-width="200" show-overflow-tooltip>
          <template #default="{ row }">
            <span class="policy-name">{{ row.name }}</span>
            <el-tag v-if="row.is_system" type="warning" size="small" style="margin-left: 8px">系统</el-tag>
            <el-tag v-if="!row.enabled" type="info" size="small" style="margin-left: 8px">禁用</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="scope" label="作用域" width="120">
          <template #default="{ row }">
            <el-tag :type="getScopeTagType(row.scope)" size="small">
              {{ getScopeName(row.scope) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
        <el-table-column prop="updated_at" label="更新时间" width="180" />
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleView(row)">查看</el-button>
            <el-button size="small" @click="handleEdit(row)" :disabled="!row.can_update">编辑</el-button>
            <el-button
              v-if="row.can_delete"
              size="small"
              type="danger"
              @click="handleDelete(row)"
            >
              删除
            </el-button>
            <el-tag v-else size="small" type="info">不可删除</el-tag>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :total="pagination.total"
          :page-sizes="[50, 100, 170]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="loadPolicies"
          @current-change="loadPolicies"
        />
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑策略' : '创建策略'" width="800px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="120px">
        <el-form-item label="策略名称" prop="name">
          <el-input v-model="form.name" placeholder="例如：domain-vm-admin" />
          <div class="form-tip">命名规范：{scope}-{service}-{role}，如 domain-vm-admin</div>
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="form.description" type="textarea" :rows="3" placeholder="策略描述" />
        </el-form-item>
        <el-form-item label="作用域" prop="scope">
          <el-select v-model="form.scope" placeholder="请选择作用域" style="width: 100%">
            <el-option label="系统级 (system)" value="system" />
            <el-option label="域级 (domain)" value="domain" />
            <el-option label="项目级 (project)" value="project" />
          </el-select>
        </el-form-item>
        <el-form-item label="策略内容" prop="policy">
          <div class="policy-editor">
            <el-alert
              title="策略语法说明"
              type="info"
              :closable="false"
              style="margin-bottom: 15px"
            >
              <template #default>
                <div>策略格式：{ service: { resource: { action: "allow" | "deny" } } }</div>
                <div>例如：{ "compute": { "servers": { "*": "allow" } } } 表示允许管理所有虚拟机</div>
              </template>
            </el-alert>
            <el-table :data="form.statements" border style="margin-bottom: 10px">
              <el-table-column label="服务" width="150">
                <template #default="{ $index }">
                  <el-select v-model="form.statements[$index].service" size="small" style="width: 100%">
                    <el-option label="计算 (compute)" value="compute" />
                    <el-option label="镜像 (image)" value="image" />
                    <el-option label="网络 (network)" value="network" />
                    <el-option label="存储 (storage)" value="storage" />
                    <el-option label="监控 (monitor)" value="monitor" />
                    <el-option label="计费 (meter)" value="meter" />
                    <el-option label="身份 (identity)" value="identity" />
                    <el-option label="容器 (k8s)" value="k8s" />
                    <el-option label="云账号 (compute/cloudaccounts)" value="compute_cloudaccounts" />
                    <el-option label="全部 (*)" value="*" />
                  </el-select>
                </template>
              </el-table-column>
              <el-table-column label="资源" width="150">
                <template #default="{ $index }">
                  <el-input v-model="form.statements[$index].resource" size="small" placeholder="如：servers" />
                </template>
              </el-table-column>
              <el-table-column label="效果" width="120">
                <template #default="{ $index }">
                  <el-select v-model="form.statements[$index].effect" size="small" style="width: 100%">
                    <el-option label="允许" value="allow" />
                    <el-option label="拒绝" value="deny" />
                  </el-select>
                </template>
              </el-table-column>
              <el-table-column label="操作">
                <template #default="{ $index }">
                  <el-select v-model="form.statements[$index].actions" size="small" style="width: 100%" multiple>
                    <el-option label="所有操作 (*)" value="*" />
                    <el-option label="查看列表 (list)" value="list" />
                    <el-option label="查看详情 (get)" value="get" />
                    <el-option label="创建 (create)" value="create" />
                    <el-option label="更新 (update)" value="update" />
                    <el-option label="删除 (delete)" value="delete" />
                    <el-option label="执行操作 (perform)" value="perform" />
                  </el-select>
                </template>
              </el-table-column>
              <el-table-column label="操作" width="80">
                <template #default="{ $index }">
                  <el-button size="small" type="danger" @click="removeStatement($index)">删除</el-button>
                </template>
              </el-table-column>
            </el-table>
            <el-button type="primary" plain @click="addStatement">添加策略语句</el-button>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="viewVisible" title="策略详情" width="900px">
      <div v-if="currentPolicy" class="policy-detail">
        <el-descriptions :column="3" border>
          <el-descriptions-item label="策略 ID" :span="3">{{ currentPolicy.id }}</el-descriptions-item>
          <el-descriptions-item label="策略名称" :span="2">
            <span class="policy-name">{{ currentPolicy.name }}</span>
            <el-tag v-if="currentPolicy.is_system" type="warning" size="small" style="margin-left: 8px">系统</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="作用域">
            <el-tag :type="getScopeTagType(currentPolicy.scope)" size="small">
              {{ getScopeName(currentPolicy.scope) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="描述" :span="3">{{ currentPolicy.description }}</el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ currentPolicy.created_at }}</el-descriptions-item>
          <el-descriptions-item label="更新时间">{{ currentPolicy.updated_at }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="currentPolicy.enabled ? 'success' : 'info'" size="small">
              {{ currentPolicy.enabled ? '已启用' : '已禁用' }}
            </el-tag>
          </el-descriptions-item>
        </el-descriptions>

        <h4 class="section-title">策略内容</h4>
        <el-card class="policy-content">
          <pre>{{ JSON.stringify(currentPolicy.policy, null, 2) }}</pre>
        </el-card>

        <h4 class="section-title">删除限制</h4>
        <el-alert
          v-if="!currentPolicy.can_delete && currentPolicy.delete_fail_reason"
          :title="currentPolicy.delete_fail_reason.details"
          type="warning"
          :closable="false"
        />
        <el-alert
          v-else-if="currentPolicy.can_delete"
          title="此策略可以被删除"
          type="success"
          :closable="false"
        />
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'

const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const viewVisible = ref(false)
const isEdit = ref(false)
const formRef = ref<FormInstance>()
const currentPolicy = ref<any>(null)

const policies = ref<any[]>([])

const form = reactive({
  id: 0,
  name: '',
  description: '',
  type: 'custom',
  statements: [] as any[]
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入策略名称', trigger: 'blur' }],
  description: [{ required: true, message: '请输入描述', trigger: 'blur' }]
}

const loadPolicies = () => {
  loading.value = true
  // 模拟数据
  policies.value = [
    {
      id: 1,
      name: 'AdministratorAccess',
      description: '管理员拥有所有资源的完全访问权限',
      type: 'system',
      user_count: 1,
      role_count: 1,
      created_at: '2024-01-01 00:00:00',
      statements: [
        { effect: 'Allow', resource: '*:*', action: ['*'] }
      ]
    },
    {
      id: 2,
      name: 'CloudAccountReadOnly',
      description: '只读访问云账户资源',
      type: 'system',
      user_count: 5,
      role_count: 2,
      created_at: '2024-01-01 00:00:00',
      statements: [
        { effect: 'Allow', resource: 'cloud_account:*', action: ['list', 'get'] }
      ]
    },
    {
      id: 3,
      name: 'VMManager',
      description: '虚拟机管理员，可以管理所有虚拟机资源',
      type: 'custom',
      user_count: 3,
      role_count: 1,
      created_at: '2024-01-15 10:30:00',
      statements: [
        { effect: 'Allow', resource: 'vm:*', action: ['list', 'get', 'create', 'update', 'delete', 'action'] }
      ]
    },
    {
      id: 4,
      name: 'UserViewer',
      description: '只能查看用户信息',
      type: 'custom',
      user_count: 10,
      role_count: 0,
      created_at: '2024-02-01 14:20:00',
      statements: [
        { effect: 'Allow', resource: 'user:*', action: ['list', 'get'] }
      ]
    }
  ]
  loading.value = false
}

const handleCreate = () => {
  isEdit.value = false
  form.id = 0
  form.name = ''
  form.description = ''
  form.statements = [{ effect: 'Allow', resource: '', action: [] }]
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  isEdit.value = true
  form.id = row.id
  form.name = row.name
  form.description = row.description
  form.statements = JSON.parse(JSON.stringify(row.statements))
  dialogVisible.value = true
}

const handleView = (row: any) => {
  currentPolicy.value = row
  viewVisible.value = true
}

const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm('确定要删除该策略吗？', '提示', { type: 'warning' })
    ElMessage.success('删除成功')
    loadPolicies()
  } catch (e: any) {
    if (e !== 'cancel') {
      console.error(e)
    }
  }
}

const addStatement = () => {
  form.statements.push({ effect: 'Allow', resource: '', action: [] })
}

const removeStatement = (index: number) => {
  if (form.statements.length > 1) {
    form.statements.splice(index, 1)
  } else {
    ElMessage.warning('至少保留一条策略语句')
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    if (form.statements.length === 0) {
      ElMessage.warning('请至少添加一条策略语句')
      return
    }
    submitting.value = true
    try {
      ElMessage.success(isEdit.value ? '更新成功' : '创建成功')
      dialogVisible.value = false
      loadPolicies()
    } catch (e) {
      console.error(e)
    } finally {
      submitting.value = false
    }
  })
}

onMounted(() => {
  loadPolicies()
})
</script>

<style scoped>
.policy-page {
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

.policy-editor {
  border: 1px solid #dcdfe6;
  padding: 15px;
  border-radius: 4px;
}

.policy-detail h4 {
  margin: 20px 0 10px;
  font-size: 14px;
  color: #333;
}
</style>
