<template>
  <div class="elasticsearch-container">
    <div class="page-header">
      <h2>Elasticsearch实例列表</h2>
      <el-button type="primary" @click="handleCreate">创建实例</el-button>
    </div>

    <el-card class="filter-card">
      <el-form :inline="true">
        <el-form-item label="云账号">
          <CloudAccountSelector v-model="queryForm.account_id" @change="loadESInstances" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="queryForm.status" placeholder="全部" clearable @change="loadESInstances">
            <el-option label="运行中" value="Running" />
            <el-option label="已停止" value="Stopped" />
            <el-option label="创建中" value="Creating" />
            <el-option label="异常" value="Error" />
          </el-select>
        </el-form-item>
        <el-form-item label="版本">
          <el-select v-model="queryForm.version" placeholder="全部" clearable @change="loadESInstances">
            <el-option label="7.10" value="7.10" />
            <el-option label="8.5" value="8.5" />
            <el-option label="6.8" value="6.8" />
          </el-select>
        </el-form-item>
      </el-form>
    </el-card>

    <el-table :data="esInstances" v-loading="loading" row-key="id">
      <el-table-column label="名称" width="180">
        <template #default="{ row }">
          <el-link @click="viewDetail(row)">{{ row.name }}</el-link>
        </template>
      </el-table-column>
      <el-table-column prop="label" label="标签" width="120" />
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)">{{ row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="instance_type" label="类型" width="120" />
      <el-table-column prop="config" label="配置" width="150" />
      <el-table-column prop="version" label="版本" width="100" />
      <el-table-column prop="storage" label="存储" width="120" />
      <el-table-column prop="billing_method" label="计费方式" width="120" />
      <el-table-column prop="platform" label="平台" width="100" />
      <el-table-column prop="account_name" label="云账号" width="150" />
      <el-table-column prop="project_id" label="项目" width="120" />
      <el-table-column prop="region_id" label="区域" width="120" />
      <el-table-column label="操作" width="150">
        <template #default="{ row }">
          <el-button size="small" @click="handleEdit(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      class="pagination"
      layout="total, sizes, prev, pager, next"
      :total="total"
      :page-size="pageSize"
      :page-sizes="[10, 20, 50]"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
    />

    <!-- Detail Modal -->
    <el-dialog v-model="detailDialogVisible" title="Elasticsearch实例详情" width="700px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ selectedES?.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ selectedES?.name }}</el-descriptions-item>
        <el-descriptions-item label="标签">{{ selectedES?.label }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ selectedES?.status }}</el-descriptions-item>
        <el-descriptions-item label="类型">{{ selectedES?.instance_type }}</el-descriptions-item>
        <el-descriptions-item label="配置">{{ selectedES?.config }}</el-descriptions-item>
        <el-descriptions-item label="版本">{{ selectedES?.version }}</el-descriptions-item>
        <el-descriptions-item label="存储">{{ selectedES?.storage }}</el-descriptions-item>
        <el-descriptions-item label="计费方式">{{ selectedES?.billing_method }}</el-descriptions-item>
        <el-descriptions-item label="平台">{{ selectedES?.platform }}</el-descriptions-item>
        <el-descriptions-item label="云账号">{{ selectedES?.account_name }}</el-descriptions-item>
        <el-descriptions-item label="项目">{{ selectedES?.project_id }}</el-descriptions-item>
        <el-descriptions-item label="区域">{{ selectedES?.region_id }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ selectedES?.created_at }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import CloudAccountSelector from '@/components/common/CloudAccountSelector.vue'
import { listElasticsearch, deleteElasticsearch, ElasticsearchInstance } from '@/api/middleware'

const esInstances = ref<ElasticsearchInstance[]>([])
const loading = ref(false)
const detailDialogVisible = ref(false)
const selectedES = ref<ElasticsearchInstance | null>(null)
const total = ref(0)
const pageSize = ref(10)
const currentPage = ref(1)

const queryForm = ref({
  account_id: 0,
  status: '',
  version: ''
})

const getStatusType = (status: string) => {
  switch (status.toLowerCase()) {
    case 'running':
    case 'active':
      return 'success'
    case 'creating':
    case 'pending':
      return 'warning'
    case 'stopped':
    case 'error':
      return 'danger'
    default:
      return 'info'
  }
}

const loadESInstances = async () => {
  if (!queryForm.value.account_id) {
    esInstances.value = []
    return
  }

  loading.value = true
  try {
    const data = await listElasticsearch({
      account_id: queryForm.value.account_id,
      status: queryForm.value.status || undefined,
      version: queryForm.value.version || undefined
    })
    esInstances.value = data || []
    total.value = esInstances.value.length
  } catch (e: any) {
    console.error(e)
    ElMessage.error(e.message || '加载Elasticsearch实例失败')
    esInstances.value = []
  } finally {
    loading.value = false
  }
}

const viewDetail = (row: ElasticsearchInstance) => {
  selectedES.value = row
  detailDialogVisible.value = true
}

const handleCreate = () => {
  ElMessage.info('创建Elasticsearch实例功能开发中')
}

const handleEdit = (row: ElasticsearchInstance) => {
  ElMessage.info(`编辑Elasticsearch实例: ${row.name}`)
}

const handleDelete = async (row: ElasticsearchInstance) => {
  try {
    await ElMessageBox.confirm(`确认删除Elasticsearch实例 ${row.name}？`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await deleteElasticsearch(queryForm.value.account_id, row.id)
    esInstances.value = esInstances.value.filter(e => e.id !== row.id)
    ElMessage.success('删除成功')
  } catch (e: any) {
    if (e !== 'cancel') {
      console.error(e)
      ElMessage.error(e.message || '删除失败')
    }
  }
}

const handleSizeChange = (size: number) => {
  pageSize.value = size
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
}

onMounted(() => {
  // 页面加载时等待用户选择云账号
})
</script>

<style scoped>
.elasticsearch-container {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-header h2 {
  font-size: 18px;
  font-weight: bold;
  margin: 0;
}

.filter-card {
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  text-align: right;
}
</style>