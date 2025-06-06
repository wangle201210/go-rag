<template>
  <div class="knowledge-selector">
    <el-popover
      placement="bottom"
      :width="300"
      trigger="click"
      v-model:visible="popoverVisible"
      :close-on-click-outside="false"
    >
      <template #reference>
        <el-button type="info" plain size="small">
          <el-icon><Folder /></el-icon>
          知识库设置
        </el-button>
      </template>
      
      <div class="selector-content">
        <h4>知识库设置</h4>
        <el-form>
          <el-form-item label="选择知识库">
            <el-select 
              v-model="selectedKnowledgeId"
              placeholder="请选择知识库"
              size="small"
              filterable
              :loading="loading"
              @change="handleKnowledgeChange"
              style="width: 100%"
              :popper-append-to-body="false"
              :teleported="false"
              popper-class="knowledge-select-dropdown"
            >
              <el-option
                v-for="item in knowledgeBaseList"
                :key="item.name"
                :label="item.name"
                :value="item.name"
                :disabled="item.status === 2"
              >
                <div class="knowledge-option">
                  <span>{{ item.name }}</span>
                  <el-tag size="small" :type="item.status === 2 ? 'danger' : 'success'">
                    {{ item.status === 2 ? '禁用' : '启用' }}
                  </el-tag>
                </div>
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item v-if="selectedKnowledge" label="知识库信息">
            <div class="knowledge-info">
              <p><strong>描述：</strong>{{ selectedKnowledge.description }}</p>
              <p v-if="selectedKnowledge.category"><strong>分类：</strong>{{ selectedKnowledge.category }}</p>
            </div>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="saveKnowledgeSelection">确认</el-button>
            <el-button @click="popoverVisible = false">取消</el-button>
          </el-form-item>
        </el-form>
      </div>
    </el-popover>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { Folder } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import axios from 'axios'

// 组件状态
const popoverVisible = ref(false)
const knowledgeBaseList = ref([])
const selectedKnowledgeId = ref('')
const loading = ref(false)

// 计算属性：获取当前选中的知识库详细信息
const selectedKnowledge = computed(() => {
  return knowledgeBaseList.value.find(item => item.name === selectedKnowledgeId.value) || null
})

// 本地存储键名
const STORAGE_KEY = 'go_rag_selected_knowledge_id'

// 初始化：获取知识库列表和已保存的选择
onMounted(async () => {
  // 从本地存储获取已保存的知识库ID
  const savedKnowledgeId = localStorage.getItem(STORAGE_KEY)
  console.log("11111:", savedKnowledgeId)

  if (savedKnowledgeId) {
    selectedKnowledgeId.value = savedKnowledgeId
  }
  
  // 获取知识库列表
  await fetchKnowledgeBaseList()
})

// 获取知识库列表
const fetchKnowledgeBaseList = async () => {
  loading.value = true
  try {
    const response = await axios.get('/v1/kb')
    knowledgeBaseList.value = response.data.data.list || []
    console.log("知识库id:", selectedKnowledgeId.value)
    console.log("knowledgeBaseList:", knowledgeBaseList.value)
    // 如果有选中的知识库ID但在列表中不存在或已禁用，则清空选择
    if (selectedKnowledgeId.value) {
      const selected = knowledgeBaseList.value.find(
        item => item.name === selectedKnowledgeId.value && item.status !== 2
      )
      // if (!selected) {
      //   selectedKnowledgeId.value = ''
      //   localStorage.removeItem(STORAGE_KEY)
      // }
    }
    
    // 如果没有选中的知识库，且列表中有可用的知识库，则自动选择第一个
    if (!selectedKnowledgeId.value && knowledgeBaseList.value.length > 0) {
      const firstAvailable = knowledgeBaseList.value.find(item => item.status !== 2)
      if (firstAvailable) {
        selectedKnowledgeId.value = firstAvailable.name
      }
    }
  } catch (error) {
    console.error('获取知识库列表失败:', error)
    ElMessage.error('获取知识库列表失败: ' + (error.response?.data?.message || '未知错误'))
  } finally {
    loading.value = false
  }
}

// 处理知识库选择变更
const handleKnowledgeChange = (value) => {
  selectedKnowledgeId.value = value
  // 注意：这里只更新选择，但不保存到localStorage，保存操作在点击确认按钮时进行
  // 确保 popover 保持打开状态
  setTimeout(() => {
    popoverVisible.value = true
  }, 0)
}

// 保存知识库选择
const saveKnowledgeSelection = () => {
  if (!selectedKnowledgeId.value) {
    ElMessage.warning('请选择一个知识库')
    return
  }
  
  localStorage.setItem(STORAGE_KEY, selectedKnowledgeId.value)
  ElMessage.success('知识库设置已保存')
  popoverVisible.value = false
}

// 对外暴露获取当前选中知识库ID的方法
const getSelectedKnowledgeId = () => {
  return selectedKnowledgeId.value
}

// 导出组件方法供外部使用
defineExpose({
  getSelectedKnowledgeId,
  fetchKnowledgeBaseList
})
</script>

<style scoped>
.knowledge-selector {
  display: inline-block;
}

.selector-content {
  padding: 10px;
}

.selector-content h4 {
  margin-top: 0;
  margin-bottom: 15px;
  color: #606266;
}

.knowledge-option {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.knowledge-info {
  font-size: 12px;
  color: #606266;
  background-color: #f5f7fa;
  padding: 8px;
  border-radius: 4px;
}

.knowledge-info p {
  margin: 5px 0;
}

/* 确保下拉菜单正确显示在 popover 内部 */
:deep(.knowledge-select-dropdown) {
  z-index: 3000 !important; /* 确保下拉菜单在 popover 上方显示 */
}

/* 防止下拉菜单与 popover 的交互冲突 */
:deep(.el-select-dropdown) {
  position: static !important;
  margin-top: 5px !important;
}
</style>