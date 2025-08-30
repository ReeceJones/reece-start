import { z } from 'zod';
import { API_TYPES } from './api';

// Request schemas
export const createOrganizationRequestSchema = z.object({
	data: z.object({
		type: z.literal(API_TYPES.organization),
		attributes: z.object({
			name: z.string().min(1).max(100),
			description: z.string().optional(),
			logo: z.string().optional()
		})
	})
});

export const updateOrganizationRequestSchema = z.object({
	data: z.object({
		type: z.literal(API_TYPES.organization),
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
	type: z.literal(API_TYPES.organization),
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

export const organizationFormSchema = z.object({
	name: z.string().min(1).max(100),
	description: z.string(),
	logo: z.custom<FileList>().nullish()
});

export type OrganizationFormData = z.infer<typeof organizationFormSchema>;

export function getFormDataFromOrganization(organization: Organization): OrganizationFormData {
	return {
		name: organization.data.attributes.name,
		description: organization.data.attributes.description || '',
		logo: undefined
	};
}
