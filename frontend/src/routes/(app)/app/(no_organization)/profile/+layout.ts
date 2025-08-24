import type { LayoutLoad } from './$types';
import { error } from '@sveltejs/kit';
import { ApiError, get } from '$lib';
import { getSelfUserResponseSchema } from '$lib/schemas/user';

export const load: LayoutLoad = async ({ fetch }) => {
	try {
		const user = await get('/api/users/me', {
			fetch,
			responseSchema: getSelfUserResponseSchema
		});

		return {
			user
		};
	} catch (apiError) {
		if (apiError instanceof ApiError) {
			error(apiError.code, apiError.message);
		}

		error(500, 'An unknown error ocurred processing your request, please try again later.');
	}
};
