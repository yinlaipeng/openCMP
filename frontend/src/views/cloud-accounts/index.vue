<template>
  <div class="cloud-accounts-container">
    <div class="page-header">
      <h2>云账户管理</h2>
      <div class="tool-bar">
        <el-button @click="loadAccounts" :loading="loading">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
        <el-button type="primary" @click="showWizard = true">
          <el-icon><Plus /></el-icon>
          新建
        </el-button>
        <el-button @click="handleBatchOperation" :disabled="selectedIds.length === 0">
          <el-icon><Operation /></el-icon>
          批量操作
        </el-button>
        <el-button @click="handleExport">
          <el-icon><Download /></el-icon>
          导出
        </el-button>
        <el-button @click="handleSettings">
          <el-icon><Setting /></el-icon>
          设置
        </el-button>
      </div>
    </div>

    <!-- 顶部Tab分类 -->
    <el-tabs v-model="activeCategory" class="top-tabs" @tab-change="handleCategoryChange">
      <el-tab-pane label="全部" name="all" />
      <el-tab-pane label="公有云" name="public" />
    </el-tabs>

    <!-- 轻量搜索栏 -->
    <div class="search-bar">
      <!-- 字段选择器 -->
      <el-dropdown trigger="click" @command="handleFieldChange" class="field-selector">
        <el-button>
          {{ currentFieldLabel }}
          <el-icon class="el-icon--right"><ArrowDown /></el-icon>
        </el-button>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="name" :class="{ active: searchField === 'name' }">
              <el-icon><Search /></el-icon> 名称
            </el-dropdown-item>
            <el-dropdown-item command="id" :class="{ active: searchField === 'id' }">
              <el-icon><Ticket /></el-icon> ID
            </el-dropdown-item>
            <el-dropdown-item command="remarks" :class="{ active: searchField === 'remarks' }">
              <el-icon><Document /></el-icon> 备注
            </el-dropdown-item>
            <el-dropdown-item command="provider_type" :class="{ active: searchField === 'provider_type' }">
              <el-icon><Cloudy /></el-icon> 平台
            </el-dropdown-item>
            <el-dropdown-item command="status" :class="{ active: searchField === 'status' }">
              <el-icon><Connection /></el-icon> 状态
            </el-dropdown-item>
            <el-dropdown-item command="enabled" :class="{ active: searchField === 'enabled' }">
              <el-icon><Switch /></el-icon> 启用状态
            </el-dropdown-item>
            <el-dropdown-item command="health_status" :class="{ active: searchField === 'health_status' }">
              <el-icon><CircleCheck /></el-icon> 健康状态
            </el-dropdown-item>
            <el-dropdown-item command="account_number" :class="{ active: searchField === 'account_number' }">
              <el-icon><User /></el-icon> 账号
            </el-dropdown-item>
            <el-dropdown-item command="domain_id" :class="{ active: searchField === 'domain_id' }">
              <el-icon><OfficeBuilding /></el-icon> 域
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>

      <!-- 搜索输入区 -->
      <div class="search-input-wrapper" v-if="isTextField">
        <el-input
          v-model="searchKeyword"
          :placeholder="getPlaceholder()"
          clearable
          @keyup.enter="handleSearch"
          style="width: 300px"
        >
          <template #suffix>
            <el-icon class="search-hint-icon" @click="showSearchTip">
              <QuestionFilled />
            </el-icon>
          </template>
        </el-input>
      </div>

      <!-- 下拉选择区（枚举字段） -->
      <div class="search-select-wrapper" v-else>
        <el-select
          v-model="searchSelectValue"
          :placeholder="getPlaceholder()"
          clearable
          style="width: 200px"
        >
          <el-option v-for="opt in getFieldOptions()" :key="opt.value" :label="opt.label" :value="opt.value" />
        </el-select>
      </div>

      <!-- 搜索提示 -->
      <el-tooltip
        v-model:visible="showTip"
        content="默认为名称搜索，自动匹配 IP 或 ID 搜索项，IP 或 ID 多个搜索用英文竖线(|)隔开"
        placement="top"
        effect="light"
      >
        <span></span>
      </el-tooltip>

      <!-- 操作按钮 -->
      <el-button type="primary" @click="handleSearch">
        <el-icon><Search /></el-icon>
        查询
      </el-button>
      <el-button @click="handleResetSearch">
        重置
      </el-button>
    </div>

    <!-- 数据表格 -->
    <el-table
      :data="accounts"
      v-loading="loading"
      style="width: 100%"
      row-key="id"
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="55" />
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="名称" min-width="150">
        <template #default="{ row }">
          <el-link type="primary" @click="openDetailDrawer(row)">{{ row.name }}</el-link>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="连接状态" width="100">
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
      <el-table-column label="上次同步耗时" width="120">
        <template #default="{ row }">
          {{ row.sync_duration ? row.sync_duration + 's' : '-' }}
        </template>
      </el-table-column>
      <el-table-column prop="last_sync" label="同步时间" width="160" show-overflow-tooltip />
      <el-table-column prop="domain_id" label="所属域" width="100">
        <template #default="{ row }">
          {{ getDomainName(row.domain_id) }}
        </template>
      </el-table-column>
      <el-table-column label="资源归属方式" width="140">
        <template #default="{ row }">
          {{ getResourceAssignmentText(row.resource_assignment_method) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="160" fixed="right">
        <template #default="{ row }">
          <el-button size="small" type="primary" @click="handleSyncClick(row)">同步</el-button>
          <el-dropdown trigger="click" style="margin-left: 8px" @command="(cmd: string) => handleDropdownCommand(cmd, row)">
            <el-button size="small">
              更多 <el-icon class="el-icon--right"><ArrowDown /></el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="statusSetting">状态设置</el-dropdown-item>
                <el-dropdown-item command="attributeSetting">属性设置</el-dropdown-item>
                <el-dropdown-item divided command="enable" :disabled="row.enabled">启用</el-dropdown-item>
                <el-dropdown-item command="disable" :disabled="!row.enabled">禁用</el-dropdown-item>
                <el-dropdown-item command="test">连接测试</el-dropdown-item>
                <el-dropdown-item command="update">更新账号</el-dropdown-item>
                <el-dropdown-item divided command="setSyncAttribution">设置同步归属策略</el-dropdown-item>
                <el-dropdown-item command="setSyncPolicy">设置同步策略</el-dropdown-item>
                <el-dropdown-item command="setProxy">设置代理</el-dropdown-item>
                <el-dropdown-item command="setPasswordless">设置免密登录</el-dropdown-item>
                <el-dropdown-item command="setReadOnly">只读模式</el-dropdown-item>
                <el-dropdown-item divided command="delete" :disabled="row.enabled">
                  <span :style="{ color: row.enabled ? 'var(--el-text-color-disabled)' : 'var(--el-color-danger)' }">删除</span>
                  <el-tooltip v-if="row.enabled" content="账号为启用状态，请先禁用后再删除" placement="right">
                    <el-icon style="margin-left: 4px; color: var(--el-text-color-secondary)"><QuestionFilled /></el-icon>
                  </el-tooltip>
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      v-model:current-page="currentPage"
      v-model:page-size="pageSize"
      :total="total"
      :page-sizes="[10, 20, 50, 100]"
      layout="total, sizes, prev, pager, next, jumper"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
      class="pagination"
    />

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

        <!-- 公有云 -->
        <div class="provider-category">
          <h4>公有云</h4>
          <div class="provider-grid">
            <el-card
              v-for="provider in publicCloudProviders"
              :key="provider.id"
              class="provider-card"
              :class="{ selected: selectedProvider === provider.id }"
              @click="selectProvider(provider.id)"
            >
              <div class="provider-item">
                <el-icon :size="40"><component :is="provider.icon" /></el-icon>
                <h4>{{ provider.displayName }}</h4>
              </div>
            </el-card>
          </div>
        </div>

        <!-- 私有云 & 虚拟化平台 -->
        <div class="provider-category">
          <h4>私有云 & 虚拟化平台</h4>
          <div class="provider-grid">
            <el-card
              v-for="provider in privateCloudProviders"
              :key="provider.id"
              class="provider-card"
              :class="{ selected: selectedProvider === provider.id }"
              @click="selectProvider(provider.id)"
            >
              <div class="provider-item">
                <el-icon :size="40"><component :is="provider.icon" /></el-icon>
                <h4>{{ provider.displayName }}</h4>
              </div>
            </el-card>
          </div>
        </div>

        <!-- 对象存储 -->
        <div class="provider-category">
          <h4>对象存储</h4>
          <div class="provider-grid">
            <el-card
              v-for="provider in storageProviders"
              :key="provider.id"
              class="provider-card"
              :class="{ selected: selectedProvider === provider.id }"
              @click="selectProvider(provider.id)"
            >
              <div class="provider-item">
                <el-icon :size="40"><component :is="provider.icon" /></el-icon>
                <h4>{{ provider.displayName }}</h4>
              </div>
            </el-card>
          </div>
        </div>

        <!-- 已选择提示 -->
        <div class="selected-indicator" v-if="selectedProvider">
          <el-tag type="success" size="large">
            已选择：{{ getProviderDisplayName(selectedProvider) }}
          </el-tag>
        </div>
      </div>

      <!-- 步骤2: 配置云账号 -->
      <div v-if="wizardStep === 1" class="wizard-step-content">
        <el-form :model="wizardForm" :rules="wizardRules" ref="wizardFormRef" label-width="150px">
          <el-form-item label="名称" prop="name">
            <el-input v-model="wizardForm.name" placeholder="请输入云账户名称" />
          </el-form-item>

          <el-form-item label="备注">
            <el-input v-model="wizardForm.remarks" type="textarea" :rows="2" placeholder="请输入备注" />
          </el-form-item>

          <el-form-item label="账号类型">
            <el-select v-model="wizardForm.accountType" placeholder="请选择账号类型">
              <el-option label="主账号" value="primary" />
              <el-option label="子账号" value="sub" />
            </el-select>
          </el-form-item>

          <el-form-item label="密钥ID" prop="accessKeyId">
            <el-input v-model="wizardForm.accessKeyId" placeholder="请输入Access Key ID" />
          </el-form-item>

          <el-form-item label="密码/Secret" prop="accessKeySecret">
            <el-input v-model="wizardForm.accessKeySecret" type="password" show-password placeholder="请输入Access Key Secret" />
          </el-form-item>

          <el-divider content-position="left">资源归属设置</el-divider>

          <el-form-item label="资源归属方式">
            <el-checkbox-group v-model="wizardForm.resourceAssignmentMethods">
              <el-checkbox label="tag_mapping">根据同步策略归属</el-checkbox>
              <el-checkbox label="project_mapping">根据云上项目归属</el-checkbox>
              <el-checkbox label="subscription_mapping">根据云订阅归属</el-checkbox>
              <el-checkbox label="specify_project">指定项目</el-checkbox>
            </el-checkbox-group>
            <!-- 动态生成的归属说明 -->
            <div class="assignment-description" v-if="wizardForm.resourceAssignmentMethods.length > 0">
              <el-icon><InfoFilled /></el-icon>
              <span>{{ resourceAssignmentDescription }}</span>
            </div>
            <!-- 校验提示 -->
            <div class="assignment-hint" v-if="resourceAssignmentValidationHint">
              <el-icon><WarningFilled /></el-icon>
              <span>{{ resourceAssignmentValidationHint }}</span>
            </div>
          </el-form-item>

          <!-- 同步策略选择器（根据勾选状态联动显示） -->
          <el-form-item label="同步策略" v-if="resourceAssignmentControls.showSyncPolicySelector">
            <el-select v-model="wizardForm.syncPolicyId" placeholder="请选择同步策略" clearable>
              <el-option v-for="policy in syncPolicies" :key="policy.id" :label="policy.name" :value="policy.id" />
            </el-select>
            <span class="form-hint">选择同步策略后，资源将按策略规则自动归属</span>
          </el-form-item>

          <!-- 同步策略生效范围（勾选"根据同步策略归属"时显示） -->
          <el-form-item label="策略生效范围" v-if="resourceAssignmentControls.showSyncScopeSelector">
            <el-checkbox-group v-model="wizardForm.syncScope">
              <el-checkbox label="resource_tag">资源标签</el-checkbox>
              <el-checkbox label="project_tag">项目标签</el-checkbox>
            </el-checkbox-group>
          </el-form-item>

          <!-- 指定项目选择器（勾选"指定项目"时必填显示） -->
          <el-form-item
            label="目标项目"
            v-if="resourceAssignmentControls.showSpecifyProjectSelector"
            :required="true"
          >
            <el-select v-model="wizardForm.specifyProjectId" placeholder="请选择目标项目" filterable>
              <el-option v-for="project in projects" :key="project.id" :label="project.name" :value="project.id" />
            </el-select>
            <span class="form-hint">所有资源将归属到此项目</span>
          </el-form-item>

          <!-- 缺省项目选择器（存在兜底逻辑时显示） -->
          <el-form-item
            label="缺省项目"
            v-if="resourceAssignmentControls.showDefaultProjectSelector"
            :required="resourceAssignmentControls.needsDefaultProject"
          >
            <el-select v-model="wizardForm.defaultProjectId" placeholder="请选择缺省项目" clearable filterable>
              <el-option v-for="project in projects" :key="project.id" :label="project.name" :value="project.id" />
            </el-select>
            <span class="form-hint">未匹配任何归属规则的资源将归属到此项目</span>
          </el-form-item>

          <el-divider content-position="left">高级设置</el-divider>

          <el-form-item label="屏蔽同步资源">
            <el-switch v-model="wizardForm.blockSyncResources" />
          </el-form-item>

          <el-form-item label="代理">
            <el-switch v-model="wizardForm.enableProxy" />
            <el-select v-if="wizardForm.enableProxy" v-model="wizardForm.proxyId" placeholder="选择代理" style="margin-left: 10px; width: 200px">
              <el-option label="默认代理" value="default" />
              <el-option label="自定义代理" value="custom" />
            </el-select>
          </el-form-item>

          <el-form-item label="开启免密登录">
            <el-switch v-model="wizardForm.enablePasswordless" />
          </el-form-item>

          <el-form-item label="只读模式">
            <el-switch v-model="wizardForm.readOnlyMode" />
          </el-form-item>

          <el-form-item label="连接测试">
            <el-button @click="testConnectionInWizard" :loading="wizardForm.testConnectionStatus === 'testing'">
              {{ wizardForm.testConnectionStatus === 'testing' ? '测试中...' : '测试连接' }}
            </el-button>
            <span v-if="wizardForm.testConnectionResult" :class="wizardForm.testConnectionResult.includes('成功') ? 'text-success' : 'text-danger'">
              {{ wizardForm.testConnectionResult }}
            </span>
          </el-form-item>
        </el-form>
      </div>

      <!-- 步骤3: 配置同步区域 -->
      <div v-if="wizardStep === 2" class="wizard-step-content">
        <p class="region-tip">
          请选择要同步资源的区域，默认同步所有区域
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

          <el-form-item label="触发频次">
            <el-select v-model="scheduleForm.frequency">
              <el-option label="每天" value="daily" />
              <el-option label="每周" value="weekly" />
              <el-option label="每月" value="monthly" />
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
        </el-form>
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
    <el-dialog v-model="showSyncDialog" title="同步云账号" width="700px">
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
          <el-radio-group v-model="syncForm.syncAll" style="margin-bottom: 10px;">
            <el-radio value="all">全部资源类型</el-radio>
            <el-radio value="specific">指定资源类型</el-radio>
          </el-radio-group>
          <div v-if="syncForm.syncAll === 'specific'" style="border: 1px solid #dcdfe6; padding: 10px; border-radius: 4px; max-height: 300px; overflow-y: auto;">
            <el-checkbox-group v-model="syncForm.syncResourceTypes">
              <el-row>
                <el-col :span="8" v-for="rt in supportedResourceTypes" :key="rt.id">
                  <el-checkbox :value="rt.id" :label="rt.id">{{ rt.name }}</el-checkbox>
                </el-col>
              </el-row>
            </el-checkbox-group>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="confirmSync" type="primary" :loading="syncing">确认同步</el-button>
        <el-button @click="showSyncDialog = false">取消</el-button>
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

    <!-- 设置同步归属策略弹窗 -->
    <SetSyncAttributionDialog
      v-model="showSetSyncAttributionDialog"
      :account-id="editAccountId"
      @saved="loadAccounts"
    />

    <!-- 状态设置弹窗 -->
    <StatusSettingDialog
      v-model="showStatusSettingDialog"
      :account-id="editAccountId"
      @saved="loadAccounts"
    />

    <!-- 属性设置弹窗 -->
    <AttributeSettingDialog
      v-model="showAttributeSettingDialog"
      :account-id="editAccountId"
      @saved="loadAccounts"
    />

    <!-- 设置代理弹窗 -->
    <SetProxyDialog
      v-model="showSetProxyDialog"
      :account-id="editAccountId"
      @saved="loadAccounts"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance } from 'element-plus'
import { Plus, Cloudy, ArrowDown, Loading, Refresh, Operation, Download, Setting, QuestionFilled, Search, Document, Ticket, Connection, CircleCheck, User, OfficeBuilding, InfoFilled, WarningFilled } from '@element-plus/icons-vue'
import { getCloudAccounts, createCloudAccount, updateCloudAccount, deleteCloudAccount, syncCloudAccount, testConnection as testConnectionAPI, updateCloudAccountStatus, updateCloudAccountAttributes, getSupportedResourceTypes, type ResourceType, type CloudAccountSearchParams } from '@/api/cloud-account'
import { getSyncPolicies } from '@/api/sync-policy'
import { getProjects } from '@/api/project'
import { getDomains } from '@/api/iam'
import { getLatestSyncLog } from '@/api/sync-log'
import { createScheduledTask } from '@/api/scheduled-task'
import type { CloudAccount, CreateCloudAccountRequest, Project } from '@/types'
import type { SyncPolicy } from '@/types/sync-policy'
import EditAccountDialog from './components/EditAccountDialog.vue'
import CloudAccountDetailDialog from './components/CloudAccountDetailDialog.vue'
import SetSyncAttributionDialog from './components/SetSyncAttributionDialog.vue'
import StatusSettingDialog from './components/StatusSettingDialog.vue'
import { useResourceAssignmentDescription } from './composables/useResourceAssignmentDescription'
import AttributeSettingDialog from './components/AttributeSettingDialog.vue'
import SetProxyDialog from './components/SetProxyDialog.vue'

const accounts = ref<CloudAccount[]>([])
const syncPolicies = ref<SyncPolicy[]>([])
const projects = ref<Project[]>([])
const domains = ref<{ id: number; name: string }[]>([])
const loading = ref(false)
const syncStatusMap = ref<Map<number, { status: string; startTime: Date; duration?: number }>>(new Map())
const showWizard = ref(false)
const wizardStep = ref(0)
const submitting = ref(false)
const wizardFormRef = ref<FormInstance>()

// 搜索相关状态
const searchField = ref('name')          // 当前搜索字段
const searchKeyword = ref('')            // 文本搜索关键字
const searchSelectValue = ref<any>(undefined)  // 下拉选择值
const showTip = ref(false)               // 搜索提示显示状态
const selectedIds = ref<number[]>([])    // 表格选中ID列表

// 顶部Tab
const activeCategory = ref('all')

// 分页数据
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

// 搜索字段配置
const fieldLabels: Record<string, string> = {
  name: '名称',
  id: 'ID',
  remarks: '备注',
  provider_type: '平台',
  status: '状态',
  enabled: '启用状态',
  health_status: '健康状态',
  account_number: '账号',
  domain_id: '域'
}

// 文本类型字段
const textFieldTypes = ['name', 'id', 'remarks', 'account_number']

// 计算当前字段标签
const currentFieldLabel = computed(() => fieldLabels[searchField.value] || '名称')

// 判断当前字段是否为文本输入
const isTextField = computed(() => textFieldTypes.includes(searchField.value))

// 获取placeholder
const getPlaceholder = () => {
  if (searchField.value === 'name') {
    return '默认为名称搜索，自动匹配 IP 或 ID...'
  }
  if (searchField.value === 'id') {
    return '输入ID，多个用|分隔...'
  }
  return `输入${currentFieldLabel.value}搜索...`
}

// 获取下拉选项
const getFieldOptions = () => {
  switch (searchField.value) {
    case 'provider_type':
      return [
        { label: '阿里云', value: 'alibaba' },
        { label: '腾讯云', value: 'tencent' },
        { label: 'AWS', value: 'aws' },
        { label: 'Azure', value: 'azure' },
        { label: '华为云', value: 'huawei' },
        { label: 'Google', value: 'google' }
      ]
    case 'status':
      return [
        { label: '已连接', value: 'active' },
        { label: '未连接', value: 'inactive' },
        { label: '错误', value: 'error' }
      ]
    case 'enabled':
      return [
        { label: '启用', value: true },
        { label: '禁用', value: false }
      ]
    case 'health_status':
      return [
        { label: '健康', value: 'healthy' },
        { label: '不健康', value: 'unhealthy' }
      ]
    case 'domain_id':
      return domains.value.map(d => ({ label: d.name, value: d.id }))
    default:
      return []
  }
}

// 字段切换处理
const handleFieldChange = (field: string) => {
  searchField.value = field
  searchKeyword.value = ''
  searchSelectValue.value = undefined
}

// 搜索提示显示
const showSearchTip = () => {
  showTip.value = true
  setTimeout(() => {
    showTip.value = false
  }, 3000)
}

// 行搜索
const handleSearch = () => {
  currentPage.value = 1
  loadAccounts()
}

// 重置搜索
const handleResetSearch = () => {
  searchField.value = 'name'
  searchKeyword.value = ''
  searchSelectValue.value = undefined
  currentPage.value = 1
  activeCategory.value = 'all'
  loadAccounts()
}

// 查询表单（已废弃，保留兼容）
const queryForm = reactive({
  name: '',
  provider_type: '',
  status: '',
  enabled: undefined as boolean | undefined
})

// 供应商分类选项
const selectedProvider = ref('')
const publicCloudProviders = ref([
  { id: 'alibaba', displayName: '阿里云', icon: Cloudy },
  { id: 'tencent', displayName: '腾讯云', icon: Cloudy },
  { id: 'aws', displayName: 'AWS', icon: Cloudy },
  { id: 'azure', displayName: 'Azure', icon: Cloudy },
  { id: 'huawei', displayName: '华为云', icon: Cloudy },
  { id: 'google', displayName: 'Google', icon: Cloudy },
  { id: 'qcloud', displayName: '天翼云', icon: Cloudy },
  { id: 'ksyun', displayName: '金山云', icon: Cloudy },
  { id: 'volcengine', displayName: '火山引擎', icon: Cloudy },
  { id: 'oracle', displayName: 'OracleCloud', icon: Cloudy }
])
const privateCloudProviders = ref([
  { id: 'vmware', displayName: 'VMware', icon: Cloudy },
  { id: 'openstack', displayName: 'OpenStack', icon: Cloudy },
  { id: 'cloudpods', displayName: 'Cloudpods', icon: Cloudy },
  { id: 'nutanix', displayName: 'Nutanix', icon: Cloudy },
  { id: 'proxmox', displayName: 'Proxmox', icon: Cloudy }
])
const storageProviders = ref([
  { id: 's3', displayName: 'S3', icon: Cloudy },
  { id: 'ceph', displayName: 'Ceph', icon: Cloudy },
  { id: 'xsky', displayName: 'XSKY', icon: Cloudy }
])

const getProviderDisplayName = (providerId: string) => {
  const allProviders = [...publicCloudProviders.value, ...privateCloudProviders.value, ...storageProviders.value]
  const provider = allProviders.find(p => p.id === providerId)
  return provider?.displayName || providerId
}

// 向导表单数据
const wizardForm = reactive({
  name: '',
  remarks: '',
  accountType: 'primary',
  accessKeyId: '',
  accessKeySecret: '',
  resourceAssignmentMethods: ['tag_mapping'] as string[],
  specifyProjectId: null as number | null,
  syncPolicyId: null as number | null,
  syncScope: [] as string[],
  defaultProjectId: null as number | null,
  blockSyncResources: false,
  enableProxy: false,
  proxyId: 'default',
  enablePasswordless: false,
  readOnlyMode: false,
  testConnectionStatus: '',
  testConnectionResult: '',
  selectedRegions: [] as string[]
})

// 资源归属方式说明组合函数
const resourceAssignmentMethodsRef = computed(() => wizardForm.resourceAssignmentMethods)
const {
  description: resourceAssignmentDescription,
  visibleControls: resourceAssignmentControls,
  validationHint: resourceAssignmentValidationHint
} = useResourceAssignmentDescription(resourceAssignmentMethodsRef)

// 定时任务表单
const scheduleForm = reactive({
  name: '',
  frequency: 'daily',
  triggerTime: '02:00'
})

// 可用区域列表
const availableRegions = ref([
  { name: '华北1（青岛）', status: 'available' },
  { name: '华北2（北京）', status: 'available' },
  { name: '华东1（杭州）', status: 'available' },
  { name: '华东2（上海）', status: 'available' },
  { name: '华南1（深圳）', status: 'available' },
  { name: '西南1（成都）', status: 'available' },
  { name: '中国（香港）', status: 'available' },
  { name: '新加坡', status: 'available' }
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

// 同步相关
const showSyncDialog = ref(false)
const showUpdateDialog = ref(false)
const showAccountDetailDialog = ref(false)
const showSetSyncAttributionDialog = ref(false)
const showStatusSettingDialog = ref(false)
const showAttributeSettingDialog = ref(false)
const showSetProxyDialog = ref(false)
const editAccountId = ref<number | null>(null)
const accountDetailId = ref<number | null>(null)
const accountDetailTab = ref('detail')
const syncing = ref(false)
const supportedResourceTypes = ref<ResourceType[]>([])

const syncForm = ref({
  name: '',
  environment: '',
  connectionStatus: false,
  syncMode: 'full',
  syncAll: 'all',
  syncResourceTypes: [] as string[]
})

const currentAccount = ref<any>(null)

// 加载云账户列表
const loadAccounts = async () => {
  loading.value = true
  try {
    const params: CloudAccountSearchParams = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    // 根据当前搜索字段构建参数
    if (searchKeyword.value || searchSelectValue.value) {
      switch (searchField.value) {
        case 'id':
          params.id = searchKeyword.value
          break
        case 'name':
          params.name = searchKeyword.value
          break
        case 'remarks':
          params.remarks = searchKeyword.value
          break
        case 'provider_type':
          params.provider_type = searchSelectValue.value
          break
        case 'status':
          params.status = searchSelectValue.value
          break
        case 'enabled':
          if (searchSelectValue.value !== '') {
            params.enabled = searchSelectValue.value === 'true'
          }
          break
        case 'health_status':
          params.health_status = searchSelectValue.value
          break
        case 'account_number':
          params.account_number = searchKeyword.value
          break
        case 'domain_id':
          if (searchSelectValue.value) {
            params.domain_id = Number(searchSelectValue.value)
          }
          break
      }
    }

    const res = await getCloudAccounts(params)
    accounts.value = res.items || []
    total.value = res.total || 0
  } catch (e) {
    console.error(e)
    ElMessage.error('加载云账户失败')
  } finally {
    loading.value = false
  }
}

const loadSyncPolicies = async () => {
  try {
    const res = await getSyncPolicies()
    syncPolicies.value = res.items || []
  } catch (e) {
    console.error('Failed to load sync policies:', e)
  }
}

const loadProjects = async () => {
  try {
    const res = await getProjects()
    projects.value = res.items || []
  } catch (e) {
    console.error('Failed to load projects:', e)
  }
}

const loadDomains = async () => {
  try {
    const res = await getDomains()
    domains.value = (res.items || []).map(d => ({ id: d.id, name: d.name }))
  } catch (e) {
    console.error('Failed to load domains:', e)
  }
}

const loadResourceTypes = async () => {
  try {
    const res = await getSupportedResourceTypes()
    supportedResourceTypes.value = res.items || []
  } catch (e) {
    console.error('Failed to load resource types:', e)
  }
}

const resetQuery = () => {
  handleResetSearch()
}

const handleSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
  loadAccounts()
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  loadAccounts()
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
    alibaba: 'warning',
    tencent: 'primary',
    aws: 'danger',
    azure: 'success'
  }
  return map[type] || 'info'
}

const getStatusType = (status: string) => {
  const statusMap: Record<string, string> = {
    active: 'success',
    connected: 'success',
    inactive: 'info',
    disconnected: 'danger',
    error: 'danger',
    pending: 'warning',
    checking: 'warning'
  }
  return statusMap[status] || 'info'
}

const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    active: '已连接',
    connected: '已连接',
    inactive: '未连接',
    disconnected: '连接断开',
    error: '连接错误',
    pending: '连接中',
    checking: '检测中'
  }
  return statusMap[status] || status
}

