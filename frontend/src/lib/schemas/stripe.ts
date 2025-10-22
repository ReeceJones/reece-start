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

// Subscription schemas
export const subscriptionAttributesSchema = z.object({
	plan: z.enum(['free', 'pro']),
	billingPeriodStart: z.string().nullable(),
	billingPeriodEnd: z.string().nullable(),
	billingAmount: z.number()
});

export const subscriptionDataSchema = z.object({
	type: z.literal('subscription'),
	id: z.string().optional(),
	attributes: subscriptionAttributesSchema
});

export const getSubscriptionResponseSchema = z.object({
	data: subscriptionDataSchema
});

export type GetSubscriptionResponse = z.infer<typeof getSubscriptionResponseSchema>;
export type SubscriptionAttributes = z.infer<typeof subscriptionAttributesSchema>;

// Checkout session schemas
export const createCheckoutSessionRequestSchema = z.object({
	successUrl: z.string().url(),
	cancelUrl: z.string().url()
});

export const checkoutSessionAttributesSchema = z.object({
	url: z.string().url()
});

export const checkoutSessionDataSchema = z.object({
	type: z.literal('checkout-session'),
	id: z.string(),
	attributes: checkoutSessionAttributesSchema
});

export const createCheckoutSessionResponseSchema = z.object({
	data: checkoutSessionDataSchema
});

export type CreateCheckoutSessionRequest = z.infer<typeof createCheckoutSessionRequestSchema>;
export type CreateCheckoutSessionResponse = z.infer<typeof createCheckoutSessionResponseSchema>;

// Billing portal session schemas
export const createBillingPortalSessionRequestSchema = z.object({
	returnUrl: z.string().url()
});

export const billingPortalSessionAttributesSchema = z.object({
	url: z.string().url()
});

export const billingPortalSessionDataSchema = z.object({
	type: z.literal('billing-portal-session'),
	id: z.string(),
	attributes: billingPortalSessionAttributesSchema
});

export const createBillingPortalSessionResponseSchema = z.object({
	data: billingPortalSessionDataSchema
});

export type CreateBillingPortalSessionRequest = z.infer<
	typeof createBillingPortalSessionRequestSchema
>;
export type CreateBillingPortalSessionResponse = z.infer<
	typeof createBillingPortalSessionResponseSchema
>;
