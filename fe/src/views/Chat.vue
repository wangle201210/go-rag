<template>
  <div class="chat-container">
    <el-row :gutter="20">
      <el-col :span="16">
        <el-card class="chat-card">
          <template #header>
            <div class="card-header">
              <el-icon class="header-icon"><ChatDotRound /></el-icon>
              <span>智能问答</span>
              <el-button 
                type="primary" 
                size="small" 
                plain 
                class="new-session-btn"
                @click="startNewSession">
                <el-icon><Plus /></el-icon> 新会话
              </el-button>
            </div>
          </template>
          
          <div class="chat-messages" ref="messagesContainer">
            <div v-if="messages.length === 0" class="empty-chat">
              <el-empty description="开始一个新的对话吧">
                <template #image>
                  <el-icon class="empty-icon"><ChatRound /></el-icon>
                </template>
              </el-empty>
            </div>
            
            <div v-else class="message-list">
              <div 
                v-for="(message, index) in messages" 
                :key="index"
                :class="['message-item', message.role === 'user' ? 'user-message' : 'ai-message']">
                <div class="message-avatar">
                  <el-avatar :icon="message.role === 'user' ? User : Service" :size="36" />
                </div>
                <div class="message-content">
                  <div class="message-text" v-if="message.role === 'user'">{{ message.content }}</div>
                  <div class="message-text markdown-content" v-else v-html="renderMarkdown(message.content)"></div>
                  <div class="message-time">{{ formatTime(message.timestamp) }}</div>
                </div>
              </div>
            </div>
            
            <div v-if="loading" class="loading-message">
              <el-skeleton :rows="3" animated />
            </div>
          </div>
          
          <div class="chat-input">
            <el-form @submit.prevent="sendMessage">
              <el-input
                v-model="inputMessage"
                type="textarea"
                :rows="3"
                placeholder="请输入您的问题..."
                :disabled="loading"
                @keyup.enter.exact="sendMessage"
                @keyup.ctrl.enter="sendMessage">
              </el-input>
              <div class="input-actions">
                <el-tooltip content="高级设置" placement="top">
                  <el-button 
                    type="info" 
                    plain 
                    circle 
                    @click="showSettings = !showSettings">
                    <el-icon><Setting /></el-icon>
                  </el-button>
                </el-tooltip>
                <el-button 
                  type="primary" 
                  :loading="loading"
                  :disabled="!inputMessage.trim()"
                  @click="sendMessage">
                  发送 <el-icon class="el-icon--right"><Position /></el-icon>
                </el-button>
              </div>
            </el-form>
            
            <el-collapse-transition>
              <div v-show="showSettings" class="settings-panel">
                <el-form :model="chatSettings" label-position="left" label-width="180px">
                  <el-form-item label="参考文档返回结果数量">
                    <el-input-number
                      v-model="chatSettings.top_k"
                      :min="1"
                      :max="10"
                      controls-position="right"
                      size="small"
                    />
                  </el-form-item>
                  <el-form-item label="相似度阈值">
                    <el-slider
                      v-model="chatSettings.score"
                      :min="0"
                      :max="1"
                      :step="0.05"
                      :format-tooltip="(val) => val.toFixed(2)"
                      size="small"
                    />
                  </el-form-item>
                </el-form>
              </div>
            </el-collapse-transition>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="8">
        <el-card class="references-card">
          <template #header>
            <div class="card-header">
              <el-icon class="header-icon"><Document /></el-icon>
              <span>会话信息</span>
            </div>
          </template>
          <div class="session-info">
            <div class="session-id">
              <span class="label">会话ID:</span>
              <el-tag size="small" type="info">{{ sessionId }}</el-tag>
              <el-tooltip content="复制会话ID" placement="top">
                <el-button
                    type="primary"
                    link
                    size="small"
                    @click="copySessionId">
                  <el-icon><CopyDocument /></el-icon>
                </el-button>
              </el-tooltip>
            </div>
            <div class="message-count">
              <span class="label">消息数:</span>
              <span>{{ messages.length }}</span>
            </div>
          </div>


          
          <div class="references-content">
            <el-divider content-position="left">参考文档</el-divider>

            <div v-if="references.length === 0" class="empty-references">
              <el-empty description="暂无参考文档" />
            </div>
            
            <div v-else class="reference-list">
              <el-collapse accordion>
                <el-collapse-item 
                  v-for="(ref, index) in references" 
                  :key="index"
                  :title="`文档片段 #${index + 1} (相似度: ${ref.meta_data._score.toFixed(2)})`"
                  :name="index">
                  <div class="reference-content">
                    <div class="source-info">
                      <el-tag size="small">{{ ref.meta_data.ext._file_name || '未知来源' }}</el-tag>
                    </div>
                    <div class="content-text markdown-content" v-html="renderMarkdown(ref.content)"></div>
                  </div>
                </el-collapse-item>
              </el-collapse>
            </div>
          </div>
          
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, nextTick, watch } from 'vue'
import { 
  ChatDotRound, 
  ChatRound, 
  User, 
  Service, 
  Document, 
  Setting, 
  Position, 
  Plus, 
  CopyDocument 
} from '@element-plus/icons-vue'
import axios from 'axios'
import { ElMessage } from 'element-plus'
import { v4 as uuidv4 } from 'uuid'
import { marked } from 'marked'
import DOMPurify from 'dompurify'
import hljs from 'highlight.js'
import 'highlight.js/styles/github.css'
import { getKnowledgeName } from '../utils/knowledgeStore'

