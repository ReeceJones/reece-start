import { fail } from '@sveltejs/kit';
import { ApiError, base64Encode, patch } from '$lib';
import type { Actions } from './$types';
import { authenticate } from '$lib/server/auth';
import {
	updateOrganizationRequestSchema,
	organizationResponseSchema
} from '$lib/schemas/organization';

export const load = async () => {
	authenticate();
};

export const actions = {
	default: async ({ request, fetch, params }) => {
		const data = await request.formData();
		const organizationId = data.get('organizationId') as string;
		const name = data.get('name') as string;
		const description = data.get('description') as string;
		const logo = data.get('logo') as File;
		let logoData: string | undefined = undefined;

		if (!organizationId || !name) {
			return fail(400, {
				success: false,
				message: 'Please fill out all the required fields correctly.'
			});
		}

		// Verify the organizationId matches the URL param
		if (organizationId !== params.organizationId) {
			return fail(400, { success: false, message: 'Invalid organization ID.' });
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
			const organization = await patch(
				`/api/organizations/${organizationId}`,
				{
					data: {
						type: 'organization',
						attributes: {
							name,
							description: description || undefined,
							logo: logoData
						}
					}
				},
				{
					fetch,
					responseSchema: organizationResponseSchema,
					requestSchema: updateOrganizationRequestSchema
				}
			);

			return {
				success: true,
				message: 'Organization updated successfully',
				data: organization
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
