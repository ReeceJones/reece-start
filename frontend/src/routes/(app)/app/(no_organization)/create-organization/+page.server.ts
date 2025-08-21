import { fail, redirect } from '@sveltejs/kit';
import { ApiError, post } from '$lib';
import type { Actions } from './$types';
import { z } from 'zod';

const createOrganizationRequestSchema = z.object({
	data: z.object({
		type: z.literal('organization'),
		attributes: z.object({
			name: z.string().min(1).max(100),
			description: z.string().optional()
		})
	})
});

const createOrganizationResponseSchema = z.object({
	data: z.object({
		type: z.literal('organization'),
		id: z.string(),
		attributes: z.object({
			name: z.string(),
			description: z.string().optional()
		})
	})
});

export const actions = {
	default: async ({ request, fetch }) => {
		const data = await request.formData();

		const name = data.get('name') as string;
		const description = data.get('description') as string;
		if (!name) {
			return fail(400, { success: false, message: 'Please fill out all the fields correctly.' });
		}

		let organization: z.infer<typeof createOrganizationResponseSchema>;
		try {
			organization = await post(
				'/api/organizations',
				{
					data: {
						type: 'organization',
						attributes: {
							name,
							description
						}
					}
				},
				{
					fetch,
					responseSchema: createOrganizationResponseSchema,
					requestSchema: createOrganizationRequestSchema
				}
			);
		} catch (apiError) {
			if (apiError instanceof ApiError) {
				return fail(apiError.code, { success: false, message: apiError.message });
			}

			return fail(500, {
				success: false,
				message: 'An unknown error ocurred processing your request, please try again later.'
			});
		}

		redirect(302, `/app/${organization.data.id}`);
	}
} satisfies Actions;
