<template>
  <el-dialog
    v-model="visible"
    title="设置同步归属策略"
    width="600px"
    :close-on-click-modal="false"
  >
    <el-form :model="form" label-width="150px" v-loading="loading">
      <el-form-item label="资源归属方式">
        <el-checkbox-group v-model="form.resourceAssignmentMethods">
          <el-checkbox label="tag_mapping">根据同步策略归属</el-checkbox>
          <el-checkbox label="project_mapping">根据云上项目归属</el-checkbox>
          <el-checkbox label="subscription_mapping">根据云订阅归属</el-checkbox>
          <el-checkbox label="specify_project">指定项目</el-checkbox>
        </el-checkbox-group>
      </el-form-item>

      <el-form-item label="指定项目" v-if="form.resourceAssignmentMethods.includes('specify_project')">
        <el-select v-model="form.specifyProjectId" placeholder="请选择项目" filterable>
          <el-option v-for="project in projects" :key="project.id" :label="project.name" :value="project.id" />
        </el-select>
      </el-form-item>

      <el-form-item label="屏蔽同步资源">
        <el-switch v-model="form.blockSyncResources" />
      </el-form-item>

      <el-form-item label="屏蔽资源类型" v-if="form.blockSyncResources">
        <el-checkbox-group v-model="form.blockedResourceTypes">
          <el-checkbox label="vm">虚拟机</el-checkbox>
          <el-checkbox label="disk">磁盘</el-checkbox>
          <el-checkbox label="network">网络</el-checkbox>
          <el-checkbox label="database">数据库</el-checkbox>
          <el-checkbox label="storage">存储桶</el-checkbox>
        </el-checkbox-group>
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button type="primary" @click="handleSave" :loading="saving">保存</el-button>
      <el-button @click="visible = false">取消</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { getCloudAccount, updateCloudAccountAttributes } from '@/api/cloud-account'
import { getProjects } from '@/api/project'
import type { Project } from '@/types'

interface Props {
  modelValue: boolean
  accountId: number | null
}

interface Emits {
  (e: 'update:modelValue', value: boolean): void
  (e: 'saved'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const loading = ref(false)
const saving = ref(false)
const projects = ref<Project[]>([])

const form = reactive({
  resourceAssignmentMethods: ['tag_mapping'] as string[],
  specifyProjectId: null as number | null,
  blockSyncResources: false,
  blockedResourceTypes: [] as string[]
})

// 加载项目列表
const loadProjects = async () => {
  try {
    const res = await getProjects()
    projects.value = res.items || []
  } catch (error) {
    console.error('Failed to load projects:', error)
  }
}

// 加载云账号当前设置
const loadAccountSettings = async () => {
  if (!props.accountId) return
  loading.value = true
  try {
    const account = await getCloudAccount(props.accountId)
    // 根据当前设置初始化表单
    if (account.resource_assignment_method) {
      form.resourceAssignmentMethods = [account.resource_assignment_method]
    }
  } catch (error) {
    ElMessage.error('加载云账号设置失败')
  } finally {
    loading.value = false
  }
}

// 保存设置
const handleSave = async () => {
  if (!props.accountId) return
  saving.value = true
  try {
    // 确定主要的资源归属方式
    const primaryMethod = form.resourceAssignmentMethods[0] || 'tag_mapping'
    await updateCloudAccountAttributes(props.accountId, {
      resource_assignment_method: primaryMethod
    })
    ElMessage.success('同步归属策略设置成功')
    emit('saved')
    visible.value = false
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

watch(() => props.accountId, (id) => {
  if (id && visible.value) {
    loadAccountSettings()
    loadProjects()
  }
}, { immediate: true })

// 需要导入reactive
import { reactive } from 'vue'
</script>

<style scoped>
.el-checkbox-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
</style>