// Plugins
import vue from '@vitejs/plugin-vue'
import vuetify, { transformAssetUrls } from 'vite-plugin-vuetify'

// Utilities
import { defineConfig } from 'vite'
import { fileURLToPath, URL } from 'node:url'

export default defineConfig({
  base: '',
  plugins: [
    vue({
      template: { transformAssetUrls },
    }),
    vuetify({
      autoImport: true,
      styles: {
        configFile: 'src/styles/settings.scss',
      },
    })
  ],
  build: {
    manifest: false,
    outDir: 'dist',
    rollupOptions: {
      output: {
        entryFileNames: 'assets/ui-[name]-[hash].js',
        chunkFileNames: 'assets/ui-[name]-[hash].js',
        assetFileNames: 'assets/ui-[name]-[hash][extname]',
      },
    }
  },
  define: { 'process.env': {} },
  optimizeDeps: {
    exclude: ['vue-i18n'],
  },
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
    extensions: ['.js', '.json', '.jsx', '.mjs', '.ts', '.tsx', '.vue'],
  },
  server: {
    port: 3000,
    proxy: {
      '/app/api': {
        target: 'http://localhost:2095',
        changeOrigin: false,
      },
    },
  }
})
