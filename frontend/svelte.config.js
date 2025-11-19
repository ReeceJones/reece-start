import adapter from '@sveltejs/adapter-node';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';
import { env } from 'process';

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
			// checkOrigin: false
			trustedOrigins: env.PUBLIC_ALLOWED_ORIGINS?.split(',') || []
		},
		paths: {
			relative: false // Required for PostHog session replay to work correctly
		}
	}
};

export default config;
