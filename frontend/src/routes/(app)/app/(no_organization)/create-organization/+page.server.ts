// import { fail, redirect } from '@sveltejs/kit';
// import { ApiError, post } from '$lib';
// import { base64Encode } from '$lib/base64';
// import {
// 	createOrganizationRequestSchema,
// 	organizationResponseSchema
// } from '$lib/schemas/organization';
// import type { Actions } from './$types';
// import { z } from 'zod';

// export const actions = {
// 	default: async ({ request, fetch }) => {
// 		const data = await request.formData();

// 		const name = data.get('name') as string;
// 		const description = data.get('description') as string;
// 		const logoFile = data.get('logo') as File | undefined;

// 		if (!name) {
// 			return fail(400, { success: false, message: 'Please fill out all the fields correctly.' });
// 		}

// 		// Convert logo file to base64 if provided
// 		let logoBase64 = undefined;
// 		if (logoFile && logoFile.size > 0) {
// 			try {
// 				const arrayBuffer = await logoFile.arrayBuffer();
// 				logoBase64 = base64Encode(arrayBuffer);
// 			} catch (error) {
// 				return fail(400, { success: false, message: 'Error processing logo file.' });
// 			}
// 		}

// 		let organization: z.infer<typeof organizationResponseSchema>;
// 		try {
// 			organization = await post(
// 				'/api/organizations',
// 				{
// 					data: {
// 						type: 'organization',
// 						attributes: {
// 							name,
// 							description,
// 							...(logoBase64 && { logo: logoBase64 })
// 						}
// 					}
// 				},
// 				{
// 					fetch,
// 					responseSchema: organizationResponseSchema,
// 					requestSchema: createOrganizationRequestSchema
// 				}
// 			);
// 		} catch (apiError) {
// 			if (apiError instanceof ApiError) {
// 				return fail(apiError.code, { success: false, message: apiError.message });
// 			}

// 			return fail(500, {
// 				success: false,
// 				message: 'An unknown error ocurred processing your request, please try again later.'
// 			});
// 		}

// 		redirect(302, `/app/${organization.data.id}`);
// 	}
// } satisfies Actions;
