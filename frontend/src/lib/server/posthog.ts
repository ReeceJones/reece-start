import posthog, { PostHog } from 'posthog-node';
import { PUBLIC_POSTHOG_KEY, PUBLIC_POSTHOG_HOST } from '$env/static/public';

let _client: PostHog | null = null;

export async function withPosthog(fn: (client: PostHog) => Promise<void>) {
	const client = getPostHogClient();
	try {
		await fn(client);
	} finally {
		await client.shutdown();
	}
}

function getPostHogClient() {
	if (!_client) {
		_client = new posthog.PostHog(PUBLIC_POSTHOG_KEY, {
			host: PUBLIC_POSTHOG_HOST
		});
	}
	return _client;
}
