import { describe, it, expect } from 'vitest';
import { fail } from '@sveltejs/kit';
import { z } from 'zod';
import { parseFormData, isParseSuccess } from './schema';

describe('schema', () => {
	describe('parseFormData', () => {
		it('should parse valid form data', async () => {
			const schema = z.object({
				name: z.string(),
				email: z.string().email(),
				age: z.string().transform((val) => parseInt(val, 10))
			});

			const formData = new FormData();
			formData.append('name', 'John Doe');
			formData.append('email', 'john@example.com');
			formData.append('age', '30');

			const request = new Request('http://localhost/test', {
				method: 'POST',
				body: formData
			});

			const result = await parseFormData(request, schema);

			expect(isParseSuccess(result)).toBe(true);
			if (isParseSuccess(result)) {
				expect(result.name).toBe('John Doe');
				expect(result.email).toBe('john@example.com');
				expect(result.age).toBe(30);
			}
		});

		it('should return ActionFailure for invalid form data', async () => {
			const schema = z.object({
				name: z.string(),
				email: z.string().email()
			});

			const formData = new FormData();
			formData.append('name', 'John Doe');
			formData.append('email', 'invalid-email');

			const request = new Request('http://localhost/test', {
				method: 'POST',
				body: formData
			});

			const result = await parseFormData(request, schema);

			expect(isParseSuccess(result)).toBe(false);
			if (!isParseSuccess(result)) {
				expect(result.status).toBe(400);
				expect(result.data.success).toBe(false);
				expect(result.data.message).toBeTruthy();
			}
		});

		it('should return ActionFailure for missing required fields', async () => {
			const schema = z.object({
				name: z.string(),
				email: z.string().email()
			});

			const formData = new FormData();
			formData.append('name', 'John Doe');
			// Missing email

			const request = new Request('http://localhost/test', {
				method: 'POST',
				body: formData
			});

			const result = await parseFormData(request, schema);

			expect(isParseSuccess(result)).toBe(false);
			if (!isParseSuccess(result)) {
				expect(result.status).toBe(400);
				expect(result.data.success).toBe(false);
			}
		});

		it('should handle empty form data', async () => {
			const schema = z.object({
				name: z.string().optional(),
				email: z.string().email().optional()
			});

			const formData = new FormData();

			const request = new Request('http://localhost/test', {
				method: 'POST',
				body: formData
			});

			const result = await parseFormData(request, schema);

			expect(isParseSuccess(result)).toBe(true);
			if (isParseSuccess(result)) {
				expect(result.name).toBeUndefined();
				expect(result.email).toBeUndefined();
			}
		});

		it('should handle nested schema', async () => {
			const schema = z.object({
				user: z.object({
					name: z.string(),
					email: z.string().email()
				})
			});

			const formData = new FormData();
			formData.append('user', JSON.stringify({ name: 'John', email: 'john@example.com' }));

			const request = new Request('http://localhost/test', {
				method: 'POST',
				body: formData
			});

			const result = await parseFormData(request, schema);

			// Note: FormData with nested objects as JSON strings won't parse correctly
			// This test demonstrates the limitation
			expect(isParseSuccess(result)).toBe(false);
		});

		it('should handle number validation', async () => {
			const schema = z.object({
				age: z.string().transform((val) => {
					const num = parseInt(val, 10);
					if (isNaN(num)) {
						throw new Error('Invalid number');
					}
					return num;
				})
			});

			const formData = new FormData();
			formData.append('age', '25');

			const request = new Request('http://localhost/test', {
				method: 'POST',
				body: formData
			});

			const result = await parseFormData(request, schema);

			expect(isParseSuccess(result)).toBe(true);
			if (isParseSuccess(result)) {
				expect(result.age).toBe(25);
			}
		});
	});

	describe('isParseSuccess', () => {
		it('should return true for successful parse result', () => {
			const result = { name: 'John', email: 'john@example.com' };
			expect(isParseSuccess(result)).toBe(true);
		});

		it('should return false for ActionFailure', () => {
			const result = fail(400, {
				success: false,
				message: 'Validation error'
			});
			expect(isParseSuccess(result)).toBe(false);
		});

		it('should correctly type guard successful results', () => {
			const result: { name: string } | ReturnType<typeof fail> = { name: 'John' };

			if (isParseSuccess(result)) {
				// TypeScript should know result is { name: string } here
				expect(result.name).toBe('John');
			}
		});
	});
});
