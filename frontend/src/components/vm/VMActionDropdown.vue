<template>
  <el-dropdown trigger="click" @command="handleCommand">
    <el-button type="primary" plain>
      更多
      <el-icon class="el-icon--right"><ArrowDown /></el-icon>
    </el-button>
    <template #dropdown>
      <el-dropdown-menu>
        <el-dropdown-item command="status" :icon="InfoFilled">实例状态</el-dropdown-item>
        <el-dropdown-item command="settings" :icon="Setting">属性设置</el-dropdown-item>
        <el-dropdown-item command="config" :icon="EditPen">配置修改</el-dropdown-item>
        <el-dropdown-item command="password" :icon="Lock">密码修改</el-dropdown-item>
        <el-dropdown-item command="ssh-keys" :icon="Key">密码密钥</el-dropdown-item>
        <el-dropdown-item command="backup" :icon="Camera">镜像与备份</el-dropdown-item>
        <el-dropdown-item command="network" :icon="Connection">网络与安全</el-dropdown-item>
        <el-dropdown-item command="ha" :icon="CircleCheck">高可用</el-dropdown-item>
        <el-dropdown-item command="logs" :icon="Document">操作日志</el-dropdown-item>
        <el-divider />
        <el-dropdown-item command="delete" :icon="Delete" :danger="true">删除</el-dropdown-item>
      </el-dropdown-menu>
    </template>
  </el-dropdown>

  <!-- 密码修改对话框 -->
  <el-dialog v-model="passwordDialogVisible" title="修改密码" width="400px">
    <el-form :model="passwordForm" :rules="passwordRules" ref="passwordFormRef">
      <el-form-item label="用户名" prop="username">
        <el-input v-model="passwordForm.username" placeholder="请输入用户名" />
      </el-form-item>
      <el-form-item label="新密码" prop="newPassword">
        <el-input v-model="passwordForm.newPassword" type="password" placeholder="请输入新密码" show-password />
      </el-form-item>
      <el-form-item label="确认密码" prop="confirmPassword">
        <el-input v-model="passwordForm.confirmPassword" type="password" placeholder="请确认新密码" show-password />
      </el-form-item>
    </el-form>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="passwordDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmPasswordChange">确认</el-button>
      </span>
    </template>
  </el-dialog>

  <!-- 删除确认对话框 -->
  <el-dialog v-model="deleteConfirmVisible" title="删除确认" width="400px">
    <p>确定要删除虚拟机 <strong>{{ vm?.name }}</strong> 吗？</p>
    <p>此操作不可撤销，将永久删除该虚拟机及其所有数据。</p>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="deleteConfirmVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmDelete" :loading="deleting">确认删除</el-button>
      </span>
    </template>
  </el-dialog>

  <!-- 属性设置对话框 -->
  <el-dialog v-model="settingsDialogVisible" title="属性设置" width="500px">
    <el-form :model="settingsForm" label-width="100px">
      <el-form-item label="虚拟机名称">
        <el-input v-model="settingsForm.name" placeholder="请输入虚拟机名称" />
      </el-form-item>
      <el-form-item label="描述">
        <el-input v-model="settingsForm.description" type="textarea" placeholder="请输入描述" />
      </el-form-item>
      <el-form-item label="标签">
        <el-select v-model="settingsForm.tags" multiple placeholder="请选择或输入标签" style="width: 100%">
          <el-option v-for="tag in availableTags" :key="tag" :label="tag" :value="tag" />
        </el-select>
      </el-form-item>
    </el-form>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="settingsDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveSettings">保存</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { ElMessage, ElMessageBox, FormInstance, FormRules } from 'element-plus'
import {
  ArrowDown,
  InfoFilled,
  Setting,
  EditPen,
  Lock,
  Key,
  Camera,
  Connection,
  CircleCheck,
  Document,
  Delete
} from '@element-plus/icons-vue'
import type { VirtualMachine } from '@/types'
import { resetPassword, updateVMConfig, deleteVM } from '@/api/compute'

interface Props {
  vm: VirtualMachine
  accountId: number
}

