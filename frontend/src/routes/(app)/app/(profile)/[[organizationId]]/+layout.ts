import type { LayoutLoad } from './$types';
import { error } from '@sveltejs/kit';
import { ApiError, get } from '$lib';
import { getSelfUserResponseSchema } from '$lib/schemas/user';
import { organizationResponseSchema } from '$lib/schemas/organization';

export const load: LayoutLoad = async ({ fetch, params }) => {
	try {
		const user = await get('/api/users/me', {
			fetch,
			responseSchema: getSelfUserResponseSchema
		});

		let organization = undefined;
		if (params.organizationId) {
			organization = await get(`/api/organizations/${params.organizationId}`, {
				fetch,
				responseSchema: organizationResponseSchema
			});
		}

		return {
			user,
			organization
		};
	} catch (apiError) {
		if (apiError instanceof ApiError) {
			error(apiError.code, apiError.message);
		}

		error(500, 'An unknown error ocurred processing your request, please try again later.');
	}
};
