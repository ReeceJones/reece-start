import { fail } from '@sveltejs/kit';
import { ApiError, base64Encode, patch } from '$lib';
import type { Actions } from './$types';
import { authenticate } from '$lib/server/auth';
import {
	updateOrganizationRequestSchema,
	organizationResponseSchema,
	updateOrganizationFormSchema
} from '$lib/schemas/organization';
import { isParseSuccess, parseFormData } from '$lib/server/schema';

export const load = async () => {
	authenticate();
};

export const actions = {
	default: async ({ request, fetch, params }) => {
		const formData = await parseFormData(request, updateOrganizationFormSchema);

		if (!isParseSuccess(formData)) {
			return formData;
		}

		const { organizationId } = params;
		const { name, description, logo } = formData;

		let logoData: string | undefined = undefined;

		if (!organizationId || !name) {
			return fail(400, {
				success: false,
				message: 'Please fill out all the required fields correctly.'
			});
		}

		if (logo) {
			const logoFile = logo as unknown as File;

			if (logoFile.size > 3_000_000) {
				return fail(400, { success: false, message: 'Logo must be less than 3MB.' });
			}

			if (logoFile.size > 0) {
				// need to base64 encode the logo content
				const logoBuffer = await logoFile.arrayBuffer();
				logoData = base64Encode(logoBuffer);
			}
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
