<template>
  <div class="disks-container">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">硬盘</span>
        </div>
      </template>

      <el-table
        :data="disks"
        v-loading="loading"
        style="width: 100%"
        row-key="id"
      >
        <el-table-column prop="name" label="名称" width="150" fixed="left" />
        <el-table-column prop="disk_capacity" label="磁盘容量" width="120" />
        <el-table-column prop="max_iops" label="最大IOPS" width="120" />
        <el-table-column prop="format" label="格式" width="100" />
        <el-table-column prop="storage_type" label="存储类型" width="120" />
        <el-table-column prop="disk_type" label="磁盘类型" width="120" />
        <el-table-column prop="attached" label="挂载状态" width="120">
          <template #default="{ row }">
            <el-tag :type="row.attached ? 'success' : 'info'">
              {{ row.attached ? '已挂载' : '未挂载' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="vm" label="虚拟机" width="150" />
        <el-table-column prop="device_name" label="设备名" width="120" />
        <el-table-column prop="primary_storage" label="主存储" width="120" />
        <el-table-column prop="creation_time" label="创建时间" width="150" />
        <el-table-column prop="platform" label="平台" width="120" />
        <el-table-column prop="region" label="区域" width="120" />
        <el-table-column prop="billing_method" label="计费方式" width="120" />
        <el-table-column prop="status" label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="project" label="项目" width="150" />
        <el-table-column prop="cloud_account" label="云账户" width="150" />
        <el-table-column prop="media_type" label="媒体类型" width="120" />
        <el-table-column prop="shutdown_auto_reset" label="关机自动重置" width="150" />
        <el-table-column label="操作" fixed="right" width="200">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" @click="handleAttach(row)">挂载</el-button>
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
const disks = ref<any[]>([])

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
    case 'Deleting':
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
    disks.value = [
      {
        id: 'disk-1',
        name: 'disk-data-01',
        disk_capacity: 100,
        max_iops: 5000,
        format: 'RAW',
        storage_type: 'SSD',
        disk_type: 'System',
        attached: true,
        vm: 'vm-1001',
        device_name: '/dev/sda1',
        primary_storage: 'Yes',
        creation_time: '2024-01-15 10:30:00',
        platform: 'Aliyun',
        region: 'cn-hangzhou',
        billing_method: 'Pay-As-You-Go',
        status: 'In-Use',
        project: 'Project-A',
        cloud_account: 'Aliyun-Account-1',
        media_type: 'SSD',
        shutdown_auto_reset: 'No'
      },
      {
        id: 'disk-2',
        name: 'disk-backup-01',
        disk_capacity: 500,
        max_iops: 3000,
        format: 'Qcow2',
        storage_type: 'HDD',
        disk_type: 'Data',
        attached: false,
        vm: '-',
        device_name: '-',
        primary_storage: 'No',
        creation_time: '2024-01-16 14:20:00',
        platform: 'Tencent',
        region: 'ap-beijing',
        billing_method: 'Subscription',
        status: 'Available',
        project: 'Project-B',
        cloud_account: 'Tencent-Account-1',
        media_type: 'HDD',
        shutdown_auto_reset: 'Yes'
      }
    ]
    pagination.total = disks.value.length
  } catch (error) {
    console.error('Failed to fetch disks:', error)
    ElMessage.error('获取硬盘列表失败')
  } finally {
    loading.value = false
  }
}

const handleEdit = (row: any) => {
  console.log('Edit disk:', row)
}

const handleAttach = (row: any) => {
  console.log('Attach disk:', row)
}

const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除硬盘 "${row.name}" 吗？`,
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
    console.error('Failed to delete disk:', error)
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
.disks-container {
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