const getHealthStatusType = (status: string) => {
  const statusMap: Record<string, string> = {
    healthy: 'success',
    normal: 'success',
    unhealthy: 'danger',
    abnormal: 'danger',
    warning: 'warning'
  }
  return statusMap[status] || 'info'
}

const getHealthStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    healthy: '正常',
    normal: '正常',
    unhealthy: '异常',
    abnormal: '异常',
    warning: '警告'
  }
  return statusMap[status] || status
}

const getDomainName = (domainId: number) => {
  const domain = domains.value.find(d => d.id === domainId)
  return domain ? domain.name : `域${domainId}`
}

const isSyncing = (accountId: number) => {
  const status = syncStatusMap.value.get(accountId)
  return status?.status === 'running'
}

const getSyncTimeDisplay = (row: CloudAccount) => {
  const status = syncStatusMap.value.get(row.id)
  if (status?.status === 'running') {
    const elapsed = Math.floor((Date.now() - status.startTime.getTime()) / 1000)
    return `已同步 ${elapsed}s`
  }
  if (status?.duration) {
    return `耗时 ${status.duration}s`
  }
  return row.last_sync ? formatDate(row.last_sync) : '-'
}

const formatDate = (dateString: string | undefined) => {
  if (!dateString) return '-'
  return new Date(dateString).toLocaleString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' })
}

