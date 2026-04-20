<template>
  <div class="group-projects-tab">
    <!-- 工具栏 -->
    <div class="tab-toolbar">
      <span class="toolbar-title">用户组关联的项目</span>
      <el-button size="small" type="primary" @click="handleJoinProject">
        <el-icon><Plus /></el-icon>
        加入项目
      </el-button>
    </div>

    <!-- 项目表格 -->
    <el-table :data="projects" v-loading="loading" border stripe>
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="项目名称" min-width="150" show-overflow-tooltip />
      <el-table-column prop="domain_name" label="所属域" width="120" />
      <el-table-column label="操作" width="100" fixed="right">
        <template #default="{ row }">
          <el-button size="small" type="danger" link @click="handleRemove(row)">移除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 加入项目弹窗 -->
    <el-dialog v-model="joinDialogVisible" title="加入项目" width="500px" append-to-body>
      <el-form :model="joinForm" :rules="joinRules" ref="joinFormRef" label-width="100px">
        <el-form-item label="项目" prop="project_id">
          <el-select v-model="joinForm.project_id" placeholder="请选择项目" style="width: 100%" filterable>
            <el-option v-for="item in allProjects" :key="item.id" :label="item.name" :value="item.id" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="joinDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleJoinSubmit" :loading="joinSubmitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { getProjects, addGroupToProject, removeGroupFromProject } from '@/api/iam'

interface Props {
  groupId: number
  loading: boolean
  projects: any[]
}

interface Emits {
  (e: 'refresh'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const joinDialogVisible = ref(false)
const joinSubmitting = ref(false)
const joinFormRef = ref<FormInstance>()
const joinForm = reactive({
  project_id: 0
})
const joinRules: FormRules = {
  project_id: [{ required: true, message: '请选择项目', trigger: 'change' }]
}

const allProjects = ref<any[]>([])

const loadAllProjects = async () => {
  try {
    const res = await getProjects({ limit: 100 })
    allProjects.value = res.items || []
  } catch (e: any) {
    console.error(e)
  }
}

const handleJoinProject = async () => {
  await loadAllProjects()
  joinForm.project_id = 0
  joinDialogVisible.value = true
}

const handleJoinSubmit = async () => {
  if (!joinFormRef.value) return
  await joinFormRef.value.validate(async (valid) => {
    if (!valid) return
    joinSubmitting.value = true
    try {
      await addGroupToProject(props.groupId, joinForm.project_id)
      ElMessage.success('加入项目成功')
      joinDialogVisible.value = false
      emit('refresh')
    } catch (e: any) {
      ElMessage.error(e.message || '加入项目失败')
    } finally {
      joinSubmitting.value = false
    }
  })
}

const handleRemove = async (row: any) => {
  try {
    await ElMessageBox.confirm('确定要移除该项目吗？', '提示', { type: 'warning' })
    await removeGroupFromProject(props.groupId, row.id)
    ElMessage.success('移除项目成功')
    emit('refresh')
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '移除项目失败')
    }
  }
}

onMounted(() => {
  loadAllProjects()
})
</script>

<style scoped>
.group-projects-tab {
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
</style>