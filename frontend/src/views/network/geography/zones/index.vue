<template>
  <div class="zones-container">
    <div class="page-header">
      <h2>可用区</h2>
      <div class="toolbar">
        <el-button link type="primary" @click="handleView">
          <el-icon><View /></el-icon>
          查看
        </el-button>
        <el-button type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>
          创建
        </el-button>
        <el-dropdown trigger="click" @command="handleBatchCommand" :disabled="selectedZones.length === 0">
          <el-button :disabled="selectedZones.length === 0">
            批量操作<el-icon class="el-icon--right"><ArrowDown /></el-icon>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="delete">批量删除</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>

    <!-- 搜索表单 -->
    <el-card class="filter-card">
      <el-form :inline="true" :model="filters" @submit.prevent="fetchData">
        <el-form-item label="名称">
          <el-input v-model="filters.name" placeholder="搜索可用区名称" clearable style="width: 180px" />
        </el-form-item>
        <el-form-item label="区域">
          <el-input v-model="filters.region_id" placeholder="区域ID" clearable style="width: 120px" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchData">查询</el-button>
          <el-button @click="resetFilters">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="table-card">
      <el-table
        :data="zones"
        v-loading="loading"
        style="width: 100%"
        row-key="id"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <!-- Cloudpods 表头顺序: Name → L2 Networks → Hosts/Enabled Hosts → Physical Machines -->
        <el-table-column prop="name" label="名称" width="200">
          <template #default="{ row }">
            <el-link type="primary" :underline="false" @click="handleDetails(row)">
              {{ row.name }}
            </el-link>
          </template>
        </el-table-column>
        <el-table-column prop="l2_network_count" label="二层网络" width="120" />
        <el-table-column label="宿主机/可用宿主机" width="180">
          <template #default="{ row }">
            {{ row.host_count }}/{{ row.enabled_host_count }}
          </template>
        </el-table-column>
        <el-table-column label="物理机/可用物理机" width="180">
          <template #default="{ row }">
            {{ row.physical_host_count }}/{{ row.enabled_physical_host_count }}
          </template>
        </el-table-column>
        <el-table-column label="操作" fixed="right" width="100">
          <template #default="{ row }">
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

    <!-- 创建可用区弹窗 -->
    <el-dialog
      title="创建可用区"
      v-model="createDialogVisible"
      width="400px"
    >
      <el-form :model="createForm" :rules="createRules" ref="createFormRef" label-width="100px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="createForm.name" placeholder="请输入可用区名称" />
        </el-form-item>
        <el-form-item label="区域" prop="region_id">
          <el-select v-model="createForm.region_id" placeholder="选择区域" style="width: 100%">
            <el-option v-for="r in regions" :key="r.id" :label="r.name" :value="r.id" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmCreate" :loading="creating">创建</el-button>
      </template>
    </el-dialog>

    <!-- 详情弹窗 -->
    <el-dialog
      title="可用区详情"
      v-model="detailDialogVisible"
      width="700px"
    >
      <el-tabs v-model="activeTab">
        <el-tab-pane label="详情" name="detail">
          <el-descriptions :column="2" border v-if="selectedZone">
            <el-descriptions-item label="ID">{{ selectedZone.id }}</el-descriptions-item>
            <el-descriptions-item label="名称">{{ selectedZone.name }}</el-descriptions-item>
            <el-descriptions-item label="区域">{{ selectedZone.region_id }}</el-descriptions-item>
            <el-descriptions-item label="二层网络数">{{ selectedZone.l2_network_count }}</el-descriptions-item>
            <el-descriptions-item label="宿主机">{{ selectedZone.host_count }}/{{ selectedZone.enabled_host_count }}</el-descriptions-item>
            <el-descriptions-item label="物理机">{{ selectedZone.physical_host_count }}/{{ selectedZone.enabled_physical_host_count }}</el-descriptions-item>
          </el-descriptions>
        </el-tab-pane>
        <el-tab-pane label="宿主机" name="hosts">
          <el-table :data="hosts" v-loading="hostsLoading">
            <el-table-column prop="id" label="宿主机ID" width="180" />
            <el-table-column prop="name" label="宿主机名称" width="150" />
            <el-table-column prop="status" label="状态" width="100" />
            <el-table-column prop="ip" label="IP地址" width="150" />
          </el-table>
        </el-tab-pane>
      </el-tabs>
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { View, Plus, ArrowDown } from '@element-plus/icons-vue'
import { getZones, createZone, deleteZone, batchDeleteZones, Zone } from '@/api/networkSync'

