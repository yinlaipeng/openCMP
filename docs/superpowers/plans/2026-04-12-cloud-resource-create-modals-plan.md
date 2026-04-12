# 云资源创建弹窗实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 为 openCMP 多云管理平台添加云资源创建弹窗，包括创建 VM 分步向导、创建 VPC/Subnet 弹窗。

**Architecture:** 使用 Vue 3 Composition API + Element Plus 组件。CreateVMModal 采用分步向导模式（el-steps），CreateVPCModal/CreateSubnetModal 采用单页表单。所有组件支持 v-model 双向绑定 visible 属性。

**Tech Stack:** Vue 3、TypeScript、Element Plus、Vite

---

## 文件结构

| 文件路径 | 负责内容 |
|---------|---------|
| `frontend/src/utils/cidr.ts` | CIDR 格式校验工具函数 |
| `frontend/src/components/common/CloudAccountSelector.vue` | 云账号选择器（复用组件） |
| `frontend/src/components/network/CreateVPCModal.vue` | 创建 VPC 弹窗 |
| `frontend/src/components/network/CreateSubnetModal.vue` | 创建 Subnet 弹窗 |
| `frontend/src/components/vm/CreateVMModal.vue` | 创建 VM 分步向导弹窗 |
| `frontend/src/views/compute/vms/index.vue` | 引入 CreateVMModal（修改） |
| `frontend/src/views/network/vpcs/index.vue` | 添加创建按钮和弹窗（修改） |
| `frontend/src/views/network/subnets/index.vue` | 添加创建按钮和弹窗（修改） |

---

## Task 1: CIDR 校验工具函数

**Files:**
- Create: `frontend/src/utils/cidr.ts`

- [ ] **Step 1: 创建 CIDR 校验工具文件**

```typescript
// frontend/src/utils/cidr.ts

/**
 * 校验 CIDR 格式是否有效
 * @param cidr CIDR 字符串，如 "10.0.0.0/16"
 * @returns 是否有效
 */
export function validateCIDR(cidr: string): boolean {
  if (!cidr) return false
  
  const pattern = /^(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})\/(\d{1,2})$/
  const match = cidr.match(pattern)
  
  if (!match) return false
  
  // 校验每个 IP 段在 0-255 范围内
  for (let i = 1; i <= 4; i++) {
    const segment = parseInt(match[i], 10)
    if (segment < 0 || segment > 255) return false
  }
  
  // 校验前缀长度在 0-32 范围内
  const prefix = parseInt(match[5], 10)
  if (prefix < 0 || prefix > 32) return false
  
  return true
}

/**
 * 将 CIDR 转换为 IP 起始范围和掩码
 * @param cidr CIDR 字符串
 * @returns { start: number, mask: number } 或 null
 */
export function parseCIDR(cidr: string): { start: number; mask: number } | null {
  if (!validateCIDR(cidr)) return null
  
  const [ip, prefixStr] = cidr.split('/')
  const prefix = parseInt(prefixStr, 10)
  
  const segments = ip.split('.').map(s => parseInt(s, 10))
  const ipNum = (segments[0] << 24) + (segments[1] << 16) + (segments[2] << 8) + segments[3]
  const mask = ~((1 << (32 - prefix)) - 1)
  
  return { start: ipNum & mask, mask }
}

/**
 * 校验子网 CIDR 是否属于 VPC CIDR 范围
 * @param subnetCIDR 子网 CIDR
 * @param vpcCIDR VPC CIDR
 * @returns 是否属于 VPC 范围
 */
export function isSubnetInVPC(subnetCIDR: string, vpcCIDR: string): boolean {
  const subnet = parseCIDR(subnetCIDR)
  const vpc = parseCIDR(vpcCIDR)
  
  if (!subnet || !vpc) return false
  
  // 子网掩码必须比 VPC 更精确（前缀更长）
  const subnetPrefix = parseInt(subnetCIDR.split('/')[1], 10)
  const vpcPrefix = parseInt(vpcCIDR.split('/')[1], 10)
  
  if (subnetPrefix <= vpcPrefix) return false
  
  // 子网起始 IP 必须在 VPC 范围内
  return (subnet.start & vpc.mask) === vpc.start
}

/**
 * 获取 CIDR 的默认掩码显示
 * @param cidr CIDR 字符串
 * @returns 格式化后的字符串
 */
export function formatCIDR(cidr: string): string {
  if (!validateCIDR(cidr)) return cidr
  return cidr
}
```

- [ ] **Step 2: 提交 CIDR 工具函数**

```bash
git add frontend/src/utils/cidr.ts
git commit -m "feat: add CIDR validation utility functions"
```

---

## Task 2: CloudAccountSelector 复用组件

**Files:**
- Create: `frontend/src/components/common/CloudAccountSelector.vue`

- [ ] **Step 1: 创建 CloudAccountSelector 组件**

