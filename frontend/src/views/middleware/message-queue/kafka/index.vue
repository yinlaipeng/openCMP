<template>
  <div class="kafka-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">Kafka实例列表</span>
        </div>
      </template>

      <el-table :data="kafkaInstances" v-loading="loading">
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
        <el-table-column prop="account" label="云账号" width="150" />
        <el-table-column prop="project" label="项目" width="120" />
        <el-table-column prop="region" label="区域" width="120" />
        <el-table-column label="操作" width="150">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

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
        <el-descriptions-item label="云账号">{{ selectedKafka?.account }}</el-descriptions-item>
        <el-descriptions-item label="项目">{{ selectedKafka?.project }}</el-descriptions-item>
        <el-descriptions-item label="区域">{{ selectedKafka?.region }}</el-descriptions-item>
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

interface KafkaInstance {
  id: string
  name: string
  label: string
  status: string
  version: string
  storage: string
  bandwidth: string
  endpoint: string
  retention: string
  billing_method: string
  platform: string
  account: string
  project: string
  region: string
  created_at: string
}

const kafkaInstances = ref<KafkaInstance[]>([])
const loading = ref(false)
const detailDialogVisible = ref(false)
const selectedKafka = ref<KafkaInstance | null>(null)

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
  loading.value = true
  try {
    // Mock data
    kafkaInstances.value = [
      {
        id: 'kafka-1',
        name: 'prod-mq',
        label: 'production',
        status: 'Running',
        version: '2.8',
        storage: '500GB',
        bandwidth: '50MB/s',
        endpoint: 'kafka-cn-hangzhou.aliyuncs.com:9092',
        retention: '7天',
        billing_method: '包年包月',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        project: 'Project A',
        region: 'cn-hangzhou',
        created_at: '2024-01-01 10:00:00'
      },
      {
        id: 'kafka-2',
        name: 'dev-mq',
        label: 'development',
        status: 'Running',
        version: '3.0',
        storage: '200GB',
        bandwidth: '20MB/s',
        endpoint: 'kafka-cn-shanghai.aliyuncs.com:9092',
        retention: '3天',
        billing_method: '按量付费',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        project: 'Project B',
        region: 'cn-shanghai',
        created_at: '2024-01-02 10:00:00'
      },
      {
        id: 'kafka-3',
        name: 'test-mq',
        label: 'test',
        status: 'Creating',
        version: '2.6',
        storage: '100GB',
        bandwidth: '10MB/s',
        endpoint: '-',
        retention: '1天',
        billing_method: '按量付费',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        project: 'Project A',
        region: 'cn-beijing',
        created_at: '2024-01-03 10:00:00'
      }
    ]
  } catch (e) {
    console.error(e)
    kafkaInstances.value = []
  } finally {
    loading.value = false
  }
}

const viewDetail = (row: KafkaInstance) => {
  selectedKafka.value = row
  detailDialogVisible.value = true
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
    kafkaInstances.value = kafkaInstances.value.filter(k => k.id !== row.id)
    ElMessage.success('删除成功')
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadKafkaInstances()
})
</script>

<style scoped>
.kafka-page {
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
</style>