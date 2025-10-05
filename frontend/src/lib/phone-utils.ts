import { countries, getEmojiFlag } from 'countries-list';

export interface PhoneCodeOption {
	code: string;
	countryCode: string;
	countryName: string;
	flag: string;
}

export function getPhoneCodeOptions(): PhoneCodeOption[] {
	return Object.entries(countries)
		.filter(([, country]) => country.phone && country.phone.length > 0)
		.map(([countryCode, country]) => ({
			code: country.phone[0].toString(),
			countryCode,
			countryName: country.name,
			flag: getEmojiFlag(countryCode as keyof typeof countries)
		}))
		.sort((a, b) => {
			// Sort by country name for better UX
			return a.countryName.localeCompare(b.countryName);
		});
}

export function sanitizePhoneNumber(phoneNumber: string): string {
	// Remove all non-numeric characters except the leading +
	if (phoneNumber.startsWith('+')) {
		return '+' + phoneNumber.slice(1).replace(/\D/g, '');
	}
	return phoneNumber.replace(/\D/g, '');
}

export function getPhoneExtensionFromCountryCode(countryCode: string): string {
	const country = countries[countryCode as keyof typeof countries];
	return country?.phone?.[0]?.toString() || '';
}

export function formatPhoneNumberWithCountryCode(phoneNumber: string, countryCode: string): string {
	const extension = getPhoneExtensionFromCountryCode(countryCode);
	const sanitizedPhone = sanitizePhoneNumber(phoneNumber);

	if (!extension || !sanitizedPhone) {
		return sanitizedPhone;
	}

	// Return in E.164 format: +{country_code}{phone_number}
	return `+${extension}${sanitizedPhone}`;
}

export function formatPhoneNumberWithExtension(phoneNumber: string, extension: string): string {
	const sanitizedPhone = sanitizePhoneNumber(phoneNumber);
	return extension && extension !== '1' ? `${extension}${sanitizedPhone}` : sanitizedPhone;
}
