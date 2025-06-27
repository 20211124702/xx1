<template>
  <el-container class="dashboard-layout">
    <!-- 侧边菜单 -->
    <el-aside width="200px" class="dashboard-aside no-padding">
      <el-menu
        :default-active="activeMenu"
        class="el-menu-vertical-demo no-padding"
        @select="handleMenuSelect"
        background-color="#2d3a4b"
        text-color="#fff"
        active-text-color="#409EFF"
      >
        <el-menu-item index="/home">
          <el-icon><House /></el-icon>
          <span>首页</span>
        </el-menu-item>
        <el-menu-item index="/home/user">
          <el-icon><User /></el-icon>
          <span>用户管理</span>
        </el-menu-item>
        <el-menu-item index="/home/role">
          <el-icon><Setting /></el-icon>
          <span>角色管理</span>
        </el-menu-item>
        <el-menu-item index="/home/menu">
          <el-icon><Menu /></el-icon>
          <span>菜单管理</span>
        </el-menu-item>
        <el-menu-item index="/home/table">
          <el-icon><List /></el-icon>
          <span>数据表格</span>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-container>
      <!-- 顶部栏 -->
      <el-header class="dashboard-header-bar no-padding">
        <div class="header-title">通用后台管理系统</div>
        <div class="header-user">
          <el-dropdown>
            <span class="el-dropdown-link">
              <el-avatar :size="32" :src="user.avatar || defaultAvatar" />
              <span class="role-label">{{ user.role?.name || '普通用户' }}</span>
              <el-icon><ArrowDown /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item disabled>
                  角色：{{ user.role?.name || '普通用户' }}
                </el-dropdown-item>
                <el-dropdown-item>
                  <span @click="showProfile">个人信息</span>
                </el-dropdown-item>
                <el-dropdown-item divided @click="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <!-- 主内容区 -->
      <el-main class="dashboard-main-content no-padding">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { House, User, Setting, Menu, List, ArrowDown } from '@element-plus/icons-vue'
import defaultAvatar from '../assets/vue.svg'

const router = useRouter()
const route = useRoute()
const user = ref(JSON.parse(localStorage.getItem('user') || '{}'))
const activeMenu = computed(() => route.path)

const handleMenuSelect = (index: string) => {
  
  router.push(index)
}

const showProfile = () => {
  ElMessage.info('这里可以弹出个人信息弹窗或跳转到个人中心页面')
}

const logout = () => {
  localStorage.removeItem('token')
  localStorage.removeItem('user')
  ElMessage.success('已退出登录')
  router.push('/login')
}
</script>

<style scoped>
.dashboard-layout {
  height: 100vh;
  width: 100vw;
  min-width: 0;
  margin: 0;
  padding: 0;
}
.dashboard-aside {
  background: #2d3a4b;
  color: #fff;
  min-height: 100vh;
  box-shadow: 2px 0 8px 0 rgba(0,0,0,0.04);
  padding: 0 !important;
  margin: 0 !important;
  border: none;
}
.el-menu-vertical-demo {
  border-right: none;
  padding: 0 !important;
  margin: 0 !important;
}
.no-padding {
  padding: 0 !important;
  margin: 0 !important;
}
.dashboard-header-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 60px;
  background: #fff;
  box-shadow: 0 2px 8px 0 rgba(0,0,0,0.03);
  padding: 0 32px;
  margin: 0;
}
.header-title {
  position: absolute;
  left: 0;
  right: 0;
  margin: auto;
  text-align: center;
  font-size: 20px;
  font-weight: bold;
  color: #2d3a4b;
  pointer-events: none;
}
.header-user {
  position: absolute;
  right: 32px;
  top: 0;
  height: 60px;
  display: flex;
  align-items: center;
}
.role-label {
  margin: 0 8px;
  color: #409EFF;
  font-weight: 500;
}
.dashboard-main-content {
  background: linear-gradient(120deg, #f5f7fa 0%, #c9e7ff 100%);
  min-height: calc(100vh - 60px);
  padding: 0 !important;
  margin: 0 !important;
  overflow: auto;
}
</style> 