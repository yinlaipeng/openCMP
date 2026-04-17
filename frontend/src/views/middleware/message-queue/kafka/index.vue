<template>
  <div class="kafka-container">
    <div class="page-header">
      <h2>Kafka实例列表</h2>
      <el-button type="primary" @click="handleCreate">创建实例</el-button>
    </div>

    <el-card class="filter-card">
      <el-form :inline="true">
        <el-form-item label="云账号">
          <CloudAccountSelector v-model="queryForm.account_id" @change="loadKafkaInstances" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="queryForm.status" placeholder="全部" clearable @change="loadKafkaInstances">
            <el-option label="运行中" value="Running" />
            <el-option label="已停止" value="Stopped" />
            <el-option label="创建中" value="Creating" />
            <el-option label="异常" value="Error" />
          </el-select>
        </el-form-item>
        <el-form-item label="版本">
          <el-select v-model="queryForm.version" placeholder="全部" clearable @change="loadKafkaInstances">
            <el-option label="2.8" value="2.8" />
            <el-option label="3.0" value="3.0" />
            <el-option label="2.6" value="2.6" />
          </el-select>
        </el-form-item>
      </el-form>
    </el-card>

    <el-table :data="kafkaInstances" v-loading="loading" row-key="id">
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
      <el-table-column prop="version" label="版本" width="100" />
      <el-table-column prop="storage" label="存储" width="120" />
      <el-table-column prop="bandwidth" label="带宽" width="100" />
      <el-table-column prop="endpoint" label="连接端点" width="200" />
      <el-table-column prop="retention" label="消息保留时长" width="120" />
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
    <el-dialog v-model="detailDialogVisible" title="Kafka实例详情" width="700px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ selectedKafka?.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ selectedKafka?.name }}</el-descriptions-item>
        <el-descriptions-item label="标签">{{ selectedKafka?.label }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ selectedKafka?.status }}</el-descriptions-item>
        <el-descriptions-item label="版本">{{ selectedKafka?.version }}</el-descriptions-item>
        <el-descriptions-item label="存储">{{ selectedKafka?.storage }}</el-descriptions-item>
        <el-descriptions-item label="带宽">{{ selectedKafka?.bandwidth }}</el-descriptions-item>
        <el-descriptions-item label="连接端点">{{ selectedKafka?.endpoint }}</el-descriptions-item>
        <el-descriptions-item label="消息保留时长">{{ selectedKafka?.retention }}</el-descriptions-item>
        <el-descriptions-item label="计费方式">{{ selectedKafka?.billing_method }}</el-descriptions-item>
        <el-descriptions-item label="平台">{{ selectedKafka?.platform }}</el-descriptions-item>
        <el-descriptions-item label="云账号">{{ selectedKafka?.account_name }}</el-descriptions-item>
        <el-descriptions-item label="项目">{{ selectedKafka?.project_id }}</el-descriptions-item>
        <el-descriptions-item label="区域">{{ selectedKafka?.region_id }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ selectedKafka?.created_at }}</el-descriptions-item>
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
import { listKafka, deleteKafka, KafkaInstance } from '@/api/middleware'

const kafkaInstances = ref<KafkaInstance[]>([])
const loading = ref(false)
const detailDialogVisible = ref(false)
const selectedKafka = ref<KafkaInstance | null>(null)
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

const loadKafkaInstances = async () => {
  if (!queryForm.value.account_id) {
    kafkaInstances.value = []
    return
  }

  loading.value = true
  try {
    const data = await listKafka({
      account_id: queryForm.value.account_id,
      status: queryForm.value.status || undefined,
      version: queryForm.value.version || undefined
    })
    kafkaInstances.value = data || []
    total.value = kafkaInstances.value.length
  } catch (e: any) {
    console.error(e)
    ElMessage.error(e.message || '加载Kafka实例失败')
    kafkaInstances.value = []
  } finally {
    loading.value = false
  }
}

const viewDetail = (row: KafkaInstance) => {
  selectedKafka.value = row
  detailDialogVisible.value = true
}

const handleCreate = () => {
  ElMessage.info('创建Kafka实例功能开发中')
}

const handleEdit = (row: KafkaInstance) => {
  ElMessage.info(`编辑Kafka实例: ${row.name}`)
}

const handleDelete = async (row: KafkaInstance) => {
  try {
    await ElMessageBox.confirm(`确认删除Kafka实例 ${row.name}？`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await deleteKafka(queryForm.value.account_id, row.id)
    kafkaInstances.value = kafkaInstances.value.filter(k => k.id !== row.id)
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
.kafka-container {
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