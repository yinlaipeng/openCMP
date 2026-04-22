<template>
  <el-dialog
    v-model="visible"
    title="创建虚拟机"
    width="800px"
    :close-on-click-modal="false"
    :destroy-on-close="true"
    @close="handleClose"
  >
    <!-- Steps wizard -->
    <el-steps :active="currentStep" finish-status="success" simple class="create-steps">
      <el-step title="基本配置" :icon="InfoFilled" />
      <el-step title="计算配置" :icon="Cpu" />
      <el-step title="网络配置" :icon="Connection" />
      <el-step title="存储配置" :icon="FolderOpened" />
      <el-step title="确认创建" :icon="Check" />
    </el-steps>

    <!-- Step content -->
    <div class="step-content" v-loading="loading">
      <!-- Step 1: Basic Configuration -->
      <div v-show="currentStep === 0" class="step-panel">
        <el-form
          ref="formRef1"
          :model="formData"
          :rules="rules1"
          label-width="120px"
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

          <el-form-item label="创建模式" prop="createMode">
            <el-radio-group v-model="formData.createMode">
              <el-radio value="template">使用模板</el-radio>
              <el-radio value="custom">自定义配置</el-radio>
            </el-radio-group>
          </el-form-item>

          <el-form-item
            v-if="formData.createMode === 'template'"
            label="主机模板"
            prop="templateId"
          >
            <el-select
              v-model="formData.templateId"
              placeholder="请选择主机模板"
              filterable
              :loading="templatesLoading"
              @change="handleTemplateChange"
            >
              <el-option
                v-for="template in hostTemplates"
                :key="template.id"
                :label="template.name"
                :value="template.id"
              >
                <div class="template-option">
                  <span>{{ template.name }}</span>
                  <el-tag size="small" type="info">{{ template.instance_type }}</el-tag>
                </div>
              </el-option>
            </el-select>
          </el-form-item>

          <el-form-item label="名称" prop="name">
            <el-input
              v-model="formData.name"
              placeholder="请输入虚拟机名称"
              maxlength="64"
              show-word-limit
            />
          </el-form-item>

          <el-form-item label="描述" prop="description">
            <el-input
              v-model="formData.description"
              type="textarea"
              placeholder="请输入描述信息"
              maxlength="256"
              show-word-limit
              :rows="3"
            />
          </el-form-item>

          <el-form-item label="计费类型" prop="billingType">
            <el-radio-group v-model="formData.billingType">
              <el-radio value="postpaid">按量付费</el-radio>
              <el-radio value="prepaid">包年包月</el-radio>
            </el-radio-group>
            <div class="form-tip">按量付费按小时计费，包年包月按月/年计费</div>
          </el-form-item>

          <el-form-item label="创建数量" prop="count">
            <el-input-number
              v-model="formData.count"
              :min="1"
              :max="100"
              :step="1"
            />
            <span class="form-tip">可创建 1-100 台虚拟机</span>
          </el-form-item>

          <el-form-item label="标签" prop="tags">
            <el-select
              v-model="formData.tags"
              placeholder="选择或输入标签"
              multiple
              filterable
              allow-create
              default-first-option
            >
              <el-option
                v-for="tag in availableTags"
                :key="tag"
                :label="tag"
                :value="tag"
              />
            </el-select>
            <div class="form-tip">可输入自定义标签</div>
          </el-form-item>
        </el-form>
      </div>

      <!-- Step 2: Compute Configuration -->
      <div v-show="currentStep === 1" class="step-panel">
        <el-form
          ref="formRef2"
          :model="formData"
          :rules="rules2"
          label-width="120px"
          label-position="right"
        >
          <el-form-item label="区域" prop="regionId">
            <el-select
              v-model="formData.regionId"
              placeholder="请选择区域"
              filterable
              :loading="regionsLoading"
              :disabled="formData.createMode === 'template'"
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

          <el-form-item label="可用区" prop="zoneId">
            <el-select
              v-model="formData.zoneId"
              placeholder="请选择可用区"
              filterable
              :loading="zonesLoading"
              :disabled="formData.createMode === 'template'"
              @change="handleZoneChange"
            >
              <el-option
                v-for="zone in zones"
                :key="zone.id"
                :label="zone.name"
                :value="zone.id"
              />
            </el-select>
          </el-form-item>

          <el-form-item label="镜像" prop="imageId">
            <el-select
              v-model="formData.imageId"
              placeholder="请选择镜像"
              filterable
              :loading="imagesLoading"
              :disabled="formData.createMode === 'template'"
            >
              <el-option
                v-for="image in images"
                :key="image.id"
                :label="`${image.name} (${image.os_name})`"
                :value="image.id"
              >
                <div class="image-option">
                  <span>{{ image.name }}</span>
                  <el-tag size="small">{{ image.os_name }}</el-tag>
                </div>
              </el-option>
            </el-select>
          </el-form-item>

          <el-form-item label="实例规格" prop="instanceType">
            <el-select
              v-model="formData.instanceType"
              placeholder="请选择实例规格"
              filterable
              :disabled="formData.createMode === 'template'"
            >
              <el-option
                v-for="type in instanceTypes"
                :key="type.value"
                :label="type.label"
                :value="type.value"
              >
                <div class="instance-type-option">
                  <span>{{ type.label }}</span>
                  <span class="type-spec">{{ type.spec }}</span>
                </div>
              </el-option>
            </el-select>
          </el-form-item>

          <el-form-item label="密钥对" prop="keypairId">
            <el-select
              v-model="formData.keypairId"
              placeholder="选择密钥对(可选)"
              filterable
              clearable
            >
              <el-option
                v-for="keypair in keypairs"
                :key="keypair.id"
                :label="keypair.name"
                :value="keypair.id"
              />
            </el-select>
            <div class="form-tip">密钥对用于 SSH 登录认证</div>
          </el-form-item>
        </el-form>
      </div>

      <!-- Step 3: Network Configuration -->
      <div v-show="currentStep === 2" class="step-panel">
        <el-form
          ref="formRef3"
          :model="formData"
          :rules="rules3"
          label-width="120px"
          label-position="right"
        >
          <el-form-item label="VPC" prop="vpcId">
            <el-select
              v-model="formData.vpcId"
              placeholder="请选择 VPC"
              filterable
              :loading="vpcsLoading"
              :disabled="formData.createMode === 'template'"
              @change="handleVpcChange"
            >
              <el-option
                v-for="vpc in vpcs"
                :key="vpc.id"
                :label="`${vpc.name} (${vpc.cidr})`"
                :value="vpc.id"
              />
            </el-select>
          </el-form-item>

          <el-form-item label="子网" prop="subnetId">
            <el-select
              v-model="formData.subnetId"
              placeholder="请选择子网"
              filterable
              :loading="subnetsLoading"
              :disabled="formData.createMode === 'template'"
            >
              <el-option
                v-for="subnet in subnets"
                :key="subnet.id"
                :label="`${subnet.name} (${subnet.cidr})`"
                :value="subnet.id"
              />
            </el-select>
          </el-form-item>

          <el-form-item label="安全组" prop="securityGroups">
            <el-select
              v-model="formData.securityGroups"
              placeholder="请选择安全组"
              multiple
              filterable
              :loading="sgLoading"
              :disabled="formData.createMode === 'template'"
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
            <el-switch
              v-model="formData.enablePublicIp"
              :disabled="formData.createMode === 'template'"
            />
          </el-form-item>

          <el-form-item
            v-if="formData.enablePublicIp"
            label="带宽"
            prop="bandwidth"
          >
            <el-input-number
              v-model="formData.bandwidth"
              :min="1"
              :max="100"
              :step="1"
            />
            <span class="form-unit">Mbps</span>
          </el-form-item>
        </el-form>
      </div>

      <!-- Step 4: Storage Configuration -->
      <div v-show="currentStep === 3" class="step-panel">
        <el-form
          ref="formRef4"
          :model="formData"
          :rules="rules4"
          label-width="120px"
          label-position="right"
        >
          <el-form-item label="系统盘大小" prop="systemDiskSize">
            <el-input-number
              v-model="formData.systemDiskSize"
              :min="20"
              :max="500"
              :step="10"
              :disabled="formData.createMode === 'template'"
            />
            <span class="form-unit">GB</span>
            <span class="form-tip">范围: 20-500 GB</span>
          </el-form-item>

          <el-form-item label="系统盘类型" prop="systemDiskType">
            <el-select
              v-model="formData.systemDiskType"
              placeholder="请选择系统盘类型"
              :disabled="formData.createMode === 'template'"
            >
              <el-option label="高效云盘" value="cloud_efficiency" />
              <el-option label="SSD 云盘" value="cloud_ssd" />
              <el-option label="ESSD 云盘" value="cloud_essd" />
              <el-option label="普通云盘" value="cloud" />
            </el-select>
          </el-form-item>

          <el-form-item label="数据盘">
            <div class="data-disks-container">
              <div
                v-for="(disk, index) in formData.dataDisks"
                :key="index"
                class="data-disk-item"
              >
                <el-input-number
                  v-model="disk.size"
                  :min="20"
                  :max="2000"
                  :step="10"
                  placeholder="大小"
                />
                <span class="form-unit">GB</span>
                <el-select
                  v-model="disk.type"
                  placeholder="类型"
                  style="width: 120px"
                >
                  <el-option label="高效云盘" value="cloud_efficiency" />
                  <el-option label="SSD 云盘" value="cloud_ssd" />
                  <el-option label="ESSD 云盘" value="cloud_essd" />
                </el-select>
                <el-button
                  type="danger"
                  :icon="Delete"
                  circle
                  size="small"
                  @click="removeDataDisk(index)"
                  :disabled="formData.createMode === 'template'"
                />
              </div>
              <el-button
                type="primary"
                :icon="Plus"
                size="small"
                @click="addDataDisk"
                :disabled="formData.createMode === 'template'"
              >
                添加数据盘
              </el-button>
            </div>
          </el-form-item>
        </el-form>
      </div>

      <!-- Step 5: Summary -->
      <div v-show="currentStep === 4" class="step-panel summary-panel">
        <el-alert
          type="info"
          :closable="false"
          show-icon
          class="summary-alert"
        >
          <template #title>
            请确认以下配置信息, 点击"创建"开始部署虚拟机
          </template>
        </el-alert>

        <el-descriptions :column="2" border class="summary-table">
          <el-descriptions-item label="云账号">
            {{ getAccountName() }}
          </el-descriptions-item>
          <el-descriptions-item label="创建模式">
            {{ formData.createMode === 'template' ? '使用模板' : '自定义配置' }}
          </el-descriptions-item>

          <el-descriptions-item v-if="formData.createMode === 'template'" label="模板">
            {{ getTemplateName() }}
          </el-descriptions-item>
          <el-descriptions-item label="名称">
            {{ formData.name }}
          </el-descriptions-item>
          <el-descriptions-item label="描述">
            {{ formData.description || '无' }}
          </el-descriptions-item>
          <el-descriptions-item label="计费类型">
            {{ formData.billingType === 'postpaid' ? '按量付费' : '包年包月' }}
          </el-descriptions-item>
          <el-descriptions-item label="创建数量">
            {{ formData.count }} 台
          </el-descriptions-item>
          <el-descriptions-item label="标签">
            {{ formData.tags.length ? formData.tags.join(', ') : '无' }}
          </el-descriptions-item>

          <el-descriptions-item label="区域">
            {{ getRegionName() }}
          </el-descriptions-item>
          <el-descriptions-item label="可用区">
            {{ getZoneName() }}
          </el-descriptions-item>
          <el-descriptions-item label="镜像">
            {{ getImageName() }}
          </el-descriptions-item>
          <el-descriptions-item label="实例规格">
            {{ formData.instanceType }}
          </el-descriptions-item>
          <el-descriptions-item label="密钥对">
            {{ formData.keypairId || '不使用' }}
          </el-descriptions-item>

          <el-descriptions-item label="VPC">
            {{ getVpcName() }}
          </el-descriptions-item>
          <el-descriptions-item label="子网">
            {{ getSubnetName() }}
          </el-descriptions-item>
          <el-descriptions-item label="安全组">
            {{ getSecurityGroupNames() }}
          </el-descriptions-item>
          <el-descriptions-item label="公网 IP">
            {{ formData.enablePublicIp ? `是 (${formData.bandwidth} Mbps)` : '否' }}
          </el-descriptions-item>

          <el-descriptions-item label="系统盘">
            {{ formData.systemDiskSize }} GB ({{ getDiskTypeName(formData.systemDiskType) }})
          </el-descriptions-item>
          <el-descriptions-item label="数据盘">
            {{ getDataDisksSummary() }}
          </el-descriptions-item>
        </el-descriptions>
      </div>
    </div>

    <!-- Footer with navigation buttons -->
    <template #footer>
      <div class="dialog-footer">
        <el-button
          v-if="currentStep > 0"
          @click="prevStep"
        >
          上一步
        </el-button>
        <el-button
          v-if="currentStep < 4"
          type="primary"
          :disabled="!canNextStep"
          @click="nextStep"
        >
          下一步
        </el-button>
        <el-button
          v-if="currentStep === 4"
          type="primary"
          :loading="submitting"
          @click="handleCreate"
        >
          创建
        </el-button>
        <el-button @click="handleClose">取消</el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import {
  InfoFilled,
  Cpu,
  Connection,
  FolderOpened,
  Check,
  Plus,
  Delete
} from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'

