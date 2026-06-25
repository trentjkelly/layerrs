import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig, type PluginOption } from 'vite';

export default defineConfig({
	plugins: [sveltekit()] as PluginOption[],
	server: {
		port: 3000,
		host: true,
		strictPort: true
	},
	optimizeDeps: {
    	exclude: ['fsevents']
  	},
  	build: {
    	rollupOptions: {
      	external: ['fsevents']
    	}
  	}
});
