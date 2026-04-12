<template>
  <div class="system-images-container">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">系统镜像</span>
        </div>
      </template>

      <el-table
        :data="images"
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
        <el-table-column prop="format" label="格式" width="100" />
        <el-table-column prop="os_name" label="操作系统" width="150" />
        <el-table-column prop="size" label="大小(GB)" width="100" />
        <el-table-column prop="share_scope" label="共享范围" width="120" />
        <el-table-column prop="project_id" label="项目" width="150" />
        <el-table-column prop="cpu_arch" label="CPU架构" width="120" />
        <el-table-column prop="image_type" label="镜像类型" width="120" />
        <el-table-column label="操作" fixed="right" width="200">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
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
import type { Image } from '@/api/compute'

// Data
const loading = ref(false)
const images = ref<Image[]>([])

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
    case 'Unavailable':
      return 'danger'
    case 'Creating':
      return 'warning'
    default:
      return 'info'
  }
}

const fetchData = async () => {
  loading.value = true
  try {
    // This would be the real API call
    // const response = await getSystemImages({
    //   page: pagination.page,
    //   page_size: pagination.pageSize
    // })

    // Mock data for now
    images.value = [
      {
        id: 'img-1',
        name: 'Ubuntu Server 20.04 LTS',
        description: 'Official Ubuntu Server 20.04 LTS image',
        os_name: 'Ubuntu',
        status: 'Available',
        size: 10
      },
      {
        id: 'img-2',
        name: 'CentOS 7.9',
        description: 'Official CentOS 7.9 image',
        os_name: 'CentOS',
        status: 'Available',
        size: 8
      },
      {
        id: 'img-3',
        name: 'Windows Server 2019',
        description: 'Official Windows Server 2019 image',
        os_name: 'Windows',
        status: 'Creating',
        size: 40
      }
    ]
    pagination.total = images.value.length
  } catch (error) {
    console.error('Failed to fetch system images:', error)
    ElMessage.error('获取系统镜像列表失败')
  } finally {
    loading.value = false
  }
}

const handleEdit = (row: Image) => {
  console.log('Edit image:', row)
}

const handleDelete = async (row: Image) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除镜像 "${row.name}" 吗？`,
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    // Real API call would go here
    // await deleteImage(row.id)
    ElMessage.success('删除成功')
    fetchData()
  } catch (error) {
    console.error('Failed to delete image:', error)
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
.system-images-container {
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