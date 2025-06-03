<template>
  <div class="document-list">
    <div class="header">
      <div class="actions">
        <el-button type="primary" @click="handleAdd">上传文档</el-button>
      </div>
    </div>

    <el-card class="list-card">
      <el-table :data="documentList" v-loading="loading" border>
        <el-table-column prop="documentName" label="文档名称" min-width="200" />
        <el-table-column prop="documentType" label="文档类型" width="100">
          <template #default="{ row }">
            <el-tag :type="getDocumentTypeTag(row.documentType)">
              {{ row.documentType }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="documentSize" label="文档大小" width="120">
          <template #default="{ row }">
            {{ formatFileSize(row.documentSize) }}
          </template>
        </el-table-column>
        <el-table-column prop="chunkCount" label="分块数量" width="100" />
        <el-table-column prop="createTime" label="创建时间" width="180" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button-group>
              <el-button type="primary" link @click="handleView(row)">查看</el-button>
              <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
              <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
            </el-button-group>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 上传文档对话框 -->
    <el-dialog
      v-model="uploadDialogVisible"
      title="上传文档"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form ref="uploadFormRef" :model="uploadForm" :rules="uploadRules" label-width="100px">
        <el-form-item label="知识库" prop="knowledgeBaseId">
          <el-select v-model="uploadForm.knowledgeBaseId" placeholder="请选择知识库">
            <el-option
              v-for="item in knowledgeBaseList"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="文档" prop="file">
          <el-upload
            class="upload-demo"
            drag
            action="/api/v1/document/upload"
            :headers="uploadHeaders"
            :data="uploadData"
            :on-success="handleUploadSuccess"
            :on-error="handleUploadError"
            :before-upload="beforeUpload"
          >
            <el-icon class="el-icon--upload"><upload-filled /></el-icon>
            <div class="el-upload__text">
              将文件拖到此处，或<em>点击上传</em>
            </div>
            <template #tip>
              <div class="el-upload__tip">
                支持 .md、.pdf、.html 格式文件
              </div>
            </template>
          </el-upload>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="uploadDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleUploadSubmit">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { UploadFilled } from '@element-plus/icons-vue'
import { getDocumentList, deleteDocument } from '@/api/document'
import { getKnowledgeBaseList } from '@/api/knowledge'
import { formatFileSize } from '@/utils/format'

// 数据列表
const documentList = ref([])
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 上传相关
const uploadDialogVisible = ref(false)
const uploadFormRef = ref(null)
const uploadForm = ref({
  knowledgeBaseId: '',
  file: null
})
const uploadRules = {
  knowledgeBaseId: [{ required: true, message: '请选择知识库', trigger: 'change' }]
}
const knowledgeBaseList = ref([])
const uploadHeaders = {
  Authorization: `Bearer ${localStorage.getItem('token')}`
}
const uploadData = ref({})

// 获取文档列表
const fetchDocumentList = async () => {
  loading.value = true
  try {
    const res = await getDocumentList({
      page: currentPage.value,
      pageSize: pageSize.value
    })
    documentList.value = res.list
    total.value = res.total
  } catch (error) {
    ElMessage.error('获取文档列表失败')
  } finally {
    loading.value = false
  }
}

// 获取知识库列表
const fetchKnowledgeBaseList = async () => {
  try {
    const res = await getKnowledgeBaseList()
    knowledgeBaseList.value = res.list
  } catch (error) {
    ElMessage.error('获取知识库列表失败')
  }
}

// 处理分页
const handleSizeChange = (val) => {
  pageSize.value = val
  fetchDocumentList()
}

const handleCurrentChange = (val) => {
  currentPage.value = val
  fetchDocumentList()
}

// 处理文档操作
const handleAdd = () => {
  uploadForm.value = {
    knowledgeBaseId: '',
    file: null
  }
  uploadDialogVisible.value = true
}

const handleView = (row) => {
  // TODO: 实现查看文档功能
}

const handleEdit = (row) => {
  // TODO: 实现编辑文档功能
}

const handleDelete = (row) => {
  ElMessageBox.confirm('确定要删除该文档吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await deleteDocument(row.id)
      ElMessage.success('删除成功')
      fetchDocumentList()
    } catch (error) {
      ElMessage.error('删除失败')
    }
  })
}

// 处理上传
const beforeUpload = (file) => {
  const isValidType = ['.md', '.pdf', '.html'].includes(file.name.toLowerCase())
  if (!isValidType) {
    ElMessage.error('只能上传 .md、.pdf、.html 格式的文件')
    return false
  }
  uploadData.value = {
    knowledgeBaseId: uploadForm.value.knowledgeBaseId
  }
  return true
}

const handleUploadSuccess = (response) => {
  if (response.code === 0) {
    ElMessage.success('上传成功')
    uploadDialogVisible.value = false
    fetchDocumentList()
  } else {
    ElMessage.error(response.message || '上传失败')
  }
}

const handleUploadError = () => {
  ElMessage.error('上传失败')
}

const handleUploadSubmit = async () => {
  if (!uploadFormRef.value) return
  await uploadFormRef.value.validate((valid) => {
    if (valid) {
      // 上传组件会自动处理上传
    }
  })
}

// 获取文档类型标签样式
const getDocumentTypeTag = (type) => {
  const map = {
    md: '',
    pdf: 'success',
    html: 'warning'
  }
  return map[type] || 'info'
}

onMounted(() => {
  fetchDocumentList()
  fetchKnowledgeBaseList()
})
</script>

<style scoped>
.document-list {
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.list-card {
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.upload-demo {
  width: 100%;
}
</style> 