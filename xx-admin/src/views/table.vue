<template>
  <div>
    <h2>数据表格</h2>
    <el-table :data="tableData">
      <el-table-column prop="id" label="ID" />
      <el-table-column prop="name" label="名称" />
      <el-table-column prop="value" label="值" />
    </el-table>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
const tableData = ref([
  { id: 1, name: 'A', value: 100 },
  { id: 2, name: 'B', value: 200 }
])
</script> 