const getRegionStatusType = (status: string) => {
  const statusMap: Record<string, string> = {
    available: 'success',
    unavailable: 'danger'
  }
  return statusMap[status] || 'info'
}

const getRegionStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    available: '可用',
    unavailable: '不可用'
  }
  return statusMap[status] || status
}

const selectProvider = (providerId: string) => {
  selectedProvider.value = providerId
  ElMessage.success(`已选择 ${getProviderDisplayName(providerId)}`)
}

const nextStep = async () => {
  if (wizardStep.value === 0) {
    if (!selectedProvider.value) {
      ElMessage.warning('请先选择云平台')
      return
    }
    wizardStep.value = 1
    return
  }
  if (wizardStep.value === 1) {
    if (!wizardFormRef.value) return
    const valid = await wizardFormRef.value.validate().catch(() => false)
    if (!valid) return
  }
  if (wizardStep.value < 3) {
    wizardStep.value++
  }
}

const previousStep = () => {
  if (wizardStep.value > 0) {
    wizardStep.value--
  }
}

const selectAllRegions = () => {
  wizardForm.selectedRegions = availableRegions.value
    .filter(region => region.status === 'available')
    .map(region => region.name)
}

const clearAllRegions = () => {
  wizardForm.selectedRegions = []
}

const handleRegionSelectionChange = (selection: any[]) => {
  wizardForm.selectedRegions = selection.map(item => item.name)
}