// Data
const loading = ref(false)
const creating = ref(false)
const zones = ref<Zone[]>([])
const selectedZones = ref<Zone[]>([])
const regions = ref<any[]>([])

const createDialogVisible = ref(false)
const detailDialogVisible = ref(false)
const selectedZone = ref<Zone | null>(null)
const activeTab = ref('detail')

const hosts = ref<any[]>([])
const hostsLoading = ref(false)

const filters = reactive({
  name: '',
  region_id: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const createForm = reactive({
  name: '',
  region_id: ''
})

const createRules = {
  name: [{ required: true, message: '请输入可用区名称', trigger: 'blur' }],
  region_id: [{ required: true, message: '请选择区域', trigger: 'change' }]
}

// Methods
const handleSelectionChange = (selection: Zone[]) => {
  selectedZones.value = selection
}

const resetFilters = () => {
  filters.name = ''
  filters.region_id = ''
  pagination.page = 1
  fetchData()
}

const fetchData = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize,
      ...filters
    }
    const res = await getZones(params)
    zones.value = res.items || []
    pagination.total = res.total || 0
  } catch (error) {
    console.error('Failed to fetch zones:', error)
    ElMessage.error('获取可用区列表失败')
  } finally {
    loading.value = false
  }
}

const handleView = () => {
  ElMessage.info('视图切换功能开发中')
}

const handleCreate = () => {
  Object.assign(createForm, {
    name: '',
    region_id: ''
  })
  createDialogVisible.value = true
}

const confirmCreate = async () => {
  creating.value = true
  try {
    await createZone({
      name: createForm.name,
      region_id: createForm.region_id
    })
    ElMessage.success('可用区创建成功')
    createDialogVisible.value = false
    fetchData()
  } catch (error) {
    ElMessage.error('创建失败')
  } finally {
    creating.value = false
  }
}

const handleBatchCommand = async (command: string) => {
  if (selectedZones.value.length === 0) return
  try {
    await ElMessageBox.confirm(
      `确定要批量删除选中的 ${selectedZones.value.length} 个可用区吗？`,
      '批量操作确认',
      { type: 'warning' }
    )
    if (command === 'delete') {
      const ids = selectedZones.value.map(z => z.id)
      await batchDeleteZones(ids)
      ElMessage.success('批量删除成功')
    }
    selectedZones.value = []
    fetchData()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('批量操作失败')
  }
}

const handleDetails = (row: Zone) => {
  selectedZone.value = row
  detailDialogVisible.value = true
  activeTab.value = 'detail'
  loadZoneDetails(row.id)
}

const loadZoneDetails = async (zoneId: string) => {
  hostsLoading.value = true
  hosts.value = []
  hostsLoading.value = false
}

const handleDelete = async (row: Zone) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除可用区 "${row.name}" 吗？`,
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    await deleteZone(row.id)
    ElMessage.success('删除成功')
    fetchData()
  } catch (error) {
    console.error('Failed to delete zone:', error)
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
.zones-container {
  padding: 20px;
}

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

.toolbar {
  display: flex;
  gap: 10px;
  align-items: center;
}

.filter-card {
  margin-bottom: 20px;
}

.table-card {
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  justify-content: flex-end;
}
</style>