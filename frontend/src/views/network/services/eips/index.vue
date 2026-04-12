<template>
  <div class="eips-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">弹性公网IP列表</span>
        </div>
      </template>

      <el-tabs v-model="activeTab">
        <el-tab-pane label="全部" name="all">
          <el-table :data="allEips" v-loading="loading">
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
            <el-table-column prop="tags" label="标签" width="150" />
            <el-table-column prop="ip" label="IP" width="150" />
            <el-table-column prop="bandwidth" label="带宽" width="100">
              <template #default="{ row }">
                {{ row.bandwidth }} Mbps
              </template>
            </el-table-column>
            <el-table-column prop="billing_method" label="计费方式" width="120" />
            <el-table-column prop="platform" label="平台" width="100" />
            <el-table-column prop="account" label="云账号" width="150" />
            <el-table-column prop="bound_resource" label="绑定资源" width="150" />
            <el-table-column prop="project" label="项目" width="120" />
            <el-table-column prop="region" label="区域" width="120" />
            <el-table-column label="操作" width="200">
              <template #default="{ row }">
                <el-button size="small" @click="handleEdit(row)">编辑</el-button>
                <el-button size="small" type="primary" @click="handleBind(row)">绑定</el-button>
                <el-button size="small" @click="handleUnbind(row)">解绑</el-button>
                <el-button size="small" type="danger" @click="handleDelete(row)">释放</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="本地idc" name="local_idc">
          <el-table :data="localIdcEips" v-loading="loading">
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
            <el-table-column prop="tags" label="标签" width="150" />
            <el-table-column prop="ip" label="IP" width="150" />
            <el-table-column prop="bandwidth" label="带宽" width="100">
              <template #default="{ row }">
                {{ row.bandwidth }} Mbps
              </template>
            </el-table-column>
            <el-table-column prop="billing_method" label="计费方式" width="120" />
            <el-table-column prop="platform" label="平台" width="100" />
            <el-table-column prop="account" label="云账号" width="150" />
            <el-table-column prop="bound_resource" label="绑定资源" width="150" />
            <el-table-column prop="project" label="项目" width="120" />
            <el-table-column prop="region" label="区域" width="120" />
            <el-table-column label="操作" width="200">
              <template #default="{ row }">
                <el-button size="small" @click="handleEdit(row)">编辑</el-button>
                <el-button size="small" type="primary" @click="handleBind(row)">绑定</el-button>
                <el-button size="small" @click="handleUnbind(row)">解绑</el-button>
                <el-button size="small" type="danger" @click="handleDelete(row)">释放</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="公有云" name="public_cloud">
          <el-table :data="publicCloudEips" v-loading="loading">
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
            <el-table-column prop="tags" label="标签" width="150" />
            <el-table-column prop="ip" label="IP" width="150" />
            <el-table-column prop="bandwidth" label="带宽" width="100">
              <template #default="{ row }">
                {{ row.bandwidth }} Mbps
              </template>
            </el-table-column>
            <el-table-column prop="billing_method" label="计费方式" width="120" />
            <el-table-column prop="platform" label="平台" width="100" />
            <el-table-column prop="account" label="云账号" width="150" />
            <el-table-column prop="bound_resource" label="绑定资源" width="150" />
            <el-table-column prop="project" label="项目" width="120" />
            <el-table-column prop="region" label="区域" width="120" />
            <el-table-column label="操作" width="200">
              <template #default="{ row }">
                <el-button size="small" @click="handleEdit(row)">编辑</el-button>
                <el-button size="small" type="primary" @click="handleBind(row)">绑定</el-button>
                <el-button size="small" @click="handleUnbind(row)">解绑</el-button>
                <el-button size="small" type="danger" @click="handleDelete(row)">释放</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- Detail Modal -->
    <el-dialog v-model="detailDialogVisible" title="弹性公网IP详情" width="600px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ selectedEip?.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ selectedEip?.name }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ selectedEip?.status }}</el-descriptions-item>
        <el-descriptions-item label="标签">{{ selectedEip?.tags }}</el-descriptions-item>
        <el-descriptions-item label="IP地址">{{ selectedEip?.ip }}</el-descriptions-item>
        <el-descriptions-item label="带宽">{{ selectedEip?.bandwidth }} Mbps</el-descriptions-item>
        <el-descriptions-item label="计费方式">{{ selectedEip?.billing_method }}</el-descriptions-item>
        <el-descriptions-item label="平台">{{ selectedEip?.platform }}</el-descriptions-item>
        <el-descriptions-item label="云账号">{{ selectedEip?.account }}</el-descriptions-item>
        <el-descriptions-item label="绑定资源">{{ selectedEip?.bound_resource }}</el-descriptions-item>
        <el-descriptions-item label="项目">{{ selectedEip?.project }}</el-descriptions-item>
        <el-descriptions-item label="区域">{{ selectedEip?.region }}</el-descriptions-item>
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

