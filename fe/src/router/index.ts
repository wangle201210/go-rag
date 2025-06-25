import type { RouteRecordRaw } from 'vue-router'
import { createRouter, createWebHashHistory } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/knowledge-base',
  },
  {
    path: '/knowledge-base',
    name: 'KnowledgeBase',
    component: () => import('~/pages/rag/knowledge-base.vue'),
    meta: {
      title: '知识库管理',
      icon: 'FolderOpened',
      showInMenu: true,
    },
  },
  {
    path: '/indexer',
    name: 'Indexer',
    component: () => import('~/pages/rag/indexer.vue'),
    meta: {
      title: '文档索引',
      icon: 'Upload',
      showInMenu: true,
    },
  },
  {
    path: '/knowledge-documents',
    name: 'KnowledgeDocuments',
    component: () => import('~/pages/rag/knowledge-documents.vue'),
    meta: {
      title: '文档管理',
      icon: 'Files',
      showInMenu: true,
    },
  },
  {
    path: '/retriever',
    name: 'Retriever',
    component: () => import('~/pages/rag/retriever.vue'),
    meta: {
      title: '文档检索',
      icon: 'Search',
      showInMenu: true,
    },
  },
  {
    path: '/chat',
    name: 'Chat',
    component: () => import('~/pages/rag/chat.vue'),
    meta: {
      title: '智能问答',
      icon: 'ChatDotRound',
      showInMenu: true,
    },
  },
  {
    path: '/chunk-details/:documentId',
    name: 'ChunkDetails',
    component: () => import('~/pages/rag/chunk-details/[documentId].vue'),
  },
]

const router = createRouter({
  history: createWebHashHistory(),
  routes,
})

export default router