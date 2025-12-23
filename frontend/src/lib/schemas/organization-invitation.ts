import { z } from 'zod';
import { API_TYPES } from './api';
import { organizationMembershipRole } from './organization-membership';

const organizationInvitationStatus = z.enum([
	'pending',
	'accepted',
	'declined',
	'expired',
	'revoked'
]);

const organizationInvitationAttributesSchema = z.object({
	email: z.string(),
	role: organizationMembershipRole,
	status: organizationInvitationStatus
});

const organizationInvitationRelationshipsSchema = z.object({
	organization: z.object({
		data: z.object({
			id: z.string(),
			type: z.literal(API_TYPES.organization)
		})
	}),
	invitingUser: z.object({
		data: z.object({
			id: z.string(),
			type: z.literal(API_TYPES.user)
		})
	})
});

const organizationInvitationDataSchema = z.object({
	id: z.string(),
	type: z.literal(API_TYPES.organizationInvitation),
	attributes: organizationInvitationAttributesSchema,
	relationships: organizationInvitationRelationshipsSchema
});

export type OrganizationInvitation = z.infer<typeof organizationInvitationDataSchema>;

export const inviteToOrganizationRequestSchema = z.object({
	data: z.object({
		type: z.literal(API_TYPES.organizationInvitation),
		attributes: organizationInvitationAttributesSchema.pick({
			email: true,
			role: true
		}),
		relationships: organizationInvitationRelationshipsSchema.pick({
			organization: true
		})
	})
});

export const inviteToOrganizationResponseSchema = z.object({
	data: organizationInvitationDataSchema
});

export const inviteToOrganizationFormSchema = z.object({
	email: z.string(),
	role: organizationMembershipRole
});

export const getOrganizationInvitationsQuerySchema = z.object({
	organizationId: z.string()
});

export const getOrganizationInvitationsResponseSchema = z.object({
	data: z.array(organizationInvitationDataSchema)
});

export const deleteInvitationFormSchema = z.object({
	invitationId: z.string()
});

const organizationInvitationIncludedSchema = z.array(
	z.discriminatedUnion('type', [
		z.object({
			id: z.string(),
			type: z.literal(API_TYPES.organization),
			attributes: z.object({
				name: z.string(),
				description: z.string().optional()
			}),
			meta: z.object({
				logoDistributionUrl: z.string().optional()
			})
		}),
		z.object({
			id: z.string(),
			type: z.literal(API_TYPES.user),
			attributes: z.object({
				name: z.string(),
				email: z.string()
			}),
			meta: z.object({
				logoDistributionUrl: z.string().optional()
			})
		})
	])
);

export const getOrganizationInvitationResponseSchema = z.object({
	data: organizationInvitationDataSchema,
	included: organizationInvitationIncludedSchema
});

export const acceptOrganizationInvitationRequestSchema = z.object({
	data: z.object({
		id: z.string(),
		type: z.literal(API_TYPES.organizationInvitation)
	})
});

export const acceptOrganizationInvitationResponseSchema = z.object({
	data: organizationInvitationDataSchema,
	included: organizationInvitationIncludedSchema
});

export const declineOrganizationInvitationRequestSchema = z.object({
	data: z.object({
		id: z.string(),
		type: z.literal(API_TYPES.organizationInvitation)
	})
});

export const declineOrganizationInvitationResponseSchema = z.object({
	data: organizationInvitationDataSchema,
	included: organizationInvitationIncludedSchema
});
