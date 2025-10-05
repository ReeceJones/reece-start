import { z } from 'zod';
import { API_TYPES } from './api';

// Request to create a Stripe onboarding link for an organization
export const createStripeOnboardingLinkRequestSchema = z.object({
	data: z.object({
		type: z.literal(API_TYPES.stripeAccountLink),
		relationships: z.object({
			organization: z.object({
				data: z.object({
					id: z.string(),
					type: z.literal(API_TYPES.organization)
				})
			})
		})
	})
});

// Response for Stripe account link creation
export const stripeAccountLinkAttributesSchema = z.object({
	url: z.string(),
	expiresAt: z.iso.datetime(),
	livemode: z.boolean(),
	accountId: z.string(),
	createdAt: z.iso.datetime()
});

export const stripeAccountLinkDataSchema = z.object({
	type: z.literal(API_TYPES.stripeAccountLink),
	attributes: stripeAccountLinkAttributesSchema
});

export const createStripeOnboardingLinkResponseSchema = z.object({
	data: stripeAccountLinkDataSchema
});

export type CreateStripeOnboardingLinkResponse = z.infer<
	typeof createStripeOnboardingLinkResponseSchema
>;
