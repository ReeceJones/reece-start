import { redirect, type Actions } from '@sveltejs/kit';
import { z } from 'zod';
import { post } from '$lib/api';
import { createUserResponseSchema } from '$lib/schemas/user';
import { setTokenInCookies } from '$lib/server/auth';
import { googleOAuthCallbackFormSchema } from '$lib/schemas/user.server';
import { isParseSuccess, parseFormData } from '$lib/server/schema';

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
		const formData = await parseFormData(request, googleOAuthCallbackFormSchema);

		if (!isParseSuccess(formData)) {
			return formData;
		}

		const { code, state, redirect: redirectUrl } = formData;

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
