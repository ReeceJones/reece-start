import { z } from 'zod';
import { userDataSchema } from './user';
import { API_TYPES } from './api';

export const organizationMembershipRole = z.enum(['admin', 'member']);

export type OrganizationMembershipRole = z.infer<typeof organizationMembershipRole>;

const organizationMembershipAttributesSchema = z.object({
	role: organizationMembershipRole
});

const organizationMembershipRelationshipsSchema = z.object({
	user: z.object({
		data: z.object({
			id: z.string(),
			type: z.literal(API_TYPES.user)
		})
	}),
	organization: z.object({
		data: z.object({
			id: z.string(),
			type: z.literal(API_TYPES.organization)
		})
	})
});

const organizationMembershipIncludedSchema = z.array(userDataSchema);

const organizationMembershipDataSchema = z.object({
	id: z.string(),
	type: z.literal(API_TYPES.organizationMembership),
	attributes: organizationMembershipAttributesSchema,
	relationships: organizationMembershipRelationshipsSchema
});

export const getOrganizationMembershipResponseSchema = z.object({
	data: organizationMembershipDataSchema,
	included: organizationMembershipIncludedSchema
});

export const getOrganizationMembershipsResponseSchema = z.object({
	data: z.array(organizationMembershipDataSchema),
	included: organizationMembershipIncludedSchema
});

export const getOrganizationMembershipsQuerySchema = z.object({
	organizationId: z.string()
});

export const updateOrganizationMembershipResponseSchema =
	getOrganizationMembershipResponseSchema.omit({
		included: true
	});

export const updateOrganizationMembershipRequestSchema = z.object({
	data: z.object({
		type: z.literal(API_TYPES.organizationMembership),
		id: z.string(),
		attributes: organizationMembershipAttributesSchema.partial()
	})
});

export const updateOrganizationMembershipFormSchema = z.object({
	role: organizationMembershipRole
});
