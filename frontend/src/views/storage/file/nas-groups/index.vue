<template>
  <div class="nas-groups-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">NAS权限组列表</span>
          <div class="toolbar">
            <el-button @click="handleView">查看</el-button>
            <el-button type="primary" @click="handleCreate">创建</el-button>
            <el-dropdown :disabled="!hasSelection" @command="handleBatchCommand">
              <el-button>
                批量操作 <el-icon class="el-icon--right"><arrow-down /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="delete">批量删除</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </div>
      </template>

      <!-- Search Filters -->
      <el-form :inline="true" class="search-form">
        <el-form-item label="名称">
          <el-input v-model="searchParams.name" placeholder="请输入名称" clearable @clear="loadNASGroups" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchParams.status" placeholder="请选择状态" clearable @clear="loadNASGroups">
            <el-option label="正常" value="active" />
            <el-option label="禁用" value="disabled" />
          </el-select>
        </el-form-item>
        <el-form-item label="平台">
          <el-select v-model="searchParams.platform" placeholder="请选择平台" clearable @clear="loadNASGroups">
            <el-option label="阿里云" value="aliyun" />
            <el-option label="腾讯云" value="tencent" />
            <el-option label="AWS" value="aws" />
            <el-option label="Azure" value="azure" />
          </el-select>
        </el-form-item>
        <el-form-item label="区域">
          <el-input v-model="searchParams.region" placeholder="请输入区域" clearable @clear="loadNASGroups" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadNASGroups">查询</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="nasGroups" v-loading="loading" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" />
        <el-table-column label="名称" width="180">
          <template #default="{ row }">
            <el-link @click="viewDetail(row)">{{ row.name }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="mount_count" label="挂载点数量" width="100">
          <template #default="{ row }">
            {{ row.mount_count }} 个
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="platform" label="平台" width="100" />
        <el-table-column prop="account" label="云账号" width="150" />
        <el-table-column prop="region" label="区域" width="120" />
        <el-table-column prop="domain" label="所属域" width="150" />
        <el-table-column label="操作" width="150">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
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
        @size-change="loadNASGroups"
        @current-change="loadNASGroups"
      />
    </el-card>

    <!-- Detail Modal -->
    <el-dialog v-model="detailDialogVisible" title="NAS权限组详情" width="600px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ selectedNASGroup?.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ selectedNASGroup?.name }}</el-descriptions-item>
        <el-descriptions-item label="挂载点数量">{{ selectedNASGroup?.mount_count }} 个</el-descriptions-item>
        <el-descriptions-item label="状态">{{ selectedNASGroup?.status }}</el-descriptions-item>
        <el-descriptions-item label="平台">{{ selectedNASGroup?.platform }}</el-descriptions-item>
        <el-descriptions-item label="云账号">{{ selectedNASGroup?.account }}</el-descriptions-item>
        <el-descriptions-item label="区域">{{ selectedNASGroup?.region }}</el-descriptions-item>
        <el-descriptions-item label="所属域">{{ selectedNASGroup?.domain }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ selectedNASGroup?.created_at }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- Create Modal -->
    <el-dialog v-model="createDialogVisible" title="创建NAS权限组" width="500px">
      <el-form :model="createForm" label-width="100px">
        <el-form-item label="名称">
          <el-input v-model="createForm.name" placeholder="请输入NAS权限组名称" />
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

interface NASGroup {
  id: string
  name: string
  mount_count: number
  status: string
  platform: string
  account: string
  region: string
  domain: string
  created_at: string
}

const nasGroups = ref<NASGroup[]>([])
const loading = ref(false)
const total = ref(0)
const selectedRows = ref<NASGroup[]>([])
const detailDialogVisible = ref(false)
const createDialogVisible = ref(false)
const selectedNASGroup = ref<NASGroup | null>(null)

const searchParams = ref({
  page: 1,
  page_size: 10,
  name: '',
  status: '',
  platform: '',
  region: ''
})

const createForm = ref({
  name: ''
})

const hasSelection = computed(() => selectedRows.value.length > 0)

const handleSelectionChange = (rows: NASGroup[]) => {
  selectedRows.value = rows
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
  loadNASGroups()
}

const handleView = () => {
  if (selectedRows.value.length === 1) {
    viewDetail(selectedRows.value[0])
  } else {
    ElMessage.warning('请选择一条记录查看')
  }
}

const handleCreate = () => {
  createForm.value = { name: '' }
  createDialogVisible.value = true
}

const submitCreate = async () => {
  try {
    ElMessage.success('创建成功')
    createDialogVisible.value = false
    loadNASGroups()
  } catch (e) {
    console.error(e)
    ElMessage.error('创建失败')
  }
}

const handleBatchCommand = async (command: string) => {
  if (!hasSelection.value) return
  if (command === 'delete') {
    try {
      await ElMessageBox.confirm(`确认删除选中的 ${selectedRows.value.length} 个NAS权限组？`, '警告', { type: 'warning' })
      nasGroups.value = nasGroups.value.filter(n => !selectedRows.value.some(r => r.id === n.id))
      ElMessage.success('删除成功')
      loadNASGroups()
    } catch (e) {
      console.error(e)
    }
  }
}

const getStatusType = (status: string) => {
  switch (status.toLowerCase()) {
    case 'active':
    case 'normal':
      return 'success'
    case 'pending':
    case 'creating':
      return 'warning'
    case 'error':
    case 'disabled':
      return 'danger'
    default:
      return 'info'
  }
}

const loadNASGroups = async () => {
  loading.value = true
  try {
    // Mock data
    nasGroups.value = [
      {
        id: 'nas-group-1',
        name: 'NAS权限组 1',
        mount_count: 3,
        status: 'Active',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        region: 'cn-hangzhou',
        domain: 'Default Domain',
        created_at: '2024-01-01 10:00:00'
      },
      {
        id: 'nas-group-2',
        name: 'NAS权限组 2',
        mount_count: 1,
        status: 'Active',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        region: 'cn-shanghai',
        domain: 'Default Domain',
        created_at: '2024-01-02 10:00:00'
      },
      {
        id: 'nas-group-3',
        name: 'NAS权限组 3',
        mount_count: 0,
        status: 'Disabled',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        region: 'cn-beijing',
        domain: 'Domain A',
        created_at: '2024-01-03 10:00:00'
      }
    ]
    total.value = 3
  } catch (e) {
    console.error(e)
    nasGroups.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

const viewDetail = (row: NASGroup) => {
  selectedNASGroup.value = row
  detailDialogVisible.value = true
}

const handleEdit = (row: NASGroup) => {
  ElMessage.info(`编辑NAS权限组: ${row.name}`)
}

const handleDelete = async (row: NASGroup) => {
  try {
    await ElMessageBox.confirm(`确认删除NAS权限组 ${row.name}？`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    nasGroups.value = nasGroups.value.filter(n => n.id !== row.id)
    ElMessage.success('删除成功')
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadNASGroups()
})
</script>

<style scoped>
.nas-groups-page {
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