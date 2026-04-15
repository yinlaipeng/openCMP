<template>
  <div class="detail-tab" v-loading="loading">
    <!-- 基本信息 -->
    <div class="section">
      <div class="section-title">基本信息</div>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="名称">{{ account.name }}</el-descriptions-item>
        <el-descriptions-item label="平台">
          <el-tag :type="getProviderType(account.provider_type)">
            {{ getProviderName(account.provider_type) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(account.status)">
            {{ getStatusText(account.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="启用状态">
          <el-tag :type="account.enabled ? 'success' : 'info'">
            {{ account.enabled ? '启用' : '禁用' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="健康状态">
          <el-tag :type="getHealthStatusType(account.health_status)">
            {{ getHealthStatusText(account.health_status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ formatTime(account.created_at) }}</el-descriptions-item>
        <el-descriptions-item label="上次同步">{{ account.last_sync ? formatTime(account.last_sync) : '未同步' }}</el-descriptions-item>
        <el-descriptions-item label="备注">{{ account.description || '无' }}</el-descriptions-item>
      </el-descriptions>
    </div>

    <!-- 账号信息 -->
    <div class="section">
      <div class="section-title">账号信息</div>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="账号ID">{{ account.account_number || '未知' }}</el-descriptions-item>
        <el-descriptions-item label="账户余额">
          <span class="balance">¥ {{ (account.balance || 0).toFixed(2) }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="默认区域">{{ getRegionFromCredentials() }}</el-descriptions-item>
        <el-descriptions-item label="所属域">{{ getDomainName(account.domain_id) }}</el-descriptions-item>
        <el-descriptions-item label="资源归属方式">{{ getResourceAssignmentMethodName(account.resource_assignment_method) }}</el-descriptions-item>
      </el-descriptions>
    </div>

    <!-- 云平台权限 -->
    <div class="section">
      <div class="section-title">
        云平台权限
        <el-button size="small" @click="refreshPermissions" :loading="permissionsLoading" style="float: right">
          刷新权限
        </el-button>
      </div>
      <el-table :data="permissions" style="width: 100%" v-if="permissions.length > 0">
        <el-table-column prop="name" label="权限名称" width="200" />
        <el-table-column prop="description" label="权限描述" />
        <el-table-column prop="status" label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="row.granted ? 'success' : 'info'">
              {{ row.granted ? '已授权' : '未授权' }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
      <el-empty v-else description="暂无权限信息，请点击刷新权限获取" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import type { CloudAccount } from '@/types'
import { getCloudAccount } from '@/api/cloud-account'

interface Props {
  account: CloudAccount
  loading?: boolean
}

interface Emits {
  (e: 'refresh'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const permissionsLoading = ref(false)
const permissions = ref<any[]>([])

// 模拟权限数据（后续会从云厂商API获取）
const mockPermissions = [
  { name: 'AliyunECSFullAccess', description: '云服务器ECS完全管理权限', granted: true },
  { name: 'AliyunVPCFullAccess', description: '专有网络VPC完全管理权限', granted: true },
  { name: 'AliyunSLBFullAccess', description: '负载均衡SLB完全管理权限', granted: true },
  { name: 'AliyunRDSFullAccess', description: '云数据库RDS完全管理权限', granted: false },
  { name: 'AliyunOSSFullAccess', description: '对象存储OSS完全管理权限', granted: true },
]

onMounted(() => {
  // 暂时使用模拟数据
  permissions.value = mockPermissions
})

async function refreshPermissions() {
  permissionsLoading.value = true
  try {
    // TODO: 调用真实API获取权限
    await new Promise(resolve => setTimeout(resolve, 1000))
    permissions.value = mockPermissions
    ElMessage.success('权限信息已刷新')
  } catch (error) {
    ElMessage.error('刷新权限失败')
  } finally {
    permissionsLoading.value = false
  }
}

function getRegionFromCredentials(): string {
  try {
    const creds = JSON.parse(JSON.stringify(props.account.credentials))
    return creds.region_id || '未知'
  } catch {
    return '未知'
  }
}

function formatTime(time: string | Date | null): string {
  if (!time) return '未知'
  const date = new Date(time)
  return date.toLocaleString('zh-CN')
}

function getProviderType(type: string): string {
  const types: Record<string, string> = {
    alibaba: '',
    tencent: 'warning',
    aws: 'success',
    azure: 'info'
  }
  return types[type] || 'info'
}

function getProviderName(type: string): string {
  const names: Record<string, string> = {
    alibaba: '阿里云',
    tencent: '腾讯云',
    aws: 'AWS',
    azure: 'Azure'
  }
  return names[type] || type
}

function getStatusType(status: string): string {
  const types: Record<string, string> = {
    active: 'success',
    inactive: 'info',
    error: 'danger'
  }
  return types[status] || 'info'
}

function getStatusText(status: string): string {
  const texts: Record<string, string> = {
    active: '已连接',
    inactive: '未连接',
    error: '错误'
  }
  return texts[status] || status
}

function getHealthStatusType(status: string): string {
  const types: Record<string, string> = {
    healthy: 'success',
    unhealthy: 'danger',
    no_permission: 'warning'
  }
  return types[status] || 'info'
}

function getHealthStatusText(status: string): string {
  const texts: Record<string, string> = {
    healthy: '正常',
    unhealthy: '异常',
    no_permission: '无权限'
  }
  return texts[status] || status
}

function getDomainName(domainId: number): string {
  // TODO: 从域列表获取名称
  return domainId === 1 ? '默认域' : `域 ${domainId}`
}

function getResourceAssignmentMethodName(method: string): string {
  const names: Record<string, string> = {
    tag_mapping: '根据标签映射',
    project_mapping: '根据项目映射',
    manual_assignment: '手动分配'
  }
  return names[method] || method
}
</script>

<style scoped>
.detail-tab {
  padding: 10px;
}

.section {
  margin-bottom: 24px;
}

.section-title {
  font-weight: 600;
  font-size: 15px;
  color: #303133;
  margin-bottom: 16px;
  padding-bottom: 8px;
  border-bottom: 1px solid #e4e7ed;
}

.balance {
  font-size: 16px;
  font-weight: 600;
  color: #409eff;
}
</style>