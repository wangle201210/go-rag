import { createRouter, createWebHistory } from 'vue-router'
import IndexerView from '../views/IndexerView.vue'
import RetrieverView from '../views/RetrieverView.vue'
import ChatView from '../views/ChatView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/indexer'
    },
    {
      path: '/indexer',
      name: 'indexer',
      component: IndexerView
    },
    {
      path: '/retriever',
      name: 'retriever',
      component: RetrieverView
    },
    {
      path: '/chat',
      name: 'chat',
      component: ChatView
    }
  ]
})

export default router