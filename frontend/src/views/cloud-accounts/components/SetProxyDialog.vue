<template>
  <el-dialog
    v-model="visible"
    title="设置代理"
    width="500px"
    :close-on-click-modal="false"
  >
    <el-form :model="form" label-width="120px" v-loading="loading">
      <el-form-item label="启用代理">
        <el-switch v-model="form.enableProxy" />
      </el-form-item>

      <el-form-item label="代理地址" v-if="form.enableProxy">
        <el-input v-model="form.proxyAddress" placeholder="http://proxy.example.com:8080" />
      </el-form-item>

      <el-form-item label="代理类型" v-if="form.enableProxy">
        <el-select v-model="form.proxyType">
          <el-option label="HTTP" value="http" />
          <el-option label="HTTPS" value="https" />
          <el-option label="SOCKS5" value="socks5" />
        </el-select>
      </el-form-item>

      <el-form-item label="需要认证" v-if="form.enableProxy">
        <el-checkbox v-model="form.needAuth" />
      </el-form-item>

      <el-form-item label="用户名" v-if="form.enableProxy && form.needAuth">
        <el-input v-model="form.proxyUsername" />
      </el-form-item>

      <el-form-item label="密码" v-if="form.enableProxy && form.needAuth">
        <el-input v-model="form.proxyPassword" type="password" show-password />
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button type="primary" @click="handleSave" :loading="saving">保存</el-button>
      <el-button @click="visible = false">取消</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { getCloudAccount, updateCloudAccount } from '@/api/cloud-account'

interface Props {
  modelValue: boolean
  accountId: number | null
}

interface Emits {
  (e: 'update:modelValue', value: boolean): void
  (e: 'saved'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const loading = ref(false)
const saving = ref(false)

const form = reactive({
  enableProxy: false,
  proxyAddress: '',
  proxyType: 'http',
  needAuth: false,
  proxyUsername: '',
  proxyPassword: ''
})

const loadAccount = async () => {
  if (!props.accountId) return
  loading.value = true
  try {
    const account = await getCloudAccount(props.accountId)
    // 从credentials中解析代理配置
    if (account.credentials) {
      const creds = JSON.parse(JSON.stringify(account.credentials))
      form.enableProxy = creds.proxy_enabled || false
      form.proxyAddress = creds.proxy_address || ''
      form.proxyType = creds.proxy_type || 'http'
      form.needAuth = creds.proxy_need_auth || false
      form.proxyUsername = creds.proxy_username || ''
    }
  } catch (error) {
    ElMessage.error('加载云账号失败')
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  if (!props.accountId) return
  saving.value = true
  try {
    // 获取当前账号信息
    const account = await getCloudAccount(props.accountId)
    const credentials = JSON.parse(JSON.stringify(account.credentials || {}))

    // 更新代理配置
    credentials.proxy_enabled = form.enableProxy
    credentials.proxy_address = form.proxyAddress
    credentials.proxy_type = form.proxyType
    credentials.proxy_need_auth = form.needAuth
    credentials.proxy_username = form.proxyUsername
    if (form.needAuth && form.proxyPassword) {
      credentials.proxy_password = form.proxyPassword
    }

    await updateCloudAccount(props.accountId, { credentials })

    ElMessage.success('代理设置成功')
    emit('saved')
    visible.value = false
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

watch(() => props.accountId, (id) => {
  if (id && visible.value) {
    loadAccount()
  }
}, { immediate: true })
</script>