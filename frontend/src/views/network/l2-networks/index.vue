<template>
  <div class="l2-networks-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">二层网络列表</span>
        </div>
      </template>

      <el-table :data="l2Networks" v-loading="loading">
        <el-table-column label="名称" width="180">
          <template #default="{ row }">
            <el-link @click="viewDetail(row)">{{ row.name }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="bandwidth" label="带宽" width="100">
          <template #default="{ row }">
            {{ row.bandwidth }} Mbps
          </template>
        </el-table-column>
        <el-table-column prop="vpc_id" label="vpc" width="180" />
        <el-table-column prop="network_count" label="网络数量" width="100" />
        <el-table-column prop="platform" label="平台" width="120" />
        <el-table-column prop="domain" label="所属域" width="150" />
        <el-table-column prop="region_id" label="区域" width="150" />
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-button size="small" type="primary" @click="modifyAttributes(row)">修改属性</el-button>
            <el-dropdown>
              <el-button size="small">
                更多 <el-icon class="el-icon--right"><arrow-down /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="changeDomain(row)">更改域</el-dropdown-item>
                  <el-dropdown-item @click="handleDelete(row)">删除</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- Detail Modal -->
    <el-dialog v-model="detailDialogVisible" title="二层网络详情" width="800px">
      <el-tabs v-model="activeTab">
        <el-tab-pane label="详情" name="detail">
          <el-descriptions :column="2" border>
            <el-descriptions-item label="ID">{{ selectedL2Network?.id }}</el-descriptions-item>
            <el-descriptions-item label="名称">{{ selectedL2Network?.name }}</el-descriptions-item>
            <el-descriptions-item label="状态">{{ selectedL2Network?.status }}</el-descriptions-item>
            <el-descriptions-item label="VLAN ID">{{ selectedL2Network?.vlan_id }}</el-descriptions-item>
            <el-descriptions-item label="带宽">{{ selectedL2Network?.bandwidth }} Mbps</el-descriptions-item>
            <el-descriptions-item label="VPC">{{ selectedL2Network?.vpc_id }}</el-descriptions-item>
            <el-descriptions-item label="网络数量">{{ selectedL2Network?.network_count }}</el-descriptions-item>
            <el-descriptions-item label="MTU">{{ selectedL2Network?.mtu || 1500 }}</el-descriptions-item>
            <el-descriptions-item label="创建时间">{{ selectedL2Network?.created_at }}</el-descriptions-item>
          </el-descriptions>
        </el-tab-pane>
        <el-tab-pane label="ip子网" name="subnets">
          <el-table :data="subnets">
            <el-table-column prop="id" label="子网ID" width="200" />
            <el-table-column prop="name" label="子网名称" width="150" />
            <el-table-column prop="status" label="状态" width="100" />
            <el-table-column prop="cidr" label="CIDR" width="180" />
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="虚拟机" name="vms">
          <el-table :data="vms">
            <el-table-column prop="id" label="虚拟机ID" width="200" />
            <el-table-column prop="name" label="虚拟机名称" width="150" />
            <el-table-column prop="status" label="状态" width="100" />
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="物理机" name="physical_hosts">
          <el-table :data="physicalHosts">
            <el-table-column prop="id" label="物理机ID" width="200" />
            <el-table-column prop="name" label="物理机名称" width="150" />
            <el-table-column prop="status" label="状态" width="100" />
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="宿主机" name="hosts">
          <el-table :data="hosts">
            <el-table-column prop="id" label="宿主机ID" width="200" />
            <el-table-column prop="name" label="宿主机名称" width="150" />
            <el-table-column prop="status" label="状态" width="100" />
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="资源统计" name="stats">
          <el-row :gutter="20">
            <el-col :span="6">
              <el-statistic title="总虚拟机数" :value="vms.length" />
            </el-col>
            <el-col :span="6">
              <el-statistic title="总IP子网数" :value="subnets.length" />
            </el-col>
            <el-col :span="6">
              <el-statistic title="总物理机数" :value="physicalHosts.length" />
            </el-col>
            <el-col :span="6">
              <el-statistic title="总宿主机数" :value="hosts.length" />
            </el-col>
          </el-row>
        </el-tab-pane>
        <el-tab-pane label="拓扑" name="topology">
          <div class="topology-placeholder">拓扑图待实现</div>
        </el-tab-pane>
        <el-tab-pane label="监控" name="monitoring">
          <div class="monitoring-placeholder">监控图表待实现</div>
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

    <!-- Modify Attributes Dialog -->
    <el-dialog v-model="modifyDialogVisible" title="修改属性" width="400px">
      <el-form :model="modifyForm" label-width="80px">
        <el-form-item label="名称">
          <el-input v-model="modifyForm.name" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="modifyForm.status">
            <el-option label="Active" value="active" />
            <el-option label="Inactive" value="inactive" />
          </el-select>
        </el-form-item>
        <el-form-item label="带宽">
          <el-input-number v-model="modifyForm.bandwidth" :min="1" :max="10000" />
        </el-form-item>
        <el-form-item label="MTU">
          <el-input-number v-model="modifyForm.mtu" :min="64" :max="9000" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="modifyDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitModify">确定</el-button>
      </template>
    </el-dialog>

    <!-- Change Domain Dialog -->
    <el-dialog v-model="changeDomainDialogVisible" title="更改域" width="400px">
      <el-form label-width="80px">
        <el-form-item label="目标域">
          <el-select v-model="targetDomain" placeholder="请选择域">
            <el-option label="Default Domain" value="1" />
            <el-option label="Domain A" value="2" />
            <el-option label="Domain B" value="3" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="changeDomainDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitChangeDomain">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowDown } from '@element-plus/icons-vue'

