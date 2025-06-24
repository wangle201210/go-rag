<script lang="ts" setup>
import {
  ChatDotRound,
  Files,
  FolderOpened,
  Search,
  Upload,
} from '@element-plus/icons-vue'
import { useRoute } from 'vue-router'

const route = useRoute()

// const isCollapse = ref(true)
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
  >
    <el-menu-item index="/knowledge-base">
      <el-icon>
        <FolderOpened />
      </el-icon>
      <template #title>
        知识库管理
      </template>
    </el-menu-item>
    <el-menu-item index="/indexer">
      <el-icon>
        <Upload />
      </el-icon>
      <template #title>
        索引管理
      </template>
    </el-menu-item>
    <el-menu-item index="/knowledge-documents">
      <el-icon>
        <Files />
      </el-icon>
      <template #title>
        文档管理
      </template>
    </el-menu-item>
    <el-menu-item index="/retriever">
      <el-icon>
        <Search />
      </el-icon>
      <template #title>
        文档检索
      </template>
    </el-menu-item>
    <el-menu-item index="/chat">
      <el-icon>
        <ChatDotRound />
      </el-icon>
      <template #title>
        智能问答
      </template>
    </el-menu-item>
<!--    <el-sub-menu index="demo">-->
<!--      <template #title>-->
<!--        <el-icon>-->
<!--          <Location />-->
<!--        </el-icon>-->
<!--        <span>示例页面</span>-->
<!--      </template>-->
<!--      <el-menu-item index="/nav/2">-->
<!--        <el-icon>-->
<!--          <IconMenu />-->
<!--        </el-icon>-->
<!--        <template #title>-->
<!--          Navigator Two-->
<!--        </template>-->
<!--      </el-menu-item>-->
<!--      <el-menu-item index="/nav/4">-->
<!--        <el-icon>-->
<!--          <Setting />-->
<!--        </el-icon>-->
<!--        <template #title>-->
<!--          Navigator Four-->
<!--        </template>-->
<!--      </el-menu-item>-->
<!--    </el-sub-menu>-->
  </el-menu>
</template>
