import { createRouter, createWebHashHistory } from 'vue-router'
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'

const routes = [
  {
    path: '/',
    redirect: '/indexer'
  },
  {
    path: '/indexer',
    name: 'Indexer',
    component: () => import('../views/Indexer.vue'),
    meta: {
      title: '索引管理'
    }
  },
  {
    path: '/retriever',
    name: 'Retriever',
    component: () => import('../views/Retriever.vue'),
    meta: {
      title: '检索测试'
    }
  },
  {
    path: '/chat',
    name: 'Chat',
    component: () => import('../views/Chat.vue'),
    meta: {
      title: '对话'
    }
  },
  {
    path: '/knowledge-base',
    name: 'KnowledgeBase',
    component: () => import('../views/KnowledgeBase.vue'),
    meta: {
      title: '知识库管理'
    }
  },
  {
    path: '/knowledge-documents',
    name: 'KnowledgeDocuments',
    component: () => import('../views/KnowledgeDocuments.vue'),
    meta: {
      title: '文档管理'
    }
  },
  {
    path: '/knowledge-documents/:documentId',
    name: 'ChunkDetails',
    component: () => import('../views/ChunkDetails.vue'),
    props: true,
    meta: {
      title: '文档块详情'
    }
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  // 开始加载进度条
  NProgress.start()
  
  // 设置页面标题
  document.title = to.meta.title ? `${to.meta.title} - Go-RAG` : 'Go-RAG'
  
  // 这里可以添加登录验证等逻辑
  next()
})

router.afterEach(() => {
  // 结束加载进度条
  NProgress.done()
})

// 路由错误处理
router.onError((error) => {
  console.error('路由错误:', error)
  NProgress.done()
})

export default router