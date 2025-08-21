import { fail } from '@sveltejs/kit';
import { ApiError, base64Encode, put } from '$lib';
import type { Actions } from './$types';
import { z } from 'zod';

const updateUserRequestSchema = z.object({
	data: z.object({
		id: z.string(),
		type: z.literal('user'),
		attributes: z.object({
			name: z.string().optional(),
			email: z.string().optional(),
			password: z.string().optional(),
			logo: z.string().optional()
		})
	})
});

const updateUserResponseSchema = z.object({
	data: z.object({
		id: z.string(),
		type: z.literal('user'),
		attributes: z.object({
			name: z.string(),
			email: z.string()
		})
	})
});

export const actions = {
	default: async ({ request, fetch }) => {
		const data = await request.formData();
		const userId = data.get('userId') as string;
		const name = data.get('name') as string;
		const email = data.get('email') as string;
		const password = data.get('password') as string;
		const confirmPassword = data.get('confirmPassword') as string;
		const logo = data.get('logo') as File;
		let logoData: string | undefined = undefined;

		if (!userId || !name || !email) {
			return fail(400, { success: false, message: 'Please fill out all the fields correctly.' });
		}

		if (password && password !== confirmPassword) {
			return fail(400, { success: false, message: 'Passwords do not match.' });
		}

		if (password && password.length < 8) {
			return fail(400, { success: false, message: 'Password must be at least 8 characters long.' });
		}

		if (logo.size > 3_000_000) {
			return fail(400, { success: false, message: 'Logo must be less than 3MB.' });
		}

		if (logo.size > 0) {
			// need to base64 encode the logo content
			const logoBuffer = await logo.arrayBuffer();
			logoData = base64Encode(logoBuffer);
		}

		try {
			const user = await put(
				`/api/users/${userId}`,
				{
					data: {
						id: userId,
						type: 'user',
						attributes: {
							name,
							email,
							password,
							logo: logoData
						}
					}
				},
				{
					fetch,
					responseSchema: updateUserResponseSchema,
					requestSchema: updateUserRequestSchema
				}
			);

			return {
				success: true,
				message: 'User updated successfully',
				data: user
			};
		} catch (apiError) {
			if (apiError instanceof ApiError) {
				return fail(apiError.code, { success: false, message: apiError.message });
			}

			console.error(apiError);

			return fail(500, {
				success: false,
				message: 'An unknown error ocurred processing your request, please try again later.'
			});
		}
	}
} satisfies Actions;
