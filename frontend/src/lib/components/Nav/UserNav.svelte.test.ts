import { render, screen, cleanup } from '@testing-library/svelte';
import { expect, describe, it, beforeEach, afterEach, vi } from 'vitest';
import UserNav from './UserNav.svelte';
import { getSelfUserResponseSchema } from '$lib/schemas/user';
import type { z } from 'zod';
import { API_TYPES } from '$lib/schemas/api';
import { setIsImpersonatingUser } from '$lib/auth';

// Mock $app/forms
vi.mock('$app/forms', async () => {
	const { createMockEnhance } = await import('$lib/test-utils');
	return createMockEnhance();
});

describe('UserNav', () => {
	const mockUser: z.infer<typeof getSelfUserResponseSchema> = {
		data: {
			id: 'user-123',
			type: API_TYPES.user,
			attributes: {
				name: 'Test User',
				email: 'test@example.com'
			},
			meta: {
				logoDistributionUrl: undefined,
				tokenRevocation: undefined
			}
		}
	};

	const mockUserWithLogo: z.infer<typeof getSelfUserResponseSchema> = {
		...mockUser,
		data: {
			...mockUser.data,
			meta: {
				...mockUser.data.meta,
				logoDistributionUrl: 'https://example.com/user-logo.png'
			}
		}
	};

	beforeEach(() => {
		setIsImpersonatingUser(false);
	});

	afterEach(() => {
		cleanup();
		vi.clearAllMocks();
	});

	describe('rendering', () => {
		it('should render the navigation menu', () => {
			const { container } = render(UserNav, {
				props: { user: mockUser }
			});
			const menu = container.querySelector('ul.menu');
			expect(menu).toBeTruthy();
		});

		it('should display the user name', () => {
			render(UserNav, {
				props: { user: mockUser }
			});
			const userName = screen.getByText('Test User');
			expect(userName).toBeTruthy();
		});

		it('should display default profile title when name is missing', () => {
			const userWithoutName: z.infer<typeof getSelfUserResponseSchema> = {
				...mockUser,
				data: {
					...mockUser.data,
					attributes: {
						...mockUser.data.attributes,
						name: ''
					}
				}
			};
			render(UserNav, {
				props: { user: userWithoutName }
			});
			const profileText = screen.getByText('Profile');
			expect(profileText).toBeTruthy();
		});

		it('should render User icon when logo is not available', () => {
			const { container } = render(UserNav, {
				props: { user: mockUser }
			});
			const icon = container.querySelector('svg');
			expect(icon).toBeTruthy();
		});

		it('should render user logo when available', () => {
			render(UserNav, {
				props: { user: mockUserWithLogo }
			});
			const logo = screen.getByAltText('User logo');
			expect(logo).toBeTruthy();
			expect(logo.getAttribute('src')).toBe('https://example.com/user-logo.png');
		});
	});

	describe('user dropdown', () => {
		it('should render user dropdown button', () => {
			const { container } = render(UserNav, {
				props: { user: mockUser }
			});
			const dropdown = container.querySelector('.dropdown');
			expect(dropdown).toBeTruthy();
		});

		it('should render settings link in dropdown', () => {
			render(UserNav, {
				props: { user: mockUser }
			});
			const settingsLink = screen.getByRole('link', { name: /settings/i });
			expect(settingsLink).toBeTruthy();
		});

		it('should render logout button in dropdown', () => {
			render(UserNav, {
				props: { user: mockUser }
			});
			const logoutButton = screen.getByRole('button', { name: /logout/i });
			expect(logoutButton).toBeTruthy();
		});

		it('should render Settings icon in settings link', () => {
			render(UserNav, {
				props: { user: mockUser }
			});
			const settingsLink = screen.getByRole('link', { name: /settings/i });
			const icon = settingsLink.querySelector('svg');
			expect(icon).toBeTruthy();
		});

		it('should render LogOut icon in logout button', () => {
			render(UserNav, {
				props: { user: mockUser }
			});
			const logoutButton = screen.getByRole('button', { name: /logout/i });
			const icon = logoutButton.querySelector('svg');
			expect(icon).toBeTruthy();
		});
	});

	describe('forms', () => {
		it('should render signout form', () => {
			const { container } = render(UserNav, {
				props: { user: mockUser }
			});
			const signoutForm = container.querySelector('#signout-form');
			expect(signoutForm).toBeTruthy();
			expect(signoutForm?.getAttribute('action')).toBe('/app?/signout');
			expect(signoutForm?.getAttribute('method')).toBe('POST');
		});

		it('should render stop impersonation form when impersonating', () => {
			setIsImpersonatingUser(true);
			const { container } = render(UserNav, {
				props: { user: mockUser }
			});
			const stopImpersonationForm = container.querySelector('#stop-impersonation-form');
			expect(stopImpersonationForm).toBeTruthy();
			expect(stopImpersonationForm?.getAttribute('action')).toBe('/app?/stopImpersonation');
			expect(stopImpersonationForm?.getAttribute('method')).toBe('POST');
		});

		it('should not render stop impersonation form when not impersonating', () => {
			setIsImpersonatingUser(false);
			const { container } = render(UserNav, {
				props: { user: mockUser }
			});
			const stopImpersonationForm = container.querySelector('#stop-impersonation-form');
			// Form should still exist but button should not be visible
			expect(stopImpersonationForm).toBeTruthy();
		});
	});

	describe('impersonation', () => {
		it('should not show stop impersonation button when not impersonating', () => {
			setIsImpersonatingUser(false);
			render(UserNav, {
				props: { user: mockUser }
			});
			const stopImpersonationButton = screen.queryByRole('button', {
				name: /stop impersonation/i
			});
			expect(stopImpersonationButton).toBeNull();
		});

		it('should show stop impersonation button when impersonating', () => {
			setIsImpersonatingUser(true);
			render(UserNav, {
				props: { user: mockUser }
			});
			const stopImpersonationButton = screen.getByRole('button', {
				name: /stop impersonation/i
			});
			expect(stopImpersonationButton).toBeTruthy();
		});

		it('should render EyeOff icon in stop impersonation button', () => {
			setIsImpersonatingUser(true);
			render(UserNav, {
				props: { user: mockUser }
			});
			const stopImpersonationButton = screen.getByRole('button', {
				name: /stop impersonation/i
			});
			const icon = stopImpersonationButton.querySelector('svg');
			expect(icon).toBeTruthy();
		});
	});

	describe('profile link', () => {
		it('should link to /app/profile when no organizationId in params', () => {
			// The default mock in setupTest.ts sets params to {}
			render(UserNav, {
				props: { user: mockUser }
			});
			const settingsLink = screen.getByRole('link', { name: /settings/i });
			expect(settingsLink.getAttribute('href')).toBe('/app/profile');
		});
	});
});
