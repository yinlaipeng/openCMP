<template>
  <div class="vpcs-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">VPC列表</span>
          <el-button type="primary" @click="handleCreate">
            <el-icon><Plus /></el-icon>
            创建 VPC
          </el-button>
        </div>
      </template>

      <el-tabs v-model="activeTab">
        <el-tab-pane label="全部" name="all">
          <el-table :data="allVpcs" v-loading="loading">
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
            <el-table-column prop="cidr" label="ipv4目标网段" width="180" />
            <el-table-column prop="ipv6_cidr" label="ipv6目标网段" width="180" />
            <el-table-column prop="allow_internet_access" label="允许外网访问" width="120">
              <template #default="{ row }">
                <el-switch v-model="row.allow_internet_access" disabled />
              </template>
            </el-table-column>
            <el-table-column prop="subnet_count" label="ip子网数量" width="100" />
            <el-table-column prop="platform" label="平台" width="120" />
            <el-table-column prop="account" label="云账号" width="150" />
            <el-table-column prop="domain" label="所属域" width="150" />
            <el-table-column prop="region_id" label="区域" width="150" />
            <el-table-column label="操作" width="200">
              <template #default="{ row }">
                <el-button size="small" type="success" @click="syncStatus(row)">同步状态</el-button>
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
        </el-tab-pane>
        <el-tab-pane label="本地idc" name="local_idc">
          <el-table :data="localIdcVpcs" v-loading="loading">
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
            <el-table-column prop="cidr" label="ipv4目标网段" width="180" />
            <el-table-column prop="ipv6_cidr" label="ipv6目标网段" width="180" />
            <el-table-column prop="allow_internet_access" label="允许外网访问" width="120">
              <template #default="{ row }">
                <el-switch v-model="row.allow_internet_access" disabled />
              </template>
            </el-table-column>
            <el-table-column prop="subnet_count" label="ip子网数量" width="100" />
            <el-table-column prop="platform" label="平台" width="120" />
            <el-table-column prop="account" label="云账号" width="150" />
            <el-table-column prop="domain" label="所属域" width="150" />
            <el-table-column prop="region_id" label="区域" width="150" />
            <el-table-column label="操作" width="200">
              <template #default="{ row }">
                <el-button size="small" type="success" @click="syncStatus(row)">同步状态</el-button>
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
        </el-tab-pane>
        <el-tab-pane label="公有云" name="public_cloud">
          <el-table :data="publicCloudVpcs" v-loading="loading">
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
            <el-table-column prop="cidr" label="ipv4目标网段" width="180" />
            <el-table-column prop="ipv6_cidr" label="ipv6目标网段" width="180" />
            <el-table-column prop="allow_internet_access" label="允许外网访问" width="120">
              <template #default="{ row }">
                <el-switch v-model="row.allow_internet_access" disabled />
              </template>
            </el-table-column>
            <el-table-column prop="subnet_count" label="ip子网数量" width="100" />
            <el-table-column prop="platform" label="平台" width="120" />
            <el-table-column prop="account" label="云账号" width="150" />
            <el-table-column prop="domain" label="所属域" width="150" />
            <el-table-column prop="region_id" label="区域" width="150" />
            <el-table-column label="操作" width="200">
              <template #default="{ row }">
                <el-button size="small" type="success" @click="syncStatus(row)">同步状态</el-button>
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
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- Detail Modal -->
    <el-dialog v-model="detailDialogVisible" title="VPC详情" width="800px">
      <el-tabs v-model="detailTab">
        <el-tab-pane label="详情" name="detail">
          <el-descriptions :column="2" border>
            <el-descriptions-item label="ID">{{ selectedVPC?.id }}</el-descriptions-item>
            <el-descriptions-item label="名称">{{ selectedVPC?.name }}</el-descriptions-item>
            <el-descriptions-item label="状态">{{ selectedVPC?.status }}</el-descriptions-item>
            <el-descriptions-item label="IPv4目标网段">{{ selectedVPC?.cidr }}</el-descriptions-item>
            <el-descriptions-item label="IPv6目标网段">{{ selectedVPC?.ipv6_cidr }}</el-descriptions-item>
            <el-descriptions-item label="允许外网访问">{{ selectedVPC?.allow_internet_access ? '是' : '否' }}</el-descriptions-item>
            <el-descriptions-item label="IP子网数量">{{ selectedVPC?.subnet_count }}</el-descriptions-item>
            <el-descriptions-item label="平台">{{ selectedVPC?.platform }}</el-descriptions-item>
            <el-descriptions-item label="云账号">{{ selectedVPC?.account }}</el-descriptions-item>
            <el-descriptions-item label="所属域">{{ selectedVPC?.domain }}</el-descriptions-item>
            <el-descriptions-item label="区域">{{ selectedVPC?.region_id }}</el-descriptions-item>
            <el-descriptions-item label="创建时间">{{ selectedVPC?.created_at }}</el-descriptions-item>
          </el-descriptions>
        </el-tab-pane>
        <el-tab-pane label="安全组" name="security_groups">
          <el-table :data="securityGroups" v-loading="securityGroupsLoading">
            <el-table-column prop="id" label="安全组ID" width="200" />
            <el-table-column prop="name" label="安全组名称" width="150" />
            <el-table-column prop="status" label="状态" width="100" />
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="ip子网" name="subnets">
          <el-table :data="subnets" v-loading="subnetsLoading">
            <el-table-column prop="id" label="子网ID" width="200" />
            <el-table-column prop="name" label="子网名称" width="150" />
            <el-table-column prop="status" label="状态" width="100" />
            <el-table-column prop="cidr" label="CIDR" width="180" />
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="拓扑" name="topology">
          <div class="topology-placeholder">拓扑图待实现</div>
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

    <!-- 创建 VPC 模态框 -->
    <CreateVPCModal
      v-model:visible="createModalVisible"
      @success="handleCreateSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowDown, Plus } from '@element-plus/icons-vue'