```vue
<!-- frontend/src/components/common/CloudAccountSelector.vue -->
<template>
  <el-select
    v-model="selectedAccountId"
    :loading="loading"
    :disabled="disabled"
    :placeholder="placeholder"
    filterable
    clearable
    @change="handleChange"
  >
    <el-option
      v-for="account in accounts"
      :key="account.id"
      :label="account.name"
      :value="account.id"
    >
      <div class="account-option">
        <span class="account-name">{{ account.name }}</span>
        <el-tag
          :type="getAccountTagType(account.provider_type)"
          size="small"
          class="account-provider"
        >
          {{ account.provider_type }}
        </el-tag>
        <el-tag
          v-if="account.health_status === 'healthy'"
          type="success"
          size="small"
        >
          正常
        </el-tag>
        <el-tag
          v-else-if="account.health_status === 'unhealthy'"
          type="danger"
          size="small"
        >
          异常
        </el-tag>
      </div>
    </el-option>
  </el-select>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getCloudAccounts } from '@/api/cloud-account'
import type { CloudAccount } from '@/types'

interface Props {
  value?: number | null
  disabled?: boolean
  placeholder?: string
}

interface Emits {
  (e: 'change', accountId: number | null, account: CloudAccount | null): void
  (e: 'update:value', accountId: number | null): void
}

const props = withDefaults(defineProps<Props>(), {
  value: null,
  disabled: false,
  placeholder: '请选择云账号'
})

const emit = defineEmits<Emits>()

const accounts = ref<CloudAccount[]>([])
const selectedAccountId = ref<number | null>(props.value)
const loading = ref(false)

const getAccountTagType = (provider: string) => {
  const types: Record<string, string> = {
    alibaba: 'primary',
    tencent: 'warning',
    aws: 'success',
    azure: 'info'
  }
  return types[provider] || ''
}

const loadAccounts = async () => {
  loading.value = true
  try {
    const res = await getCloudAccounts()
    accounts.value = res.items || []
  } catch (e) {
    console.error(e)
    ElMessage.error('加载云账号列表失败')
  } finally {
    loading.value = false
  }
}

const handleChange = (val: number | null) => {
  emit('update:value', val)
  const account = accounts.value.find(a => a.id === val) || null
  emit('change', val, account)
}

watch(() => props.value, (val) => {
  selectedAccountId.value = val
})

onMounted(() => {
  loadAccounts()
})
</script>

<style scoped>
.account-option {
  display: flex;
  align-items: center;
  gap: 8px;
}

.account-name {
  flex: 1;
}

.account-provider {
  margin-left: 4px;
}
</style>
```

- [ ] **Step 2: 提交 CloudAccountSelector 组件**

```bash
git add frontend/src/components/common/CloudAccountSelector.vue
git commit -m "feat: add CloudAccountSelector reusable component"
```

---

## Task 3: CreateVPCModal 创建 VPC 弹窗

**Files:**
- Create: `frontend/src/components/network/CreateVPCModal.vue`

- [ ] **Step 1: 创建 CreateVPCModal 组件**

```vue
<!-- frontend/src/components/network/CreateVPCModal.vue -->
<template>
  <el-dialog
    v-model="dialogVisible"
    title="创建 VPC"
    width="500px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="formData"
      :rules="formRules"
      label-width="100px"
    >
      <el-form-item label="云账号" prop="account_id">
        <CloudAccountSelector
          v-model:value="formData.account_id"
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
          placeholder="如 10.0.0.0/16"
        >
          <template #append>
            <el-button @click="showCIDRHelp">帮助</el-button>
          </template>
        </el-input>
      </el-form-item>

      <el-form-item label="IPv6 CIDR">
        <el-input
          v-model="formData.ipv6_cidr"
          placeholder="可选，如 2001:db8::/32"
        />
      </el-form-item>

      <el-form-item label="描述">
        <el-input
          v-model="formData.description"
          type="textarea"
          :rows="3"
          placeholder="可选，最多 256 字符"
          maxlength="256"
          show-word-limit
        />
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" :loading="submitting" @click="handleSubmit">
        创建
      </el-button>
    </template>
  </el-dialog>

  <!-- CIDR 帮助弹窗 -->
  <el-dialog
    v-model="cidrHelpVisible"
    title="CIDR 格式说明"
    width="400px"
  >
    <div class="cidr-help">
      <p>CIDR (Classless Inter-Domain Routing) 格式用于表示 IP 地址范围：</p>
      <ul>
        <li><strong>格式</strong>: IP地址/前缀长度</li>
        <li><strong>示例</strong>: 10.0.0.0/16 表示 10.0.0.0 - 10.0.255.255</li>
        <li><strong>常用前缀</strong>:</li>
        <li>/8 - 约 1677 万地址（A类网络）</li>
        <li>/16 - 约 65536 地址（B类网络）</li>
        <li>/24 - 约 256 地址（C类网络）</li>
      </ul>
    </div>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, watch, computed } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import CloudAccountSelector from '@/components/common/CloudAccountSelector.vue'
import { createVPC } from '@/api/network'
import { validateCIDR } from '@/utils/cidr'
import type { VPC } from '@/types'

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

const dialogVisible = computed({
  get: () => props.visible,
  set: (val) => emit('update:visible', val)
})

const formRef = ref<FormInstance>()
const submitting = ref(false)
const cidrHelpVisible = ref(false)

const formData = reactive({
  account_id: null as number | null,
  name: '',
  cidr: '',
  ipv6_cidr: '',
  description: ''
})

const formRules: FormRules = {
  account_id: [
    { required: true, message: '请选择云账号', trigger: 'change' }
  ],
  name: [
    { required: true, message: '请输入名称', trigger: 'blur' },
    { min: 1, max: 64, message: '名称长度 1-64 字符', trigger: 'blur' },
    { pattern: /^[a-zA-Z0-9_\-\u4e00-\u9fa5]+$/, message: '名称只能包含字母、数字、下划线、横线和中文', trigger: 'blur' }
  ],
  cidr: [
    { required: true, message: '请输入 IPv4 CIDR', trigger: 'blur' },
    { validator: (rule, value, callback) => {
      if (!validateCIDR(value)) {
        callback(new Error('CIDR 格式无效，正确格式如 10.0.0.0/16'))
      } else {
        callback()
      }
    }, trigger: 'blur' }
  ]
}

const handleAccountChange = () => {
  // 云账号变更时清空其他字段
  formData.name = ''
  formData.cidr = ''
  formData.ipv6_cidr = ''
  formData.description = ''
}

const showCIDRHelp = () => {
  cidrHelpVisible.value = true
}

const handleClose = () => {
  formRef.value?.resetFields()
  emit('update:visible', false)
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
  } catch {
    return
  }
  
  if (!formData.account_id) {
    ElMessage.warning('请选择云账号')
    return
  }
  
  submitting.value = true
  try {
    const vpc = await createVPC({
      account_id: formData.account_id,
      name: formData.name,
      cidr: formData.cidr,
      description: formData.description
    })
    ElMessage.success('VPC 创建成功')
    emit('success', vpc)
    handleClose()
  } catch (e: any) {
    console.error(e)
    ElMessage.error(`创建失败: ${e.message || '未知错误'}`)
  } finally {
    submitting.value = false
  }
}

watch(() => props.accountId, (val) => {
  if (val) {
    formData.account_id = val
  }
}, { immediate: true })
</script>

<style scoped>
.cidr-help {
  line-height: 1.8;
}

.cidr-help ul {
  padding-left: 20px;
  margin-top: 10px;
}

.cidr-help li {
  margin: 4px 0;
}
</style>
```

