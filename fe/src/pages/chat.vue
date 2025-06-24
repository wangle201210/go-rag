<script setup>
import {
  ChatDotRound,
  ChatRound,
  CopyDocument,
  Document,
  Plus,
  Position,
  Service,
  Setting,
  User,
} from '@element-plus/icons-vue';
import { ElMessage, ElNotification } from 'element-plus';
import { nextTick, onMounted, reactive, ref } from 'vue';
import { v4 as uuidv4 } from 'uuid';
import KnowledgeSelector from '../components/KnowledgeSelector.vue';
import '~/styles/markdown.css';
import { renderMarkdown } from '~/utils/markdown';

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
  score: 0.5,
});

function scrollToBottom() {
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight;
  }
}

// 发送消息
async function sendMessage() {
  const message = inputMessage.value.trim();
  if (!message || loading.value) {
    return;
  }

  // 添加用户消息
  messages.value.push({
    role: 'user',
    content: message,
    timestamp: new Date(),
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
    timestamp: new Date(),
  });

  // 滚动到底部
  await nextTick();
  scrollToBottom();

  try {
    // 使用fetch API进行流式请求
    references.value = [];
    const response = await fetch('/api/v1/chat/stream', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        question: message,
        top_k: chatSettings.top_k,
        score: chatSettings.score,
        conv_id: sessionId.value,
        knowledge_name: knowledgeSelectorRef.value?.getSelectedKnowledgeId() || '',
      }),
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const reader = response.body.getReader();
    const decoder = new TextDecoder();

    // eslint-disable-next-line no-constant-condition
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
            // eslint-disable-next-line no-console
            console.error('解析流数据失败:', e);
          }
        }

        if (line.startsWith('documents:')) {
          const data = line.slice(10).trim();

          try {
            const parsedData = JSON.parse(data);
            if (parsedData.document) {
              references.value.push(...parsedData.document);
              // eslint-disable-next-line no-console
              console.log('references', references.value);
            }
          } catch (e) {
            // eslint-disable-next-line no-console
            console.error('解析流数据失败:', e);
          }
        }
      }
    }
  } catch (error) {
    // eslint-disable-next-line no-console
    console.error('发送消息失败:', error);
    ElNotification({
      title: '错误',
      message: '发送消息失败，请稍后重试',
      type: 'error',
    });

    // 移除最后一条消息（AI回复）
    if (messages.value.length > 0 && messages.value[messages.value.length - 1].role === 'assistant') {
      messages.value.pop();
    }
  } finally {
    loading.value = false;
  }
}

// 处理键盘事件
function handleKeyDown(e) {
  // 只有在按下Enter键且没有同时按下Shift键时才发送消息
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault(); // 阻止默认行为
    sendMessage();
  }
}

// 开始新会话
function startNewSession() {
  if (messages.value.length > 0) {
    ElMessage({
      message: '已开始新的会话',
      type: 'success',
    });
  }

  messages.value = [];
  references.value = [];
  sessionId.value = uuidv4();
}

// 复制会话ID
function copySessionId() {
  navigator.clipboard.writeText(sessionId.value)
    .then(() => {
      ElMessage({
        message: '会话ID已复制到剪贴板',
        type: 'success',
      });
    })
    .catch(() => {
      ElMessage({
        message: '复制失败，请手动复制',
        type: 'error',
      });
    });
}

// 格式化时间
function formatTime(timestamp) {
  const date = new Date(timestamp);
  return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
}

// 组件挂载后滚动到底部
onMounted(() => {
  scrollToBottom();
});
</script>

<template>
  <div class="chat-container">
    <el-row>
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
                  @click="startNewSession"
                >
                  <el-icon><Plus /></el-icon> 新会话
                </el-button>
              </div>
            </div>
          </template>

          <div ref="messagesContainer" class="chat-messages">
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
                class="message-item"
                :class="[message.role === 'user' ? 'user-message' : 'ai-message']"
              >
                <div class="message-avatar">
                  <el-avatar :icon="message.role === 'user' ? User : Service" :size="36" />
                </div>
                <div class="message-content">
                  <div v-if="message.role === 'user'" class="message-text">{{ message.content }}</div>
                  <div v-else class="message-text markdown-content" v-html="renderMarkdown(message.content)" />
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
                @keydown="handleKeyDown"
              />
              <div class="input-actions">
                <el-tooltip content="高级设置" placement="top">
                  <el-button
                    type="info"
                    plain
                    circle
                    @click="showSettings = !showSettings"
                  >
                    <el-icon><Setting /></el-icon>
                  </el-button>
                </el-tooltip>
                <el-button
                  type="primary"
                  :loading="loading"
                  :disabled="!inputMessage.trim()"
                  @click="sendMessage"
                >
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
                  @click="copySessionId"
                >
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
                  :name="index"
                >
                  <div class="reference-content">
                    <div class="source-info">
                      <el-tag size="small">{{ ref.meta_data.ext._file_name || '未知来源' }}</el-tag>
                    </div>
                    <div class="content-text markdown-content" v-html="renderMarkdown(ref.content)" />
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
  margin: 10px;
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





.chat-input {
  margin-top: auto;
}

.references-content {
  flex: 1;
  overflow-y: auto;
}


/* 页面特定的Markdown样式扩展 */
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