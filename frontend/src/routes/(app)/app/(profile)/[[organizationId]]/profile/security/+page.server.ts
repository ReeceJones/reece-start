import { fail } from '@sveltejs/kit';
import { ApiError, patch } from '$lib';
import type { Actions } from './$types';
import { authenticate } from '$lib/server/auth';
import { updateUserRequestSchema, updateUserResponseSchema } from '$lib/schemas/user';
import { updateUserSecurityFormSchema } from '$lib/schemas/user.server';
import { isParseSuccess, parseFormData } from '$lib/server/schema';

export const load = async () => {
	authenticate();
};

export const actions = {
	default: async ({ request, fetch }) => {
		const formData = await parseFormData(request, updateUserSecurityFormSchema);

		if (!isParseSuccess(formData)) {
			return formData;
		}

		const { userId, email, password, confirmPassword } = formData;

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
