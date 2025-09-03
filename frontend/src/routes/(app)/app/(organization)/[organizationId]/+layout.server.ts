import { authenticate, getMembershipScopes } from '$lib/server/auth';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async () => {
	authenticate();

	const membershipScopes = getMembershipScopes();

	return {
		membershipScopes
	};
};
