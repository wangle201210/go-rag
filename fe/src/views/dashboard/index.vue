<template>
  <div class="dashboard-container">
    <!-- 统计卡片 -->
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <template #header>
            <div class="card-header">
              <span>知识库总数</span>
              <el-icon><Folder /></el-icon>
            </div>
          </template>
          <div class="card-body">
            <div class="stat-value">{{ stats.knowledgeBaseCount }}</div>
            <div class="stat-label">个</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <template #header>
            <div class="card-header">
              <span>文档总数</span>
              <el-icon><Document /></el-icon>
            </div>
          </template>
          <div class="card-body">
            <div class="stat-value">{{ stats.documentCount }}</div>
            <div class="stat-label">个</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <template #header>
            <div class="card-header">
              <span>文档块总数</span>
              <el-icon><Files /></el-icon>
            </div>
          </template>
          <div class="card-body">
            <div class="stat-value">{{ stats.chunkCount }}</div>
            <div class="stat-label">个</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <template #header>
            <div class="card-header">
              <span>存储空间</span>
              <el-icon><DataLine /></el-icon>
            </div>
          </template>
          <div class="card-body">
            <div class="stat-value">{{ formatFileSize(stats.storageSize) }}</div>
            <div class="stat-label">已使用</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 最近文档 -->
    <el-card class="recent-documents" shadow="hover">
      <template #header>
        <div class="card-header">
          <span>最近文档</span>
          <el-button type="primary" link @click="goToDocuments">
            查看更多
            <el-icon class="el-icon--right"><ArrowRight /></el-icon>
          </el-button>
        </div>
      </template>
      <el-table :data="recentDocuments" style="width: 100%">
        <el-table-column prop="documentName" label="文档名称" />
        <el-table-column prop="documentType" label="类型" width="120" />
        <el-table-column prop="documentSize" label="大小" width="120">
          <template #default="{ row }">
            {{ formatFileSize(row.documentSize) }}
          </template>
        </el-table-column>
        <el-table-column prop="createTime" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDateTime(row.createTime) }}
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Folder, Document, Files, DataLine, ArrowRight } from '@element-plus/icons-vue'
import { formatFileSize, formatDateTime } from '@/utils/format'

const router = useRouter()

// 统计数据
const stats = ref({
  knowledgeBaseCount: 0,
  documentCount: 0,
  chunkCount: 0,
  storageSize: 0
})

// 最近文档
const recentDocuments = ref([])

// 获取统计数据
const getStats = async () => {
  try {
    // TODO: 调用获取统计数据的接口
    // const res = await getDashboardStats()
    // stats.value = res.data
    
    // 模拟数据
    stats.value = {
      knowledgeBaseCount: 5,
      documentCount: 20,
      chunkCount: 100,
      storageSize: 1024 * 1024 * 10 // 10MB
    }
  } catch (error) {
    console.error('获取统计数据失败:', error)
  }
}

// 获取最近文档
const getRecentDocuments = async () => {
  try {
    // TODO: 调用获取最近文档的接口
    // const res = await getRecentDocuments()
    // recentDocuments.value = res.data
    
    // 模拟数据
    recentDocuments.value = [
      {
        documentName: '示例文档1.md',
        documentType: 'markdown',
        documentSize: 1024 * 10,
        createTime: new Date()
      },
      {
        documentName: '示例文档2.pdf',
        documentType: 'pdf',
        documentSize: 1024 * 100,
        createTime: new Date()
      }
    ]
  } catch (error) {
    console.error('获取最近文档失败:', error)
  }
}

// 跳转到文档列表
const goToDocuments = () => {
  router.push('/document')
}

onMounted(() => {
  getStats()
  getRecentDocuments()
})
</script>

<style scoped>
.dashboard-container {
  padding: 20px;
}

.stat-card {
  height: 180px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-body {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100px;
}

.stat-value {
  font-size: 36px;
  font-weight: bold;
  color: #303133;
}

.stat-label {
  font-size: 14px;
  color: #909399;
  margin-top: 8px;
}

.recent-documents {
  margin-top: 20px;
}
</style> 