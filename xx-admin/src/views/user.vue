<template>
  <div>
    <h2>用户管理</h2>
    <el-table :data="users">
      <el-table-column prop="id" label="ID" />
      <el-table-column prop="name" label="用户名" />
      <el-table-column prop="role" label="角色" />
    </el-table>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
const users = ref([
  { id: 1, name: '张三', role: '管理员' },
  { id: 2, name: '李四', role: '普通用户' }
])
</script> 