import { authenticate, getUserScopes } from '$lib/server/auth';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async () => {
	authenticate();

	const userScopes = getUserScopes();

	return {
		userScopes
	};
};
