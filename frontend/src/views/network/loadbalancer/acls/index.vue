<template>
  <div class="lb-acls-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">访问控制</span>
          <div class="toolbar">
            <el-button @click="handleView">查看</el-button>
            <el-button type="primary" @click="handleCreate">创建</el-button>
            <el-button :disabled="!hasSelection" @click="handleBatchDelete">删除</el-button>
          </div>
        </div>
      </template>

      <!-- Search Filters -->
      <el-form :inline="true" class="search-form">
        <el-form-item label="名称">
          <el-input v-model="searchParams.name" placeholder="请输入名称" clearable @clear="loadACLs" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchParams.status" placeholder="请选择状态" clearable @clear="loadACLs">
            <el-option label="启用" value="active" />
            <el-option label="禁用" value="disabled" />
          </el-select>
        </el-form-item>
        <el-form-item label="平台">
          <el-select v-model="searchParams.platform" placeholder="请选择平台" clearable @clear="loadACLs">
            <el-option label="阿里云" value="aliyun" />
            <el-option label="腾讯云" value="tencent" />
            <el-option label="AWS" value="aws" />
            <el-option label="Azure" value="azure" />
          </el-select>
        </el-form-item>
        <el-form-item label="区域">
          <el-input v-model="searchParams.region" placeholder="请输入区域" clearable @clear="loadACLs" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadACLs">查询</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="acls" v-loading="loading" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" />
        <el-table-column label="名称" prop="name" width="180">
          <template #default="{ row }">
            <el-link @click="viewDetail(row)">{{ row.name }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="address_source" label="源地址 | 备注" width="200">
          <template #default="{ row }">
            {{ row.address_source }} {{ row.remarks ? `(${row.remarks})` : '' }}
          </template>
        </el-table-column>
        <el-table-column prop="listener_count" label="监听器" width="100">
          <template #default="{ row }">
            {{ row.listener_count }} 个
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="provider_type" label="平台" width="100" />
        <el-table-column prop="account_name" label="云账号" width="150" />
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column prop="updated_at" label="更新时间" width="180" />
        <el-table-column prop="shared_scope" label="共享范围" width="100" />
        <el-table-column prop="region" label="区域" width="120" />
        <el-table-column prop="project_name" label="项目" width="120" />
        <el-table-column label="操作" width="150" fixed="right">
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
        @size-change="loadACLs"
        @current-change="loadACLs"
      />
    </el-card>

    <!-- Detail Modal -->
    <el-dialog v-model="detailDialogVisible" title="访问控制详情" width="600px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ selectedACL?.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ selectedACL?.name }}</el-descriptions-item>
        <el-descriptions-item label="源地址">{{ selectedACL?.address_source }}</el-descriptions-item>
        <el-descriptions-item label="备注">{{ selectedACL?.remarks }}</el-descriptions-item>
        <el-descriptions-item label="监听器">{{ selectedACL?.listener_count }} 个</el-descriptions-item>
        <el-descriptions-item label="状态">{{ selectedACL?.status }}</el-descriptions-item>
        <el-descriptions-item label="平台">{{ selectedACL?.provider_type }}</el-descriptions-item>
        <el-descriptions-item label="云账号">{{ selectedACL?.account_name }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ selectedACL?.created_at }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ selectedACL?.updated_at }}</el-descriptions-item>
        <el-descriptions-item label="共享范围">{{ selectedACL?.shared_scope }}</el-descriptions-item>
        <el-descriptions-item label="区域">{{ selectedACL?.region }}</el-descriptions-item>
        <el-descriptions-item label="项目">{{ selectedACL?.project_name }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- Create Modal -->
    <el-dialog v-model="createDialogVisible" title="创建访问控制策略" width="500px">
      <el-form :model="createForm" label-width="100px">
        <el-form-item label="名称">
          <el-input v-model="createForm.name" placeholder="请输入名称" />
        </el-form-item>
        <el-form-item label="源地址">
          <el-input v-model="createForm.address_source" placeholder="请输入IP地址段" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="createForm.remarks" placeholder="请输入备注" />
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
import { getLBACLs, createLBACL, deleteLBACL, batchDeleteLBACLs, type LBACL } from '@/api/networkSync'

const acls = ref<LBACL[]>([])
const loading = ref(false)
const total = ref(0)
const selectedRows = ref<LBACL[]>([])
const detailDialogVisible = ref(false)
const createDialogVisible = ref(false)
const selectedACL = ref<LBACL | null>(null)

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
  address_source: '',
  remarks: ''
})

const hasSelection = computed(() => selectedRows.value.length > 0)

const getStatusType = (status: string) => {
  switch (status.toLowerCase()) {
    case 'active':
    case 'enabled':
      return 'success'
    case 'pending':
      return 'warning'
    case 'disabled':
    case 'error':
      return 'danger'
    default:
      return 'info'
  }
}

const loadACLs = async () => {
  loading.value = true
  try {
    const res = await getLBACLs(searchParams.value)
    acls.value = res.items || []
    total.value = res.total || 0
  } catch (e) {
    console.error(e)
    ElMessage.error('加载失败')
    acls.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

const handleSelectionChange = (rows: LBACL[]) => {
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
  loadACLs()
}

const handleView = () => {
  if (selectedRows.value.length === 1) {
    viewDetail(selectedRows.value[0])
  } else {
    ElMessage.warning('请选择一条记录查看')
  }
}

const viewDetail = (row: LBACL) => {
  selectedACL.value = row
  detailDialogVisible.value = true
}

const handleCreate = () => {
  createForm.value = { name: '', address_source: '', remarks: '' }
  createDialogVisible.value = true
}

const submitCreate = async () => {
  try {
    await createLBACL(createForm.value)
    ElMessage.success('创建成功')
    createDialogVisible.value = false
    loadACLs()
  } catch (e) {
    console.error(e)
    ElMessage.error('创建失败')
  }
}

const handleBatchDelete = async () => {
  if (!hasSelection.value) return
  try {
    await ElMessageBox.confirm(`确认删除选中的 ${selectedRows.value.length} 个策略？`, '警告', { type: 'warning' })
    const ids = selectedRows.value.map(r => Number(r.id))
    await batchDeleteLBACLs(ids)
    ElMessage.success('删除成功')
    loadACLs()
  } catch (e) {
    console.error(e)
  }
}

const handleEdit = (row: LBACL) => {
  ElMessage.info(`编辑访问控制策略: ${row.name}`)
}

const handleDelete = async (row: LBACL) => {
  try {
    await ElMessageBox.confirm(`确认删除访问控制策略 ${row.name}？`, '警告', { type: 'warning' })
    await deleteLBACL(row.id)
    ElMessage.success('删除成功')
    loadACLs()
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadACLs()
})
</script>

<style scoped>
.lb-acls-page {
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