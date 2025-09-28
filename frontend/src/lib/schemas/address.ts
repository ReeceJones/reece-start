import { z } from 'zod';

export const addressSchema = z.object({
	line1: z.string().optional(),
	line2: z.string().optional(),
	city: z.string().optional(),
	stateOrProvince: z.string().optional(),
	zip: z.string().optional(),
	country: z.string().optional()
});
