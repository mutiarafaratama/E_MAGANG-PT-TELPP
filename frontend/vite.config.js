import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'
import path from 'path'

const port = Number(process.env.PORT ?? 5174)
const basePath = process.env.BASE_PATH ?? '/'
const backendPort = process.env.BACKEND_PORT ?? '8080'
const isProd = process.env.NODE_ENV === 'production'

export default defineConfig({
  base: basePath,
  plugins: [vue(), tailwindcss()],
  resolve: {
    alias: {
      '@': path.resolve(import.meta.dirname, 'src'),
    },
  },
  root: path.resolve(import.meta.dirname),
  build: {
    outDir: path.resolve(import.meta.dirname, 'dist/public'),
    emptyOutDir: true,
    // Matikan source map di production agar kode Vue tidak bisa dibaca publik
    sourcemap: false,
    // Minifikasi + obfuscation nama variabel
    minify: 'esbuild',
  },
  server: {
    port,
    strictPort: true,
    host: '0.0.0.0',
    allowedHosts: true,
    headers: isProd
      ? {
          'X-Frame-Options': 'DENY',
          'X-Content-Type-Options': 'nosniff',
          'Referrer-Policy': 'strict-origin-when-cross-origin',
        }
      : {
          // Dev: longgar agar bisa di-embed di Replit preview iframe
          'X-Frame-Options': 'ALLOWALL',
          'Content-Security-Policy': 'frame-ancestors *',
        },
    proxy: {
      '/api': {
        target: `http://localhost:${backendPort}`,
        changeOrigin: true,
        ws: true,
      },
      // /uploads tidak lagi di-proxy ke backend secara langsung —
      // dokumen pribadi diakses via /api/dokumen/:id/download (sudah ada auth)
      // Landing page images (/uploads/hero*, /uploads/alur*, dst.) tetap aman
      // karena backend hanya serve folder non-UUID via handler khusus
      '/uploads': {
        target: `http://localhost:${backendPort}`,
        changeOrigin: true,
      },
      '/__mockup': {
        target: 'http://localhost:8000',
        changeOrigin: true,
      },
    },
  },
  preview: {
    port,
    host: '0.0.0.0',
    allowedHosts: true,
  },
})
