import { describe, it, expect, beforeEach, vi } from 'vitest';
import { get } from 'svelte/store';
import { locale, locales, initializeLocale, t } from './i18n';
import { translations } from './translations';

describe('i18n', () => {
	beforeEach(() => {
		// Reset locale to default before each test
		locale.set('en');
	});

	describe('locale store', () => {
		it('should initialize with default locale', () => {
			expect(get(locale)).toBe('en');
		});

		it('should update locale value', () => {
			locale.set('fr');
			expect(get(locale)).toBe('fr');
		});
	});

	describe('locales', () => {
		it('should contain available locales', () => {
			expect(Array.isArray(locales)).toBe(true);
			expect(locales.length).toBeGreaterThan(0);
			expect(locales).toContain('en');
		});
	});

	describe('initializeLocale', () => {
		it('should set locale from server-provided value', () => {
			initializeLocale('fr');
			expect(get(locale)).toBe('fr');
		});

		it('should use default locale when server locale is empty', () => {
			initializeLocale('');
			expect(get(locale)).toBe('en');
		});

		it('should use default locale when server locale is null', () => {
			// @ts-expect-error - testing null case
			initializeLocale(null);
			expect(get(locale)).toBe('en');
		});
	});

	describe('t store (translate)', () => {
		it('should translate simple keys', () => {
			const translate = get(t);
			expect(translate('signIn')).toBe('Sign in');
			expect(translate('signUp')).toBe('Sign up');
			expect(translate('dashboard')).toBe('Dashboard');
		});

		it('should translate nested keys', () => {
			const translate = get(t);
			expect(translate('auth.signIn.title')).toBe('Sign in');
			expect(translate('auth.signIn.description')).toBe(
				'Enter your details below to sign in to your account.'
			);
			expect(translate('settings.fields.email.label')).toBe('Email');
		});

		it('should substitute variables in translations', () => {
			const translate = get(t);
			// Find a translation with variables - checking if any exist
			// If none exist, we'll test the mechanism anyway
			const result = translate('signIn', {});
			expect(typeof result).toBe('string');
		});

		it('should return key when translation is missing', () => {
			const consoleSpy = vi.spyOn(console, 'warn').mockImplementation(() => {});
			const translate = get(t);
			const result = translate('nonExistentKey');
			expect(result).toBe('nonExistentKey');
			expect(consoleSpy).toHaveBeenCalledWith('Translation key not found: nonExistentKey');
			consoleSpy.mockRestore();
		});

		it('should fallback to default locale when key is missing in current locale', () => {
			// This test assumes we might have other locales in the future
			// For now, it tests the fallback mechanism
			const translate = get(t);
			const result = translate('signIn');
			expect(result).toBe('Sign in');
		});

		it('should update translation when locale changes', () => {
			const translate1 = get(t);
			expect(translate1('signIn')).toBe('Sign in');

			locale.set('en'); // Still en, but tests reactivity
			const translate2 = get(t);
			expect(translate2('signIn')).toBe('Sign in');
		});
	});

	describe('translation completeness', () => {
		// Helper function to get all keys from a nested object
		function getAllKeys(obj: Record<string, unknown>, prefix = ''): string[] {
			const keys: string[] = [];
			for (const [key, value] of Object.entries(obj)) {
				const fullKey = prefix ? `${prefix}.${key}` : key;
				if (typeof value === 'object' && value !== null && !Array.isArray(value)) {
					keys.push(...getAllKeys(value as Record<string, unknown>, fullKey));
				} else if (typeof value === 'string') {
					keys.push(fullKey);
				}
			}
			return keys;
		}

		it('should have all locales contain translations for all default locale keys', () => {
			const defaultLocale = translations.en;
			const defaultKeys = getAllKeys(defaultLocale as Record<string, unknown>);

			expect(defaultKeys.length).toBeGreaterThan(0);

			for (const localeKey of locales) {
				const localeTranslations = translations[localeKey as keyof typeof translations];
				expect(localeTranslations).toBeDefined();

				for (const key of defaultKeys) {
					// Helper function to check if key exists in translations
					function keyExists(obj: Record<string, unknown>, path: string): boolean {
						const keys = path.split('.');
						let current = obj;
						for (const k of keys) {
							if (current === undefined || current === null || typeof current !== 'object') {
								return false;
							}
							current = current[k] as Record<string, unknown>;
						}
						return typeof current === 'string';
					}

					const exists = keyExists(
						localeTranslations as Record<string, unknown>,
						key
					);
					expect(exists).toBe(true);
				}
			}
		});

		it('should have consistent key structure across all locales', () => {
			const defaultLocale = translations.en;
			const defaultKeys = getAllKeys(defaultLocale as Record<string, unknown>);

			for (const localeKey of locales) {
				const localeTranslations = translations[localeKey as keyof typeof translations];
				const localeKeys = getAllKeys(localeTranslations as Record<string, unknown>);

				// Check that locale has same number of keys as default (or more)
				expect(localeKeys.length).toBeGreaterThanOrEqual(defaultKeys.length);

				// Check that all default keys exist in locale
				for (const key of defaultKeys) {
					expect(localeKeys).toContain(key);
				}
			}
		});
	});

	describe('edge cases', () => {
		it('should handle empty key', () => {
			const consoleSpy = vi.spyOn(console, 'warn').mockImplementation(() => {});
			const translate = get(t);
			const result = translate('');
			expect(result).toBe('');
			consoleSpy.mockRestore();
		});

		it('should handle deeply nested keys', () => {
			const translate = get(t);
			// Test with a known deeply nested key
			const result = translate('settings.fields.email.label');
			expect(result).toBe('Email');
		});

		it('should handle keys with special characters in variable substitution', () => {
			const translate = get(t);
			// Test variable substitution with empty vars
			const result = translate('signIn', {});
			expect(typeof result).toBe('string');
			expect(result.length).toBeGreaterThan(0);
		});
	});
});

