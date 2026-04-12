<template>
  <div class="dns-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">DNS解析列表</span>
        </div>
      </template>

      <el-table :data="dnsRecords" v-loading="loading">
        <el-table-column label="域名" width="200">
          <template #default="{ row }">
            <el-link @click="viewDetail(row)">{{ row.domain }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="tags" label="标签" width="150" />
        <el-table-column prop="platform" label="平台" width="100" />
        <el-table-column prop="type" label="解析域类型" width="120" />
        <el-table-column prop="records" label="记录数" width="100" />
        <el-table-column prop="associated_vpc" label="关联vpc" width="180" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="share_scope" label="共享范围" width="100" />
        <el-table-column prop="account" label="云账号" width="150" />
        <el-table-column prop="project" label="项目" width="120" />
        <el-table-column label="操作" width="150">
          <template #default="{ row }">
            <el-button size="small" @click="associateVpc(row)">关联vpc</el-button>
            <el-dropdown>
              <el-button size="small">
                更多 <el-icon class="el-icon--right"><arrow-down /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="syncStatus(row)">同步状态</el-dropdown-item>
                  <el-dropdown-item @click="changeDomain(row)">更改域</el-dropdown-item>
                  <el-dropdown-item @click="setShare(row)">设置共享</el-dropdown-item>
                  <el-dropdown-item divided @click="handleDelete(row)">删除</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- Detail Modal -->
    <el-dialog v-model="detailDialogVisible" title="DNS解析详情" width="800px">
      <el-tabs v-model="detailTab">
        <el-tab-pane label="详情" name="detail">
          <el-descriptions :column="2" border>
            <el-descriptions-item label="ID">{{ selectedDNS?.id }}</el-descriptions-item>
            <el-descriptions-item label="域名">{{ selectedDNS?.domain }}</el-descriptions-item>
            <el-descriptions-item label="标签">{{ selectedDNS?.tags }}</el-descriptions-item>
            <el-descriptions-item label="平台">{{ selectedDNS?.platform }}</el-descriptions-item>
            <el-descriptions-item label="解析域类型">{{ selectedDNS?.type }}</el-descriptions-item>
            <el-descriptions-item label="记录数">{{ selectedDNS?.records }}</el-descriptions-item>
            <el-descriptions-item label="关联VPC">{{ selectedDNS?.associated_vpc }}</el-descriptions-item>
            <el-descriptions-item label="状态">{{ selectedDNS?.status }}</el-descriptions-item>
            <el-descriptions-item label="共享范围">{{ selectedDNS?.share_scope }}</el-descriptions-item>
            <el-descriptions-item label="云账号">{{ selectedDNS?.account }}</el-descriptions-item>
            <el-descriptions-item label="项目">{{ selectedDNS?.project }}</el-descriptions-item>
            <el-descriptions-item label="创建时间">{{ selectedDNS?.created_at }}</el-descriptions-item>
          </el-descriptions>
        </el-tab-pane>
        <el-tab-pane label="记录" name="records">
          <el-table :data="dnsRecordList" v-loading="recordsLoading">
            <el-table-column prop="name" label="记录名称" width="150" />
            <el-table-column prop="type" label="类型" width="100" />
            <el-table-column prop="value" label="值" width="200" />
            <el-table-column prop="ttl" label="TTL" width="100" />
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="row.status === 'Active' ? 'success' : 'info'">{{ row.status }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="100">
              <template #default="{ row }">
                <el-button size="small" type="danger" @click="deleteRecord(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="操作日志" name="logs">
          <el-table :data="logs">
            <el-table-column prop="operation" label="操作" width="150" />
            <el-table-column prop="operator" label="操作员" width="150" />
            <el-table-column prop="result" label="结果" width="100" />
            <el-table-column prop="timestamp" label="时间" width="180" />
          </el-table>
        </el-tab-pane>
      </el-tabs>
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- Associate VPC Dialog -->
    <el-dialog v-model="associateVpcDialogVisible" title="关联VPC" width="400px">
      <el-form label-width="80px">
        <el-form-item label="VPC">
          <el-select v-model="selectedVpc" placeholder="请选择VPC" multiple>
            <el-option label="VPC 1" value="vpc-1" />
            <el-option label="VPC 2" value="vpc-2" />
            <el-option label="VPC 3" value="vpc-3" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="associateVpcDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitAssociateVpc">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowDown } from '@element-plus/icons-vue'

interface DNS {
  id: string
  domain: string
  tags: string
  platform: string
  type: string
  records: number
  associated_vpc: string
  status: string
  share_scope: string
  account: string
  project: string
  created_at: string
}

interface DNSRecord {
  name: string
  type: string
  value: string
  ttl: number
  status: string
}

const dnsRecords = ref<DNS[]>([])
const loading = ref(false)
const detailDialogVisible = ref(false)
const associateVpcDialogVisible = ref(false)
const selectedDNS = ref<DNS | null>(null)
const detailTab = ref('detail')
const selectedVpc = ref<string[]>([])

// Detail modal data
const dnsRecordList = ref<DNSRecord[]>([])
const recordsLoading = ref(false)
const logs = ref<any[]>([])

const getStatusType = (status: string) => {
  switch (status.toLowerCase()) {
    case 'active':
    case 'normal':
      return 'success'
    case 'pending':
    case 'configuring':
      return 'warning'
    case 'error':
    case 'failed':
      return 'danger'
    default:
      return 'info'
  }
}

const loadDNSRecords = async () => {
  loading.value = true
  try {
    // Mock data
    dnsRecords.value = [
      {
        id: 'dns-1',
        domain: 'example.com',
        tags: 'prod',
        platform: '阿里云',
        type: '公网解析',
        records: 10,
        associated_vpc: 'vpc-1, vpc-2',
        status: 'Active',
        share_scope: '私有',
        account: 'Aliyun Account 1',
        project: 'Project A',
        created_at: '2024-01-01 10:00:00'
      },
      {
        id: 'dns-2',
        domain: 'internal.example.cn',
        tags: 'internal',
        platform: '阿里云',
        type: '内网解析',
        records: 5,
        associated_vpc: 'vpc-3',
        status: 'Active',
        share_scope: '共享',
        account: 'Aliyun Account 1',
        project: 'Project B',
        created_at: '2024-01-02 10:00:00'
      },
      {
        id: 'dns-3',
        domain: 'api.example.com',
        tags: 'api',
        platform: '阿里云',
        type: '公网解析',
        records: 3,
        associated_vpc: '-',
        status: 'Pending',
        share_scope: '私有',
        account: 'Aliyun Account 1',
        project: 'Project A',
        created_at: '2024-01-03 10:00:00'
      }
    ]
  } catch (e) {
    console.error(e)
    dnsRecords.value = []
  } finally {
    loading.value = false
  }
}

const viewDetail = (row: DNS) => {
  selectedDNS.value = row
  detailDialogVisible.value = true
  detailTab.value = 'detail'

  // Mock records data
  recordsLoading.value = true
  dnsRecordList.value = [
    { name: '@', type: 'A', value: '192.168.1.100', ttl: 600, status: 'Active' },
    { name: 'www', type: 'A', value: '192.168.1.101', ttl: 600, status: 'Active' },
    { name: 'api', type: 'CNAME', value: 'api.example.com', ttl: 300, status: 'Active' },
    { name: 'mail', type: 'MX', value: 'mail.example.com', ttl: 3600, status: 'Active' }
  ]
  recordsLoading.value = false

  logs.value = [
    { operation: '创建', operator: 'admin', result: '成功', timestamp: '2024-01-01 10:00:00' },
    { operation: '添加记录', operator: 'admin', result: '成功', timestamp: '2024-01-02 12:00:00' }
  ]
}

const associateVpc = (row: DNS) => {
  selectedDNS.value = row
  selectedVpc.value = row.associated_vpc.split(',').map(v => v.trim()).filter(v => v !== '-')
  associateVpcDialogVisible.value = true
}

const submitAssociateVpc = () => {
  ElMessage.success('关联VPC成功')
  associateVpcDialogVisible.value = false
}

const syncStatus = (row: DNS) => {
  ElMessage.info(`正在同步DNS解析 ${row.domain} 的状态`)
}

const changeDomain = (row: DNS) => {
  ElMessage.info(`更改域功能开发中: ${row.domain}`)
}

const setShare = (row: DNS) => {
  ElMessage.info(`设置共享功能开发中: ${row.domain}`)
}

const deleteRecord = async (row: DNSRecord) => {
  try {
    await ElMessageBox.confirm(`确认删除DNS记录 ${row.name}？`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    dnsRecordList.value = dnsRecordList.value.filter(r => r.name !== row.name)
    ElMessage.success('删除成功')
  } catch (e) {
    console.error(e)
  }
}

const handleDelete = async (row: DNS) => {
  try {
    await ElMessageBox.confirm(`确认删除DNS解析 ${row.domain}？`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    dnsRecords.value = dnsRecords.value.filter(d => d.id !== row.id)
    ElMessage.success('删除成功')
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadDNSRecords()
})
</script>

<style scoped>
.dns-page {
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