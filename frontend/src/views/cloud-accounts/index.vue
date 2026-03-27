<template>
  <div class="cloud-accounts-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">云账户管理</span>
          <el-button type="primary" @click="handleCreate">
            <el-icon><Plus /></el-icon>
            添加云账户
          </el-button>
        </div>
      </template>
      
      <el-table :data="accounts" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="名称" />
        <el-table-column prop="provider_type" label="云厂商" width="120">
          <template #default="{ row }">
            <el-tag :type="getProviderType(row.provider_type)">
              {{ getProviderName(row.provider_type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'danger'">
              {{ row.status === 'active' ? '正常' : row.status === 'inactive' ? '未激活' : '错误' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="描述" />
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleVerify(row)">验证</el-button>
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    
    <!-- 创建/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑云账户' : '添加云账户'"
      width="500px"
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入账户名称" />
        </el-form-item>
        <el-form-item label="云厂商" prop="provider_type">
          <el-select v-model="form.provider_type" placeholder="请选择云厂商" style="width: 100%">
            <el-option label="阿里云" value="alibaba" />
            <el-option label="腾讯云" value="tencent" />
            <el-option label="AWS" value="aws" />
            <el-option label="Azure" value="azure" />
          </el-select>
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="form.description" type="textarea" placeholder="请输入描述" />
        </el-form-item>
        <el-form-item label="Access Key" prop="credentials.access_key_id" v-if="form.provider_type === 'alibaba'">
          <el-input v-model="form.credentials.access_key_id" placeholder="请输入 Access Key ID" />
        </el-form-item>
        <el-form-item label="Secret" prop="credentials.access_key_secret" v-if="form.provider_type === 'alibaba'">
          <el-input v-model="form.credentials.access_key_secret" type="password" placeholder="请输入 Access Key Secret" />
        </el-form-item>
        <el-form-item label="Secret ID" prop="credentials.secret_id" v-if="form.provider_type === 'tencent'">
          <el-input v-model="form.credentials.secret_id" placeholder="请输入 Secret ID" />
        </el-form-item>
        <el-form-item label="Secret Key" prop="credentials.secret_key" v-if="form.provider_type === 'tencent'">
          <el-input v-model="form.credentials.secret_key" type="password" placeholder="请输入 Secret Key" />
        </el-form-item>
        <el-form-item label="Access Key" prop="credentials.access_key_id" v-if="form.provider_type === 'aws'">
          <el-input v-model="form.credentials.access_key_id" placeholder="请输入 Access Key ID" />
        </el-form-item>
        <el-form-item label="Secret Key" prop="credentials.access_key_secret" v-if="form.provider_type === 'aws'">
          <el-input v-model="form.credentials.access_key_secret" type="password" placeholder="请输入 Secret Access Key" />
        </el-form-item>
        <el-form-item label="Tenant ID" prop="credentials.tenant_id" v-if="form.provider_type === 'azure'">
          <el-input v-model="form.credentials.tenant_id" placeholder="请输入 Tenant ID" />
        </el-form-item>
        <el-form-item label="Client ID" prop="credentials.client_id" v-if="form.provider_type === 'azure'">
          <el-input v-model="form.credentials.client_id" placeholder="请输入 Client ID" />
        </el-form-item>
        <el-form-item label="Client Secret" prop="credentials.client_secret" v-if="form.provider_type === 'azure'">
          <el-input v-model="form.credentials.client_secret" type="password" placeholder="请输入 Client Secret" />
        </el-form-item>
        <el-form-item label="Subscription" prop="credentials.subscription_id" v-if="form.provider_type === 'azure'">
          <el-input v-model="form.credentials.subscription_id" placeholder="请输入 Subscription ID" />
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
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { getCloudAccounts, createCloudAccount, updateCloudAccount, deleteCloudAccount, verifyCloudAccount } from '@/api/cloud-account'
import type { CloudAccount, CreateCloudAccountRequest } from '@/types'

const accounts = ref<CloudAccount[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const formRef = ref<FormInstance>()

const form = reactive<CreateCloudAccountRequest>({
  name: '',
  provider_type: 'alibaba',
  credentials: {},
  description: ''
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
  provider_type: [{ required: true, message: '请选择云厂商', trigger: 'change' }],
  'credentials.access_key_id': [{ required: true, message: '请输入 Access Key', trigger: 'blur' }],
  'credentials.access_key_secret': [{ required: true, message: '请输入 Secret', trigger: 'blur' }]
}

const getProviderName = (type: string) => {
  const map: Record<string, string> = {
    alibaba: '阿里云',
    tencent: '腾讯云',
    aws: 'AWS',
    azure: 'Azure'
  }
  return map[type] || type
}

const getProviderType = (type: string) => {
  const map: Record<string, any> = {
    alibaba: 'orange',
    tencent: 'cyan',
    aws: 'warning',
    azure: 'success'
  }
  return map[type] || ''
}

const loadAccounts = async () => {
  loading.value = true
  try {
    const res = await getCloudAccounts()
    accounts.value = res.items || res
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const handleCreate = () => {
  isEdit.value = false
  form.name = ''
  form.provider_type = 'alibaba'
  form.credentials = {}
  form.description = ''
  dialogVisible.value = true
}

const handleEdit = (row: CloudAccount) => {
  isEdit.value = true
  form.name = row.name
  form.provider_type = row.provider_type
  form.credentials = {}
  form.description = row.description || ''
  dialogVisible.value = true
}

const handleDelete = async (row: CloudAccount) => {
  try {
    await ElMessageBox.confirm('确定要删除该云账户吗？', '提示', { type: 'warning' })
    await deleteCloudAccount(row.id)
    ElMessage.success('删除成功')
    loadAccounts()
  } catch (e: any) {
    if (e !== 'cancel') {
      console.error(e)
    }
  }
}

const handleVerify = async (row: CloudAccount) => {
  try {
    await verifyCloudAccount(row.id)
    ElMessage.success('验证成功')
  } catch (e) {
    ElMessage.error('验证失败')
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitting.value = true
    try {
      if (isEdit.value) {
        // await updateCloudAccount()
      } else {
        await createCloudAccount(form)
      }
      ElMessage.success(isEdit.value ? '更新成功' : '创建成功')
      dialogVisible.value = false
      loadAccounts()
    } catch (e) {
      console.error(e)
    } finally {
      submitting.value = false
    }
  })
}

onMounted(() => {
  loadAccounts()
})
</script>

<style scoped>
.cloud-accounts-page {
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
</style>