const testConnectionInWizard = async () => {
  if (!wizardForm.accessKeyId || !wizardForm.accessKeySecret) {
    ElMessage.warning('请先填写密钥ID和密钥Secret')
    return
  }

  wizardForm.testConnectionStatus = 'testing'
  wizardForm.testConnectionResult = '正在测试连接...'

  try {
    // 使用唯一的临时名称避免与现有账号冲突
    const tempName = `temp-connection-test-${Date.now()}`
    const tempAccountData: CreateCloudAccountRequest = {
      name: tempName,
      provider_type: 'alibaba',
      credentials: {
        access_key_id: wizardForm.accessKeyId,
        access_key_secret: wizardForm.accessKeySecret
      },
      description: 'Temporary account for connection testing',
      remarks: 'Temporary account for connection testing',
      enabled: false, // 禁用状态以便可直接删除清理
      health_status: 'healthy',
      domain_id: 1,
      resource_assignment_method: 'tag_mapping'
    }

    let tempAccountId: number | null = null
    try {
      const tempAccount = await createCloudAccount(tempAccountData)
      tempAccountId = tempAccount.id
      const response = await testConnectionAPI(tempAccount.id)
      wizardForm.testConnectionStatus = response.connected ? 'success' : 'error'
      wizardForm.testConnectionResult = response.message || (response.connected ? '连接成功' : '连接失败')
    } finally {
      // 确保临时账号被清理
      if (tempAccountId) {
        await deleteCloudAccount(tempAccountId).catch(() => {})
      }
    }
  } catch (error: any) {
    wizardForm.testConnectionStatus = 'error'
    wizardForm.testConnectionResult = error.message || '连接测试失败'
  }
}

