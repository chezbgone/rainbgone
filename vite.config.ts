import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
  plugins: [tailwindcss(), sveltekit()],
  resolve: {
    alias: {
      $lib: '/frontend/src/lib',
      $src: '/frontend/src',
    },
  },
  server: {
    fs: {
      allow: ['./frontend'],
    },
  },
});
