import { describe, it, expect, vi, beforeEach } from 'vitest';
import { withPosthog } from './posthog';

// Create a shared mock client that persists across tests using vi.hoisted
const { mockClient, mockEnv } = vi.hoisted(() => {
	const mockClient = {
		identify: vi.fn(),
		shutdown: vi.fn().mockResolvedValue(undefined)
	};
	const mockEnv = {
		PUBLIC_POSTHOG_KEY: undefined as string | undefined,
		PUBLIC_POSTHOG_HOST: undefined as string | undefined
	};
	return { mockClient, mockEnv };
});

// Mock the posthog-node module
vi.mock('posthog-node', () => {
	// Create a proper constructor function
	function PostHogConstructor() {
		return mockClient;
	}

	return {
		default: {
			PostHog: PostHogConstructor
		},
		PostHog: PostHogConstructor
	};
});

// Mock $env/dynamic/public
vi.mock('$env/dynamic/public', () => ({
	env: mockEnv
}));

describe('posthog', () => {
	beforeEach(() => {
		// Reset call counts but preserve the methods
		vi.mocked(mockClient.identify).mockClear();
		vi.mocked(mockClient.shutdown).mockClear();
		vi.mocked(mockClient.shutdown).mockResolvedValue(undefined);
	});

	describe('withPosthog', () => {
		beforeEach(() => {
			// Reset env values before each test
			mockEnv.PUBLIC_POSTHOG_KEY = undefined;
			mockEnv.PUBLIC_POSTHOG_HOST = undefined;
		});

		it('should not log posthog events in local or dev environment', async () => {
			process.env.NODE_ENV = 'development';
			const fn = vi.fn();

			await withPosthog(fn);

			expect(fn).not.toHaveBeenCalled();
		});

		it('should log posthog events in production environment', async () => {
			process.env.NODE_ENV = 'production';
			mockEnv.PUBLIC_POSTHOG_KEY = 'phc_123';
			const fn = vi.fn();

			await withPosthog(async (client) => {
				fn(client);
			});

			expect(fn).toHaveBeenCalledWith(mockClient);
			expect(mockClient.shutdown).toHaveBeenCalled();
		});
	});
});
