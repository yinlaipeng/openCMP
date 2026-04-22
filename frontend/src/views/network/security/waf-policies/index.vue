<template>
  <div class="waf-policies-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">WAF策略列表</span>
          <div class="toolbar">
          <el-button link @click="handleViewMode">
          <el-icon><View /></el-icon>
           查看
          </el-button>
          <el-button :disabled="selectedRows.length === 0" @click="handleBatchDelete">
           删除
          </el-button>
          <el-button :disabled="selectedRows.length === 0" @click="handleSetTags">
           设置标签
           </el-button>
				<el-button @click="handleRefresh">
					<el-icon><Refresh /></el-icon>
				</el-button>
				<el-button @click="handleTags">
					标签
				</el-button>
			</div>
        </div>
      </template>

      <!-- 搜索表单 -->
      <el-form :inline="true" :model="queryParams" class="search-form" @submit.prevent="loadData">
        <el-form-item label="名称">
          <el-input v-model="queryParams.name" placeholder="搜索名称" clearable style="width: 200px" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="queryParams.status" placeholder="选择状态" clearable style="width: 140px">
            <el-option label="正常" value="normal" />
            <el-option label="创建中" value="creating" />
            <el-option label="删除中" value="deleting" />
            <el-option label="错误" value="error" />
          </el-select>
        </el-form-item>
        <el-form-item label="平台">
          <el-select v-model="queryParams.platform" placeholder="选择平台" clearable style="width: 140px">
            <el-option label="阿里云" value="alibaba" />
            <el-option label="腾讯云" value="tencent" />
            <el-option label="AWS" value="aws" />
            <el-option label="Azure" value="azure" />
          </el-select>
        </el-form-item>
        <el-form-item label="云账户">
          <el-select v-model="queryParams.cloud_account_id" placeholder="选择云账户" clearable style="width: 160px">
            <el-option v-for="account in cloudAccounts" :key="account.id" :label="account.name" :value="account.id" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadData">查询</el-button>
          <el-button @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table
        :data="wafPolicies"
        v-loading="loading"
        @selection-change="handleSelectionChange"
        row-key="id"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column label="名称" min-width="150">
          <template #default="{ row }">
            <el-link type="primary" :underline="false" @click="viewDetail(row)">{{ row.name }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="tags" label="标签" width="150">
          <template #default="{ row }">
            <span v-if="row.tags">{{ formatTags(row.tags) }}</span>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusName(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="type" label="类型" width="120">
          <template #default="{ row }">
            <span>{{ row.type || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="platform" label="平台" width="100">
          <template #default="{ row }">
            <el-tag size="small" :type="getPlatformType(row.platform)" effect="plain">
              {{ getPlatformLabel(row.platform) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="云账户" width="150">
          <template #default="{ row }">
            <span>{{ row.cloud_account?.name || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="domain_id" label="所属域" width="120">
          <template #default="{ row }">
            <span>{{ row.domain_id || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="region_id" label="区域" width="120">
          <template #default="{ row }">
            <span>{{ row.region_id || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button size="small" link @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" link type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="queryParams.page"
        v-model:page-size="queryParams.page_size"
        :total="total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next"
        @size-change="loadData"
        @current-change="loadData"
        style="margin-top: 16px; justify-content: flex-end"
      />
    </el-card>

    <!-- Detail Modal -->
    <el-dialog v-model="detailDialogVisible" title="WAF策略详情" width="700px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ selectedWAFPolicy?.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ selectedWAFPolicy?.name }}</el-descriptions-item>
        <el-descriptions-item label="标签">{{ formatTags(selectedWAFPolicy?.tags) || '-' }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(selectedWAFPolicy?.status || '')">
            {{ getStatusName(selectedWAFPolicy?.status || '') }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="类型">{{ selectedWAFPolicy?.type || '-' }}</el-descriptions-item>
        <el-descriptions-item label="平台">{{ getPlatformLabel(selectedWAFPolicy?.platform) }}</el-descriptions-item>
        <el-descriptions-item label="云账户">{{ selectedWAFPolicy?.cloud_account?.name || '-' }}</el-descriptions-item>
        <el-descriptions-item label="所属域">{{ selectedWAFPolicy?.domain_id || '-' }}</el-descriptions-item>
        <el-descriptions-item label="区域">{{ selectedWAFPolicy?.region_id || '-' }}</el-descriptions-item>
        <el-descriptions-item label="外部ID">{{ selectedWAFPolicy?.external_id || '-' }}</el-descriptions-item>
        <el-descriptions-item label="描述">{{ selectedWAFPolicy?.description || '-' }}</el-descriptions-item>
        <el-descriptions-item label="同步时间">{{ selectedWAFPolicy?.sync_time || '-' }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ selectedWAFPolicy?.created_at }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- Edit Modal -->
    <el-dialog v-model="editDialogVisible" title="编辑WAF策略" width="500px">
      <el-form :model="editForm" label-width="80px">
        <el-form-item label="名称">
          <el-input v-model="editForm.name" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="editForm.description" type="textarea" :rows="3" />
        </el-form-item>
        <el-form-item label="启用">
          <el-switch v-model="editForm.enabled" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitEdit" :loading="editLoading">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Refresh, View } from '@element-plus/icons-vue'
import {
  getWAFList,
  getWAFDetail,
  updateWAF,
  deleteWAF,
  batchDeleteWAF,
  syncWAFStatus,
  WAFInstance,
  WAFListParams,
  UpdateWAFParams
} from '@/api/waf'
import { getCloudAccounts } from '@/api/cloud-account'

const wafPolicies = ref<WAFInstance[]>([])
const cloudAccounts = ref<{ id: number; name: string }[]>([])
const loading = ref(false)
const total = ref(0)
const selectedRows = ref<WAFInstance[]>([])
const detailDialogVisible = ref(false)
const editDialogVisible = ref(false)
const editLoading = ref(false)
const selectedWAFPolicy = ref<WAFInstance | null>(null)

const queryParams = ref<WAFListParams>({
  page: 1,
  page_size: 20,
  name: '',
  status: '',
  platform: '',
  cloud_account_id: '',
  domain_id: ''
})

const editForm = ref<UpdateWAFParams>({
  name: '',
  description: '',
  enabled: true
})

const getStatusType = (status: string) => {
  switch (status?.toLowerCase()) {
    case 'normal':
      return 'success'
    case 'creating':
    case 'configuring':
      return 'warning'
    case 'deleting':
    case 'error':
      return 'danger'
    default:
      return 'info'
  }
}

const getStatusName = (status: string) => {
  switch (status?.toLowerCase()) {
    case 'normal':
      return '正常'
    case 'creating':
      return '创建中'
    case 'deleting':
      return '删除中'
    case 'error':
      return '错误'
    default:
      return status || '未知'
  }
}

const getPlatformType = (platform: string) => {
  switch (platform?.toLowerCase()) {
    case 'alibaba':
    case 'aliyun':
      return ''
    case 'tencent':
      return 'success'
    case 'aws':
      return 'warning'
    case 'azure':
      return 'info'
    default:
      return ''
  }
}

const getPlatformLabel = (platform: string) => {
  switch (platform?.toLowerCase()) {
    case 'alibaba':
    case 'aliyun':
      return '阿里云'
    case 'tencent':
      return '腾讯云'
    case 'aws':
      return 'AWS'
    case 'azure':
      return 'Azure'
    default:
      return platform || '-'
  }
}

const formatTags = (tags: Record<string, string> | null | undefined) => {
  if (!tags) return '-'
  return Object.entries(tags).map(([k, v]) => `${k}:${v}`).join(', ')
}

const loadData = async () => {
  loading.value = true
  try {
    const res = await getWAFList(queryParams.value)
    wafPolicies.value = res.data.items || []
    total.value = res.data.total || 0
  } catch (e: any) {
    console.error(e)
    ElMessage.error(e.message || '加载WAF策略列表失败')
    wafPolicies.value = []
  } finally {
    loading.value = false
  }
}

const loadCloudAccounts = async () => {
  try {
    const res = await getCloudAccounts()
    cloudAccounts.value = res.data.items || []
  } catch (e) {
    console.error(e)
  }
}

const resetQuery = () => {
  queryParams.value = {
    page: 1,
    page_size: 20,
    name: '',
    status: '',
    platform: '',
    cloud_account_id: '',
    domain_id: ''
  }
  loadData()
}

const handleSelectionChange = (rows: WAFInstance[]) => {
  selectedRows.value = rows
}

const viewDetail = async (row: WAFInstance) => {
  try {
    const res = await getWAFDetail(row.id)
    selectedWAFPolicy.value = res.data
    detailDialogVisible.value = true
  } catch (e: any) {
    ElMessage.error(e.message || '获取详情失败')
  }
}

const handleEdit = (row: WAFInstance) => {
  editForm.value = {
    name: row.name,
    description: row.description,
    enabled: row.enabled
  }
  selectedWAFPolicy.value = row
  editDialogVisible.value = true
}

const submitEdit = async () => {
  if (!selectedWAFPolicy.value) return
  editLoading.value = true
  try {
    await updateWAF(selectedWAFPolicy.value.id, editForm.value)
    ElMessage.success('更新成功')
    editDialogVisible.value = false
    loadData()
  } catch (e: any) {
    ElMessage.error(e.message || '更新失败')
  } finally {
    editLoading.value = false
  }
}

const handleDelete = async (row: WAFInstance) => {
  try {
    await ElMessageBox.confirm(`确认删除WAF策略 ${row.name}？`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await deleteWAF(row.id)
    ElMessage.success('删除成功')
    loadData()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '删除失败')
    }
  }
}

const handleBatchDelete = async () => {
  if (selectedRows.value.length === 0) return
  try {
    await ElMessageBox.confirm(`确认删除选中的 ${selectedRows.value.length} 个WAF策略？`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await batchDeleteWAF({ ids: selectedRows.value.map(r => r.id) })
    ElMessage.success('批量删除成功')
    loadData()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '批量删除失败')
    }
  }
}

const handleSetTags = () => {
  ElMessage.info('标签设置功能待实现')
}

const handleRefresh = () => {
loadData()
}

const handleViewMode = () => {
	ElMessage.info('视图切换功能待实现')
}

const handleTags = () => {
	ElMessage.info('标签筛选功能待实现')
}

onMounted(() => {
  loadCloudAccounts()
  loadData()
})
</script>

<style scoped>
.waf-policies-page {
  height: 100%;
  padding: 16px;
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

.toolbar {
  display: flex;
  gap: 8px;
}

.search-form {
  margin-bottom: 16px;
}
</style>