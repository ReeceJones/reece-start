import { authenticate } from '$lib/server/auth';
import { redirect } from '@sveltejs/kit';
import type { Actions, PageLoad } from './$types';

export const load: PageLoad = async () => {
	authenticate();
};

export const actions = {
	signout: async ({ cookies }) => {
		cookies.delete('app-session-token', { path: '/' });
		redirect(302, '/signin');
	}
} satisfies Actions;
