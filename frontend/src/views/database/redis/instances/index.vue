<template>
  <div class="redis-instances-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">Redis实例列表</span>
        </div>
      </template>

      <el-table :data="redisInstances" v-loading="loading">
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
      </el-table>
    </el-card>

    <!-- Detail Modal -->
    <el-dialog v-model="detailDialogVisible" title="Redis实例详情" width="600px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ selectedRedis?.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ selectedRedis?.name }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ selectedRedis?.status }}</el-descriptions-item>
        <el-descriptions-item label="实例类型">{{ selectedRedis?.instance_type }}</el-descriptions-item>
        <el-descriptions-item label="类型版本">{{ selectedRedis?.version }}</el-descriptions-item>
        <el-descriptions-item label="密码">已设置</el-descriptions-item>
        <el-descriptions-item label="链接地址">{{ selectedRedis?.connection_address }}</el-descriptions-item>
        <el-descriptions-item label="端口">{{ selectedRedis?.port }}</el-descriptions-item>
        <el-descriptions-item label="安全组">{{ selectedRedis?.security_group }}</el-descriptions-item>
        <el-descriptions-item label="计费方式">{{ selectedRedis?.billing_method }}</el-descriptions-item>
        <el-descriptions-item label="平台">{{ selectedRedis?.platform }}</el-descriptions-item>
        <el-descriptions-item label="云账号">{{ selectedRedis?.account }}</el-descriptions-item>
        <el-descriptions-item label="项目">{{ selectedRedis?.project }}</el-descriptions-item>
        <el-descriptions-item label="区域">{{ selectedRedis?.region }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ selectedRedis?.created_at }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

interface RedisInstance {
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

const redisInstances = ref<RedisInstance[]>([])
const loading = ref(false)
const detailDialogVisible = ref(false)
const selectedRedis = ref<RedisInstance | null>(null)

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

const loadRedisInstances = async () => {
  loading.value = true
  try {
    // Mock data
    redisInstances.value = [
      {
        id: 'redis-1',
        name: 'prod-cache',
        status: 'Running',
        instance_type: '标准版',
        version: '6.0',
        password: '******',
        connection_address: 'r-cn-hangzhou.redis.rds.aliyuncs.com',
        port: 6379,
        security_group: 'sg-redis',
        billing_method: '包年包月',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        project: 'Project A',
        region: 'cn-hangzhou',
        created_at: '2024-01-01 10:00:00'
      },
      {
        id: 'redis-2',
        name: 'dev-session',
        status: 'Running',
        instance_type: '基础版',
        version: '7.0',
        password: '******',
        connection_address: 'r-cn-shanghai.redis.rds.aliyuncs.com',
        port: 6379,
        security_group: 'sg-redis-dev',
        billing_method: '按量付费',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        project: 'Project B',
        region: 'cn-shanghai',
        created_at: '2024-01-02 10:00:00'
      },
      {
        id: 'redis-3',
        name: 'test-cache',
        status: 'Creating',
        instance_type: '集群版',
        version: '5.0',
        password: '******',
        connection_address: '-',
        port: 6379,
        security_group: 'sg-redis-test',
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
    redisInstances.value = []
  } finally {
    loading.value = false
  }
}

const viewDetail = (row: RedisInstance) => {
  selectedRedis.value = row
  detailDialogVisible.value = true
}

onMounted(() => {
  loadRedisInstances()
})
</script>

<style scoped>
.redis-instances-page {
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