const submitWizard = async () => {
  if (!selectedProvider.value) {
    ElMessage.warning('请先选择云平台')
    return
  }
  submitting.value = true
  try {
    const formData: CreateCloudAccountRequest = {
      name: wizardForm.name,
      provider_type: selectedProvider.value,
      credentials: {
        access_key_id: wizardForm.accessKeyId,
        access_key_secret: wizardForm.accessKeySecret,
        account_type: wizardForm.accountType,
        proxy_enabled: String(wizardForm.enableProxy),
        proxy_id: wizardForm.proxyId || '',
        passwordless_enabled: String(wizardForm.enablePasswordless),
        read_only_mode: String(wizardForm.readOnlyMode)
      },
      description: wizardForm.remarks,
      remarks: wizardForm.remarks,
      enabled: true,
      health_status: 'healthy',
      account_number: wizardForm.accessKeyId.substring(0, 10) + '...',
      domain_id: 1,
      resource_assignment_method: wizardForm.resourceAssignmentMethods[0] || 'tag_mapping',
      sync_policy_id: wizardForm.syncPolicyId,
      default_project_id: wizardForm.defaultProjectId
    }

    // 创建云账号
    const accountResponse = await createCloudAccount(formData)
    const accountId = accountResponse?.id || accountResponse?.data?.id

    ElMessage.success('云账户添加成功')

    // 如果设置了定时任务，创建定时同步任务
    if (scheduleForm.name && accountId) {
      try {
        // 构建 cron 表达式
        let cronExpression = ''
        const [hour, minute] = scheduleForm.triggerTime.split(':')
        switch (scheduleForm.frequency) {
          case 'hourly':
            cronExpression = `${minute} * * * *`
            break
          case 'daily':
            cronExpression = `${minute} ${hour} * * *`
            break
          case 'weekly':
            cronExpression = `${minute} ${hour} * * 1`
            break
          default:
            cronExpression = `${minute} ${hour} * * *`
        }

        await createScheduledTask({
          name: scheduleForm.name,
          task_type: 'resource_sync',
          target_id: accountId,
          cron_expression: cronExpression,
          enabled: true,
          config: {
            account_id: accountId,
            account_name: wizardForm.name,
            provider: selectedProvider.value,
            sync_mode: 'incremental'
          }
        })
        ElMessage.success('定时同步任务创建成功')
      } catch (scheduleError) {
        console.error('创建定时任务失败:', scheduleError)
        ElMessage.warning('定时同步任务创建失败，可手动创建')
      }
    }

    resetWizard()
    showWizard.value = false
    loadAccounts()
  } catch (error) {
    console.error(error)
    ElMessage.error('添加云账户失败')
  } finally {
    submitting.value = false
  }
}

