import adapter from '@sveltejs/adapter-auto';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
  // Consult https://svelte.dev/docs/kit/integrations
  // for more information about preprocessors
  preprocess: vitePreprocess(),

  kit: {
    adapter: adapter(),
    files: {
      assets: 'frontend/static',
      hooks: {
        client: 'frontend/src/hooks.client',
        server: 'frontend/src/hooks.server',
        universal: 'frontend/src/hooks',
      },
      lib: 'frontend/src/lib',
      params: 'frontend/src/params',
      routes: 'frontend/src/routes',
      serviceWorker: 'frontend/src/service-worker',
      appTemplate: 'frontend/src/app.html',
      errorTemplate: 'frontend/src/error.svelte',
    }
  }
};

export default config;
