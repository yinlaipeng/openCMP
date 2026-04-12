<template>
  <div class="lb-acls-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">访问控制列表</span>
        </div>
      </template>

      <el-table :data="acls" v-loading="loading">
        <el-table-column label="策略组名称" width="180">
          <template #default="{ row }">
            <el-link @click="viewDetail(row)">{{ row.name }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="source" label="源地址/备注" width="200" />
        <el-table-column prop="listeners_count" label="关联监听(数量)" width="150">
          <template #default="{ row }">
            {{ row.listeners_count }} 个
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="platform" label="平台" width="100" />
        <el-table-column prop="account" label="云账号" width="150" />
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column prop="updated_at" label="更新时间" width="180" />
        <el-table-column prop="share_scope" label="共享范围" width="100" />
        <el-table-column prop="region" label="区域" width="120" />
        <el-table-column prop="project" label="项目" width="120" />
        <el-table-column label="操作" width="150">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- Detail Modal -->
    <el-dialog v-model="detailDialogVisible" title="访问控制详情" width="600px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ selectedACL?.id }}</el-descriptions-item>
        <el-descriptions-item label="策略组名称">{{ selectedACL?.name }}</el-descriptions-item>
        <el-descriptions-item label="源地址/备注">{{ selectedACL?.source }}</el-descriptions-item>
        <el-descriptions-item label="关联监听">{{ selectedACL?.listeners_count }} 个</el-descriptions-item>
        <el-descriptions-item label="状态">{{ selectedACL?.status }}</el-descriptions-item>
        <el-descriptions-item label="平台">{{ selectedACL?.platform }}</el-descriptions-item>
        <el-descriptions-item label="云账号">{{ selectedACL?.account }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ selectedACL?.created_at }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ selectedACL?.updated_at }}</el-descriptions-item>
        <el-descriptions-item label="共享范围">{{ selectedACL?.share_scope }}</el-descriptions-item>
        <el-descriptions-item label="区域">{{ selectedACL?.region }}</el-descriptions-item>
        <el-descriptions-item label="项目">{{ selectedACL?.project }}</el-descriptions-item>
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

interface ACL {
  id: string
  name: string
  source: string
  listeners_count: number
  status: string
  platform: string
  account: string
  created_at: string
  updated_at: string
  share_scope: string
  region: string
  project: string
}

const acls = ref<ACL[]>([])
const loading = ref(false)
const detailDialogVisible = ref(false)
const selectedACL = ref<ACL | null>(null)

const getStatusType = (status: string) => {
  switch (status.toLowerCase()) {
    case 'active':
    case 'enabled':
      return 'success'
    case 'pending':
      return 'warning'
    case 'disabled':
    case 'error':
      return 'danger'
    default:
      return 'info'
  }
}

const loadACLs = async () => {
  loading.value = true
  try {
    // Mock data
    acls.value = [
      {
        id: 'acl-1',
        name: '访问控制策略组 1',
        source: '192.168.1.0/24 (内网)',
        listeners_count: 3,
        status: 'Active',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        created_at: '2024-01-01 10:00:00',
        updated_at: '2024-01-05 12:00:00',
        share_scope: '私有',
        region: 'cn-hangzhou',
        project: 'Project A'
      },
      {
        id: 'acl-2',
        name: '访问控制策略组 2',
        source: '0.0.0.0/0 (全部)',
        listeners_count: 1,
        status: 'Active',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        created_at: '2024-01-02 10:00:00',
        updated_at: '2024-01-02 10:00:00',
        share_scope: '共享',
        region: 'cn-shanghai',
        project: 'Project B'
      },
      {
        id: 'acl-3',
        name: '访问控制策略组 3',
        source: '10.0.0.0/8 (测试网段)',
        listeners_count: 0,
        status: 'Disabled',
        platform: '阿里云',
        account: 'Aliyun Account 1',
        created_at: '2024-01-03 10:00:00',
        updated_at: '2024-01-10 15:00:00',
        share_scope: '私有',
        region: 'cn-beijing',
        project: 'Project A'
      }
    ]
  } catch (e) {
    console.error(e)
    acls.value = []
  } finally {
    loading.value = false
  }
}

const viewDetail = (row: ACL) => {
  selectedACL.value = row
  detailDialogVisible.value = true
}

const handleEdit = (row: ACL) => {
  ElMessage.info(`编辑访问控制策略: ${row.name}`)
}

const handleDelete = async (row: ACL) => {
  try {
    await ElMessageBox.confirm(`确认删除访问控制策略 ${row.name}？`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    acls.value = acls.value.filter(a => a.id !== row.id)
    ElMessage.success('删除成功')
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadACLs()
})
</script>

<style scoped>
.lb-acls-page {
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