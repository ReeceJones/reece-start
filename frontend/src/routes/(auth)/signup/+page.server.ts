import { fail, redirect } from '@sveltejs/kit';
import type { Actions } from './$types';
import { post } from '$lib';
import { env } from '$env/dynamic/private';
import { createUserRequestSchema, createUserResponseSchema } from '$lib/schemas/user';
import { performGoogleOAuth } from '$lib/server/oauth';

export const actions = {
	signup: async ({ cookies, request, fetch }) => {
		const data = await request.formData();
		const name = data.get('name') as string;
		const email = data.get('email') as string;
		const password = data.get('password') as string;

		const searchParams = new URLSearchParams(request.url.slice(request.url.indexOf('?')));
		const redirectUrl = searchParams.get('redirect') ?? '/app';

		if (!name || !email || !password) {
			return fail(400, { success: false, message: 'Please fill out all the fields correctly.' });
		}

		const userWithToken = await post(
			`/api/users`,
			{
				data: {
					attributes: {
						name,
						email,
						password
					}
				}
			},
			{
				fetch,
				requestSchema: createUserRequestSchema,
				responseSchema: createUserResponseSchema
			}
		);

		cookies.set('app-session-token', userWithToken.data.meta.token, {
			httpOnly: true,
			secure: env.NODE_ENV === 'production',
			sameSite: 'strict',
			path: '/',
			maxAge: 60 * 60 * 24 * 30 // 30 days
		});

		redirect(302, redirectUrl);
	},
	oauthGoogle: async ({ cookies, request, fetch }) => {
		const data = await request.formData();
		const redirectUrl = data.get('redirect') as string | undefined;

		performGoogleOAuth(redirectUrl ?? '/app');
	}
} satisfies Actions;