const resetWizard = () => {
  wizardStep.value = 0
  selectedProvider.value = ''
  Object.assign(wizardForm, {
    name: '',
    remarks: '',
    accountType: 'primary',
    accessKeyId: '',
    accessKeySecret: '',
    resourceAssignmentMethods: ['tag_mapping'],
    specifyProjectId: null,
    syncPolicyId: null,
    syncScope: [],
    defaultProjectId: null,
    blockSyncResources: false,
    enableProxy: false,
    proxyId: 'default',
    enablePasswordless: false,
    readOnlyMode: false,
    testConnectionStatus: '',
    testConnectionResult: '',
    selectedRegions: []
  })
  Object.assign(scheduleForm, {
    name: '',
    frequency: 'daily',
    triggerTime: '02:00'
  })
}

const handleSyncClick = (row: any) => {
  currentAccount.value = row
  syncForm.value.name = row.name
  syncForm.value.environment = getProviderName(row.provider_type)
  syncForm.value.connectionStatus = row.status === 'active'
  showSyncDialog.value = true
}

const handleDropdownCommand = async (command: string, row: any) => {
  currentAccount.value = row

  switch (command) {
    case 'statusSetting':
      editAccountId.value = row.id
      showStatusSettingDialog.value = true
      break
    case 'attributeSetting':
      editAccountId.value = row.id
      showAttributeSettingDialog.value = true
      break
    case 'enable':
      try {
        await updateCloudAccountStatus(row.id, true)
        ElMessage.success('已启用云账号')
        loadAccounts()
      } catch (error) {
        ElMessage.error('启用失败')
      }
      break
    case 'disable':
      try {
        await updateCloudAccountStatus(row.id, false)
        ElMessage.success('已禁用云账号')
        loadAccounts()
      } catch (error) {
        ElMessage.error('禁用失败')
      }
      break
    case 'test':
      try {
        const response = await testConnectionAPI(row.id)
        if (response.connected) {
          ElMessage.success('连接测试成功')
        } else {
          ElMessage.warning('连接测试失败: ' + (response.message || '无法连接'))
        }
      } catch (error) {
        ElMessage.error('连接测试失败')
      }
      break
    case 'update':
      editAccountId.value = row.id
      showUpdateDialog.value = true
      break
    case 'setSyncAttribution':
      editAccountId.value = row.id
      showSetSyncAttributionDialog.value = true
      break
    case 'setSyncPolicy':
      // 打开详情抽屉的定时任务tab
      accountDetailId.value = row.id
      accountDetailTab.value = 'scheduledTasks'
      showAccountDetailDialog.value = true
      break
    case 'setProxy':
      editAccountId.value = row.id
      showSetProxyDialog.value = true
      break
    case 'setPasswordless':
      try {
        const newValue = !row.passwordless_enabled
        await ElMessageBox.confirm(
          `确定要${newValue ? '开启' : '关闭'}免密登录吗？`,
          '确认操作',
          { type: 'warning' }
        )
        // TODO: 调用API更新免密登录状态
        ElMessage.success(`${newValue ? '已开启' : '已关闭'}免密登录`)
        loadAccounts()
      } catch (e: any) {
        if (e !== 'cancel') {
          ElMessage.error('操作失败')
        }
      }
      break
    case 'setReadOnly':
      try {
        const newValue = !row.read_only_mode
        await ElMessageBox.confirm(
          `确定要${newValue ? '开启' : '关闭'}只读模式吗？${newValue ? '开启后将无法对云资源进行修改操作。' : ''}`,
          '确认操作',
          { type: 'warning' }
        )
        // TODO: 调用API更新只读模式状态
        ElMessage.success(`${newValue ? '已开启' : '已关闭'}只读模式`)
        loadAccounts()
      } catch (e: any) {
        if (e !== 'cancel') {
          ElMessage.error('操作失败')
        }
      }
      break
    case 'delete':
      // 检查启用状态
      if (row.enabled) {
        ElMessage.warning('账号为启用状态，请先禁用后再删除')
        return
      }
      try {
        await ElMessageBox.confirm(`确定要删除云账户 "${row.name}" 吗？`, '提示', { type: 'warning' })
        await deleteCloudAccount(row.id)
        ElMessage.success('删除成功')
        loadAccounts()
      } catch (e: any) {
        if (e !== 'cancel') {
          // 显示后端返回的具体错误信息
          ElMessage.error(e.response?.data?.error || e.message || '删除失败')
        }
      }
      break
  }
}

