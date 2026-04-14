<template>
  <div class="subnets-page">
    <el-card class="page-card glass">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <span class="title">IP子网列表</span>
            <el-tag size="small" type="info" class="count-tag">共 {{ subnets.length }} 个</el-tag>
          </div>
          <el-button type="primary" class="create-btn" @click="handleCreate">
            <el-icon><Plus /></el-icon>
            创建子网
          </el-button>
        </div>
      </template>

      <!-- Empty State -->
      <el-empty v-if="!loading && subnets.length === 0" description="暂无子网数据">
        <el-button type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>
          创建第一个子网
        </el-button>
      </el-empty>

      <el-table v-if="subnets.length > 0 || loading" :data="subnets" v-loading="loading" class="subnet-table">
        <el-table-column label="名称" width="180">
          <template #default="{ row }">
            <el-link type="primary" @click="viewDetail(row)">{{ row.name }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" class="status-tag">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag size="small" effect="plain">{{ row.type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="auto_schedule" label="自动调度" width="100">
          <template #default="{ row }">
            <el-tag :type="row.auto_schedule ? 'success' : 'info'" size="small" effect="plain">
              {{ row.auto_schedule ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="ip_address" label="IP地址" width="150">
          <template #default="{ row }">
            <span class="ip-address font-mono">{{ row.ip_address }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="ipv6_address" label="IPv6地址" width="150">
          <template #default="{ row }">
            <span v-if="row.ipv6_address" class="ip-address font-mono">{{ row.ipv6_address }}</span>
            <span v-else class="no-data">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="usage" label="使用情况" width="100">
          <template #default="{ row }">
            <el-progress
              :percentage="parseInt(row.usage) || 0"
              :stroke-width="6"
              :show-text="false"
              :color="row.usage > '80%' ? 'var(--color-warning)' : 'var(--color-accent)'"
            />
            <span class="usage-text">{{ row.usage }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="platform" label="平台" width="100">
          <template #default="{ row }">
            <el-tag size="small" effect="plain">{{ row.platform }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="region_id" label="区域" width="100">
          <template #default="{ row }">
            <span class="region font-mono">{{ row.region_id }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <div class="operation-buttons">
              <el-dropdown trigger="click">
                <el-button size="small" type="primary">
                  <el-icon><Setting /></el-icon>
                  操作
                  <el-icon class="el-icon--right"><ArrowDown /></el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item @click="adjustScheduleTags(row)">
                      <el-icon><PriceTag /></el-icon>调整调度标签
                    </el-dropdown-item>
                    <el-dropdown-item @click="modifyAttributes(row)">修改属性</el-dropdown-item>
                    <el-dropdown-item @click="changeProject(row)">更改项目</el-dropdown-item>
                    <el-dropdown-item @click="syncStatus(row)">
                      <el-icon><Refresh /></el-icon>同步状态
                    </el-dropdown-item>
                    <el-dropdown-item divided @click="handleDelete(row)">
                      <el-icon color="var(--color-danger)"><Delete /></el-icon>删除
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- Detail Modal -->
    <el-dialog v-model="detailDialogVisible" title="IP子网详情" width="800px">
      <el-tabs v-model="detailTab">
        <el-tab-pane label="详情" name="detail">
          <el-descriptions :column="2" border>
            <el-descriptions-item label="ID">{{ selectedSubnet?.id }}</el-descriptions-item>
            <el-descriptions-item label="名称">{{ selectedSubnet?.name }}</el-descriptions-item>
            <el-descriptions-item label="状态">{{ selectedSubnet?.status }}</el-descriptions-item>
            <el-descriptions-item label="类型">{{ selectedSubnet?.type }}</el-descriptions-item>
            <el-descriptions-item label="自动调度">{{ selectedSubnet?.auto_schedule ? '是' : '否' }}</el-descriptions-item>
            <el-descriptions-item label="IP地址">{{ selectedSubnet?.ip_address }}</el-descriptions-item>
            <el-descriptions-item label="IPv6地址">{{ selectedSubnet?.ipv6_address }}</el-descriptions-item>
            <el-descriptions-item label="使用情况">{{ selectedSubnet?.usage }}</el-descriptions-item>
            <el-descriptions-item label="调度标签">{{ selectedSubnet?.schedule_tags }}</el-descriptions-item>
            <el-descriptions-item label="平台">{{ selectedSubnet?.platform }}</el-descriptions-item>
            <el-descriptions-item label="项目">{{ selectedSubnet?.project }}</el-descriptions-item>
            <el-descriptions-item label="区域">{{ selectedSubnet?.region_id }}</el-descriptions-item>
            <el-descriptions-item label="创建时间">{{ selectedSubnet?.created_at }}</el-descriptions-item>
          </el-descriptions>
        </el-tab-pane>
        <el-tab-pane label="ip使用情况" name="ip_usage">
          <div class="ip-usage-chart">
            <el-row :gutter="20">
              <el-col :span="8">
                <el-statistic title="总IP数" :value="256" />
              </el-col>
              <el-col :span="8">
                <el-statistic title="已使用" :value="128" />
              </el-col>
              <el-col :span="8">
                <el-statistic title="使用率" :value="50" suffix="%" />
              </el-col>
            </el-row>
          </div>
        </el-tab-pane>
        <el-tab-pane label="预留ip" name="reserved_ips">
          <el-table :data="reservedIPs">
            <el-table-column prop="ip" label="IP地址" width="150" />
            <el-table-column prop="type" label="类型" width="100" />
            <el-table-column prop="usage" label="用途" width="200" />
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="mac映射表" name="mac_mappings">
          <el-table :data="macMappings">
            <el-table-column prop="mac" label="MAC地址" width="200" />
            <el-table-column prop="ip" label="IP地址" width="150" />
            <el-table-column prop="device" label="设备" width="150" />
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="操作日志" name="logs">
          <el-table :data="logs">
            <el-table-column prop="operation" label="操作" width="150" />
            <el-table-column prop="operator" label="操作员" width="150" />
            <el-table-column prop="result" label="结果" width="100" />
            <el-table-column prop="timestamp" label="时间" width="180" />
          </el-table>
        </el-tab-pane>
      </el-tabs>
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- Adjust Schedule Tags Dialog -->
    <el-dialog v-model="adjustTagsDialogVisible" title="调整调度标签" width="400px">
      <el-form label-width="100px">
        <el-form-item label="调度标签">
          <el-select v-model="scheduleTags" multiple placeholder="请选择调度标签">
            <el-option label="tag1" value="tag1" />
            <el-option label="tag2" value="tag2" />
            <el-option label="tag3" value="tag3" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="adjustTagsDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitAdjustTags">确定</el-button>
      </template>
    </el-dialog>

    <!-- 创建子网模态框 -->
    <CreateSubnetModal
      v-model:visible="createModalVisible"
      @success="handleCreateSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowDown, Plus, Setting, PriceTag, Refresh, Delete } from '@element-plus/icons-vue'
import CreateSubnetModal from '@/components/network/CreateSubnetModal.vue'

interface ExtendedSubnet {
  id: string
  name: string
  status: string
  type: string
  auto_schedule: boolean
  ip_address: string
  ipv6_address: string
  usage: string
  schedule_tags: string
  platform: string
  project: string
  region_id: string
  created_at?: string
}

const subnets = ref<ExtendedSubnet[]>([])
const loading = ref(false)
const detailDialogVisible = ref(false)
const adjustTagsDialogVisible = ref(false)
const createModalVisible = ref(false)
const selectedSubnet = ref<ExtendedSubnet | null>(null)
const detailTab = ref('detail')
const scheduleTags = ref<string[]>([])

// Detail modal data
const reservedIPs = ref<any[]>([])
const macMappings = ref<any[]>([])
const logs = ref<any[]>([])

const getStatusType = (status: string) => {
  switch (status.toLowerCase()) {
    case 'available':
    case 'active':
      return 'success'
    case 'pending':
      return 'warning'
    case 'failed':
      return 'danger'
    default:
      return 'info'
  }
}

const loadSubnets = async () => {
  loading.value = true
  try {
    // Mock data - API requires valid cloud account credentials
    subnets.value = [
      {
        id: 'subnet-1',
        name: 'Subnet 1',
        status: 'Available',
        type: 'private',
        auto_schedule: true,
        ip_address: '10.0.1.0/24',
        ipv6_address: '2001:db8::/64',
        usage: '50%',
        schedule_tags: 'tag1,tag2',
        platform: '阿里云',
        project: 'Project A',
        region_id: 'cn-hangzhou',
        created_at: '2024-01-01 10:00:00'
      },
      {
        id: 'subnet-2',
        name: 'Subnet 2',
        status: 'Available',
        type: 'private',
        auto_schedule: true,
        ip_address: '10.0.2.0/24',
        ipv6_address: '2001:db8::/64',
        usage: '30%',
        schedule_tags: 'tag1',
        platform: '阿里云',
        project: 'Project B',
        region_id: 'cn-hangzhou',
        created_at: '2024-01-01 10:00:00'
      },
      {
        id: 'subnet-3',
        name: 'Subnet 3',
        status: 'Available',
        type: 'public',
        auto_schedule: false,
        ip_address: '10.0.3.0/24',
        ipv6_address: '2001:db8::/64',
        usage: '20%',
        schedule_tags: 'tag2,tag3',
        platform: '阿里云',
        project: 'Project A',
        region_id: 'cn-shanghai',
        created_at: '2024-01-01 10:00:00'
      }
    ]
  } catch (e) {
    console.error(e)
    subnets.value = []
  } finally {
    loading.value = false
  }
}

const viewDetail = (row: ExtendedSubnet) => {
  selectedSubnet.value = row
  detailDialogVisible.value = true
  detailTab.value = 'detail'

  // Load detail data (mock)
  reservedIPs.value = [
    { ip: '10.0.1.1', type: 'gateway', usage: 'Gateway' },
    { ip: '10.0.1.2', type: 'dns', usage: 'DNS Server' }
  ]
  macMappings.value = [
    { mac: 'AA:BB:CC:DD:EE:FF', ip: '10.0.1.10', device: 'VM-01' },
    { mac: '11:22:33:44:55:66', ip: '10.0.1.11', device: 'VM-02' }
  ]
  logs.value = [
    { operation: '创建', operator: 'admin', result: '成功', timestamp: '2024-01-01 10:00:00' }
  ]
}

const adjustScheduleTags = (row: ExtendedSubnet) => {
  selectedSubnet.value = row
  scheduleTags.value = row.schedule_tags.split(',')
  adjustTagsDialogVisible.value = true
}

const submitAdjustTags = () => {
  ElMessage.success('调整调度标签成功')
  adjustTagsDialogVisible.value = false
  loadSubnets()
}

const modifyAttributes = (row: ExtendedSubnet) => {
  ElMessage.info(`修改属性功能开发中: ${row.name}`)
}

const changeProject = (row: ExtendedSubnet) => {
  ElMessage.info(`更改项目功能开发中: ${row.name}`)
}

const splitSubnet = (row: ExtendedSubnet) => {
  ElMessage.info(`分割IP子网功能开发中: ${row.name}`)
}

const reserveIP = (row: ExtendedSubnet) => {
  ElMessage.info(`预留IP功能开发中: ${row.name}`)
}

const toggleAutoSchedule = (row: ExtendedSubnet) => {
  ElMessage.info(`设置自动调度功能开发中: ${row.name}`)
}

const switchLayer2Network = (row: ExtendedSubnet) => {
  ElMessage.info(`更换二层网络功能开发中: ${row.name}`)
}

const syncStatus = (row: ExtendedSubnet) => {
  ElMessage.info(`正在同步子网 ${row.name} 的状态`)
}

const handleDelete = async (row: ExtendedSubnet) => {
  try {
    await ElMessageBox.confirm(`确认删除IP子网 ${row.name}？`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    // Mock deletion - API requires valid cloud account credentials
    subnets.value = subnets.value.filter(s => s.id !== row.id)
    ElMessage.success('删除成功')
  } catch (e) {
    console.error(e)
  }
}

const handleCreate = () => {
  createModalVisible.value = true
}

const handleCreateSuccess = (subnet: ExtendedSubnet) => {
  ElMessage.success(`${subnet.name} 创建成功`)
  loadSubnets()
}

onMounted(() => {
  loadSubnets()
})
</script>

<style scoped>
.subnets-page {
  height: 100%;
  padding: var(--space-4);
}

.page-card {
  height: 100%;
}

.page-card.glass {
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(16px);
  border-radius: var(--radius-lg);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.title {
  font-size: var(--font-size-xl);
  font-weight: var(--font-weight-semibold);
  color: var(--color-foreground);
}

.count-tag {
  font-family: var(--font-mono);
}

.create-btn {
  transition: all var(--duration-fast) var(--ease-out);
}

.create-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(34, 197, 94, 0.2);
}

/* Table */
.subnet-table {
  border-radius: var(--radius-md);
}

.status-tag {
  font-weight: var(--font-weight-medium);
}

.ip-address {
  font-size: var(--font-size-sm);
}

.no-data {
  color: var(--color-muted);
}

.usage-text {
  font-size: var(--font-size-xs);
  color: var(--color-muted);
  margin-top: var(--space-1);
}

.region {
  font-size: var(--font-size-sm);
}

.operation-buttons {
  display: flex;
  gap: var(--space-2);
  align-items: center;
}

/* Detail Modal */
.ip-usage-chart {
  padding: var(--space-6);
}

/* Responsive */
@media (max-width: 768px) {
  .subnets-page {
    padding: var(--space-2);
  }

  .card-header {
    flex-direction: column;
    gap: var(--space-2);
    align-items: flex-start;
  }

  .create-btn {
    width: 100%;
  }
}
</style>