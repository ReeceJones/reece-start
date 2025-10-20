import type { JwtPayload } from 'jwt-decode';
import type { OrganizationMembershipRole } from './organization-membership';
import type { UserRole } from './user';

export enum UserScope {
	// Organization
	OrganizationRead = 'organization:read',
	OrganizationUpdate = 'organization:update',
	OrganizationDelete = 'organization:delete',
	OrganizationMembershipsList = 'organization:memberships:list',
	OrganizationMembershipsRead = 'organization:memberships:read',
	OrganizationMembershipsCreate = 'organization:memberships:create',
	OrganizationMembershipsUpdate = 'organization:memberships:update',
	OrganizationMembershipsDelete = 'organization:memberships:delete',
	OrganizationInvitationsList = 'organization:invitations:list',
	OrganizationInvitationsRead = 'organization:invitations:read',
	OrganizationInvitationsCreate = 'organization:invitations:create',
	OrganizationInvitationsUpdate = 'organization:invitations:update',
	OrganizationInvitationsDelete = 'organization:invitations:delete',
	OrganizationStripeUpdate = 'organization:stripe:update',

	// Admin
	Admin = 'admin',
	AdminUsersList = 'admin:users:list',
	AdminUsersRead = 'admin:users:read',
	AdminUsersImpersonate = 'admin:users:impersonate'
}

interface OrganizationClaims {
	organization_id?: string;
	organization_role?: OrganizationMembershipRole;
	scopes?: UserScope[];
	role?: UserRole;
	is_impersonating?: boolean;
	impersonating_user_id?: string;
}

export type JwtClaims = JwtPayload & OrganizationClaims;
