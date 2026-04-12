<template>
  <div class="rds-instances-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">RDS实例列表</span>
        </div>
      </template>

      <el-table :data="rdsInstances" v-loading="loading">
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
        <el-table-column prop="type" label="类型" width="100" />
        <el-table-column prop="engine" label="数据库引擎" width="150">
          <template #default="{ row }">
            {{ row.engine }} {{ row.engine_version }}
          </template>
        </el-table-column>
        <el-table-column prop="connection_address" label="链接地址" width="200" />
        <el-table-column prop="port" label="数据库端口号" width="120" />
        <el-table-column prop="storage_type" label="存储类型" width="120" />
        <el-table-column prop="security_group" label="安全组" width="150" />
        <el-table-column prop="billing_method" label="计费方式" width="120" />
        <el-table-column prop="platform" label="平台" width="100" />
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
    <el-dialog v-model="detailDialogVisible" title="RDS实例详情" width="700px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ selectedRDS?.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ selectedRDS?.name }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ selectedRDS?.status }}</el-descriptions-item>
        <el-descriptions-item label="类型">{{ selectedRDS?.type }}</el-descriptions-item>
        <el-descriptions-item label="数据库引擎">{{ selectedRDS?.engine }} {{ selectedRDS?.engine_version }}</el-descriptions-item>
        <el-descriptions-item label="链接地址">{{ selectedRDS?.connection_address }}</el-descriptions-item>
        <el-descriptions-item label="端口号">{{ selectedRDS?.port }}</el-descriptions-item>
        <el-descriptions-item label="存储类型">{{ selectedRDS?.storage_type }}</el-descriptions-item>
        <el-descriptions-item label="安全组">{{ selectedRDS?.security_group }}</el-descriptions-item>
        <el-descriptions-item label="计费方式">{{ selectedRDS?.billing_method }}</el-descriptions-item>
        <el-descriptions-item label="平台">{{ selectedRDS?.platform }}</el-descriptions-item>
        <el-descriptions-item label="项目">{{ selectedRDS?.project }}</el-descriptions-item>
        <el-descriptions-item label="区域">{{ selectedRDS?.region }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ selectedRDS?.created_at }}</el-descriptions-item>
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

interface RDSInstance {
  id: string
  name: string
  status: string
  type: string
  engine: string
  engine_version: string
  connection_address: string
  port: number
  storage_type: string
  security_group: string
  billing_method: string
  platform: string
  project: string
  region: string
  created_at: string
}

const rdsInstances = ref<RDSInstance[]>([])
const loading = ref(false)
const detailDialogVisible = ref(false)
const selectedRDS = ref<RDSInstance | null>(null)

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

const loadRDSInstances = async () => {
  loading.value = true
  try {
    // Mock data
    rdsInstances.value = [
      {
        id: 'rds-1',
        name: 'prod-mysql',
        status: 'Running',
        type: '高可用版',
        engine: 'MySQL',
        engine_version: '5.7',
        connection_address: 'rm-cn-hangzhou.mysql.rds.aliyuncs.com',
        port: 3306,
        storage_type: 'SSD云盘',
        security_group: 'sg-mysql',
        billing_method: '包年包月',
        platform: '阿里云',
        project: 'Project A',
        region: 'cn-hangzhou',
        created_at: '2024-01-01 10:00:00'
      },
      {
        id: 'rds-2',
        name: 'dev-postgres',
        status: 'Running',
        type: '基础版',
        engine: 'PostgreSQL',
        engine_version: '14',
        connection_address: 'pg-cn-shanghai.pg.rds.aliyuncs.com',
        port: 5432,
        storage_type: 'SSD云盘',
        security_group: 'sg-pg',
        billing_method: '按量付费',
        platform: '阿里云',
        project: 'Project B',
        region: 'cn-shanghai',
        created_at: '2024-01-02 10:00:00'
      },
      {
        id: 'rds-3',
        name: 'test-sqlserver',
        status: 'Creating',
        type: '高可用版',
        engine: 'SQLServer',
        engine_version: '2019',
        connection_address: '-',
        port: 1433,
        storage_type: 'SSD云盘',
        security_group: 'sg-sqlserver',
        billing_method: '按量付费',
        platform: '阿里云',
        project: 'Project A',
        region: 'cn-beijing',
        created_at: '2024-01-03 10:00:00'
      }
    ]
  } catch (e) {
    console.error(e)
    rdsInstances.value = []
  } finally {
    loading.value = false
  }
}

const viewDetail = (row: RDSInstance) => {
  selectedRDS.value = row
  detailDialogVisible.value = true
}

const handleEdit = (row: RDSInstance) => {
  ElMessage.info(`编辑RDS实例: ${row.name}`)
}

const handleDelete = async (row: RDSInstance) => {
  try {
    await ElMessageBox.confirm(`确认删除RDS实例 ${row.name}？`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    rdsInstances.value = rdsInstances.value.filter(r => r.id !== row.id)
    ElMessage.success('删除成功')
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadRDSInstances()
})
</script>

<style scoped>
.rds-instances-page {
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