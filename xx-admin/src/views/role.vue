<template>
  <div>
    <h2>权限管理</h2>
    <el-table :data="roles">
      <el-table-column prop="id" label="ID" />
      <el-table-column prop="name" label="角色名" />
    </el-table>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
const roles = ref([
  { id: 1, name: '管理员' },
  { id: 2, name: '普通用户' }
])
</script> 