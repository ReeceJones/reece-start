const DEFAULT_LOCALE = 'en';
const LOCALE_COOKIE_NAME = 'app-locale';

export const load = async ({ cookies }) => {
	const locale: string = cookies.get(LOCALE_COOKIE_NAME) || DEFAULT_LOCALE;
	return {
		locale
	};
};
