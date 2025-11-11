import { describe, it, expect, vi } from 'vitest';
import { z } from 'zod';
import { ApiError, post, put, patch, get, del } from './api';

describe('api', () => {
	describe('ApiError', () => {
		it('should create an ApiError with message and code', () => {
			const error = new ApiError('Test error', 404);
			expect(error).toBeInstanceOf(Error);
			expect(error.message).toBe('Test error');
			expect(error.code).toBe(404);
		});

		it('should be throwable', () => {
			expect(() => {
				throw new ApiError('Error', 500);
			}).toThrow('Error');
		});
	});

	describe('post', () => {
		const requestSchema = z.object({
			name: z.string(),
			age: z.number()
		});

		const responseSchema = z.object({
			id: z.string(),
			name: z.string(),
			age: z.number()
		});

		it('should make a POST request with correct headers and body', async () => {
			const mockFetch = vi.fn().mockResolvedValue({
				ok: true,
				json: async () => ({ id: '123', name: 'John', age: 30 })
			});

			const result = await post(
				'/api/users',
				{ name: 'John', age: 30 },
				{
					fetch: mockFetch,
					requestSchema,
					responseSchema
				}
			);

			expect(mockFetch).toHaveBeenCalledWith('/api/users', {
				method: 'POST',
				body: JSON.stringify({ name: 'John', age: 30 }),
				headers: {
					'Content-Type': 'application/json'
				}
			});
			expect(result).toEqual({ id: '123', name: 'John', age: 30 });
		});

		it('should parse and validate response data', async () => {
			const mockFetch = vi.fn().mockResolvedValue({
				ok: true,
				json: async () => ({ id: '123', name: 'John', age: 30 })
			});

			const result = await post(
				'/api/users',
				{ name: 'John', age: 30 },
				{
					fetch: mockFetch,
					requestSchema,
					responseSchema
				}
			);

			expect(result).toEqual({ id: '123', name: 'John', age: 30 });
		});

		it('should warn on invalid request data but still make request', async () => {
			const consoleSpy = vi.spyOn(console, 'warn').mockImplementation(() => {});
			const mockFetch = vi.fn().mockResolvedValue({
				ok: true,
				json: async () => ({ id: '123', name: 'John', age: 30 })
			});

			await post(
				'/api/users',
				{ name: 'John', age: 'invalid' } as unknown as z.infer<typeof requestSchema>,
				{
					fetch: mockFetch,
					requestSchema,
					responseSchema
				}
			);

			expect(consoleSpy).toHaveBeenCalled();
			expect(mockFetch).toHaveBeenCalled();
			consoleSpy.mockRestore();
		});

		it('should throw ApiError on non-ok response with message', async () => {
			const mockFetch = vi.fn().mockResolvedValue({
				ok: false,
				status: 400,
				json: async () => ({ message: 'Bad request' })
			});

			await expect(
				post(
					'/api/users',
					{ name: 'John', age: 30 },
					{
						fetch: mockFetch,
						requestSchema,
						responseSchema
					}
				)
			).rejects.toThrow(ApiError);

			try {
				await post(
					'/api/users',
					{ name: 'John', age: 30 },
					{
						fetch: mockFetch,
						requestSchema,
						responseSchema
					}
				);
			} catch (error) {
				expect(error).toBeInstanceOf(ApiError);
				expect((error as ApiError).message).toBe('Bad request');
				expect((error as ApiError).code).toBe(400);
			}
		});

		it('should throw ApiError with default message when no message in response', async () => {
			const mockFetch = vi.fn().mockResolvedValue({
				ok: false,
				status: 500,
				json: async () => ({})
			});

			await expect(
				post(
					'/api/users',
					{ name: 'John', age: 30 },
					{
						fetch: mockFetch,
						requestSchema,
						responseSchema
					}
				)
			).rejects.toThrow(ApiError);

			try {
				await post(
					'/api/users',
					{ name: 'John', age: 30 },
					{
						fetch: mockFetch,
						requestSchema,
						responseSchema
					}
				);
			} catch (error) {
				expect(error).toBeInstanceOf(ApiError);
				expect((error as ApiError).message).toBe('Request failed with invalid status code: 500');
				expect((error as ApiError).code).toBe(500);
			}
		});

		it('should throw error when response schema validation fails', async () => {
			const mockFetch = vi.fn().mockResolvedValue({
				ok: true,
				json: async () => ({ invalid: 'data' })
			});

			await expect(
				post(
					'/api/users',
					{ name: 'John', age: 30 },
					{
						fetch: mockFetch,
						requestSchema,
						responseSchema
					}
				)
			).rejects.toThrow();
		});
	});

	describe('put', () => {
		const requestSchema = z.object({
			name: z.string()
		});

		const responseSchema = z.object({
			id: z.string(),
			name: z.string()
		});

		it('should make a PUT request with correct headers and body', async () => {
			const mockFetch = vi.fn().mockResolvedValue({
				ok: true,
				json: async () => ({ id: '123', name: 'John' })
			});

			const result = await put(
				'/api/users/123',
				{ name: 'John' },
				{
					fetch: mockFetch,
					requestSchema,
					responseSchema
				}
			);

			expect(mockFetch).toHaveBeenCalledWith('/api/users/123', {
				method: 'PUT',
				body: JSON.stringify({ name: 'John' }),
				headers: {
					'Content-Type': 'application/json'
				}
			});
			expect(result).toEqual({ id: '123', name: 'John' });
		});

		it('should throw ApiError on non-ok response', async () => {
			const mockFetch = vi.fn().mockResolvedValue({
				ok: false,
				status: 404,
				json: async () => ({})
			});

			await expect(
				put(
					'/api/users/123',
					{ name: 'John' },
					{
						fetch: mockFetch,
						requestSchema,
						responseSchema
					}
				)
			).rejects.toThrow(ApiError);

			try {
				await put(
					'/api/users/123',
					{ name: 'John' },
					{
						fetch: mockFetch,
						requestSchema,
						responseSchema
					}
				);
			} catch (error) {
				expect(error).toBeInstanceOf(ApiError);
				expect((error as ApiError).code).toBe(404);
			}
		});
	});

	describe('patch', () => {
		const requestSchema = z.object({
			name: z.string().optional()
		});

		const responseSchema = z.object({
			id: z.string(),
			name: z.string()
		});

		it('should make a PATCH request with correct headers and body', async () => {
			const mockFetch = vi.fn().mockResolvedValue({
				ok: true,
				json: async () => ({ id: '123', name: 'Jane' })
			});

			const result = await patch(
				'/api/users/123',
				{ name: 'Jane' },
				{
					fetch: mockFetch,
					requestSchema,
					responseSchema
				}
			);

			expect(mockFetch).toHaveBeenCalledWith('/api/users/123', {
				method: 'PATCH',
				body: JSON.stringify({ name: 'Jane' }),
				headers: {
					'Content-Type': 'application/json'
				}
			});
			expect(result).toEqual({ id: '123', name: 'Jane' });
		});

		it('should throw ApiError on non-ok response', async () => {
			const mockFetch = vi.fn().mockResolvedValue({
				ok: false,
				status: 400,
				json: async () => ({})
			});

			await expect(
				patch(
					'/api/users/123',
					{ name: 'Jane' },
					{
						fetch: mockFetch,
						requestSchema,
						responseSchema
					}
				)
			).rejects.toThrow(ApiError);
		});
	});

	describe('get', () => {
		const responseSchema = z.object({
			id: z.string(),
			name: z.string()
		});

		it('should make a GET request with correct headers', async () => {
			const mockFetch = vi.fn().mockResolvedValue({
				ok: true,
				json: async () => ({ id: '123', name: 'John' })
			});

			const result = await get('/api/users/123', {
				fetch: mockFetch,
				responseSchema
			});

			expect(mockFetch).toHaveBeenCalledWith('/api/users/123', {
				headers: {
					'Content-Type': 'application/json'
				}
			});
			expect(result).toEqual({ id: '123', name: 'John' });
		});

		it('should handle path with existing query parameters', async () => {
			const mockFetch = vi.fn().mockResolvedValue({
				ok: true,
				json: async () => ({ id: '123', name: 'John' })
			});

			const paramsSchema = z.object({
				page: z.number()
			});

			await get('/api/users?existing=param', {
				fetch: mockFetch,
				responseSchema,
				paramsSchema,
				params: { page: 1 }
			});

			expect(mockFetch).toHaveBeenCalledWith(
				expect.stringContaining('/api/users'),
				expect.any(Object)
			);
			const callUrl = mockFetch.mock.calls[0][0];
			expect(callUrl).toContain('existing=param');
			expect(callUrl).toContain('page=1');
		});

		it('should convert params to query string format', async () => {
			const mockFetch = vi.fn().mockResolvedValue({
				ok: true,
				json: async () => ({ id: '123', name: 'John' })
			});

			const paramsSchema = z.object({
				page: z.number(),
				filter: z.object({
					name: z.string()
				})
			});

			await get('/api/users', {
				fetch: mockFetch,
				responseSchema,
				paramsSchema,
				params: { page: 1, filter: { name: 'John' } }
			});

			const callUrl = mockFetch.mock.calls[0][0];
			expect(callUrl).toContain('page=1');
			// URLSearchParams encodes brackets, so check for encoded version
			expect(decodeURIComponent(callUrl)).toContain('filter[name]=John');
		});

		it('should handle array params', async () => {
			const mockFetch = vi.fn().mockResolvedValue({
				ok: true,
				json: async () => ({ id: '123', name: 'John' })
			});

			const paramsSchema = z.object({
				ids: z.array(z.string())
			});

			await get('/api/users', {
				fetch: mockFetch,
				responseSchema,
				paramsSchema,
				params: { ids: ['1', '2', '3'] }
			});

			const callUrl = mockFetch.mock.calls[0][0];
			const decodedUrl = decodeURIComponent(callUrl);
			expect(decodedUrl).toContain('ids[0]=1');
			expect(decodedUrl).toContain('ids[1]=2');
			expect(decodedUrl).toContain('ids[2]=3');
		});

		it('should validate params with paramsSchema', async () => {
			const mockFetch = vi.fn().mockResolvedValue({
				ok: true,
				json: async () => ({ id: '123', name: 'John' })
			});

			const paramsSchema = z.object({
				page: z.number()
			});

			await expect(
				get('/api/users', {
					fetch: mockFetch,
					responseSchema,
					paramsSchema,
					params: { page: 'invalid' } as unknown as z.infer<typeof paramsSchema>
				})
			).rejects.toThrow();
		});

		it('should throw ApiError on non-ok response', async () => {
			const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
			const mockFetch = vi.fn().mockResolvedValue({
				ok: false,
				status: 404,
				json: async () => ({ error: 'Not found' })
			});

			await expect(
				get('/api/users/123', {
					fetch: mockFetch,
					responseSchema
				})
			).rejects.toThrow(ApiError);

			expect(consoleErrorSpy).toHaveBeenCalled();
			consoleErrorSpy.mockRestore();
		});

		it('should log parsed data on success', async () => {
			const consoleLogSpy = vi.spyOn(console, 'log').mockImplementation(() => {});
			const mockFetch = vi.fn().mockResolvedValue({
				ok: true,
				json: async () => ({ id: '123', name: 'John' })
			});

			await get('/api/users/123', {
				fetch: mockFetch,
				responseSchema
			});

			expect(consoleLogSpy).toHaveBeenCalledWith('Parsed data:', { id: '123', name: 'John' });
			consoleLogSpy.mockRestore();
		});

		it('should handle null and undefined params', async () => {
			const mockFetch = vi.fn().mockResolvedValue({
				ok: true,
				json: async () => ({ id: '123', name: 'John' })
			});

			const paramsSchema = z.object({
				page: z.number().optional(),
				filter: z.string().optional()
			});

			await get('/api/users', {
				fetch: mockFetch,
				responseSchema,
				paramsSchema,
				params: { page: 1, filter: undefined }
			});

			const callUrl = mockFetch.mock.calls[0][0];
			expect(callUrl).toContain('page=1');
			expect(callUrl).not.toContain('filter');
		});
	});

	describe('del', () => {
		it('should make a DELETE request with correct headers', async () => {
			const mockFetch = vi.fn().mockResolvedValue({
				ok: true
			});

			await del('/api/users/123', {
				fetch: mockFetch
			});

			expect(mockFetch).toHaveBeenCalledWith('/api/users/123', {
				method: 'DELETE',
				headers: {
					'Content-Type': 'application/json'
				}
			});
		});

		it('should throw ApiError on non-ok response', async () => {
			const mockFetch = vi.fn().mockResolvedValue({
				ok: false,
				status: 404
			});

			await expect(
				del('/api/users/123', {
					fetch: mockFetch
				})
			).rejects.toThrow(ApiError);

			try {
				await del('/api/users/123', {
					fetch: mockFetch
				});
			} catch (error) {
				expect(error).toBeInstanceOf(ApiError);
				expect((error as ApiError).code).toBe(404);
			}
		});

		it('should return void on success', async () => {
			const mockFetch = vi.fn().mockResolvedValue({
				ok: true
			});

			const result = await del('/api/users/123', {
				fetch: mockFetch
			});

			expect(result).toBeUndefined();
		});
	});

	describe('query string conversion', () => {
		it('should handle nested objects in get params', async () => {
			const mockFetch = vi.fn().mockResolvedValue({
				ok: true,
				json: async () => ({})
			});

			const paramsSchema = z.object({
				user: z.object({
					name: z.string(),
					age: z.number()
				})
			});

			await get('/api/search', {
				fetch: mockFetch,
				responseSchema: z.object({}),
				paramsSchema,
				params: { user: { name: 'John', age: 30 } }
			});

			const callUrl = mockFetch.mock.calls[0][0];
			const decodedUrl = decodeURIComponent(callUrl);
			expect(decodedUrl).toContain('user[name]=John');
			expect(decodedUrl).toContain('user[age]=30');
		});

		it('should handle deeply nested objects', async () => {
			const mockFetch = vi.fn().mockResolvedValue({
				ok: true,
				json: async () => ({})
			});

			const paramsSchema = z.object({
				filter: z.object({
					user: z.object({
						name: z.string()
					})
				})
			});

			await get('/api/search', {
				fetch: mockFetch,
				responseSchema: z.object({}),
				paramsSchema,
				params: { filter: { user: { name: 'John' } } }
			});

			const callUrl = mockFetch.mock.calls[0][0];
			const decodedUrl = decodeURIComponent(callUrl);
			expect(decodedUrl).toContain('filter[user][name]=John');
		});

		it('should skip null and undefined values', async () => {
			const mockFetch = vi.fn().mockResolvedValue({
				ok: true,
				json: async () => ({})
			});

			const paramsSchema = z.object({
				name: z.string().optional(),
				age: z.number().nullable().optional()
			});

			await get('/api/search', {
				fetch: mockFetch,
				responseSchema: z.object({}),
				paramsSchema,
				params: { name: 'John', age: null }
			});

			const callUrl = mockFetch.mock.calls[0][0];
			expect(callUrl).toContain('name=John');
			expect(callUrl).not.toContain('age');
		});
	});
});