interface L2Network {
  id: string
  name: string
  status: string
  bandwidth: number
  vlan_id?: string
  vpc_id: string
  network_count?: number
  platform?: string
  domain?: string
  region_id?: string
  mtu?: number
  created_at?: string
}

const l2Networks = ref<L2Network[]>([])
const loading = ref(false)
const detailDialogVisible = ref(false)
const modifyDialogVisible = ref(false)
const changeDomainDialogVisible = ref(false)
const selectedL2Network = ref<L2Network | null>(null)
const activeTab = ref('detail')
const targetDomain = ref('')

// Modify form
const modifyForm = ref({
  name: '',
  status: 'active',
  bandwidth: 1000,
  mtu: 1500
})

// Detail modal data
const subnets = ref<any[]>([])
const vms = ref<any[]>([])
const physicalHosts = ref<any[]>([])
const hosts = ref<any[]>([])
const logs = ref<any[]>([])

const getStatusType = (status: string) => {
  switch (status.toLowerCase()) {
    case 'active':
    case 'available':
      return 'success'
    case 'pending':
      return 'warning'
    case 'inactive':
    case 'failed':
      return 'danger'
    default:
      return 'info'
  }
}

const loadL2Networks = async () => {
  loading.value = true
  try {
    // Mock data - API requires valid cloud account credentials
    l2Networks.value = [
      {
        id: 'l2-1',
        name: '二层网络 1',
        status: 'Active',
        bandwidth: 1000,
        vlan_id: '100',
        vpc_id: 'vpc-1',
        network_count: 3,
        platform: '阿里云',
        domain: 'Default Domain',
        region_id: 'cn-hangzhou',
        mtu: 1500,
        created_at: '2024-01-01 10:00:00'
      },
      {
        id: 'l2-2',
        name: '二层网络 2',
        status: 'Active',
        bandwidth: 2000,
        vlan_id: '200',
        vpc_id: 'vpc-2',
        network_count: 2,
        platform: '阿里云',
        domain: 'Default Domain',
        region_id: 'cn-hangzhou',
        mtu: 1500,
        created_at: '2024-01-01 10:00:00'
      },
      {
        id: 'l2-3',
        name: '二层网络 3',
        status: 'Pending',
        bandwidth: 500,
        vlan_id: '300',
        vpc_id: 'vpc-1',
        network_count: 1,
        platform: '阿里云',
        domain: 'Default Domain',
        region_id: 'cn-shanghai',
        mtu: 9000,
        created_at: '2024-01-01 10:00:00'
      }
    ]
  } catch (e) {
    console.error(e)
    l2Networks.value = []
  } finally {
    loading.value = false
  }
}

const viewDetail = async (row: L2Network) => {
  selectedL2Network.value = row
  detailDialogVisible.value = true
  activeTab.value = 'detail'

  // Mock data for subnets (API requires valid cloud account)
  subnets.value = [
    { id: 'subnet-1', name: 'Subnet 1', status: 'Available', cidr: '10.0.1.0/24' },
    { id: 'subnet-2', name: 'Subnet 2', status: 'Available', cidr: '10.0.2.0/24' }
  ]

  // Mock data for other tabs
  vms.value = [{ id: 'vm-1', name: 'VM 1', status: 'Running' }]
  physicalHosts.value = [{ id: 'ph-1', name: 'Physical Host 1', status: 'Running' }]
  hosts.value = [{ id: 'host-1', name: 'Host 1', status: 'Running' }]
  logs.value = [{ operation: '创建', operator: 'admin', result: '成功', timestamp: '2024-01-01 10:00:00' }]
}

const modifyAttributes = (row: L2Network) => {
  selectedL2Network.value = row
  modifyForm.value = {
    name: row.name,
    status: row.status.toLowerCase(),
    bandwidth: row.bandwidth,
    mtu: 1500
  }
  modifyDialogVisible.value = true
}

const submitModify = () => {
  ElMessage.success('修改成功')
  modifyDialogVisible.value = false
  loadL2Networks()
}

const changeDomain = (row: L2Network) => {
  selectedL2Network.value = row
  changeDomainDialogVisible.value = true
}

const submitChangeDomain = () => {
  ElMessage.success('更改域成功')
  changeDomainDialogVisible.value = false
  loadL2Networks()
}

const handleDelete = async (row: L2Network) => {
  try {
    await ElMessageBox.confirm(`确认删除二层网络 ${row.name}？`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    ElMessage.success('删除成功')
    loadL2Networks()
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadL2Networks()
})
</script>

<style scoped>
.l2-networks-page {
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

.topology-placeholder,
.monitoring-placeholder {
  height: 400px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f5f5f5;
  border-radius: 4px;
}
</style>