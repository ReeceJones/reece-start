import { baseLocale, cookieName } from '$lib/paraglide/runtime';

export const load = async ({ cookies }) => {
	const locale: string = cookies.get(cookieName) || baseLocale;
	return {
		locale
	};
};
