<template>
  <div class="proxies-container">
    <div class="page-header">
      <h2>代理管理</h2>
      <div class="toolbar">
        <el-button @click="loadProxies" :loading="loading">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
        <el-button type="primary" @click="showCreateDialog">
          <el-icon><Plus /></el-icon>
          新建
        </el-button>
        <el-dropdown trigger="click" @command="handleBatchCommand" :disabled="selectedIds.length === 0">
          <el-button>
            <el-icon><Operation /></el-icon>
            批量操作
            <el-badge v-if="selectedIds.length > 0" :value="selectedIds.length" type="primary" />
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="batchSetSharing" :icon="Share">批量设置共享</el-dropdown-item>
              <el-dropdown-item command="batchDelete" :icon="Delete" divided>
                <span style="color: var(--el-color-danger)">批量删除</span>
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>

    <!-- 筛选区 -->
    <el-card class="filter-card">
      <el-form :inline="true" :model="queryForm" @submit.prevent="loadProxies">
        <el-form-item label="名称">
          <el-input v-model="queryForm.name" placeholder="支持名称搜索" clearable style="width: 200px">
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="所属域">
          <el-select v-model="queryForm.domain_id" placeholder="全部" clearable style="width: 140px">
            <el-option v-for="domain in domains" :key="domain.id" :label="domain.name" :value="domain.id" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadProxies">查询</el-button>
          <el-button @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 数据表格 -->
    <el-table
      ref="tableRef"
      :data="proxies"
      v-loading="loading"
      style="width: 100%"
      row-key="id"
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="50" />
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="名称" min-width="150">
        <template #default="{ row }">
          <el-button link type="primary" @click="handleView(row)" class="name-link">
            {{ row.name }}
          </el-button>
        </template>
      </el-table-column>
      <el-table-column prop="https_proxy" label="HTTPS代理" min-width="200">
        <template #default="{ row }">
          {{ row.https_proxy || '-' }}
        </template>
      </el-table-column>
      <el-table-column prop="http_proxy" label="HTTP代理" min-width="200">
        <template #default="{ row }">
          {{ row.http_proxy || '-' }}
        </template>
      </el-table-column>
      <el-table-column prop="no_proxy" label="不走代理" min-width="150">
        <template #default="{ row }">
          {{ row.no_proxy || '-' }}
        </template>
      </el-table-column>
      <el-table-column prop="owner_domain" label="所属域" width="120">
        <template #default="{ row }">
          {{ getDomainName(row.domain_id) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="160" fixed="right">
        <template #default="{ row }">
          <el-button size="small" type="primary" @click="handleView(row)">查看</el-button>
          <el-dropdown trigger="click" @command="(cmd: string) => handleDropdownCommand(cmd, row)" style="margin-left: 8px">
            <el-button size="small" link>
              更多 <el-icon class="el-icon--right"><ArrowDown /></el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="edit" :icon="EditPen">编辑属性</el-dropdown-item>
                <el-dropdown-item command="setSharing" :icon="Share">设置共享</el-dropdown-item>
                <el-dropdown-item command="delete" :icon="Delete" divided>
                  <span style="color: var(--el-color-danger)">删除</span>
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      v-model:current-page="currentPage"
      v-model:page-size="pageSize"
      :total="total"
      :page-sizes="[10, 20, 50, 100]"
      layout="total, sizes, prev, pager, next, jumper"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
      class="pagination"
    />

    <!-- 创建/编辑对话框 -->
    <el-dialog
      v-model="showDialog"
      :title="isEdit ? '编辑代理' : '新建代理'"
      width="600px"
      destroy-on-close
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入代理名称" />
        </el-form-item>

        <el-form-item label="所属域" prop="domain_id">
          <el-select v-model="form.domain_id" placeholder="请选择所属域" style="width: 100%">
            <el-option v-for="domain in domains" :key="domain.id" :label="domain.name" :value="domain.id" />
          </el-select>
        </el-form-item>

        <el-form-item label="HTTPS代理">
          <el-input v-model="form.https_proxy" placeholder="如：https://proxy.example.com:443" />
        </el-form-item>

        <el-form-item label="HTTP代理">
          <el-input v-model="form.http_proxy" placeholder="如：http://proxy.example.com:8080" />
        </el-form-item>

        <el-form-item label="不走代理">
          <el-input v-model="form.no_proxy" placeholder="多个地址用逗号分隔，如：localhost,127.0.0.1" />
        </el-form-item>

        <el-form-item label="共享范围">
          <el-select v-model="form.shared_scope" placeholder="请选择共享范围" style="width: 100%" clearable>
            <el-option label="私有" value="private" />
            <el-option label="域内共享" value="domain" />
            <el-option label="全局共享" value="global" />
          </el-select>
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>

    <!-- 设置共享对话框 -->
    <el-dialog v-model="showSharingDialog" title="设置共享" width="400px">
      <el-form label-width="80px">
        <el-form-item label="代理名称">
          <span>{{ currentProxy?.name }}</span>
        </el-form-item>
        <el-form-item label="共享范围">
          <el-select v-model="sharingForm.shared_scope" placeholder="请选择共享范围" style="width: 100%">
            <el-option label="私有" value="private" />
            <el-option label="域内共享" value="domain" />
            <el-option label="全局共享" value="global" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showSharingDialog = false">取消</el-button>
        <el-button type="primary" @click="handleSetSharing" :loading="settingSharing">确定</el-button>
      </template>
    </el-dialog>

    <!-- 详情抽屉 -->
    <el-drawer
      v-model="showDetailDrawer"
      title="代理详情"
      size="50%"
      direction="rtl"
    >
      <div class="drawer-header" v-if="currentProxy">
        <div class="proxy-icon">
          <el-avatar :size="48" :style="{ backgroundColor: '#67C23A' }">
            <el-icon :size="32"><Connection /></el-icon>
          </el-avatar>
        </div>
        <div class="proxy-info">
          <h3>{{ currentProxy.name }}</h3>
          <div class="proxy-tags">
            <el-tag size="small" type="info">
              {{ getSharedScopeText(currentProxy.shared_scope) }}
            </el-tag>
          </div>
        </div>
        <div class="quick-actions">
          <el-button size="small" @click="handleEdit(currentProxy)">
            <el-icon><EditPen /></el-icon>
            编辑
          </el-button>
          <el-button size="small" @click="handleSetSharingDialog(currentProxy)">
            <el-icon><Share /></el-icon>
            设置共享
          </el-button>
          <el-button size="small" type="danger" @click="handleDelete(currentProxy)">
            <el-icon><Delete /></el-icon>
            删除
          </el-button>
        </div>
      </div>

      <!-- 基本信息 -->
      <el-descriptions :column="1" border style="margin-top: 20px" title="基本信息">
        <el-descriptions-item label="ID">{{ currentProxy?.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ currentProxy?.name }}</el-descriptions-item>
        <el-descriptions-item label="HTTPS代理">{{ currentProxy?.https_proxy || '-' }}</el-descriptions-item>
        <el-descriptions-item label="HTTP代理">{{ currentProxy?.http_proxy || '-' }}</el-descriptions-item>
        <el-descriptions-item label="不走代理">{{ currentProxy?.no_proxy || '-' }}</el-descriptions-item>
        <el-descriptions-item label="共享范围">{{ getSharedScopeText(currentProxy?.shared_scope) }}</el-descriptions-item>
        <el-descriptions-item label="所属域">{{ getDomainName(currentProxy?.domain_id) }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ formatDate(currentProxy?.created_at) }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ formatDate(currentProxy?.updated_at) }}</el-descriptions-item>
      </el-descriptions>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance } from 'element-plus'
