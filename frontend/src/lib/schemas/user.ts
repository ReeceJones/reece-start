import { z } from 'zod';
import { API_TYPES } from './api';

const userAttributesSchema = z.object({
	name: z.string(),
	email: z.string()
});

export const userDataSchema = z.object({
	id: z.string(),
	type: z.literal(API_TYPES.user),
	attributes: z.object({
		name: z.string(),
		email: z.string()
	}),
	meta: z.object({
		logoDistributionUrl: z.string().optional()
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
