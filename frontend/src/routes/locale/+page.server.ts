import { fail } from '@sveltejs/kit';
import type { Actions } from './$types';
import { env } from '$env/dynamic/private';
import { locales } from '$lib/i18n';

const LOCALE_COOKIE_NAME = 'app-locale';

export const actions = {
	setLocale: async ({ cookies, request }) => {
		const data = await request.formData();
		const locale = data.get('locale') as string;

		if (!locale) {
			return fail(400, {
				success: false,
				message: 'Locale is required'
			});
		}

		// Validate locale is one of the supported locales
		if (!locales.includes(locale)) {
			return fail(400, {
				success: false,
				message: 'Unsupported locale'
			});
		}

		// Set the locale cookie with httpOnly, secure, and sameSite settings
		cookies.set(LOCALE_COOKIE_NAME, locale, {
			path: '/',
			httpOnly: true,
			secure: env.NODE_ENV === 'production',
			sameSite: 'strict',
			maxAge: 60 * 60 * 24 * 365 // 1 year
		});

		return {
			success: true,
			locale
		};
	}
} satisfies Actions;
