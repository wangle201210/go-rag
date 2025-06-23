<template>
  <div class="documents-container">
    <el-card class="documents-card">
      <template #header>
        <div class="card-header">
          <el-icon class="header-icon"><Document /></el-icon>
          <span>文档数据集管理</span>
          <div class="header-actions">
            <el-select 
              v-model="selectedKnowledgeBase" 
              placeholder="请选择知识库" 
              style="width: 200px;"
              @change="handleKnowledgeBaseChange">
              <el-option
                v-for="kb in knowledgeBaseList"
                :key="kb.id"
                :label="kb.name"
                :value="kb.name" />
            </el-select>
          </div>
        </div>
      </template>
      
      <!-- 文档列表 -->
      <div class="documents-list" v-if="selectedKnowledgeBase">
        <el-table 
          v-loading="loading" 
          :data="documentsList" 
          style="width: 100%"
          border>
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column prop="fileName" label="文件名" min-width="200">
            <template #default="scope">
              <router-link 
                :to="{ 
                  name: 'ChunkDetails', 
                  params: { documentId: scope.row.id }, 
                  state: { 
                    document: {
                      id: scope.row.id,
                      fileName: scope.row.fileName,
                      knowledgeBaseName: scope.row.knowledgeBaseName
                    } 
                  } 
                }" 
                class="document-link">
                {{ scope.row.fileName }}
              </router-link>
            </template>
          </el-table-column>
          <el-table-column prop="knowledgeBaseName" label="所属知识库"/>
          <el-table-column prop="status" label="状态" width="120">
            <template #default="scope">
              <el-tag :type="getStatusType(scope.row.status)">
                {{ getStatusText(scope.row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="createdAt" label="创建时间">
            <template #default="scope">
              {{ formatDate(scope.row.createdAt) }}
            </template>
          </el-table-column>
          <el-table-column prop="updatedAt" label="更新时间">
            <template #default="scope">
              {{ formatDate(scope.row.updatedAt) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" fixed="right">
            <template #default="scope">
              <el-button 
                size="small" 
                type="danger" 
                @click="confirmDelete(scope.row)"
                plain>
                删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>
        
        <!-- 分页 -->
        <div class="pagination-container" v-if="total > 0">
          <el-pagination
            :current-page="currentPage"
            :page-size="pageSize"
            :page-sizes="[10, 20, 50, 100]"
            :total="total"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange" />
        </div>
      </div>
      
      <!-- 空状态 -->
      <div v-if="!loading && (!selectedKnowledgeBase || documentsList.length === 0)" class="empty-documents">
        <el-empty :description="selectedKnowledgeBase ? '暂无文档数据' : '请先选择知识库'">
          <template #image>
            <el-icon class="empty-icon"><Document /></el-icon>
          </template>
        </el-empty>
      </div>
    </el-card>
    
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Document } from '@element-plus/icons-vue'
import request from '../utils/request'
import { getStatusType, getStatusText, formatDate } from '../utils/format'

const KEY_LAST_KB = 'last_selected_kb'

// 知识库列表
const knowledgeBaseList = ref([])
// 选中的知识库
const selectedKnowledgeBase = ref('')
// 文档列表
const documentsList = ref([])
// 加载状态
const loading = ref(false)
// 分页相关
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 页面加载时获取知识库列表
onMounted(async () => {
  await fetchKnowledgeBaseList()
  
  const lastSelectedKB = sessionStorage.getItem(KEY_LAST_KB)
  if (lastSelectedKB && knowledgeBaseList.value.some(kb => kb.name === lastSelectedKB)) {
    selectedKnowledgeBase.value = lastSelectedKB
    await fetchDocumentsList()
  }
})

// 获取知识库列表
const fetchKnowledgeBaseList = async () => {
  try {
    const response = await request.get('/v1/kb')
    knowledgeBaseList.value = response.data.list || []
  } catch (error) {
    console.error('获取知识库列表失败:', error)
    ElMessage.error('获取知识库列表失败: ' + (error.response?.message || '未知错误'))
  }
}

// 知识库选择变化
const handleKnowledgeBaseChange = () => {
  currentPage.value = 1
  fetchDocumentsList()
  sessionStorage.setItem(KEY_LAST_KB, selectedKnowledgeBase.value)
}

// 获取文档列表
const fetchDocumentsList = async () => {
  if (!selectedKnowledgeBase.value) {
    documentsList.value = []
    total.value = 0
    return
  }
  
  loading.value = true
  try {
    const response = await request.get('/v1/documents', {
      params: {
        knowledge_name: selectedKnowledgeBase.value,
        page: currentPage.value,
        size: pageSize.value
      }
    })
    documentsList.value = response.data.data || []
    total.value = response.data.total || 0
  } catch (error) {
    console.error('获取文档列表失败:', error)
    ElMessage.error('获取文档列表失败: ' + (error.response?.data?.message || '未知错误'))
  } finally {
    loading.value = false
  }
}

// 确认删除
const confirmDelete = async (document) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除文档 "${document.fileName}" 吗？此操作将一并删除该文档下的所有数据分块，且不可恢复。`,
      '确认删除',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await request.delete('/v1/documents', { params: { document_id: document.id } })
    
    ElMessage.success(`文档 "${document.fileName}" 删除成功`)
    fetchDocumentsList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
      // 错误消息已由 request 拦截器统一处理
    }
  }
}

// 分页大小变化
const handleSizeChange = (size) => {
  pageSize.value = size
  currentPage.value = 1
  fetchDocumentsList()
}

// 当前页变化
const handleCurrentChange = (page) => {
  currentPage.value = page
  fetchDocumentsList()
}

</script>

<style scoped>
.documents-container {
  height: 100%;
}
.documents-card {
  height: 100%;
}
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.header-icon {
  margin-right: 8px;
  font-size: 18px;
}
.header-actions {
  display: flex;
  align-items: center;
}
.documents-list {
  margin-top: 20px;
}
.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}
.empty-documents {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 300px;
}
.empty-icon {
  font-size: 60px;
  color: #c0c4cc;
}
.document-link {
  color: #409eff;
  text-decoration: none;
}
.document-link:hover {
  text-decoration: underline;
}
</style> 