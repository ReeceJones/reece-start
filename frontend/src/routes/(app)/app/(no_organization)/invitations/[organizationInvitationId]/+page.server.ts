import { ApiError, post } from '$lib';
import {
	acceptOrganizationInvitationRequestSchema,
	acceptOrganizationInvitationResponseSchema,
	declineOrganizationInvitationRequestSchema,
	declineOrganizationInvitationResponseSchema
} from '$lib/schemas/organization-invitation';
import { authenticate } from '$lib/server/auth';
import type { Actions } from './$types';
import { API_TYPES } from '$lib/schemas/api';
import { fail, redirect } from '@sveltejs/kit';

export const load = async () => {
	authenticate();
};

export const actions = {
	accept: async ({ fetch, params }) => {
		let response;
		try {
			response = await post(
				`/api/organization-invitations/${params.organizationInvitationId}/accept`,
				{
					data: {
						type: API_TYPES.organizationInvitation,
						id: params.organizationInvitationId
					}
				},
				{
					fetch,
					requestSchema: acceptOrganizationInvitationRequestSchema,
					responseSchema: acceptOrganizationInvitationResponseSchema
				}
			);
		} catch (apiError) {
			if (apiError instanceof ApiError) {
				return fail(apiError.code, { success: false, message: apiError.message });
			}

			console.error(apiError);

			return fail(500, {
				success: false,
				message: 'An unknown error ocurred processing your request, please try again later'
			});
		}

		redirect(303, `/app/${response.data.relationships.organization.data.id}`);
	},
	decline: async ({ fetch, params }) => {
		try {
			await post(
				`/api/organization-invitations/${params.organizationInvitationId}/decline`,
				{
					data: {
						type: API_TYPES.organizationInvitation,
						id: params.organizationInvitationId
					}
				},
				{
					fetch,
					requestSchema: declineOrganizationInvitationRequestSchema,
					responseSchema: declineOrganizationInvitationResponseSchema
				}
			);
		} catch (apiError) {
			if (apiError instanceof ApiError) {
				return fail(apiError.code, { success: false, message: apiError.message });
			}

			console.error(apiError);

			return fail(500, {
				success: false,
				message: 'An unknown error ocurred processing your request, please try again later'
			});
		}

		redirect(303, '/app');
	}
} satisfies Actions;
