import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  
  // 添加开发服务器代理配置
  server: {
    proxy: {
      // 将所有 /api 开头的请求代理到 8080 端口
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        secure: false,
        // 如果后端接口不包含 /api 前缀，可以通过 rewrite 去掉
        // rewrite: (path) => path.replace(/^\/api/, '')
      }
    }
  }
})
