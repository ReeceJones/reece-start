import { ApiError, del, post } from '$lib';
import { API_TYPES } from '$lib/schemas/api';
import {
	deleteInvitationFormSchema,
	inviteToOrganizationFormSchema,
	inviteToOrganizationRequestSchema,
	inviteToOrganizationResponseSchema
} from '$lib/schemas/organization-invitation.js';
import { authenticate } from '$lib/server/auth';
import { fail } from '@sveltejs/kit';
import type { Actions } from './$types';
import { isParseSuccess, parseFormData } from '$lib/server/schema';

export const load = async () => {
	authenticate();
};

export const actions = {
	invite: async ({ request, fetch, params }) => {
		const formData = await parseFormData(request, inviteToOrganizationFormSchema);

		if (!isParseSuccess(formData)) {
			return formData;
		}

		const { email, role } = formData;
		const { organizationId } = params;

		try {
			const response = await post(
				`/api/organization-invitations`,
				{
					data: {
						type: API_TYPES.organizationInvitation,
						attributes: {
							email,
							role
						},
						relationships: {
							organization: {
								data: {
									type: API_TYPES.organization,
									id: organizationId
								}
							}
						}
					}
				},
				{
					fetch,
					requestSchema: inviteToOrganizationRequestSchema,
					responseSchema: inviteToOrganizationResponseSchema
				}
			);

			return {
				success: true,
				message: 'Organization invitation sent successfully',
				data: response
			};
		} catch (error) {
			if (error instanceof ApiError) {
				return fail(error.code, { success: false, message: error.message });
			}

			console.error(error);

			return fail(500, {
				success: false,
				message: 'An unknown error ocurred processing your request, please try again later'
			});
		}
	},
	deleteInvitation: async ({ request, fetch, params }) => {
		const formData = await parseFormData(request, deleteInvitationFormSchema);

		if (!isParseSuccess(formData)) {
			return formData;
		}

		const { invitationId } = formData;

		try {
			await del(`/api/organization-invitations/${invitationId}`, {
				fetch
			});
		} catch (error) {
			if (error instanceof ApiError) {
				return fail(error.code, { success: false, message: error.message });
			}

			console.error(error);

			return fail(500, {
				success: false,
				message: 'An unknown error ocurred processing your request, please try again later'
			});
		}
	}
} satisfies Actions;
