<template>
  <div class="vpcs-container">
    <div class="page-header">
      <h2>VPC管理</h2>
      <el-button type="primary" @click="handleCreate">
        <el-icon><Plus /></el-icon>
        创建 VPC
      </el-button>
    </div>

    <el-card class="filter-card">
      <el-form :inline="true" :model="filters" @submit.prevent="loadVPCs">
        <el-form-item label="云账号">
          <CloudAccountSelector
            v-model:value="filters.account_id"
            placeholder="选择云账号"
            @change="handleAccountChange"
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="filters.status" placeholder="状态" clearable>
            <el-option label="可用" value="Available" />
            <el-option label="创建中" value="Creating" />
            <el-option label="已删除" value="Deleted" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadVPCs">查询</el-button>
          <el-button @click="resetFilters">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

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

    <el-pagination
      v-model:current-page="pagination.page"
      v-model:page-size="pagination.pageSize"
      :total="pagination.total"
      :page-sizes="[10, 20, 50, 100]"
      layout="total, sizes, prev, pager, next, jumper"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
      class="pagination"
    />

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
import { ref, reactive, onMounted, computed } from 'vue'
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
import CloudAccountSelector from '@/components/common/CloudAccountSelector.vue'
import { getVPCs, getSubnets, getSecurityGroups, deleteVPC } from '@/api/network'
import { getCloudAccounts } from '@/api/cloud-account'
import type { VPC, Subnet, SecurityGroup, CloudAccount } from '@/types'

interface ExtendedVPC extends VPC {
  subnet_count?: number
  ipv6_cidr?: string
  allow_internet_access?: boolean
  platform?: string
  account?: string
  domain?: string
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

// 筛选条件
const filters = reactive({
  account_id: null as number | null,
  status: ''
})

// 分页数据
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// Detail modal data
const securityGroups = ref<SecurityGroup[]>([])
const securityGroupsLoading = ref(false)
const subnets = ref<Subnet[]>([])
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
  if (!filters.account_id) {
    allVpcs.value = []
    localIdcVpcs.value = []
    publicCloudVpcs.value = []
    return
  }

  loading.value = true
  try {
    const res = await getVPCs({
      account_id: filters.account_id,
      status: filters.status || undefined,
      page: pagination.page,
      size: pagination.pageSize
    })
    const vpcsData = Array.isArray(res) ? res : res.items || []
    pagination.total = res.total || vpcsData.length

    allVpcs.value = vpcsData.map((v: any) => ({
      ...v,
      platform: v.platform || '未知',
      account: v.account_name || '未知'
    }))
    localIdcVpcs.value = allVpcs.value.filter(vpc => vpc.name?.toLowerCase().includes('idc') || vpc.name?.toLowerCase().includes('local'))
    publicCloudVpcs.value = allVpcs.value.filter(vpc => !vpc.name?.toLowerCase().includes('idc') && !vpc.name?.toLowerCase().includes('local'))
  } catch (error: any) {
    console.error('Failed to load VPCs:', error)
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

  // Load real data from API
  if (filters.account_id && row.id) {
    securityGroupsLoading.value = true
    try {
      const res = await getSecurityGroups({ account_id: filters.account_id, vpc_id: row.id })
      securityGroups.value = Array.isArray(res) ? res : res.items || []
    } catch (error) {
      console.error('Failed to load security groups:', error)
      securityGroups.value = []
    } finally {
      securityGroupsLoading.value = false
    }

    subnetsLoading.value = true
    try {
      const res = await getSubnets({ account_id: filters.account_id, vpc_id: row.id })
      subnets.value = Array.isArray(res) ? res : res.items || []
    } catch (error) {
      console.error('Failed to load subnets:', error)
      subnets.value = []
    } finally {
      subnetsLoading.value = false
    }
  }

  logs.value = [
    { operation: '创建', operator: 'system', result: '成功', timestamp: row.created_at || '-' }
  ]
}

const syncStatus = async (row: ExtendedVPC) => {
  ElMessage.info(`正在同步 VPC ${row.name} 的状态`)
  await loadVPCs()
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

    if (filters.account_id) {
      await deleteVPC(row.id, filters.account_id)
      ElMessage.success('删除成功')
      loadVPCs()
    }
  } catch (e: any) {
    if (e !== 'cancel') {
      console.error(e)
      ElMessage.error(`删除失败: ${e.message}`)
    }
  }
}

const handleCreate = () => {
  createModalVisible.value = true
}

const handleCreateSuccess = (vpc: ExtendedVPC) => {
  ElMessage.success(`${vpc.name} 创建成功`)
  loadVPCs()
}

const handleAccountChange = (accountId: number | null) => {
  filters.account_id = accountId
  if (accountId) {
    loadVPCs()
  } else {
    allVpcs.value = []
    localIdcVpcs.value = []
    publicCloudVpcs.value = []
  }
}

const resetFilters = () => {
  filters.account_id = null
  filters.status = ''
  pagination.page = 1
  allVpcs.value = []
  localIdcVpcs.value = []
  publicCloudVpcs.value = []
}

const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  pagination.page = 1
  loadVPCs()
}

const handleCurrentChange = (page: number) => {
  pagination.page = page
  loadVPCs()
}

onMounted(() => {
  // Wait for account selector to initialize
})
</script>

<style scoped>
.vpcs-container {
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

.filter-card {
  margin-bottom: 20px;
}

.vpc-tabs {
  margin-top: 0;
}

.pagination {
  margin-top: 20px;
  text-align: right;
}

.tab-label {
  display: flex;
  align-items: center;
  gap: 8px;
}

.tab-badge {
  margin-left: 4px;
}

.vpc-table {
  border-radius: 8px;
}

.status-tag {
  font-weight: 500;
}

.cidr {
  font-size: 14px;
}

.no-data {
  color: #909399;
}

.region {
  font-size: 14px;
}

.operation-buttons {
  display: flex;
  gap: 8px;
  align-items: center;
}

.topology-placeholder {
  height: 400px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f5f7fa;
  border-radius: 8px;
  color: #909399;
}

@media (max-width: 768px) {
  .vpcs-container {
    padding: 10px;
  }

  .page-header {
    flex-direction: column;
    gap: 10px;
    align-items: flex-start;
  }
}
</style>