<template>
  <div class="error-container">
    <div class="error-content">
      <h1>404</h1>
      <h2>抱歉，您访问的页面不存在</h2>
      <el-button type="primary" @click="goHome">返回首页</el-button>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'

const router = useRouter()

const goHome = () => {
  router.push('/')
}
</script>

<style scoped>
.error-container {
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f0f2f5;
}

.error-content {
  text-align: center;
}

.error-content h1 {
  font-size: 120px;
  color: #409EFF;
  margin: 0;
  line-height: 1.2;
}

.error-content h2 {
  font-size: 24px;
  color: #606266;
  margin: 20px 0 30px;
}
</style> 