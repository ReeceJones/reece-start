import { getRequestEvent } from '$app/server';
import { env } from '$env/dynamic/public';
import { redirect } from '@sveltejs/kit';

// Google OAuth configuration
const GOOGLE_OAUTH_SCOPE = 'email profile';
const GOOGLE_OAUTH_RESPONSE_TYPE = 'code';
const GOOGLE_OAUTH_URL = 'https://accounts.google.com/o/oauth2/v2/auth';

export function performGoogleOAuth(redirect: string | undefined) {
	const googleClientId = env.PUBLIC_GOOGLE_OAUTH_CLIENT_ID || '';

	performOAuth({
		clientId: googleClientId,
		successRedirectUrl: redirect ?? '/app',
		scope: GOOGLE_OAUTH_SCOPE,
		responseType: GOOGLE_OAUTH_RESPONSE_TYPE,
		accessType: 'offline',
		state: generateState(),
		oauthUrl: GOOGLE_OAUTH_URL,
		oauthCallbackUrl: getGoogleOAuthCallbackUrl()
	});
}

export function verifyOAuth({ state }: { state: string }) {
	const { cookies } = getRequestEvent();
	const storedState = cookies.get('oauth_state');
	const successRedirectUrl = cookies.get('oauth_success_redirect');

	if (!state || !successRedirectUrl) {
		return false;
	}

	if (state !== storedState) {
		return false;
	}

	return true;
}

export function deleteOAuthCookies() {
	const { cookies } = getRequestEvent();
	cookies.delete('oauth_state', { path: '/' });
	cookies.delete('oauth_success_redirect', { path: '/' });
}

function performOAuth({
	clientId,
	successRedirectUrl,
	scope,
	responseType,
	accessType,
	state,
	oauthUrl,
	oauthCallbackUrl
}: {
	clientId: string;
	successRedirectUrl: string;
	scope: string;
	responseType: string;
	accessType: 'offline';
	state: string;
	oauthUrl: string;
	oauthCallbackUrl: string;
}) {
	// generate the url
	const oauthRedirectUrl = generateOAuthUrl({
		clientId,
		oauthUrl,
		oauthCallbackUrl,
		scope,
		responseType,
		state,
		accessType
	});

	// save the challenge in the cookies, so we can verify it later
	const { cookies } = getRequestEvent();
	// store the challenge
	cookies.set('oauth_state', state, {
		path: '/',
		httpOnly: true,
		secure: true,
		sameSite: 'strict',
		maxAge: 60 * 5 // 5 minutes
	});
	// store the redirect url used on success
	cookies.set('oauth_success_redirect', successRedirectUrl, {
		path: '/',
		httpOnly: true,
		secure: true,
		sameSite: 'strict',
		maxAge: 60 * 5 // 5 minutes
	});

	redirect(302, oauthRedirectUrl);
}

function generateState(): string {
	const array = new Uint8Array(32);
	crypto.getRandomValues(array);
	return Array.from(array, (byte) => byte.toString(16).padStart(2, '0')).join('');
}

function getGoogleOAuthCallbackUrl(): string {
	const { url } = getRequestEvent();

	return `${url.origin}/oauth/google/callback`;
}

function generateOAuthUrl({
	clientId,
	oauthUrl,
	oauthCallbackUrl,
	scope,
	responseType,
	state,
	accessType
}: {
	clientId: string;
	oauthUrl: string;
	oauthCallbackUrl: string;
	scope: string;
	responseType: string;
	state: string;
	accessType: 'offline';
}) {
	const params = new URLSearchParams({
		client_id: clientId,
		redirect_uri: oauthCallbackUrl,
		response_type: responseType,
		scope: scope,
		state: state,
		access_type: accessType
	});

	return `${oauthUrl}?${params.toString()}`;
}
