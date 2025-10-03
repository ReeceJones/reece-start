import { z } from 'zod';

export const addressSchema = z.object({
	line1: z.string(),
	line2: z.string().optional(),
	city: z.string(),
	stateOrProvince: z.string().optional(),
	zip: z.string(), // rename this to postalCode
	country: z.string().min(2)
});
