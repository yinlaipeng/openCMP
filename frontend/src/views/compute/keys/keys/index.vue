<template>
  <div class="ssh-keys-container">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">密钥</span>
        </div>
      </template>

      <el-table
        :data="sshKeys"
        v-loading="loading"
        style="width: 100%"
        row-key="id"
      >
        <el-table-column prop="name" label="名称" width="200" fixed="left" />
        <el-table-column prop="share_scope" label="共享范围" width="120" />
        <el-table-column prop="public_key_content" label="公钥内容" width="300" show-overflow-tooltip />
        <el-table-column prop="fingerprint" label="指纹" width="200" show-overflow-tooltip />
        <el-table-column prop="type" label="类型" width="100" />
        <el-table-column prop="associated_vm_count" label="关联虚拟机数" width="130" />
        <el-table-column label="操作" fixed="right" width="250">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" @click="handleView(row)">查看</el-button>
            <el-button size="small" @click="handleImport(row)">导入</el-button>
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
const sshKeys = ref<any[]>([])

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

// Methods
const fetchData = async () => {
  loading.value = true
  try {
    // Mock data for now
    sshKeys.value = [
      {
        id: 'key-1',
        name: 'web-server-key',
        share_scope: 'Project',
        public_key_content: 'ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQD...',
        fingerprint: 'SHA256:a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2',
        type: 'RSA',
        associated_vm_count: 5
      },
      {
        id: 'key-2',
        name: 'app-server-key',
        share_scope: 'Domain',
        public_key_content: 'ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQJ...',
        fingerprint: 'SHA256:z9y8x7w6v5u4t3s2r1q0p9o8n7m6l5k4j3i2h1g0f9e8',
        type: 'RSA',
        associated_vm_count: 10
      }
    ]
    pagination.total = sshKeys.value.length
  } catch (error) {
    console.error('Failed to fetch SSH keys:', error)
    ElMessage.error('获取密钥列表失败')
  } finally {
    loading.value = false
  }
}

const handleEdit = (row: any) => {
  console.log('Edit SSH key:', row)
}

const handleView = (row: any) => {
  console.log('View SSH key:', row)
}

const handleImport = (row: any) => {
  console.log('Import SSH key:', row)
}

const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除SSH密钥 "${row.name}" 吗？`,
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
    console.error('Failed to delete SSH key:', error)
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
.ssh-keys-container {
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