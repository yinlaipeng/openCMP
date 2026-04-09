<template>
  <div class="robots-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">机器人管理</span>
          <el-button type="primary" @click="handleCreate">
            <el-icon><Plus /></el-icon>
            新建机器人
          </el-button>
        </div>
      </template>

      <el-form :inline="true" :model="filterForm" class="filter-form">
        <el-form-item label="类型">
          <el-select v-model="filterForm.type" placeholder="全部" clearable style="width: 120px">
            <el-option label="钉钉" value="dingtalk" />
            <el-option label="飞书" value="feishu" />
            <el-option label="企业微信" value="wechat" />
            <el-option label="Lark" value="lark" />
            <el-option label="Webhook" value="webhook" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadData">查询</el-button>
          <el-button @click="resetFilter">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="robots" v-loading="loading" border stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="名称" min-width="150" show-overflow-tooltip />
        <el-table-column prop="enabled" label="启用状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.enabled ? 'success' : 'info'">{{ row.enabled ? '启用' : '禁用' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="type" label="类型" width="120">
          <template #default="{ row }">
            <el-tag>{{ typeLabel(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
        <el-table-column label="操作" width="300" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" type="primary" @click="handleTest(row)">测试</el-button>
            <el-button size="small" :type="row.enabled ? 'warning' : 'success'" @click="handleToggle(row)">
              {{ row.enabled ? '禁用' : '启用' }}
            </el-button>
            <el-popconfirm title="确定删除该机器人吗？" @confirm="handleDelete(row)">
              <template #reference>
                <el-button size="small" type="danger">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next"
        class="pagination"
        @size-change="loadData"
        @current-change="loadData"
      />
    </el-card>

    <!-- 创建/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="800px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="120px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入机器人名称" />
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <el-select v-model="form.type" placeholder="请选择类型" style="width: 100%">
            <el-option label="钉钉" value="dingtalk" />
            <el-option label="飞书" value="feishu" />
            <el-option label="企业微信" value="wechat" />
            <el-option label="Lark" value="lark" />
            <el-option label="Webhook" value="webhook" />
          </el-select>
        </el-form-item>
        <el-form-item label="Webhook地址" prop="webhook_url">
          <el-input v-model="form.webhook_url" type="textarea" :rows="4" placeholder="请输入Webhook地址" />
        </el-form-item>
        <el-form-item label="密钥（如有）" prop="secret">
          <el-input v-model="form.secret" type="password" placeholder="请输入密钥" show-password />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item label="启用">
          <el-switch v-model="form.enabled" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, type FormInstance } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import {
  getRobots,
  createRobot,
  updateRobot,
  deleteRobot,
  enableRobot,
  disableRobot,
  testRobot
} from '@/api/message'

const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('新建机器人')
const robots = ref<any[]>([])
const editingId = ref<number | null>(null)
const formRef = ref<FormInstance>()

const filterForm = reactive({ type: '' })
const pagination = reactive({ page: 1, pageSize: 20, total: 0 })

const form = reactive({
  name: '',
  type: '',
  webhook_url: '',
  secret: '',
  description: '',
  enabled: true
})

const rules = {
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择类型', trigger: 'change' }],
  webhook_url: [
    { required: true, message: '请输入Webhook地址', trigger: 'blur' },
    { type: 'url', message: '请输入有效的URL', trigger: 'blur' }
  ]
}

const typeLabel = (type: string) => {
  const map: Record<string, string> = {
    dingtalk: '钉钉',
    feishu: '飞书',
    wechat: '企业微信',
    lark: 'Lark',
    webhook: 'Webhook'
  }
  return map[type] || type
}

const loadData = async () => {
  loading.value = true
  try {
    const res = await getRobots({
      type: filterForm.type || undefined,
      page: pagination.page,
      page_size: pagination.pageSize
    })
    robots.value = res.items || []
    pagination.total = res.total || 0
  } catch (e: any) {
    ElMessage.error(e.message || '加载失败')
  } finally {
    loading.value = false
  }
}

const resetFilter = () => {
  filterForm.type = ''
  pagination.page = 1
  loadData()
}

const handleCreate = () => {
  editingId.value = null
  dialogTitle.value = '新建机器人'
  form.name = ''
  form.type = ''
  form.webhook_url = ''
  form.secret = ''
  form.description = ''
  form.enabled = true
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  editingId.value = row.id
  dialogTitle.value = '编辑机器人'
  form.name = row.name
  form.type = row.type
  form.webhook_url = row.webhook_url
  form.secret = row.secret
  form.description = row.description || ''
  form.enabled = row.enabled
  dialogVisible.value = true
}

const handleSubmit = async () => {
  await formRef.value?.validate()
  submitting.value = true
  try {
    // Prepare data for submission
    const data = {
      name: form.name,
      type: form.type,
      webhook_url: form.webhook_url,
      secret: form.secret,
      description: form.description,
      enabled: form.enabled
    }

    if (editingId.value) {
      await updateRobot(editingId.value, data)
      ElMessage.success('更新成功')
    } else {
      await createRobot(data)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    loadData()
  } catch (e: any) {
    ElMessage.error(e.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

const handleTest = async (row: any) => {
  try {
    await testRobot(row.id)
    ElMessage.success('测试消息已发送')
  } catch (e: any) {
    ElMessage.error(e.message || '测试失败')
  }
}

const handleToggle = async (row: any) => {
  try {
    if (row.enabled) {
      await disableRobot(row.id)
      ElMessage.success('已禁用')
    } else {
      await enableRobot(row.id)
      ElMessage.success('已启用')
    }
    loadData()
  } catch (e: any) {
    ElMessage.error(e.message || '操作失败')
  }
}

const handleDelete = async (row: any) => {
  try {
    await deleteRobot(row.id)
    ElMessage.success('删除成功')
    loadData()
  } catch (e: any) {
    ElMessage.error(e.message || '删除失败')
  }
}

onMounted(loadData)
</script>

<style scoped>
.card-header { display: flex; justify-content: space-between; align-items: center; }
.title { font-size: 16px; font-weight: bold; }
.filter-form { margin-bottom: 16px; }
.pagination { margin-top: 16px; justify-content: flex-end; }
</style>
