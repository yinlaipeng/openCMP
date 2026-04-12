<template>
  <div class="snapshot-policies-container">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">自动快照策略</span>
        </div>
      </template>

      <el-table
        :data="snapshotPolicies"
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
        <el-table-column prop="platform" label="平台" width="120" />
        <el-table-column prop="cloud_account" label="云账户" width="150" />
        <el-table-column prop="resource_type" label="资源类型" width="120" />
        <el-table-column prop="associated_resource_count" label="关联资源数" width="130" />
        <el-table-column prop="policy" label="策略" width="200" />
        <el-table-column prop="region" label="区域" width="120" />
        <el-table-column prop="project" label="项目" width="150" />
        <el-table-column label="操作" fixed="right" width="200">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" @click="handleEnable(row)" :type="row.status === 'Active' ? 'warning' : 'success'">
              {{ row.status === 'Active' ? '禁用' : '启用' }}
            </el-button>
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
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

// Data
const loading = ref(false)
const snapshotPolicies = ref<any[]>([])

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

// Methods
const getStatusType = (status: string) => {
  switch (status) {
    case 'Active':
      return 'success'
    case 'Inactive':
      return 'info'
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
    snapshotPolicies.value = [
      {
        id: 'pol-1',
        name: 'Daily-Snap-Policy',
        status: 'Active',
        platform: 'Aliyun',
        cloud_account: 'Aliyun-Account-1',
        resource_type: 'Disk',
        associated_resource_count: 10,
        policy: '每天凌晨2点执行',
        region: 'cn-hangzhou',
        project: 'Project-A'
      },
      {
        id: 'pol-2',
        name: 'Weekly-Host-Snap-Policy',
        status: 'Inactive',
        platform: 'Tencent',
        cloud_account: 'Tencent-Account-1',
        resource_type: 'Host',
        associated_resource_count: 5,
        policy: '每周日凌晨3点执行',
        region: 'ap-beijing',
        project: 'Project-B'
      }
    ]
    pagination.total = snapshotPolicies.value.length
  } catch (error) {
    console.error('Failed to fetch snapshot policies:', error)
    ElMessage.error('获取自动快照策略列表失败')
  } finally {
    loading.value = false
  }
}

const handleEdit = (row: any) => {
  console.log('Edit snapshot policy:', row)
}

const handleEnable = async (row: any) => {
  try {
    const newStatus = row.status === 'Active' ? 'Inactive' : 'Active'
    const actionText = newStatus === 'Active' ? '启用' : '禁用'

    await ElMessageBox.confirm(
      `确定要${actionText}自动快照策略 "${row.name}" 吗？`,
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    row.status = newStatus
    ElMessage.success(`${actionText}成功`)
  } catch (error) {
    console.error(`Failed to ${row.status === 'Active' ? 'disable' : 'enable'} snapshot policy:`, error)
    if (error !== 'cancel') {
      ElMessage.error(`${row.status === 'Active' ? '禁用' : '启用'}失败`)
    }
  }
}

const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除自动快照策略 "${row.name}" 吗？`,
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
    console.error('Failed to delete snapshot policy:', error)
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
.snapshot-policies-container {
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