- [ ] **Step 2: 提交 CreateVPCModal 组件**

```bash
git add frontend/src/components/network/CreateVPCModal.vue
git commit -m "feat: add CreateVPCModal component"
```

---

## Task 4: CreateSubnetModal 创建子网弹窗

**Files:**
- Create: `frontend/src/components/network/CreateSubnetModal.vue`

- [ ] **Step 1: 创建 CreateSubnetModal 组件**

```vue
<!-- frontend/src/components/network/CreateSubnetModal.vue -->
<template>
  <el-dialog
    v-model="dialogVisible"
    title="创建子网"
    width="500px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="formData"
      :rules="formRules"
      label-width="100px"
    >
      <el-form-item label="云账号" prop="account_id">
        <CloudAccountSelector
          v-model:value="formData.account_id"
          @change="handleAccountChange"
        />
      </el-form-item>

      <el-form-item label="VPC" prop="vpc_id">
        <el-select
          v-model="formData.vpc_id"
          :loading="loadingVPCs"
          :disabled="!formData.account_id"
          placeholder="请选择 VPC"
          filterable
          @change="handleVPCChange"
        >
          <el-option
            v-for="vpc in vpcs"
            :key="vpc.id"
            :label="vpc.name"
            :value="vpc.id"
          >
            <div class="vpc-option">
              <span>{{ vpc.name }}</span>
              <span class="vpc-cidr">{{ vpc.cidr }}</span>
            </div>
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
          placeholder="如 10.0.1.0/24，需在 VPC 范围内"
        />
        <div v-if="selectedVPC" class="cidr-hint">
          VPC CIDR: {{ selectedVPC.cidr }}，子网需在此范围内
        </div>
      </el-form-item>

      <el-form-item label="可用区" prop="zone_id">
        <el-select
          v-model="formData.zone_id"
          :loading="loadingZones"
          :disabled="!formData.account_id"
          placeholder="请选择可用区"
          filterable
        >
          <el-option
            v-for="zone in zones"
            :key="zone.id"
            :label="zone.name"
            :value="zone.id"
          />
        </el-select>
      </el-form-item>

      <el-form-item label="描述">
        <el-input
          v-model="formData.description"
          type="textarea"
          :rows="3"
          placeholder="可选"
        />
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" :loading="submitting" @click="handleSubmit">
        创建
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, watch, computed } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import CloudAccountSelector from '@/components/common/CloudAccountSelector.vue'
import { getVPCs, createSubnet, getZones } from '@/api/network'
import { validateCIDR, isSubnetInVPC } from '@/utils/cidr'
import type { VPC, Subnet, Zone } from '@/types'

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

const dialogVisible = computed({
  get: () => props.visible,
  set: (val) => emit('update:visible', val)
})

const formRef = ref<FormInstance>()
const submitting = ref(false)
const loadingVPCs = ref(false)
const loadingZones = ref(false)

const vpcs = ref<VPC[]>([])
const zones = ref<Zone[]>([])
const selectedVPC = ref<VPC | null>(null)

const formData = reactive({
  account_id: null as number | null,
  vpc_id: '',
  name: '',
  cidr: '',
  zone_id: '',
  description: ''
})

const formRules: FormRules = {
  account_id: [
    { required: true, message: '请选择云账号', trigger: 'change' }
  ],
  vpc_id: [
    { required: true, message: '请选择 VPC', trigger: 'change' }
  ],
  name: [
    { required: true, message: '请输入名称', trigger: 'blur' },
    { min: 1, max: 64, message: '名称长度 1-64 字符', trigger: 'blur' }
  ],
  cidr: [
    { required: true, message: '请输入 CIDR', trigger: 'blur' },
    { validator: (rule, value, callback) => {
      if (!validateCIDR(value)) {
        callback(new Error('CIDR 格式无效'))
        return
      }
      if (selectedVPC.value && !isSubnetInVPC(value, selectedVPC.value.cidr)) {
        callback(new Error(`子网 CIDR 必须在 VPC CIDR (${selectedVPC.value.cidr}) 范围内`))
        return
      }
      callback()
    }, trigger: 'blur' }
  ],
  zone_id: [
    { required: true, message: '请选择可用区', trigger: 'change' }
  ]
}

const loadVPCs = async () => {
  if (!formData.account_id) return
  
  loadingVPCs.value = true
  try {
    const res = await getVPCs({ account_id: formData.account_id })
    vpcs.value = Array.isArray(res) ? res : []
  } catch (e) {
    console.error(e)
    ElMessage.error('加载 VPC 列表失败')
  } finally {
    loadingVPCs.value = false
  }
}

const loadZones = async () => {
  if (!formData.account_id) return
  
  loadingZones.value = true
  try {
    const res = await getZones({ account_id: formData.account_id })
    zones.value = Array.isArray(res) ? res : []
  } catch (e) {
    console.error(e)
    ElMessage.error('加载可用区列表失败')
  } finally {
    loadingZones.value = false
  }
}

const handleAccountChange = () => {
  formData.vpc_id = ''
  formData.zone_id = ''
  vpcs.value = []
  zones.value = []
  selectedVPC.value = null
  
  if (formData.account_id) {
    loadVPCs()
    loadZones()
  }
}

const handleVPCChange = (val: string) => {
  selectedVPC.value = vpcs.value.find(v => v.id === val) || null
  formData.cidr = ''
}

const handleClose = () => {
  formRef.value?.resetFields()
  vpcs.value = []
  zones.value = []
  selectedVPC.value = null
  emit('update:visible', false)
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
  } catch {
    return
  }
  
  if (!formData.account_id) {
    ElMessage.warning('请选择云账号')
    return
  }
  
  submitting.value = true
  try {
    const subnet = await createSubnet({
      account_id: formData.account_id,
      name: formData.name,
      vpc_id: formData.vpc_id,
      cidr: formData.cidr,
      zone_id: formData.zone_id,
      description: formData.description
    })
    ElMessage.success('子网创建成功')
    emit('success', subnet)
    handleClose()
  } catch (e: any) {
    console.error(e)
    ElMessage.error(`创建失败: ${e.message || '未知错误'}`)
  } finally {
    submitting.value = false
  }
}

watch(() => props.accountId, (val) => {
  if (val) {
    formData.account_id = val
    loadVPCs()
    loadZones()
  }
}, { immediate: true })

watch(() => props.vpcId, (val) => {
  if (val) {
    formData.vpc_id = val
    selectedVPC.value = vpcs.value.find(v => v.id === val) || null
  }
}, { immediate: true })
</script>

<style scoped>
.vpc-option {
  display: flex;
  justify-content: space-between;
  width: 100%;
}

.vpc-cidr {
  color: #909399;
  font-size: 12px;
}

.cidr-hint {
  color: #909399;
  font-size: 12px;
  margin-top: 4px;
}
</style>
```

