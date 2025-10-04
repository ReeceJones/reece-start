import {
	createOrganizationFormSchema,
	type CreateOrganizationFormData
} from './schemas/organization';

export function isBasicInformationValid(state: CreateOrganizationFormData) {
	const basicInformationSchema = createOrganizationFormSchema.pick({
		name: true,
		description: true,
		logo: true
	});

	const result = basicInformationSchema.safeParse(state);

	return result.success;
}

export function isContactInformationValid(state: CreateOrganizationFormData) {
	const contactInformationSchema = createOrganizationFormSchema.pick({
		contactEmail: true,
		contactPhone: true,
		contactPhoneCountry: true
	});

	const result = contactInformationSchema.safeParse(state);

	return result.success;
}

export function isAddressValid(state: CreateOrganizationFormData) {
	const addressSchema = createOrganizationFormSchema.pick({
		addressLine1: true,
		addressLine2: true,
		addressCity: true,
		addressStateOrProvince: true,
		addressZip: true,
		addressCountry: true
	});

	const result = addressSchema.safeParse(state);

	return result.success;
}

export function isBusinessDetailsValid(state: CreateOrganizationFormData) {
	const businessDetailsSchema = createOrganizationFormSchema.pick({
		entityType: true,
		locale: true
	});

	const result = businessDetailsSchema.safeParse(state);

	return result.success;
}

export function isAllInformationValid(state: CreateOrganizationFormData) {
	const result = createOrganizationFormSchema.safeParse(state);

	return result.success;
}

export function getValidationErrors(state: CreateOrganizationFormData) {
	const result = createOrganizationFormSchema.safeParse(state);

	if (!result.success) {
		return result.error.issues.map((issue) => ({
			field: issue.path.join('.'),
			message: issue.message
		}));
	}

	return [];
}

export function formatUrl(url: string) {
	if (!url) {
		return url;
	}

	if (!url.startsWith('https://') || !url.startsWith('http://')) {
		return `https://${url}`;
	}

	return url;
}
