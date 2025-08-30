import { z } from 'zod';
import { API_TYPES } from './api';

const organizationInvitationAttributesSchema = z.object({
	email: z.string()
});

const organizationInvitationRelationshipsSchema = z.object({
	organization: z.object({
		data: z.object({
			id: z.string(),
			type: z.literal(API_TYPES.organization)
		})
	})
});

const organizationInvitationDataSchema = z.object({
	id: z.string(),
	type: z.literal(API_TYPES.organizationInvitation),
	attributes: organizationInvitationAttributesSchema,
	relationships: organizationInvitationRelationshipsSchema
});

export const inviteToOrganizationRequestSchema = z.object({
	data: z.object({
		type: z.literal(API_TYPES.organizationInvitation),
		attributes: organizationInvitationAttributesSchema,
		relationships: organizationInvitationRelationshipsSchema
	})
});

export const inviteToOrganizationResponseSchema = z.object({
	data: organizationInvitationDataSchema
});
