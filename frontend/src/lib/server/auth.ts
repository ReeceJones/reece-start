import { getRequestEvent } from '$app/server';
import { redirect, type RequestEvent } from '@sveltejs/kit';
import { jwtDecode } from 'jwt-decode';

export function authenticate() {
	const requestEvent = getRequestEvent();
	performAuthenticationCheck(requestEvent);
}

export function performAuthenticationCheck(requestEvent: RequestEvent) {
	const { url, cookies } = requestEvent;

	if (url.pathname.startsWith('/app')) {
		// get the token from the cookies and check if the exp claim is still valid
		const token = cookies.get('app-session-token');
		const currentPath = url.href;
		const redirectUrl = `/signin?redirect=${encodeURIComponent(currentPath)}`;

		if (!token) {
			console.log('User token not found, redirecting');
			redirect(303, redirectUrl);
		}

		// parse the token
		const { exp } = jwtDecode(token);

		if (exp) {
			const tokenExpirationDate = new Date(exp * 1000);

			if (tokenExpirationDate < new Date()) {
				console.log('User token expired, redirecting');
				redirect(303, redirectUrl);
			}
		}
	}
}
