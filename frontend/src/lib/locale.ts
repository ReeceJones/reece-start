export function localeToLanguageName(currentLocale: string, targetLocale: string) {
	const displayNames = new Intl.DisplayNames(currentLocale, { type: 'language' });

	return displayNames.of(targetLocale);
}
