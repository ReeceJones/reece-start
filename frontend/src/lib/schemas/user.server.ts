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

export const impersonateUserFormSchema = z.object({
	impersonatedUserId: z.string()
});

export const signinFormSchema = z.object({
	email: z.string().min(1),
	password: z.string().min(1)
});

export const signupFormSchema = z.object({
	name: z.string().min(1),
	email: z.string().min(1),
	password: z.string().min(1)
});

export const updateUserProfileFormSchema = z.object({
	userId: z.string().min(1),
	name: z.string().min(1),
	logo: z
		.instanceof(File)
		.optional()
		.refine((file) => !file || file.size === 0 || file.size <= 3_000_000, {
			message: 'Logo must be less than 3MB.'
		})
});

export const updateUserSecurityFormSchema = z.object({
	userId: z.string().min(1),
	email: z.string().min(1),
	password: z.string().optional(),
	confirmPassword: z.string().optional()
});

export const googleOAuthCallbackFormSchema = z.object({
	code: z.string().min(1),
	state: z.string().min(1),
	redirect: z.string().optional()
});
