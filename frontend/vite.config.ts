import { playwright } from '@vitest/browser-playwright';
import devtoolsJson from 'vite-plugin-devtools-json';
import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { configDefaults, defineConfig } from 'vitest/config';
import { svelteTesting } from '@testing-library/svelte/vite';

export default defineConfig({
	plugins: [tailwindcss(), sveltekit(), devtoolsJson(), svelteTesting()],
	define: {
		// Eliminate in-source test code
		'import.meta.vitest': 'undefined'
	},
	test: {
		// jest like globals
		globals: true,
		// in-source testing
		includeSource: ['src/**/*.{js,ts,svelte}'],
		// Add @testing-library/jest-dom matchers & mocks of SvelteKit modules
		coverage: { exclude: ['setupTest.ts'] },
		projects: [
			{
				extends: './vite.config.ts',
				test: {
					name: 'client',
					environment: 'jsdom',
					setupFiles: ['./setupTest.ts'],
					include: ['src/**/*.svelte.{test,spec}.{js,ts}'],
					exclude: [...configDefaults.exclude, 'src/lib/server/**']
				}
			},
			{
				extends: './vite.config.ts',
				test: {
					name: 'server',
					environment: 'node',
					include: ['src/**/*.{test,spec}.{js,ts}'],
					exclude: [...configDefaults.exclude, 'src/**/*.svelte.{test,spec}.{js,ts}']
				}
			}
		]
	},
	resolve: process.env.VITEST ? { conditions: ['browser'] } : undefined
});
