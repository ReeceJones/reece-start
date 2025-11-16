import { z } from 'zod';

export const API_TYPES = {
	user: 'user',
	organization: 'organization',
	token: 'token',
	organizationMembership: 'organization-membership',
	organizationInvitation: 'organization-invitation',
	stripeAccountLink: 'stripe-account-link'
} as const;

export const apiErrorSchema = z.object({
	message: z.string()
});

export type ApiErrorResponse = z.infer<typeof apiErrorSchema>;
