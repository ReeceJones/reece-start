import { fail, redirect } from '@sveltejs/kit';
import type { Actions } from './$types';
import { z } from 'zod';
import { post } from '$lib';
import { env } from '$env/dynamic/private';

const createUserRequestSchema = z.object({
	data: z.object({
		attributes: z.object({
			name: z.string().min(1),
			email: z.string().min(1),
			password: z.string().min(1)
		})
	})
});

const createUserResponseSchema = z.object({
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
	default: async ({ cookies, request, fetch }) => {
		const data = await request.formData();
		const name = data.get('name') as string;
		const email = data.get('email') as string;
		const password = data.get('password') as string;

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

		console.log('Created user', userWithToken);

		redirect(302, '/app');
	}
} satisfies Actions;
