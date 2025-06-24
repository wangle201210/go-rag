<script setup>
import { Delete, Edit, CopyDocument } from '@element-plus/icons-vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { computed, onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import '~/styles/markdown.css';
import { formatDate } from '~/utils/format.js';
import { renderMarkdown } from '~/utils/markdown.js';
import request from '~/utils/request.js';

const route = useRoute();
const router = useRouter();

const documentId = ref(route.params.documentId);
const documentInfo = ref(null);

const chunksList = ref([]);
const chunksLoading = ref(false);
const chunksCurrentPage = ref(1);
const chunksPageSize = ref(6);
const chunksTotal = ref(0);

const editingChunkId = ref(null);
const editedContent = ref('');
const isSaving = ref(false);

const pageTitle = computed(() => {
  if (documentInfo.value) {
    return `文档 "${documentInfo.value.fileName}" 的分块详情`;
  }
  return '文档分块详情';
});

function goBack() {
  router.back();
}

// 获取分块列表
async function fetchChunksList() {
  if (!documentId.value) return;
  chunksLoading.value = true;
  try {
    const response = await request.get('/v1/chunks', {
      params: {
        knowledge_doc_id: documentId.value,
        page: chunksCurrentPage.value,
        size: chunksPageSize.value,
      },
    });
    chunksList.value = response.data.data || [];
    chunksTotal.value = response.data.total || 0;
  } catch (error) {
    // eslint-disable-next-line no-console
    console.error('获取分块列表失败:', error);
  } finally {
    chunksLoading.value = false;
  }
}

function handleChunksPageChange(page) {
  chunksCurrentPage.value = page;
  fetchChunksList();
}

async function copyChunkContent(content) {
  if (!content) {
    ElMessage.warning('没有内容可复制');
    return;
  }
  try {
    await navigator.clipboard.writeText(content);
    ElMessage.success('内容已复制到剪贴板');
  } catch (error) {
    ElMessage.error('复制失败，请手动复制');
  }
}

function handleEdit(chunk) {
  editingChunkId.value = chunk.id;
  editedContent.value = chunk.content;
}

function handleCancelEdit() {
  editingChunkId.value = null;
  editedContent.value = '';
}

async function handleSaveEdit(chunk) {
  isSaving.value = true;
  try {
    await request.put('/v1/chunks_content', {
      id: chunk.id,
      content: editedContent.value,
    });
    // 更新前端数据
    const chunkToUpdate = chunksList.value.find((c) => c.id === chunk.id);
    if (chunkToUpdate) {
      chunkToUpdate.content = editedContent.value;
    }
    ElMessage.success('内容更新成功！');
    handleCancelEdit();
  } catch (error) {
    ElMessage.error('更新失败，请重试。');
    // eslint-disable-next-line no-console
    console.error('更新分块内容失败:', error);
  } finally {
    isSaving.value = false;
  }
}

async function handleDeleteChunk(chunk) {
  try {
    await ElMessageBox.confirm(
      `确定要删除分块ID为 ${chunk.chunkId} 的内容吗？此操作不可恢复。`,
      '确认删除',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning',
      },
    );
    await request.delete('/v1/chunks', { params: { id: chunk.id } });
    ElMessage.success('分块删除成功！');
    fetchChunksList();
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('分块删除失败，请重试。');
      // eslint-disable-next-line no-console
      console.error('分块删除失败:', error);
    }
  }
}

onMounted(() => {
  const docFromStorage = localStorage.getItem(`document-${documentId.value}`);
  // eslint-disable-next-line no-console
  console.log('docFromStorage', docFromStorage);
  // eslint-disable-next-line no-console
  console.log('route.params', route.params);
  if (docFromStorage) {
    try {
      const docData = JSON.parse(docFromStorage);
      if (docData && docData.id == documentId.value) {
        documentInfo.value = docData;
      } else {
        ElMessage.warning('文档信息不匹配，请从文档列表页进入。');
        router.push('/knowledge-documents');
      }
    } catch (error) {
      // eslint-disable-next-line no-console
      console.error('解析文档数据失败:', error);
      ElMessage.warning('文档信息格式错误，请从文档列表页进入。');
      router.push('/knowledge-documents');
    }
  } else {
    ElMessage.warning('文档信息不完整，请从文档列表页进入。');
    router.push('/knowledge-documents');
  }
  fetchChunksList();
});
</script>

<template>
  <div class="chunk-details-container">
    <el-page-header @back="goBack" class="page-header">
      <template #content>
        <span class="text-large font-600 mr-3"> {{ pageTitle }} </span>
      </template>
    </el-page-header>
    <div v-if="chunksLoading" class="loading-container">
      <el-skeleton :rows="5" animated />
    </div>
    <div v-else-if="chunksList.length > 0">
      <el-card v-for="chunk in chunksList" :key="chunk.id" class="chunk-item-card">
        <template #header>
          <div class="chunk-card-header">
            <span>Chunk ID: {{ chunk.chunkId }}</span>
            <el-space>
              <el-button
                text
                size="small"
                :icon="CopyDocument"
                @click="copyChunkContent(chunk.content)">
                复制
              </el-button>
              <el-button 
                text
                size="small"
                :icon="Edit"
                @click="handleEdit(chunk)">
                编辑
              </el-button>
              <el-button
                text
                size="small"
                type="danger"
                :icon="Delete"
                @click="handleDeleteChunk(chunk)"
              >
                删除
              </el-button>
            </el-space>
          </div>
        </template>
        <el-input
          v-if="editingChunkId === chunk.id"
          v-model="editedContent"
          type="textarea"
          :rows="8"
          class="chunk-content-textarea"
        />
        <el-scrollbar v-else class="chunk-content-scrollbar">
          <div class="markdown-content chunk-content-pre" v-html="renderMarkdown(chunk.content)"></div>
        </el-scrollbar>

        <div class="chunk-card-footer">
          <div v-if="editingChunkId === chunk.id" class="edit-actions">
            <el-button @click="handleCancelEdit">取消</el-button>
            <el-button type="primary" @click="handleSaveEdit(chunk)" :loading="isSaving">保存</el-button>
          </div>
          <span v-else>创建于: {{ formatDate(chunk.createdAt) }}</span>
        </div>
      </el-card>
      <div class="pagination-container" v-if="chunksTotal > chunksPageSize">
        <el-pagination
          :current-page="chunksCurrentPage"
          :page-size="chunksPageSize"
          :total="chunksTotal"
          layout="total, prev, pager, next"
          @current-change="handleChunksPageChange" />
      </div>
    </div>
    <el-empty v-else description="该文档下暂无分块数据"></el-empty>

  </div>
</template>

<style scoped>
.chunk-details-container {
   margin: 10px;
}

.chunk-item-card {
  box-shadow: var(--el-box-shadow-light);
  margin-bottom: 20px;
}

.chunk-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 14px;
  color: #606266;
}

.chunk-content-scrollbar {
  height: 200px;
  text-align: left;
}

.chunk-content-pre {
  white-space: normal;
  word-wrap: break-word;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 13px;
  line-height: 1.6;
  background-color: #f8f9fa;
  padding: 10px;
  border-radius: 4px;
  color: #495057;
  margin: 0;
}

.chunk-content-textarea {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 13px;
  line-height: 1.6;
}

.chunk-card-footer {
  margin-top: 15px;
  text-align: right;
  font-size: 12px;
  color: #909399;
}
</style>