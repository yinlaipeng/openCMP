<template>
  <div class="system-images-container">
    <div class="page-header">
      <h2>系统镜像</h2>
      <div class="toolbar">
        <el-button link type="primary" @click="handleView">
          <el-icon><View /></el-icon>
          查看
        </el-button>
        <el-button type="primary" @click="handleUpload">
          <el-icon><Upload /></el-icon>
          上传
        </el-button>
        <el-button @click="handleCommunityMirror">
          <el-icon><Link /></el-icon>
          社区镜像
        </el-button>
        <el-dropdown trigger="click" @command="handleBatchCommand" :disabled="selectedImages.length === 0">
          <el-button :disabled="selectedImages.length === 0">
            批量操作<el-icon class="el-icon--right"><ArrowDown /></el-icon>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="share">批量共享</el-dropdown-item>
              <el-dropdown-item command="unshare">取消共享</el-dropdown-item>
              <el-dropdown-item command="delete" divided>批量删除</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        <el-button @click="handleTags">
          <el-icon><PriceTag /></el-icon>
          标签
        </el-button>
      </div>
    </div>

    <!-- 公有云/私有云 Tabs -->
    <el-tabs v-model="activeTab" class="image-tabs" @tab-change="handleTabChange">
      <el-tab-pane label="公有云" name="public">
        <el-card class="filter-card">
      <el-form :inline="true" :model="filters" @submit.prevent="fetchData">
          <el-form-item label="名称">
            <el-input v-model="filters.name" placeholder="搜索镜像名称" clearable style="width: 180px" />
          </el-form-item>
          <el-form-item label="操作系统">
            <el-select v-model="filters.os_name" placeholder="选择操作系统" clearable style="width: 140px">
              <el-option label="CentOS" value="CentOS" />
              <el-option label="Ubuntu" value="Ubuntu" />
              <el-option label="Windows" value="Windows" />
              <el-option label="Debian" value="Debian" />
              <el-option label="RedHat" value="RedHat" />
              <el-option label="Alibaba Cloud" value="Alibaba" />
            </el-select>
          </el-form-item>
          <el-form-item label="平台">
            <el-select v-model="filters.platform" placeholder="选择平台" clearable style="width: 120px">
              <el-option label="阿里云" value="aliyun" />
              <el-option label="腾讯云" value="tencent" />
              <el-option label="AWS" value="aws" />
              <el-option label="Azure" value="azure" />
            </el-select>
          </el-form-item>
          <el-form-item label="格式">
            <el-select v-model="filters.format" placeholder="选择格式" clearable style="width: 120px">
              <el-option label="qcow2" value="qcow2" />
              <el-option label="raw" value="raw" />
              <el-option label="vhd" value="vhd" />
              <el-option label="iso" value="iso" />
            </el-select>
          </el-form-item>
          <el-form-item label="状态">
            <el-select v-model="filters.status" placeholder="选择状态" clearable style="width: 120px">
              <el-option label="可用" value="Available" />
              <el-option label="创建中" value="Creating" />
              <el-option label="不可用" value="Unavailable" />
            </el-select>
          </el-form-item>
          <el-form-item label="架构">
            <el-select v-model="filters.cpu_arch" placeholder="选择架构" clearable style="width: 120px">
              <el-option label="x86_64" value="x86_64" />
              <el-option label="arm64" value="arm64" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="fetchData">查询</el-button>
            <el-button @click="resetFilters">重置</el-button>
          </el-form-item>
        </el-form>
    </el-card>

        <el-table
          :data="images"
          v-loading="loading"
          style="width: 100%"
          row-key="id"
          @selection-change="handleSelectionChange"
        >
          <el-table-column type="selection" width="55" />
          <el-table-column label="名称" min-width="200">
            <template #default="{ row }">
              <el-link type="primary" :underline="false" @click="handleDetails(row)">
                {{ row.name }}
              </el-link>
            </template>
          </el-table-column>

          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="getStatusType(row.status)">{{ getStatusLabel(row.status) }}</el-tag>
            </template>
          </el-table-column>

          <el-table-column prop="format" label="格式" width="80" />

          <el-table-column label="操作系统" width="150">
            <template #default="{ row }">
              <div class="os-cell">
                <span class="os-name">{{ row.os_name }}</span>
                <span v-if="row.os_version" class="os-version">{{ row.os_version }}</span>
              </div>
            </template>
          </el-table-column>

          <el-table-column label="大小" width="100">
            <template #default="{ row }">
              {{ formatSize(row.size) }}
            </template>
          </el-table-column>

          <el-table-column prop="cpu_arch" label="架构" width="100" />

          <el-table-column label="平台/云账号" width="150">
            <template #default="{ row }">
              <div class="platform-cell">
                <el-tag size="small" :type="getPlatformType(row.platform)">
                  {{ getPlatformLabel(row.platform) }}
                </el-tag>
                <span class="account-name">{{ row.account_name || '-' }}</span>
              </div>
            </template>
          </el-table-column>

          <el-table-column prop="region_id" label="区域" width="120" />

          <el-table-column prop="image_type" label="类型" width="100">
            <template #default="{ row }">
              <el-tag size="small" type="info">{{ row.image_type || '系统' }}</el-tag>
            </template>
          </el-table-column>

          <el-table-column prop="share_scope" label="共享范围" width="100" />

          <el-table-column prop="project_name" label="项目" width="120" />

          <el-table-column label="操作" width="160" fixed="right">
            <template #default="{ row }">
              <el-dropdown trigger="click" @command="(cmd: string) => handleActionCommand(cmd, row)">
                <el-button size="small" link type="primary">
                  操作<el-icon class="el-icon--right"><ArrowDown /></el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="details">查看详情</el-dropdown-item>
                    <el-dropdown-item command="edit">编辑</el-dropdown-item>
                    <el-dropdown-item command="share">共享</el-dropdown-item>
                    <el-dropdown-item command="unshare">取消共享</el-dropdown-item>
                    <el-dropdown-item divided command="delete">删除</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
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
      </el-tab-pane>
      <el-tab-pane label="私有云" name="private">
        <el-card class="filter-card">
          <el-alert type="info" :closable="false" show-icon>
            <template #title>私有云镜像管理功能开发中</template>
          </el-alert>
        </el-card>
      </el-tab-pane>
    </el-tabs>

    <!-- 上传镜像弹窗 -->
    <el-dialog
      title="上传镜像"
      v-model="uploadDialogVisible"
      width="600px"
    >
      <el-form :model="uploadForm" :rules="uploadRules" ref="uploadFormRef" label-width="120px">
        <el-form-item label="镜像名称" prop="name">
          <el-input v-model="uploadForm.name" placeholder="请输入镜像名称" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="uploadForm.description" type="textarea" :rows="3" placeholder="请输入镜像描述" />
        </el-form-item>
        <el-form-item label="镜像文件" prop="file">
          <el-upload
            ref="uploadRef"
            :auto-upload="false"
            :limit="1"
            accept=".qcow2,.raw,.vhd,.iso"
          >
            <el-button type="primary">选择文件</el-button>
            <template #tip>
              <div class="el-upload__tip">支持 qcow2/raw/vhd/iso 格式，最大 50GB</div>
            </template>
          </el-upload>
        </el-form-item>
        <el-form-item label="操作系统" prop="os_name">
          <el-select v-model="uploadForm.os_name" placeholder="请选择操作系统" style="width: 100%">
            <el-option label="CentOS" value="CentOS" />
            <el-option label="Ubuntu" value="Ubuntu" />
            <el-option label="Windows" value="Windows" />
            <el-option label="Debian" value="Debian" />
            <el-option label="Other" value="Other" />
          </el-select>
        </el-form-item>
        <el-form-item label="版本" prop="os_version">
          <el-input v-model="uploadForm.os_version" placeholder="如 20.04, 7.9" />
        </el-form-item>
        <el-form-item label="架构" prop="cpu_arch">
          <el-select v-model="uploadForm.cpu_arch" placeholder="请选择架构" style="width: 100%">
            <el-option label="x86_64" value="x86_64" />
            <el-option label="arm64" value="arm64" />
          </el-select>
        </el-form-item>
        <el-form-item label="格式" prop="format">
          <el-select v-model="uploadForm.format" placeholder="请选择格式" style="width: 100%">
            <el-option label="qcow2" value="qcow2" />
            <el-option label="raw" value="raw" />
            <el-option label="vhd" value="vhd" />
            <el-option label="iso" value="iso" />
          </el-select>
        </el-form-item>
        <el-form-item label="项目" prop="project_id">
          <el-select v-model="uploadForm.project_id" placeholder="请选择项目" style="width: 100%">
            <el-option v-for="p in projects" :key="p.id" :label="p.name" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="标签" prop="tags">
          <el-select v-model="uploadForm.tags" placeholder="选择或输入标签" multiple filterable allow-create style="width: 100%">
            <el-option label="production" value="production" />
            <el-option label="development" value="development" />
            <el-option label="test" value="test" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="uploadDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmUpload" :loading="uploading">上传</el-button>
      </template>
    </el-dialog>

    <!-- 详情弹窗 -->
    <el-dialog
      title="镜像详情"
      v-model="detailDialogVisible"
      width="800px"
    >
      <el-tabs v-model="detailTab" v-if="selectedImage">
        <el-tab-pane label="基础信息" name="basic">
          <el-descriptions :column="2" border>
            <el-descriptions-item label="ID">{{ selectedImage.id }}</el-descriptions-item>
            <el-descriptions-item label="名称">{{ selectedImage.name }}</el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="getStatusType(selectedImage.status)">{{ getStatusLabel(selectedImage.status) }}</el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="格式">{{ selectedImage.format }}</el-descriptions-item>
            <el-descriptions-item label="操作系统">{{ selectedImage.os_name }}</el-descriptions-item>
            <el-descriptions-item label="版本">{{ selectedImage.os_version || '-' }}</el-descriptions-item>
            <el-descriptions-item label="大小">{{ formatSize(selectedImage.size) }}</el-descriptions-item>
            <el-descriptions-item label="架构">{{ selectedImage.cpu_arch }}</el-descriptions-item>
            <el-descriptions-item label="平台">
              <el-tag size="small" :type="getPlatformType(selectedImage.platform)">
                {{ getPlatformLabel(selectedImage.platform) }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="云账号">{{ selectedImage.account_name || '-' }}</el-descriptions-item>
            <el-descriptions-item label="区域">{{ selectedImage.region_id || '-' }}</el-descriptions-item>
            <el-descriptions-item label="镜像类型">{{ selectedImage.image_type || '系统' }}</el-descriptions-item>
            <el-descriptions-item label="共享范围">{{ selectedImage.share_scope || '私有' }}</el-descriptions-item>
            <el-descriptions-item label="项目">{{ selectedImage.project_name || '-' }}</el-descriptions-item>
            <el-descriptions-item label="创建时间">{{ selectedImage.created_at }}</el-descriptions-item>
            <el-descriptions-item label="描述" :span="2">{{ selectedImage.description || '无' }}</el-descriptions-item>
          </el-descriptions>
        </el-tab-pane>
        <el-tab-pane label="标签" name="tags">
          <el-table :data="selectedImageTags" style="width: 100%">
            <el-table-column prop="key" label="标签键" width="200" />
            <el-table-column prop="value" label="标签值" />
            <el-table-column label="操作" width="100">
              <template #default="{ row }">
                <el-button size="small" link type="danger" @click="removeTag(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
          <el-button type="primary" size="small" style="margin-top: 10px" @click="addTag">添加标签</el-button>
        </el-tab-pane>
        <el-tab-pane label="操作日志" name="logs">
          <el-table :data="operationLogs" style="width: 100%">
            <el-table-column prop="time" label="时间" width="180" />
            <el-table-column prop="action" label="操作" width="120" />
            <el-table-column prop="operator" label="操作人" width="120" />
            <el-table-column prop="result" label="结果">
              <template #default="{ row }">
                <el-tag size="small" :type="row.result === '成功' ? 'success' : 'danger'">{{ row.result }}</el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </el-dialog>

    <!-- 编辑弹窗 -->
    <el-dialog
      title="编辑镜像"
      v-model="editDialogVisible"
      width="500px"
    >
      <el-form :model="editForm" ref="editFormRef" label-width="100px">
        <el-form-item label="名称">
          <el-input v-model="editForm.name" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="editForm.description" type="textarea" :rows="3" />
        </el-form-item>
        <el-form-item label="操作系统">
          <el-select v-model="editForm.os_name" style="width: 100%">
            <el-option label="CentOS" value="CentOS" />
            <el-option label="Ubuntu" value="Ubuntu" />
            <el-option label="Windows" value="Windows" />
            <el-option label="Debian" value="Debian" />
            <el-option label="Other" value="Other" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmEdit">确定</el-button>
      </template>
    </el-dialog>

    <!-- 社区镜像弹窗 -->
    <el-dialog
      title="社区镜像"
      v-model="communityDialogVisible"
      width="800px"
    >
      <el-alert type="info" :closable="false" show-icon class="community-tip">
        <template #title>从社区镜像库选择并导入到您的项目中</template>
      </el-alert>
      <el-table :data="communityImages" style="width: 100%; margin-top: 15px">
        <el-table-column prop="name" label="镜像名称" width="200" />
        <el-table-column prop="os_name" label="操作系统" width="120" />
        <el-table-column prop="os_version" label="版本" width="100" />
        <el-table-column prop="size" label="大小" width="80">
          <template #default="{ row }">{{ formatSize(row.size) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-button size="small" type="primary" @click="importCommunityImage(row)">导入</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { View, Upload, Link, ArrowDown, PriceTag } from '@element-plus/icons-vue'
import { getImages, createImage, deleteImage, updateImage } from '@/api/image'
import { getProjects } from '@/api/iam'

interface SystemImage {
  id: string
  name: string
  description?: string
  status: string
  format: string
  os_name: string
  os_version?: string
  size: number
  cpu_arch: string
  image_type?: string
  share_scope?: string
  project_id?: string
  project_name?: string
  platform?: string
  account_name?: string
  region_id?: string
  created_at?: string
}

// Data
const loading = ref(false)
const uploading = ref(false)
const images = ref<SystemImage[]>([])
const selectedImages = ref<SystemImage[]>([])
const projects = ref<any[]>([])
const activeTab = ref('public')
const detailTab = ref('basic')
const uploadDialogVisible = ref(false)
const detailDialogVisible = ref(false)
const editDialogVisible = ref(false)
const communityDialogVisible = ref(false)
const selectedImage = ref<SystemImage | null>(null)
const selectedImageTags = ref<{ key: string; value: string }[]>([])
const operationLogs = ref<{ time: string; action: string; operator: string; result: string }[]>([])

const uploadFormRef = ref()
const editFormRef = ref()

const filters = reactive({
  name: '',
  os_name: '',
  platform: '',
  format: '',
  status: '',
  cpu_arch: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

const uploadForm = reactive({
  name: '',
  description: '',
  os_name: '',
  os_version: '',
  cpu_arch: 'x86_64',
  format: 'qcow2',
  project_id: '',
  tags: [] as string[]
})

const editForm = reactive({
  id: '',
  name: '',
  description: '',
  os_name: ''
})

const uploadRules = {
  name: [{ required: true, message: '请输入镜像名称', trigger: 'blur' }],
  os_name: [{ required: true, message: '请选择操作系统', trigger: 'change' }],
  format: [{ required: true, message: '请选择格式', trigger: 'change' }]
}

// Mock community images
const communityImages = ref([
  { id: 'c1', name: 'Ubuntu 22.04 LTS', os_name: 'Ubuntu', os_version: '22.04', size: 3221225472 },
  { id: 'c2', name: 'CentOS 8 Stream', os_name: 'CentOS', os_version: '8', size: 2147483648 },
  { id: 'c3', name: 'Debian 12', os_name: 'Debian', os_version: '12', size: 2147483648 },
])

// Methods
const getStatusType = (status: string) => {
  switch (status) {
    case 'Available': return 'success'
    case 'Creating': return 'warning'
    case 'Unavailable': return 'danger'
    default: return 'info'
  }
}

const getStatusLabel = (status: string) => {
  switch (status) {
    case 'Available': return '可用'
    case 'Creating': return '创建中'
    case 'Unavailable': return '不可用'
    default: return status
  }
}

const getPlatformLabel = (platform: string) => {
  const labels: Record<string, string> = {
    aliyun: '阿里云',
    alibaba: '阿里云',
    tencent: '腾讯云',
    aws: 'AWS',
    azure: 'Azure'
  }
  return labels[platform] || platform || '未知'
}

const getPlatformType = (platform: string) => {
  const types: Record<string, 'primary' | 'warning' | 'success' | 'info'> = {
    aliyun: 'primary',
    alibaba: 'primary',
    tencent: 'warning',
    aws: 'success',
    azure: 'info'
  }
  return types[platform] || 'info'
}

const formatSize = (size?: number): string => {
  if (!size) return '-'
  const gb = size / 1024 / 1024 / 1024
  if (gb >= 1) return gb.toFixed(2) + ' GB'
  const mb = size / 1024 / 1024
  return mb.toFixed(2) + ' MB'
}

const handleSelectionChange = (selection: SystemImage[]) => {
  selectedImages.value = selection
}

const handleView = () => {
  ElMessage.info('视图切换功能开发中')
}

const handleTags = () => {
  ElMessage.info('标签管理功能开发中')
}

const resetFilters = () => {
  filters.name = ''
  filters.os_name = ''
  filters.platform = ''
  filters.format = ''
  filters.status = ''
  filters.cpu_arch = ''
  pagination.page = 1
  fetchData()
}

const handleTabChange = (tab: string) => {
  pagination.page = 1
  fetchData()
}

const handleDetails = (row: SystemImage) => {
  selectedImage.value = row
  detailTab.value = 'basic'
  // Mock tags data
  selectedImageTags.value = [
    { key: 'environment', value: 'production' },
    { key: 'os-type', value: row.os_name || 'linux' }
  ]
  // Mock operation logs
  operationLogs.value = [
    { time: row.created_at || '2024-01-01 10:00:00', action: '创建镜像', operator: 'admin', result: '成功' },
    { time: '2024-01-02 14:30:00', action: '修改名称', operator: 'admin', result: '成功' }
  ]
  detailDialogVisible.value = true
}

const addTag = () => {
  ElMessage.info('添加标签功能开发中')
}

const removeTag = (tag: { key: string; value: string }) => {
  const index = selectedImageTags.value.findIndex(t => t.key === tag.key)
  if (index > -1) {
    selectedImageTags.value.splice(index, 1)
    ElMessage.success('标签已删除')
  }
}

const fetchData = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.pageSize,
      details: true,
      is_guest_image: false,
      image_type: activeTab.value === 'public' ? 'public' : 'private',
      ...filters
    }
    const res = await getImages(params)
    images.value = res.items || res.data?.items || []
    pagination.total = res.total || res.data?.pagination?.total || 0
  } catch (e) {
    console.error(e)
    ElMessage.error('获取镜像列表失败')
  } finally {
    loading.value = false
  }
}

