/**
 * 知识库名称管理服务
 * 用于在不同页面之间共享knowledge_name参数
 * 使用localStorage存储，确保页面刷新后仍然可用
 */

import { v4 as uuidv4 } from 'uuid'

const STORAGE_KEY = 'go_rag_knowledge_name'

/**
 * 获取当前知识库名称
 * 如果不存在，则自动生成一个新的
 */
export function getKnowledgeName() {
  let knowledgeName = localStorage.getItem(STORAGE_KEY)
  
  // 如果不存在，则生成一个新的知识库名称
  if (!knowledgeName) {
    knowledgeName = generateKnowledgeName()
    localStorage.setItem(STORAGE_KEY, knowledgeName)
  }
  
  return knowledgeName
}

/**
 * 设置知识库名称
 */
export function setKnowledgeName(name) {
  localStorage.setItem(STORAGE_KEY, name)
  return name
}

/**
 * 生成一个新的知识库名称
 * 使用UUID确保唯一性
 */
export function generateKnowledgeName() {
  return `knowledge_${uuidv4().substring(0, 8)}`
}