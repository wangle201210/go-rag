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
    component: () => import('~/pages/knowledge-base.vue'),
  },
  {
    path: '/knowledge-documents',
    name: 'KnowledgeDocuments',
    component: () => import('~/pages/knowledge-documents.vue'),
  },
  {
    path: '/indexer',
    name: 'Indexer',
    component: () => import('~/pages/indexer.vue'),
  },
  {
    path: '/chat',
    name: 'Chat',
    component: () => import('~/pages/chat.vue'),
  },
  {
    path: '/chunk-details/:documentId',
    name: 'ChunkDetails',
    component: () => import('~/pages/chunk-details/[documentId].vue'),
  },
  {
    path: '/retriever',
    name: 'Retriever',
    component: () => import('~/pages/retriever.vue'),
  },
]

const router = createRouter({
  history: createWebHashHistory(),
  routes,
})

export default router