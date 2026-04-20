<template>
  <el-dialog
    v-model="visible"
    title="属性设置"
    width="700px"
    :close-on-click-modal="false"
  >
    <el-tabs v-model="activeTab">
      <el-tab-pane label="基本属性" name="basic">
        <el-form :model="form" label-width="120px" v-loading="loading">
          <el-form-item label="备注">
            <el-input v-model="form.description" type="textarea" :rows="3" placeholder="请输入备注信息" />
          </el-form-item>

          <el-form-item label="所属域">
            <el-select v-model="form.domainId" placeholder="请选择所属域">
              <el-option v-for="domain in domains" :key="domain.id" :label="domain.name" :value="domain.id" />
            </el-select>
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <el-tab-pane label="同步属性" name="sync">
        <el-form :model="form" label-width="120px">
          <el-form-item label="自动同步">
            <el-switch v-model="form.autoSync" />
          </el-form-item>

          <el-form-item label="同步间隔">
            <el-input-number v-model="form.syncInterval" :min="1" :max="24" />
            <span style="margin-left: 8px">小时</span>
          </el-form-item>

          <el-form-item label="同步资源类型">
            <el-checkbox-group v-model="form.syncResourceTypes">
              <el-checkbox label="vm">虚拟机</el-checkbox>
              <el-checkbox label="disk">磁盘</el-checkbox>
              <el-checkbox label="network">网络</el-checkbox>
              <el-checkbox label="database">数据库</el-checkbox>
              <el-checkbox label="storage">存储桶</el-checkbox>
              <el-checkbox label="image">镜像</el-checkbox>
            </el-checkbox-group>
          </el-form-item>

          <el-form-item label="绑定同步策略">
            <el-select v-model="form.syncPolicyId" placeholder="请选择同步策略" clearable>
              <el-option v-for="policy in syncPolicies" :key="policy.id" :label="policy.name" :value="policy.id" />
            </el-select>
          </el-form-item>
        </el-form>
      </el-tab-pane>
    </el-tabs>

    <template #footer>
      <el-button type="primary" @click="handleSave" :loading="saving">保存</el-button>
      <el-button @click="visible = false">取消</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { getCloudAccount, updateCloudAccount, updateCloudAccountAttributes } from '@/api/cloud-account'
import { getSyncPolicies } from '@/api/sync-policy'
import { getDomains } from '@/api/iam'
import type { SyncPolicy } from '@/types/sync-policy'

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
const activeTab = ref('basic')
const domains = ref<{ id: number; name: string }[]>([])
const syncPolicies = ref<SyncPolicy[]>([])

const form = reactive({
  description: '',
  domainId: 1,
  autoSync: false,
  syncInterval: 1,
  syncResourceTypes: ['vm', 'disk', 'network'] as string[],
  syncPolicyId: null as number | null
})

const loadAccount = async () => {
  if (!props.accountId) return
  loading.value = true
  try {
    const account = await getCloudAccount(props.accountId)
    form.description = account.description || ''
    form.domainId = account.domain_id
    form.syncPolicyId = account.sync_policy_id
  } catch (error) {
    ElMessage.error('加载云账号失败')
  } finally {
    loading.value = false
  }
}

const loadDomains = async () => {
  try {
    const res = await getDomains()
    domains.value = (res.items || []).map(d => ({ id: d.id, name: d.name }))
  } catch (error) {
    console.error('Failed to load domains:', error)
  }
}

const loadSyncPolicies = async () => {
  try {
    const res = await getSyncPolicies()
    syncPolicies.value = res.items || []
  } catch (error) {
    console.error('Failed to load sync policies:', error)
  }
}

const handleSave = async () => {
  if (!props.accountId) return
  saving.value = true
  try {
    await updateCloudAccount(props.accountId, {
      description: form.description,
      domain_id: form.domainId
    })

    await updateCloudAccountAttributes(props.accountId, {
      sync_policy_id: form.syncPolicyId,
      sync_resource_types: form.syncResourceTypes
    })

    ElMessage.success('属性设置成功')
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
    loadAccount()
    loadDomains()
    loadSyncPolicies()
  }
}, { immediate: true })
</script>