import CreateVPCModal from '@/components/network/CreateVPCModal.vue'

interface ExtendedVPC {
  id: string
  name: string
  status: string
  cidr: string
  ipv6_cidr?: string
  subnet_count?: number
  allow_internet_access?: boolean
  platform?: string
  account?: string
  domain?: string
  region_id?: string
  created_at?: string
}

const allVpcs = ref<ExtendedVPC[]>([])
const localIdcVpcs = ref<ExtendedVPC[]>([])
const publicCloudVpcs = ref<ExtendedVPC[]>([])
const loading = ref(false)
const activeTab = ref('all')
const detailDialogVisible = ref(false)
const changeDomainDialogVisible = ref(false)
const createModalVisible = ref(false)
const selectedVPC = ref<ExtendedVPC | null>(null)
const detailTab = ref('detail')
const targetDomain = ref('')

// Detail modal data
const securityGroups = ref<any[]>([])
const securityGroupsLoading = ref(false)
const subnets = ref<any[]>([])
const subnetsLoading = ref(false)
const logs = ref<any[]>([])

const getStatusType = (status: string) => {
  switch (status.toLowerCase()) {
    case 'available':
    case 'running':
    case 'active':
      return 'success'
    case 'pending':
      return 'warning'
    case 'failed':
    case 'error':
      return 'danger'
    default:
      return 'info'
  }
}

