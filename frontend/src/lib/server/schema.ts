import { fail, isActionFailure, type ActionFailure } from '@sveltejs/kit';
import { z } from 'zod';

export async function parseFormData<T>(
	request: Request,
	schema: z.ZodSchema<T>
): Promise<T | ActionFailure<{ success: boolean; message: string }>> {
	const formData = await request.formData();
	const result = schema.safeParse(Object.fromEntries(formData));

	if (!result.success) {
		return fail(400, {
			success: false,
			message: result.error.message
		});
	}

	return result.data;
}

export function isParseSuccess<T>(
	result: T | ActionFailure<{ success: boolean; message: string }>
): result is T {
	return !isActionFailure(result);
}
