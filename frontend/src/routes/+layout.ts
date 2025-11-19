import posthog from 'posthog-js';
import { browser } from '$app/environment';
import { env } from '$env/dynamic/public';

export const load = async ({ data }) => {
	if (browser) {
		const posthogKey = env.PUBLIC_POSTHOG_KEY;

		if (posthogKey) {
			posthog.init(posthogKey, {
				api_host: env.PUBLIC_POSTHOG_HOST,
				capture_exceptions: true
			});
		}
	}
	return data;
};
