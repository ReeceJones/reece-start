import { browser } from '$app/environment';
import { jwtDecode } from 'jwt-decode';
import type { JwtClaims, OrganizationScope } from './schemas/jwt';
import { getRequestEvent } from '$app/server';

export function hasScope(scope: OrganizationScope): boolean {
	if (browser) {
		// Get cookies using browser API
		const cookies = document.cookie;
		const token = cookies
			.split('; ')
			.find((row) => row.startsWith('app-session-token='))
			?.split('=')[1];
		if (!token) {
			return false;
		}

		const claims = jwtDecode<JwtClaims>(token);
		return claims.organization_scopes?.includes(scope) ?? false;
	} else {
		// Get cookies using server API
		const requestEvent = getRequestEvent();
		const token = requestEvent.cookies.get('app-session-token');
		if (!token) {
			return false;
		}
		const claims = jwtDecode<JwtClaims>(token);
		return claims.organization_scopes?.includes(scope) ?? false;
	}
}
