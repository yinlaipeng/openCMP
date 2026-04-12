<template>
  <div class="mongodb-instances-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">MongoDB实例列表</span>
        </div>
      </template>

      <el-table :data="mongodbInstances" v-loading="loading">
        <el-table-column label="名称" width="180">
          <template #default="{ row }">
            <el-link @click="viewDetail(row)">{{ row.name }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="instance_type" label="实例类型" width="120" />
        <el-table-column prop="version" label="类型版本" width="100" />
        <el-table-column prop="password" label="密码" width="100">
          <template #default="{ row }">
            <el-tag>已设置</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="connection_address" label="链接地址" width="200" />
        <el-table-column prop="port" label="端口" width="100" />
        <el-table-column prop="security_group" label="安全组" width="150" />
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
    <el-dialog v-model="detailDialogVisible" title="MongoDB实例详情" width="600px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ selectedMongoDB?.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ selectedMongoDB?.name }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ selectedMongoDB?.status }}</el-descriptions-item>
        <el-descriptions-item label="实例类型">{{ selectedMongoDB?.instance_type }}</el-descriptions-item>
        <el-descriptions-item label="类型版本">{{ selectedMongoDB?.version }}</el-descriptions-item>
        <el-descriptions-item label="密码">已设置</el-descriptions-item>
        <el-descriptions-item label="链接地址">{{ selectedMongoDB?.connection_address }}</el-descriptions-item>
        <el-descriptions-item label="端口">{{ selectedMongoDB?.port }}</el-descriptions-item>
        <el-descriptions-item label="安全组">{{ selectedMongoDB?.security_group }}</el-descriptions-item>
        <el-descriptions-item label="计费方式">{{ selectedMongoDB?.billing_method }}</el-descriptions-item>
        <el-descriptions-item label="平台">{{ selectedMongoDB?.platform }}</el-descriptions-item>
        <el-descriptions-item label="云账号">{{ selectedMongoDB?.account }}</el-descriptions-item>
        <el-descriptions-item label="项目">{{ selectedMongoDB?.project }}</el-descriptions-item>
        <el-descriptions-item label="区域">{{ selectedMongoDB?.region }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ selectedMongoDB?.created_at }}</el-descriptions-item>
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

interface MongoDBInstance {
  id: string
  name: string
  status: string
  instance_type: string
  version: string
  password: string
  connection_address: string
  port: number
  security_group: string
  billing_method: string
  platform: string
  account: string
  project: string
  region: string
  created_at: string
}

const mongodbInstances = ref<MongoDBInstance[]>([])
const loading = ref(false)
const detailDialogVisible = ref(false)
const selectedMongoDB = ref<MongoDBInstance | null>(null)

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

const loadMongoDBInstances = async () => {
  loading.value = true
  try {
    // Mock data
    mongodbInstances.value = [
      {
        id: 'mongodb-1',
        name: 'prod-nosql',
        status: 'Running',
        instance_type: '副本集',
        version: '4.4',
        password: '******',
        connection_address: 'dds-cn-hangzhou.mongodb.rds.aliyuncs.com',
        port: 27017,
        security_group: 'sg-mongodb',
        billing_method: '包年包月',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        project: 'Project A',
        region: 'cn-hangzhou',
        created_at: '2024-01-01 10:00:00'
      },
      {
        id: 'mongodb-2',
        name: 'dev-nosql',
        status: 'Running',
        instance_type: '单节点',
        version: '5.0',
        password: '******',
        connection_address: 'dds-cn-shanghai.mongodb.rds.aliyuncs.com',
        port: 27017,
        security_group: 'sg-mongodb-dev',
        billing_method: '按量付费',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        project: 'Project B',
        region: 'cn-shanghai',
        created_at: '2024-01-02 10:00:00'
      },
      {
        id: 'mongodb-3',
        name: 'test-nosql',
        status: 'Creating',
        instance_type: '分片集群',
        version: '6.0',
        password: '******',
        connection_address: '-',
        port: 27017,
        security_group: 'sg-mongodb-test',
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
    mongodbInstances.value = []
  } finally {
    loading.value = false
  }
}

const viewDetail = (row: MongoDBInstance) => {
  selectedMongoDB.value = row
  detailDialogVisible.value = true
}

const handleEdit = (row: MongoDBInstance) => {
  ElMessage.info(`编辑MongoDB实例: ${row.name}`)
}

const handleDelete = async (row: MongoDBInstance) => {
  try {
    await ElMessageBox.confirm(`确认删除MongoDB实例 ${row.name}？`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    mongodbInstances.value = mongodbInstances.value.filter(m => m.id !== row.id)
    ElMessage.success('删除成功')
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadMongoDBInstances()
})
</script>

<style scoped>
.mongodb-instances-page {
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