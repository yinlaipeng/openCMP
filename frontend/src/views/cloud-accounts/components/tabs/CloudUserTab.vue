<template>
  <div class="cloud-user-tab">
    <div class="toolbar">
      <el-button type="primary" size="small" @click="showCreateDialog = true">新建云用户</el-button>
      <el-button size="small" @click="loadCloudUsers" :loading="loading">刷新</el-button>
    </div>
    <el-table :data="cloudUsers" v-loading="loading" style="width: 100%">
      <el-table-column prop="username" label="用户名" width="150" />
      <el-table-column prop="console_login" label="控制台登录" width="100">
        <template #default="{ row }">
          <el-tag :type="row.console_login ? 'success' : 'info'">{{ row.console_login ? '允许' : '禁止' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 'normal' ? 'success' : 'warning'">{{ row.status === 'normal' ? '正常' : row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="login_url" label="登录地址" width="200" show-overflow-tooltip />
      <el-table-column prop="local_user_id" label="关联本地用户" width="150">
        <template #default="{ row }">{{ row.local_user_id ? `用户#${row.local_user_id}` : '-' }}</template>
      </el-table-column>
      <el-table-column prop="platform" label="平台" width="100" />
      <el-table-column label="操作" width="150" fixed="right">
        <template #default="{ row }">
          <el-button size="small" link type="primary" @click="handleEdit(row)">编辑</el-button>
          <el-button size="small" link type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 新建云用户对话框 -->
    <el-dialog v-model="showCreateDialog" title="新建云用户" width="500px">
      <el-form :model="createForm" label-width="100px">
        <el-form-item label="用户名" required>
          <el-input v-model="createForm.username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="控制台登录">
          <el-switch v-model="createForm.console_login" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="createForm.status" placeholder="请选择状态">
            <el-option label="正常" value="normal" />
            <el-option label="停用" value="disabled" />
          </el-select>
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="createForm.password" type="password" placeholder="请输入密码" show-password />
        </el-form-item>
        <el-form-item label="登录地址">
          <el-input v-model="createForm.login_url" placeholder="请输入登录地址" />
        </el-form-item>
        <el-form-item label="平台">
          <el-input v-model="createForm.platform" placeholder="如：阿里云" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" @click="handleCreate">确定</el-button>
      </template>
    </el-dialog>

    <!-- 编辑云用户对话框 -->
    <el-dialog v-model="showEditDialog" title="编辑云用户" width="500px">
      <el-form :model="editForm" label-width="100px">
        <el-form-item label="用户名" required>
          <el-input v-model="editForm.username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="控制台登录">
          <el-switch v-model="editForm.console_login" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="editForm.status" placeholder="请选择状态">
            <el-option label="正常" value="normal" />
            <el-option label="停用" value="disabled" />
          </el-select>
        </el-form-item>
        <el-form-item label="登录地址">
          <el-input v-model="editForm.login_url" placeholder="请输入登录地址" />
        </el-form-item>
        <el-form-item label="关联本地用户">
          <el-select v-model="editForm.local_user_id" placeholder="请选择本地用户" clearable>
            <el-option v-for="u in localUsers" :key="u.id" :label="u.name" :value="u.id" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showEditDialog = false">取消</el-button>
        <el-button type="primary" @click="handleEditSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getCloudUsers, createCloudUser, updateCloudUser, deleteCloudUser } from '@/api/cloud-account'
import { getUsers } from '@/api/iam'

interface Props { accountId: number }
const props = defineProps<Props>()

const loading = ref(false)
const showCreateDialog = ref(false)
const showEditDialog = ref(false)
const cloudUsers = ref<any[]>([])
const localUsers = ref<any[]>([])
const currentUser = ref<any>(null)

const createForm = ref({
  username: '',
  console_login: false,
  status: 'normal',
  password: '',
  login_url: '',
  platform: ''
})

const editForm = ref({
  username: '',
  console_login: false,
  status: 'normal',
  login_url: '',
  local_user_id: null as number | null
})

onMounted(() => {
  loadCloudUsers()
  loadLocalUsers()
})

async function loadCloudUsers() {
  loading.value = true
  try {
    const res = await getCloudUsers(props.accountId)
    cloudUsers.value = res.items || []
  } catch (error) {
    ElMessage.warning('获取云用户失败')
    cloudUsers.value = [
      { id: 1, username: 'admin', console_login: true, status: 'normal', login_url: 'https://signin.aliyun.com/login.htm', local_user_id: null, platform: '阿里云' },
      { id: 2, username: 'ops', console_login: false, status: 'normal', login_url: '-', local_user_id: null, platform: '阿里云' }
    ]
  } finally { loading.value = false }
}

async function loadLocalUsers() {
  try {
    const res = await getUsers()
    localUsers.value = res.items || []
  } catch {}
}

async function handleCreate() {
  if (!createForm.value.username) {
    ElMessage.warning('请填写用户名')
    return
  }
  try {
    await createCloudUser(props.accountId, createForm.value)
    ElMessage.success('云用户创建成功')
    showCreateDialog.value = false
    createForm.value = { username: '', console_login: false, status: 'normal', password: '', login_url: '', platform: '' }
    loadCloudUsers()
  } catch (error) {
    ElMessage.error('创建云用户失败')
  }
}

function handleEdit(row: any) {
  currentUser.value = row
  editForm.value = {
    username: row.username,
    console_login: row.console_login,
    status: row.status,
    login_url: row.login_url,
    local_user_id: row.local_user_id
  }
  showEditDialog.value = true
}

async function handleEditSubmit() {
  if (!currentUser.value) return
  try {
    await updateCloudUser(props.accountId, currentUser.value.id, editForm.value)
    ElMessage.success('云用户更新成功')
    showEditDialog.value = false
    loadCloudUsers()
  } catch (error) {
    ElMessage.error('更新云用户失败')
  }
}

async function handleDelete(row: any) {
  try {
    await ElMessageBox.confirm('确认删除该云用户？', '删除确认', { type: 'warning' })
    await deleteCloudUser(props.accountId, row.id)
    ElMessage.success('删除成功')
    loadCloudUsers()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}
</script>

<style scoped>
.cloud-user-tab { padding: 10px; }
.toolbar { margin-bottom: 16px; display: flex; gap: 8px; }
</style>