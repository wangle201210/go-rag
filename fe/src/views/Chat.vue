<template>
  <div class="chat-container">
    <el-row :gutter="20">
      <el-col :span="16">
        <el-card class="chat-card">
          <template #header>
            <div class="card-header">
              <el-icon class="header-icon"><ChatDotRound /></el-icon>
              <span>智能问答</span>
              <div class="header-actions">
                <KnowledgeSelector ref="knowledgeSelectorRef" class="knowledge-selector" />
                <el-button 
                  type="primary" 
                  size="small" 
                  plain 
                  class="new-session-btn"
                  @click="startNewSession">
                  <el-icon><Plus /></el-icon> 新会话
                </el-button>
              </div>
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
              <el-skeleton :rows="1" animated />
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
                @keydown="handleKeyDown">
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
import { ref, reactive, onMounted, nextTick } from 'vue'
import { ElMessage, ElNotification } from 'element-plus'
import { User, Service, ChatDotRound, ChatRound, Plus, Position, Setting, Document, CopyDocument } from '@element-plus/icons-vue'
import axios from 'axios'
import { marked } from 'marked'
import hljs from 'highlight.js'
import 'highlight.js/styles/github.css'
import { v4 as uuidv4 } from 'uuid'
import KnowledgeSelector from '../components/KnowledgeSelector.vue'

// 初始化marked配置
marked.setOptions({
  highlight: function (code, lang) {
    const language = hljs.getLanguage(lang) ? lang : 'plaintext';
    return hljs.highlight(code, { language }).value;
  },
  langPrefix: 'hljs language-',
  gfm: true,
  breaks: true
});

// 聊天消息列表
const messages = ref([]);
// 输入框消息
const inputMessage = ref('');
// 加载状态
const loading = ref(false);
// 消息容器引用
const messagesContainer = ref(null);
// 会话ID
const sessionId = ref(uuidv4());
// 知识库选择器引用
const knowledgeSelectorRef = ref(null);
// 参考文档
const references = ref([]);
// 显示设置面板
const showSettings = ref(false);
// 当前正在流式传输的消息
const currentStreamingMessage = ref('');
// 是否正在流式传输中
const isStreaming = ref(false);

// 聊天设置
const chatSettings = reactive({
  top_k: 3,
  score: 0.5
});

// 处理键盘事件
const handleKeyDown = (e) => {
  // 只有在按下Enter键且没有同时按下Shift键时才发送消息
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault(); // 阻止默认行为
    sendMessage();
  }
};

// 发送消息
const sendMessage = async () => {
  const message = inputMessage.value.trim();
  if (!message || loading.value) return;
  
  // 添加用户消息
  messages.value.push({
    role: 'user',
    content: message,
    timestamp: new Date()
  });
  
  // 清空输入框
  inputMessage.value = '';
  
  // 设置加载状态
  loading.value = true;
  currentStreamingMessage.value = '';
  isStreaming.value = true;
  
  // 添加AI消息占位
  messages.value.push({
    role: 'assistant',
    content: '',
    timestamp: new Date()
  });
  
  // 滚动到底部
  await nextTick();
  scrollToBottom();
  
  try {
    // 使用fetch API进行流式请求
    references.value = [];
    const response = await fetch('/v1/chat/stream', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        question: message,
        top_k: chatSettings.top_k,
        score: chatSettings.score,
        conv_id: sessionId.value,
        knowledge_name: knowledgeSelectorRef.value?.getSelectedKnowledgeId() || ''
      }),
    });
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    const reader = response.body.getReader();
    const decoder = new TextDecoder();
    
    // 读取流数据
    while (true) {
      const { value, done } = await reader.read();
      
      if (done) {
        break;
      }
      
      // 解码数据
      const chunk = decoder.decode(value, { stream: true });
      const lines = chunk.split('\n');
      
      for (const line of lines) {
        if (line.startsWith('data:')) {
          const data = line.slice(5).trim();

          if (data === '[DONE]') {
            // 流结束
            isStreaming.value = false;
            // 确保最后一次完整渲染
            messages.value[messages.value.length - 1].content = currentStreamingMessage.value;
            await nextTick();
            scrollToBottom();
            break;
          }

          try {
            const parsedData = JSON.parse(data);
            if (parsedData.content) {
              currentStreamingMessage.value += parsedData.content;
              // 更新最后一条消息的内容
              messages.value[messages.value.length - 1].content = currentStreamingMessage.value;
              await nextTick();
              scrollToBottom();
            }
          } catch (e) {
            console.error('解析流数据失败:', e);
          }
        }

        if (line.startsWith('documents:')) {
          const data = line.slice(10).trim();

          try {
            const parsedData = JSON.parse(data);
            if (parsedData.document) {
              references.value.push(...parsedData.document);
              console.log("references",references.value);
            }
          } catch (e) {
            console.error('解析流数据失败:', e);
          }
        }
      }
    }
    
    // // 获取参考文档
    // const refsResponse = await axios.post('/api/v1/chat', {
    //   session_id: sessionId.value,
    //   query: message,
    //   top_k: chatSettings.top_k,
    //   score: chatSettings.score
    // });
    //
    // references.value = refsResponse.data.references || [];
    
  } catch (error) {
    console.error('发送消息失败:', error);
    ElNotification({
      title: '错误',
      message: '发送消息失败，请稍后重试',
      type: 'error'
    });
    
    // 移除最后一条消息（AI回复）
    if (messages.value.length > 0 && messages.value[messages.value.length - 1].role === 'assistant') {
      messages.value.pop();
    }
  } finally {
    loading.value = false;
  }
};

