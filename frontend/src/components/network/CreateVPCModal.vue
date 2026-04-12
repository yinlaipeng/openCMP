<template>
  <el-dialog
    v-model="visible"
    title="创建 VPC"
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

      <el-form-item label="名称" prop="name">
        <el-input
          v-model="formData.name"
          placeholder="请输入 VPC 名称"
          maxlength="64"
          show-word-limit
        />
      </el-form-item>

      <el-form-item label="IPv4 CIDR" prop="cidr">
        <el-input
          v-model="formData.cidr"
          placeholder="例如: 10.0.0.0/16"
        >
          <template #append>
            <el-button :icon="QuestionFilled" @click="showCidrHelp" />
          </template>
        </el-input>
      </el-form-item>

      <el-form-item label="IPv6 CIDR" prop="ipv6_cidr">
        <el-input
          v-model="formData.ipv6_cidr"
          placeholder="可选, 例如: 2001:db8::/32"
        />
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

  <!-- CIDR 帮助对话框 -->
  <el-dialog
    v-model="cidrHelpVisible"
    title="CIDR 说明"
    width="500px"
    append-to-body
  >
    <div class="cidr-help-content">
      <p><strong>CIDR (无类别域间路由)</strong> 是一种用于分配 IP 地址和路由的方法。</p>
      <el-divider />
      <h4>格式说明</h4>
      <p>CIDR 格式为: <code>IP地址/前缀长度</code></p>
      <p>例如: <code>10.0.0.0/16</code></p>
      <el-divider />
      <h4>常用 CIDR 示例</h4>
      <el-table :data="cidrExamples" border size="small">
        <el-table-column prop="cidr" label="CIDR" width="150" />
        <el-table-column prop="range" label="IP 范围" />
        <el-table-column prop="count" label="可用 IP 数" width="120" />
      </el-table>
      <el-divider />
      <h4>规划建议</h4>
      <ul>
        <li>VPC 建议使用较大的网段, 如 /16 或 /8</li>
        <li>子网使用较小的网段, 如 /24</li>
        <li>确保子网 CIDR 在 VPC CIDR 范围内</li>
      </ul>
    </div>
    <template #footer>
      <el-button type="primary" @click="cidrHelpVisible = false">知道了</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { QuestionFilled } from '@element-plus/icons-vue'
import CloudAccountSelector from '@/components/common/CloudAccountSelector.vue'
import { createVPC } from '@/api/network'
import { validateCIDR } from '@/utils/cidr'
import type { VPC } from '@/types'
import type { FormInstance, FormRules } from 'element-plus'

interface Props {
  visible: boolean
  accountId?: number
}

