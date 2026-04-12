<template>
  <div class="ip-subnets-container">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">IP子网</span>
        </div>
      </template>

      <el-table
        :data="ipSubnets"
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
        <el-table-column prop="type" label="类型" width="100" />
        <el-table-column prop="auto_allocation" label="自动分配" width="120">
          <template #default="{ row }">
            <el-tag :type="row.auto_allocation ? 'success' : 'info'">
              {{ row.auto_allocation ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="ip_address" label="IP地址段" width="150" />
        <el-table-column prop="usage" label="使用情况" width="120" />
        <el-table-column prop="schedule_labels" label="调度标签" width="150" />
        <el-table-column prop="platform" label="平台" width="120" />
        <el-table-column prop="project" label="项目" width="150" />
        <el-table-column prop="region" label="区域" width="120" />
        <el-table-column label="操作" fixed="right" width="200">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" @click="handleAllocation(row)">分配</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
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
const ipSubnets = ref<any[]>([])

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
    case 'Occupied':
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
    ipSubnets.value = [
      {
        id: 'subnet-1',
        name: 'subnet-web-01',
        status: 'Available',
        type: 'IPv4',
        auto_allocation: true,
        ip_address: '192.168.1.0/24',
        usage: '50%',
        schedule_labels: 'web-zone',
        platform: 'Aliyun',
        project: 'Project-A',
        region: 'cn-hangzhou'
      },
      {
        id: 'subnet-2',
        name: 'subnet-db-01',
        status: 'Occupied',
        type: 'IPv4',
        auto_allocation: true,
        ip_address: '192.168.2.0/24',
        usage: '90%',
        schedule_labels: 'db-zone',
        platform: 'Tencent',
        project: 'Project-B',
        region: 'ap-beijing'
      }
    ]
    pagination.total = ipSubnets.value.length
  } catch (error) {
    console.error('Failed to fetch IP subnets:', error)
    ElMessage.error('获取IP子网列表失败')
  } finally {
    loading.value = false
  }
}

const handleEdit = (row: any) => {
  console.log('Edit IP subnet:', row)
}

const handleAllocation = (row: any) => {
  console.log('Allocate IP subnet:', row)
}

const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除IP子网 "${row.name}" 吗？`,
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    ElMessage.success('删除成功')
    fetchData()
  } catch (error) {
    console.error('Failed to delete IP subnet:', error)
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
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
.ip-subnets-container {
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