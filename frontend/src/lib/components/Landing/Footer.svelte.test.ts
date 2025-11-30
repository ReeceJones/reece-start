import { render, screen, within, cleanup } from '@testing-library/svelte';
import { expect, describe, it, beforeEach, afterEach } from 'vitest';
import Footer from './Footer.svelte';

describe('Footer', () => {
	let container: HTMLElement;

	beforeEach(() => {
		const result = render(Footer);
		container = result.container;
	});

	afterEach(() => {
		cleanup();
	});

	it('should render the footer component', () => {
		const footer = screen.getByRole('contentinfo');
		expect(footer).toBeTruthy();
	});

	it('should display the brand name', () => {
		const brandName = screen.getByText('reece-start');
		expect(brandName).toBeTruthy();
	});

	it('should display the footer description', () => {
		const description = screen.getByText(/Production-ready SvelteKit \+ Go starter template/);
		expect(description).toBeTruthy();
	});

	it('should display the copyright text', () => {
		const copyright = screen.getByText(/Copyright © 2025 - All rights reserved/);
		expect(copyright).toBeTruthy();
	});

	it('should render the pricing link with correct href', () => {
		const pricingLink = screen.getByRole('link', { name: /Pricing/i });
		expect(pricingLink).toBeTruthy();
		expect(pricingLink.getAttribute('href')).toBe('/pricing');
	});

	it('should render the FAQ link with correct href', () => {
		const faqLink = screen.getByRole('link', { name: /FAQ/i });
		expect(faqLink).toBeTruthy();
		expect(faqLink.getAttribute('href')).toBe('/faq');
	});

	it('should render the GitHub link with correct href', () => {
		const githubLink = screen.getByRole('link', { name: /GitHub/i });
		expect(githubLink).toBeTruthy();
		expect(githubLink.getAttribute('href')).toBe('https://github.com/reecejones/reece-start');
	});

	it('should have all navigation links', () => {
		const footer = screen.getByRole('contentinfo');
		const links = within(footer).getAllByRole('link');
		expect(links).toHaveLength(3);
	});

	it('should have correct CSS classes for styling', () => {
		const footer = container.querySelector('footer');
		expect(footer).toBeTruthy();
		expect(footer?.className).toContain('footer');
		expect(footer?.className).toContain('bg-base-200');
		expect(footer?.className).toContain('p-16');
		expect(footer?.className).toContain('text-base-content');
	});
});
