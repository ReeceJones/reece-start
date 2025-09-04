import { authenticate, getUserAndValidateToken, getUserScopes } from '$lib/server/auth.js';

export const load = async () => {
	authenticate();
	// This won't run on every request, but it will be called on page reloads in the worst case so it should be fine
	const { user } = await getUserAndValidateToken();
	const userScopes = getUserScopes();

	return {
		user,
		userScopes
	};
};
