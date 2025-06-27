<template>
  <div class="login-bg">
    <el-card class="login-card">
      <div class="login-title">
        <el-icon size="32"><UserFilled /></el-icon>
        <span>后台管理系统登录</span>
      </div>
      <el-form :model="loginForm" :rules="rules" ref="loginFormRef">
        <el-form-item prop="username">
          <el-input v-model="loginForm.username" placeholder="用户名">
            <template #prefix>
              <el-icon><User /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item prop="password">
          <el-input v-model="loginForm.password" type="password" placeholder="密码" @keyup.enter="handleLogin">
            <template #prefix>
              <el-icon><Lock /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" style="width:100%" :loading="loading" @click="handleLogin">登 录</el-button>
        </el-form-item>
        <el-form-item>
          <el-link @click="$router.push('/register')">没有账号？去注册</el-link>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>
<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock, UserFilled } from '@element-plus/icons-vue'
import { login } from '../api/auth'

const router = useRouter()
const loginFormRef = ref()
const loading = ref(false)

const loginForm = reactive({
  username: '',
  password: ''
})

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' }
  ]
}

const handleLogin = async () => {
  if (!loginFormRef.value) return
  
  try {
    await loginFormRef.value.validate()
    loading.value = true
    
    const response = await login(loginForm)
    localStorage.setItem('token', response.data.token)
    localStorage.setItem('user', JSON.stringify(response.data.user))
    
    ElMessage.success('登录成功')
    router.push('/home')
  } catch (error: any) {
    ElMessage.error(error.message || '登录失败')
  } finally {
    loading.value = false
  }
}
</script>
<style scoped>
.login-bg {
  width: 100vw;
  height: 100vh;
  background: linear-gradient(135deg, #74ebd5 0%, #ACB6E5 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0;
  padding: 0;
  position: fixed;
  top: 0;
  left: 0;
}
.login-card {
  width: 350px;
  padding: 30px 20px 20px 20px;
  border-radius: 16px;
  box-shadow: 0 8px 32px 0 rgba(31, 38, 135, 0.2);
}
.login-title {
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 22px;
  font-weight: bold;
  margin-bottom: 24px;
  color: #409eff;
  gap: 10px;
}
</style>
