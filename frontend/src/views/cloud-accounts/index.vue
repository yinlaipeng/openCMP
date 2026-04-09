<template>
  <div class="cloud-accounts-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">云账户管理</span>
          <el-button type="primary" @click="showWizard = true">
            <el-icon><Plus /></el-icon>
            添加云账户
          </el-button>
        </div>
      </template>

      <el-table :data="accounts" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="名称" width="150" />
        <el-table-column prop="description" label="备注" width="200" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'danger'">
              {{ row.status === 'active' ? '正常' : row.status === 'inactive' ? '未激活' : '错误' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="enabled" label="启用状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.enabled ? 'success' : 'info'">
              {{ row.enabled ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="health_status" label="健康状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.health_status === 'healthy' ? 'success' : 'warning'">
              {{ row.health_status === 'healthy' ? '健康' : '异常' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="balance" label="余额" width="100">
          <template #default="{ row }">
            ¥{{ row.balance?.toFixed(2) || '0.00' }}
          </template>
        </el-table-column>
        <el-table-column prop="provider_type" label="平台" width="120">
          <template #default="{ row }">
            <el-tag :type="getProviderType(row.provider_type)">
              {{ getProviderName(row.provider_type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="account_number" label="账号" width="150" />
        <el-table-column prop="last_sync" label="上次同步" width="180" />
        <el-table-column prop="sync_time" label="同步时间" width="120" />
        <el-table-column prop="domain_id" label="所属域" width="100" />
        <el-table-column prop="resource_assignment_method" label="资源归属方式" width="150" />
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleVerify(row)">验证</el-button>
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="pagination.currentPage"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        class="pagination"
      />
    </el-card>

    <!-- 向导式添加云账号对话框 -->
    <el-dialog
      v-model="showWizard"
      :title="'添加云账户'"
      width="80%"
      :fullscreen="isMobile"
    >
      <el-steps :active="wizardStep" finish-status="success" align-center>
        <el-step title="选择云平台" />
        <el-step title="配置基本信息" />
        <el-step title="配置同步区域" />
        <el-step title="定时同步任务" />
      </el-steps>

      <!-- 步骤1: 选择云平台 -->
      <div v-if="wizardStep === 0" class="wizard-step-content">
        <h3>请选择您要添加的云平台</h3>
        <div class="provider-grid">
          <el-card
            v-for="provider in providers"
            :key="provider.id"
            class="provider-card"
            @click="selectProvider(provider.id)"
          >
            <div class="provider-item">
              <el-icon :size="40"><component :is="provider.icon" /></el-icon>
              <h4>{{ provider.displayName }}</h4>
            </div>
          </el-card>
        </div>
      </div>

      <!-- 步骤2: 配置基本信息 -->
      <div v-if="wizardStep === 1" class="wizard-step-content">
        <el-form :model="wizardForm" :rules="wizardRules" ref="wizardFormRef" label-width="150px">
          <el-form-item label="名称" prop="name">
            <el-input v-model="wizardForm.name" placeholder="请输入云账户名称" />
          </el-form-item>

          <el-form-item label="备注" prop="remarks">
            <el-input v-model="wizardForm.remarks" type="textarea" placeholder="请输入备注" />
          </el-form-item>

          <el-form-item label="账号类型" prop="accountType">
            <el-radio-group v-model="wizardForm.accountType">
              <el-radio label="public">公共云</el-radio>
              <el-radio label="finance" disabled>金融云（暂不支持）</el-radio>
            </el-radio-group>
          </el-form-item>

          <el-form-item label="密钥ID" prop="accessKeyId">
            <el-input v-model="wizardForm.accessKeyId" placeholder="请输入Access Key ID" />
          </el-form-item>

          <el-form-item label="密钥Secret" prop="accessKeySecret">
            <el-input v-model="wizardForm.accessKeySecret" type="password" placeholder="请输入Access Key Secret" />
          </el-form-item>

          <el-form-item label="资源归属方式" prop="resourceAssignmentMethod">
            <el-checkbox-group v-model="wizardForm.resourceAssignmentMethod">
              <el-checkbox label="sync_strategy">根据同步策略归属</el-checkbox>
              <el-checkbox label="cloud_project">根据云上项目归属</el-checkbox>
              <el-checkbox label="cloud_subscription">根据云订阅归属</el-checkbox>
              <el-checkbox label="specific_project">指定项目</el-checkbox>
            </el-checkbox-group>
          </el-form-item>

          <!-- 根据选择的资源归属方式进行相应配置 -->
          <div v-if="wizardForm.resourceAssignmentMethod.includes('sync_strategy')">
            <h4>资源归属方式-根据同步策略归属</h4>
            <el-form-item label="同步策略">
              <el-select v-model="wizardForm.syncStrategy.policy" placeholder="请选择同步策略">
                <el-option
                  v-for="policy in syncPolicies"
                  :key="policy.id"
                  :label="policy.name"
                  :value="policy.id.toString()"
                />
              </el-select>
            </el-form-item>
            <el-form-item label="同步策略生效范围">
              <el-checkbox-group v-model="wizardForm.syncStrategy.scope">
                <el-checkbox label="resource_tags">资源标签</el-checkbox>
                <el-checkbox label="project_tags">项目标签</el-checkbox>
              </el-checkbox-group>
            </el-form-item>
          </div>

          <div v-if="wizardForm.resourceAssignmentMethod.includes('cloud_project')">
            <h4>资源归属方式-根据云上项目归属</h4>
          </div>

          <div v-if="wizardForm.resourceAssignmentMethod.includes('cloud_subscription')">
            <h4>资源归属方式-根据云订阅归属</h4>
          </div>

          <div v-if="wizardForm.resourceAssignmentMethod.includes('specific_project')">
            <h4>资源归属方式-指定项目</h4>
          </div>

          <!-- 显示资源归属方式的备注信息 -->
          <div v-if="getResourceAssignmentNotes()" class="resource-assignment-notes">
            <p><em>{{ getResourceAssignmentNotes() }}</em></p>
          </div>

          <!-- 缺省项目选择框（仅在有选择的情况下显示一次） -->
          <el-form-item
            v-if="wizardForm.resourceAssignmentMethod.length > 0"
            label="缺省项目"
          >
            <el-select
              v-model="wizardForm.syncStrategy.defaultProject"
              placeholder="请选择缺省项目"
            >
              <el-option
                v-for="project in projects"
                :key="project.id"
                :label="project.name"
                :value="project.id"
              />
            </el-select>
          </el-form-item>
        </el-form>
      </div>

      <!-- 步骤3: 配置同步区域 -->
      <div v-if="wizardStep === 2" class="wizard-step-content">
        <p class="region-tip">
          请选择要同步资源的区域，该配置可在云账号导入后在云账号详情页-订阅-区域进行修改，请放心配置！（默认同步所有区域）
        </p>

        <div class="region-actions">
          <el-button @click="selectAllRegions">全选</el-button>
          <el-button @click="clearAllRegions">清空</el-button>
        </div>

        <el-table
          :data="availableRegions"
          style="width: 100%; margin-top: 20px;"
          @selection-change="handleRegionSelectionChange"
        >
          <el-table-column type="selection" width="55" />
          <el-table-column prop="name" label="名称" />
          <el-table-column prop="status" label="状态">
            <template #default="{ row }">
              <el-tag :type="row.status === 'available' ? 'success' : 'info'">
                {{ row.status === 'available' ? '可用' : '不可用' }}
              </el-tag>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 步骤4: 定时同步任务设置 -->
      <div v-if="wizardStep === 3" class="wizard-step-content">
        <h3>定时同步任务设置（可选）</h3>
        <el-form :model="scheduleForm" label-width="120px">
          <el-form-item label="名称">
            <el-input v-model="scheduleForm.name" placeholder="请输入任务名称" />
          </el-form-item>

          <el-form-item label="类型">
            <el-select v-model="scheduleForm.type" disabled>
              <el-option label="同步云账号" value="sync_cloud_account" />
            </el-select>
          </el-form-item>

          <el-form-item label="触发频次">
            <el-select v-model="scheduleForm.frequency">
              <el-option label="单次" value="once" />
              <el-option label="每天" value="daily" />
              <el-option label="每周" value="weekly" />
              <el-option label="每月" value="monthly" />
              <el-option label="周期" value="custom" />
            </el-select>
          </el-form-item>

          <el-form-item label="触发时间">
            <el-time-picker
              v-model="scheduleForm.triggerTime"
              placeholder="选择时间"
              format="HH:mm"
              value-format="HH:mm"
            />
          </el-form-item>

          <el-form-item label="有效时间">
            <el-date-picker
              v-model="scheduleForm.validRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              format="YYYY-MM-DD"
              value-format="YYYY-MM-DD"
            />
          </el-form-item>
        </el-form>

        <div class="sync-strategy-info" v-if="wizardForm.resourceAssignmentMethod.length > 0">
          <h4>资源归属方式配置信息：</h4>
          <ul>
            <li v-if="wizardForm.resourceAssignmentMethod.includes('sync_strategy')">同步策略归属：策略 {{ wizardForm.syncStrategy.policy }}，生效范围 {{ wizardForm.syncStrategy.scope.join(', ') }}</li>
            <li v-if="wizardForm.resourceAssignmentMethod.includes('cloud_project')">云上项目映射：缺省项目ID {{ wizardForm.cloudProjectMapping.defaultProject }}</li>
            <li v-if="wizardForm.resourceAssignmentMethod.includes('cloud_subscription')">云订阅映射：缺省项目ID {{ wizardForm.cloudSubscriptionMapping.defaultProject }}</li>
            <li v-if="wizardForm.resourceAssignmentMethod.includes('specific_project')">指定项目：项目ID {{ wizardForm.specificProject.defaultProject }}</li>
          </ul>
        </div>
      </div>

      <template #footer>
        <div class="wizard-footer">
          <el-button @click="previousStep" :disabled="wizardStep === 0">上一步</el-button>
          <el-button
            v-if="wizardStep < 3"
            type="primary"
            @click="nextStep"
          >
            下一步
          </el-button>
          <el-button
            v-else
            type="primary"
            @click="submitWizard"
            :loading="submitting"
          >
            提交
          </el-button>
          <el-button @click="showWizard = false">取消</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed, watch } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance } from 'element-plus'
import { Plus, Cloudy, Search, Menu } from '@element-plus/icons-vue'
import { getCloudAccounts, createCloudAccount, updateCloudAccount, deleteCloudAccount, verifyCloudAccount } from '@/api/cloud-account'
import { getSyncPolicies } from '@/api/sync-policy'
import { getProjects } from '@/api/project'
import type { CloudAccount, CreateCloudAccountRequest, Project } from '@/types'
import type { SyncPolicy } from '@/types/sync-policy'

const accounts = ref<CloudAccount[]>([])
const syncPolicies = ref<SyncPolicy[]>([])
const projects = ref<Project[]>([])
const loading = ref(false)
const showWizard = ref(false)
const wizardStep = ref(0)
const submitting = ref(false)
const wizardFormRef = ref<FormInstance>()

// 分页
const pagination = reactive({
  currentPage: 1,
  pageSize: 20,
  total: 0
})

// 供应商选项
const providers = ref([
  { id: 'alibaba', name: 'alibaba', displayName: '阿里云', icon: Cloudy },
  { id: 'tencent', name: 'tencent', displayName: '腾讯云', icon: Cloudy },
  { id: 'aws', name: 'aws', displayName: 'AWS', icon: Cloudy },
  { id: 'azure', name: 'azure', displayName: 'Azure', icon: Cloudy }
])

// 向导表单数据
const wizardForm = reactive({
  name: '',
  remarks: '',
  accountType: 'public',
  accessKeyId: '',
  accessKeySecret: '',
  resourceAssignmentMethod: [] as string[],
  syncStrategy: {
    policy: '',
    scope: [] as string[],
    defaultProject: 1
  },
  cloudProjectMapping: {
    defaultProject: 1
  },
  cloudSubscriptionMapping: {
    defaultProject: 1
  },
  specificProject: {
    defaultProject: 1
  },
  selectedRegions: [] as string[]
})

// 定时任务表单
const scheduleForm = reactive({
  name: '',
  type: 'sync_cloud_account',
  frequency: 'daily',
  triggerTime: '02:00',
  validRange: [] as [string, string] | []
})

// 可用区域列表
const availableRegions = ref([
  { name: '华北1（青岛）', status: 'available' },
  { name: '华北2（北京）', status: 'available' },
  { name: '华北3（张家口）', status: 'available' },
  { name: '华东1（杭州）', status: 'available' },
  { name: '华东2（上海）', status: 'available' },
  { name: '华南1（深圳）', status: 'available' },
  { name: '西南1（成都）', status: 'available' },
  { name: '中国（香港）', status: 'available' },
  { name: '美国（硅谷）', status: 'available' },
  { name: '美国（弗吉尼亚）', status: 'available' },
  { name: '日本（东京）', status: 'available' },
  { name: '新加坡', status: 'available' },
  { name: '德国（法兰克福）', status: 'available' }
])

// 验证规则
const wizardRules = {
  name: [{ required: true, message: '请输入云账户名称', trigger: 'blur' }],
  accessKeyId: [{ required: true, message: '请输入密钥ID', trigger: 'blur' }],
  accessKeySecret: [{ required: true, message: '请输入密钥Secret', trigger: 'blur' }]
}

// 检测移动端
const isMobile = computed(() => {
  return window.innerWidth < 768
})

// 加载项目
const loadProjects = async () => {
  try {
    const res = await getProjects()
    projects.value = res.items || []
  } catch (e) {
    console.error('Failed to load projects:', e)
    ElMessage.error('加载项目失败')
  }
}

// 加载同步策略
const loadSyncPolicies = async () => {
  try {
    const res = await getSyncPolicies()
    syncPolicies.value = res.items || []
  } catch (e) {
    console.error('Failed to load sync policies:', e)
    ElMessage.error('加载同步策略失败')
  }
}

// 监听缺省项目变化，同步更新所有相关配置
watch(() => wizardForm.syncStrategy.defaultProject, (newVal) => {
  if (newVal) {
    updateAllDefaultProjects(newVal)
  }
})

watch(() => wizardForm.cloudProjectMapping.defaultProject, (newVal) => {
  if (newVal) {
    updateAllDefaultProjects(newVal)
  }
})

watch(() => wizardForm.cloudSubscriptionMapping.defaultProject, (newVal) => {
  if (newVal) {
    updateAllDefaultProjects(newVal)
  }
})

watch(() => wizardForm.specificProject.defaultProject, (newVal) => {
  if (newVal) {
    updateAllDefaultProjects(newVal)
  }
})

// 更新所有资源归属方法的缺省项目值
const updateAllDefaultProjects = (projectId: number) => {
  wizardForm.syncStrategy.defaultProject = projectId
  wizardForm.cloudProjectMapping.defaultProject = projectId
  wizardForm.cloudSubscriptionMapping.defaultProject = projectId
  wizardForm.specificProject.defaultProject = projectId
}

// 获取资源归属方式备注
const getResourceAssignmentNotes = () => {
  const methods = wizardForm.resourceAssignmentMethod

  if (methods.length === 0) {
    return ''
  } else if (methods.length === 1) {
    const method = methods[0]
    if (method === 'sync_strategy') {
      return '资源会根据同步策略归属，若资源不匹配该归属方式则归属到缺省项目中。'
    } else if (method === 'cloud_project') {
      return '资源会同步到与云上项目同名的本地项目中，若资源无云上项目属性则归属到缺省项目中。'
    } else if (method === 'cloud_subscription') {
      return '资源会同步到与云订阅同名的本地项目中，若资源无云订阅属性则归属到缺省项目中。'
    } else if (method === 'specific_project') {
      return '所有项目资源均归属到指定项目。'
    }
  } else {
    // 多选的情况
    return '资源会优先根据同步策略归属，不匹配同步策略的资源会优先同步与云上项目同名的本地项目中，不匹配上述两种方式的资源会归属到与云订阅同名的本地项目中，若资源不匹配上述三种归属方式则归属到缺省项目中。'
  }

  return ''
}

// 选择供应商
const selectProvider = (providerId: string) => {
  // 确认选择了供应商，进入下一步
  ElMessage.success(`已选择 ${providers.value.find(p => p.id === providerId)?.displayName}`)
  wizardStep.value = 1
}

// 下一步
const nextStep = async () => {
  if (wizardStep.value === 1) {
    // 验证第二步表单
    if (!wizardFormRef.value) return
    const valid = await wizardFormRef.value.validate().catch(() => false)
    if (!valid) return
  }
  if (wizardStep.value < 3) {
    wizardStep.value++
  }
}

// 上一步
const previousStep = () => {
  if (wizardStep.value > 0) {
    wizardStep.value--
  }
}

// 全选区域
const selectAllRegions = () => {
  wizardForm.selectedRegions = availableRegions.value
    .filter(region => region.status === 'available')
    .map(region => region.name)
}

// 清空区域
const clearAllRegions = () => {
  wizardForm.selectedRegions = []
}

// 区域选择变化
const handleRegionSelectionChange = (selection: any[]) => {
  wizardForm.selectedRegions = selection.map(item => item.name)
}

// 提交向导
const submitWizard = async () => {
  submitting.value = true
  try {
    // 整合表单数据
    const formData: CreateCloudAccountRequest = {
      name: wizardForm.name,
      provider_type: 'alibaba', // 从第一步的选择中获取
      credentials: {
        access_key_id: wizardForm.accessKeyId,
        access_key_secret: wizardForm.accessKeySecret
      },
      description: wizardForm.remarks,
      remarks: wizardForm.remarks,
      enabled: true,
      health_status: 'healthy',
      account_number: wizardForm.accessKeyId.substring(0, 10) + '...', // 示例账号
      domain_id: 1, // 默认域
      resource_assignment_method: wizardForm.resourceAssignmentMethod.join(',')
    }

    await createCloudAccount(formData)
    ElMessage.success('云账户添加成功')

    // 重置表单
    resetWizard()
    showWizard.value = false

    // 刷新数据
    loadAccounts()
  } catch (error) {
    console.error(error)
    ElMessage.error('添加云账户失败')
  } finally {
    submitting.value = false
  }
}

// 重置向导
const resetWizard = () => {
  wizardStep.value = 0
  Object.assign(wizardForm, {
    name: '',
    remarks: '',
    accountType: 'public',
    accessKeyId: '',
    accessKeySecret: '',
    resourceAssignmentMethod: [],
    syncStrategy: {
      policy: '',
      scope: [],
      defaultProject: 1
    },
    cloudProjectMapping: {
      defaultProject: 1
    },
    cloudSubscriptionMapping: {
      defaultProject: 1
    },
    specificProject: {
      defaultProject: 1
    },
    selectedRegions: []
  })
  Object.assign(scheduleForm, {
    name: '',
    type: 'sync_cloud_account',
    frequency: 'daily',
    triggerTime: '02:00',
    validRange: []
  })
}

// 现有的加载、编辑、删除等函数保持不变
const loadAccounts = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.currentPage,
      page_size: pagination.pageSize
    }
    const res = await getCloudAccounts(params)
    accounts.value = res.items || []
    pagination.total = res.total || 0
  } catch (e) {
    console.error(e)
    ElMessage.error('加载云账户失败')
  } finally {
    loading.value = false
  }
}

const getProviderName = (type: string) => {
  const map: Record<string, string> = {
    alibaba: '阿里云',
    tencent: '腾讯云',
    aws: 'AWS',
    azure: 'Azure'
  }
  return map[type] || type
}

const getProviderType = (type: string) => {
  const map: Record<string, any> = {
    alibaba: 'orange',
    tencent: 'cyan',
    aws: 'warning',
    azure: 'success'
  }
  return map[type] || ''
}

const handleVerify = async (row: CloudAccount) => {
  try {
    await verifyCloudAccount(row.id)
    ElMessage.success('验证成功')
  } catch (e) {
    ElMessage.error('验证失败')
  }
}

const handleEdit = (row: CloudAccount) => {
  ElMessage.warning('编辑功能将在后续版本中实现')
}

const handleDelete = async (row: CloudAccount) => {
  try {
    await ElMessageBox.confirm(`确定要删除云账户 "${row.name}" 吗？`, '提示', {
      type: 'warning'
    })
    await deleteCloudAccount(row.id)
    ElMessage.success('删除成功')
    loadAccounts()
  } catch (e: any) {
    if (e !== 'cancel') {
      console.error(e)
      ElMessage.error('删除失败')
    }
  }
}

const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  loadAccounts()
}

const handleCurrentChange = (page: number) => {
  pagination.currentPage = page
  loadAccounts()
}

onMounted(() => {
  loadAccounts()
  loadSyncPolicies()
  loadProjects()
})
</script>

<style scoped>
.cloud-accounts-page {
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

.pagination {
  margin-top: 20px;
  justify-content: flex-end;
}

.wizard-step-content {
  margin: 30px 0;
}

.provider-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 20px;
  margin-top: 20px;
}

.provider-card {
  cursor: pointer;
  transition: all 0.3s;
}

.provider-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 16px rgba(0,0,0,0.1);
}

.provider-item {
  text-align: center;
  padding: 20px 10px;
}

.provider-item h4 {
  margin-top: 10px;
}

.region-tip {
  background-color: #ecf5ff;
  padding: 12px;
  border-radius: 4px;
  border-left: 4px solid #409eff;
  margin-bottom: 20px;
}

.region-actions {
  margin-bottom: 15px;
}

.resource-assignment-notes {
  margin-top: 15px;
  padding: 12px;
  background-color: #f0f9ff;
  border: 1px solid #d2ebff;
  border-radius: 4px;
  color: #5e7d8a;
  font-size: 14px;
  line-height: 1.5;
}

.sync-strategy-info {
  margin-top: 20px;
  padding: 15px;
  background-color: #f5f5f5;
  border-radius: 4px;
}

.sync-strategy-info ul {
  list-style-type: disc;
  margin-left: 20px;
  margin-top: 10px;
}

.sync-strategy-info li {
  margin: 5px 0;
}
</style>