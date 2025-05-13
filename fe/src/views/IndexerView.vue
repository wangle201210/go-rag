<template>
  <div class="indexer-page">
    <el-card class="indexer-card">
      <template #header>
        <div class="card-header">
          <h2>文档索引</h2>
          <p class="description">上传文件并向量化到ES，支持多种文件格式</p>
        </div>
      </template>
      
      <div class="upload-container">
        <el-upload
          class="upload-area"
          drag
          action="#"
          :auto-upload="false"
          :on-change="handleFileChange"
          :limit="1"
          :file-list="fileList"
        >
          <el-icon class="el-icon--upload"><i-ep-upload-filled /></el-icon>
          <div class="el-upload__text">
            将文件拖到此处，或<em>点击上传</em>
          </div>
        </el-upload>
        
        <el-button 
          type="primary" 
          :disabled="!selectedFile || isUploading" 
          @click="uploadFile"
          :loading="isUploading"
        >
          {{ isUploading ? '上传中...' : '开始索引' }}
        </el-button>
      </div>
      
      <el-alert
        v-if="uploadStatus"
        :title="uploadStatus.message"
        :type="uploadStatus.success ? 'success' : 'error'"
        show-icon
        :closable="false"
        style="margin-top: 20px;"
      />
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import axios from 'axios'
import { UploadFilled } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const selectedFile = ref(null)
const isUploading = ref(false)
const uploadStatus = ref(null)
const fileList = ref([])

const handleFileChange = (file) => {
  // 只保留最新上传的文件
  fileList.value = [file]
  selectedFile.value = file.raw
  uploadStatus.value = null
}

const uploadFile = async () => {
  if (!selectedFile.value) return
  
  isUploading.value = true
  uploadStatus.value = null
  
  const formData = new FormData()
  formData.append('file', selectedFile.value)
  
  try {
    const response = await axios.post('http://localhost:8000/v1/indexer', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
    console.log('上传成功:', response.data.data.doc_ids)
    uploadStatus.value = {
      success: true,
      message: '文件索引成功！'
    }
    ElMessage.success('文件索引成功！')
  } catch (error) {
    console.error('上传失败:', error)
    uploadStatus.value = {
      success: false,
      message: `索引失败: ${error.response?.data?.message || error.message || '未知错误'}`
    }
    ElMessage.error(`索引失败: ${error.response?.data?.message || error.message || '未知错误'}`)
  } finally {
    isUploading.value = false
  }
}
</script>

<style scoped>
.indexer-page {
  padding: 20px;
}

.indexer-card {
  margin: 0 auto;
}

.card-header {
  display: flex;
  flex-direction: column;
  align-items: center;
}


.description {
  color: var(--el-text-color-secondary);
  font-size: 0.9rem;
  margin-bottom: 0;
}

.upload-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1.5rem;
  padding: 20px 0;
}

.upload-area {
  width: 100%;
  max-width: 500px;
}

:deep(.el-upload-dragger) {
  width: 100%;
  height: 200px;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

:deep(.el-icon--upload) {
  font-size: 48px;
  color: var(--el-color-primary);
  margin-bottom: 16px;
}

:deep(.el-upload__text) {
  color: var(--el-text-color-regular);
}

:deep(.el-upload__text em) {
  color: var(--el-color-primary);
  font-style: normal;
  cursor: pointer;
}
</style>