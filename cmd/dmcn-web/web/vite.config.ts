import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig({
  plugins: [react()],
  build: {
    outDir: 'dist',
    emptyOutDir: true,
  },
  server: {
    proxy: {
      '/api': {
        target: 'https://localhost:8443',
        secure: false,
      },
      '/ws': {
        target: 'wss://localhost:8443',
        secure: false,
        ws: true,
      },
    },
  },
});
