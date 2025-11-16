import adapter from '@sveltejs/adapter-node';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	// Consult https://svelte.dev/docs/kit/integrations
	// for more information about preprocessors
	preprocess: vitePreprocess(),
	kit: {
		adapter: adapter(),
		// fix reverse proxy issues with csrf (TODO: investigate if this is still needed)
		csrf: {
			// This is deprecated. Use trustedOrigins instead.
			checkOrigin: false
		}
	}
};

export default config;
