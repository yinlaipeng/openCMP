<template>
  <div class="lb-certificates-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">证书</span>
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
            <el-button :disabled="!hasSelection" @click="handleTags">标签</el-button>
            <el-button :disabled="!hasSelection" @click="handleChangeProject">更改项目</el-button>
            <el-button :disabled="!hasSelection" @click="handleSetupSharing">设置共享</el-button>
            <el-button :disabled="!hasSelection" type="danger" @click="handleBatchDelete">删除</el-button>
          </div>
        </div>
      </template>

      <!-- Search Filters -->
      <el-form :inline="true" class="search-form">
        <el-form-item label="名称">
          <el-input v-model="searchParams.name" placeholder="请输入名称" clearable @clear="loadCertificates" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchParams.status" placeholder="请选择状态" clearable @clear="loadCertificates">
            <el-option label="正常" value="normal" />
            <el-option label="即将过期" value="expiring" />
            <el-option label="已过期" value="expired" />
          </el-select>
        </el-form-item>
        <el-form-item label="平台">
          <el-select v-model="searchParams.platform" placeholder="请选择平台" clearable @clear="loadCertificates">
            <el-option label="阿里云" value="aliyun" />
            <el-option label="腾讯云" value="tencent" />
            <el-option label="AWS" value="aws" />
            <el-option label="Azure" value="azure" />
          </el-select>
        </el-form-item>
        <el-form-item label="区域">
          <el-input v-model="searchParams.region" placeholder="请输入区域" clearable @clear="loadCertificates" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadCertificates">查询</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="certificates" v-loading="loading" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" />
        <el-table-column label="名称" prop="name" width="180">
          <template #default="{ row }">
            <el-link @click="viewDetail(row)">{{ row.name }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="tags" label="标签" width="150">
          <template #default="{ row }">
            <el-tag v-for="tag in row.tags" :key="tag.key" size="small" class="tag-item">{{ tag.key }}: {{ tag.value }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="domain_name" label="域名" width="180" />
        <el-table-column prop="expiration" label="过期时间" width="180" />
        <el-table-column prop="subject_alternative_names" label="扩展域名" width="200" />
        <el-table-column prop="listener_count" label="监听器" width="100">
          <template #default="{ row }">
            {{ row.listener_count }} 个
          </template>
        </el-table-column>
        <el-table-column prop="shared_scope" label="共享范围" width="100" />
        <el-table-column prop="project_name" label="项目" width="120" />
        <el-table-column prop="provider_type" label="平台" width="100" />
        <el-table-column prop="account_name" label="云账号" width="150" />
        <el-table-column prop="region" label="区域" width="120" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleChangeProjectSingle(row)">更改项目</el-button>
            <el-button size="small" @click="handleSetupSharingSingle(row)">设置共享</el-button>
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
        @size-change="loadCertificates"
        @current-change="loadCertificates"
      />
    </el-card>

    <!-- Detail Modal -->
    <el-dialog v-model="detailDialogVisible" title="证书详情" width="600px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ selectedCertificate?.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ selectedCertificate?.name }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ selectedCertificate?.status }}</el-descriptions-item>
        <el-descriptions-item label="域名">{{ selectedCertificate?.domain_name }}</el-descriptions-item>
        <el-descriptions-item label="过期时间">{{ selectedCertificate?.expiration }}</el-descriptions-item>
        <el-descriptions-item label="扩展域名">{{ selectedCertificate?.subject_alternative_names }}</el-descriptions-item>
        <el-descriptions-item label="监听器">{{ selectedCertificate?.listener_count }} 个</el-descriptions-item>
        <el-descriptions-item label="共享范围">{{ selectedCertificate?.shared_scope }}</el-descriptions-item>
        <el-descriptions-item label="项目">{{ selectedCertificate?.project_name }}</el-descriptions-item>
        <el-descriptions-item label="平台">{{ selectedCertificate?.provider_type }}</el-descriptions-item>
        <el-descriptions-item label="云账号">{{ selectedCertificate?.account_name }}</el-descriptions-item>
        <el-descriptions-item label="区域">{{ selectedCertificate?.region }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ selectedCertificate?.created_at }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- Create Modal -->
    <el-dialog v-model="createDialogVisible" title="创建证书" width="500px">
      <el-form :model="createForm" label-width="100px">
        <el-form-item label="名称">
          <el-input v-model="createForm.name" placeholder="请输入名称" />
        </el-form-item>
        <el-form-item label="域名">
          <el-input v-model="createForm.domain_name" placeholder="请输入域名" />
        </el-form-item>
        <el-form-item label="证书内容">
          <el-input v-model="createForm.cert_content" type="textarea" rows="5" placeholder="请输入证书内容" />
        </el-form-item>
        <el-form-item label="私钥内容">
          <el-input v-model="createForm.private_key" type="textarea" rows="5" placeholder="请输入私钥内容" />
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
import { getLBCertificates, createLBCertificate, deleteLBCertificate, batchDeleteLBCertificates, changeCertificateProject, setupCertificateSharing, type LBCertificate } from '@/api/networkSync'

const certificates = ref<LBCertificate[]>([])
const loading = ref(false)
const total = ref(0)
const selectedRows = ref<LBCertificate[]>([])
const detailDialogVisible = ref(false)
const createDialogVisible = ref(false)
const selectedCertificate = ref<LBCertificate | null>(null)

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
  domain_name: '',
  cert_content: '',
  private_key: ''
})

