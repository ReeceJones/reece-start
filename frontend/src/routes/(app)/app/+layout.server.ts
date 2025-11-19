import {
	authenticate,
	getIsImpersonatingUser,
	getUserAndValidateToken,
	getUserScopes
} from '$lib/server/auth.js';
import { withPosthog } from '$lib/server/posthog';

export const load = async () => {
	authenticate();
	// This won't run on every request, but it will be called on page reloads in the worst case so it should be fine
	const { user } = await getUserAndValidateToken();
	const userScopes = getUserScopes();
	const isImpersonatingUser = getIsImpersonatingUser();

	// Identify the authenticated user in Posthog
	await withPosthog(async (client) => {
		client.identify({
			distinctId: user.data.id,
			properties: {
				...user.data.attributes,
				scopes: userScopes
			}
		});
	});

	return {
		user,
		userScopes,
		isImpersonatingUser
	};
};
