<template>
  <div>
    <h2>菜单管理</h2>
    <el-table :data="menus">
      <el-table-column prop="id" label="ID" />
      <el-table-column prop="name" label="菜单名" />
      <el-table-column prop="path" label="路径" />
    </el-table>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
const menus = ref([
  { id: 1, name: '首页', path: '/' },
  { id: 2, name: '用户管理', path: '/user' }
])
</script> 