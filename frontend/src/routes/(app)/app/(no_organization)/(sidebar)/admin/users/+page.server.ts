import { impersonateUserFormSchema } from '$lib/schemas/user.server';
import { authenticate } from '$lib/server/auth';
import { impersonateUser } from '$lib/server/auth';
import { isParseSuccess, parseFormData } from '$lib/server/schema';
import { redirect, fail, type Actions } from '@sveltejs/kit';

export const load = async () => {
	authenticate();
};

export const actions: Actions = {
	impersonate: async (event) => {
		authenticate();

		const formData = await parseFormData(event.request, impersonateUserFormSchema);

		if (!isParseSuccess(formData)) {
			return fail(400, { message: 'Missing impersonated user id' });
		}

		await impersonateUser(event, formData.impersonatedUserId);
		throw redirect(303, '/app');
	}
};
