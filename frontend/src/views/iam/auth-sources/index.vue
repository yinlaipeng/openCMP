<template>
  <div class="auth-sources-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span class="title">认证源管理</span>
          <el-button type="primary" @click="handleCreate">
            <el-icon><Plus /></el-icon>
            添加认证源
          </el-button>
        </div>
      </template>
      
      <el-table :data="sources" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="名称" />
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="getTypeTag(row.type)">
              {{ getTypeName(row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="enabled" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.enabled ? 'success' : 'info'">
              {{ row.enabled ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="auto_create" label="自动创建" width="80">
          <template #default="{ row }">
            {{ row.auto_create ? '是' : '否' }}
          </template>
        </el-table-column>
        <el-table-column prop="description" label="描述" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleTest(row)">测试</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    
    <el-dialog v-model="dialogVisible" title="添加认证源" width="600px">
      <el-form :model="form" label-width="120px">
        <el-form-item label="名称" required>
          <el-input v-model="form.name" placeholder="请输入名称" />
        </el-form-item>
        <el-form-item label="类型" required>
          <el-select v-model="form.type" placeholder="请选择类型" style="width: 100%">
            <el-option label="LDAP" value="ldap" />
            <el-option label="OIDC" value="oidc" />
            <el-option label="SAML" value="saml" />
          </el-select>
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" />
        </el-form-item>
        <el-form-item label="自动创建用户">
          <el-switch v-model="form.auto_create" />
        </el-form-item>
        <el-form-item label="启用">
          <el-switch v-model="form.enabled" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getAuthSources, createAuthSource, testAuthSource } from '@/api/iam'

const sources = ref([])
const loading = ref(false)
const dialogVisible = ref(false)

const form = reactive({
  name: '',
  type: 'ldap',
  description: '',
  auto_create: false,
  enabled: true
})

const getTypeName = (type: string) => {
  const map: Record<string, string> = {
    ldap: 'LDAP',
    oidc: 'OIDC',
    saml: 'SAML',
    local: '本地'
  }
  return map[type] || type
}

const getTypeTag = (type: string) => {
  const map: Record<string, any> = {
    ldap: 'warning',
    oidc: 'success',
    saml: 'primary',
    local: 'info'
  }
  return map[type] || ''
}

const loadSources = async () => {
  loading.value = true
  try {
    const res = await getAuthSources()
    sources.value = res.items || res
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const handleCreate = () => {
  form.name = ''
  form.type = 'ldap'
  form.description = ''
  form.auto_create = false
  form.enabled = true
  dialogVisible.value = true
}

const handleTest = async (row: any) => {
  try {
    await testAuthSource(row.id)
    ElMessage.success('连接测试成功')
  } catch (e) {
    ElMessage.error('连接测试失败')
  }
}

const handleSubmit = async () => {
  try {
    await createAuthSource(form)
    ElMessage.success('创建成功')
    dialogVisible.value = false
    loadSources()
  } catch (e) {
    console.error(e)
  }
}

const handleDelete = async (row: any) => {
  ElMessage.info('删除功能开发中')
}

onMounted(() => {
  loadSources()
})
</script>

<style scoped>
.auth-sources-page {
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
