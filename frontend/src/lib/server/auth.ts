import { getRequestEvent } from '$app/server';
import { ApiError, post } from '$lib/api';
import { API_TYPES } from '$lib/schemas/api';
import type { JwtClaims } from '$lib/schemas/jwt';
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

	// Validate token works for authenticated organization
	if (organizationId === undefined) {
		return;
	}

	await validateTokenOrganization(requestEvent);
}

export async function refreshUserToken(requestEvent: RequestEvent) {
	const { params, fetch } = requestEvent;
	const { organizationId } = params;

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
		httpOnly: false,
		secure: true,
		sameSite: 'strict'
	});
}
