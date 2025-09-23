import { z } from 'zod';
import { API_TYPES } from './api';

export const createAuthenticatedUserTokenRequestSchema = z.object({
	data: z.object({
		type: z.literal(API_TYPES.token),
		relationships: z.object({
			impersonatedUser: z
				.object({
					data: z.object({
						id: z.string(),
						type: z.literal(API_TYPES.user)
					})
				})
				.optional(),
			organization: z
				.object({
					data: z.object({
						id: z.string(),
						type: z.literal(API_TYPES.organization)
					})
				})
				.optional()
		}),
		meta: z
			.object({
				stopImpersonating: z.boolean().optional()
			})
			.optional()
	})
});

export const createAuthenticatedUserTokenResponseSchema = z.object({
	data: z.object({
		type: z.literal(API_TYPES.token),
		relationships: z.object({
			organization: z
				.object({
					data: z.object({
						id: z.string(),
						type: z.literal(API_TYPES.organization)
					})
				})
				.optional()
				.nullable(),
			impersonatedUser: z
				.object({
					data: z.object({
						id: z.string(),
						type: z.literal(API_TYPES.user)
					})
				})
				.optional()
				.nullable()
		}),
		meta: z.object({
			token: z.string()
		})
	})
});