const fetchProjects = async () => {
  try {
    const res = await getProjects()
    projects.value = res.items || []
  } catch (e) {
    projects.value = []
  }
}

const handleUpload = () => {
  Object.assign(uploadForm, {
    name: '', description: '', os_name: '', os_version: '',
    cpu_arch: 'x86_64', format: 'qcow2', project_id: '', tags: []
  })
  uploadDialogVisible.value = true
}

const handleCommunityMirror = () => {
  communityDialogVisible.value = true
}

const confirmUpload = async () => {
  if (!uploadFormRef.value) return
  try {
    await uploadFormRef.value.validate()
    uploading.value = true
    await createImage(uploadForm)
    ElMessage.success('镜像上传任务已创建')
    uploadDialogVisible.value = false
    fetchData()
  } catch (e) {
    ElMessage.error('上传失败')
  } finally {
    uploading.value = false
  }
}

const importCommunityImage = async (img: any) => {
  try {
    await ElMessageBox.confirm(`确定要导入 "${img.name}" 吗？`, '导入确认', { type: 'info' })
    await createImage({
      name: img.name,
      os_name: img.os_name,
      os_version: img.os_version,
      format: 'qcow2',
      cpu_arch: 'x86_64'
    })
    ElMessage.success('导入任务已创建')
    communityDialogVisible.value = false
    fetchData()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error('导入失败')
  }
}

