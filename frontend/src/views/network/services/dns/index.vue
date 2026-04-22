<template>
  <div class="dns-zones-container">
    <!-- 工具栏 -->
    <div class="toolbar">
      <div class="toolbar-left">
        <el-dropdown trigger="click" @command="handleBatchCommand">
          <el-button :disabled="selectedRows.length === 0">
            批量操作 <el-icon class="el-icon--right"><ArrowDown /></el-icon>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="tags" :disabled="selectedRows.length === 0">设置标签</el-dropdown-item>
              <el-dropdown-item command="sync" :disabled="selectedRows.length === 0">同步状态</el-dropdown-item>
              <el-dropdown-item command="delete" :disabled="selectedRows.length === 0" divided>删除</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
      <div class="toolbar-right">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索域名"
          clearable
          style="width: 200px"
          @keyup.enter="loadDNSZones"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
      </div>
    </div>

    <!-- 过滤栏 -->
    <el-card class="filter-card" shadow="never">
      <el-form :inline="true" :model="filters" @submit.prevent="loadDNSZones">
        <el-form-item label="云账号">
          <CloudAccountSelector
            v-model:value="filters.cloud_account_id"
            placeholder="全部云账号"
            @change="handleAccountChange"
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="filters.status" placeholder="全部状态" clearable style="width: 120px">
            <el-option label="可用" value="Available" />
            <el-option label="创建中" value="Creating" />
            <el-option label="错误" value="Error" />
          </el-select>
        </el-form-item>
        <el-form-item label="区域">
          <el-select v-model="filters.region" placeholder="全部区域" clearable style="width: 150px">
            <el-option v-for="r in regions" :key="r.id" :label="r.name" :value="r.id" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadDNSZones">查询</el-button>
          <el-button @click="resetFilters">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 表格 -->
    <el-table
      :data="dnsZones"
      v-loading="loading"
      style="width: 100%"
      row-key="id"
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="50" />
      <el-table-column label="名称/域名" min-width="200">
        <template #default="{ row }">
          <el-link type="primary" @click="viewDetail(row)">{{ row.name || row.dns_zone_id }}</el-link>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)" size="small">
            {{ getStatusLabel(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="标签" width="120">
        <template #default="{ row }">
          <template v-if="row.tags && row.tags.length > 0">
            <el-tag v-for="tag in row.tags.slice(0, 2)" :key="tag.key" size="small" class="tag-item">
              {{ tag.key }}: {{ tag.value }}
            </el-tag>
          </template>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column label="VPC数" width="80">
        <template #default="{ row }">
          {{ row.vpc_count }}
        </template>
      </el-table-column>
      <el-table-column label="归属范围" width="100">
        <template #default="{ row }">
          {{ row.attribution_scope || '-' }}
        </template>
      </el-table-column>
      <el-table-column label="平台/云账号" width="180">
        <template #default="{ row }">
          <div class="platform-cell">
            <el-tag size="small" :type="getPlatformType(row.provider_type)" effect="plain">
              {{ getPlatformLabel(row.provider_type) }}
            </el-tag>
            <span class="account-name">{{ row.account_name || '-' }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="region" label="区域" width="120" />
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button size="small" link type="primary" @click="associateVpc(row)">关联VPC</el-button>
          <el-button size="small" link type="primary" @click="syncStatus(row)">同步状态</el-button>
          <el-button size="small" link type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
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

    <!-- 详情对话框 -->
    <el-dialog v-model="detailVisible" title="DNS解析详情" width="800px">
      <el-tabs v-model="detailTab">
        <el-tab-pane label="详情" name="detail">
          <el-descriptions :column="2" border>
            <el-descriptions-item label="ID">{{ selectedZone?.id }}</el-descriptions-item>
            <el-descriptions-item label="名称">{{ selectedZone?.name }}</el-descriptions-item>
            <el-descriptions-item label="DNS Zone ID">{{ selectedZone?.dns_zone_id }}</el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="getStatusType(selectedZone?.status)">
                {{ getStatusLabel(selectedZone?.status) }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="VPC数">{{ selectedZone?.vpc_count }}</el-descriptions-item>
            <el-descriptions-item label="归属范围">{{ selectedZone?.attribution_scope || '-' }}</el-descriptions-item>
            <el-descriptions-item label="区域">{{ selectedZone?.region }}</el-descriptions-item>
            <el-descriptions-item label="平台">{{ getPlatformLabel(selectedZone?.provider_type) }}</el-descriptions-item>
            <el-descriptions-item label="云账号">{{ selectedZone?.account_name }}</el-descriptions-item>
            <el-descriptions-item label="标签" :span="2">
              <template v-if="selectedZone?.tags && selectedZone.tags.length > 0">
                <el-tag v-for="tag in selectedZone.tags" :key="tag.key" class="tag-item">
                  {{ tag.key }}: {{ tag.value }}
                </el-tag>
              </template>
              <span v-else>-</span>
            </el-descriptions-item>
            <el-descriptions-item label="创建时间" :span="2">{{ selectedZone?.created_at }}</el-descriptions-item>
          </el-descriptions>
        </el-tab-pane>
        <el-tab-pane label="解析记录" name="records">
          <div class="records-toolbar">
            <el-button type="primary" size="small" @click="showAddRecord">添加记录</el-button>
          </div>
          <el-table :data="dnsRecords" v-loading="recordsLoading" style="width: 100%">
            <el-table-column prop="name" label="记录名称" width="150" />
            <el-table-column prop="type" label="类型" width="100" />
            <el-table-column prop="value" label="值" width="200" />
            <el-table-column prop="ttl" label="TTL" width="100" />
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="row.status === 'active' ? 'success' : 'info'" size="small">{{ row.status }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="100">
              <template #default="{ row }">
                <el-button size="small" link type="danger" @click="deleteRecord(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- 关联VPC对话框 -->
    <el-dialog v-model="associateVpcVisible" title="关联VPC" width="500px">
      <el-form label-width="80px">
        <el-form-item label="VPC">
          <el-select v-model="selectedVpcs" placeholder="请选择VPC" multiple style="width: 100%">
            <el-option v-for="vpc in vpcs" :key="vpc.vpc_id" :label="vpc.name" :value="vpc.vpc_id" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="associateVpcVisible = false">取消</el-button>
        <el-button type="primary" @click="submitAssociateVpc">确定</el-button>
      </template>
    </el-dialog>

    <!-- 添加记录对话框 -->
    <el-dialog v-model="addRecordVisible" title="添加DNS记录" width="500px">
      <el-form :model="addRecordForm" label-width="100px">
        <el-form-item label="记录名称">
          <el-input v-model="addRecordForm.name" placeholder="如 @ 或 www" />
        </el-form-item>
        <el-form-item label="记录类型">
          <el-select v-model="addRecordForm.type" style="width: 100%">
            <el-option label="A" value="A" />
            <el-option label="CNAME" value="CNAME" />
            <el-option label="MX" value="MX" />
            <el-option label="TXT" value="TXT" />
            <el-option label="NS" value="NS" />
            <el-option label="AAAA" value="AAAA" />
          </el-select>
        </el-form-item>
        <el-form-item label="记录值">
          <el-input v-model="addRecordForm.value" placeholder="IP地址或域名" />
        </el-form-item>
        <el-form-item label="TTL">
          <el-input-number v-model="addRecordForm.ttl" :min="60" :max="86400" style="width: 200px" />
        </el-form-item>
        <el-form-item label="优先级" v-if="addRecordForm.type === 'MX'">
          <el-input-number v-model="addRecordForm.priority" :min="1" :max="100" style="width: 200px" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="addRecordVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmAddRecord">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowDown, Search } from '@element-plus/icons-vue'
import CloudAccountSelector from '@/components/common/CloudAccountSelector.vue'
import {
  getDNSZones,
  getDNSZone,
  createDNSZone,
  deleteDNSZone,
  batchDeleteDNSZones,
  getDNSRecords,
  createDNSRecord,
  deleteDNSRecord,
  type DNSZone,
  type DNSRecord
} from '@/api/networkSync'
import { getVPCs, type VPC } from '@/api/networkSync'

const dnsZones = ref<DNSZone[]>([])
const loading = ref(false)
const detailVisible = ref(false)
const associateVpcVisible = ref(false)
const addRecordVisible = ref(false)
const selectedZone = ref<DNSZone | null>(null)
const selectedRows = ref<DNSZone[]>([])
const recordsLoading = ref(false)

// 分页
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// 搜索
const searchKeyword = ref('')

// 过滤条件
const filters = reactive({
  cloud_account_id: null as number | null,
  status: '',
  region: ''
})

// 可用资源
const regions = ref<{ id: string; name: string }[]>([])
const vpcs = ref<VPC[]>([])
const selectedVpcs = ref<string[]>([])

// DNS记录
const dnsRecords = ref<DNSRecord[]>([])
const detailTab = ref('detail')

// 添加记录表单
const addRecordForm = reactive({
  name: '',
  type: 'A',
  value: '',
  ttl: 600,
  priority: 10
})

// 平台类型映射
const platformLabels: Record<string, string> = {
  alibaba: '阿里云',
  tencent: '腾讯云',
  aws: 'AWS',
  azure: 'Azure',
  huawei: '华为云',
  google: 'GCP',
  onpremise: '私有云'
}

const platformTypes: Record<string, '' | 'success' | 'warning' | 'info' | 'danger'> = {
  alibaba: 'warning',
  tencent: 'success',
  aws: 'info',
  azure: '',
  huawei: 'danger',
  google: 'success',
  onpremise: 'info'
}

const getPlatformLabel = (platform: string): string => {
  return platformLabels[platform] || platform || '未知'
}

const getPlatformType = (platform: string): '' | 'success' | 'warning' | 'info' | 'danger' => {
  return platformTypes[platform] || 'info'
}

// 状态映射
const statusLabels: Record<string, string> = {
  Available: '可用',
  Creating: '创建中',
  Deleting: '删除中',
  Error: '错误',
  available: '可用',
  creating: '创建中'
}

const getStatusLabel = (status: string): string => {
  return statusLabels[status] || status
}

const getStatusType = (status: string): '' | 'success' | 'warning' | 'info' | 'danger' => {
  if (status === 'Available' || status === 'available') return 'success'
  if (status === 'Creating' || status === 'creating') return 'warning'
  if (status === 'Deleting') return 'warning'
  if (status === 'Error' || status === 'error') return 'danger'
  return 'info'
}

const handleAccountChange = (accountId: number | null) => {
  filters.cloud_account_id = accountId
}

const handleSelectionChange = (rows: DNSZone[]) => {
  selectedRows.value = rows
}

const handleBatchCommand = (command: string) => {
  if (selectedRows.value.length === 0) {
    ElMessage.warning('请先选择要操作的DNS Zone')
    return
  }
  if (command === 'delete') {
    handleBatchDelete()
  } else {
    ElMessage.info(`${command}功能开发中`)
  }
}

const handleBatchDelete = async () => {
  try {
    await ElMessageBox.confirm(`确认批量删除 ${selectedRows.value.length} 个DNS Zone？此操作不可恢复`, '警告', { type: 'warning' })
    const ids = selectedRows.value.map(z => Number(z.id))
    await batchDeleteDNSZones(ids)
    ElMessage.success('批量删除成功')
    loadDNSZones()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error('批量删除失败')
  }
}

const loadDNSZones = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.pageSize
    }
    if (filters.cloud_account_id) params.cloud_account_id = filters.cloud_account_id
    if (filters.status) params.status = filters.status
    if (filters.region) params.region = filters.region
    if (searchKeyword.value) params.name = searchKeyword.value

    const res = await getDNSZones(params)
    dnsZones.value = res.items || []
    pagination.total = res.total || 0

    // 提取区域列表
    const regionMap = new Map<string, string>()
    dnsZones.value.forEach(zone => {
      if (zone.region_id && zone.region) {
        regionMap.set(zone.region_id, zone.region)
      }
    })
    regions.value = Array.from(regionMap.entries()).map(([id, name]) => ({ id, name }))
  } catch (e) {
    console.error(e)
    ElMessage.error('加载DNS Zone列表失败')
  } finally {
    loading.value = false
  }
}

const loadVPCs = async () => {
  try {
    const res = await getVPCs({ page: 1, page_size: 100 })
    vpcs.value = res.items || []
  } catch (e) {
    console.error('加载VPC失败', e)
  }
}

const viewDetail = async (row: DNSZone) => {
  selectedZone.value = row
  detailVisible.value = true
  detailTab.value = 'detail'
  recordsLoading.value = true

  try {
    const res = await getDNSZone(row.id)
    dnsRecords.value = res.records || []
  } catch (e) {
    console.error('加载DNS记录失败', e)
    dnsRecords.value = []
  } finally {
    recordsLoading.value = false
  }
}

const associateVpc = (row: DNSZone) => {
  selectedZone.value = row
  selectedVpcs.value = []
  associateVpcVisible.value = true
  loadVPCs()
}

const submitAssociateVpc = () => {
  ElMessage.success('关联VPC成功')
  associateVpcVisible.value = false
}

const syncStatus = (row: DNSZone) => {
  ElMessage.info(`正在同步DNS Zone ${row.name} 的状态`)
}

const showAddRecord = () => {
  addRecordForm.name = ''
  addRecordForm.type = 'A'
  addRecordForm.value = ''
  addRecordForm.ttl = 600
  addRecordForm.priority = 10
  addRecordVisible.value = true
}

const confirmAddRecord = async () => {
  if (!addRecordForm.name || !addRecordForm.value) {
    ElMessage.warning('请填写记录名称和值')
    return
  }
  try {
    await createDNSRecord(selectedZone.value!.id, addRecordForm)
    ElMessage.success('添加记录成功')
    addRecordVisible.value = false
    viewDetail(selectedZone.value!)
  } catch (e) {
    ElMessage.error('添加记录失败')
  }
}

const deleteRecord = async (row: DNSRecord) => {
  try {
    await ElMessageBox.confirm(`确认删除DNS记录 "${row.name}"？`, '警告', { type: 'warning' })
    await deleteDNSRecord(selectedZone.value!.id, row.id)
    ElMessage.success('删除成功')
    viewDetail(selectedZone.value!)
  } catch (e) {
    if (e !== 'cancel') ElMessage.error('删除失败')
  }
}

const handleDelete = async (row: DNSZone) => {
  try {
    await ElMessageBox.confirm(`确认删除DNS Zone "${row.name}"？此操作不可恢复`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await deleteDNSZone(row.id)
    ElMessage.success('删除成功')
    loadDNSZones()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error('删除失败')
  }
}

const resetFilters = () => {
  filters.cloud_account_id = null
  filters.status = ''
  filters.region = ''
  searchKeyword.value = ''
  pagination.page = 1
  loadDNSZones()
}

const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  pagination.page = 1
  loadDNSZones()
}

const handleCurrentChange = (page: number) => {
  pagination.page = page
  loadDNSZones()
}

onMounted(() => {
  loadDNSZones()
})
</script>

<style scoped>
.dns-zones-container {
  padding: 20px;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.toolbar-left {
  display: flex;
  gap: 8px;
}

.filter-card {
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  justify-content: flex-end;
}

.platform-cell {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.account-name {
  font-size: 12px;
  color: #909399;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.tag-item {
  margin-right: 4px;
  margin-bottom: 2px;
}

.records-toolbar {
  margin-bottom: 16px;
}
</style>