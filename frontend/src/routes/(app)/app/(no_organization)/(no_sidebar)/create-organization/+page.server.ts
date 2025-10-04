import { ApiError, base64Encode, post } from '$lib';
import { formatUrl } from '$lib/organization-onboarding';
import { formatPhoneNumberWithCountryCode } from '$lib/phone-utils';
import {
	createOrganizationFormSchema,
	createOrganizationRequestSchema,
	organizationResponseSchema
} from '$lib/schemas/organization';
import { authenticate } from '$lib/server/auth';
import { isParseSuccess, parseFormData } from '$lib/server/schema';
import { fail, redirect } from '@sveltejs/kit';
import type z from 'zod';

export const load = async ({ locals }) => {
	authenticate();
	redirect(302, '/app/create-organization/basic-information');
};

export const actions = {
	default: async ({ request, fetch }) => {
		const formData = await parseFormData(request, createOrganizationFormSchema);

		if (!isParseSuccess(formData)) {
			return formData;
		}

		const { logo } = formData;
		let logoBase64 = undefined;

		if (logo) {
			const logoFile = logo as unknown as File;
			if (logoFile.size > 3_000_000) {
				return fail(400, { success: false, message: 'Logo must be less than 3MB.' });
			}

			if (logoFile.size > 0) {
				const logoBuffer = await logoFile.arrayBuffer();
				logoBase64 = base64Encode(logoBuffer);
			}
		}

		let organization: z.infer<typeof organizationResponseSchema>;
		try {
			organization = await post(
				`/api/organizations`,
				{
					data: {
						type: 'organization',
						attributes: {
							name: formData.name,
							description: formData.description,
							...(logoBase64 && { logo: logoBase64 }),
							contactEmail: formData.contactEmail,
							contactPhone: formatPhoneNumberWithCountryCode(
								formData.contactPhone || '',
								formData.contactPhoneCountry || ''
							),
							contactPhoneCountry: formData.contactPhoneCountry || '',
							address: {
								city: formData.addressCity || '',
								stateOrProvince: formData.addressStateOrProvince || '',
								zip: formData.addressZip || '',
								country: formData.addressCountry || '',
								line1: formData.addressLine1 || '',
								line2: formData.addressLine2 || ''
							},
							locale: formData.locale,
							entityType: formData.entityType
						}
					}
				},
				{
					fetch,
					responseSchema: organizationResponseSchema,
					requestSchema: createOrganizationRequestSchema
				}
			);
		} catch (apiError) {
			if (apiError instanceof ApiError) {
				return fail(apiError.code, { success: false, message: apiError.message });
			}

			return fail(500, {
				success: false,
				message: 'An unknown error ocurred processing your request, please try again later.'
			});
		}

		redirect(302, `/app/${organization.data.id}`);
	}
};
