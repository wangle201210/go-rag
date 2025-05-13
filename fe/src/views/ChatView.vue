<template>
  <div class="chat-page">
    <el-card class="chat-card">
      <template #header>
        <div class="card-header">
          <h2>智能问答</h2>
          <p class="description">检索文档后回答您的问题</p>
        </div>
      </template>
      
      <div class="chat-interface">
        <div class="messages-container" ref="messagesContainer">
          <el-empty v-if="messages.length === 0" description="开始提问，AI将基于检索到的文档为您解答" />
          
          <div v-for="(message, index) in messages" :key="index" class="message-wrapper">
            <el-card 
              class="message-card" 
              :class="message.role"
              :shadow="'hover'"
            >
              <div class="message-content">
                <div v-html="md.render(message.content)"></div>
              </div>
              
              <div v-if="message.documents && message.documents.length > 0" class="message-documents">
                <el-collapse>
                  <el-collapse-item>
                    <template #title>
                      <el-badge :value="message.documents.length" type="primary">
                        查看参考文档
                      </el-badge>
                    </template>
                    <div class="documents-list">
                      <el-card 
                        v-for="(doc, docIndex) in message.documents" 
                        :key="docIndex" 
                        class="document-item"
                        shadow="never"
                      >
                        <div v-html="md.render(doc.content)"></div>
                        <el-tag v-if="doc.metadata" size="small" type="info">来源: {{ doc.metadata.source || '未知' }}</el-tag>
                      </el-card>
                    </div>
                  </el-collapse-item>
                </el-collapse>
              </div>
            </el-card>
          </div>
        </div>
        
        <div class="input-container">
          <el-form :inline="true" class="params-form">
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
            <el-form-item label="会话ID:">
              <div class="conv-id-container">
                <el-input v-model="convId" size="small" />
                <el-button 
                  @click="generateNewConvId" 
                  type="primary" 
                  :icon="RefreshRight" 
                  circle 
                  size="small"
                  title="生成新会话ID"
                />
              </div>
            </el-form-item>
          </el-form>
          
          <div class="message-input">
            <el-input 
              v-model="userInput" 
              type="textarea" 
              :rows="3"
              placeholder="请输入您的问题" 
              @keyup.enter.ctrl="sendMessage"
              :disabled="isProcessing"
              resize="none"
            />
            <el-button 
              type="primary" 
              :disabled="!userInput || isProcessing" 
              @click="sendMessage"
              :loading="isProcessing"
            >
              {{ isProcessing ? '处理中...' : '发送' }}
            </el-button>
          </div>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, watch, nextTick } from 'vue'
import axios from 'axios'
import { RefreshRight } from '@element-plus/icons-vue'
import MarkdownIt from 'markdown-it'

const messagesContainer = ref(null)
const userInput = ref('')
const messages = ref([])
const isProcessing = ref(false)
const topK = ref(5)
const score = ref(0.2)
const convId = ref(generateUUID())

// 创建markdown解析器
const md = new MarkdownIt()

// 生成UUID
function generateUUID() {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
    const r = Math.random() * 16 | 0
    const v = c === 'x' ? r : (r & 0x3 | 0x8)
    return v.toString(16)
  })
}

// 生成新的会话ID
function generateNewConvId() {
  convId.value = generateUUID()
  messages.value = []
}

// 发送消息
async function sendMessage() {
  if (!userInput.value || isProcessing.value) return
  
  const question = userInput.value.trim()
  userInput.value = ''
  
  // 添加用户消息
  messages.value.push({
    role: 'user',
    content: question
  })
  
  isProcessing.value = true
  
  try {
    const response = await axios.post('http://localhost:8000/v1/chat', {
      question,
      top_k: topK.value,
      score: score.value,
      conv_id: convId.value
    })
    
    // 添加AI回复
    messages.value.push({
      role: 'assistant',
      content: response.data?.data?.answer || '抱歉，我无法回答这个问题。',
      documents: response.data.documents || []
    })
  } catch (error) {
    console.error('请求失败:', error)
    messages.value.push({
      role: 'system',
      content: `错误: ${error.response?.data?.message || error.message || '未知错误'}`
    })
  } finally {
    isProcessing.value = false
  }
}

// 监听消息变化，自动滚动到底部
watch(messages, () => {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  })
}, { deep: true })

// 组件挂载时，生成新的会话ID
onMounted(() => {
  convId.value = generateUUID()
})
</script>

<style scoped>
.chat-page {
  padding: 20px;
}

.chat-card {
  margin: 0 auto;
  max-width: 1000px;
}

.card-header {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.description {
  margin: 8px 0 0;
  color: var(--el-text-color-secondary);
  font-size: 0.9rem;
}

.chat-interface {
  display: flex;
  flex-direction: column;
  height: calc(100vh - 250px);
  max-height: 700px;
  overflow: hidden;
}

.messages-container {
  flex: 1;
  overflow-y: auto;
  padding: 10px;
  margin-bottom: 20px;
}

.message-wrapper {
  margin-bottom: 16px;
}

.message-card {
  max-width: 80%;
  margin-bottom: 10px;
  line-height: 1.5;
}

.message-card.user {
  margin-left: auto;
  background-color: var(--el-color-primary-light-9);
}

.message-card.assistant {
  margin-right: auto;
  background-color: white;
}

.message-card.system {
  margin: 0 auto;
  background-color: var(--el-color-warning-light-9);
  color: var(--el-color-warning-dark-2);
  max-width: 100%;
  text-align: left;
}

.message-content :deep(p) {
  margin: 0 0 0.5em 0;
  white-space: pre-wrap;
}

.message-content :deep(pre) {
  background-color: #f6f8fa;
  border-radius: 3px;
  padding: 12px;
  overflow: auto;
}

.message-content :deep(code) {
  font-family: monospace;
  background-color: rgba(0, 0, 0, 0.05);
  padding: 2px 4px;
  border-radius: 3px;
}

.message-content :deep(blockquote) {
  border-left: 4px solid #dfe2e5;
  padding-left: 16px;
  margin-left: 0;
  color: #6a737d;
}

.message-content :deep(ul), .message-content :deep(ol) {
  padding-left: 2em;
}

.message-content :deep(table) {
  border-collapse: collapse;
  margin: 12px 0;
}

.message-content :deep(th), .message-content :deep(td) {
  border: 1px solid #dfe2e5;
  padding: 6px 13px;
}

.message-content :deep(th) {
  background-color: #f6f8fa;
}

.message-documents {
  margin-top: 12px;
}

.documents-list {
  margin-top: 10px;
  max-height: 200px;
  overflow-y: auto;
}

.document-item {
  margin-bottom: 10px;
}

.document-item :deep(p) {
  margin: 0 0 8px 0;
}

.document-item :deep(pre) {
  background-color: #f6f8fa;
  border-radius: 3px;
  padding: 12px;
  overflow: auto;
}

.document-item :deep(code) {
  font-family: monospace;
  background-color: rgba(0, 0, 0, 0.05);
  padding: 2px 4px;
  border-radius: 3px;
}

.input-container {
  padding: 16px;
  border-top: 1px solid var(--el-border-color-lighter);
}

.params-form {
  margin-bottom: 16px;
}

.conv-id-container {
  display: flex;
  align-items: center;
  gap: 8px;
}

.message-input {
  display: flex;
  gap: 12px;
  align-items: flex-start;
}

.message-input .el-textarea {
  flex: 1;
}

.message-input .el-button {
  margin-top: 8px;
}
</style>