const handleBatchCommand = async (command: string) => {
  if (selectedImages.value.length === 0) return

  const actionNames = { share: '共享', unshare: '取消共享', delete: '删除' }
  try {
    await ElMessageBox.confirm(
      `确定要批量${actionNames[command]}选中的 ${selectedImages.value.length} 个镜像吗？`,
      '批量操作确认',
      { type: 'warning' }
    )

    for (const img of selectedImages.value) {
      if (command === 'delete') {
        await deleteImage(img.id)
      } else if (command === 'share') {
        await updateImage(img.id, { share_scope: 'public' })
      } else if (command === 'unshare') {
        await updateImage(img.id, { share_scope: 'private' })
      }
    }
    ElMessage.success(`批量${actionNames[command]}完成`)
    selectedImages.value = []
    fetchData()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error('批量操作失败')
  }
}

const handleActionCommand = (command: string, row: SystemImage) => {
  switch (command) {
    case 'details': handleDetails(row); break
    case 'edit': handleEdit(row); break
    case 'share': handleShare(row); break
    case 'unshare': handleUnshare(row); break
    case 'delete': handleDelete(row); break
  }
}

const handleEdit = (row: SystemImage) => {
  Object.assign(editForm, {
    id: row.id,
    name: row.name,
    description: row.description || '',
    os_name: row.os_name
  })
  editDialogVisible.value = true
}

