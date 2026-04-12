<template>
  <div class="lb-instances-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">负载均衡实例列表</span>
        </div>
      </template>

      <el-tabs v-model="activeTab">
        <el-tab-pane label="全部" name="all">
          <el-table :data="allInstances" v-loading="loading">
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
            <el-table-column prop="address" label="服务地址" width="180" />
            <el-table-column prop="spec" label="规格" width="120" />
            <el-table-column prop="vpc" label="vpc" width="180" />
            <el-table-column prop="security_group" label="安全组" width="150" />
            <el-table-column prop="billing_method" label="计费方式" width="120" />
            <el-table-column prop="platform" label="平台" width="100" />
            <el-table-column prop="project" label="项目" width="120" />
            <el-table-column prop="region" label="区域" width="120" />
            <el-table-column label="操作" width="150">
              <template #default="{ row }">
                <el-button size="small" @click="handleEdit(row)">编辑</el-button>
                <el-button size="small" @click="manageListeners(row)">管理监听</el-button>
                <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="idc" name="idc">
          <el-table :data="idcInstances" v-loading="loading">
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
            <el-table-column prop="address" label="服务地址" width="180" />
            <el-table-column prop="spec" label="规格" width="120" />
            <el-table-column prop="vpc" label="vpc" width="180" />
            <el-table-column prop="security_group" label="安全组" width="150" />
            <el-table-column prop="billing_method" label="计费方式" width="120" />
            <el-table-column prop="platform" label="平台" width="100" />
            <el-table-column prop="project" label="项目" width="120" />
            <el-table-column prop="region" label="区域" width="120" />
            <el-table-column label="操作" width="150">
              <template #default="{ row }">
                <el-button size="small" @click="handleEdit(row)">编辑</el-button>
                <el-button size="small" @click="manageListeners(row)">管理监听</el-button>
                <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="公有云" name="public_cloud">
          <el-table :data="publicCloudInstances" v-loading="loading">
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
            <el-table-column prop="address" label="服务地址" width="180" />
            <el-table-column prop="spec" label="规格" width="120" />
            <el-table-column prop="vpc" label="vpc" width="180" />
            <el-table-column prop="security_group" label="安全组" width="150" />
            <el-table-column prop="billing_method" label="计费方式" width="120" />
            <el-table-column prop="platform" label="平台" width="100" />
            <el-table-column prop="project" label="项目" width="120" />
            <el-table-column prop="region" label="区域" width="120" />
            <el-table-column label="操作" width="150">
              <template #default="{ row }">
                <el-button size="small" @click="handleEdit(row)">编辑</el-button>
                <el-button size="small" @click="manageListeners(row)">管理监听</el-button>
                <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- Detail Modal -->
    <el-dialog v-model="detailDialogVisible" title="负载均衡实例详情" width="600px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ selectedInstance?.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ selectedInstance?.name }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ selectedInstance?.status }}</el-descriptions-item>
        <el-descriptions-item label="服务地址">{{ selectedInstance?.address }}</el-descriptions-item>
        <el-descriptions-item label="规格">{{ selectedInstance?.spec }}</el-descriptions-item>
        <el-descriptions-item label="VPC">{{ selectedInstance?.vpc }}</el-descriptions-item>
        <el-descriptions-item label="安全组">{{ selectedInstance?.security_group }}</el-descriptions-item>
        <el-descriptions-item label="计费方式">{{ selectedInstance?.billing_method }}</el-descriptions-item>
        <el-descriptions-item label="平台">{{ selectedInstance?.platform }}</el-descriptions-item>
        <el-descriptions-item label="项目">{{ selectedInstance?.project }}</el-descriptions-item>
        <el-descriptions-item label="区域">{{ selectedInstance?.region }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ selectedInstance?.created_at }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- Listeners Modal -->
    <el-dialog v-model="listenersDialogVisible" title="管理监听规则" width="800px">
      <el-table :data="listeners" v-loading="listenersLoading">
        <el-table-column prop="name" label="监听名称" width="150" />
        <el-table-column prop="protocol" label="协议" width="100" />
        <el-table-column prop="port" label="端口" width="100" />
        <el-table-column prop="backend_port" label="后端端口" width="100" />
        <el-table-column prop="scheduler" label="调度算法" width="120" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'Active' ? 'success' : 'info'">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-button size="small" type="danger" @click="deleteListener(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <template #footer>
        <el-button @click="listenersDialogVisible = false">关闭</el-button>
        <el-button type="primary" @click="addListener">添加监听</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

interface LBInstance {
  id: string
  name: string
  status: string
  address: string
  spec: string
  vpc: string
  security_group: string
  billing_method: string
  platform: string
  project: string
  region: string
  created_at: string
}

interface Listener {
  name: string
  protocol: string
  port: number
  backend_port: number
  scheduler: string
  status: string
}

const allInstances = ref<LBInstance[]>([])
const idcInstances = ref<LBInstance[]>([])
const publicCloudInstances = ref<LBInstance[]>([])
const loading = ref(false)
const activeTab = ref('all')
const detailDialogVisible = ref(false)
const listenersDialogVisible = ref(false)
const listenersLoading = ref(false)
const selectedInstance = ref<LBInstance | null>(null)
const listeners = ref<Listener[]>([])

const getStatusType = (status: string) => {
  switch (status.toLowerCase()) {
    case 'active':
    case 'running':
      return 'success'
    case 'pending':
    case 'creating':
      return 'warning'
    case 'error':
    case 'stopped':
      return 'danger'
    default:
      return 'info'
  }
}

const loadInstances = async () => {
  loading.value = true
  try {
    // Mock data
    const instancesData: LBInstance[] = [
      {
        id: 'lb-1',
        name: '负载均衡 1',
        status: 'Active',
        address: '192.168.1.100',
        spec: '性能保障型',
        vpc: 'vpc-1',
        security_group: 'sg-1',
        billing_method: '按量付费',
        platform: '阿里云',
        project: 'Project A',
        region: 'cn-hangzhou',
        created_at: '2024-01-01 10:00:00'
      },
      {
        id: 'lb-2',
        name: '负载均衡 2 (IDC)',
        status: 'Active',
        address: '10.10.10.100',
        spec: '标准型',
        vpc: 'vpc-2',
        security_group: 'sg-2',
        billing_method: '包年包月',
        platform: '本地IDC',
        project: 'Project B',
        region: '本地机房',
        created_at: '2024-01-02 10:00:00'
      },
      {
        id: 'lb-3',
        name: '负载均衡 3',
        status: 'Creating',
        address: '-',
        spec: '性能保障型',
        vpc: 'vpc-3',
        security_group: 'sg-3',
        billing_method: '按量付费',
        platform: '阿里云',
        project: 'Project A',
        region: 'cn-shanghai',
        created_at: '2024-01-03 10:00:00'
      }
    ]
    allInstances.value = instancesData
    idcInstances.value = instancesData.filter(i => i.platform.toLowerCase().includes('idc') || i.region.includes('本地'))
    publicCloudInstances.value = instancesData.filter(i => !i.platform.toLowerCase().includes('idc') && !i.region.includes('本地'))
  } catch (e) {
    console.error(e)
    allInstances.value = []
    idcInstances.value = []
    publicCloudInstances.value = []
  } finally {
    loading.value = false
  }
}

const viewDetail = (row: LBInstance) => {
  selectedInstance.value = row
  detailDialogVisible.value = true
}

const manageListeners = (row: LBInstance) => {
  selectedInstance.value = row
  listenersDialogVisible.value = true
  listenersLoading.value = true
  listeners.value = [
    { name: 'HTTP监听', protocol: 'HTTP', port: 80, backend_port: 8080, scheduler: '轮询', status: 'Active' },
    { name: 'HTTPS监听', protocol: 'HTTPS', port: 443, backend_port: 8443, scheduler: '轮询', status: 'Active' }
  ]
  listenersLoading.value = false
}

const addListener = () => {
  ElMessage.info('添加监听功能开发中')
}

const deleteListener = async (row: Listener) => {
  try {
    await ElMessageBox.confirm(`确认删除监听 ${row.name}？`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    listeners.value = listeners.value.filter(l => l.name !== row.name)
    ElMessage.success('删除成功')
  } catch (e) {
    console.error(e)
  }
}

const handleEdit = (row: LBInstance) => {
  ElMessage.info(`编辑负载均衡实例: ${row.name}`)
}

const handleDelete = async (row: LBInstance) => {
  try {
    await ElMessageBox.confirm(`确认删除负载均衡实例 ${row.name}？`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    allInstances.value = allInstances.value.filter(i => i.id !== row.id)
    idcInstances.value = idcInstances.value.filter(i => i.id !== row.id)
    publicCloudInstances.value = publicCloudInstances.value.filter(i => i.id !== row.id)
    ElMessage.success('删除成功')
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadInstances()
})
</script>

<style scoped>
.lb-instances-page {
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