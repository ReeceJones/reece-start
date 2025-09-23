import { authenticate } from '$lib/server/auth';
import { impersonateUser } from '$lib/server/auth';
import { redirect, fail, type Actions } from '@sveltejs/kit';

export const load = async () => {
	authenticate();
};

export const actions: Actions = {
	impersonate: async (event) => {
		authenticate();

		const formData = await event.request.formData();
		const userId = String(formData.get('impersonatedUserId') ?? '');

		if (!userId) {
			return fail(400, { message: 'Missing impersonated user id' });
		}

		await impersonateUser(event, userId);
		throw redirect(303, '/app');
	}
};
