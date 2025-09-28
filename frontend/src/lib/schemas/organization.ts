import { z } from 'zod';
import { API_TYPES } from './api';
import { addressSchema } from './address';

const organizationAttributesSchema = z.object({
	name: z.string().min(1).max(100),
	description: z.string().optional(),
	logo: z.string().optional(),
	address: addressSchema,
	currency: z.string(),
	locale: z.string(),
	contactEmail: z.email().or(z.literal('')).optional(),
	contactPhone: z.string().or(z.literal('')).optional(),
	websiteUrl: z.url().or(z.literal('')).optional()
});

const organizationOnboardingAttributesSchema = z.object({
	entityType: z.string(),
	residingCountry: z.string(),
	registeredBusinessName: z.string().optional(),
	structure: z.string().optional(),
	firstName: z.string().optional(),
	lastName: z.string().optional()
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
		logoDistributionUrl: z.string().optional()
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

export const organizationFormSchema = z.object({
	name: z.string().min(1).max(100),
	description: z.string(),
	logo: z.custom<FileList>().nullish(),
	contactEmail: z.email().or(z.literal('')).optional(),
	contactPhone: z.string().or(z.literal('')).optional(),
	websiteUrl: z.url().or(z.literal('')).optional(),
	addressLine1: z.string().optional(),
	addressLine2: z.string().optional(),
	addressCity: z.string().optional(),
	addressStateOrProvince: z.string().optional(),
	addressZip: z.string().optional(),
	addressCountry: z.string().optional(),
	currency: z.string(),
	locale: z.string()
});

export type OrganizationFormData = z.infer<typeof organizationFormSchema>;

export const createOrganizationFormSchema = organizationFormSchema.extend({
	entityType: z.string(),
	residingCountry: z.string(),
	registeredBusinessName: z.string().optional(),
	firstName: z.string().optional(),
	lastName: z.string().optional()
});

export function getFormDataFromOrganization(organization: Organization): OrganizationFormData {
	return {
		name: organization.data.attributes.name,
		description: organization.data.attributes.description || '',
		logo: undefined,
		contactEmail: organization.data.attributes.contactEmail,
		contactPhone: organization.data.attributes.contactPhone,
		websiteUrl: organization.data.attributes.websiteUrl,
		addressLine1: organization.data.attributes.address.line1,
		addressLine2: organization.data.attributes.address.line2,
		addressCity: organization.data.attributes.address.city,
		addressStateOrProvince: organization.data.attributes.address.stateOrProvince,
		addressZip: organization.data.attributes.address.zip,
		addressCountry: organization.data.attributes.address.country,
		currency: organization.data.attributes.currency,
		locale: organization.data.attributes.locale
	};
}
