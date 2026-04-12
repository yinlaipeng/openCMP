<template>
  <div class="disk-snapshots-container">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">硬盘快照</span>
        </div>
      </template>

      <el-table
        :data="diskSnapshots"
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
        <el-table-column prop="snapshot_size" label="快照大小" width="120" />
        <el-table-column prop="disk_type" label="磁盘类型" width="120" />
        <el-table-column prop="disk" label="磁盘" width="150" />
        <el-table-column prop="vm" label="虚拟机" width="150" />
        <el-table-column prop="platform" label="平台" width="120" />
        <el-table-column prop="cloud_account" label="云账户" width="150" />
        <el-table-column prop="project" label="项目" width="150" />
        <el-table-column prop="region" label="区域" width="120" />
        <el-table-column label="操作" fixed="right" width="200">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" @click="handleRestore(row)">恢复</el-button>
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
const diskSnapshots = ref<any[]>([])

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

// Methods
const getStatusType = (status: string) => {
  switch (status) {
    case 'Ready':
      return 'success'
    case 'Creating':
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
    diskSnapshots.value = [
      {
        id: 'snap-1',
        name: 'disk-snap-20240101',
        status: 'Ready',
        snapshot_size: 100,
        disk_type: 'System',
        disk: 'disk-system-01',
        vm: 'vm-1001',
        platform: 'Aliyun',
        cloud_account: 'Aliyun-Account-1',
        project: 'Project-A',
        region: 'cn-hangzhou'
      },
      {
        id: 'snap-2',
        name: 'disk-snap-20240102',
        status: 'Creating',
        snapshot_size: 500,
        disk_type: 'Data',
        disk: 'disk-data-01',
        vm: 'vm-1002',
        platform: 'Tencent',
        cloud_account: 'Tencent-Account-1',
        project: 'Project-B',
        region: 'ap-beijing'
      }
    ]
    pagination.total = diskSnapshots.value.length
  } catch (error) {
    console.error('Failed to fetch disk snapshots:', error)
    ElMessage.error('获取硬盘快照列表失败')
  } finally {
    loading.value = false
  }
}

const handleEdit = (row: any) => {
  console.log('Edit disk snapshot:', row)
}

const handleRestore = (row: any) => {
  console.log('Restore disk snapshot:', row)
}

const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除硬盘快照 "${row.name}" 吗？`,
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
    console.error('Failed to delete disk snapshot:', error)
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
.disk-snapshots-container {
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