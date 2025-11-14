// setupTest.ts
/* eslint-disable @typescript-eslint/no-empty-function */
// import matchers from '@testing-library/jest-dom/matchers';
import { vi } from 'vitest';
import type { Navigation, Page } from '@sveltejs/kit';
import { readable } from 'svelte/store';
import * as environment from '$app/environment';
import * as navigation from '$app/navigation';
import * as stores from '$app/stores';

// Mock Web Animations API for Svelte transitions
// jsdom doesn't support element.animate() by default
if (typeof Element !== 'undefined' && !Element.prototype.animate) {
	Element.prototype.animate = function () {
		return {
			finished: Promise.resolve(),
			cancel: vi.fn(),
			play: vi.fn(),
			pause: vi.fn(),
			reverse: vi.fn(),
			updatePlaybackRate: vi.fn(),
			addEventListener: vi.fn(),
			removeEventListener: vi.fn(),
			dispatchEvent: vi.fn()
		} as unknown as Animation;
	};
}

// Mock Svelte context API for testing
// This allows setContext/getContext to work outside of component initialization
vi.mock('svelte', async () => {
	const actual = await import('svelte');
	// Create a persistent context map that survives across test runs
	const contextMap = new Map<symbol | string, unknown>();
	return {
		...actual,
		setContext: <T>(key: symbol | string, context: T): T => {
			contextMap.set(key, context);
			return context;
		},
		getContext: <T>(key: symbol | string): T | undefined => {
			return contextMap.get(key) as T | undefined;
		}
	};
});

// Mock SvelteKit runtime module $app/environment
vi.mock('$app/environment', (): typeof environment => ({
	browser: false,
	dev: true,
	building: false,
	version: 'any'
}));

// Mock SvelteKit runtime module $app/navigation
vi.mock('$app/navigation', (): typeof navigation => ({
	afterNavigate: () => {},
	beforeNavigate: () => {},
	disableScrollHandling: () => {},
	goto: () => Promise.resolve(),
	invalidate: () => Promise.resolve(),
	invalidateAll: () => Promise.resolve(),
	onNavigate: () => {},
	preloadData: () =>
		Promise.resolve({
			type: 'loaded' as const,
			status: 200,
			data: {}
		}),
	preloadCode: () => Promise.resolve(),
	pushState: () => {},
	refreshAll: () => Promise.resolve(),
	replaceState: () => {}
}));

// Mock SvelteKit runtime module $app/stores
vi.mock('$app/stores', (): typeof stores => {
	const getStores: typeof stores.getStores = () => {
		const navigating = readable<Navigation | null>(null);
		const page = readable<Page>({
			url: new URL('http://localhost'),
			params: {},
			route: {
				id: null
			},
			status: 200,
			error: null,
			data: {},
			form: undefined,
			state: {}
		});
		const updated = { subscribe: readable(false).subscribe, check: async () => false };

		return { navigating, page, updated };
	};

	const page: typeof stores.page = {
		subscribe(fn) {
			return getStores().page.subscribe(fn);
		}
	};
	const navigating: typeof stores.navigating = {
		subscribe(fn) {
			return getStores().navigating.subscribe(fn);
		}
	};
	const updated: typeof stores.updated = {
		subscribe(fn) {
			return getStores().updated.subscribe(fn);
		},
		check: async () => false
	};

	return {
		getStores,
		navigating,
		page,
		updated
	};
});

// Mock SvelteKit runtime module $app/state (Svelte 5 runes)
vi.mock('$app/state', () => {
	const pageState = {
		url: new URL('http://localhost/app'),
		params: {},
		route: { id: null },
		status: 200,
		error: null,
		data: {},
		form: undefined,
		state: {}
	};

	return {
		page: pageState
	};
});