import { Plus, Delete, Connection, ArrowDown, EditPen, Refresh, Operation, Search, Share } from '@element-plus/icons-vue'
import { getProxies, getProxy, createProxy, updateProxy, deleteProxy, batchDeleteProxies, setProxySharing } from '@/api/proxy'
import { getDomains } from '@/api/iam'
import type { Proxy, CreateProxyRequest } from '@/api/proxy'
import EmptyState from '@/components/common/EmptyState.vue'

const proxies = ref<Proxy[]>([])
const domains = ref<{ id: number; name: string }[]>([])
const loading = ref(false)
const showDialog = ref(false)
const showDetailDrawer = ref(false)
const showSharingDialog = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const settingSharing = ref(false)
const formRef = ref<FormInstance>()
const tableRef = ref()
const currentProxy = ref<Proxy | null>(null)
const selectedIds = ref<number[]>([])

// 分页数据
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

// 查询表单
const queryForm = reactive({
  name: '',
  domain_id: undefined as number | undefined
})

const form = reactive<CreateProxyRequest>({
  name: '',
  domain_id: 1,
  https_proxy: '',
  http_proxy: '',
  no_proxy: '',
  shared_scope: ''
})

const sharingForm = reactive({
  shared_scope: ''
})

const rules = {
  name: [{ required: true, message: '请输入代理名称', trigger: 'blur' }],
  domain_id: [{ required: true, message: '请选择所属域', trigger: 'change' }]
}

const loadProxies = async () => {
  loading.value = true
  try {
    const params: any = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    if (queryForm.name) params.name = queryForm.name
    if (queryForm.domain_id) params.domain_id = queryForm.domain_id

    const res = await getProxies(params)
    proxies.value = res.items || []
    total.value = res.total || 0
  } catch (e) {
    console.error(e)
    ElMessage.error('加载代理列表失败')
  } finally {
    loading.value = false
  }
}

const loadDomains = async () => {
  try {
    const res = await getDomains()
    domains.value = (res.items || []).map(d => ({ id: d.id, name: d.name }))
  } catch (e) {
    console.error(e)
  }
}

const handleSelectionChange = (selection: Proxy[]) => {
  selectedIds.value = selection.map(p => p.id)
}

