<template>
  <div class="policy-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">策略管理</span>
          <el-button type="primary" @click="handleCreate">
            <el-icon><Plus /></el-icon>
            创建策略
          </el-button>
        </div>
      </template>
      
      <el-table :data="policies" v-loading="loading">
        <el-table-column prop="id" label="策略 ID" width="100" />
        <el-table-column prop="name" label="策略名称" width="200" />
        <el-table-column prop="description" label="描述" />
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="row.type === 'system' ? 'warning' : ''">
              {{ row.type === 'system' ? '系统' : '自定义' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="user_count" label="关联用户" width="100" />
        <el-table-column prop="role_count" label="关联角色" width="100" />
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleView(row)">查看</el-button>
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
    </el-card>
    
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑策略' : '创建策略'" width="800px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="120px">
        <el-form-item label="策略名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入策略名称" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="form.description" type="textarea" :rows="3" placeholder="策略描述" />
        </el-form-item>
        <el-form-item label="策略内容" prop="statement">
          <div class="policy-editor">
            <el-table :data="form.statements" border style="margin-bottom: 10px">
              <el-table-column label="效果" width="120">
                <template #default="{ row, $index }">
                  <el-select v-model="row.effect" size="small" style="width: 100%">
                    <el-option label="允许" value="Allow" />
                    <el-option label="拒绝" value="Deny" />
                  </el-select>
                </template>
              </el-table-column>
              <el-table-column label="资源类型">
                <template #default="{ row, $index }">
                  <el-select v-model="row.resource" size="small" style="width: 100%">
                    <el-option label="云账户" value="cloud_account:*" />
                    <el-option label="虚拟机" value="vm:*" />
                    <el-option label="网络资源" value="network:*" />
                    <el-option label="存储" value="storage:*" />
                    <el-option label="用户" value="user:*" />
                    <el-option label="角色" value="role:*" />
                    <el-option label="所有资源" value="*:*" />
                  </el-select>
                </template>
              </el-table-column>
              <el-table-column label="操作">
                <template #default="{ row, $index }">
                  <el-select v-model="row.action" size="small" style="width: 100%" multiple>
                    <el-option label="查看列表" value="list" />
                    <el-option label="查看详情" value="get" />
                    <el-option label="创建" value="create" />
                    <el-option label="更新" value="update" />
                    <el-option label="删除" value="delete" />
                    <el-option label="所有操作" value="*" />
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
    
    <el-dialog v-model="viewVisible" title="策略详情" width="700px">
      <div v-if="currentPolicy" class="policy-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="策略 ID">{{ currentPolicy.id }}</el-descriptions-item>
          <el-descriptions-item label="类型">
            <el-tag :type="currentPolicy.type === 'system' ? 'warning' : ''">
              {{ currentPolicy.type === 'system' ? '系统' : '自定义' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="策略名称" :span="2">{{ currentPolicy.name }}</el-descriptions-item>
          <el-descriptions-item label="描述" :span="2">{{ currentPolicy.description }}</el-descriptions-item>
          <el-descriptions-item label="关联用户">{{ currentPolicy.user_count }}</el-descriptions-item>
          <el-descriptions-item label="关联角色">{{ currentPolicy.role_count }}</el-descriptions-item>
          <el-descriptions-item label="创建时间" :span="2">{{ currentPolicy.created_at }}</el-descriptions-item>
        </el-descriptions>
        
        <h4>策略内容</h4>
        <el-table :data="currentPolicy.statements" border>
          <el-table-column prop="effect" label="效果">
            <template #default="{ row }">
              <el-tag :type="row.effect === 'Allow' ? 'success' : 'danger'">
                {{ row.effect === 'Allow' ? '允许' : '拒绝' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="resource" label="资源" />
          <el-table-column prop="action" label="操作">
            <template #default="{ row }">
              <el-tag v-for="action in row.action" :key="action" size="small" style="margin-right: 5px">
                {{ action }}
              </el-tag>
            </template>
          </el-table-column>
        </el-table>
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
