import type { Actions } from './$types';
import { env } from '$env/dynamic/private';
import { locales } from '$lib/i18n';
import { z } from 'zod';
import { isParseSuccess, parseFormData } from '$lib/server/schema';

const LOCALE_COOKIE_NAME = 'app-locale';

const setLocaleFormSchema = z.object({
	locale: z.string().refine((val) => locales.includes(val), {
		message: 'Unsupported locale'
	})
});

export const actions = {
	setLocale: async ({ cookies, request }) => {
		const formData = await parseFormData(request, setLocaleFormSchema);

		if (!isParseSuccess(formData)) {
			return formData;
		}

		const { locale } = formData;

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
