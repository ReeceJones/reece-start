import { get } from '$lib';
import {
	getOrganizationMembershipsQuerySchema,
	getOrganizationMembershipsResponseSchema
} from '$lib/schemas/organization-membership';

export const load = async ({ params, fetch }) => {
	const { organizationId } = params;

	const memberships = await get(`/api/organization-memberships`, {
		fetch: fetch,
		responseSchema: getOrganizationMembershipsResponseSchema,
		paramsSchema: getOrganizationMembershipsQuerySchema,
		params: {
			organizationId
		}
	});

	return {
		organizationId,
		memberships: memberships
	};
};
