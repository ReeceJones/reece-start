import { z } from 'zod';
import { API_TYPES } from './api';
import { addressSchema } from './address';
import { State } from 'country-state-city';

const organizationAttributesSchema = z.object({
	name: z.string().min(1).max(100),
	description: z.string().optional(),
	logo: z.string().optional(),
	address: addressSchema,
	locale: z.string(),
	contactEmail: z.email(),
	contactPhone: z.string().min(1),
	contactPhoneCountry: z.string().min(2)
});

const organizationOnboardingAttributesSchema = z.object({
	entityType: z.string()
});

// Request schemas
export const createOrganizationRequestSchema = z.object({
	data: z.object({
		type: z.literal(API_TYPES.organization),
		attributes: organizationAttributesSchema.extend(organizationOnboardingAttributesSchema.shape)
	})
});

export const updateOrganizationRequestSchema = z.object({
	data: z.object({
		type: z.literal(API_TYPES.organization),
		attributes: organizationAttributesSchema.partial()
	})
});

// Response schemas
export const organizationDataSchema = z.object({
	id: z.string(),
	type: z.literal(API_TYPES.organization),
	attributes: organizationAttributesSchema,
	meta: z.object({
		logoDistributionUrl: z.string().optional(),
		onboardingStatus: z.enum(['pending', 'completed', 'in_progress']),
		stripe: z.object({
			hasPendingRequirements: z.boolean().optional(),
			onboardingStatus: z.enum([
				'pending',
				'completed',
				'missing_requirements',
				'missing_capabilities'
			])
		})
	})
});

export type OrganizationData = z.infer<typeof organizationDataSchema>;

export const organizationResponseSchema = z.object({
	data: organizationDataSchema
});

export type Organization = z.infer<typeof organizationResponseSchema>;

export const organizationsResponseSchema = z.object({
	data: z.array(organizationDataSchema)
});

export const organizationFormSchema = z
	.object({
		name: z.string().min(1).max(100),
		description: z.string(),
		logo: z.custom<FileList>().nullish(),
		contactEmail: z.email(),
		contactPhone: z.string().min(1),
		contactPhoneCountry: z.string().min(2),
		addressLine1: z.string().min(1),
		addressLine2: z.string().optional(),
		addressCity: z.string().min(1),
		addressStateOrProvince: z.string().optional(),
		addressZip: z.string().min(1),
		addressCountry: z.string().min(2),
		locale: z.string()
	})
	.superRefine((data, ctx) => {
		// If phone number is provided, country code must also be provided
		if (
			data.contactPhone &&
			data.contactPhone.trim() !== '' &&
			(!data.contactPhoneCountry || data.contactPhoneCountry.trim() === '')
		) {
			ctx.addIssue({
				code: 'custom',
				message: 'Country code is required when phone number is provided',
				path: ['contactPhoneCountry']
			});
		}
	})
	.superRefine((data, ctx) => {
		// If country requires state and state is not provided
		if (!data.addressCountry) {
			return;
		}

		const states = State.getStatesOfCountry(data.addressCountry);

		if (states.length > 0) {
			if (!data.addressStateOrProvince) {
				ctx.addIssue({
					code: 'custom',
					message: 'State/Province is required for ' + data.addressCountry,
					path: ['addressStateOrProvince']
				});
			}
		}
	});

export type OrganizationFormData = z.infer<typeof organizationFormSchema>;

export const createOrganizationFormSchema = organizationFormSchema.extend({
	entityType: z.string()
});

export type CreateOrganizationFormData = z.infer<typeof createOrganizationFormSchema>;

export function getFormDataFromOrganization(organization: Organization): OrganizationFormData {
	return {
		name: organization.data.attributes.name,
		description: organization.data.attributes.description || '',
		logo: undefined,
		contactEmail: organization.data.attributes.contactEmail,
		contactPhone: organization.data.attributes.contactPhone,
		contactPhoneCountry: organization.data.attributes.contactPhoneCountry,
		addressLine1: organization.data.attributes.address.line1,
		addressLine2: organization.data.attributes.address.line2,
		addressCity: organization.data.attributes.address.city,
		addressStateOrProvince: organization.data.attributes.address.stateOrProvince,
		addressZip: organization.data.attributes.address.zip,
		addressCountry: organization.data.attributes.address.country,
		locale: organization.data.attributes.locale
	};
}

// Stripe dashboard link schema
export const createStripeDashboardLinkResponseSchema = z.object({
	data: z.object({
		type: z.literal('stripe-dashboard-link'),
		attributes: z.object({
			url: z.string()
		})
	})
});

export type CreateStripeDashboardLinkResponse = z.infer<
	typeof createStripeDashboardLinkResponseSchema
>;
