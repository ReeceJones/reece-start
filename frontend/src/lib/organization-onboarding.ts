import {
	createOrganizationFormSchema,
	type CreateOrganizationFormData
} from './schemas/organization';

export function isBasicInformationValid(state: CreateOrganizationFormData) {
	return !!state.name;
}

export function isAddressValid(state: CreateOrganizationFormData) {
	return (
		!!state.addressLine1 &&
		!!state.addressCity &&
		!!state.addressStateOrProvince &&
		!!state.addressZip &&
		!!state.addressCountry
	);
}

export function isBusinessDetailsValid(state: CreateOrganizationFormData) {
	return !!state.entityType && !!state.locale;
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
