const DEFAULT_LOCALE = 'en';
const LOCALE_COOKIE_NAME = 'app-locale';

export const load = async ({ cookies }: { cookies: any }) => {
	const locale = cookies.get(LOCALE_COOKIE_NAME) || DEFAULT_LOCALE;
	return {
		locale
	};
};
