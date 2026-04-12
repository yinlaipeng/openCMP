<template>
  <div class="subnets-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">IP子网列表</span>
        </div>
      </template>

      <el-table :data="subnets" v-loading="loading">
        <el-table-column label="名称" width="180">
          <template #default="{ row }">
            <el-link @click="viewDetail(row)">{{ row.name }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="type" label="类型" width="100" />
        <el-table-column prop="auto_schedule" label="自动调度" width="100">
          <template #default="{ row }">
            <el-tag :type="row.auto_schedule ? 'success' : 'info'">
              {{ row.auto_schedule ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="ip_address" label="ip地址" width="150" />
        <el-table-column prop="ipv6_address" label="ipv6地址" width="150" />
        <el-table-column prop="usage" label="使用情况" width="100" />
        <el-table-column prop="schedule_tags" label="调度标签" width="150" />
        <el-table-column prop="platform" label="平台" width="100" />
        <el-table-column prop="project" label="项目" width="120" />
        <el-table-column prop="region_id" label="区域" width="100" />
        <el-table-column label="操作" width="150">
          <template #default="{ row }">
            <el-button size="small" @click="adjustScheduleTags(row)">调整调度标签</el-button>
            <el-dropdown>
              <el-button size="small">
                更多 <el-icon class="el-icon--right"><arrow-down /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="modifyAttributes(row)">修改属性</el-dropdown-item>
                  <el-dropdown-item @click="changeProject(row)">更改项目</el-dropdown-item>
                  <el-dropdown-item @click="splitSubnet(row)">分割ip子网</el-dropdown-item>
                  <el-dropdown-item @click="reserveIP(row)">预留ip</el-dropdown-item>
                  <el-dropdown-item @click="toggleAutoSchedule(row)">设置自动调度</el-dropdown-item>
                  <el-dropdown-item @click="switchLayer2Network(row)">更换二层网络</el-dropdown-item>
                  <el-dropdown-item @click="syncStatus(row)">同步状态</el-dropdown-item>
                  <el-dropdown-item divided @click="handleDelete(row)">删除</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
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
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowDown } from '@element-plus/icons-vue'

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

onMounted(() => {
  loadSubnets()
})
</script>

<style scoped>
.subnets-page {
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

.ip-usage-chart {
  padding: 20px;
}
</style>