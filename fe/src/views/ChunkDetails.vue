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
              <span>ES Chunk ID: {{ chunk.ChunkId }}</span>
              <el-space>
                <el-button 
                  text 
                  :icon="CopyDocument" 
                  @click="copyChunkContent(chunk.content)">
                  复制
                </el-button>
                <el-button 
                  text 
                  :icon="Edit" 
                  @click="handleEdit(chunk)">
                  编辑
                </el-button>
                <el-button
                  text
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
            <pre class="chunk-content-pre">{{ chunk.content }}</pre>
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

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { CopyDocument, Grid, Edit, Delete } from '@element-plus/icons-vue'
import request from '../utils/request'
import { formatDate } from '../utils/format'

const route = useRoute()
const router = useRouter()

const documentId = ref(route.params.documentId)
const documentInfo = ref(null)

const chunksList = ref([])
const chunksLoading = ref(false)
const chunksCurrentPage = ref(1)
const chunksPageSize = ref(6)
const chunksTotal = ref(0)

const editingChunkId = ref(null)
const editedContent = ref('')
const isSaving = ref(false)

const pageTitle = computed(() => {
  if (documentInfo.value) {
    return `文档 "${documentInfo.value.fileName}" 的分块详情`
  }
  return '文档分块详情'
})

const goBack = () => {
  router.back()
}

// 获取分块列表
const fetchChunksList = async () => {
  if (!documentId.value) return
  chunksLoading.value = true
  try {
    const response = await request.get('/v1/chunks', {
      params: {
        knowledge_doc_id: documentId.value,
        page: chunksCurrentPage.value,
        size: chunksPageSize.value
      }
    })
    chunksList.value = response.data.data || []
    chunksTotal.value = response.data.total || 0
  } catch (error) {
    console.error('获取分块列表失败:', error)
  } finally {
    chunksLoading.value = false
  }
}

const handleChunksPageChange = (page) => {
  chunksCurrentPage.value = page
  fetchChunksList()
}

const copyChunkContent = async (content) => {
  if (!content) {
    ElMessage.warning('没有内容可复制')
    return
  }
  try {
    await navigator.clipboard.writeText(content)
    ElMessage.success('内容已复制到剪贴板')
  } catch (error) {
    ElMessage.error('复制失败，请手动复制')
  }
}

const handleEdit = (chunk) => {
  editingChunkId.value = chunk.id
  editedContent.value = chunk.content
}

const handleCancelEdit = () => {
  editingChunkId.value = null
  editedContent.value = ''
}

const handleSaveEdit = async (chunk) => {
  isSaving.value = true
  try {
    await request.put('/v1/chunks_content', {
      id: chunk.id,
      content: editedContent.value
    })
    // 更新前端数据
    const chunkToUpdate = chunksList.value.find(c => c.id === chunk.id)
    if (chunkToUpdate) {
      chunkToUpdate.content = editedContent.value
    }
    ElMessage.success('内容更新成功！')
    handleCancelEdit()
  } catch (error) {
    ElMessage.error('更新失败，请重试。')
    console.error('更新分块内容失败:', error)
  } finally {
    isSaving.value = false
  }
}

const handleDeleteChunk = async (chunk) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除分块ID为 ${chunk.ChunkId} 的内容吗？此操作不可恢复。`,
      '确认删除',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    await request.delete('/v1/chunks', { params: { id: chunk.id } })
    ElMessage.success('分块删除成功！')
    fetchChunksList()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('分块删除失败，请重试。')
      console.error('分块删除失败:', error)
    }
  }
}

onMounted(() => {
  const docFromState = history.state.document
  if (docFromState && docFromState.id == documentId.value) {
    documentInfo.value = docFromState
  } else {
    ElMessage.warning('文档信息不完整，请从文档列表页进入。')
    router.push('/knowledge-documents')
  }
  fetchChunksList()
})

</script>

<style scoped>
.chunk-details-container {
  /* max-width: 1200px; */
  /* margin: 0 auto; */
}
.page-header {
  margin-bottom: 20px;
  background-color: #fff;
  padding: 16px;
  border-radius: 4px;
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
}
.chunk-content-pre {
  white-space: pre-wrap;
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
.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}
.edit-actions {
  width: 100%;
}
</style> 