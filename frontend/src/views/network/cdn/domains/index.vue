<template>
  <div class="cdn-domains-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">CDN域名</span>
          <div class="toolbar">
            <el-button @click="handleView">查看</el-button>
            <el-button type="primary" @click="handleCreate">创建</el-button>
            <el-button :disabled="!hasSelection" @click="handleSyncStatus">同步状态</el-button>
            <el-dropdown :disabled="!hasSelection" @command="handleBatchCommand">
              <el-button>
                批量操作 <el-icon class="el-icon--right"><arrow-down /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="sync">同步状态</el-dropdown-item>
                  <el-dropdown-item command="delete">批量删除</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
            <el-button @click="handleTags">标签</el-button>
          </div>
        </div>
      </template>

      <!-- Search Filters -->
      <el-form :inline="true" class="search-form">
        <el-form-item label="名称">
          <el-input v-model="searchParams.name" placeholder="请输入名称" clearable @clear="loadCDNDomains" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchParams.status" placeholder="请选择状态" clearable @clear="loadCDNDomains">
            <el-option label="在线" value="online" />
            <el-option label="配置中" value="configuring" />
            <el-option label="离线" value="offline" />
          </el-select>
        </el-form-item>
        <el-form-item label="加速区域">
          <el-select v-model="searchParams.area" placeholder="请选择区域" clearable @clear="loadCDNDomains">
            <el-option label="中国大陆" value="china" />
            <el-option label="全球" value="global" />
            <el-option label="亚太地区" value="asia" />
          </el-select>
        </el-form-item>
        <el-form-item label="平台">
          <el-select v-model="searchParams.platform" placeholder="请选择平台" clearable @clear="loadCDNDomains">
            <el-option label="阿里云" value="aliyun" />
            <el-option label="腾讯云" value="tencent" />
            <el-option label="AWS" value="aws" />
            <el-option label="Azure" value="azure" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadCDNDomains">查询</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="cdnDomains" v-loading="loading" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" />
        <el-table-column label="名称" prop="name" width="200">
          <template #default="{ row }">
            <el-link @click="viewDetail(row)">{{ row.name }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="tags" label="标签" width="150">
          <template #default="{ row }">
            <el-tag v-for="tag in row.tags" :key="tag.key" size="small" class="tag-item">{{ tag.key }}: {{ tag.value }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="area" label="加速区域" width="150" />
        <el-table-column prop="cname" label="CNAME" width="200" />
        <el-table-column prop="service_type" label="服务类型" width="120" />
        <el-table-column prop="provider_type" label="平台" width="100" />
        <el-table-column prop="account_name" label="云账号" width="150" />
        <el-table-column prop="project_name" label="项目" width="120" />
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="refreshPreheat(row)">刷新预热</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="searchParams.page"
        v-model:page-size="searchParams.page_size"
        :total="total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next"
        @size-change="loadCDNDomains"
        @current-change="loadCDNDomains"
      />
    </el-card>

    <!-- Detail Modal -->
    <el-dialog v-model="detailDialogVisible" title="CDN域名详情" width="600px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ selectedCDNDomain?.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ selectedCDNDomain?.name }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ selectedCDNDomain?.status }}</el-descriptions-item>
        <el-descriptions-item label="加速区域">{{ selectedCDNDomain?.area }}</el-descriptions-item>
        <el-descriptions-item label="CNAME">{{ selectedCDNDomain?.cname }}</el-descriptions-item>
        <el-descriptions-item label="服务类型">{{ selectedCDNDomain?.service_type }}</el-descriptions-item>
        <el-descriptions-item label="平台">{{ selectedCDNDomain?.provider_type }}</el-descriptions-item>
        <el-descriptions-item label="云账号">{{ selectedCDNDomain?.account_name }}</el-descriptions-item>
        <el-descriptions-item label="项目">{{ selectedCDNDomain?.project_name }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ selectedCDNDomain?.created_at }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- Create Modal -->
    <el-dialog v-model="createDialogVisible" title="创建CDN域名" width="500px">
      <el-form :model="createForm" label-width="100px">
        <el-form-item label="名称">
          <el-input v-model="createForm.name" placeholder="请输入域名" />
        </el-form-item>
        <el-form-item label="加速区域">
          <el-select v-model="createForm.area" placeholder="请选择加速区域">
            <el-option label="中国大陆" value="china" />
            <el-option label="全球" value="global" />
            <el-option label="亚太地区" value="asia" />
          </el-select>
        </el-form-item>
        <el-form-item label="服务类型">
          <el-select v-model="createForm.service_type" placeholder="请选择服务类型">
            <el-option label="全站加速" value="full" />
            <el-option label="静态加速" value="static" />
            <el-option label="视频点播" value="video" />
          </el-select>
        </el-form-item>
        <el-form-item label="源站地址">
          <el-input v-model="createForm.origin_address" placeholder="请输入源站地址" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitCreate">创建</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowDown } from '@element-plus/icons-vue'
import { getCDNDomains, createCDNDomain, deleteCDNDomain, batchDeleteCDNDomains, syncCDNStatus, batchSyncCDNStatus, type CDNDomain } from '@/api/networkSync'

const cdnDomains = ref<CDNDomain[]>([])
const loading = ref(false)
const total = ref(0)
const selectedRows = ref<CDNDomain[]>([])
const detailDialogVisible = ref(false)
const createDialogVisible = ref(false)
const selectedCDNDomain = ref<CDNDomain | null>(null)

const searchParams = ref({
  page: 1,
  page_size: 10,
  name: '',
  status: '',
  area: '',
  platform: ''
})

const createForm = ref({
  name: '',
  area: '',
  service_type: '',
  origin_address: ''
})

const hasSelection = computed(() => selectedRows.value.length > 0)

const getStatusType = (status: string) => {
  switch (status.toLowerCase()) {
    case 'online':
    case 'active':
      return 'success'
    case 'pending':
    case 'configuring':
      return 'warning'
    case 'offline':
    case 'error':
      return 'danger'
    default:
      return 'info'
  }
}

const loadCDNDomains = async () => {
  loading.value = true
  try {
    const res = await getCDNDomains(searchParams.value)
    cdnDomains.value = res.items || []
    total.value = res.total || 0
  } catch (e) {
    console.error(e)
    ElMessage.error('加载失败')
    cdnDomains.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

const handleSelectionChange = (rows: CDNDomain[]) => {
  selectedRows.value = rows
}

const resetSearch = () => {
  searchParams.value = {
    page: 1,
    page_size: 10,
    name: '',
    status: '',
    area: '',
    platform: ''
  }
  loadCDNDomains()
}

const handleView = () => {
  if (selectedRows.value.length === 1) {
    viewDetail(selectedRows.value[0])
  } else {
    ElMessage.warning('请选择一条记录查看')
  }
}

const viewDetail = (row: CDNDomain) => {
  selectedCDNDomain.value = row
  detailDialogVisible.value = true
}

const handleCreate = () => {
  createForm.value = { name: '', area: '', service_type: '', origin_address: '' }
  createDialogVisible.value = true
}

const submitCreate = async () => {
  try {
    await createCDNDomain(createForm.value)
    ElMessage.success('创建成功')
    createDialogVisible.value = false
    loadCDNDomains()
  } catch (e) {
    console.error(e)
    ElMessage.error('创建失败')
  }
}

const handleSyncStatus = async () => {
  if (!hasSelection.value) return
  try {
    await ElMessageBox.confirm('确认同步选中域名的状态？', '提示', { type: 'info' })
    const ids = selectedRows.value.map(r => Number(r.id))
    await batchSyncCDNStatus(ids)
    ElMessage.success('同步成功')
    loadCDNDomains()
  } catch (e) {
    console.error(e)
  }
}

const handleBatchCommand = async (command: string) => {
  if (!hasSelection.value) return
  if (command === 'delete') {
    handleBatchDelete()
  } else if (command === 'sync') {
    handleSyncStatus()
  }
}

const handleTags = () => {
  if (!hasSelection.value) {
    ElMessage.warning('请选择需要设置标签的域名')
    return
  }
  ElMessage.info('标签功能开发中')
}

const refreshPreheat = (row: CDNDomain) => {
  ElMessage.info(`刷新预热CDN域名: ${row.name}`)
}

const handleBatchDelete = async () => {
  if (!hasSelection.value) return
  try {
    await ElMessageBox.confirm(`确认删除选中的 ${selectedRows.value.length} 个域名？`, '警告', { type: 'warning' })
    const ids = selectedRows.value.map(r => Number(r.id))
    await batchDeleteCDNDomains(ids)
    ElMessage.success('删除成功')
    loadCDNDomains()
  } catch (e) {
    console.error(e)
  }
}

const handleDelete = async (row: CDNDomain) => {
  try {
    await ElMessageBox.confirm(`确认删除CDN域名 ${row.name}？`, '警告', { type: 'warning' })
    await deleteCDNDomain(row.id)
    ElMessage.success('删除成功')
    loadCDNDomains()
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadCDNDomains()
})
</script>

<style scoped>
.cdn-domains-page {
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

.toolbar {
  display: flex;
  gap: 8px;
}

.search-form {
  margin-bottom: 16px;
}

.tag-item {
  margin-right: 4px;
}
</style>