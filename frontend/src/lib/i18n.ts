import { derived, writable } from 'svelte/store';
import { translations } from './translations';

const DEFAULT_LOCALE = 'en';

export const locale = writable(DEFAULT_LOCALE);
export const locales = Object.keys(translations);

// Helper function to get nested value from object using dot notation
function getNestedValue(obj: any, path: string): string | undefined {
	const keys = path.split('.');
	let current = obj;

	for (const key of keys) {
		if (current === undefined || current === null) {
			return undefined;
		}
		current = current[key];
	}

	return typeof current === 'string' ? current : undefined;
}

function translate(
	locale: keyof typeof translations,
	key: string,
	args?: Record<string, string>
): string {
	const defaultLocaleTranslations = translations[DEFAULT_LOCALE];

	if (!defaultLocaleTranslations) {
		throw new Error('Default locale not found in translations');
	}

	const localeTranslations = translations[locale] ?? defaultLocaleTranslations;

	if (!localeTranslations) {
		throw new Error(`Locale ${locale} not found in translations`);
	}

	let text = getNestedValue(localeTranslations, key);

	if (!text) {
		text = getNestedValue(defaultLocaleTranslations, key);
	}

	if (!text) {
		console.warn(`Translation key not found: ${key}`);
		return key;
	}

	Object.keys(args ?? {}).map((k) => {
		const regex = new RegExp(`{{${k}}}`, 'g');
		text = text!.replace(regex, args?.[k] ?? '');
	});

	return text;
}

export const t = derived(
	locale,
	($locale) =>
		(key: string, vars: Record<string, string> = {}) =>
			translate($locale as keyof typeof translations, key, vars)
);
