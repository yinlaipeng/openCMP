<template>
  <div class="nat-gateways-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">NAT网关列表</span>
        </div>
      </template>

      <el-table :data="natGateways" v-loading="loading">
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
        <el-table-column prop="type" label="类型" width="120" />
        <el-table-column prop="tags" label="标签" width="150" />
        <el-table-column prop="rules" label="规则" width="100">
          <template #default="{ row }">
            <el-link @click="manageRules(row)">{{ row.rules }} 条</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="billing_method" label="计费方式" width="120" />
        <el-table-column prop="platform" label="平台" width="100" />
        <el-table-column prop="account" label="云账号" width="150" />
        <el-table-column prop="vpc" label="所属vpc" width="180" />
        <el-table-column prop="region" label="区域" width="120" />
        <el-table-column label="操作" width="150">
          <template #default="{ row }">
            <el-button size="small" @click="manageRules(row)">管理规则</el-button>
            <el-dropdown>
              <el-button size="small">
                更多 <el-icon class="el-icon--right"><arrow-down /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="handleEdit(row)">编辑</el-dropdown-item>
                  <el-dropdown-item divided @click="handleDelete(row)">删除</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- Detail Modal -->
    <el-dialog v-model="detailDialogVisible" title="NAT网关详情" width="600px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ selectedNATGateway?.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ selectedNATGateway?.name }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ selectedNATGateway?.status }}</el-descriptions-item>
        <el-descriptions-item label="类型">{{ selectedNATGateway?.type }}</el-descriptions-item>
        <el-descriptions-item label="标签">{{ selectedNATGateway?.tags }}</el-descriptions-item>
        <el-descriptions-item label="规则数">{{ selectedNATGateway?.rules }} 条</el-descriptions-item>
        <el-descriptions-item label="计费方式">{{ selectedNATGateway?.billing_method }}</el-descriptions-item>
        <el-descriptions-item label="平台">{{ selectedNATGateway?.platform }}</el-descriptions-item>
        <el-descriptions-item label="云账号">{{ selectedNATGateway?.account }}</el-descriptions-item>
        <el-descriptions-item label="所属VPC">{{ selectedNATGateway?.vpc }}</el-descriptions-item>
        <el-descriptions-item label="区域">{{ selectedNATGateway?.region }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ selectedNATGateway?.created_at }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- Rules Modal -->
    <el-dialog v-model="rulesDialogVisible" title="管理NAT规则" width="800px">
      <el-table :data="natRules" v-loading="rulesLoading">
        <el-table-column prop="name" label="规则名称" width="150" />
        <el-table-column prop="type" label="类型" width="100" />
        <el-table-column prop="external_ip" label="外部IP" width="150" />
        <el-table-column prop="internal_ip" label="内部IP" width="150" />
        <el-table-column prop="port" label="端口" width="100" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'Active' ? 'success' : 'info'">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-button size="small" type="danger" @click="deleteRule(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <template #footer>
        <el-button @click="rulesDialogVisible = false">关闭</el-button>
        <el-button type="primary" @click="addRule">添加规则</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowDown } from '@element-plus/icons-vue'

interface NATGateway {
  id: string
  name: string
  status: string
  type: string
  tags: string
  rules: number
  billing_method: string
  platform: string
  account: string
  vpc: string
  region: string
  created_at: string
}

interface NATRule {
  name: string
  type: string
  external_ip: string
  internal_ip: string
  port: string
  status: string
}

const natGateways = ref<NATGateway[]>([])
const loading = ref(false)
const detailDialogVisible = ref(false)
const rulesDialogVisible = ref(false)
const rulesLoading = ref(false)
const selectedNATGateway = ref<NATGateway | null>(null)
const natRules = ref<NATRule[]>([])

const getStatusType = (status: string) => {
  switch (status.toLowerCase()) {
    case 'available':
    case 'active':
      return 'success'
    case 'pending':
    case 'creating':
      return 'warning'
    case 'error':
    case 'failed':
      return 'danger'
    default:
      return 'info'
  }
}

const loadNATGateways = async () => {
  loading.value = true
  try {
    // Mock data
    natGateways.value = [
      {
        id: 'nat-1',
        name: 'NAT网关 1',
        status: 'Available',
        type: '公网NAT',
        tags: 'prod',
        rules: 5,
        billing_method: '按量付费',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        vpc: 'vpc-1',
        region: 'cn-hangzhou',
        created_at: '2024-01-01 10:00:00'
      },
      {
        id: 'nat-2',
        name: 'NAT网关 2',
        status: 'Available',
        type: '公网NAT',
        tags: 'dev',
        rules: 3,
        billing_method: '包年包月',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        vpc: 'vpc-2',
        region: 'cn-shanghai',
        created_at: '2024-01-02 10:00:00'
      },
      {
        id: 'nat-3',
        name: 'NAT网关 3',
        status: 'Creating',
        type: '私网NAT',
        tags: 'test',
        rules: 0,
        billing_method: '按量付费',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        vpc: 'vpc-3',
        region: 'cn-beijing',
        created_at: '2024-01-03 10:00:00'
      }
    ]
  } catch (e) {
    console.error(e)
    natGateways.value = []
  } finally {
    loading.value = false
  }
}

const viewDetail = (row: NATGateway) => {
  selectedNATGateway.value = row
  detailDialogVisible.value = true
}

const manageRules = (row: NATGateway) => {
  selectedNATGateway.value = row
  rulesDialogVisible.value = true
  // Mock rules data
  natRules.value = [
    { name: 'SNAT规则 1', type: 'SNAT', external_ip: '47.98.123.45', internal_ip: '10.0.1.0/24', port: '-', status: 'Active' },
    { name: 'DNAT规则 1', type: 'DNAT', external_ip: '47.98.123.46', internal_ip: '10.0.1.100', port: '80', status: 'Active' },
    { name: 'DNAT规则 2', type: 'DNAT', external_ip: '47.98.123.47', internal_ip: '10.0.1.101', port: '443', status: 'Active' }
  ]
}

const addRule = () => {
  ElMessage.info('添加NAT规则功能开发中')
}

const deleteRule = async (row: NATRule) => {
  try {
    await ElMessageBox.confirm(`确认删除规则 ${row.name}？`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    natRules.value = natRules.value.filter(r => r.name !== row.name)
    ElMessage.success('删除成功')
  } catch (e) {
    console.error(e)
  }
}

const handleEdit = (row: NATGateway) => {
  ElMessage.info(`编辑NAT网关: ${row.name}`)
}

const handleDelete = async (row: NATGateway) => {
  try {
    await ElMessageBox.confirm(`确认删除NAT网关 ${row.name}？`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    natGateways.value = natGateways.value.filter(n => n.id !== row.id)
    ElMessage.success('删除成功')
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadNATGateways()
})
</script>

<style scoped>
.nat-gateways-page {
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