const handleBatchCommand = async (command: string) => {
  if (selectedIds.value.length === 0) {
    ElMessage.warning('请先选择代理')
    return
  }

  try {
    if (command === 'batchDelete') {
      await ElMessageBox.confirm(`确定要批量删除 ${selectedIds.value.length} 个代理吗？此操作不可恢复。`, '提示', { type: 'warning' })
      await batchDeleteProxies(selectedIds.value)
      ElMessage.success('批量删除成功')
    } else if (command === 'batchSetSharing') {
      await ElMessageBox.prompt('请输入共享范围（private/domain/global）', '批量设置共享', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        inputPattern: /^(private|domain|global)$/,
        inputErrorMessage: '请输入正确的共享范围'
      }).then(({ value }) => {
        for (const id of selectedIds.value) {
          setProxySharing(id, value)
        }
        ElMessage.success('批量设置共享成功')
      })
    }
    tableRef.value?.clearSelection()
    loadProxies()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error('批量操作失败')
    }
  }
}

const getSharedScopeText = (scope?: string) => {
  const map: Record<string, string> = {
    private: '私有',
    domain: '域内共享',
    global: '全局共享'
  }
  return map[scope || ''] || scope || '-'
}

const getDomainName = (domainId?: number) => {
  if (!domainId) return '-'
  const domain = domains.value.find(d => d.id === domainId)
  return domain?.name || `域${domainId}`
}

const formatDate = (date?: string) => {
  if (!date) return ''
  return new Date(date).toLocaleString('zh-CN')
}

const showCreateDialog = () => {
  isEdit.value = false
  resetForm()
  showDialog.value = true
}

const resetForm = () => {
  form.name = ''
  form.domain_id = domains.value[0]?.id || 1
  form.https_proxy = ''
  form.http_proxy = ''
  form.no_proxy = ''
  form.shared_scope = ''
}

const handleView = async (row: Proxy) => {
  try {
    const detail = await getProxy(row.id)
    currentProxy.value = detail
    showDetailDrawer.value = true
  } catch (e) {
    console.error(e)
    ElMessage.error('获取详情失败')
  }
}

const handleEdit = async (row: Proxy) => {
  try {
    const detail = await getProxy(row.id)
    isEdit.value = true
    form.name = detail.name
    form.domain_id = detail.domain_id || 1
    form.https_proxy = detail.https_proxy || ''
    form.http_proxy = detail.http_proxy || ''
    form.no_proxy = detail.no_proxy || ''
    form.shared_scope = detail.shared_scope || ''
    currentProxy.value = detail
    showDialog.value = true
  } catch (e) {
    console.error(e)
    ElMessage.error('获取代理详情失败')
  }
}

const handleSetSharingDialog = (row: Proxy) => {
  currentProxy.value = row
  sharingForm.shared_scope = row.shared_scope || 'private'
  showSharingDialog.value = true
}

const handleSetSharing = async () => {
  if (!currentProxy.value) return
  settingSharing.value = true
  try {
    await setProxySharing(currentProxy.value.id, sharingForm.shared_scope)
    ElMessage.success('设置共享成功')
    showSharingDialog.value = false
    loadProxies()
  } catch (e) {
    console.error(e)
    ElMessage.error('设置共享失败')
  } finally {
    settingSharing.value = false
  }
}

const handleDelete = async (row: Proxy) => {
  try {
    await ElMessageBox.confirm(`确定要删除代理 "${row.name}" 吗？此操作不可恢复。`, '提示', { type: 'warning' })
    await deleteProxy(row.id)
    ElMessage.success('删除成功')
    showDetailDrawer.value = false
    loadProxies()
  } catch (e: any) {
    if (e !== 'cancel') {
      console.error(e)
      ElMessage.error('删除失败')
    }
  }
}

const handleDropdownCommand = (command: string, row: Proxy) => {
  switch (command) {
    case 'edit':
      handleEdit(row)
      break
    case 'setSharing':
      handleSetSharingDialog(row)
      break
    case 'delete':
      handleDelete(row)
      break
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitting.value = true
    try {
      if (isEdit.value && currentProxy.value) {
        await updateProxy(currentProxy.value.id, form)
        ElMessage.success('更新成功')
      } else {
        await createProxy(form)
        ElMessage.success('创建成功')
      }
      showDialog.value = false
      loadProxies()
    } catch (e) {
      console.error(e)
      ElMessage.error(isEdit.value ? '更新失败' : '创建失败')
    } finally {
      submitting.value = false
    }
  })
}

const resetQuery = () => {
  queryForm.name = ''
  queryForm.domain_id = undefined
  currentPage.value = 1
  loadProxies()
}

const handleSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
  loadProxies()
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  loadProxies()
}

onMounted(() => {
  loadProxies()
  loadDomains()
})
</script>

<style scoped>
.proxies-container {
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

.toolbar {
  display: flex;
  gap: 8px;
  align-items: center;
}

.filter-card {
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  text-align: right;
}

.name-link {
  font-weight: 500;
}

/* 抽屉顶部区域 */
.drawer-header {
  display: flex;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid var(--el-border-color);
}

.proxy-icon {
  margin-right: 16px;
}

.proxy-info h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
}

.proxy-tags {
  display: flex;
  gap: 8px;
  margin-top: 8px;
  align-items: center;
}

.quick-actions {
  margin-left: auto;
  display: flex;
  gap: 8px;
}
</style>