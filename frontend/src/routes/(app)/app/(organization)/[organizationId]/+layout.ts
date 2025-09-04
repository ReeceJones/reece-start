import type { LayoutLoad } from './$types';
import { error } from '@sveltejs/kit';
import { ApiError, get } from '$lib';
import { organizationResponseSchema } from '$lib/schemas/organization';

export const load: LayoutLoad = async ({ fetch, params, data }) => {
	const { organizationId } = params;

	try {
		const organization = await get(`/api/organizations/${organizationId}`, {
			fetch,
			responseSchema: organizationResponseSchema
		});

		return {
			...data,
			organization
		};
	} catch (apiError) {
		if (apiError instanceof ApiError) {
			error(apiError.code, apiError.message);
		}

		error(500, 'An unknown error ocurred processing your request, please try again later.');
	}
};
