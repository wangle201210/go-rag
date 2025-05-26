import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    redirect: '/knowledge'
  },
  {
    path: '/knowledge',
    name: 'Knowledge',
    component: () => import('@/views/knowledge/List.vue'),
    meta: {
      title: '知识库管理'
    }
  },
  {
    path: '/indexer',
    name: 'Indexer',
    component: () => import('../views/Indexer.vue')
  },
  {
    path: '/retriever',
    name: 'Retriever',
    component: () => import('../views/Retriever.vue')
  },
  {
    path: '/chat',
    name: 'Chat',
    component: () => import('../views/Chat.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  document.title = to.meta.title || 'Agentic RAG 系统'
  next()
})

export default router