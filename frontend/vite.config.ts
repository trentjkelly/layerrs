import { sveltekit } from '@sveltejs/kit/vite';
import basicSsl from '@vitejs/plugin-basic-ssl';
import { defineConfig, type PluginOption } from 'vite';

export default defineConfig({
	plugins: [basicSsl(), sveltekit()] as PluginOption[],
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
