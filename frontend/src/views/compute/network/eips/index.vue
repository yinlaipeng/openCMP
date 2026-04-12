<template>
  <div class="elastic-ips-container">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">弹性公网IP</span>
        </div>
      </template>

      <el-table
        :data="elasticIPs"
        v-loading="loading"
        style="width: 100%"
        row-key="id"
      >
        <el-table-column prop="name" label="名称" width="200" fixed="left" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="ip" label="IP地址" width="150" />
        <el-table-column prop="bandwidth" label="带宽(Mbps)" width="130" />
        <el-table-column prop="billing_method" label="计费方式" width="120" />
        <el-table-column prop="platform" label="平台" width="120" />
        <el-table-column prop="cloud_account" label="云账户" width="150" />
        <el-table-column prop="bound_resource" label="绑定资源" width="150" />
        <el-table-column prop="project" label="项目" width="150" />
        <el-table-column prop="region" label="区域" width="120" />
        <el-table-column label="操作" fixed="right" width="200">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" @click="handleBind(row)">绑定</el-button>
            <el-button size="small" @click="handleUnbind(row)">解绑</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">释放</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="pagination.total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        class="pagination"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

// Data
const loading = ref(false)
const elasticIPs = ref<any[]>([])

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

// Methods
const getStatusType = (status: string) => {
  switch (status) {
    case 'Available':
      return 'success'
    case 'In-Use':
      return 'primary'
    case 'Binding':
      return 'warning'
    case 'Error':
      return 'danger'
    default:
      return 'info'
  }
}

const fetchData = async () => {
  loading.value = true
  try {
    // Mock data for now
    elasticIPs.value = [
      {
        id: 'eip-1',
        name: 'eip-web-01',
        status: 'In-Use',
        ip: '47.98.123.45',
        bandwidth: 100,
        billing_method: 'Pay-As-You-Go',
        platform: 'Aliyun',
        cloud_account: 'Aliyun-Account-1',
        bound_resource: 'vm-1001',
        project: 'Project-A',
        region: 'cn-hangzhou'
      },
      {
        id: 'eip-2',
        name: 'eip-app-01',
        status: 'Available',
        ip: '119.28.123.67',
        bandwidth: 200,
        billing_method: 'Subscription',
        platform: 'Tencent',
        cloud_account: 'Tencent-Account-1',
        bound_resource: '-',
        project: 'Project-B',
        region: 'ap-beijing'
      }
    ]
    pagination.total = elasticIPs.value.length
  } catch (error) {
    console.error('Failed to fetch elastic IPs:', error)
    ElMessage.error('获取弹性公网IP列表失败')
  } finally {
    loading.value = false
  }
}

const handleEdit = (row: any) => {
  console.log('Edit elastic IP:', row)
}

const handleBind = (row: any) => {
  console.log('Bind elastic IP:', row)
}

const handleUnbind = (row: any) => {
  console.log('Unbind elastic IP:', row)
}

const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm(
      `确定要释放弹性公网IP "${row.name}" 吗？`,
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    ElMessage.success('释放成功')
    fetchData()
  } catch (error) {
    console.error('Failed to release elastic IP:', error)
    if (error !== 'cancel') {
      ElMessage.error('释放失败')
    }
  }
}

const handleSizeChange = (val: number) => {
  pagination.pageSize = val
  pagination.page = 1
  fetchData()
}

const handleCurrentChange = (val: number) => {
  pagination.page = val
  fetchData()
}

// Lifecycle
onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.elastic-ips-container {
  padding: 20px;
}

.page-card {
  min-height: calc(100vh - 120px);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-title {
  font-size: 18px;
  font-weight: bold;
}

.pagination {
  margin-top: 20px;
  text-align: center;
}
</style>