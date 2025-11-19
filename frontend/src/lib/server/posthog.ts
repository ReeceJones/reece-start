import { PostHog } from 'posthog-node';
import { env } from '$env/dynamic/public';

export async function withPosthog(fn: (client: PostHog) => Promise<void>) {
	// only log posthog events in prod
	if (process.env.NODE_ENV !== 'production') {
		console.log('Not logging posthog events in non-production environment');
		return;
	}

	if (!env.PUBLIC_POSTHOG_KEY) {
		console.log('No PostHog key found');
		return;
	}

	const client = getPostHogClient();
	try {
		await fn(client);
	} finally {
		await client.shutdown();
	}
}

export function getPostHogClient() {
	return new PostHog(env.PUBLIC_POSTHOG_KEY ?? '', {
		host: env.PUBLIC_POSTHOG_HOST
	});
}
