import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'
import { resolve } from 'node:path'
import seoPrerender from 'vite-plugin-seo-prerender'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    vueDevTools(),
    seoPrerender({
      renderTarget: '#app',
      routes: ['/', '/user', '/register', '/register_result', '/login', "/upload"],
    }),
  ],
  build: {
    rollupOptions: {
      input: {
        index: resolve(__dirname, 'index.html'),
        user: resolve(__dirname, 'user.html'),
        register: resolve(__dirname, 'register.html'),
        register_result: resolve(__dirname, 'register_result.html'),
        login: resolve(__dirname, 'login.html'),
        upload: resolve(__dirname, 'upload.html'),
      }
    },
  },
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
  }
})
