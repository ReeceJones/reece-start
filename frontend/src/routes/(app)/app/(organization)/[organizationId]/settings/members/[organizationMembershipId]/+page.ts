import { get } from '$lib';
import { getOrganizationMembershipResponseSchema } from '$lib/schemas/organization-membership';

export const load = async ({ params, fetch }) => {
	const { organizationMembershipId } = params;

	const organizationMembership = await get(
		`/api/organization-memberships/${organizationMembershipId}`,
		{
			fetch,
			responseSchema: getOrganizationMembershipResponseSchema
		}
	);

	return {
		organizationMembership
	};
};
