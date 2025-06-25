<template>
  <div class="kb-container">
    <el-card class="kb-card">
      <template #header>
        <div class="card-header">
          <el-icon class="header-icon"><Folder /></el-icon>
          <span>知识库管理</span>
          <el-button 
            type="primary" 
            size="small" 
            plain 
            class="add-kb-btn"
            @click="showAddDialog">
            <el-icon><Plus /></el-icon> 新建知识库
          </el-button>
        </div>
      </template>
      
      <!-- 知识库列表 -->
      <div class="kb-list">
        <el-table 
          v-loading="loading" 
          :data="knowledgeBaseList" 
          style="width: 100%"
          border>
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column prop="name" label="知识库名称" width="180" />
          <el-table-column prop="description" label="描述" />
          <el-table-column prop="category" label="分类" width="120" />
          <el-table-column prop="status" label="状态" width="100">
            <template #default="scope">
              <el-tag :type="scope.row.status === 2 ? 'danger': 'success'">
                {{ scope.row.status === 2 ? '禁用': '启用' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="200">
            <template #default="scope">
              <el-button 
                size="small" 
                type="primary" 
                @click="showEditDialog(scope.row)"
                plain>
                编辑
              </el-button>
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
      </div>
      
      <!-- 空状态 -->
      <div v-if="!loading && knowledgeBaseList.length === 0" class="empty-kb">
        <el-empty description="暂无知识库，请点击右上角新建">
          <template #image>
            <el-icon class="empty-icon"><Folder /></el-icon>
          </template>
        </el-empty>
      </div>
    </el-card>
    
    <!-- 新建/编辑知识库对话框 -->
    <el-dialog 
      v-model="dialogVisible" 
      :title="isEdit ? '编辑知识库' : '新建知识库'"
      width="500px">
      <el-form 
        :model="kbForm" 
        :rules="rules" 
        ref="kbFormRef" 
        label-width="100px">
        <el-form-item label="知识库名称" prop="name">
          <el-input v-model="kbForm.name" placeholder="请输入知识库名称" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input 
            v-model="kbForm.description" 
            type="textarea" 
            :rows="3" 
            placeholder="请输入知识库描述" />
        </el-form-item>
        <el-form-item label="分类" prop="category">
          <el-input v-model="kbForm.category" placeholder="请输入知识库分类" />
        </el-form-item>
        <el-form-item label="状态" prop="status" v-if="isEdit">
          <el-radio-group v-model="kbForm.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="2">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitForm" :loading="submitting">
            确认
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Folder, Plus } from '@element-plus/icons-vue'
import request from '../../utils/request.js'

// 知识库列表
const knowledgeBaseList = ref([])
// 加载状态
const loading = ref(false)
// 提交状态
const submitting = ref(false)
// 对话框可见性
const dialogVisible = ref(false)
// 是否为编辑模式
const isEdit = ref(false)
// 表单引用
const kbFormRef = ref(null)

// 表单数据
const kbForm = reactive({
  id: null,
  name: '',
  description: '',
  category: '',
  status: 1
})

// 表单验证规则
const rules = {
  name: [
    { required: true, message: '请输入知识库名称', trigger: 'blur' },
    { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  description: [
    { required: true, message: '请输入知识库描述', trigger: 'blur' },
    { min: 3, max: 200, message: '长度在 3 到 200 个字符', trigger: 'blur' }
  ],
  category: [
    { min: 3, max: 10, message: '长度在 3 到 10 个字符', trigger: 'blur' }
  ]
}

// 页面加载时获取知识库列表
onMounted(() => {
  fetchKnowledgeBaseList()
})

// 获取知识库列表
const fetchKnowledgeBaseList = async () => {
  loading.value = true
  try {
    const response = await request.get('/v1/kb')
    knowledgeBaseList.value = response.data.list || []
  } catch (error) {
    console.error('获取知识库列表失败:', error)
    ElMessage.error('获取知识库列表失败: ' + (error.response?.data?.message || '未知错误'))
  } finally {
    loading.value = false
  }
}

// 显示新建对话框
const showAddDialog = () => {
  isEdit.value = false
  resetForm()
  dialogVisible.value = true
}

// 显示编辑对话框
const showEditDialog = (row) => {
  isEdit.value = true
  resetForm()
  // 复制数据到表单
  Object.assign(kbForm, row)
  dialogVisible.value = true
}

// 重置表单
const resetForm = () => {
  kbForm.id = null
  kbForm.name = ''
  kbForm.description = ''
  kbForm.category = ''
  kbForm.status = 1
  
  // 重置表单验证
  if (kbFormRef.value) {
    kbFormRef.value.resetFields()
  }
}

// 提交表单
const submitForm = async () => {
  if (!kbFormRef.value) return
  
  await kbFormRef.value.validate(async (valid) => {
    if (!valid) return
    
    submitting.value = true
    try {
      if (isEdit.value) {
        // 编辑知识库
        await request.put(`/v1/kb/${kbForm.id}`, {
          name: kbForm.name,
          description: kbForm.description,
          category: kbForm.category,
          status: kbForm.status
        })
        ElMessage.success('知识库更新成功')
      } else {
        // 创建知识库
        await request.post('/v1/kb', {
          name: kbForm.name,
          description: kbForm.description,
          category: kbForm.category,
        })
        ElMessage.success('知识库创建成功')
      }
      
      // 关闭对话框并刷新列表
      dialogVisible.value = false
      fetchKnowledgeBaseList()
    } catch (error) {
      console.error('操作失败:', error)
      ElMessage.error('操作失败: ' + (error.response?.data?.message || '未知错误'))
    } finally {
      submitting.value = false
    }
  })
}

// 确认删除
const confirmDelete = (row) => {
  ElMessageBox.confirm(
    `确定要删除知识库 "${row.name}" 吗？此操作不可恢复。`,
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(async () => {
    try {
      await request.delete(`/v1/kb/${row.id}`)
      ElMessage.success('知识库删除成功')
      fetchKnowledgeBaseList()
    } catch (error) {
      console.error('删除失败:', error)
      ElMessage.error('删除失败: ' + (error.response?.data?.message || '未知错误'))
    }
  }).catch(() => {
    // 用户取消删除
  })
}
</script>

<style scoped>
.kb-container {
  margin: 10px;
}

.kb-card {
  margin-bottom: 20px;
}

.add-kb-btn {
  margin-left: auto;
}

.kb-list {
  margin-top: 20px;
}

.empty-kb {
  height: 300px;
}
</style>