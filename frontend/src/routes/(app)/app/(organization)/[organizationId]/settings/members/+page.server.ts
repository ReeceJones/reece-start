import { ApiError, post } from '$lib';
import { API_TYPES } from '$lib/schemas/api';
import {
	inviteToOrganizationRequestSchema,
	inviteToOrganizationResponseSchema
} from '$lib/schemas/organization-invitation.js';
import { authenticate } from '$lib/server/auth';
import { fail } from '@sveltejs/kit';
import type { Actions } from './$types';

export const load = async () => {
	authenticate();
};

export const actions = {
	invite: async ({ request, fetch, params }) => {
		const data = await request.formData();
		const email = data.get('email') as string;

		try {
			const response = await post(
				`/api/organization-invitations`,
				{
					data: {
						type: API_TYPES.organizationInvitation,
						attributes: {
							email
						},
						relationships: {
							organization: {
								data: {
									type: API_TYPES.organization,
									id: params.organizationId
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
	}
} satisfies Actions;
