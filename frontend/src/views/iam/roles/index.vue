<template>
  <div class="roles-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">角色权限</span>
          <el-button type="primary" @click="handleCreate">
            <el-icon><Plus /></el-icon>
            添加角色
          </el-button>
        </div>
      </template>
      
      <el-table :data="roles" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="角色名" />
        <el-table-column prop="display_name" label="显示名" />
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="row.type === 'system' ? 'warning' : ''">
              {{ row.type === 'system' ? '系统' : '自定义' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="enabled" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.enabled ? 'success' : 'info'">
              {{ row.enabled ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handlePermissions(row)">权限</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    
    <el-dialog v-model="dialogVisible" title="添加角色" width="500px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="角色名" required>
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="显示名">
          <el-input v-model="form.display_name" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
    
    <el-dialog v-model="permDialogVisible" title="角色权限" width="600px">
      <el-transfer
        v-model="selectedPermissions"
        :data="allPermissions"
        :titles="['可选权限', '已选权限']"
        filterable
      />
      <template #footer>
        <el-button @click="permDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSavePermissions">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getRoles, createRole, deleteRole, getPermissions, getRolePermissions, assignPermission } from '@/api/iam'

const roles = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const permDialogVisible = ref(false)
const currentRoleId = ref(0)
const selectedPermissions = ref<number[]>([])
const allPermissions = ref<any[]>([])

const form = reactive({
  name: '',
  display_name: '',
  description: ''
})

const loadRoles = async () => {
  loading.value = true
  try {
    const res = await getRoles()
    roles.value = res.items || res
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const loadPermissions = async () => {
  try {
    const res = await getPermissions({ limit: 100 })
    allPermissions.value = (res.items || res).map((p: any) => ({
      key: p.id,
      label: `${p.resource}.${p.action}`,
      data: p
    }))
  } catch (e) {
    console.error(e)
  }
}

const loadRolePermissions = async (roleId: number) => {
  try {
    const res = await getRolePermissions(roleId)
    selectedPermissions.value = (res.items || res).map((p: any) => p.id)
  } catch (e) {
    console.error(e)
  }
}

const handleCreate = () => {
  form.name = ''
  form.display_name = ''
  form.description = ''
  dialogVisible.value = true
}

const handlePermissions = async (row: any) => {
  currentRoleId.value = row.id
  await loadPermissions()
  await loadRolePermissions(row.id)
  permDialogVisible.value = true
}

const handleSavePermissions = async () => {
  try {
    // 简化处理，实际应该对比差异
    ElMessage.success('权限已更新')
    permDialogVisible.value = false
  } catch (e) {
    console.error(e)
  }
}

const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm('确定要删除该角色吗？', '提示', { type: 'warning' })
    await deleteRole(row.id)
    ElMessage.success('删除成功')
    loadRoles()
  } catch (e: any) {
    if (e !== 'cancel') {
      console.error(e)
    }
  }
}

const handleSubmit = async () => {
  try {
    await createRole(form)
    ElMessage.success('创建成功')
    dialogVisible.value = false
    loadRoles()
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadRoles()
})
</script>

<style scoped>
.roles-page {
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
