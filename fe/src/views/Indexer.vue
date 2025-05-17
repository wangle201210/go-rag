<template>
  <div class="indexer-container">
    <el-card class="indexer-card">
      <template #header>
        <div class="card-header">
          <el-icon class="header-icon"><Upload /></el-icon>
          <span>文档索引</span>
        </div>
      </template>
      <div class="upload-area">
        <el-upload
          class="upload-component"
          drag
          action="/v1/indexer"
          :on-success="handleUploadSuccess"
          :on-error="handleUploadError"
          :before-upload="beforeUpload"
          :data="{ knowledge_name: getKnowledgeName() }"
          :show-file-list="true"
          multiple>
          <el-icon class="el-icon--upload"><Upload /></el-icon>
          <div class="el-upload__text">
            拖拽文件到此处或 <em>点击上传</em>
          </div>
          <template #tip>
            <div class="el-upload__tip">
              支持上传 PDF、Markdown、HTML 等文档文件
            </div>
          </template>
        </el-upload>
      </div>
      
      <div class="process-info" v-if="processingInfo">
        <el-alert
          :title="processingInfo.title"
          :type="processingInfo.type"
          :description="processingInfo.description"
          show-icon
          :closable="false">
        </el-alert>
      </div>
    </el-card>
    
    <el-card class="indexer-info-card" v-if="indexResult">
      <template #header>
        <div class="card-header">
          <el-icon class="header-icon"><InfoFilled /></el-icon>
          <span>索引结果</span>
        </div>
      </template>
      <el-descriptions :column="1" border>
        <el-descriptions-item label="文档片段数">{{ indexResult.chunks }}</el-descriptions-item>
        <el-descriptions-item label="索引状态">
          <el-tag :type="indexResult.status === 'success' ? 'success' : 'danger'">
            {{ indexResult.status === 'success' ? '成功' : '失败' }}
          </el-tag>
        </el-descriptions-item>
      </el-descriptions>
    </el-card>
  </div>
</template>

<script setup>
import { getKnowledgeName } from '../utils/knowledgeStore'
import { ref } from 'vue'
import { ElMessage } from 'element-plus'

const processingInfo = ref(null)
const indexResult = ref(null)

const beforeUpload = (file) => {
  // 检查文件类型
  const allowedTypes = ['application/pdf', 'text/markdown', 'text/html', 'text/plain']
  const isAllowed = allowedTypes.includes(file.type)
  
  if (!isAllowed) {
    ElMessage.error('只支持 PDF、Markdown、HTML 和文本文件!')
    return false
  }
  
  // 显示处理中信息
  processingInfo.value = {
    title: '文档处理中',
    type: 'info',
    description: `正在处理文件: ${file.name}，请稍候...`
  }
  
  return true
}

const handleUploadSuccess = (response) => {
  processingInfo.value = {
    title: '文档处理完成',
    type: 'success',
    description: '文档已成功索引到系统中'
  }
  // 显示索引结果
  indexResult.value = {
    chunks: response.data.doc_ids.length,
    status: 'success'
  }
  
  ElMessage.success('文档索引成功!')
}

const handleUploadError = (error) => {
  processingInfo.value = {
    title: '文档处理失败',
    type: 'error',
    description: '文档索引过程中发生错误，请重试'
  }
  
  indexResult.value = {
    status: 'error'
  }
  
  ElMessage.error('文档索引失败: ' + (error.message || '未知错误'))
}
</script>

<style scoped>
.indexer-container {
  max-width: 800px;
  margin: 0 auto;
}

.indexer-card, .indexer-info-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  align-items: center;
  font-size: 16px;
  font-weight: bold;
}

.header-icon {
  margin-right: 8px;
  font-size: 18px;
}

.upload-area {
  display: flex;
  justify-content: center;
  padding: 20px 0;
}

.upload-component {
  width: 100%;
}

.process-info {
  margin-top: 20px;
}
</style>