import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    component: () => import('@/views/layout/Layout.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: { title: '仪表盘', icon: 'dashboard' }
      },
      {
        path: 'knowledge-base',
        name: 'KnowledgeBase',
        component: () => import('@/views/knowledge-base/List.vue'),
        meta: { title: '知识库管理', icon: 'folder' }
      },
      {
        path: 'document',
        name: 'Document',
        component: () => import('@/views/document/List.vue'),
        meta: { title: '文档管理', icon: 'document' }
      }
    ]
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue'),
    meta: { title: '登录' }
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/error/404.vue'),
    meta: { title: '404' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  // 设置页面标题
  document.title = to.meta.title ? `${to.meta.title} - ${import.meta.env.VITE_APP_TITLE || 'GO-RAG'}` : (import.meta.env.VITE_APP_TITLE || 'GO-RAG')

  // 获取 token
  const token = localStorage.getItem('token')

  // 如果访问登录页，直接放行
  if (to.path === '/login') {
    if (token) {
      next('/')
    } else {
      next()
    }
    return
  }

  // 如果没有 token，跳转到登录页
  if (!token) {
    next('/login')
    return
  }

  next()
})

export default router