import CloudAccountSelector from '@/components/common/CloudAccountSelector.vue'
import {
  getRegions,
  getZones,
  getVPCs,
  getSubnets,
  getSecurityGroups
} from '@/api/network'
import { getImages, createVM, getHostTemplates } from '@/api/compute'
import { getCloudAccounts } from '@/api/cloud-account'
import type {
  VirtualMachine,
  Image,
  HostTemplate,
  Region,
  Zone,
  VPC,
  Subnet,
  SecurityGroup,
  CloudAccount
} from '@/types'

interface Props {
  visible: boolean
  accountId?: number
  templateId?: string
}

interface Emits {
  (e: 'update:visible', val: boolean): void
  (e: 'success', vm: VirtualMachine): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// Internal state for VM creation
interface VMCreateState {
  account_id: number | null
  createMode: 'template' | 'custom'
  templateId: string | null
  name: string
  description: string
  billingType: 'postpaid' | 'prepaid'
  count: number
  tags: string[]
  regionId: string
  zoneId: string
  imageId: string
  instanceType: string
  keypairId: string
  vpcId: string
  subnetId: string
  securityGroups: string[]
  enablePublicIp: boolean
  bandwidth: number
  systemDiskSize: number
  systemDiskType: string
  dataDisks: Array<{ size: number; type: string }>
}

// Form refs for each step
const formRef1 = ref<FormInstance>()
const formRef2 = ref<FormInstance>()
const formRef3 = ref<FormInstance>()
const formRef4 = ref<FormInstance>()

// Current step
const currentStep = ref(0)

// Form data
const formData = reactive<VMCreateState>({
  account_id: null,
  createMode: 'custom',
  templateId: null,
  name: '',
  description: '',
  billingType: 'postpaid',
  count: 1,
  tags: [],
  regionId: '',
  zoneId: '',
  imageId: '',
  instanceType: '',
  keypairId: '',
  vpcId: '',
  subnetId: '',
  securityGroups: [],
  enablePublicIp: false,
  bandwidth: 10,
  systemDiskSize: 40,
  systemDiskType: 'cloud_ssd',
  dataDisks: []
})

// Available tags for selection
const availableTags = ref<string[]>(['production', 'development', 'test', 'staging', 'important', 'temporary'])

// Loading states
const loading = ref(false)
const submitting = ref(false)
const regionsLoading = ref(false)
const zonesLoading = ref(false)
const imagesLoading = ref(false)
const vpcsLoading = ref(false)
const subnetsLoading = ref(false)
const sgLoading = ref(false)
const templatesLoading = ref(false)

// Data sources
const regions = ref<Region[]>([])
const zones = ref<Zone[]>([])
const images = ref<Image[]>([])
const vpcs = ref<VPC[]>([])
const subnets = ref<Subnet[]>([])
const securityGroups = ref<SecurityGroup[]>([])
const hostTemplates = ref<HostTemplate[]>([])
const cloudAccounts = ref<CloudAccount[]>([])
const keypairs = ref<any[]>([])

// Instance types preset list
const instanceTypes = [
  { value: 'ecs.t5-lc2m2.nano', label: 'ecs.t5-lc2m2.nano', spec: '1核 2GB' },
  { value: 'ecs.t5-lc1m2.small', label: 'ecs.t5-lc1m2.small', spec: '1核 2GB' },
  { value: 'ecs.t5-c1m2.large', label: 'ecs.t5-c1m2.large', spec: '2核 4GB' },
  { value: 'ecs.c6.large', label: 'ecs.c6.large', spec: '2核 4GB' },
  { value: 'ecs.c6.xlarge', label: 'ecs.c6.xlarge', spec: '4核 8GB' },
  { value: 'ecs.c6.2xlarge', label: 'ecs.c6.2xlarge', spec: '8核 16GB' },
  { value: 'ecs.c6.3xlarge', label: 'ecs.c6.3xlarge', spec: '12核 24GB' },
  { value: 'ecs.c6.4xlarge', label: 'ecs.c6.4xlarge', spec: '16核 32GB' },
  { value: 'ecs.g6.large', label: 'ecs.g6.large', spec: '2核 8GB' },
  { value: 'ecs.g6.xlarge', label: 'ecs.g6.xlarge', spec: '4核 16GB' },
  { value: 'ecs.g6.2xlarge', label: 'ecs.g6.2xlarge', spec: '8核 32GB' },
  { value: 'ecs.r6.large', label: 'ecs.r6.large', spec: '2核 16GB' },
  { value: 'ecs.r6.xlarge', label: 'ecs.r6.xlarge', spec: '4核 32GB' },
  { value: 'ecs.r6.2xlarge', label: 'ecs.r6.2xlarge', spec: '8核 64GB' },
  { value: 'ecs.sn1ne.large', label: 'ecs.sn1ne.large', spec: '2核 8GB' },
  { value: 'ecs.sn2ne.large', label: 'ecs.sn2ne.large', spec: '2核 8GB' }
]

// Form validation rules for Step 1
const rules1: FormRules = {
  account_id: [
    { required: true, message: '请选择云账号', trigger: 'change' }
  ],
  createMode: [
    { required: true, message: '请选择创建模式', trigger: 'change' }
  ],
  templateId: [
    {
      required: true,
      message: '请选择主机模板',
      trigger: 'change',
      validator: (_rule, value, callback) => {
        if (formData.createMode === 'template' && !value) {
          callback(new Error('请选择主机模板'))
        } else {
          callback()
        }
      }
    }
  ],
  name: [
    { required: true, message: '请输入虚拟机名称', trigger: 'blur' },
    { min: 2, max: 64, message: '名称长度为 2-64 个字符', trigger: 'blur' }
  ],
  count: [
    { required: true, message: '请输入创建数量', trigger: 'change' }
  ]
}

// Form validation rules for Step 2
const rules2: FormRules = {
  regionId: [
    { required: true, message: '请选择区域', trigger: 'change' }
  ],
  zoneId: [
    { required: true, message: '请选择可用区', trigger: 'change' }
  ],
  imageId: [
    { required: true, message: '请选择镜像', trigger: 'change' }
  ],
  instanceType: [
    { required: true, message: '请选择实例规格', trigger: 'change' }
  ]
}

// Form validation rules for Step 3
const rules3: FormRules = {
  vpcId: [
    { required: true, message: '请选择 VPC', trigger: 'change' }
  ],
  subnetId: [
    { required: true, message: '请选择子网', trigger: 'change' }
  ],
  bandwidth: [
    {
      validator: (_rule, _value, callback) => {
        if (formData.enablePublicIp && formData.bandwidth < 1) {
          callback(new Error('带宽至少为 1 Mbps'))
        } else {
          callback()
        }
      },
      trigger: 'change'
    }
  ]
}

// Form validation rules for Step 4
const rules4: FormRules = {
  systemDiskSize: [
    { required: true, message: '请输入系统盘大小', trigger: 'change' }
  ],
  systemDiskType: [
    { required: true, message: '请选择系统盘类型', trigger: 'change' }
  ]
}

// Computed visibility
const visible = computed({
  get: () => props.visible,
  set: (val: boolean) => emit('update:visible', val)
})

// Computed: can proceed to next step
const canNextStep = computed(() => {
  switch (currentStep.value) {
    case 0:
      if (!formData.account_id) return false
      if (formData.createMode === 'template' && !formData.templateId) return false
      if (!formData.name) return false
      return true
    case 1:
      if (!formData.regionId || !formData.zoneId) return false
      if (!formData.imageId || !formData.instanceType) return false
      return true
    case 2:
      if (!formData.vpcId || !formData.subnetId) return false
      return true
    case 3:
      if (!formData.systemDiskSize || !formData.systemDiskType) return false
      return true
    default:
      return true
  }
})

// Watch accountId prop
watch(() => props.accountId, (newVal) => {
  if (newVal) {
    formData.account_id = newVal
  }
}, { immediate: true })

// Watch templateId prop
watch(() => props.templateId, (newVal) => {
  if (newVal) {
    formData.templateId = newVal
    formData.createMode = 'template'
  }
}, { immediate: true })

// Watch visible prop - reset and load data when opened
watch(() => props.visible, async (newVal) => {
  if (newVal) {
    resetForm()
    currentStep.value = 0
    if (props.accountId) {
      formData.account_id = props.accountId
    }
    if (props.templateId) {
      formData.templateId = props.templateId
      formData.createMode = 'template'
    }
    await loadCloudAccounts()
    if (formData.account_id) {
      await loadRegions()
      await loadImages()
      await loadHostTemplates()
    }
  }
})

// Reset form data
const resetForm = () => {
  Object.assign(formData, {
    account_id: props.accountId || null,
    createMode: props.templateId ? 'template' : 'custom',
    templateId: props.templateId || null,
    name: '',
    description: '',
    billingType: 'postpaid',
    count: 1,
    tags: [],
    regionId: '',
    zoneId: '',
    imageId: '',
    instanceType: '',
    keypairId: '',
    vpcId: '',
    subnetId: '',
    securityGroups: [],
    enablePublicIp: false,
    bandwidth: 10,
    systemDiskSize: 40,
    systemDiskType: 'cloud_ssd',
    dataDisks: []
  })
  regions.value = []
  zones.value = []
  images.value = []
  vpcs.value = []
  subnets.value = []
  securityGroups.value = []
}

// Handle account change - cascading load
const handleAccountChange = async (accountId: number | null) => {
  formData.account_id = accountId
  // Clear dependent fields
  formData.regionId = ''
  formData.zoneId = ''
  formData.imageId = ''
  formData.vpcId = ''
  formData.subnetId = ''
  formData.securityGroups = []
  zones.value = []
  vpcs.value = []
  subnets.value = []
  securityGroups.value = []

  if (accountId) {
    loading.value = true
    try {
      await Promise.all([
        loadRegions(),
        loadImages(),
        loadHostTemplates()
      ])
    } finally {
      loading.value = false
    }
  }
}

// Handle region change - cascading load
const handleRegionChange = async () => {
  formData.zoneId = ''
  formData.vpcId = ''
  formData.subnetId = ''
  formData.securityGroups = []
  subnets.value = []
  securityGroups.value = []

  if (formData.regionId && formData.account_id) {
    zonesLoading.value = true
    try {
      await loadZones()
      await loadVPCs()
    } finally {
      zonesLoading.value = false
    }
  }
}

// Handle zone change
const handleZoneChange = () => {
  // Zone change doesn't need to load additional data
}

// Handle VPC change - cascading load
const handleVpcChange = async () => {
  formData.subnetId = ''
  formData.securityGroups = []

  if (formData.vpcId && formData.account_id) {
    subnetsLoading.value = true
    sgLoading.value = true
    try {
      await Promise.all([
        loadSubnets(),
        loadSecurityGroups()
      ])
    } finally {
      subnetsLoading.value = false
      sgLoading.value = false
    }
  }
}

// Handle template selection - auto-fill config
const handleTemplateChange = async (templateId: string) => {
  if (!templateId) return

  const template = hostTemplates.value.find(t => t.id === templateId)
  if (!template) return

  // Auto-fill from template
  formData.regionId = template.region_id || ''
  formData.zoneId = template.zone_id || ''
  formData.imageId = template.image_id || ''
  formData.instanceType = template.instance_type || ''
  formData.vpcId = template.vpc_id || ''
  formData.subnetId = template.subnet_id || ''
  formData.systemDiskSize = template.disk_size || 40

  // Load dependent data based on template values
  if (formData.account_id && formData.regionId) {
    await loadZones()
    await loadVPCs()
  }
  if (formData.vpcId) {
    await loadSubnets()
    await loadSecurityGroups()
  }
}

// Load cloud accounts
const loadCloudAccounts = async () => {
  try {
    const res = await getCloudAccounts()
    cloudAccounts.value = res.items || []
  } catch (error) {
    console.error('Failed to load cloud accounts:', error)
    cloudAccounts.value = []
  }
}

// Load regions
const loadRegions = async () => {
  if (!formData.account_id) return
  regionsLoading.value = true
  try {
    const res = await getRegions({ account_id: formData.account_id })
    regions.value = res || []
  } catch (error) {
    console.error('Failed to load regions:', error)
    regions.value = []
  } finally {
    regionsLoading.value = false
  }
}

// Load zones
const loadZones = async () => {
  if (!formData.account_id || !formData.regionId) return
  try {
    const res = await getZones({
      account_id: formData.account_id,
      region_id: formData.regionId
    })
    zones.value = res || []
  } catch (error) {
    console.error('Failed to load zones:', error)
    zones.value = []
  }
}

// Load images
const loadImages = async () => {
  if (!formData.account_id) return
  imagesLoading.value = true
  try {
    const res = await getImages({ account_id: formData.account_id })
    images.value = res || []
  } catch (error) {
    console.error('Failed to load images:', error)
    images.value = []
  } finally {
    imagesLoading.value = false
  }
}

// Load VPCs
const loadVPCs = async () => {
  if (!formData.account_id || !formData.regionId) return
  vpcsLoading.value = true
  try {
    const res = await getVPCs({
      account_id: formData.account_id,
      region_id: formData.regionId
    })
    vpcs.value = res || []
  } catch (error) {
    console.error('Failed to load VPCs:', error)
    vpcs.value = []
  } finally {
    vpcsLoading.value = false
  }
}

// Load subnets
const loadSubnets = async () => {
  if (!formData.account_id || !formData.vpcId) return
  try {
    const res = await getSubnets({
      account_id: formData.account_id,
      vpc_id: formData.vpcId
    })
    subnets.value = res || []
  } catch (error) {
    console.error('Failed to load subnets:', error)
    subnets.value = []
  }
}

// Load security groups
const loadSecurityGroups = async () => {
  if (!formData.account_id || !formData.vpcId) return
  try {
    const res = await getSecurityGroups({
      account_id: formData.account_id,
      vpc_id: formData.vpcId
    })
    securityGroups.value = res || []
  } catch (error) {
    console.error('Failed to load security groups:', error)
    securityGroups.value = []
  }
}

// Load host templates
const loadHostTemplates = async () => {
  if (!formData.account_id) return
  templatesLoading.value = true
  try {
    const res = await getHostTemplates()
    hostTemplates.value = res.items || []
  } catch (error) {
    console.error('Failed to load host templates:', error)
    hostTemplates.value = []
  } finally {
    templatesLoading.value = false
  }
}

// Add data disk
const addDataDisk = () => {
  formData.dataDisks.push({ size: 20, type: 'cloud_ssd' })
}

// Remove data disk
const removeDataDisk = (index: number) => {
  formData.dataDisks.splice(index, 1)
}

// Navigation methods
const prevStep = () => {
  if (currentStep.value > 0) {
    currentStep.value -= 1
  }
}

const nextStep = async () => {
  // Validate current step
  const formRef = [formRef1, formRef2, formRef3, formRef4][currentStep.value]
  if (formRef.value) {
    try {
      await formRef.value.validate()
    } catch {
      return
    }
  }

  if (currentStep.value < 4) {
    currentStep.value += 1
  }
}

// Close dialog
const handleClose = () => {
  visible.value = false
}

// Summary helper methods
const getAccountName = (): string => {
  const account = cloudAccounts.value.find(a => a.id === formData.account_id)
  return account?.name || '未知'
}

const getTemplateName = (): string => {
  const template = hostTemplates.value.find(t => t.id === formData.templateId)
  return template?.name || '未知'
}

const getRegionName = (): string => {
  const region = regions.value.find(r => r.id === formData.regionId)
  return region?.name || formData.regionId || '未知'
}

const getZoneName = (): string => {
  const zone = zones.value.find(z => z.id === formData.zoneId)
  return zone?.name || formData.zoneId || '未知'
}

const getImageName = (): string => {
  const image = images.value.find(i => i.id === formData.imageId)
  return image?.name || formData.imageId || '未知'
}

const getVpcName = (): string => {
  const vpc = vpcs.value.find(v => v.id === formData.vpcId)
  return vpc ? `${vpc.name} (${vpc.cidr})` : formData.vpcId || '未知'
}

const getSubnetName = (): string => {
  const subnet = subnets.value.find(s => s.id === formData.subnetId)
  return subnet ? `${subnet.name} (${subnet.cidr})` : formData.subnetId || '未知'
}

const getSecurityGroupNames = (): string => {
  if (!formData.securityGroups.length) return '无'
  const names = formData.securityGroups.map(id => {
    const sg = securityGroups.value.find(s => s.id === id)
    return sg?.name || id
  })
  return names.join(', ')
}

const getDiskTypeName = (type: string): string => {
  const typeMap: Record<string, string> = {
    cloud_efficiency: '高效云盘',
    cloud_ssd: 'SSD 云盘',
    cloud_essd: 'ESSD 云盘',
    cloud: '普通云盘'
  }
  return typeMap[type] || type
}

const getDataDisksSummary = (): string => {
  if (!formData.dataDisks.length) return '无'
  return formData.dataDisks
    .map(d => `${d.size}GB (${getDiskTypeName(d.type)})`)
    .join(', ')
}

// Handle create
const handleCreate = async () => {
  if (!formData.account_id) {
    ElMessage.error('请选择云账号')
    return
  }

  submitting.value = true
  try {
    // Create VM request
    const requestData: {
      account_id: number
      name: string
      instance_type: string
      image_id: string
      vpc_id: string
      subnet_id: string
      security_groups?: string[]
      disk_size?: number
    } = {
      account_id: formData.account_id,
      name: formData.name,
      instance_type: formData.instanceType,
      image_id: formData.imageId,
      vpc_id: formData.vpcId,
      subnet_id: formData.subnetId,
      security_groups: formData.securityGroups,
      disk_size: formData.systemDiskSize
    }

    const vm = await createVM(requestData)
    ElMessage.success(`虚拟机创建成功`)
    emit('success', vm)
    handleClose()
  } catch (error: any) {
    console.error('Failed to create VM:', error)
    ElMessage.error(error.response?.data?.message || '虚拟机创建失败')
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
/* Dialog */
.el-dialog {
  border-radius: var(--radius-xl);
}

/* Steps Wizard */
.create-steps {
  margin-bottom: var(--space-6);
  padding: var(--space-4);
  background: rgba(248, 250, 252, 0.5);
  backdrop-filter: blur(8px);
  border-radius: var(--radius-lg);
}

.create-steps .el-step__title {
  font-weight: var(--font-weight-medium);
}

.create-steps .el-step.is-success .el-step__icon {
  background: rgba(34, 197, 94, 0.15);
  border-color: var(--color-accent);
}

.create-steps .el-step.is-process .el-step__icon {
  background: rgba(34, 197, 94, 0.25);
  border-color: var(--color-accent);
  box-shadow: 0 0 12px rgba(34, 197, 94, 0.3);
}

/* Step Content */
.step-content {
  min-height: 320px;
  padding: var(--space-4) 0;
}

.step-panel {
  max-height: 420px;
  overflow-y: auto;
  padding: var(--space-2);
}

.summary-panel {
  max-height: none;
}

/* Summary Page */
.summary-alert {
  margin-bottom: var(--space-4);
  border-radius: var(--radius-md);
}

.summary-table {
  margin-top: var(--space-3);
  background: rgba(248, 250, 252, 0.4);
  border-radius: var(--radius-lg);
}

.summary-table .el-descriptions__label {
  font-weight: var(--font-weight-medium);
  background: rgba(248, 250, 252, 0.6);
}

.summary-table .el-descriptions__content {
  font-family: var(--font-mono);
}

/* Footer */
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: var(--space-2);
}

/* Select Options */
.template-option,
.image-option {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.instance-type-option {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.type-spec {
  color: var(--color-muted);
  font-size: var(--font-size-xs);
}

/* Form Helpers */
.form-tip {
  margin-left: var(--space-2);
  color: var(--color-muted);
  font-size: var(--font-size-xs);
}

.form-unit {
  margin-left: var(--space-2);
  color: var(--color-muted);
  font-size: var(--font-size-sm);
}

/* Data Disks */
.data-disks-container {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.data-disk-item {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2);
  background: rgba(248, 250, 252, 0.3);
  border-radius: var(--radius-md);
}

/* Responsive */
@media (max-width: 768px) {
  .el-dialog {
    width: 95% !important;
  }

  .create-steps {
    padding: var(--space-2);
  }

  .step-panel {
    max-height: 300px;
  }

  .dialog-footer {
    flex-wrap: wrap;
  }
}
</style>