const confirmEdit = async () => {
  try {
    await updateImage(editForm.id, editForm)
    ElMessage.success('更新成功')
    editDialogVisible.value = false
    fetchData()
  } catch (e) {
    ElMessage.error('更新失败')
  }
}

const handleShare = async (row: SystemImage) => {
  try {
    await ElMessageBox.confirm(`确定要共享镜像 "${row.name}" 吗？`, '共享确认', { type: 'info' })
    await updateImage(row.id, { share_scope: 'public' })
    ElMessage.success('已共享')
    fetchData()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error('共享失败')
  }
}

const handleUnshare = async (row: SystemImage) => {
  try {
    await ElMessageBox.confirm(`确定要取消共享镜像 "${row.name}" 吗？`, '取消共享确认', { type: 'warning' })
    await updateImage(row.id, { share_scope: 'private' })
    ElMessage.success('已取消共享')
    fetchData()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error('取消共享失败')
  }
}

const handleDelete = async (row: SystemImage) => {
  try {
    await ElMessageBox.confirm(`确定要删除镜像 "${row.name}" 吗？`, '删除警告', { type: 'warning' })
    await deleteImage(row.id)
    ElMessage.success('删除成功')
    fetchData()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error('删除失败')
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

onMounted(() => {
  fetchData()
  fetchProjects()
})
</script>

<style scoped>
.system-images-container { padding: 20px; }

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-header h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
}

.toolbar { display: flex; gap: 10px; align-items: center; }
.filter-card { margin-bottom: 20px; }
.image-tabs { margin-bottom: 20px; }

.os-cell {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.os-name { font-weight: 500; }
.os-version { font-size: 12px; color: var(--el-text-color-secondary); }

.platform-cell {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.account-name { font-size: 12px; color: var(--el-text-color-secondary); }

.pagination { margin-top: 20px; text-align: right; justify-content: flex-end; }
.community-tip { margin-bottom: 15px; }
</style>