<template>
  <div class="joined-groups-container">
    <!-- 工具栏 -->
    <div class="toolbar">
      <el-button size="small" @click="handleRefresh" :loading="loading">
        <el-icon><Refresh /></el-icon>
        刷新
      </el-button>
      <el-button size="small" @click="handleJoinGroup">
        <el-icon><Plus /></el-icon>
        加入组
      </el-button>
      <el-button
        size="small"
        :disabled="selectedGroups.length === 0"
        @click="handleLeave"
      >
        退出组
      </el-button>
    </div>

    <!-- 组列表 -->
    <el-table
      :data="groups"
      v-loading="loading"
      @selection-change="handleSelectionChange"
      row-key="id"
      border
      stripe
      class="group-table"
    >
      <el-table-column type="selection" width="55" />
      <el-table-column prop="name" label="组" min-width="150" show-overflow-tooltip />
      <el-table-column prop="domain_name" label="域" width="120" />
      <el-table-column label="操作" width="80" fixed="right">
        <template #default="{ row }">
          <el-button size="small" type="danger" link @click="handleLeaveSingle(row)">
            退出
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 空状态 -->
    <EmptyState
      v-if="!loading && groups.length === 0"
      icon="UserFilled"
      title="暂无数据"
      description="该用户尚未加入任何组"
      action-text="加入组"
      @action="handleJoinGroup"
    />

    <!-- 加入组弹窗 -->
    <el-dialog v-model="joinDialogVisible" title="加入组" width="400px">
      <el-form :model="joinForm" label-width="80px">
        <el-form-item label="选择组">
          <el-select v-model="joinForm.group_id" placeholder="请选择组" style="width: 100%">
            <el-option v-for="g in allGroups" :key="g.id" :label="g.name" :value="g.id" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="joinDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleJoinSubmit" :loading="joinLoading">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Refresh, Plus } from '@element-plus/icons-vue'
import { getUserGroups, joinGroup, leaveGroup, getGroups } from '@/api/iam'
import type { Group } from '@/types/iam'
import EmptyState from '@/components/common/EmptyState.vue'

const props = defineProps<{
  userId?: number
}>()

const emit = defineEmits<{
  (e: 'refresh'): void
}>()

// State
const loading = ref(false)
const joinLoading = ref(false)
const groups = ref<Group[]>([])
const selectedGroups = ref<Group[]>([])
const allGroups = ref<Group[]>([])
const joinDialogVisible = ref(false)

const joinForm = ref({
  group_id: 0
})

// Load data
watch(() => props.userId, (id) => {
  if (id) loadGroups()
}, { immediate: true })

const loadGroups = async () => {
  if (!props.userId) return
  loading.value = true
  try {
    const res = await getUserGroups(props.userId)
    groups.value = res.items || []
  } catch (e: any) {
    ElMessage.error(e.message || '加载组列表失败')
  } finally {
    loading.value = false
  }
}

const loadAllGroups = async () => {
  try {
    const res = await getGroups({ limit: 100 })
    allGroups.value = res.items || []
  } catch (e) {
    console.error('加载组失败', e)
  }
}

// Handlers
const handleRefresh = () => loadGroups()

const handleSelectionChange = (selection: Group[]) => {
  selectedGroups.value = selection
}

const handleJoinGroup = async () => {
  await loadAllGroups()
  joinForm.value.group_id = 0
  joinDialogVisible.value = true
}

const handleJoinSubmit = async () => {
  if (!props.userId || !joinForm.value.group_id) {
    ElMessage.error('请选择组')
    return
  }
  joinLoading.value = true
  try {
    await joinGroup(props.userId, joinForm.value.group_id)
    ElMessage.success('加入组成功')
    joinDialogVisible.value = false
    loadGroups()
    emit('refresh')
  } catch (e: any) {
    ElMessage.error(e.message || '加入组失败')
  } finally {
    joinLoading.value = false
  }
}

const handleLeave = async () => {
  try {
    await ElMessageBox.confirm(`确定要退出选中的 ${selectedGroups.value.length} 个组吗？`, '提示', { type: 'warning' })
    for (const g of selectedGroups.value) {
      await leaveGroup(props.userId!, g.id)
    }
    ElMessage.success('退出成功')
    loadGroups()
    emit('refresh')
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '退出失败')
    }
  }
}

const handleLeaveSingle = async (row: Group) => {
  try {
    await ElMessageBox.confirm('确定要退出该组吗？', '提示', { type: 'warning' })
    await leaveGroup(props.userId!, row.id)
    ElMessage.success('退出成功')
    loadGroups()
    emit('refresh')
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '退出失败')
    }
  }
}

onMounted(() => loadAllGroups())
</script>

<style scoped>
.joined-groups-container {
  padding: var(--space-4, 16px);
}

.toolbar {
  display: flex;
  gap: var(--space-2, 8px);
  margin-bottom: var(--space-3, 12px);
}

.group-table {
  margin-bottom: var(--space-4, 16px);
}
</style>