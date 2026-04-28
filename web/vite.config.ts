import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    port: 3000,
    proxy: {
      '/job/': process.env.VITE_API_URL || 'http://localhost:5000',
      '/execution/': process.env.VITE_API_URL || 'http://localhost:5000',
      '/health': process.env.VITE_API_URL || 'http://localhost:5000',
    },
  },
  define: {
    'import.meta.env.VITE_API_URL': JSON.stringify(process.env.VITE_API_URL || 'http://localhost:5000'),
  },
})