interface EIP {
  id: string
  name: string
  status: string
  tags: string
  ip: string
  bandwidth: number
  billing_method: string
  platform: string
  account: string
  bound_resource: string
  project: string
  region: string
}

const allEips = ref<EIP[]>([])
const localIdcEips = ref<EIP[]>([])
const publicCloudEips = ref<EIP[]>([])
const loading = ref(false)
const activeTab = ref('all')
const detailDialogVisible = ref(false)
const selectedEip = ref<EIP | null>(null)

const getStatusType = (status: string) => {
  switch (status.toLowerCase()) {
    case 'available':
      return 'success'
    case 'in-use':
      return 'primary'
    case 'binding':
    case 'pending':
      return 'warning'
    case 'error':
      return 'danger'
    default:
      return 'info'
  }
}

const loadEips = async () => {
  loading.value = true
  try {
    // Mock data
    const eipsData: EIP[] = [
      {
        id: 'eip-1',
        name: 'eip-web-01',
        status: 'In-Use',
        tags: 'web,prod',
        ip: '47.98.123.45',
        bandwidth: 100,
        billing_method: '按量付费',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        bound_resource: 'vm-web-01',
        project: 'Project A',
        region: 'cn-hangzhou'
      },
      {
        id: 'eip-2',
        name: 'eip-app-01 (本地IDC)',
        status: 'Available',
        tags: 'app,idc',
        ip: '10.10.10.100',
        bandwidth: 200,
        billing_method: '包年包月',
        platform: '本地IDC',
        account: 'Local IDC Account',
        bound_resource: '-',
        project: 'Project B',
        region: '本地机房'
      },
      {
        id: 'eip-3',
        name: 'eip-db-01',
        status: 'In-Use',
        tags: 'db',
        ip: '119.28.123.67',
        bandwidth: 50,
        billing_method: '按量付费',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        bound_resource: 'vm-db-01',
        project: 'Project A',
        region: 'cn-shanghai'
      },
      {
        id: 'eip-4',
        name: 'eip-api-01',
        status: 'Available',
        tags: 'api',
        ip: '120.55.45.89',
        bandwidth: 100,
        billing_method: '包年包月',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        bound_resource: '-',
        project: 'Project C',
        region: 'cn-beijing'
      }
    ]
    allEips.value = eipsData
    localIdcEips.value = eipsData.filter(e => e.platform.toLowerCase().includes('idc') || e.region.includes('本地'))
    publicCloudEips.value = eipsData.filter(e => !e.platform.toLowerCase().includes('idc') && !e.region.includes('本地'))
  } catch (e) {
    console.error(e)
    allEips.value = []
    localIdcEips.value = []
    publicCloudEips.value = []
  } finally {
    loading.value = false
  }
}

const viewDetail = (row: EIP) => {
  selectedEip.value = row
  detailDialogVisible.value = true
}

const handleEdit = (row: EIP) => {
  ElMessage.info(`编辑弹性IP: ${row.name}`)
}

const handleBind = (row: EIP) => {
  ElMessage.info(`绑定弹性IP: ${row.name}`)
}

const handleUnbind = (row: EIP) => {
  ElMessage.info(`解绑弹性IP: ${row.name}`)
}

const handleDelete = async (row: EIP) => {
  try {
    await ElMessageBox.confirm(`确认释放弹性公网IP ${row.name}？`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    allEips.value = allEips.value.filter(e => e.id !== row.id)
    localIdcEips.value = localIdcEips.value.filter(e => e.id !== row.id)
    publicCloudEips.value = publicCloudEips.value.filter(e => e.id !== row.id)
    ElMessage.success('释放成功')
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadEips()
})
</script>

<style scoped>
.eips-page {
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