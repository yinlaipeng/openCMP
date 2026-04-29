<template>
  <div class="login-container">
    <div class="login-box">
      <div class="login-header">
        <el-icon :size="40"><Lock /></el-icon>
        <h1>openCMP</h1>
        <p>多云管理平台</p>
        <p class="domain-info" v-if="domain">{{ domain }} 域</p>
      </div>

      <el-form :model="form" :rules="rules" ref="formRef" class="login-form">
        <el-form-item prop="username">
          <el-input
            v-model="form.username"
            placeholder="请输入用户名"
            size="large"
            prefix-icon="User"
          />
        </el-form-item>
        <el-form-item prop="password">
          <el-input
            v-model="form.password"
            type="password"
            placeholder="请输入密码"
            size="large"
            prefix-icon="Lock"
            show-password
            @keyup.enter="handleLogin"
          />
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            size="large"
            style="width: 100%"
            :loading="loading"
            @click="handleLogin"
          >
            登 录
          </el-button>
        </el-form-item>
      </el-form>

      <div class="login-footer">
        <p>默认管理员账号：admin / admin@123</p>
        <p class="chooser-link" @click="goToChooser">选择其他账号</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import request from '@/utils/request'

const router = useRouter()
const route = useRoute()
const formRef = ref<FormInstance>()
const loading = ref(false)
const domain = ref('')

const form = reactive({
  username: '',
  password: ''
})

const rules: FormRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少 6 位', trigger: 'blur' }
  ]
}

// 处理 URL 参数
onMounted(() => {
  const usernameParam = route.query.username as string
  const domainParam = route.query.fd_domain as string

  if (usernameParam) {
    form.username = usernameParam
  }
  if (domainParam) {
    domain.value = domainParam
  }
})

const handleLogin = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    loading.value = true
    try {
      const res = await request({
        url: '/auth/login',
        method: 'post',
        data: {
          username: form.username,
          password: form.password
        }
      })

      // 先保存 token 和用户信息到 localStorage
      if (res && res.token) {
        localStorage.setItem('token', res.token)
        localStorage.setItem('user', JSON.stringify(res.user))
        console.log('登录成功，token已保存')

        // 然后获取用户权限（异步，不阻塞跳转）
        request({
          url: '/auth/permissions',
          method: 'post'
        }).then(permRes => {
          localStorage.setItem('permissions', JSON.stringify(permRes.permissions || []))
        }).catch(e => {
          console.log('获取权限失败:', e)
        })

        // 获取用户信息（异步，不阻塞跳转）
        request({
          url: '/auth/user',
          method: 'get'
        }).then(userRes => {
          localStorage.setItem('userInfo', JSON.stringify(userRes))
        }).catch(e => {
          console.log('获取用户信息失败:', e)
        })

        ElMessage.success('登录成功')
        router.push('/dashboard')
      } else {
        console.error('登录响应数据不正确:', res)
        ElMessage.error('登录响应数据不正确')
      }
    } catch (e: any) {
      console.error('登录失败:', e)
      const status = e.response?.status
      if (status === 401 || status === 403) {
        ElMessage.error('用户名或密码错误')
      }
    } finally {
      loading.value = false
    }
  })
}

const goToChooser = () => {
  router.push('/auth/login/chooser')
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-box {
  width: 400px;
  padding: 40px;
  background: #fff;
  border-radius: 10px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.2);
}

.login-header {
  text-align: center;
  margin-bottom: 30px;
}

.login-header .el-icon {
  color: #667eea;
  margin-bottom: 10px;
}

.login-header h1 {
  font-size: 28px;
  color: #333;
  margin: 10px 0 5px;
}

.login-header p {
  color: #999;
  font-size: 14px;
}

.login-header .domain-info {
  color: #667eea;
  font-size: 12px;
  margin-top: 5px;
}

.login-form {
  margin-top: 20px;
}

.login-footer {
  margin-top: 20px;
  text-align: center;
  color: #999;
  font-size: 12px;
}

.chooser-link {
  color: #667eea;
  cursor: pointer;
  margin-top: 10px;
}

.chooser-link:hover {
  text-decoration: underline;
}
</style>