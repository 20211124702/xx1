<template>
  <div class="register-bg">
    <el-card class="register-card">
      <div class="register-title">
        <el-icon size="32"><UserFilled /></el-icon>
        <span>用户注册</span>
      </div>
      <el-form :model="form" :rules="rules" ref="formRef">
        <el-form-item prop="username">
          <el-input v-model="form.username" placeholder="用户名" />
        </el-form-item>
        <el-form-item prop="password">
          <el-input v-model="form.password" type="password" placeholder="密码" />
        </el-form-item>
        <el-form-item prop="email">
          <el-input v-model="form.email" placeholder="邮箱" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" style="width:100%" :loading="loading" @click="handleRegister">注册</el-button>
        </el-form-item>
        <el-form-item>
          <el-link @click="$router.push('/')">已有账号？去登录</el-link>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { UserFilled } from '@element-plus/icons-vue'
import { register } from '../api/auth'

const router = useRouter()
const formRef = ref()
const loading = ref(false)
const form = reactive({
  username: '',
  password: '',
  email: ''
})

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
  email: [{ required: true, message: '请输入邮箱', trigger: 'blur' }]
}

const handleRegister = async () => {
  if (!formRef.value) return
  try {
    await formRef.value.validate()
    loading.value = true
    await register(form)
    ElMessage.success('注册成功，请登录')
    router.push('/login')
  } catch (error: any) {
    ElMessage.error(error.message || '注册失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.register-bg {
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
.register-card {
  width: 350px;
  padding: 30px 20px 20px 20px;
  border-radius: 16px;
  box-shadow: 0 8px 32px 0 rgba(31, 38, 135, 0.2);
}
.register-title {
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