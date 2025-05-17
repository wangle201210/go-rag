import { createRouter, createWebHashHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    redirect: '/indexer'
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
  history: createWebHashHistory(),
  routes
})

export default router