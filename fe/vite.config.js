import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

// https://vitejs.dev/config/
export default defineConfig(({ command, mode }) => {
  // 根据当前工作目录中的 `mode` 加载 .env 文件
  const env = loadEnv(mode, process.cwd())
  
  return {
    plugins: [vue()],
    resolve: {
      alias: {
        '@': path.resolve(__dirname, 'src'),
      },
    },
    server: {
      port: env.VITE_DEV_PORT || 5173,
      proxy: {
        '/api': {
          target: env.VITE_DEV_PROXY_TARGET || 'http://localhost:8000',
          changeOrigin: true,
          // rewrite: (path) => path.replace(/^\/api/, '')
        }
      }
    },
    build: {
      // 生产环境构建配置
      outDir: 'dist',
      assetsDir: 'assets',
      // 小于此阈值的资源将被内联为base64编码
      // assetsInlineLimit: 4096,
      // 启用/禁用CSS代码拆分
      // cssCodeSplit: true,
      // 构建后是否生成source map文件
      // sourcemap: false,
      // // 自定义底层的Rollup打包配置
      // rollupOptions: {
      //   output: {
      //     // 用于控制chunks的拆分
      //     manualChunks: {
      //       'element-plus': ['element-plus'],
      //       'vue-vendor': ['vue', 'vue-router', 'pinia']
      //     }
      //   }
      // }
      // 设置最小化混淆
      // minify: 'terser',
      // terserOptions: {
      //   compress: {
      //     // 生产环境时移除console
      //     drop_console: true,
      //     drop_debugger: true
      //   }
      // }
    }
  }
})