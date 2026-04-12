<template>
  <el-dialog
    v-model="visible"
    title="创建子网"
    width="550px"
    :close-on-click-modal="false"
    :destroy-on-close="true"
    @close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="formData"
      :rules="rules"
      label-width="100px"
      label-position="right"
    >
      <el-form-item label="云账号" prop="account_id">
        <CloudAccountSelector
          v-model:value="formData.account_id"
          :disabled="!!accountId"
          placeholder="请选择云账号"
          @change="handleAccountChange"
        />
      </el-form-item>

      <el-form-item label="VPC" prop="vpc_id">
        <el-select
          v-model="formData.vpc_id"
          :disabled="!!vpcId"
          :loading="vpcLoading"
          placeholder="请选择 VPC"
          filterable
          clearable
          @change="handleVPCChange"
        >
          <el-option
            v-for="vpc in vpcs"
            :key="vpc.id"
            :label="vpc.name"
            :value="vpc.id"
          >
            <span>{{ vpc.name }}</span>
            <span class="vpc-cidr-hint">{{ vpc.cidr }}</span>
          </el-option>
        </el-select>
      </el-form-item>

      <el-form-item label="名称" prop="name">
        <el-input
          v-model="formData.name"
          placeholder="请输入子网名称"
          maxlength="64"
          show-word-limit
        />
      </el-form-item>

      <el-form-item label="CIDR" prop="cidr">
        <el-input
          v-model="formData.cidr"
          placeholder="例如: 10.0.1.0/24"
        />
        <div v-if="selectedVpcCidr" class="cidr-hint">
          <el-tag size="small" type="info">VPC CIDR: {{ selectedVpcCidr }}</el-tag>
        </div>
      </el-form-item>

      <el-form-item label="可用区" prop="zone_id">
        <el-select
          v-model="formData.zone_id"
          :loading="zoneLoading"
          placeholder="请选择可用区"
          filterable
          clearable
        >
          <el-option
            v-for="zone in zones"
            :key="zone.id"
            :label="zone.name"
            :value="zone.id"
          />
        </el-select>
      </el-form-item>

      <el-form-item label="描述" prop="description">
        <el-input
          v-model="formData.description"
          type="textarea"
          :rows="3"
          placeholder="请输入描述信息"
          maxlength="256"
          show-word-limit
        />
      </el-form-item>
    </el-form>

    <template #footer>
      <span class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">
          确定
        </el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import CloudAccountSelector from '@/components/common/CloudAccountSelector.vue'
import { getVPCs, getZones, createSubnet } from '@/api/network'
import { validateCIDR, isSubnetInVPC } from '@/utils/cidr'
import type { VPC, Subnet, Zone } from '@/types'
import type { FormInstance, FormRules } from 'element-plus'

interface Props {
  visible: boolean
  accountId?: number
  vpcId?: string
}