const messages = ref([])
const inputMessage = ref('')
const loading = ref(false)
const references = ref([])
const showSettings = ref(false)
const messagesContainer = ref(null)
const sessionId = ref(uuidv4())

const chatSettings = reactive({
  top_k: 5,
  score: 0.2
})

// 格式化时间
const formatTime = (timestamp) => {
  const date = new Date(timestamp)
  return date.toLocaleTimeString()
}

// 滚动到底部
const scrollToBottom = async () => {
  await nextTick()
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}

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

// 发送消息
const sendMessage = async () => {
  const message = inputMessage.value.trim()
  if (!message || loading.value) return
  
  // 添加用户消息
  messages.value.push({
    role: 'user',
    content: message,
    timestamp: Date.now()
  })
  
  inputMessage.value = ''
  loading.value = true
  scrollToBottom()
  
  try {
    const response = await axios.post('/v1/chat', {
      question: message,
      top_k: chatSettings.top_k,
      score: chatSettings.score,
      conv_id: sessionId.value,
      knowledge_name: getKnowledgeName()
    })
    
    // 添加AI回复
    messages.value.push({
      role: 'assistant',
      content: response.data.data.answer || '抱歉，我无法回答这个问题。',
      timestamp: Date.now()
    })
    
    // 更新参考文档
    references.value = response.data.data.references || []
  } catch (error) {
    console.error('发送消息失败:', error)
    ElMessage.error('发送消息失败: ' + (error.response?.data?.message || '未知错误'))
    
    // 添加错误消息
    messages.value.push({
      role: 'assistant',
      content: '抱歉，发生了错误，请稍后重试。',
      timestamp: Date.now()
    })
  } finally {
    loading.value = false
    await scrollToBottom()
  }
}

// 开始新会话
const startNewSession = () => {
  sessionId.value = uuidv4()
  messages.value = []
  references.value = []
  ElMessage.success('已开始新会话')
}

// 复制会话ID
const copySessionId = () => {
  navigator.clipboard.writeText(sessionId.value)
    .then(() => ElMessage.success('会话ID已复制'))
    .catch(() => ElMessage.error('复制失败'))
}

// 监听消息变化，自动滚动到底部
watch(messages, () => {
  scrollToBottom()
}, { deep: true })

onMounted(() => {
  scrollToBottom()
})
</script>

<style scoped>
.chat-container {
  height: calc(100vh - 140px);
  max-height: 800px;
  min-height: 500px;
}

.chat-card, .references-card {
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
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

.new-session-btn {
  margin-left: auto;
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 10px;
  background-color: #f9f9f9;
  border-radius: 4px;
  margin-bottom: 15px;
  min-height: 300px;
  max-height: calc(100vh - 350px);
}

.empty-chat {
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}

.empty-icon {
  font-size: 60px;
  color: #909399;
}

.message-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.message-item {
  display: flex;
  margin-bottom: 15px;
}

.user-message {
  flex-direction: row-reverse;
}

.message-avatar {
  margin: 0 10px;
}

.message-content {
  max-width: 70%;
  padding: 10px 15px;
  border-radius: 8px;
  position: relative;
}

.user-message .message-content {
  background-color: #ecf5ff;
  border: 1px solid #d9ecff;
  text-align: right;
}

.ai-message .message-content {
  background-color: #fff;
  border: 1px solid #ebeef5;
  text-align: left;
}

.message-text {
  word-break: break-word;
  line-height: 1.5;
}

.message-time {
  font-size: 12px;
  color: #909399;
  margin-top: 5px;
}

.loading-message {
  padding: 10px;
  background-color: #fff;
  border-radius: 8px;
  margin: 10px 0;
  border: 1px solid #ebeef5;
}

.chat-input {
  margin-top: auto;
}

.input-actions {
  display: flex;
  justify-content: space-between;
  margin-top: 10px;
}

.settings-panel {
  margin-top: 15px;
  padding: 15px;
  background-color: #f9f9f9;
  border-radius: 4px;
  border: 1px solid #ebeef5;
  position: absolute;
  width: calc(100% - 30px);
  z-index: 10;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.references-content {
  flex: 1;
  overflow-y: auto;
  height: calc(100% - 100px);
}

.empty-references {
  height: calc(100% - 100px);
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 20px 0;
}

.reference-content {
  padding: 10px;
  background-color: #f9f9f9;
  border-radius: 4px;
}

.source-info {
  margin-bottom: 10px;
}

.content-text {
  white-space: pre-wrap;
  line-height: 1.6;
}

.session-info {
  margin-top: auto;
  padding: 15px 0;
}

.session-id, .message-count {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
}

.label {
  font-weight: bold;
  margin-right: 10px;
  color: #606266;
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