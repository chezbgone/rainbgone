import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [tailwindcss(), sveltekit()],
	resolve: {
		alias: {
			$lib: '/src/lib',
			$src: '/src'
		}
	},
	server: {
		watch: {
			ignored: ['../server/**']
		}
	}
});