interface Emits {
  (e: 'update:visible', val: boolean): void
  (e: 'success', subnet: Subnet): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// Form reference
const formRef = ref<FormInstance>()

// Form data
const formData = ref({
  account_id: null as number | null,
  vpc_id: '',
  name: '',
  cidr: '',
  zone_id: '',
  description: ''
})

// Loading states
const vpcLoading = ref(false)
const zoneLoading = ref(false)
const submitting = ref(false)

// Data
const vpcs = ref<VPC[]>([])
const zones = ref<Zone[]>([])

// Get selected VPC's CIDR for hint display
const selectedVpcCidr = computed(() => {
  if (!formData.value.vpc_id) return ''
  const vpc = vpcs.value.find(v => v.id === formData.value.vpc_id)
  return vpc?.cidr || ''
})

// Name validation rule
const validateName = (_rule: any, value: string, callback: (error?: Error) => void) => {
  if (!value) {
    callback(new Error('请输入子网名称'))
    return
  }
  if (value.length < 1 || value.length > 64) {
    callback(new Error('名称长度为 1-64 个字符'))
    return
  }
  const pattern = /^[a-zA-Z0-9_\-\u4e00-\u9fa5]+$/
  if (!pattern.test(value)) {
    callback(new Error('名称只能包含字母、数字、下划线、中划线和中文'))
    return
  }
  callback()
}

// CIDR validation rule
const validateCidrField = (_rule: any, value: string, callback: (error?: Error) => void) => {
  if (!value) {
    callback(new Error('请输入 CIDR'))
    return
  }
  if (!validateCIDR(value)) {
    callback(new Error('CIDR 格式不正确, 请使用如 10.0.1.0/24 的格式'))
    return
  }
  // If VPC is selected, validate subnet CIDR is within VPC CIDR
  if (selectedVpcCidr.value) {
    if (!isSubnetInVPC(value, selectedVpcCidr.value)) {
      callback(new Error(`子网 CIDR 必须在 VPC CIDR (${selectedVpcCidr.value}) 范围内`))
      return
    }
  }
  callback()
}

// Form validation rules
const rules: FormRules = {
  account_id: [
    { required: true, message: '请选择云账号', trigger: 'change' }
  ],
  vpc_id: [
    { required: true, message: '请选择 VPC', trigger: 'change' }
  ],
  name: [
    { required: true, validator: validateName, trigger: 'blur' }
  ],
  cidr: [
    { required: true, validator: validateCidrField, trigger: 'blur' }
  ],
  zone_id: [
    { required: true, message: '请选择可用区', trigger: 'change' }
  ],
  description: [
    { max: 256, message: '描述不能超过 256 个字符', trigger: 'blur' }
  ]
}

// Two-way binding for visible
const visible = computed({
  get: () => props.visible,
  set: (val: boolean) => emit('update:visible', val)
})

// Load VPCs when account changes
const loadVPCs = async (accountId: number) => {
  vpcLoading.value = true
  try {
    const res = await getVPCs({ account_id: accountId })
    vpcs.value = res || []
  } catch (error) {
    console.error('Failed to load VPCs:', error)
    ElMessage.error('加载 VPC 列表失败')
    vpcs.value = []
  } finally {
    vpcLoading.value = false
  }
}

// Load Zones when account changes
const loadZones = async (accountId: number) => {
  zoneLoading.value = true
  try {
    const res = await getZones({ account_id: accountId })
    zones.value = res || []
  } catch (error) {
    console.error('Failed to load zones:', error)
    ElMessage.error('加载可用区列表失败')
    zones.value = []
  } finally {
    zoneLoading.value = false
  }
}

// Handle account change - cascade load VPCs and Zones
const handleAccountChange = (accountId: number | null) => {
  formData.value.account_id = accountId
  // Clear dependent fields
  formData.value.vpc_id = ''
  formData.value.zone_id = ''
  formData.value.cidr = ''
  vpcs.value = []
  zones.value = []

  if (accountId) {
    loadVPCs(accountId)
    loadZones(accountId)
  }
}

// Handle VPC change
const handleVPCChange = (vpcId: string) => {
  formData.value.vpc_id = vpcId
  // Re-validate CIDR if already entered
  if (formData.value.cidr) {
    formRef.value?.validateField('cidr')
  }
}

// Watch accountId prop
watch(() => props.accountId, (newVal) => {
  if (newVal && props.visible) {
    formData.value.account_id = newVal
    loadVPCs(newVal)
    loadZones(newVal)
  }
}, { immediate: true })

// Watch vpcId prop
watch(() => props.vpcId, (newVal) => {
  if (newVal && props.visible) {
    formData.value.vpc_id = newVal
  }
}, { immediate: true })

// Watch visible change - reset form when opening
watch(() => props.visible, (newVal) => {
  if (newVal) {
    resetForm()
    if (props.accountId) {
      formData.value.account_id = props.accountId
      loadVPCs(props.accountId)
      loadZones(props.accountId)
    }
    if (props.vpcId) {
      formData.value.vpc_id = props.vpcId
    }
  }
})

// Reset form
const resetForm = () => {
  formData.value = {
    account_id: props.accountId || null,
    vpc_id: props.vpcId || '',
    name: '',
    cidr: '',
    zone_id: '',
    description: ''
  }
  vpcs.value = []
  zones.value = []
  formRef.value?.clearValidate()
}

// Close dialog
const handleClose = () => {
  visible.value = false
}

// Submit form
const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
  } catch {
    return
  }

  if (!formData.value.account_id) {
    ElMessage.error('请选择云账号')
    return
  }

  submitting.value = true
  try {
    const data: {
      account_id: number
      name: string
      vpc_id: string
      cidr: string
      zone_id: string
      description?: string
    } = {
      account_id: formData.value.account_id,
      name: formData.value.name,
      vpc_id: formData.value.vpc_id,
      cidr: formData.value.cidr,
      zone_id: formData.value.zone_id,
      description: formData.value.description || undefined
    }

    const subnet = await createSubnet(data)
    ElMessage.success('子网创建成功')
    emit('success', subnet)
    handleClose()
  } catch (error: any) {
    console.error('Failed to create subnet:', error)
    ElMessage.error(error.response?.data?.message || '子网创建失败')
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.vpc-cidr-hint {
  margin-left: 8px;
  color: #909399;
  font-size: 12px;
}

.cidr-hint {
  margin-top: 4px;
}
</style>