- [ ] **Step 2: 提交 CreateSubnetModal 组件**

```bash
git add frontend/src/components/network/CreateSubnetModal.vue
git commit -m "feat: add CreateSubnetModal component"
```

---

## Task 5: CreateVMModal 创建虚拟机弹窗（核心）

**Files:**
- Create: `frontend/src/components/vm/CreateVMModal.vue`

这是最复杂的组件，分步实现。

- [ ] **Step 1: 创建 CreateVMModal 组件框架和状态管理**

```vue
<!-- frontend/src/components/vm/CreateVMModal.vue -->
<template>
  <el-dialog
    v-model="dialogVisible"
    title="创建虚拟机"
    width="800px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <!-- 步骤导航 -->
    <el-steps :active="currentStep" finish-status="success" align-center>
      <el-step title="基本配置" />
      <el-step title="计算配置" />
      <el-step title="网络配置" />
      <el-step title="存储配置" />
      <el-step title="确认创建" />
    </el-steps>

    <!-- 步骤内容 -->
    <div class="step-content">
      <!-- Step 1: 基本配置 -->
      <div v-show="currentStep === 0" class="step-panel">
        <el-form :model="formData" label-width="120px">
          <el-form-item label="云账号">
            <CloudAccountSelector
              v-model:value="formData.account_id"
              @change="handleAccountChange"
            />
          </el-form-item>

          <el-form-item label="创建方式">
            <el-radio-group v-model="formData.createMode">
              <el-radio value="template">使用模板</el-radio>
              <el-radio value="custom">自定义配置</el-radio>
            </el-radio-group>
          </el-form-item>

          <el-form-item v-if="formData.createMode === 'template'" label="主机模板">
            <el-select
              v-model="formData.templateId"
              :loading="loadingTemplates"
              placeholder="请选择主机模板"
              filterable
              @change="handleTemplateChange"
            >
              <el-option
                v-for="tpl in templates"
                :key="tpl.id"
                :label="tpl.name"
                :value="tpl.id"
              />
            </el-select>
          </el-form-item>

          <el-form-item label="虚拟机名称">
            <el-input
              v-model="formData.name"
              placeholder="请输入名称"
              maxlength="64"
              show-word-limit
            />
          </el-form-item>

          <el-form-item label="创建数量">
            <el-input-number
              v-model="formData.count"
              :min="1"
              :max="100"
            />
          </el-form-item>
        </el-form>
      </div>

      <!-- Step 2: 计算配置 -->
      <div v-show="currentStep === 1" class="step-panel">
        <el-form :model="formData" label-width="120px">
          <el-form-item label="区域">
            <el-select
              v-model="formData.regionId"
              :loading="loadingRegions"
              :disabled="!formData.account_id"
              placeholder="请选择区域"
              filterable
              @change="handleRegionChange"
            >
              <el-option
                v-for="region in regions"
                :key="region.id"
                :label="region.name"
                :value="region.id"
              />
            </el-select>
          </el-form-item>

          <el-form-item label="可用区">
            <el-select
              v-model="formData.zoneId"
              :loading="loadingZones"
              :disabled="!formData.regionId"
              placeholder="请选择可用区"
              filterable
            >
              <el-option
                v-for="zone in zones"
                :key="zone.id"
                :label="zone.name"
                :value="zone.id"
              />
            </el-select>
          </el-form-item>

          <el-form-item label="镜像">
            <el-select
              v-model="formData.imageId"
              :loading="loadingImages"
              :disabled="!formData.account_id"
              placeholder="请选择镜像"
              filterable
            >
              <el-option
                v-for="image in images"
                :key="image.id"
                :label="image.name"
                :value="image.id"
              >
                <div class="image-option">
                  <span>{{ image.name }}</span>
                  <span class="image-os">{{ image.os_name }}</span>
                </div>
              </el-option>
            </el-select>
          </el-form-item>

          <el-form-item label="实例规格">
            <el-select
              v-model="formData.instanceType"
              :disabled="!formData.account_id"
              placeholder="请选择规格"
              filterable
            >
              <el-option
                v-for="spec in instanceTypes"
                :key="spec.value"
                :label="spec.label"
                :value="spec.value"
              />
            </el-select>
          </el-form-item>

          <el-form-item label="密钥对">
            <el-select
              v-model="formData.keypairId"
              :disabled="!formData.account_id"
              placeholder="可选，选择 SSH 密钥"
              clearable
              filterable
            >
              <el-option
                v-for="kp in keypairs"
                :key="kp.id"
                :label="kp.name"
                :value="kp.id"
              />
            </el-select>
          </el-form-item>
        </el-form>
      </div>

      <!-- Step 3: 网络配置 -->
      <div v-show="currentStep === 2" class="step-panel">
        <el-form :model="formData" label-width="120px">
          <el-form-item label="VPC">
            <el-select
              v-model="formData.vpcId"
              :loading="loadingVPCs"
              :disabled="!formData.account_id || !formData.regionId"
              placeholder="请选择 VPC"
              filterable
              @change="handleVPCChange"
            >
              <el-option
                v-for="vpc in vpcs"
                :key="vpc.id"
                :label="vpc.name"
                :value="vpc.id"
              />
            </el-select>
          </el-form-item>

          <el-form-item label="子网">
            <el-select
              v-model="formData.subnetId"
              :loading="loadingSubnets"
              :disabled="!formData.vpcId"
              placeholder="请选择子网"
              filterable
            >
              <el-option
                v-for="subnet in subnets"
                :key="subnet.id"
                :label="subnet.name"
                :value="subnet.id"
              >
                <div class="subnet-option">
                  <span>{{ subnet.name }}</span>
                  <span class="subnet-cidr">{{ subnet.cidr }}</span>
                </div>
              </el-option>
            </el-select>
          </el-form-item>

          <el-form-item label="安全组">
            <el-select
              v-model="formData.securityGroups"
              :loading="loadingSecurityGroups"
              :disabled="!formData.vpcId"
              placeholder="请选择安全组"
              multiple
              filterable
            >
              <el-option
                v-for="sg in securityGroups"
                :key="sg.id"
                :label="sg.name"
                :value="sg.id"
              />
            </el-select>
          </el-form-item>

          <el-form-item label="公网 IP">
            <el-switch v-model="formData.enablePublicIp" />
          </el-form-item>

          <el-form-item v-if="formData.enablePublicIp" label="带宽">
            <el-input-number
              v-model="formData.bandwidth"
              :min="1"
              :max="1000"
            />
            <span class="bandwidth-unit">Mbps</span>
          </el-form-item>
        </el-form>
      </div>

      <!-- Step 4: 存储配置 -->
      <div v-show="currentStep === 3" class="step-panel">
        <el-form :model="formData" label-width="120px">
          <el-form-item label="系统盘大小">
            <el-input-number
              v-model="formData.systemDiskSize"
              :min="20"
              :max="500"
            />
            <span class="disk-unit">GB</span>
          </el-form-item>

          <el-form-item label="系统盘类型">
            <el-select v-model="formData.systemDiskType" placeholder="请选择">
              <el-option label="高效云盘" value="cloud_efficiency" />
              <el-option label="SSD 云盘" value="cloud_ssd" />
              <el-option label="ESSD 云盘" value="cloud_essd" />
            </el-select>
          </el-form-item>

          <el-form-item label="数据盘">
            <div class="data-disks">
              <div v-for="(disk, idx) in formData.dataDisks" :key="idx" class="data-disk-item">
                <el-input-number
                  v-model="disk.size"
                  :min="20"
                  :max="2000"
                  placeholder="大小"
                />
                <span class="disk-unit">GB</span>
                <el-select v-model="disk.type" placeholder="类型" style="width: 120px">
                  <el-option label="高效云盘" value="cloud_efficiency" />
                  <el-option label="SSD 云盘" value="cloud_ssd" />
                  <el-option label="ESSD 云盘" value="cloud_essd" />
                </el-select>
                <el-button type="danger" link @click="removeDataDisk(idx)">
                  删除
                </el-button>
              </div>
              <el-button type="primary" link @click="addDataDisk">
                + 添加数据盘
              </el-button>
            </div>
          </el-form-item>
        </el-form>
      </div>

      <!-- Step 5: 确认创建 -->
      <div v-show="currentStep === 4" class="step-panel confirm-panel">
        <h4>配置汇总</h4>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="云账号">{{ getAccountName() }}</el-descriptions-item>
          <el-descriptions-item label="名称">{{ formData.name }}</el-descriptions-item>
          <el-descriptions-item label="数量">{{ formData.count }} 台</el-descriptions-item>
          <el-descriptions-item label="区域">{{ getRegionName() }}</el-descriptions-item>
          <el-descriptions-item label="可用区">{{ getZoneName() }}</el-descriptions-item>
          <el-descriptions-item label="镜像">{{ getImageName() }}</el-descriptions-item>
          <el-descriptions-item label="规格">{{ formData.instanceType }}</el-descriptions-item>
          <el-descriptions-item label="VPC">{{ getVPCName() }}</el-descriptions-item>
          <el-descriptions-item label="子网">{{ getSubnetName() }}</el-descriptions-item>
          <el-descriptions-item label="安全组">{{ getSecurityGroupNames() }}</el-descriptions-item>
          <el-descriptions-item label="公网 IP">{{ formData.enablePublicIp ? '是' : '否' }}</el-descriptions-item>
          <el-descriptions-item label="系统盘">{{ formData.systemDiskSize }}GB {{ formData.systemDiskType }}</el-descriptions-item>
        </el-descriptions>
        
        <div class="confirm-tip">
          <el-icon><InfoFilled /></el-icon>
          创建过程可能需要几分钟时间，请在虚拟机列表中查看创建状态
        </div>
      </div>
    </div>

    <!-- 底部按钮 -->
    <template #footer>
      <div class="dialog-footer">
        <el-button v-if="currentStep > 0" @click="prevStep">上一步</el-button>
        <el-button v-if="currentStep < 4" type="primary" :disabled="!canNextStep" @click="nextStep">
          下一步
        </el-button>
        <el-button v-if="currentStep === 4" type="primary" :loading="submitting" @click="handleSubmit">
          确认创建
        </el-button>
        <el-button @click="handleClose">取消</el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, watch, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { InfoFilled } from '@element-plus/icons-vue'
import CloudAccountSelector from '@/components/common/CloudAccountSelector.vue'
import {
  getRegions,
  getZones,
  getVPCs,
  getSubnets,
  getSecurityGroups,
  createVPC,
  createSubnet
} from '@/api/network'
import { getImages, createVM, getHostTemplates } from '@/api/compute'
import { getCloudAccounts } from '@/api/cloud-account'
import type { VPC, Subnet, SecurityGroup, Image, Region, Zone } from '@/types'
import type { HostTemplate } from '@/types'
import type { CloudAccount } from '@/types'

interface Props {
  visible: boolean
  accountId?: number
  templateId?: string
}

interface Emits {
  (e: 'update:visible', val: boolean): void
  (e: 'success', vm: any): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const dialogVisible = computed({
  get: () => props.visible,
  set: (val) => emit('update:visible', val)
})

const currentStep = ref(0)
const submitting = ref(false)

// 加载状态
const loadingTemplates = ref(false)
const loadingRegions = ref(false)
const loadingZones = ref(false)
const loadingImages = ref(false)
const loadingVPCs = ref(false)
const loadingSubnets = ref(false)
const loadingSecurityGroups = ref(false)

// 数据列表
const accounts = ref<CloudAccount[]>([])
const templates = ref<HostTemplate[]>([])
const regions = ref<Region[]>([])
const zones = ref<Zone[]>([])
const images = ref<Image[]>([])
const vpcs = ref<VPC[]>([])
const subnets = ref<Subnet[]>([])
const securityGroups = ref<SecurityGroup[]>([])
const keypairs = ref<any[]>([])

// 实例规格（通用列表）
const instanceTypes = ref([
  { value: 'ecs.t5-lc1m1.small', label: '1核1G' },
  { value: 'ecs.t5-lc1m2.small', label: '1核2G' },
  { value: 'ecs.c5.large', label: '2核4G' },
  { value: 'ecs.c5.xlarge', label: '4核8G' },
  { value: 'ecs.r5.large', label: '2核16G' },
  { value: 'ecs.g5.large', label: '2核8G GPU' }
])

// 表单数据
const formData = reactive({
  account_id: null as number | null,
  createMode: 'custom' as 'template' | 'custom',
  templateId: '',
  name: '',
  count: 1,
  regionId: '',
  zoneId: '',
  imageId: '',
  instanceType: '',
  keypairId: '',
  vpcId: '',
  subnetId: '',
  securityGroups: [] as string[],
  enablePublicIp: false,
  bandwidth: 10,
  systemDiskSize: 40,
  systemDiskType: 'cloud_ssd',
  dataDisks: [] as Array<{ size: number; type: string }>
})

// 当前步骤是否可以前进
const canNextStep = computed(() => {
  switch (currentStep.value) {
    case 0:
      return formData.account_id && formData.name
    case 1:
      return formData.regionId && formData.zoneId && formData.imageId && formData.instanceType
    case 2:
      return formData.vpcId && formData.subnetId && formData.securityGroups.length > 0
    case 3:
      return formData.systemDiskSize >= 20
    default:
      return true
  }
})

// 获取显示名称的方法
const getAccountName = () => {
  const acc = accounts.value.find(a => a.id === formData.account_id)
  return acc?.name || '-'
}

const getRegionName = () => {
  const r = regions.value.find(r => r.id === formData.regionId)
  return r?.name || formData.regionId || '-'
}

const getZoneName = () => {
  const z = zones.value.find(z => z.id === formData.zoneId)
  return z?.name || formData.zoneId || '-'
}

const getImageName = () => {
  const img = images.value.find(i => i.id === formData.imageId)
  return img?.name || formData.imageId || '-'
}

const getVPCName = () => {
  const v = vpcs.value.find(v => v.id === formData.vpcId)
  return v?.name || '-'
}

const getSubnetName = () => {
  const s = subnets.value.find(s => s.id === formData.subnetId)
  return s?.name || '-'
}

const getSecurityGroupNames = () => {
  return formData.securityGroups.map(id => {
    const sg = securityGroups.value.find(s => s.id === id)
    return sg?.name || id
  }).join(', ') || '-'
}

// 数据加载方法
const loadTemplates = async () => {
  if (!formData.account_id) return
  loadingTemplates.value = true
  try {
    const res = await getHostTemplates()
    templates.value = res.items || []
  } catch (e) {
    console.error(e)
  } finally {
    loadingTemplates.value = false
  }
}

const loadRegions = async () => {
  if (!formData.account_id) return
  loadingRegions.value = true
  try {
    const res = await getRegions({ account_id: formData.account_id })
    regions.value = Array.isArray(res) ? res : []
  } catch (e) {
    console.error(e)
    ElMessage.error('加载区域失败')
  } finally {
    loadingRegions.value = false
  }
}

const loadZones = async () => {
  if (!formData.account_id || !formData.regionId) return
  loadingZones.value = true
  try {
    const res = await getZones({
      account_id: formData.account_id,
      region_id: formData.regionId
    })
    zones.value = Array.isArray(res) ? res : []
  } catch (e) {
    console.error(e)
  } finally {
    loadingZones.value = false
  }
}

const loadImages = async () => {
  if (!formData.account_id) return
  loadingImages.value = true
  try {
    const res = await getImages({ account_id: formData.account_id })
    images.value = Array.isArray(res) ? res : []
  } catch (e) {
    console.error(e)
    ElMessage.error('加载镜像失败')
  } finally {
    loadingImages.value = false
  }
}

const loadVPCs = async () => {
  if (!formData.account_id) return
  loadingVPCs.value = true
  try {
    const res = await getVPCs({
      account_id: formData.account_id,
      region_id: formData.regionId
    })
    vpcs.value = Array.isArray(res) ? res : []
  } catch (e) {
    console.error(e)
  } finally {
    loadingVPCs.value = false
  }
}

const loadSubnets = async () => {
  if (!formData.account_id || !formData.vpcId) return
  loadingSubnets.value = true
  try {
    const res = await getSubnets({
      account_id: formData.account_id,
      vpc_id: formData.vpcId
    })
    subnets.value = Array.isArray(res) ? res : []
  } catch (e) {
    console.error(e)
  } finally {
    loadingSubnets.value = false
  }
}

const loadSecurityGroups = async () => {
  if (!formData.account_id || !formData.vpcId) return
  loadingSecurityGroups.value = true
  try {
    const res = await getSecurityGroups({
      account_id: formData.account_id,
      vpc_id: formData.vpcId
    })
    securityGroups.value = Array.isArray(res) ? res : []
  } catch (e) {
    console.error(e)
  } finally {
    loadingSecurityGroups.value = false
  }
}

const loadAccounts = async () => {
  try {
    const res = await getCloudAccounts()
    accounts.value = res.items || []
  } catch (e) {
    console.error(e)
  }
}

// 事件处理
const handleAccountChange = (id: number | null, acc: CloudAccount | null) => {
  // 清空依赖数据
  formData.regionId = ''
  formData.zoneId = ''
  formData.imageId = ''
  formData.vpcId = ''
  formData.subnetId = ''
  formData.securityGroups = []
  regions.value = []
  zones.value = []
  images.value = []
  vpcs.value = []
  subnets.value = []
  securityGroups.value = []
  
  if (id) {
    loadRegions()
    loadImages()
    if (formData.createMode === 'template') {
      loadTemplates()
    }
  }
}

const handleTemplateChange = (id: string) => {
  const tpl = templates.value.find(t => t.id === id)
  if (tpl) {
    // 自动填充模板配置
    formData.regionId = tpl.region_id
    formData.zoneId = tpl.zone_id
    formData.imageId = tpl.image_id
    formData.instanceType = tpl.instance_type
    formData.vpcId = tpl.vpc_id
    formData.subnetId = tpl.subnet_id
    formData.systemDiskSize = tpl.disk_size
  }
}

const handleRegionChange = () => {
  formData.zoneId = ''
  zones.value = []
  formData.vpcId = ''
  vpcs.value = []
  if (formData.regionId) {
    loadZones()
    loadVPCs()
  }
}

const handleVPCChange = () => {
  formData.subnetId = ''
  formData.securityGroups = []
  subnets.value = []
  securityGroups.value = []
  if (formData.vpcId) {
    loadSubnets()
    loadSecurityGroups()
  }
}

const addDataDisk = () => {
  formData.dataDisks.push({ size: 100, type: 'cloud_ssd' })
}

const removeDataDisk = (idx: number) => {
  formData.dataDisks.splice(idx, 1)
}

const prevStep = () => {
  if (currentStep.value > 0) {
    currentStep.value--
  }
}

const nextStep = () => {
  if (currentStep.value < 4 && canNextStep.value) {
    currentStep.value++
    // 进入步骤时加载必要数据
    if (currentStep.value === 2) {
      loadVPCs()
    }
  }
}

const handleClose = () => {
  currentStep.value = 0
  emit('update:visible', false)
}

const handleSubmit = async () => {
  if (!formData.account_id) {
    ElMessage.warning('请选择云账号')
    return
  }
  
  submitting.value = true
  try {
    const vm = await createVM({
      account_id: formData.account_id,
      name: formData.name,
      instance_type: formData.instanceType,
      image_id: formData.imageId,
      vpc_id: formData.vpcId,
      subnet_id: formData.subnetId,
      security_groups: formData.securityGroups,
      disk_size: formData.systemDiskSize
    })
    ElMessage.success('虚拟机创建成功')
    emit('success', vm)
    handleClose()
  } catch (e: any) {
    console.error(e)
    ElMessage.error(`创建失败: ${e.message || '未知错误'}`)
  } finally {
    submitting.value = false
  }
}

// 初始化
watch(() => props.visible, (val) => {
  if (val) {
    loadAccounts()
    if (props.accountId) {
      formData.account_id = props.accountId
      handleAccountChange(props.accountId, null)
    }
    if (props.templateId) {
      formData.createMode = 'template'
      formData.templateId = props.templateId
    }
  }
})
</script>

<style scoped>
.step-content {
  margin-top: 20px;
  min-height: 300px;
}

.step-panel {
  padding: 10px 0;
}

.confirm-panel h4 {
  margin-bottom: 16px;
}

.confirm-tip {
  margin-top: 16px;
  display: flex;
  align-items: center;
  gap: 8px;
  color: #909399;
  font-size: 13px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.image-option,
.subnet-option {
  display: flex;
  justify-content: space-between;
}

.image-os,
.subnet-cidr {
  color: #909399;
  font-size: 12px;
}

.disk-unit,
.bandwidth-unit {
  margin-left: 8px;
  color: #909399;
}

.data-disks {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.data-disk-item {
  display: flex;
  align-items: center;
  gap: 8px;
}
</style>
```

