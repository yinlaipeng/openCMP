<template>
  <div class="group-users-tab">
    <!-- 工具栏 -->
    <div class="tab-toolbar">
      <span class="toolbar-title">用户组成员列表</span>
      <el-button size="small" type="primary" @click="handleAddUser">
        <el-icon><Plus /></el-icon>
        添加用户
      </el-button>
    </div>

    <!-- 用户表格 -->
    <el-table :data="users" v-loading="loading" border stripe>
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="用户名" min-width="120" show-overflow-tooltip />
      <el-table-column prop="display_name" label="显示名" min-width="120" />
      <el-table-column prop="enabled" label="启用状态" width="100">
        <template #default="{ row }">
          <span class="status-dot" :data-status="row.enabled ? 'enabled' : 'disabled'">
            {{ row.enabled ? '启用' : '禁用' }}
          </span>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="100" fixed="right">
        <template #default="{ row }">
          <el-button size="small" type="danger" link @click="handleRemove(row)">移除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 添加用户弹窗 -->
    <el-dialog v-model="addDialogVisible" title="添加用户" width="500px" append-to-body>
      <el-form :model="addForm" :rules="addRules" ref="addFormRef" label-width="100px">
        <el-form-item label="用户" prop="user_id">
          <el-select v-model="addForm.user_id" placeholder="请选择用户" style="width: 100%" filterable>
            <el-option
              v-for="item in allUsers"
              :key="item.id"
              :label="`${item.name} (${item.display_name || '-'})`"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="addDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleAddSubmit" :loading="addSubmitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { getUsers, addUserToGroup, removeUserFromGroup } from '@/api/iam'

interface Props {
  groupId: number
  loading: boolean
  users: any[]
}

interface Emits {
  (e: 'refresh'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const addDialogVisible = ref(false)
const addSubmitting = ref(false)
const addFormRef = ref<FormInstance>()
const addForm = reactive({
  user_id: 0
})
const addRules: FormRules = {
  user_id: [{ required: true, message: '请选择用户', trigger: 'change' }]
}

const allUsers = ref<any[]>([])

const loadAllUsers = async () => {
  try {
    const res = await getUsers({ limit: 100 })
    allUsers.value = res.items || []
  } catch (e: any) {
    console.error(e)
  }
}

const handleAddUser = async () => {
  await loadAllUsers()
  addForm.user_id = 0
  addDialogVisible.value = true
}

const handleAddSubmit = async () => {
  if (!addFormRef.value) return
  await addFormRef.value.validate(async (valid) => {
    if (!valid) return
    addSubmitting.value = true
    try {
      await addUserToGroup(props.groupId, addForm.user_id)
      ElMessage.success('添加用户成功')
      addDialogVisible.value = false
      emit('refresh')
    } catch (e: any) {
      ElMessage.error(e.message || '添加用户失败')
    } finally {
      addSubmitting.value = false
    }
  })
}

const handleRemove = async (row: any) => {
  try {
    await ElMessageBox.confirm('确定要移除该用户吗？', '提示', { type: 'warning' })
    await removeUserFromGroup(props.groupId, row.id)
    ElMessage.success('移除用户成功')
    emit('refresh')
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '移除用户失败')
    }
  }
}

onMounted(() => {
  loadAllUsers()
})
</script>

<style scoped>
.group-users-tab {
  padding: 16px 0;
}

.tab-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.toolbar-title {
  font-size: 14px;
  color: #606266;
}

/* 状态圆点样式 */
.status-dot {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.status-dot::before {
  content: '';
  width: 8px;
  height: 8px;
  border-radius: 50%;
  display: inline-block;
}

.status-dot[data-status="enabled"]::before {
  background: #22C55E;
}

.status-dot[data-status="disabled"]::before {
  background: #64748B;
}
</style>