const confirmSync = async () => {
  syncing.value = true
  try {
    const resourceTypes = syncForm.value.syncAll === 'all' ? ['all'] : syncForm.value.syncResourceTypes

    await syncCloudAccount(currentAccount.value.id, {
      mode: syncForm.value.syncMode,
      resource_types: resourceTypes
    })

    ElMessage.success(`已启动对 ${syncForm.value.name} 的同步`)
    showSyncDialog.value = false

    pollSyncStatus(currentAccount.value.id)
  } catch (error) {
    ElMessage.error('同步启动失败')
  } finally {
    syncing.value = false
  }
}

const pollSyncStatus = async (accountId: number) => {
  syncStatusMap.value.set(accountId, {
    status: 'running',
    startTime: new Date()
  })

  const pollInterval = setInterval(async () => {
    try {
      const log = await getLatestSyncLog(accountId)
      if (log.status !== 'running') {
        clearInterval(pollInterval)
        syncStatusMap.value.set(accountId, {
          status: log.status,
          startTime: new Date(log.sync_start_time),
          duration: log.sync_duration
        })
        loadAccounts()
      }
    } catch (e) {
      console.error('Failed to poll sync status:', e)
    }
  }, 3000)

  setTimeout(() => {
    clearInterval(pollInterval)
  }, 5 * 60 * 1000)
}

