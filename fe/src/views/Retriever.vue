<template>
  <div class="retriever-container">
    <el-card class="retriever-card">
      <template #header>
        <div class="card-header">
          <el-icon class="header-icon"><Search /></el-icon>
          <span>文档检索</span>
          <div class="header-actions">
            <KnowledgeSelector ref="knowledgeSelectorRef" />
          </div>
        </div>
      </template>
      
      <div class="search-area">
        <el-form :model="searchForm" label-position="top">
          <el-form-item label="搜索问题">
            <el-input
              v-model="searchForm.question"
              placeholder="请输入您想要检索的问题"
              clearable
              @keyup.enter="handleSearch">
              <template #append>
                <el-button :icon="Search" @click="handleSearch">检索</el-button>
              </template>
            </el-input>
          </el-form-item>
          
          <el-form-item>
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="返回结果数量">
                  <el-input-number
                    v-model="searchForm.top_k"
                    :min="1"
                    :max="10"
                    controls-position="right"
                  />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="相似度阈值">
                  <el-slider
                    v-model="searchForm.score"
                    :min="0"
                    :max="1"
                    :step="0.05"
                    :format-tooltip="(val) => val.toFixed(2)"
                  />
                </el-form-item>
              </el-col>
            </el-row>
          </el-form-item>
        </el-form>
      </div>
      
      <div class="loading-area" v-if="loading">
        <el-skeleton :rows="5" animated />
      </div>
      
      <div class="result-area" v-if="!loading && searchResults.length > 0">
        <div class="result-header">
          <el-divider content-position="left">
            <el-icon><Document /></el-icon>
            检索结果
          </el-divider>
        </div>
        
        <el-collapse v-model="activeNames">
          <el-collapse-item 
            v-for="(result, index) in searchResults" 
            :key="index"
            :title="`文档片段 #${index + 1} (相似度: ${result.meta_data._score.toFixed(2)})`"
            :name="index">
            <div class="result-content">
              <el-card shadow="never" class="content-card">
                <div class="source-info">
                  <el-tag size="small">{{ result.meta_data.ext._file_name || '未知来源' }}</el-tag>
                </div>
                <div class="content-text markdown-content" v-html="renderMarkdown(result.content)"></div>
              </el-card>
            </div>
          </el-collapse-item>
        </el-collapse>
      </div>
      
      <div class="empty-result" v-if="!loading && searchResults.length === 0 && searched">
        <el-empty description="未找到相关文档" />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { Search, Document } from '@element-plus/icons-vue'
import axios from 'axios'
import { ElMessage } from 'element-plus'
import { marked } from 'marked'
import DOMPurify from 'dompurify'
import hljs from 'highlight.js'
import 'highlight.js/styles/github.css'
import KnowledgeSelector from '../components/KnowledgeSelector.vue'
import request from "../utils/request";

// 配置Marked和代码高亮
marked.setOptions({
  highlight: function(code, lang) {
    if (lang && hljs.getLanguage(lang)) {
      return hljs.highlight(code, { language: lang }).value;
    }
    return hljs.highlightAuto(code).value;
  },
  breaks: true
});

// Markdown渲染函数
const renderMarkdown = (content) => {
  if (!content) return '';
  try {
    const html = marked(content);
    return DOMPurify.sanitize(html);
  } catch (error) {
    console.error('Markdown渲染错误:', error);
    return content;
  }
};

const searchForm = reactive({
  question: '',
  top_k: 5,
  score: 0.2
})

const loading = ref(false)
const searchResults = ref([])
const activeNames = ref([0]) // 默认展开第一个结果
const searched = ref(false)
const knowledgeSelectorRef = ref(null)

const handleSearch = async () => {
  if (!searchForm.question) {
    ElMessage.warning('请输入搜索问题')
    return
  }
  
  loading.value = true
  searched.value = true
  
  try {
    const response = await request.post('/v1/retriever', {
      question: searchForm.question,
      top_k: searchForm.top_k,
      score: searchForm.score,
      knowledge_name: knowledgeSelectorRef.value?.getSelectedKnowledgeId() || ''
    })
    searchResults.value = response.data.document || []
    
    if (searchResults.value.length === 0) {
      ElMessage.info('未找到相关文档')
    }
  } catch (error) {
    console.error('检索失败:', error)
    ElMessage.error('检索失败: ' + (error.response?.data?.message || '未知错误'))
    searchResults.value = []
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.retriever-container {
  max-width: 800px;
  margin: 0 auto;
}

.retriever-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  align-items: center;
  font-size: 16px;
  font-weight: bold;
}

.header-actions {
  margin-left: auto;
}

.header-icon {
  margin-right: 8px;
  font-size: 18px;
}

.search-area {
  margin-bottom: 20px;
}

.result-header {
  margin: 20px 0;
}

.result-content {
  padding: 10px 0;
}

.content-card {
  background-color: #f9f9f9;
}

.source-info {
  margin-bottom: 10px;
}

.content-text {
  white-space: pre-wrap;
  line-height: 1.6;
}

.empty-result {
  padding: 40px 0;
}

/* Markdown 样式 */
.markdown-content {
  text-align: left;
}

.markdown-content :deep(h1),
.markdown-content :deep(h2),
.markdown-content :deep(h3),
.markdown-content :deep(h4),
.markdown-content :deep(h5),
.markdown-content :deep(h6) {
  margin-top: 16px;
  margin-bottom: 8px;
  font-weight: 600;
  line-height: 1.25;
}

.markdown-content :deep(h1) {
  font-size: 1.5em;
}

.markdown-content :deep(h2) {
  font-size: 1.3em;
}

.markdown-content :deep(h3) {
  font-size: 1.2em;
}

.markdown-content :deep(p) {
  margin-top: 0;
  margin-bottom: 10px;
}

.markdown-content :deep(ul),
.markdown-content :deep(ol) {
  padding-left: 20px;
  margin-bottom: 10px;
}

.markdown-content :deep(pre) {
  padding: 12px;
  overflow: auto;
  font-size: 85%;
  line-height: 1.45;
  background-color: #f6f8fa;
  border-radius: 3px;
  margin-bottom: 10px;
}

.markdown-content :deep(code) {
  padding: 0.2em 0.4em;
  margin: 0;
  font-size: 85%;
  background-color: rgba(27, 31, 35, 0.05);
  border-radius: 3px;
}

.markdown-content :deep(pre code) {
  padding: 0;
  background-color: transparent;
}

.markdown-content :deep(blockquote) {
  padding: 0 1em;
  color: #6a737d;
  border-left: 0.25em solid #dfe2e5;
  margin-bottom: 10px;
}

.markdown-content :deep(table) {
  border-collapse: collapse;
  width: 100%;
  margin-bottom: 10px;
}

.markdown-content :deep(th),
.markdown-content :deep(td) {
  padding: 6px 13px;
  border: 1px solid #dfe2e5;
}

.markdown-content :deep(tr:nth-child(2n)) {
  background-color: #f6f8fa;
}
</style>