- [ ] **Step 2: 提交 CreateVMModal 组件**

```bash
git add frontend/src/components/vm/CreateVMModal.vue
git commit -m "feat: add CreateVMModal step wizard component"
```

---

## Task 6: 集成 CreateVMModal 到 VM 列表页面

**Files:**
- Modify: `frontend/src/views/compute/vms/index.vue`

- [ ] **Step 1: 修改 VM 列表页面引入创建弹窗**

需要修改的行：
- 导入 CreateVMModal
- 添加弹窗状态变量
- 修改 handleCreate 方法
- 添加弹窗组件到模板

```vue
<!-- 在 <script setup> 部分添加导入 -->
import CreateVMModal from '@/components/vm/CreateVMModal.vue'

<!-- 添加弹窗状态 -->
const createModalVisible = ref(false)

<!-- 修改 handleCreate 方法 -->
const handleCreate = () => {
  createModalVisible.value = true
}

<!-- 添加创建成功处理 -->
const handleCreateSuccess = (vm: VirtualMachine) => {
  ElMessage.success(`虚拟机 ${vm.name} 创建成功`)
  loadVMs()
}

<!-- 在模板最后添加弹窗组件 -->
<CreateVMModal
  v-model:visible="createModalVisible"
  :account-id="parseInt(queryForm.account_id) || undefined"
  @success="handleCreateSuccess"
/>
```

