import type { JwtPayload } from 'jwt-decode';
import type { OrganizationMembershipRole } from './organization-membership';

export enum OrganizationScope {
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
	OrganizationInvitationsDelete = 'organization:invitations:delete'
}

interface OrganizationClaims {
	organization_id?: string;
	organization_role?: OrganizationMembershipRole;
	organization_scopes?: OrganizationScope[];
}

export type JwtClaims = JwtPayload & OrganizationClaims;