// 批量选择处理
const handleSelectionChange = (selection: CloudAccount[]) => {
  selectedIds.value = selection.map(item => item.id)
}

// Tab分类切换
const handleCategoryChange = (tab: string) => {
  // 根据分类筛选
  if (tab === 'public') {
    queryForm.provider_type = 'alibaba'
  } else {
    queryForm.provider_type = ''
  }
  loadAccounts()
}

// 批量操作
const handleBatchOperation = async () => {
  if (selectedIds.value.length === 0) {
    ElMessage.warning('请先选择要操作的云账号')
    return
  }
  try {
    const action = await ElMessageBox.confirm(
      `已选择 ${selectedIds.value.length} 个云账号，请确认批量同步操作`,
      '批量操作',
      {
        confirmButtonText: '批量同步',
        cancelButtonText: '取消',
        type: 'info'
      }
    )
    // 执行批量同步
    for (const id of selectedIds.value) {
      await syncCloudAccount(id, { mode: 'full', resource_types: ['all'] })
    }
    ElMessage.success('批量同步已启动')
    selectedIds.value = []
    loadAccounts()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error('批量操作失败')
    }
  }
}

// 导出
const handleExport = async () => {
  try {
    // 导出CSV
    const csvContent = [
      ['ID', '名称', '平台', '状态', '启用状态', '余额', '账号', '所属域'].join(','),
      ...accounts.value.map(row => [
        row.id,
        row.name,
        getProviderName(row.provider_type),
        getStatusText(row.status),
        row.enabled ? '启用' : '禁用',
        row.balance?.toFixed(2) || '0',
        row.account_number || '',
        getDomainName(row.domain_id)
      ].join(','))
    ].join('\n')

    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
    const link = document.createElement('a')
    link.href = URL.createObjectURL(blob)
    link.download = `cloud_accounts_${new Date().toISOString().slice(0,10)}.csv`
    link.click()
    ElMessage.success('导出成功')
  } catch (error) {
    ElMessage.error('导出失败')
  }
}

// 设置
const handleSettings = () => {
  ElMessage.info('全局设置功能开发中')
}

// 打开详情抽屉
const openDetailDrawer = (row: CloudAccount) => {
  accountDetailId.value = row.id
  accountDetailTab.value = 'detail'
  showAccountDetailDialog.value = true
}

// 获取资源归属方式文本
const getResourceAssignmentText = (method: string) => {
  const map: Record<string, string> = {
    'tag_mapping': '根据同步策略归属',
    'project_mapping': '根据云上项目归属',
    'subscription_mapping': '根据云订阅归属',
    'specify_project': '指定项目',
    'manual_assignment': '手动分配'
  }
  return map[method] || method || '-'
}

onMounted(() => {
  loadAccounts()
  loadSyncPolicies()
  loadProjects()
  loadDomains()
  loadResourceTypes()
})
</script>

<style scoped>
.cloud-accounts-container {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-header h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
}

.tool-bar {
  display: flex;
  gap: 8px;
}

.top-tabs {
  margin-bottom: 16px;
}

.filter-card {
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  text-align: right;
}

.wizard-step-content {
  margin: 30px 0;
  min-height: 400px;
}

.provider-category {
  margin-bottom: 24px;
}

.provider-category h4 {
  margin: 0 0 16px 0;
  font-size: 14px;
  font-weight: 600;
  color: var(--el-text-color-secondary);
}

.provider-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
  gap: 16px;
}

.provider-card {
  cursor: pointer;
  transition: all 0.2s ease;
}

.provider-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.provider-card.selected {
  border-color: var(--el-color-primary);
  background-color: var(--el-color-primary-light-9);
}

.provider-item {
  text-align: center;
  padding: 16px 10px;
}

.provider-item h4 {
  margin-top: 8px;
  font-size: 13px;
}

.selected-indicator {
  margin-top: 20px;
  text-align: center;
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

.wizard-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.text-success {
  color: #67c23a;
}

.text-danger {
  color: #f56c6c;
}

.search-bar {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
  padding: 12px 16px;
  background: var(--el-bg-color);
  border-radius: 8px;
  border: 1px solid var(--el-border-color-light);
}

.field-selector {
  width: 120px;
}

.search-input-area {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
}

.search-input {
  max-width: 300px;
}

.search-select {
  width: 200px;
}

.search-tip {
  margin-left: 8px;
  color: var(--el-text-color-secondary);
  font-size: 12px;
}

.search-actions {
  display: flex;
  gap: 8px;
}

.assignment-description {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  margin-top: 8px;
  padding: 12px 16px;
  background-color: var(--el-color-primary-light-9);
  border-radius: 8px;
  border-left: 4px solid var(--el-color-primary);
  color: var(--el-text-color-regular);
  font-size: 13px;
  line-height: 1.6;
}

.assignment-description .el-icon {
  color: var(--el-color-primary);
  margin-top: 2px;
}

.assignment-hint {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-top: 6px;
  color: var(--el-color-warning);
  font-size: 12px;
}

.assignment-hint .el-icon {
  margin-top: 1px;
}

.form-hint {
  display: block;
  margin-top: 4px;
  color: var(--el-text-color-secondary);
  font-size: 12px;
}
</style>