- [ ] **Step 2: 提交 VM 列表页面修改**

```bash
git add frontend/src/views/compute/vms/index.vue
git commit -m "feat: integrate CreateVMModal into VM list page"
```

---

## Task 7: 集成 CreateVPCModal 到 VPC 列表页面

**Files:**
- Modify: `frontend/src/views/network/vpcs/index.vue`

- [ ] **Step 1: 修改 VPC 列表页面添加创建功能**

```vue
<!-- 添加导入 -->
import CreateVPCModal from '@/components/network/CreateVPCModal.vue'
import { Plus } from '@element-plus/icons-vue'

<!-- 添加状态 -->
const createModalVisible = ref(false)

<!-- 添加创建方法 -->
const handleCreate = () => {
  createModalVisible.value = true
}

const handleCreateSuccess = (vpc: VPC) => {
  ElMessage.success(`VPC ${vpc.name} 创建成功`)
  loadVPCs()
}

<!-- 在 card-header 添加创建按钮 -->
<el-button type="primary" @click="handleCreate">
  <el-icon><Plus /></el-icon>
  创建 VPC
</el-button>

<!-- 添加弹窗组件 -->
<CreateVPCModal
  v-model:visible="createModalVisible"
  @success="handleCreateSuccess"
/>
```

- [ ] **Step 2: 提交 VPC 列表页面修改**

