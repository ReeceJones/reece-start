import { authenticate } from '$lib/server/auth';
import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

// TODO: This should be a form action, not a page load
export const load: PageServerLoad = async ({ cookies }) => {
	authenticate();

	cookies.delete('app-session-token', { path: '/' });

	redirect(302, '/signin');
};
