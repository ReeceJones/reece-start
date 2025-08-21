import { get } from '$lib';
import type { PageLoad } from './$types';
import { z } from 'zod';

const getOrganizationsResponseSchema = z.object({
	data: z.array(
		z.object({
			id: z.string(),
			type: z.literal('organization'),
			attributes: z.object({
				name: z.string(),
				description: z.string().optional()
			})
		})
	)
});

export const load: PageLoad = async ({ fetch }) => {
	const organizations = await get('/api/organizations', {
		fetch,
		responseSchema: getOrganizationsResponseSchema
	});

	return {
		organizations
	};
};
