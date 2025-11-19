import { describe, it, expect, vi, beforeEach } from 'vitest';
import { withPosthog } from './posthog';
import type { PostHog } from 'posthog-node';

// Mock the posthog-node module
vi.mock('posthog-node', () => {
	return {
		default: {
			PostHog: vi.fn()
		},
		PostHog: vi.fn()
	};
});

describe('posthog', () => {
	beforeEach(() => {
		vi.clearAllMocks();
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
			const fn = vi.fn();
			// mock the posthog client
			const mockClient = {
				identify: vi.fn(),
				shutdown: vi.fn().mockResolvedValue(undefined)
			} as unknown as PostHog;

			// Mock the PostHog constructor to return our mock client
			const { default: posthogModule } = await import('posthog-node');
			vi.mocked(posthogModule.PostHog).mockImplementation(function (this: unknown) {
				return mockClient;
			});

			await withPosthog(async (client) => {
				fn(client);
			});

			expect(fn).toHaveBeenCalledWith(mockClient);
			expect(mockClient.shutdown).toHaveBeenCalled();
		});
	});
});
