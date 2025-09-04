import type { UserScope } from './schemas/jwt';
import { getContext, setContext } from 'svelte';

const scopesKey = 'membership-scopes';

export function setScopes(scopes: UserScope[]) {
	setContext(scopesKey, scopes);
}

export function getScopes() {
	return getContext<UserScope[]>(scopesKey) ?? [];
}

export function hasScope(scope: UserScope): boolean {
	return getScopes().includes(scope) ?? false;
}
