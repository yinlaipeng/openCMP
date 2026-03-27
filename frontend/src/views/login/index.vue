<template>
  <div class="login-container">
    <div class="login-box">
      <div class="login-header">
        <el-icon :size="40"><Lock /></el-icon>
        <h1>openCMP</h1>
        <p>多云管理平台</p>
      </div>
      
      <el-form :model="form" :rules="rules" ref="formRef" class="login-form">
        <el-form-item prop="username">
          <el-input
            v-model="form.username"
            placeholder="用户名"
            size="large"
            prefix-icon="User"
          />
        </el-form-item>
        <el-form-item prop="password">
          <el-input
            v-model="form.password"
            type="password"
            placeholder="密码"
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
            登录
          </el-button>
        </el-form-item>
      </el-form>
      
      <div class="login-footer">
        <p>默认管理员账号：admin / admin123</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import request from '@/utils/request'

const router = useRouter()
const formRef = ref<FormInstance>()
const loading = ref(false)

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

const handleLogin = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    
    loading.value = true
    try {
      // 调用登录 API
      const res = await request({
        url: '/auth/login',
        method: 'post',
        data: {
          username: form.username,
          password: form.password
        }
      })
      
      // 保存 token
      localStorage.setItem('token', res.token || 'demo-token')
      localStorage.setItem('user', JSON.stringify(res.user || {
        id: 1,
        name: form.username,
        display_name: '管理员'
      }))
      
      ElMessage.success('登录成功')
      router.push('/')
    } catch (e: any) {
      console.error(e)
      // 演示模式：允许任意 admin 用户登录
      if (form.username === 'admin') {
        localStorage.setItem('token', 'demo-token')
        localStorage.setItem('user', JSON.stringify({
          id: 1,
          name: 'admin',
          display_name: '超级管理员',
          email: 'admin@example.com'
        }))
        ElMessage.success('登录成功（演示模式）')
        router.push('/')
      } else {
        ElMessage.error('用户名或密码错误')
      }
    } finally {
      loading.value = false
    }
  })
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

.login-form {
  margin-top: 20px;
}

.login-footer {
  margin-top: 20px;
  text-align: center;
  color: #999;
  font-size: 12px;
}
</style>
