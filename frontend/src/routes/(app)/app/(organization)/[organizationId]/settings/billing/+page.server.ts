import { ApiError, post } from '$lib';
import {
	createBillingPortalSessionRequestSchema,
	createBillingPortalSessionResponseSchema,
	createCheckoutSessionRequestSchema,
	createCheckoutSessionResponseSchema,
	type CreateBillingPortalSessionRequest,
	type CreateBillingPortalSessionResponse,
	type CreateCheckoutSessionRequest,
	type CreateCheckoutSessionResponse
} from '$lib/schemas/stripe';
import { authenticate } from '$lib/server/auth';
import { fail, redirect } from '@sveltejs/kit';

export const load = async () => {
	authenticate();
};

export const actions = {
	checkout: async ({ fetch, params, url }) => {
		const { organizationId } = params;

		const requestBody: CreateCheckoutSessionRequest = {
			successUrl: url.toString(),
			cancelUrl: url.toString()
		};

		let response: CreateCheckoutSessionResponse;
		try {
			response = await post(`/api/organizations/${organizationId}/checkout-session`, requestBody, {
				fetch: fetch,
				requestSchema: createCheckoutSessionRequestSchema,
				responseSchema: createCheckoutSessionResponseSchema
			});
		} catch (error) {
			if (error instanceof ApiError) {
				return fail(error.code, { success: false, message: error.message });
			}

			console.error(error);

			return fail(500, {
				success: false,
				message: 'An unknown error ocurred processing your request, please try again later'
			});
		}

		return redirect(303, response.data.attributes.url);
	},
	portal: async ({ fetch, params, url }) => {
		const { organizationId } = params;

		const requestBody: CreateBillingPortalSessionRequest = {
			returnUrl: url.toString()
		};

		let response: CreateBillingPortalSessionResponse;
		try {
			response = await post(
				`/api/organizations/${organizationId}/billing-portal-session`,
				requestBody,
				{
					fetch: fetch,
					requestSchema: createBillingPortalSessionRequestSchema,
					responseSchema: createBillingPortalSessionResponseSchema
				}
			);
		} catch (error) {
			if (error instanceof ApiError) {
				return fail(error.code, { success: false, message: error.message });
			}

			console.error(error);

			return fail(500, {
				success: false,
				message: 'An unknown error ocurred processing your request, please try again later'
			});
		}

		return redirect(303, response.data.attributes.url);
	}
};
