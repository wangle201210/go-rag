<script lang="ts" setup>
import { getCurrentInstance, ref } from 'vue'
import { useRouter } from 'vue-router'

import * as ElementPlusIconsVue from '@element-plus/icons-vue'
const isCollapse = ref(false)

const router = useRouter()
const route = useRouter().currentRoute

const menuRoutes = router.getRoutes().filter(route => 
  route.meta?.showInMenu,
)

// 注册所有图标组件
const app = getCurrentInstance()
if (app) {
  for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.appContext.app.component(key, component)
  }
}

function handleOpen(key: string, keyPath: string[]) {
  // eslint-disable-next-line no-console
  console.log(key, keyPath)
}
function handleClose(key: string, keyPath: string[]) {
  // eslint-disable-next-line no-console
  console.log(key, keyPath)
}
</script>

<template>
  <el-menu
    router
    :default-active="route.path"
    class="el-menu-vertical-demo"
    @open="handleOpen"
    @close="handleClose"
    :collapse="isCollapse"
  >
    <el-menu-item
      v-for="route in menuRoutes"
      :key="route.path"
      :index="route.path"
    >
      <el-icon>
        <component :is="ElementPlusIconsVue[route.meta?.icon as string]" />
      </el-icon>
      <template #title>
        {{ route.meta?.title }}
      </template>
    </el-menu-item>
  </el-menu>
</template>