// 开始新会话
const startNewSession = () => {
  if (messages.value.length > 0) {
    ElMessage({
      message: '已开始新的会话',
      type: 'success'
    });
  }
  
  messages.value = [];
  references.value = [];
  sessionId.value = uuidv4();
};

// 复制会话ID
const copySessionId = () => {
  navigator.clipboard.writeText(sessionId.value)
    .then(() => {
      ElMessage({
        message: '会话ID已复制到剪贴板',
        type: 'success'
      });
    })
    .catch(() => {
      ElMessage({
        message: '复制失败，请手动复制',
        type: 'error'
      });
    });
};

// 格式化时间
const formatTime = (timestamp) => {
  const date = new Date(timestamp);
  return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
};

// 渲染Markdown
const renderMarkdown = (text) => {
  if (!text) return '';
  try {
    // 尝试渲染markdown
    return marked(text);
  } catch (error) {
    console.error('Markdown渲染错误:', error);
    // 如果渲染失败，返回原始文本
    return text;
  }
};

// 滚动到底部
const scrollToBottom = () => {
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight;
  }
};

// 组件挂载后滚动到底部
onMounted(() => {
  scrollToBottom();
});
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

.header-actions {
  display: flex;
  align-items: center;
  margin-left: auto;
}

.knowledge-selector {
  margin-right: 10px;
}

.new-session-btn {
  margin-left: 5px;
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
  padding: 12px;
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
  margin-top: 16px;
  padding: 16px;
  background-color: var(--el-color-info-light-9);
  border-radius: 8px;
}

.session-info {
  margin-bottom: 16px;
  padding: 12px;
  background-color: var(--el-color-info-light-9);
  border-radius: 8px;
}

.session-id, .message-count {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
}

.session-id:last-child, .message-count:last-child {
  margin-bottom: 0;
}

.label {
  font-weight: bold;
  margin-right: 8px;
}

.references-content {
  flex: 1;
  overflow-y: auto;
}

.empty-references {
  padding: 20px;
  text-align: center;
}

.reference-list {
  margin-top: 12px;
}

.reference-content {
  padding: 8px;
}

.source-info {
  margin-bottom: 8px;
}

.content-text {
  white-space: pre-wrap;
  word-break: break-word;
  line-height: 1.5;
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

.markdown-content blockquote {
  border-left: 4px solid #d0d7de;
  padding-left: 1em;
  color: #57606a;
  margin: 1em 0;
}

/* 打字效果的光标动画 */
@keyframes cursor-blink {
  0%, 100% { opacity: 1; }
  50% { opacity: 0; }
}

/* 为最后一条AI消息添加光标效果，但仅在流式传输时显示 */
.ai-message:last-child .message-text:after {
  content: '|';
  display: inline-block;
  color: var(--el-color-primary);
  animation: cursor-blink 0.8s infinite;
  font-weight: bold;
  margin-left: 2px;
  /* 仅在流式传输时显示光标 */
  display: v-bind(isStreaming ? 'inline-block' : 'none');
}
</style>