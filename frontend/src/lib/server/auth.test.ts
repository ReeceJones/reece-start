import { describe, it, expect, vi, beforeEach } from 'vitest';
import type { RequestEvent, Cookies } from '@sveltejs/kit';
import { jwtDecode } from 'jwt-decode';
import type { JwtClaims } from '$lib/schemas/jwt';
import { ApiError, get, post } from '$lib/api';
import { getRequestEvent } from '$app/server';
import * as authModule from './auth';

// Mock SvelteKit server functions
vi.mock('$app/server', () => ({
	getRequestEvent: vi.fn()
}));

vi.mock('@sveltejs/kit', () => ({
	error: vi.fn((code: number, message: string) => {
		const err = new Error(message) as Error & { status: number };
		err.status = code;
		throw err;
	}),
	redirect: vi.fn((status: number, url: string) => {
		const err = new Error(`redirect to ${url}`) as Error & { status: number; location: string };
		err.status = status;
		err.location = url;
		throw err;
	})
}));

vi.mock('jwt-decode', () => ({
	jwtDecode: vi.fn()
}));

vi.mock('$lib/api', () => ({
	ApiError: class ApiError extends Error {
		constructor(
			message: string,
			public code: number
		) {
			super(message);
		}
	},
	get: vi.fn(),
	post: vi.fn()
}));

