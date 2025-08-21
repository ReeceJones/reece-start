import type { PageLoad } from './$types';
import { error } from '@sveltejs/kit';
import { z } from 'zod';
import { ApiError, get } from '$lib';

const getSelfUserResponseSchema = z.object({
	data: z.object({
		id: z.string(),
		type: z.literal('user'),
		attributes: z.object({
			name: z.string(),
			email: z.string()
		})
	})
});

export const load: PageLoad = async ({ fetch }) => {
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
