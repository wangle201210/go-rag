<template>
  <div class="knowledge-list">
    <div class="header">
      <el-input
        v-model="searchForm.keyword"
        placeholder="请输入知识库名称"
        class="search-input"
        @keyup.enter="handleSearch"
      />
      <el-select v-model="searchForm.category" placeholder="请选择分类" class="category-select">
        <el-option
          v-for="item in categoryOptions"
          :key="item.value"
          :label="item.label"
          :value="item.value"
        />
      </el-select>
      <el-button type="primary" @click="handleSearch">搜索</el-button>
      <el-button type="success" @click="handleAdd">新增</el-button>
    </div>

    <el-table :data="tableData" border style="width: 100%">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="知识库名称" />
      <el-table-column prop="description" label="描述" show-overflow-tooltip />
      <el-table-column prop="category" label="分类" width="120" />
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'danger'">
            {{ row.status === 1 ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="createTime" label="创建时间" width="180" />
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
          <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination">
      <el-pagination
        v-model:current-page="searchForm.page"
        v-model:page-size="searchForm.pageSize"
        :total="total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>

    <el-dialog
      v-model="dialogVisible"
      :title="dialogType === 'add' ? '新增知识库' : '编辑知识库'"
      width="500px"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="知识库名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入知识库名称" />
        </el-form-item>
        <el-form-item label="知识库描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            placeholder="请输入知识库描述"
          />
        </el-form-item>
        <el-form-item label="知识库分类" prop="category">
          <el-select v-model="form.category" placeholder="请选择分类">
            <el-option
              v-for="item in categoryOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="状态" prop="status" v-if="dialogType === 'edit'">
          <el-switch
            v-model="form.status"
            :active-value="1"
            :inactive-value="0"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getKnowledgeBaseList, createKnowledgeBase, updateKnowledgeBase, deleteKnowledgeBase } from '@/api/knowledge'

const searchForm = reactive({
  page: 1,
  pageSize: 10,
  keyword: '',
  category: ''
})

const tableData = ref([])
const total = ref(0)
const dialogVisible = ref(false)
const dialogType = ref('add')
const formRef = ref(null)

const form = reactive({
  id: null,
  name: '',
  description: '',
  category: '',
  status: 1
})

const rules = {
  name: [
    { required: true, message: '请输入知识库名称', trigger: 'blur' },
    { min: 1, max: 50, message: '长度在 1 到 50 个字符', trigger: 'blur' }
  ],
  description: [
    { required: true, message: '请输入知识库描述', trigger: 'blur' },
    { min: 1, max: 200, message: '长度在 1 到 200 个字符', trigger: 'blur' }
  ],
  category: [
    { required: true, message: '请选择知识库分类', trigger: 'change' }
  ]
}

const categoryOptions = [
  { label: '技术文档', value: '技术文档' },
  { label: '产品文档', value: '产品文档' },
  { label: '用户手册', value: '用户手册' },
  { label: '其他', value: '其他' }
]

const loadData = async () => {
  try {
    const res = await getKnowledgeBaseList(searchForm)
    tableData.value = res.data.list
    total.value = res.data.total
  } catch (error) {
    ElMessage.error('获取知识库列表失败')
  }
}

const handleSearch = () => {
  searchForm.page = 1
  loadData()
}

const handleSizeChange = (val) => {
  searchForm.pageSize = val
  loadData()
}

const handleCurrentChange = (val) => {
  searchForm.page = val
  loadData()
}

const resetForm = () => {
  form.id = null
  form.name = ''
  form.description = ''
  form.category = ''
  form.status = 1
}

const handleAdd = () => {
  dialogType.value = 'add'
  resetForm()
  dialogVisible.value = true
}

const handleEdit = (row) => {
  dialogType.value = 'edit'
  Object.assign(form, row)
  dialogVisible.value = true
}

const handleDelete = (row) => {
  ElMessageBox.confirm('确认删除该知识库吗？', '提示', {
    type: 'warning'
  }).then(async () => {
    try {
      await deleteKnowledgeBase(row.id)
      ElMessage.success('删除成功')
      loadData()
    } catch (error) {
      ElMessage.error('删除失败')
    }
  })
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        if (dialogType.value === 'add') {
          await createKnowledgeBase(form)
          ElMessage.success('创建成功')
        } else {
          await updateKnowledgeBase(form)
          ElMessage.success('更新成功')
        }
        dialogVisible.value = false
        loadData()
      } catch (error) {
        ElMessage.error(dialogType.value === 'add' ? '创建失败' : '更新失败')
      }
    }
  })
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.knowledge-list {
  padding: 20px;
}

.header {
  margin-bottom: 20px;
  display: flex;
  gap: 10px;
}

.search-input {
  width: 200px;
}

.category-select {
  width: 150px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style> 