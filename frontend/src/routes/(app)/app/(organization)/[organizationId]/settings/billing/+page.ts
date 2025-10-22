import { get } from '$lib';
import { getSubscriptionResponseSchema } from '$lib/schemas/stripe';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params, fetch }) => {
	const subscription = await get(`/api/organizations/${params.organizationId}/subscription`, {
		fetch,
		responseSchema: getSubscriptionResponseSchema
	});

	return {
		subscription
	};
};

