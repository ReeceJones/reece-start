import type { Actions } from './$types';
import { env } from '$env/dynamic/private';
import { cookieName, locales } from '$lib/paraglide/runtime';
import { z } from 'zod';
import { isParseSuccess, parseFormData } from '$lib/server/schema';

const setLocaleFormSchema = z.object({
	locale: z.enum(locales)
});

export const actions = {
	setLocale: async ({ cookies, request }) => {
		const formData = await parseFormData(request, setLocaleFormSchema);

		if (!isParseSuccess(formData)) {
			return formData;
		}

		const { locale } = formData;

		// Set the locale cookie with httpOnly, secure, and sameSite settings
		cookies.set(cookieName, locale, {
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
