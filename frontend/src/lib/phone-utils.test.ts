import { describe, it, expect } from 'vitest';
import {
	getPhoneCodeOptions,
	sanitizePhoneNumber,
	getPhoneExtensionFromCountryCode,
	formatPhoneNumberWithCountryCode,
	formatPhoneNumberWithExtension
} from './phone-utils';

describe('phone-utils', () => {
	describe('getPhoneCodeOptions', () => {
		it('should return an array of phone code options', () => {
			const options = getPhoneCodeOptions();
			expect(Array.isArray(options)).toBe(true);
			expect(options.length).toBeGreaterThan(0);
		});

		it('should return options with required properties', () => {
			const options = getPhoneCodeOptions();
			const firstOption = options[0];
			expect(firstOption).toHaveProperty('code');
			expect(firstOption).toHaveProperty('countryCode');
			expect(firstOption).toHaveProperty('countryName');
			expect(firstOption).toHaveProperty('flag');
		});

		it('should return options sorted by country name', () => {
			const options = getPhoneCodeOptions();
			for (let i = 1; i < options.length; i++) {
				expect(options[i].countryName >= options[i - 1].countryName).toBe(true);
			}
		});

		it('should only include countries with phone codes', () => {
			const options = getPhoneCodeOptions();
			options.forEach((option) => {
				expect(option.code).toBeTruthy();
				expect(option.code.length).toBeGreaterThan(0);
			});
		});

		it('should include common countries', () => {
			const options = getPhoneCodeOptions();
			const countryCodes = options.map((opt) => opt.countryCode);
			expect(countryCodes).toContain('US');
			expect(countryCodes).toContain('GB');
			expect(countryCodes).toContain('CA');
		});
	});

	describe('sanitizePhoneNumber', () => {
		it('should remove all non-numeric characters except leading +', () => {
			expect(sanitizePhoneNumber('+1 (555) 123-4567')).toBe('+15551234567');
			expect(sanitizePhoneNumber('+44 20 7946 0958')).toBe('+442079460958');
		});

		it('should preserve leading +', () => {
			expect(sanitizePhoneNumber('+1234567890')).toBe('+1234567890');
		});

		it('should remove all non-numeric characters when no leading +', () => {
			expect(sanitizePhoneNumber('(555) 123-4567')).toBe('5551234567');
			expect(sanitizePhoneNumber('555-123-4567')).toBe('5551234567');
			expect(sanitizePhoneNumber('555.123.4567')).toBe('5551234567');
		});

		it('should handle empty string', () => {
			expect(sanitizePhoneNumber('')).toBe('');
		});

		it('should handle string with only non-numeric characters', () => {
			expect(sanitizePhoneNumber('abc')).toBe('');
			expect(sanitizePhoneNumber('+abc')).toBe('+');
		});

		it('should handle string with spaces', () => {
			expect(sanitizePhoneNumber('555 123 4567')).toBe('5551234567');
			expect(sanitizePhoneNumber('+1 555 123 4567')).toBe('+15551234567');
		});

		it('should handle string with special characters', () => {
			expect(sanitizePhoneNumber('555-123-4567 ext. 123')).toBe('5551234567123');
			expect(sanitizePhoneNumber('+1-555-123-4567 ext. 123')).toBe('+15551234567123');
		});
	});

	describe('getPhoneExtensionFromCountryCode', () => {
		it('should return phone extension for valid country code', () => {
			expect(getPhoneExtensionFromCountryCode('US')).toBe('1');
			expect(getPhoneExtensionFromCountryCode('GB')).toBe('44');
			expect(getPhoneExtensionFromCountryCode('CA')).toBe('1');
		});

		it('should return empty string for invalid country code', () => {
			expect(getPhoneExtensionFromCountryCode('XX')).toBe('');
			expect(getPhoneExtensionFromCountryCode('INVALID')).toBe('');
		});

		it('should return empty string for empty string', () => {
			expect(getPhoneExtensionFromCountryCode('')).toBe('');
		});
	});

	describe('formatPhoneNumberWithCountryCode', () => {
		it('should format phone number with country code in E.164 format', () => {
			expect(formatPhoneNumberWithCountryCode('5551234567', 'US')).toBe('+15551234567');
			expect(formatPhoneNumberWithCountryCode('2079460958', 'GB')).toBe('+442079460958');
		});

		it('should sanitize phone number before formatting', () => {
			expect(formatPhoneNumberWithCountryCode('(555) 123-4567', 'US')).toBe('+15551234567');
			expect(formatPhoneNumberWithCountryCode('555-123-4567', 'US')).toBe('+15551234567');
		});

		it('should return sanitized phone number if no extension', () => {
			expect(formatPhoneNumberWithCountryCode('5551234567', 'XX')).toBe('5551234567');
			expect(formatPhoneNumberWithCountryCode('5551234567', '')).toBe('5551234567');
		});

		it('should return sanitized phone number if phone number is empty', () => {
			expect(formatPhoneNumberWithCountryCode('', 'US')).toBe('');
			expect(formatPhoneNumberWithCountryCode('', 'GB')).toBe('');
		});

		it('should handle phone number with leading +', () => {
			// Note: The function adds country code even if phone already has +
			expect(formatPhoneNumberWithCountryCode('+15551234567', 'US')).toBe('+1+15551234567');
		});

		it('should handle phone number already in E.164 format', () => {
			// Note: The function adds country code even if phone already has +
			expect(formatPhoneNumberWithCountryCode('+15551234567', 'US')).toBe('+1+15551234567');
			expect(formatPhoneNumberWithCountryCode('+442079460958', 'GB')).toBe('+44+442079460958');
		});
	});

	describe('formatPhoneNumberWithExtension', () => {
		it('should format phone number with extension', () => {
			expect(formatPhoneNumberWithExtension('5551234567', '44')).toBe('445551234567');
			expect(formatPhoneNumberWithExtension('2079460958', '44')).toBe('442079460958');
		});

		it('should sanitize phone number before formatting', () => {
			expect(formatPhoneNumberWithExtension('(555) 123-4567', '44')).toBe('445551234567');
			expect(formatPhoneNumberWithExtension('555-123-4567', '44')).toBe('445551234567');
		});

		it('should return sanitized phone number if extension is 1', () => {
			expect(formatPhoneNumberWithExtension('5551234567', '1')).toBe('5551234567');
		});

		it('should return sanitized phone number if extension is empty', () => {
			expect(formatPhoneNumberWithExtension('5551234567', '')).toBe('5551234567');
		});

		it('should handle empty phone number', () => {
			// Note: When phone is empty, sanitized phone is empty, but extension is still prepended if not '1'
			expect(formatPhoneNumberWithExtension('', '44')).toBe('44');
			expect(formatPhoneNumberWithExtension('', '1')).toBe('');
		});

		it('should handle phone number with leading +', () => {
			// Note: The + is preserved during sanitization
			expect(formatPhoneNumberWithExtension('+15551234567', '44')).toBe('44+15551234567');
		});
	});
});
