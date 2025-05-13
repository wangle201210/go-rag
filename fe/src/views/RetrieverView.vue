<template>
  <div class="retriever-page">
    <el-card class="retriever-card">
      <template #header>
        <div class="card-header">
          <h2>文档检索</h2>
          <p class="description">根据问题检索相关文档</p>
        </div>
      </template>
      
      <div class="search-container">
        <el-input
          v-model="question"
          placeholder="请输入您的问题"
          clearable
          @keyup.enter="searchDocuments"
        >
          <template #append>
            <el-button 
              type="primary" 
              :disabled="!question || isSearching" 
              @click="searchDocuments"
              :loading="isSearching"
            >
              {{ isSearching ? '检索中' : '检索文档' }}
            </el-button>
          </template>
        </el-input>
        
        <div class="search-params">
          <el-form :inline="true">
            <el-form-item label="Top K:">
              <el-input-number
                v-model="topK"
                :min="1"
                :max="20"
                :controls="true"
                size="small"
              />
            </el-form-item>
            <el-form-item label="最低相关度:">
              <el-input-number
                v-model="score"
                :min="0"
                :max="1"
                :step="0.1"
                :precision="1"
                size="small"
              />
            </el-form-item>
          </el-form>
        </div>
      </div>
      
      <div v-if="searchResults.length > 0" class="results-container">
        <h3>检索结果</h3>
        <el-divider />
        <div class="results-list">
          <el-card 
            v-for="(result, index) in searchResults" 
            :key="index" 
            class="result-item"
            shadow="hover"
          >
            <template #header>
              <div class="result-header">
                <el-tag type="primary" size="small" effect="plain">{{ index + 1 }}</el-tag>
                <el-tag type="info" size="small">相关度: {{ result.meta_data?._score.toFixed(2) || '未知' }}</el-tag>
              </div>
            </template>
            <div class="result-content">
              <div v-html="md.render(result.content)"></div>
            </div>
          </el-card>
        </div>
      </div>
      
      <el-empty 
        v-else-if="hasSearched && !isSearching" 
        description="没有找到相关文档，请尝试其他问题或调整参数"
      />
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import axios from 'axios'
import { ElMessage } from 'element-plus'
import MarkdownIt from 'markdown-it'

const question = ref('')
const topK = ref(5)
const score = ref(0.2)
const isSearching = ref(false)
const searchResults = ref([])
const hasSearched = ref(false)

// 创建markdown解析器
const md = new MarkdownIt()

const searchDocuments = async () => {
  if (!question.value || isSearching.value) return
  
  isSearching.value = true
  searchResults.value = []
  
  try {
    const response = await axios.post('http://localhost:8000/v1/retriever', {
      question: question.value,
      top_k: topK.value,
      score: score.value
    })
    
    searchResults.value = response.data.data?.document || []
    hasSearched.value = true
    
    if (searchResults.value.length === 0) {
      ElMessage.info('没有找到相关文档，请尝试其他问题或调整参数')
    } else {
      ElMessage.success(`成功检索到 ${searchResults.value.length} 条相关文档`)
    }
  } catch (error) {
    console.error('检索失败:', error)
    ElMessage.error(`检索失败: ${error.response?.data?.message || error.message || '未知错误'}`)
  } finally {
    isSearching.value = false
  }
}
</script>

<style scoped>
.retriever-page {
  padding: 20px;
}

.retriever-card {
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

.search-container {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.search-params {
  margin-top: 1rem;
}

.results-container {
  margin-top: 2rem;
}

.results-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.result-item {
  margin-bottom: 1rem;
}

.result-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.result-content {
  margin-top: 0.5rem;
  line-height: 1.6;
  color: var(--el-text-color-primary);
}

.result-content :deep(p) {
  margin: 0 0 0.5em 0;
}

.result-content :deep(pre) {
  background-color: #f6f8fa;
  border-radius: 3px;
  padding: 12px;
  overflow: auto;
}

.result-content :deep(code) {
  font-family: monospace;
  background-color: rgba(0, 0, 0, 0.05);
  padding: 2px 4px;
  border-radius: 3px;
}

.result-content :deep(blockquote) {
  border-left: 4px solid #dfe2e5;
  padding-left: 16px;
  margin-left: 0;
  color: #6a737d;
}

.result-content :deep(ul), .result-content :deep(ol) {
  padding-left: 2em;
}

.result-content :deep(table) {
  border-collapse: collapse;
  margin: 12px 0;
}

.result-content :deep(th), .result-content :deep(td) {
  border: 1px solid #dfe2e5;
  padding: 6px 13px;
}

.result-content :deep(th) {
  background-color: #f6f8fa;
}
</style>