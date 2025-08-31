import { fail } from '@sveltejs/kit';
import { ApiError, patch } from '$lib';
import type { Actions } from './$types';
import { authenticate } from '$lib/server/auth';
import { updateUserRequestSchema, updateUserResponseSchema } from '$lib/schemas/user';

export const load = async () => {
	authenticate();
};

export const actions = {
	default: async ({ request, fetch }) => {
		const data = await request.formData();
		const userId = data.get('userId') as string;
		const email = data.get('email') as string;
		const password = data.get('password') as string;
		const confirmPassword = data.get('confirmPassword') as string;

		if (!userId || !email) {
			return fail(400, { success: false, message: 'Please fill out all the fields correctly.' });
		}

		if (password && password !== confirmPassword) {
			return fail(400, { success: false, message: 'Passwords do not match.' });
		}

		if (password && password.length < 8) {
			return fail(400, { success: false, message: 'Password must be at least 8 characters long.' });
		}

		try {
			const user = await patch(
				`/api/users/${userId}`,
				{
					data: {
						id: userId,
						type: 'user',
						attributes: {
							email,
							password
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
