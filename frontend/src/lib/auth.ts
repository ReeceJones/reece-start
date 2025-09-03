import type { OrganizationScope } from './schemas/jwt';
import { getContext, setContext } from 'svelte';

const scopesKey = 'membership-scopes';

export function setScopes(scopes: OrganizationScope[]) {
	setContext(scopesKey, scopes);
}

export function getScopes() {
	return getContext<OrganizationScope[]>(scopesKey) ?? [];
}

export function hasScope(scope: OrganizationScope): boolean {
	return getScopes().includes(scope) ?? false;
}