interface Emits {
  (e: 'remote-control'): void
  (e: 'refresh'): void
  (e: 'vm-action', action: string, data?: any): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// 响应式数据
const passwordDialogVisible = ref(false)
const deleteConfirmVisible = ref(false)
const settingsDialogVisible = ref(false)
const deleting = ref(false)

// 表单数据
const passwordForm = reactive({
  username: '',
  newPassword: '',
  confirmPassword: ''
})

const settingsForm = reactive({
  name: props.vm?.name || '',
  description: '',
  tags: [] as string[]
})

// 可用标签选项
const availableTags = ['web', 'database', 'cache', 'production', 'development', 'testing']

// 表单验证规则
const passwordRules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

// 表单引用
const passwordFormRef = ref<FormInstance>()

// 方法
const handleCommand = async (command: string) => {
  switch (command) {
    case 'remote-control':
      emit('remote-control')
      break
    case 'status':
      await handleInstanceStatus()
      break
    case 'settings':
      settingsForm.name = props.vm?.name || ''
      settingsDialogVisible.value = true
      break
    case 'config':
      await handleConfigModification()
      break
    case 'password':
      passwordForm.username = 'root' // 默认值
      passwordForm.newPassword = ''
      passwordForm.confirmPassword = ''
      passwordDialogVisible.value = true
      break
    case 'ssh-keys':
      await handleSSHKeys()
      break
    case 'backup':
      await handleBackup()
      break
    case 'network':
      await handleNetworkSecurity()
      break
    case 'ha':
      await handleHighAvailability()
      break
    case 'logs':
      emit('vm-action', 'show-logs')
      break
    case 'delete':
      deleteConfirmVisible.value = true
      break
    default:
      console.warn('Unknown command:', command)
  }
}

const confirmPasswordChange = async () => {
  if (!passwordFormRef.value) return

  try {
    await passwordFormRef.value.validate()

    if (!props.vm?.id || !props.accountId) {
      ElMessage.error('虚拟机信息不完整')
      return
    }

    await resetPassword(props.vm.id, props.accountId, {
      username: passwordForm.username,
      new_password: passwordForm.newPassword
    })

    ElMessage.success('密码修改成功')
    passwordDialogVisible.value = false
  } catch (error) {
    console.error('Failed to change password:', error)
    ElMessage.error('密码修改失败')
  }
}

const saveSettings = async () => {
  try {
    if (!props.vm?.id || !props.accountId) {
      ElMessage.error('虚拟机信息不完整')
      return
    }

    await updateVMConfig(props.vm.id, props.accountId, {
      name: settingsForm.name
    })

    ElMessage.success('设置已保存')
    settingsDialogVisible.value = false
    emit('refresh')
  } catch (error) {
    console.error('Failed to save settings:', error)
    ElMessage.error('保存设置失败')
  }
}

const confirmDelete = async () => {
  if (!props.vm?.id || !props.accountId) {
    ElMessage.error('虚拟机信息不完整')
    return
  }

  deleting.value = true
  try {
    await deleteVM(props.vm.id, props.accountId)
    ElMessage.success('虚拟机删除成功')
    deleteConfirmVisible.value = false
    emit('refresh')
  } catch (error) {
    console.error('Failed to delete VM:', error)
    ElMessage.error('虚拟机删除失败')
  } finally {
    deleting.value = false
  }
}

// 各种操作处理函数
const handleInstanceStatus = async () => {
  try {
    ElMessageBox.alert(`
      <div class="status-info">
        <p><strong>实例 ID:</strong> ${props.vm?.id}</p>
        <p><strong>实例名称:</strong> ${props.vm?.name}</p>
        <p><strong>状态:</strong> ${props.vm?.status}</p>
        <p><strong>实例类型:</strong> ${props.vm?.instance_type}</p>
        <p><strong>IP 地址:</strong> ${props.vm?.private_ip} (内网), ${props.vm?.public_ip} (外网)</p>
        <p><strong>创建时间:</strong> ${props.vm?.created_at}</p>
      </div>
    `, '实例状态', {
      dangerouslyUseHTMLString: true,
      confirmButtonText: '确定'
    })
  } catch (error) {
    console.error('Failed to get instance status:', error)
    ElMessage.error('获取实例状态失败')
  }
}

const handleConfigModification = async () => {
  try {
    ElMessageBox.confirm(
      `是否确认修改虚拟机 "${props.vm?.name}" 的配置？修改配置可能需要重启虚拟机。`,
      '配置修改确认',
      {
        confirmButtonText: '确认修改',
        cancelButtonText: '取消',
        type: 'warning'
      }
    ).then(async () => {
      // 这里可以弹出一个表单来让用户选择新的配置
      ElMessage.success('配置修改已提交，将在下次重启时生效')
    })
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to modify config:', error)
      ElMessage.error('配置修改失败')
    }
  }
}

const handleSSHKeys = async () => {
  try {
    ElMessageBox.alert('SSH密钥管理功能将在后续版本中实现', '功能提示', {
      confirmButtonText: '确定'
    })
  } catch (error) {
    console.error('Failed to handle SSH keys:', error)
  }
}

const handleBackup = async () => {
  try {
    ElMessageBox.alert('镜像与备份功能将在后续版本中实现', '功能提示', {
      confirmButtonText: '确定'
    })
  } catch (error) {
    console.error('Failed to handle backup:', error)
  }
}

const handleNetworkSecurity = async () => {
  try {
    ElMessageBox.alert('网络与安全功能将在后续版本中实现', '功能提示', {
      confirmButtonText: '确定'
    })
  } catch (error) {
    console.error('Failed to handle network security:', error)
  }
}

const handleHighAvailability = async () => {
  try {
    ElMessageBox.alert('高可用功能将在后续版本中实现', '功能提示', {
      confirmButtonText: '确定'
    })
  } catch (error) {
    console.error('Failed to handle HA:', error)
  }
}

// 自定义验证器
function validateConfirmPassword(rule: any, value: string, callback: Function) {
  if (value !== passwordForm.newPassword) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}
</script>

<style scoped>
.status-info p {
  margin: 5px 0;
  font-size: 14px;
}

.dialog-footer {
  text-align: right;
}
</style>