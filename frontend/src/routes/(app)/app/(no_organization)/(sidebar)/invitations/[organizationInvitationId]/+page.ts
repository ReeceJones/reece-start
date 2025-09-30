import { get } from '$lib';
import { getOrganizationInvitationResponseSchema } from '$lib/schemas/organization-invitation';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, params }) => {
	const { organizationInvitationId } = params;

	const invitation = await get(`/api/organization-invitations/${organizationInvitationId}`, {
		fetch,
		responseSchema: getOrganizationInvitationResponseSchema
	});

	return {
		invitation: invitation
	};
};
