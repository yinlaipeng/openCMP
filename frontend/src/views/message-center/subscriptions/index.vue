<template>
  <div class="subscriptions-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">消息订阅</span>
          <el-button type="primary" @click="handleCreate">
            <el-icon><Plus /></el-icon>
            新建订阅
          </el-button>
        </div>
      </template>

      <el-form :inline="true" :model="filterForm" class="filter-form">
        <el-form-item label="用户 ID">
          <el-input v-model="filterForm.user_id" placeholder="用户 ID" clearable style="width: 120px" @keyup.enter="loadData" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadData">查询</el-button>
          <el-button @click="resetFilter">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="subscriptions" v-loading="loading" border stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="user_id" label="用户 ID" width="100" />
        <el-table-column prop="message_type_id" label="消息类型 ID" width="130" />
        <el-table-column label="通知渠道" min-width="200">
          <template #default="{ row }">
            <el-tag v-if="row.inbox" size="small" class="channel-tag">站内信</el-tag>
            <el-tag v-if="row.email" size="small" type="success" class="channel-tag">邮件</el-tag>
            <el-tag v-if="row.wechat" size="small" type="warning" class="channel-tag">企业微信</el-tag>
            <el-tag v-if="row.dingtalk" size="small" type="danger" class="channel-tag">钉钉</el-tag>
            <el-tag v-if="row.webhook" size="small" type="info" class="channel-tag">Webhook</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-popconfirm title="确定删除该订阅吗？" @confirm="handleDelete(row)">
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

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="110px">
        <el-form-item label="用户 ID" prop="user_id">
          <el-input-number v-model="form.user_id" :min="1" style="width: 100%" />
        </el-form-item>
        <el-form-item label="消息类型" prop="message_type_id">
          <el-select v-model="form.message_type_id" placeholder="请选择消息类型" style="width: 100%">
            <el-option
              v-for="mt in messageTypes"
              :key="mt.id"
              :label="mt.display_name || mt.name"
              :value="mt.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="通知渠道">
          <el-checkbox v-model="form.inbox">站内信</el-checkbox>
          <el-checkbox v-model="form.email">邮件</el-checkbox>
          <el-checkbox v-model="form.wechat">企业微信</el-checkbox>
          <el-checkbox v-model="form.dingtalk">钉钉</el-checkbox>
          <el-checkbox v-model="form.webhook">Webhook</el-checkbox>
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
import { getSubscriptions, createSubscription, updateSubscription, deleteSubscription, getMessageTypes } from '@/api/message'

const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('新建订阅')
const subscriptions = ref<any[]>([])
const messageTypes = ref<any[]>([])
const editingId = ref<number | null>(null)
const formRef = ref<FormInstance>()

const filterForm = reactive({ user_id: '' })
const pagination = reactive({ page: 1, pageSize: 20, total: 0 })

const form = reactive({
  user_id: 0,
  message_type_id: null as number | null,
  inbox: true,
  email: false,
  wechat: false,
  dingtalk: false,
  webhook: false
})

const rules = {
  user_id: [{ required: true, type: 'number', min: 1, message: '请输入用户 ID', trigger: 'blur' }],
  message_type_id: [{ required: true, message: '请选择消息类型', trigger: 'change' }]
}

const loadData = async () => {
  loading.value = true
  try {
    const params: any = { page: pagination.page, page_size: pagination.pageSize }
    if (filterForm.user_id) params.user_id = Number(filterForm.user_id)
    const res = await getSubscriptions(params)
    subscriptions.value = res.items || []
    pagination.total = res.total || 0
  } catch (e: any) {
    ElMessage.error(e.message || '加载失败')
  } finally {
    loading.value = false
  }
}

const resetFilter = () => { filterForm.user_id = ''; pagination.page = 1; loadData() }

const loadMessageTypes = async () => {
  try {
    const res = await getMessageTypes()
    messageTypes.value = res.items || []
  } catch { /* ignore */ }
}

const handleCreate = () => {
  editingId.value = null
  dialogTitle.value = '新建订阅'
  Object.assign(form, { user_id: 0, message_type_id: null, inbox: true, email: false, wechat: false, dingtalk: false, webhook: false })
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  editingId.value = row.id
  dialogTitle.value = '编辑订阅'
  Object.assign(form, {
    user_id: row.user_id,
    message_type_id: row.message_type_id,
    inbox: row.inbox ?? true,
    email: row.email ?? false,
    wechat: row.wechat ?? false,
    dingtalk: row.dingtalk ?? false,
    webhook: row.webhook ?? false
  })
  dialogVisible.value = true
}

const handleSubmit = async () => {
  await formRef.value?.validate()
  submitting.value = true
  try {
    if (editingId.value) {
      await updateSubscription(editingId.value, form)
      ElMessage.success('更新成功')
    } else {
      await createSubscription(form)
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

const handleDelete = async (row: any) => {
  try { await deleteSubscription(row.id); ElMessage.success('删除成功'); loadData() }
  catch (e: any) { ElMessage.error(e.message || '删除失败') }
}

onMounted(() => { loadData(); loadMessageTypes() })
</script>

<style scoped>
.card-header { display: flex; justify-content: space-between; align-items: center; }
.title { font-size: 16px; font-weight: bold; }
.filter-form { margin-bottom: 16px; }
.pagination { margin-top: 16px; justify-content: flex-end; }
.channel-tag { margin-right: 4px; }
</style>
