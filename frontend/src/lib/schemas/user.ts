import { z } from 'zod';

export const getSelfUserResponseSchema = z.object({
	data: z.object({
		id: z.string(),
		type: z.literal('user'),
		attributes: z.object({
			name: z.string(),
			email: z.string()
		}),
		meta: z.object({
			logoDistributionUrl: z.string()
		})
	})
});

export type User = z.infer<typeof getSelfUserResponseSchema>;

export const updateUserRequestSchema = z.object({
	data: z.object({
		id: z.string(),
		type: z.literal('user'),
		attributes: z.object({
			name: z.string().optional(),
			email: z.string().optional(),
			password: z.string().optional(),
			logo: z.string().optional()
		})
	})
});

export const updateUserResponseSchema = z.object({
	data: z.object({
		id: z.string(),
		type: z.literal('user'),
		attributes: z.object({
			name: z.string(),
			email: z.string()
		})
	})
});
