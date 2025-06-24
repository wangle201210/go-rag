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
          <el-icon><Edit /></el-icon>
          当前知识库：{{ selectedKnowledgeId }}
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
            <el-button type="primary" size="small" @click="saveKnowledgeSelection">确认</el-button>
            <el-button @click="popoverVisible = false" size="small">取消</el-button>
          </el-form-item>
        </el-form>
      </div>
    </el-popover>
  </div>
</template>
<script setup>
import { computed, onMounted, ref } from 'vue'
import { Edit } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import request from '../utils/request'

// 定义事件
const emit = defineEmits(['knowledge-changed', 'change'])

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

// 获取知识库列表
function fetchKnowledgeBaseList() {
  loading.value = true
  return new Promise((resolve) => {
    request.get('/v1/kb')
        .then((response) => {
          knowledgeBaseList.value = response.data?.list || []
          // 如果当前选中的知识库不在列表中，清空选择
          if (selectedKnowledgeId.value) {
            const selected = knowledgeBaseList.value.find(
                item => item.name === selectedKnowledgeId.value && item.status !== 2,
            )
            if (!selected) {
              selectedKnowledgeId.value = ''
              localStorage.removeItem(STORAGE_KEY)
            }
          }
          // 如果没有选中的知识库，且列表中有可用的知识库，则自动选择第一个
          if (!selectedKnowledgeId.value && knowledgeBaseList.value.length > 0) {
            const firstAvailable = knowledgeBaseList.value.find(item => item.status !== 2)
            if (firstAvailable) {
              selectedKnowledgeId.value = firstAvailable.name
              // 自动选择后触发change事件
              emit('change', selectedKnowledgeId.value)
            }
          }
          resolve()
        })
        .catch((error) => {
          ElMessage.error('获取知识库列表失败')
          resolve()
        })
        .finally(() => {
          loading.value = false
        })
  })
}

// 处理知识库选择变化
function handleKnowledgeChange(value) {
  // 知识库选择变化处理
}

// 保存知识库选择
function saveKnowledgeSelection() {
  if (!selectedKnowledgeId.value) {
    ElMessage.warning('请选择一个知识库')
    return
  }

  // 保存到本地存储
  localStorage.setItem(STORAGE_KEY, selectedKnowledgeId.value)

  ElMessage.success(`已选择知识库: ${selectedKnowledgeId.value}`)
  popoverVisible.value = false

  // 触发自定义事件，通知父组件
  emit('knowledge-changed', selectedKnowledgeId.value)
}

// 初始化：获取知识库列表和已保存的选择
onMounted(async () => {
  // 从本地存储获取已保存的知识库ID
  const savedKnowledgeId = localStorage.getItem(STORAGE_KEY)

  if (savedKnowledgeId) {
    selectedKnowledgeId.value = savedKnowledgeId
  }

  // 获取知识库列表
  await fetchKnowledgeBaseList()
})

// 暴露方法给父组件
defineExpose({
  fetchKnowledgeBaseList,
  getSelectedKnowledgeId: () => selectedKnowledgeId.value,
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
  color: #909399;
  line-height: 1.5;
}

.knowledge-info p {
  margin: 5px 0;
}

:deep(.knowledge-select-dropdown) {
  z-index: 9999 !important;
}
</style>