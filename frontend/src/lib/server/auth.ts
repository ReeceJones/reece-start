import { getRequestEvent } from '$app/server';
import { ApiError, get, post } from '$lib/api';
import { API_TYPES } from '$lib/schemas/api';
import type { JwtClaims } from '$lib/schemas/jwt';
import { getSelfUserResponseSchema } from '$lib/schemas/user';
import {
	createAuthenticatedUserTokenRequestSchema,
	createAuthenticatedUserTokenResponseSchema
} from '$lib/schemas/user.server';
import { error, redirect, type RequestEvent } from '@sveltejs/kit';
import { jwtDecode } from 'jwt-decode';

export function authenticate() {
	const requestEvent = getRequestEvent();
	performAuthenticationCheck(requestEvent);
}

export async function performAuthenticationCheck(requestEvent: RequestEvent) {
	const { url, params } = requestEvent;
	const { organizationId } = params;

	// Only protect /app*
	if (!url.pathname.startsWith('/app')) {
		return;
	}

	// Validate token is present and is not expired
	validateTokenExpiration(requestEvent);

	if (url.pathname.startsWith('/app/admin')) {
		validateTokenAdmin(requestEvent);
	}

	// Validate token works for authenticated organization
	if (organizationId === undefined) {
		return;
	}

	await validateTokenOrganization(requestEvent);
}

export async function getUserAndValidateToken() {
	const requestEvent = getRequestEvent();
	const token = getDefinedToken(requestEvent);

	// get the issued at time from the token
	const { iat } = jwtDecode<JwtClaims>(token);

	if (!iat) {
		error(401, 'Invalid token: missing iat claim.');
	}

	let user;

	try {
		user = await get('/api/users/me', {
			fetch: requestEvent.fetch,
			responseSchema: getSelfUserResponseSchema
		});
	} catch (apiError) {
		if (apiError instanceof ApiError) {
			error(apiError.code, apiError.message);
		}

		console.error(apiError);

		error(500, 'An unknown error ocurred processing your request, please try again later.');
	}

	// If the token is older than the last issued at time, refresh the token
	if (
		user.data.meta.tokenRevocation.lastIssuedAt &&
		iat > user.data.meta.tokenRevocation.lastIssuedAt
	) {
		console.log('Token expired, refreshing');
		// If the token cannot be refreshed, sign the user out and redirect to the signin page
		if (!user.data.meta.tokenRevocation.canRefresh) {
			requestEvent.cookies.delete('app-session-token', { path: '/' });
			redirect(303, getRedirectUrl(requestEvent));
		}

		await refreshUserToken(requestEvent);
	}

	return {
		user
	};
}

export async function refreshUserToken(requestEvent: RequestEvent) {
	const { params, fetch } = requestEvent;
	const { organizationId } = params;

	console.log('Refreshing user token');

	let newToken: string;
	try {
		const response = await post(
			'/api/users/me/token',
			{
				data: {
					type: API_TYPES.token,
					relationships: {
						...(organizationId
							? {
									organization: {
										data: {
											id: organizationId,
											type: API_TYPES.organization
										}
									}
								}
							: {})
					}
				}
			},
			{
				fetch,
				requestSchema: createAuthenticatedUserTokenRequestSchema,
				responseSchema: createAuthenticatedUserTokenResponseSchema
			}
		);

		newToken = response.data.meta.token;
	} catch (apiError) {
		if (apiError instanceof ApiError) {
			console.error('Error validating token organization', apiError.code, apiError.message);

			if (apiError.code === 404) {
				// Organization membership not found, return auth error
				error(404, 'Organization not found');
			}
		}

		throw error;
	}

	// Set the token in the cookies
	setTokenInCookies(requestEvent, newToken);
}

export function getUserScopes() {
	const requestEvent = getRequestEvent();
	const token = getDefinedToken(requestEvent);
	if (!token) {
		return [];
	}
	const claims = jwtDecode<JwtClaims>(token);
	return claims.scopes ?? [];
}

function validateTokenAdmin(requestEvent: RequestEvent) {
	const token = getDefinedToken(requestEvent);
	const claims = jwtDecode<JwtClaims>(token);
	if (claims.role !== 'admin') {
		error(403, 'Forbidden');
	}
}

function validateTokenExpiration(requestEvent: RequestEvent) {
	const token = getDefinedToken(requestEvent);

	// parse the token
	const { exp } = jwtDecode(token);

	if (exp) {
		const tokenExpirationDate = new Date(exp * 1000);

		if (tokenExpirationDate < new Date()) {
			console.log('User token expired, redirecting');
			redirect(303, getRedirectUrl(requestEvent));
		}
	}
}

async function validateTokenOrganization(requestEvent: RequestEvent) {
	const { params } = requestEvent;
	const { organizationId } = params;

	// Skip if not under an organization route
	if (organizationId === undefined) {
		return;
	}

	// Parse token and validate organization matches route param
	const token = getDefinedToken(requestEvent);
	const claims = jwtDecode<JwtClaims>(token);

	if (claims.organization_id === organizationId) {
		return;
	}

	// Attempt to generate new token with the organization id
	await refreshUserToken(requestEvent);
}

function getDefinedToken(requestEvent: RequestEvent) {
	const { cookies } = requestEvent;
	const token = cookies.get('app-session-token');

	if (!token) {
		console.log('User token not found, redirecting');
		redirect(303, getRedirectUrl(requestEvent));
	}

	return token;
}

function getRedirectUrl(requestEvent: RequestEvent) {
	const { url } = requestEvent;
	const currentPath = url.href;
	return `/signin?redirect=${encodeURIComponent(currentPath)}`;
}

function setTokenInCookies(requestEvent: RequestEvent, token: string) {
	const { cookies } = requestEvent;
	cookies.set('app-session-token', token, {
		path: '/',
		httpOnly: true,
		secure: true,
		sameSite: 'strict',
		maxAge: 60 * 60 * 24 * 30 // 30 days
	});
}
