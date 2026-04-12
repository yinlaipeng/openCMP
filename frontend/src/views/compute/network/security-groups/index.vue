<template>
  <div class="security-groups-container">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">安全组</span>
        </div>
      </template>

      <el-table
        :data="securityGroups"
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
        <el-table-column prop="associated_instances" label="关联实例数" width="130" />
        <el-table-column prop="share_scope" label="共享范围" width="120" />
        <el-table-column prop="project" label="项目" width="150" />
        <el-table-column prop="region" label="区域" width="120" />
        <el-table-column prop="vpc" label="VPC" width="150" />
        <el-table-column label="操作" fixed="right" width="200">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" @click="handleRules(row)">规则</el-button>
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
const securityGroups = ref<any[]>([])

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
    securityGroups.value = [
      {
        id: 'sg-1',
        name: 'web-security-group',
        status: 'Available',
        platform: 'Aliyun',
        cloud_account: 'Aliyun-Account-1',
        associated_instances: 5,
        share_scope: 'Project',
        project: 'Project-A',
        region: 'cn-hangzhou',
        vpc: 'vpc-web-01'
      },
      {
        id: 'sg-2',
        name: 'db-security-group',
        status: 'Available',
        platform: 'Tencent',
        cloud_account: 'Tencent-Account-1',
        associated_instances: 2,
        share_scope: 'Domain',
        project: 'Project-B',
        region: 'ap-beijing',
        vpc: 'vpc-db-01'
      }
    ]
    pagination.total = securityGroups.value.length
  } catch (error) {
    console.error('Failed to fetch security groups:', error)
    ElMessage.error('获取安全组列表失败')
  } finally {
    loading.value = false
  }
}

const handleEdit = (row: any) => {
  console.log('Edit security group:', row)
}

const handleRules = (row: any) => {
  console.log('View security group rules:', row)
}

const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除安全组 "${row.name}" 吗？`,
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
    console.error('Failed to delete security group:', error)
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
.security-groups-container {
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