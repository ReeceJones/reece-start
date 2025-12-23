import { render, screen, cleanup } from '@testing-library/svelte';
import { expect, describe, it, beforeEach, afterEach, vi } from 'vitest';
import IndexNav from './IndexNav.svelte';
import { UserScope } from '$lib/schemas/jwt';
import { setScopes } from '$lib/auth';

describe('IndexNav', () => {
	beforeEach(() => {
		setScopes(() => []);
	});

	afterEach(() => {
		cleanup();
		vi.clearAllMocks();
	});

	describe('rendering', () => {
		it('should render the navigation menu', () => {
			const { container } = render(IndexNav);
			const menu = container.querySelector('ul.menu');
			expect(menu).toBeTruthy();
		});

		it('should display the application menu title', () => {
			render(IndexNav);
			const menuTitle = screen.getByText('Application');
			expect(menuTitle).toBeTruthy();
		});

		it('should render the home link', () => {
			render(IndexNav);
			const homeLink = screen.getByRole('link', { name: /home/i });
			expect(homeLink).toBeTruthy();
			expect(homeLink.getAttribute('href')).toBe('/app');
		});

		it('should render House icon in home link', () => {
			render(IndexNav);
			const homeLink = screen.getByRole('link', { name: /home/i });
			const icon = homeLink.querySelector('svg');
			expect(icon).toBeTruthy();
		});
	});

	describe('admin section', () => {
		it('should not render admin section when user is not admin', () => {
			setScopes(() => []);
			render(IndexNav);
			const adminSection = screen.queryByText('Admin');
			expect(adminSection).toBeNull();
		});

		it('should render admin section when user is admin', () => {
			setScopes(() => [UserScope.Admin]);
			render(IndexNav);
			const adminSection = screen.getByText('Admin');
			expect(adminSection).toBeTruthy();
		});

		it('should render users link in admin section when user is admin', () => {
			setScopes(() => [UserScope.Admin]);
			render(IndexNav);
			const usersLink = screen.getByRole('link', { name: /users/i });
			expect(usersLink).toBeTruthy();
			expect(usersLink.getAttribute('href')).toBe('/app/admin/users');
		});

		it('should render debug link in admin section when user is admin', () => {
			setScopes(() => [UserScope.Admin]);
			render(IndexNav);
			const debugLink = screen.getByRole('link', { name: /debug/i });
			expect(debugLink).toBeTruthy();
			expect(debugLink.getAttribute('href')).toBe('/app/admin/debug');
		});

		it('should render Shield icon in admin section', () => {
			setScopes(() => [UserScope.Admin]);
			render(IndexNav);
			const adminSection = screen.getByText('Admin');
			const summary = adminSection.closest('summary');
			const icon = summary?.querySelector('svg');
			expect(icon).toBeTruthy();
		});

		it('should render Users icon in users link', () => {
			setScopes(() => [UserScope.Admin]);
			render(IndexNav);
			const usersLink = screen.getByRole('link', { name: /users/i });
			const icon = usersLink.querySelector('svg');
			expect(icon).toBeTruthy();
		});

		it('should render Bug icon in debug link', () => {
			setScopes(() => [UserScope.Admin]);
			render(IndexNav);
			const debugLink = screen.getByRole('link', { name: /debug/i });
			const icon = debugLink.querySelector('svg');
			expect(icon).toBeTruthy();
		});
	});

	describe('active state', () => {
		it('should apply active class to home link when on /app route', () => {
			render(IndexNav);
			const homeLink = screen.getByRole('link', { name: /home/i });
			// The default mock in setupTest.ts sets url to 'http://localhost/app'
			// so pathname should be '/app' and active class should be applied
			expect(homeLink.className).toContain('bg-base-300');
		});
	});
});
