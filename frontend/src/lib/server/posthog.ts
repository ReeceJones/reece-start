import posthog, { PostHog } from 'posthog-node';
import { PUBLIC_POSTHOG_KEY, PUBLIC_POSTHOG_HOST } from '$env/static/public';

export async function withPosthog(fn: (client: PostHog) => Promise<void>) {
	// only log posthog events in prod
	if (process.env.NODE_ENV !== 'production') {
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
	return new posthog.PostHog(PUBLIC_POSTHOG_KEY, {
		host: PUBLIC_POSTHOG_HOST
	});
}