const hasSelection = computed(() => selectedRows.value.length > 0)

const getStatusType = (status: string) => {
  switch (status.toLowerCase()) {
    case 'normal':
    case 'active':
      return 'success'
    case 'expiring':
      return 'warning'
    case 'expired':
    case 'error':
      return 'danger'
    default:
      return 'info'
  }
}

const loadCertificates = async () => {
  loading.value = true
  try {
    const res = await getLBCertificates(searchParams.value)
    certificates.value = res.items || []
    total.value = res.total || 0
  } catch (e) {
    console.error(e)
    ElMessage.error('加载失败')
    certificates.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

const handleSelectionChange = (rows: LBCertificate[]) => {
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
  loadCertificates()
}

const handleView = () => {
  if (selectedRows.value.length === 1) {
    viewDetail(selectedRows.value[0])
  } else {
    ElMessage.warning('请选择一条记录查看')
  }
}

const viewDetail = (row: LBCertificate) => {
  selectedCertificate.value = row
  detailDialogVisible.value = true
}

const handleCreate = () => {
  createForm.value = { name: '', domain_name: '', cert_content: '', private_key: '' }
  createDialogVisible.value = true
}

const submitCreate = async () => {
  try {
    await createLBCertificate(createForm.value)
    ElMessage.success('创建成功')
    createDialogVisible.value = false
    loadCertificates()
  } catch (e) {
    console.error(e)
    ElMessage.error('创建失败')
  }
}

const handleBatchCommand = async (command: string) => {
  if (!hasSelection.value) return
  if (command === 'delete') {
    handleBatchDelete()
  }
}

const handleTags = () => {
  if (!hasSelection.value) {
    ElMessage.warning('请选择需要设置标签的证书')
    return
  }
  ElMessage.info('标签功能开发中')
}

const handleChangeProject = async () => {
  if (!hasSelection.value) return
  ElMessage.info('批量更改项目功能开发中')
}

const handleChangeProjectSingle = async (row: LBCertificate) => {
  ElMessage.info(`更改证书项目: ${row.name}`)
}

const handleSetupSharing = async () => {
  if (!hasSelection.value) return
  ElMessage.info('批量设置共享功能开发中')
}

const handleSetupSharingSingle = async (row: LBCertificate) => {
  ElMessage.info(`设置证书共享: ${row.name}`)
}

const handleBatchDelete = async () => {
  if (!hasSelection.value) return
  try {
    await ElMessageBox.confirm(`确认删除选中的 ${selectedRows.value.length} 个证书？`, '警告', { type: 'warning' })
    const ids = selectedRows.value.map(r => Number(r.id))
    await batchDeleteLBCertificates(ids)
    ElMessage.success('删除成功')
    loadCertificates()
  } catch (e) {
    console.error(e)
  }
}

const handleDelete = async (row: LBCertificate) => {
  try {
    await ElMessageBox.confirm(`确认删除证书 ${row.name}？`, '警告', { type: 'warning' })
    await deleteLBCertificate(row.id)
    ElMessage.success('删除成功')
    loadCertificates()
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadCertificates()
})
</script>

<style scoped>
.lb-certificates-page {
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