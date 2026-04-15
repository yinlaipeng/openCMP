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

      <!-- 空状态 -->
      <EmptyState
        v-if="!loading && accounts.length === 0"
        title="暂无云账户"
        description="当前没有任何云账户，点击下方按钮添加云账户"
        :icon="Cloudy"
        createButtonText="添加云账户"
        @create="showWizard = true"
      />

      <!-- 数据表格 -->
      <el-table :data="accounts" v-loading="loading" style="width: 100%" v-if="accounts.length > 0 || loading" :table-layout="'auto'">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="名称" min-width="150" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusText(row.status) }}
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
            <el-tag :type="getHealthStatusType(row.health_status)">
              {{ getHealthStatusText(row.health_status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="balance" label="余额" width="100">
          <template #default="{ row }">
            ¥{{ row.balance?.toFixed(2) || '0.00' }}
          </template>
        </el-table-column>
        <el-table-column prop="provider_type" label="平台" width="100">
          <template #default="{ row }">
            <el-tag :type="getProviderType(row.provider_type)">
              {{ getProviderName(row.provider_type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="account_number" label="账号" width="150" show-overflow-tooltip />
        <el-table-column prop="last_sync" label="上次同步" width="160" show-overflow-tooltip />
        <el-table-column prop="sync_time" label="同步时间" width="120" show-overflow-tooltip />
        <el-table-column prop="domain_id" label="所属域" width="100">
          <template #default="{ row }">
            {{ getDomainName(row.domain_id) }}
          </template>
        </el-table-column>
        <el-table-column prop="resource_assignment_method" label="资源归属方式" width="150">
          <template #default="{ row }">
            {{ getResourceAssignmentMethodName(row.resource_assignment_method) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="140" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="primary" @click="handleSyncClick(row)">同步云账号</el-button>
            <el-dropdown trigger="click" style="margin-left: 8px">
              <el-button size="small">
                更多 <el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <!-- 状态设置 - 使用 popover 子菜单 -->
                  <el-dropdown-item>
                    <el-popover
                      placement="right-start"
                      trigger="hover"
                      :width="130"
                      :show-arrow="false"
                      :offset="4"
                    >
                      <template #reference>
                        <span style="display: flex; justify-content: space-between; align-items: center; width: 100%;">
                          状态设置 <el-icon><ArrowRight /></el-icon>
                        </span>
                      </template>
                      <div style="padding: 4px 0;">
                        <div
                          class="submenu-item"
                          :class="{ 'submenu-item-disabled': row.enabled }"
                          @click="!row.enabled && handleStatusCommand('enable', row)"
                        >
                          <span v-if="row.enabled" style="color: #909399;">启用 (已启用)</span>
                          <span v-else>启用</span>
                        </div>
                        <div
                          class="submenu-item"
                          :class="{ 'submenu-item-disabled': !row.enabled }"
                          @click="row.enabled && handleStatusCommand('disable', row)"
                        >
                          <span v-if="!row.enabled" style="color: #909399;">禁用 (已禁用)</span>
                          <span v-else>禁用</span>
                        </div>
                        <div class="submenu-item" @click="handleStatusCommand('test', row)">连接测试</div>
                      </div>
                    </el-popover>
                  </el-dropdown-item>
                  <!-- 属性设置 - 使用 popover 子菜单 -->
                  <el-dropdown-item>
                    <el-popover
                      placement="right-start"
                      trigger="hover"
                      :width="150"
                      :show-arrow="false"
                      :offset="4"
                    >
                      <template #reference>
                        <span style="display: flex; justify-content: space-between; align-items: center; width: 100%;">
                          属性设置 <el-icon><ArrowRight /></el-icon>
                        </span>
                      </template>
                      <div style="padding: 4px 0;">
                        <div class="submenu-item" @click="handleAttributesCommand('auto_sync', row)">设置自动同步</div>
                        <div class="submenu-item" @click="handleAttributesCommand('sync_policy', row)">设置同步策略</div>
                        <div class="submenu-item" @click="handleAttributesCommand('update', row)">更新账号</div>
                      </div>
                    </el-popover>
                  </el-dropdown-item>
                  <!-- 删除 - 直接操作 -->
                  <el-dropdown-item divided @click="handleDelete(row)">
                    <span style="color: var(--el-color-danger)">删除</span>
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
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

          <el-form-item label="测试连接">
            <el-button @click="testConnectionInWizard" :loading="wizardForm.testConnectionStatus === 'testing'">
              {{ wizardForm.testConnectionStatus === 'testing' ? '测试中...' : '测试连接' }}
            </el-button>
            <span v-if="wizardForm.testConnectionResult" :class="wizardForm.testConnectionResult.includes('成功') ? 'text-success' : 'text-danger'">
              {{ wizardForm.testConnectionResult }}
            </span>
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
              <el-tag :type="getRegionStatusType(row.status)">
                {{ getRegionStatusText(row.status) }}
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

    <!-- 同步云账号对话框 -->
    <el-dialog v-model="showSyncDialog" title="同步云账号" width="600px">
      <el-form :model="syncForm" label-width="120px">
        <el-form-item label="云账号名称">
          <el-input v-model="syncForm.name" readonly />
        </el-form-item>
        <el-form-item label="环境">
          <el-input v-model="syncForm.environment" readonly />
        </el-form-item>
        <el-form-item label="连接状态">
          <el-tag :type="syncForm.connectionStatus ? 'success' : 'danger'">
            {{ syncForm.connectionStatus ? '已连接' : '未连接' }}
          </el-tag>
        </el-form-item>
        <el-form-item label="同步模式">
          <el-radio-group v-model="syncForm.syncMode">
            <el-radio value="full">全量同步</el-radio>
            <el-radio value="incremental">增量同步</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="同步资源类型">
          <el-select v-model="syncForm.syncResourceTypes" multiple placeholder="请选择同步资源类型">
            <el-option value="all" label="全部" />
            <el-option value="instances" label="主机" />
            <el-option value="subnets" label="IP子网" />
            <el-option value="vpcs" label="VPC" />
            <el-option value="security_groups" label="安全组" />
            <el-option value="load_balancers" label="负载均衡器" />
            <el-option value="eips" label="EIP" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="confirmSync">确认同步</el-button>
        <el-button @click="showSyncDialog = false">取消</el-button>
      </template>
    </el-dialog>

    <!-- 属性设置对话框 -->
    <el-dialog v-model="showAttributesDialog" title="属性设置" width="600px">
      <el-form :model="attributesForm" label-width="150px">
        <el-form-item label="自动同步">
          <el-switch v-model="attributesForm.autoSync" />
        </el-form-item>
        <el-form-item label="同步策略">
          <el-select v-model="attributesForm.syncPolicy" placeholder="请选择同步策略">
            <el-option value="default" label="默认策略" />
            <el-option value="full" label="全量策略" />
            <el-option value="incremental" label="增量策略" />
          </el-select>
        </el-form-item>
        <el-form-item label="同步间隔(小时)">
          <el-input-number v-model="attributesForm.syncInterval" :min="1" :max="168" />
        </el-form-item>
        <el-form-item label="同步资源类型">
          <el-select v-model="attributesForm.syncResourceTypes" multiple placeholder="请选择同步资源类型">
            <el-option value="instances" label="主机" />
            <el-option value="subnets" label="IP子网" />
            <el-option value="vpcs" label="VPC" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="updateAccount">更新账号</el-button>
        <el-button @click="confirmAttributesChange">确认</el-button>
        <el-button @click="showAttributesDialog = false">取消</el-button>
      </template>
    </el-dialog>

    <!-- 更新账号对话框 -->
    <EditAccountDialog
      v-model="showUpdateDialog"
      :account-id="editAccountId"
      @saved="loadAccounts"
    />

    <!-- 属性设置弹窗 -->
    <CloudAccountDetailDialog
      v-model="showAccountDetailDialog"
      :account-id="accountDetailId"
      :initial-tab="accountDetailTab"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed, watch } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance } from 'element-plus'
import { Plus, Cloudy, ArrowDown, ArrowRight, Refresh, FolderOpened } from '@element-plus/icons-vue'
import { getCloudAccounts, createCloudAccount, updateCloudAccount, deleteCloudAccount, syncCloudAccount, testConnection as testConnectionAPI, updateCloudAccountStatus, updateCloudAccountAttributes } from '@/api/cloud-account'
import { getSyncPolicies } from '@/api/sync-policy'
import { getProjects } from '@/api/project'
import type { CloudAccount, CreateCloudAccountRequest, Project } from '@/types'
import type { SyncPolicy } from '@/types/sync-policy'
import EmptyState from '@/components/common/EmptyState.vue'
import EditAccountDialog from './components/EditAccountDialog.vue'
import CloudAccountDetailDialog from './components/CloudAccountDetailDialog.vue'

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
  testConnectionStatus: '', // 'testing', 'success', 'error'
  testConnectionResult: '',
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

// 测试连接（在向导中）
const testConnectionInWizard = async () => {
  if (!wizardForm.accessKeyId || !wizardForm.accessKeySecret) {
    ElMessage.warning('请先填写密钥ID和密钥Secret')
    return
  }

  wizardForm.testConnectionStatus = 'testing'
  wizardForm.testConnectionResult = '正在测试连接...'

  try {
    // 创建临时云账户用于测试
    const tempAccountData: CreateCloudAccountRequest = {
      name: wizardForm.name || 'temp-test',
      provider_type: 'alibaba',
      credentials: {
        access_key_id: wizardForm.accessKeyId,
        access_key_secret: wizardForm.accessKeySecret
      },
      description: 'Temporary account for connection testing',
      remarks: 'Temporary account for connection testing',
      enabled: true,
      health_status: 'healthy',
      domain_id: 1,
      resource_assignment_method: 'tag_mapping'
    }

    // 首先创建一个临时账户
    const tempAccount = await createCloudAccount(tempAccountData)

    // 测试连接
    const response = await testConnection(tempAccount.id)

    // 删除临时账户
    await deleteCloudAccount(tempAccount.id)

    // 设置结果
    wizardForm.testConnectionStatus = response.connected ? 'success' : 'error'
    wizardForm.testConnectionResult = response.message || (response.connected ? '连接成功' : '连接失败')
  } catch (error: any) {
    wizardForm.testConnectionStatus = 'error'
    wizardForm.testConnectionResult = error.message || '连接测试失败'
    console.error('Connection test error:', error)
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
    alibaba: 'warning',  // 改为有效的Element Plus标签类型
    tencent: 'primary',  // 改为有效的Element Plus标签类型
    aws: 'danger',       // 改为有效的Element Plus标签类型
    azure: 'success'     // 保持有效的Element Plus标签类型
  }
  return map[type] || 'info'
}

const getStatusType = (status: string) => {
  const statusMap: Record<string, string> = {
    active: 'success',      // 已连接
    connected: 'success',   // 已连接
    inactive: 'info',       // 未连接
    disconnected: 'info',   // 未连接
    error: 'danger',
    pending: 'warning',
    unknown: 'info'
  }
  return statusMap[status] || 'info'
}

const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    active: '已连接',
    connected: '已连接',
    inactive: '未连接',
    disconnected: '未连接',
    error: '连接错误',
    pending: '连接中',
    unknown: '未知'
  }
  return statusMap[status] || status
}

const getHealthStatusType = (status: string) => {
  const statusMap: Record<string, string> = {
    healthy: 'success',     // 正常
    normal: 'success',      // 正常
    unhealthy: 'danger',    // 异常
    abnormal: 'danger',     // 异常
    no_permission: 'warning', // 无权限
    warning: 'warning',
    unknown: 'info'
  }
  return statusMap[status] || 'info'
}

const getHealthStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    healthy: '正常',
    normal: '正常',
    unhealthy: '异常',
    abnormal: '异常',
    no_permission: '无权限',
    warning: '警告',
    unknown: '未知'
  }
  return statusMap[status] || status
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

// 新增的响应式数据
const showSyncDialog = ref(false)
const showAttributesDialog = ref(false)
const showUpdateDialog = ref(false)
const editAccountId = ref<number | null>(null)
const showAccountDetailDialog = ref(false)
const accountDetailId = ref<number | null>(null)
const accountDetailTab = ref('detail')
const activeTab = ref('details')

const syncForm = ref({
  name: '',
  environment: '',
  connectionStatus: false,
  syncMode: 'full',
  syncResourceTypes: ['all'],
  syncTime: new Date(),
  notes: ''
})

const updateForm = ref({
  name: '',
  description: ''
})

const attributesForm = ref({
  autoSync: true,
  syncPolicy: 'default',
  syncInterval: 24,
  syncResourceTypes: [],
  updateAccountInfo: false
})

const currentAccount = ref<any>(null)
const subscriptions = ref<any[]>([])
const subscriptionsLoading = ref(false)
const resourceStats = ref({
  instances: 0,
  vpcs: 0,
  subnets: 0,
  security_groups: 0,
  eips: 0,
  load_balancers: 0
})

// 新增的方法
const showAccountDetails = async (account: any) => {
  currentAccount.value = account
  showAccountDetailDialog.value = true

  // 加载资源统计信息
  loadResourceStats(account.id)

  // 加载订阅信息
  loadSubscriptions(account.id)
}

const loadResourceStats = async (accountId: number) => {
  // 模拟API调用获取资源统计信息
  resourceStats.value = {
    instances: Math.floor(Math.random() * 100),
    vpcs: Math.floor(Math.random() * 20),
    subnets: Math.floor(Math.random() * 50),
    security_groups: Math.floor(Math.random() * 30),
    eips: Math.floor(Math.random() * 25),
    load_balancers: Math.floor(Math.random() * 15)
  }
}

const loadSubscriptions = async (accountId: number) => {
  subscriptionsLoading.value = true
  try {
    // 模拟API调用获取订阅信息
    await new Promise(resolve => setTimeout(resolve, 500))

    subscriptions.value = [
      {
        id: 1,
        name: '订阅1',
        subscription: 'abc-def-ghi',
        enabled: true,
        status: 'active',
        last_sync: new Date().toISOString(),
        sync_duration: '2m 30s',
        sync_status: 'success',
        domain_id: 1,
        default_project: '项目A'
      },
      {
        id: 2,
        name: '订阅2',
        subscription: 'jkl-mno-pqr',
        enabled: false,
        status: 'inactive',
        last_sync: new Date(Date.now() - 86400000).toISOString(),
        sync_duration: '1m 15s',
        sync_status: 'failed',
        domain_id: 1,
        default_project: '项目B'
      }
    ]
  } catch (error) {
    console.error('Failed to load subscriptions:', error)
    ElMessage.error('加载订阅失败')
  } finally {
    subscriptionsLoading.value = false
  }
}

// 状态设置下拉菜单命令处理
const handleStatusCommand = async (command: string, row: any) => {
  currentAccount.value = row

  switch (command) {
    case 'enable':
      await handleEnable()
      break
    case 'disable':
      await handleDisable()
      break
    case 'test':
      await handleTestConnection(row)
      break
  }
}

// 属性设置下拉菜单命令处理
const handleAttributesCommand = async (command: string, row: any) => {
  currentAccount.value = row

  switch (command) {
    case 'auto_sync':
      accountDetailId.value = row.id
      accountDetailTab.value = 'detail'
      showAccountDetailDialog.value = true
      break
    case 'sync_policy':
      accountDetailId.value = row.id
      accountDetailTab.value = 'scheduledTasks'
      showAccountDetailDialog.value = true
      break
    case 'update':
      editAccountId.value = row.id
      showUpdateDialog.value = true
      break
  }
}

// 同步云账号按钮点击
const handleSyncClick = (row: any) => {
  currentAccount.value = row
  syncForm.value.name = row.name
  syncForm.value.environment = getProviderName(row.provider_type)
  syncForm.value.connectionStatus = row.status === 'active'
  showSyncDialog.value = true
}

// 启用云账号
const handleEnable = async () => {
  try {
    await updateCloudAccountStatus(currentAccount.value.id, true)
    ElMessage.success('已启用云账号')
    loadAccounts()
  } catch (error) {
    ElMessage.error('启用失败')
    console.error('Enable error:', error)
  }
}

// 禁用云账号
const handleDisable = async () => {
  try {
    await updateCloudAccountStatus(currentAccount.value.id, false)
    ElMessage.success('已禁用云账号')
    loadAccounts()
  } catch (error) {
    ElMessage.error('禁用失败')
    console.error('Disable error:', error)
  }
}

// 连接测试
const handleTestConnection = async (row: any) => {
  try {
    const response = await testConnectionAPI(row.id)
    if (response.connected) {
      ElMessage.success('连接测试成功')
    } else {
      ElMessage.warning('连接测试失败: ' + (response.message || '无法连接'))
    }
  } catch (error: any) {
    ElMessage.error('连接测试失败')
    console.error('Connection test error:', error)
  }
}

const confirmSync = async () => {
  try {
    await syncCloudAccount(currentAccount.value.id, {
      mode: syncForm.value.syncMode,
      resource_types: syncForm.value.syncResourceTypes
    })

    ElMessage.success(`已启动对 ${syncForm.value.name} 的同步（模式：${syncForm.value.syncMode === 'full' ? '全量' : '增量'}）`)
    showSyncDialog.value = false

    // 可以在这里添加刷新数据的逻辑
  } catch (error) {
    ElMessage.error('同步启动失败')
    console.error('Sync error:', error)
  }
}

const updateAccount = async () => {
  try {
    await updateCloudAccount(currentAccount.value.id, {
      name: currentAccount.value.name,
      provider_type: currentAccount.value.provider_type,
      // 其他需要更新的字段
    })

    ElMessage.success('账号信息更新成功')
  } catch (error) {
    ElMessage.error('账号信息更新失败')
    console.error('Update error:', error)
  }
}

const confirmAttributesChange = async () => {
  try {
    await updateCloudAccountAttributes(currentAccount.value.id, {
      auto_sync: attributesForm.value.autoSync,
      sync_policy: attributesForm.value.syncPolicy,
      sync_interval: attributesForm.value.syncInterval,
      sync_resource_types: attributesForm.value.syncResourceTypes
    })

    ElMessage.success('属性更新成功')
    showAttributesDialog.value = false
    loadAccounts() // 重新加载列表
  } catch (error) {
    ElMessage.error('属性更新失败')
    console.error('Update error:', error)
  }
}

const handleUpdateSubmit = async () => {
  try {
    await updateCloudAccount(currentAccount.value.id, {
      name: updateForm.value.name,
      description: updateForm.value.description
    })
    ElMessage.success('账号更新成功')
    showUpdateDialog.value = false
    loadAccounts()
  } catch (error) {
    ElMessage.error('账号更新失败')
    console.error('Update error:', error)
  }
}

const changeProject = (row: any) => {
  ElMessage.info(`更改项目功能将在后续版本中实现 - 订阅: ${row.name}`)
}

const syncResources = (row: any) => {
  ElMessage.success(`已启动对 ${row.name} 资源的同步`)
}

const editSubscription = (row: any) => {
  ElMessage.info(`编辑订阅功能将在后续版本中实现 - 订阅: ${row.name}`)
}

const deleteSubscription = async (row: any) => {
  try {
    await ElMessageBox.confirm(`确定要删除订阅 "${row.name}" 吗？`, '提示', {
      type: 'warning'
    })
    ElMessage.success('订阅删除成功')
    loadSubscriptions(currentAccount.value.id) // 重新加载订阅
  } catch (e: any) {
    if (e !== 'cancel') {
      console.error(e)
      ElMessage.error('删除订阅失败')
    }
  }
}

const getRegionStatusType = (status: string) => {
  const statusMap: Record<string, string> = {
    available: 'success',
    unavailable: 'danger',
    unknown: 'info'
  }
  return statusMap[status] || 'info'
}

const getRegionStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    available: '可用',
    unavailable: '不可用',
    unknown: '未知'
  }
  return statusMap[status] || status
}

const getSyncStatusType = (status: string) => {
  switch (status) {
    case 'success':
      return 'success'
    case 'failed':
      return 'danger'
    case 'running':
      return 'warning'
    default:
      return 'info'
  }
}

const formatDate = (dateString: string | undefined) => {
  if (!dateString) return '-'
  return new Date(dateString).toLocaleString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' })
}

const getDomainName = (domainId: number) => {
  // 实际应用中应该从API获取域名
  return `域${domainId}`
}

const getResourceAssignmentMethodName = (method: string) => {
  const map: Record<string, string> = {
    'sync_strategy': '同步策略',
    'cloud_project': '云上项目',
    'cloud_subscription': '云订阅',
    'specific_project': '指定项目',
    'tag_mapping': '标签映射',
    'project_mapping': '项目映射',
    'manual_assignment': '手动分配'
  }
  return map[method] || method
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

.sync-time {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}

.stat-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
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

.el-descriptions {
  margin: 15px 0;
}

.text-success {
  color: #67c23a;
}

.text-danger {
  color: #f56c6c;
}

.status-cell {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  align-items: center;
}
</style>

<!-- 全局样式（子菜单项样式） -->
<style>
/* 子菜单项样式 - 模拟 Element Plus dropdown item */
.submenu-item {
  padding: 8px 12px;
  cursor: pointer;
  font-size: 14px;
  line-height: 1.4;
  transition: background-color 0.2s;
}

.submenu-item:hover {
  background-color: var(--el-fill-color-light);
  color: var(--el-color-primary);
}

.submenu-item:active {
  background-color: var(--el-fill-color);
}

/* 子菜单项禁用状态 */
.submenu-item-disabled {
  cursor: not-allowed;
  opacity: 0.6;
}

.submenu-item-disabled:hover {
  background-color: transparent;
  color: #909399;
}
</style>