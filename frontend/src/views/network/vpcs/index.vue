<template>
  <div class="vpcs-page">
    <el-card class="page-card glass">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <span class="title">VPC列表</span>
            <el-tag size="small" type="info" class="count-tag">共 {{ allVpcs.length }} 个</el-tag>
          </div>
          <el-button type="primary" class="create-btn" @click="handleCreate">
            <el-icon><Plus /></el-icon>
            创建 VPC
          </el-button>
        </div>
      </template>

      <el-tabs v-model="activeTab" class="vpc-tabs">
        <el-tab-pane name="all">
          <template #label>
            <span class="tab-label">
              <el-icon><Grid /></el-icon>
              全部
              <el-badge :value="allVpcs.length" :max="99" class="tab-badge" />
            </span>
          </template>
          <el-table :data="allVpcs" v-loading="loading" class="vpc-table">
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
            <el-table-column prop="cidr" label="IPv4网段" width="180">
              <template #default="{ row }">
                <span class="cidr font-mono">{{ row.cidr }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="ipv6_cidr" label="IPv6网段" width="180">
              <template #default="{ row }">
                <span v-if="row.ipv6_cidr" class="cidr font-mono">{{ row.ipv6_cidr }}</span>
                <span v-else class="no-data">-</span>
              </template>
            </el-table-column>
            <el-table-column prop="subnet_count" label="子网数" width="80">
              <template #default="{ row }">
                <span class="font-mono">{{ row.subnet_count || 0 }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="platform" label="平台" width="100">
              <template #default="{ row }">
                <el-tag size="small" effect="plain">{{ row.platform }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="region_id" label="区域" width="120">
              <template #default="{ row }">
                <span class="region font-mono">{{ row.region_id }}</span>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="180" fixed="right">
              <template #default="{ row }">
                <div class="operation-buttons">
                  <el-button size="small" type="success" @click="syncStatus(row)">
                    <el-icon><Refresh /></el-icon>
                    同步
                  </el-button>
                  <el-dropdown trigger="click">
                    <el-button size="small">
                      更多<el-icon class="el-icon--right"><ArrowDown /></el-icon>
                    </el-button>
                    <template #dropdown>
                      <el-dropdown-menu>
                        <el-dropdown-item @click="changeDomain(row)">
                          <el-icon><SwitchButton /></el-icon>更改域
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
        </el-tab-pane>
        <el-tab-pane name="local_idc">
          <template #label>
            <span class="tab-label">
              <el-icon><OfficeBuilding /></el-icon>
              本地IDC
              <el-badge :value="localIdcVpcs.length" :max="99" class="tab-badge" />
            </span>
          </template>
          <el-table :data="localIdcVpcs" v-loading="loading" class="vpc-table">
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
            <el-table-column prop="cidr" label="IPv4网段" width="180">
              <template #default="{ row }">
                <span class="cidr font-mono">{{ row.cidr }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="ipv6_cidr" label="IPv6网段" width="180">
              <template #default="{ row }">
                <span v-if="row.ipv6_cidr" class="cidr font-mono">{{ row.ipv6_cidr }}</span>
                <span v-else class="no-data">-</span>
              </template>
            </el-table-column>
            <el-table-column prop="subnet_count" label="子网数" width="80">
              <template #default="{ row }">
                <span class="font-mono">{{ row.subnet_count || 0 }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="platform" label="平台" width="100">
              <template #default="{ row }">
                <el-tag size="small" effect="plain">{{ row.platform }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="region_id" label="区域" width="120">
              <template #default="{ row }">
                <span class="region font-mono">{{ row.region_id }}</span>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="180" fixed="right">
              <template #default="{ row }">
                <div class="operation-buttons">
                  <el-button size="small" type="success" @click="syncStatus(row)">
                    <el-icon><Refresh /></el-icon>
                    同步
                  </el-button>
                  <el-dropdown trigger="click">
                    <el-button size="small">
                      更多<el-icon class="el-icon--right"><ArrowDown /></el-icon>
                    </el-button>
                    <template #dropdown>
                      <el-dropdown-menu>
                        <el-dropdown-item @click="changeDomain(row)">更改域</el-dropdown-item>
                        <el-dropdown-item divided @click="handleDelete(row)">删除</el-dropdown-item>
                      </el-dropdown-menu>
                    </template>
                  </el-dropdown>
                </div>
              </template>
            </el-table-column>
          </el-table>
          <el-empty v-if="!loading && localIdcVpcs.length === 0" description="暂无本地IDC VPC" :image-size="60" />
        </el-tab-pane>
        <el-tab-pane name="public_cloud">
          <template #label>
            <span class="tab-label">
              <el-icon><Cloudy /></el-icon>
              公有云
              <el-badge :value="publicCloudVpcs.length" :max="99" class="tab-badge" />
            </span>
          </template>
          <el-table :data="publicCloudVpcs" v-loading="loading" class="vpc-table">
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
            <el-table-column prop="cidr" label="IPv4网段" width="180">
              <template #default="{ row }">
                <span class="cidr font-mono">{{ row.cidr }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="ipv6_cidr" label="IPv6网段" width="180">
              <template #default="{ row }">
                <span v-if="row.ipv6_cidr" class="cidr font-mono">{{ row.ipv6_cidr }}</span>
                <span v-else class="no-data">-</span>
              </template>
            </el-table-column>
            <el-table-column prop="subnet_count" label="子网数" width="80">
              <template #default="{ row }">
                <span class="font-mono">{{ row.subnet_count || 0 }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="platform" label="平台" width="100">
              <template #default="{ row }">
                <el-tag size="small" effect="plain">{{ row.platform }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="region_id" label="区域" width="120">
              <template #default="{ row }">
                <span class="region font-mono">{{ row.region_id }}</span>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="180" fixed="right">
              <template #default="{ row }">
                <div class="operation-buttons">
                  <el-button size="small" type="success" @click="syncStatus(row)">
                    <el-icon><Refresh /></el-icon>
                    同步
                  </el-button>
                  <el-dropdown trigger="click">
                    <el-button size="small">
                      更多<el-icon class="el-icon--right"><ArrowDown /></el-icon>
                    </el-button>
                    <template #dropdown>
                      <el-dropdown-menu>
                        <el-dropdown-item @click="changeDomain(row)">更改域</el-dropdown-item>
                        <el-dropdown-item divided @click="handleDelete(row)">删除</el-dropdown-item>
                      </el-dropdown-menu>
                    </template>
                  </el-dropdown>
                </div>
              </template>
            </el-table-column>
          </el-table>
          <el-empty v-if="!loading && publicCloudVpcs.length === 0" description="暂无公有云 VPC" :image-size="60" />
        </el-tab-pane>
      </el-tabs>

      <!-- Empty State -->
      <el-empty v-if="!loading && allVpcs.length === 0" description="暂无 VPC 数据">
        <el-button type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>
          创建第一个 VPC
        </el-button>
      </el-empty>
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
import {
  ArrowDown,
  Plus,
  Grid,
  OfficeBuilding,
  Cloudy,
  Refresh,
  SwitchButton,
  Delete
} from '@element-plus/icons-vue'
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

/* Tabs Styling */
.vpc-tabs {
  margin-top: var(--space-2);
}

.tab-label {
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.tab-badge {
  margin-left: var(--space-1);
}

/* Table */
.vpc-table {
  border-radius: var(--radius-md);
}

.status-tag {
  font-weight: var(--font-weight-medium);
}

.cidr {
  font-size: var(--font-size-sm);
}

.no-data {
  color: var(--color-muted);
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
.topology-placeholder {
  height: 400px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: var(--color-background);
  border-radius: var(--radius-md);
  color: var(--color-muted);
}

/* Responsive */
@media (max-width: 768px) {
  .vpcs-page {
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

  .tab-label {
    font-size: var(--font-size-sm);
  }
}
</style>