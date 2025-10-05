import { redirect, error, type Actions } from '@sveltejs/kit';
import { authenticate } from '$lib/server/auth';
import { post } from '$lib/api';
import { API_TYPES } from '$lib/schemas/api';
import {
	createStripeOnboardingLinkRequestSchema,
	createStripeOnboardingLinkResponseSchema
} from '$lib/schemas/stripe';

export const load = async () => {
	authenticate();
};

export const actions: Actions = {
	default: async ({ fetch, params }) => {
		const { organizationId } = params;
		if (!organizationId) {
			throw error(400, 'Missing organizationId');
		}

		const response = await post(
			`/api/organizations/${organizationId}/stripe-onboarding-link`,
			{
				data: {
					type: API_TYPES.stripeAccountLink,
					relationships: {
						organization: {
							data: {
								id: organizationId,
								type: API_TYPES.organization
							}
						}
					}
				}
			},
			{
				fetch,
				requestSchema: createStripeOnboardingLinkRequestSchema,
				responseSchema: createStripeOnboardingLinkResponseSchema
			}
		);

		const url = response.data.attributes.url;
		throw redirect(302, url);
	}
};