```bash
git add frontend/src/views/network/vpcs/index.vue
git commit -m "feat: add VPC create functionality to VPC list page"
```

---

## Task 8: 集成 CreateSubnetModal 到 Subnet 列表页面

**Files:**
- Modify: `frontend/src/views/network/subnets/index.vue`

- [ ] **Step 1: 修改 Subnet 列表页面添加创建功能**

```vue
<!-- 添加导入 -->
import CreateSubnetModal from '@/components/network/CreateSubnetModal.vue'
import { Plus } from '@element-plus/icons-vue'

<!-- 添加状态 -->
const createModalVisible = ref(false)

<!-- 添加创建方法 -->
const handleCreate = () => {
  createModalVisible.value = true
}

const handleCreateSuccess = (subnet: Subnet) => {
  ElMessage.success(`子网 ${subnet.name} 创建成功`)
  loadSubnets()
}

<!-- 在 card-header 添加创建按钮 -->
<el-button type="primary" @click="handleCreate">
  <el-icon><Plus /></el-icon>
  创建子网
</el-button>

<!-- 添加弹窗组件 -->
<CreateSubnetModal
  v-model:visible="createModalVisible"
  @success="handleCreateSuccess"
/>
```

- [ ] **Step 2: 提交 Subnet 列表页面修改**

```bash
git add frontend/src/views/network/subnets/index.vue
git commit -m "feat: add Subnet create functionality to Subnet list page"
```

---

## Task 9: 前端构建验证

- [ ] **Step 1: 验证前端构建**

```bash
cd frontend && npm run build 2>&1
```

Expected: 构建成功，无错误

- [ ] **Step 2: 提交所有改动并推送**

```bash
git status
git log --oneline -10
```

---

## Spec Coverage 自检

| Spec 要求 | 对应 Task |
|----------|-----------|
| CIDR 校验函数 | Task 1 |
| CloudAccountSelector | Task 2 |
| CreateVPCModal | Task 3 |
| CreateSubnetModal | Task 4 |
| CreateVMModal 分步向导 | Task 5 |
| VM 页面集成 | Task 6 |
| VPC 页面集成 | Task 7 |
| Subnet 页面集成 | Task 8 |
| Props/Emits 结构 | 各组件定义 |
| 字段联动逻辑 | Task 5 handleAccountChange/handleRegionChange |
| API 调用时机 | Task 5 各 load 方法 |
| 错误处理 | 各组件 try/catch + ElMessage |

---

**Plan complete. Ready for implementation.**