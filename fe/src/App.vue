<template>
  <div class="app-container">
    <el-container>
      <el-aside width="200px">
        <div class="logo-container">
          <h2>Go-RAG</h2>
        </div>
        <el-menu
          router
          :default-active="$route.path"
          class="el-menu-vertical"
          background-color="#001529"
          text-color="#fff"
          active-text-color="#409EFF">
          <el-menu-item index="/knowledge-base">
            <el-icon><Folder /></el-icon>
            <span>知识库管理</span>
          </el-menu-item>
          <el-menu-item index="/indexer">
            <el-icon><Upload /></el-icon>
            <span>文档索引</span>
          </el-menu-item>
          <el-menu-item index="/retriever">
            <el-icon><Search /></el-icon>
            <span>文档检索</span>
          </el-menu-item>
          <el-menu-item index="/chat">
            <el-icon><ChatDotRound /></el-icon>
            <span>智能问答</span>
          </el-menu-item>
        </el-menu>
      </el-aside>
      <el-container>
        <el-header height="60px">
          <div class="header-content">
            <div class="header-title">
              <h2>基于检索增强生成的智能问答系统</h2>
            </div>
          </div>
        </el-header>
        <el-main>
          <router-view />
        </el-main>
        <el-footer height="40px">
          <div class="footer-content">
            © 2025 Go-RAG Project
          </div>
        </el-footer>
      </el-container>
    </el-container>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { Folder } from '@element-plus/icons-vue'
</script>

<style scoped>
.app-container {
  height: 100vh;
  width: 100%;
}

.el-container {
  height: 100%;
}

.el-aside {
  background-color: #001529;
  color: white;
  height: 100%;
}

.logo-container {
  height: 60px;
  display: flex;
  justify-content: center;
  align-items: center;
  color: #409EFF;
}

.el-header {
  background-color: #fff;
  color: #333;
  line-height: 60px;
  border-bottom: 1px solid #eee;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-title {
  text-align: center;
}

.header-actions {
  display: flex;
  align-items: center;
}

.el-main {
  padding: 20px;
  background-color: #f5f7fa;
  height: calc(100vh - 100px);
  overflow-y: auto;
}

.el-footer {
  background-color: #fff;
  color: #999;
  text-align: center;
  line-height: 40px;
  font-size: 12px;
  border-top: 1px solid #eee;
}

.footer-content {
  text-align: center;
}
</style>