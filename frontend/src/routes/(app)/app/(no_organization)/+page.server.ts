import { authenticate, stopImpersonatingUser } from '$lib/server/auth';
import { redirect } from '@sveltejs/kit';
import type { Actions, PageLoad } from './$types';

export const load: PageLoad = async () => {
	authenticate();
};

export const actions = {
	signout: async ({ cookies }) => {
		cookies.delete('app-session-token', { path: '/' });
		redirect(302, '/signin');
	},
	stopImpersonation: async (event) => {
		await stopImpersonatingUser(event);
		redirect(302, '/app');
	}
} satisfies Actions;
