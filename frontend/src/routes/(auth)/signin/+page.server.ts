import { fail, redirect } from '@sveltejs/kit';
import type { Actions } from './$types';
import { z } from 'zod';
import { post, ApiError } from '$lib';
import { env } from '$env/dynamic/private';
import { performGoogleOAuth } from '$lib/server/oauth';

const loginUserRequestSchema = z.object({
	data: z.object({
		attributes: z.object({
			email: z.string().min(1),
			password: z.string().min(1)
		})
	})
});

const loginUserResponseSchema = z.object({
	data: z.object({
		id: z.string(),
		type: z.literal('user'),
		attributes: z.object({
			name: z.string(),
			email: z.string()
		}),
		meta: z.object({
			token: z.string()
		})
	})
});

export const actions = {
	signin: async ({ cookies, request, fetch }) => {
		const data = await request.formData();
		const email = data.get('email') as string;
		const password = data.get('password') as string;
		const searchParams = new URLSearchParams(request.url.slice(request.url.indexOf('?')));
		const redirectUrl = searchParams.get('redirect') ?? '/app';

		if (!email || !password) {
			return fail(400, { success: false, message: 'Please fill out all the fields correctly.' });
		}

		try {
			const userWithToken = await post(
				`/api/users/login`,
				{
					data: {
						attributes: {
							email,
							password
						}
					}
				},
				{
					fetch,
					requestSchema: loginUserRequestSchema,
					responseSchema: loginUserResponseSchema
				}
			);

			// set session token cookie
			cookies.set('app-session-token', userWithToken.data.meta.token, {
				httpOnly: true,
				secure: env.NODE_ENV === 'production',
				sameSite: 'strict',
				path: '/',
				maxAge: 60 * 60 * 24 * 30 // 30 days
			});
		} catch (error) {
			if (error instanceof ApiError) {
				if (error.code === 401) {
					return fail(401, { success: false, message: 'Invalid email or password' });
				}

				return fail(500, { success: false, message: error.message });
			}

			return fail(500, {
				success: false,
				message: 'An unknown error ocurred processing your request, please try again later.'
			});
		}

		redirect(302, redirectUrl);
	},
	oauthGoogle: async ({ request }) => {
		const data = await request.formData();
		const redirectUrl = data.get('redirect') as string | undefined;

		performGoogleOAuth(redirectUrl ?? '/app');
	}
} satisfies Actions;
