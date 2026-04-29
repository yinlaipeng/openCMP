<template>
  <div class="login-chooser-container">
    <div class="login-chooser-box">
      <div class="login-chooser-header">
        <el-icon :size="40"><UserFilled /></el-icon>
        <h1>openCMP</h1>
        <p>多云管理平台 - 选择账号</p>
      </div>

      <div class="account-list">
        <div class="account-item" @click="selectAccount('Default', 'admin')">
          <div class="account-icon">
            <el-icon><OfficeBuilding /></el-icon>
          </div>
          <div class="account-info">
            <div class="account-name">Default Domain</div>
            <div class="account-desc">默认域 - 系统管理员</div>
          </div>
        </div>

        <div class="account-divider">
          <span>或者</span>
        </div>

        <div class="manual-login">
          <el-button type="primary" @click="goToLogin">手动输入账号</el-button>
        </div>
      </div>

      <div class="login-chooser-footer">
        <p>选择账号后将跳转到登录页面</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter, useRoute } from 'vue-router'
import { UserFilled, OfficeBuilding } from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()

const selectAccount = (domain: string, username: string) => {
  router.push({
    path: '/login',
    query: {
      fd_domain: domain,
      username: username
    }
  })
}

const goToLogin = () => {
  router.push('/login')
}
</script>

<style scoped>
.login-chooser-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-chooser-box {
  width: 400px;
  padding: 40px;
  background: #fff;
  border-radius: 10px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.2);
}

.login-chooser-header {
  text-align: center;
  margin-bottom: 30px;
}

.login-chooser-header .el-icon {
  color: #667eea;
  margin-bottom: 10px;
}

.login-chooser-header h1 {
  font-size: 28px;
  color: #333;
  margin: 10px 0 5px;
}

.login-chooser-header p {
  color: #999;
  font-size: 14px;
}

.account-list {
  margin-top: 20px;
}

.account-item {
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 15px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
  background: #f5f7fa;
  margin-bottom: 10px;
}

.account-item:hover {
  background: #ecf5ff;
  border: 1px solid #409eff;
}

.account-icon {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.account-info {
  flex: 1;
}

.account-name {
  font-size: 16px;
  font-weight: bold;
  color: #303133;
}

.account-desc {
  font-size: 12px;
  color: #909399;
}

.account-divider {
  text-align: center;
  margin: 20px 0;
  position: relative;
}

.account-divider::before,
.account-divider::after {
  content: '';
  position: absolute;
  top: 50%;
  width: 40%;
  height: 1px;
  background: #e4e7ed;
}

.account-divider::before {
  left: 0;
}

.account-divider::after {
  right: 0;
}

.account-divider span {
  padding: 0 10px;
  color: #909399;
  font-size: 12px;
}

.manual-login {
  text-align: center;
}

.login-chooser-footer {
  margin-top: 20px;
  text-align: center;
  color: #999;
  font-size: 12px;
}
</style>