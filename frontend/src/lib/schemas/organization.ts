import { z } from 'zod';

// Request schemas
export const createOrganizationRequestSchema = z.object({
	data: z.object({
		type: z.literal('organization'),
		attributes: z.object({
			name: z.string().min(1).max(100),
			description: z.string().optional(),
			logo: z.string().optional()
		})
	})
});

export const updateOrganizationRequestSchema = z.object({
	data: z.object({
		type: z.literal('organization'),
		attributes: z.object({
			name: z.string().min(1).max(100).optional(),
			description: z.string().optional(),
			logo: z.string().optional()
		})
	})
});

// Response schemas
const organizationSchema = z.object({
	id: z.string(),
	type: z.literal('organization'),
	attributes: z.object({
		name: z.string(),
		description: z.string().optional()
	}),
	meta: z.object({
		logoDistributionUrl: z.string().optional()
	})
});

export const organizationResponseSchema = z.object({
	data: organizationSchema
});

export type Organization = z.infer<typeof organizationResponseSchema>;

export const organizationsResponseSchema = z.object({
	data: z.array(organizationSchema)
});
