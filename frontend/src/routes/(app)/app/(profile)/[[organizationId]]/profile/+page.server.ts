import { fail } from '@sveltejs/kit';
import { ApiError, base64Encode, patch } from '$lib';
import type { Actions } from './$types';
import { authenticate } from '$lib/server/auth';
import { updateUserRequestSchema, updateUserResponseSchema } from '$lib/schemas/user';
import { updateUserProfileFormSchema } from '$lib/schemas/user.server';
import { isParseSuccess, parseFormData } from '$lib/server/schema';

export const load = async () => {
	authenticate();
};

export const actions = {
	default: async ({ request, fetch }) => {
		const formData = await parseFormData(request, updateUserProfileFormSchema);

		if (!isParseSuccess(formData)) {
			return formData;
		}

		const { userId, name, logo } = formData;
		let logoData: string | undefined = undefined;

		if (logo && logo.size > 0) {
			// need to base64 encode the logo content
			const logoBuffer = await logo.arrayBuffer();
			logoData = base64Encode(logoBuffer);
		}

		try {
			const user = await patch(
				`/api/users/${userId}`,
				{
					data: {
						id: userId,
						type: 'user',
						attributes: {
							name,
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
