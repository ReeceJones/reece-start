import { render, screen, cleanup, waitFor } from '@testing-library/svelte';
import { expect, describe, it, beforeEach, afterEach, vi } from 'vitest';
import userEvent from '@testing-library/user-event';
import InvitationRow from './InvitationRow.svelte';
import type { OrganizationInvitation } from '$lib/schemas/organization-invitation';
import { API_TYPES } from '$lib/schemas/api';
import { UserScope } from '$lib/schemas/jwt';
import { setScopes } from '$lib/auth';
import * as clipboard from '$lib/clipboard';

// Mock dependencies
vi.mock('$app/forms', async () => {
	const { createMockEnhance } = await import('$lib/test-utils');
	return createMockEnhance();
});

vi.mock('$lib/clipboard');

describe('InvitationRow', () => {
	const mockInvitation: OrganizationInvitation = {
		id: 'inv-123',
		type: API_TYPES.organizationInvitation,
		attributes: {
			email: 'test@example.com',
			role: 'member',
			status: 'pending'
		},
		relationships: {
			organization: {
				data: {
					id: 'org-123',
					type: API_TYPES.organization
				}
			},
			invitingUser: {
				data: {
					id: 'user-123',
					type: API_TYPES.user
				}
			}
		}
	};

	beforeEach(() => {
		// Mock window.location.origin
		Object.defineProperty(window, 'location', {
			value: {
				origin: 'http://localhost:5173'
			},
			writable: true
		});

		// Mock copyToClipboard by default
		vi.spyOn(clipboard, 'copyToClipboard').mockResolvedValue(undefined);
	});

	afterEach(() => {
		cleanup();
		vi.clearAllMocks();
	});

	describe('rendering', () => {
		it('should render the invitation row', () => {
			setScopes([UserScope.OrganizationInvitationsDelete]);
			const { container } = render(InvitationRow, {
				props: { invitation: mockInvitation }
			});

			const row = container.querySelector('tr');
			expect(row).toBeTruthy();
			expect(row?.className).toContain('hover:bg-base-300');
		});

		it('should display the invitation email as a mailto link', () => {
			setScopes([UserScope.OrganizationInvitationsDelete]);
			render(InvitationRow, {
				props: { invitation: mockInvitation }
			});

			const emailLink = screen.getByRole('link', { name: 'test@example.com' });
			expect(emailLink).toBeTruthy();
			expect(emailLink.getAttribute('href')).toBe('mailto:test@example.com');
		});

		it('should render the copy invitation link button', () => {
			setScopes([UserScope.OrganizationInvitationsDelete]);
			render(InvitationRow, {
				props: { invitation: mockInvitation }
			});

			const copyButton = screen.getByRole('button', { name: /copy invitation link/i });
			expect(copyButton).toBeTruthy();
		});

		it('should render the delete button', () => {
			setScopes([UserScope.OrganizationInvitationsDelete]);
			const { container } = render(InvitationRow, {
				props: { invitation: mockInvitation }
			});

			const form = container.querySelector('form');
			expect(form).toBeTruthy();
			const deleteButton = form?.querySelector('button');
			expect(deleteButton).toBeTruthy();
		});

		it('should have correct form action for delete', () => {
			setScopes([UserScope.OrganizationInvitationsDelete]);
			const { container } = render(InvitationRow, {
				props: { invitation: mockInvitation }
			});

			const form = container.querySelector('form');
			expect(form?.getAttribute('action')).toBe('/app/org-123/settings/members?/deleteInvitation');
		});

		it('should have hidden input with invitation ID', () => {
			setScopes([UserScope.OrganizationInvitationsDelete]);
			const { container } = render(InvitationRow, {
				props: { invitation: mockInvitation }
			});

			const hiddenInput = container.querySelector('input[type="hidden"][name="invitationId"]');
			expect(hiddenInput).toBeTruthy();
			expect(hiddenInput?.getAttribute('value')).toBe('inv-123');
		});
	});

	describe('copy invitation link', () => {
		it('should copy invitation link to clipboard when button is clicked', async () => {
			const copyToClipboardSpy = vi
				.spyOn(clipboard, 'copyToClipboard')
				.mockResolvedValue(undefined);
			setScopes([UserScope.OrganizationInvitationsDelete]);
			const user = userEvent.setup();
			render(InvitationRow, {
				props: { invitation: mockInvitation }
			});

			const copyButton = screen.getByRole('button', { name: /copy invitation link/i });
			await user.click(copyButton);

			expect(copyToClipboardSpy).toHaveBeenCalledWith(
				'http://localhost:5173/app/invitations/inv-123'
			);
		});

		it('should show toast notification when link is copied', async () => {
			setScopes([UserScope.OrganizationInvitationsDelete]);
			const user = userEvent.setup();
			const { container } = render(InvitationRow, {
				props: { invitation: mockInvitation }
			});

			const copyButton = screen.getByRole('button', { name: /copy invitation link/i });
			await user.click(copyButton);

			await waitFor(() => {
				const toast = container.querySelector('.toast');
				expect(toast).toBeTruthy();
			});
		});

		it('should show Check icon when link is copied', async () => {
			setScopes([UserScope.OrganizationInvitationsDelete]);
			const user = userEvent.setup();
			const { container } = render(InvitationRow, {
				props: { invitation: mockInvitation }
			});

			const copyButton = screen.getByRole('button', { name: /copy invitation link/i });
			await user.click(copyButton);

			await waitFor(() => {
				const checkIcon = container.querySelector('svg');
				expect(checkIcon).toBeTruthy();
			});
		});
	});

	describe('delete invitation', () => {
		it('should disable delete button when user does not have delete scope', () => {
			setScopes([]);
			const { container } = render(InvitationRow, {
				props: { invitation: mockInvitation }
			});

			const form = container.querySelector('form');
			const deleteButton = form?.querySelector('button');
			expect(deleteButton?.hasAttribute('disabled')).toBe(true);
		});

		it('should enable delete button when user has delete scope', () => {
			setScopes([UserScope.OrganizationInvitationsDelete]);
			const { container } = render(InvitationRow, {
				props: { invitation: mockInvitation }
			});

			const form = container.querySelector('form');
			const deleteButton = form?.querySelector('button');
			expect(deleteButton?.hasAttribute('disabled')).toBe(false);
		});
	});

	describe('permissions', () => {
		it('should show delete button enabled when user has OrganizationInvitationsDelete scope', () => {
			setScopes([UserScope.OrganizationInvitationsDelete]);
			const { container } = render(InvitationRow, {
				props: { invitation: mockInvitation }
			});

			const form = container.querySelector('form');
			const deleteButton = form?.querySelector('button');
			expect(deleteButton?.hasAttribute('disabled')).toBe(false);
		});

		it('should show delete button disabled when user does not have OrganizationInvitationsDelete scope', () => {
			setScopes([UserScope.OrganizationInvitationsList]);
			const { container } = render(InvitationRow, {
				props: { invitation: mockInvitation }
			});

			const form = container.querySelector('form');
			const deleteButton = form?.querySelector('button');
			expect(deleteButton?.hasAttribute('disabled')).toBe(true);
		});
	});
});
