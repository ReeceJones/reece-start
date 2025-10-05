import { redirect, type Actions } from '@sveltejs/kit';
import { z } from 'zod';
import { post } from '$lib/api';
import { createUserResponseSchema } from '$lib/schemas/user';
import { setTokenInCookies } from '$lib/server/auth';

// Schema for Google OAuth callback request
const GoogleOAuthCallbackRequest = z.object({
	data: z.object({
		attributes: z.object({
			code: z.string(),
			state: z.string(),
			redirectUri: z.string()
		})
	})
});

export const actions: Actions = {
	default: async (requestEvent) => {
		const { request, fetch, url } = requestEvent;
		const formData = await request.formData();
		const code = formData.get('code') as string;
		const state = formData.get('state') as string;
		const redirectUrl = formData.get('redirect') as string;

		if (!code || !state) {
			return {
				success: false,
				message: 'Missing OAuth parameters'
			};
		}

		try {
			// Get the full redirect URI
			const origin = url.origin;
			const redirectUri = `${origin}/oauth/google/callback`;

			// Call the backend OAuth endpoint
			const response = await post(
				'/api/oauth/google/callback',
				{
					data: {
						attributes: {
							code,
							state,
							redirectUri
						}
					}
				},
				{
					fetch,
					requestSchema: GoogleOAuthCallbackRequest,
					responseSchema: createUserResponseSchema
				}
			);

			// Set the token using the existing auth function
			if (response.data.meta.token) {
				setTokenInCookies(requestEvent, response.data.meta.token);
			}
		} catch (error) {
			console.error('OAuth callback error:', error);

			return {
				success: false,
				message: error instanceof Error ? error.message : 'Authentication failed. Please try again.'
			};
		}

		// Redirect to the intended destination or dashboard
		const finalRedirect = redirectUrl || '/app';
		throw redirect(302, finalRedirect);
	}
};
