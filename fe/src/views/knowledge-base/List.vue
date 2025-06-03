<template>
  <div class="knowledge-base-container">
    <div class="page-header">
      <h2>知识库管理</h2>
      <el-button type="primary" @click="openCreateDialog">创建知识库</el-button>
    </div>

    <el-table :data="knowledgeBaseList" style="width: 100%" v-loading="loading">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="知识库名称" />
      <el-table-column prop="description" label="描述" />
      <el-table-column prop="created_at" label="创建时间" />
      <el-table-column label="操作" width="200">
        <template #default="scope">
          <el-button size="small" @click="viewDocuments(scope.row)">查看文档</el-button>
          <el-button size="small" type="danger" @click="deleteKnowledgeBase(scope.row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 创建知识库对话框 -->
    <el-dialog v-model="dialogVisible" title="创建知识库" width="500px">
      <el-form :model="form" label-width="120px">
        <el-form-item label="知识库名称" required>
          <el-input v-model="form.name" placeholder="请输入知识库名称" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" placeholder="请输入知识库描述" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="createKnowledgeBase">创建</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useRouter } from 'vue-router'

const router = useRouter()
const knowledgeBaseList = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const form = ref({
  name: '',
  description: ''
})

// 获取知识库列表
const fetchKnowledgeBaseList = async () => {
  loading.value = true
  try {
    // 这里应该替换为实际的API调用
    // const response = await fetch('/api/knowledge-base')
    // const data = await response.json()
    // knowledgeBaseList.value = data
    
    // 模拟数据
    setTimeout(() => {
      knowledgeBaseList.value = [
        { id: 1, name: '通用知识库', description: '包含常见问题和解答', created_at: '2023-05-01' },
        { id: 2, name: '产品文档', description: '产品使用说明和文档', created_at: '2023-05-10' }
      ]
      loading.value = false
    }, 500)
  } catch (error) {
    console.error('获取知识库列表失败:', error)
    ElMessage.error('获取知识库列表失败')
    loading.value = false
  }
}

// 打开创建对话框
const openCreateDialog = () => {
  form.value = {
    name: '',
    description: ''
  }
  dialogVisible.value = true
}

// 创建知识库
const createKnowledgeBase = async () => {
  if (!form.value.name) {
    ElMessage.warning('请输入知识库名称')
    return
  }
  
  try {
    // 这里应该替换为实际的API调用
    // const response = await fetch('/api/knowledge-base', {
    //   method: 'POST',
    //   headers: {
    //     'Content-Type': 'application/json'
    //   },
    //   body: JSON.stringify(form.value)
    // })
    // const data = await response.json()
    
    // 模拟创建成功
    setTimeout(() => {
      dialogVisible.value = false
      ElMessage.success('创建知识库成功')
      fetchKnowledgeBaseList()
    }, 500)
  } catch (error) {
    console.error('创建知识库失败:', error)
    ElMessage.error('创建知识库失败')
  }
}

// 查看文档
const viewDocuments = (row) => {
  router.push({
    path: '/document',
    query: { knowledgeBaseId: row.id }
  })
}

// 删除知识库
const deleteKnowledgeBase = (row) => {
  ElMessageBox.confirm(
    `确定要删除知识库「${row.name}」吗？删除后将无法恢复。`,
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(async () => {
    try {
      // 这里应该替换为实际的API调用
      // await fetch(`/api/knowledge-base/${row.id}`, {
      //   method: 'DELETE'
      // })
      
      // 模拟删除成功
      setTimeout(() => {
        ElMessage.success('删除知识库成功')
        fetchKnowledgeBaseList()
      }, 500)
    } catch (error) {
      console.error('删除知识库失败:', error)
      ElMessage.error('删除知识库失败')
    }
  }).catch(() => {
    // 取消删除
  })
}

onMounted(() => {
  fetchKnowledgeBaseList()
})
</script>

<style scoped>
.knowledge-base-container {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
}
</style>