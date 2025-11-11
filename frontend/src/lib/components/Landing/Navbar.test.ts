import { render, screen, cleanup, waitFor } from '@testing-library/svelte';
import { expect, describe, it, beforeEach, afterEach } from 'vitest';
import userEvent from '@testing-library/user-event';
import Navbar from './Navbar.svelte';

describe('Navbar', () => {
	afterEach(() => {
		cleanup();
	});

	describe('when user is not logged in', () => {
		let container: HTMLElement;

		beforeEach(() => {
			const result = render(Navbar, { props: { isLoggedIn: false } });
			container = result.container;
		});

		it('should render the navbar component', () => {
			const navbar = container.querySelector('.navbar');
			expect(navbar).toBeTruthy();
		});

		it('should display the brand name', () => {
			const brandName = screen.getByText('reece-start');
			expect(brandName).toBeTruthy();
		});

		it('should render the Rocket icon', () => {
			const icon = container.querySelector('svg');
			expect(icon).toBeTruthy();
		});

		it('should render brand link with correct href', () => {
			const brandLink = screen.getByRole('link', { name: /reece-start/i });
			expect(brandLink).toBeTruthy();
			expect(brandLink.getAttribute('href')).toBe('/');
		});

		describe('desktop navigation', () => {
			it('should render FAQ link', () => {
				const faqLink = screen.getByRole('link', { name: /FAQ/i });
				expect(faqLink).toBeTruthy();
				expect(faqLink.getAttribute('href')).toBe('/faq');
			});

			it('should render Pricing link', () => {
				const pricingLink = screen.getByRole('link', { name: /Pricing/i });
				expect(pricingLink).toBeTruthy();
				expect(pricingLink.getAttribute('href')).toBe('/pricing');
			});

			it('should render Sign In link', () => {
				const signInLink = screen.getByRole('link', { name: /Sign In/i });
				expect(signInLink).toBeTruthy();
				expect(signInLink.getAttribute('href')).toBe('/signin');
			});

			it('should render Get Started link', () => {
				const getStartedLink = screen.getByRole('link', { name: /Get Started/i });
				expect(getStartedLink).toBeTruthy();
				expect(getStartedLink.getAttribute('href')).toBe('/signup');
			});

			it('should not render Dashboard link', () => {
				const dashboardLinks = screen.queryAllByRole('link', { name: /Dashboard/i });
				expect(dashboardLinks).toHaveLength(0);
			});
		});

		describe('mobile navigation', () => {
			it('should render mobile menu toggle button', () => {
				const toggleButton = screen.getByRole('button', { name: /Toggle mobile menu/i });
				expect(toggleButton).toBeTruthy();
			});

			it('should not show mobile menu by default', () => {
				// Mobile menu dropdown should not be present initially
				const mobileMenu = container.querySelector('.border-t.border-base-300.bg-base-100');
				expect(mobileMenu).toBeNull();
			});

			it('should have correct CSS classes for desktop visibility', () => {
				// Desktop navigation should have lg:flex classes
				const desktopNav = container.querySelector('.hidden.lg\\:flex');
				expect(desktopNav).toBeTruthy();
			});

			it('should have hamburger button with lg:hidden class', () => {
				// Mobile toggle button container should be hidden on desktop
				const mobileButtonContainer = container.querySelector('.lg\\:hidden');
				expect(mobileButtonContainer).toBeTruthy();
			});

			it('should render Menu icon initially', () => {
				const toggleButton = screen.getByRole('button', { name: /Toggle mobile menu/i });
				const svgs = toggleButton.querySelectorAll('svg');
				// Should have exactly one SVG (the Menu icon)
				expect(svgs.length).toBeGreaterThan(0);
			});

			it('should open mobile menu when toggle button is clicked', async () => {
				const user = userEvent.setup();
				const toggleButton = screen.getByRole('button', { name: /Toggle mobile menu/i });

				await user.click(toggleButton);

				// After opening, the mobile menu dropdown should be present
				await waitFor(() => {
					const mobileMenu = container.querySelector('.border-t.border-base-300.bg-base-100');
					expect(mobileMenu).toBeTruthy();
				});
			});

			it('should close mobile menu when toggle button is clicked again', async () => {
				const user = userEvent.setup();
				const toggleButton = screen.getByRole('button', { name: /Toggle mobile menu/i });

				// Open the menu
				await user.click(toggleButton);

				await waitFor(() => {
					const mobileMenu = container.querySelector('.border-t.border-base-300.bg-base-100');
					expect(mobileMenu).toBeTruthy();
				});

				// Close the menu
				await user.click(toggleButton);

				await waitFor(() => {
					const mobileMenu = container.querySelector('.border-t.border-base-300.bg-base-100');
					expect(mobileMenu).toBeNull();
				});
			});

			it('should show Sign In and Get Started links in mobile menu when opened', async () => {
				const user = userEvent.setup();
				const toggleButton = screen.getByRole('button', { name: /Toggle mobile menu/i });

				await user.click(toggleButton);

				await waitFor(() => {
					const mobileMenu = container.querySelector('.border-t.border-base-300.bg-base-100');
					expect(mobileMenu).toBeTruthy();

					if (mobileMenu) {
						const signInLink = mobileMenu.querySelector('a[href="/signin"]');
						const getStartedLink = mobileMenu.querySelector('a[href="/signup"]');
						expect(signInLink).toBeTruthy();
						expect(getStartedLink).toBeTruthy();
					}
				});
			});

			it('should close mobile menu when a link is clicked', async () => {
				const user = userEvent.setup();
				const toggleButton = screen.getByRole('button', { name: /Toggle mobile menu/i });

				// Open the menu
				await user.click(toggleButton);

				await waitFor(() => {
					const mobileMenu = container.querySelector('.border-t.border-base-300.bg-base-100');
					expect(mobileMenu).toBeTruthy();
				});

				// Click a link in the mobile menu
				const mobileMenu = container.querySelector('.border-t.border-base-300.bg-base-100');
				const faqLink = mobileMenu?.querySelector('a[href="/faq"]');
				expect(faqLink).toBeTruthy();

				if (faqLink) {
					await user.click(faqLink);
				}

				// Menu should close
				await waitFor(() => {
					const closedMenu = container.querySelector('.border-t.border-base-300.bg-base-100');
					expect(closedMenu).toBeNull();
				});
			});

			it('should toggle icon from Menu to X when opened', async () => {
				const user = userEvent.setup();
				const toggleButton = screen.getByRole('button', { name: /Toggle mobile menu/i });

				// Initially should have Menu icon (with path d="M4 12h16")
				let svg = toggleButton.querySelector('svg');
				let menuPath = svg?.querySelector('path[d="M4 12h16"]');
				expect(menuPath).toBeTruthy();

				// Click to open
				await user.click(toggleButton);

				await waitFor(() => {
					// Should now have X icon (with path d="M18 6 6 18")
					const updatedSvg = toggleButton.querySelector('svg');
					const xPath = updatedSvg?.querySelector('path[d="M18 6 6 18"]');
					expect(xPath).toBeTruthy();
				});
			});
		});
	});

	describe('when user is logged in', () => {
		let container: HTMLElement;

		beforeEach(() => {
			const result = render(Navbar, { props: { isLoggedIn: true } });
			container = result.container;
		});

		it('should render the navbar component', () => {
			const navbar = container.querySelector('.navbar');
			expect(navbar).toBeTruthy();
		});

		describe('desktop navigation', () => {
			it('should render FAQ link', () => {
				const faqLink = screen.getByRole('link', { name: /FAQ/i });
				expect(faqLink).toBeTruthy();
				expect(faqLink.getAttribute('href')).toBe('/faq');
			});

			it('should render Pricing link', () => {
				const pricingLink = screen.getByRole('link', { name: /Pricing/i });
				expect(pricingLink).toBeTruthy();
				expect(pricingLink.getAttribute('href')).toBe('/pricing');
			});

			it('should render Dashboard link', () => {
				const dashboardLink = screen.getByRole('link', { name: /Dashboard/i });
				expect(dashboardLink).toBeTruthy();
				expect(dashboardLink.getAttribute('href')).toBe('/app');
			});

			it('should not render Sign In link', () => {
				const signInLinks = screen.queryAllByRole('link', { name: /Sign In/i });
				expect(signInLinks).toHaveLength(0);
			});

			it('should not render Get Started link', () => {
				const getStartedLinks = screen.queryAllByRole('link', { name: /Get Started/i });
				expect(getStartedLinks).toHaveLength(0);
			});

			it('should render DoorOpen icon in Dashboard link', () => {
				// Dashboard link should contain an SVG icon
				const dashboardLink = screen.getByRole('link', { name: /Dashboard/i });
				const icon = dashboardLink.querySelector('svg');
				expect(icon).toBeTruthy();
			});
		});

		describe('mobile navigation', () => {
			it('should render mobile menu toggle button', () => {
				const toggleButton = screen.getByRole('button', { name: /Toggle mobile menu/i });
				expect(toggleButton).toBeTruthy();
			});

			it('should not show mobile menu by default', () => {
				// Mobile menu dropdown should not be present initially
				const mobileMenu = container.querySelector('.border-t.border-base-300.bg-base-100');
				expect(mobileMenu).toBeNull();
			});

			it('should have correct CSS classes for desktop visibility', () => {
				// Desktop navigation should have lg:flex classes
				const desktopNav = container.querySelector('.hidden.lg\\:flex');
				expect(desktopNav).toBeTruthy();
			});

			it('should have hamburger button with lg:hidden class', () => {
				// Mobile toggle button container should be hidden on desktop
				const mobileButtonContainer = container.querySelector('.lg\\:hidden');
				expect(mobileButtonContainer).toBeTruthy();
			});

			it('should render Menu icon initially', () => {
				const toggleButton = screen.getByRole('button', { name: /Toggle mobile menu/i });
				const svgs = toggleButton.querySelectorAll('svg');
				// Should have exactly one SVG (the Menu icon)
				expect(svgs.length).toBeGreaterThan(0);
			});

			it('should open mobile menu when toggle button is clicked', async () => {
				const user = userEvent.setup();
				const toggleButton = screen.getByRole('button', { name: /Toggle mobile menu/i });

				await user.click(toggleButton);

				// After opening, the mobile menu dropdown should be present
				await waitFor(() => {
					const mobileMenu = container.querySelector('.border-t.border-base-300.bg-base-100');
					expect(mobileMenu).toBeTruthy();
				});
			});

			it('should show Dashboard link in mobile menu when opened', async () => {
				const user = userEvent.setup();
				const toggleButton = screen.getByRole('button', { name: /Toggle mobile menu/i });

				await user.click(toggleButton);

				await waitFor(() => {
					const mobileMenu = container.querySelector('.border-t.border-base-300.bg-base-100');
					expect(mobileMenu).toBeTruthy();

					if (mobileMenu) {
						const dashboardLink = mobileMenu.querySelector('a[href="/app"]');
						expect(dashboardLink).toBeTruthy();
					}
				});
			});

			it('should not show Sign In and Get Started links in mobile menu', async () => {
				const user = userEvent.setup();
				const toggleButton = screen.getByRole('button', { name: /Toggle mobile menu/i });

				await user.click(toggleButton);

				await waitFor(() => {
					const mobileMenu = container.querySelector('.border-t.border-base-300.bg-base-100');
					expect(mobileMenu).toBeTruthy();

					if (mobileMenu) {
						const signInLink = mobileMenu.querySelector('a[href="/signin"]');
						const getStartedLink = mobileMenu.querySelector('a[href="/signup"]');
						expect(signInLink).toBeNull();
						expect(getStartedLink).toBeNull();
					}
				});
			});

			it('should show DoorOpen icon in mobile Dashboard link', async () => {
				const user = userEvent.setup();
				const toggleButton = screen.getByRole('button', { name: /Toggle mobile menu/i });

				await user.click(toggleButton);

				await waitFor(() => {
					const mobileMenu = container.querySelector('.border-t.border-base-300.bg-base-100');
					expect(mobileMenu).toBeTruthy();

					if (mobileMenu) {
						const dashboardLink = mobileMenu.querySelector('a[href="/app"]');
						expect(dashboardLink).toBeTruthy();
						const icon = dashboardLink?.querySelector('svg');
						expect(icon).toBeTruthy();
					}
				});
			});

			it('should close mobile menu when toggle button is clicked again', async () => {
				const user = userEvent.setup();
				const toggleButton = screen.getByRole('button', { name: /Toggle mobile menu/i });

				// Open the menu
				await user.click(toggleButton);

				await waitFor(() => {
					const mobileMenu = container.querySelector('.border-t.border-base-300.bg-base-100');
					expect(mobileMenu).toBeTruthy();
				});

				// Close the menu
				await user.click(toggleButton);

				await waitFor(() => {
					const mobileMenu = container.querySelector('.border-t.border-base-300.bg-base-100');
					expect(mobileMenu).toBeNull();
				});
			});
		});
	});

	describe('mobile menu toggle button', () => {
		it('should render toggle button with correct aria-label', () => {
			render(Navbar, { props: { isLoggedIn: false } });
			const toggleButton = screen.getByRole('button', { name: /Toggle mobile menu/i });
			expect(toggleButton).toBeTruthy();
			expect(toggleButton.getAttribute('aria-label')).toBe('Toggle mobile menu');
		});

		it('should have an icon in the toggle button', () => {
			render(Navbar, { props: { isLoggedIn: false } });
			const toggleButton = screen.getByRole('button', { name: /Toggle mobile menu/i });
			const icon = toggleButton.querySelector('svg');
			expect(icon).toBeTruthy();
		});

		it('should have correct button styling classes', () => {
			render(Navbar, { props: { isLoggedIn: false } });
			const toggleButton = screen.getByRole('button', { name: /Toggle mobile menu/i });
			expect(toggleButton.className).toContain('btn');
			expect(toggleButton.className).toContain('btn-ghost');
			expect(toggleButton.className).toContain('btn-sm');
		});
	});
});
