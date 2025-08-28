import { get } from '$lib';
import type { PageLoad } from './$types';
import { z } from 'zod';
import { organizationsResponseSchema } from '$lib/schemas/organization';

export const load: PageLoad = async ({ fetch }) => {
	const organizations = await get('/api/organizations', {
		fetch,
		responseSchema: organizationsResponseSchema
	});

	return {
		organizations
	};
};
