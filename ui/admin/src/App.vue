<script setup lang="ts">
// OA管理系统根组件
</script>

<template>
  <div id="app">
    <router-view />
  </div>
</template>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html, body, #app {
  height: 100%;
  font-family: 'PingFang SC', 'Microsoft YaHei', 'Source Han Sans CN', 'Helvetica Neue', Helvetica, Arial, sans-serif;
}

#app {
  background-color: #f0f2f5;
}

/* 全局滚动条样式 */
::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 3px;
}

::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}

/* Element Plus 全局样式调整 */
.el-button {
  font-weight: 500;
}

.el-card {
  border: 1px solid #ebeef5;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
}

.el-table {
  border-radius: 8px;
}

.el-pagination {
  justify-content: center;
}

/* 表单样式 */
.el-form-item__label {
  font-weight: 500;
  color: #333;
}

/* 加载状态 */
.el-loading-mask {
  background-color: rgba(255, 255, 255, 0.8);
}
</style>
