<template>
  <div class="kubernetes-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">Kubernetes集群列表</span>
        </div>
      </template>

      <el-table :data="k8sClusters" v-loading="loading">
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
        <el-table-column prop="management_scale" label="管理规模" width="100" />
        <el-table-column prop="cluster_type" label="集群类型" width="120" />
        <el-table-column prop="k8s_version" label="Kubernetes版本" width="150" />
        <el-table-column prop="node_count" label="节点数" width="80" />
        <el-table-column prop="cpu" label="CPU" width="80" />
        <el-table-column prop="memory" label="内存" width="80" />
        <el-table-column prop="created_at" label="创建时间" width="160" />
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
    <el-dialog v-model="detailDialogVisible" title="Kubernetes集群详情" width="800px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ selectedK8s?.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ selectedK8s?.name }}</el-descriptions-item>
        <el-descriptions-item label="标签">{{ selectedK8s?.label }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ selectedK8s?.status }}</el-descriptions-item>
        <el-descriptions-item label="管理规模">{{ selectedK8s?.management_scale }}</el-descriptions-item>
        <el-descriptions-item label="集群类型">{{ selectedK8s?.cluster_type }}</el-descriptions-item>
        <el-descriptions-item label="Kubernetes版本">{{ selectedK8s?.k8s_version }}</el-descriptions-item>
        <el-descriptions-item label="节点数">{{ selectedK8s?.node_count }}</el-descriptions-item>
        <el-descriptions-item label="CPU">{{ selectedK8s?.cpu }}</el-descriptions-item>
        <el-descriptions-item label="内存">{{ selectedK8s?.memory }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ selectedK8s?.created_at }}</el-descriptions-item>
        <el-descriptions-item label="平台">{{ selectedK8s?.platform }}</el-descriptions-item>
        <el-descriptions-item label="云账号">{{ selectedK8s?.account }}</el-descriptions-item>
        <el-descriptions-item label="项目">{{ selectedK8s?.project }}</el-descriptions-item>
        <el-descriptions-item label="区域">{{ selectedK8s?.region }}</el-descriptions-item>
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

interface K8sCluster {
  id: string
  name: string
  label: string
  status: string
  management_scale: string
  cluster_type: string
  k8s_version: string
  node_count: number
  cpu: string
  memory: string
  created_at: string
  platform: string
  account: string
  project: string
  region: string
}

const k8sClusters = ref<K8sCluster[]>([])
const loading = ref(false)
const detailDialogVisible = ref(false)
const selectedK8s = ref<K8sCluster | null>(null)

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

const loadK8sClusters = async () => {
  loading.value = true
  try {
    // Mock data
    k8sClusters.value = [
      {
        id: 'k8s-1',
        name: 'prod-cluster',
        label: 'production',
        status: 'Running',
        management_scale: '标准',
        cluster_type: '托管集群',
        k8s_version: 'v1.24.6',
        node_count: 5,
        cpu: '20核',
        memory: '40GB',
        created_at: '2024-01-01 10:00:00',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        project: 'Project A',
        region: 'cn-hangzhou'
      },
      {
        id: 'k8s-2',
        name: 'dev-cluster',
        label: 'development',
        status: 'Running',
        management_scale: '小型',
        cluster_type: '专有集群',
        k8s_version: 'v1.26.3',
        node_count: 3,
        cpu: '8核',
        memory: '16GB',
        created_at: '2024-01-02 10:00:00',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        project: 'Project B',
        region: 'cn-shanghai'
      },
      {
        id: 'k8s-3',
        name: 'test-cluster',
        label: 'test',
        status: 'Creating',
        management_scale: '大型',
        cluster_type: '托管集群',
        k8s_version: 'v1.28.0',
        node_count: 10,
        cpu: '40核',
        memory: '80GB',
        created_at: '2024-01-03 10:00:00',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        project: 'Project A',
        region: 'cn-beijing'
      }
    ]
  } catch (e) {
    console.error(e)
    k8sClusters.value = []
  } finally {
    loading.value = false
  }
}

const viewDetail = (row: K8sCluster) => {
  selectedK8s.value = row
  detailDialogVisible.value = true
}

const handleEdit = (row: K8sCluster) => {
  ElMessage.info(`编辑Kubernetes集群: ${row.name}`)
}

const handleDelete = async (row: K8sCluster) => {
  try {
    await ElMessageBox.confirm(`确认删除Kubernetes集群 ${row.name}？`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    k8sClusters.value = k8sClusters.value.filter(k => k.id !== row.id)
    ElMessage.success('删除成功')
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadK8sClusters()
})
</script>

<style scoped>
.kubernetes-page {
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