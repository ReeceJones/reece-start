import { fail, redirect } from '@sveltejs/kit';
import type { Actions } from './$types';
import { post } from '$lib';
import { env } from '$env/dynamic/private';
import { createUserRequestSchema, createUserResponseSchema } from '$lib/schemas/user';
import { performGoogleOAuth } from '$lib/server/oauth';
import { signupFormSchema } from '$lib/schemas/user.server';
import { isParseSuccess, parseFormData } from '$lib/server/schema';

export const actions = {
	signup: async ({ cookies, request, fetch }) => {
		const formData = await parseFormData(request, signupFormSchema);

		if (!isParseSuccess(formData)) {
			return formData;
		}

		const { name, email, password } = formData;
		const isSignInDisabled = env.PUBLIC_DISABLE_SIGNIN === 'true';

		const searchParams = new URLSearchParams(request.url.slice(request.url.indexOf('?')));
		const redirectUrl = searchParams.get('redirect') ?? '/app';

		if (isSignInDisabled) {
			return fail(403, { success: false, message: 'Sign in is disabled' });
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
	oauthGoogle: async ({ request }) => {
		const data = await request.formData();
		const redirectUrl = data.get('redirect') as string | undefined;

		performGoogleOAuth(redirectUrl ?? '/app');
	}
} satisfies Actions;
