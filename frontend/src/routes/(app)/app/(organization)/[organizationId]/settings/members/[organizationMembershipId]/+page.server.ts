import { ApiError, del, patch } from '$lib';
import { API_TYPES } from '$lib/schemas/api.js';
import {
	updateOrganizationMembershipRequestSchema,
	updateOrganizationMembershipResponseSchema,
	type OrganizationMembershipRole
} from '$lib/schemas/organization-membership.js';
import { authenticate } from '$lib/server/auth';
import { fail, redirect } from '@sveltejs/kit';
import type { Actions } from './$types';

export const load = async () => {
	authenticate();
};

export const actions = {
	update: async ({ request, fetch, params }) => {
		const data = await request.formData();
		const { organizationMembershipId } = params;
		const role = data.get('role') as OrganizationMembershipRole;

		if (!role) {
			return fail(400, {
				success: false,
				message: 'Please fill out all the required fields correctly'
			});
		}

		try {
			const organizationMembership = await patch(
				`/api/organization-memberships/${organizationMembershipId}`,
				{
					data: {
						type: API_TYPES.organizationMembership,
						id: organizationMembershipId,
						attributes: {
							role: role
						}
					}
				},
				{
					fetch,
					responseSchema: updateOrganizationMembershipResponseSchema,
					requestSchema: updateOrganizationMembershipRequestSchema
				}
			);

			return {
				success: true,
				message: 'Organization membership updated successfully',
				data: organizationMembership
			};
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
	},
	delete: async ({ fetch, params }) => {
		const { organizationMembershipId } = params;

		try {
			await del(`/api/organization-memberships/${organizationMembershipId}`, {
				fetch
			});
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

		redirect(303, `/app/${params.organizationId}/settings/members`);
	}
} satisfies Actions;
