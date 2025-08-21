import { redirect, type Handle } from '@sveltejs/kit';
import { jwtDecode } from 'jwt-decode';

export const handle: Handle = async ({ event, resolve }) => {
	if (event.url.pathname.startsWith('/app')) {
		// get the token from the cookies and check if the exp claim is still valid
		const token = event.cookies.get('app-session-token');

		if (!token) {
			console.log('User token not found, redirecting');
			redirect(303, '/signin');
		}

		// parse the token
		const { exp } = jwtDecode(token);

		if (exp) {
			const tokenExpirationDate = new Date(exp * 1000);

			if (tokenExpirationDate < new Date()) {
				console.log('User token expired, redirecting');
				redirect(303, '/signin');
			}
		}
	}

	const response = await resolve(event);
	return response;
};
