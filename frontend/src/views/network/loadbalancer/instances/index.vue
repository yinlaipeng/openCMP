<template>
  <div class="lb-instances-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">负载均衡实例</span>
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
          <el-input v-model="searchParams.name" placeholder="请输入名称" clearable @clear="loadInstances" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchParams.status" placeholder="请选择状态" clearable @clear="loadInstances">
            <el-option label="运行中" value="active" />
            <el-option label="已停止" value="stopped" />
            <el-option label="创建中" value="creating" />
            <el-option label="错误" value="error" />
          </el-select>
        </el-form-item>
        <el-form-item label="平台">
          <el-select v-model="searchParams.platform" placeholder="请选择平台" clearable @clear="loadInstances">
            <el-option label="阿里云" value="aliyun" />
            <el-option label="腾讯云" value="tencent" />
            <el-option label="AWS" value="aws" />
            <el-option label="Azure" value="azure" />
          </el-select>
        </el-form-item>
        <el-form-item label="区域">
          <el-input v-model="searchParams.region" placeholder="请输入区域" clearable @clear="loadInstances" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadInstances">查询</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>

      <el-tabs v-model="activeTab" @tab-change="handleTabChange">
        <el-tab-pane label="全部" name="all">
          <el-table :data="instances" v-loading="loading" @selection-change="handleSelectionChange">
            <el-table-column type="selection" width="55" />
            <el-table-column label="名称" prop="name" width="180">
              <template #default="{ row }">
                <el-link @click="viewDetail(row)">{{ row.name }}</el-link>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="getStatusType(row.status)">{{ row.status }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="address" label="服务地址" width="180" />
            <el-table-column prop="specification" label="规格" width="120" />
            <el-table-column prop="vpc" label="VPC" width="180" />
            <el-table-column prop="security_group" label="安全组" width="150" />
            <el-table-column prop="charging_method" label="计费方式" width="120" />
            <el-table-column prop="provider_type" label="平台" width="100" />
            <el-table-column prop="project_name" label="项目" width="120" />
            <el-table-column prop="region" label="区域" width="120" />
            <el-table-column label="操作" width="180" fixed="right">
              <template #default="{ row }">
                <el-button size="small" @click="manageListeners(row)">管理监听</el-button>
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
            @size-change="loadInstances"
            @current-change="loadInstances"
          />
        </el-tab-pane>
        <el-tab-pane label="idc" name="idc">
          <el-table :data="instances" v-loading="loading" @selection-change="handleSelectionChange">
            <el-table-column type="selection" width="55" />
            <el-table-column label="名称" prop="name" width="180">
              <template #default="{ row }">
                <el-link @click="viewDetail(row)">{{ row.name }}</el-link>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="getStatusType(row.status)">{{ row.status }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="address" label="服务地址" width="180" />
            <el-table-column prop="specification" label="规格" width="120" />
            <el-table-column prop="vpc" label="VPC" width="180" />
            <el-table-column prop="security_group" label="安全组" width="150" />
            <el-table-column prop="charging_method" label="计费方式" width="120" />
            <el-table-column prop="provider_type" label="平台" width="100" />
            <el-table-column prop="project_name" label="项目" width="120" />
            <el-table-column prop="region" label="区域" width="120" />
            <el-table-column label="操作" width="180" fixed="right">
              <template #default="{ row }">
                <el-button size="small" @click="manageListeners(row)">管理监听</el-button>
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
            @size-change="loadInstances"
            @current-change="loadInstances"
          />
        </el-tab-pane>
        <el-tab-pane label="公有云" name="public_cloud">
          <el-table :data="instances" v-loading="loading" @selection-change="handleSelectionChange">
            <el-table-column type="selection" width="55" />
            <el-table-column label="名称" prop="name" width="180">
              <template #default="{ row }">
                <el-link @click="viewDetail(row)">{{ row.name }}</el-link>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="getStatusType(row.status)">{{ row.status }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="address" label="服务地址" width="180" />
            <el-table-column prop="specification" label="规格" width="120" />
            <el-table-column prop="vpc" label="VPC" width="180" />
            <el-table-column prop="security_group" label="安全组" width="150" />
            <el-table-column prop="charging_method" label="计费方式" width="120" />
            <el-table-column prop="provider_type" label="平台" width="100" />
            <el-table-column prop="project_name" label="项目" width="120" />
            <el-table-column prop="region" label="区域" width="120" />
            <el-table-column label="操作" width="180" fixed="right">
              <template #default="{ row }">
                <el-button size="small" @click="manageListeners(row)">管理监听</el-button>
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
            @size-change="loadInstances"
            @current-change="loadInstances"
          />
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- Detail Modal -->
    <el-dialog v-model="detailDialogVisible" title="负载均衡实例详情" width="600px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ selectedInstance?.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ selectedInstance?.name }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ selectedInstance?.status }}</el-descriptions-item>
        <el-descriptions-item label="服务地址">{{ selectedInstance?.address }}</el-descriptions-item>
        <el-descriptions-item label="规格">{{ selectedInstance?.specification }}</el-descriptions-item>
        <el-descriptions-item label="VPC">{{ selectedInstance?.vpc }}</el-descriptions-item>
        <el-descriptions-item label="安全组">{{ selectedInstance?.security_group }}</el-descriptions-item>
        <el-descriptions-item label="计费方式">{{ selectedInstance?.charging_method }}</el-descriptions-item>
        <el-descriptions-item label="平台">{{ selectedInstance?.provider_type }}</el-descriptions-item>
        <el-descriptions-item label="项目">{{ selectedInstance?.project_name }}</el-descriptions-item>
        <el-descriptions-item label="区域">{{ selectedInstance?.region }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ selectedInstance?.created_at }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- Listeners Modal -->
    <el-dialog v-model="listenersDialogVisible" title="管理监听规则" width="800px">
      <el-table :data="listeners" v-loading="listenersLoading">
        <el-table-column prop="name" label="监听名称" width="150" />
        <el-table-column prop="protocol" label="协议" width="100" />
        <el-table-column prop="port" label="端口" width="100" />
        <el-table-column prop="backend_port" label="后端端口" width="100" />
        <el-table-column prop="scheduler" label="调度算法" width="120" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'Active' ? 'success' : 'info'">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-button size="small" type="danger" @click="deleteListener(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <template #footer>
        <el-button @click="listenersDialogVisible = false">关闭</el-button>
        <el-button type="primary" @click="addListener">添加监听</el-button>
      </template>
    </el-dialog>

    <!-- Create Modal -->
    <el-dialog v-model="createDialogVisible" title="创建负载均衡实例" width="500px">
      <el-form :model="createForm" label-width="100px">
        <el-form-item label="名称">
          <el-input v-model="createForm.name" placeholder="请输入名称" />
        </el-form-item>
        <el-form-item label="规格">
          <el-select v-model="createForm.specification" placeholder="请选择规格">
            <el-option label="性能保障型" value="performance" />
            <el-option label="标准型" value="standard" />
          </el-select>
        </el-form-item>
        <el-form-item label="VPC">
          <el-input v-model="createForm.vpc_id" placeholder="请输入VPC ID" />
        </el-form-item>
        <el-form-item label="计费方式">
          <el-select v-model="createForm.charging_method" placeholder="请选择计费方式">
            <el-option label="按量付费" value="postpaid" />
            <el-option label="包年包月" value="prepaid" />
          </el-select>
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
import { getLBInstances, createLBInstance, deleteLBInstance, batchDeleteLBInstances, syncLBStatus, type LBInstance } from '@/api/networkSync'

interface Listener {
  name: string
  protocol: string
  port: number
  backend_port: number
  scheduler: string
  status: string
}

const instances = ref<LBInstance[]>([])
const loading = ref(false)
const activeTab = ref('all')
const total = ref(0)
const selectedRows = ref<LBInstance[]>([])
const detailDialogVisible = ref(false)
const listenersDialogVisible = ref(false)
const createDialogVisible = ref(false)
const listenersLoading = ref(false)
const selectedInstance = ref<LBInstance | null>(null)
const listeners = ref<Listener[]>([])

const searchParams = ref({
  page: 1,
  page_size: 10,
  name: '',
  status: '',
  platform: '',
  region: ''
})

const createForm = ref({
  name: '',
  specification: '',
  vpc_id: '',
  charging_method: ''
})

const hasSelection = computed(() => selectedRows.value.length > 0)

const getStatusType = (status: string) => {
  switch (status.toLowerCase()) {
    case 'active':
    case 'running':
      return 'success'
    case 'pending':
    case 'creating':
      return 'warning'
    case 'error':
    case 'stopped':
      return 'danger'
    default:
      return 'info'
  }
}

const loadInstances = async () => {
  loading.value = true
  try {
    const params = { ...searchParams.value }
    if (activeTab.value === 'idc') {
      params.platform = 'idc'
    } else if (activeTab.value === 'public_cloud') {
      params.platform = searchParams.value.platform || ''
    }
    const res = await getLBInstances(params)
    instances.value = res.items || []
    total.value = res.total || 0
  } catch (e) {
    console.error(e)
    ElMessage.error('加载失败')
    instances.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

const handleSelectionChange = (rows: LBInstance[]) => {
  selectedRows.value = rows
}

const handleTabChange = () => {
  searchParams.value.page = 1
  loadInstances()
}

const resetSearch = () => {
  searchParams.value = {
    page: 1,
    page_size: 10,
    name: '',
    status: '',
    platform: '',
    region: ''
  }
  loadInstances()
}

const handleView = () => {
  if (selectedRows.value.length === 1) {
    viewDetail(selectedRows.value[0])
  } else {
    ElMessage.warning('请选择一条记录查看')
  }
}

const viewDetail = (row: LBInstance) => {
  selectedInstance.value = row
  detailDialogVisible.value = true
}

const handleCreate = () => {
  createForm.value = { name: '', specification: '', vpc_id: '', charging_method: '' }
  createDialogVisible.value = true
}

const submitCreate = async () => {
  try {
    await createLBInstance(createForm.value)
    ElMessage.success('创建成功')
    createDialogVisible.value = false
    loadInstances()
  } catch (e) {
    console.error(e)
    ElMessage.error('创建失败')
  }
}

const handleSyncStatus = async () => {
  if (!hasSelection.value) return
  try {
    await ElMessageBox.confirm('确认同步选中实例的状态？', '提示', { type: 'info' })
    const ids = selectedRows.value.map(r => Number(r.id))
    await batchSyncLBStatus(ids)
    ElMessage.success('同步成功')
    loadInstances()
  } catch (e) {
    console.error(e)
  }
}

const handleBatchCommand = async (command: string) => {
  if (!hasSelection.value) return
  if (command === 'delete') {
    try {
      await ElMessageBox.confirm(`确认删除选中的 ${selectedRows.value.length} 个实例？`, '警告', { type: 'warning' })
      const ids = selectedRows.value.map(r => Number(r.id))
      await batchDeleteLBInstances(ids)
      ElMessage.success('删除成功')
      loadInstances()
    } catch (e) {
      console.error(e)
    }
  } else if (command === 'sync') {
    handleSyncStatus()
  }
}

const handleTags = () => {
  if (!hasSelection.value) {
    ElMessage.warning('请选择需要设置标签的实例')
    return
  }
  ElMessage.info('标签功能开发中')
}

const manageListeners = (row: LBInstance) => {
  selectedInstance.value = row
  listenersDialogVisible.value = true
  listenersLoading.value = true
  listeners.value = [
    { name: 'HTTP监听', protocol: 'HTTP', port: 80, backend_port: 8080, scheduler: '轮询', status: 'Active' },
    { name: 'HTTPS监听', protocol: 'HTTPS', port: 443, backend_port: 8443, scheduler: '轮询', status: 'Active' }
  ]
  listenersLoading.value = false
}

const addListener = () => {
  ElMessage.info('添加监听功能开发中')
}

const deleteListener = async (row: Listener) => {
  try {
    await ElMessageBox.confirm(`确认删除监听 ${row.name}？`, '警告', { type: 'warning' })
    listeners.value = listeners.value.filter(l => l.name !== row.name)
    ElMessage.success('删除成功')
  } catch (e) {
    console.error(e)
  }
}

const handleDelete = async (row: LBInstance) => {
  try {
    await ElMessageBox.confirm(`确认删除负载均衡实例 ${row.name}？`, '警告', { type: 'warning' })
    await deleteLBInstance(row.id)
    ElMessage.success('删除成功')
    loadInstances()
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadInstances()
})
</script>

<style scoped>
.lb-instances-page {
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
</style>