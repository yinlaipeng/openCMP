<template>
  <div class="cloud-user-group-tab">
    <div class="toolbar">
      <el-button type="primary" size="small" @click="showCreateDialog = true">新建云用户组</el-button>
      <el-button size="small" @click="loadCloudUserGroups" :loading="loading">刷新</el-button>
    </div>
    <el-table :data="cloudUserGroups" v-loading="loading" style="width: 100%">
      <el-table-column prop="name" label="名称" width="150" />
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 'normal' ? 'success' : 'warning'">{{ row.status === 'normal' ? '正常' : row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="permissions" label="权限" width="150" />
      <el-table-column prop="platform" label="平台" width="100" />
      <el-table-column prop="domain_id" label="所属域" width="100">
        <template #default="{ row }">{{ row.domain_id ? `域#${row.domain_id}` : '默认域' }}</template>
      </el-table-column>
      <el-table-column label="操作" width="150" fixed="right">
        <template #default="{ row }">
          <el-button size="small" link type="primary" @click="handleEdit(row)">编辑</el-button>
          <el-button size="small" link type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 新建云用户组对话框 -->
    <el-dialog v-model="showCreateDialog" title="新建云用户组" width="500px">
      <el-form :model="createForm" label-width="100px">
        <el-form-item label="名称" required>
          <el-input v-model="createForm.name" placeholder="请输入名称" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="createForm.status" placeholder="请选择状态">
            <el-option label="正常" value="normal" />
            <el-option label="停用" value="disabled" />
          </el-select>
        </el-form-item>
        <el-form-item label="权限">
          <el-input v-model="createForm.permissions" placeholder="如：管理员、运维" />
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

    <!-- 编辑云用户组对话框 -->
    <el-dialog v-model="showEditDialog" title="编辑云用户组" width="500px">
      <el-form :model="editForm" label-width="100px">
        <el-form-item label="名称" required>
          <el-input v-model="editForm.name" placeholder="请输入名称" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="editForm.status" placeholder="请选择状态">
            <el-option label="正常" value="normal" />
            <el-option label="停用" value="disabled" />
          </el-select>
        </el-form-item>
        <el-form-item label="权限">
          <el-input v-model="editForm.permissions" placeholder="如：管理员、运维" />
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
import { getCloudUserGroups, createCloudUserGroup, updateCloudUserGroup, deleteCloudUserGroup } from '@/api/cloud-account'

interface Props { accountId: number }
const props = defineProps<Props>()

const loading = ref(false)
const showCreateDialog = ref(false)
const showEditDialog = ref(false)
const cloudUserGroups = ref<any[]>([])
const currentGroup = ref<any>(null)

const createForm = ref({
  name: '',
  status: 'normal',
  permissions: '',
  platform: ''
})

const editForm = ref({
  name: '',
  status: 'normal',
  permissions: ''
})

onMounted(() => { loadCloudUserGroups() })

async function loadCloudUserGroups() {
  loading.value = true
  try {
    const res = await getCloudUserGroups(props.accountId)
    cloudUserGroups.value = res.items || []
  } catch (error) {
    ElMessage.warning('获取云用户组失败')
    cloudUserGroups.value = [
      { id: 1, name: 'Admins', status: 'normal', permissions: '管理员', platform: '阿里云', domain_id: 1 },
      { id: 2, name: 'DevOps', status: 'normal', permissions: '运维', platform: '阿里云', domain_id: 1 }
    ]
  } finally { loading.value = false }
}

async function handleCreate() {
  if (!createForm.value.name) {
    ElMessage.warning('请填写名称')
    return
  }
  try {
    await createCloudUserGroup(props.accountId, createForm.value)
    ElMessage.success('云用户组创建成功')
    showCreateDialog.value = false
    createForm.value = { name: '', status: 'normal', permissions: '', platform: '' }
    loadCloudUserGroups()
  } catch (error) {
    ElMessage.error('创建云用户组失败')
  }
}

function handleEdit(row: any) {
  currentGroup.value = row
  editForm.value = {
    name: row.name,
    status: row.status,
    permissions: row.permissions
  }
  showEditDialog.value = true
}

async function handleEditSubmit() {
  if (!currentGroup.value) return
  try {
    await updateCloudUserGroup(props.accountId, currentGroup.value.id, editForm.value)
    ElMessage.success('云用户组更新成功')
    showEditDialog.value = false
    loadCloudUserGroups()
  } catch (error) {
    ElMessage.error('更新云用户组失败')
  }
}

async function handleDelete(row: any) {
  try {
    await ElMessageBox.confirm('确认删除该云用户组？', '删除确认', { type: 'warning' })
    await deleteCloudUserGroup(props.accountId, row.id)
    ElMessage.success('删除成功')
    loadCloudUserGroups()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}
</script>

<style scoped>
.cloud-user-group-tab { padding: 10px; }
.toolbar { margin-bottom: 16px; display: flex; gap: 8px; }
</style>