describe('auth', () => {
	const mockCookies = {
		get: vi.fn(),
		set: vi.fn(),
		delete: vi.fn(),
		getAll: vi.fn(),
		serialize: vi.fn()
	} as Cookies & {
		get: ReturnType<typeof vi.fn>;
		set: ReturnType<typeof vi.fn>;
		delete: ReturnType<typeof vi.fn>;
		getAll: ReturnType<typeof vi.fn>;
		serialize: ReturnType<typeof vi.fn>;
	};

	const createMockRequestEvent = (overrides: Partial<RequestEvent> = {}): RequestEvent => {
		return {
			cookies: mockCookies,
			url: new URL('https://example.com/app'),
			params: {},
			fetch: vi.fn(),
			getClientAddress: vi.fn(),
			locals: {},
			platform: undefined,
			request: new Request('https://example.com/app'),
			route: { id: null },
			setHeaders: vi.fn(),
			isDataRequest: false,
			isSubRequest: false,
			...overrides
		} as RequestEvent;
	};

	beforeEach(() => {
		vi.clearAllMocks();
		vi.mocked(getRequestEvent).mockReturnValue(createMockRequestEvent());
	});

	describe('isLoggedIn', () => {
		it('should return false when token is not present', () => {
			mockCookies.get.mockReturnValue(undefined);
			vi.mocked(jwtDecode).mockReturnValue({ exp: undefined } as Partial<JwtClaims>);

			const result = authModule.isLoggedIn();
			expect(result).toBe(false);
		});

		it('should return false when token is expired', () => {
			mockCookies.get.mockReturnValue('valid-token');
			const expiredTime = Math.floor(Date.now() / 1000) - 1000; // 1000 seconds ago
			vi.mocked(jwtDecode).mockReturnValue({ exp: expiredTime } as Partial<JwtClaims>);

			const result = authModule.isLoggedIn();
			expect(result).toBe(false);
		});

		it('should return true when token is valid and not expired', () => {
			mockCookies.get.mockReturnValue('valid-token');
			const futureTime = Math.floor(Date.now() / 1000) + 3600; // 1 hour from now
			vi.mocked(jwtDecode).mockReturnValue({ exp: futureTime } as Partial<JwtClaims>);

			const result = authModule.isLoggedIn();
			expect(result).toBe(true);
		});

		it('should return false when exp is missing', () => {
			mockCookies.get.mockReturnValue('valid-token');
			vi.mocked(jwtDecode).mockReturnValue({} as Partial<JwtClaims>);

			const result = authModule.isLoggedIn();
			expect(result).toBe(false);
		});
	});

	describe('authenticate', () => {
		it('should call performAuthenticationCheck', async () => {
			const mockRequestEvent = createMockRequestEvent({
				url: new URL('https://example.com/app')
			});
			mockCookies.get.mockReturnValue('valid-token');
			vi.mocked(jwtDecode).mockReturnValue({
				exp: Math.floor(Date.now() / 1000) + 3600
			} as Partial<JwtClaims>);

			vi.mocked(getRequestEvent).mockReturnValue(mockRequestEvent);

			await authModule.authenticate();
			// Should not throw if path doesn't require organization validation
		});
	});

	describe('performAuthenticationCheck', () => {
		it('should return early for non-app paths', async () => {
			const mockRequestEvent = createMockRequestEvent({
				url: new URL('https://example.com/public')
			});

			await authModule.performAuthenticationCheck(mockRequestEvent);
			// Should not throw or call validation functions
		});

		it('should validate token expiration for app paths', async () => {
			const mockRequestEvent = createMockRequestEvent({
				url: new URL('https://example.com/app/dashboard')
			});
			mockCookies.get.mockReturnValue('valid-token');
			vi.mocked(jwtDecode).mockReturnValue({
				exp: Math.floor(Date.now() / 1000) + 3600
			} as Partial<JwtClaims>);

			await authModule.performAuthenticationCheck(mockRequestEvent);
			// Should not throw for valid token
		});

		it('should validate admin role for admin paths', async () => {
			const mockRequestEvent = createMockRequestEvent({
				url: new URL('https://example.com/app/admin')
			});
			mockCookies.get.mockReturnValue('admin-token');
			vi.mocked(jwtDecode).mockReturnValue({
				exp: Math.floor(Date.now() / 1000) + 3600,
				role: 'admin'
			} as Partial<JwtClaims>);

			await authModule.performAuthenticationCheck(mockRequestEvent);
			// Should not throw for admin role
		});

		it('should skip organization validation when organizationId is undefined', async () => {
			const mockRequestEvent = createMockRequestEvent({
				url: new URL('https://example.com/app'),
				params: {}
			});
			mockCookies.get.mockReturnValue('valid-token');
			vi.mocked(jwtDecode).mockReturnValue({
				exp: Math.floor(Date.now() / 1000) + 3600
			} as Partial<JwtClaims>);

			await authModule.performAuthenticationCheck(mockRequestEvent);
			// Should not throw
		});
	});

	describe('getUserAndValidateToken', () => {
		it('should return user when token is valid', async () => {
			const mockRequestEvent = createMockRequestEvent();
			mockCookies.get.mockReturnValue('valid-token');
			const iat = Math.floor(Date.now() / 1000) - 100;
			vi.mocked(jwtDecode).mockReturnValue({ iat } as Partial<JwtClaims>);

			const mockUser = {
				data: {
					id: 'user-123',
					type: 'user',
					attributes: {
						name: 'John Doe',
						email: 'john@example.com'
					},
					meta: {
						tokenRevocation: {
							lastIssuedAt: iat + 50, // lastIssuedAt is AFTER iat, so iat < lastIssuedAt, no refresh
							canRefresh: true
						}
					}
				}
			};

			vi.mocked(get).mockResolvedValue(mockUser);

			vi.mocked(getRequestEvent).mockReturnValue(mockRequestEvent);

			const result = await authModule.getUserAndValidateToken();
			expect(result.user).toEqual(mockUser);
			expect(post).not.toHaveBeenCalled();
		});

		it('should throw error when iat is missing', async () => {
			const mockRequestEvent = createMockRequestEvent();
			mockCookies.get.mockReturnValue('valid-token');
			vi.mocked(jwtDecode).mockReturnValue({} as Partial<JwtClaims>);

			vi.mocked(getRequestEvent).mockReturnValue(mockRequestEvent);

			await expect(authModule.getUserAndValidateToken()).rejects.toThrow();
		});

		it('should throw ApiError when API call fails', async () => {
			const mockRequestEvent = createMockRequestEvent();
			mockCookies.get.mockReturnValue('valid-token');
			vi.mocked(jwtDecode).mockReturnValue({
				iat: Math.floor(Date.now() / 1000) - 100
			} as Partial<JwtClaims>);

			vi.mocked(get).mockRejectedValue(new ApiError('Not found', 404));

			vi.mocked(getRequestEvent).mockReturnValue(mockRequestEvent);

			await expect(authModule.getUserAndValidateToken()).rejects.toThrow();
		});

		it('should refresh token when token is older than lastIssuedAt', async () => {
			const mockRequestEvent = createMockRequestEvent();
			mockCookies.get.mockReturnValue('valid-token');
			const iat = Math.floor(Date.now() / 1000) - 100;
			vi.mocked(jwtDecode).mockReturnValue({ iat } as Partial<JwtClaims>);

			const mockUser = {
				data: {
					id: 'user-123',
					type: 'user',
					attributes: {
						name: 'John Doe',
						email: 'john@example.com'
					},
					meta: {
						tokenRevocation: {
							lastIssuedAt: iat - 200, // Token is older
							canRefresh: true
						}
					}
				}
			};

			vi.mocked(get).mockResolvedValue(mockUser);
			vi.mocked(post).mockResolvedValue({
				data: {
					type: 'token',
					meta: { token: 'new-token' }
				}
			} as Awaited<ReturnType<typeof post>>);

			vi.mocked(getRequestEvent).mockReturnValue(mockRequestEvent);

			await authModule.getUserAndValidateToken();
			expect(post).toHaveBeenCalled();
		});

		it('should redirect when token cannot be refreshed', async () => {
			const mockRequestEvent = createMockRequestEvent();
			mockCookies.get.mockReturnValue('valid-token');
			const iat = Math.floor(Date.now() / 1000) - 100;
			vi.mocked(jwtDecode).mockReturnValue({ iat } as Partial<JwtClaims>);

			const mockUser = {
				data: {
					id: 'user-123',
					type: 'user',
					attributes: {
						name: 'John Doe',
						email: 'john@example.com'
					},
					meta: {
						tokenRevocation: {
							lastIssuedAt: iat - 200,
							canRefresh: false // Cannot refresh
						}
					}
				}
			};

			vi.mocked(get).mockResolvedValue(mockUser);

			vi.mocked(getRequestEvent).mockReturnValue(mockRequestEvent);

			await expect(authModule.getUserAndValidateToken()).rejects.toThrow('redirect');
			expect(mockCookies.delete).toHaveBeenCalledWith('app-session-token', { path: '/' });
		});
	});

	describe('refreshUserToken', () => {
		it('should refresh token successfully', async () => {
			const mockRequestEvent = createMockRequestEvent({
				params: { organizationId: 'org-123' }
			});

			vi.mocked(post).mockResolvedValue({
				data: {
					type: 'token',
					meta: { token: 'new-token-123' }
				}
			} as Awaited<ReturnType<typeof post>>);

			await authModule.refreshUserToken(mockRequestEvent);

			expect(post).toHaveBeenCalledWith(
				'/api/users/me/token',
				expect.objectContaining({
					data: expect.objectContaining({
						type: 'token',
						relationships: expect.objectContaining({
							organization: expect.any(Object)
						})
					})
				}),
				expect.any(Object)
			);
			expect(mockCookies.set).toHaveBeenCalledWith(
				'app-session-token',
				'new-token-123',
				expect.objectContaining({
					path: '/',
					httpOnly: true,
					secure: true,
					sameSite: 'strict',
					maxAge: 60 * 60 * 24 * 30
				})
			);
		});

		it('should include impersonatedUserId in request when provided', async () => {
			const mockRequestEvent = createMockRequestEvent({
				params: { organizationId: 'org-123' }
			});

			vi.mocked(post).mockResolvedValue({
				data: {
					type: 'token',
					meta: { token: 'new-token' }
				}
			} as Awaited<ReturnType<typeof post>>);

			await authModule.refreshUserToken(mockRequestEvent, {
				impersonatedUserId: 'user-456'
			});

			expect(post).toHaveBeenCalledWith(
				'/api/users/me/token',
				expect.objectContaining({
					data: expect.objectContaining({
						relationships: expect.objectContaining({
							impersonatedUser: expect.any(Object)
						})
					})
				}),
				expect.any(Object)
			);
		});

		it('should include stopImpersonating flag when provided', async () => {
			const mockRequestEvent = createMockRequestEvent({
				params: { organizationId: 'org-123' }
			});

			vi.mocked(post).mockResolvedValue({
				data: {
					type: 'token',
					meta: { token: 'new-token' }
				}
			} as Awaited<ReturnType<typeof post>>);

			await authModule.refreshUserToken(mockRequestEvent, {
				stopImpersonating: true
			});

			expect(post).toHaveBeenCalledWith(
				'/api/users/me/token',
				expect.objectContaining({
					data: expect.objectContaining({
						meta: expect.objectContaining({
							stopImpersonating: true
						})
					})
				}),
				expect.any(Object)
			);
		});

		it('should handle 404 error for organization not found', async () => {
			const mockRequestEvent = createMockRequestEvent({
				params: { organizationId: 'org-123' }
			});

			vi.mocked(post).mockRejectedValue(new ApiError('Not found', 404));

			await expect(authModule.refreshUserToken(mockRequestEvent)).rejects.toThrow();
		});
	});

	describe('impersonateUser', () => {
		it('should call refreshUserToken with impersonatedUserId', async () => {
			const mockRequestEvent = createMockRequestEvent();
			vi.mocked(post).mockResolvedValue({
				data: {
					type: 'token',
					meta: { token: 'new-token' }
				}
			} as Awaited<ReturnType<typeof post>>);

			await authModule.impersonateUser(mockRequestEvent, 'user-123');

			expect(post).toHaveBeenCalledWith(
				'/api/users/me/token',
				expect.objectContaining({
					data: expect.objectContaining({
						relationships: expect.objectContaining({
							impersonatedUser: expect.any(Object)
						})
					})
				}),
				expect.any(Object)
			);
		});
	});

	describe('stopImpersonatingUser', () => {
		it('should call refreshUserToken with stopImpersonating flag', async () => {
			const mockRequestEvent = createMockRequestEvent();
			vi.mocked(post).mockResolvedValue({
				data: {
					type: 'token',
					meta: { token: 'new-token' }
				}
			} as Awaited<ReturnType<typeof post>>);

			await authModule.stopImpersonatingUser(mockRequestEvent);

			expect(post).toHaveBeenCalledWith(
				'/api/users/me/token',
				expect.objectContaining({
					data: expect.objectContaining({
						meta: expect.objectContaining({
							stopImpersonating: true
						})
					})
				}),
				expect.any(Object)
			);
		});
	});

	describe('getUserScopes', () => {
		it('should return scopes from token', () => {
			const mockRequestEvent = createMockRequestEvent();
			mockCookies.get.mockReturnValue('valid-token');
			vi.mocked(jwtDecode).mockReturnValue({
				scopes: ['scope1', 'scope2'] as unknown as JwtClaims['scopes']
			} as Partial<JwtClaims>);

			vi.mocked(getRequestEvent).mockReturnValue(mockRequestEvent);

			const scopes = authModule.getUserScopes();
			expect(scopes).toEqual(['scope1', 'scope2']);
		});

		it('should return empty array when scopes are missing', () => {
			const mockRequestEvent = createMockRequestEvent();
			mockCookies.get.mockReturnValue('valid-token');
			vi.mocked(jwtDecode).mockReturnValue({} as Partial<JwtClaims>);

			vi.mocked(getRequestEvent).mockReturnValue(mockRequestEvent);

			const scopes = authModule.getUserScopes();
			expect(scopes).toEqual([]);
		});
	});

	describe('getIsImpersonatingUser', () => {
		it('should return true when impersonating', () => {
			const mockRequestEvent = createMockRequestEvent();
			mockCookies.get.mockReturnValue('valid-token');
			vi.mocked(jwtDecode).mockReturnValue({
				is_impersonating: true
			} as Partial<JwtClaims>);

			vi.mocked(getRequestEvent).mockReturnValue(mockRequestEvent);

			const result = authModule.getIsImpersonatingUser();
			expect(result).toBe(true);
		});

		it('should return false when not impersonating', () => {
			const mockRequestEvent = createMockRequestEvent();
			mockCookies.get.mockReturnValue('valid-token');
			vi.mocked(jwtDecode).mockReturnValue({
				is_impersonating: false
			} as Partial<JwtClaims>);

			vi.mocked(getRequestEvent).mockReturnValue(mockRequestEvent);

			const result = authModule.getIsImpersonatingUser();
			expect(result).toBe(false);
		});

		it('should return false when is_impersonating is missing', () => {
			const mockRequestEvent = createMockRequestEvent();
			mockCookies.get.mockReturnValue('valid-token');
			vi.mocked(jwtDecode).mockReturnValue({} as Partial<JwtClaims>);

			vi.mocked(getRequestEvent).mockReturnValue(mockRequestEvent);

			const result = authModule.getIsImpersonatingUser();
			expect(result).toBe(false);
		});
	});

	describe('setTokenInCookies', () => {
		it('should set token in cookies with correct options', () => {
			const mockRequestEvent = createMockRequestEvent();
			const token = 'test-token-123';

			authModule.setTokenInCookies(mockRequestEvent, token);

			expect(mockCookies.set).toHaveBeenCalledWith('app-session-token', token, {
				path: '/',
				httpOnly: true,
				secure: true,
				sameSite: 'strict',
				maxAge: 60 * 60 * 24 * 30 // 30 days
			});
		});
	});
});
