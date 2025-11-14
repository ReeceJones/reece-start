import { render, screen, cleanup } from '@testing-library/svelte';
import { expect, describe, it, beforeEach, afterEach } from 'vitest';
import OrganizationNav from './OrganizationNav.svelte';
import type { Organization } from '$lib/schemas/organization';
import { API_TYPES } from '$lib/schemas/api';

describe('OrganizationNav', () => {
	const mockOrganization: Organization = {
		data: {
			id: 'org-123',
			type: API_TYPES.organization,
			attributes: {
				name: 'Test Organization',
				description: 'Test Description',
				address: {
					line1: '123 Test St',
					line2: '',
					city: 'Test City',
					stateOrProvince: 'TS',
					zip: '12345',
					country: 'US'
				},
				locale: 'en',
				contactEmail: 'test@example.com',
				contactPhone: '1234567890',
				contactPhoneCountry: 'US'
			},
			meta: {
				logoDistributionUrl: undefined,
				onboardingStatus: 'completed',
				stripe: {
					onboardingStatus: 'completed'
				}
			}
		}
	};

	const mockOrganizationWithLogo: Organization = {
		...mockOrganization,
		data: {
			...mockOrganization.data,
			meta: {
				...mockOrganization.data.meta,
				logoDistributionUrl: 'https://example.com/logo.png'
			}
		}
	};

	beforeEach(() => {
		// No setup needed
	});

	afterEach(() => {
		cleanup();
	});

	describe('rendering', () => {
		it('should render the navigation menu', () => {
			const { container } = render(OrganizationNav, {
				props: { organization: mockOrganization }
			});
			const menu = container.querySelector('ul.menu');
			expect(menu).toBeTruthy();
		});

		it('should display the organization name', () => {
			render(OrganizationNav, {
				props: { organization: mockOrganization }
			});
			const orgName = screen.getByText('Test Organization');
			expect(orgName).toBeTruthy();
		});

		it('should display default organization text when name is missing', () => {
			const orgWithoutName: Organization = {
				...mockOrganization,
				data: {
					...mockOrganization.data,
					attributes: {
						...mockOrganization.data.attributes,
						name: ''
					}
				}
			};
			render(OrganizationNav, {
				props: { organization: orgWithoutName }
			});
			const orgText = screen.getByText('Organization');
			expect(orgText).toBeTruthy();
		});

		it('should render Building2 icon when logo is not available', () => {
			const { container } = render(OrganizationNav, {
				props: { organization: mockOrganization }
			});
			const icon = container.querySelector('svg');
			expect(icon).toBeTruthy();
		});

		it('should render organization logo when available', () => {
			render(OrganizationNav, {
				props: { organization: mockOrganizationWithLogo }
			});
			const logo = screen.getByAltText('Organization logo');
			expect(logo).toBeTruthy();
			expect(logo.getAttribute('src')).toBe('https://example.com/logo.png');
		});

		it('should display the application menu title', () => {
			render(OrganizationNav, {
				props: { organization: mockOrganization }
			});
			const menuTitle = screen.getByText('Application');
			expect(menuTitle).toBeTruthy();
		});
	});

	describe('navigation links', () => {
		it('should render dashboard link', () => {
			render(OrganizationNav, {
				props: { organization: mockOrganization }
			});
			const dashboardLink = screen.getByRole('link', { name: /dashboard/i });
			expect(dashboardLink).toBeTruthy();
			expect(dashboardLink.getAttribute('href')).toBe('/app/org-123');
		});

		it('should render foo link', () => {
			render(OrganizationNav, {
				props: { organization: mockOrganization }
			});
			const fooLink = screen.getByRole('link', { name: /foo/i });
			expect(fooLink).toBeTruthy();
			expect(fooLink.getAttribute('href')).toBe('/app/org-123/foo');
		});

		it('should render bar link', () => {
			render(OrganizationNav, {
				props: { organization: mockOrganization }
			});
			const barLink = screen.getByRole('link', { name: /bar/i });
			expect(barLink).toBeTruthy();
			expect(barLink.getAttribute('href')).toBe('/app/org-123/bar');
		});

		it('should render settings link', () => {
			render(OrganizationNav, {
				props: { organization: mockOrganization }
			});
			const settingsLinks = screen.getAllByRole('link', { name: /settings/i });
			expect(settingsLinks.length).toBeGreaterThan(0);
			// Check that at least one settings link points to the correct URL
			const mainSettingsLink = settingsLinks.find(
				(link) => link.getAttribute('href') === '/app/org-123/settings'
			);
			expect(mainSettingsLink).toBeTruthy();
		});

		it('should render House icon in dashboard link', () => {
			render(OrganizationNav, {
				props: { organization: mockOrganization }
			});
			const dashboardLink = screen.getByRole('link', { name: /dashboard/i });
			const icon = dashboardLink.querySelector('svg');
			expect(icon).toBeTruthy();
		});

		it('should render Folder icon in foo link', () => {
			render(OrganizationNav, {
				props: { organization: mockOrganization }
			});
			const fooLink = screen.getByRole('link', { name: /foo/i });
			const icon = fooLink.querySelector('svg');
			expect(icon).toBeTruthy();
		});

		it('should render DollarSign icon in bar link', () => {
			render(OrganizationNav, {
				props: { organization: mockOrganization }
			});
			const barLink = screen.getByRole('link', { name: /bar/i });
			const icon = barLink.querySelector('svg');
			expect(icon).toBeTruthy();
		});
	});

	describe('organization dropdown', () => {
		it('should render organization dropdown button', () => {
			const { container } = render(OrganizationNav, {
				props: { organization: mockOrganization }
			});
			const dropdown = container.querySelector('.dropdown');
			expect(dropdown).toBeTruthy();
		});

		it('should render settings link in dropdown', () => {
			render(OrganizationNav, {
				props: { organization: mockOrganization }
			});
			const settingsLinks = screen.getAllByRole('link', { name: /settings/i });
			expect(settingsLinks.length).toBeGreaterThan(0);
		});

		it('should render switch organization link in dropdown', () => {
			render(OrganizationNav, {
				props: { organization: mockOrganization }
			});
			const switchLink = screen.getByRole('link', { name: /switch organization/i });
			expect(switchLink).toBeTruthy();
			expect(switchLink.getAttribute('href')).toBe('/app');
		});
	});
});
