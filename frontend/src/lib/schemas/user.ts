import { z } from 'zod';
import { API_TYPES } from './api';

const userAttributesSchema = z.object({
	name: z.string(),
	email: z.string()
});

export type UserRole = 'admin' | 'default';

export const userDataSchema = z.object({
	id: z.string(),
	type: z.literal(API_TYPES.user),
	attributes: z.object({
		name: z.string(),
		email: z.string()
	}),
	meta: z.object({
		logoDistributionUrl: z.string().optional(),
		tokenRevocation: z
			.object({
				lastIssuedAt: z.string().optional(),
				canRefresh: z.boolean().optional()
			})
			.optional()
	})
});

export type UserData = z.infer<typeof userDataSchema>;

export const getSelfUserResponseSchema = z.object({
	data: userDataSchema
});

export const createUserRequestSchema = z.object({
	data: z.object({
		attributes: userAttributesSchema.extend({
			password: z.string()
		})
	})
});

export const createUserResponseSchema = z.object({
	data: userDataSchema.extend({
		meta: z.object({
			token: z.string()
		})
	})
});

export const updateUserRequestSchema = z.object({
	data: z.object({
		id: z.string(),
		type: z.literal('user'),
		attributes: userAttributesSchema.partial().extend({
			password: z.string().optional(),
			logo: z.string().optional()
		})
	})
});

export const updateUserResponseSchema = z.object({
	data: userDataSchema
});

// Pagination schemas for users list
export const paginationLinksSchema = z.object({
	first: z.string().optional(),
	last: z.string().optional(),
	prev: z.string().optional(),
	next: z.string().optional()
});

export const getUsersParamsSchema = z.object({
	search: z.string().optional(),
	page: z.object({
		cursor: z.string().optional(),
		size: z.number().optional()
	})
});

export const getUsersResponseSchema = z.object({
	data: z.array(userDataSchema),
	links: paginationLinksSchema
});

export type GetUsersResponse = z.infer<typeof getUsersResponseSchema>;
export type PaginationLinks = z.infer<typeof paginationLinksSchema>;
