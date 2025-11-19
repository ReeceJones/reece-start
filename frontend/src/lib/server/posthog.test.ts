import { describe, it, expect, vi, beforeEach } from 'vitest';
import { withPosthog } from './posthog';

// Create a shared mock client that persists across tests using vi.hoisted
const { mockClient } = vi.hoisted(() => {
	const mockClient = {
		identify: vi.fn(),
		shutdown: vi.fn().mockResolvedValue(undefined)
	};
	return { mockClient };
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

describe('posthog', () => {
	beforeEach(() => {
		// Reset call counts but preserve the methods
		vi.mocked(mockClient.identify).mockClear();
		vi.mocked(mockClient.shutdown).mockClear();
		vi.mocked(mockClient.shutdown).mockResolvedValue(undefined);
	});

	describe('withPosthog', () => {
		it('should not log posthog events in local or dev environment', async () => {
			process.env.NODE_ENV = 'development';
			const fn = vi.fn();

			await withPosthog(fn);

			expect(fn).not.toHaveBeenCalled();
		});

		it('should log posthog events in production environment', async () => {
			process.env.NODE_ENV = 'production';
			process.env.PUBLIC_POSTHOG_KEY = 'phc_123';
			const fn = vi.fn();

			await withPosthog(async (client) => {
				fn(client);
			});

			expect(fn).toHaveBeenCalledWith(mockClient);
			expect(mockClient.shutdown).toHaveBeenCalled();
		});
	});
});
