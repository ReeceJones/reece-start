import { fail } from '@sveltejs/kit';
import { ApiError, base64Encode, patch } from '$lib';
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
		const name = data.get('name') as string;
		const logo = data.get('logo') as File;
		let logoData: string | undefined = undefined;

		if (!userId || !name) {
			return fail(400, { success: false, message: 'Please fill out all the fields correctly.' });
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