interface Emits {
  (e: 'update:visible', val: boolean): void
  (e: 'success', vpc: VPC): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// 表单引用
const formRef = ref<FormInstance>()

// 表单数据
const formData = ref({
  account_id: null as number | null,
  name: '',
  cidr: '',
  ipv6_cidr: '',
  description: ''
})

// 提交状态
const submitting = ref(false)

// CIDR 帮助对话框
const cidrHelpVisible = ref(false)

// CIDR 示例数据
const cidrExamples = [
  { cidr: '10.0.0.0/8', range: '10.0.0.0 - 10.255.255.255', count: '16,777,214' },
  { cidr: '10.0.0.0/16', range: '10.0.0.0 - 10.0.255.255', count: '65,534' },
  { cidr: '172.16.0.0/12', range: '172.16.0.0 - 172.31.255.255', count: '1,048,574' },
  { cidr: '192.168.0.0/16', range: '192.168.0.0 - 192.168.255.255', count: '65,534' },
  { cidr: '192.168.1.0/24', range: '192.168.1.0 - 192.168.1.255', count: '254' }
]

// 名称验证规则
const validateName = (_rule: any, value: string, callback: (error?: Error) => void) => {
  if (!value) {
    callback(new Error('请输入 VPC 名称'))
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

// IPv4 CIDR 验证规则
const validateCidrField = (_rule: any, value: string, callback: (error?: Error) => void) => {
  if (!value) {
    callback(new Error('请输入 IPv4 CIDR'))
    return
  }
  if (!validateCIDR(value)) {
    callback(new Error('CIDR 格式不正确, 请使用如 10.0.0.0/16 的格式'))
    return
  }
  callback()
}

// IPv6 CIDR 验证规则 (可选)
const validateIpv6Cidr = (_rule: any, value: string, callback: (error?: Error) => void) => {
  if (!value) {
    callback()
    return
  }
  // 简单的 IPv6 CIDR 验证
  const ipv6Pattern = /^([0-9a-fA-F]{1,4}:){2,7}[0-9a-fA-F]{1,4}\/\d{1,3}$/
  if (!ipv6Pattern.test(value)) {
    callback(new Error('IPv6 CIDR 格式不正确, 请使用如 2001:db8::/32 的格式'))
    return
  }
  callback()
}

// 描述验证规则
const validateDescription = (_rule: any, value: string, callback: (error?: Error) => void) => {
  if (value && value.length > 256) {
    callback(new Error('描述不能超过 256 个字符'))
    return
  }
  callback()
}

// 表单验证规则
const rules: FormRules = {
  account_id: [
    { required: true, message: '请选择云账号', trigger: 'change' }
  ],
  name: [
    { required: true, validator: validateName, trigger: 'blur' }
  ],
  cidr: [
    { required: true, validator: validateCidrField, trigger: 'blur' }
  ],
  ipv6_cidr: [
    { validator: validateIpv6Cidr, trigger: 'blur' }
  ],
  description: [
    { validator: validateDescription, trigger: 'blur' }
  ]
}

// 双向绑定 visible
const visible = computed({
  get: () => props.visible,
  set: (val: boolean) => emit('update:visible', val)
})

// 监听 accountId prop 变化
watch(() => props.accountId, (newVal) => {
  if (newVal) {
    formData.value.account_id = newVal
  }
}, { immediate: true })

// 监听 visible 变化, 打开时重置表单
watch(() => props.visible, (newVal) => {
  if (newVal) {
    resetForm()
    if (props.accountId) {
      formData.value.account_id = props.accountId
    }
  }
})

// 重置表单
const resetForm = () => {
  formData.value = {
    account_id: props.accountId || null,
    name: '',
    cidr: '',
    ipv6_cidr: '',
    description: ''
  }
  formRef.value?.clearValidate()
}

// 云账号变更
const handleAccountChange = (accountId: number | null) => {
  formData.value.account_id = accountId
}

// 显示 CIDR 帮助
const showCidrHelp = () => {
  cidrHelpVisible.value = true
}

// 关闭对话框
const handleClose = () => {
  visible.value = false
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
  } catch {
    return
  }

  submitting.value = true
  try {
    const data: {
      account_id: number
      name: string
      cidr: string
      ipv6_cidr?: string
      description?: string
    } = {
      account_id: formData.value.account_id!,
      name: formData.value.name,
      cidr: formData.value.cidr,
      ipv6_cidr: formData.value.ipv6_cidr || undefined,
      description: formData.value.description || undefined
    }

    const vpc = await createVPC(data)
    ElMessage.success('VPC 创建成功')
    emit('success', vpc)
    handleClose()
  } catch (error: any) {
    console.error('Failed to create VPC:', error)
    ElMessage.error(error.response?.data?.message || 'VPC 创建失败')
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

.cidr-help-content {
  line-height: 1.8;
}

.cidr-help-content h4 {
  margin-top: 10px;
  margin-bottom: 10px;
  color: #303133;
}

.cidr-help-content p {
  color: #606266;
}

.cidr-help-content code {
  padding: 2px 6px;
  background-color: #f5f7fa;
  border-radius: 4px;
  font-family: monospace;
}

.cidr-help-content ul {
  padding-left: 20px;
  color: #606266;
}

.cidr-help-content li {
  margin: 5px 0;
}

.cidr-help-content .el-table {
  margin: 10px 0;
}
</style>