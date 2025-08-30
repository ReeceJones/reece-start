import { z } from 'zod';
import { API_TYPES } from './api';
import { organizationMembershipRole } from './organization-membership';

const organizationInvitationAttributesSchema = z.object({
	email: z.string(),
	role: organizationMembershipRole
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
		attributes: organizationInvitationAttributesSchema,
		relationships: organizationInvitationRelationshipsSchema.omit({
			invitingUser: true
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
