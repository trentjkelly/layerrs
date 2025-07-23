import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [sveltekit()],
	server: {
		port: 3000,
		host: true,
		strictPort: true,
		allowedHosts: ['localhost', '127.0.0.1', '0.0.0.0']
	}
});
