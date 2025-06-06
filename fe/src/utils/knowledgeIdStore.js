/**
 * 知识库ID管理服务
 * 用于在不同页面之间共享选中的知识库ID
 * 使用localStorage存储，确保页面刷新后仍然可用
 */

const STORAGE_KEY = 'go_rag_selected_knowledge_id'

/**
 * 获取当前选中的知识库ID
 * 如果不存在，则返回空字符串
 */
export function getSelectedKnowledgeId() {
  return localStorage.getItem(STORAGE_KEY) || ''
}

/**
 * 设置选中的知识库ID
 */
export function setSelectedKnowledgeId(id) {
  localStorage.setItem(STORAGE_KEY, id)
  return id
}

/**
 * 清除选中的知识库ID
 */
export function clearSelectedKnowledgeId() {
  localStorage.removeItem(STORAGE_KEY)
}