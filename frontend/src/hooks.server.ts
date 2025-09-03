import { performAuthenticationCheck } from '$lib/server/auth';
import { type Handle } from '@sveltejs/kit';

export const handle: Handle = async ({ event, resolve }) => {
	await performAuthenticationCheck(event);
	const response = await resolve(event);
	return response;
};