const loadVPCs = async () => {
  loading.value = true
  try {
    // Mock data - API requires valid cloud account credentials
    const vpcsData: ExtendedVPC[] = [
      {
        id: 'vpc-1',
        name: 'VPC 1',
        status: 'Available',
        cidr: '10.0.0.0/16',
        ipv6_cidr: '2001:db8::/64',
        subnet_count: 5,
        allow_internet_access: true,
        platform: '阿里云',
        account: 'Aliyun Account 1',
        domain: 'Default Domain',
        region_id: 'cn-hangzhou',
        created_at: '2024-01-01 10:00:00'
      },
      {
        id: 'vpc-2',
        name: 'VPC 2 (本地IDC)',
        status: 'Available',
        cidr: '10.1.0.0/16',
        ipv6_cidr: '2001:db8::/64',
        subnet_count: 3,
        allow_internet_access: false,
        platform: '阿里云',
        account: 'Aliyun Account 1',
        domain: 'Default Domain',
        region_id: 'cn-hangzhou',
        created_at: '2024-01-01 10:00:00'
      },
      {
        id: 'vpc-3',
        name: 'VPC 3',
        status: 'Pending',
        cidr: '10.2.0.0/16',
        ipv6_cidr: '2001:db8::/64',
        subnet_count: 2,
        allow_internet_access: true,
        platform: '阿里云',
        account: 'Aliyun Account 1',
        domain: 'Default Domain',
        region_id: 'cn-shanghai',
        created_at: '2024-01-01 10:00:00'
      }
    ]
    allVpcs.value = vpcsData
    localIdcVpcs.value = vpcsData.filter(vpc => vpc.name.toLowerCase().includes('idc') || vpc.name.toLowerCase().includes('local'))
    publicCloudVpcs.value = vpcsData.filter(vpc => !vpc.name.toLowerCase().includes('idc') && !vpc.name.toLowerCase().includes('local'))
  } catch (e) {
    console.error(e)
    allVpcs.value = []
    localIdcVpcs.value = []
    publicCloudVpcs.value = []
  } finally {
    loading.value = false
  }
}

const viewDetail = async (row: ExtendedVPC) => {
  selectedVPC.value = row
  detailDialogVisible.value = true
  detailTab.value = 'detail'

  // Mock data for security groups and subnets (API requires valid cloud account)
  securityGroupsLoading.value = true
  securityGroups.value = [
    { id: 'sg-1', name: 'Security Group 1', status: 'Active' },
    { id: 'sg-2', name: 'Security Group 2', status: 'Active' }
  ]
  securityGroupsLoading.value = false

  subnetsLoading.value = true
  subnets.value = [
    { id: 'subnet-1', name: 'Subnet 1', status: 'Available', cidr: '10.0.1.0/24' },
    { id: 'subnet-2', name: 'Subnet 2', status: 'Available', cidr: '10.0.2.0/24' }
  ]
  subnetsLoading.value = false

  // Load logs (mock)
  logs.value = [
    { operation: '创建', operator: 'admin', result: '成功', timestamp: '2024-01-01 10:00:00' }
  ]
}

const syncStatus = (row: ExtendedVPC) => {
  ElMessage.info(`正在同步 VPC ${row.name} 的状态`)
}

const changeDomain = (row: ExtendedVPC) => {
  selectedVPC.value = row
  changeDomainDialogVisible.value = true
}

const submitChangeDomain = () => {
  ElMessage.success('更改域成功')
  changeDomainDialogVisible.value = false
  loadVPCs()
}

const handleDelete = async (row: ExtendedVPC) => {
  try {
    await ElMessageBox.confirm(`确认删除 VPC ${row.name}？`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    // Mock deletion - API requires valid cloud account credentials
    allVpcs.value = allVpcs.value.filter(v => v.id !== row.id)
    localIdcVpcs.value = localIdcVpcs.value.filter(v => v.id !== row.id)
    publicCloudVpcs.value = publicCloudVpcs.value.filter(v => v.id !== row.id)
    ElMessage.success('删除成功')
  } catch (e) {
    console.error(e)
  }
}

const handleCreate = () => {
  createModalVisible.value = true
}

const handleCreateSuccess = (vpc: ExtendedVPC) => {
  ElMessage.success(`${vpc.name} 创建成功`)
  loadVPCs()
}

onMounted(() => {
  loadVPCs()
})
</script>

<style scoped>
.vpcs-page {
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

.topology-placeholder {
  height: 400px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f5f5f5;
  border-radius: 4px;
}
</style>